package core

import "testing"

// A car came from nowhere. The player pressed a button at the motor works, the
// money left their pocket, and a car existed — nobody sold it and nobody was
// paid for it. The city has a forecourt now, and the money goes to whoever
// holds it, which is the whole idea: the places in this city should be where
// the things in it actually come from.

func TestOnlyAForecourtSellsCars(t *testing.T) {
	t.Parallel()
	lots := 0
	for _, l := range Locations {
		if CarSource(l.ID) {
			lots++
			if l.Kind != "dealer" {
				t.Errorf("%s sells cars and is a %q", l.ID, l.Kind)
			}
		}
	}
	if lots == 0 {
		t.Fatal("nowhere in this city sells a car")
	}
	// A motor works repairs cars. It does not sell them.
	for _, l := range Locations {
		if l.Kind == "garage" && CarSource(l.ID) {
			t.Errorf("%s is a garage and sells cars", l.ID)
		}
	}
}

// The money has to arrive somewhere. Buying a car from a family's forecourt
// pays that family; buying from your own pays you back the margin.
func TestBuyingACarPaysWhoeverHoldsTheForecourt(t *testing.T) {
	t.Parallel()
	lot := ""
	for _, l := range Locations {
		if l.Kind == "dealer" {
			lot = l.ID
		}
	}
	if lot == "" {
		t.Fatal("no forecourt")
	}

	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect = 40000, 200
	w.Player.Location = lot
	w.Properties[lot].Owner = "bellandi"
	house := w.faction("bellandi")
	if house == nil {
		t.Fatal("no bellandi")
	}
	before := house.Cash
	if err := w.BuyVehicle(); err != nil {
		t.Fatalf("buying a car: %v", err)
	}
	if house.Cash <= before {
		t.Errorf("a car was sold off Bellandi's forecourt and they took $%d", house.Cash-before)
	}
	t.Logf("a family holding the forecourt took $%d on one sale", house.Cash-before)
}

// And the player's own forecourt gives them the margin back rather than paying
// a rival for their own car.
func TestYourOwnForecourtSellsYouACarCheaply(t *testing.T) {
	t.Parallel()
	lot := ""
	for _, l := range Locations {
		if l.Kind == "dealer" {
			lot = l.ID
		}
	}
	spend := func(owner string) int {
		w := New(53)
		w.District = 2
		w.Player.Cash, w.Player.Respect = 40000, 200
		w.Player.Location = lot
		w.Properties[lot].Owner = owner
		before := w.Player.Cash
		if err := w.BuyVehicle(); err != nil {
			t.Fatalf("buying a car from %s: %v", owner, err)
		}
		return before - w.Player.Cash
	}
	theirs, mine := spend("bellandi"), spend("player:1")
	if mine >= theirs {
		t.Errorf("a car costs %d off your own forecourt and %d off a rival's", mine, theirs)
	}
	t.Logf("a car costs $%d off a rival's forecourt and $%d off your own", theirs, mine)
}

// One question was doing two jobs. `CarSource` answered both "where is a car
// sold" and "where is a car worked on", which nobody noticed while the answer
// to both was the motor works. Moving sales to a forecourt moved servicing with
// them, and a garage could no longer touch a car.
func TestAGarageWorksOnCarsAndAForecourtSellsThem(t *testing.T) {
	t.Parallel()
	garages, lots := 0, 0
	for _, l := range Locations {
		works, sells := CarWorkshop(l.ID), CarSource(l.ID)
		switch l.Kind {
		case "garage":
			garages++
			if !works {
				t.Errorf("%s is a garage and nobody works on cars there", l.ID)
			}
			if sells {
				t.Errorf("%s is a garage and sells cars", l.ID)
			}
		case "dealer":
			lots++
			if !sells {
				t.Errorf("%s is a forecourt and sells nothing", l.ID)
			}
		default:
			if works || sells {
				t.Errorf("%s is a %q and works on or sells cars", l.ID, l.Kind)
			}
		}
	}
	if garages < 2 {
		t.Errorf("the city has %d garages, so 'any garage' is still not provable", garages)
	}
	if lots == 0 {
		t.Fatal("no forecourt")
	}
}

// And now that there are two garages, the rule I recorded last slice as latent
// and unprovable can be caught failing.
func TestEitherGarageWorksOnYourCar(t *testing.T) {
	t.Parallel()
	for _, l := range Locations {
		if l.Kind != "garage" {
			continue
		}
		w := New(53)
		w.District = 2
		w.Player.Cash, w.Player.Respect = 40000, 200
		w.Player.Location = "dealer"
		if err := w.BuyVehicle(); err != nil {
			t.Fatal(err)
		}
		w.Player.CarWear = 40
		w.Player.Location = l.ID
		if reason := w.ServiceReadiness(l.ID); reason != "" {
			t.Errorf("%s is a garage and will not work on the car: %q", l.ID, reason)
		}
	}
}
