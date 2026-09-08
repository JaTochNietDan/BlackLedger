package core

import (
	"strings"
	"testing"
)

func walker(t *testing.T) *World {
	t.Helper()
	w := New(127)
	w.MigrateLivingWorld()
	w.Player.Location, w.Player.Health = "bar", 100
	for i := range w.Conflicts {
		w.Conflicts[i].State, w.Conflicts[i].Hostility = "cold", 10
	}
	w.Grudges = nil
	w.Player.Heat = 0
	return w
}

func TestAQuietCityHasQuietStreets(t *testing.T) {
	w := walker(t)
	for seed := uint32(1); seed <= 300; seed++ {
		w.WorldRNG = seed * 2654435761
		if _, ok := w.OnTheWay("bar", "docks"); ok {
			t.Fatal("something happened on the way through a city where nothing was happening")
		}
	}
}

func TestAWarReachesTheStreet(t *testing.T) {
	const runs = 400
	seen, stray := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := walker(t)
		w.WorldRNG = seed * 2654435761
		w.Conflicts[0].State, w.Conflicts[0].Hostility = "war", 80
		s, ok := w.OnTheWay("bar", "docks")
		if !ok {
			continue
		}
		seen++
		if !strings.Contains(s.Title, "Shooting") {
			t.Fatalf("a war produced %q", s.Title)
		}
		if s.Stray > 0 {
			stray++
		}
	}
	if seen == 0 {
		t.Fatalf("an open war never reached the street in %d journeys", runs)
	}
	if seen == runs {
		t.Fatal("every single journey during a war was a shooting")
	}
	if stray == 0 {
		t.Fatal("nobody passing a war was ever caught by it")
	}
	if stray == seen {
		t.Fatal("everybody passing a war was caught by it")
	}
	t.Logf("%d journeys during an open war: %d crossed a shooting, and %d of those were close enough to be hit", runs, seen, stray)
}

func TestWhatYouUnderstandDependsOnWhoYouKnow(t *testing.T) {
	quiet := walker(t)
	quiet.Conflicts[0].State, quiet.Conflicts[0].Hostility = "war", 80
	quiet.Player.Contacts = 0
	wired := walker(t)
	wired.Conflicts[0].State, wired.Conflicts[0].Hostility = "war", 80
	wired.Player.Contacts = 3

	for seed := uint32(1); seed <= 400; seed++ {
		quiet.WorldRNG, wired.WorldRNG = seed*2654435761, seed*2654435761
		a, okA := quiet.OnTheWay("bar", "docks")
		b, okB := wired.OnTheWay("bar", "docks")
		if !okA || !okB {
			continue
		}
		if a.Explained {
			t.Fatal("a stranger with no contacts understood what they were looking at")
		}
		if !b.Explained {
			t.Fatal("a player with three contacts was told nothing")
		}
		if a.Body == b.Body {
			t.Fatal("both accounts read the same")
		}
		return
	}
	t.Fatal("no journey in 400 crossed anything")
}

func TestAGrievanceIsSomethingYouCanWalkPast(t *testing.T) {
	const runs = 400
	seen := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := walker(t)
		w.WorldRNG = seed * 2654435761
		people := w.People()
		w.Resent(people[3].ID, people[4].ID, GrudgeCap, "an old debt")
		s, ok := w.OnTheWay("bar", "docks")
		if !ok {
			continue
		}
		seen++
		if !strings.Contains(s.Title, "doorway") {
			t.Fatalf("a grievance produced %q", s.Title)
		}
		if s.Stray != 0 {
			t.Fatal("two men arguing in a doorway shot a passer-by")
		}
	}
	if seen == 0 {
		t.Fatalf("no journey in %d crossed a live grievance", runs)
	}
	t.Logf("%d journeys past a live grievance: %d of them saw it", runs, seen)
}

func TestNothingSmallEnoughToIgnoreIsShown(t *testing.T) {
	w := walker(t)
	people := w.People()
	w.Resent(people[3].ID, people[4].ID, 5, "a slight")
	for seed := uint32(1); seed <= 300; seed++ {
		w.WorldRNG = seed * 2654435761
		if _, ok := w.OnTheWay("bar", "docks"); ok {
			t.Fatal("a five-weight grievance was staged in the street")
		}
	}
}

func TestThePoliceAreWhereThePoliceHaveAReasonToBe(t *testing.T) {
	w := walker(t)
	w.Player.Heat = RaidThreshold / 2
	found := false
	for seed := uint32(1); seed <= 400 && !found; seed++ {
		w.WorldRNG = seed * 2654435761
		if s, ok := w.OnTheWay("bar", "docks"); ok {
			found = true
			if !strings.Contains(s.Title, "watching") {
				t.Fatalf("attention produced %q", s.Title)
			}
		}
	}
	if !found {
		t.Fatal("nobody was ever watching a player at half the raid threshold")
	}
	// And nobody watches a player nobody is interested in.
	cool := walker(t)
	for seed := uint32(1); seed <= 300; seed++ {
		cool.WorldRNG = seed * 2654435761
		if _, ok := cool.OnTheWay("bar", "docks"); ok {
			t.Fatal("the police watched somebody with no attention on them")
		}
	}
}

func TestPassingThroughReachesTheRecordAndCanHurt(t *testing.T) {
	w := walker(t)
	w.Conflicts[0].State, w.Conflicts[0].Hostility = "war", 80
	hurt := false
	for seed := uint32(1); seed <= 400 && !hurt; seed++ {
		probe := walker(t)
		probe.Conflicts[0].State, probe.Conflicts[0].Hostility = "war", 80
		probe.WorldRNG, probe.RNG = seed*2654435761, seed*2654435761
		before := len(probe.History)
		probe.PassThrough("bar", "docks")
		if len(probe.History) > before && probe.Player.Health < 100 {
			hurt = true
			if !probe.hasRecord("Shooting on the way to Pier 14") {
				t.Fatal("nothing was recorded")
			}
		}
	}
	if !hurt {
		t.Fatal("nobody in 400 journeys past a war was ever hurt by it")
	}
}
