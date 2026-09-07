package sim

import (
	"blackledger/core"
	"reflect"
	"testing"
)

func TestCampaignsAreReproducible(t *testing.T) {
	for _, strategy := range []string{"worker", "investor", "reckless"} {
		for _, director := range []string{"authored", "fixture"} {
			a, b := Run(27, strategy, director, 100, false), Run(27, strategy, director, 100, false)
			if a.Error != "" || !reflect.DeepEqual(a, b) {
				t.Fatalf("%s/%s not reproducible: %+v / %+v", strategy, director, a, b)
			}
		}
	}
}
func TestPoliciesExerciseTheirIntendedBehavior(t *testing.T) {
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
