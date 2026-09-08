package core

import "testing"

// Skimming has to be a real bargain rather than an obviously correct choice or
// a trap nobody would take. Measured over many campaigns rather than argued.
func TestSkimmingIsABargainNotAFreeLunch(t *testing.T) {
	const campaigns, days = 120, 40
	type outcome struct{ cash, heat, condition, demands int }
	results := map[string]outcome{}

	for _, mode := range []string{"clean", "standard", "hard"} {
		var total outcome
		for seed := uint32(1); seed <= campaigns; seed++ {
			w := New(seed)
			w.Properties["laundry"].Owner = "player:1"
			w.Properties["garage"].Owner = "player:1"
			for _, id := range []string{"laundry", "garage"} {
				if mode != "standard" {
					if err := w.SetMode(id, mode); err != nil {
						t.Fatal(err)
					}
				}
			}
			start := w.Player.Cash
			for day := 0; day < days; day++ {
				w.Advance(1440)
				if w.Event != nil {
					if w.Event.Kind == "business_pressure" {
						total.demands++
					}
					w.Event = nil // the measurement is of standing pressure, not of choices
				}
				if !w.Player.Alive {
					break
				}
			}
			total.cash += w.Player.Cash - start
			total.heat += w.Player.Heat
			total.condition += w.Properties["laundry"].Condition
		}
		results[mode] = outcome{total.cash / campaigns, total.heat / campaigns,
			total.condition / campaigns, total.demands}
	}

	for _, mode := range []string{"clean", "standard", "hard"} {
		r := results[mode]
		t.Logf("%-8s median-ish over %d campaigns of %d days: earned %d, heat %d, laundry condition %d, demands %d",
			mode, campaigns, days, r.cash, r.heat, r.condition, r.demands)
	}

	if results["hard"].cash <= results["clean"].cash {
		t.Fatal("skimming earns no more than running clean; nobody would ever take the risk")
	}
	if results["hard"].heat <= results["clean"].heat {
		t.Fatal("skimming costs no more attention than running clean")
	}
	if results["hard"].condition >= results["clean"].condition {
		t.Fatal("skimming does not wear premises faster than running clean")
	}
	// Family attention is asserted directly and deterministically in
	// TestSkimmingBringsAFamilyDemandSooner; demands are only logged here.
	// The upside has to be worth something, or the decision is only a penalty.
	if results["hard"].cash < results["clean"].cash*5/4 {
		t.Fatalf("skimming pays too little to be worth its costs: %d vs %d",
			results["hard"].cash, results["clean"].cash)
	}
}
