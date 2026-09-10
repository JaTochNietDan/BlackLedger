package core

import "testing"

// "It may be worth adding a couple more casino games."
//
// Craps is the one the period actually had, and it is a different thing to
// offer than the wheel: the wheel is one decision and then nothing, a hand of
// cards is a decision every card, and this is one decision followed by a
// sequence of rolls you cannot affect but have to sit through.

func shooter(t *testing.T) *World {
	t.Helper()
	w := gambler(t)
	w.Event, w.District = nil, 9
	w.Player.Location = "club"
	return w
}

// The edge is the real one: 1.414% on the pass line, over a sample big enough
// for the claim. At 200,000 decisions the standard error is 0.22%, which is
// wide enough that a perfectly fair table reads anywhere from 1.0% to 1.9% —
// and a run that came out at 1.75% sent me looking for a fault in the dice that
// was not there. Two million puts the error at 0.07%.
func TestThePassLineKeepsItsRealEdge(t *testing.T) {
	t.Parallel()
	const rounds = 2000000
	w := shooter(t)
	w.Player.Cash = 1000000000
	staked, returned := 0, 0
	for i := 0; i < rounds; i++ {
		before := w.Player.Cash
		if err := w.PlayDice("club", "pass", 10); err != nil {
			t.Fatal(err)
		}
		for w.Dice != nil && !w.Dice.Done {
			if err := w.RollAgain(); err != nil {
				t.Fatal(err)
			}
		}
		staked += 10
		returned += w.Player.Cash - before + 10
	}
	edge := float64(staked-returned) / float64(staked) * 100
	t.Logf("over %d pass-line decisions: staked $%d, returned $%d, edge %.3f%%", rounds, staked, returned, edge)
	// Four standard errors either side of the truth.
	if edge < 1.13 || edge > 1.70 {
		t.Fatalf("the pass line kept %.3f%% and should keep about 1.414%%", edge)
	}
}

func TestSevenOnTheComeOutWinsAndTwoLoses(t *testing.T) {
	t.Parallel()
	w := shooter(t)
	wins, losses, points := 0, 0, 0
	for i := 0; i < 4000; i++ {
		if err := w.PlayDice("club", "pass", 5); err != nil {
			t.Fatal(err)
		}
		d := w.Dice
		switch {
		case d.Point > 0:
			points++
			if d.Done {
				t.Fatalf("a point of %d settled on the come-out", d.Point)
			}
			for !w.Dice.Done {
				if err := w.RollAgain(); err != nil {
					t.Fatal(err)
				}
			}
		case d.Total == 7 || d.Total == 11:
			wins++
			if !d.Won || !d.Done {
				t.Fatalf("a %d on the come-out did not pay", d.Total)
			}
		case d.Total == 2 || d.Total == 3 || d.Total == 12:
			losses++
			if d.Won || !d.Done {
				t.Fatalf("a %d on the come-out did not lose", d.Total)
			}
		default:
			t.Fatalf("a total of %d was neither a decision nor a point", d.Total)
		}
	}
	t.Logf("of 4000 come-outs: %d naturals, %d craps, %d points set", wins, losses, points)
	if wins == 0 || losses == 0 || points == 0 {
		t.Fatal("the come-out never did one of the three things it can do")
	}
}

func TestThePointHasToBeMadeBeforeTheSeven(t *testing.T) {
	t.Parallel()
	w := shooter(t)
	made, sevened := 0, 0
	for i := 0; i < 2000; i++ {
		if err := w.PlayDice("club", "pass", 5); err != nil {
			t.Fatal(err)
		}
		if w.Dice.Point == 0 {
			continue
		}
		point := w.Dice.Point
		for !w.Dice.Done {
			if err := w.RollAgain(); err != nil {
				t.Fatal(err)
			}
		}
		d := w.Dice
		if d.Won {
			made++
			if d.Total != point {
				t.Fatalf("the point was %d and it paid on a %d", point, d.Total)
			}
		} else {
			sevened++
			if d.Total != 7 {
				t.Fatalf("the point was %d and it lost on a %d", point, d.Total)
			}
		}
	}
	t.Logf("of the points set: %d made, %d sevened out", made, sevened)
	if made == 0 || sevened == 0 {
		t.Fatal("a point only ever went one way")
	}
}

func TestTheDontIsTheOtherSideOfTheSameGame(t *testing.T) {
	t.Parallel()
	const rounds = 2000000
	w := shooter(t)
	w.Player.Cash = 1000000000
	staked, returned := 0, 0
	for i := 0; i < rounds; i++ {
		before := w.Player.Cash
		if err := w.PlayDice("club", "dont", 10); err != nil {
			t.Fatal(err)
		}
		for w.Dice != nil && !w.Dice.Done {
			if err := w.RollAgain(); err != nil {
				t.Fatal(err)
			}
		}
		staked += 10
		returned += w.Player.Cash - before + 10
	}
	edge := float64(staked-returned) / float64(staked) * 100
	t.Logf("over %d don't-pass decisions: edge %.3f%%", rounds, edge)
	if edge < 1.08 || edge > 1.65 {
		t.Fatalf("don't pass kept %.3f%% and should keep about 1.364%%", edge)
	}
}

func TestTheRoomOffersTheDice(t *testing.T) {
	t.Parallel()
	w := shooter(t)
	a := actionByID(w.Actions("club"), "dice")
	if a == nil {
		t.Fatal("a casino has no dice")
	}
	if a.Disabled {
		t.Fatalf("the dice are refused: %s", a.Reason)
	}
	if err := w.PlayDice("club", "pass", 5); err != nil {
		t.Fatal(err)
	}
	if w.Dice.Point > 0 {
		if actionByID(w.Actions("club"), "roll") == nil {
			t.Fatal("a point was set and there is no way to roll again")
		}
	}
	d := w.DiceDescription()
	if d["playing"] != true && d["settled"] != true {
		t.Fatalf("the table says nothing about the game just played: %+v", d)
	}
}
