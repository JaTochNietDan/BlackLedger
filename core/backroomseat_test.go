package core

import "testing"

// "It should again be a separate scene that takes up the screen when you're
// playing it and you have to leave it rather than right now it just lives in a
// small box above the action bar. That's silly stuff. We need to stop doing
// that in future and always dedicate these games to their own screen."
//
// The screen is the view's half of it. The core's half is that sitting in on a
// game is a seat the world knows about, exactly as it is at the tables: you are
// there until you get up, you cannot get up in the middle of a hand, and
// walking out of the room ends it.

func TestTheBackRoomIsASeatLikeAnyOther(t *testing.T) {
	t.Parallel()
	w, _ := backroom(t)
	if !Playable(BackRoom) {
		t.Fatal("there is a card game behind the poolhall and nothing to sit down to")
	}
	if reason := w.SitReadiness(BackRoom); reason != "" {
		t.Fatalf("taking a seat in the back room was refused: %s", reason)
	}
	if err := w.Sit(BackRoom, Backroom); err != nil {
		t.Fatalf("sitting down failed: %v", err)
	}
	if w.Seated != BackRoom {
		t.Fatalf("the player sat down and the world says they are at %q", w.Seated)
	}
	if err := w.Rise(); err != nil {
		t.Fatalf("getting up from an empty table was refused: %v", err)
	}
	if w.Seated != "" {
		t.Fatal("the player got up and the world still has them at the table")
	}
}

func TestYouCannotWalkOutOnMoneyYouHavePutIn(t *testing.T) {
	t.Parallel()
	w, _ := backroom(t)
	if err := w.Sit(BackRoom, Backroom); err != nil {
		t.Fatal(err)
	}
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatal(err)
	}
	if reason := w.RiseReadiness(); reason == "" {
		t.Fatal("a player with $50 in the middle of the table was free to walk away from it")
	}
	if err := w.Rise(); err == nil {
		t.Fatal("getting up mid-hand was allowed")
	}
	playOut(t, w)
	if reason := w.RiseReadiness(); reason != "" {
		t.Fatalf("the hand is over and getting up is still refused: %s", reason)
	}
	if err := w.Rise(); err != nil {
		t.Fatal(err)
	}
	// And the table is cleared behind them, so nobody walks up to somebody
	// else's finished hand.
	if w.Game != nil {
		t.Fatal("the player left and their hand is still lying on the table")
	}
}

func TestWalkingOutOfThePoolhallEndsTheSitting(t *testing.T) {
	t.Parallel()
	w, _ := backroom(t)
	if err := w.Sit(BackRoom, Backroom); err != nil {
		t.Fatal(err)
	}
	w.Player.Location = "bar"
	w.LeaveTable()
	if w.Seated != "" {
		t.Fatalf("the player is at the bar and seated at %q", w.Seated)
	}
}

// The seat is offered where the game is, and getting up is offered once you
// have taken it.
func TestTheSeatIsOfferedInTheBackRoom(t *testing.T) {
	t.Parallel()
	w, _ := backroom(t)
	find := func(id string) *Action {
		for _, a := range w.Actions(w.Player.Location) {
			if a.ID == id {
				return &a
			}
		}
		return nil
	}
	seat := find("sit:" + Backroom)
	if seat == nil || seat.Disabled {
		t.Fatal("there is no way to sit down to the game behind the poolhall")
	}
	if err := w.apply(Command{Kind: "sit:" + Backroom, Target: BackRoom, RequestID: "sitinthebackroom1"}); err != nil {
		t.Fatalf("sitting down through the same path as everything else failed: %v", err)
	}
	if find("rise") == nil {
		t.Fatal("the player is sitting at a table with no way to get up from it")
	}
}

// Two rooms in this city have both a wall of machines and a room behind the
// room. A seat that only knew the address could not tell them apart, so the
// screen picked which game to draw from what the room held rather than from
// what the player sat down to: "I click play the machines in Saint Agnes but I
// only get the option to play poker."
func TestARoomWithBothOffersBoth(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	both := 0
	for _, l := range Locations {
		if !HasBackRoom(l.ID) || !(HasTables(l.ID) || HasMachines(l.ID)) {
			continue
		}
		both++
		w.Player.Location, w.Event = l.ID, nil
		floor, back := false, false
		for _, a := range w.Actions(l.ID) {
			if a.ID == "sit:"+Floor {
				floor = true
			}
			if a.ID == "sit:"+Backroom {
				back = true
			}
			if a.ID == "sit" {
				t.Fatalf("%s offers one seat that has to guess which game the player meant", l.ID)
			}
		}
		if !floor || !back {
			t.Fatalf("%s has a floor and a room behind it and offers floor=%v back=%v",
				l.ID, floor, back)
		}
	}
	if both < 2 {
		t.Fatalf("%d rooms have both a floor and a room behind them, so this measures nothing", both)
	}
}

// And sitting down to one is not sitting down to the other.
func TestSittingDownToTheMachinesIsNotSittingDownToCards(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Location = "bar"
	if err := w.Sit("bar", Floor); err != nil {
		t.Fatal(err)
	}
	if w.SeatedTo != Floor {
		t.Fatalf("sat down to the machines and the world says %q", w.SeatedTo)
	}
	if err := w.Sit("bar", Backroom); err != nil {
		t.Fatal(err)
	}
	if w.SeatedTo != Backroom {
		t.Fatalf("went through to the back room and the world says %q", w.SeatedTo)
	}
	// And a room with no room behind it cannot be sat down to that way.
	w.Player.Location = "casino"
	if err := w.Sit("casino", Backroom); err == nil {
		t.Fatal("the casino has a back room now")
	}
	// Getting up forgets which it was.
	w.Player.Location = "bar"
	if err := w.Sit("bar", Floor); err != nil {
		t.Fatal(err)
	}
	if err := w.Rise(); err != nil {
		t.Fatal(err)
	}
	if w.SeatedTo != "" {
		t.Fatalf("got up and the world still has them at the %s", w.SeatedTo)
	}
}
