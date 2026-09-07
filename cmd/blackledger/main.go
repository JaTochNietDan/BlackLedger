package main

import (
	"blackledger/core"
	"blackledger/store"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
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

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

type app struct {
	s      *store.Store
	ai     sync.Mutex
	port   string
	client *http.Client
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
			reply(w, 200, map[string]any{"ok": true, "game": "Black Ledger", "core": "go"})
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
		root := env("BLACK_LEDGER_WEB", "web")
		if _, e := os.Stat(filepath.Join("dist", "index.html")); e == nil {
			root = "dist"
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
	case "/api/speech":
		var q struct {
			Event string `json:"event"`
		}
		if e := body(r, &q); e != nil {
			fail(w, 400, e)
			return
		}
		s, e := a.s.Read()
		if e != nil {
			fail(w, 500, e)
			return
		}
		if s.Event == nil || q.Event != s.Event.ID {
			fail(w, 409, fmt.Errorf("conversation ended"))
			return
		}
		person := s.NPC(s.Event.Speaker)
		if person == nil {
			fail(w, 409, fmt.Errorf("unknown speaker"))
			return
		}
		payload, _ := json.Marshal(map[string]string{"text": s.Event.Body, "speaker": person.Name, "voice": "warm"})
		ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
		defer cancel()
		req, _ := http.NewRequestWithContext(ctx, "POST", env("AFTERLIGHT_DIRECTOR_URL", "http://127.0.0.1:8787")+"/speech", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		res, e := a.client.Do(req)
		if e != nil {
			fail(w, 503, fmt.Errorf("voice service unavailable; you can continue reading"))
			return
		}
		defer res.Body.Close()
		data, e := io.ReadAll(io.LimitReader(res.Body, 8000001))
		if e != nil || res.StatusCode != 200 || len(data) > 8000000 || !bytes.HasPrefix(data, []byte("RIFF")) {
			fail(w, 503, fmt.Errorf("voice not ready; retry or continue reading"))
			return
		}
		w.Header().Set("Content-Type", "audio/wav")
		w.Header().Set("Cache-Control", "private, max-age=3600")
		_, _ = w.Write(data)
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
				}
				return nil
			})
		}
	}()
	return true
}
func (a *app) generate(snapshot *core.World) error {
	contextData := map[string]any{"life": snapshot.Life, "minute": snapshot.Minute, "player": snapshot.Player, "factions": snapshot.Factions, "npcs": snapshot.NPCs, "dead": snapshot.Dead, "places": core.Locations, "properties": snapshot.Properties, "recent_history": snapshot.History[max(0, len(snapshot.History)-12):]}
	b, _ := json.Marshal(contextData)
	payload, _ := json.Marshal(map[string]any{"model": env("BLACK_LEDGER_MODEL", "qwen3:14b"), "stream": false, "think": false, "format": "json", "messages": []map[string]string{{"role": "system", "content": prompt}, {"role": "user", "content": string(b)}}, "options": map[string]any{"temperature": .8, "num_predict": 700}})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if e = json.NewDecoder(io.LimitReader(res.Body, 20000)).Decode(&answer); e != nil {
		return e
	}
	var proposal core.Proposal
	if e = json.Unmarshal([]byte(answer.Message.Content), &proposal); e != nil {
		return e
	}
	return a.s.Change(func(w *core.World) error {
		if w.ID != snapshot.ID || w.Life != snapshot.Life || !w.Player.Alive {
			return nil
		}
		scene, e := w.ValidateProposal(proposal)
		if e != nil {
			return e
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
