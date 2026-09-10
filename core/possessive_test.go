package core

import (
	"strings"
	"testing"
)

// Reading the Ledger screen in a browser: "Sofia Doyle helped themselves to
// $133 of Franca Sabbatini's people money at Bluebird Laundry."
//
// Splinter families are named after the person who broke away — "Franca
// Sabbatini's people" — so a sentence that puts a possessive after an
// organization's name stacks one possessive on another. The seeded families
// read fine, which is why this survived: "$133 of Bellandi Family money" is
// ordinary attributive English. It only breaks on the names the living world
// makes for itself, and making families is most of what it does.

func TestATheftReadsForEveryKindOfFamilyName(t *testing.T) {
	t.Parallel()
	for _, name := range []string{
		"Bellandi Family",
		"Russo Outfit",
		"Franca Sabbatini's people",
		"the Duarte Brothers",
	} {
		got := theftLine("Sofia Doyle", 133, name, "Bluebird Laundry")
		if strings.Contains(got, "people money") || strings.Contains(got, "Brothers money") {
			t.Errorf("%q stacks a possessive: %q", name, got)
		}
		if !strings.Contains(got, name) {
			t.Errorf("%q is not named at all: %q", name, got)
		}
		// The verb has to follow the name too: people take, a Family keeps.
		if strings.HasSuffix(name, "people") || strings.HasSuffix(name, "Brothers") {
			if !strings.Contains(got, name+" keep ") {
				t.Errorf("a plural name takes a plural verb: %q", got)
			}
		} else if !strings.Contains(got, name+" keeps ") {
			t.Errorf("a singular name takes a singular verb: %q", got)
		}
		t.Logf("%s", got)
	}
}
