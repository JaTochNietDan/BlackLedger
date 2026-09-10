package core

import "testing"

// A trader who buys blind and sells blind should roughly break even: the price
// moves around a base, so profit has to come from judgement, not from the
// existence of a market. A trader who waits for a good price should do better,
// and holding stock should cost attention either way.
func TestTheTradeRewardsJudgementRatherThanExistence(t *testing.T) {
	heavy(t)
	const campaigns, cycles = 150, 30

	blind, patient := 0, 0
	heatHeld := 0
	for seed := uint32(1); seed <= campaigns; seed++ {
		// Buys and sells on a fixed rhythm, ignoring the price entirely.
		w := New(seed)
		w.Player.Location = "market"
		w.Player.Cash = 100000
		start := w.Player.Cash
		for i := 0; i < cycles; i++ {
			_ = w.Buy("moonshine", 0)
			w.MarketPrices()
			_ = w.Sell("moonshine", 0)
			w.MarketPrices()
		}
		blind += w.Player.Cash - start

		// Buys under the base and sells over it, which is the whole skill.
		p := New(seed)
		p.Player.Location = "market"
		p.Player.Cash = 100000
		pStart := p.Player.Cash
		for i := 0; i < cycles*2; i++ {
			g := p.Good("moonshine")
			switch {
			case p.Holding("moonshine") == 0 && g.Price <= g.Base*9/10:
				_ = p.Buy("moonshine", 0)
			case p.Holding("moonshine") > 0 && g.Price >= g.Base*11/10:
				_ = p.Sell("moonshine", 0)
			}
			p.MarketPrices()
			if p.Carrying() > 0 {
				heatHeld++
			}
		}
		_ = p.Sell("moonshine", 0)
		patient += p.Player.Cash - pStart
	}
	blind /= campaigns
	patient /= campaigns
	t.Logf("over %d campaigns: buying and selling blind nets %d, waiting for a price nets %d, stock held on %d ticks",
		campaigns, blind, patient, heatHeld)

	if patient <= blind {
		t.Fatalf("waiting for a good price is not rewarded: patient %d against blind %d", patient, blind)
	}
	if blind > 2000 {
		t.Fatalf("trading blind is reliably profitable (%d); the market is free money", blind)
	}
	if patient <= 0 {
		t.Fatalf("even a disciplined trader loses money (%d); nobody would use the market", patient)
	}
	if heatHeld == 0 {
		t.Fatal("stock was never actually held, so the risk was never taken")
	}
}
