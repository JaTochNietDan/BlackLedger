package core

import "testing"

// Information is only a resource if most players do not have all of it. This
// measures what a campaign actually knows about the city it is living in.

func TestMostOfTheCityIsSomethingYouHaveNotAskedAbout(t *testing.T) {
	t.Parallel()
	const runs, days = 200, 60
	levels := map[int]int{}
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		// A campaign that has been played rather than one that has not: some
		// contacts, and some people dealt with along the way.
		w.Player.Contacts = int(seed % 4)
		for i := range w.NPCs {
			if w.NPCs[i].Faction != "" && seed%5 == 0 {
				w.NPCs[i].Trust = 1
				break
			}
		}
		for day := 0; day < days; day++ {
			w.FactionTurn()
			w.Minute += 1440
		}
		for i := range w.Factions {
			levels[w.Intelligence(w.Factions[i].ID)]++
		}
	}
	total := 0
	for _, n := range levels {
		total += n
	}
	if levels[0] == 0 {
		t.Fatal("every campaign knew something about everybody")
	}
	if levels[3]*4 > total {
		t.Fatalf("%d of %d organizations were fully known without anybody asking", levels[3], total)
	}
	t.Logf("%d organizations across %d campaigns: %d known not at all, %d by reputation, %d through somebody inside, %d completely",
		total, runs, levels[0], levels[1], levels[2], levels[3])
}
