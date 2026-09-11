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
// Getting it to work took six passes and every one of them was a fault in the
// policy rather than a fact about the game. It starved, then it died in a day and a half, and then it turned out
// that both of its rules were reading an action id that carries a person's name
// after a colon: the list of things that get you killed never matched
// `strike:person-8`, and "taken fewest times" counted every person in the room
// as a separate thing to try. Reading the id's root instead took it from dying
// in every run at 2.2 days to surviving eight in ten past a month, and from
// nineteen kinds of action to thirty-six — more on its own than the eight
// policies with plans reach between them.
//
// This pins the breadth so it cannot quietly fall. Raising it is how a policy
// that reaches further gets recognised; a drop is a policy that stopped
// reaching, which is how a whole area of the game goes unpriced without
// anybody noticing.
const kindsTheHarnessPlays = 77

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
