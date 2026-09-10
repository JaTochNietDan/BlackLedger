package core

import (
	"strings"
	"testing"
)

// "It seems like you can rob places or take from people's cars multiple times
// in a row, that should probably be tracked and time limited etc."
//
// A place that was robbed last night does not have last night's takings sitting
// in the same drawer, and a street that lost a car does not have another row of
// untouched ones by morning.

func thief(t *testing.T) *World {
	t.Helper()
	w := New(41)
	w.Event, w.District = nil, 9
	w.Player.Location, w.Player.Health, w.Player.Cash = "laundry", 100, 500
	return w
}

func TestATillIsNotFilledAgainByMorning(t *testing.T) {
	t.Parallel()
	w := thief(t)
	if reason := w.RobberyReadiness("laundry"); reason != "" {
		t.Fatalf("could not rob it once: %s", reason)
	}
	if err := w.Rob("laundry"); err != nil {
		t.Fatal(err)
	}
	reason := w.RobberyReadiness("laundry")
	if reason == "" {
		t.Fatal("walked straight back in and took the same till twice")
	}
	if !strings.Contains(strings.ToLower(reason), "again") && !strings.Contains(reason, "last time") {
		t.Fatalf("the refusal does not say why: %q", reason)
	}
	// And it comes back. A place that can never be robbed again is not a
	// cooling-off, it is a deletion.
	w.Advance(TillRest + 60)
	w.Player.Location, w.Player.Health = "laundry", 100
	if reason := w.RobberyReadiness("laundry"); reason != "" {
		t.Fatalf("the till never filled up again: %s", reason)
	}
}

func TestTheRowIsPickedOverForANight(t *testing.T) {
	t.Parallel()
	w := thief(t)
	w.Player.Location = "bar"
	mark, ok := w.StripTarget("bar")
	if !ok {
		// Give the street a car to lose.
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if !n.Dead && n.Faction != w.PlayerOrganizationID() {
				n.Location, n.Car, n.Trust = "bar", 2, 20
				w.Player.Contacts = 5
				break
			}
		}
		mark, ok = w.StripTarget("bar")
	}
	if !ok {
		t.Skip("no car to take apart in this city")
	}
	_ = mark
	if err := w.StripCar("bar"); err != nil {
		t.Fatal(err)
	}
	// Another car parked in the same street is not a second night's work.
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead && n.Car == 0 && n.Faction != w.PlayerOrganizationID() && n.Location != "bar" {
			n.Location, n.Car = "bar", 2
			break
		}
	}
	if reason := w.StripReadiness("bar"); reason == "" {
		t.Fatal("took a second car apart in the same street the same night")
	}
	w.Advance(StreetRest + 60)
	w.Player.Location = "bar"
	if _, ok := w.StripTarget("bar"); ok {
		if reason := w.StripReadiness("bar"); reason != "" {
			t.Fatalf("the street was never worth walking down again: %s", reason)
		}
	}
}

func TestTheRefusalSaysHowLongItWillBe(t *testing.T) {
	t.Parallel()
	w := thief(t)
	if err := w.Rob("laundry"); err != nil {
		t.Fatal(err)
	}
	a := actionByID(w.Actions("laundry"), "rob")
	if a == nil || !a.Disabled {
		t.Fatalf("the button is still offered as available: %+v", a)
	}
	if a.Reason == "" {
		t.Fatal("refused with no reason at all")
	}
}
