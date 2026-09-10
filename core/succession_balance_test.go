package core

import "testing"

// A city that eats its own leadership every week is not dramatic, it is
// unstable. This measures a season of cities running on their own.

func TestUpheavalIsAnEventNotAClimate(t *testing.T) {
	t.Parallel()
	const runs, days = 200, 120
	coups, cities, collapsed, spawned := 0, 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		started := len(w.Factions)
		seen := 0
		// FactionTurn is deliberately not called: it rolls the same chance
		// itself, and counting both would report twice the rate the city
		// actually runs at.
		for day := 0; day < days; day++ {
			for i := range w.Factions {
				if w.movedThisDay(&w.Factions[i]) {
					seen++
				}
			}
			w.FamilyDay()
			w.PeopleDay()
			w.GrudgeDay()
			w.SettleGrudges()
			w.Minute += 1440
		}
		coups += seen
		if seen > 0 {
			cities++
		}
		if len(w.Factions) > started {
			spawned++
		}
		holding := 0
		for i := range w.Factions {
			if len(w.FamilyHoldings(w.Factions[i].ID)) > 0 {
				holding++
			}
		}
		if holding < 2 {
			collapsed++
		}
	}
	if cities == 0 {
		t.Fatalf("no organization in %d cities ever turned on its own leader in %d days", runs, days)
	}
	if cities == runs {
		t.Fatal("every city had a coup, so upheaval is the weather")
	}
	if collapsed*3 > runs {
		t.Fatalf("%d of %d cities ended with fewer than two organizations holding anything", collapsed, runs)
	}
	t.Logf("%d of %d cities saw an organization turn on its own leadership in %d days (%d moves in total); %d ended with more organizations than they started with, and %d with fewer than two holding anything",
		cities, runs, days, coups, spawned, collapsed)
}
