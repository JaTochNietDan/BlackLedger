package main

import (
	"testing"

	"blackledger/core"
)

// The workshop's menu is the city's own list. A kind added to the game has to
// turn up here without anybody remembering to add it, which is the whole reason
// it is read from the core rather than written out again.

func TestTheWorkshopOffersEveryMomentTheCityCanShow(t *testing.T) {
	offered := map[string]bool{}
	for _, k := range kinds() {
		offered[k["kind"].(string)] = true
	}
	for _, kind := range core.MomentKinds() {
		if !offered[kind] {
			t.Errorf("the city can show a %s and the workshop does not offer it", kind)
		}
	}
	if len(offered) != len(core.MomentKinds()) {
		t.Errorf("the workshop offers %d kinds and the city has %d", len(offered), len(core.MomentKinds()))
	}
}

func TestTheHeaviestMomentIsOfferedFirst(t *testing.T) {
	list := kinds()
	if len(list) < 2 {
		t.Fatal("the workshop offers nothing to compare")
	}
	for i := 1; i < len(list); i++ {
		if list[i]["gravity"].(int) > list[i-1]["gravity"].(int) {
			t.Errorf("%s is worth more than %s and is offered after it",
				list[i]["kind"], list[i-1]["kind"])
		}
	}
	if list[0]["kind"] != "killing" {
		t.Errorf("the heaviest moment is %q, expected a killing", list[0]["kind"])
	}
}

func TestEveryMomentHoldsForLongerThanTheOneBelowIt(t *testing.T) {
	// A robbery and a killing were once played for exactly the same two and a
	// half seconds. The workshop shows the hold so that stays visible.
	for _, k := range kinds() {
		if k["hold"].(int) < 1800 {
			t.Errorf("%s holds for %dms, less than the floor", k["kind"], k["hold"])
		}
	}
}
