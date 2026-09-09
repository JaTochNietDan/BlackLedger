package core

import "testing"

// The last line of the inbox entry: "more car repairs to be made from broken
// windows from theft". Stripping and raids took a car away outright, which
// gives a garage parts and a forecourt a sale but never gives a garage the
// thing it actually lives on — somebody at the counter paying to be put back
// on the road.

func TestStrippingACarGoesThroughTheRestOfTheRow(t *testing.T) {
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
	w := New(61)
	w.Properties["garage"].Owner = "player:1"
	var mark *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && w.WouldDrive(n) {
			n.Car, n.Hurt, n.Purse = 1, true, 500
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
			n.Car, n.Hurt = 1, true
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
