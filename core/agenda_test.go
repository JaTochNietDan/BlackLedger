package core

import "testing"

func TestPeopleActOnTheirOwnAmbition(t *testing.T) {
	acted := 0
	for i := uint32(1); i <= 300; i++ {
		w := New(i * 2654435761)
		before := len(w.History)
		for day := 0; day < 20; day++ {
			w.Minute += 1440
			w.PeopleDay()
		}
		if len(w.History) > before {
			acted++
		}
	}
	t.Logf("of 300 cities over 20 days, %d saw somebody act on their own account", acted)
	if acted == 0 {
		t.Fatal("nobody in the city ever does anything of their own")
	}
	if acted == 300 {
		t.Fatal("something happens in every city every time, which is not a city, it is a treadmill")
	}
}

func TestTheOneAtTheTopDoesNotDoThisPersonally(t *testing.T) {
	w := New(601)
	leader := w.Members("bellandi")[0]
	if leader.Rank != RankLeader {
		t.Fatal("the test did not find a leader")
	}
	before := *leader
	for day := 0; day < 400; day++ {
		w.Minute += 1440
		w.PeopleDay()
	}
	after := w.NPC(before.ID)
	if after.Location != before.Location || after.Role != before.Role {
		t.Fatal("the head of a family went out robbing premises personally")
	}
}

func TestNobodyRobsTheirOwnOrganization(t *testing.T) {
	w := New(603)
	holdings := w.FamilyHoldings("bellandi")
	if len(holdings) == 0 {
		t.Fatal("no holdings to protect")
	}
	member := w.Members("bellandi")[1]
	member.Ambition, member.Skill = 100, 100
	before := map[string]int{}
	for _, id := range holdings {
		before[id] = w.Properties[id].Condition
	}
	for day := 0; day < 200; day++ {
		w.takeFromSomebody(member)
	}
	for _, id := range holdings {
		if w.Properties[id].Condition < before[id] {
			t.Fatalf("%s robbed premises belonging to their own organization", member.Name)
		}
	}
}

func TestThePlayersBusinessIsATargetLikeAnyOther(t *testing.T) {
	robbed, protectedRobbed := 0, 0
	for i := uint32(1); i <= 300; i++ {
		w := New(i * 2654435761)
		w.Properties["laundry"].Owner = "player:1"
		w.Player.Cash = 5000
		for day := 0; day < 40; day++ {
			w.Minute += 1440
			w.PeopleDay()
		}
		if w.Player.Cash < 5000 {
			robbed++
		}

		// The same city, with somebody watching the place.
		guarded := New(i * 2654435761)
		guarded.Properties["laundry"].Owner = "player:1"
		guarded.Player.Cash = 5000
		guarded.Player.Security = 3
		guarded.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 90}}
		for day := 0; day < 40; day++ {
			guarded.Minute += 1440
			guarded.PeopleDay()
		}
		if guarded.Player.Cash < 5000 {
			protectedRobbed++
		}
	}
	t.Logf("of 300 cities over 40 days: the player was robbed in %d, and in %d with security and crew", robbed, protectedRobbed)
	if robbed == 0 {
		t.Fatal("the player's business is never a target")
	}
	if robbed == 300 {
		t.Fatal("owning a business means being robbed in every campaign")
	}
	if protectedRobbed >= robbed {
		t.Fatal("security and a loyal crew did not make the player a harder mark")
	}
}

func TestSomebodyWithNothingCanMakeThemselvesSomebody(t *testing.T) {
	w := New(607)
	// An unaffiliated person with the makings of an operator, and premises
	// standing without an owner.
	w.Properties["garage"].Owner = "former:Alex Varga"
	w.Properties["garage"].Income = 24
	person := w.NPC("leo")
	person.Faction, person.Ambition, person.Skill = "", 80, 70
	if !w.claimPremises(person) {
		t.Fatal("a capable unaffiliated person could not take unheld premises")
	}
	if person.Rank < RankSoldier {
		t.Fatal("taking over premises did not raise their standing")
	}
	// Never the player's, and never an organization's.
	w.Properties["laundry"].Owner = "player:1"
	other := w.NPC("mara")
	other.Faction, other.Ambition, other.Skill = "", 90, 90
	for i := 0; i < 50; i++ {
		w.claimPremises(other)
	}
	if !w.Own("laundry") {
		t.Fatal("somebody took the player's business by walking into it")
	}
	if w.Properties["club"].Owner != "bellandi" {
		t.Fatal("somebody walked into a family's premises and took them")
	}
}
