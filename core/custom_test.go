package core

import (
	"fmt"
	"testing"
)

func shopkeeper(t *testing.T) *World {
	t.Helper()
	w := New(191)
	w.MigrateLivingWorld()
	w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
	w.Player.Location, w.Player.Cash, w.Player.Heat = "laundry", 20000, 30
	return w
}

func TestAPlaceTakenOnHasATradeAlready(t *testing.T) {
	t.Parallel()
	w := shopkeeper(t)
	if w.Custom("laundry") != CustomStart {
		t.Fatalf("a laundry came with %d%% trade", w.Custom("laundry"))
	}
	if w.TradeMultiplier("laundry") != 1 {
		t.Fatalf("half a reputation was worth %.2f", w.TradeMultiplier("laundry"))
	}
	// Somewhere that is not a trading business has no trade to speak of.
	if w.CustomDescription("club") != nil {
		t.Fatal("a casino floor had a laundry's custom")
	}
}

func TestRunningItProperlyBuildsATradeAndNeglectLosesIt(t *testing.T) {
	t.Parallel()
	w := shopkeeper(t)
	for i := 0; i < 30; i++ {
		w.CustomDay()
	}
	built := w.Custom("laundry")
	if built <= CustomStart {
		t.Fatalf("a month of running it properly left trade at %d%%", built)
	}
	if w.TradeMultiplier("laundry") <= 1 {
		t.Fatal("a better reputation earned no more")
	}
	w.Properties["laundry"].Trouble = true
	for i := 0; i < 10; i++ {
		w.CustomDay()
	}
	if w.Custom("laundry") >= built {
		t.Fatalf("ten days of something being wrong left trade at %d%%", w.Custom("laundry"))
	}
	// And it cannot go past the ceiling or below nothing.
	w.Properties["laundry"].Trouble = false
	for i := 0; i < 200; i++ {
		w.CustomDay()
	}
	if w.Custom("laundry") != CustomCeiling {
		t.Fatalf("trade settled at %d%%", w.Custom("laundry"))
	}
	w.Properties["laundry"].Trouble = true
	for i := 0; i < 200; i++ {
		w.CustomDay()
	}
	if w.Custom("laundry") != CustomFloor {
		t.Fatalf("trade bottomed out at %d%%", w.Custom("laundry"))
	}
}

func TestSkimmingHardCostsTheShop(t *testing.T) {
	t.Parallel()
	clean := shopkeeper(t)
	hard := shopkeeper(t)
	hard.Properties["laundry"].Mode = "hard"
	for i := 0; i < 20; i++ {
		clean.CustomDay()
		hard.CustomDay()
	}
	if hard.Custom("laundry") >= clean.Custom("laundry") {
		t.Fatalf("skimming left trade at %d%% against %d%% running it as usual", hard.Custom("laundry"), clean.Custom("laundry"))
	}
}

func TestUsingTheBooksIsWhatTradeIsSpentOn(t *testing.T) {
	t.Parallel()
	w := shopkeeper(t)
	before := w.Custom("laundry")
	if err := w.Launder("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.Custom("laundry") != before-CustomLaunderLoss {
		t.Fatalf("a round through the books cost %d trade", before-w.Custom("laundry"))
	}
}

func TestASearchInDaylightIsTheEndOfAShopsStanding(t *testing.T) {
	t.Parallel()
	w := shopkeeper(t)
	w.Properties["laundry"].Income = 20
	before := w.Custom("laundry")
	w.Player.Heat = 90
	w.Raid()
	if w.Custom("laundry") >= before {
		t.Fatalf("a raid left trade at %d%% against %d%%", w.Custom("laundry"), before)
	}
}

func TestNobodyOffersStandingWorkToAPlaceWithNoStanding(t *testing.T) {
	t.Parallel()
	w := shopkeeper(t)
	if w.OrderReadiness("laundry") == "" {
		t.Fatal("a place at half trade was offered a standing order")
	}
	w.Properties["laundry"].Custom = OrderCustom
	if w.OrderReadiness("laundry") != "" {
		t.Fatal("a place at the stated trade was refused:", w.OrderReadiness("laundry"))
	}
	if err := w.TakeOrder("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.OrderReadiness("laundry") == "" {
		t.Fatal("took a second order on the same books")
	}
	// It pays while it is filled.
	cash := w.Player.Cash
	w.OrderDay()
	if w.Player.Cash != cash+OrderBonus {
		t.Fatalf("a filled order paid $%d", w.Player.Cash-cash)
	}
}

func TestAnOrderThatCannotBeFilledIsLostAndCostsMore(t *testing.T) {
	t.Parallel()
	w := shopkeeper(t)
	w.Properties["laundry"].Custom = 80
	w.TakeOrder("laundry")
	before := w.Custom("laundry")
	w.Properties["laundry"].Staff = 0
	cash := w.Player.Cash
	w.OrderDay()
	if w.Properties["laundry"].Order {
		t.Fatal("an order nobody could fill stayed on the books")
	}
	if w.Player.Cash != cash {
		t.Fatal("an unfilled order still paid")
	}
	if w.Custom("laundry") != before-OrderLoss {
		t.Fatalf("losing it cost %d trade", before-w.Custom("laundry"))
	}
	if !w.hasRecord("The order is cancelled at Bluebird Laundry") {
		t.Fatal("nobody was told")
	}
}

func TestTradeReachesWhatTheBusinessEarns(t *testing.T) {
	t.Parallel()
	const days = 20
	rate := func(custom int) float64 {
		w := shopkeeper(t)
		w.Properties["laundry"].Custom = custom
		w.Player.Cash, w.Player.Earned = 50000, 0
		start := w.Minute
		w.Advance(days * 1440)
		return float64(w.Player.Earned) / float64(max(1, w.Minute-start))
	}
	poor, rich := rate(10), rate(CustomCeiling)
	if rich <= poor {
		t.Fatalf("a laundry with a reputation earned $%.3f a minute against $%.3f without one", rich, poor)
	}
	t.Logf("a laundry earns $%.3f a minute at 10%% trade and $%.3f at 100%%", poor, rich)
}
