package core

import (
	"strings"
	"testing"
)

// The Monarch is a business now.
//
// Every address that can be bought is a business — hands, supplies, trouble,
// cover, somebody behind a counter worth asking. Every address that can only be
// taken by force was a number. That is backwards: thirty people drink at the
// bar of an evening and twenty-eight at the club, and a room you fight a war
// for was the only kind in the city that ran on nothing.

const theClub = "club"

func clubbing(t *testing.T, hold bool) *World {
	t.Helper()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 90
	if hold {
		own(w, theClub)
	}
	for d := 0; d < 4; d++ {
		w.Advance(1440)
		w.Event = nil
	}
	return w
}

func TestASeatIsABusinessAndNotANumber(t *testing.T) {
	t.Parallel()
	w := clubbing(t, false)
	trade, runs := TradeOf(theClub)
	if !runs {
		t.Fatal("The Monarch runs on nothing")
	}
	prop := w.Properties[theClub]
	if prop.Staff != trade.Hands {
		t.Fatalf("it opens with %d of %d hands", prop.Staff, trade.Hands)
	}
	if prop.Supply != trade.RestockAmount {
		t.Fatalf("it opens with %d of %d supplies", prop.Supply, trade.RestockAmount)
	}
	// A seat has no price on purpose. It is a family's, and it changes hands
	// the way a family's ground changes hands.
	place, _ := PlaceByID(theClub)
	if place.Cost > 0 {
		t.Fatalf("a family's seat is for sale at $%d", place.Cost)
	}
	if !w.SeatOf(prop.Owner, theClub) {
		t.Fatal("nobody holds The Monarch as their own ground")
	}
}

func TestTakingTheSeatGivesYouEverythingThatRunsInIt(t *testing.T) {
	t.Parallel()
	mine := clubbing(t, true)
	trade, _ := TradeOf(theClub)
	prop := mine.Properties[theClub]
	if !mine.Own(theClub) {
		t.Fatal("it is not yours")
	}
	mine.Player.Location = theClub
	offered := map[string]bool{}
	for _, a := range mine.Actions(theClub) {
		offered[a.ID] = true
	}
	// The things every other business of yours offers.
	for _, card := range []string{"inspect", "hire", "restock"} {
		if !offered[card] {
			t.Errorf("a room you took offers no %q", card)
		}
	}
	// Somebody behind the counter, worth asking.
	if len(prop.Hands) == 0 {
		t.Fatalf("%d staff and nobody in the list", prop.Staff)
	}
	if trade.Cover <= 0 {
		t.Fatal("a room full of cash and strangers explains nothing")
	}
	if mine.launderCapacity(theClub) <= 0 {
		t.Fatal("the books will not take anything")
	}
}

func TestAGoodNightIsWorthStandingRatherThanMoney(t *testing.T) {
	t.Parallel()
	theirs, mine := clubbing(t, false), clubbing(t, true)
	// The room has to be busy enough to be worth something before "somebody
	// else's room pays you nothing" means anything. Without this the assertion
	// passed on a quiet night with the ownership check deleted.
	if theirs.TheEveningCrowd(theClub) < ClubPerHead {
		t.Fatalf("only %d drink there, so this proves nothing", theirs.TheEveningCrowd(theClub))
	}
	if got := theirs.StandingFromTheDoor(theClub); got != 0 {
		t.Fatalf("somebody else's room paid you %d standing", got)
	}
	got := mine.StandingFromTheDoor(theClub)
	if got <= 0 {
		t.Fatalf("%d people drink in a room of yours and it is worth nothing",
			mine.TheEveningCrowd(theClub))
	}
	// And it is paid, in standing rather than cash.
	respect, cash := mine.Player.Respect, mine.Player.Cash
	mine.ClubNight()
	if mine.Player.Respect != respect+got {
		t.Fatalf("standing went from %d to %d where %d was earned",
			respect, mine.Player.Respect, got)
	}
	if mine.Player.Cash != cash {
		t.Fatal("a night at the door paid in money")
	}
	said := false
	for _, r := range mine.History {
		if strings.Contains(r.Title, "A good night") {
			said = true
		}
	}
	if !said {
		t.Fatal("the night was paid and nobody was told")
	}
}

func TestAStandingReputationIsBuiltAndNotHandedOver(t *testing.T) {
	t.Parallel()
	mine := clubbing(t, true)
	// However full the room, one night is one night.
	if got := mine.StandingFromTheDoor(theClub); got > ClubNightly {
		t.Fatalf("one night at the door was worth %d standing", got)
	}
	// A room with something wrong in the back is a room people talk about for
	// the wrong reason, and so is a room with a queue and nobody serving it.
	prop := mine.Properties[theClub]
	prop.Trouble = true
	if got := mine.StandingFromTheDoor(theClub); got != 0 {
		t.Fatalf("a room with trouble in the back earned %d standing", got)
	}
	prop.Trouble = false
	trade, _ := TradeOf(theClub)
	prop.Staff = trade.Hands - 1
	if got := mine.StandingFromTheDoor(theClub); got != 0 {
		t.Fatalf("short-handed at %d of %d and still earning standing",
			prop.Staff, trade.Hands)
	}
}

func TestAnEmptyRoomMakesNobodysName(t *testing.T) {
	t.Parallel()
	mine := clubbing(t, true)
	// Nobody in it, nobody talking about it. The whole point is that standing
	// is made in front of people.
	// Nobody whose evening this room is.
	for i := range mine.NPCs {
		mine.NPCs[i].Dead = true
	}
	if got := mine.StandingFromTheDoor(theClub); got != 0 {
		t.Fatalf("an empty room was worth %d standing", got)
	}
}
