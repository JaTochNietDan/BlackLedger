package core

import "testing"

// Each comfort is bought once and paid for every day after, so each one has to
// be worth that across a run of campaigns rather than in a single anecdote.

// TestADoorIsWorthPayingForEveryDay measures what a door that holds is worth
// when somebody comes for the player at home.
func TestADoorIsWorthPayingForEveryDay(t *testing.T) {
	heavy(t)
	const runs = 500
	openDoor, heldDoor := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		plain := New(seed * 2654435761)
		plain.Player.Health = 100
		plain.Attack(Plot{ID: ID(), Kind: "hit", Actor: "bellandi", Life: plain.Life, Due: plain.Minute})
		if plain.Player.Alive {
			openDoor++
		}

		barred := New(seed * 2654435761)
		barred.Player.Health, barred.Player.Cash = 100, 8000
		barred.Fit(barred.Player.Home, "door")
		barred.Attack(Plot{ID: ID(), Kind: "hit", Actor: "bellandi", Life: barred.Life, Due: barred.Minute})
		if barred.Player.Alive {
			heldDoor++
		}
	}
	if heldDoor <= openDoor {
		t.Fatalf("%d of %d survived behind a door against %d without one", heldDoor, runs, openDoor)
	}
	t.Logf("attacked at home across %d campaigns: %d survive with nothing on the door, %d with steel behind it", runs, openDoor, heldDoor)
}

// TestASafeIsWorthPayingForEveryDay measures what stays in the player's hands
// across a run of police raids at the attention where raids happen.
func TestASafeIsWorthPayingForEveryDay(t *testing.T) {
	heavy(t)
	const runs = 400
	exposed, behind := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		plain := New(seed * 2654435761)
		plain.Player.Cash, plain.Player.Heat = 1400, 95
		plain.Raid()
		exposed += plain.Player.Cash

		safe := New(seed * 2654435761)
		safe.Player.Cash, safe.Player.Heat = 1400+comfortCost("safe"), 95
		safe.Fit(safe.Player.Home, "safe")
		safe.Raid()
		behind += safe.Player.Cash
	}
	if behind <= exposed {
		t.Fatalf("a safe left $%d against $%d without one", behind/runs, exposed/runs)
	}
	kept := (behind - exposed) / runs
	if kept < comfortCost("safe")/4 {
		t.Fatalf("a safe saved $%d a raid, which will never repay $%d and $1 a day", kept, comfortCost("safe"))
	}
	t.Logf("raided at 95 attention across %d campaigns: $%d left in hand without a safe, $%d with one", runs, exposed/runs, behind/runs)
}

// TestACellarIsWorthPayingForEveryDay measures the attention a hidden hoard
// does not draw, and confirms the warrant that finds it is the price.
func TestACellarIsWorthPayingForEveryDay(t *testing.T) {
	heavy(t)
	const runs, days = 200, 20
	open, hidden := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		plain := New(seed * 2654435761)
		plain.Player.Stock = map[string]int{"moonshine": 20}
		for d := 0; d < days; d++ {
			plain.ContrabandDay()
		}
		open += plain.Player.Heat

		cellar := New(seed * 2654435761)
		cellar.Player.Cash = 8000
		cellar.Fit(cellar.Player.Home, "cellar")
		cellar.Player.Stock = map[string]int{"moonshine": 20}
		for d := 0; d < days; d++ {
			cellar.ContrabandDay()
		}
		hidden += cellar.Player.Heat
	}
	if hidden >= open {
		t.Fatalf("a cellar drew %d attention against %d in the open", hidden, open)
	}
	t.Logf("holding 20 units for %d days across %d campaigns: %d attention in the open, %d in a cellar", days, runs, open, hidden)
}

func comfortCost(id string) int {
	c, _ := ComfortByID(id)
	return c.Cost
}

// TestATelephoneIsWorthPayingForEveryDay measures how often somebody living
// alone with one contact gets a scene instead of a killing.
func TestATelephoneIsWorthPayingForEveryDay(t *testing.T) {
	heavy(t)
	const runs = 500
	silent, wired := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		plain := New(seed * 2654435761)
		plain.Player.Health, plain.Player.Contacts = 100, 1
		plain.Attack(Plot{ID: ID(), Kind: "hit", Actor: "bellandi", Life: plain.Life, Due: plain.Minute})
		if plain.Event != nil {
			silent++
		}

		hall := New(seed * 2654435761)
		hall.Player.Health, hall.Player.Contacts, hall.Player.Cash = 100, 1, 8000
		hall.Fit(hall.Player.Home, "telephone")
		hall.Attack(Plot{ID: ID(), Kind: "hit", Actor: "bellandi", Life: hall.Life, Due: hall.Minute})
		if hall.Event != nil {
			wired++
		}
	}
	if wired <= silent {
		t.Fatalf("a telephone produced a warning in %d of %d against %d without one", wired, runs, silent)
	}
	t.Logf("attacked at home with one contact across %d campaigns: warned %d times without a telephone, %d with one", runs, silent, wired)
}
