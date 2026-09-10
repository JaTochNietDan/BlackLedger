package core

import "testing"

// A still should pay well for somebody who moves the stock and manages the
// attention, and punish somebody who lets it pile up. It must never be free
// money, because it is the highest-yield thing in the game.
func TestAStillPaysForMovingTheStockNotForOwningIt(t *testing.T) {
	t.Parallel()
	const campaigns, days = 120, 40

	type outcome struct{ cash, heat, raided, seized int }
	run := func(still, sell bool) outcome {
		var total outcome
		for i := uint32(1); i <= campaigns; i++ {
			w := New(i * 2654435761)
			w.Properties["laundry"].Owner = "player:1"
			w.Player.Cash = 3000
			w.Player.Location = "laundry"
			if still {
				if err := w.BuildStill("laundry"); err != nil {
					t.Fatal(err)
				}
			}
			start := w.Player.Cash
			for day := 0; day < days; day++ {
				w.Advance(1440)
				if w.Event != nil {
					w.Event = nil
				}
				if !w.Player.Alive {
					break
				}
				w.Player.Location = "laundry"
				if w.RestockReadiness("laundry") == "" && w.Properties["laundry"].Supply < 15 {
					_ = w.Restock("laundry")
				}
				if sell && w.Holding("moonshine") > 0 {
					// Move it: the whole point is that stock has to go.
					w.Player.Location = "market"
					_ = w.Sell("moonshine", 0)
					w.Player.Location = "laundry"
				}
			}
			total.cash += w.Player.Cash - start
			total.heat += w.Player.Heat
			if !w.Properties["laundry"].Still && still {
				total.seized++
			}
		}
		return outcome{total.cash / campaigns, total.heat / campaigns, total.raided, total.seized}
	}

	none := run(false, false)
	hoarding := run(true, false)
	moving := run(true, true)
	for name, r := range map[string]outcome{"no still": none, "still, hoarding": hoarding, "still, selling": moving} {
		t.Logf("%-16s earned %6d, heat %3d, still lost in %d of %d", name, r.cash, r.heat, r.seized, campaigns)
	}

	if moving.cash <= none.cash {
		t.Fatalf("running a still and moving the stock pays no better than not having one: %d against %d", moving.cash, none.cash)
	}
	if moving.cash <= hoarding.cash {
		t.Fatalf("moving the stock is worth no more than letting it pile up: %d against %d", moving.cash, hoarding.cash)
	}
	if hoarding.heat <= none.heat {
		t.Fatal("letting stock pile up draws no more attention than having no still")
	}
	if hoarding.seized == 0 {
		t.Fatal("a hoarder never loses the still, so the risk is not real")
	}
}
