package core

import "testing"

// The tests that skip themselves.
//
// A door test of mine looked for one of the player's own people to post, found
// none, and skipped — reporting a pass in every summary while proving nothing,
// and passing with the rule it guarded removed. That is a whole family of tests
// here: "no lieutenant in this world", "this city has no forecourt", "no
// breakaway for this seed", "everybody in this city is known". Each looks for
// something in a generated world and quietly stands down when it is not there.
//
// They are not wrong to be written that way — establishing a lieutenant by hand
// is not the thing under test — but it means the day the city stops making
// lieutenants, a dozen tests start proving nothing and nothing says so. This
// says so. It asserts the preconditions themselves, in the same worlds those
// tests use, so a change to the city that empties them is one loud failure
// rather than twelve silent passes.

func TestTheCityStillMakesWhatItsTestsLookFor(t *testing.T) {
	t.Parallel()
	// The seeds those suites actually use, so this is about the same worlds.
	for _, seed := range []uint32{4, 27, 61, 91, 92, 404} {
		w := New(seed)
		lieutenants, known, players := 0, 0, 0
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if n.Dead {
				continue
			}
			players++
			// Below leader, or this cannot fail: a family always has a head,
			// and a head satisfies "rank at or above lieutenant" — the first
			// version of this guard counted them and passed with every
			// lieutenant in the city turned into a soldier.
			if n.Rank >= RankLieutenant && n.Rank < RankLeader && n.Faction != "" {
				lieutenants++
			}
			if w.Known(n) {
				known++
			}
		}
		if players == 0 {
			t.Fatalf("seed %d: nobody lives in this city", seed)
		}
		if lieutenants == 0 {
			t.Errorf("seed %d: no family has a lieutenant, and five tests stand down when that is true", seed)
		}
		// "Everybody in this city is known" is the other half of two of them:
		// there has to be somebody the player has never dealt with.
		if known >= players {
			t.Errorf("seed %d: every one of the %d people here is already known to the player", seed, players)
		}
	}
	// And the addresses those suites ask for by name.
	for _, id := range []string{"dealer", "garage", "scrapyard", "pumps", "filling"} {
		if _, ok := PlaceByID(id); !ok {
			t.Errorf("this city has no %s, and the tests that want one skip themselves", id)
		}
	}
}
