package main

import (
	"blackledger/core"
	"testing"
)

func cityAtWar() *core.World {
	w := core.New(211)
	w.Antagonize("bellandi", "russo", 100)
	return w
}

func TestConflictWorkMustReferToTheConflict(t *testing.T) {
	t.Parallel()
	w := cityAtWar()
	if w.Conflict("bellandi", "russo").State != "war" {
		t.Fatal("the test world is not at war")
	}
	// The observed failure: a war-derived warning written as an errand.
	ungrounded := core.Proposal{Operation: "warning", Body: "The Mariner's bar is where this has to be done. Speak with the man at the bar and hand him a message. He'll know what it is."}
	if validateSituationalGrounding(w, ungrounded) == nil {
		t.Fatal("a war job that never mentions the war was accepted")
	}
	for _, body := range []string{
		"The Bellandi Family will not let this pass. Carry the message and leave.",
		"Say it to Russo's people directly, and do not stay to argue.",
		"Vittorio Bellandi needs to hear it from somebody standing in front of him.",
	} {
		t.Run(body, func(t *testing.T) {
			if err := validateSituationalGrounding(w, core.Proposal{Operation: "warning", Body: body}); err != nil {
				t.Fatal("a grounded war job was rejected:", err)
			}
		})
	}
}

func TestOrdinaryWorkIsNotAskedToMentionAnything(t *testing.T) {
	t.Parallel()
	w := cityAtWar()
	for _, operation := range []string{"courier", "mediation", "collection"} {
		if err := validateSituationalGrounding(w, core.Proposal{Operation: operation, Body: "Two staff cannot agree about the loading bay."}); err != nil {
			t.Fatalf("%s was required to mention a conflict: %v", operation, err)
		}
	}
	// And in a quiet city there is nothing to ground against.
	quiet := core.New(213)
	if err := validateSituationalGrounding(quiet, core.Proposal{Operation: "warning", Body: "Carry it over and come back."}); err != nil {
		t.Fatal("a quiet city demanded grounding:", err)
	}
}
