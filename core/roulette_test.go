package core

import (
	"fmt"
	"strings"
	"testing"
)

// The wheel has to be a real wheel. If the payouts are shaded to make the house
// win, the game is lying about what it is offering; the edge should be the
// zero and nothing else, and that is a property you can measure.

func TestTheWheelIsThirtySevenPocketsAndTheZeroBelongsToNobody(t *testing.T) {
	t.Parallel()
	if Pockets != 37 {
		t.Fatalf("the wheel has %d pockets", Pockets)
	}
	reds, blacks := 0, 0
	for p := 0; p <= 36; p++ {
		if Red(p) {
			reds++
		}
		if Black(p) {
			blacks++
		}
		if Red(p) && Black(p) {
			t.Errorf("%d is both colours", p)
		}
	}
	if reds != 18 || blacks != 18 {
		t.Errorf("%d red and %d black pockets", reds, blacks)
	}
	if Red(Zero) || Black(Zero) {
		t.Error("the zero has a colour, so it is not the house's pocket")
	}
	// And no outside bet takes it.
	for _, bet := range RouletteBets() {
		if bet.ID == "number:0" {
			continue
		}
		if bet.Wins(Zero) {
			t.Errorf("%s wins on the nought, so the house has no edge at all", bet.Label)
		}
	}
}

// The payouts are the true ones. Over every pocket exactly once, a dollar on
// any bet returns thirty-six dollars for thirty-seven staked — the edge is one
// pocket in thirty-seven and it is the same for every bet on the cloth, which
// is what tells you nothing has been shaded.
func TestEveryBetOnTheClothHasTheSameEdgeAndItIsTheZero(t *testing.T) {
	t.Parallel()
	for _, bet := range RouletteBets() {
		staked, returned := 0, 0
		for pocket := 0; pocket < Pockets; pocket++ {
			staked += 1
			if bet.Wins(pocket) {
				returned += bet.Pays + 1
			}
		}
		if returned != 36 {
			t.Errorf("%s returns %d for %d staked over the whole wheel, so its odds are not true",
				bet.Label, returned, staked)
		}
	}
}

// And the money moves once, for the right amount. A dollar on a number that
// comes in is thirty-six dollars back, not thirty-five.
func TestAWinningNumberPaysThirtyFiveToOneAndReturnsTheStake(t *testing.T) {
	t.Parallel()
	for pocket := 0; pocket < Pockets; pocket++ {
		w := wheelRoom(t)
		before := w.Player.Cash
		// Drive the wheel onto a known pocket by walking its stream.
		w.RNG = seedForPocket(t, w, pocket)
		if err := w.PlayWheel("casino", fmt.Sprintf("number:%d", pocket), 50); err != nil {
			t.Fatalf("pocket %d: %v", pocket, err)
		}
		if w.Spin.Pocket != pocket {
			t.Fatalf("wanted pocket %d and the ball dropped in %d", pocket, w.Spin.Pocket)
		}
		if !w.Spin.Won {
			t.Fatalf("backing %d and the ball dropped in %d did not win", pocket, w.Spin.Pocket)
		}
		stake, _ := tableStake("small")
		if won := w.Player.Cash - before; won != stake.Amount*35 {
			t.Fatalf("pocket %d: $%d staked came back as $%d net, wanted $%d",
				pocket, stake.Amount, won, stake.Amount*35)
		}
	}
}

// A losing spin takes the stake once and no more.
func TestALosingSpinTakesTheStakeOnce(t *testing.T) {
	t.Parallel()
	w := wheelRoom(t)
	stake, _ := tableStake("small")
	before := w.Player.Cash
	// Back a number, and whichever pocket comes up it is wrong 36 times in 37.
	if err := w.PlayWheel("casino", "number:7", 50); err != nil {
		t.Fatal(err)
	}
	if w.Spin.Won {
		t.Skip("the seven came in, which is not the case this test is about")
	}
	if lost := before - w.Player.Cash; lost != stake.Amount {
		t.Errorf("a losing spin on a $%d stake took $%d", stake.Amount, lost)
	}
}

// The room cannot have a hand of cards and a wheel going at once.
func TestTheWheelWaitsForTheCardsToBeFinished(t *testing.T) {
	t.Parallel()
	w := wheelRoom(t)
	if err := w.Deal("casino", 50); err != nil {
		t.Fatalf("dealing: %v", err)
	}
	if w.Hand == nil || w.Hand.Done {
		t.Skip("the hand settled itself, so there is nothing on the table to clash with")
	}
	if err := w.PlayWheel("casino", "red", 50); err == nil {
		t.Error("the wheel was played with a hand still on the table")
	}
}

