package core

import "testing"

func traveller(t *testing.T) *World {
	t.Helper()
	w := New(19)
	w.Player.Location = TripDepart
	w.Player.Cash = 6000
	w.Player.Health = 100
	return w
}

func TestJourneysAreBookedAtTheExchangeAndNowhereElse(t *testing.T) {
	t.Parallel()
	w := traveller(t)
	if w.TripReadiness("halloway") != "" {
		t.Fatal("the exchange would not sell a ticket:", w.TripReadiness("halloway"))
	}
	for _, elsewhere := range []string{"bar", "docks", "club", "room"} {
		w.Player.Location = elsewhere
		if w.TripReadiness("halloway") == "" {
			t.Fatalf("booked a journey at %s", elsewhere)
		}
		if err := w.Trip("halloway"); err == nil {
			t.Fatalf("travelled from %s", elsewhere)
		}
	}
}

func TestTheCityRunsWhileYouAreNotInIt(t *testing.T) {
	t.Parallel()
	w := traveller(t)
	minute := w.Minute
	if err := w.Trip("halloway"); err != nil {
		t.Fatal(err)
	}
	d, _ := DestinationByID("halloway")
	if w.Minute-minute < d.Days*1440 {
		t.Fatalf("four days away took %d minutes", w.Minute-minute)
	}
	if w.Player.Location != TripDepart {
		t.Fatalf("came back to %q", w.Player.Location)
	}
	if !w.hasRecord("Back in Bellwether") {
		t.Fatal("nobody noticed they came back")
	}
}

func TestBeingNowhereAnybodyIsLookingLowersAttention(t *testing.T) {
	t.Parallel()
	w := traveller(t)
	w.Player.Heat = 60
	if err := w.Trip("halloway"); err != nil {
		t.Fatal(err)
	}
	d, _ := DestinationByID("halloway")
	if w.Player.Heat > 60-d.Days*d.Relief {
		t.Fatalf("attention only fell to %d", w.Player.Heat)
	}

	// And it cannot fall below nothing, however long the trip.
	quiet := traveller(t)
	quiet.Player.Heat = 2
	if err := quiet.Trip("halloway"); err != nil {
		t.Fatal(err)
	}
	if quiet.Player.Heat != 0 {
		t.Fatalf("attention went to %d", quiet.Player.Heat)
	}
}

func TestCountryPricesAreWorthTheFareOnlyIfYouCanHideTheLoad(t *testing.T) {
	t.Parallel()
	// The crates are bought at the other end and the run home is the risk, so
	// the honest measure is what a season of these trips is worth rather than
	// whether any one of them arrives.
	const runs = 300
	run := func(car int) (int, int, int) {
		spent, landed, jumped := 0, 0, 0
		for seed := uint32(1); seed <= runs; seed++ {
			w := traveller(t)
			w.WorldRNG = seed * 2654435761
			if car > 0 {
				w.Player.Car, w.Player.CarWear = car, 100
			}
			price := w.Good("moonshine").Price
			cost, before := w.TripCost("rockridge"), w.Carrying()
			if err := w.Trip("rockridge"); err != nil {
				t.Fatal(err)
			}
			spent += cost
			landed += (w.Carrying() - before) * price
			if w.Carrying()-before < w.CarryLimit() {
				jumped++
			}
		}
		return spent, landed, jumped
	}

	pocketSpent, pocketLanded, pocketJumped := run(0)
	if pocketJumped == 0 {
		t.Fatal("a load in your pockets came home safe every single time, so the run is not a risk")
	}
	if pocketLanded > pocketSpent {
		t.Fatalf("carrying it in your pockets turned $%d into $%d, so a hiding place buys nothing", pocketSpent, pocketLanded)
	}

	floorSpent, floorLanded, floorJumped := run(2)
	if floorLanded <= floorSpent {
		t.Fatalf("even with a false floor, $%d of trips landed $%d, so there is no reason to go", floorSpent, floorLanded)
	}
	t.Logf("%d trips to Rockridge in your pockets: $%d spent, $%d landed, jumped %d times", runs, pocketSpent, pocketLanded, pocketJumped)
	t.Logf("%d trips with a false floor: $%d spent, $%d landed, jumped %d times", runs, floorSpent, floorLanded, floorJumped)
}

