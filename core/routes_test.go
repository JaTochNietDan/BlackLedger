package core

import "testing"

// Layer 6 of docs/LIVING_WORLD.md asks for "sources, routes, and buyers", and
// what existed was one price for the whole city. Measured last slice: a policy
// that trades made 152 purchases and 26 sales over sixty campaigns and a
// quarter of the investor's money, because buying at the docks and selling at
// the market was the same price. There was no route, only a wait.
//
// A route is two addresses where the same crate is worth different money, and
// the risk is the walk between them — which is the one stretch of this city
// where a hit can catch you cold, and the reason anybody plates a car.

func TestTheSameCrateIsWorthDifferentMoneyInDifferentPlaces(t *testing.T) {
	w := New(79)
	w.Event, w.District = nil, 9
	found := false
	for _, g := range w.Goods {
		docks, market := w.PriceAt("docks", g.ID), w.PriceAt("market", g.ID)
		if docks <= 0 || market <= 0 {
			continue
		}
		if docks != market {
			found = true
		}
	}
	if !found {
		t.Fatal("every good is worth the same on the waterfront as it is on the exchange floor")
	}
	// The waterfront is where it comes ashore, so it is cheaper there.
	if w.PriceAt("docks", "moonshine") >= w.PriceAt("market", "moonshine") {
		t.Fatalf("moonshine costs $%d at the docks and $%d at the market",
			w.PriceAt("docks", "moonshine"), w.PriceAt("market", "moonshine"))
	}
}

func TestARouteIsWorthWalking(t *testing.T) {
	w := New(79)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 5000, 100
	w.Player.Location = "docks"
	cash := w.Player.Cash
	if err := w.Buy("moonshine", 4); err != nil {
		t.Fatal(err)
	}
	spent := cash - w.Player.Cash
	w.Player.Location = "market"
	before := w.Player.Cash
	if err := w.Sell("moonshine", 4); err != nil {
		t.Fatal(err)
	}
	got := w.Player.Cash - before
	if got <= spent {
		t.Fatalf("carrying four crates across the city turned $%d into $%d", spent, got)
	}
}

func TestTheCardQuotesThePriceWhereYouAreStanding(t *testing.T) {
	w := New(79)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 5000, 100
	for _, where := range []string{"docks", "market"} {
		w.Player.Location = where
		a := actionByID(w.Actions(where), "buy:moonshine")
		if a == nil {
			t.Fatalf("no moonshine for sale at %s", where)
		}
		if a.Sum == nil {
			t.Fatal("nowhere to say how many")
		}
		// The field's ceiling has to be worked out at this floor's price, not
		// at some average of the city: a card that offers more than the rule
		// will take is a card that refuses you after you have pressed it.
		if err := w.Buy("moonshine", a.Sum.Most); err != nil {
			t.Fatalf("at %s the card offered %d and the rule refused it: %v", where, a.Sum.Most, err)
		}
		w.Player.Stock["moonshine"] = 0
		w.Player.Cash = 5000
	}
}
