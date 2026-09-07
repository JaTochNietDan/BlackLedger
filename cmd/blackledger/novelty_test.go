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

func TestRenamedPaymentStoryStillCountsAsRepetition(t *testing.T) {
	old := "The Bluebird Laundry has a debt to settle. A supplier left a sealed payment under the counter last night, but the owner is too nervous to collect it himself. I need someone who can take it without drawing attention. You've handled quiet jobs before, and I trust you to keep it that way."
	renamed := "The Blue Hour needs a payment collected discreetly. A bookmaker left a sealed envelope under the counter last night, but the owner is too nervous to collect it himself. I need someone who can take it without drawing attention. You’ve handled quiet jobs before, and I trust you to keep it that way."
	if !recycledPassage(old, renamed) {
		t.Fatal("location substitution bypassed repeated-passage check")
	}
	distinct := "Two delivery firms arrive at the garage together every morning. Their drivers block the entrance while arguing about who unloads first. Ask them to agree on separate slots so the mechanics can keep working. I need a practical agreement they can both follow, without a scene that brings the police to our door."
	if recycledPassage(old, distinct) {
		t.Fatal("ordinary mafia vocabulary rejected as repetition")
	}
	if recycledPassage("Carry the package for me.", "Carry the package for me.") {
		t.Fatal("short passage should use exact check only")
	}
}

func TestPendingAndCurrentOffersAlsoPreventRepetition(t *testing.T) {
	w := core.New(27)
	p := core.Proposal{Title: "A pending favor", Body: "Move the records discreetly."}
	scene := &core.Scene{Title: p.Title, Body: p.Body, Kind: "proposal"}
	w.Offers = []core.Offer{{Event: scene}}
	if repeatedProposal(w, p) == nil {
		t.Fatal("queued duplicate accepted")
	}
	w.Offers = nil
	w.Event = scene
	if repeatedProposal(w, p) == nil {
		t.Fatal("current duplicate accepted")
	}
}
