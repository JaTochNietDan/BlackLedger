package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func pressureWorld() *World {
	w := New(27)
	w.Properties["laundry"].Owner = "player:1"
	w.Player.Cash = 200
	w.NextPressure = w.Minute + 30
	return w
}
func TestPressureInterruptsAtScheduledMinute(t *testing.T) {
	w := pressureWorld()
	w.Advance(120)
	if w.Minute != 510 || w.Event == nil || w.Event.Kind != "business_pressure" || w.Event.Target != "laundry" {
		t.Fatal("demand missed boundary or wrong business")
	}
	data, _ := json.Marshal(w.Public())
	if strings.Contains(string(data), "next_pressure") {
		t.Fatal("private schedule disclosed")
	}
}
func TestPressurePaymentIsSpecificNotBlanketProtection(t *testing.T) {
	w := pressureWorld()
	w.Retaliation()
	w.Advance(30)
	before := w.Player.Cash
	next, err := Execute(w, Command{Kind: "choice", Event: w.Event.ID, Choice: "pay", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Cash != before-60 || next.Factions[0].Goodwill != 12 || next.NextPressure != next.Minute+1440 || len(next.Plots) != 1 {
		t.Fatal("payment rules or personal threat changed")
	}
}
func TestRefusingCreatesHiddenConsequences(t *testing.T) {
	w := pressureWorld()
	w.Advance(30)
	next, err := Execute(w, Command{Kind: "choice", Event: w.Event.ID, Choice: "resist", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Factions[0].Goodwill != -20 || next.Player.Respect != 2 || len(next.Plots) != 1 {
		t.Fatal("refusal consequences absent")
	}
	data, _ := json.Marshal(next.Public())
	if strings.Contains(string(data), "sabotage") {
		t.Fatal("hidden retaliation leaked")
	}
	next.Advance(120)
	if next.Properties["laundry"].Condition != 65 || len(next.Plots) != 0 {
		t.Fatal("sabotage did not resolve once")
	}
}
func TestBusinessDamageUsesAvailableCrewNotHomeGuards(t *testing.T) {
	for _, tc := range []struct {
		name string
		crew bool
		busy bool
		want int
	}{{"unprotected", false, false, 65}, {"available crew", true, false, 85}, {"crew on job", true, true, 65}} {
		t.Run(tc.name, func(t *testing.T) {
			w := pressureWorld()
			w.Player.Security = 3
			if tc.crew {
				w.Player.Crew = []Crew{{"leo", "Leo", 65}}
			}
			if tc.busy {
				w.Tasks = []Task{{ID(), "Collections", 999}}
			}
			w.ResolveSabotage(Plot{Target: "laundry", Strength: 35})
			if w.Properties["laundry"].Condition != tc.want {
				t.Fatal("wrong protection source")
			}
		})
	}
}
func TestNoDemandWithoutBusiness(t *testing.T) {
	w := New(27)
	w.NextPressure = 500
	w.Advance(60)
	if w.Event != nil || w.NextPressure != 0 || w.Minute != 540 {
		t.Fatal("business demand without business")
	}
}

func TestBusinessPlotDoesNotPreventPersonalHit(t *testing.T) {
	w := pressureWorld()
	w.Plots = []Plot{{ID: ID(), Kind: "sabotage", Life: w.Life, Actor: "russo", Target: "laundry", Due: 900}}
	w.Retaliation()
	w.Retaliation()
	if len(w.Plots) != 2 || w.Plots[1].Kind != "hit" {
		t.Fatal("different plots incorrectly deduplicated")
	}
}
func TestInvestigationNamesRealActorAndTarget(t *testing.T) {
	w := pressureWorld()
	w.Plots = []Plot{{ID: ID(), Kind: "sabotage", Life: w.Life, Actor: "russo", Target: "laundry", Due: 900}}
	w.Investigate()
	record := w.History[len(w.History)-1]
	if !strings.Contains(record.Text, "Russo Outfit") || !strings.Contains(record.Text, "Bluebird Laundry") || strings.Contains(record.Text, "Bellandi") || !w.Plots[0].Known {
		t.Fatal("investigation invented actor or target")
	}
}
func TestOpportunityDoesNotRevealHiddenPlans(t *testing.T) {
	w := New(27)
	before, _ := json.Marshal(w.NextOpportunity())
	w.Retaliation()
	after, _ := json.Marshal(w.NextOpportunity())
	if string(before) != string(after) {
		t.Fatal("guidance leaked private threat")
	}
	w.Player.JobCount = 3
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["laundry"].Condition = 65
	if w.NextOpportunity().Target != "laundry" {
		t.Fatal("damage has no actionable guidance")
	}
}

func TestPoliceStopDefersJobRewardUntilDecision(t *testing.T) {
	for _, decision := range []string{"pay", "abandon"} {
		t.Run(decision, func(t *testing.T) {
			w := New(27)
			w.Player.Heat = 14
			var err error
			w.Event, err = w.ValidateProposal(Proposal{"A risky delivery", "Please deliver this sealed package.", "mara", "courier", "Delivered.", ""})
			if err != nil {
				t.Fatal(err)
			}
			choice(t, &w, "accept")
			if w.Event == nil || w.Event.Kind != "police_stop" || w.Player.Cash != 90 || w.Minute != 525 || w.Player.Respect != 0 {
				t.Fatal("police stop failed to defer reward")
			}
			oldID := w.Event.ID
			choice(t, &w, decision)
			if decision == "pay" && (w.Player.Cash != 125 || w.Player.Respect != 3 || w.Player.Heat != 7) {
				t.Fatal("paid completion incorrect")
			}
			if decision == "abandon" && (w.Player.Cash != 90 || w.Player.Respect != 0 || w.Player.Heat != 8) {
				t.Fatal("abandoned job paid reward")
			}
			if _, err = Execute(w, Command{Kind: "choice", Event: oldID, Choice: decision, Revision: w.Revision}); err == nil {
				t.Fatal("police choice applied twice")
			}
		})
	}
}

func TestGeneratedFactionWorkHasValidatedPoliticalEffect(t *testing.T) {
	w := New(27)
	if _, err := w.ValidateProposal(Proposal{Title: "A job", Body: "Help them.", Speaker: "mara", Operation: "courier", Outcome: "Done.", Beneficiary: "invented-family"}); err == nil {
		t.Fatal("unknown faction accepted")
	}
	scene, err := w.ValidateProposal(Proposal{Title: "A Russo favor", Body: "Deliver these papers for Russo.", Speaker: "mara", Operation: "courier", Outcome: "Done.", Beneficiary: "russo"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(scene.Choices[0].Detail, "Russo Outfit standing +6") {
		t.Fatal("political stakes hidden")
	}
	w.Event = scene
	choice(t, &w, "accept")
	if w.Factions[1].Goodwill != 6 || w.Factions[0].Goodwill != -3 {
		t.Fatal("completed faction job has no politics")
	}
	w.Player.Heat = 14
	w.Event = scene
	w.Event.ID = ID()
	choice(t, &w, "accept")
	if w.Event == nil || w.Event.Kind != "police_stop" || w.Event.Beneficiary != "russo" {
		t.Fatal("political context lost on interruption")
	}
	before := w.Factions[1].Goodwill
	choice(t, &w, "abandon")
	if w.Factions[1].Goodwill != before {
		t.Fatal("abandoned work granted favor")
	}
}
