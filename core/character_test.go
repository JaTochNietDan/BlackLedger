package core

import "testing"

func TestATemperamentIsTheSamePersonEveryTime(t *testing.T) {
	a := New(7)
	b := New(7)
	for i := range a.NPCs {
		if TemperamentOf(&a.NPCs[i]).ID != TemperamentOf(&b.NPCs[i]).ID {
			t.Fatalf("%s was two different people in two loads of the same city", a.NPCs[i].Name)
		}
	}
	// And it is not stored, so nothing about it can drift out of a save.
	one := a.NPCs[0]
	before := TemperamentOf(&one).ID
	one.Ambition, one.Skill, one.Rank, one.Trust = 1, 1, 1, 1
	if TemperamentOf(&one).ID != before {
		t.Fatal("who somebody is changed when their numbers did")
	}
}

func TestTheCityIsNotAllOneKindOfPerson(t *testing.T) {
	seen := map[string]int{}
	for seed := uint32(1); seed <= 200; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		for i := range w.NPCs {
			seen[TemperamentOf(&w.NPCs[i]).ID]++
		}
	}
	if len(seen) < len(temperaments) {
		t.Fatalf("only %d of %d kinds of person ever appeared: %v", len(seen), len(temperaments), seen)
	}
	total := 0
	for _, n := range seen {
		total += n
	}
	for id, n := range seen {
		if n*2 > total {
			t.Fatalf("%s people were %d of %d, so the city is mostly one kind of person", id, n, total)
		}
	}
	t.Logf("across 200 cities, the people in them were: %v", seen)
}

func TestTemperamentChangesWhatSomebodyDoes(t *testing.T) {
	w := New(3)
	base := &NPC{ID: "x", Name: "Test Person", Ambition: 80, Skill: 60}
	hot, careful := *base, *base
	// Names chosen for the temperament they produce, then asserted rather than
	// assumed, so this test fails loudly if the derivation changes.
	hot.Name, careful.Name = nameWith(t, "hot"), nameWith(t, "careful")
	if w.Nerve(&hot) <= w.Nerve(&careful) {
		t.Fatalf("a hot-headed person acted less readily (%.3f) than a careful one (%.3f)", w.Nerve(&hot), w.Nerve(&careful))
	}
	if w.Poise(&hot) >= w.Poise(&careful) {
		t.Fatalf("a hot-headed person did it better (%d) than a careful one (%d)", w.Poise(&hot), w.Poise(&careful))
	}
	if TemperamentOf(&hot).Greed >= TemperamentOf(&NPC{Name: nameWith(t, "greedy")}).Greed {
		t.Fatal("a grasping person took no more than a hot-headed one")
	}
}

// nameWith finds a name that produces a given temperament, so the tests can
// speak about kinds of people without hard-coding a hash.
func nameWith(t *testing.T, id string) string {
	t.Helper()
	for i := 0; i < 4000; i++ {
		name := "Probe " + itoa(i)
		if TemperamentOf(&NPC{Name: name}).ID == id {
			return name
		}
	}
	t.Fatalf("no short name produced a %s person", id)
	return ""
}

func TestALoyalPersonWillNotMoveAgainstTheirOwn(t *testing.T) {
	w, a, _ := quarrel(t)
	// Give the quarrel to two people in the same organization, one of them
	// loyal, and let it run long enough that anybody else would have acted.
	var peer *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Faction == a.Faction && n.ID != a.ID && !n.Dead {
			peer = n
			break
		}
	}
	if peer == nil {
		t.Fatal("the organization had one person in it")
	}
	a.Name = nameWith(t, "loyal")
	a.Ambition = 100
	w.Resent(a.ID, peer.ID, GrudgeCap, "an old debt")
	for i := 0; i < 200; i++ {
		w.SettleGrudges()
	}
	if peer.Dead || a.Dead {
		t.Fatal("a loyal person settled it with their own people anyway")
	}
	if settledBetween(w, a, peer) {
		t.Fatal("the grievance was acted on")
	}

	// The same grievance against somebody outside their organization is a
	// different matter entirely.
	outsider := (*NPC)(nil)
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Faction != "" && n.Faction != a.Faction && !n.Dead {
			outsider = n
			break
		}
	}
	w.Resent(a.ID, outsider.ID, GrudgeCap, "an older debt")
	for i := 0; i < 200 && !settledBetween(w, a, outsider); i++ {
		w.SettleGrudges()
	}
	if !settledBetween(w, a, outsider) {
		t.Fatal("a loyal person would not act against anybody at all")
	}
}

func TestAStrangerIsAStranger(t *testing.T) {
	w := New(29)
	w.MigrateLivingWorld()
	w.Player.Contacts = 0
	var soldier *NPC
	for _, n := range w.People() {
		if n.Rank < RankLeader && n.Trust == 0 {
			soldier = n
			break
		}
	}
	if soldier == nil {
		t.Fatal("everybody in the city was already somebody")
	}
	if w.Known(soldier) || w.Dossier(soldier.ID) != nil {
		t.Fatalf("a stranger came with a file: %v", w.Dossier(soldier.ID))
	}
	// Whoever is at the top is a name everybody knows.
	for _, n := range w.People() {
		if n.Rank >= RankLeader {
			if !w.Known(n) || w.Dossier(n.ID) == nil {
				t.Fatalf("nobody had heard of %s, who runs an organization", n.Name)
			}
		}
	}
	w.MeetPerson(soldier.ID)
	entry := w.Dossier(soldier.ID)
	if entry == nil || entry["name"] != soldier.Name || entry["temperament"] == "" {
		t.Fatalf("meeting somebody told the player %v", entry)
	}
	if !w.hasRecord("You know " + soldier.Name + " now") {
		t.Fatal("meeting somebody was never mentioned")
	}
}

func TestWhatYouKnowIsWhatTheCityRecorded(t *testing.T) {
	w := New(29)
	w.MigrateLivingWorld()
	var person *NPC
	for _, n := range w.People() {
		if n.Rank >= RankLeader {
			person = n
			break
		}
	}
	entry := w.Dossier(person.ID)
	if len(entry["known_for"].([]string)) != 0 {
		t.Fatalf("a person nothing had happened to was known for %v", entry["known_for"])
	}
	w.Log(person.Name+" takes over something", "It happened.", "politics")
	entry = w.Dossier(person.ID)
	known := entry["known_for"].([]string)
	if len(known) != 1 || known[0] != person.Name+" takes over something" {
		t.Fatalf("what the player knows is %v", known)
	}
}

func TestTheCastIsOnlyPeopleYouKnow(t *testing.T) {
	w := New(29)
	w.MigrateLivingWorld()
	w.Player.Contacts = 0
	cast := w.Cast()
	if len(cast) == 0 {
		t.Fatal("the player had heard of nobody at all, not even the heads of families")
	}
	if len(cast) == len(w.People()) {
		t.Fatal("the player knew everybody in the city on their first morning")
	}
	// Enough contacts and the city stops being full of strangers.
	w.Player.Contacts = 5
	if len(w.Cast()) <= len(cast) {
		t.Fatal("building a network taught the player nothing about anybody")
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	out := ""
	for n > 0 {
		out = string(rune('0'+n%10)) + out
		n /= 10
	}
	return out
}
