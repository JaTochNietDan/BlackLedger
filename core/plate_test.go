package core

import (
	"strings"
	"testing"
)

// "You should probably also be able to outfit your car with protection like
// armor etc at a garage which will help you survive attacks when traversing out
// in the streets."
//
// The top of the range was already called an armoured Packard and nothing in
// the rules had ever read that word. Now plating is a thing you buy, a thing
// that slows the car down, and a thing that is worth something on the one
// stretch of the city where nothing else is.

func plater(t *testing.T) *World {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	// A car that is actually on the road: in order, and with something in it.
	w.Player.Car, w.Player.CarWear = 1, 100
	w.Player.Fuel, w.Player.Fuelled = FuelFull, max(1, w.Minute)
	w.Player.Location = w.theGarage()
	if w.Player.Location == "" {
		t.Skip("this city has no garage")
	}
	return w
}

func TestPlatingIsFittedAtAGarageAndNowhereElse(t *testing.T) {
	t.Parallel()
	w := plater(t)
	garage := w.Player.Location
	a := actionByID(w.Actions(garage), "plate")
	if a == nil {
		t.Fatal("a garage does not fit plate")
	}
	if a.Disabled {
		t.Fatalf("refused: %s", a.Reason)
	}
	before := w.Plating()
	if err := w.FitPlate(garage); err != nil {
		t.Fatal(err)
	}
	if w.Plating() != before+1 {
		t.Fatalf("plating went from %d to %d", before, w.Plating())
	}
	// Not at the laundry, and not without a car.
	w.Player.Location = "laundry"
	if w.PlateReadiness("laundry") == "" {
		t.Fatal("a laundry fitted plate")
	}
	w.Player.Location, w.Player.Car = garage, 0
	if reason := w.PlateReadiness(garage); !strings.Contains(strings.ToLower(reason), "car") {
		t.Fatalf("plating was offered to somebody on foot: %q", reason)
	}
}

func TestPlateIsWorthSomethingOnTheStreetAndNowhereElse(t *testing.T) {
	t.Parallel()
	w := plater(t)
	garage := w.Player.Location
	w.Player.Location = "transit"
	bare := w.AmbushOddsHere()
	w.Player.Location = garage
	for w.Plating() < PlateStages {
		if err := w.FitPlate(garage); err != nil {
			t.Fatal(err)
		}
	}
	w.Player.Location = "transit"
	plated := w.AmbushOddsHere()
	if plated >= bare {
		t.Fatalf("plate changed nothing on the street: %.3f against %.3f", plated, bare)
	}
	// Indoors it is a car parked outside and worth nothing at all.
	w.Player.Location = "bar"
	clearRoom(w, "bar")
	inside := w.Cover()
	w.Player.Plate = 0
	if inside != w.Cover() {
		t.Fatal("plating on a car outside protected somebody standing in a bar")
	}
}

func TestPlateSurvivesTheJourneyAndNotTheCar(t *testing.T) {
	t.Parallel()
	w := plater(t)
	garage := w.Player.Location
	if err := w.FitPlate(garage); err != nil {
		t.Fatal(err)
	}
	plated := w.Plating()
	if plated == 0 {
		t.Fatal("nothing was fitted")
	}
	// A different car is a different car. Plate is fitted to the one you have
	// and stays with it; every car off the lot is bare, the Packard included.
	w.Player.Location, w.Player.Cash = w.theForecourt(), 20000
	if w.Player.Location == "" {
		t.Skip("this city has no forecourt")
	}
	if err := w.BuyVehicle(3); err != nil {
		t.Skip("nothing to buy: " + err.Error())
	}
	if w.Plating() != 0 {
		t.Fatalf("plate came off the old car and onto the new one: %d", w.Plating())
	}
}

func TestPlateCostsYouSpeed(t *testing.T) {
	t.Parallel()
	w := plater(t)
	garage := w.Player.Location
	quick := w.Pace()
	if err := w.FitPlate(garage); err != nil {
		t.Fatal(err)
	}
	if w.Pace() <= quick {
		t.Fatalf("a plated car was no slower: %.3f against %.3f", w.Pace(), quick)
	}
}

// The inverse of the rule this used to guard. The top of the range was called
// an armoured Packard and the rules had started reading that word, so the only
// way to get cover was to buy the car that came with it — buying a car and
// arming a car were one purchase and there was nothing to decide.
func TestNoCarComesOffTheLotWithPlateOnIt(t *testing.T) {
	t.Parallel()
	w := plater(t)
	w.Player.Location, w.Player.Cash = w.theForecourt(), 20000
	w.Player.Car, w.Player.CarWear, w.Player.Plate = 0, 0, 0
	if w.Player.Location == "" {
		t.Skip("this city has no forecourt")
	}
	for tier := 1; tier < len(vehicles); tier++ {
		if err := w.BuyVehicle(tier); err != nil {
			t.Fatalf("buying tier %d: %v", tier, err)
		}
		if w.Plating() != 0 {
			t.Fatalf("%s came off the lot carrying %d of plate", VehicleByTier(tier).Label, w.Plating())
		}
	}
	// And a save written while the Packard was armoured keeps what it had,
	// rather than losing two stages overnight to a change in the counting.
	old := plater(t)
	old.Player.Car, old.Player.Plate = 3, 0
	old.MigrateLivingWorld()
	if old.Plating() != PlateStages {
		t.Fatalf("an older save's armoured Packard came back with %d of plate", old.Plating())
	}
}
