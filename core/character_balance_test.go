package core

import "testing"

// A temperament has to change outcomes, not just labels. These put the same
// grievance in the hands of different kinds of person and measure what happens
// to them.

func TestWhoYouAreDecidesHowItEnds(t *testing.T) {
	heavy(t)
	const runs, days = 400, 20
	type result struct{ acted, won, died int }
	results := map[string]*result{}

	for _, id := range []string{"hot", "careful", "greedy", "vain"} {
		r := &result{}
		results[id] = r
		for seed := uint32(1); seed <= runs; seed++ {
			w, a, b := quarrel(t)
			w.WorldRNG = seed * 2654435761
			a.Name = nameWith(t, id)
			a.Ambition, a.Skill = 70, 60
			b.Skill = 60
			w.Resent(a.ID, b.ID, GrudgeCap, "an old debt")
			for day := 0; day < days && !settledBetween(w, a, b); day++ {
				w.SettleGrudges()
			}
			if settledBetween(w, a, b) {
				r.acted++
			}
			if b.Dead {
				r.won++
			}
			if a.Dead {
				r.died++
			}
		}
	}

	hot, careful := results["hot"], results["careful"]
	if hot.acted <= careful.acted {
		t.Fatalf("hot-headed people acted %d times against %d careful, inside %d days", hot.acted, careful.acted, days)
	}
	// The careful ones act less often and get more out of it when they do.
	hotRate := float64(hot.won) / float64(max(1, hot.acted))
	carefulRate := float64(careful.won) / float64(max(1, careful.acted))
	if carefulRate <= hotRate {
		t.Fatalf("careful people succeeded %.0f%% of the time they moved against %.0f%% for the hot-headed", carefulRate*100, hotRate*100)
	}
	for id, r := range results {
		t.Logf("%-8s acted in %d of %d cities inside %d days, killed their mark %d times, died trying %d",
			id, r.acted, runs, days, r.won, r.died)
	}
}

func TestGraspingPeopleTakeMore(t *testing.T) {
	heavy(t)
	const runs = 300
	take := map[string]int{}
	for _, id := range []string{"greedy", "careful"} {
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(seed)
			w.MigrateLivingWorld()
			w.WorldRNG = seed * 2654435761
			w.Properties["laundry"].Owner = "independent"
			w.Properties["laundry"].Income = 40
			var thief *NPC
			for i := range w.NPCs {
				if n := &w.NPCs[i]; !n.Dead && n.Rank < RankLeader {
					thief = n
					break
				}
			}
			thief.Name, thief.Skill = nameWith(t, id), 100
			before := w.faction("bellandi").Cash + w.faction("russo").Cash
			w.takeFromSomebody(thief)
			take[id] += max(0, before-(w.faction("bellandi").Cash+w.faction("russo").Cash))
		}
	}
	if take["greedy"] <= take["careful"] {
		t.Fatalf("grasping people took $%d against $%d for careful ones across %d robberies", take["greedy"], take["careful"], runs)
	}
	t.Logf("across %d robberies each: grasping people took $%d, careful ones took $%d", runs, take["greedy"], take["careful"])
}
