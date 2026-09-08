package main

import (
	"blackledger/core"
	"blackledger/store"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//go:embed director.txt
var prompt string

//go:embed director-focused.txt
var focusedPrompt string

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

type app struct {
	voiceMu    sync.Mutex
	voiceCache map[string][]byte
	voiceOrder []string
	s          *store.Store
	ai         sync.Mutex
	port       string
	client     *http.Client
}

func reply(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func fail(w http.ResponseWriter, code int, e error) {
	reply(w, code, map[string]string{"error": e.Error()})
}
func body(r *http.Request, v any) error {
	defer r.Body.Close()
	d := json.NewDecoder(io.LimitReader(r.Body, 12001))
	return d.Decode(v)
}
func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		switch r.URL.Path {
		case "/api/health":
			reply(w, 200, map[string]any{"ok": true, "game": "Black Ledger", "core": "go", "build": runningBuild})
			return
		case "/api/state":
			s, e := a.s.Read()
			if e != nil {
				fail(w, 500, e)
			} else {
				reply(w, 200, s.Public())
			}
			return
		}
		root := os.Getenv("BLACK_LEDGER_WEB")
		if root == "" {
			root = "web"
			if _, e := os.Stat(filepath.Join("dist", "index.html")); e == nil {
				root = "dist"
			}
		}
		http.FileServer(http.Dir(root)).ServeHTTP(w, r)
		return
	}
	if r.Method != "POST" {
		fail(w, 405, fmt.Errorf("method not allowed"))
		return
	}
	origin := r.Header.Get("Origin")
	if origin != "" && origin != "http://127.0.0.1:"+a.port && origin != "http://localhost:"+a.port {
		fail(w, 403, fmt.Errorf("origin rejected"))
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		fail(w, 415, fmt.Errorf("JSON required"))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 12000)
	switch r.URL.Path {
	case "/api/action":
		var c core.Command
		if e := body(r, &c); e != nil {
			fail(w, 400, e)
			return
		}
		out, e := a.s.Command(c)
		if e != nil {
			fail(w, 409, e)
		} else {
			reply(w, 200, out)
		}
	case "/api/director":
		reply(w, 200, map[string]bool{"started": a.prepare()})
	case "/api/speech", "/api/speech/prepare":
		a.speech(w, r)

	default:
		fail(w, 404, fmt.Errorf("not found"))
	}
}
func (a *app) prepare() bool {
	if !a.ai.TryLock() {
		return false
	}
	snapshot, e := a.s.Read()
	if e != nil || !snapshot.Player.Alive || len(snapshot.Offers) >= 2 {
		a.ai.Unlock()
		return false
	}
	if e = a.s.Change(func(w *core.World) error {
		w.Director = core.Director{Status: "writing", Detail: "Preparing a new arrangement.", LastRequest: w.Minute}
		return nil
	}); e != nil {
		a.ai.Unlock()
		return false
	}
	go func() {
		defer a.ai.Unlock()
		e := a.generate(snapshot)
		if e != nil {
			log.Println("Director:", e)
			_ = a.s.Change(func(w *core.World) error {
				if w.ID == snapshot.ID && w.Life == snapshot.Life {
					w.Director.Status = "offline"
					w.Director.Detail = "Local AI unavailable or proposal rejected. Authored play remains available."
					if errors.Is(e, errDirectorContextChanged) {
						w.Director.Status = "available"
						w.Director.Detail = "The situation changed while this encounter was being prepared. A new arrangement can be requested."
						if len(w.Offers) > 0 {
							w.Director.Status = "ready"
							w.Director.Detail = "An earlier encounter is ready. The outdated new draft was discarded."
						}
					}
				}
				return nil
			})
		}
	}()
	return true
}

type proposalRejected struct{ error }

