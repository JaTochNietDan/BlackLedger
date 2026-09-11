package core

import (
	"strings"
	"testing"
)

// A still and a bar, joined.
//
// Both halves have been in this game a long time and there was no road between
// them: the still made crates, the only buyer was the market, and the bar was
// restocked with cash. A player who owned both carried crates past his own
// cellar to sell them to somebody else.
func TestARoomOfYoursCanBeRunOffYourOwnStill(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.Event, w.District, w.Player.Cash = nil, 9, 40000

	// A room that sells drink, read out of the trade table rather than named
	// here, so a room that starts selling drink tomorrow is covered too.
	bar := ""
	for _, l := range Locations {
		if trade, ok := TradeOf(l.ID); ok && trade.Drink > 0 {
			bar = l.ID
			break
		}
	}
	if bar == "" {
		t.Skip("nowhere in this city sells drink")
	}
	trade, _ := TradeOf(bar)
	w.Properties[bar].Owner = "player:1"
	w.Properties[bar].Supply = 0
	w.Player.Location = bar

	// With nothing in hand the card is offered and refused, with the reason on
	// it rather than an error behind it.
	card := func() Action {
		for _, a := range w.Actions(bar) {
			if a.ID == "own_cellar" {
				return a
			}
		}
		return Action{}
	}
	empty := card()
	if empty.ID == "" {
		t.Fatal("a room that sells drink does not offer to use its own cellar")
	}
	if !empty.Disabled || !strings.Contains(empty.Reason, "holding") {
		t.Fatalf("a player holding nothing is not told so: disabled %v, %q",
			empty.Disabled, empty.Reason)
	}

	// The still's own output, put in hand the way the still puts it there.
	w.Player.Stock = map[string]int{"moonshine": trade.Drink + 2}
	live := card()
	if live.Disabled {
		t.Fatalf("a player holding enough is refused: %s", live.Reason)
	}

	cash, held := w.Player.Cash, w.Holding("moonshine")
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision,
		Kind: "own_cellar", Target: bar})
	if err != nil {
		t.Fatal(err)
	}
	if got := next.Properties[bar].Supply; got != trade.RestockAmount {
		t.Fatalf("the room is at %d of %d after a delivery", got, trade.RestockAmount)
	}
	if got := next.Holding("moonshine"); got != held-trade.Drink {
		t.Fatalf("%d crates went in and %d left your hands", trade.Drink, held-got)
	}
	// The whole point of the link: no cash. The clock still moves, so what is
	// asked is that the restock's own price did not go out.
	if spent := cash - next.Player.Cash; spent >= w.RestockCost(bar) {
		t.Fatalf("running it off your own cellar took $%d, and buying it in costs $%d",
			spent, w.RestockCost(bar))
	}
}

// And a room nobody drinks in is not offered it at all.
func TestARoomNobodyDrinksInHasNoCellarToRunOff(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.Event, w.District, w.Player.Cash = nil, 9, 40000
	w.Player.Stock = map[string]int{"moonshine": 40}

	dry := ""
	for _, l := range Locations {
		if trade, ok := TradeOf(l.ID); ok && trade.Drink == 0 {
			dry = l.ID
			break
		}
	}
	if dry == "" {
		t.Skip("every trade in this city sells drink")
	}
	w.Properties[dry].Owner = "player:1"
	w.Properties[dry].Supply = 0
	w.Player.Location = dry
	for _, a := range w.Actions(dry) {
		if a.ID == "own_cellar" {
			t.Fatalf("%s does not sell drink and is offered a cellar to run off", dry)
		}
	}
	if reason := w.OwnCellarReadiness(dry); reason == "" {
		t.Fatalf("%s does not sell drink and would take the crates anyway", dry)
	}
}
