package core

import (
	"strings"
	"testing"
)

// "It probably would make sense that we don't always know everyone's location
// and that's a service that we have to pay for to gain that information from...
// This would also act as a way of making it harder to make attempts on people's
// lives in the game as it would be a drawn out task more than anything."

func watcher(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 5000, 100
	w.Player.Location = "bar"
	var away *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead && n.Location != "bar" && n.Faction != w.PlayerOrganizationID() && !IsOfficial(n.ID) {
			away = n
			break
		}
	}
	if away == nil {
		t.Fatal("everybody in this city is in the bar")
	}
	return w, away
}

func TestYouKnowWhereSomebodyIsBecauseYouSawThem(t *testing.T) {
	t.Parallel()
	w, away := watcher(t)
	if w.KnowsWhere(away.ID) {
		t.Fatalf("%s is across the city and the player knows where without ever having seen them", away.Name)
	}
	if note := w.WhereNote(away.ID); !strings.Contains(note, "Nobody has told you") {
		t.Fatalf("somebody never seen reads as %q", note)
	}
	// Walk into the room they are standing in.
	w.Player.Location = away.Location
	w.SeeTheRoom()
	if !w.KnowsWhere(away.ID) {
		t.Fatalf("the player is standing in front of %s and does not know where they are", away.Name)
	}
	// Walk away. What you know is what you saw, and it is still true while they
	// have not moved.
	w.Player.Location = "bar"
	if !w.KnowsWhere(away.ID) {
		t.Fatal("what the player saw stopped being true the moment they left the room")
	}
	// They move, and it stops being true.
	away.Location = "docks"
	if w.KnowsWhere(away.ID) {
		t.Fatalf("%s moved and the player followed them without looking", away.Name)
	}
	if note := w.WhereNote(away.ID); !strings.Contains(note, "Last seen") {
		t.Fatalf("a stale sighting reads as %q", note)
	}
}

func TestASightingGoesStale(t *testing.T) {
	t.Parallel()
	w, away := watcher(t)
	w.Player.Location = away.Location
	w.SeeTheRoom()
	// Away from them, or standing in front of somebody is knowing where they
	// are however long ago you first saw them.
	w.Player.Location = "bar"
	w.Minute += SightingLasts + 60
	if w.KnowsWhere(away.ID) {
		t.Fatalf("a day later the player still knows where %s was", away.Name)
	}
}

// Your own people and the people you employ are findable, because they work for
// you and you can send for them.
func TestYourOwnPeopleAreAlwaysFindable(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	hand := w.Properties["laundry"].Hands[0]
	w.Player.Location = "bar"
	if !w.KnowsWhere(hand) {
		t.Fatalf("%s works behind your own counter and you cannot find them", w.NPC(hand).Name)
	}
}

// The screen that lists the city says what the player knows rather than where
// everybody is, which is where the difficulty actually lives: going after
// somebody already needs them standing in front of you, so what was easy was
// reading an address off a list and walking to it.
func TestTheCityListSaysWhatYouKnowRatherThanWhereEverybodyIs(t *testing.T) {
	t.Parallel()
	w, away := watcher(t)
	find := func() Presence {
		for _, p := range w.Everyone() {
			if p.ID == away.ID {
				return p
			}
		}
		t.Fatalf("%s is not on the list at all", away.Name)
		return Presence{}
	}
	unseen := find()
	if unseen.WhereID != "" || unseen.Where != "" {
		t.Fatalf("%s has never been seen and the list gives their address as %q", away.Name, unseen.Where)
	}
	if unseen.Lost == "" {
		t.Fatalf("%s cannot be found and the list says nothing about it", away.Name)
	}
	w.Player.Location = away.Location
	w.SeeTheRoom()
	seen := find()
	if seen.WhereID != away.Location {
		t.Fatalf("the player is standing in front of %s and the list says %q", away.Name, seen.WhereID)
	}
	if seen.Lost != "" {
		t.Fatalf("somebody in the same room reads as lost: %q", seen.Lost)
	}
	// They move, and the address goes with them.
	away.Location = "docks"
	if moved := find(); moved.WhereID != "" {
		t.Fatalf("%s moved and the list followed them to %q", away.Name, moved.WhereID)
	} else if !strings.Contains(moved.Lost, "Last seen") {
		t.Fatalf("a stale sighting reads as %q", moved.Lost)
	}
}
