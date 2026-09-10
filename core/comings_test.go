package core

import "testing"

// People move between buildings now, and a player standing in a room never
// noticed. Somebody they had been talking to walked out of the door and was
// simply not in the list any more; somebody else appeared in it with no
// explanation. A room where people arrive and leave without remark is a list
// that changes behind your back rather than a place you are standing in.

func TestYouSeeSomebodyLeaveTheRoomYouAreStandingIn(t *testing.T) {
	t.Parallel()
	w := New(4)
	w.Player.Location = "bar"
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	w.SetOut()
	if len(w.Comings) != 0 {
		t.Fatal("she was noticed leaving before she had left")
	}
	// The room notices her at the moment she walks out of it, not at the
	// moment she decides to.
	walksOut(w, n)
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
	t.Parallel()
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
	t.Parallel()
	w := New(4)
	// The Bellwether Herald: an address with nobody in it who has anywhere to
	// be, so anything reported here came from another room.
	w.Player.Location = "herald"
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	w.See(n)
	w.SetOut()
	w.Minute = n.Arrives
	w.Arrivals()
	// And seen again where she arrived, because the question here is what the
	// list says about somebody standing still rather than how she was found.
	w.See(n)
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
	t.Parallel()
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

// The People screen reads a person's address straight off their record, and a
// person out walking keeps the address they set off from until they arrive. So
// the screen said Mara Bell was at Saint Agnes while she was somewhere on the
// way to the laundry, and a player could spend forty minutes crossing the city
// to a room she was not in. An interface that is confidently wrong about where
// somebody is, is worse than one that says nothing.
func TestSomebodyOutWalkingIsNotReportedAsBeingSomewhere(t *testing.T) {
	t.Parallel()
	w := New(4)
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	// Seen before she left, or the list has nothing to say about her at all:
	// what it reports is the player's knowledge now rather than the city's own
	// books.
	w.See(n)
	w.SetOut()
	walksOut(w, n)

	var her Presence
	for _, p := range w.Everyone() {
		if p.ID == n.ID {
			her = p
		}
	}
	if her.ID == "" {
		t.Fatal("she is not in the city at all")
	}
	// The city list is the player's knowledge now, not the city's own books, so
	// somebody they have never laid eyes on has no address at all. She was seen
	// at Saint Agnes before she walked out of it.
	if !her.Walking {
		t.Fatal("she is on the street and the city says she is standing somewhere")
	}
	if her.Where == "Saint Agnes" {
		t.Fatal("she is reported at the address she walked out of")
	}
	// Where the player is pointed has to be somewhere they could meet her, and
	// the only such place is where she is going.
	if her.WhereID != "laundry" {
		t.Fatalf("the player is pointed at %q, where she is not and is not going", her.WhereID)
	}
	if her.Minutes <= 0 || her.Doing == "" {
		t.Fatalf("she is %d minutes out, doing %q", her.Minutes, her.Doing)
	}
	for _, word := range []string{"Bluebird Laundry"} {
		if !contains(her.Doing, word) {
			t.Fatalf("what she is doing does not say where she is going: %q", her.Doing)
		}
	}
}

// And once she gets there the report goes back to being ordinary.
func TestArrivingPutsSomebodyBackInARoom(t *testing.T) {
	t.Parallel()
	w := New(4)
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	w.SetOut()
	w.Minute = n.Arrives
	w.Arrivals()
	// Seen where she arrived: the list reports the player's knowledge now, and
	// the question here is what it says about somebody standing still rather
	// than how she came to be found.
	w.See(n)
	for _, p := range w.Everyone() {
		if p.ID == n.ID {
			if p.Walking || p.Where != "Bluebird Laundry" {
				t.Fatalf("she has arrived and is reported as %+v", p)
			}
			return
		}
	}
	t.Fatal("she vanished on arrival")
}

// The other half of the same lie: an action offered against somebody who has
// walked out. Buying a coffee for a man who is halfway across the city is not
// something the player should be able to press.
func TestYouCannotDealWithSomebodyWhoHasWalkedOut(t *testing.T) {
	t.Parallel()
	w := New(4)
	w.Player.Location = "bar"
	n := w.NPC("mara")
	before := 0
	for _, a := range w.Actions("bar") {
		if a.Subject == n.ID && !a.Disabled {
			before++
		}
	}
	if before == 0 {
		t.Fatal("there was nothing to do with her in the first place")
	}
	// Send her somewhere; she is now on the street, not in the bar.
	n.Role = "Runs Bluebird Laundry"
	w.SetOut()
	walksOut(w, n)
	if !w.Travelling(n) {
		t.Fatal("she did not set off")
	}
	for _, a := range w.Actions("bar") {
		if a.Subject == n.ID && !a.Disabled {
			t.Fatalf("%q is offered against somebody who is out on the street", a.Label)
		}
	}
}
