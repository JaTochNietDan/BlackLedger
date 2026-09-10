package core

import "testing"

// You can walk somebody off a rival's counter. Nothing walks anybody off yours,
// and nobody behind one ever decides they have had enough — so employment was a
// thing that only ever happened in one direction. A business the city can take
// people out of is a business worth defending.

func staffedUp(t *testing.T) (*World, string, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	prop := w.Properties["laundry"]
	if len(prop.Hands) == 0 {
		t.Fatal("nobody works at the laundry")
	}
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	return w, "laundry", rival
}

func TestSomebodyCarryingSomethingAgainstYouWalksOut(t *testing.T) {
	t.Parallel()
	w, id, _ := staffedUp(t)
	prop := w.Properties[id]
	filled, who := prop.Staff, prop.Hands[0]
	// Something you did, rather than a low opinion: everybody in this city
	// starts at nothing and thinks nothing of a stranger, so a trust score on
	// its own is every employee in the game.
	w.Aggrieve(who, SoreAtCards, "the night you took them at cards")
	// Days of it, because somebody who has had enough leaves on a day of their
	// own choosing rather than the moment a number crosses a line.
	for day := 0; day < 30 && len(prop.Hands) == filled; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if len(prop.Hands) == filled {
		t.Fatalf("%s thinks nothing of you and stood behind your counter for a month", w.NPC(who).Name)
	}
	if prop.Staff != len(prop.Hands) {
		t.Fatalf("somebody left and the count says %d of %d", prop.Staff, len(prop.Hands))
	}
	if w.EmployerOf(who) == id {
		t.Fatalf("%s walked out and is still on the books", w.NPC(who).Name)
	}
	// And the gap stays open for a couple of days. Without this the city
	// handed the position straight back the same morning, so losing somebody
	// cost nothing at all.
	short := len(prop.Hands)
	if prop.Shorthanded <= w.Minute {
		t.Fatal("somebody walked out and the counter is ready to hire the same morning")
	}
	w.Event = nil
	w.Advance(1440)
	w.Event = nil
	if len(prop.Hands) > short {
		t.Fatalf("the counter filled the gap the next day: %d against %d", len(prop.Hands), short)
	}
}

func TestSomebodyWithNothingAgainstYouStays(t *testing.T) {
	t.Parallel()
	w, id, _ := staffedUp(t)
	prop := w.Properties[id]
	filled := prop.Staff
	for day := 0; day < 30; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if len(prop.Hands) < filled {
		t.Fatalf("people with nothing against you left anyway: %d of %d", len(prop.Hands), filled)
	}
}

// And a rival takes them, which is the same move the player has and the reason
// to keep your own people happy.
func TestARivalCanTakeSomebodyOffYourCounter(t *testing.T) {
	t.Parallel()
	w, id, rival := staffedUp(t)
	prop := w.Properties[id]
	// A rival with a position going and the money to fill it.
	w.Properties["butcher"].Owner = rival
	w.Properties["butcher"].Staff, w.Properties["butcher"].Hands = 0, nil
	if f := w.faction(rival); f != nil {
		f.Cash = 20000
	}
	for _, who := range prop.Hands {
		w.Aggrieve(who, SoreAtCards, "the night you took them at cards")
	}
	filled, taken := prop.Staff, false
	for day := 0; day < 60 && !taken; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
		for _, who := range w.Properties["butcher"].Hands {
			for _, was := range prop.Hands {
				taken = taken || who == was
			}
		}
		if len(w.Properties["butcher"].Hands) > 0 {
			taken = true
		}
	}
	if !taken {
		t.Fatalf("a rival with an empty counter and money never took one of the %d you had", filled)
	}
}
