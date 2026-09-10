package core

import "testing"

// Now that a room's takings depend on who is standing in it, an owner needs
// something to do about an empty room. Putting a night on is the oldest answer
// there is: pay for a band and a barrel, and the people who would have drunk
// somewhere else drink here instead.
//
// It has to move actual people. A number that goes up without anybody walking
// through the door would be the same abstraction the takings used to be.

func host(t *testing.T) *World {
	t.Helper()
	w := New(47)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 5000, 100
	own(w, "club")
	w.Player.Location = "club"
	// Let the city settle into its own habits first. Nobody has a place they
	// usually drink until they have had a day of standing somewhere, and a
	// night is a reason to go somewhere other than usual.
	for day := 0; day < 3; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	w.Player.Cash, w.Player.Location = 5000, "club"
	return w
}

func TestPuttingOnANightFillsTheRoom(t *testing.T) {
	w := host(t)
	a := actionByID(w.Actions("club"), "night")
	if a == nil {
		t.Fatal("a room of yours cannot put a night on")
	}
	if a.Disabled {
		t.Fatalf("refused: %s", a.Reason)
	}
	cash := w.Player.Cash
	if err := w.PutOnANight("club"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash >= cash {
		t.Fatal("the band played for nothing")
	}
	if !w.NightOn("club") {
		t.Fatal("the night is not on")
	}

	// Against the same evening without one. Measuring the room before and after
	// the evening comes measures the evening: the bar, the club and the casino
	// fill up every night whether anybody paid for a band or not, and a test
	// that watched one room fill passed with the draw deleted.
	quiet := host(t)
	evening := func(x *World) int {
		for i := 0; i < 40 && !Evening(x.Minute); i++ {
			x.Advance(60)
		}
		x.Advance(180) // long enough for the walk
		return x.Footfall("club")
	}
	loud := evening(w)
	ordinary := evening(quiet)
	if loud <= ordinary {
		t.Fatalf("a night drew %d where an ordinary evening drew %d", loud, ordinary)
	}
}

func TestANightIsOneNight(t *testing.T) {
	w := host(t)
	if err := w.PutOnANight("club"); err != nil {
		t.Fatal(err)
	}
	if w.NightReadiness("club") == "" {
		t.Fatal("two nights were put on at once")
	}
	w.Advance(NightLasts + 1440)
	if w.NightOn("club") {
		t.Fatal("the band never went home")
	}
	w.Player.Cash = 5000
	if reason := w.NightReadiness("club"); reason != "" {
		t.Fatalf("could not put another one on a week later: %s", reason)
	}
}

func TestYouCannotPutANightOnSomebodyElsesRoom(t *testing.T) {
	w := host(t)
	w.Player.Location = "bar"
	if w.NightReadiness("bar") == "" {
		t.Fatal("paid for a band in a room that is not yours")
	}
	w.Player.Location = "laundry"
	w.Properties["laundry"].Owner = w.Properties["club"].Owner
	if w.NightReadiness("laundry") == "" {
		t.Fatal("a laundry put a night on")
	}
}
