package core

import "testing"

// "Sells gas that cars need." A car cost money every day and never needed
// anything: it drove for ever on nothing, so a filling station would have been
// a shop selling something no vehicle in this city consumed.

func TestACarBurnsWhatItIsDriven(t *testing.T) {
	w := New(31)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 4000, 40, 100
	w.Player.Car, w.Player.CarWear = 2, 100
	w.SettleFuel()
	full := w.Fuel()
	if full <= 0 {
		t.Fatal("a car bought with the tank filled has nothing in it")
	}
	w.Burn(120)
	if w.Fuel() >= full {
		t.Errorf("two hours of driving burned nothing: %d of %d", w.Fuel(), full)
	}
	// And a dry car is not a car you are driving.
	w.Burn(FuelFull * 100)
	if w.Fuel() != 0 {
		t.Errorf("the tank went to %d", w.Fuel())
	}
	if w.Driving() {
		t.Error("a car with a dry tank is still driving")
	}
	if w.Pace() != 1 {
		t.Errorf("a dry car still moves faster than walking: pace %v", w.Pace())
	}
}

// Somebody who never had a tank filled gets one; somebody who ran theirs dry
// does not get it back for free by reloading the game.
func TestSettlingFuelCannotRefillATankSomebodyEmptied(t *testing.T) {
	w := New(31)
	w.Player.Car = 1
	w.SettleFuel()
	if w.Fuel() != FuelFull {
		t.Fatalf("a car that has never been driven holds %d", w.Fuel())
	}
	w.Burn(FuelFull * 100)
	w.SettleFuel()
	if w.Fuel() != 0 {
		t.Error("a tank somebody ran dry filled itself when the world settled")
	}
}

func TestTheCityHasSomewhereToBuyPetrol(t *testing.T) {
	stations := 0
	for _, l := range Locations {
		if l.Kind == "filling" {
			stations++
			if PlaceIncome[l.ID] <= 0 {
				t.Errorf("%s sells petrol and earns nothing", l.ID)
			}
			if _, runs := TradeOf(l.ID); !runs {
				t.Errorf("%s is a filling station with no trade", l.ID)
			}
		}
	}
	if stations < 2 {
		t.Errorf("this city has %d filling stations, and one is not a business worth holding", stations)
	}
}

func TestFillingUpCostsMoneyAndFillsTheTank(t *testing.T) {
	w := New(31)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 4000, 40, 100
	w.Player.Car, w.Player.CarWear = 1, 100
	w.SettleFuel()
	station := ""
	for _, l := range Locations {
		if l.Kind == "filling" {
			station = l.ID
			break
		}
	}
	if station == "" {
		t.Fatal("nowhere sells petrol")
	}
	w.Player.Location = station
	w.Burn(FuelFull * 2)
	if w.Fuel() >= FuelFull {
		t.Fatal("the tank did not go down at all")
	}
	cash, before := w.Player.Cash, w.Fuel()
	if reason := w.FillReadiness(station); reason != "" {
		t.Fatalf("standing at a pump with money and a half-empty tank: %q", reason)
	}
	if err := w.FillUp(station); err != nil {
		t.Fatal(err)
	}
	if w.Fuel() != FuelFull {
		t.Errorf("filled up to %d of %d", w.Fuel(), FuelFull)
	}
	if w.Player.Cash >= cash {
		t.Error("the petrol was free")
	}
	t.Logf("from %d to %d for $%d", before, w.Fuel(), cash-w.Player.Cash)
	// A full tank is not worth paying for twice.
	if reason := w.FillReadiness(station); reason == "" {
		t.Error("a full tank can be filled again")
	}
}

// The link that makes it a business: the city's own drivers run dry and go and
// buy a tank, and the money lands with whoever holds the station.
func TestTheCitysDriversBuyTheirPetrolSomewhere(t *testing.T) {
	standing, sold, cities := 0, 0, 0
	for _, seed := range []uint32{7, 29, 53, 101, 199} {
		w := New(seed)
		w.District = 2
		cities++
		station := w.theStation()
		if station == "" {
			t.Fatal("nowhere sells petrol")
		}
		// Nobody holds it, so anybody standing there came to buy something.
		w.Properties[station].Owner = "independent"
		purses := w.cityPurses()
		for day := 0; day < 10; day++ {
			w.RNG, w.WorldRNG = seed+uint32(day), seed+uint32(day)
			w.DryDay()
			aDay(w)
			for i := range w.NPCs {
				if n := &w.NPCs[i]; !n.Dead && n.Location == station && n.Car > 0 {
					standing++
				}
			}
			w.PumpDay()
		}
		if spent := purses - w.cityPurses(); spent > 0 {
			sold += spent / FuelPrice
		}
	}
	t.Logf("over ten days in %d cities: %d driver-days on a forecourt and %d tanks paid for", cities, standing, sold)
	if standing == 0 {
		t.Error("nobody in this city ever went to a filling station")
	}
	if sold == 0 {
		t.Error("no petrol was ever paid for")
	}
}
