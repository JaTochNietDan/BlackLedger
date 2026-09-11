package core

import "testing"

// The pawnbroker's window.
//
// A ticket running out used to be a wall: the thing was "sold to somebody else"
// and that was the end of it. It is a shelf now, and a shelf is a place — the
// city's forfeits have a price on them and anybody can have them, the player
// included, and the thing that was theirs last week included.

// shop is a campaign standing at the counter with money.
func shop(t *testing.T, seed uint32) *World {
	t.Helper()
	w := New(seed)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 50000, 60
	w.Player.Location = w.thePawnshop()
	return w
}

func days(w *World, n int) {
	for i := 0; i < n; i++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
}

func TestTheCityFillsTheWindowOnItsOwn(t *testing.T) {
	t.Parallel()
	// Across several cities, because the first rule written for this asked a
	// modelled person to be broke and to own a car, and no such person exists:
	// fifteen people drive and the thinnest purse among them is over $140. It
	// passed on nothing.
	for _, seed := range []uint32{11, 29, 47, 53, 83} {
		w := shop(t, seed)
		if len(w.Window) != 0 {
			t.Fatalf("city %d: the window had something in it on day one", seed)
		}
		days(w, 21)
		if len(w.Window) < 3 {
			t.Fatalf("city %d: three weeks and the window holds %d things", seed, len(w.Window))
		}
	}
}

func TestWhatYouCouldNotRedeemIsStillOnTheShelf(t *testing.T) {
	t.Parallel()
	w := shop(t, 53)
	w.Player.Dress, w.Player.DressWear = 2, 30
	was := AttireByTier(2).Label
	if err := w.Pawn("dress"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Dress != 0 {
		t.Fatal("it is still on their back")
	}
	days(w, PawnDays+2)
	var mine *Shelf
	for i := range w.Window {
		if w.Window[i].Yours {
			mine = &w.Window[i]
		}
	}
	if mine == nil {
		t.Fatal("the ticket ran out and the suit went nowhere")
	}
	if mine.What() != was {
		t.Fatalf("the window holds %q where %q went in", mine.What(), was)
	}
	if mine.Wear != 30 {
		t.Fatalf("it came back %d%% worn where it went in 30%%", mine.Wear)
	}
	// And it can be had again, which is the whole point: a bad trade rather
	// than a wall.
	w.Player.Location = w.thePawnshop()
	if reason := w.WindowReadiness(mine.ID); reason != "" {
		t.Fatalf("refused: %s", reason)
	}
	cash, id := w.Player.Cash, mine.ID
	if err := w.BuyFromWindow(id); err != nil {
		t.Fatal(err)
	}
	if w.Player.Dress != 2 {
		t.Fatal("paid for it and did not get it")
	}
	if w.Player.Cash >= cash {
		t.Fatal("it came off the shelf for nothing")
	}
	if w.shelfByID(id) != nil {
		t.Fatal("it is on the shelf and on their back at once")
	}
}

func TestTheWindowAsksMoreThanItLentAndLessThanNew(t *testing.T) {
	t.Parallel()
	w := shop(t, 53)
	days(w, 21)
	if len(w.Window) == 0 {
		t.Fatal("nothing in the window to price")
	}
	for _, s := range w.Window {
		if s.Ask <= s.Lent {
			t.Fatalf("%s: asks $%d having lent $%d, which is not a trade", s.What(), s.Ask, s.Lent)
		}
		if s.Ask >= worthNew(s.Kind, s.Tier) {
			t.Fatalf("%s: asks $%d where a new one is $%d, so the window is worth nothing",
				s.What(), s.Ask, worthNew(s.Kind, s.Tier))
		}
	}
}

func TestYourOwnCounterSellsYouStockAtWhatItLent(t *testing.T) {
	t.Parallel()
	w := shop(t, 53)
	days(w, 21)
	if len(w.Window) == 0 {
		t.Fatal("nothing in the window")
	}
	s := w.Window[0]
	asked := w.WindowPrice(s.ID)
	if asked != s.Ask {
		t.Fatalf("a counter that is not yours wants $%d where it asks $%d", asked, s.Ask)
	}
	own(w, w.thePawnshop())
	mine := w.WindowPrice(s.ID)
	if mine != s.Lent {
		t.Fatalf("your own counter wants $%d where it lent $%d", mine, s.Lent)
	}
	if mine >= asked {
		t.Fatalf("holding the counter saved nothing: $%d against $%d", mine, asked)
	}
	// And the saving is the trade itself, not a rounding.
	if asked-mine < asked/4 {
		t.Fatalf("holding the counter saved $%d off $%d, which is not a margin worth owning",
			asked-mine, asked)
	}
}

func TestNoTwoThingsInTheWindowReadTheSame(t *testing.T) {
	t.Parallel()
	// Two cards in a room under one name is the oldest fault here, and a shelf
	// of second-hand suits is where it would happen.
	for _, seed := range []uint32{11, 29, 47, 53, 83} {
		w := shop(t, seed)
		days(w, 40)
		seen := map[string]bool{}
		for _, a := range w.Actions(w.thePawnshop()) {
			if seen[a.Label] {
				t.Fatalf("city %d: two cards at the counter both say %q", seed, a.Label)
			}
			seen[a.Label] = true
		}
	}
}

func TestTheWindowIsNotAWarehouse(t *testing.T) {
	t.Parallel()
	// The daily stocking stops at the shelf's size, so it can never overfill on
	// its own and a test that only let the days run proved nothing. What can
	// overfill it is forfeiture: a ticket running out puts a thing in the
	// window whether there is room or not, because the shop does not decline to
	// take what it has already lent against.
	w := shop(t, 53)
	days(w, 30)
	for len(w.Window) < WindowHolds {
		w.shelve(Shelf{Kind: "dress", Tier: 1, Wear: 50, Ask: 40, Lent: 20})
	}
	oldest := w.Window[0].ID
	w.shelve(Shelf{Kind: "car", Tier: 2, Wear: 10, Ask: 900, Lent: 400, Yours: true})
	if len(w.Window) != WindowHolds {
		t.Fatalf("the shelf holds %d where it has room for %d", len(w.Window), WindowHolds)
	}
	if w.shelfByID(oldest) != nil {
		t.Fatal("the shelf took another thing and let nothing go")
	}
	if w.Window[len(w.Window)-1].Kind != "car" {
		t.Fatal("the thing that came in is not on the shelf")
	}
	// And a long campaign never grows one.
	days(w, 200)
	if len(w.Window) > WindowHolds {
		t.Fatalf("two hundred days and the shelf holds %d things", len(w.Window))
	}
}
