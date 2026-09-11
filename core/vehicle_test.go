package core

import (
	"fmt"
	"strings"
	"testing"
)

func driver(t *testing.T) *World {
	t.Helper()
	w := New(53)
	// The forecourt, not the motor works. Cars used to be bought at the garage
	// because it was the only place in the city that had any; the city has a
	// dealership now, and a garage repairs what you already have.
	w.Player.Location = "dealer"
	w.Player.Cash = 20000
	return w
}

// Renamed and repointed with the rule it guards: this asserted that the motor
// works sold cars, which was true of a city that had no forecourt in it.
func TestCarsAreSoldAtTheForecourtAndNowhereElse(t *testing.T) {
	t.Parallel()
	w := driver(t)
	if w.CarReadiness() != "" {
		t.Fatal("the forecourt refused to sell:", w.CarReadiness())
	}
	for _, elsewhere := range []string{"bar", "market", "docks", "room", "garage"} {
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	w := driver(t)
	w.Player.Car, w.Player.CarWear = 2, 0
	w.MigrateLivingWorld()
	if !w.Driving() || w.CarCondition() != 100 {
		t.Fatalf("migration left the car at %d", w.CarCondition())
	}
}

// "When I go to the car dealership I don't see any option to buy cars."
//
// The button lived in the garage's own case in the room switch, and the rule
// behind it asks for a forecourt: at a garage it was offered and permanently
// refused, and at the forecourt it was never offered at all. A rule and its
// button asking two different questions, which is the same fault that made
// CarSource and CarWorkshop one function in the first place.
func TestTheForecourtOffersTheCarsAndTheGarageDoesNot(t *testing.T) {
	t.Parallel()
	w := New(11)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 4000, 40, 100
	offered := func(place string) (Action, bool) {
		w.Player.Location = place
		for _, a := range w.Actions(place) {
			if a.ID == "car" {
				return a, true
			}
		}
		return Action{}, false
	}
	lot := ""
	for _, l := range Locations {
		if l.Kind == "dealer" {
			lot = l.ID
			break
		}
	}
	if lot == "" {
		t.Fatal("this city sells no cars anywhere")
	}
	a, ok := offered(lot)
	if !ok {
		t.Fatal("the forecourt does not offer a car")
	}
	if a.Disabled {
		t.Errorf("standing on the forecourt with $4,000: %q", a.Reason)
	}
	if _, ok := offered("garage"); ok {
		t.Error("a garage offers cars for sale; a garage works on the one you have")
	}
	// And the one thing a garage does have is still there.
	w.Player.Location = "garage"
	work := false
	for _, act := range w.Actions("garage") {
		if act.ID == "service" {
			work = true
		}
	}
	if !work {
		t.Error("the garage stopped working on cars")
	}
}

// A yard of cabs is a business whose whole trade is moving somebody across
// town, so what it does for the rest of what the player holds writes itself:
// they always have a ride. It was an address that paid and did nothing else,
// which is true of nine of this city's thirteen trades — the garage and the
// haulier were the only two that reached past their own income.
func TestAYardOfCabsIsAlwaysARide(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 30, 200000
	w.Player.Location = "room"
	walk := w.Journey("room", "docks")
	if w.RidingWithTheCabs() {
		t.Fatal("riding in cabs nobody holds")
	}

	w.Player.Location, w.Event = "cabstand", nil
	if err := w.apply(Command{Kind: "acquire", Target: "cabstand", RequestID: "buythecabs"}); err != nil {
		t.Fatal(err)
	}
	w.Player.Location = "room"
	if !w.RidingWithTheCabs() {
		t.Fatal("the player holds the cab yard and is still walking")
	}
	riding := w.Journey("room", "docks")
	t.Logf("on foot %d minutes, in your own cabs %d", walk, riding)
	if riding >= walk {
		t.Fatalf("a yard of cabs is worth nothing: %d against %d on foot", riding, walk)
	}

	// And a car of your own is still better than a cab, or there would be no
	// reason to buy one.
	w.Player.Car, w.Player.CarWear, w.Player.Fuelled = 3, 100, 0
	if !w.Driving() {
		t.Fatal("the car will not start")
	}
	if w.RidingWithTheCabs() {
		t.Fatal("the cabs are taking somebody who has their own car running")
	}
	own := w.Journey("room", "docks")
	if own >= riding {
		t.Fatalf("your own car is no better than a cab: %d against %d", own, riding)
	}
	t.Logf("in your own car %d minutes", own)

	// A dry tank puts them back in a cab rather than on the pavement, which is
	// the whole point of holding the yard.
	w.Player.Fuel, w.Player.Fuelled = 0, max(1, w.Minute)
	if w.Driving() {
		t.Fatal("a dry car is still driving")
	}
	if !w.RidingWithTheCabs() {
		t.Fatal("the car ran dry and nobody sent a cab")
	}
	if dry := w.Journey("room", "docks"); dry != riding {
		t.Fatalf("the cab took %d minutes with a dry car and %d without one", dry, riding)
	}
}

// And the city says which it is, because a journey the player did not choose
// the speed of is one they cannot plan around.
func TestTheCitySaysWhenACabIsTakingYou(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 30, 200000
	w.Player.Location, w.Event = "cabstand", nil
	if err := w.apply(Command{Kind: "acquire", Target: "cabstand", RequestID: "buythecabs2"}); err != nil {
		t.Fatal(err)
	}
	w.Player.Location, w.Event = "room", nil
	said := false
	for _, a := range w.Actions("docks") {
		if a.ID == "travel" {
			said = strings.Contains(a.Detail, "in one of your own cabs")
		}
	}
	if !said {
		t.Fatal("a cab is taking the player across town and the card says they are driving")
	}
	note, _ := w.Crossing("room", "docks")["note"].(string)
	if !strings.Contains(note, "own drivers") {
		t.Fatalf("the street says %q while one of the player's cabs has them", note)
	}
}
