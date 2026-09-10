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
