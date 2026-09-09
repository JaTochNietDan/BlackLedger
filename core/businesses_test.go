package core

import "testing"

// Four businesses were added to a city that had three. The balance simulation
// barely moved, and that is not evidence they work: those campaigns run seven
// to thirteen days and never get rich enough to buy a $980 haulage yard. So
// every trade is exercised directly instead, and the test is written over the
// whole table rather than over the four, so a fifth cannot be added without it.

func TestEveryTradeInTheCityCanActuallyBeRun(t *testing.T) {
	if len(trades) < 7 {
		t.Fatalf("the city has %d trades, so this is testing less than it was written for", len(trades))
	}
	for id := range trades {
		place, ok := PlaceByID(id)
		if !ok {
			t.Errorf("%s has a trade and is not a place in the city", id)
			continue
		}
		trade, _ := TradeOf(id)
		w := New(97)
		w.District = 2
		w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 200, 100
		prop := w.Properties[id]
		if prop == nil {
			t.Errorf("%s has no premises", id)
			continue
		}
		if prop.Income <= 0 {
			t.Errorf("%s (%s) is a business that earns nothing, so nobody can be starved of it", id, place.Name)
		}
		if place.Cost <= 0 {
			t.Errorf("%s (%s) is a business nobody can buy", id, place.Name)
		}
		// A trading business is already running before anybody buys it.
		if prop.Staff != trade.Hands || prop.Supply != trade.RestockAmount {
			t.Errorf("%s opens with %d of %d hands and %d of %d supplies",
				id, prop.Staff, trade.Hands, prop.Supply, trade.RestockAmount)
		}
		if w.Capacity(id) <= 0 {
			t.Errorf("%s works at nothing on the day it opens", id)
		}

		// Run it down and see it suffer, then put it right.
		prop.Owner = "player:1"
		prop.Supply, prop.Staff = 0, 0
		empty := w.Capacity(id)
		if empty >= 1 {
			t.Errorf("%s with nobody in it and nothing to work with runs at %.2f", id, empty)
		}
		prop.Trouble = true
		troubled := w.Capacity(id)
		if troubled > empty {
			t.Errorf("%s earns more in trouble (%.2f) than out of it (%.2f)", id, troubled, empty)
		}
		if trade.Trouble == "" || trade.Remedy == "" || trade.RemedyCost <= 0 {
			t.Errorf("%s has no trouble of its own to deal with", id)
		}
	}
}

// The point of adding them is that they are not the same business four times.
// Every trade has its own trouble and its own words for what it runs on.
func TestNoTwoTradesHaveTheSameTrouble(t *testing.T) {
	troubles, supplies := map[string]string{}, map[string]string{}
	for id, trade := range trades {
		if other, same := troubles[trade.Trouble]; same {
			t.Errorf("%s and %s go wrong in exactly the same way", other, id)
		}
		troubles[trade.Trouble] = id
		if other, same := supplies[trade.Supplies]; same {
			t.Errorf("%s and %s run on exactly the same thing", other, id)
		}
		supplies[trade.Supplies] = id
	}
}
