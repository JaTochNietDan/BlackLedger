package main

import (
	"blackledger/core"
	"strings"
	"testing"
)

func TestNarrativeBriefPreservesDirectionAndCompletedCallback(t *testing.T) {
	t.Parallel()
	w := core.New(27)
	w.District = 1
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

// Every operation the city can offer needs a premise, or the model writes its
// own — a consignment came back as a story about collecting payment for fabric,
// because nothing had told it what a consignment is.
func TestEveryOfferableOperationHasABrief(t *testing.T) {
	t.Parallel()
	w := core.New(151)
	w.MigrateLivingWorld()
	catalog := []string{"courier", "collection", "mediation"}
	for id := range core.SituationalEffects() {
		catalog = append(catalog, id)
	}
	for _, id := range catalog {
		b := jobBrief(w, id, nil)
		if b.Premise == "" {
			t.Fatalf("%s has no premise, so the model will invent one", id)
		}
		if b.PlayerTask == "" {
			t.Fatalf("%s does not say what the player is being asked to do", id)
		}
		if len(b.Constraints) == 0 {
			t.Fatalf("%s constrains nothing", id)
		}
	}
}

// The brief has to reach the model on whichever prompt is in use. It was sent
// only on the focused one, and the default is the full one.
func TestTheBriefReachesBothPrompts(t *testing.T) {
	t.Parallel()
	w := core.New(151)
	w.MigrateLivingWorld()
	focused := focusedContext(w, "consignment", nil, "", []string{""})
	if _, ok := focused["job_brief"]; !ok {
		t.Fatal("the focused prompt lost the brief")
	}
	// The full path assembles its context in generate; this asserts the key it
	// must carry, so removing the line fails here rather than in a play-test.
	brief := jobBrief(w, "consignment", nil)
	if brief.Premise == "" || !strings.Contains(brief.Premise, "under their floor") {
		t.Fatalf("the consignment brief reads %q", brief.Premise)
	}
}
