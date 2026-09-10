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
	if err := w.Sit(BackRoom); err != nil {
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
	if err := w.Sit(BackRoom); err != nil {
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
	if err := w.Sit(BackRoom); err != nil {
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
	seat := find("sit")
	if seat == nil || seat.Disabled {
		t.Fatal("there is no way to sit down to the game behind the poolhall")
	}
	if err := w.apply(Command{Kind: "sit", Target: BackRoom, RequestID: "sitinthebackroom1"}); err != nil {
		t.Fatalf("sitting down through the same path as everything else failed: %v", err)
	}
	if find("rise") == nil {
		t.Fatal("the player is sitting at a table with no way to get up from it")
	}
}
