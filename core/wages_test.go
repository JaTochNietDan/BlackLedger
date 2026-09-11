package core

import "testing"

// The people behind your counters can walk out over something you did, and
// there is nothing you can do about it except not do it. A wage is the oldest
// answer to that: pay over the rate and they think better of you, pay under it
// and they think less, and either way the books say what it costs.

func payroll(t *testing.T) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	return w, "laundry"
}

func TestTheWageIsYoursToSet(t *testing.T) {
	t.Parallel()
	w, id := payroll(t)
	trade, _ := TradeOf(id)
	if w.WageAt(id) != trade.Wage {
		t.Fatalf("a business nobody has touched pays $%d against a rate of $%d", w.WageAt(id), trade.Wage)
	}
	if reason := w.PayReadiness(id, trade.Wage*2); reason != "" {
		t.Fatalf("paying over the rate was refused: %s", reason)
	}
	if err := w.SetWage(id, trade.Wage*2); err != nil {
		t.Fatalf("setting the wage failed: %v", err)
	}
	if w.WageAt(id) != trade.Wage*2 {
		t.Fatalf("the wage was set to $%d and reads $%d", trade.Wage*2, w.WageAt(id))
	}
	// The books say what it costs, or a decision about money is invisible.
	bill := w.Wages()
	if err := w.SetWage(id, trade.Wage); err != nil {
		t.Fatal(err)
	}
	if w.Wages() >= bill {
		t.Fatalf("halving the wage did not change the wage bill: $%d against $%d", w.Wages(), bill)
	}
	// And nobody works for nothing, or for a figure that would ruin you.
	if w.PayReadiness(id, 0) == "" {
		t.Fatal("a business can pay nothing")
	}
	if w.PayReadiness(id, trade.Wage*20) == "" {
		t.Fatal("a business can promise anything")
	}
}

func TestPayingOverTheRateIsRemembered(t *testing.T) {
	t.Parallel()
	generous, id := payroll(t)
	mean, _ := payroll(t)
	trade, _ := TradeOf(id)
	if err := generous.SetWage(id, trade.Wage*2); err != nil {
		t.Fatal(err)
	}
	if err := mean.SetWage(id, trade.Wage/2); err != nil {
		t.Fatal(err)
	}
	for day := 0; day < 20; day++ {
		for _, w := range []*World{generous, mean} {
			w.Event = nil
			w.Advance(1440)
			w.Event = nil
		}
	}
	well, badly := 0, 0
	for _, who := range generous.Properties[id].Hands {
		well += generous.NPC(who).Trust
	}
	for _, who := range mean.Properties[id].Hands {
		badly += mean.NPC(who).Trust
	}
	t.Logf("after twenty days: paid double they think of you at %d all told, paid half at %d", well, badly)
	if well <= badly {
		t.Fatalf("paying double bought nothing: %d against %d", well, badly)
	}
}

// And it is offered where the business is, in whatever figure the player names.
func TestTheWageIsOfferedAtTheBusiness(t *testing.T) {
	t.Parallel()
	w, id := payroll(t)
	w.Player.Location = id
	var pay *Action
	for _, a := range w.Actions(id) {
		if a.ID == "wage" {
			pay = &a
		}
	}
	if pay == nil {
		t.Fatal("a business of yours offers no way to change what it pays")
	}
	if pay.Disabled {
		t.Fatalf("setting your own wage is refused: %s", pay.Reason)
	}
	if pay.Sum == nil {
		t.Fatal("the figure is not the player's to name")
	}
	trade, _ := TradeOf(id)
	if err := w.apply(Command{Kind: "wage", Target: id, Amount: trade.Wage + 3,
		RequestID: "setthewagehere1"}); err != nil {
		t.Fatalf("setting the wage through the same path as everything else failed: %v", err)
	}
	if w.WageAt(id) != trade.Wage+3 {
		t.Fatalf("the wage reads $%d", w.WageAt(id))
	}
}

