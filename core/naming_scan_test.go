package core

import (
	"regexp"
	"strings"
	"testing"
)

// The agreement sweep proved the method: build a city with an awkward name in
// it, run it forward, then read everything it wrote. This generalises that to
// the other name shapes the world can produce, and to two hazards rather than
// one.
//
// A name that carries its own article — "the Duarte Brothers" — begins a
// sentence in lower case unless something capitalises it. core/names.go has
// Leads for exactly that. And a name ending in s — "Otto Reiss" — takes a
// possessive the code has to spell one way or the other.

var lowerStart = regexp.MustCompile(`(?:^|\. )(the [A-Z])`)

// sentences the city wrote, split for reading one clause at a time.
func writings(w *World) []string {
	out := []string{}
	for _, r := range w.History {
		out = append(out, r.Title, r.Text)
	}
	for _, s := range w.News {
		out = append(out, s.Headline, s.Body)
	}
	for _, c := range w.Commissions {
		out = append(out, c.Brief)
	}
	for i := range Locations {
		id := Locations[i].ID
		w.Player.Location = id
		for _, a := range w.Actions(id) {
			out = append(out, a.Label, a.Detail, a.Reason)
		}
	}
	return out
}

// live runs a city forward so the prose is what the world actually produced,
// not what a fixture posed for.
func live(t *testing.T, w *World, halfHours int) {
	t.Helper()
	for i := 0; i < halfHours; i++ {
		w.Advance(30)
		if !w.Player.Alive {
			w.Player.Alive, w.Player.Health = true, 100
		}
	}
}

func withFamily(seed uint32, id, name, leader string) *World {
	w := New(seed)
	w.Factions = append(w.Factions, Faction{
		ID: id, Name: name, Leader: leader,
		Power: 55, Cash: 3000, Peak: 55, Reported: 55,
	})
	w.NPCs = append(w.NPCs, NPC{ID: "boss-" + id, Name: leader, Role: "Head of " + name,
		Faction: id, Location: "market", Rank: RankLeader, Ambition: 70, Skill: 70})
	w.Properties["laundry"].Owner = id
	w.Properties["casino"].Owner = id
	w.District = 2
	w.Player.Cash = 6000
	w.Player.Respect = 80
	w.Antagonize("bellandi", id, 90)
	w.Antagonize("russo", id, 90)
	return w
}

// A name that carries "the " must not begin a sentence in lower case.
func TestNoSentenceBeginsWithALowerCaseArticle(t *testing.T) {
	t.Parallel()
	// Several cities, because one run's luck is not evidence: a sentence only
	// gets written when the world happens to do the thing that writes it.
	total := 0
	for _, seed := range []uint32{31, 47, 88, 103, 219} {
		w := withFamily(seed, "splinter-Duarte", "the Duarte Brothers", "Bruno Duarte")
		live(t, w, 4000)
		found := 0
		for _, text := range writings(w) {
			if !strings.Contains(text, "Duarte") {
				continue
			}
			found++
			for _, sentence := range strings.Split(text, ". ") {
				if strings.HasPrefix(sentence, "the ") {
					t.Errorf("seed %d: a sentence begins in lower case: %q\n  in: %s", seed, sentence, text)
				}
			}
		}
		total += found
	}
	if total == 0 {
		t.Fatal("no city ever mentioned the family, so this proved nothing")
	}
	t.Logf("read %d passages naming the family across five cities", total)
}

// A name ending in s takes one possessive form, consistently.
func TestAPossessiveOnANameEndingInSIsSpelledOneWay(t *testing.T) {
	t.Parallel()
	w := withFamily(33, "estate:reiss", "Otto Reiss's people", "Otto Reiss")
	live(t, w, 4000)
	apostropheOnly, withS, found := 0, 0, 0
	for _, text := range writings(w) {
		if !strings.Contains(text, "Reiss") {
			continue
		}
		found++
		withS += strings.Count(text, "Reiss's")
		// "Reiss' " with no trailing s, which is the other convention.
		apostropheOnly += strings.Count(text, "Reiss' ")
	}
	if found == 0 {
		t.Fatal("the city never mentioned them, so this proved nothing")
	}
	if apostropheOnly > 0 && withS > 0 {
		t.Errorf("the city spells the same possessive both ways: %d as Reiss's, %d as Reiss'", withS, apostropheOnly)
	}
	t.Logf("read %d passages, %d possessives", found, withS+apostropheOnly)
}

// And the player's own organization, which appears in prose a rival's does not.
func TestThePlayersOwnOrganizationReadsInEverySentence(t *testing.T) {
	t.Parallel()
	w := New(37)
	w.District = 2
	w.Player.Cash = 8000
	w.Player.Respect = 120
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	w.Incorporate()
	f := w.PlayerOrganization()
	if f == nil {
		t.Fatalf("the player did not become an organization: ready=%v strength=%d",
			w.OrganizationReady(), w.PlayerStrength())
	}
	live(t, w, 3000)
	found := 0
	for _, text := range writings(w) {
		if !strings.Contains(text, f.Name) {
			continue
		}
		found++
		for _, sentence := range strings.Split(text, ". ") {
			if strings.HasPrefix(sentence, "the ") {
				t.Errorf("a sentence begins in lower case: %q", sentence)
			}
			i := strings.Index(sentence, f.Name)
			if i < 0 || strings.Contains(strings.ToLower(sentence[:max(0, i)]), "between ") {
				continue
			}
			if m := disagreements.FindString(sentence[i:]); m != "" {
				t.Errorf("the player's own organization reads %q:\n  %s", m, sentence)
			}
		}
	}
	if found == 0 {
		t.Fatal("the city never named the player's organization, so this proved nothing")
	}
	t.Logf("read %d passages naming %q", found, f.Name)
}
