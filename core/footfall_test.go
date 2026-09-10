package core

import "testing"

// A room with people in it is a room doing business.
//
// Measured before this was written, over twenty days of one city: at ten in the
// morning the market holds twelve people and the bar seven; at nine at night
// the bar, the club and the casino hold twenty-five each and the docks and the
// butcher are empty. The city has a real rhythm, and none of it reached the
// ledger — a butcher earned the same at three in the morning as at noon, and a
// bar earned the same at noon as when it was full.

func keeper(t *testing.T) *World {
	t.Helper()
	w := New(43)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 5000, 100
	own(w, "laundry")
	return w
}

func TestAFullRoomEarnsMoreThanAnEmptyOne(t *testing.T) {
	w := keeper(t)
	clearRoom(w, "laundry")
	empty := w.Footfall("laundry")
	if empty != 0 {
		t.Fatalf("an emptied room holds %d people", empty)
	}
	quiet := w.RoomTrade("laundry")
	fill(w, "laundry", 20)
	busy := w.RoomTrade("laundry")
	if busy <= quiet {
		t.Fatalf("twenty people through the door were worth %.2f against %.2f empty", busy, quiet)
	}
	// It modulates rather than decides: a full room is worth more, not
	// everything, and an empty one still has a door on it.
	if quiet <= 0 || busy > 2 {
		t.Fatalf("the factor runs %.2f..%.2f, which is not a modulation", quiet, busy)
	}
}

func TestYourOwnManOnTheDoorIsNotACustomer(t *testing.T) {
	w := keeper(t)
	clearRoom(w, "laundry")
	// Somebody of yours, standing in your own laundry.
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || IsOfficial(n.ID) {
			continue
		}
		n.Faction, n.Rank, n.Location = w.PlayerOrganizationID(), RankSoldier, "laundry"
		break
	}
	if w.Footfall("laundry") != 0 {
		t.Fatalf("your own man counted as %d customers", w.Footfall("laundry"))
	}
}

func TestTheRoomReachesTheLedger(t *testing.T) {
	w := keeper(t)
	clearRoom(w, "laundry")
	quiet := w.Trading("laundry")
	fill(w, "laundry", 20)
	if w.Trading("laundry") <= quiet {
		t.Fatalf("the room filled up and the books say %.2f against %.2f", w.Trading("laundry"), quiet)
	}
}
