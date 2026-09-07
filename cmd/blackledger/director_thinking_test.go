package main

import (
	"blackledger/core"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDirectorThinkingIsExplicitAndStillQueuesValidatedProposal(t *testing.T) {
	for _, mode := range []string{"", "1"} {
		t.Run("mode="+mode, func(t *testing.T) {
			a := testApp(t)
			snapshot, err := a.s.Read()
			if err != nil {
				t.Fatal(err)
			}
			model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request struct {
					Think   bool `json:"think"`
					Options struct {
						Limit int `json:"num_predict"`
					} `json:"options"`
				}
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					t.Error(err)
				}
				limit := 700
				if mode == "1" {
					limit = 4096
				}
				if request.Think != (mode == "1") || request.Options.Limit != limit {
					t.Errorf("unexpected inference options: %+v", request)
				}
				b, _ := json.Marshal(core.Proposal{Location: "bar", Title: "A private delivery", Body: "Carry a sealed parcel at Saint Agnes.", Speaker: "mara", Operation: snapshot.NextDirectorOperation(), Outcome: "Delivered."})
				// Only the final content is parsed. Reasoning is not dialogue.
				draft := "Private draft analysis, not player dialogue."
				if mode == "1" {
					draft = strings.Repeat("draft ", 5000)
				}
				json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"thinking": draft, "content": string(b)}})
			}))
			defer model.Close()
			t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
			t.Setenv("BLACK_LEDGER_DIRECTOR_THINK", mode)
			if err := a.generate(snapshot); err != nil {
				t.Fatal(err)
			}
			current, err := a.s.Read()
			if err != nil {
				t.Fatal(err)
			}
			if len(current.Offers) != 1 || current.Offers[0].Event.Body != "Carry a sealed parcel at Saint Agnes." || current.Minute != snapshot.Minute || current.Player.Cash != snapshot.Player.Cash {
				t.Fatal("reasoning leaked, proposal missing or generation advanced gameplay")
			}
		})
	}
}
