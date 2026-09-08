package core

import "testing"

// People move between buildings now, and a player standing in a room never
// noticed. Somebody they had been talking to walked out of the door and was
// simply not in the list any more; somebody else appeared in it with no
// explanation. A room where people arrive and leave without remark is a list
// that changes behind your back rather than a place you are standing in.

func TestYouSeeSomebodyLeaveTheRoomYouAreStandingIn(t *testing.T) {
	w := New(4)
	w.Player.Location = "bar"
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	w.SetOut()
	if len(w.Comings) != 1 {
		t.Fatalf("she walked out of the room and the city noted %d things", len(w.Comings))
	}
	saw := w.Comings[0]
	if saw.ID != n.ID || saw.Leaving != true || saw.Where != "bar" {
		t.Fatalf("what was noticed was %+v", saw)
	}
	if saw.Note == "" || saw.To != "Bluebird Laundry" {
		t.Fatalf("the room cannot say where she went: %+v", saw)
	}
}

func TestYouSeeSomebodyComeIn(t *testing.T) {
	w := New(4)
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	w.SetOut()
	// Stand where she is going, then let the walk finish.
	w.Player.Location = "laundry"
	w.Comings = nil
	w.Minute = n.Arrives
	w.Arrivals()
	if len(w.Comings) != 1 {
		t.Fatalf("she walked in and the city noted %d things", len(w.Comings))
	}
	saw := w.Comings[0]
	if saw.ID != n.ID || saw.Leaving || saw.Where != "laundry" || saw.Note == "" {
		t.Fatalf("what was noticed was %+v", saw)
	}
}

func TestNothingIsRemarkedOnInARoomYouAreNotIn(t *testing.T) {
	w := New(4)
	// The Bellwether Herald: an address with nobody in it who has anywhere to
	// be, so anything reported here came from another room.
	w.Player.Location = "herald"
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	w.SetOut()
	w.Minute = n.Arrives
	w.Arrivals()
	for _, c := range w.Comings {
		if c.Where != w.Player.Location {
			t.Fatalf("the player was at the Herald and was told about %s: %+v", c.Where, c)
		}
	}
	if len(w.Comings) != 0 {
		t.Fatalf("nobody should have come or gone here, and %d did: %+v", len(w.Comings), w.Comings)
	}
}

// The player has to be able to read this after the fact: it belongs on the
// result of the command that the time passed during, beside what they did.
func TestWhatHappenedInTheRoomIsReportedWithTheAction(t *testing.T) {
	w := New(4)
	w.Player.Location = "bar"
	n := w.NPC("mara")
	n.Role = "Runs Bluebird Laundry"
	n.Location = "bar"
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "courier", Target: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if next.LastResult == nil {
		t.Fatal("nothing was reported at all")
	}
	// Whether anybody moved during this particular command depends on where the
	// clock lands, so this asserts the channel exists rather than the traffic.
	for _, c := range next.LastResult.Comings {
		if c.Note == "" || c.Where == "" {
			t.Fatalf("a coming and going was reported with nothing in it: %+v", c)
		}
	}
}
