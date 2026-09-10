package core

import "testing"

// A casino was a room the player played in and a number that ran itself. The
// city's own people never sat down at a table, so the money behind the tables
// came from nowhere the city could feel: nobody in Bellwether was ever poorer
// on a Tuesday because of a Monday night.
//
// They play now, and what they lose is what the house takes.

func floor(t *testing.T, seed uint32) (*World, *NPC) {
	t.Helper()
	w := New(seed)
	w.District = 2
	var punter *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead {
			continue
		}
		if punter == nil {
			n.Location, n.Purse = "casino", 900
			punter = n
			continue
		}
		if n.Location == "casino" {
			n.Location = "bar" // one named person on the floor, so it can be followed
		}
	}
	if punter == nil {
		t.Fatal("nobody alive to play")
	}
	return w, punter
}

func TestTheFloorIsThePeopleStandingInTheRoom(t *testing.T) {
	t.Parallel()
	w, punter := floor(t, 31)
	on := w.OnTheFloor("casino")
	if len(on) != 1 || on[0].ID != punter.ID {
		t.Fatalf("the floor holds %d people and the room holds one", len(on))
	}
	if len(w.OnTheFloor("laundry")) != 0 {
		t.Error("a laundry has a gaming floor")
	}
}

// What the city loses is exactly what the house takes. Money must not appear.
func TestWhatTheCityLosesAtTheTablesIsWhatTheHouseTakes(t *testing.T) {
	t.Parallel()
	w, _ := floor(t, 31)
	w.Properties["casino"].Owner = "bellandi"
	house := w.faction("bellandi")
	if house == nil {
		t.Fatal("nobody holds the room")
	}
	purses, cash := w.cityPurses(), house.Cash
	nights := 0
	for day := 0; day < 60; day++ {
		w.WorldRNG = 1000 + uint32(day)
		w.Minute += 1440
		before := w.cityPurses()
		w.TableNight()
		if w.cityPurses() != before {
			nights++
		}
	}
	if nights == 0 {
		t.Fatal("sixty nights and nobody sat down")
	}
	lost, took := purses-w.cityPurses(), house.Cash-cash
	if lost != took {
		t.Errorf("the city is $%d poorer and the house is $%d richer", lost, took)
	}
	if took <= 0 {
		t.Errorf("sixty nights against the house edge and the house is down $%d", -took)
	}
	t.Logf("%d nights played, the city lost $%d and the house took $%d", nights, lost, took)
}

// And nobody plays with money they do not have. An empty purse is the obvious
// case; the one worth guarding is the purse that is not empty but is not worth
// sitting down with, because that is the rule a stake calculation quietly drops
// and nothing else notices.
func TestNobodyPlaysWithMoneyTheyCannotStandToLose(t *testing.T) {
	t.Parallel()
	for _, purse := range []int{0, 4 * TableFloor / 2} {
		w, punter := floor(t, 31)
		w.Properties["casino"].Owner = "bellandi"
		punter.Purse = purse
		if stake := w.TableStakeFor(punter); stake != 0 {
			t.Errorf("somebody holding $%d was put down for $%d", purse, stake)
		}
		for day := 0; day < 20; day++ {
			w.WorldRNG = uint32(day)
			w.Minute += 1440
			w.TableNight()
		}
		if punter.Purse != purse {
			t.Errorf("somebody holding $%d played anyway and came away with $%d", purse, punter.Purse)
		}
	}
}

// cityPurses is every penny in the city's pockets, for proving money is moved
// rather than made.
func (w *World) cityPurses() int {
	total := 0
	for i := range w.NPCs {
		total += w.NPCs[i].Purse
	}
	return total
}
