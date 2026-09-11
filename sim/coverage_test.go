package sim

import "testing"

// How much of the game the harness actually plays.
//
// Every balance figure this project prints comes out of these campaigns, so
// what the campaigns never do is never priced. Measured across eight hundred
// runs of the eight policies that had plans: they took 78 of the game's 116
// kinds of action, and ninety-four kinds were never touched by any of them.
// That does not mean those are unreachable — it means nobody wrote a policy
// that wanted them, and the difference matters, because a window at a
// pawnbroker, a boat at a pier, a night at a room you host and a meeting
// between two families had all been built and balanced without one campaign
// ever taking any of them.
//
// The magpie exists for that. It has no plan: it works the docks until it can
// afford to be curious, then takes whichever card in the room it has taken
// fewest times, and moves on when the room has nothing new. It plays badly on
// purpose and its cash column should be read as what happens to somebody who
// does everything once.
//
// It is short-lived and that is stated rather than hidden: a median of a bit
// over two game days, and it dies in every run. Tuning bought some of that
// back — waiting until it is unhurt before taking an attempt on anybody raised
// the median from 1.4 days and the ids it reaches from 28 to 45 — and then
// stopped buying much, so the tuning stopped. What it is worth is breadth in
// the early game, which is nine kinds of action nothing else here ever took.
//
// This pins the breadth so it cannot quietly fall. Raising it is how a policy
// that reaches further gets recognised; a drop is a policy that stopped
// reaching, which is how a whole area of the game goes unpriced without
// anybody noticing.
const kindsTheHarnessPlays = 29

func TestTheHarnessPlaysEnoughOfTheGame(t *testing.T) {
	t.Parallel()
	every := map[string]bool{}
	for _, strategy := range []string{
		"worker", "investor", "defiant", "reckless", "thief",
		"smuggler", "racketeer", "publican", "magpie",
	} {
		for seed := uint32(1); seed <= 2; seed++ {
			r := Run(seed*2654435761, strategy, "fixture", 300, false)
			for id := range r.Actions {
				for i := 0; i < len(id); i++ {
					if id[i] == ':' {
						id = id[:i]
						break
					}
				}
				every[id] = true
			}
		}
	}
	if len(every) < kindsTheHarnessPlays {
		kinds := []string{}
		for k := range every {
			kinds = append(kinds, k)
		}
		t.Fatalf("the harness plays %d kinds of action and this expects %d: %v",
			len(every), kindsTheHarnessPlays, kinds)
	}
	t.Logf("the harness plays %d kinds of action", len(every))
}
