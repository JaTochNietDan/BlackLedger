package core

import (
	"strings"
	"testing"
)

// A family can now put one of the people behind your counter off coming in, and
// the answer to that is the oldest one there is: somebody of yours standing at
// the door. Putting a man on it did nothing about this at all — it was worth
// something against a raid and nothing against somebody walking in and having a
// quiet word.

func minded(t *testing.T, seed uint32) (*World, string, *Faction) {
	t.Helper()
	w := New(seed)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	rival := &w.Factions[0]
	if rival.ID == w.PlayerOrganizationID() {
		rival = &w.Factions[1]
	}
	rival.Goodwill = -80
	return w, "laundry", rival
}

// somebodyOnTheDoor puts one of the player's own at an address, the way the
// posting action does. A player who has taken nobody on has nobody to post, so
// somebody is signed on first — which is the real precondition and not a reason
// to skip the test.
func somebodyOnTheDoor(w *World, id string) *NPC {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || IsOfficial(n.ID) || n.Faction != "" {
			continue
		}
		n.Faction, n.Location, n.Heading = w.PlayerOrganizationID(), id, ""
		w.Properties[id].Posted = n.ID
		return n
	}
	return nil
}

// Somebody of your own standing in the room is the answer to a quiet word.
//
// This used to run one city for ninety days and ask that nobody at all be put
// off a counter with a door on it. That is not the mechanic and never was: what
// minding a room does is make leaning on the people in it refuse, which is a
// gate rather than a rate. The ninety-day claim held on the one campaign number
// it was written with and, once those numbers stopped opening the stream in the
// same eighth of its range, lost two people — and compared against a counter
// nobody was watching at all, over twelve cities, a minded one lost seven and an
// unwatched one six. The door was never protecting anybody from the city; it
// stops somebody walking in and having a word.
func TestSomebodyOnTheDoorStopsTheQuietWord(t *testing.T) {
	t.Parallel()
	w, id, rival := minded(t, spread(1))
	w.Properties[id].Owner = rival.ID
	w.Player.Location, w.Player.Health = id, 100
	if len(w.Properties[id].Hands) == 0 {
		t.Skip("there is nobody behind that counter to put off")
	}
	if reason := w.FrightenReadiness(id); reason != "" {
		t.Fatalf("nobody is minding the room and a quiet word is refused: %s", reason)
	}

	// And with one of theirs standing in it.
	watcher := ""
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || IsOfficial(n.ID) || n.Faction != "" {
			continue
		}
		n.Faction, n.Location, n.Heading = rival.ID, id, ""
		watcher = n.Name
		break
	}
	if watcher == "" {
		t.Skip("nobody in this city is free to stand in a doorway")
	}
	reason := w.FrightenReadiness(id)
	if reason == "" {
		t.Fatalf("%s is standing in the room and a quiet word goes ahead anyway", watcher)
	}
	if !strings.Contains(reason, watcher) {
		t.Fatalf("the refusal does not say who is standing there: %q", reason)
	}
}

func TestYouCannotLeanOnAMindedCounter(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health, w.Player.Respect = 5000, 100, 30
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	w.Properties["butcher"].Owner = rival
	w.EmptyChairs()
	w.Player.Location = "butcher"
	if reason := w.FrightenReadiness("butcher"); reason != "" {
		t.Fatalf("a counter with nobody minding it was already refused: %s", reason)
	}
	// One of theirs standing in it, which is what a family minding its own
	// holding already looks like in this city.
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Faction == rival {
			n.Location, n.Heading = "butcher", ""
			break
		}
	}
	if w.FrightenReadiness("butcher") == "" {
		t.Fatal("somebody of theirs is standing at the door and you can still have a quiet word with the staff")
	}
}
