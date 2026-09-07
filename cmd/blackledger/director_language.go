package main

import (
	"blackledger/core"
	"fmt"
	"unicode"
)

// The current interface is English. Catch the observed mixed-script choice
// labels without stripping accented names or punctuation from generated text.
// This is a script check, not a general language classifier.
func validateChoiceScript(p core.Proposal) error {
	for _, a := range p.Approaches {
		for _, r := range a.Label {
			if unicode.IsLetter(r) && !unicode.In(r, unicode.Latin) {
				return fmt.Errorf("write approach labels in English using Latin-script text; replace the mixed-script label %q", a.Label)
			}
		}
	}
	return nil
}
