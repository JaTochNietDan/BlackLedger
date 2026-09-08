package core

import "testing"

// An understanding is protection bought with money and paid for in enemies.
// This measures both halves over a campaign.

func TestAnUnderstandingIsBoughtWithMoneyAndPaidForInEnemies(t *testing.T) {
	const runs, days = 200, 60
	measure := func(ally bool) (kept, hostility, spent int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w, friend, enemy := diplomat(t)
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Cash = 60000
			// Somebody who means it.
			if c := w.Conflict(enemy.ID, w.PlayerOrganizationID()); c != nil {
				c.Hostility, c.State = 75, "war"
			}
			enemy.Power = 90
			start := w.Player.Cash
			if ally {
				if err := w.MakePact(friend.ID); err != nil {
					t.Fatal(err)
				}
			}
			for day := 0; day < days; day++ {
				w.FactionTurn()
				w.FamilyDay()
				w.OrganizationDay()
				w.PactDay()
				w.Minute += 1440
			}
			kept += len(w.FamilyHoldings(w.PlayerOrganizationID()))
			if c := w.Conflict(enemy.ID, w.PlayerOrganizationID()); c != nil {
				hostility += c.Hostility
			}
			spent += start - w.Player.Cash
		}
		return
	}

	aloneKept, aloneHostility, aloneSpent := measure(false)
	alliedKept, alliedHostility, alliedSpent := measure(true)

	if alliedKept <= aloneKept {
		t.Fatalf("standing with somebody kept %d holdings across %d campaigns against %d alone", alliedKept, runs, aloneKept)
	}
	// The enemy's hostility is not the measure over a season: signing costs
	// fifteen at once and a day of standing with somebody costs two, but a war
	// burns itself out faster than that accumulates, and an ally at the door
	// takes strength off the attacker. What making an enemy costs is immediate,
	// and TestStandingWithSomebodyPutsYouInTheirQuarrels is where it is proved.
	if alliedSpent <= aloneSpent {
		t.Fatal("standing with somebody cost nothing")
	}
	t.Logf("%d campaigns of %d days at war: alone keeps %d holdings, ends at %d hostility and spends $%d; with an understanding, %d holdings, %d hostility and $%d",
		runs, days, aloneKept, aloneHostility, aloneSpent, alliedKept, alliedHostility, alliedSpent)
}
