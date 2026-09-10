package core

import (
	"fmt"
	"testing"
)

// An arrangement with a man who has a salary is a standing cost against a
// standing benefit, so the only honest question is where it starts paying.

func TestAMayorOnlyPaysForHimselfWithAPortfolio(t *testing.T) {
	heavy(t)
	const days = 60
	run := func(places []string, retained bool) int {
		w := New(37)
		w.MigrateLivingWorld()
		w.Player.Location, w.Player.Cash, w.Player.Respect = CityHall, 40000, 40
		for _, id := range places {
			w.Properties[id].Owner = fmt.Sprintf("player:%d", w.Life)
		}
		if retained {
			if err := w.Retain("mayor"); err != nil {
				t.Fatal(err)
			}
		}
		start := w.Player.Cash
		w.Advance(days * 1440)
		return w.Player.Cash - start
	}

	one := []string{"laundry"}
	all := []string{"laundry", "garage", "casino"}
	if run(one, true) >= run(one, false) {
		t.Fatal("a mayor paid for himself on the strength of one laundry, so there is no decision here")
	}
	if run(all, true) <= run(all, false) {
		t.Fatalf("a mayor cost more than he was worth across three businesses: $%d against $%d", run(all, true), run(all, false))
	}
	t.Logf("over %d days: one laundry nets $%d without the mayor and $%d with him; three businesses net $%d and $%d",
		days, run(one, false), run(one, true), run(all, false), run(all, true))
}

func TestACommissionerBuysYouRoomToRunHot(t *testing.T) {
	heavy(t)
	const runs, days = 300, 30
	raided := func(retained bool) int {
		hit := 0
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(seed)
			w.MigrateLivingWorld()
			w.WorldRNG = seed * 2654435761
			w.Player.Location, w.Player.Cash, w.Player.Respect = CityHall, 40000, 40
			w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
			if retained {
				if err := w.Retain("commissioner"); err != nil {
					t.Fatal(err)
				}
			}
			for day := 0; day < days; day++ {
				w.Player.Heat = 60 // held there deliberately
				w.PoliceDay()
				if w.hasRecord("Turned over at Bluebird Laundry") || w.hasRecord("They came to the door") {
					hit++
					break
				}
			}
		}
		return hit
	}
	exposed, covered := raided(false), raided(true)
	if covered >= exposed {
		t.Fatalf("a commissioner was raided %d times in %d against %d without one", covered, runs, exposed)
	}
	t.Logf("held at 60 attention for %d days across %d campaigns: raided %d times without a commissioner, %d with one", days, runs, exposed, covered)
}
