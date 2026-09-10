package core

import "testing"

// The point of private history is that it becomes public. These measure a city
// left entirely to itself: nobody plays it, and it still produces killings,
// quarrels and occasionally a war nobody was told the reason for.

func TestACityLeftAloneKillsItsOwnPeople(t *testing.T) {
	heavy(t)
	const runs, days = 200, 120
	killings, wars, quiet := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		wars0 := 0
		for _, c := range w.Conflicts {
			if c.State == "war" {
				wars0++
			}
		}
		dead := 0
		for day := 0; day < days; day++ {
			w.FactionTurn()
			w.FamilyDay()
			w.PeopleDay()
			w.GrudgeDay()
			w.SettleGrudges()
			w.Minute += 1440
		}
		for i := range w.NPCs {
			if w.NPCs[i].Dead {
				dead++
			}
		}
		killings += dead
		if dead == 0 {
			quiet++
		}
		for _, c := range w.Conflicts {
			if c.State == "war" {
				wars++
			}
		}
		_ = wars0
	}
	if killings == 0 {
		t.Fatal("a city left alone for four months killed nobody")
	}
	if quiet == 0 {
		t.Fatal("every city killed somebody, so this is a certainty rather than a risk")
	}
	t.Logf("%d cities left alone for %d days: %d people dead across them, %d cities where nobody died, %d live wars at the end", runs, days, killings, quiet, wars)
}

// TestGrudgesComeFromThingsThatHappened checks the source rather than the
// outcome: nobody should be resenting anybody in a city where nothing has
// happened yet.
func TestGrudgesComeFromThingsThatHappened(t *testing.T) {
	heavy(t)
	fresh := New(5)
	fresh.MigrateLivingWorld()
	if len(fresh.Grudges) != 0 {
		t.Fatalf("a city on its first morning already had %d grievances in it", len(fresh.Grudges))
	}

	const runs, days = 200, 60
	withHistory := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		for day := 0; day < days; day++ {
			w.PeopleDay()
			w.GrudgeDay()
			w.Minute += 1440
		}
		if len(w.Grudges) > 0 {
			withHistory++
		}
	}
	if withHistory == 0 {
		t.Fatalf("no city in %d produced a single grievance in %d days of people going about their business", runs, days)
	}
	t.Logf("%d of %d cities carried at least one live grievance after %d days of nobody but its own people acting", withHistory, runs, days)
}

// TestTheSaveNeverRunsAway is the bound that matters: a hundred and twenty days
// of a city resenting itself has to still fit in a save file.
func TestTheSaveNeverRunsAway(t *testing.T) {
	heavy(t)
	worst := 0
	for seed := uint32(1); seed <= 100; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		for day := 0; day < 120; day++ {
			w.FactionTurn()
			w.PeopleDay()
			w.GrudgeDay()
			w.SettleGrudges()
			w.Minute += 1440
		}
		worst = max(worst, len(w.Grudges))
	}
	if worst > MaxGrudges {
		t.Fatalf("a city carried %d grievances, past the %d cap", worst, MaxGrudges)
	}
	t.Logf("the most grievances any of 100 cities carried after 120 days: %d of a possible %d", worst, MaxGrudges)
}
