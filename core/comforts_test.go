package core

import (
	"fmt"
	"testing"
)

func resident(t *testing.T) *World {
	t.Helper()
	w := New(67)
	w.Player.Location = w.Player.Home
	w.Player.Cash = 8000
	return w
}

func TestComfortsCanOnlyBeFittedWhereYouLive(t *testing.T) {
	t.Parallel()
	w := resident(t)
	if w.FitReadiness(w.Player.Home, "door") != "" {
		t.Fatal("could not fit a door at home:", w.FitReadiness(w.Player.Home, "door"))
	}
	for _, elsewhere := range []string{"apartment", "estate", "bar", "laundry"} {
		if elsewhere == w.Player.Home {
			continue
		}
		if w.FitReadiness(elsewhere, "door") == "" {
			t.Fatalf("fitted a door at %s", elsewhere)
		}
		if err := w.Fit(elsewhere, "door"); err == nil {
			t.Fatalf("fitted a door at %s", elsewhere)
		}
	}
	if err := w.Fit(w.Player.Home, "door"); err != nil {
		t.Fatal(err)
	}
	if w.FitReadiness(w.Player.Home, "door") == "" {
		t.Fatal("fitted a second door to the same doorway")
	}
}

func TestComfortsStayWithTheBuilding(t *testing.T) {
	t.Parallel()
	w := resident(t)
	w.Fit(w.Player.Home, "door")
	w.Fit(w.Player.Home, "telephone")
	room := w.Player.Home
	if !w.Fitted("door") || w.Reach() != w.Player.Contacts+1 {
		t.Fatal("what was fitted did not take effect")
	}
	upkeep := w.ComfortUpkeep()
	if upkeep == 0 {
		t.Fatal("comforts cost nothing to keep")
	}

	w.Player.Home = "apartment"
	if w.Fitted("door") || w.Fitted("telephone") {
		t.Fatal("the fittings followed the player to a different building")
	}
	if w.ComfortUpkeep() != 0 {
		t.Fatalf("still paying $%d a day for somewhere they no longer live", w.ComfortUpkeep())
	}
	w.Player.Home = room
	if !w.Fitted("door") || w.ComfortUpkeep() != upkeep {
		t.Fatal("moving back found the building stripped")
	}
}

func TestADoorThatHoldsCountsAsAGuard(t *testing.T) {
	t.Parallel()
	w := resident(t)
	before := w.Guard()
	w.Fit(w.Player.Home, "door")
	if w.Guard() != before+1 {
		t.Fatalf("guard went %d to %d", before, w.Guard())
	}
}

func TestATelephoneCountsAsAContact(t *testing.T) {
	t.Parallel()
	w := resident(t)
	w.Player.Contacts = 1
	if w.Reach() != 1 {
		t.Fatalf("reach was %d without a telephone", w.Reach())
	}
	w.Fit(w.Player.Home, "telephone")
	if w.Reach() != 2 {
		t.Fatalf("reach was %d with one", w.Reach())
	}
}

func TestASafeKeepsMoneyOutOfReach(t *testing.T) {
	t.Parallel()
	w := resident(t)
	w.Player.Cash = 3000
	if w.Reachable() != 3000 || w.Sheltered() != 0 {
		t.Fatal("money was sheltered without a safe")
	}
	w.Fit(w.Player.Home, "safe")
	cash := w.Player.Cash
	if w.Sheltered() != SafeShelter || w.Reachable() != cash-SafeShelter {
		t.Fatalf("sheltered %d of %d", w.Sheltered(), cash)
	}
	// A fine cannot take what it cannot reach.
	w.Player.Heat = 90
	w.Raid()
	if w.Player.Cash < SafeShelter {
		t.Fatalf("a raid left $%d, less than the safe holds", w.Player.Cash)
	}

	// And the safe is not a vault: it shelters what it holds and no more.
	poor := resident(t)
	poor.Fit(poor.Player.Home, "safe")
	poor.Player.Cash = 400
	if poor.Reachable() != 0 || poor.Sheltered() != 400 {
		t.Fatalf("sheltered %d, reachable %d", poor.Sheltered(), poor.Reachable())
	}
}

