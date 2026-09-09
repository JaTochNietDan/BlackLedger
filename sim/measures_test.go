package sim

import (
	"testing"

	"blackledger/core"
)

// The living-world measures are only worth having if they count the city and
// not the player. The first version counted the player naming their own outfit
// as the city creating a new family, which made it read as roughly one new
// organization per campaign for the strategies that form one and almost none
// for the strategies that do not — a number about the strategies, not the city.

func TestThePlayersOwnOrganizationIsNotTheCityMoving(t *testing.T) {
	w := core.New(101)
	eyes := watch(w)
	var got WorldMeasures

	// The player forms an organization.
	w.Factions = append(w.Factions, core.Faction{ID: w.PlayerOrganizationID(), Name: "Ward's people"})
	eyes.changed(w, &got)
	if got.FactionsCreated != 0 {
		t.Errorf("the player naming their outfit counted as %d new families in the city", got.FactionsCreated)
	}

	// Somebody else forms one.
	w.Factions = append(w.Factions, core.Faction{ID: "alcaraz", Name: "The Alcaraz people"})
	eyes.changed(w, &got)
	if got.FactionsCreated != 1 {
		t.Errorf("a family splitting off counted as %d", got.FactionsCreated)
	}

	// And the player's own dissolving is not the city losing a family either.
	// Removed by id: slicing off the front took an original family instead and
	// the test failed for the right reason on the wrong thing.
	kept := w.Factions[:0]
	for _, f := range w.Factions {
		if f.ID != w.PlayerOrganizationID() {
			kept = append(kept, f)
		}
	}
	w.Factions = kept
	eyes.changed(w, &got)
	if got.FactionsDestroyed != 0 {
		t.Errorf("the player's outfit ending counted as %d families destroyed", got.FactionsDestroyed)
	}
}

func TestAWarIsCountedOnceAndPlacedCorrectly(t *testing.T) {
	w := core.New(102)
	eyes := watch(w)
	var got WorldMeasures

	// A war between two other families.
	w.Conflicts = append(w.Conflicts, core.Conflict{A: "bellandi", B: "russo", State: "war"})
	eyes.changed(w, &got)
	eyes.changed(w, &got) // still running; must not count again
	if got.WarsStarted != 1 {
		t.Errorf("one war was counted %d times", got.WarsStarted)
	}
	if got.WarsElsewhere != 1 {
		t.Errorf("a war between two other families was not counted as elsewhere (%d)", got.WarsElsewhere)
	}

	// And one the player is in is a war, but not one happening elsewhere.
	w.Conflicts = append(w.Conflicts, core.Conflict{A: w.PlayerOrganizationID(), B: "bellandi", State: "war"})
	eyes.changed(w, &got)
	if got.WarsStarted != 2 {
		t.Errorf("the player's own war was not counted as a war (%d)", got.WarsStarted)
	}
	if got.WarsElsewhere != 1 {
		t.Errorf("the player's own war was counted as happening elsewhere (%d)", got.WarsElsewhere)
	}
}

func TestBuyingPremisesIsNotTheCityMovingThem(t *testing.T) {
	w := core.New(103)
	eyes := watch(w)
	var got WorldMeasures

	var first string
	for _, l := range core.Locations {
		if w.Properties[l.ID] != nil {
			first = l.ID
			break
		}
	}
	if first == "" {
		t.Skip("no premises in this world")
	}

	// The player takes it: that is the player playing.
	was := w.Properties[first].Owner
	w.Properties[first].Owner = w.PlayerOrganizationID()
	eyes.changed(w, &got)
	if got.HoldingsChangedHands != 0 {
		t.Errorf("the player taking premises counted as the city moving them (%d)", got.HoldingsChangedHands)
	}

	// One family taking it off another is the city moving.
	w.Properties[first].Owner = was
	eyes.changed(w, &got)
	if got.HoldingsChangedHands != 0 {
		t.Errorf("premises coming back off the player counted as the city moving them (%d)", got.HoldingsChangedHands)
	}
	w.Properties[first].Owner = "somebody-else"
	eyes.changed(w, &got)
	if got.HoldingsChangedHands != 1 {
		t.Errorf("a holding changing hands between families counted as %d", got.HoldingsChangedHands)
	}
}
