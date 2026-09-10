package core

import "testing"

// A war the player can only survive is not a war. This measures whether it can
// now be won, and what winning it takes.

func TestAWarCanNowBeWonAndCostsToWin(t *testing.T) {
	heavy(t)
	const runs, days = 200, 60
	measure := func(fight bool) (won, lost, bled, hurt int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(seed)
			w.MigrateLivingWorld()
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Cash, w.Player.Health = 60000, 100
			own(w, "laundry", "garage")
			w.Player.Respect = OrganizationStanding
			w.OrganizationDay()
			holder := &w.Factions[0]
			if c := w.Conflict(holder.ID, w.PlayerOrganizationID()); c != nil {
				c.Hostility, c.State = 80, "war"
			}
			mine := len(w.FamilyHoldings(w.PlayerOrganizationID()))
			for day := 0; day < days && w.Player.Alive; day++ {
				if fight {
					for _, id := range w.FamilyHoldings(holder.ID) {
						w.Player.Location = id
						if w.MoveOnReadiness(id) == "" {
							_ = w.MoveOn(id)
						}
						break
					}
				}
				w.FactionTurn()
				w.FamilyDay()
				w.OrganizationDay()
				w.Minute += 1440
			}
			// Winning is holding more than you started with. Whether the
			// rival shrank is not the measure: the city takes ground off them
			// on its own, and a campaign that ends early sees less of that.
			if len(w.FamilyHoldings(w.PlayerOrganizationID())) > mine {
				won++
			}
			if len(w.FamilyHoldings(w.PlayerOrganizationID())) < 2 {
				lost++
			}
			// A move on premises does not kill the player outright — the
			// readiness gate keeps them off it below forty health — so what it
			// costs is measured in the people who went with them.
			if w.Player.Health < 100 {
				hurt++
			}
			if len(w.OwnPeople()) == 0 && len(w.Player.Crew) == 0 {
				bled++
			}
		}
		return
	}

	passiveWon, passiveLost, _, _ := measure(false)
	fightWon, fightLost, _, fightHurt := measure(true)

	if fightWon <= passiveWon {
		t.Fatalf("fighting ended with more ground in %d campaigns against %d sitting still", fightWon, passiveWon)
	}
	if fightHurt == 0 {
		t.Fatal("nobody who fought a war was ever hurt in one")
	}
	t.Logf("%d campaigns of %d days at war: sitting still ends with more ground in %d and less in %d; fighting ends with more in %d, less in %d, and hurts you in %d",
		runs, days, passiveWon, passiveLost, fightWon, fightLost, fightHurt)
}
