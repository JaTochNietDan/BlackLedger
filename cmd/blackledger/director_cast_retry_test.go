package main

import (
	"blackledger/core"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDirectorRejectsUnestablishedSpeakerBeforeSaving(t *testing.T) {
	a := testApp(t)
	calls := 0
	model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		speaker := "vittorio"
		if calls > 1 {
			speaker = "mara"
		}
		proposal, _ := json.Marshal(core.Proposal{Location: "bar", Title: "Shared loading hours", Body: "Let us negotiate shared access to the loading area. At Saint Agnes.", Speaker: speaker, Operation: "mediation", Outcome: "Negotiated access."})
		json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"content": string(proposal)}})
	}))
	defer model.Close()
	t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
	snapshot, _ := a.s.Read()
	if err := a.generate(snapshot); err != nil {
		t.Fatal(err)
	}
	saved, _ := a.s.Read()
	if calls != 2 || len(saved.Offers) != 1 || saved.Offers[0].Event.Speaker != "mara" {
		t.Fatal("unestablished speaker bypassed validation")
	}
}
