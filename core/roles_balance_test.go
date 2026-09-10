package core

import "testing"

// The point is not that these three die often — with the city's pace scaled to
// its population, any one person acts rarely, and singling them out would be
// the opposite of what this is for. The point is that nothing protects them,
// and that when the city does reach them it keeps speaking afterwards.

func TestACityKeepsTalkingAfterItLosesItsTalkers(t *testing.T) {
	heavy(t)
	const runs, days = 120, 60
	recovered, gaps, sameFace := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		started := map[string]string{}
		for _, r := range roles {
			started[r.ID] = w.Holder(r.ID).ID
		}

		// Something the city does anyway, at the place two of them stand: a
		// charge under Saint Agnes.
		w.Properties["bar"].Condition = 100
		for _, r := range roles {
			if h := w.Holder(r.ID); h != nil && h.Location == "bar" {
				w.KillBy(h.ID, nil, "A charge went off under it.")
			}
		}

		for day := 0; day < days; day++ {
			w.FactionTurn()
			w.PeopleDay()
			w.GrudgeDay()
			w.SettleGrudges()
			w.RecruitDay()
			w.FillRoles()
			w.PrunePeople()
			w.Minute += 1440
		}

		for _, r := range roles {
			holder := w.Holder(r.ID)
			if holder == nil {
				gaps++
				continue
			}
			if holder.ID != started[r.ID] {
				recovered++
			} else {
				sameFace++
			}
		}
	}
	if gaps > 0 {
		t.Fatalf("%d jobs were left with nobody doing them", gaps)
	}
	if recovered == 0 {
		t.Fatal("nobody was ever replaced, so these three are still fixed")
	}
	if sameFace == 0 {
		t.Fatal("everybody was replaced, so the opening cast means nothing")
	}
	t.Logf("%d cities: %d jobs passed to somebody new after a charge at Saint Agnes, %d stayed with the person who started with them, and none was left unfilled",
		runs, recovered, sameFace)
}

// And nothing about being replaced is remembered: a stranger doing the fixer's
// job does not inherit what the last one thought of the player.
func TestASuccessorStartsAsAStranger(t *testing.T) {
	heavy(t)
	w := New(179)
	w.MigrateLivingWorld()
	before := w.Holder("fixer")
	before.Trust = 40
	w.KillBy(before.ID, nil, "Shot.")
	w.FillRoles()
	after := w.Holder("fixer")
	if after == nil || after.ID == before.ID {
		t.Fatal("nobody took the job")
	}
	if after.Trust != 0 {
		t.Fatalf("the new fixer arrived trusting the player %d", after.Trust)
	}
	if w.Known(after) && w.Reach() < 3 {
		t.Fatal("a stranger was known on arrival")
	}
}
