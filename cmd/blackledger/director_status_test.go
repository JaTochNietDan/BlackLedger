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
	line := firstSentence(as.Error())
	if strings.Contains(line, ";") {
		t.Fatalf("the status line carries the whole instruction to the model: %q", line)
	}
	if !strings.HasPrefix(line, "Proposal repeats") || !strings.HasSuffix(line, ".") {
		t.Fatalf("the status line reads %q", line)
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
