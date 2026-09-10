package core

import (
	"regexp"
	"testing"
)

// Reading every action's description side by side turned up the city counting
// out loud and getting it wrong: "$164 due in 1 days" on three buttons at the
// casino, and "1 stories in today's paper are about you" at the Herald. The
// game already has grammar helpers keyed to code rather than to save state —
// Agree, spelled, withArticle — and nothing that agrees a noun with a number.

func TestCountedAgreesTheNounWithTheNumber(t *testing.T) {
	for _, c := range []struct {
		n    int
		want string
	}{
		{0, "0 days"}, {1, "1 day"}, {2, "2 days"}, {21, "21 days"},
	} {
		if got := counted(c.n, "day", "days"); got != c.want {
			t.Errorf("counted(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

// And the property, over everything the player can actually read in a room: a
// description that counts must not put a plural noun after one.
var singularOne = regexp.MustCompile(`\b1 [a-z]+s\b`)

func TestNoDescriptionCountsOneOfSomethings(t *testing.T) {
	w := New(6)
	w.District = 2
	w.Player.Cash = 5000
	w.Player.Respect = 60
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 65}}
	// One loan, due tomorrow, so the collection buttons have to say "1 day".
	w.NPCs = append(w.NPCs, NPC{ID: "borrower", Name: "Otto Reiss", Location: "casino"})
	w.Loans = []Loan{{ID: ID(), Debtor: "borrower", Life: w.Life, Principal: 140, Owed: 164, Due: w.Minute + 1440}}
	checked := 0
	for i := range Locations {
		id := Locations[i].ID
		w.Player.Location = id
		for _, a := range w.Actions(id) {
			checked++
			for _, text := range []string{a.Label, a.Detail, a.Reason} {
				if m := singularOne.FindString(text); m != "" {
					t.Errorf("%s at %s counts %q:\n  %s", a.ID, id, m, text)
				}
			}
		}
	}
	t.Logf("read %d actions", checked)
}

// The Herald counts stories, and the verb has to follow the number too.
func TestOneStoryIsAboutYouAndTwoStoriesAre(t *testing.T) {
	w := New(6)
	w.District = 2
	w.Player.Location = "herald"
	read := func() string {
		for _, a := range w.Actions("herald") {
			if a.ID == "spike" {
				return a.Detail
			}
		}
		t.Fatal("the Herald does not offer to pull a story")
		return ""
	}
	if got := read(); !contains(got, "Nothing in today's paper is about you") {
		t.Fatalf("a paper with nothing in it does not say so: %q", got)
	}
	// One story about the player in today's paper, then two. A police story is
	// spikeable by kind, so this does not depend on what the player owns.
	w.News = append(w.News, Story{ID: ID(), Minute: w.Minute, Life: w.Life,
		Headline: "Arrested at the docks", Body: "The police took somebody in.", Kind: "police"})
	if n := len(w.Spikeable()); n != 1 {
		t.Fatalf("expected one spikeable story, got %d", n)
	}
	if one := read(); !contains(one, "1 story in today's paper is about you") {
		t.Fatalf("one story does not read as one: %q", one)
	}
	w.News = append(w.News, Story{ID: ID(), Minute: w.Minute, Life: w.Life,
		Headline: "Questioned about a killing", Body: "The police called again.", Kind: "police"})
	if n := len(w.Spikeable()); n != 2 {
		t.Fatalf("expected two spikeable stories, got %d", n)
	}
	if two := read(); !contains(two, "2 stories in today's paper are about you") {
		t.Fatalf("two stories do not read as two: %q", two)
	}
}
