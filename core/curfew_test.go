package core

import "testing"

// "The player may be a bystander, a beneficiary, or collateral."
//
// Measured: over a hundred campaigns the investor was party to no war, watched
// seven or nine break out elsewhere, and was never touched by any of them —
// lowest health a hundred, nothing taken. A city can be at war around somebody
// who owns half of it and cost them nothing at all.
//
// What a war does to people who are not in it is keep them at home. Now that a
// room's takings are who is standing in it, that is a cost the player feels
// without anybody having to come for them.

func townsman(t *testing.T) *World {
	t.Helper()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 5000, 100
	own(w, "club")
	// A room with nothing behind the tables covers nothing whoever is in it.
	w.Properties["club"].Bankroll = 3000
	w.Properties["club"].Condition = 100
	for day := 0; day < 3; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	return w
}

// atWar puts two families who are not the player's at each other's throats.
func atWar(w *World) {
	for i := range w.Conflicts {
		c := &w.Conflicts[i]
		if c.A == w.PlayerOrganizationID() || c.B == w.PlayerOrganizationID() {
			continue
		}
		c.State, c.Hostility = "war", 100
		return
	}
}

func TestAWarKeepsPeopleAtHome(t *testing.T) {
	t.Parallel()
	evening := func(w *World) int {
		for i := 0; i < 40 && !Evening(w.Minute); i++ {
			w.Advance(60)
		}
		w.Advance(180)
		out := 0
		for _, id := range haunts {
			out += w.Footfall(id)
		}
		return out
	}
	peace := townsman(t)
	war := townsman(t)
	atWar(war)
	if !war.CityAtWar() {
		t.Fatal("the city is not at war after two families went to war")
	}
	if peace.CityAtWar() {
		t.Fatal("a city at peace says it is fighting")
	}
	quiet, ordinary := evening(war), evening(peace)
	t.Logf("of an evening: %d people out with a war on, %d with the city at peace", quiet, ordinary)
	if quiet >= ordinary {
		t.Fatalf("a war on and the same %d people went out drinking", quiet)
	}
}

func TestAWarCostsAnOwnerWhoIsNotInIt(t *testing.T) {
	t.Parallel()
	peace := townsman(t)
	war := townsman(t)
	atWar(war)
	for _, w := range []*World{peace, war} {
		for i := 0; i < 40 && !Evening(w.Minute); i++ {
			w.Advance(60)
		}
		w.Advance(180)
	}
	// A casino's takings are the action it can cover, and that is the people in
	// the room. The club is the one business in this city people go out to.
	if war.NightHandleAt("club") >= peace.NightHandleAt("club") {
		t.Fatalf("a war emptied the streets and the club covered $%d against $%d",
			war.NightHandleAt("club"), peace.NightHandleAt("club"))
	}
}

// keepWar holds one fight that is not the player's open for the whole run.
func keepWar(w *World) {
	for i := range w.Conflicts {
		c := &w.Conflicts[i]
		if c.A == w.PlayerOrganizationID() || c.B == w.PlayerOrganizationID() {
			continue
		}
		c.State, c.Hostility = "war", 100
		return
	}
}

// findable reports the share of person-hours where one place accounts for three
// quarters of somebody's appearances at that hour — the same measure the habit
// guard uses, so the two numbers can be read against each other.
func findable(moved map[string]int, place map[string]map[int]map[string]int) float64 {
	reliable, samples := 0, 0
	for id, hours := range place {
		if moved[id] == 0 {
			continue
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
		return 0
	}
	return float64(reliable) / float64(samples)
}

// A war is meant to be legible, not to turn the city into noise. The habit
// guard holds the city at peace; this one holds a war open for the whole
// fortnight and says the rhythm bends without breaking.
func TestAWarThinsTheStreetWithoutErasingIt(t *testing.T) {
	t.Parallel()
	quiet, quietHours, quietPlace := weekUnder(t, 404, 14, keepPeace)
	loud, loudHours, loudPlace := weekUnder(t, 404, 14, keepWar)
	peace, war := findable(quiet, quietPlace), findable(loud, loudPlace)
	t.Logf("a face is where it usually is on %.0f%% of person-hours at peace, %.0f%% with a war on", peace*100, war*100)
	// Staying home can make a resident more predictable. Measure the actual
	// evening crowd, rather than requiring war to reduce predictability.
	peaceCrowd, warCrowd := 0, 0
	for _, id := range haunts {
		peaceCrowd += quietHours[21][id]
		warCrowd += loudHours[21][id]
	}
	if warCrowd >= peaceCrowd {
		t.Fatalf("war did not thin the evening crowd: %d against %d", warCrowd, peaceCrowd)
	}
	if war < .7 {
		t.Fatalf("a war made the city unlearnable: a face is findable on only %.0f%% of person-hours", war*100)
	}
}