func TestSomewhereToPutItIsWhatDecidesTheLoadAndWhatSurvives(t *testing.T) {
	t.Parallel()
	w := traveller(t)
	plain := w.CarryLimit()
	w.Player.Car, w.Player.CarWear = 2, 100
	if driving := w.CarryLimit(); driving <= plain {
		t.Fatalf("a false floor carried %d against %d", driving, plain)
	}
	w.Properties[w.Player.Home].Comforts = []string{"cellar"}
	if both := w.CarryLimit(); both <= w.Concealed() {
		t.Fatalf("a cellar and a car carried %d", both)
	}

	// What is hidden is what survives being jumped, so a hiding place is worth
	// its price on the road as well as in a raid.
	const runs = 300
	openLoad, hiddenLoad := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		pockets := traveller(t)
		pockets.WorldRNG = seed * 2654435761
		pockets.Trip("rockridge")
		openLoad += pockets.Carrying()

		floor := traveller(t)
		floor.WorldRNG = seed * 2654435761
		floor.Player.Car, floor.Player.CarWear = 2, 100
		floor.Trip("rockridge")
		hiddenLoad += floor.Carrying()
	}
	if hiddenLoad <= openLoad {
		t.Fatalf("a false floor landed %d crates against %d in pockets", hiddenLoad, openLoad)
	}
	t.Logf("%d runs home: %d crates land carrying them in your pockets, %d with a false floor", runs, openLoad, hiddenLoad)
}

func TestTheBankIsCheaperInPerson(t *testing.T) {
	t.Parallel()
	w := traveller(t)
	w.Offshore = 2000
	if w.Player.Offshore {
		t.Fatal("the account already answered to them")
	}
	cost := w.TripCost("kingsport")
	d, _ := DestinationByID("kingsport")
	if cost != d.Fare+KingsportAccess {
		t.Fatalf("the trip was quoted at $%d", cost)
	}
	if cost >= AccessCost {
		t.Fatalf("going in person cost $%d against $%d by wire, so there is no reason to go", cost, AccessCost)
	}
	if err := w.Trip("kingsport"); err != nil {
		t.Fatal(err)
	}
	if !w.Player.Offshore {
		t.Fatal("the account still did not answer to them")
	}

	// With nothing out there and nothing to arrange, there is nothing to go for.
	empty := traveller(t)
	empty.Player.Offshore, empty.Offshore = true, 0
	if empty.TripReadiness("kingsport") == "" {
		t.Fatal("sold a ticket to arrange nothing")
	}
}

func TestHallowayIsWhereYouMeetSomebody(t *testing.T) {
	t.Parallel()
	w := traveller(t)
	w.Player.Contacts = 1
	if err := w.Trip("halloway"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Contacts != 2 {
		t.Fatalf("came back with %d contacts", w.Player.Contacts)
	}
	// And there is a ceiling, so it is not an infinite supply.
	full := traveller(t)
	full.Player.Contacts = 5
	if err := full.Trip("halloway"); err != nil {
		t.Fatal(err)
	}
	if full.Player.Contacts != 5 {
		t.Fatalf("contacts went to %d", full.Player.Contacts)
	}
}

func TestYouCannotTravelInNoConditionToTravel(t *testing.T) {
	t.Parallel()
	w := traveller(t)
	w.Player.Health = 20
	for _, id := range []string{"rockridge", "kingsport", "halloway"} {
		if w.TripReadiness(id) == "" {
			t.Fatalf("sold a ticket to %s to somebody who could barely stand", id)
		}
	}
}

func TestTheQuotedPriceIsThePriceCharged(t *testing.T) {
	t.Parallel()
	for _, id := range []string{"rockridge", "kingsport", "halloway"} {
		w := traveller(t)
		w.Offshore = 2000
		quoted, cash := w.TripCost(id), w.Player.Cash
		if err := w.Trip(id); err != nil {
			t.Fatal(id, err)
		}
		// Days away cost money too — bills fall due — so the check is that the
		// fare and the goods came to exactly what was quoted, before the city
		// took its own.
		if cash-w.Player.Cash < quoted {
			t.Fatalf("%s quoted $%d and took $%d", id, quoted, cash-w.Player.Cash)
		}
	}
}

// TestLeavingTownIsAWayToSurviveAWeek measures the one thing a journey buys
// that nothing else in the game does, and what it costs: whatever was arranged
// for the player happens to the house instead of to them.
func TestLeavingTownIsAWayToSurviveAWeek(t *testing.T) {
	t.Parallel()
	const runs = 300
	stayed, left, wrecked := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		home := traveller(t)
		home.WorldRNG = seed * 2654435761
		home.Player.Location = home.Player.Home
		home.RetaliationFrom("bellandi")
		home.Advance(4 * 1440)
		if home.Player.Alive {
			stayed++
		}

		away := traveller(t)
		away.WorldRNG = seed * 2654435761
		condition := away.Properties[away.Player.Home].Condition
		away.RetaliationFrom("bellandi")
		if err := away.Trip("halloway"); err != nil {
			t.Fatal(err)
		}
		if away.Player.Alive {
			left++
		}
		if away.Properties[away.Player.Home].Condition < condition {
			wrecked++
		}
	}
	if left <= stayed {
		t.Fatalf("%d of %d survived leaving against %d staying, so a journey buys nothing", left, runs, stayed)
	}
	if wrecked == 0 {
		t.Fatal("leaving town cost the house nothing, so it is free")
	}
	t.Logf("with a hit arranged and nobody watching the door: %d of %d survive staying home, %d survive leaving town, and the house is damaged in %d of those",
		stayed, runs, left, wrecked)
}
