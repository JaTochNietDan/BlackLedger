package core

import "testing"

// A journey is the most common thing a player does, so what happens on one has
// to be rare enough to stay an event over a whole campaign.

func TestCrossingTheCityIsMostlyUneventful(t *testing.T) {
	const runs, journeys = 120, 60
	sightings, injuries, deaths := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG, w.RNG = seed*2654435761, seed*2654435761
		w.Player.Health = 100
		// A city at war, which is the worst case rather than the usual one.
		w.Conflicts[0].State, w.Conflicts[0].Hostility = "war", 80
		before := len(w.History)
		for i := 0; i < journeys && w.Player.Alive; i++ {
			health := w.Player.Health
			w.PassThrough("bar", "docks")
			if w.Player.Health < health {
				injuries++
			}
			w.Player.Health = 100 // rested between journeys, so this counts events not attrition
		}
		sightings += len(w.History) - before
		if !w.Player.Alive {
			deaths++
		}
	}
	rate := float64(sightings) / float64(runs*journeys)
	if rate == 0 {
		t.Fatal("nothing ever happened on any journey in a city at war")
	}
	if rate > .25 {
		t.Fatalf("%.0f%% of journeys were an incident, which is a shooting gallery", rate*100)
	}
	t.Logf("%d journeys across %d cities at war: %d incidents (%.0f%% of journeys), %d of them injuring the player, %d killing them",
		runs*journeys, runs, sightings, rate*100, injuries, deaths)
}
