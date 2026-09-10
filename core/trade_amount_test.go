package core

import "testing"

// The last place in this game where money moves in somebody else's units.
//
// "It should be as dynamic and user settable as possible." That was said about
// funding a casino and it is true of the one high-variance income path in the
// game: contraband moved in lots of five, so the decision the whole system is
// built around — how much do you dare carry — was made by a constant.

func runner(t *testing.T) (*World, string) {
	t.Helper()
	w := New(73)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	for _, g := range w.Goods {
		for _, l := range Locations {
			if TradesAt(l.ID, g.ID) {
				w.Player.Location = l.ID
				return w, g.ID
			}
		}
	}
	t.Fatal("nowhere in this city trades anything")
	return nil, ""
}

func TestYouBuyTheNumberYouName(t *testing.T) {
	w, good := runner(t)
	g := w.Good(good)
	cash := w.Player.Cash
	if err := w.Buy(good, 7); err != nil {
		t.Fatal(err)
	}
	if w.Holding(good) != 7 {
		t.Fatalf("bought 7 and hold %d", w.Holding(good))
	}
	if w.Player.Cash != cash-g.Price*7 {
		t.Fatalf("7 at $%d cost $%d", g.Price, cash-w.Player.Cash)
	}
	// Nothing named is still the lot, so every existing caller is unchanged.
	if err := w.Buy(good, 0); err != nil {
		t.Fatal(err)
	}
	if w.Holding(good) != 7+Lot {
		t.Fatalf("an unnamed amount bought %d", w.Holding(good)-7)
	}
}

func TestYouSellTheNumberYouName(t *testing.T) {
	w, good := runner(t)
	if err := w.Buy(good, 12); err != nil {
		t.Fatal(err)
	}
	g := w.Good(good)
	cash := w.Player.Cash
	if err := w.Sell(good, 5); err != nil {
		t.Fatal(err)
	}
	if w.Holding(good) != 7 {
		t.Fatalf("sold 5 of 12 and hold %d", w.Holding(good))
	}
	if w.Player.Cash != cash+g.Price*5 {
		t.Fatalf("5 at $%d fetched $%d", g.Price, w.Player.Cash-cash)
	}
	// Nothing named still sells the lot of it, which is what selling meant.
	if err := w.Sell(good, 0); err != nil {
		t.Fatal(err)
	}
	if w.Holding(good) != 0 {
		t.Fatalf("an unnamed sale left %d", w.Holding(good))
	}
}

func TestTheTradeRefusesWhatYouCannotDo(t *testing.T) {
	w, good := runner(t)
	g := w.Good(good)
	w.Player.Cash = g.Price * 3
	if w.TradeReadiness(good, "buy", 4) == "" {
		t.Fatal("bought four with the price of three")
	}
	if reason := w.TradeReadiness(good, "buy", 3); reason != "" {
		t.Fatalf("refused three with the price of three: %s", reason)
	}
	w.Player.Cash = 20000
	if err := w.Buy(good, 2); err != nil {
		t.Fatal(err)
	}
	if w.TradeReadiness(good, "sell", 3) == "" {
		t.Fatal("sold three while carrying two")
	}
}

func TestTheCardSaysWhatYouCanCarry(t *testing.T) {
	w, good := runner(t)
	buy := actionByID(w.Actions(w.Player.Location), "buy:"+good)
	if buy == nil {
		t.Fatal("nothing for sale where something is traded")
	}
	if buy.Sum == nil {
		t.Fatal("there is nowhere to say how much")
	}
	// What a person can carry is a real limit and the field has to know it.
	if buy.Sum.Most > w.CarryLimit() {
		t.Fatalf("the field offers %d units and a person can carry %d", buy.Sum.Most, w.CarryLimit())
	}
	if err := w.Buy(good, 3); err != nil {
		t.Fatal(err)
	}
	sell := actionByID(w.Actions(w.Player.Location), "sell:"+good)
	if sell == nil || sell.Sum == nil {
		t.Fatal("no way to name how much to sell")
	}
	if sell.Sum.Most != 3 {
		t.Fatalf("carrying 3 and the field offers %d", sell.Sum.Most)
	}
}
