package core

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Read out of a real ledger: "Stella Iordan was not as easy as he looked." The
// city hands out names of every kind — Stella, Franca, Perla, Elena, Mara, Ida
// sit beside Nico, Leo and Otto — and nothing anywhere records a gender,
// because the game has never had a reason to. The prose assumed one anyway, in
// forty-nine places, while the obituaries two paragraphs away said "They were a
// soldier of Russo Outfit."
//
// That is the game contradicting itself about the same person in the same
// issue, which is the only argument this test needs. The convention was already
// there; it just was not applied.
//
// This reads the source rather than the output because the faults are spread
// across forty files and most of them need a state no single campaign reaches.

// Nouns count as well as pronouns. "Nobody makes an arrangement like this with
// a man" and "3 men" for a crew of three were the same fault wearing a
// different word, and the first version of this test could not see either.
var gendered = regexp.MustCompile(`\b(he|him|his|himself|she|her|hers|herself|man|men|woman|women|boy|boys|girl|girls|gentleman|gentlemen|lady|ladies)\b`)

func TestNobodyInThisCityHasAGenderTheGameNeverGaveThem(t *testing.T) {
	names, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	offenders := []string{}
	for _, name := range names {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		source, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for i, line := range strings.Split(string(source), "\n") {
			trimmed := strings.TrimSpace(line)
			// Comments are prose about the code, not prose the player reads.
			if strings.HasPrefix(trimmed, "//") || !strings.Contains(line, `"`) {
				continue
			}
			// A few string literals are data the game matches against rather
			// than words it shows anybody — the endings that make an
			// organization's name plural, for one. They carry a marker, and
			// the marker has to be justified where it is written.
			if strings.Contains(line, "// not prose") {
				continue
			}
			// Only what is inside a string literal reaches a player.
			for _, part := range strings.Split(line, `"`)[1:] {
				if gendered.MatchString(strings.ToLower(part)) {
					offenders = append(offenders,
						name+":"+itoa(i+1)+" "+strings.TrimSpace(line))
					break
				}
			}
		}
	}
	if len(offenders) > 0 {
		t.Fatalf("%d line(s) of player-facing prose assume a gender nothing recorded:\n  %s",
			len(offenders), strings.Join(offenders, "\n  "))
	}
}