func TestACellarHidesStockAndAWarrantFindsIt(t *testing.T) {
	t.Parallel()
	w := resident(t)
	w.Fit(w.Player.Home, "cellar")
	if w.Concealed() != CellarHold {
		t.Fatalf("a cellar hid %d units", w.Concealed())
	}
	w.Player.Stock = map[string]int{"moonshine": CellarHold - 5}
	heat := w.Player.Heat
	w.ContrabandDay()
	if w.Player.Heat != heat {
		t.Fatal("stock in a cellar still drew attention")
	}
	found := w.CellarFound()
	if found != CellarHold-5 || w.Carrying() != 0 {
		t.Fatalf("a warrant found %d of %d units", found, CellarHold-5)
	}
}

func TestACellarAndAFalseFloorAreDifferentHidingPlaces(t *testing.T) {
	t.Parallel()
	w := resident(t)
	w.Fit(w.Player.Home, "cellar")
	w.Player.Car, w.Player.CarWear = 2, 100
	floor := VehicleByTier(2).Compartment
	if w.Concealed() != CellarHold+floor {
		t.Fatalf("together they hid %d units", w.Concealed())
	}
	w.Player.Stock = map[string]int{"moonshine": CellarHold + floor}
	if found := w.CellarFound(); found != CellarHold {
		t.Fatalf("a search of the cellar took %d units", found)
	}
	if w.Carrying() != floor {
		t.Fatalf("%d units left under the floor of the car", w.Carrying())
	}
}

func TestComfortsJoinTheDailyBill(t *testing.T) {
	t.Parallel()
	w := resident(t)
	before := w.DailyCost()
	total := 0
	for _, c := range comforts {
		if err := w.Fit(w.Player.Home, c.ID); err != nil {
			t.Fatal(err)
		}
		total += c.Upkeep
	}
	if w.DailyCost() != before+total {
		t.Fatalf("the bill went %d to %d for $%d of upkeep", before, w.DailyCost(), total)
	}
}

func TestNobodyInheritsAFittedHome(t *testing.T) {
	t.Parallel()
	w := resident(t)
	w.Fit(w.Player.Home, "safe")
	home := w.Player.Home
	w.Player.Alive = false
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	// The building keeps what was fitted to it. The next arrival starts in the
	// same rented room, which is exactly the kind of inheritance this city
	// allows: the thing is there, the money in it is not.
	if len(next.Comforts(home)) == 0 {
		t.Fatal("the safe left the building with its owner")
	}
	if next.Player.Cash != 90 {
		t.Fatalf("the next arrival started with $%d", next.Player.Cash)
	}
}

func TestAnOlderSaveHasAnEmptyHome(t *testing.T) {
	t.Parallel()
	w := resident(t)
	w.MigrateLivingWorld()
	if len(w.Comforts(w.Player.Home)) != 0 || w.ComfortUpkeep() != 0 {
		t.Fatal("migration fitted out a home nobody paid for")
	}
	if w.Reach() != w.Player.Contacts || w.Sheltered() != 0 {
		t.Fatal("migration handed out perks")
	}
}

func TestEveryComfortIsWorthWhatItCosts(t *testing.T) {
	t.Parallel()
	// Each one has to change something measurable, or it is decoration with a
	// price on it.
	for _, c := range comforts {
		w := resident(t)
		w.Player.Contacts, w.Player.Cash = 1, 8000
		w.Player.Stock = map[string]int{"moonshine": 10}
		before := fmt.Sprint(w.Guard(), w.Reach(), w.Sheltered(), w.Concealed())
		if err := w.Fit(w.Player.Home, c.ID); err != nil {
			t.Fatal(err)
		}
		if after := fmt.Sprint(w.Guard(), w.Reach(), w.Sheltered(), w.Concealed()); after == before {
			t.Fatalf("%s changed nothing measurable", c.ID)
		}
		if c.Upkeep <= 0 || c.Cost <= 0 {
			t.Fatalf("%s is free", c.ID)
		}
	}
}
