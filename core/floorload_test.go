package core

import (
	"strings"
	"testing"
)

// What a false floor is worth depends on there being something to lose.
//
// `scrap` — weighing the car in at a yard — was on the list of actions the
// harness has never played, and walking it found the thing the yard was only
// half doing. Both the card and the log said "anything that was under the floor
// of it went with the car", and both of them took nothing: stock in this game
// is one pool with a hiding budget, so a car going to the crane left every
// crate where it was and only made it visible. A false floor with nothing at
// stake is a discount on attention, not a risk anybody weighs.
func TestWhatRidesUnderTheFloorGoesWithTheCar(t *testing.T) {
	t.Parallel()
	w := New(23)
	w.Event, w.District, w.Player.Cash = nil, 9, 400000

	yard := ""
	for _, l := range Locations {
		if ScrapYard(l.ID) {
			yard = l.ID
		}
	}
	if yard == "" {
		t.Skip("this city has no scrapyard")
	}

	// A car with a floor under it, and more stock than fits in a coat.
	w.Player.Car, w.Player.CarWear = 2, 0
	if VehicleByTier(w.Player.Car).Compartment == 0 {
		t.Fatal("the car under test has no false floor, so this measures nothing")
	}
	good := w.Goods[0].ID
	load := PocketLoad + 12
	w.Player.Stock = map[string]int{good: load}
	if w.Carrying() != load {
		t.Fatalf("the player is holding %d of an intended %d", w.Carrying(), load)
	}
	if floor := w.InTheFloor(); floor != 12 {
		t.Fatalf("%d units ride under the floor, not the %d that did not fit in a coat", floor, 12)
	}

	// The card says so before the crane picks it up, because a crate you did
	// not know you were losing is not a decision.
	w.Player.Location = yard
	var card Action
	for _, a := range w.Actions(yard) {
		if a.ID == "scrap" {
			card = a
		}
	}
	if card.ID == "" || card.Disabled {
		t.Fatalf("the yard will not take the car: %q", card.Reason)
	}
	if !strings.Contains(card.Detail, "12 crates") {
		t.Fatalf("the card does not say what goes with the car: %q", card.Detail)
	}

	scrapped, err := Execute(w, Command{Kind: "scrap", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if left := scrapped.Holding(good); left != PocketLoad {
		t.Fatalf("%d units left after weighing the car in, and %d fit in a coat", left, PocketLoad)
	}
	if scrapped.Player.Car != 0 {
		t.Fatal("the car is still there")
	}
}

// And the same thing when somebody else decides the car is going.
func TestWhatRidesUnderTheFloorGoesWhenTheCarIsBurned(t *testing.T) {
	t.Parallel()
	w := New(23)
	w.Event = nil
	w.Player.Car, w.Player.CarWear = 3, 0
	good := w.Goods[0].ID
	w.Player.Stock = map[string]int{good: PocketLoad + 5}
	if !w.LoseCar("Somebody put a match to it.") {
		t.Fatal("there was no car to lose")
	}
	if left := w.Holding(good); left != PocketLoad {
		t.Fatalf("%d units survived the fire, and %d fit in a coat", left, PocketLoad)
	}
	// And a man carrying nothing but what is in his pockets loses nothing, so
	// the rule is about the floor rather than about the car.
	w.Player.Car, w.Player.Stock = 3, map[string]int{good: 4}
	w.LoseCar("And the next one.")
	if left := w.Holding(good); left != 4 {
		t.Fatalf("%d of 4 units in a coat pocket went with the car", left)
	}
}
