package core

import (
	"fmt"
	"testing"
)

// The point of a float is that it is a decision with money on both sides. These
// measure the three ways somebody can own a room — leaving the tables dark,
// funding them thinly, and funding them properly — rather than asserting that
// one is better.

func room(seed uint32) *World {
	w := New(seed)
	w.Properties["casino"].Owner = fmt.Sprintf("player:%d", w.Life)
	trade, _ := TradeOf("casino")
	w.Properties["casino"].Staff, w.Properties["casino"].Supply = trade.Hands, trade.RestockAmount
	w.Player.Cash = 6000
	return w
}

// keepFloat runs a month at a room, topping the float up to a target each day
// and drawing off anything above it, which is how an owner actually banks what
// the tables make. Reports what reached the player's hands and whether the
// house was ever unable to pay.
func keepFloat(seed uint32, target, days int) (int, bool) {
	w := room(seed)
	w.WorldRNG = seed * 2654435761
	banked := 0
	for day := 0; day < days; day++ {
		for w.Properties["casino"].Bankroll+BankrollLot <= target && w.Player.Cash >= BankrollLot {
			if err := w.Bankroll("casino", 0); err != nil {
				break
			}
		}
		w.CasinoDay()
		for w.Properties["casino"].Bankroll-BankrollLot >= target {
			before := w.Player.Cash
			if err := w.Draw("casino", 0); err != nil {
				break
			}
			banked += w.Player.Cash - before
		}
	}
	// Whatever is still behind the tables is the owner's too, if they live.
	return banked + w.Properties["casino"].Bankroll - target, w.hasRecord("The tables could not cover it at The Blue Hour")
}

func TestFundingTheTablesIsADecisionWithMoneyOnBothSides(t *testing.T) {
	heavy(t)
	const runs, days = 120, 40
	dark, thin, deep := 0, 0, 0
	thinBust, deepBust := 0, 0
	for i := uint32(1); i <= runs; i++ {
		w := room(i)
		w.WorldRNG = i * 2654435761
		for day := 0; day < days; day++ {
			w.CasinoDay()
		}
		dark += w.Properties["casino"].Bankroll

		earned, bust := keepFloat(i, 500, days)
		thin += earned
		if bust {
			thinBust++
		}
		earned, bust = keepFloat(i, BankrollFull*2, days)
		deep += earned
		if bust {
			deepBust++
		}
	}
	if dark != 0 {
		t.Fatalf("a room with no float made $%d out of nothing", dark)
	}
	if thin >= deep {
		t.Fatalf("a $500 float earned $%d against $%d for a $3,000 one, so funding a room buys nothing", thin/runs, deep/runs)
	}
	if thinBust <= deepBust {
		t.Fatalf("thin floats went bust %d times against %d deep, so under-funding is free", thinBust, deepBust)
	}
	t.Logf("over %d days across %d campaigns: dark tables earn nothing; a $500 float earns $%d and cannot pay in %d; a $3,000 float earns $%d and cannot pay in %d",
		days, runs, thin/runs, thinBust, deep/runs, deepBust)
}

// TestARoomIsWorthLessThanItWasAndMoreThanItCosts records the deliberate
// rebalance: the casino's hourly number used to include the tables, and now it
// does not. An owner who never touches the float earns strictly less than
// before; one who runs it properly earns most of it back and carries the risk.
func TestARoomIsWorthLessThanItWasAndMoreThanItCosts(t *testing.T) {
	heavy(t)
	const days = 40
	floorTake := New(1).Properties["casino"].Income * 24 * days
	deep, _ := keepFloat(9, BankrollFull*2, days)
	if floorTake <= 0 {
		t.Fatal("the room stopped earning at all")
	}
	if deep <= floorTake/4 {
		t.Fatalf("the tables added $%d against $%d of floor take, which is not worth the risk", deep, floorTake)
	}
	t.Logf("over %d days: floor take $%d, tables $%d on top", days, floorTake, deep)
}
