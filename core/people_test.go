package core

import "testing"

func TestEveryOrganizationHasPeopleWhoCouldReplaceItsLeader(t *testing.T) {
	t.Parallel()
	w := New(1)
	for _, f := range w.Factions {
		members := w.Members(f.ID)
		if len(members) < 2 {
			t.Fatalf("%s has %d people; nobody could take over", f.Name, len(members))
		}
		if members[0].Rank != RankLeader || members[0].Name != f.Leader {
			t.Fatalf("%s: strongest standing is %q, but the leader is %q", f.Name, members[0].Name, f.Leader)
		}
		for _, m := range members {
			if m.Name == "" || m.Location == "" {
				t.Fatalf("%s has a person who is not really anywhere: %+v", f.Name, m)
			}
		}
	}
	// Nobody shares a name with anybody, including the player.
	seen := map[string]bool{w.Player.Name: true}
	for _, n := range w.NPCs {
		if seen[n.Name] {
			t.Fatal("two people share the name", n.Name)
		}
		seen[n.Name] = true
	}
}

func TestKillingALeaderPromotesSomebody(t *testing.T) {
	t.Parallel()
	w := New(2)
	before := *w.faction("bellandi") // copy: faction() returns a live pointer
	deputy := w.Members("bellandi")[1]
	if !w.Kill("vittorio", "Shot leaving the club.") {
		t.Fatal("the head of a family could not be killed")
	}
	if !w.NPC("vittorio").Dead {
		t.Fatal("the dead are still listed as living")
	}
	after := w.faction("bellandi")
	if after.Leader != deputy.Name {
		t.Fatalf("succession did not promote the strongest survivor: leader is %q, expected %q", after.Leader, deputy.Name)
	}
	if deputy.Rank != RankLeader {
		t.Fatal("the successor did not take the standing that goes with leading")
	}
	if after.Power >= before.Power {
		t.Fatal("a change at the top cost the organization nothing")
	}
	// The city hears about it.
	found := false
	for _, r := range w.History {
		if r.Title == "Vittorio Bellandi is dead" {
			found = true
		}
	}
	if !found {
		t.Fatal("a leader died and the city never heard")
	}
}

func TestAnOrganizationCanBeLeftWithNobody(t *testing.T) {
	t.Parallel()
	w := New(3)
	f := w.faction("russo")
	power := f.Power
	for _, m := range w.Members("russo") {
		w.Kill(m.ID, "Killed in the same night.")
	}
	if got := len(w.Members("russo")); got != 0 {
		t.Fatalf("%d people survived a killing of everyone", got)
	}
	if w.faction("russo").Power >= power {
		t.Fatal("losing everyone cost the organization nothing")
	}
}

func TestTheDeadDoNotActOrDieTwice(t *testing.T) {
	t.Parallel()
	w := New(4)
	if !w.Kill("elena", "An accident that was not one.") {
		t.Fatal("could not kill a living person")
	}
	if w.Kill("elena", "Again.") {
		t.Fatal("a dead person was killed a second time")
	}
	for _, p := range w.People() {
		if p.ID == "elena" {
			t.Fatal("the dead are still counted among the living")
		}
	}
	if w.Kill("nobody-at-all", "Nothing.") {
		t.Fatal("killed somebody who does not exist")
	}
}

func TestLosingOrdinaryMembersStillCostsStrength(t *testing.T) {
	t.Parallel()
	w := New(5)
	f := w.faction("bellandi")
	power := f.Power
	soldier := w.Members("bellandi")[2]
	w.Kill(soldier.ID, "Found at the docks.")
	if w.faction("bellandi").Power >= power {
		t.Fatal("losing people cost the organization nothing")
	}
	// Killing someone who leads nothing does not change who leads.
	if w.faction("bellandi").Leader != "Vittorio Bellandi" {
		t.Fatal("an ordinary death changed the leadership")
	}
}

func TestNobodyMovesOnALeaderWhoIsWinning(t *testing.T) {
	t.Parallel()
	w := New(6) // full strength
	if w.ConsiderInternalMove(w.faction("bellandi")) {
		t.Fatal("a deputy moved on a leader whose organization was still strong")
	}
}

