package core

import "testing"

// A machine needs a wall, not a floor, a dealer and somebody watching the
// dealer. That is why there was one in every bar in America and why a poolhall
// can take money off people all day without being a casino.

func TestTheMachineKeepsAKnownShareOfWhatGoesThroughIt(t *testing.T) {
	// Worked out over every one of the eight thousand lines the drums can show,
	// from the strip and the paytable, so the number on the wall cannot drift
	// away from the machine underneath it.
	edge := MachineEdge()
	t.Logf("the machine keeps %d in every hundred", edge)
	if edge <= 0 {
		t.Fatalf("a machine that gives money away is not a machine: edge %d", edge)
	}
	// Harsher than a wheel, because a bandit always was, and not so harsh that
	// nobody would put a nickel in twice.
	if edge <= HouseEdge {
		t.Errorf("the machine keeps %d and the wheel keeps %d; a bandit is the worse bet", edge, HouseEdge)
	}
	if edge > 25 {
		t.Errorf("the machine keeps %d in every hundred, which nobody would play twice", edge)
	}
}

func TestThreeOfAKindPaysWhatTheStripSays(t *testing.T) {
	for _, s := range ReelStrip() {
		if got := MachinePays([3]Symbol{s, s, s}); got != s.Pays {
			t.Errorf("three %s paid %d and the strip says %d", s.ID, got, s.Pays)
		}
	}
	// Cherries pay on their own, and nothing else does.
	var cherry, lemon, bell Symbol
	for _, s := range ReelStrip() {
		switch s.ID {
		case "cherry":
			cherry = s
		case "lemon":
			lemon = s
		case "bell":
			bell = s
		}
	}
	if got := MachinePays([3]Symbol{cherry, lemon, bell}); got != OneCherry {
		t.Errorf("one cherry paid %d", got)
	}
	if got := MachinePays([3]Symbol{cherry, cherry, bell}); got != TwoCherries {
		t.Errorf("two cherries paid %d", got)
	}
	if got := MachinePays([3]Symbol{lemon, bell, plumOf(t)}); got != 0 {
		t.Errorf("three different faces paid %d", got)
	}
}

func plumOf(t *testing.T) Symbol {
	t.Helper()
	for _, s := range ReelStrip() {
		if s.ID == "plum" {
			return s
		}
	}
	t.Fatal("no plum on the strip")
	return Symbol{}
}

// The money, both ways: what the player puts in leaves them, what comes back
// comes out of whoever holds the room.
func TestAPullTakesTheStakeAndTheHousePaysTheWins(t *testing.T) {
	w := New(17)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 3000, 40, 100
	w.Player.Location = "bar"
	if !HasMachines("bar") {
		t.Fatal("there is no machine in the bar")
	}
	stake, _ := slotStake("nickel")
	pulls, won, spent := 0, 0, 0
	for i := 0; i < 60; i++ {
		cash := w.Player.Cash
		if err := w.PullHandle("bar", 5); err != nil {
			t.Fatal(err)
		}
		pulls++
		spent += stake.Amount
		if w.Player.Cash > cash {
			won++
		}
		if w.Reels == nil || len(w.MachineDescription()["line"].([]string)) != 3 {
			t.Fatal("the machine does not say where it stopped")
		}
	}
	back := w.Player.Cash - (3000 - spent)
	t.Logf("%d pulls at $%d: %d of them paid something, and $%d of $%d came back",
		pulls, stake.Amount, won, back, spent)
	if won == 0 {
		t.Error("sixty pulls and the machine never paid anything at all")
	}
	if back >= spent {
		t.Errorf("the player is up $%d after sixty pulls; the house edge is not being taken", back-spent)
	}
}

// And you cannot play your own machine, which is the difference between owning
// one and standing in front of one.
func TestYouDoNotPlayYourOwnMachine(t *testing.T) {
	w := New(17)
	w.District = 2
	w.Player.Cash = 3000
	w.Properties["bar"].Owner = "player:1"
	stake, _ := slotStake("nickel")
	if reason := w.PullReadiness("bar", stake); reason == "" {
		t.Error("the machine on your own wall will take your own money")
	}
	if reason := w.PullReadiness("laundry", stake); reason == "" {
		t.Error("a laundry has a bandit in it")
	}
}

// The wiring: the button is offered where there is a machine, it is not offered
// where there is not, and pressing it plays the machine.
func TestTheMachineCanActuallyBePlayedFromTheRoom(t *testing.T) {
	w := New(17)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 3000, 40, 100
	w.Player.Location = "bar"
	found := false
	for _, a := range w.Actions("bar") {
		if a.ID == "pull" {
			found = true
			if a.Disabled {
				t.Errorf("the machine in the bar is refused: %s", a.Reason)
			}
		}
	}
	if !found {
		t.Fatal("the bar has a machine and does not offer it")
	}
	for _, a := range w.Actions("laundry") {
		if a.ID == "pull" {
			t.Error("a laundry offers a machine")
		}
	}
	cash := w.Player.Cash
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "pull", Amount: 5})
	if err != nil {
		t.Fatal(err)
	}
	if next.Reels == nil {
		t.Fatal("the handle was pulled and the drums never moved")
	}
	if next.Player.Cash == cash {
		t.Error("a pull cost nothing at all")
	}
	if next.Minute == w.Minute {
		t.Error("a pull took no time")
	}
}

// The point of a machine is that it does not need a casino. If every room with
// a bandit in it were a casino, there would be no reason for the city to have
// anything but casinos.
func TestMachinesStandInRoomsThatAreNotCasinos(t *testing.T) {
	casinos, machines, both := 0, 0, 0
	for _, l := range Locations {
		tables, bandits := HasTables(l.ID), HasMachines(l.ID)
		if tables {
			casinos++
		}
		if bandits {
			machines++
		}
		if bandits && !tables {
			both++
		}
		if tables && !bandits {
			t.Errorf("%s runs tables and has no machine against the wall", l.ID)
		}
	}
	t.Logf("%d rooms with tables, %d with machines, %d with machines and no tables", casinos, machines, both)
	if both == 0 {
		t.Error("every machine in this city is in a casino, so a machine is no reason for anywhere else")
	}
}
