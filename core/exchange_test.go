package core

import (
	"strings"
	"testing"
)

// Knowing what a price means.
//
// A trader can see what a good costs here and what the other floor pays. That
// says which way to walk and nothing else — whether this is a week to buy or a
// week to sit still is not on any card, because what a thing is normally worth
// lives in the core and has never been shown to anybody. A floor of your own is
// where that is common knowledge.

const theFloor = "market"

func trading(t *testing.T, hold bool) *World {
	t.Helper()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 90
	if hold {
		own(w, theFloor)
	}
	w.Player.Location = "docks"
	w.Event = nil
	return w
}

func TestAFloorOfYourOwnSaysWhatAPriceMeans(t *testing.T) {
	t.Parallel()
	theirs, mine := trading(t, false), trading(t, true)
	for _, g := range theirs.Goods {
		if word := theirs.MarketWord(g.ID); word != "" {
			t.Fatalf("somebody else's floor told you about %s: %q", g.ID, word)
		}
		if mine.MarketWord(g.ID) == "" {
			t.Fatalf("your own floor says nothing about %s", g.ID)
		}
	}
}

func TestTheFloorSaysWhichWayThePriceIsLeaning(t *testing.T) {
	t.Parallel()
	w := trading(t, true)
	g := w.Good("moonshine")
	high, low, even := "", "", ""
	g.Price = g.Base * 2
	high = w.MarketWord("moonshine")
	g.Price = g.Base / 2
	low = w.MarketWord("moonshine")
	g.Price = g.Base
	even = w.MarketWord("moonshine")
	if high == low || high == even || low == even {
		t.Fatalf("dear, cheap and ordinary all read the same:\n  %s\n  %s\n  %s", high, low, even)
	}
	if !strings.Contains(high, "over") {
		t.Fatalf("at twice what it is worth the floor says: %s", high)
	}
	if !strings.Contains(low, "under") {
		t.Fatalf("at half what it is worth the floor says: %s", low)
	}
	// And the figure is the real one rather than a word.
	if !strings.Contains(high, "100%") {
		t.Fatalf("at twice what it is worth the floor did not say by how much: %s", high)
	}
}

func TestTheFloorTellsYouAboutTheWholeMarketNotTheRoom(t *testing.T) {
	t.Parallel()
	// Holding the sell side is meant to inform the buy side across town. If it
	// only spoke in its own room it would be worth nothing, because you already
	// know what a price is when you are standing on it.
	w := trading(t, true)
	w.Player.Location = "docks"
	said := ""
	for _, a := range w.Actions("docks") {
		if strings.HasPrefix(a.ID, "buy:") {
			said += a.Detail
		}
	}
	if !strings.Contains(said, "Your floor says") {
		t.Fatalf("the exchange said nothing at the waterfront: %s", said)
	}
	theirs := trading(t, false)
	theirs.Player.Location = "docks"
	for _, a := range theirs.Actions("docks") {
		if strings.Contains(a.Detail, "Your floor says") {
			t.Fatalf("a floor you do not hold spoke anyway: %s", a.Detail)
		}
	}
}

func TestAFloorNobodyIsStandingOnSaysNothing(t *testing.T) {
	t.Parallel()
	w := trading(t, true)
	trade, _ := TradeOf(theFloor)
	prop := w.Properties[theFloor]
	prop.Trouble = true
	if w.TheExchange() != "" {
		t.Fatal("a floor whose weights are condemned was still settling prices")
	}
	prop.Trouble = false
	prop.Staff = trade.Hands - 1
	if w.TheExchange() != "" {
		t.Fatalf("short-handed at %d of %d and still settling prices", prop.Staff, trade.Hands)
	}
	prop.Staff = trade.Hands
	if w.TheExchange() != theFloor {
		t.Fatal("a working floor of yours settles nothing")
	}
}

func TestTheFloorKnowsAWarIsHoldingPricesUp(t *testing.T) {
	t.Parallel()
	w := trading(t, true)
	for i := range w.Conflicts {
		w.Conflicts[i].State, w.Conflicts[i].Hostility = "peace", 0
	}
	if quiet := w.TheWarPremium(); quiet != "" {
		t.Fatalf("a city at peace: %q", quiet)
	}
	for i := range w.Conflicts {
		w.Conflicts[i].State, w.Conflicts[i].Hostility = "war", 80
	}
	if !w.CityAtWar() {
		t.Skip("this city will not go to war")
	}
	if w.TheWarPremium() == "" {
		t.Fatal("the shooting is on and the floor has not noticed")
	}
	theirs := trading(t, false)
	for i := range theirs.Conflicts {
		theirs.Conflicts[i].State, theirs.Conflicts[i].Hostility = "war", 80
	}
	if theirs.TheWarPremium() != "" {
		t.Fatal("a floor you do not hold told you about the war premium")
	}
}
