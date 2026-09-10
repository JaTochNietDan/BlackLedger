package core

import (
	"strings"
	"testing"
)

// The last line of the inbox entry: "more car repairs to be made from broken
// windows from theft". Stripping and raids took a car away outright, which
// gives a garage parts and a forecourt a sale but never gives a garage the
// thing it actually lives on — somebody at the counter paying to be put back
// on the road.

func TestStrippingACarGoesThroughTheRestOfTheRow(t *testing.T) {
	t.Parallel()
	w, mark := stripper(t)
	// Somebody else parked in the same street who is not the one being taken.
	var other *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.ID == mark.ID || !w.WouldDrive(n) {
			continue
		}
		n.Location, n.Car, n.Hurt = "bar", 1, false
		other = n
		break
	}
	if other == nil {
		t.Fatal("only one person in this city drives")
	}
	if err := w.StripCar("bar"); err != nil {
		t.Fatal(err)
	}
	if other.Car == 0 {
		t.Fatalf("%s lost the car outright; the point is that it is still there", other.Name)
	}
	if !other.Hurt {
		t.Errorf("a car was taken apart beside %s and their own was left untouched", other.Name)
	}
}

// The money side. A broken car is a bill, and somebody in this city is paid it.
func TestABrokenCarPutsWorkAndMoneyThroughAGarage(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Properties["garage"].Owner = "player:1"
	var mark *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && w.WouldDrive(n) {
			// Standing at the bench, because that is now the whole of it: they
			// brought it in. This test used to leave them wherever the city had
			// put them, which passed while a repair was a figure moving once a
			// day rather than somebody at a counter.
			n.Car, n.Hurt, n.Purse, n.Location = 1, true, 500, "garage"
			mark = n
			break
		}
	}
	if mark == nil {
		t.Fatal("nobody in this city drives")
	}
	cash, custom, purse := w.Player.Cash, w.Custom("garage"), mark.Purse
	w.RepairsDay()
	if mark.Hurt {
		t.Fatal("somebody with the fee in their pocket drove past the garage")
	}
	if mark.Purse != purse-GlassCost {
		t.Errorf("the job cost them %d, not %d", purse-mark.Purse, GlassCost)
	}
	if w.Player.Cash != cash+GlassCost {
		t.Errorf("a job through your own garage earned %d", w.Player.Cash-cash)
	}
	if w.Custom("garage") <= custom {
		t.Errorf("custom is %d, was %d", w.Custom("garage"), custom)
	}
	t.Logf("$%d through the bench, custom %d up from %d", GlassCost, w.Custom("garage"), custom)
}

// And the link running the other way: a garage in a city with no money in it
// has the same crimes and less work, because the work goes unpaid for.
func TestSomebodyWhoCannotFindTheFeeKeepsDrivingItBroken(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Properties["garage"].Owner = "player:1"
	var mark *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead {
			continue
		}
		n.Purse = 0 // nobody in this city can pay for anything
		if w.WouldDrive(n) && mark == nil {
			n.Car, n.Hurt, n.Location = 1, true, "garage"
			mark = n
		}
	}
	if mark == nil {
		t.Fatal("nobody in this city drives")
	}
	cash, custom := w.Player.Cash, w.Custom("garage")
	w.RepairsDay()
	if !mark.Hurt {
		t.Error("somebody with nothing in their pocket had the car put right anyway")
	}
	if w.Player.Cash != cash || w.Custom("garage") != custom {
		t.Errorf("a garage was paid %d for a job nobody could afford", w.Player.Cash-cash)
	}
}

// A raid that does not burn what is parked outside goes through it instead.
// Same scaffold as the burning test, and the same rule: follow named people who
// are alive at both ends, because the dead leave the member list.
func TestARaidLeavesCarsNeedingWork(t *testing.T) {
	t.Parallel()
	broken, watched := 0, 0
	for _, seed := range []uint32{11, 41, 77, 109, 233, 311} {
		w := New(seed)
		w.District = 2
		a, b := w.faction("bellandi"), w.faction("russo")
		if a == nil || b == nil {
			continue
		}
		held := w.FamilyHoldings(b.ID)
		if len(held) == 0 {
			continue
		}
		had := map[string]bool{}
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if n.Faction != b.ID || n.Dead {
				continue
			}
			n.Location = held[0]
			n.Car, n.Drove = 1, 1
			had[n.ID] = true
		}
		for day := 0; day < 30; day++ {
			w.RNG, w.WorldRNG = seed*2654435761+uint32(day), seed*2654435761+uint32(day)
			w.Antagonize(a.ID, b.ID, 100)
			w.FactionTurn()
		}
		for id := range had {
			n := w.NPC(id)
			if n == nil || n.Dead || n.Car == 0 {
				continue // dead, or the car was taken outright: neither is a repair
			}
			watched++
			if n.Hurt {
				broken++
			}
		}
	}
	if watched == 0 {
		t.Fatal("nobody came through it still driving, so this proves nothing")
	}
	if broken == 0 {
		t.Errorf("a month of raids in six cities, %d cars still on the road, and not one needs a garage", watched)
	}
	t.Logf("%d of %d cars that survived need work", broken, watched)
}

