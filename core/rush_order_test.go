package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRushOrderPaysQuotedFeeAndReservesThePremisesOrder(t *testing.T) {
	w := New(81)
	w.Player.Location = "laundry"
	w.Event = nil
	w.Plots = nil
	w.Tasks = nil
	initialCash, initialSupply := w.Player.Cash, w.Properties["laundry"].Supply
	var offered Action
	for _, a := range w.Actions("laundry") {
		if a.ID == "rushorder" {
			offered = a
		}
	}
	if offered.Disabled || offered.Minutes != RushOrderMinutes || offered.Group != "work" || !strings.Contains(offered.Detail, "$80") {
		t.Fatalf("wrong terms: %+v", offered)
	}
	next, err := Execute(w, Command{Kind: "rushorder", Target: "laundry", RequestID: ID(), Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Cash != initialCash+RushOrderPay || next.Minute != w.Minute+RushOrderMinutes || next.Properties["laundry"].Supply != initialSupply-RushOrderSupply || next.Player.Respect != 1 {
		t.Fatalf("incorrect payment/time/resources: %+v", next.Player)
	}
	if _, err = Execute(next, Command{Kind: "rushorder", Target: "laundry", RequestID: ID(), Revision: next.Revision}); err == nil {
		t.Fatal("daily order paid twice")
	}
	decoded := next.Clone()
	if decoded.RushOrderReadiness("laundry") == "" {
		t.Fatal("save roundtrip reset quota")
	}
	decoded.Life++ // The reservation belongs to the premises, not the life.
	if decoded.RushOrderReadiness("laundry") == "" {
		t.Fatal("another protagonist took same order")
	}
	decoded.Minute += 1440
	if why := decoded.RushOrderReadiness("laundry"); why != "" {
		t.Fatal("next day's order unavailable:", why)
	}
}

func TestRushOrderRefusesClosedUnstaffedUnstockedOwnedAndWrongPlaces(t *testing.T) {
	for _, scenario := range []string{"too-early", "too-late", "unstaffed", "supply", "damaged", "trouble", "owned", "wrong-place"} {
		t.Run(scenario, func(t *testing.T) {
			w := New(81)
			w.Event = nil
			w.Player.Location = "laundry"
			target := "laundry"
			switch scenario {
			case "too-early":
				w.Minute = 7 * 60
			case "too-late":
				w.Minute = 17*60 + 1
			case "unstaffed":
				w.Properties[target].Staff = 0
			case "supply":
				w.Properties[target].Supply = 1
			case "damaged":
				w.Properties[target].Condition = 39
			case "trouble":
				w.Properties[target].Trouble = true
			case "owned":
				w.Properties[target].Owner = "player:1"
			case "wrong-place":
				target = "bar"
				w.Player.Location = "bar"
			}
			before, _ := json.Marshal(w)
			if _, err := Execute(w, Command{Kind: "rushorder", Target: target, RequestID: ID(), Revision: w.Revision}); err == nil {
				t.Fatal("invalid order accepted")
			}
			after, _ := json.Marshal(w)
			if string(before) != string(after) {
				t.Fatal("refusal changed world")
			}
		})
	}
	w := New(81)
	w.Minute = 17 * 60
	if why := w.RushOrderReadiness("laundry"); why != "" {
		t.Fatal("last full hour refused:", why)
	}
	w.Player.Location = "room"
	if _, err := Execute(w, Command{Kind: "rushorder", Target: "laundry", RequestID: ID(), Revision: w.Revision}); err == nil {
		t.Fatal("job accepted without visiting")
	}
}

func TestRushOrderInterruptedByWarningDoesNotPayOrReleaseReservation(t *testing.T) {
	w := New(27)
	w.Player.Location = "laundry"
	w.Player.Contacts = 2
	w.Event = nil
	w.Plots = []Plot{{ID: "rush-danger", Kind: "hit", Actor: "bellandi", Life: w.Life, Due: w.Minute + 120}}
	next, err := Execute(w, Command{Kind: "rushorder", Target: "laundry", RequestID: ID(), Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Event == nil || next.Event.Kind != "warning" {
		t.Fatal("fixture did not interrupt work")
	}
	if next.Player.Cash != w.Player.Cash || next.Properties["laundry"].RushOrderDay == 0 || next.Properties["laundry"].Supply != w.Properties["laundry"].Supply-RushOrderSupply {
		t.Fatal("interruption paid or unreserved work")
	}
}
