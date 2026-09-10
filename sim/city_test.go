package sim

import "testing"

// A city left alone for a season must go on being a city.
//
// This is the standard docs/LIVING_WORLD.md sets and the campaign runner cannot
// meet: it follows one protagonist for a median of a few days, which is why
// every report the harness has printed says no organization has ever fallen.
// They fall. Nobody was watching long enough to see it.

func TestACityLeftAloneKeepsMoving(t *testing.T) {
	const cities, days = 12, 60
	fell, formed, wars, changed := 0, 0, 0, 0
	for seed := uint32(1); seed <= cities; seed++ {
		r := City(seed*2654435761, days)
		fell += r.Fell
		formed += r.Formed
		wars += r.Wars
		changed += r.Changed
	}
	t.Logf("%d cities over %d days: %d organizations fell, %d formed, %d wars started, %d holdings changed hands",
		cities, days, fell, formed, wars, changed)
	if fell == 0 {
		t.Fatal("no organization in any city ever ended, which is the layer this measure exists for")
	}
	if formed == 0 {
		t.Fatal("nobody ever split off and made something of their own")
	}
	if changed == 0 {
		t.Fatal("no property changed hands in a season anywhere")
	}
}

func TestACityDoesNotCollapseIntoOneOwner(t *testing.T) {
	const cities, days = 12, 60
	worst, empty, dead := 0, 0, 0
	for seed := uint32(1); seed <= cities; seed++ {
		r := City(seed*2654435761, days)
		if r.Biggest > worst {
			worst = r.Biggest
		}
		if r.Ended == 0 {
			empty++
		}
		if r.Living == 0 {
			dead++
		}
	}
	t.Logf("after %d days the largest holding in any of %d cities is %d%% of it; %d cities ended with no organizations and %d with nobody alive",
		days, cities, worst, empty, dead)
	if worst >= 90 {
		t.Fatalf("a city ended with %d%% of it in one pair of hands", worst)
	}
	if empty > 0 {
		t.Fatalf("%d cities ended with no organization in them at all", empty)
	}
	if dead > 0 {
		t.Fatalf("%d cities ended with nobody alive in them", dead)
	}
}
