package main

import (
	"blackledger/core"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestFocusedBriefCorrectionRequiresAcknowledgement(t *testing.T) {
	a := testApp(t)
	if err := a.s.Change(func(w *core.World) error {
		w.Arrangements = []core.ArrangementMemory{{ID: "completed-job", Life: 1, Title: "Tool dispute", Offer: "A disagreement at Russo Motor Works", Speaker: "mara", Operation: "mediation", Status: "completed", Result: "Mediation completed without violence."}}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	var calls atomic.Int32
	model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Format map[string]any `json:"format"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
		}
		if request.Format["type"] != "object" {
			t.Error("focused request omitted response schema")
		}
		n := calls.Add(1)
		body := "Collect a customer's payment and bring it to the manager."
		if n == 2 {
			body = "You handled our last mediation without violence. " + body
		}
		proposal := core.Proposal{Location: "bar", Title: "Customer settlement", Body: body + " At Saint Agnes.", Speaker: "mara", Operation: "collection", Outcome: "Payment collected."}
		b, _ := json.Marshal(proposal)
		json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"content": string(b)}})
	}))
	defer model.Close()
	t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
	t.Setenv("BLACK_LEDGER_DIRECTOR_BRIEF", "focused")
	snapshot, err := a.s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if err = a.generate(snapshot); err != nil {
		t.Fatal(err)
	}
	after, err := a.s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 || len(after.Offers) != 1 || after.Offers[0].Event.Connection == nil {
		t.Fatal("callback correction or preserved link failed")
	}
}
