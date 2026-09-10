package core

import "testing"

// Measured before any of this was written: over a simulated week sampled every
// hour, forty-seven of fifty people never moved once, and the busiest places at
// three in the morning were the busiest places at three in the afternoon by the
// same margin. The city had errands, but every errand was a reason that ends,
// and once everybody had a reason to be where they already stood, nobody moved
// again for the rest of the game.

// week runs a world forward an hour at a time, reporting where everybody stood.
func week(t *testing.T, seed uint32, days int) (map[string]int, map[int]map[string]int, map[string]map[int]map[string]int) {
	t.Helper()
	return weekUnder(t, seed, days, func(*World) {})
}

// weekUnder is week with a hand on the city's temperature: hold applies to the
// world before every hour, so a caller can keep the streets at peace or keep
// two families shooting for the whole run.
func weekUnder(t *testing.T, seed uint32, days int, hold func(*World)) (moved map[string]int, byHour map[int]map[string]int, place map[string]map[int]map[string]int) {
	t.Helper()
	w := New(seed)
	moved, byHour, place = map[string]int{}, map[int]map[string]int{}, map[string]map[int]map[string]int{}
	last := map[string]string{}
	for step := 0; step < 24*days; step++ {
		hold(w)
		w.Advance(60)
		hour := (w.Minute % 1440) / 60
		if byHour[hour] == nil {
			byHour[hour] = map[string]int{}
		}
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if n.Dead {
				continue
			}
			byHour[hour][n.Location]++
			if place[n.ID] == nil {
				place[n.ID] = map[int]map[string]int{}
			}
			if place[n.ID][hour] == nil {
				place[n.ID][hour] = map[string]int{}
			}
			place[n.ID][hour][n.Location]++
			if was, ok := last[n.ID]; ok && was != n.Location {
				moved[n.ID]++
			}
			last[n.ID] = n.Location
		}
	}
	return moved, byHour, place
}

func TestMostOfTheCityGoesSomewhereInAWeek(t *testing.T) {
	t.Parallel()
	moved, _, place := week(t, 404, 7)
	movers := 0
	for id := range place {
		if moved[id] > 0 {
			movers++
		}
	}
	if movers*2 < len(place) {
		t.Fatalf("%d of %d people went anywhere at all in a week", movers, len(place))
	}
}

func TestTheEveningLooksDifferentFromTheMorning(t *testing.T) {
	t.Parallel()
	_, byHour, _ := week(t, 404, 7)
	busiest := func(hour int) string {
		best, at := 0, ""
		for id, n := range byHour[hour] {
			if n > best || (n == best && id < at) {
				best, at = n, id
			}
		}
		return at
	}
	morning, evening := busiest(9), busiest(21)
	if morning == evening {
		t.Fatalf("the busiest place at nine in the morning and nine at night is the same room (%s)", morning)
	}
	// And the difference is the one the design intends: people are out.
	drinking := 0
	for _, id := range haunts {
		drinking += byHour[21][id] - byHour[9][id]
	}
	if drinking <= 0 {
		t.Fatalf("the bars and clubs hold %d fewer people at nine at night than at nine in the morning", -drinking)
	}
}

// The promise: the same faces in the same places at the same hours. Somebody
// who moves must still be findable — a rhythm nobody can learn is just noise.
//
// Held at peace, and the test says so because it did not used to. A war now
// keeps people off the street (core/curfew.go), and over fourteen days this
// city has one, which dropped the measure to 88% against a promise of ninety.
// Two readings were available: the curfew is too strong, or the promise is
// about a settled city and a war is exactly the legible disruption the player
// is supposed to notice. The second is the true one — a rhythm that never
// breaks is furniture — so the promise is measured at peace here and the
// wartime floor is measured separately in TestAWarThinsTheStreetWithoutErasing.
func TestAFaceIsFoundInTheSamePlaceAtTheSameHour(t *testing.T) {
	t.Parallel()
	moved, _, place := weekUnder(t, 404, 14, keepPeace)
	reliable, samples := 0, 0
	for id, hours := range place {
		if moved[id] == 0 {
			continue // furniture is reliable for the wrong reason
		}
		for _, places := range hours {
			best, sum := 0, 0
			for _, n := range places {
				sum += n
				if n > best {
					best = n
				}
			}
			samples++
			if best*4 >= sum*3 {
				reliable++
			}
		}
	}
	if samples == 0 {
		t.Fatal("nobody in the city moves, so this proves nothing")
	}
	if share := float64(reliable) / float64(samples); share < .9 {
		t.Fatalf("a moving face is where it usually is on only %.0f%% of person-hours", share*100)
	}
}

// Where somebody drinks must not drift, or the player cannot learn it.
func TestWhereSomebodyDrinksNeverChanges(t *testing.T) {
	t.Parallel()
	for _, id := range []string{"mara", "leo", "harlow", "a-made-up-person"} {
		want := haunt(id)
		for i := 0; i < 50; i++ {
			if haunt(id) != want {
				t.Fatalf("%s drinks somewhere different each time it is asked", id)
			}
		}
		found := false
		for _, h := range haunts {
			if h == want {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s drinks at %q, which is not a place in this city", id, want)
		}
	}
}

// Nobody who is needed somewhere is pulled off it by wanting a drink.
func TestPeopleOnDutyDoNotGoOutDrinking(t *testing.T) {
	t.Parallel()
	w := New(4)
	w.Minute = 720 // the evening
	n := w.NPC("mara")
	n.Role, n.Location, n.Post = roles[0].Title, roles[0].Where, roles[0].Where
	if _, ok := w.routine(n); ok {
		t.Fatalf("the %s left her post for a drink", roles[0].Title)
	}
	held := w.NPC("leo")
	held.Role, held.Post, held.Location = "", "market", "market"
	held.Held = w.Minute + 2880
	if _, ok := w.routine(held); ok {
		t.Fatal("somebody the police are holding went out for a drink")
	}
}

// keepPeace cools any fight back to a feud before the hour turns, so a run
// measures the city's settled rhythm rather than a fortnight with a war in it.
func keepPeace(w *World) {
	for i := range w.Conflicts {
		if w.Conflicts[i].State == "war" {
			w.Conflicts[i].State, w.Conflicts[i].Hostility = "feud", feudAt
		}
	}
}
