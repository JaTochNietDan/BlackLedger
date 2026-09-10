package core

import "testing"

// The city's whole temperature is derived from what the paper carried. This
// measures whether a man who decides that is worth $38 a day, and what using
// him against somebody else does to the city everybody lives in.

func TestAManAtTheHeraldIsWorthWhatTheCityDoesNotRead(t *testing.T) {
	heavy(t)
	const runs, days = 200, 40
	measure := func(paid bool) (scrutiny, heat, crackdowns, cash int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w := proprietor(t)
			own(w, "laundry", "garage")
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Cash, w.Player.Respect = 30000, 60
			w.Player.Location = HeraldPlace
			w.ensureOfficials()
			if paid {
				if err := w.Retain("editor"); err != nil {
					t.Fatal(err)
				}
			}
			// A loud city: something is in the paper most days.
			for i := range w.Factions {
				w.Factions[i].Power = 85
			}
			w.Antagonize("bellandi", "russo", 90)
			crackdown := false
			for day := 0; day < days; day++ {
				w.Player.Location = HeraldPlace
				if paid && w.SpikeReadiness() == "" {
					w.Spike()
				}
				w.Advance(1440)
				if !w.Player.Alive {
					break
				}
				if w.UnderCrackdown() {
					crackdown = true
				}
			}
			if !w.Player.Alive {
				continue
			}
			scrutiny += w.Scrutiny()
			heat += w.Player.Heat
			cash += w.Player.Cash
			if crackdown {
				crackdowns++
			}
		}
		return
	}

	quietScrutiny, quietHeat, quietCrackdowns, quietCash := measure(false)
	paidScrutiny, paidHeat, paidCrackdowns, paidCash := measure(true)

	if paidScrutiny >= quietScrutiny {
		t.Fatalf("%d campaigns of %d days: paying the Herald ended at %d scrutiny against %d for reading it like everybody else", runs, days, paidScrutiny, quietScrutiny)
	}
	if paidCash >= quietCash {
		t.Fatal("a man on a daily retainer cost nothing")
	}
	t.Logf("%d campaigns of %d days in a loud city: unpaid ends at %d scrutiny, %d attention, %d crackdowns and $%d; a man at the Herald ends at %d scrutiny, %d attention, %d crackdowns and $%d",
		runs, days, quietScrutiny, quietHeat, quietCrackdowns, quietCash, paidScrutiny, paidHeat, paidCrackdowns, paidCash)
}

func TestRunningSomethingAboutSomebodyElseHeatsTheWholeCity(t *testing.T) {
	heavy(t)
	// The point of the mechanic: it works, and it works on a city you also live
	// in. Measured over the same campaigns with and without.
	const runs, days = 200, 30
	measure := func(smear bool) (theirs, scrutiny, goodwill int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w := proprietor(t)
			own(w, "laundry")
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Cash, w.Player.Respect = 60000, 60
			w.Player.Location = HeraldPlace
			w.ensureOfficials()
			if err := w.Retain("editor"); err != nil {
				t.Fatal(err)
			}
			for day := 0; day < days; day++ {
				w.Player.Location = HeraldPlace
				if smear && w.SmearReadiness("bellandi") == "" {
					w.Smear("bellandi")
				}
				w.Advance(1440)
				if !w.Player.Alive {
					break
				}
			}
			if !w.Player.Alive {
				continue
			}
			// A family can be destroyed inside a month now that money presses on
			// a quarrel, and a city that no longer has one has nothing to say
			// about what it thinks of the player. This used to reach straight
			// into the lookup and crashed the first time a war finished one
			// off.
			them := w.faction("bellandi")
			if them == nil {
				continue
			}
			for _, id := range w.FamilyHoldings("bellandi") {
				theirs += w.Custom(id)
			}
			scrutiny += w.Scrutiny()
			goodwill += them.Goodwill
		}
		return
	}

	quietTheirs, quietScrutiny, quietGoodwill := measure(false)
	loudTheirs, loudScrutiny, loudGoodwill := measure(true)

	if loudTheirs >= quietTheirs {
		t.Fatalf("a month of stories left their trade at %d against %d untouched", loudTheirs, quietTheirs)
	}
	if loudScrutiny <= quietScrutiny {
		t.Fatalf("a month of crime reporting left the city at %d scrutiny against %d", loudScrutiny, quietScrutiny)
	}
	if loudGoodwill >= quietGoodwill {
		t.Fatal("nobody ever worked out who was paying for it")
	}
	t.Logf("%d campaigns of %d days: leaving them alone ends with their trade at %d, the city at %d scrutiny and them at %d standing; running stories ends with their trade at %d, the city at %d scrutiny and them at %d standing",
		runs, days, quietTheirs, quietScrutiny, quietGoodwill, loudTheirs, loudScrutiny, loudGoodwill)
}
