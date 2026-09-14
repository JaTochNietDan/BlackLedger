package core

import "testing"

// Four businesses were added to a city that had three. The balance simulation
// barely moved, and that is not evidence they work: those campaigns run seven
// to thirteen days and never get rich enough to buy a $980 haulage yard. So
// every trade is exercised directly instead, and the test is written over the
// whole table rather than over the four, so a fifth cannot be added without it.

func TestEveryTradeInTheCityCanActuallyBeRun(t *testing.T) {
	t.Parallel()
	if len(trades) < 9 {
		t.Fatalf("the city knows %d kinds of business, so this is testing less than it was written for", len(trades))
	}
	// Over the addresses, not over the trade table. A trade belongs to a KIND
	// of business now, and several addresses can share one, so walking the
	// table asks about "cabs" rather than about the cab company on Ordway
	// Street. The property is still that every business in the city can be
	// run — there are simply more of them than there are kinds.
	for _, l := range Locations {
		if l.Kind == "" {
			continue
		}
		id := l.ID
		place := l
		trade, running := TradeOf(id)
		if !running {
			t.Errorf("%s is a %s and has no trade", id, l.Kind)
			continue
		}
		w := New(97)
		w.District = 2
		_ = place
		w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 200, 100
		prop := w.Properties[id]
		if prop == nil {
			t.Errorf("%s has no premises", id)
			continue
		}
		if prop.Income <= 0 {
			t.Errorf("%s (%s) is a business that earns nothing, so nobody can be starved of it", id, place.Name)
		}
		// Held by the player somehow: bought, or taken.
		//
		// This used to say bought, full stop, and that was the same assumption
		// that left the city's four busiest rooms as numbers. A family's seat
		// has no price — "this is not somewhere that changes hands" — because
		// it changes hands by force, and a room you fight a war for should be
		// the best business in the city rather than the only address in it that
		// runs on nothing. So the rule is that every trade can end up yours,
		// and a seat proves it the other way.
		if AcquisitionCost(w, id) <= 0 {
			if w.faction(prop.Owner) == nil {
				t.Errorf("%s (%s) has no price and belongs to no family, so nobody can ever hold it", id, place.Name)
			}
			if !w.SeatOf(prop.Owner, id) {
				t.Errorf("%s (%s) has no price and is not a family's holding either", id, place.Name)
			}
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
	t.Parallel()
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

// "How do I buy businesses, I don't see where I can buy any of them."
//
// Because the whole business block — buying, hiring, restocking, the operating
// modes, repairs, the books — was gated on three addresses by name: the laundry,
// the garage and the casino, which is every business this city had when that
// line was written. Everything added since could be walked into and robbed and
// never bought. The rule underneath was general the whole time; only the button
// was not.
func TestEverySomewhereThatEarnsCanBeTakenOver(t *testing.T) {
	t.Parallel()
	w := New(9)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 90, 100
	missing, offered := []string{}, 0
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || prop.Income <= 0 {
			continue
		}
		w.Player.Location = l.ID
		found := false
		for _, a := range w.Actions(l.ID) {
			if a.ID == "acquire" {
				found, offered = true, offered+1
				if a.Disabled && a.Reason == "There is nothing here to take over" {
					t.Errorf("%s earns $%d an hour and says there is nothing to take over", l.ID, prop.Income)
				}
			}
		}
		if !found && !w.Own(l.ID) {
			missing = append(missing, l.ID)
		}
	}
	t.Logf("%d earning addresses offer a way in", offered)
	if len(missing) > 0 {
		t.Errorf("these earn and cannot be bought anywhere: %v", missing)
	}
}

// And once it is yours, it is a business you can actually run.
func TestABusinessYouHoldCanBeRun(t *testing.T) {
	t.Parallel()
	w := New(9)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 90, 100
	for _, id := range []string{"restaurant", "butcher", "cabstand", "pumps"} {
		if w.Properties[id] == nil {
			continue
		}
		w.Properties[id].Owner = "player:1"
		w.Player.Location = id
		want := map[string]bool{"inspect": false, "restock": false, "operate:hard": false, "repair": false}
		for _, a := range w.Actions(id) {
			if _, ours := want[a.ID]; ours {
				want[a.ID] = true
			}
		}
		for kind, there := range want {
			if !there {
				t.Errorf("%s is yours and offers no %s", id, kind)
			}
		}
	}
}

// A place with no price is not for sale. The bar, the docks, the exchange and
// the Bellandi club all earn and all carry a cost of nothing, because none of
// them was ever meant to be bought — and "Establish protection, $0" on the
// exchange is the interface offering a business for free.
func TestSomewhereWithNoPriceIsNotForSale(t *testing.T) {
	t.Parallel()
	w := New(9)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 90, 100
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || prop.Income <= 0 || AcquisitionCost(w, l.ID) > 0 {
			continue
		}
		// Nobody holds it, so the only thing left that can refuse is the
		// missing price. In a fresh city all four of these are family-held, and
		// "belongs to another organization" would answer for them whether the
		// price rule existed or not — which is a guard that cannot fail.
		prop.Owner = "independent"
		w.Player.Location = l.ID
		for _, a := range w.Actions(l.ID) {
			if a.ID != "acquire" {
				continue
			}
			if !a.Disabled {
				t.Errorf("%s costs nothing and is offered for sale anyway", l.ID)
			}
			if a.Cost == 0 && !a.Disabled {
				t.Errorf("%s is a business for free", l.ID)
			}
		}
		if reason := w.AcquireReadiness(l.ID); reason == "" {
			t.Errorf("%s has no price and the rules allow taking it over", l.ID)
		}
	}
}
