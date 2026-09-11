package core

import (
	"strings"
	"testing"
)

// A family that stops asking and starts taking.
//
// Nothing in this city ever cost a legitimate operator anything but money, and
// the money was a week of one shop's takings at a time. A publican who buys,
// hires and restocks ends the richest policy in the game with no police
// attention and no seizures. That is a coherent arrangement and it is also the
// end of the story: a family watching somebody assemble half a district does
// not go on sending for an envelope.

func claimant(t *testing.T, holdings int) (*World, string, string) {
	t.Helper()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 90
	// A family with ground of its own, thinking poorly of the player.
	f := w.faction("bellandi")
	if f == nil {
		t.Skip("no Bellandi Family in this city")
	}
	f.Goodwill = ClaimStanding - 5
	district := -1
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Owner == f.ID {
			district = l.District
		}
	}
	if district < 0 {
		t.Skip("that family holds nothing")
	}
	// And the player holding that many earning addresses on the same streets.
	mine := ""
	taken := 0
	for _, l := range Locations {
		if l.District != district || taken >= holdings {
			continue
		}
		prop := w.Properties[l.ID]
		if prop == nil || prop.Income <= 0 || prop.Owner == f.ID {
			continue
		}
		own(w, l.ID)
		if mine == "" {
			mine = l.ID
		}
		taken++
	}
	if taken < holdings {
		t.Skipf("only %d earning addresses on those streets", taken)
	}
	return w, mine, f.ID
}

func TestAFamilyAsksForTheShareUntilYouHoldTooMuch(t *testing.T) {
	t.Parallel()
	// One or two places and they want their share.
	few, id, family := claimant(t, ClaimHoldings-1)
	if _, wanted := few.TheirClaim(family, id); wanted {
		t.Fatalf("they came for the place over %d holdings", ClaimHoldings-1)
	}
	// Enough of the district, and they want the place.
	many, id, family := claimant(t, ClaimHoldings)
	buyout, wanted := many.TheirClaim(family, id)
	if !wanted {
		t.Fatalf("%d holdings on their own streets and they still asked for an envelope", ClaimHoldings)
	}
	if buyout <= many.TheirShare(id) {
		t.Fatalf("buying them off costs $%d against an ordinary share of $%d", buyout, many.TheirShare(id))
	}
	// Somebody they are fond of gets the envelope, however much they hold.
	many.faction(family).Goodwill = ClaimStanding
	if _, still := many.TheirClaim(family, id); still {
		t.Fatal("a family that thinks well of you came for the place anyway")
	}
}

func TestHandingItOverIsTheOneThingPayingNeverDid(t *testing.T) {
	t.Parallel()
	w, id, family := claimant(t, ClaimHoldings)
	prop := w.Properties[id]
	hands := len(prop.Hands)
	before := w.faction(family).Goodwill
	if err := w.HandOver(id, family); err != nil {
		t.Fatal(err)
	}
	if w.Own(id) {
		t.Fatal("the deed is still yours")
	}
	if prop.Owner != family {
		t.Fatalf("it went to %q", prop.Owner)
	}
	// The people behind the counter stay. They work there; who is being paid
	// is not their decision.
	if len(prop.Hands) != hands {
		t.Fatalf("%d behind the counter where there were %d", len(prop.Hands), hands)
	}
	if w.faction(family).Goodwill <= before {
		t.Fatal("giving them a business bought no standing at all")
	}
	if w.NextPressure <= w.Minute {
		t.Fatal("they came back for another the same afternoon")
	}
	// And it cannot be given twice.
	if err := w.HandOver(id, family); err == nil {
		t.Fatal("handed over a place that was no longer yours")
	}
}

func TestTheSceneAsksForThePlaceAndNotAShare(t *testing.T) {
	t.Parallel()
	w, id, family := claimant(t, ClaimHoldings)
	// The demand itself, rather than waiting for the schedule to deliver one.
	// Calling the retaliation and skipping when nothing arrived is the oldest
	// fault in this suite: a test that stands down reports a pass and proves
	// nothing.
	w.BusinessPressure()
	if w.Event == nil || w.Event.Kind != "business_pressure" {
		t.Fatal("holding half a district brought no demand at all")
	}
	if w.Event.Target != id {
		// They came for a different one of the player's places on the same
		// streets, which is the same claim about a different address.
		id = w.Event.Target
		family = w.Event.Actor
	}
	if !strings.Contains(w.Event.Body, "I am asking for") {
		t.Fatalf("they asked for a share: %s", w.Event.Body)
	}
	offered := map[string]bool{}
	for _, c := range w.Event.Choices {
		offered[c.ID] = true
	}
	for _, want := range []string{"hand", "pay", "resist"} {
		if !offered[want] {
			t.Errorf("no way to %q when a family comes for the place", want)
		}
	}
	// Buying them off costs what the claim says, not what a share costs.
	buyout, _ := w.TheirClaim(family, id)
	for _, c := range w.Event.Choices {
		if c.ID == "pay" && c.Cost != buyout {
			t.Fatalf("the card says $%d and the claim is $%d", c.Cost, buyout)
		}
	}
}
