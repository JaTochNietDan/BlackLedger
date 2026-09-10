package main

import (
	"blackledger/core"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDirectorLocationConstraints(t *testing.T) {
	t.Parallel()
	w := core.New(27)
	for _, location := range []string{"", "missing", "garage"} {
		if validateJobLocation(w, core.Proposal{Location: location}) == nil {
			t.Fatal("invalid venue", location)
		}
	}
	p := core.Proposal{Location: "laundry", Body: "Meet at Russo Motor Works."}
	if validateJobLocation(w, p) == nil {
		t.Fatal("explicit inaccessible site in dialogue accepted")
	}
	p.Body = "Meet at Bluebird Laundry."
	if err := validateJobLocation(w, p); err != nil {
		t.Fatal(err)
	}
	w.District = 1
	p.Location = "garage"
	p.Body = "Meet at Russo Motor Works."
	if err := validateJobLocation(w, p); err != nil {
		t.Fatal("expanded district unavailable", err)
	}
	for _, loc := range proposalSchema(core.New(27), "courier", nil)["properties"].(map[string]any)["location"].(map[string]any)["enum"].([]string) {
		place, _ := core.PlaceByID(loc)
		if place.District > 0 {
			t.Fatal("schema exposes locked venue")
		}
	}
}

func TestStructuredMemoryWinsOverAmbiguousOldDialogue(t *testing.T) {
	t.Parallel()
	w := core.New(27)
	w.District = 1
	m := &core.ArrangementMemory{Location: "laundry", Offer: "A call from Russo Motor Works", Operation: "courier"}
	if jobBrief(w, "collection", m).Location != "Bluebird Laundry" {
		t.Fatal("prose replaced saved venue")
	}
	w.District = 0
	m.Location = "garage"
	if jobBrief(w, "collection", m).Location == "Russo Motor Works" {
		t.Fatal("legacy memory unlocked a district")
	}
}

func TestLocationCorrectionCannotQueueInvalidJob(t *testing.T) {
	for _, badLocation := range []string{"", "garage", "laundry"} {
		t.Run("invalid-"+badLocation, func(t *testing.T) {
			a := testApp(t)
			calls := 0
			model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				p := core.Proposal{Location: badLocation, Title: "A quiet dispute", Body: "Settle a scheduling dispute at Saint Agnes.", Speaker: "mara", Operation: "mediation", Outcome: "Settled."}
				if calls > 1 {
					p.Location = "bar"
				}
				b, _ := json.Marshal(p)
				json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"content": string(b)}})
			}))
			defer model.Close()
			t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
			snapshot, _ := a.s.Read()
			if err := a.generate(snapshot); err != nil {
				t.Fatal(err)
			}
			saved, _ := a.s.Read()
			if calls != 2 || len(saved.Offers) != 1 || saved.Offers[0].Event.Target != "bar" {
				t.Fatal("invalid location escaped correction")
			}
		})
	}
}
