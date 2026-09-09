package core

import "testing"

// Every authored scene in the game names its speaker by the job — the fixer,
// the detective — so that when the person holding it dies the next one speaks.
// Two did not. The audience and the business demand hardcoded "vittorio" and
// "elena", the two people who happened to lead the two families on the first
// morning. Families change hands: a leader is killed, the strongest survivor
// takes over, and the city goes on. Those two scenes did not go on with it.

func leaderless(t *testing.T, faction string) *World {
	t.Helper()
	w := New(9)
	// Give the family somebody who can inherit it, then kill the man at the top.
	f := w.faction(faction)
	if f == nil {
		t.Fatalf("no such family %q", faction)
	}
	w.NPCs = append(w.NPCs, NPC{
		ID: "heir", Name: "Rosa Marchetti", Role: "Soldier", Faction: faction,
		Location: "club", Rank: RankLieutenant, Ambition: 60, Skill: 70,
	})
	old := w.Leader(faction)
	if old == nil {
		t.Fatal("the family had no leader to begin with")
	}
	w.Kill(old.ID, "Shot on the steps.")
	if now := w.Leader(faction); now == nil || now.ID == old.ID {
		t.Fatalf("nobody took over %s: %+v", faction, now)
	}
	return w
}

func TestAnAudienceIsWithWhoeverLeadsTheFamilyNow(t *testing.T) {
	for _, tc := range []struct{ faction, where string }{
		{"bellandi", "club"},
		{"russo", "garage"},
	} {
		w := leaderless(t, tc.faction)
		want := w.Leader(tc.faction)
		w.OpenAudience(tc.where)
		if w.Event == nil {
			t.Fatalf("%s: no audience opened", tc.faction)
		}
		if w.Event.Speaker != want.ID {
			spoke := w.NPC(w.Event.Speaker)
			name, dead := w.Event.Speaker, false
			if spoke != nil {
				name, dead = spoke.Name, spoke.Dead
			}
			t.Fatalf("%s leads %s now, but the chair holds %s (dead: %v)",
				want.Name, tc.faction, name, dead)
		}
	}
}

func TestADemandComesFromWhoeverLeadsTheFamilyNow(t *testing.T) {
	w := leaderless(t, "bellandi")
	// Own something visibly earning, and let a demand land. Which family
	// collects is decided by who holds ground; what this asserts is that
	// whoever it turns out to be speaks with their current leader's voice,
	// not a predecessor's.
	w.Properties["club"].Owner = "player:1"
	w.Properties["club"].Income = 30
	for i := range w.Factions {
		w.Factions[i].Goodwill = 0
	}
	w.NextPressure = w.Minute
	w.BusinessPressure()
	if w.Event == nil || w.Event.Kind != "business_pressure" {
		t.Skip("no demand landed in this arrangement")
	}
	collecting := w.faction(w.Event.Actor)
	if collecting == nil {
		t.Fatalf("the demand comes from %q, which is not a family", w.Event.Actor)
	}
	want := w.Leader(collecting.ID)
	if want == nil {
		t.Fatalf("%s is collecting with nobody to lead it", collecting.Name)
	}
	if w.Event.Speaker != want.ID {
		spoke := w.NPC(w.Event.Speaker)
		name, dead := w.Event.Speaker, false
		if spoke != nil {
			name, dead = spoke.Name, spoke.Dead
		}
		t.Fatalf("%s leads the family now, but the demand comes from %s (dead: %v)",
			want.Name, name, dead)
	}
}

// A family with nobody left to lead it grants no audiences. The view falls back
// to the first person in the city when it cannot find a scene's speaker, so an
// audience that opened without a leader would seat a stranger at the table and
// let them set a family's terms.
func TestAFamilyWithNobodyLeftGrantsNoAudience(t *testing.T) {
	w := New(9)
	for i := range w.NPCs {
		if w.NPCs[i].Faction == "bellandi" {
			w.Kill(w.NPCs[i].ID, "The whole table, in one night.")
		}
	}
	if w.Leader("bellandi") != nil {
		t.Fatal("somebody still leads the family")
	}
	if w.AudienceReadiness("club") == "" {
		t.Fatal("the audience is still offered when nobody can grant it")
	}
	w.Event = nil
	w.OpenAudience("club")
	if w.Event != nil {
		t.Fatalf("a scene opened with nobody to speak: %+v", w.Event)
	}
}
