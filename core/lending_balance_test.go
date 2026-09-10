package core

import "testing"

// Lending is easy and collecting is where it becomes a decision. This measures
// the three answers against each other over a season, and whether money on the
// street beats money in a pocket at all.

func TestMoneyOnTheStreetBeatsMoneyInAPocket(t *testing.T) {
	t.Parallel()
	const runs, days = 200, 45
	measure := func(lend bool) (cash, respect, heat, bad int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w, _ := lender(t)
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Cash = 4000
			for day := 0; day < days; day++ {
				if lend {
					// Walked, not teleported. Every conversation costs the
					// journey to it and the half hour in it, out of the same
					// day the other arm spends doing nothing — otherwise this
					// measures free money rather than a day's work.
					for _, n := range w.Civilians() {
						if IsOfficial(n.ID) || w.isRoleHolder(n) || !w.Player.Alive || w.Held() {
							continue
						}
						travel := TravelMinutes(w.Player.Location, n.Location)
						w.Player.Location = n.Location
						w.Advance(travel)
						if w.LendReadiness(n.ID) == "" {
							w.Lend(n.ID)
							w.Advance(LendMinutes)
						}
						// Whatever comes back overdue gets collected, which is
						// the trade as it is actually practised.
						if w.LeanReadiness(n.ID) == "" {
							w.Lean(n.ID)
							w.Advance(LeanMinutes)
						}
					}
				}
				w.Advance(1440)
				if !w.Player.Alive {
					break
				}
			}
			if !w.Player.Alive {
				continue
			}
			cash += w.Player.Cash + w.OutOnLoan()
			respect += w.Player.Respect
			heat += w.Player.Heat
			bad += len(w.BadDebt())
		}
		return
	}

	idleCash, idleRespect, idleHeat, _ := measure(false)
	lentCash, lentRespect, lentHeat, lentBad := measure(true)

	if lentCash <= idleCash {
		t.Fatalf("%d campaigns of %d days: lending ended with $%d against $%d for leaving it in a pocket", runs, days, lentCash, idleCash)
	}
	if lentHeat <= idleHeat {
		t.Fatalf("a season of collections drew %d attention against %d for doing nothing", lentHeat, idleHeat)
	}
	t.Logf("%d campaigns of %d days: a pocket ends with $%d, %d respect and %d attention; the street ends with $%d, %d respect, %d attention and %d debts standing",
		runs, days, idleCash, idleRespect, idleHeat, lentCash, lentRespect, lentHeat, lentBad)
}

func TestTheThreeAnswersToAManWhoCannotPay(t *testing.T) {
	t.Parallel()
	const runs = 300
	measure := func(answer string) (cash, respect, heat, trust, sore int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w, n := lender(t)
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			// Somebody worth lending to: a lieutenant of a family that is
			// doing well, who can service a sum worth the conversation.
			n.Rank, n.Faction, n.Ambition = RankLieutenant, "bellandi", 10
			if f := w.faction("bellandi"); f != nil {
				f.Cash = 200000
			}
			w.Player.Cash, w.Player.Respect = 40000, 80
			if err := w.Lend(n.ID); err != nil {
				t.Fatal(err)
			}
			start := w.Player.Cash
			// His position collapses while he is carrying it — the family he
			// stood behind stops paying, or he was never anybody to begin with.
			// That is a real default rather than a man who chose not to pay,
			// and it is the case the three answers exist for.
			n.Rank, n.Faction, n.Skill = 0, "", 5
			w.Minute = w.LoanTo(n.ID).Due
			w.LoanDay()
			if w.LoanTo(n.ID) == nil {
				t.Fatal("a man with nothing paid in full")
			}
			switch answer {
			case "lean":
				w.Lean(n.ID)
			case "extend":
				w.Extend(n.ID)
			case "forgive":
				w.Forgive(n.ID)
			}
			cash += w.Player.Cash - start
			respect += w.Player.Respect
			heat += w.Player.Heat
			trust += n.Trust
			sore += n.Sore
		}
		return
	}

	leanCash, leanRespect, leanHeat, leanTrust, leanSore := measure("lean")
	extCash, extRespect, extHeat, extTrust, extSore := measure("extend")
	forCash, forRespect, forHeat, forTrust, forSore := measure("forgive")

	if leanCash <= extCash || leanCash <= forCash {
		t.Fatalf("leaning recovered $%d against $%d for extending and $%d for writing it off", leanCash, extCash, forCash)
	}
	if leanRespect <= forRespect {
		t.Fatalf("collecting was worth %d respect against %d for writing it off", leanRespect, forRespect)
	}
	if leanHeat <= extHeat || leanHeat <= forHeat {
		t.Fatal("a collection in the street drew no more attention than a conversation")
	}
	if forTrust <= leanTrust || forTrust <= extTrust {
		t.Fatalf("writing it off bought %d trust against %d for leaning and %d for extending", forTrust, leanTrust, extTrust)
	}
	if leanSore <= extSore || leanSore <= forSore {
		t.Fatalf("a man put against a wall held %d against the player, against %d for one given another week", leanSore, extSore)
	}
	t.Logf("%d defaults each: leaning recovers $%d, ends at %d respect, %d attention, %d trust, %d held against you; extending $%d, %d respect, %d attention, %d trust, %d held; writing off $%d, %d respect, %d attention, %d trust, %d held",
		runs, leanCash, leanRespect, leanHeat, leanTrust, leanSore,
		extCash, extRespect, extHeat, extTrust, extSore,
		forCash, forRespect, forHeat, forTrust, forSore)
}
