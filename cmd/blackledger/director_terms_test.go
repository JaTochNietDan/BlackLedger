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

func TestSpokenTermsKeepGeneratedDialogueConsistentWithChoiceDetails(t *testing.T) {
	for _, text := range []string{
		"Collect $25 from the mechanic; leave it before noon.", "Collect twenty-five dollars.", "Bring USD 25.", "Pay a hundred bucks.", "A €50 favor.",
		"Finish in fifteen minutes.", "Take 2 hours.", "Wait one-day for a response.",
		"Bring it before sundown.", "Before the market closes, collect it.", "Finish before the week ends.", "Be there by 4:30.", "Be there at 5 pm.", "Meet at five o'clock.", "I need the money before the end of the day.", "Finish by close of business.", "Bring 25 USD.",
	} {
		t.Run(text, func(t *testing.T) {
			if validateSpokenTerms(core.Proposal{Body: text}) == nil {
				t.Fatal("unsupported terms accepted")
			}
		})
	}
	for _, text := range []string{
		"At Pier 14, negotiate between two crews.", "They disagree about night shift hours.", "Listen before making a proposal.", "Check the sealed envelope before delivery.", "Collect the agreed payment at Russo Motor Works.", "The grand entrance is blocked.",
	} {
		if err := validateSpokenTerms(core.Proposal{Body: text}); err != nil {
			t.Fatal("ordinary task rejected", text, err)
		}
	}
	if validateSpokenTerms(core.Proposal{Approaches: []core.Approach{{Method: "careful", Label: "Wait until midnight"}}}) == nil {
		t.Fatal("choice deadline accepted")
	}
	if validateSpokenTerms(core.Proposal{Title: "The $25 favor"}) == nil {
		t.Fatal("title amount accepted")
	}
}

func TestUnsupportedTermsCorrectOnceWithoutChangingCampaign(t *testing.T) {
	for _, corrected := range []bool{true, false} {
		t.Run(map[bool]string{true: "corrected", false: "still invalid"}[corrected], func(t *testing.T) {
			a := testApp(t)
			var calls atomic.Int32
			model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				n := calls.Add(1)
				if n == 2 {
					var request struct {
						Messages []struct {
							Content string `json:"content"`
						} `json:"messages"`
					}
					if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
						t.Error(err)
					}
					if len(request.Messages) != 2 || !strings.Contains(request.Messages[1].Content, "unsupported monetary amount") {
						t.Error("missing actionable correction")
					}
				}
				body := "At Saint Agnes, mediate the dispute for $25 before noon."
				if corrected && n == 2 {
					body = "At Saint Agnes, help two staff members agree on access to a workbench. Hear both sides quietly."
				}
				b, _ := json.Marshal(core.Proposal{Title: "A sealed request", Location: "bar", Speaker: "mara", Operation: "mediation", Body: body, Outcome: "Delivered."})
				json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"content": string(b)}})
			}))
			defer model.Close()
			t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
			before, _ := a.s.Read()
			err := a.generate(before)
			after, _ := a.s.Read()
			if (err == nil) != corrected || calls.Load() != 2 {
				t.Fatal("correction not bounded", err, calls.Load())
			}
			expected := 0
			if corrected {
				expected = 1
			}
			if len(after.Offers) != expected || after.Minute != before.Minute || after.Player.Cash != before.Player.Cash {
				t.Fatal("invalid queue or changed campaign")
			}
		})
	}
}
