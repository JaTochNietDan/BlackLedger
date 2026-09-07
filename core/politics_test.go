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
