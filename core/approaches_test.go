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

// The scene modal asks the player to compare two or three ways of doing the
// same job. The terms of each were written into one prose string —
// "$140 · 75 minutes · +6 respect · +4 heat · At 15 heat, police may stop
// completion." — so the only decision the scene exists to pose was the one
// thing on it that could not be scanned. The numbers belong on the choice as
// numbers, and they have to be the numbers the command will actually apply.
func TestEveryWayOfDoingAJobStatesItsTermsAsFigures(t *testing.T) {
	w := New(27)
	scene, err := w.ValidateProposal(approachProposal())
	if err != nil {
		t.Fatal(err)
	}
	terms := map[string]Choice{}
	for _, c := range scene.Choices {
		terms[c.ID] = c
	}
	for _, id := range []string{"accept", "approach:careful", "approach:press"} {
		c, ok := terms[id]
		if !ok {
			t.Fatalf("%s was not offered", id)
		}
		if c.Minutes <= 0 || c.Pay <= 0 {
			t.Fatalf("%s offers $%d for %d minutes", id, c.Pay, c.Minutes)
		}
		// The figures on the button must be the figures the command applies,
		// or the scene is quoting a price the city will not honour.
		want := scene.Effect
		if fx, ok := scene.Alternatives[id]; ok {
			want = fx
		}
		if c.Pay != want.Reward || c.Minutes != want.Minutes || c.Respect != want.Respect || c.Heat != want.Heat {
			t.Fatalf("%s advertises $%d/%dm/+%d/+%d and the job pays $%d/%dm/+%d/+%d",
				id, c.Pay, c.Minutes, c.Respect, c.Heat, want.Reward, want.Minutes, want.Respect, want.Heat)
		}
	}
	// And the comparison the modal is for has to be real in both directions.
	careful, press := terms["approach:careful"], terms["approach:press"]
	if !(press.Pay > careful.Pay && press.Heat > careful.Heat && press.Minutes < careful.Minutes) {
		t.Fatalf("pressing pays $%d/%d heat/%dm and being careful pays $%d/%d heat/%dm: the choice is not a trade",
			press.Pay, press.Heat, press.Minutes, careful.Pay, careful.Heat, careful.Minutes)
	}
	// The prose that is left should be the part that is not a number.
	if strings.Contains(terms["accept"].Detail, "$") || strings.Contains(terms["accept"].Detail, "minutes") {
		t.Fatalf("the figures are still buried in prose: %q", terms["accept"].Detail)
	}
}

// A choice the player cannot take says so in the city's words. The interface
// used to append "Not enough cash" to every disabled choice, which was a guess
// that happened to be right; the core is the only thing that knows.
func TestAChoiceOutOfReachSaysWhyItself(t *testing.T) {
	w := New(27)
	var err error
	w.Event, err = w.ValidateProposal(approachProposal())
	if err != nil {
		t.Fatal(err)
	}
	w.Event.Choices = append(w.Event.Choices, Choice{ID: "buy", Label: "Put money on the table", Cost: 5000})
	w.Player.Cash = 40
	choices := w.Public()["event"].(map[string]any)["choices"].([]Choice)
	var buy Choice
	for _, c := range choices {
		if c.ID == "buy" {
			buy = c
		}
	}
	if !buy.Disabled {
		t.Fatal("a $5,000 choice was offered to somebody holding $40")
	}
	if buy.Reason == "" {
		t.Fatal("the choice is refused and says nothing about why")
	}
	for _, c := range choices {
		if !c.Disabled && c.Reason != "" {
			t.Fatalf("%q is available and carries the refusal %q", c.ID, c.Reason)
		}
	}
}

// Where the job is, what the police will do at 15 heat, and who gains standing
// by it are true of the job however it is done. They were printed on every
// button, so a scene with three ways of doing one thing repeated the same
// ninety characters three times and buried the words that actually differed.
// Conditions belong to the scene; only what changes belongs on a choice.
func TestWhatIsTrueOfTheJobIsSaidOnceNotOnEveryButton(t *testing.T) {
	w := New(27)
	proposal := approachProposal()
	proposal.Location = "bar"
	scene, err := w.ValidateProposal(proposal)
	if err != nil {
		t.Fatal(err)
	}
	if scene.Conditions == "" {
		t.Fatal("the scene states no conditions at all")
	}
	for _, want := range []string{"Saint Agnes", "police", "standing"} {
		if !strings.Contains(scene.Conditions, want) {
			t.Fatalf("the conditions do not mention %q: %q", want, scene.Conditions)
		}
	}
	for _, c := range scene.Choices {
		for _, moved := range []string{"Saint Agnes", "police may stop", "standing +6"} {
			if strings.Contains(c.Detail, moved) {
				t.Fatalf("%q still repeats %q on the button: %q", c.ID, moved, c.Detail)
			}
		}
	}
}
