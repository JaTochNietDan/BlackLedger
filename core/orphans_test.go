package core

import (
	"strings"
	"testing"
)

// When an organization ends, the city writes an obituary for it and forgets it.
// The people who answered to it are not touched: their record still names a
// family that no longer exists in the world. Everything that looks one up finds
// nothing, and the men themselves go on standing in rooms belonging to nobody.
//
// That is the whole of the idea worth building here — a family falling should
// put people on the street — but the first thing to establish is whether it
// currently leaves a dangling name behind at all.

func endTheFamily(t *testing.T) (*World, []string) {
	t.Helper()
	w := New(41)
	w.District = 2
	// Somebody has to be left standing or dissolve refuses: it will not empty
	// a city down to fewer than two organizations.
	w.Factions = append(w.Factions, Faction{ID: "spare", Name: "Alvarez Company",
		Leader: "Ilse Alvarez", Power: 40, Peak: 40, Cash: 5000})
	f := w.faction("bellandi")
	if f == nil {
		t.Fatal("no bellandi")
	}
	people := []string{}
	for _, n := range w.Members("bellandi") {
		people = append(people, n.ID)
	}
	if len(people) < 2 {
		t.Fatalf("bellandi has %d people, so there is nobody to strand", len(people))
	}
	for _, id := range w.FamilyHoldings("bellandi") {
		w.Properties[id].Owner = "independent"
	}
	f.Power, f.Cash = 10, 0
	w.dissolve()
	if w.faction("bellandi") != nil {
		t.Fatal("the family did not end")
	}
	return w, people
}

func TestWhenAFamilyEndsNobodyStillAnswersToIt(t *testing.T) {
	t.Parallel()
	w, people := endTheFamily(t)
	stranded := []string{}
	for _, id := range people {
		n := w.NPC(id)
		if n == nil || n.Dead {
			continue
		}
		if n.Faction == "bellandi" {
			stranded = append(stranded, n.Name)
		}
	}
	if len(stranded) > 0 {
		t.Errorf("%d people still answer to an organization the city has buried: %v", len(stranded), stranded)
	}
}

// And nothing anywhere should be able to find a member of a family that is not
// in the world.
func TestNobodyIsAMemberOfAnOrganizationThatIsGone(t *testing.T) {
	t.Parallel()
	w, _ := endTheFamily(t)
	if left := w.Members("bellandi"); len(left) > 0 {
		t.Errorf("%d people are still on the books of a family nobody can look up", len(left))
	}
	for _, n := range w.People() {
		if n.Faction == "" {
			continue
		}
		if w.faction(n.Faction) == nil {
			t.Errorf("%s answers to %q, which is not an organization in this city", n.Name, n.Faction)
		}
	}
}

// The city should say what happened to them, and say it correctly for one man
// as well as for nine.
func TestTheCityCountsWhoIsPutOnTheStreet(t *testing.T) {
	t.Parallel()
	w, people := endTheFamily(t)
	said := ""
	for _, r := range w.History {
		if r.Title == "An organization ends" {
			said = r.Text
		}
	}
	if said == "" {
		t.Fatal("a family ended and the city wrote nothing about it")
	}
	if !strings.Contains(said, counted(len(people), "person is", "people are")) {
		t.Errorf("%d people were put out and the city said: %q", len(people), said)
	}
	// The singular has to read as English too, and nobody in this city has a
	// gender the game ever recorded.
	if got := counted(1, "person is", "people are"); got != "1 person is" {
		t.Errorf("one person reads as %q", got)
	}
}
