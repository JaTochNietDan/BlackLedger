package core

import "testing"

// The second career now has a top. This measures what standing at it takes and
// what it is worth against the career of buying premises one at a time.

func TestTakingItIsTheFastestWayUpAndTheMostLikelyToEndYou(t *testing.T) {
	const runs = 500
	took, thrownOut, killed := 0, 0, 0
	holdings := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, f, leader := lieutenant(t)
		w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
		// The danger is not in trying it well; it is in trying it at the
		// bottom of what the rules allow, which is where a player who has just
		// come up from something else usually is.
		w.Player.Health = 50 + int(seed%40)
		if err := w.TakeOver(); err != nil {
			t.Fatal(err)
		}
		switch {
		case leader.Dead:
			took++
			holdings += len(w.FamilyHoldings(w.PlayerOrganizationID()))
		case !w.Player.Alive:
			killed++
		default:
			thrownOut++
		}
		_ = f
	}
	if took == 0 || thrownOut == 0 {
		t.Fatalf("%d attempts took it and %d were thrown out", took, thrownOut)
	}
	if killed == 0 {
		t.Fatal("nobody who tried it ever died in the room")
	}
	t.Logf("%d attempts by a lieutenant on a weakened organization: %d took it and came away with %d premises in total, %d were thrown out, %d did not leave the room",
		runs, took, holdings, thrownOut, killed)
}

// And what it takes to be worth following: the same attempt by somebody with a
// name and by somebody without.
func TestAManWithANameTakesItAndAManWithoutDoesNot(t *testing.T) {
	const runs = 400
	rate := func(respect, weapon int) int {
		took := 0
		for seed := uint32(1); seed <= runs; seed++ {
			w, _, leader := lieutenant(t)
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Respect, w.Player.Weapon = respect, weapon
			if w.TakeoverReadiness() != "" {
				continue
			}
			w.TakeOver()
			if leader.Dead {
				took++
			}
		}
		return took
	}
	bare, armed := rate(TakeoverStanding, 0), rate(100, 3)
	if armed <= bare {
		t.Fatalf("a man with a name and a gun took it %d times against %d for a man with neither", armed, bare)
	}
	t.Logf("%d attempts each: at the bare minimum it works %d times; with a hundred respect and a Thompson, %d", runs, bare, armed)
}
