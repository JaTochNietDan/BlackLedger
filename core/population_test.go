package core

import "testing"

func TestTheCityIsFullOfPeople(t *testing.T) {
	w := New(17)
	if len(w.People()) < 30 {
		t.Fatalf("a new city had %d people in it", len(w.People()))
	}
	if len(w.Civilians()) < StreetCount-2 {
		t.Fatalf("%d people in the city answered to nobody", len(w.Civilians()))
	}
	for i := range w.Factions {
		if n := len(w.Members(w.Factions[i].ID)); n < FamilySize {
			t.Fatalf("%s had %d people in it", w.Factions[i].Name, n)
		}
		if lieutenants(w.Members(w.Factions[i].ID)) < 2 {
			t.Fatalf("%s had %d lieutenants", w.Factions[i].Name, lieutenants(w.Members(w.Factions[i].ID)))
		}
	}
}

func TestEverybodyHasTheirOwnName(t *testing.T) {
	w := New(17)
	seen := map[string]bool{}
	for _, n := range w.People() {
		if n.Name == "" {
			t.Fatalf("%s had no name", n.ID)
		}
		if seen[n.Name] {
			t.Fatalf("two people were both called %s", n.Name)
		}
		seen[n.Name] = true
		if n.Location == "" || n.Role == "" {
			t.Fatalf("%s was nowhere doing nothing", n.Name)
		}
	}
}

func TestPopulatingIsIdempotent(t *testing.T) {
	w := New(17)
	before := len(w.NPCs)
	w.Populate()
	w.Populate()
	if len(w.NPCs) != before {
		t.Fatalf("the city grew from %d to %d by being counted", before, len(w.NPCs))
	}
}

func TestAnOrganizationShortOfPeopleTakesThemOffTheStreet(t *testing.T) {
	w := New(17)
	var seed uint32 = 17
	w.WorldRNG = seed * 2654435761
	f := w.faction("bellandi")
	// Cut them down to three.
	for _, m := range w.Members(f.ID)[3:] {
		m.Faction, m.Rank, m.Role = "", RankAssociate, "Docker"
		m.Location = "docks"
	}
	before, street := len(w.Members(f.ID)), len(w.Civilians())
	recruited := false
	for day := 0; day < 200 && !recruited; day++ {
		w.RecruitDay()
		recruited = len(w.Members(f.ID)) > before
	}
	if !recruited {
		t.Fatal("an organization at three men never took anybody on")
	}
	if len(w.Civilians()) >= street && len(w.NPCs) <= MaxPeople {
		t.Fatal("the recruit came out of thin air rather than off the street")
	}
	// And it stops when it is up to strength.
	for day := 0; day < 500; day++ {
		w.RecruitDay()
	}
	if n := len(w.Members(f.ID)); n > FamilySize {
		t.Fatalf("an organization recruited itself to %d", n)
	}
}

func TestTheSaveStaysBoundedHoweverLongItRuns(t *testing.T) {
	const runs, days = 60, 200
	worst := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.WorldRNG = seed * 2654435761
		for day := 0; day < days; day++ {
			w.FactionTurn()
			w.FamilyDay()
			w.PeopleDay()
			w.GrudgeDay()
			w.SettleGrudges()
			w.RecruitDay()
			w.PrunePeople()
			w.Minute += 1440
		}
		worst = max(worst, len(w.NPCs))
	}
	if worst > MaxPeople {
		t.Fatalf("a city held %d people, past the %d cap", worst, MaxPeople)
	}
	t.Logf("the most people any of %d cities held after %d days: %d of a possible %d", runs, days, worst, MaxPeople)
}

func TestTheDeadAreForgottenOnlyWhenNothingNeedsThem(t *testing.T) {
	w := New(17)
	// Fill the city so pruning is allowed to run at all.
	for len(w.NPCs) <= MaxPeople*3/4 {
		if w.AddCivilian() == nil {
			break
		}
	}
	// Identities are held as ids rather than pointers: pruning compacts the
	// list in place, so a pointer taken before it means a different person
	// afterwards.
	victim, other := w.Civilians()[0].ID, w.Civilians()[1].ID
	w.Resent(other, victim, 40, "something")
	w.NPC(victim).Dead = true
	w.PrunePeople()
	if w.NPC(victim) == nil {
		t.Fatal("somebody with a grievance against them was forgotten")
	}
	w.Grudges = nil
	w.PrunePeople()
	if w.NPC(victim) != nil {
		t.Fatal("a dead person nothing referred to was kept")
	}
	// The record of the death survives them.
	if w.isCrew("leo") != true {
		t.Fatal("the player's own people are not protected from pruning")
	}
}

func TestAStrangerIsStillAStrangerInACrowd(t *testing.T) {
	w := New(17)
	w.Player.Contacts = 0
	known := len(w.Cast())
	if known == 0 {
		t.Fatal("the player had heard of nobody at all")
	}
	if known > len(w.People())/3 {
		t.Fatalf("a stranger knew %d of %d people on their first morning", known, len(w.People()))
	}
}
