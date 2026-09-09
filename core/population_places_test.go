package core

import "testing"

// "A lot of characters living in this city" is not a bigger number on its own.
// The street trades — the people who are not in this business and live here
// anyway — name the addresses they work at, and that list was written when the
// city had ten addresses. It has twenty. Half of them had nobody in them at
// all: a butcher with no butcher, a cab company with no drivers, a revue bar
// with nobody on the stage.
//
// So the test is about the city being inhabited, not about the count.

func TestSomebodyWorksAtEveryAddressWorthWorkingAt(t *testing.T) {
	worked := map[string]bool{}
	for _, trade := range streetTrades {
		worked[trade.place] = true
	}
	empty := []string{}
	for _, l := range Locations {
		// The police station and the newspaper have their own standing jobs
		// and are not places the street works out of.
		if l.Type == "civic" {
			continue
		}
		if !worked[l.ID] {
			empty = append(empty, l.ID)
		}
	}
	if len(empty) > 0 {
		t.Errorf("%d addresses have nobody who works there: %v", len(empty), empty)
	}
}

// And every one of those roles has to be somewhere a person can actually be.
func TestNobodyWorksAtAnAddressTheCityDoesNotHave(t *testing.T) {
	for _, trade := range streetTrades {
		if _, ok := PlaceByID(trade.place); !ok {
			t.Errorf("a %s works at %q, which is not a place in this city", trade.role, trade.place)
		}
	}
}

// Two people in the same room doing the same job is a shorter city than it
// looks. Roles repeat across addresses — there is more than one barman — but
// no address should have the same job twice.
func TestNoAddressHasTheSameJobTwice(t *testing.T) {
	seen := map[string]bool{}
	for _, trade := range streetTrades {
		key := trade.place + "/" + trade.role
		if seen[key] {
			t.Errorf("%s has two people doing the same job: %s", trade.place, trade.role)
		}
		seen[key] = true
	}
}

// A city that means to hold sixty-two people on the street has to actually
// manage it. Naming somebody is sixty attempts at a random first name and
// surname, and a pool that was ample for twenty-six people is not necessarily
// ample for sixty-two: the failure is silent, because AddCivilian returns
// nothing and the city is simply smaller than it meant to be.
func TestTheCityFillsTheStreetItMeansTo(t *testing.T) {
	short := 0
	for _, seed := range []uint32{7, 31, 88, 149, 219, 401} {
		w := New(seed)
		loose := 0
		for _, n := range w.People() {
			if n.Faction == "" && !IsOfficial(n.ID) {
				loose++
			}
		}
		if loose < StreetCount {
			short++
			t.Logf("seed %d: %d people on the street against %d intended", seed, loose, StreetCount)
		}
	}
	if short > 0 {
		t.Errorf("%d of 6 cities came up short of the street they meant to have", short)
	}
}

// And nobody in the city shares a name with anybody else, however many there
// are. Two Otto Reisses is two people the player cannot tell apart.
func TestNoTwoPeopleInTheCityShareAName(t *testing.T) {
	for _, seed := range []uint32{7, 31, 219} {
		w := New(seed)
		live(t, w, 3000)
		seen := map[string]string{}
		for _, n := range w.People() {
			if other, taken := seen[n.Name]; taken {
				t.Errorf("seed %d: %s and %s are both called %s", seed, other, n.ID, n.Name)
			}
			seen[n.Name] = n.ID
		}
		t.Logf("seed %d: %d people, all named apart", seed, len(seen))
	}
}
