package core

import (
	"fmt"
	"strings"
	"testing"
)

// Four more states the reading harness had not been pointed at. Same checks,
// same several-seeds discipline: a state is only exercised when the world
// happens to enter it, so one city proves little.

// In a cell. Everywhere but the precinct offers nothing at all, which is
// deliberate, so this reads what the precinct and the ledger say instead.
func TestTheCityReadsWhileThePlayerIsHeld(t *testing.T) {
	read := 0
	for _, seed := range []uint32{31, 88, 219} {
		w := withFamily(seed, "estate:held", "Vera Kohl's people", "Vera Kohl")
		live(t, w, 1500)
		w.Player.Location = "precinct"
		w.Player.HeldUntil = w.Minute + 3*1440
		if !w.Held() {
			t.Fatalf("seed %d: the player is not in custody", seed)
		}
		cell := 0
		for _, a := range w.Actions("precinct") {
			cell++
			for _, flaw := range flaws(a.Label + ". " + a.Detail + ". " + a.Reason) {
				t.Errorf("seed %d, in custody: %s\n  %s / %s", seed, flaw, a.Label, a.Detail)
			}
		}
		if cell == 0 {
			t.Fatalf("seed %d: a cell with no way out of it", seed)
		}
		read += scan(t, w, fmt.Sprintf("seed %d, in custody", seed))
	}
	t.Logf("read %d passages", read)
}

// Dead, with whatever the city says about it afterwards.
func TestTheCityReadsAfterThePlayerDies(t *testing.T) {
	read := 0
	for _, seed := range []uint32{31, 88, 219} {
		w := withFamily(seed, "estate:gone", "Otto Reiss's people", "Otto Reiss")
		w.Properties["laundry"].Owner = "player:1"
		live(t, w, 1500)
		w.Player.Alive = true
		w.Die("Shot on the steps of the Mariner.")
		if w.Player.Alive {
			t.Fatalf("seed %d: the player did not die", seed)
		}
		for _, text := range w.Epitaph() {
			if s, ok := text.(string); ok {
				for _, flaw := range flaws(s) {
					t.Errorf("seed %d, the epitaph: %s\n  %s", seed, flaw, s)
				}
			}
		}
		read += scan(t, w, fmt.Sprintf("seed %d, after death", seed))
	}
	t.Logf("read %d passages", read)
}

// A crowded board: as many organizations as the city allows.
func TestACrowdedCityStillReads(t *testing.T) {
	read := 0
	for _, seed := range []uint32{31, 88, 219} {
		w := New(seed)
		w.District = 2
		w.Player.Cash, w.Player.Respect = 6000, 80
		names := []struct{ id, name, leader string }{
			{"estate:kohl", "Vera Kohl's people", "Vera Kohl"},
			{"splinter-Duarte", "the Duarte Brothers", "Bruno Duarte"},
			{"splinter-Amato", "Amato Crew", "Sal Amato"},
		}
		for _, n := range names {
			w.Factions = append(w.Factions, Faction{ID: n.id, Name: n.name, Leader: n.leader,
				Power: 50, Cash: 2500, Peak: 50, Reported: 50})
			w.NPCs = append(w.NPCs, NPC{ID: "boss-" + n.id, Name: n.leader, Role: "Head of " + n.name,
				Faction: n.id, Location: "market", Rank: RankLeader, Ambition: 70, Skill: 65})
		}
		if len(w.Factions) != 5 {
			t.Fatalf("seed %d: expected five organizations, got %d", seed, len(w.Factions))
		}
		for i := range w.Factions {
			for j := i + 1; j < len(w.Factions); j++ {
				w.Antagonize(w.Factions[i].ID, w.Factions[j].ID, 80)
			}
		}
		live(t, w, 4000)
		read += scan(t, w, fmt.Sprintf("seed %d, a crowded city", seed))
	}
	t.Logf("read %d passages", read)
}

// Two names that collide when the same sentence has to hold both.
func TestTwoFamiliesWithSimilarNamesStillRead(t *testing.T) {
	read := 0
	for _, seed := range []uint32{31, 88, 219} {
		w := New(seed)
		w.District = 2
		w.Player.Cash, w.Player.Respect = 6000, 80
		for _, n := range []struct{ id, name, leader string }{
			{"splinter-A", "the Duarte Brothers", "Bruno Duarte"},
			{"splinter-B", "the Duarte Boys", "Cleo Duarte"},
		} {
			w.Factions = append(w.Factions, Faction{ID: n.id, Name: n.name, Leader: n.leader,
				Power: 50, Cash: 2500, Peak: 50, Reported: 50})
			w.NPCs = append(w.NPCs, NPC{ID: "boss-" + n.id, Name: n.leader, Role: "Head of " + n.name,
				Faction: n.id, Location: "market", Rank: RankLeader, Ambition: 70, Skill: 65})
		}
		w.Antagonize("splinter-A", "splinter-B", 95)
		live(t, w, 4000)
		named := 0
		for _, text := range writings(w) {
			if strings.Contains(text, "Duarte") {
				named++
			}
		}
		if named == 0 {
			t.Fatalf("seed %d: neither family was ever mentioned", seed)
		}
		read += scan(t, w, fmt.Sprintf("seed %d, colliding names", seed))
	}
	t.Logf("read %d passages", read)
}

// The other half. Leads only touches a name that carries "the "; a name that
// does not must come through untouched, or every seeded family would be
// capitalised twice or mangled.
func TestLeadsLeavesOrdinaryNamesAlone(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"Bellandi Family", "Bellandi Family"},
		{"Russo Outfit", "Russo Outfit"},
		{"Otto Reiss's people", "Otto Reiss's people"},
		{"Amato Crew", "Amato Crew"},
		{"the Duarte Brothers", "The Duarte Brothers"},
		{"the Duarte Boys", "The Duarte Boys"},
	} {
		if got := Leads(c.in); got != c.want {
			t.Errorf("Leads(%q) = %q, want %q", c.in, got, c.want)
		}
	}
	// And a seeded family still reads correctly in a sentence that begins with
	// it, end to end.
	w := New(17)
	w.Player.Serves = "bellandi"
	w.Player.Service = 5
	if err := w.LeaveService(); err != nil {
		t.Fatalf("could not leave: %v", err)
	}
	last := w.History[len(w.History)-1]
	if !contains(last.Text, "Bellandi Family does not keep people who leave") {
		t.Fatalf("a seeded family was mangled: %q", last.Text)
	}
}