func (a *app) generate(snapshot *core.World) error {
	err := a.generateAttempt(snapshot, "")
	var rejected proposalRejected
	if !errors.As(err, &rejected) {
		return err
	}
	current, readErr := a.s.Read()
	if readErr != nil {
		return readErr
	}
	if current.ID != snapshot.ID || current.Life != snapshot.Life || !current.Player.Alive {
		return nil
	}
	log.Printf("Director correcting rejected proposal: %s", err)
	return a.generateAttempt(snapshot, "Your last proposal failed validation: "+err.Error()+". Return a corrected complete proposal under the same constraints. Do not mention this correction in character dialogue.")
}
func (a *app) generateAttempt(snapshot *core.World, feedback string) error {
	operation := snapshot.NextDirectorOperation()
	connection := directorConnection(snapshot)
	beneficiaries := []string{""}
	for _, faction := range snapshot.Factions {
		beneficiaries = append(beneficiaries, faction.ID)
	}
	contextData := map[string]any{"allowed_beneficiary_ids": beneficiaries, "current_clock": fmt.Sprintf("Day %d, %02d:%02d", snapshot.Minute/1440+1, snapshot.Minute%1440/60, snapshot.Minute%60), "validation_feedback": feedback, "required_operation": operation, "required_connection": connection, "recent_arrangements": arrangementBriefs(snapshot), "life": snapshot.Life, "minute": snapshot.Minute, "player": snapshot.Player, "factions": snapshot.Factions, "npcs": snapshot.NPCs, "dead": snapshot.Dead, "places": core.Locations, "properties": snapshot.Properties, "recent_history": recentWorldChanges(snapshot)}
	responseFormat := proposalSchema(snapshot, operation, connection)
	activePrompt := prompt
	if env("BLACK_LEDGER_DIRECTOR_BRIEF", "full") == "focused" {
		activePrompt = focusedPrompt
		contextData = focusedContext(snapshot, operation, connection, feedback, beneficiaries)
	}
	attributeDirectorContext(contextData, snapshot, connection)
	contextData["active_business_ceasefires_until_minute"] = snapshot.ActiveBusinessTruces()
	contextData["allowed_speaker_ids"] = directorSpeakers(snapshot, connection)
	affiliation := map[string][]string{}
	for _, speaker := range directorSpeakers(snapshot, connection) {
		affiliation[speaker] = speakerBeneficiaries(snapshot, speaker)
	}
	contextData["speaker_beneficiary_ids"] = affiliation
	contextData["accessible_job_locations"] = accessibleJobLocations(snapshot)
	b, _ := json.Marshal(contextData)
	// Experimental opt-in: reasoning shares the bounded generation budget with
	// the final JSON. Requests remain asynchronous and bounded; structural
	// validation and the save boundary still own what can enter the game.
	thinking := env("BLACK_LEDGER_DIRECTOR_THINK", "0") == "1"
	predictionLimit := 700
	requestTimeout := 100 * time.Second
	responseLimit := int64(20000)
	if thinking {
		predictionLimit = 4096
		requestTimeout = 180 * time.Second
		responseLimit = 128 * 1024
	}
	payload, _ := json.Marshal(map[string]any{"model": env("BLACK_LEDGER_MODEL", "qwen3:14b"), "stream": false, "think": thinking, "format": responseFormat, "messages": []map[string]string{{"role": "system", "content": activePrompt}, {"role": "user", "content": string(b)}}, "options": map[string]any{"temperature": .8, "num_predict": predictionLimit}})
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "POST", env("BLACK_LEDGER_OLLAMA", "http://127.0.0.1:11435")+"/api/chat", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res, e := a.client.Do(req)
	if e != nil {
		return e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("model returned %d", res.StatusCode)
	}
	var answer struct {
		DoneReason string `json:"done_reason"`
		Message    struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if e = json.NewDecoder(io.LimitReader(res.Body, responseLimit)).Decode(&answer); e != nil {
		return e
	}
	if answer.DoneReason == "length" {
		// A token budget failure is not a bad story proposal. Repeating the same
		// bounded request as a "correction" wastes another model call and can
		// still never produce a complete offer.
		return fmt.Errorf("model exhausted its %d-token response budget; no offer was queued", predictionLimit)
	}
	var proposal core.Proposal
	if e = json.Unmarshal([]byte(answer.Message.Content), &proposal); e != nil {
		return proposalRejected{fmt.Errorf("response must be a complete JSON object matching the proposal schema")}
	}
	if err := repeatedProposal(snapshot, proposal); err != nil {
		return proposalRejected{err}
	}
	if err := validateChoiceScript(proposal); err != nil {
		return proposalRejected{err}
	}
	if err := validateSceneTitle(proposal); err != nil {
		return proposalRejected{err}
	}
	if err := validateSpokenTerms(proposal); err != nil {
		return proposalRejected{err}
	}
	if err := validateInWorldVoice(proposal); err != nil {
		return proposalRejected{err}
	}
	if err := validatePlayerRole(snapshot, proposal); err != nil {
		return proposalRejected{err}
	}
	if err := validateEarnedFamiliarity(snapshot, proposal); err != nil {
		return proposalRejected{err}
	}
	if err := validateSpeakerBeneficiary(snapshot, proposal.Speaker, proposal.Beneficiary); err != nil {
		return proposalRejected{err}
	}
	if err := validateBeneficiaryMention(snapshot, proposal); err != nil {
		return proposalRejected{err}
	}
	if err := validateJobLocation(snapshot, proposal); err != nil {
		return proposalRejected{err}
	}
	allowedSpeaker := false
	for _, id := range directorSpeakers(snapshot, connection) {
		allowedSpeaker = allowedSpeaker || id == proposal.Speaker
	}
	if !allowedSpeaker {
		return proposalRejected{fmt.Errorf("speaker must be one of the available contacts %q", directorSpeakers(snapshot, connection))}
	}
	if env("BLACK_LEDGER_DIRECTOR_BRIEF", "full") == "focused" {
		if err := validateBriefOpening(proposal.Body, jobBrief(snapshot, operation, connection)); err != nil {
			return proposalRejected{err}
		}
	}
	if connection != nil && (proposal.Speaker != connection.Speaker || proposal.Beneficiary != connection.Beneficiary) {
		return proposalRejected{fmt.Errorf("director ignored the established contact connection: speaker must be %q and beneficiary must be %q (empty means neutral), not speaker %q / beneficiary %q. The job location does not decide allegiance", connection.Speaker, connection.Beneficiary, proposal.Speaker, proposal.Beneficiary)}
	}
	if proposal.Operation != operation {
		return proposalRejected{fmt.Errorf("director ignored operation brief: wanted %s", operation)}
	}
	return a.s.Change(func(w *core.World) error {
		if w.ID != snapshot.ID || w.Life != snapshot.Life || !w.Player.Alive {
			return nil
		}
		if err := validateDirectorFreshness(snapshot, w, proposal, connection); err != nil {
			return err
		}
		if err := repeatedProposal(w, proposal); err != nil {
			return proposalRejected{err}
		}
		scene, e := w.ValidateProposal(proposal)
		if e != nil {
			return proposalRejected{e}
		}
		if connection != nil {
			// Reference the completed record, never the model's claimed past outcome.
			for _, m := range w.Arrangements {
				if m.ID == connection.ID && m.Life == w.Life && m.Status == "completed" {
					scene.Connection = &core.StoryConnection{ID: m.ID, Title: m.Title, Result: m.Result}
					break
				}
			}
		}
		w.Offers = append(w.Offers, core.Offer{Ready: w.Minute + 30, Event: scene})
		w.Director.Status = "ready"
		w.Director.Detail = "A local AI encounter is ready. It can arrive after your next activity."
		return nil
	})
}
func main() {
	port := env("BLACK_LEDGER_PORT", "8791")
	s, e := store.Open(env("BLACK_LEDGER_DB", ".runtime/campaign.sqlite3"))
	if e != nil {
		log.Fatal(e)
	}
	defer s.DB.Close()
	_ = s.Change(func(w *core.World) error {
		if w.Director.Status == "writing" {
			w.Director.Status = "available"
			w.Director.Detail = "Previous preparation interrupted. Ready to retry."
		}
		return nil
	})
	a := &app{s: s, port: port, client: &http.Client{}}
	log.Printf("Black Ledger Go core · http://127.0.0.1:%s", port)
	log.Fatal((&http.Server{Addr: "127.0.0.1:" + port, Handler: a, ReadHeaderTimeout: 5 * time.Second}).ListenAndServe())
}
