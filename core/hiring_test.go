package core

import "testing"

// A business the player buys comes with the people already working in it —
// they were behind that counter yesterday and nobody has told them to go. What
// must not happen is the city quietly filling a position that empties: the
// counter refilled itself two days after somebody was frightened off, for
// nothing, so losing a pair of hands cost the player neither money nor a
// decision, and `hire` was an action no policy ever needed in two hundred
// commands.

func hiring(t *testing.T) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	w.Player.Location = "laundry"
	return w, "laundry"
}

func TestABusinessComesWithItsPeople(t *testing.T) {
	t.Parallel()
	// Bought through the same path a player uses, and looked at straight away:
	// it said three positions filled and had nobody standing in it until the
	// next business day, so there was nobody to ask, nobody to put in charge,
	// and a wage bill for nobody.
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health, w.Player.Respect = 200000, 100, 60
	w.Player.Location = "laundry"
	if err := w.apply(Command{Kind: "acquire", Target: "laundry", RequestID: "buyalaundrynow1"}); err != nil {
		t.Fatal(err)
	}
	id := "laundry"
	prop := w.Properties[id]
	trade, _ := TradeOf(id)
	if prop.Staff != trade.Hands || len(prop.Hands) != prop.Staff {
		t.Fatalf("a laundry just taken over has %d of %d positions filled and %d names",
			prop.Staff, trade.Hands, len(prop.Hands))
	}
}

func TestAPositionThatEmptiesStaysEmptyUntilYouFillIt(t *testing.T) {
	t.Parallel()
	w, id := hiring(t)
	prop := w.Properties[id]
	filled := prop.Staff
	who := prop.Hands[0]
	w.walkOut(who, id, "")
	for day := 0; day < 20; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if prop.Staff == filled {
		t.Fatalf("the city filled the gap for nothing: %d of %d", prop.Staff, filled)
	}
	// And filling it is the player's decision, and it costs.
	cash := w.Player.Cash
	if reason := w.HireReadiness(id); reason != "" {
		t.Fatalf("a short-handed business of yours cannot hire: %s", reason)
	}
	if err := w.Hire(id); err != nil {
		t.Fatal(err)
	}
	if prop.Staff != filled {
		t.Fatalf("somebody was hired and the counter says %d of %d", prop.Staff, filled)
	}
	if w.Player.Cash >= cash {
		t.Fatal("hiring somebody cost nothing")
	}
}

// A rival's counter still fills itself: their family has people and does not
// need the player to think about it.
func TestARivalsCounterStillFillsItself(t *testing.T) {
	t.Parallel()
	w, _ := hiring(t)
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	w.Properties["butcher"].Owner = rival
	w.Properties["butcher"].Staff, w.Properties["butcher"].Hands = 0, nil
	w.EmptyChairs()
	if len(w.Properties["butcher"].Hands) == 0 {
		t.Fatal("a family's own butcher has nobody in it and nobody is coming")
	}
}
