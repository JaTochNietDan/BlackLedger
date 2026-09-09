package core

import (
	"regexp"
	"strings"
	"testing"
)

// Shape found by reading the Ledger: prose that reads correctly for the two
// families every campaign starts with, and breaks on the names the living world
// makes for itself. An organization inherited by somebody is called "<Person>'s
// people" and takes a plural verb; the seeded "Bellandi Family" and "Russo
// Outfit" take a singular one, so a sentence written against them alone is
// wrong for every family the world creates by that path.
//
// This drives a world with a plural-named family in it and reads everything the
// city writes down, rather than trusting a grep.

var disagreements = regexp.MustCompile(`(?i)\bpeople (has|is|does|knows|wants|believes|holds|keeps|comes|goes|thinks|makes|takes|was|owns)\b`)

// A plural name inside a prepositional phrase is not the subject of the verb
// that follows it. "Violence between Bellandi Family and Otto Reiss's people
// has escalated" is correct — the violence has escalated — and the first
// version of this scan reported it. Only the clause the name governs counts.
func disagrees(text string) string {
	for _, sentence := range strings.Split(text, ". ") {
		m := disagreements.FindStringIndex(sentence)
		if m == nil {
			continue
		}
		if strings.Contains(strings.ToLower(sentence[:m[0]]), "between ") {
			continue
		}
		return sentence[m[0]:m[1]]
	}
	return ""
}

func TestNothingTheCityWritesDisagreesWithAPluralFamily(t *testing.T) {
	w := New(21)
	// A family that answers to a person's name, exactly as succession makes one.
	w.Factions = append(w.Factions, Faction{
		ID: "estate:otto", Name: "Otto Reiss's people", Leader: "Otto Reiss",
		Power: 55, Cash: 3000, Peak: 55, Reported: 55,
	})
	w.NPCs = append(w.NPCs, NPC{ID: "otto", Name: "Otto Reiss", Role: "Head of Otto Reiss's people",
		Faction: "estate:otto", Location: "market", Rank: RankLeader, Ambition: 70, Skill: 70})
	w.Properties["laundry"].Owner = "estate:otto"
	w.Properties["casino"].Owner = "estate:otto"
	w.District = 2
	w.Player.Cash = 6000
	w.Player.Respect = 80

	read := func(where string) {
		for _, r := range w.History {
			for _, text := range []string{r.Title, r.Text} {
				if m := disagrees(text); m != "" {
					t.Errorf("%s: the city wrote %q:\n  %s", where, m, text)
				}
			}
		}
		for _, s := range w.News {
			for _, text := range []string{s.Headline, s.Body} {
				if m := disagrees(text); m != "" {
					t.Errorf("%s: the paper wrote %q:\n  %s", where, m, text)
				}
			}
		}
		for i := range Locations {
			id := Locations[i].ID
			w.Player.Location = id
			for _, a := range w.Actions(id) {
				for _, text := range []string{a.Label, a.Detail, a.Reason} {
					if m := disagrees(text); m != "" {
						t.Errorf("%s: %s at %s reads %q:\n  %s", where, a.ID, id, m, text)
					}
				}
			}
		}
		for _, c := range w.Commissions {
			if m := disagrees(c.Brief); m != "" {
				t.Errorf("%s: a commission reads %q:\n  %s", where, m, c.Brief)
			}
		}
	}

	read("at the start")
	// Let the city live: wars, successions, takeovers, commissions.
	w.Antagonize("bellandi", "estate:otto", 90)
	w.Antagonize("russo", "estate:otto", 90)
	for i := 0; i < 4000; i++ {
		w.Advance(30)
		if !w.Player.Alive {
			w.Player.Alive = true
			w.Player.Health = 100
		}
	}
	read("after 4000 half-hours")
	t.Logf("%d ledger entries, %d stories, %d families", len(w.History), len(w.News), len(w.Factions))
}

// The other half: a singular name must keep its singular verb. Fixing agreement
// by making everything plural would read just as wrong for the two families
// every campaign starts with.
func TestASingularFamilyKeepsItsSingularVerb(t *testing.T) {
	for _, c := range []struct{ name, want string }{
		{"Bellandi Family", "does not"},
		{"Russo Outfit", "does not"},
		{"Otto Reiss's people", "do not"},
		{"the Duarte Brothers", "do not"},
	} {
		if got := Agree(c.name, "does not", "do not"); got != c.want {
			t.Errorf("%q takes %q, want %q", c.name, got, c.want)
		}
	}
	w := New(21)
	w.Player.Serves = "bellandi"
	w.Player.Service = 5
	if err := w.LeaveService(); err != nil {
		t.Fatalf("could not leave: %v", err)
	}
	last := w.History[len(w.History)-1]
	if !contains(last.Text, "Bellandi Family does not keep people who leave") {
		t.Fatalf("a singular family lost its singular verb: %q", last.Text)
	}
}
