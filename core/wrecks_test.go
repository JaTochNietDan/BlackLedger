package core

import "testing"

// "The business could get more business if cars are destroyed during
// operations." The player could take a car apart and nobody else could, so
// every car lost in this city was lost to the protagonist. A raid already
// reaches the people standing at a place; it should reach what they drive.
//
// And a car has to end somewhere. A scrapyard is where.

func TestARaidReachesWhatPeopleDrive(t *testing.T) {
	// My first version counted cars among a family's surviving members before
	// and after, and passed with twenty-two "lost" before a line of this
	// existed — because somebody killed in a raid leaves the member list, so
	// it was counting DEATHS. This follows named people who are alive at both
	// ends, so the only thing it can see is a car taken off somebody who
	// lived through it.
	lost, watched := 0, 0
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
			if n == nil || n.Dead {
				continue // died; that is not a car being taken
			}
			watched++
			if n.Car == 0 {
				lost++
			}
		}
	}
	if watched == 0 {
		t.Fatal("nobody survived to be measured, so this proves nothing")
	}
	if lost == 0 {
		t.Errorf("six cities at war for a month, %d people came through it, and not one lost a car", watched)
	}
	t.Logf("%d of %d survivors lost what they drove", lost, watched)
}

// A car does not simply evaporate. Wherever one ends, it ends at a scrapyard,
// and a scrapyard does better the more of them there are.
func TestTheCityHasSomewhereACarEnds(t *testing.T) {
	yards := 0
	for _, l := range Locations {
		if l.Kind == "scrapyard" {
			yards++
			if PlaceIncome[l.ID] <= 0 {
				t.Errorf("%s is a scrapyard that earns nothing", l.ID)
			}
			if _, runs := TradeOf(l.ID); !runs {
				t.Errorf("%s is a scrapyard with no trade", l.ID)
			}
		}
	}
	if yards == 0 {
		t.Fatal("there is nowhere in this city a car ends")
	}
}

// And a wreck is worth something to the yard. Stripping a car should reach the
// scrapyard the way it reaches the garages: more wrecks, more trade.
func TestAWreckIsWorthSomethingToTheYard(t *testing.T) {
	w, _ := stripper(t)
	before := map[string]int{}
	for _, l := range Locations {
		if l.Kind == "scrapyard" {
			before[l.ID] = w.Custom(l.ID)
		}
	}
	if len(before) == 0 {
		t.Fatal("no scrapyard")
	}
	if err := w.StripCar("bar"); err != nil {
		t.Fatal(err)
	}
	for id, was := range before {
		if w.Custom(id) <= was {
			t.Errorf("a car went to pieces and %s saw no more work", id)
		}
	}
}
