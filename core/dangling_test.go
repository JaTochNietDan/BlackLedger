package core

import "testing"

// Nothing in a city that has been running points at a name that is not there.
//
// Written after two nights of finding these one at a time: the player still
// answering to a family that had ended, and an understanding with one, both
// reachable only because spreading the campaign numbers showed a city that
// actually buries organizations. A city churns, and every fall, death and change
// of hands leaves an id somewhere that used to mean something.
func TestNothingInTheCityPointsAtANameThatIsGone(t *testing.T) {
	t.Parallel()
	cities, found := 0, 0
	for n := uint32(1); n <= 24; n++ {
		w := New(spread(n))
		w.Event = nil
		w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 200, 100
		// Things to dangle: an understanding and a ceasefire with everybody,
		// and somebody to answer to.
		w.BusinessTruces = map[string]int{}
		for _, f := range w.Factions {
			if f.ID == w.PlayerOrganizationID() {
				continue
			}
			w.BusinessTruces[f.ID] = w.Minute + 400*1440
			w.Pacts = append(w.Pacts, Pact{With: f.ID, Since: w.Minute, Life: w.Life})
			if w.Player.Serves == "" {
				w.Player.Serves = f.ID
			}
		}
		if bad := w.Dangling(); len(bad) > 0 {
			t.Fatalf("a city points at nothing before it has even run: %v", bad)
		}
		w.Advance(180 * 1440)
		cities++
		for _, bad := range w.Dangling() {
			found++
			if found <= 6 {
				t.Errorf("city %d: %s", n, bad)
			}
		}
	}
	t.Logf("%d cities of half a year, %d references pointing at nothing", cities, found)
	if cities < 20 {
		t.Fatalf("only %d cities ran, so this measures nothing", cities)
	}
}

// And the reader can see one. A sweep that finds nothing is either a clean city
// or a broken reader, and there is no way to tell from the outside.
func TestTheDanglingReaderCanSeeOne(t *testing.T) {
	t.Parallel()
	w := New(spread(1))
	w.Event = nil
	if bad := w.Dangling(); len(bad) > 0 {
		t.Fatalf("a fresh city already points at nothing: %v", bad)
	}
	w.Player.Serves = "a-family-that-never-was"
	if bad := w.Dangling(); len(bad) != 1 {
		t.Fatalf("the player answers to nobody and the reader found %v", bad)
	}
	w.Player.Serves = ""
	w.Pacts = append(w.Pacts, Pact{With: "another-one", Since: w.Minute, Life: w.Life})
	w.BusinessTruces = map[string]int{"a-third": w.Minute + 1440}
	if bad := w.Dangling(); len(bad) != 2 {
		t.Fatalf("two agreements point at nothing and the reader found %v", bad)
	}
}
