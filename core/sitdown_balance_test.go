package core

import "testing"

// A room that always settles is a button, and one that always ends in gunfire
// is a trap with extra steps. These measure the spread across cities the
// simulation produced on its own.

func TestSomeRoomsAreSafeAndSomeAreNot(t *testing.T) {
	const runs, days = 300, 90
	quarrels, traps, warned := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		w.Player.Contacts = 2
		for day := 0; day < days; day++ {
			w.FactionTurn()
			w.FamilyDay()
			w.PeopleDay()
			w.GrudgeDay()
			w.SettleGrudges()
			w.Minute += 1440
		}
		q, ok := w.OpenQuarrel()
		if !ok {
			continue
		}
		quarrels++
		if q.Trap {
			traps++
		}
		if q.Suspected {
			warned++
		}
	}
	if quarrels == 0 {
		t.Fatalf("no city in %d had a quarrel worth mediating after %d days", runs, days)
	}
	if traps == 0 {
		t.Fatal("no quarrel in any city was ever an ambush")
	}
	if traps == quarrels {
		t.Fatal("every quarrel was an ambush, so mediating is never worth trying")
	}
	if warned != traps {
		t.Fatalf("with two contacts, %d of %d ambushes were flagged", warned, traps)
	}
	t.Logf("%d of %d cities had a quarrel worth calling a room over after %d days; %d of those rooms had somebody in them who had already decided", quarrels, runs, days, traps)
}

func TestSettlingAWarIsWorthTheEvening(t *testing.T) {
	const runs = 200
	settled, ended, hurt := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
		w.Player.Location, w.Player.Cash, w.Player.Respect, w.Player.Health = SitdownGround, 5000, 60, 100
		c := &w.Conflicts[0]
		c.Hostility = 55 + int(w.WorldRandom()*40)
		c.State = classify(c)
		if w.SitdownReadiness() != "" {
			continue
		}
		if err := w.CallSitdown(); err != nil {
			t.Fatal(err)
		}
		state := w.Conflict(w.Event.Actor, w.Event.Target).State
		if err := w.ResolveSitdown(w.Event, "press"); err != nil {
			t.Fatal(err)
		}
		settled++
		if w.Conflict(w.Event.Actor, w.Event.Target).State != state {
			ended++
		}
		if w.Player.Health < 100 {
			hurt++
		}
	}
	if settled == 0 {
		t.Fatal("no room was ever called")
	}
	if ended == 0 {
		t.Fatal("pressing never changed the state of a quarrel")
	}
	if hurt == 0 {
		t.Fatal("pressing was never dangerous")
	}
	if hurt == settled {
		t.Fatal("pressing was always dangerous, so nobody would ever do it")
	}
	t.Logf("%d rooms pressed to settle: the quarrel changed state in %d, and the player was hurt in %d", settled, ended, hurt)
}
