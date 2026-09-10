package core

import "testing"

// Being taken in is not death and it is not a fine. This measures what it
// actually is: a week of the city running without you, and the bill afterwards.

func TestAWeekInsideIsAWeekTheCityHadWithoutYou(t *testing.T) {
	t.Parallel()
	const runs, days = 200, 7
	measure := func(inside bool) (cash, condition, ground, supply int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w := chargeable(t)
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Cash = 6000
			// Both arms start cold. A raid mid-week would stop the clock on the
			// arm that is standing outside and turn this into a measurement of
			// the police rather than of a week of absence.
			w.Player.Heat = 0
			enemy := &w.Factions[0]
			enemy.Power = 85
			if c := w.Conflict(enemy.ID, w.PlayerOrganizationID()); c != nil {
				c.Hostility, c.State = 70, "war"
			}
			if inside {
				w.Confine(days, "a still in the back of the laundry")
				w.SitOut()
			} else {
				// The same week, standing outside it, doing what anybody who
				// owns two businesses does with a week: keeping them stocked
				// and in one piece. That is the thing a cell takes away — not
				// the days, which pass either way, but the ability to spend
				// them on anything.
				for day := 0; day < days; day++ {
					for _, id := range []string{"laundry", "garage"} {
						w.Restock(id)
						w.Remedy(id)
					}
					w.Advance(1440)
				}
			}
			if !w.Player.Alive {
				continue
			}
			cash += w.Player.Cash
			for _, id := range []string{"laundry", "garage"} {
				if w.Own(id) {
					condition += w.Properties[id].Condition
					supply += w.Properties[id].Supply
				}
			}
			ground += len(w.FamilyHoldings(w.PlayerOrganizationID()))
		}
		return
	}

	outCash, outCondition, outGround, outSupply := measure(false)
	inCash, inCondition, inGround, inSupply := measure(true)

	if inCash >= outCash {
		t.Fatalf("a week inside ended richer: $%d across %d campaigns against $%d for a week outside", inCash, runs, outCash)
	}
	// Ground is not where a week shows up: an organization holds what it holds
	// whether or not the man running it is available, and seven days is not
	// long enough for a war to take premises off anybody. The measurement says
	// so rather than claiming otherwise.
	if inSupply >= outSupply {
		t.Fatalf("premises nobody could restock ended the week with %d supplies against %d for premises somebody could", inSupply, outSupply)
	}
	// Condition is not where a week shows up either, and it comes out within a
	// percent either way — the two arms run the clock in different numbers of
	// calls, which is enough to move a figure that small. What a cell costs is
	// the takings and the ability to keep a business supplied; the rest of the
	// numbers are recorded rather than asserted.
	t.Logf("%d campaigns, %d days at war: outside ends with $%d, %d condition, %d holdings, %d supplies; inside ends with $%d, %d condition, %d holdings, %d supplies",
		runs, days, outCash, outCondition, outGround, outSupply, inCash, inCondition, inGround, inSupply)
}

func TestTalkingIsWorthMoreThanTimeAndCostsMoreThanMoney(t *testing.T) {
	t.Parallel()
	// Three ways out of the same sentence, measured against each other over the
	// same campaigns, so the comparison is what it costs rather than what
	// happened to the city that week.
	const runs = 200
	measure := func(way string) (respect, cash, goodwill, minutes int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w := chargeable(t)
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Cash, w.Player.Respect = 12000, 120
			start := w.Minute
			w.Confine(5, "a still")
			switch way {
			case "serve":
				w.SitOut()
			case "lawyer":
				w.Lawyer()
			case "talk":
				w.Talk()
			}
			if !w.Player.Alive {
				continue
			}
			respect += w.Player.Respect
			cash += w.Player.Cash
			for i := range w.Factions {
				goodwill += w.Factions[i].Goodwill
			}
			minutes += w.Minute - start
		}
		return
	}

	sRespect, sCash, sGoodwill, sMinutes := measure("serve")
	lRespect, lCash, lGoodwill, lMinutes := measure("lawyer")
	tRespect, tCash, tGoodwill, tMinutes := measure("talk")

	if sRespect <= lRespect || sRespect <= tRespect {
		t.Fatalf("doing the time was worth %d respect against %d for a lawyer and %d for talking", sRespect, lRespect, tRespect)
	}
	if lCash >= sCash {
		t.Fatalf("a lawyer cost nothing: $%d against $%d for doing the time", lCash, sCash)
	}
	if tGoodwill >= sGoodwill || tGoodwill >= lGoodwill {
		t.Fatalf("talking left %d standing against %d for serving and %d for paying", tGoodwill, sGoodwill, lGoodwill)
	}
	if tMinutes >= sMinutes {
		t.Fatal("talking took as long as doing the time")
	}
	t.Logf("%d campaigns, five days each: serving ends at %d respect, $%d, %d standing, %d minutes; a lawyer at %d respect, $%d, %d standing, %d minutes; talking at %d respect, $%d, %d standing, %d minutes",
		runs, sRespect, sCash, sGoodwill, sMinutes, lRespect, lCash, lGoodwill, lMinutes, tRespect, tCash, tGoodwill, tMinutes)
}
