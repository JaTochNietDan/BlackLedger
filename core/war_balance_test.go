package core

import "testing"

// The living world is only worth having if it produces varied outcomes. A city
// where war is certain is as dead as one where it is impossible, and a war that
// always annihilates someone removes the world the next life inherits. These
// bounds are deliberately wide: they catch a tuning change that makes the city
// degenerate, not small drift.
func TestCityConflictStaysVaried(t *testing.T) {
	const campaigns, days = 200, 180
	sawWar, seizures, wiped := 0, 0, 0
	emptied, grew, born := 0, 0, 0
	warBy := map[int]int{15: 0, 30: 0}
	for seed := 1; seed <= campaigns; seed++ {
		w := New(uint32(seed))
		firstWar := 0
		for day := 1; day <= days; day++ {
			owners := map[string]string{}
			for id, prop := range w.Properties {
				owners[id] = prop.Owner
			}
			w.Minute += 720
			w.FactionTurn()
			w.Minute += 720
			w.FactionTurn()
			w.FamilyDay()
			for id, prop := range w.Properties {
				if owners[id] != prop.Owner {
					seizures++
				}
			}
			if firstWar == 0 {
				for _, c := range w.Conflicts {
					if c.State == "war" {
						firstWar = day
					}
				}
			}
		}
		if firstWar > 0 {
			sawWar++
			for limit := range warBy {
				if firstWar <= limit {
					warBy[limit]++
				}
			}
		}
		viable := 0
		for _, f := range w.Factions {
			if len(w.FamilyHoldings(f.ID)) > 0 {
				viable++
			} else {
				wiped++
			}
		}
		if viable < 2 {
			emptied++
		}
		if len(w.Factions) > 2 {
			grew++
		}
		born += len(w.Factions) - 2
	}
	t.Logf("%d campaigns over %d days: war in %d, by day15=%d day30=%d, seizures=%d, holding nothing=%d, new organizations=%d in %d campaigns, cities left with fewer than two organizations=%d",
		campaigns, days, sawWar, warBy[15], warBy[30], seizures, wiped, born, grew, emptied)

	if sawWar < campaigns/10 {
		t.Fatalf("the city almost never goes to war (%d of %d); nothing happens without the player", sawWar, campaigns)
	}
	if sawWar > campaigns*9/10 {
		t.Fatalf("war is effectively certain (%d of %d); it stops being a surprise", sawWar, campaigns)
	}
	if warBy[15] > campaigns/2 {
		t.Fatalf("war arrives too early to be provoked or avoided (%d of %d by day 15)", warBy[15], campaigns)
	}
	if seizures == 0 {
		t.Fatal("wars never change who holds anything, so they cost nobody ground")
	}
	// Organizations failing is the point; a city with nobody left in it is not.
	if emptied > campaigns/10 {
		t.Fatalf("%d of %d cities ended with fewer than two organizations holding anything", emptied, campaigns)
	}
	if born == 0 {
		t.Fatal("no organization was ever created; the city can only shrink")
	}
}
