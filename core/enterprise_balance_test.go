package core

import "testing"

// Skimming is meant to be a bargain for someone who handles the consequences,
// not free money and not a trap. Three players are measured over the same
// campaigns: one runs clean and does nothing else, one skims and ignores the
// attention it brings, and one skims and launders it away. The middle one is
// supposed to lose.
func TestSkimmingPaysOnlyIfTheAttentionIsManaged(t *testing.T) {
	heavy(t)
	const campaigns, days = 120, 40

	type outcome struct{ cash, heat, condition, raids, kept int }
	run := func(mode string, manage bool) outcome {
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
			raids := 0
			for day := 0; day < days; day++ {
				w.Advance(1440)
				if w.Event != nil {
					w.Event = nil // measuring standing pressure, not choices
				}
				if !w.Player.Alive {
					break
				}
				for _, r := range w.History {
					if r.Minute > w.Minute-1440 && (r.Title == "Turned over at Bluebird Laundry" || r.Title == "They took it") {
						raids++
					}
				}
				// A player who handles it: clear attention through the books
				// and keep the premises standing.
				if manage {
					w.Player.Location = "laundry"
					if w.LaunderReadiness("laundry") == "" {
						_ = w.Launder("laundry")
					}
					if w.Properties["laundry"].Condition < 70 && w.Player.Cash > 50 {
						w.Player.Cash -= 50
						w.Properties["laundry"].Condition = min(100, w.Properties["laundry"].Condition+40)
					}
					// Managing a business now means running it as well as
					// laundering through it: keeping it stocked, keeping it
					// staffed, and dealing with what goes wrong inside it.
					for _, id := range []string{"laundry", "garage"} {
						w.Player.Location = id
						if w.RestockReadiness(id) == "" && w.Properties[id].Supply < 15 {
							_ = w.Restock(id)
						}
						if w.RemedyReadiness(id) == "" {
							_ = w.Remedy(id)
						}
						if w.HireReadiness(id) == "" {
							_ = w.Hire(id)
						}
					}
				}
			}
			total.cash += w.Player.Cash - start
			total.heat += w.Player.Heat
			total.condition += w.Properties["laundry"].Condition
			total.raids += raids
			if w.Own("laundry") {
				total.kept++
			}
		}
		return outcome{total.cash / campaigns, total.heat / campaigns,
			total.condition / campaigns, total.raids, total.kept}
	}

	clean := run("clean", false)
	careless := run("hard", false)
	careful := run("hard", true)
	for name, r := range map[string]outcome{"clean": clean, "hard, unmanaged": careless, "hard, managed": careful} {
		t.Logf("%-16s earned %6d, heat %3d, laundry condition %3d, kept the business in %d of %d campaigns",
			name, r.cash, r.heat, r.condition, r.kept, campaigns)
	}

	if careful.cash <= clean.cash {
		t.Fatalf("skimming and handling the consequences pays no better than running clean: %d against %d", careful.cash, clean.cash)
	}
	if careless.cash >= careful.cash {
		t.Fatalf("ignoring the attention costs nothing: careless %d, careful %d", careless.cash, careful.cash)
	}
	if careless.heat <= clean.heat {
		t.Fatal("skimming draws no more attention than running clean")
	}
	if careless.kept >= campaigns {
		t.Fatal("a careless skimmer never loses a business, so forfeiture is not a real risk")
	}
	if careful.kept <= careless.kept {
		t.Fatal("handling the attention does not protect the business")
	}
}
