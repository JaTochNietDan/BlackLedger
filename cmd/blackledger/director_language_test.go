package main

import (
	"blackledger/core"
	"testing"
)

func TestMixedScriptChoiceIsRejectedWithoutDamagingAccents(t *testing.T) {
	for _, label := range []string{"Use the back通道", "Пройти quietly"} {
		if validateChoiceScript(core.Proposal{Approaches: []core.Approach{{Label: label}}}) == nil {
			t.Fatal("mixed script accepted", label)
		}
	}
	for _, label := range []string{"Meet at the café", "Use Mara’s entrance", "Take the back passage"} {
		if err := validateChoiceScript(core.Proposal{Approaches: []core.Approach{{Label: label}}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSceneTitleNamesTheSituationRatherThanRepeatingAChoice(t *testing.T) {
	// Observed on qwen3.5:35b-a3b: a new-life mediation was titled with its own
	// approach label, "Quietly listen to both sides".
	repeated := core.Proposal{Title: "Quietly listen to both sides", Approaches: []core.Approach{{Method: "careful", Label: "Quietly listen to both sides"}}}
	if validateSceneTitle(repeated) == nil {
		t.Fatal("title repeating an approach label accepted")
	}
	// Punctuation and casing must not let the same label through as a title.
	if validateSceneTitle(core.Proposal{Title: "Quietly Listen to Both Sides!", Approaches: []core.Approach{{Method: "careful", Label: "quietly listen to both sides"}}}) == nil {
		t.Fatal("restyled duplicate title accepted")
	}
	if validateSceneTitle(core.Proposal{Approaches: []core.Approach{{Method: "careful", Label: "Quietly listen"}}}) == nil {
		t.Fatal("missing title accepted")
	}
	for _, p := range []core.Proposal{
		{Title: "A Dispute at The Mariner", Approaches: []core.Approach{{Method: "careful", Label: "Quietly listen to both staff"}}},
		{Title: "Payment at Russo Motor Works", Approaches: []core.Approach{{Method: "careful", Label: "Secure the cash quietly"}}},
		{Title: "Negotiating Access at The Mariner", Approaches: []core.Approach{{Method: "careful", Label: "Quietly settle their dispute"}}},
	} {
		if err := validateSceneTitle(p); err != nil {
			t.Fatal("ordinary scene title rejected:", p.Title, err)
		}
	}
}

func TestEveryApproachNeedsALabelThePlayerCanPress(t *testing.T) {
	// Observed live on qwen3.5:35b-a3b: a war-derived offer came back with an
	// empty first approach label, which renders as an unpressable button.
	if validateApproachLabels(core.Proposal{Approaches: []core.Approach{{Method: "careful", Label: ""}}}) == nil {
		t.Fatal("an empty approach label was accepted")
	}
	if validateApproachLabels(core.Proposal{Approaches: []core.Approach{{Method: "press", Label: "  "}}}) == nil {
		t.Fatal("a whitespace approach label was accepted")
	}
	if validateApproachLabels(core.Proposal{Approaches: []core.Approach{{Method: "press", Label: "Go"}}}) == nil {
		t.Fatal("a label too short to read was accepted")
	}
	if err := validateApproachLabels(core.Proposal{Approaches: []core.Approach{
		{Method: "careful", Label: "Speak with him quietly"},
		{Method: "press", Label: "Walk in and make it clear"},
	}}); err != nil {
		t.Fatal("ordinary approach labels rejected:", err)
	}
}

func TestAnApproachLabelDoesNotStopMidPhrase(t *testing.T) {
	// Observed live: the model wrote up to the 45-character limit and the label
	// was cut at "Drive directly to the exchange and hand it to".
	for _, label := range []string{
		"Drive directly to the exchange and hand it to",
		"Wait for the buyer and",
		"Take it round the",
		"Leave it with a",
	} {
		t.Run(label, func(t *testing.T) {
			if validateApproachLabels(core.Proposal{Approaches: []core.Approach{{Method: "careful", Label: label}}}) == nil {
				t.Fatal("a label stopping mid-phrase was accepted")
			}
		})
	}
	for _, label := range []string{
		"Use the back entrance",
		"Hand it over quietly",
		"Drive it across town",
		"Wait until after dark.",
	} {
		t.Run("clean/"+label, func(t *testing.T) {
			if err := validateApproachLabels(core.Proposal{Approaches: []core.Approach{{Method: "careful", Label: label}}}); err != nil {
				t.Fatal("a finished label was rejected:", err)
			}
		})
	}
}
