package core

import "testing"

// A raid burns a car and the yards feel it, but the person who lost it only got
// another when the day's one sale happened to reach them — first eligible name
// in the city's own order, so somebody far enough down the list walked for the
// rest of the campaign. Losing what you drive is a reason to go somewhere, the
// same way a windscreen is.

// forecourtCity is a city with a lot somebody holds and a driver who has just
// lost their car.
func forecourtCity(t *testing.T, seed uint32) (*World, *NPC) {
	t.Helper()
	w := New(seed)
	w.District = 2
	lot := ""
	for _, l := range Locations {
		if l.Kind == "dealer" {
			lot = l.ID
		}
	}
	if lot == "" {
		t.Fatal("this city has nowhere to buy a car")
	}
	w.Properties[lot].Owner = "bellandi"
	// The last person in the city's own order who drives, because the fault
	// being measured is that the queue was read from the top every day.
	var mark *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && w.WouldDrive(n) {
			mark = n
		}
	}
	if mark == nil {
		t.Fatal("nobody in this city drives")
	}
	mark.Car, mark.Drove, mark.Purse = 0, max(1, w.Minute), 4000
	return w, mark
}

func TestLosingACarIsAReasonToGoToTheForecourt(t *testing.T) {
	t.Parallel()
	w, mark := forecourtCity(t, 41)
	w.SetOut()
	if mark.Heading == "" {
		t.Fatal("somebody with four thousand dollars and nothing to drive is going nowhere")
	}
	place, _ := PlaceByID(mark.Heading)
	if place.Kind != "dealer" {
		t.Fatalf("%s is heading for %s rather than a forecourt", mark.Name, place.Name)
	}
	t.Logf("%s: %s", mark.Name, mark.Errand)
}

// The first version of this measured how long a replacement takes and PASSED
// before a line of the errand existed — one person with nothing to drive is the
// only name the day's sale can reach, so it reached them the next morning. That
// is not what was wrong. What was wrong is that nobody went anywhere: a sale
// happened to whoever stood first in the city's own order, wherever they were,
// and the forecourt was an empty room with a takings figure attached.
//
// So this measures the traffic, and holds the replacement time as a floor so
// the fix cannot buy a room full of people at the price of a city on foot.
func TestAForecourtHasPeopleInItAfterCarsAreLost(t *testing.T) {
	t.Parallel()
	standing, slowest, never, cities := 0, 0, 0, 0
	for _, seed := range []uint32{41, 77, 109, 233, 311} {
		w := New(seed)
		w.District = 2
		cities++
		lot := ""
		for _, l := range Locations {
			if l.Kind == "dealer" {
				lot = l.ID
			}
		}
		// Nobody holds the lot. A family that held it would send its own people
		// to mind it, and this is counting who came to BUY — the first version
		// of this test gave it to the Bellandis and counted their doorman.
		w.Properties[lot].Owner = "independent"
		// A bad week: eight people lose what they drive, which is what a month
		// of raids does to a family.
		lost := map[string]bool{}
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if n.Dead || n.Car == 0 || len(lost) >= 8 {
				continue
			}
			n.Car, n.Drove, n.Purse = 0, max(1, w.Minute), 4000
			lost[n.ID] = true
		}
		if len(lost) == 0 {
			t.Fatal("nobody in this city drives")
		}
		days := 0
		for ; days < 30; days++ {
			w.RNG, w.WorldRNG = seed+uint32(days), seed+uint32(days)
			aDay(w)
			for i := range w.NPCs {
				if n := &w.NPCs[i]; !n.Dead && n.Location == lot && lost[n.ID] {
					standing++
				}
			}
			w.CarTrade()
			done := true
			for id := range lost {
				if n := w.NPC(id); n != nil && n.Car == 0 {
					done = false
				}
			}
			if done {
				break
			}
		}
		for id := range lost {
			if n := w.NPC(id); n != nil && n.Car == 0 {
				never++
			}
		}
		if days > slowest {
			slowest = days
		}
	}
	t.Logf("across %d cities: %d people standing at a forecourt, %d still walking after a month, slowest city took %d days",
		cities, standing, never, slowest)
	if standing == 0 {
		t.Error("eight people lost their cars and not one of them went to a forecourt")
	}
	if never > 0 {
		t.Errorf("%d people never replaced what they lost", never)
	}
}

// aDay runs one ordinary day of the city's movement the way the clock does it:
// people are sent out twice, and arrivals are checked on every step. A harness
// that advances twelve hours in one jump and only then looks for arrivals
// cancels every journey longer than the gap between its own steps — which is
// what the first version of this measurement did, and it left somebody standing
// at the docks for a month with the money for a car in their pocket.
func aDay(w *World) {
	for step := 0; step < 24; step++ {
		if step == 0 || step == 12 {
			w.SetOut()
		}
		w.Minute += 60
		w.Arrivals()
	}
}
