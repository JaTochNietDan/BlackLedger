package core

import "testing"

// "Everything anybody in this city does is done by the same rules you use.
// There is no separate arithmetic for them."
//
// The RNG is deliberately two streams — one reserved for events the player is
// not party to, so families quarrelling off-screen cannot shift the odds of a
// decision being made. That is separate randomness for a stated reason, not
// separate arithmetic, and it is fine.
//
// The arithmetic is a different question. A family moving on a rival's premises
// runs contestAt; the player moving on the same premises runs SabotageBy. They
// are two functions. This measures whether they agree about what an attack on a
// building costs — because if they do not, the sentence above is not true.

func damageBy(t *testing.T, runs int, attack func(w *World, place string) int) (int, int, float64) {
	t.Helper()
	low, high, total, counted := 1000, 0, 0, 0
	for i := 0; i < runs; i++ {
		w := New(uint32(1000 + i*7))
		w.District = 2
		w.Player.Cash, w.Player.Respect, w.Player.Health = 20000, 200, 100
		w.Player.Crew = []Crew{{"leo", "Leo Carver", 100}}
		w.Player.Location = "club"
		w.Properties["club"].Condition = 100
		before := w.Properties["club"].Condition
		if !attackRan(w, attack, "club") {
			continue
		}
		hit := before - w.Properties["club"].Condition
		if hit <= 0 {
			continue
		}
		counted++
		total += hit
		if hit < low {
			low = hit
		}
		if hit > high {
			high = hit
		}
	}
	if counted == 0 {
		t.Fatal("no attack ever landed, so this measured nothing")
	}
	return low, high, float64(total) / float64(counted)
}

func attackRan(w *World, attack func(w *World, place string) int, place string) bool {
	return attack(w, place) == 0
}

func TestAFamilyAndThePlayerAgreeWhatAnAttackCosts(t *testing.T) {
	t.Parallel()
	player := func(w *World, place string) int {
		if err := w.SabotageBy(place, w.OwnHands()); err != nil {
			return 1
		}
		return 0
	}
	family := func(w *World, place string) int {
		attacker, defender := w.faction("russo"), w.faction("bellandi")
		if attacker == nil || defender == nil || w.Properties[place].Owner != defender.ID {
			return 1
		}
		attacker.Power, defender.Power = 90, 58
		w.contestAt(attacker, defender, place)
		return 0
	}
	pLow, pHigh, pMean := damageBy(t, 400, player)
	fLow, fHigh, fMean := damageBy(t, 400, family)
	t.Logf("the player: %d to %d, mean %.1f", pLow, pHigh, pMean)
	t.Logf("a family:   %d to %d, mean %.1f", fLow, fHigh, fMean)

	// Not identical code, so not identical draws. But if the two disagree by
	// more than a fifth about what an attack on a building costs, the city has
	// separate arithmetic for the people in it and the rule says it does not.
	apart(t, "condition off a building", pMean, fMean, 0.2)
}

// apart reports when the two disagree by more than a share of the larger. The
// first version of this had the comparison inside an else, so it was skipped
// whenever the family's figure was the larger of the two — a check that could
// only fail in one direction, which is half a check.
func apart(t *testing.T, what string, player, family, tolerance float64) {
	t.Helper()
	gap := player - family
	if gap < 0 {
		gap = -gap
	}
	larger := bigger(player, family)
	if larger == 0 {
		t.Fatalf("%s: both are zero, so this measured nothing", what)
	}
	share := gap / larger
	t.Logf("%s: %.1f by the player, %.1f by a family, %.0f%% apart", what, player, family, 100*share)
	if share > tolerance {
		t.Errorf("%s differs by %.0f%% between the player and a family, and the rules say there is no separate arithmetic",
			what, 100*share)
	}
}

// The same question about money. The player's attack takes damage*20 off the
// family that owns the premises; a family's raid takes damage*15 off the family
// it raided. Those are different constants, so this measures what each one
// actually costs.
func TestAFamilyAndThePlayerAgreeWhatAnAttackTakes(t *testing.T) {
	t.Parallel()
	taken := func(runs int, attack func(w *World) (before, after int, ok bool)) float64 {
		total, counted := 0, 0
		for i := 0; i < runs; i++ {
			w := New(uint32(2000 + i*7))
			w.District = 2
			w.Player.Cash, w.Player.Respect, w.Player.Health = 20000, 200, 100
			w.Player.Crew = []Crew{{"leo", "Leo Carver", 100}}
			w.Player.Location = "club"
			w.Properties["club"].Condition = 100
			before, after, ok := attack(w)
			if !ok || before-after <= 0 {
				continue
			}
			counted++
			total += before - after
		}
		if counted == 0 {
			t.Fatal("nothing was ever taken, so this measured nothing")
		}
		return float64(total) / float64(counted)
	}
	byPlayer := taken(400, func(w *World) (int, int, bool) {
		f := w.faction("bellandi")
		before := f.Cash
		if err := w.SabotageBy("club", w.OwnHands()); err != nil {
			return 0, 0, false
		}
		return before, f.Cash, true
	})
	byFamily := taken(400, func(w *World) (int, int, bool) {
		attacker, defender := w.faction("russo"), w.faction("bellandi")
		attacker.Power, defender.Power = 90, 58
		before := defender.Cash
		w.contestAt(attacker, defender, "club")
		return before, defender.Cash, true
	})
	apart(t, "money off a family", byPlayer, byFamily, 0.2)
}

func bigger(a, b float64) float64 {
	if b > a {
		return b
	}
	return a
}

// The comparison above must work in both directions. The first version had it
// inside an else, so whichever way round the two figures came it could only
// ever fail one way — half a check, and it happened to be the half that was
// green.
func TestTheComparisonWorksInBothDirections(t *testing.T) {
	t.Parallel()
	if bigger(10, 100) != 100 {
		t.Error("a larger family figure is not taken as the larger")
	}
	if bigger(100, 10) != 100 {
		t.Error("a larger player figure is not taken as the larger")
	}
	if bigger(7, 7) != 7 {
		t.Error("equal figures do not agree")
	}
}
