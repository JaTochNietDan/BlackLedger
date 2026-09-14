package core

import "testing"

func warlord(t *testing.T) (*World, *Faction, string) {
	t.Helper()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect = OrganizationStanding
	w.Incorporate()
	w.OrganizationDay()
	holder := &w.Factions[0]
	held := w.FamilyHoldings(holder.ID)
	if len(held) == 0 {
		t.Fatal("the fixture's rival held nothing")
	}
	target := held[0]
	c := w.Conflict(holder.ID, w.PlayerOrganizationID())
	c.Hostility, c.State = 80, "war"
	w.Player.Location, w.Player.Health = target, 100
	return w, holder, target
}

func TestAManDoesNotTakeGround(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry")
	w.Player.Respect = 5
	w.OrganizationDay()
	w.Player.Location = "club"
	if _, ok := w.MoveTarget("club"); ok {
		t.Fatal("a proprietor was offered somebody else's ground")
	}
	if w.MoveOnReadiness("club") == "" {
		t.Fatal("a proprietor could move on a family")
	}
}

func TestYouOnlyMoveOnPeopleYouAreAtOddsWith(t *testing.T) {
	t.Parallel()
	w, holder, target := warlord(t)
	if w.MoveOnReadiness(target) != "" {
		t.Fatal("could not move on a rival at war:", w.MoveOnReadiness(target))
	}
	c := w.Conflict(holder.ID, w.PlayerOrganizationID())
	c.Hostility, c.State = 10, "cold"
	if w.MoveOnReadiness(target) == "" {
		t.Fatal("moved on somebody the player was on civil terms with")
	}
	c.State = "feud"
	if w.MoveOnReadiness(target) != "" {
		t.Fatal("a feud was not enough:", w.MoveOnReadiness(target))
	}
	// Nor on your own premises.
	if _, ok := w.MoveTarget("laundry"); ok {
		t.Fatal("the player was offered their own laundry")
	}
}

func TestYouHaveToBeThereAndAbleToStand(t *testing.T) {
	t.Parallel()
	w, _, target := warlord(t)
	w.Player.Location = "room"
	if w.MoveOnReadiness(target) == "" {
		t.Fatal("moved on premises across the city")
	}
	w.Player.Location = target
	w.Player.Health = 20
	if w.MoveOnReadiness(target) == "" {
		t.Fatal("moved on somebody while barely standing")
	}
	w.Player.Health = 100
	w.Player.Respect = 0
	for _, id := range w.FamilyHoldings(w.PlayerOrganizationID()) {
		w.Properties[id].Condition = 1
	}
	w.OrganizationDay()
	if w.PlayerStrength() >= MoveStanding {
		t.Skip("this fixture is still strong enough")
	}
	if w.MoveOnReadiness(target) == "" {
		t.Fatal("nobody would follow them and they went anyway")
	}
}

func TestAMoveIsTheSameRaidTheCityRuns(t *testing.T) {
	t.Parallel()
	const runs = 400
	took, repelled, hurt, died := 0, 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, holder, target := warlord(t)
		w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
		w.Properties[target].Condition = 20
		holder.Power = 20
		w.OrganizationDay()
		before := len(w.FamilyHoldings(w.PlayerOrganizationID()))
		condition := w.Properties[target].Condition
		if err := w.MoveOn(target); err != nil {
			t.Fatal(err)
		}
		switch {
		case len(w.FamilyHoldings(w.PlayerOrganizationID())) > before:
			took++
		case w.Properties[target].Condition < condition:
			// Struck but not taken.
		default:
			repelled++
		}
		if w.Player.Health < 100 {
			hurt++
		}
		if !w.Player.Alive {
			died++
		}
	}
	if took == 0 {
		t.Fatal("no move in 400 ever took ground")
	}
	if repelled == 0 {
		t.Fatal("no move in 400 was ever driven off")
	}
	t.Logf("%d moves on a weakened holding: %d took it, %d were driven off, the player was hurt %d times and killed %d", runs, took, repelled, hurt, died)
}

func TestTakingGroundIsWorthSomethingAndCostsStanding(t *testing.T) {
	t.Parallel()
	found := false
	for seed := uint32(1); seed <= 400 && !found; seed++ {
		w, holder, target := warlord(t)
		w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
		w.Properties[target].Condition = 10
		holder.Power = 15
		w.OrganizationDay()
		respect, goodwill := w.Player.Respect, holder.Goodwill
		hostility := w.Conflict(holder.ID, w.PlayerOrganizationID()).Hostility
		before := len(w.FamilyHoldings(w.PlayerOrganizationID()))
		w.MoveOn(target)
		if len(w.FamilyHoldings(w.PlayerOrganizationID())) > before {
			found = true
			if w.Player.Respect <= respect {
				t.Fatal("taking somebody's ground was worth no standing")
			}
			if holder.Goodwill >= goodwill {
				t.Fatal("taking somebody's ground did not offend them")
			}
			if w.Conflict(holder.ID, w.PlayerOrganizationID()).Hostility <= hostility {
				t.Fatal("taking somebody's ground did not harden the quarrel")
			}
			if w.Player.Heat == 0 {
				t.Fatal("a daylight move drew no attention")
			}
		}
	}
	if !found {
		t.Fatal("no move in 400 took ground")
	}
}

func TestBeingDrivenOffCostsSomebodyWhoWent(t *testing.T) {
	t.Parallel()
	found := false
	for seed := uint32(1); seed <= 500 && !found; seed++ {
		w, holder, target := warlord(t)
		w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
		holder.Power = 100
		w.OrganizationDay()
		// Somebody to lose.
		var member *NPC
		for _, n := range w.Civilians() {
			if !IsOfficial(n.ID) && !w.isRoleHolder(n) {
				member = n
				break
			}
		}
		member.Faction, member.Rank, member.Trust = w.PlayerOrganizationID(), RankSoldier, 50
		health := w.Player.Health
		w.MoveOn(target)
		if member.Dead || w.Player.Health < health {
			found = true
		}
	}
	if !found {
		t.Fatal("no move in 500 was driven off at any cost")
	}
}
