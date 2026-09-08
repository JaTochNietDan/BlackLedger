package core

import "testing"

// The city's temperature has to be something a campaign feels: a violent one
// should end up living in a harder city than a quiet one, without anybody
// having said so.

func TestAViolentCityBecomesAHarderCity(t *testing.T) {
	const runs, days = 150, 120
	measure := func(loud bool) (peak, ended, crackdowns int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(seed)
			w.MigrateLivingWorld()
			w.WorldRNG = seed * 2654435761
			if loud {
				// A city at war with itself, which is what a violent campaign
				// leaves behind whoever started it.
				for i := range w.Conflicts {
					w.Conflicts[i].Hostility, w.Conflicts[i].State = 85, "war"
				}
			} else {
				for i := range w.Conflicts {
					w.Conflicts[i].Hostility, w.Conflicts[i].State = 5, "cold"
				}
			}
			worst, seen := 0, false
			for day := 0; day < days; day++ {
				w.Minute += 1440
				w.FactionTurn()
				w.FamilyDay()
				w.PeopleDay()
				w.GrudgeDay()
				w.SettleGrudges()
				w.ScrutinyDay()
				worst = max(worst, w.Scrutiny())
				if w.UnderCrackdown() {
					seen = true
				}
			}
			peak += worst
			ended += w.Scrutiny()
			if seen {
				crackdowns++
			}
		}
		return
	}

	quietPeak, quietEnd, quietCrack := measure(false)
	loudPeak, loudEnd, loudCrack := measure(true)
	if loudPeak <= quietPeak {
		t.Fatalf("a city at war peaked at %d against %d for one at peace", loudPeak/runs, quietPeak/runs)
	}
	if loudCrack <= quietCrack {
		t.Fatalf("a city at war saw a crackdown in %d campaigns against %d", loudCrack, quietCrack)
	}
	if quietCrack*2 > runs {
		t.Fatalf("a city at peace saw a crackdown in %d of %d campaigns, so it is not about what anybody did", quietCrack, runs)
	}
	t.Logf("%d cities over %d days: at peace they peak at %d and end at %d, with a crackdown in %d; at war they peak at %d and end at %d, with a crackdown in %d",
		runs, days, quietPeak/runs, quietEnd/runs, quietCrack, loudPeak/runs, loudEnd/runs, loudCrack)
}
