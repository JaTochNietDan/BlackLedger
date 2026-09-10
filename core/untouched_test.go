package core

import "testing"

// The ordinary case has to be the safe one.
//
// Twice in two ticks a new rule punished a player who had done nothing: "below
// twenty trust" was every employee in the game, and then paying exactly the
// going rate was worth half of the floor's temptation. Both times a laundry
// emptied itself over a month with nothing having happened, and both times the
// test that caught it was written for the rule rather than for the city.
//
// This is written for the city. A player who owns two businesses, pays the
// rate, wrongs nobody and does nothing at all should still have them a month
// later. Anything that takes them away has to be something they chose.

func untouched(t *testing.T) *World {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 40000, 100
	own(w, "laundry")
	own(w, "butcher")
	w.EmptyChairs()
	return w
}

func TestAPlayerWhoDoesNothingLosesNothing(t *testing.T) {
	t.Parallel()
	w := untouched(t)
	mine := []string{"laundry", "butcher"}
	staff, held := map[string]int{}, map[string][]string{}
	for _, id := range mine {
		staff[id] = w.Properties[id].Staff
		held[id] = append([]string{}, w.Properties[id].Hands...)
		if staff[id] == 0 {
			t.Fatalf("%s starts with nobody in it", id)
		}
	}
	for day := 0; day < 45; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	for _, id := range mine {
		prop := w.Properties[id]
		if !w.Own(id) {
			t.Fatalf("%s changed hands while the player did nothing", id)
		}
		if prop.Staff < staff[id] {
			t.Errorf("%s went from %d hands to %d with nothing having happened", id, staff[id], prop.Staff)
		}
		if len(prop.Hands) != prop.Staff {
			t.Errorf("%s says %d positions filled and has %d names", id, prop.Staff, len(prop.Hands))
		}
	}
	// And the player is still standing, still solvent, and nobody has taken
	// against them for something they did not do.
	if !w.Player.Alive {
		t.Fatal("a player who did nothing at all was killed by the city")
	}
	sore := 0
	for i := range w.NPCs {
		if w.NPCs[i].Sore > 0 {
			sore++
		}
	}
	if sore > 0 {
		t.Errorf("%d people hold something against a player who has done nothing", sore)
	}
}
