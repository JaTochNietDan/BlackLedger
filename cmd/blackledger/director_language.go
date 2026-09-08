package main

import (
	"blackledger/core"
	"fmt"
	"strings"
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

// An approach is a button the player presses. A blank one is unpressable and
// reads as a missing option; observed live on qwen3.5:35b-a3b, which returned an
// empty first label.
func validateApproachLabels(p core.Proposal) error {
	for _, a := range p.Approaches {
		if len(strings.TrimSpace(a.Label)) < 4 {
			return fmt.Errorf("approach %q has no usable label; give every approach a short phrase naming how the player would do it", a.Method)
		}
	}
	return nil
}

// A title names the situation the player is looking at; an approach label names
// one way to act on it. Reusing the label as the title leaves the scene with no
// heading of its own and makes the offer list read as a row of duplicate verbs.
func validateSceneTitle(p core.Proposal) error {
	title := normalizedOffer(p.Title)
	if title == "" {
		return fmt.Errorf("give the scene a short title naming the situation and its place")
	}
	for _, a := range p.Approaches {
		if title == normalizedOffer(a.Label) {
			return fmt.Errorf("title %q repeats the approach label %q; title the situation and where it is happening, and keep the action wording in the approach", p.Title, a.Label)
		}
	}
	return nil
}
