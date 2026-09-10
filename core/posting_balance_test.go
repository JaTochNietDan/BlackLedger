package core

import "testing"

// A man on the door is one of the player's people doing nothing but standing
// somewhere. This measures whether a season of that is worth what it costs:
// the ground it keeps against the men it uses up.

func TestSomebodyOnTheDoorKeepsTheGround(t *testing.T) {
	t.Parallel()
	const runs, days = 200, 60
	measure := func(post bool) (condition, standing, buried int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w, _ := testator(t)
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Cash = 40000
			// Enough people that both arms of the measurement have somebody to
			// put on a door, and somebody who wants the premises.
			for _, n := range w.Civilians() {
				if IsOfficial(n.ID) || w.isRoleHolder(n) || len(w.OwnPeople()) >= 2 {
					continue
				}
				w.Player.Location = n.Location
				w.SignOn(n.ID)
			}
			enemy := &w.Factions[0]
			enemy.Power = 90
			if c := w.Conflict(enemy.ID, w.PlayerOrganizationID()); c != nil {
				c.Hostility, c.State = 80, "war"
			}
			started := len(w.OwnPeople())
			if post {
				w.Post("laundry")
				w.Post("garage")
			}
			for day := 0; day < days; day++ {
				w.FactionTurn()
				w.FamilyDay()
				w.OrganizationDay()
				w.PeopleDay()
				w.Minute += 1440
			}
			for _, id := range []string{"laundry", "garage"} {
				if w.Own(id) {
					condition += w.Properties[id].Condition
				}
			}
			standing += len(w.FamilyHoldings(w.PlayerOrganizationID()))
			buried += started - len(w.OwnPeople())
		}
		return
	}

	bareCondition, bareHoldings, bareBuried := measure(false)
	postedCondition, postedHoldings, postedBuried := measure(true)

	if postedCondition <= bareCondition {
		t.Fatalf("premises with somebody on the door ended a season at %d condition across %d campaigns against %d with nobody", postedCondition, runs, bareCondition)
	}
	// The men are not where the cost shows up over a season, and the
	// measurement says so rather than claiming otherwise. A man on the door
	// dies in the raids that get through — proved in
	// TestTheManOnTheDoorIsTheOneWhoPaysForIt — but he is also the reason
	// fewer of them get through, and the second effect is the larger one. What
	// posting actually costs is the man himself: somebody who is standing on a
	// door is not available for anything else.
	if postedBuried > bareBuried {
		t.Fatalf("standing on a door cost %d men across %d campaigns against %d for standing nowhere, which is not what the raid arithmetic should produce", postedBuried, runs, bareBuried)
	}
	t.Logf("%d campaigns of %d days at war: nobody on the door ends at %d condition, %d holdings and %d men lost; somebody on it ends at %d condition, %d holdings and %d men lost",
		runs, days, bareCondition, bareHoldings, bareBuried, postedCondition, postedHoldings, postedBuried)
}
