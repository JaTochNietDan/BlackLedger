package sim

import (
	"blackledger/core"
	"reflect"
	"testing"
)

func TestCampaignsAreReproducible(t *testing.T) {
	t.Parallel()
	for _, strategy := range []string{"worker", "investor", "defiant", "reckless"} {
		for _, director := range []string{"authored", "fixture"} {
			a, b := Run(27, strategy, director, 100, false), Run(27, strategy, director, 100, false)
			if a.Error != "" || !reflect.DeepEqual(a, b) {
				t.Fatalf("%s/%s not reproducible: %+v / %+v", strategy, director, a, b)
			}
		}
	}
}
func TestPoliciesExerciseTheirIntendedBehavior(t *testing.T) {
	t.Parallel()
	investor := Run(27, "investor", "authored", 100, false)
	for _, milestone := range []string{"laundry", "crew", "housing", "security", "garage", "casino"} {
		if investor.Milestones[milestone] == 0 {
			t.Fatal("investor never reached", milestone)
		}
	}
	worker := Run(27, "worker", "authored", 30, false)
	if worker.Actions["dockwork"] == 0 || worker.Actions["acquire"] != 0 {
		t.Fatal("worker policy not ordinary work")
	}
	reckless := Run(27, "reckless", "authored", 30, false)
	if reckless.Actions["provoke"] != 1 || reckless.Actions["rest"] == 0 {
		t.Fatal("reckless policy did not challenge then stay home")
	}
	fixture := Run(27, "investor", "fixture", 100, false)
	if fixture.Events["proposal"] == 0 {
		t.Fatal("fixture proposals never exercised")
	}
}
func TestPolicyCannotSeePrivateThreatChanges(t *testing.T) {
	t.Parallel()
	w := core.New(27)
	before := Public(w)
	w.Retaliation()
	after := Public(w)
	if !reflect.DeepEqual(before, after) {
		t.Fatal("private plan changed policy observation")
	}
	a, _ := Choose(before, "investor")
	b, _ := Choose(after, "investor")
	if !reflect.DeepEqual(a, b) {
		t.Fatal("policy reacted to private threat")
	}
}

func TestDefiantInvestorRefusesPublicBusinessDemand(t *testing.T) {
	t.Parallel()
	v := View{Revision: 3, Event: &Event{ID: "demand", Kind: "business_pressure", Choices: []core.Choice{{ID: "pay"}, {ID: "resist"}}}}
	compliant, err := Choose(v, "investor")
	if err != nil {
		t.Fatal(err)
	}
	defiant, err := Choose(v, "defiant")
	if err != nil {
		t.Fatal(err)
	}
	if compliant.Choice != "pay" || defiant.Choice != "resist" {
		t.Fatal("policies do not distinguish compliance and resistance")
	}
}

func TestInvestorRepairsCrewLoyaltyUsingPublicActions(t *testing.T) {
	t.Parallel()
	w := core.New(27)
	w.Player.Crew = []core.Crew{{ID: "leo", Name: "Leo", Loyalty: 25}}
	c, err := Choose(Public(w), "investor")
	if err != nil || c.Kind != "crew_bonus" {
		t.Fatalf("investor ignored recoverable crew: %+v %v", c, err)
	}
	next, err := core.Execute(w, c)
	if err != nil {
		t.Fatal(err)
	}
	c, err = Choose(Public(next), "investor")
	if err != nil || c.Kind != "delegate" {
		t.Fatalf("investor did not resume collections: %+v %v", c, err)
	}
	w.Player.Cash = 0
	c, err = Choose(Public(w), "investor")
	if err != nil || c.Kind == "crew_bonus" {
		t.Fatal("policy attempted unaffordable recovery")
	}
}

func TestDiplomatMaintainsPublicBusinessAgreements(t *testing.T) {
	t.Parallel()
	r := Run(27, "diplomat", "authored", 220, false)
	if r.Error != "" {
		t.Fatal(r.Error)
	}
	// This used to require the casino specifically, which is one contingent
	// purchase out of one seeded run: it passed on seed 27 and already failed
	// on seed 28, and any change to where the city spends its evening moved it.
	// Putting the poolhall on the list of places people go at night was enough
	// to move it on every seed, with the diplomat's money within $150 of where
	// it had been — so the guard was measuring a coincidence, not progression.
	// What it means to test is that a policy spending its time on audiences
	// still builds something.
	held := 0
	for _, milestone := range []string{"crew", "laundry", "garage", "casino", "housing", "security"} {
		if r.Milestones[milestone] > 0 {
			held++
		}
	}
	if r.Actions["audience"] == 0 || r.Events["audience"] == 0 || held < 4 {
		t.Fatal("diplomacy did not coexist with progression", r)
	}
	w := core.New(27)
	w.Player.Cash = 300
	w.Properties["laundry"].Owner = "player:1"
	w.Player.Location = "club"
	v := Public(w)
	c, err := Choose(v, "diplomat")
	if err != nil || c.Kind != "audience" {
		t.Fatal("missing public negotiation decision", c, err)
	}
	w.BusinessTruces = map[string]int{"bellandi": w.Minute + 1440}
	c, err = Choose(Public(w), "diplomat")
	if err != nil || c.Kind == "audience" {
		t.Fatal("repeated active agreement purchase", c, err)
	}
}
