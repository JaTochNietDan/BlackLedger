package main

import (
	"errors"
	"strings"
	"testing"
)

// Driven against a live model: the director wrote a scene, the validator turned
// it down for repeating an earlier arrangement, the retry repeated it again,
// and the game set the status to "offline" with "Local AI unavailable or
// proposal rejected." Two different situations behind one word, and the wrong
// one of the two: the model was running and answering. Worse, "offline" is
// terminal — nothing asked again, and a single repeated title ended the AI for
// the rest of a hundred-and-thirty-day campaign.

func TestARefusedDraftIsNotAnOutage(t *testing.T) {
	reason := errors.New(`proposal repeats a recent offer "A warning from the Bellandi"; write a distinct task and title rather than swapping names`)
	var refused error = proposalRejected{reason}
	var as proposalRejected
	if !errors.As(refused, &as) {
		t.Fatal("a rejected proposal is not recognisable as one")
	}
	// What the player is shown must not be the validator's own words. Read on
	// the Settings screen in a browser: "The last draft was turned down: This
	// work exists only because of a committed event, and the dialogue never
	// refers to it." That is prose written to instruct a model.
	note := strings.ToLower(refusalNote(as.Error()))
	for _, leak := range []string{"proposal", "body", "dialogue", "beneficiary", "operation", "required_", "committed event", "swapping names"} {
		if strings.Contains(note, leak) {
			t.Fatalf("the player is shown the validator's own words: %q", note)
		}
	}
	if !strings.Contains(note, "ask for another") {
		t.Fatalf("the player is not told they can try again: %q", note)
	}
	// The exact reason still has somewhere to go: the log, where the last five
	// faults were found.
	if line := firstSentence(as.Error()); !strings.HasPrefix(line, "Proposal repeats") {
		t.Fatalf("the reason kept for the log reads %q", line)
	}
}

func TestAStatusLineSurvivesAnEmptyReason(t *testing.T) {
	if got := firstSentence("   "); got != "no reason was given." {
		t.Fatalf("an empty reason reads as %q", got)
	}
	if got := firstSentence("something went wrong"); got != "Something went wrong." {
		t.Fatalf("a bare reason reads as %q", got)
	}
}
