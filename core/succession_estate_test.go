package core

import "testing"

func testator(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect, w.Player.Contacts = OrganizationStanding, 3
	w.Incorporate()
	w.OrganizationDay()
	var member *NPC
	for _, n := range w.Civilians() {
		if IsOfficial(n.ID) || w.isRoleHolder(n) {
			continue
		}
		w.Player.Location = n.Location
		if w.SignOn(n.ID) == nil {
			member = n
			break
		}
	}
	if member == nil {
		t.Fatal("nobody would sign on")
	}
	return w, member
}

func TestWhatYouBuiltOutlivesYouIfAnybodyIsLeft(t *testing.T) {
	t.Parallel()
	w, member := testator(t)
	name := member.Name
	// Measured before they die: afterwards the player holds nothing and their
	// strength is not the number the estate should be compared against.
	strength := w.PlayerStrength()
	w.Die("Shot on the way home.")

	f := w.faction("estate:" + member.ID)
	if f == nil {
		t.Fatal("nobody took it over")
	}
	if f.Leader != name || f.Name != name+"'s people" {
		t.Fatalf("the organization was %+v", f)
	}
	if len(w.FamilyHoldings(f.ID)) != 2 {
		t.Fatalf("it kept %d of the two premises", len(w.FamilyHoldings(f.ID)))
	}
	if member.Faction != f.ID || member.Rank != RankLeader {
		t.Fatalf("the successor was %+v", member)
	}
	if f.Power >= strength {
		t.Fatalf("it was worth %d without the man it was built around, against %d with him", f.Power, strength)
	}
	if !w.hasRecord(name + " is running it now") {
		t.Fatal("nobody was told")
	}
	if !w.hasNewsKind("politics") {
		t.Fatal("the city did not notice")
	}
}

func TestTheCitysOpinionOfTheManCarriesToTheThingHeLeft(t *testing.T) {
	t.Parallel()
	w, member := testator(t)
	me := w.PlayerOrganizationID()
	c := w.Conflict(w.Factions[0].ID, me)
	c.Hostility, c.State = 70, "war"
	rival := w.Factions[0].ID
	w.Die("Shot.")

	estate := "estate:" + member.ID
	if w.Conflict(rival, me) != nil && w.Conflict(rival, me).Hostility == 70 {
		t.Fatal("the quarrel stayed with a man who is dead")
	}
	after := w.Conflict(rival, estate)
	if after == nil || after.Hostility != 70 || after.State != "war" {
		t.Fatalf("the quarrel carried as %+v", after)
	}
}

func TestNobodyIsLeftAnsweringToNothing(t *testing.T) {
	t.Parallel()
	// Somebody who had people leaves an organization; somebody who had none
	// leaves nothing, and nobody is left pointing at an id that resolves to
	// nothing either way.
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect = OrganizationStanding
	w.Incorporate()
	w.OrganizationDay()
	me := w.PlayerOrganizationID()
	w.Die("Shot.")
	if w.faction("estate:") != nil {
		t.Fatal("an organization with nobody in it was inherited by nobody")
	}
	for i := range w.NPCs {
		if w.NPCs[i].Faction == me {
			t.Fatalf("%s still answers to a dead man", w.NPCs[i].Name)
		}
	}
	// And the premises fall to the estate as they always did.
	held := 0
	for _, l := range Locations {
		if w.Properties[l.ID].Owner == "former:"+w.Player.Name {
			held++
		}
	}
	if held != 2 {
		t.Fatalf("the estate holds %d premises", held)
	}
}

func TestTheNextArrivalMeetsWhatTheLastOneLeft(t *testing.T) {
	t.Parallel()
	w, member := testator(t)
	w.Die("Shot.")
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	estate := next.faction("estate:" + member.ID)
	if estate == nil {
		t.Fatal("the organization did not survive into the next life")
	}
	if len(next.FamilyHoldings(estate.ID)) != 2 {
		t.Fatal("it did not keep its premises")
	}
	if next.faction(w.PlayerOrganizationID()) != nil {
		t.Fatal("the dead man's own entry survived")
	}
	if next.Incorporated() {
		t.Fatal("the next arrival started as an organization")
	}
	// And it is a thing in the city like any other: it can be dealt with,
	// enquired about, and fought.
	if next.Intelligence(estate.ID) < 0 {
		t.Fatal("it could not be asked about")
	}
	if len(next.Members(estate.ID)) == 0 {
		t.Fatal("it has nobody in it")
	}
}

func TestAnInheritedOrganizationLivesLikeAnyOther(t *testing.T) {
	t.Parallel()
	const runs, days = 150, 90
	survived, fought, gone := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, member := testator(t)
		w.WorldRNG = seed * 2654435761
		w.Die("Shot.")
		next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life"})
		if err != nil {
			t.Fatal(err)
		}
		id := "estate:" + member.ID
		for day := 0; day < days; day++ {
			next.FactionTurn()
			next.FamilyDay()
			next.PeopleDay()
			next.GrudgeDay()
			next.SettleGrudges()
			next.Minute += 1440
		}
		if f := next.faction(id); f != nil {
			survived++
			if len(next.FamilyHoldings(id)) == 0 {
				fought++
			}
		} else {
			gone++
		}
	}
	if survived == 0 {
		t.Fatal("no inherited organization ever lasted a season")
	}
	if gone == 0 && fought == 0 {
		t.Fatal("an inherited organization was never touched by anything")
	}
	t.Logf("%d inherited organizations over %d days: %d still standing, %d of those holding nothing, %d gone entirely", runs, days, survived, fought, gone)
}
