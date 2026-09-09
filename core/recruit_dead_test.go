package core

import (
	"strings"
	"testing"
)

// A save driven to day 58 held a crew of one: Leo Carver, dead. The delegate
// button in that room said so — "Leo Carver is dead" — while the recruit
// button beside it refused for the opposite reason, "Leo is already in your
// crew". Both cannot be true.
//
// The label and the gate had already been taught that nobody drives forever:
// both name whoever actually holds the job. The effect had not. It appended a
// hardcoded Leo Carver whatever the button said, so once the city gave the
// wheel to somebody else the player hired a man they had buried.

func TestYouHireThePersonTheButtonNamed(t *testing.T) {
	w := New(4)
	w.Player.Location = "bar"
	w.Player.Respect = 20
	// Leo is dead and the city has found somebody else to drive.
	w.Kill("leo", "Shot twice outside the Mariner.")
	w.FillRoles()
	driver := w.Holder("driver")
	if driver == nil || driver.ID == "leo" {
		t.Fatalf("nobody else took the wheel: %+v", driver)
	}
	label := ""
	for _, a := range w.Actions("bar") {
		if a.ID == "recruit" {
			if a.Disabled {
				t.Fatalf("the new driver cannot be hired at all: %q", a.Reason)
			}
			label = a.Label
		}
	}
	if !strings.Contains(label, driver.Name) {
		t.Fatalf("the button names somebody who does not drive: %q", label)
	}
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "recruit", Target: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	w = next
	if len(w.Player.Crew) != 1 {
		t.Fatalf("nobody was hired: %+v", w.Player.Crew)
	}
	hired := w.Player.Crew[0]
	if hired.ID != driver.ID {
		t.Fatalf("the button offered %s and the game hired %s (%s)", driver.Name, hired.Name, hired.ID)
	}
	if n := w.NPC(hired.ID); n == nil || n.Dead {
		t.Fatalf("a dead man is on the books: %+v", hired)
	}
}
