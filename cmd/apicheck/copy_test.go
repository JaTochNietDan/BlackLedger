package main

import "testing"

// These patterns have been wrong twice. The count rule reported "11 people in
// here, which for this hour is a crowd" as a fault until it was given a word
// boundary, and the plural-subject rule reported "Violence between Brenner
// Company and Franca Sabbatini's people has escalated" — which is correct
// English, because the subject is the violence. A check that cries wolf gets
// ignored, and an ignored check is worse than none.

func TestThePluralSubjectRuleKnowsWhatTheSubjectIs(t *testing.T) {
	faults := []string{
		"Franca Sabbatini's people has people asking where you sleep.",
		"They were a soldier a week ago. Franca Sabbatini's people is short of people.",
		"Cesare Ferro's people holds the ground.",
	}
	fine := []string{
		"Violence between Brenner Company and Franca Sabbatini's people has escalated beyond the usual.",
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