// The wage slider had one correct notch and ten wrong ones. What somebody
// thought of the player moved by the same amount whether they were paid a
// dollar over the rate or twelve, because the rule asked whether the wage was
// above the rate and not by how much. Measured over sixty days at a laundry:
// trust 90 either way, and the generous end cost $1,048 against $58. A player
// who worked that out would never move the slider again, and nothing in the
// game told them.
func trustAfter(t *testing.T, seed uint32, over, days int) int {
	t.Helper()
	w := New(seed)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Location, w.Player.Cash = "laundry", 200000
	if err := w.apply(Command{Kind: "acquire", Target: "laundry", RequestID: "whattheythink"}); err != nil {
		t.Fatal(err)
	}
	trade, _ := TradeOf("laundry")
	if over != 0 {
		if err := w.SetWage("laundry", trade.Wage+over); err != nil {
			t.Fatal(err)
		}
	}
	for day := 0; day < days; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	trust := 0
	for _, who := range w.Properties["laundry"].Hands {
		if n := w.NPC(who); n != nil {
			trust += n.Trust
		}
	}
	return trust
}

func TestWhatYouPayOverTheRateIsWorthWhatItIs(t *testing.T) {
	t.Parallel()
	trade, _ := TradeOf("laundry")
	_, most := WageBounds(trade.Wage)
	// The same three seeds, because what the city does to a laundry over sixty
	// days swings the cash by a hundred thousand and this is not about the
	// cash.
	for _, seed := range []uint32{61, 7, 104} {
		little := trustAfter(t, seed, 1, 60)
		some := trustAfter(t, seed, (most-trade.Wage)/2, 60)
		lots := trustAfter(t, seed, most-trade.Wage, 60)
		t.Logf("seed %d: a dollar over the rate is worth %d, half way up %d, the ceiling %d",
			seed, little, some, lots)
		if little >= some || some >= lots {
			t.Fatalf("seed %d: the slider does not run: %d, %d, %d", seed, little, some, lots)
		}
	}
}

// And the other way. Somebody paid a dollar under the rate thinks a little less
// of the player; somebody paid the floor thinks a great deal less.
func TestWhatYouPayUnderTheRateIsWorthWhatItIs(t *testing.T) {
	t.Parallel()
	trade, _ := TradeOf("laundry")
	least, _ := WageBounds(trade.Wage)
	// Up from nothing first, so there is something to take away: everybody in
	// this city starts at nothing and cannot think less of anybody.
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Location, w.Player.Cash = "laundry", 200000
	if err := w.apply(Command{Kind: "acquire", Target: "laundry", RequestID: "downwards"}); err != nil {
		t.Fatal(err)
	}
	for _, who := range w.Properties["laundry"].Hands {
		if n := w.NPC(who); n != nil {
			n.Trust = 100
		}
	}
	mild, harsh := w.Clone(), w.Clone()
	if err := mild.SetWage("laundry", trade.Wage-1); err != nil {
		t.Fatal(err)
	}
	if err := harsh.SetWage("laundry", least); err != nil {
		t.Fatal(err)
	}
	fell := func(x *World) int {
		x.PayDay()
		gone := 0
		for _, who := range x.Properties["laundry"].Hands {
			if n := x.NPC(who); n != nil {
				gone += 100 - n.Trust
			}
		}
		return gone
	}
	a, b := fell(mild), fell(harsh)
	t.Logf("a day a dollar under the rate costs %d; a day at the floor costs %d", a, b)
	if a <= 0 || b <= a {
		t.Fatalf("the floor and a dollar under are worth the same: %d against %d", a, b)
	}
}

// And the rate itself costs nothing and buys nothing, which is what a going
// rate means.
func TestPayingTheRateIsWorthNothingEitherWay(t *testing.T) {
	t.Parallel()
	if trust := trustAfter(t, 61, 0, 60); trust != 0 {
		t.Fatalf("paying exactly the going rate for sixty days is worth %d", trust)
	}
}