func TestAFailingOrganizationCanLoseItsLeaderFromInside(t *testing.T) {
	t.Parallel()
	moved, succeeded, failed := 0, 0, 0
	for seed := uint32(1); seed <= 400; seed++ {
		w := New(seed)
		f := w.faction("bellandi")
		f.Power = peak(f) / 3 // failing badly
		// Whoever moves is whoever is most willing, so the test does not name
		// them: it counts the dead before and after and checks that the right
		// side of it died.
		deputy := w.Members("bellandi")[1]
		deputy.Rank, deputy.Ambition = RankLieutenant, 90
		leaderBefore := f.Leader
		living := len(w.Members("bellandi"))
		if !w.ConsiderInternalMove(f) {
			continue
		}
		moved++
		switch {
		case w.faction("bellandi").Leader != leaderBefore:
			succeeded++
			if !w.NPC("vittorio").Dead {
				t.Fatal("the old leader survived being replaced by force")
			}
		default:
			failed++
			if !w.NPC("vittorio").Dead && len(w.Members("bellandi")) >= living {
				t.Fatal("a failed challenger walked away")
			}
		}
	}
	if moved == 0 {
		t.Fatal("no organization ever lost its leader from inside")
	}
	if succeeded == 0 || failed == 0 {
		t.Fatalf("internal moves are not a real gamble: %d succeeded, %d failed of %d", succeeded, failed, moved)
	}
	t.Logf("of %d internal moves, %d succeeded and %d cost the challenger their life", moved, succeeded, failed)
}

func TestViolenceReachesEveryRank(t *testing.T) {
	t.Parallel()
	ranks := map[int]int{}
	for seed := uint32(1); seed <= 600; seed++ {
		w := New(seed)
		if victim := w.casualty("bellandi"); victim != nil {
			ranks[victim.Rank]++
		}
	}
	if len(ranks) < 2 {
		t.Fatalf("casualties only ever fall on one rank: %v", ranks)
	}
	if ranks[RankLeader] == 0 {
		t.Fatal("a leader is never caught in violence against their own organization")
	}
	if ranks[RankSoldier] <= ranks[RankLeader] {
		t.Fatalf("the people sent to do the work are not the ones most often hit: %v", ranks)
	}
}

func TestPeopleKeepDistinctVoices(t *testing.T) {
	t.Parallel()
	w := New(31)
	used := map[string][]string{}
	for _, n := range w.People() {
		if n.Voice == "" {
			t.Fatal(n.Name, "has no voice")
		}
		used[n.Voice] = append(used[n.Voice], n.Name)
	}
	// There are more people in this city than there are voices, so voices are
	// shared. What matters is that no one voice covers the room: if a third of
	// the city sounds the same, nobody sounds like themselves.
	worst, worstVoice := 0, ""
	for voice, names := range used {
		if len(names) > worst {
			worst, worstVoice = len(names), voice
		}
	}
	if worst*3 > len(w.People()) {
		t.Fatalf("%d of %d people share the voice %s", worst, len(w.People()), worstVoice)
	}
	if len(used) < 8 {
		t.Fatalf("a city of %d people used %d voices", len(w.People()), len(used))
	}
	// A person keeps their voice as the city changes around them.
	mara := w.NPC("mara")
	voice := mara.Voice
	for i := 0; i < 5; i++ {
		w.AddMember("bellandi", "Soldier", RankSoldier, "club")
	}
	if w.NPC("mara").Voice != voice {
		t.Fatal("an established character's voice changed")
	}
}

func TestABreakawayBringsItsOwnPeople(t *testing.T) {
	t.Parallel()
	w := embattled(41)
	if !w.Splinter(w.faction("bellandi")) {
		t.Skip("no breakaway for this seed")
	}
	born := w.Factions[len(w.Factions)-1]
	members := w.Members(born.ID)
	if len(members) < 2 {
		t.Fatalf("a new organization has %d people; nobody could replace its leader", len(members))
	}
	if members[0].Name != born.Leader {
		t.Fatal("the new organization's strongest standing is not its leader")
	}
	if w.NPC(members[0].ID).Voice == "" {
		t.Fatal("the new leader has no voice")
	}
}

func TestVoicesMatchHowPeopleAreNamed(t *testing.T) {
	t.Parallel()
	womens := map[string]bool{}
	for _, v := range womensVoices {
		womens[v] = true
	}
	mens := map[string]bool{}
	for _, v := range mensVoices {
		mens[v] = true
	}
	checked := 0
	for seed := uint32(1); seed <= 60; seed++ {
		w := New(seed)
		for i := 0; i < 4; i++ {
			w.AddMember("bellandi", "Soldier", RankSoldier, "club")
		}
		for _, n := range w.People() {
			first := n.Name
			if idx := len(first); idx > 0 {
				for j, r := range n.Name {
					if r == ' ' {
						first = n.Name[:j]
						break
					}
				}
			}
			feminine := false
			for _, name := range womensFirstNames {
				if name == first {
					feminine = true
				}
			}
			masculine := false
			for _, name := range mensFirstNames {
				if name == first {
					masculine = true
				}
			}
			if !feminine && !masculine {
				continue // the authored cast keeps its original casting
			}
			checked++
			if feminine && !womens[n.Voice] {
				t.Fatalf("%s was given %s", n.Name, n.Voice)
			}
			if masculine && !mens[n.Voice] {
				t.Fatalf("%s was given %s", n.Name, n.Voice)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no generated person was checked")
	}
	t.Logf("checked %d generated people", checked)
}
