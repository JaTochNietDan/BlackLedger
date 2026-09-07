package core

import (
	"strings"
	"testing"
)

func approachProposal() Proposal {
	return Proposal{Title: "A ledger before dawn", Body: "Move the ledger before the merchant's partners arrive.", Speaker: "mara", Operation: "courier", Outcome: "Delivered.", Beneficiary: "russo", Approaches: []Approach{{Method: "careful", Label: "Wait for the street to clear"}, {Method: "press", Label: "Make the delivery before closing"}}}
}
func TestApproachChangesRiskAndPaymentThroughPoliceInterruption(t *testing.T) {
	for _, method := range []string{"careful", "press"} {
		t.Run(method, func(t *testing.T) {
			w := New(27)
			w.Player.Heat = 14
			var err error
			w.Event, err = w.ValidateProposal(approachProposal())
			if err != nil {
				t.Fatal(err)
			}
			if len(w.Event.Choices) != 4 || w.Event.Choices[3].ID != "decline" {
				t.Fatal("approaches removed the standard or refusal choice")
			}
			// Clone represents the same JSON persistence boundary used by saved commands.
			w = w.Clone()
			choice(t, &w, "approach:"+method)
			if method == "careful" {
				if w.Event != nil || w.Minute != 555 || w.Player.Cash != 150 || w.Player.Heat != 14 || w.Factions[1].Goodwill != 6 {
					t.Fatal("careful route did not trade time and pay for less heat")
				}
			} else {
				if w.Event == nil || w.Event.Kind != "police_stop" || w.Minute != 510 || w.Player.Cash != 90 || w.Factions[1].Goodwill != 0 {
					t.Fatal("rushed job bypassed police or paid too early")
				}
				choice(t, &w, "pay")
				if w.Player.Cash != 145 || w.Player.Heat != 12 || w.Factions[1].Goodwill != 6 {
					t.Fatal("police lost the selected approach's terms")
				}
			}
		})
	}
}
func TestInvalidDirectorApproachesAreRejected(t *testing.T) {
	for _, options := range [][]Approach{{{Method: "kill", Label: "Attack someone"}}, {{Method: "careful", Label: "Wait"}, {Method: "careful", Label: "Wait again"}}, {{Method: "careful", Label: strings.Repeat("x", 66)}}, {{Method: "careful", Label: "Wait"}, {Method: "press", Label: "Hurry"}, {Method: "press", Label: "Again"}}} {
		p := approachProposal()
		p.Approaches = options
		if _, err := New(27).ValidateProposal(p); err == nil {
			t.Fatal("invalid approach accepted")
		}
	}
	w := New(27)
	p := approachProposal()
	p.Approaches = nil
	w.Event, _ = w.ValidateProposal(p)
	if _, err := Execute(w, Command{Kind: "choice", Event: w.Event.ID, Choice: "approach:press"}); err == nil {
		t.Fatal("unoffered alternative could be forged")
	}
}
func TestInterruptedApproachDoesNotGrantCompletion(t *testing.T) {
	w := New(27)
	w.Properties["laundry"].Owner = "player:1"
	w.NextPressure = 500
	w.Event, _ = w.ValidateProposal(approachProposal())
	choice(t, &w, "approach:careful")
	if w.Minute != 500 || w.Event == nil || w.Event.Kind != "business_pressure" || w.Player.Respect != 0 || w.Factions[1].Goodwill != 0 {
		t.Fatal("interrupted careful work falsely completed")
	}
}
