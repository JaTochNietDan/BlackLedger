package core

import "testing"

func embattled(seed uint32) *World {
	w := New(seed)
	f := w.faction("bellandi")
	f.Power = peak(f) / 2 // beaten down enough that people look elsewhere
	return w
}

func TestABreakawayTakesGroundStrengthAndMoney(t *testing.T) {
	w := embattled(5)
	parent := w.faction("bellandi")
	holdings := len(w.FamilyHoldings("bellandi"))
	power, cash := parent.Power, parent.Cash
	if !w.Splinter(parent) {
		t.Fatal("a weakened family with two holdings produced no breakaway")
	}
	parent = w.faction("bellandi")
	if got := len(w.FamilyHoldings("bellandi")); got != holdings-1 {
		t.Fatalf("the parent kept its ground: %d holdings, expected %d", got, holdings-1)
	}
	if parent.Power >= power || parent.Cash >= cash {
		t.Fatal("a breakaway cost the parent nothing")
	}
	if len(w.Factions) != 3 {
		t.Fatalf("expected a third organization, found %d", len(w.Factions))
	}
	born := w.Factions[2]
	if len(w.FamilyHoldings(born.ID)) != 1 {
		t.Fatal("the new organization holds nothing it broke away with")
	}
	if born.Name == "" || born.Leader == "" || born.Power <= 0 {
		t.Fatal("the new organization is not a real actor:", born)
	}
	// It has a leader the city and the director can address.
	if w.NPC(born.ID) == nil {
		t.Fatal("the new organization has no person at its head")
	}
	// Enemies from the first day, but not already at war.
	c := w.Conflict("bellandi", born.ID)
	if c == nil || c.Hostility < feudAt {
		t.Fatal("a betrayal did not create a quarrel")
	}
	if c.State == "war" {
		t.Fatal("a breakaway should not begin already at war")
	}
}

func TestAFamilyKeepsItsLastHolding(t *testing.T) {
	w := embattled(6)
	for _, id := range w.FamilyHoldings("bellandi")[1:] {
		w.Properties[id].Owner = "independent"
	}
	if w.Splinter(w.faction("bellandi")) {
		t.Fatal("someone walked away with the only thing the family had left")
	}
}

func TestHealthyFamiliesHoldTogether(t *testing.T) {
	w := New(7) // full strength, no war
	if w.Splinter(w.faction("bellandi")) {
		t.Fatal("a family at full strength and at peace still split")
	}
}

func TestNewOrganizationsHaveDistinctNames(t *testing.T) {
	w := embattled(8)
	seen := map[string]bool{}
	for _, f := range w.Factions {
		seen[f.Name] = true
	}
	for i := 0; i < 3; i++ {
		for j := range w.Factions {
			f := &w.Factions[j]
			f.Power = peak(f) / 2
			if w.Splinter(f) {
				break
			}
		}
	}
	for _, f := range w.Factions {
		if f.Name == "" {
			t.Fatal("an organization has no name")
		}
	}
	names := map[string]int{}
	for _, f := range w.Factions {
		names[f.Name]++
		if names[f.Name] > 1 {
			t.Fatal("two organizations share the name", f.Name)
		}
	}
}

func TestAnOrganizationHoldingNothingEventuallyEnds(t *testing.T) {
	w := embattled(9)
	if !w.Splinter(w.faction("bellandi")) {
		t.Skip("no breakaway for this seed")
	}
	born := &w.Factions[2]
	id := born.ID
	// Strip it and leave it powerless.
	for _, holding := range w.FamilyHoldings(id) {
		w.Properties[holding].Owner = "independent"
	}
	born.Power = 10
	w.dissolve()
	if w.faction(id) != nil {
		t.Fatal("an organization with nothing left is still listed")
	}
	for _, c := range w.Conflicts {
		if c.A == id || c.B == id {
			t.Fatal("a dissolved organization still has quarrels")
		}
	}
}

func TestAFamilyWithNobodyLeftToFightFractures(t *testing.T) {
	// Observed in a 1500-command campaign: Bellandi ended holding everything at
	// full strength with Russo reduced to a shell, and nothing further could
	// happen in that city.
	w := New(419)
	for _, id := range w.FamilyHoldings("russo") {
		w.Properties[id].Owner = "bellandi"
	}
	f := w.faction("bellandi")
	f.Power = peak(f) // healthy, and at peace
	if !w.splinterReady(f) {
		t.Fatal("a family that owns the whole city has no reason to fracture")
	}
	if !w.Splinter(f) {
		t.Fatal("an unopposed family did not fracture")
	}
	born := w.Factions[len(w.Factions)-1]
	if len(w.FamilyHoldings(born.ID)) == 0 {
		t.Fatal("the breakaway took no ground")
	}
	// And a healthy family with a live rival still holds together.
	rivals := New(421)
	healthy := rivals.faction("bellandi")
	healthy.Power = peak(healthy)
	if rivals.splinterReady(healthy) {
		t.Fatal("a healthy family with a living rival fractured anyway")
	}
}
