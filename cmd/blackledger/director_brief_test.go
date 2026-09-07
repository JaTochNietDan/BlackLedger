package main

import (
	"blackledger/core"
	"strings"
	"testing"
)

func TestNarrativeBriefPreservesDirectionAndCompletedCallback(t *testing.T) {
	w := core.New(27)
	m := &core.ArrangementMemory{Operation: "mediation", Offer: "A tool disagreement at Russo Motor Works", Status: "completed"}
	b := jobBrief(w, "collection", m)
	if b.Location != "Russo Motor Works" || !strings.Contains(b.SourceRole, "customer") || !strings.Contains(b.RecipientRole, "manager") {
		t.Fatal("collection roles reversed or location lost")
	}
	if validateBriefOpening("The laundry has a payment.", b) == nil {
		t.Fatal("unacknowledged callback accepted")
	}
	if err := validateBriefOpening("“"+b.Opening+" Here is another request.", b); err != nil {
		t.Fatal(err)
	}
	if jobBrief(w, "courier", nil).Opening != "" {
		t.Fatal("fresh job claims prior work")
	}
	if !strings.Contains(jobBrief(w, "mediation", nil).PlayerTask, "without violence") {
		t.Fatal("mediation brief changed mechanics")
	}
}
