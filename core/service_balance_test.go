package core

import "testing"

// Two careers: buy premises and become somebody, or go to work for somebody who
// already is. This measures what the second one is worth against the first over
// the same weeks, starting from the same ninety dollars.

func TestAWageIsASlowerLivingThanAShop(t *testing.T) {
	const runs, days = 150, 45
	measure := func(serve bool) (earned, deaths int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(seed)
			w.MigrateLivingWorld()
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Respect, w.Player.Health = ServiceRespect, 100
			f := &w.Factions[0]
			f.Goodwill, f.Cash = ServiceGoodwill, 60000
			if serve {
				w.Player.Location = w.homeOf(f.ID)
				if err := w.Serve(f.ID); err != nil {
					t.Fatal(err)
				}
			} else {
				// The other career: a shop, run properly.
				w.Properties["laundry"].Owner = w.PlayerOrganizationID()
			}
			start := w.Player.Earned
			for day := 0; day < days && w.Player.Alive; day++ {
				w.ServiceDay()
				w.CustomDay()
				w.Advance(1440)
			}
			earned += w.Player.Earned - start
			if !w.Player.Alive {
				deaths++
			}
		}
		return
	}

	wage, wageDeaths := measure(true)
	shop, shopDeaths := measure(false)
	if wage == 0 {
		t.Fatal("answering to somebody paid nothing at all")
	}
	if wage >= shop {
		t.Fatalf("a wage earned $%d against $%d for a shop, so nobody would ever buy one", wage/runs, shop/runs)
	}
	t.Logf("%d campaigns of %d days from the same start: a wage earns $%d and kills %d; a laundry earns $%d and kills %d",
		runs, days, wage/runs, wageDeaths, shop/runs, shopDeaths)
}

// And coming up is what makes it worth staying: a lieutenant takes a share.
func TestALieutenantTakesAShare(t *testing.T) {
	w, f := recruit(t)
	w.Serve(f.ID)
	bottom := w.ServicePay()
	for i := 0; i < PromotionWork*2; i++ {
		w.ServeWork(f.ID)
	}
	top := w.ServicePay()
	if top <= bottom {
		t.Fatalf("a lieutenant takes %d against %d at the bottom", top, bottom)
	}
	// And the share is of what they actually hold, so a family that loses
	// ground pays their people less.
	for _, id := range w.FamilyHoldings(f.ID) {
		w.Properties[id].Condition = 10
	}
	if w.ServicePay() >= top {
		t.Fatal("a share of a wrecked organization was worth as much as a whole one")
	}
	t.Logf("answering to %s pays $%d a day at the bottom and $%d as a lieutenant", f.Name, bottom, top)
}
