package core

import "testing"

// The whole decision is a trade: your own hands are better at it and put the
// risk on you. This measures both halves across enough robberies that neither
// is an anecdote.

func TestGoingYourselfEarnsMoreAndCostsMore(t *testing.T) {
	t.Parallel()
	const runs = 500
	type result struct{ took, respect, heat, hurt, dead, lostCrew int }
	measure := func(own bool) result {
		r := result{}
		for seed := uint32(1); seed <= runs; seed++ {
			w := withCrew(t, 100)
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Cash = 0
			hand := w.OwnHands()
			if !own {
				hand, _ = w.CrewHands()
			}
			if err := w.RobBy("club", hand); err != nil {
				t.Fatal(err)
			}
			r.took += w.Player.Cash
			r.respect += w.Player.Respect - 40
			r.heat += w.Player.Heat
			if w.Player.Health < 100 {
				r.hurt++
			}
			if !w.Player.Alive {
				r.dead++
			}
			if len(w.Player.Crew) == 0 {
				r.lostCrew++
			}
		}
		return r
	}

	own, sent := measure(true), measure(false)
	if own.took <= sent.took {
		t.Fatalf("going yourself took $%d against $%d sending somebody", own.took, sent.took)
	}
	if own.respect <= sent.respect {
		t.Fatalf("going yourself earned %d standing against %d", own.respect, sent.respect)
	}
	if own.heat <= sent.heat {
		t.Fatalf("going yourself drew %d attention against %d", own.heat, sent.heat)
	}
	if own.hurt <= sent.hurt {
		t.Fatalf("going yourself hurt the player %d times against %d", own.hurt, sent.hurt)
	}
	if sent.lostCrew <= own.lostCrew {
		t.Fatalf("sending somebody cost the crewman %d times against %d going yourself", sent.lostCrew, own.lostCrew)
	}
	t.Logf("%d robberies at the same casino, going yourself: $%d taken, %d standing, %d attention, hurt %d times, killed %d times",
		runs, own.took, own.respect, own.heat, own.hurt, own.dead)
	t.Logf("%d robberies sending Leo: $%d taken, %d standing, %d attention, the player hurt %d times, Leo lost %d times",
		runs, sent.took, sent.respect, sent.heat, sent.hurt, sent.lostCrew)
}
