package core

import "testing"

// Missing payroll costs what the people behind the counter think of the player,
// and that number decides how well the place runs. It did not decide whether
// they stay: leaving was about a grudge, or about being paid under the rate. So
// a player could go weeks paying nobody and keep every hand, thinking less of
// them each morning.
//
// Somebody who is not being paid at all stops coming in. Not because they think
// little of you — everybody here starts at nothing and a stranger thinks
// nothing of anybody — but because there is no money in it, which is a fact
// about the week rather than about them.

// A place that earns less than its counter costs is what makes an unpaid night
// happen at all. A held address in working order pays for its own staff several
// times over, so no amount of standing bills starves it: the security and the
// address are given up first and what is left clears out of the takings. A
// wrecked one still has people on the books and nothing coming in.
func unpaidRun(t *testing.T) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Location = "laundry"
	w.Player.Cash = 200000
	if err := w.apply(Command{Kind: "acquire", Target: "laundry", RequestID: "buyitthenstarve"}); err != nil {
		t.Fatal(err)
	}
	w.Properties["laundry"].Condition = 0
	w.Player.Cash = 0
	if w.Wages() == 0 {
		t.Fatal("nobody stands behind this counter, so nothing here can go unpaid")
	}
	return w, "laundry"
}

func TestAnUnpaidNightIsCountedAgainstThePlace(t *testing.T) {
	t.Parallel()
	w, id := unpaidRun(t)
	for day := 0; day < 3 && w.Properties[id].Unpaid == 0; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if w.Properties[id].Unpaid == 0 {
		t.Fatal("a night nobody covered was not counted against the place")
	}
}

// And a week of it empties the counter, which is the rule that could not be
// written while the day's bill was all or nothing.
func TestAWeekUnpaidAndTheyStopComingIn(t *testing.T) {
	t.Parallel()
	w, id := unpaidRun(t)
	filled, worst := w.Properties[id].Staff, 0
	for day := 0; day < 40; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
		if n := w.Properties[id].Unpaid; n > worst {
			worst = n
		}
	}
	t.Logf("forty days of a place that earns nothing: %d unpaid nights at worst, %d of %d hands left",
		worst, w.Properties[id].Staff, filled)
	if worst < PatienceRunsOut {
		t.Fatalf("a week of unpaid wages is still not reachable, at %d nights", worst)
	}
	if w.Properties[id].Staff >= filled {
		t.Fatal("nobody was paid for weeks and the whole counter still turned up")
	}
}

// A week of it is not the same as a day of it: somebody misses one payday and
// comes in the next morning, which is the ordinary case for a player who is
// briefly short.
func TestOneShortNightCostsNobody(t *testing.T) {
	t.Parallel()
	w, id := unpaidRun(t)
	filled := w.Properties[id].Staff
	// Days rather than a day: a scene can interrupt an advance before midnight,
	// so "one night" is the first night that actually finishes.
	for day := 0; day < 3 && w.Properties[id].Unpaid == 0; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if w.Properties[id].Unpaid == 0 {
		t.Fatal("a night nobody covered was not counted against the place")
	}
	if w.Properties[id].Staff < filled {
		t.Fatalf("one short night and somebody was already gone: %d of %d", w.Properties[id].Staff, filled)
	}
	// And paying again forgets it.
	w.Properties[id].Condition, w.Player.Cash = 100, 50000
	for day := 0; day < 3 && w.Properties[id].Unpaid != 0; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if w.Properties[id].Unpaid != 0 {
		t.Fatalf("the bills were covered and the place still counts %d unpaid nights", w.Properties[id].Unpaid)
	}
}
