package core

import "testing"

func quarrel(t *testing.T) (*World, *NPC, *NPC) {
	t.Helper()
	w := New(83)
	w.MigrateLivingWorld()
	var a, b *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Faction == "" {
			continue
		}
		if a == nil {
			a = n
			continue
		}
		if n.Faction != a.Faction {
			b = n
			break
		}
	}
	if a == nil || b == nil {
		t.Fatal("the city had no two people in different organizations")
	}
	return w, a, b
}

func TestNobodyResentsNobody(t *testing.T) {
	t.Parallel()
	w, a, _ := quarrel(t)
	w.Resent(a.ID, a.ID, 40, "themselves")
	w.Resent(a.ID, "somebody-who-does-not-exist", 40, "nothing")
	w.Resent("", a.ID, 40, "nothing")
	w.Resent(a.ID, "leo", 0, "nothing at all")
	if len(w.Grudges) != 0 {
		t.Fatalf("%d grudges out of nothing: %+v", len(w.Grudges), w.Grudges)
	}
}

func TestAGrudgeAccumulatesAndFades(t *testing.T) {
	t.Parallel()
	w, a, b := quarrel(t)
	w.Resent(a.ID, b.ID, 20, "the first thing")
	w.Resent(a.ID, b.ID, 20, "the second thing")
	if len(w.Grudges) != 1 || w.Grudges[0].Weight != 40 {
		t.Fatalf("grudges: %+v", w.Grudges)
	}
	if w.Grudges[0].Because != "the second thing" {
		t.Fatal("the grudge forgot what it was most recently about")
	}
	if w.Grudges[0].Weight > GrudgeCap {
		t.Fatal("a grudge went past the cap")
	}
	for i := 0; i < 40; i++ {
		w.GrudgeDay()
	}
	if len(w.Grudges) != 0 {
		t.Fatalf("a grudge outlasted forty days: %+v", w.Grudges)
	}
}

func TestTheDeadCarryNothingAndAreOwedNothing(t *testing.T) {
	t.Parallel()
	w, a, b := quarrel(t)
	w.Resent(a.ID, b.ID, 50, "something")
	w.Resent(b.ID, a.ID, 50, "something else")
	a.Dead = true
	w.GrudgeDay()
	if len(w.Grudges) != 0 {
		t.Fatalf("%d grudges survived a death: %+v", len(w.Grudges), w.Grudges)
	}
}

func TestSomebodyWithAHeavyGrudgeActsOnIt(t *testing.T) {
	t.Parallel()
	const runs, days = 300, 10
	acted, killed := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, a, b := quarrel(t)
		w.WorldRNG = seed * 2654435761
		a.Ambition = 100
		w.Resent(a.ID, b.ID, GrudgeCap, "an old debt")
		for day := 0; day < days && !settledBetween(w, a, b); day++ {
			w.SettleGrudges()
		}
		if settledBetween(w, a, b) {
			acted++
		}
		if a.Dead || b.Dead {
			killed++
		}
	}
	if acted < runs*4/5 {
		t.Fatalf("somebody carrying the heaviest possible grudge acted on it in %d of %d cities inside %d days", acted, runs, days)
	}
	if killed == 0 {
		t.Fatal("nobody died over any of it, so this is a mood rather than a mechanic")
	}
	t.Logf("with the heaviest possible grudge and every reason to act: settled in %d of %d cities inside %d days, somebody died in %d", acted, runs, days, killed)
}

func TestNobodyActsOnSomethingSmall(t *testing.T) {
	t.Parallel()
	w, a, b := quarrel(t)
	a.Ambition = 100
	w.Resent(a.ID, b.ID, GrudgeActs-20, "a slight")
	for i := 0; i < 50; i++ {
		w.SettleGrudges()
	}
	if len(w.Grudges) != 1 || b.Dead || a.Dead {
		t.Fatal("a slight got somebody killed")
	}
}

func TestAKillingGivesTheOrganizationsAReason(t *testing.T) {
	t.Parallel()
	const runs, days = 300, 10
	settled, escalated, warred := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, a, b := quarrel(t)
		w.WorldRNG = seed * 2654435761
		a.Ambition, a.Skill = 100, 95
		c := w.Conflict(a.Faction, b.Faction)
		if c == nil {
			t.Fatal("the two organizations had no relationship at all")
		}
		before, state := c.Hostility, c.State
		w.Resent(a.ID, b.ID, GrudgeCap, "an old debt")
		for day := 0; day < days && !settledBetween(w, a, b); day++ {
			w.SettleGrudges()
		}
		if !settledBetween(w, a, b) {
			continue
		}
		settled++
		after := w.Conflict(a.Faction, b.Faction)
		if after.Hostility > before {
			escalated++
		}
		if after.State != state {
			warred++
		}
	}
	if settled == 0 {
		t.Fatal("nothing was ever settled")
	}
	if escalated != settled {
		t.Fatalf("%d of %d settled grudges moved the organizations apart, so private history stays private", escalated, settled)
	}
	t.Logf("%d grudges settled between two organizations: every one raised the hostility between them, and %d changed the state of their relationship outright", settled, warred)
}

func TestWarIsItsOwnReasonAndStopsBeingOneWhenItEnds(t *testing.T) {
	t.Parallel()
	w, a, b := quarrel(t)
	c := w.Conflict(a.Faction, b.Faction)
	c.State = "war"
	if w.warGrudge(a, b) == 0 {
		t.Fatal("being at war was not a reason to move against somebody")
	}
	c.State = "cold"
	if w.warGrudge(a, b) != 0 {
		t.Fatal("a war that ended was still a reason")
	}
	// And never against your own people.
	for i := range w.NPCs {
		if peer := &w.NPCs[i]; peer.Faction == a.Faction && peer.ID != a.ID {
			if w.warGrudge(a, peer) != 0 {
				t.Fatal("a war was a reason to move against your own organization")
			}
			return
		}
	}
}

func TestTheSaveStaysBounded(t *testing.T) {
	t.Parallel()
	w, _, _ := quarrel(t)
	living := []string{}
	for i := range w.NPCs {
		if !w.NPCs[i].Dead {
			living = append(living, w.NPCs[i].ID)
		}
	}
	for i, holder := range living {
		for j, against := range living {
			w.Resent(holder, against, 1+(i+j)%20, "something")
		}
	}
	// Every pair in the city, repeatedly, and the list still fits in a save.
	if len(w.Grudges) > MaxGrudges {
		t.Fatalf("%d grudges stored", len(w.Grudges))
	}
}

func TestYouOnlyHearAboutThisIfSomebodyTellsYou(t *testing.T) {
	t.Parallel()
	w, a, b := quarrel(t)
	w.Resent(a.ID, b.ID, GrudgeCap, "an old debt")
	w.Player.Contacts = 0
	if len(w.GrudgeSummary()) != 0 {
		t.Fatal("a stranger with no contacts knew who wanted who dead")
	}
	w.Player.Contacts = 3
	summary := w.GrudgeSummary()
	if len(summary) != 1 || summary[0]["holder"] != a.Name || summary[0]["against"] != b.Name {
		t.Fatalf("what came back was %v", summary)
	}
}

// settledBetween reports whether one person's grievance against another has
// been acted on, which is not the same as the city having no grudges left: a
// killing gives somebody else a reason of their own.
func settledBetween(w *World, holder, against *NPC) bool {
	for _, g := range w.Grudges {
		if g.Holder == holder.ID && g.Against == against.ID {
			return false
		}
	}
	return true
}
