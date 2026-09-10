package core

import "testing"

// "You can also place multiple bets in roulette, on different numbers,
// combinations etc, like the real game by putting down chips on each one you
// want to bet on."
//
// One bet a spin was never how the game works. What a real table takes is a
// cloth covered in chips, every one of them settled against the same pocket.

func clothPlayer(t *testing.T) *World {
	t.Helper()
	w := gambler(t)
	w.Event, w.District = nil, 9
	w.Player.Location = "club"
	return w
}

func TestTheClothTakesAsManyChipsAsYouPutOnIt(t *testing.T) {
	t.Parallel()
	w := clothPlayer(t)
	cash := w.Player.Cash
	chips := []Chip{{Bet: "red", Amount: 20}, {Bet: "number:17", Amount: 5}, {Bet: "even", Amount: 10}}
	if err := w.SpinChips("club", chips); err != nil {
		t.Fatal(err)
	}
	if w.Spin == nil || len(w.Spin.Chips) != 3 {
		t.Fatalf("the cloth kept %+v", w.Spin)
	}
	// Every one of them is settled against the same pocket, and the arithmetic
	// is the sum of what each was worth.
	pocket := w.Spin.Pocket
	want := 0
	for _, c := range chips {
		bet, _ := RouletteBetByID(c.Bet)
		if bet.Wins(pocket) {
			want += c.Amount * (bet.Pays + 1)
		}
	}
	if got := w.Player.Cash - (cash - 35); got != want {
		t.Fatalf("pocket %d paid %d and should have paid %d", pocket, got, want)
	}
}

func TestEveryChipIsWeighedAgainstTheHouseLimit(t *testing.T) {
	t.Parallel()
	w := clothPlayer(t)
	limit := w.TableLimit("club")
	// The limit is per bet, the way a table's is: two chips at the limit are
	// fine, one chip over it is not.
	if reason := w.ChipsReadiness("club", []Chip{{Bet: "red", Amount: limit}, {Bet: "black", Amount: limit}}); reason != "" {
		t.Fatalf("two chips at the limit were refused: %s", reason)
	}
	if w.ChipsReadiness("club", []Chip{{Bet: "red", Amount: limit + 1}}) == "" {
		t.Fatal("a chip over the house limit was taken")
	}
	// And you cannot put down more than you have, however it is spread.
	w.Player.Cash = 30
	if w.ChipsReadiness("club", []Chip{{Bet: "red", Amount: 20}, {Bet: "black", Amount: 20}}) == "" {
		t.Fatal("the table took $40 off somebody with $30")
	}
	if w.ChipsReadiness("club", []Chip{}) == "" {
		t.Fatal("an empty cloth was spun")
	}
}

func TestTheWholeClothComesOffYourCashAtOnce(t *testing.T) {
	t.Parallel()
	w := clothPlayer(t)
	cash := w.Player.Cash
	// Red and black together: whatever the pocket, one of them comes back and
	// the other does not, unless it is the nought and both are gone. Either way
	// the money leaves the pocket before the ball drops.
	if err := w.SpinChips("club", []Chip{{Bet: "red", Amount: 25}, {Bet: "black", Amount: 25}}); err != nil {
		t.Fatal(err)
	}
	pocket := w.Spin.Pocket
	switch {
	case pocket == Zero:
		if w.Player.Cash != cash-50 {
			t.Fatalf("the nought took %d of 50", cash-w.Player.Cash)
		}
	default:
		if w.Player.Cash != cash {
			t.Fatalf("covering both colours on a %d left the player %d up", pocket, w.Player.Cash-cash)
		}
	}
}

func TestOneChipIsStillTheOldGame(t *testing.T) {
	t.Parallel()
	w := clothPlayer(t)
	cash := w.Player.Cash
	if err := w.PlayWheel("club", "red", 40); err != nil {
		t.Fatal(err)
	}
	if w.Spin == nil || w.Spin.Down != 40 || w.Spin.Bet != "red" {
		t.Fatalf("a single bet no longer reads as one: %+v", w.Spin)
	}
	won := Red(w.Spin.Pocket)
	if won && w.Player.Cash != cash+40 {
		t.Fatalf("red came in and paid %d", w.Player.Cash-cash)
	}
	if !won && w.Player.Cash != cash-40 {
		t.Fatalf("red lost and cost %d", cash-w.Player.Cash)
	}
}
