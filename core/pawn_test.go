package core

import (
	"strings"
	"testing"
)

// "There is always another business, and always more people to put in the
// city", and "find ways to link them together."
//
// A pawnbroker is the trade that sits at the join between everything already
// here: it is where what gets taken off the street turns into money, and it is
// where somebody who is short turns what they own into money and hopes to get
// it back.

func short(t *testing.T) *World {
	t.Helper()
	w := New(67)
	w.Event, w.District = nil, 9
	w.Player.Health = 100
	w.Player.Location = w.thePawnshop()
	if w.Player.Location == "" {
		t.Fatal("this city has no pawnbroker")
	}
	return w
}

func TestThePawnbrokerIsARealAddress(t *testing.T) {
	t.Parallel()
	id := ""
	for _, l := range Locations {
		if l.Kind == "pawn" {
			id = l.ID
		}
	}
	if id == "" {
		t.Fatal("no pawnbroker in this city")
	}
	if PlaceIncome[id] <= 0 {
		t.Fatal("the pawnbroker earns nothing")
	}
	if _, runs := TradeOf(id); !runs {
		t.Fatal("the pawnbroker has no trade to run")
	}
	if len(mannerByPlace[id]) == 0 {
		t.Fatal("nobody can be killed at the pawnbroker")
	}
	who := 0
	for _, s := range streetTrades {
		if s.place == id {
			who++
		}
	}
	if who == 0 {
		t.Fatal("the pawnbroker has nobody working in it")
	}
}

func TestYouCanPawnTheSuitOffYourBackAndGetItBack(t *testing.T) {
	t.Parallel()
	w := short(t)
	w.Player.Dress, w.Player.DressWear = 2, 100
	w.Player.Cash = 5
	standing := w.Presence()

	lent := w.PawnValue("dress")
	if lent <= 0 {
		t.Fatal("a tailored suit is worth nothing over the counter")
	}
	if err := w.Pawn("dress"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Dress != 0 {
		t.Fatal("the suit is still on your back")
	}
	if w.Player.Cash != 5+lent {
		t.Fatalf("paid %d for it, quoted %d", w.Player.Cash-5, lent)
	}
	if w.Presence() >= standing {
		t.Fatal("standing in working clothes was worth the same as the suit")
	}
	// Getting it back costs more than they gave you. That is the trade.
	back := w.RedeemPrice("dress")
	if back <= lent {
		t.Fatalf("redeeming cost %d against %d lent, which is a charity", back, lent)
	}
	w.Player.Cash = back
	if err := w.Redeem("dress"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Dress != 2 {
		t.Fatalf("got back a tier %d suit", w.Player.Dress)
	}
	if w.Player.Cash != 0 {
		t.Fatalf("the redemption cost %d", back-w.Player.Cash)
	}
}

func TestATicketRunsOutAndTheThingIsSold(t *testing.T) {
	t.Parallel()
	w := short(t)
	w.Player.Car, w.Player.CarWear = 2, 100
	w.Player.Fuel, w.Player.Fuelled = FuelFull, max(1, w.Minute)
	if err := w.Pawn("car"); err != nil {
		t.Fatal(err)
	}
	w.Player.Cash = w.RedeemPrice("car")
	if w.Player.Car != 0 {
		t.Fatal("still driving it")
	}
	if reason := w.RedeemReadiness("car"); reason != "" {
		t.Fatalf("could not redeem it the same day: %s", reason)
	}
	w.Advance(PawnDays*1440 + 1440)
	if reason := w.RedeemReadiness("car"); reason == "" {
		t.Fatal("the ticket never ran out")
	}
	if !strings.Contains(strings.ToLower(w.RedeemReadiness("car")), "sold") {
		t.Fatalf("the refusal does not say what happened: %q", w.RedeemReadiness("car"))
	}
}

func TestWhatIsTakenOffTheStreetPassesThroughTheShop(t *testing.T) {
	t.Parallel()
	w := New(67)
	w.Event, w.District = nil, 9
	shop := w.thePawnshop()
	before := w.Properties[shop].Custom
	w.FenceAbout(FenceTrade)
	if w.Properties[shop].Custom <= before {
		t.Fatalf("a shop that buys what is stolen did not notice: %d then %d", before, w.Properties[shop].Custom)
	}
	// And a robbery is one of the things that feeds it.
	w2 := New(67)
	w2.Event, w2.District = nil, 9
	w2.Player.Location, w2.Player.Health, w2.Player.Cash = "laundry", 100, 500
	was := w2.Properties[shop].Custom
	if err := w2.Rob("laundry"); err != nil {
		t.Fatal(err)
	}
	if w2.Properties[shop].Custom <= was {
		t.Fatal("a till went out of a laundry and nobody in this city bought any of it")
	}
}
