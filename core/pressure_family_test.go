package core

import "testing"

// The business demand chose which family was collecting by position in the
// faction list — index 0 for the player's first district, index 1 for anything
// beyond it. Dissolve removes a family from that slice, and the living world
// destroys seventeen families across twenty long campaigns. So the two indices
// stopped meaning the two families they were written for as soon as the city
// did what it is built to do, and with one family left they are not an index
// into anything.

func pressured(t *testing.T, w *World, place string) {
	t.Helper()
	w.Properties[place].Owner = "player:1"
	w.Properties[place].Income = 30
	w.Player.Location = place
	w.NextPressure = w.Minute
	w.Event = nil
}

func TestADemandComesFromAFamilyThatStillExists(t *testing.T) {
	w := New(9)
	// The city does what it does: one of the two original families is gone.
	w.Dissolve("bellandi")
	if len(w.Factions) != 1 {
		t.Fatalf("expected one family left, got %d", len(w.Factions))
	}
	pressured(t, w, "club")
	w.BusinessPressure()
	if w.Event == nil {
		t.Skip("no demand landed")
	}
	if w.faction(w.Event.Actor) == nil {
		t.Fatalf("the demand comes from %q, which is not a family in this city", w.Event.Actor)
	}
}

func TestTheLastFamilyStandingDoesNotCrashTheCity(t *testing.T) {
	w := New(9)
	w.Dissolve("bellandi")
	// A second-district premises: the old code indexed Factions[1], and there
	// is no Factions[1] any more.
	place := ""
	for i := range Locations {
		if Locations[i].District > 0 {
			place = Locations[i].ID
			break
		}
	}
	if place == "" {
		t.Fatal("no premises outside the first district")
	}
	w.District = 2
	pressured(t, w, place)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("the city panicked collecting from %s: %v", place, r)
		}
	}()
	w.BusinessPressure()
}

// Whose street it is, is a question about ground. A family that has taken a
// district over collects there; the family that used to has no claim left.
func TestTheFamilyHoldingTheStreetIsTheOneThatCollects(t *testing.T) {
	w := New(9)
	club, _ := PlaceByID("club")
	// Russo takes the whole district off Bellandi.
	for i := range Locations {
		if Locations[i].District == club.District {
			w.Properties[Locations[i].ID].Owner = "russo"
		}
	}
	if got := w.Claimants()[club.District]; got == nil || got.ID != "russo" {
		t.Fatalf("Russo holds the whole district and %v is still collecting", got)
	}
	// And with no families left, nobody collects anywhere.
	w.Dissolve("bellandi")
	w.Dissolve("russo")
	if claims := w.Claimants(); len(claims) != 0 {
		t.Fatalf("demands are still coming from %v, and there are no families", claims)
	}
}

// Buying one family off must not silence the city. Districts the player
// operates in get different claimants while there are families to spare, which
// is what keeps a truce with one from suppressing the other's demand.
func TestSeparateDistrictsHaveSeparateClaimants(t *testing.T) {
	w := New(9)
	claims := w.Claimants()
	if len(claims) < 2 {
		t.Fatalf("the city has %d districts claimed", len(claims))
	}
	seen := map[string]int{}
	for _, f := range claims {
		seen[f.ID]++
	}
	if len(seen) < 2 {
		t.Fatalf("one family claims every district: %v", seen)
	}
}
