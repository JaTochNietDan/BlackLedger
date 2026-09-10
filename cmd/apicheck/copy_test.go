package main

import (
	"strings"
	"testing"
)

// These patterns have been wrong twice. The count rule reported "11 people in
// here, which for this hour is a crowd" as a fault until it was given a word
// boundary, and the plural-subject rule reported "Violence between Brenner
// Company and Franca Sabbatini's people has escalated" — which is correct
// English, because the subject is the violence. A check that cries wolf gets
// ignored, and an ignored check is worse than none.

func TestThePluralSubjectRuleKnowsWhatTheSubjectIs(t *testing.T) {
	t.Parallel()
	faults := []string{
		"Franca Sabbatini's people has people asking where you sleep.",
		"They were a soldier a week ago. Franca Sabbatini's people is short of people.",
		"Cesare Ferro's people holds the ground.",
		"Nico Ward's people controls Saint Agnes now.",
	}
	fine := []string{
		"Violence between Brenner Company and Franca Sabbatini's people has escalated beyond the usual.",
		"Nico Ward's people moved against Saint Agnes and were driven off.",
		"Nico Ward's people came for Bluebird Laundry.",
		"Franca Sabbatini's people have taken Saint Agnes from Falcone Crew.",
		"Bellandi Family has taken Saint Agnes.",
		"Officers raided premises connected to Cesare Ferro's people.",
	}
	for _, s := range faults {
		if !pluralSubject.MatchString(s) {
			t.Fatalf("a real fault went unreported: %q", s)
		}
	}
	for _, s := range fine {
		if pluralSubject.MatchString(s) {
			t.Fatalf("correct English was reported as a fault: %q", s)
		}
	}
}

func TestTheCountRuleCountsToOne(t *testing.T) {
	t.Parallel()
	faults := []string{
		"Whatever was arranged for you happened 1 times to a locked door.",
		"$320 for the 1 days still on him.",
	}
	fine := []string{
		"11 people in here, which for this hour is a crowd.",
		"The Monarch is full tonight. 21 people, and the ones by the door are watching.",
		"Whatever was arranged for you happened once to a locked door.",
		"And 1 other story.",
	}
	for _, s := range faults {
		if !singularOne.MatchString(s) {
			t.Fatalf("a real fault went unreported: %q", s)
		}
	}
	for _, s := range fine {
		if singularOne.MatchString(s) {
			t.Fatalf("correct English was reported as a fault: %q", s)
		}
	}
}

// A 409 on an action the game had just listed as available is the only way the
// timing-drift class has ever been caught, and the report used to carry the
// message and nothing else. Finding the cause of the one instance took three
// failed hypotheses and seven hundred attempts, because the state that produced
// it was gone by the time anyone read the report.
func TestARefusalBringsItsOwnEvidence(t *testing.T) {
	t.Parallel()
	s := &snapshot{Revision: 9, Minute: 41520}
	s.Player.Location = "bar"
	s.Player.Cash, s.Player.Health, s.Player.Heat, s.Player.Respect = 512, 74, 31, 60
	s.Locations = []place{{
		ID: "bar", Name: "Saint Agnes",
		Actions: []action{
			{ID: "sitdown", Disabled: false},
			{ID: "rob", Disabled: true, Reason: "There is nothing here worth taking"},
		},
	}}
	s.Factions = []struct {
		ID       string `json:"id"`
		Goodwill int    `json:"goodwill"`
	}{{ID: "bellandi", Goodwill: 12}, {ID: "doyle", Goodwill: -26}}

	r := asOffered(s, command{Kind: "sitdown", Target: "bar"})
	if !r.Offered {
		t.Fatal("the action was listed as available and the evidence says otherwise")
	}
	if r.Location != "bar" || r.Cash != 512 || r.Health != 74 || r.Heat != 31 || r.Minute != 41520 {
		t.Fatalf("the world around the refusal was not recorded: %+v", r)
	}
	// The standing that moved under the sitdown is the thing the report never
	// showed, so it has to be there.
	if !strings.Contains(r.Standing, "doyle -26") {
		t.Fatalf("standing was not recorded: %q", r.Standing)
	}
	// And a genuinely disabled action is recorded as disabled, with its reason.
	d := asOffered(s, command{Kind: "rob", Target: "bar"})
	if d.Offered {
		t.Fatal("a disabled action was recorded as offered")
	}
	if !strings.Contains(d.Reason, "nothing here worth taking") {
		t.Fatalf("the reason shown was not recorded: %q", d.Reason)
	}
}
