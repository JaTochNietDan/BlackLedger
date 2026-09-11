package core

import "testing"

// The two honest jobs.
//
// Carrying an envelope pays $45 for forty-five minutes and two respect, with
// nothing that can go wrong. A shift on Pier 14 pays $75 for ninety minutes and
// one respect, and one shift in seven costs ten health. Measured by the minute
// — the only honest way to compare two things that take different amounts of a
// day — the envelope pays more money, four times the standing, and carries no
// risk at all.
//
// So there was never a reason to work the docks. The first job in the game was
// strictly better than the grind it is supposed to graduate into, and it was
// unlimited: the fixer had an infinite number of envelopes.
//
// A favour is not a career. There are only so many in a day, and after them the
// pier is what is left — which is the shape the fiction always had.

func TestTheEnvelopeIsAFavourAndNotAJob(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 5000, 40
	w.Player.Location = "bar"

	ran := 0
	for i := 0; i < CourierADay+3; i++ {
		w.Event = nil
		a := actionByID(w.Actions("bar"), "courier")
		if a == nil {
			t.Fatal("the bar offers no envelope at all")
		}
		if a.Disabled {
			break
		}
		if err := w.apply(Command{Kind: "courier", Target: "bar", RequestID: ID()}); err != nil {
			t.Fatal(err)
		}
		w.Event = nil
		w.Player.Location = "bar"
		ran++
	}
	if ran != CourierADay {
		t.Fatalf("carried %d envelopes in a day where %d is the most there are", ran, CourierADay)
	}
	// Refused with a reason rather than going quiet.
	a := actionByID(w.Actions("bar"), "courier")
	if a == nil || !a.Disabled || a.Reason == "" {
		t.Fatalf("the envelope after the last one was refused without saying why: %+v", a)
	}
	// And there are more tomorrow. Advanced by a whole day rather than stepped
	// to the next midnight: the envelopes leave the clock at forty-five past,
	// so stepping by the hour never lands on it.
	w.Advance(1440)
	w.Event = nil
	w.Player.Location = "bar"
	if again := actionByID(w.Actions("bar"), "courier"); again == nil || again.Disabled {
		t.Fatalf("a new day and still no envelopes: %+v", again)
	}
}

func TestThePierIsAlwaysThere(t *testing.T) {
	t.Parallel()
	// Whatever else is refused, the docks are not. It is the floor of this
	// city's economy and the only thing somebody with nothing can always do.
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 0, 0
	w.Player.Location = "docks"
	for i := 0; i < 12; i++ {
		a := actionByID(w.Actions("docks"), "dockwork")
		if a == nil || a.Disabled {
			t.Fatalf("shift %d: the pier turned somebody away: %+v", i+1, a)
		}
		w.Event = nil
		if err := w.apply(Command{Kind: "dockwork", Target: "docks", RequestID: ID()}); err != nil {
			t.Fatal(err)
		}
		w.Event = nil
		w.Player.Location, w.Player.Health = "docks", 100
	}
}
