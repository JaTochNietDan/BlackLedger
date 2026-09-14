package core

import (
	"fmt"
	"testing"
)

func proprietor(t *testing.T) *World {
	t.Helper()
	w := New(223)
	w.MigrateLivingWorld()
	w.Player.Cash, w.Player.Health = 20000, 100
	return w
}

func own(w *World, ids ...string) {
	for _, id := range ids {
		w.Properties[id].Owner = fmt.Sprintf("player:%d", w.Life)
	}
}

func TestAManWithOneShopIsStillAMan(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry")
	w.Player.Respect = 80
	w.OrganizationDay()
	if w.Incorporated() {
		t.Fatal("one shop made an organization")
	}
	own(w, "garage")
	w.Player.Respect = OrganizationStanding - 1
	w.OrganizationDay()
	if w.Incorporated() {
		t.Fatal("two shops and no name made an organization")
	}
	d := w.PlayerOrganizationDescription()
	if d["named"] != false || len(d["needs"].([]string)) == 0 {
		t.Fatalf("the interface was told %v", d)
	}
}

func TestExplicitFormationCreatesAnOrganization(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect = OrganizationStanding
	w.Incorporate()
	w.OrganizationDay()
	f := w.PlayerOrganization()
	if f == nil {
		t.Fatal("the city did not file them")
	}
	if f.ID != w.PlayerOrganizationID() || f.Leader != w.Player.Name || f.Power <= 0 {
		t.Fatalf("the entry was %+v", f)
	}
	if !w.hasRecord("They have started calling you something") {
		t.Fatal("nobody told the player")
	}
	if !w.hasNewsKind("politics") {
		t.Fatal("the city did not notice")
	}
	// And their holdings are the organization's holdings, because they always
	// were: the ownership string is the organization's id.
	if len(w.FamilyHoldings(f.ID)) != 2 {
		t.Fatalf("the organization held %d premises", len(w.FamilyHoldings(f.ID)))
	}
	// Filing them twice does not file them twice.
	before := len(w.Factions)
	w.Incorporate()
	w.OrganizationDay()
	w.Incorporate()
	w.OrganizationDay()
	if len(w.Factions) != before {
		t.Fatalf("the city listed them %d times", len(w.Factions)-before+1)
	}
}

func TestStrengthComesFromWhatTheyActuallyHave(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect = 30
	w.Incorporate()
	w.OrganizationDay()
	base := w.PlayerStrength()
	w.Player.Respect = 90
	if w.PlayerStrength() <= base {
		t.Fatal("a bigger name was worth nothing")
	}
	w.Player.Respect = 30
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 100}}
	if w.PlayerStrength() <= base {
		t.Fatal("somebody standing with them was worth nothing")
	}
	w.Player.Crew = nil
	w.Properties["laundry"].Condition = 10
	if w.PlayerStrength() >= base {
		t.Fatal("a wrecked holding was worth as much as a sound one")
	}
	if w.PlayerStrength() > 100 {
		t.Fatal("strength went past the ceiling")
	}
}

func TestEverybodyInTheCityHasAViewOnANewOrganization(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect = OrganizationStanding
	w.Factions[0].Goodwill = -60
	w.Incorporate()
	w.OrganizationDay()
	me := w.PlayerOrganizationID()
	for i := range w.Factions {
		other := w.Factions[i]
		if other.ID == me {
			continue
		}
		if w.Conflict(other.ID, me) == nil {
			t.Fatalf("%s had no relationship with the new organization", other.Name)
		}
	}
	// Somebody who already hated them starts further along than somebody who
	// did not.
	hostile := w.Conflict(w.Factions[0].ID, me)
	friendly := w.Conflict(w.Factions[1].ID, me)
	if hostile.Hostility <= friendly.Hostility {
		t.Fatalf("a family at -60 goodwill started at %d against %d", hostile.Hostility, friendly.Hostility)
	}
}

func TestARaidOnThePlayerTakesTheirOwnMoney(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect = OrganizationStanding
	w.Incorporate()
	w.OrganizationDay()
	me := w.PlayerOrganization()
	attacker := &w.Factions[0]
	attacker.Power = 100
	cash, condition := w.Player.Cash, w.Properties["laundry"].Condition
	for i := 0; i < 40 && w.Player.Cash == cash; i++ {
		w.contest(attacker, me)
	}
	if w.Player.Cash >= cash {
		t.Fatal("a raid on the player's holdings never took their money")
	}
	if w.Properties["laundry"].Condition >= condition && w.Properties["garage"].Condition >= 100 {
		t.Fatal("a raid damaged nothing")
	}
}

func TestARaidOnThePlayerCanReachTheirCrew(t *testing.T) {
	t.Parallel()
	found := false
	for seed := uint32(1); seed <= 300 && !found; seed++ {
		w := proprietor(t)
		w.WorldRNG = seed * 2654435761
		own(w, "laundry", "garage")
		w.Player.Respect = OrganizationStanding
		w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 90}}
		w.Incorporate()
		w.OrganizationDay()
		attacker := &w.Factions[0]
		attacker.Power = 100
		for i := 0; i < 20 && len(w.Player.Crew) > 0; i++ {
			w.contest(attacker, w.PlayerOrganization())
		}
		if len(w.Player.Crew) == 0 {
			found = true
			if n := w.NPC("leo"); n != nil && !n.Dead {
				t.Fatal("the crewman left the crew and lived")
			}
		}
	}
	if !found {
		t.Fatal("no raid in 300 ever reached the player's crew")
	}
}

func TestTheCityDoesNotRepairThePlayersPremisesForFree(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect = OrganizationStanding
	w.Incorporate()
	w.OrganizationDay()
	w.Properties["laundry"].Condition = 40
	for i := 0; i < 10; i++ {
		w.FamilyDay()
	}
	if w.Properties["laundry"].Condition != 40 {
		t.Fatalf("the player's laundry repaired itself to %d%%", w.Properties["laundry"].Condition)
	}
}

func TestAnOrganizationEndsWithThePersonItBelongedTo(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect = OrganizationStanding
	w.Incorporate()
	w.OrganizationDay()
	before := len(w.Factions)
	w.Player.Alive = false
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	if len(next.Factions) != before-1 {
		t.Fatalf("the city still lists %d organizations", len(next.Factions))
	}
	for _, c := range next.Conflicts {
		if c.A == w.PlayerOrganizationID() || c.B == w.PlayerOrganizationID() {
			t.Fatal("a quarrel with a dead man survived him")
		}
	}
	if next.Incorporated() {
		t.Fatal("the next arrival inherited an organization")
	}
}

func TestABadMonthDoesNotEndTheirOrganization(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect = OrganizationStanding
	w.Incorporate()
	w.OrganizationDay()
	w.Properties["laundry"].Owner = "independent"
	w.Properties["garage"].Owner = "independent"
	for i := 0; i < 30; i++ {
		w.FactionTurn()
	}
	if !w.Incorporated() {
		t.Fatal("losing everything ended an organization whose owner was still alive")
	}
}