// wheelRoom is a player standing at a casino they do not own, with money.
func wheelRoom(t *testing.T) *World {
	t.Helper()
	w := New(23)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 20000, 200, 100
	w.Player.Dress = 1
	w.Player.Location = "casino"
	if w.Own("casino") {
		t.Fatal("the player owns the room, so they cannot play in it")
	}
	stake, _ := tableStake("small")
	if reason := w.TableReadiness("casino", stake); reason != "" {
		t.Fatalf("cannot sit down: %s", reason)
	}
	return w
}

// seedForPocket finds a seed whose first draw puts the ball in a given pocket.
// The wheel is driven by the player's own stream, so this is the only honest
// way to test a named pocket without reaching inside the spin itself.
func seedForPocket(t *testing.T, w *World, pocket int) uint32 {
	t.Helper()
	for seed := uint32(1); seed < 200000; seed++ {
		probe := *w
		probe.RNG = seed
		if got := int(probe.Random() * Pockets); got == pocket {
			return seed
		}
	}
	t.Fatalf("no seed in range puts the ball in %d", pocket)
	return 0
}

// The wiring, not the function. The wheel has to be reachable as a command the
// way the interface will actually send it, and the bet has to survive the trip:
// Target names the room and is what the action lookup searches, so a bet put
// there would send the game looking for a wheel in a place called "red".
func TestTheWheelCanBePlayedAsACommand(t *testing.T) {
	t.Parallel()
	w := wheelRoom(t)
	before := w.Player.Cash
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision,
		Kind: "wheel", Target: "casino", Choice: "number:17", Amount: 50})
	if err != nil {
		t.Fatalf("playing the wheel: %v", err)
	}
	if next.Spin == nil {
		t.Fatal("the command went through and nothing was ever spun")
	}
	if next.Spin.Bet != "number:17" {
		t.Errorf("backed number:17 and the game recorded %q", next.Spin.Bet)
	}
	if next.Spin.Place != "casino" {
		t.Errorf("played at the casino and the game recorded %q", next.Spin.Place)
	}
	stake, _ := tableStake("small")
	moved := before - next.Player.Cash
	if next.Spin.Won && moved > 0 {
		t.Errorf("a winning spin left the player $%d down", moved)
	}
	if !next.Spin.Won && moved != stake.Amount {
		t.Errorf("a losing $%d spin moved $%d", stake.Amount, moved)
	}
}

// And the button is there to be pressed, naming its price, declaring no cost.
func TestTheWheelIsOfferedInARoomWithTables(t *testing.T) {
	t.Parallel()
	w := wheelRoom(t)
	var wheel *Action
	for i, a := range w.Actions("casino") {
		if a.ID == "wheel" {
			wheel = &w.Actions("casino")[i]
		}
	}
	if wheel == nil {
		t.Fatal("a casino offers no wheel")
	}
	// The wheel used to be two buttons at two fixed prices, and this asked that
	// the button named its own price. A player names the figure now, so what is
	// asked instead is that the button names no price of its own and that the
	// room says what it will take.
	if wheel.Asks != 0 {
		t.Errorf("the wheel names a price of $%d when the player picks the figure", wheel.Asks)
	}
	if !strings.Contains(wheel.Detail, "$") || !strings.Contains(wheel.Detail, fmt.Sprint(w.TableLimit("casino"))) {
		t.Errorf("the wheel does not say what the house takes: %q", wheel.Detail)
	}
	if wheel.Cost != 0 {
		t.Errorf("the wheel declares a cost of %d, so the stake would go out twice", wheel.Cost)
	}
	// And not in a room without tables.
	for _, a := range w.Actions("bar") {
		if a.ID == "wheel" {
			t.Error("there is a roulette wheel in the bar")
		}
	}
}

// The view can only show what the core sends it. The hand was in the world the
// interface reads and the wheel was not, so a spin happened and nothing outside
// the ledger could ever say what the ball did.
func TestTheWorldTellsTheInterfaceWhatTheWheelDid(t *testing.T) {
	t.Parallel()
	w := wheelRoom(t)
	quiet, ok := w.Public()["wheel"].(map[string]any)
	if !ok {
		t.Fatal("the world does not describe the wheel at all")
	}
	if quiet["spun"] != false {
		t.Error("a wheel nobody has played reads as spun")
	}
	if err := w.PlayWheel("casino", "number:17", 50); err != nil {
		t.Fatal(err)
	}
	after, _ := w.Public()["wheel"].(map[string]any)
	if after["spun"] != true {
		t.Fatal("the wheel was played and the world says it was not")
	}
	for _, field := range []string{"pocket", "colour", "bet", "won", "stake", "pays", "place"} {
		if _, there := after[field]; !there {
			t.Errorf("the wheel is described without %s, so the table cannot be drawn", field)
		}
	}
	if after["bet"] != "Straight up on 17" {
		t.Errorf("backed 17 and the world describes the bet as %q", after["bet"])
	}
}
