package core

import "testing"

// "Taking over businesses should be a lot more expensive and high level stuff
// that you build up to over time."
//
// It was neither. A laundry was $180 against a campaign that ends with twelve
// thousand in the bank, so the whole business layer opened on the second
// afternoon and every address after the first cost exactly what the first one
// did. Nothing about holding four premises made the fifth a decision.

func freeholder(t *testing.T) *World {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health, w.Player.Respect = 200000, 100, 60
	return w
}

func TestTheFirstPremisesIsAPurchaseRatherThanAnAfternoon(t *testing.T) {
	t.Parallel()
	w := freeholder(t)
	// Against what a careful campaign actually earns: the worker's median is
	// about $12,500 over a couple of weeks.
	cheapest := 0
	for _, l := range Locations {
		if l.Cost <= 0 || l.Type == "home" {
			continue
		}
		if cost := AcquisitionCost(w, l.ID); cheapest == 0 || cost < cheapest {
			cheapest = cost
		}
	}
	t.Logf("the cheapest freehold in the city is $%d", cheapest)
	// Against what a careful campaign earns in a fortnight, about $12,500: the
	// cheapest door in the city should be a few days of it rather than an
	// afternoon. Four times the old listing puts it at $720, and every one held
	// raises the next by nearly half.
	if cheapest < 600 {
		t.Fatalf("the cheapest business in the city is $%d, which is an afternoon's work", cheapest)
	}
}

func TestEachOneYouHoldMakesTheNextDearer(t *testing.T) {
	t.Parallel()
	w := freeholder(t)
	first := AcquisitionCost(w, "haulage")
	own(w, "laundry")
	second := AcquisitionCost(w, "haulage")
	own(w, "butcher")
	third := AcquisitionCost(w, "haulage")
	t.Logf("the haulage yard is $%d, then $%d, then $%d as you take the city", first, second, third)
	if second <= first || third <= second {
		t.Fatalf("holding the city changes nothing about what the next place costs: %d, %d, %d", first, second, third)
	}
	// And it is a premium rather than a wall: the fourth is dearer than the
	// first without being out of a successful campaign's reach.
	if third > first*4 {
		t.Fatalf("the third purchase is $%d against a first of $%d, which is a wall", third, first)
	}
}

// And it is high-level: a nobody is not handed the keys to anything.
func TestNobodyIsHandedTheKeysToAnything(t *testing.T) {
	t.Parallel()
	w := freeholder(t)
	w.Player.Respect = 0
	refused := 0
	for _, l := range Locations {
		if l.Cost <= 0 || l.Type == "home" {
			continue
		}
		reason := w.AcquireReadiness(l.ID)
		if reason == "" {
			t.Fatalf("%s is handed to somebody the city has never heard of", l.Name)
		}
		refused++
	}
	if refused == 0 {
		t.Fatal("nothing in this city is for sale, so this proves nothing")
	}
	// And the ladder has a bottom rung: somebody who has done a little can buy
	// something, or it is a wall rather than a thing to build up to.
	w.Player.Respect = 15
	open := 0
	for _, l := range Locations {
		if l.Cost > 0 && l.Type != "home" && w.AcquireReadiness(l.ID) == "" {
			open++
		}
	}
	t.Logf("with 15 respect and money in hand, %d of the city's freeholds are open", open)
	if open == 0 {
		t.Fatal("somebody who has made a name for themselves can buy nothing at all")
	}
}
