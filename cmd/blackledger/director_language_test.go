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
