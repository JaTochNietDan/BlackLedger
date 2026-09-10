package core

import "testing"

// A business was a set of books that all read the same. Whatever the player
// owned, it absorbed exactly the same amount of police attention — fourteen
// points scaled by condition — and the gate asked what KIND OF ROOM it was
// rather than what trade was run in it, so a casino could not launder a dollar.
// A casino is the best cash front there is; that is most of why anybody owns
// one.
//
// What a business is worth as a front is a fact about the trade.

func TestWhatABusinessHidesDependsOnWhatTradeItIs(t *testing.T) {
	t.Parallel()
	seen := map[int]string{}
	for kind, trade := range trades {
		if trade.Cover <= 0 {
			t.Errorf("a %s can hide nothing at all, so owning one is only ever a wage bill", kind)
		}
		seen[trade.Cover] = kind
	}
	if len(seen) < 4 {
		t.Errorf("nine kinds of business between them offer %d different amounts of cover", len(seen))
	}
	// The classic front hides more than the yard full of trucks.
	laundry, _ := TradeOfKind("laundry")
	haulage, _ := TradeOfKind("haulage")
	casino, _ := TradeOfKind("casino")
	if laundry.Cover <= haulage.Cover {
		t.Errorf("a laundry hides %d and a haulage yard %d", laundry.Cover, haulage.Cover)
	}
	if casino.Cover <= haulage.Cover {
		t.Errorf("a casino hides %d and a haulage yard %d", casino.Cover, haulage.Cover)
	}
}

// And the room's own books have to be the ones doing it. Two businesses of
// different kinds, in the same condition, must not absorb the same attention.
func TestTwoDifferentBusinessesDoNotLaunderAlike(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Heat = 40000, 200, 60
	for _, id := range []string{"laundry", "haulage"} {
		w.Properties[id].Owner = "player:1"
		w.Properties[id].Condition = 100
	}
	soft, hard := w.launderCapacity("laundry"), w.launderCapacity("haulage")
	if soft <= hard {
		t.Errorf("a laundry absorbs %d and a haulage yard %d, in the same condition", soft, hard)
	}
	t.Logf("in perfect condition: laundry %d, haulage %d", soft, hard)
}

// A casino is a cash business and could not launder at all, because the gate
// asked what kind of ROOM it was. It asks what trade is run there now.
func TestACasinoCanPutMoneyThroughItsOwnBooks(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Heat = 40000, 200, 40
	w.Properties["casino"].Owner = "player:1"
	w.Properties["casino"].Condition = 100
	if reason := w.LaunderReadiness("casino"); reason != "" {
		t.Errorf("a casino of the player's own will not take their money: %q", reason)
	}
	// And a place that runs no business at all still cannot.
	if reason := w.LaunderReadiness("room"); reason == "" {
		t.Error("the player's rented room launders money")
	}
}

// The wiring, which is where this kind of change goes wrong. The readiness
// function allowing a casino is worth nothing if the button is still handed out
// by what the room looks like: the player would be refused something the rules
// say they may do, and never see why.
func TestTheButtonIsOfferedWhereverTheTradeAllowsIt(t *testing.T) {
	t.Parallel()
	for _, l := range Locations {
		trade, runs := TradeOf(l.ID)
		w := New(53)
		w.District = 2
		w.Player.Cash, w.Player.Respect, w.Player.Heat = 40000, 200, 40
		w.Player.Location = l.ID
		if runs {
			w.Properties[l.ID].Owner = "player:1"
			w.Properties[l.ID].Condition = 100
		}
		offered := false
		for _, a := range w.Actions(l.ID) {
			if a.ID == "launder" {
				offered = true
			}
		}
		wanted := runs && trade.Cover > 0
		if offered != wanted {
			t.Errorf("%s: the books are %s and the trade says they should be %s",
				l.ID, offeredWord(offered), offeredWord(wanted))
		}
	}
}

func offeredWord(b bool) string {
	if b {
		return "offered"
	}
	return "not offered"
}
