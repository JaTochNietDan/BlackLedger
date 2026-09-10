package core

import "testing"

// A car is bought with money and paid for in time, so the only honest way to
// price one is to give two identical people the same number of days and see how
// much more the one who drives gets done.

// circuit works a fixed round of jobs across the city until the clock runs out,
// and reports what reached the player's hands net of every bill they paid.
func circuit(seed uint32, car, days int) (int, int) {
	w := New(seed)
	w.WorldRNG = seed * 2654435761
	w.Player.Cash = 3000
	if car > 0 {
		w.Player.Car, w.Player.CarWear = car, 100
	}
	start, budget := w.Player.Cash, days*1440
	stops := []struct{ place, work string }{{"bar", "courier"}, {"docks", "dockwork"}}
	jobs, stuck := 0, 0
	try := func(kind, target string) bool {
		next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: kind, Target: target})
		if err != nil {
			return false
		}
		w = next
		return true
	}
	for i := 0; w.Minute < budget && w.Player.Alive && stuck < 8; i++ {
		// Whatever the city puts in front of them, both people answer it the
		// same way, so the only difference between the two runs is the driving.
		if e := w.Event; e != nil {
			if !try("choice", "") {
				next, err := Execute(w, Command{RequestID: w.Event.ID, Revision: w.Revision, Kind: "choice", Event: e.ID, Choice: e.Choices[len(e.Choices)-1].ID})
				if err != nil {
					stuck++
					continue
				}
				w = next
			}
			continue
		}
		// A driver buys petrol. The car in this comparison now has a tank that
		// goes down, so the driving life includes the detour to a forecourt and
		// the price of filling it — which is the honest question: does a car
		// buy time NET of keeping it running? The walking life never does this,
		// because a person on foot has nothing to fill.
		if car > 0 && w.Fuel() < FuelFull/4 {
			station := w.theStation()
			if station != "" && w.Player.Location != station {
				if !try("travel", station) {
					stuck++
					continue
				}
			}
			if !try("fill", station) {
				stuck++
			}
			continue
		}
		stop := stops[i%len(stops)]
		if w.Player.Location != stop.place && !try("travel", stop.place) {
			stuck++
			continue
		}
		if try(stop.work, stop.place) {
			jobs++
			stuck = 0
			continue
		}
		if !try("wait", w.Player.Location) {
			stuck++
		}
	}
	// Cars are paid for out of the same pocket the jobs fill, so the comparison
	// is what is left rather than what was earned.
	return w.Player.Cash - start, jobs
}

func TestACarPaysForItselfInWorkDone(t *testing.T) {
	heavy(t)
	const runs, days = 24, 16
	walked, drove, walkJobs, driveJobs := 0, 0, 0, 0
	for i := uint32(1); i <= runs; i++ {
		cash, jobs := circuit(i, 0, days)
		walked += cash
		walkJobs += jobs
		cash, jobs = circuit(i, 2, days)
		drove += cash
		driveJobs += jobs
	}
	if driveJobs <= walkJobs {
		t.Fatalf("driving completed %d jobs against %d on foot, so a car buys no time", driveJobs, walkJobs)
	}
	if drove <= walked {
		t.Fatalf("driving netted $%d against $%d walking, so the upkeep is not worth paying", drove/runs, walked/runs)
	}
	gain := float64(driveJobs-walkJobs) / float64(walkJobs) * 100
	t.Logf("over %d days across %d campaigns: on foot %d jobs and $%d, driving a Hudson %d jobs (+%.0f%%) and $%d, after $%d a day of upkeep",
		days, runs, walkJobs, walked/runs, driveJobs, gain, drove/runs, VehicleByTier(2).Upkeep)
}

// TestTheBestCarIsNotAlwaysTheRightCar measures what the trail costs. An
// armoured Packard is the fastest thing in the city and the most described, so
// somebody who commits crimes in one pays for the speed in police attention.
func TestTheBestCarIsNotAlwaysTheRightCar(t *testing.T) {
	heavy(t)
	const runs = 300
	heat := map[int]int{}
	for _, tier := range []int{0, 1, 2, 3} {
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(seed * 2654435761)
			w.Player.Location, w.Player.Health, w.Player.Respect = "bar", 100, 40
			if tier > 0 {
				w.Player.Car, w.Player.CarWear = tier, 100
			}
			w.Rob("bar")
			heat[tier] += w.Player.Heat
		}
	}
	for tier := 1; tier <= 3; tier++ {
		if heat[tier] <= heat[tier-1] {
			t.Fatalf("a tier %d car drew %d attention against %d for tier %d", tier, heat[tier], heat[tier-1], tier-1)
		}
	}
	t.Logf("attention across %d robberies: on foot %d, Ford %d, Hudson %d, Packard %d", runs, heat[0], heat[1], heat[2], heat[3])
}
