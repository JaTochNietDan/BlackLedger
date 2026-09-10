package main

import (
	"blackledger/core"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// A small process-local cache avoids resynthesizing queued encounters and replays.
// Nothing here advances time or changes the saved campaign.
func (a *app) voiceData(ctx context.Context, scene *core.Scene, speaker, profile string) ([]byte, error) {
	a.voiceMu.Lock()
	defer a.voiceMu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	key := scene.ID + "\x00" + speaker + "\x00" + profile + "\x00" + scene.Body
	if data, ok := a.voiceCache[key]; ok {
		return data, nil
	}
	// The voice the core picked, which follows the painting the interface draws.
	// Without it the speech service hashes the name for itself and a detective
	// painted in a trench coat speaks in a woman's voice.
	payload, _ := json.Marshal(map[string]string{
		"text": scene.Body, "speaker": speaker, "voice": "warm", "profile": profile,
	})
	req, _ := http.NewRequestWithContext(ctx, "POST", env("AFTERLIGHT_DIRECTOR_URL", "http://127.0.0.1:8787")+"/speech", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("voice service unavailable; continue reading")
	}
	defer res.Body.Close()
	data, err := io.ReadAll(io.LimitReader(res.Body, 8000001))
	if err != nil || res.StatusCode != 200 || len(data) > 8000000 || !bytes.HasPrefix(data, []byte("RIFF")) {
		return nil, fmt.Errorf("voice not ready; retry or continue reading")
	}
	if a.voiceCache == nil {
		a.voiceCache = map[string][]byte{}
	}
	for len(a.voiceOrder) >= 4 {
		delete(a.voiceCache, a.voiceOrder[0])
		a.voiceOrder = a.voiceOrder[1:]
	}
	a.voiceCache[key] = data
	a.voiceOrder = append(a.voiceOrder, key)
	return data, nil
}

func (a *app) speech(w http.ResponseWriter, r *http.Request) {
	prepare := r.URL.Path == "/api/speech/prepare"
	var q struct {
		Event string `json:"event"`
	}
	if !prepare {
		if err := body(r, &q); err != nil {
			fail(w, 400, err)
			return
		}
	}
	state, err := a.s.Read()
	if err != nil {
		fail(w, 500, err)
		return
	}
	scene := state.Event
	if prepare {
		if !state.Player.Alive || len(state.Offers) == 0 {
			reply(w, 200, map[string]bool{"ready": false})
			return
		}
		scene = state.Offers[0].Event
	} else if scene == nil || q.Event != scene.ID {
		fail(w, 409, fmt.Errorf("conversation ended"))
		return
	}
	if scene == nil {
		fail(w, 409, fmt.Errorf("conversation unavailable"))
		return
	}
	person := state.NPC(scene.Speaker)
	if person == nil {
		fail(w, 409, fmt.Errorf("unknown speaker"))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()
	data, err := a.voiceData(ctx, scene, person.Name, state.VoiceOf(person.ID))
	if err != nil {
		fail(w, 503, err)
		return
	}
	if prepare {
		reply(w, 200, map[string]bool{"ready": true})
		return
	}
	// A decision made during synthesis invalidates even a successfully generated clip.
	current, err := a.s.Read()
	if err != nil || current.Event == nil || current.Event.ID != q.Event {
		fail(w, 409, fmt.Errorf("conversation ended"))
		return
	}
	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Cache-Control", "private, max-age=3600")
	_, _ = w.Write(data)
}
