package main

import (
	"blackledger/core"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDirectorUsesRememberedOutcomeAndRejectsRepeatedOperation(t *testing.T) {
	for _, wrong := range []bool{false, true} {
		t.Run(map[bool]string{false: "correct brief", true: "ignored brief"}[wrong], func(t *testing.T) {
			a := testApp(t)
			if err := a.s.Change(func(w *core.World) error {
				scene, err := w.ValidateProposal(core.Proposal{Location: "bar", Title: "An old delivery", Body: "Deliver the venue's booking book. At Saint Agnes.", Speaker: "mara", Operation: "courier", Outcome: "Done."})
				if err != nil {
					return err
				}
				w.RememberArrangement(scene, "declined")
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req struct {
					Messages []struct {
						Content string `json:"content"`
					} `json:"messages"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Error(err)
					return
				}
				var context struct {
					Operation string                   `json:"required_operation"`
					Memory    []core.ArrangementMemory `json:"recent_arrangements"`
				}
				if len(req.Messages) != 2 {
					t.Error("missing context")
					return
				}
				if err := json.Unmarshal([]byte(req.Messages[1].Content), &context); err != nil {
					t.Error(err)
					return
				}
				if context.Operation != "mediation" || len(context.Memory) != 1 || context.Memory[0].Status != "declined" || context.Memory[0].Result != "" {
					t.Error("missing truthful memory", context)
				}
				operation := context.Operation
				if wrong {
					operation = "courier"
				}
				proposal, _ := json.Marshal(core.Proposal{Location: "bar", Title: "Shared loading hours", Body: "Negotiate shared loading hours at the laundry. At Saint Agnes.", Speaker: "mara", Operation: operation, Outcome: "Negotiated."})
				json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"content": string(proposal)}})
			}))
			defer model.Close()
			t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
			snapshot, _ := a.s.Read()
			err := a.generate(snapshot)
			if (err != nil) != wrong {
				t.Fatal("operation validation result", err)
			}
			saved, _ := a.s.Read()
			expected := 1
			if wrong {
				expected = 0
			}
			if len(saved.Offers) != expected {
				t.Fatal("invalid offer entered queue")
			}
		})
	}
}
