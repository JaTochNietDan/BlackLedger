package core

import "testing"

// Sitting down at a table is something the world knows about, not a screen the
// interface opens on its own. It matters because the last hand, the last spin
// and the last three drums are saved state: without a sitting, walking up to a
// table shows you a game somebody already played, which reads as a game that
// started without you.

func sitter(t *testing.T) *World {
	t.Helper()
	// Somebody else's room: you cannot win money from your own house, so a
	// house keeper is the wrong person to test a seat with.
	w := gambler(t)
	w.Event, w.District = nil, 9
	return w
}

func TestSittingDownLeavesNothingOnTheTable(t *testing.T) {
	t.Parallel()
	w := sitter(t)
	if err := w.PullHandle("club", 5); err != nil {
		t.Fatal(err)
	}
	if err := w.PlayWheel("club", "red", 50); err != nil {
		t.Fatal(err)
	}
	if w.Reels == nil || w.Spin == nil {
		t.Fatal("nothing was played, so this proves nothing")
	}
	if err := w.Sit("club", Floor); err != nil {
		t.Fatal(err)
	}
	if w.Seated != "club" {
		t.Fatalf("nobody sat down: %q", w.Seated)
	}
	if w.Reels != nil || w.Spin != nil || w.Hand != nil {
		t.Fatal("the table still had the last player's game on it")
	}
	if w.MachineDescription()["pulled"] != false || w.WheelDescription()["spun"] != false {
		t.Fatal("the felt reported a game in progress to somebody who just sat down")
	}
}

func TestYouCannotGetUpInTheMiddleOfAHand(t *testing.T) {
	t.Parallel()
	w := sitter(t)
	if err := w.Sit("club", Floor); err != nil {
		t.Fatal(err)
	}
	if err := w.Deal("club", 50); err != nil {
		t.Fatal(err)
	}
	if w.RiseReadiness() == "" {
		t.Fatal("walked away from a live hand")
	}
	for w.Hand != nil && !w.Hand.Done {
		if err := w.Stand(); err != nil {
			t.Fatal(err)
		}
	}
	if reason := w.RiseReadiness(); reason != "" {
		t.Fatalf("could not get up from a finished hand: %s", reason)
	}
	if err := w.Rise(); err != nil {
		t.Fatal(err)
	}
	if w.Seated != "" || w.Hand != nil {
		t.Fatalf("still at the table: seated %q", w.Seated)
	}
}

func TestTheRoomOffersASeatRatherThanAGame(t *testing.T) {
	t.Parallel()
	w := sitter(t)
	sit := actionByID(w.Actions("club"), "sit:"+Floor)
	if sit == nil || sit.Disabled {
		t.Fatalf("a club with tables did not offer a seat: %+v", sit)
	}
	// A room with nothing to play in it has no seat to offer.
	w.Player.Location = "laundry"
	if actionByID(w.Actions("laundry"), "sit") != nil {
		t.Fatal("a laundry offered a seat at the tables")
	}
	if reason := w.SitReadiness("laundry"); reason == "" {
		t.Fatal("sitting down was allowed in a laundry")
	}
}

func TestASeatIsTakenWhereTheGameIs(t *testing.T) {
	t.Parallel()
	w := sitter(t)
	if err := w.Sit("club", Floor); err != nil {
		t.Fatal(err)
	}
	w.Player.Location = "bar"
	// Walking out of the room is getting up, whatever the interface thinks.
	w.Advance(1)
	if w.Seated != "" {
		t.Fatalf("still seated at the club from the bar: %q", w.Seated)
	}
}
