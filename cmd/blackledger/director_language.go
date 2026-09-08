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
// danglingEnds are words a finished phrase does not end on. A label written up
// to the length limit gets cut mid-sentence, which reads as a broken button:
// observed live as "Drive directly to the exchange and hand it to".
var danglingEnds = map[string]bool{
	"to": true, "and": true, "the": true, "a": true, "an": true, "at": true,
	"with": true, "for": true, "of": true, "in": true, "on": true, "from": true,
	"into": true, "by": true, "or": true, "then": true,
}

func validateApproachLabels(p core.Proposal) error {
	for _, a := range p.Approaches {
		label := strings.TrimSpace(a.Label)
		if len(label) < 4 {
			return fmt.Errorf("approach %q has no usable label; give every approach a short phrase naming how the player would do it", a.Method)
		}
		words := strings.Fields(strings.Trim(label, ".,;:!?"))
		if last := strings.ToLower(words[len(words)-1]); danglingEnds[last] {
			return fmt.Errorf("approach label %q stops mid-phrase on %q; write a short complete phrase of a few words that fits well inside the length limit", label, last)
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
