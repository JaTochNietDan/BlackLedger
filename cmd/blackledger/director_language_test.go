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
