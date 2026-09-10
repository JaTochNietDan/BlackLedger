package core

import "testing"

// "While managing the Green Baize I don't see its current funds or how to add
// to the funds or withdraw from the funds dynamically like we talked about."
//
// The Green Baize is the poolhall, which is a racket rather than a casino, so
// none of the float machinery reached it — and it is the one room in this city
// that now runs a game and charges for the table. The money it takes for the
// seat went straight into the player's pocket without ever being anywhere they
// could look at it.

func baize(t *testing.T) *World {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, BackRoom)
	w.Player.Location = BackRoom
	return w
}

func TestTheGreenBaizeHasFundsYouCanSee(t *testing.T) {
	w := baize(t)
	if !RunsAGame(BackRoom) {
		t.Fatal("the room with a card game behind it does not run a game")
	}
	if w.Properties[BackRoom] == nil {
		t.Fatal("the poolhall is not a property")
	}
	// The room's money is published where the room is drawn.
	for _, l := range w.Public()["locations"].([]map[string]any) {
		if l["id"] != BackRoom {
			continue
		}
		if _, shown := l["bankroll"]; !shown {
			t.Fatal("the Green Baize does not say what is in its till")
		}
		return
	}
	t.Fatal("the Green Baize is not on the map")
}

func TestYouCanPutMoneyInTheBaizeAndTakeItOut(t *testing.T) {
	w := baize(t)
	cash := w.Player.Cash
	if reason := w.BankrollReadiness(BackRoom, 600); reason != "" {
		t.Fatalf("putting money into the Green Baize was refused: %s", reason)
	}
	if err := w.Bankroll(BackRoom, 600); err != nil {
		t.Fatalf("putting money in failed: %v", err)
	}
	if w.Properties[BackRoom].Bankroll != 600 {
		t.Fatalf("$600 went in and the till holds $%d", w.Properties[BackRoom].Bankroll)
	}
	if w.Player.Cash != cash-600 {
		t.Fatalf("putting $600 in cost $%d", cash-w.Player.Cash)
	}
	// And out again, in whatever figure the player names rather than a lot the
	// room decided for them.
	if reason := w.DrawReadiness(BackRoom, 250); reason != "" {
		t.Fatalf("taking money out was refused: %s", reason)
	}
	if err := w.Draw(BackRoom, 250); err != nil {
		t.Fatalf("taking money out failed: %v", err)
	}
	if w.Properties[BackRoom].Bankroll != 350 {
		t.Fatalf("$250 came out of $600 and left $%d", w.Properties[BackRoom].Bankroll)
	}
	if w.Player.Cash != cash-350 {
		t.Fatalf("the player is $%d down after putting $600 in and taking $250 out", cash-w.Player.Cash)
	}
	// You cannot take out what is not there.
	if w.DrawReadiness(BackRoom, 5000) == "" {
		t.Fatal("a room with $350 in it would hand over $5,000")
	}
}

func TestBothWaysAreOfferedWhereTheRoomIs(t *testing.T) {
	w := baize(t)
	find := func(id string) *Action {
		for _, a := range w.Actions(BackRoom) {
			if a.ID == id {
				return &a
			}
		}
		return nil
	}
	in, out := find("bankroll"), find("draw")
	if in == nil || out == nil {
		t.Fatal("the Green Baize offers no way to put money in or take it out")
	}
	if in.Disabled {
		t.Fatalf("putting money into your own room is refused: %s", in.Reason)
	}
	if in.Sum == nil || out.Sum == nil {
		t.Fatal("the figure is not the player's to name")
	}
	if err := w.apply(Command{Kind: "bankroll", Target: BackRoom, Amount: 400, RequestID: "baizeputitin1"}); err != nil {
		t.Fatalf("putting money in through the same path as everything else failed: %v", err)
	}
	if w.Properties[BackRoom].Bankroll != 400 {
		t.Fatalf("the till holds $%d", w.Properties[BackRoom].Bankroll)
	}
}

// The seat charge is the room's money, so it goes into the room's till rather
// than into the holder's pocket. That is what makes the till worth looking at.
func TestWhatTheTableTakesGoesIntoTheTill(t *testing.T) {
	w := baize(t)
	cash, till := w.Player.Cash, w.Properties[BackRoom].Bankroll
	for day := 0; day < 30; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	took := w.Properties[BackRoom].Bankroll - till
	t.Logf("thirty days: the table took $%d, and the till holds $%d", w.BackRoomTake, w.Properties[BackRoom].Bankroll)
	if w.BackRoomTake == 0 {
		t.Fatal("nobody paid for a seat in a month")
	}
	if took < w.BackRoomTake {
		t.Fatalf("the table took $%d and the till is only $%d fuller", w.BackRoomTake, took)
	}
	// And it is not also in the player's pocket, because it is in one place.
	if w.Player.Cash-cash >= w.BackRoomTake {
		t.Fatal("the seat money was paid into the till and into the holder's hands")
	}
}
