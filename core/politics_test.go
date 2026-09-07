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
			w.Event, err = w.ValidateProposal(Proposal{"A risky delivery", "Please deliver this sealed package.", "mara", "courier", "Delivered.", "", nil})
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

func TestAllyPressureReportsBothRelationships(t *testing.T) {
	w := pressureWorld()
	w.Advance(30)
	before := w.Player.Cash
	next, err := Execute(w, Command{Kind: "choice", Event: w.Event.ID, Choice: "ally", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Cash != before-35 || next.Factions[0].Goodwill != -12 || next.Factions[1].Goodwill != 12 {
		t.Fatal("incorrect alliance cost or standing")
	}
	record := next.History[len(next.History)-1]
	for _, phrase := range []string{"Russo Outfit", "Bellandi Family", "+12", "-12", "does not guarantee protection"} {
		if !strings.Contains(record.Text, phrase) {
			t.Fatalf("missing political feedback %q: %s", phrase, record.Text)
		}
	}
	if strings.Contains(record.Text, "sabotage") {
		t.Fatal("private retaliation revealed")
	}
}

func TestBookReviewDoesNotAdvanceCity(t *testing.T) {
	w := pressureWorld()
	w.Player.Location = "laundry"
	w.Properties["laundry"].Condition = 85
	next, err := Execute(w, Command{Kind: "inspect", Target: "laundry", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Minute != w.Minute || next.Player.Cash != w.Player.Cash || next.Event != nil {
		t.Fatal("inspection advanced economics or pressure")
	}
	if !strings.Contains(next.History[len(next.History)-1].Text, "$11/hour") {
		t.Fatal("missing actual income")
	}
}

func TestFriendlyHighIncomeDistrictDoesNotSilenceHostileFamily(t *testing.T) {
	w := pressureWorld()
	w.Properties["casino"].Owner = "player:1"
	w.Factions[1].Goodwill = 40
	w.Factions[0].Goodwill = -20
	w.BusinessPressure()
	if w.Event == nil || w.Event.Actor != "bellandi" || w.Event.Target != "laundry" {
		t.Fatal("friendly casino concealed hostile laundry claim")
	}
}
func TestHostileHighIncomeDistrictStillMakesItsOwnClaim(t *testing.T) {
	w := pressureWorld()
	w.Properties["casino"].Owner = "player:1"
	w.Factions[0].Goodwill = 40
	w.Factions[1].Goodwill = -20
	w.BusinessPressure()
	if w.Event == nil || w.Event.Actor != "russo" || w.Event.Target != "casino" {
		t.Fatal("wrong territorial claimant")
	}
}
func TestGoodRelationsInAllOwnedDistrictsPreservePeace(t *testing.T) {
	w := pressureWorld()
	w.Properties["casino"].Owner = "player:1"
	w.Factions[0].Goodwill, w.Factions[1].Goodwill = 25, 25
	w.BusinessPressure()
	if w.Event != nil || len(w.Plots) != 0 || w.NextPressure <= w.Minute {
		t.Fatal("good relationships must still prevent demands")
	}
	w.Factions[0].Goodwill = 24
	w.Advance(720)
	if w.Event == nil || w.Event.Actor != "bellandi" {
		t.Fatal("claim did not resume after standing changed")
	}
}

func TestSustainedDefianceEscalatesForEitherFamily(t *testing.T) {
	for _, actor := range []string{"bellandi", "russo"} {
		w := pressureWorld()
		e := &Scene{Actor: actor, Target: "laundry"}
		if err := w.ResolvePressure(e, "resist"); err != nil {
			t.Fatal(err)
		}
		for _, p := range w.Plots {
			if p.Kind == "hit" {
				t.Fatal("first refusal became a personal hit")
			}
		}
		if err := w.ResolvePressure(e, "resist"); err != nil {
			t.Fatal(err)
		}
		hits := 0
		for _, p := range w.Plots {
			if p.Kind == "hit" {
				hits++
				if p.Actor != actor || p.Due != w.Minute+240 {
					t.Fatal("wrong personal operation")
				}
			}
		}
		if hits != 1 {
			t.Fatal("sustained feud did not escalate")
		}
		w.RetaliationFrom(actor)
		count := 0
		for _, p := range w.Plots {
			if p.Kind == "hit" {
				count++
			}
		}
		if count != 1 {
			t.Fatal("personal operation duplicated")
		}
		data, _ := json.Marshal(w.Public())
		if strings.Contains(string(data), "\"plots\"") {
			t.Fatal("private operation leaked")
		}
	}
}

func TestRussoWarningNamesActualFamilyAndCanBeNegotiated(t *testing.T) {
	w := New(27)
	w.Player.Contacts = 2
	w.RetaliationFrom("russo")
	w.Advance(200)
	if w.Event == nil || w.Event.Kind != "warning" || !strings.Contains(w.Event.Body, "Russo") || strings.Contains(w.Event.Body, "Bellandi") {
		t.Fatal("wrong family in warning")
	}
	w.Event = nil
	w.RetaliationFrom("bellandi")
	if err := w.ResolveAudience(&Scene{Actor: "russo"}, "tribute"); err != nil {
		t.Fatal(err)
	}
	if len(w.Plots) != 1 || w.Plots[0].Actor != "bellandi" {
		t.Fatal("audience cancelled wrong family's threats")
	}
}

func TestFavorCanCreateRussoFeud(t *testing.T) {
	w := New(27)
	w.Factions[1].Goodwill = -28
	w.ResolveBeneficiary("bellandi")
	if len(w.Plots) != 1 || w.Plots[0].Actor != "russo" || w.Plots[0].Kind != "hit" {
		t.Fatal("Russo feud has no consequence")
	}
}
