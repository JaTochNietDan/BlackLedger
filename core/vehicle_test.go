package core

import (
	"fmt"
	"testing"
)

func driver(t *testing.T) *World {
	t.Helper()
	w := New(53)
	w.Player.Location = "garage"
	w.Player.Cash = 20000
	return w
}

func TestCarsAreSoldAtTheMotorWorksAndNowhereElse(t *testing.T) {
	w := driver(t)
	if w.CarReadiness() != "" {
		t.Fatal("the motor works refused to sell:", w.CarReadiness())
	}
	for _, elsewhere := range []string{"bar", "market", "docks", "room"} {
		w.Player.Location = elsewhere
		if w.CarReadiness() == "" {
			t.Fatalf("%s was selling cars", elsewhere)
		}
		if err := w.BuyVehicle(); err == nil {
			t.Fatalf("bought a car at %s", elsewhere)
		}
	}
}

func TestACarBuysTimeAndNothingElseDoes(t *testing.T) {
	w := driver(t)
	walk := TravelMinutes("garage", "docks")
	if w.Journey("garage", "docks") != walk {
		t.Fatal("walking was not walking")
	}
	previous := walk
	for tier := 1; tier < len(vehicles); tier++ {
		if err := w.BuyVehicle(); err != nil {
			t.Fatal(err)
		}
		got := w.Journey("garage", "docks")
		if got >= previous {
			t.Fatalf("tier %d took %d minutes against %d", tier, got, previous)
		}
		previous = got
	}
	if w.CarReadiness() == "" {
		t.Fatal("there was something better than a Packard")
	}
}

func TestACarInPoorOrderIsWorthLessAndAWreckIsWorthNothing(t *testing.T) {
	w := driver(t)
	w.BuyVehicle()
	w.BuyVehicle() // the Hudson, with a false floor
	fast := w.Journey("garage", "docks")
	w.Player.CarWear = 60
	middling := w.Journey("garage", "docks")
	if middling <= fast {
		t.Fatalf("a worn car was as quick: %d against %d", middling, fast)
	}
	w.Damage(40) // below Wreck
	if w.Driving() {
		t.Fatal("a wreck was still a car")
	}
	if w.Journey("garage", "docks") != TravelMinutes("garage", "docks") {
		t.Fatal("a wreck still made journeys shorter")
	}
	if w.Concealed() != 0 {
		t.Fatal("a wreck still hid stock")
	}
	if !w.hasRecord("It will not start") {
		t.Fatal("nobody was told the car had stopped")
	}
}

func TestAFalseFloorHidesStockFromAttentionAndFromASearch(t *testing.T) {
	w := driver(t)
	w.BuyVehicle()
	w.BuyVehicle()
	compartment := VehicleByTier(w.Player.Car).Compartment
	w.Player.Stock = map[string]int{"moonshine": compartment - 5}
	if w.Exposed() != 0 {
		t.Fatalf("%d units were in the open with room to hide them", w.Exposed())
	}
	heat := w.Player.Heat
	w.ContrabandDay()
	if w.Player.Heat != heat {
		t.Fatal("stock under the floor still drew attention")
	}
	if lost := w.Seize("A search."); lost != 0 {
		t.Fatalf("a search found %d hidden units", lost)
	}
	if w.Holding("moonshine") != compartment-5 {
		t.Fatalf("%d units survived the search", w.Holding("moonshine"))
	}

	// More than it can hide is more than it can hide.
	w.Player.Stock["moonshine"] = compartment + 12
	if w.Exposed() != 12 {
		t.Fatalf("%d units in the open", w.Exposed())
	}
	if lost := w.Seize("A search."); lost != 12 {
		t.Fatalf("a search took %d units", lost)
	}
	if w.Holding("moonshine") != compartment {
		t.Fatalf("%d left after the search", w.Holding("moonshine"))
	}
}

func TestACarCostsSomethingEveryDayAndLessAtYourOwnGarage(t *testing.T) {
	w := driver(t)
	before := w.DailyCost()
	w.BuyVehicle()
	w.BuyVehicle()
	w.BuyVehicle()
	full := w.CarUpkeep()
	if full <= 0 || w.DailyCost() != before+full {
		t.Fatalf("upkeep %d, bill went %d to %d", full, before, w.DailyCost())
	}
	w.Properties["garage"].Owner = fmt.Sprintf("player:%d", w.Life)
	if w.CarUpkeep() >= full {
		t.Fatalf("owning the motor works cost %d against %d", w.CarUpkeep(), full)
	}
	if w.ServiceFee() != 0 {
		t.Fatal("your own people charged you to work on your own car")
	}
}

func TestServicingPutsACarBackOnTheRoad(t *testing.T) {
	w := driver(t)
	w.BuyVehicle()
	w.Player.CarWear = 20
	if w.Driving() {
		t.Fatal("a wreck was still driveable")
	}
	cash := w.Player.Cash
	if err := w.Service("garage"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-CarService || w.CarCondition() != 75 || !w.Driving() {
		t.Fatalf("cash %d, condition %d", w.Player.Cash, w.CarCondition())
	}
	if w.ServiceReadiness("bar") == "" {
		t.Fatal("a bar worked on the car")
	}
}

func TestDrivingToARobberyLeavesATrail(t *testing.T) {
	onFoot, driven := 0, 0
	for seed := uint32(1); seed <= 300; seed++ {
		walk := New(seed * 2654435761)
		walk.Player.Location, walk.Player.Health, walk.Player.Respect = "bar", 100, 40
		walk.Rob("bar")
		onFoot += walk.Player.Heat

		drive := New(seed * 2654435761)
		drive.Player.Location, drive.Player.Health, drive.Player.Respect = "bar", 100, 40
		drive.Player.Car, drive.Player.CarWear = 3, 100
		drive.Rob("bar")
		driven += drive.Player.Heat
	}
	if driven <= onFoot {
		t.Fatalf("driving drew %d attention against %d on foot across 300 robberies", driven, onFoot)
	}
	t.Logf("300 robberies: %d attention on foot, %d driving an armoured car", onFoot, driven)
}

func TestAWarrantThatFindsTheFloorTakesTheCar(t *testing.T) {
	taken, kept := 0, 0
	for seed := uint32(1); seed <= 400; seed++ {
		w := New(seed * 2654435761)
		w.Player.Car, w.Player.CarWear = 2, 100
		w.Player.Stock = map[string]int{"moonshine": 15}
		w.Player.Cash, w.Player.Heat = 3000, 60
		w.Raid()
		if w.Player.Car == 0 {
			taken++
			if w.Carrying() != 0 {
				t.Fatal("the car went and the stock under its floor did not")
			}
		} else {
			kept++
			if w.Carrying() != 15 {
				t.Fatalf("the car was kept but %d units went", 15-w.Carrying())
			}
		}
	}
	if taken == 0 || kept == 0 {
		t.Fatalf("a warrant took the car in %d of 400 raids, which is not a risk", taken)
	}
	t.Logf("a warrant found the false floor in %d of 400 raids", taken)
}

func TestNobodyInheritsACar(t *testing.T) {
	w := driver(t)
	w.BuyVehicle()
	w.Player.Alive = false
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Car != 0 || next.Driving() {
		t.Fatal("the next arrival found a car waiting")
	}
}

func TestAnOlderSaveIsDrivingSomethingThatRuns(t *testing.T) {
	w := driver(t)
	w.Player.Car, w.Player.CarWear = 2, 0
	w.MigrateLivingWorld()
	if !w.Driving() || w.CarCondition() != 100 {
		t.Fatalf("migration left the car at %d", w.CarCondition())
	}
}
