package main

import (
	"blackledger/core"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestDirectorCorrectionIsBoundedAndValidated(t *testing.T) {
	for _, scenario := range []string{"corrected", "faction-corrected", "still-invalid", "service-down"} {
		t.Run(scenario, func(t *testing.T) {
			a := testApp(t)
			var calls atomic.Int32
			model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				number := calls.Add(1)
				if scenario == "service-down" {
					w.WriteHeader(503)
					return
				}
				var input struct {
					Format   map[string]any `json:"format"`
					Messages []struct {
						Content string `json:"content"`
					} `json:"messages"`
				}
				json.NewDecoder(r.Body).Decode(&input)
				if input.Format["type"] != "object" {
					t.Error("default director request omitted schema")
				}
				if number == 2 && !strings.Contains(input.Messages[1].Content, "failed validation") {
					t.Error("correction lacked feedback")
				}
				var contextData map[string]json.RawMessage
				if err := json.Unmarshal([]byte(input.Messages[1].Content), &contextData); err != nil {
					t.Fatal(err)
				}
				if string(contextData["allowed_beneficiary_ids"]) != `["","bellandi","russo"]` || string(contextData["current_clock"]) != `"Day 1, 08:00"` {
					t.Error("missing explicit valid IDs or readable clock")
				}
				beneficiary := ""
				if scenario == "faction-corrected" && number == 1 {
					beneficiary = "neutral"
				}
				if scenario == "faction-corrected" && number == 2 && !strings.Contains(input.Messages[1].Content, "use exactly one of") {
					t.Error("missing actionable faction correction")
				}
				label := strings.Repeat("overlong ", 12)
				if (scenario == "corrected" && number == 2) || scenario == "faction-corrected" {
					label = "Hear both sides first"
				}
				proposal, _ := json.Marshal(core.Proposal{Location: "bar", Title: "Shared hours", Body: "Help us negotiate the laundry schedule. At Saint Agnes.", Speaker: "mara", Beneficiary: beneficiary, Operation: "mediation", Outcome: "Agreed.", Approaches: []core.Approach{{Method: "careful", Label: label}}})
				json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"content": string(proposal)}})
			}))
			defer model.Close()
			t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
			before, _ := a.s.Read()
			err := a.generate(before)
			after, _ := a.s.Read()
			if scenario == "corrected" || scenario == "faction-corrected" {
				if err != nil || len(after.Offers) != 1 || calls.Load() != 2 {
					t.Fatal("valid correction failed", err, calls.Load())
				}
			} else {
				expected := int32(2)
				if scenario == "service-down" {
					expected = 1
				}
				if err == nil || len(after.Offers) != 0 || calls.Load() != expected {
					t.Fatal("invalid response accepted or retry unbounded")
				}
			}
			if after.Minute != before.Minute || after.Revision != before.Revision || after.Player.Cash != before.Player.Cash {
				t.Fatal("model correction changed gameplay")
			}
		})
	}
}
