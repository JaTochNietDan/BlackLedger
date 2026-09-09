package core

import "testing"

// An action that pays its own fee has to declare Cost of nothing, or the
// command layer takes the money a second time. The panel prints the price from
// Cost, so those actions showed no price at all: calling two families to a room
// takes $220 and its button said nothing about money.
//
// Asks is that fee for the panel only. Two things have to hold: it has to be
// there, and declaring it must not charge anybody twice.

func priced(t *testing.T, w *World, place, kind string) Action {
	t.Helper()
	w.Player.Location = place
	for _, a := range w.Actions(place) {
		if a.ID == kind {
			return a
		}
	}
	t.Fatalf("%s is not offered at %s", kind, place)
	return Action{}
}

func TestWorkThatPaysItsOwnFeeStillShowsAPrice(t *testing.T) {
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 200, 100
	w.Player.Contacts = 5
	w.Offshore = 5000
	w.Player.Offshore = false
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	w.Incorporate()
	w.ensureOfficials()
	w.Player.Retainers = append(w.Player.Retainers, "editor")
	w.News = append(w.News, Story{ID: ID(), Minute: w.Minute, Life: w.Life,
		Headline: "QUESTIONED AT THE DOCKS", Body: "Police called again.", Kind: "police"})

	for _, c := range []struct{ place, kind string; fee int }{
		{"herald", "spike", SpikeCost},
		{"herald", "puff", PuffCost},
		{"market", "offshore_access", AccessCost},
		{"laundry", "still", StillCost},
	} {
		a := priced(t, w, c.place, c.kind)
		if a.Cost != 0 {
			t.Errorf("%s declares a cost of %d, and the command layer would charge it on top of its own fee", c.kind, a.Cost)
		}
		if a.Asks != c.fee {
			t.Errorf("%s takes $%d and the button offers %d", c.kind, c.fee, a.Asks)
		}
	}
}

// The other half, and the reason Cost has to stay nothing: declaring a price
// must not take the money twice.
func TestDeclaringAPriceDoesNotChargeItTwice(t *testing.T) {
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 200, 100
	w.Properties["laundry"].Owner = "player:1"
	a := priced(t, w, "laundry", "still")
	if a.Disabled {
		t.Skipf("a still cannot be built here: %s", a.Reason)
	}
	before := w.Player.Cash
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "still", Target: "laundry"})
	if err != nil {
		t.Fatalf("building it failed: %v", err)
	}
	spent := before - next.Player.Cash
	// The clock moves, so rent and income land too; the fee must be in there
	// once, not twice.
	// The clock moves while it is built, so a day's income and rent land on top
	// of the fee. What matters is that the fee itself went out once: comfortably
	// more than nothing, and nowhere near twice.
	if spent >= 2*StillCost {
		t.Errorf("a still costs $%d and pressing it took $%d — charged twice", StillCost, spent)
	}
	if spent < StillCost/2 {
		t.Errorf("a still costs $%d and pressing it took only $%d — not charged at all", StillCost, spent)
	}
	t.Logf("declared $%d, took $%d once the hours' other money had moved", a.Asks, spent)
}
