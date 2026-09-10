package core

import "testing"

// Work an organization asks for has to be work the city can actually produce,
// and it has to be worth taking. These measure both rather than assuming them.

// TestTheCityAsksForDifferentThingsAtDifferentTimes checks that the objective
// is derived from the asking organization's situation, so the purpose layer is
// not one errand with the names changed.
func TestTheCityAsksForDifferentThingsAtDifferentTimes(t *testing.T) {
	heavy(t)
	const runs = 300
	kinds := map[string]int{}
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed * 2654435761)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		// Let the city run a season on its own, so wars start and finish,
		// premises change hands and organizations get rich or poor.
		for day := 0; day < 60; day++ {
			w.FactionTurn()
			w.FamilyDay()
			w.PeopleDay()
		}
		w.Player.Location = "bar"
		// Ask a different organization each time, so the sample is the city
		// rather than whichever family happens to be listed first.
		wanted, seen := int(seed)%max(1, len(w.Factions)), 0
		asked := false
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if n.Faction == "" || n.Rank < RankLieutenant || n.Dead {
				continue
			}
			if seen != wanted {
				seen++
				continue
			}
			n.Location = "bar"
			asked = true
			break
		}
		if !asked {
			kinds["none"]++
			continue
		}
		if c, ok := w.AvailableCommission("bar"); ok {
			kinds[c.Kind]++
		} else {
			kinds["none"]++
		}
	}
	distinct := 0
	for kind, n := range kinds {
		if kind != "none" && n > 0 {
			distinct++
		}
	}
	if distinct < 3 {
		t.Fatalf("the city only ever asked for one thing: %v", kinds)
	}
	t.Logf("across %d cities run for 60 days, organizations asked for: %v", runs, kinds)
}

// TestADeliveryIsWorthTakingOnlySometimes measures the supply run against the
// market it has to be bought on, so it is a judgement rather than free money.
func TestADeliveryIsWorthTakingOnlySometimes(t *testing.T) {
	heavy(t)
	const runs = 400
	profitable, ruinous, total := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed * 2654435761)
		w.WorldRNG = seed * 2654435761
		for day := 0; day < 40; day++ {
			w.MarketPrices()
		}
		cost := w.Good("moonshine").Price * 25
		total += 1150 - cost
		if 1150 > cost {
			profitable++
		}
		if cost > 1150 {
			ruinous++
		}
	}
	if profitable == 0 || ruinous == 0 {
		t.Fatalf("a supply run paid off in %d of %d campaigns, so the market does not matter", profitable, runs)
	}
	t.Logf("buying 25 crates to fill a $1,150 delivery: worth it in %d of %d campaigns, a loss in %d, averaging $%d", profitable, runs, ruinous, total/runs)
}

// TestFailingWorkCostsMoreThanNeverTakingIt is the point of a deadline: a
// commission has to be a risk, not a free lottery ticket.
func TestFailingWorkCostsMoreThanNeverTakingIt(t *testing.T) {
	heavy(t)
	w := New(11)
	w.MigrateLivingWorld()
	f := &w.Factions[0]
	f.Goodwill = 0
	w.Commissions = append(w.Commissions, Commission{ID: ID(), PatronID: f.ID, GiverName: "Somebody",
		Kind: ObjectiveStanding, Amount: 9999, Due: w.Minute + CommissionWindow, Life: w.Life, Penalty: 10})
	w.Minute += CommissionWindow
	w.SettleCommissions()
	if f.Goodwill >= 0 {
		t.Fatalf("failing left standing at %d", f.Goodwill)
	}
	// And the standing lost is worth something: it is the same number that
	// governs whether that family moves against the player at all.
	if f.Goodwill != -10 {
		t.Fatalf("standing fell to %d rather than the stated penalty", f.Goodwill)
	}
}
