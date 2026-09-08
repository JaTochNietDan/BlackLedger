package core

import (
	"fmt"
	"testing"
)

// A laundry now has an outside as well as an inside, and the two pull against
// each other: the books are the reason to own it and using them is what costs
// it its trade. This measures which way that comes out over a campaign.

func TestTheBooksAndTheShopPullAgainstEachOther(t *testing.T) {
	const runs, days = 120, 60
	shop := func(launderEvery int) (int, int) {
		earned, custom := 0, 0
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(seed)
			w.MigrateLivingWorld()
			w.WorldRNG = seed * 2654435761
			w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
			w.Player.Location, w.Player.Cash, w.Player.Earned = "laundry", 60000, 0
			for day := 0; day < days; day++ {
				w.Player.Heat = 20 // something worth clearing, held steady
				if launderEvery > 0 && day%launderEvery == 0 {
					w.Player.LastLaunder = 0
					_ = w.Launder("laundry")
				}
				w.CustomDay()
				w.Advance(1440)
			}
			earned += w.Player.Earned
			custom += w.Custom("laundry")
		}
		return earned / runs, custom / runs
	}

	quietEarned, quietCustom := shop(0)
	washedEarned, washedCustom := shop(4)
	if washedCustom >= quietCustom {
		t.Fatalf("using the books left trade at %d%% against %d%% never using them", washedCustom, quietCustom)
	}
	if quietEarned <= 0 {
		t.Fatal("a laundry left alone earned nothing")
	}
	t.Logf("%d campaigns of %d days: never using the books leaves trade at %d%% and $%d earned; a round every four days leaves %d%% and $%d",
		runs, days, quietCustom, quietEarned, washedCustom, washedEarned)
}

// And the standing order is the reward for the other way of playing it: a shop
// kept at capacity earns more than one merely owned.
func TestAShopWorthRelyingOnEarnsMore(t *testing.T) {
	const days = 90
	run := func(mind bool) int {
		w := New(7)
		w.MigrateLivingWorld()
		w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
		w.Player.Location, w.Player.Cash, w.Player.Earned = "laundry", 60000, 0
		trade, _ := TradeOf("laundry")
		for day := 0; day < days; day++ {
			if mind {
				prop := w.Properties["laundry"]
				prop.Supply, prop.Staff = trade.RestockAmount, trade.Hands
				prop.Trouble = false
				if w.OrderReadiness("laundry") == "" {
					_ = w.TakeOrder("laundry")
				}
			}
			w.CustomDay()
			w.OrderDay()
			w.Advance(1440)
		}
		return w.Player.Earned
	}
	minded, left := run(true), run(false)
	if minded <= left {
		t.Fatalf("minding the shop earned $%d against $%d leaving it alone", minded, left)
	}
	t.Logf("over %d days: a laundry left alone earns $%d, one kept at capacity with a standing order earns $%d", days, left, minded)
}
