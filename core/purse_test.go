package core

import "testing"

// People in this city had no money. What somebody was carrying was worked out
// at the moment they were robbed, from their rank and their family's cash, and
// nothing was ever taken off the person — so the same man could be robbed every
// day for a year and be carrying the same amount every time.

func robbable(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := New(67)
	w.Player.Cash, w.Player.Respect, w.Player.Health = 500, 60, 100
	w.Player.Location = "bar"
	w.NPCs = append(w.NPCs, NPC{ID: "mark", Name: "Otto Reiss", Location: "bar",
		Faction: "bellandi", Rank: RankSoldier, Trust: 30, Skill: 5})
	w.SettlePurses() // the test builds a person by hand rather than through the city
	w.MeetPerson("mark")
	n := w.NPC("mark")
	if n == nil {
		t.Fatal("nobody to rob")
	}
	return w, n
}

func TestRobbingAManLeavesHimWithLess(t *testing.T) {
	t.Parallel()
	w, mark := robbable(t)
	before := w.Pockets(mark)
	if before <= 0 {
		t.Fatalf("the man is carrying $%d, so there is nothing to take", before)
	}
	mark.Purse = before
	took := mark.Purse
	w.TakePurse(mark)
	after := w.Pockets(mark)
	if after >= before {
		t.Errorf("a man robbed of $%d is still carrying $%d, up from $%d", took, after, before)
	}
}

// And what he is carrying is his, not a number worked out about him. Two men of
// the same rank in the same family who have had different weeks are not
// carrying the same amount.
func TestTwoMenOfTheSameRankCanCarryDifferentAmounts(t *testing.T) {
	t.Parallel()
	w, a := robbable(t)
	w.NPCs = append(w.NPCs, NPC{ID: "other", Name: "Bruno Sala", Location: "bar",
		Faction: "bellandi", Rank: RankSoldier, Trust: 30})
	b := w.NPC("other")
	a.Purse, b.Purse = 400, 40
	if w.Pockets(a) == w.Pockets(b) {
		t.Errorf("a man with $400 and a man with $40 are both carrying $%d", w.Pockets(a))
	}
}

// The point of people having money of their own: starving a family reaches the
// people in it. Take everything the family earns from, and the men who answer
// to it stop being paid and run down to nothing.
func TestStarvingAFamilyEmptiesItsPeoplesPockets(t *testing.T) {
	t.Parallel()
	w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
	f := w.faction("vasco")
	f.Cash, f.Power, f.Peak = 3000, 60, 60
	for i := 0; i < 4; i++ {
		w.AddMember("vasco", "Soldier", RankSoldier, "market")
	}
	people := w.Members("vasco")
	if len(people) < 3 {
		t.Fatalf("the family has %d people, so there is nobody to go hungry", len(people))
	}

	// A fortnight of the family paying its way.
	for day := 0; day < 14; day++ {
		w.FamilyDay()
		w.PayTheCity()
	}
	comfortable := 0
	for _, n := range people {
		comfortable += n.Purse
	}
	if comfortable == 0 {
		t.Fatal("a family that met every payday has people carrying nothing")
	}

	// Now everything it earns from is taken.
	for _, id := range w.FamilyHoldings("vasco") {
		w.Properties[id].Owner = "player:1"
	}
	for day := 0; day < 45; day++ {
		w.FamilyDay()
		w.PayTheCity()
	}
	after, broke := 0, 0
	for _, n := range people {
		after += n.Purse
		if w.Broke(n) {
			broke++
		}
	}
	if after >= comfortable {
		t.Errorf("a family's people carry $%d between them after six weeks unpaid, against $%d when it was paying", after, comfortable)
	}
	if broke == 0 {
		t.Errorf("not one of %d people is broke after six weeks of missed paydays", len(people))
	}
	t.Logf("paid: $%d between %d people. Starved: $%d, %d of them broke", comfortable, len(people), after, broke)
}

// Officials are on the city's books, not a family's, so a family's bad year
// does not stop their salary.
func TestAnOfficialIsPaidWhoeverIsStruggling(t *testing.T) {
	t.Parallel()
	w := New(67)
	w.ensureOfficials()
	var official *NPC
	for i := range w.NPCs {
		if IsOfficial(w.NPCs[i].ID) {
			official = &w.NPCs[i]
			break
		}
	}
	if official == nil {
		t.Fatal("the city has no officials")
	}
	for i := range w.Factions {
		w.Factions[i].Short = 30
	}
	official.Purse = 0
	w.PayTheCity()
	if official.Purse <= 0 {
		t.Errorf("%s went unpaid because the families were struggling", official.Name)
	}
}
