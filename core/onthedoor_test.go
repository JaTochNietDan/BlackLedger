package core

import "testing"

// A family can now put one of the people behind your counter off coming in, and
// the answer to that is the oldest one there is: somebody of yours standing at
// the door. Putting a man on it did nothing about this at all — it was worth
// something against a raid and nothing against somebody walking in and having a
// quiet word.

func minded(t *testing.T) (*World, string, *Faction) {
	t.Helper()
	w := New(61)
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

func TestSomebodyOnTheDoorStopsTheQuietWord(t *testing.T) {
	t.Parallel()
	w, id, _ := minded(t)
	if somebodyOnTheDoor(w, id) == nil {
		t.Skip("the player has nobody of their own to post")
	}
	prop := w.Properties[id]
	filled := prop.Staff
	for day := 0; day < 90; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if prop.Staff < filled {
		t.Fatalf("somebody was on the door and %d of your people were still put off", filled-prop.Staff)
	}
}

// And the same door works the other way: a rival's counter with one of theirs
// on it is not a counter you can walk into and lean on.
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
