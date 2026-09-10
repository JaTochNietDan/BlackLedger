package core

import "testing"

// Core owns the list of groups an action can belong to, complete with titles,
// and says so: "Groups is the ordered list, for anything that renders them."
// Nothing rendered them. The view kept its own copy of the list, and the copy
// was missing "people".
//
// The effect: any action in that group whose subject is not standing in the
// room disappeared from the panel entirely. Not greyed out with a reason —
// gone. Paying the crew a bonus while they are out on collections is the
// ordinary case, and the player could neither see the button nor learn why it
// was unavailable.
//
// This states the property from core's side. The view now renders from what
// core sends, so the two cannot drift again.

func TestEveryGroupAnActionCanCarryHasATitle(t *testing.T) {
	t.Parallel()
	titled := map[string]bool{}
	for _, g := range Groups() {
		if g.Title == "" || g.Blurb == "" {
			t.Errorf("the %q group has no title or blurb", g.ID)
		}
		titled[g.ID] = true
	}
	// Every action the game can offer, in every room, across the states the
	// reachability sweep builds.
	offered, _ := sweepCities(t)
	seen := map[string]bool{}
	for kind := range offered {
		seen[GroupOf(kind)] = true
	}
	for group := range seen {
		if !titled[group] {
			t.Errorf("actions are grouped %q and nothing can render it", group)
		}
	}
	if !seen["people"] {
		t.Fatal("no action was grouped as people, so this proved nothing")
	}
	t.Logf("groups in use: %v", seen)
}