// A car with the glass out of it is a reason to be somewhere. The repair used
// to happen as a silent transfer once a day, wherever the owner happened to be
// standing, which is a garage's trade with no garage in it. They take the car
// in now, so the bench has people at it who are there for a reason — and that
// is a room the player can walk into and meet somebody.
func TestABrokenCarTakesItsOwnerToTheGarage(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.District = 2
	var mark *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && w.WouldDrive(n) && n.Location != w.theGarage() {
			n.Car, n.Hurt, n.Purse = 1, true, 500
			mark = n
			break
		}
	}
	if mark == nil {
		t.Fatal("nobody in this city drives")
	}
	w.SetOut()
	if mark.Heading != w.theGarage() {
		t.Fatalf("%s is heading for %q rather than the garage", mark.Name, mark.Heading)
	}
	if mark.Errand == "" {
		t.Error("they are going there for no stated reason")
	}
	t.Logf("%s: %s", mark.Name, mark.Errand)
}

// And the money only moves when they are actually standing at the bench.
func TestNobodyIsBilledForWorkTheyNeverBroughtIn(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Properties["garage"].Owner = "player:1"
	var mark *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && w.WouldDrive(n) {
			n.Car, n.Hurt, n.Purse = 1, true, 500
			n.Location = "bar" // broken down on the other side of the city
			mark = n
			break
		}
	}
	if mark == nil {
		t.Fatal("nobody in this city drives")
	}
	// Nobody else may be standing at the bench with a broken car either.
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.ID != mark.ID && n.Location == "garage" {
			n.Hurt = false
		}
	}
	cash, purse := w.Player.Cash, mark.Purse
	w.RepairsDay()
	if mark.Purse != purse || w.Player.Cash != cash {
		t.Errorf("a car across the city was worked on: they paid %d and the bench took %d",
			purse-mark.Purse, w.Player.Cash-cash)
	}
	if !mark.Hurt {
		t.Error("a car nobody brought in came back mended")
	}
}

// What it is worth: how many people stand in a garage on an ordinary day. A
// bench nobody visits is a number in a ledger; a bench with somebody at it is a
// room worth walking into.
func TestAGarageHasPeopleInItOnAnOrdinaryDay(t *testing.T) {
	t.Parallel()
	visits, cities := 0, 0
	for _, seed := range []uint32{5, 23, 61, 97, 181} {
		w := New(seed)
		w.District = 2
		cities++
		// An ordinary week in a city where cars get gone through: the glass
		// goes somewhere every night, which is what the street does.
		for day := 0; day < 7; day++ {
			w.RNG, w.WorldRNG = seed+uint32(day), seed+uint32(day)
			for i := range w.NPCs {
				n := &w.NPCs[i]
				if !n.Dead && n.Car > 0 && !n.Hurt && n.Purse >= GlassCost && w.WorldRandom() < .1 {
					n.Hurt = true
				}
			}
			aDay(w)
			for _, n := range w.OnTheFloor(w.theGarage()) {
				_ = n
			}
			for i := range w.NPCs {
				if n := &w.NPCs[i]; !n.Dead && n.Location == w.theGarage() {
					visits++
				}
			}
			w.RepairsDay()
		}
	}
	t.Logf("%d people standing in a garage across %d cities over a week", visits, cities)
	if visits == 0 {
		t.Error("a week of broken glass and nobody ever went to a garage")
	}
}

// And the room says why they are in it. A family head standing in a garage
// reads as holding court unless somebody says otherwise, and what he is
// actually doing is waiting on a windscreen.
func TestTheRoomSaysWhyTheyAreAtTheBench(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.District = 2
	var mark *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Car > 0 {
			n.Hurt, n.Location = true, "garage"
			mark = n
			break
		}
	}
	if mark == nil {
		t.Fatal("nobody in this city drives")
	}
	said := w.doingNow(mark)
	if !strings.Contains(said, "glass") {
		t.Errorf("standing in a garage with the windscreen out, the city says: %q", said)
	}
	// And once it is put right they are not still waiting on it.
	mark.Hurt = false
	if after := w.doingNow(mark); strings.Contains(after, "glass") {
		t.Errorf("the car is mended and they are still waiting: %q", after)
	}
}
