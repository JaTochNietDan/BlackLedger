package core

import "testing"

// Missing payroll costs what the people behind the counter think of the player,
// and that number decides how well the place runs. It does not decide whether
// they stay: leaving is about a grudge, or about being paid under the rate. So
// a player could go weeks paying nobody and keep every hand, thinking less of
// them each morning.
//
// Somebody who is not being paid at all stops coming in. Not because they think
// little of you — everybody here starts at nothing and a stranger thinks
// nothing of anybody — but because there is no money in it, which is a fact
// about the week rather than about them.

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
	// A bill nothing can cover: security is charged by the day and earns
	// nothing back.
	// Nothing the laundry earns can cover this, so no night quietly clears and
	// resets the count.
	w.Player.Security, w.Player.Cash = 200, 0
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

// And nobody can be unpaid for long, which is the finding rather than the
// feature. The first night the bills do not clear strips the security and the
// address, so the day's cost falls from $2,033 to $33 and every night after
// clears out of what the business earns.
func TestNobodyStaysUnpaidForLong(t *testing.T) {
	t.Parallel()
	w, id := unpaidRun(t)
	worst := 0
	for day := 0; day < 20; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
		if n := w.Properties[id].Unpaid; n > worst {
			worst = n
		}
	}
	t.Logf("twenty days of a player with nothing: the longest unpaid run was %d nights", worst)
	if worst == 0 {
		t.Fatal("a player with nothing covered every bill, so this measures nothing")
	}
	if worst >= 7 {
		t.Fatalf("a week of unpaid wages is reachable after all, at %d nights — the rule that was "+
			"removed for being unreachable should go back in", worst)
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
	w.Player.Security, w.Player.Cash = 0, 50000
	for day := 0; day < 3 && w.Properties[id].Unpaid != 0; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if w.Properties[id].Unpaid != 0 {
		t.Fatalf("the bills were covered and the place still counts %d unpaid nights", w.Properties[id].Unpaid)
	}
}
