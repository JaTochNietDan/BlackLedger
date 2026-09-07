package main

import (
	"blackledger/core"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestQueuedVoicePreparationIsPrivateAndReused(t *testing.T) {
	a := testApp(t)
	var calls atomic.Int32
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		var payload map[string]string
		json.NewDecoder(r.Body).Decode(&payload)
		if payload["speaker"] != "Mara Bell" || payload["text"] != "Negotiate the laundry schedule." {
			t.Error("unstable speaker or text", payload)
		}
		w.Write([]byte("RIFFtest-wave"))
	}))
	defer provider.Close()
	t.Setenv("AFTERLIGHT_DIRECTOR_URL", provider.URL)
	var id string
	a.s.Change(func(w *core.World) error {
		scene, err := w.ValidateProposal(core.Proposal{Title: "Laundry schedule", Body: "Negotiate the laundry schedule.", Speaker: "mara", Operation: "mediation", Outcome: "Done."})
		if err != nil {
			return err
		}
		id = scene.ID
		w.Offers = append(w.Offers, core.Offer{Ready: w.Minute, Event: scene})
		return nil
	})
	before, _ := a.s.Read()
	response := request(a, "POST", "/api/speech/prepare", "{}")
	if response.Code != 200 || bytes.Contains(response.Body.Bytes(), []byte("RIFF")) {
		t.Fatal("queued audio exposed", response.Body.String())
	}
	after, _ := a.s.Read()
	if before.Minute != after.Minute || before.Revision != after.Revision {
		t.Fatal("preparation changed gameplay")
	}
	response = request(a, "POST", "/api/speech", `{"event":"`+id+`"}`)
	if response.Code != 409 {
		t.Fatal("queued clip available before encounter")
	}
	a.s.Change(func(w *core.World) error { w.OfferIfReady(); return nil })
	for i := 0; i < 2; i++ {
		response = request(a, "POST", "/api/speech", `{"event":"`+id+`"}`)
		if response.Code != 200 || !bytes.HasPrefix(response.Body.Bytes(), []byte("RIFF")) {
			t.Fatal("prepared audio unavailable")
		}
	}
	if calls.Load() != 1 {
		t.Fatal("replay synthesized again", calls.Load())
	}
	a.s.Change(func(w *core.World) error { w.Event = nil; return nil })
	if request(a, "POST", "/api/speech", `{"event":"`+id+`"}`).Code != 409 {
		t.Fatal("cached ended conversation accessible")
	}
}

func TestVoiceResponseRejectedWhenDecisionEndsConversation(t *testing.T) {
	a := testApp(t)
	started := make(chan struct{})
	release := make(chan struct{})
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		w.Write([]byte("RIFFtest-wave"))
	}))
	defer provider.Close()
	t.Setenv("AFTERLIGHT_DIRECTOR_URL", provider.URL)
	a.s.Change(func(w *core.World) error {
		w.Event = &core.Scene{ID: "active", Speaker: "mara", Body: "A warning."}
		return nil
	})
	finished := make(chan *httptest.ResponseRecorder, 1)
	go func() { finished <- request(a, "POST", "/api/speech", `{"event":"active"}`) }()
	<-started
	a.s.Change(func(w *core.World) error { w.Event = nil; return nil })
	close(release)
	response := <-finished
	if response.Code != 409 || bytes.Contains(response.Body.Bytes(), []byte("RIFF")) {
		t.Fatal("late audio escaped")
	}
}
