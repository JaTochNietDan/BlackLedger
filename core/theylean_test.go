package core

import "testing"

// Every actor in the city runs on the same rules. The player can put the wind
// up somebody else's counter; a family that has fallen out with them should be
// able to do the same, and until now the only thing they could do to a business
// was break it.

func leanedOn(t *testing.T) (*World, string, *Faction) {
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
	return w, "laundry", rival
}

func TestAFamilyYouHaveFallenOutWithComesForYourCounter(t *testing.T) {
	t.Parallel()
	w, id, rival := leanedOn(t)
	rival.Goodwill = -80
	prop := w.Properties[id]
	filled := prop.Staff
	if filled == 0 {
		t.Fatal("nobody works at the laundry")
	}
	for day := 0; day < 60 && prop.Staff == filled; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if prop.Staff == filled {
		t.Fatalf("a family that thinks %d of you never touched the %d people behind your counter in two months",
			rival.Goodwill, filled)
	}
	if len(prop.Hands) != prop.Staff {
		t.Fatalf("somebody was put off and the count says %d of %d", prop.Staff, len(prop.Hands))
	}
}

// And a family with no quarrel leaves them alone, which is the ordinary case
// and has to be the safe one.
func TestAFamilyWithNoQuarrelLeavesYourPeopleAlone(t *testing.T) {
	t.Parallel()
	w, id, _ := leanedOn(t)
	prop := w.Properties[id]
	filled := prop.Staff
	for day := 0; day < 60; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if prop.Staff < filled {
		t.Fatalf("a city with no quarrel with the player took %d of their people", filled-prop.Staff)
	}
}
