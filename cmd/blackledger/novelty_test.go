package main

import (
	"blackledger/core"
	"testing"
)

func TestRecentOffersCannotBePitchedAgainVerbatim(t *testing.T) {
	w := core.New(27)
	w.Arrangements = []core.ArrangementMemory{{Life: 1, Title: "A dispute at the garage", Offer: "Share the tools fairly.", Status: "declined"}}
	for _, p := range []core.Proposal{{Title: "  A DISPUTE at the garage "}, {Title: "Different title", Body: "“Share the tools fairly.”"}} {
		if repeatedProposal(w, p) == nil {
			t.Fatal("repeated offer accepted")
		}
	}
	if repeatedProposal(w, core.Proposal{Title: "A new schedule", Body: "Arrange access to the loading bay."}) != nil {
		t.Fatal("distinct task rejected")
	}
	w.Life = 2
	if repeatedProposal(w, core.Proposal{Title: "A dispute at the garage"}) != nil {
		t.Fatal("old life incorrectly blocked current offer")
	}
}
