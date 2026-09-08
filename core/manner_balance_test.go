package core

import "testing"

// A city that kills a dozen people over a season should have described a dozen
// different deaths. This measures how much of the vocabulary the simulation
// actually reaches on its own.

func TestACityDescribesItsDeathsDifferently(t *testing.T) {
	const runs, days = 300, 120
	descriptions := map[string]int{}
	deaths := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		for day := 0; day < days; day++ {
			w.FactionTurn()
			w.FamilyDay()
			w.PeopleDay()
			w.GrudgeDay()
			w.SettleGrudges()
			w.Minute += 1440
		}
		for _, r := range w.History {
			if r.Kind == "danger" && len(r.Title) > 9 && r.Title[len(r.Title)-8:] == " is dead" {
				deaths++
				descriptions[r.Text]++
			}
		}
	}
	if deaths == 0 {
		t.Fatalf("nobody died in %d cities over %d days", runs, days)
	}
	if len(descriptions) < 12 {
		t.Fatalf("%d deaths produced %d different descriptions", deaths, len(descriptions))
	}
	worst := 0
	for _, n := range descriptions {
		worst = max(worst, n)
	}
	if worst*3 > deaths {
		t.Fatalf("one description covered %d of %d deaths", worst, deaths)
	}
	t.Logf("%d cities over %d days: %d deaths in %d distinct descriptions, the commonest of them used %d times", runs, days, deaths, len(descriptions), worst)
}
