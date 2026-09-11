package core

import (
	"strings"
	"testing"
)

// Saint Agnes, and hearing things.
//
// The busiest room in the city: thirty people drink there of an evening, the
// fixer stands at one end, the envelope job starts there, there is a game
// behind it, and two families that will not speak will speak in the back. All
// of that was true whoever held the deed. What a bar is, that nothing else in
// this city is, is somewhere everything gets said out loud — so a bar of your
// own is a line into the city, the same size as a telephone in the hall.

const theSaloon = "bar"

func drinking(t *testing.T, hold bool) *World {
	t.Helper()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 90
	if hold {
		own(w, theSaloon)
	}
	w.Player.Location = theSaloon
	w.Event = nil
	return w
}

func TestABarOfYourOwnIsALineIntoTheCity(t *testing.T) {
	t.Parallel()
	theirs, mine := drinking(t, false), drinking(t, true)
	if theirs.EarsInTheRoom() != 0 {
		t.Fatal("somebody else's bar was listening for you")
	}
	if mine.EarsInTheRoom() != BarReach {
		t.Fatalf("your own bar is worth %d", mine.EarsInTheRoom())
	}
	if mine.Reach() != theirs.Reach()+BarReach {
		t.Fatalf("reach went from %d to %d", theirs.Reach(), mine.Reach())
	}
	// The same size as the telephone, and they stack, because they are two
	// different ways of hearing the same city.
	wired := drinking(t, true)
	if home := wired.Properties[wired.Player.Home]; home != nil {
		home.Comforts = append(home.Comforts, "telephone")
	}
	if !wired.Fitted("telephone") {
		t.Skip("no telephone in this build")
	}
	if wired.Reach() != mine.Reach()+1 {
		t.Fatalf("a bar and a telephone together are worth %d where the bar alone is %d",
			wired.Reach(), mine.Reach())
	}
}

func TestNothingIsSaidInARoomThatIsNotOpen(t *testing.T) {
	t.Parallel()
	w := drinking(t, true)
	trade, _ := TradeOf(theSaloon)
	prop := w.Properties[theSaloon]
	prop.Trouble = true
	if w.TheBar() != "" {
		t.Fatal("a bar with the cellar flooded was still worth listening in")
	}
	prop.Trouble = false
	prop.Staff = trade.Hands - 1
	if w.TheBar() != "" {
		t.Fatalf("short-handed at %d of %d and still hearing everything", prop.Staff, trade.Hands)
	}
	prop.Staff = trade.Hands
	if w.TheBar() != theSaloon {
		t.Fatal("a working bar of yours hears nothing")
	}
}

func TestWhatTheBarBuysIsWorthSomething(t *testing.T) {
	t.Parallel()
	// Reach is not a number for its own sake. It decides whether anybody warns
	// you that one side came to a sitdown to finish it rather than settle it,
	// and that warning is the difference between walking into a room and
	// knowing what is in it.
	quiet, mine := drinking(t, false), drinking(t, true)
	for _, w := range []*World{quiet, mine} {
		for i := range w.Conflicts {
			w.Conflicts[i].State, w.Conflicts[i].Hostility = "war", 90
		}
		w.Player.Contacts = 1
	}
	q, ok := quiet.OpenQuarrel()
	if !ok || !q.Trap {
		t.Skip("nobody in this city came to finish anything")
	}
	if q.Suspected {
		t.Fatalf("warned at a reach of %d", quiet.Reach())
	}
	if warned, _ := mine.OpenQuarrel(); !warned.Suspected {
		t.Fatalf("the same city, a bar of your own, a reach of %d, and nobody said anything",
			mine.Reach())
	}
}

func TestTheRoomSaysWhatItIsWorth(t *testing.T) {
	t.Parallel()
	mine := drinking(t, true)
	note := mine.PlaceNote(theSaloon)
	if !strings.Contains(note, "out loud") {
		t.Fatalf("a bar of your own says: %q", note)
	}
	theirs := drinking(t, false)
	if strings.Contains(theirs.PlaceNote(theSaloon), "out loud") {
		t.Fatalf("a bar that is not yours says: %q", theirs.PlaceNote(theSaloon))
	}
}
