package core

import "testing"

// Being filed with the others is the point at which everything the city does to
// an organization starts being done to the player. This measures whether that
// is a real change in a campaign rather than a line in a log.

func TestBeingSomebodyIsWorseThanBeingNobody(t *testing.T) {
	heavy(t)
	const runs, days = 150, 90
	measure := func(incorporated bool) (raids, seizures, lost int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(seed)
			w.MigrateLivingWorld()
			w.WorldRNG = seed * 2654435761
			w.Player.Cash, w.Player.Health = 40000, 100
			own(w, "laundry", "garage")
			if incorporated {
				w.Player.Respect = OrganizationStanding
				w.OrganizationDay()
				// Somebody with a reason, so the machinery has something to run.
				for i := range w.Factions {
					if w.Factions[i].ID != w.PlayerOrganizationID() {
						if c := w.Conflict(w.Factions[i].ID, w.PlayerOrganizationID()); c != nil {
							c.Hostility, c.State = 80, "war"
						}
					}
				}
			} else {
				w.Player.Respect = OrganizationStanding - 1
			}
			before := w.Player.Cash
			for day := 0; day < days; day++ {
				w.FactionTurn()
				w.FamilyDay()
				w.Minute += 1440
			}
			if w.Player.Cash < before {
				raids++
			}
			held := len(w.FamilyHoldings(w.PlayerOrganizationID()))
			if held < 2 {
				seizures++
			}
			lost += 2 - held
		}
		return
	}

	quietRaids, quietSeizures, _ := measure(false)
	namedRaids, namedSeizures, namedLost := measure(true)
	if namedRaids <= quietRaids {
		t.Fatalf("a named organization was raided in %d of %d campaigns against %d for a proprietor", namedRaids, runs, quietRaids)
	}
	if namedSeizures <= quietSeizures {
		t.Fatalf("a named organization lost premises in %d campaigns against %d", namedSeizures, quietSeizures)
	}
	t.Logf("%d campaigns of %d days: a proprietor loses money in %d and premises in %d; an organization at war loses money in %d, premises in %d, and %d holdings in total",
		runs, days, quietRaids, quietSeizures, namedRaids, namedSeizures, namedLost)
}
