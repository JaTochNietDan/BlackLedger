package core

import (
	"strings"
	"testing"
)

// The other end of taking a car apart.
//
// A bench lives on parts. The player was tearing them off a car in the street
// and the only thing that ever happened to them was a figure in the log and a
// nudge to every garage in the city — including the ones belonging to whoever
// the player is at war with. A garage of their own still had to be restocked
// for cash, off the same shelves the parts should have gone onto. Same distance
// between two halves as the still and the bar.
func stripWith(t *testing.T, fit func(w *World)) (*World, string, string) {
	t.Helper()
	w := New(53)
	w.Event, w.District, w.Player.Cash, w.Player.Health = nil, 9, 40000, 100
	fit(w)
	// Somewhere with a car worth taking apart, found rather than named.
	for _, l := range Locations {
		if l.District > w.District {
			continue
		}
		w.Player.Location = l.ID
		if w.StripReadiness(l.ID) != "" {
			continue
		}
		var card Action
		for _, a := range w.Actions(l.ID) {
			if a.ID == "strip" {
				card = a
			}
		}
		if card.ID == "" || card.Disabled {
			continue
		}
		return w, l.ID, card.Detail
	}
	t.Skip("nowhere in this city has a car worth taking apart")
	return nil, "", ""
}

func TestWhatComesOffACarGoesOntoABenchOfYours(t *testing.T) {
	t.Parallel()
	garage := ""
	for _, l := range Locations {
		if l.Kind == "garage" {
			garage = l.ID
			break
		}
	}
	if garage == "" {
		t.Skip("this city has no garage")
	}
	full := TradeRestockAmount("garage")
	if full == 0 {
		t.Fatal("a garage's store holds nothing, so this measures nothing")
	}

	w, where, detail := stripWith(t, func(w *World) {
		w.Properties[garage].Owner = "player:1"
		w.Properties[garage].Supply = 0
	})
	// The card says where they are going before the crowbar comes out.
	place, _ := PlaceByID(garage)
	if !strings.Contains(detail, place.Name) {
		t.Fatalf("the card does not say the parts have somewhere of yours to go: %q", detail)
	}

	before := w.Properties[garage].Supply
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision,
		Kind: "strip", Target: where})
	if err != nil {
		t.Fatal(err)
	}
	after := next.Properties[garage].Supply
	if after <= before {
		t.Fatalf("a car went to pieces and the bench went from %d to %d", before, after)
	}
	if after > full {
		t.Fatalf("the bench holds %d of %d", after, full)
	}
}

// And somebody with no garage loses nothing they had: the parts are still worth
// money, the trade still lifts, and no shelf anywhere fills up by itself.
func TestWithoutABenchTheCityIsUnchanged(t *testing.T) {
	t.Parallel()
	garage := ""
	for _, l := range Locations {
		if l.Kind == "garage" {
			garage = l.ID
			break
		}
	}
	if garage == "" {
		t.Skip("this city has no garage")
	}
	w, where, detail := stripWith(t, func(w *World) {
		w.Properties[garage].Supply = 5
	})
	place, _ := PlaceByID(garage)
	if strings.Contains(detail, place.Name) {
		t.Fatalf("a player who owns no garage is told the parts go to one: %q", detail)
	}
	cash, supply := w.Player.Cash, w.Properties[garage].Supply
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision,
		Kind: "strip", Target: where})
	if err != nil {
		t.Fatal(err)
	}
	if got := next.Properties[garage].Supply; got != supply {
		t.Fatalf("somebody else's bench went from %d to %d", supply, got)
	}
	if next.Player.Cash <= cash {
		t.Fatalf("a car went for parts and paid nothing: $%d then, $%d now", cash, next.Player.Cash)
	}
}
