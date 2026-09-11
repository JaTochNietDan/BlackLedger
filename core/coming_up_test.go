package core

import "testing"

// TestComingUpInAFamilyIsWalkable walks the game's longest progression with the
// game's own actions: answer to a family, do enough work for them to be made a
// lieutenant, then take the family. Every step of it was on the never-played
// list, and the one existing test of a takeover sets Service directly, so the
// earning of a promotion had never been exercised at all. What is asked here is
// only whether the path can be walked, not how fast.
func TestComingUpInAFamilyIsWalkable(t *testing.T) {
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 400
	mine := w.faction("bellandi")
	if mine == nil {
		t.Skip("no Bellandi Family in this city")
	}

	// Somebody they are at odds with, holding something that can be moved on.
	var enemy *Faction
	var target string
	for i := range w.Factions {
		other := &w.Factions[i]
		if other.ID == mine.ID || other.ID == w.PlayerOrganizationID() {
			continue
		}
		if w.Conflict(mine.ID, other.ID) == nil {
			continue
		}
		for _, id := range w.FamilyHoldings(other.ID) {
			if _, ok := w.SabotageTarget(id); ok {
				enemy, target = other, id
				break
			}
		}
		if enemy != nil {
			break
		}
	}
	if enemy == nil {
		t.Skip("nobody the Bellandi Family is at odds with holds anything")
	}

	// Answering to them is a conversation held in their own room.
	w.Player.Location = w.homeOf(mine.ID)
	mine.Goodwill = 100
	if reason := w.ServeReadiness(mine.ID); reason != "" {
		t.Fatalf("cannot go to work for them: %s", reason)
	}
	if err := w.Serve(mine.ID); err != nil {
		t.Fatal(err)
	}
	if w.ServiceRank() == RankLieutenant {
		t.Fatal("a new man is already a lieutenant")
	}

	// Somebody to go with. Recruiting is a scene choice rather than an action,
	// and it is the same driver the command appends.
	if driver := w.Holder("driver"); driver != nil && len(w.Player.Crew) == 0 {
		w.Player.Crew = append(w.Player.Crew, Crew{ID: driver.ID, Name: driver.Name, Loyalty: 65})
	}

	// Work done against their enemy is work done for them. The attack can be
	// turned away, so this is an attempt count and not a hit count.
	w.Antagonize(mine.ID, enemy.ID, 100)
	attempts := 0
	for ; attempts < 40 && w.Player.Service < PromotionWork*2; attempts++ {
		w.Player.Health, w.Player.Heat = 100, 0
		w.Tasks = nil
		attempt := attempts
		if len(w.Player.Crew) == 0 || w.Player.Crew[0].Loyalty < 40 {
			t.Fatalf("no crew to go with after %d attempts", attempt)
		}
		if reason := w.SabotageReadiness(target); reason != "" {
			t.Fatalf("cannot move on them after %d attempts: %s", attempt, reason)
		}
		if err := w.SabotageBy(target, w.OwnHands()); err != nil {
			t.Fatal(err)
		}
		w.Event = nil
	}
	if w.ServiceRank() != RankLieutenant {
		t.Fatalf("work for them never made a lieutenant: service %d, title %q", w.Player.Service, w.ServiceTitle())
	}
	t.Logf("lieutenant after %d moves against %s", attempts, enemy.Name)

	// And then the room itself.
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 400
	leader := w.Leader(mine.ID)
	if leader == nil {
		t.Skip("nobody is in the chair")
	}
	w.Player.Location = leader.Location
	mine.Power, mine.Peak = 10, 100
	if reason := w.TakeoverReadiness(); reason != "" {
		t.Fatalf("a lieutenant of theirs cannot take them: %s", reason)
	}
	if err := w.TakeOver(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Serves != "" {
		t.Fatalf("still answering to %q after taking the family", w.Player.Serves)
	}
}
