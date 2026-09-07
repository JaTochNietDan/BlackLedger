package main

import (
	"blackledger/core"
	"testing"
)

func TestBriefsPreserveCanonicalResultsWithoutMutatingStoryMemory(t *testing.T) {
	w := core.New(27)
	w.Arrangements = []core.ArrangementMemory{{Life: 1, Title: "Old task", Offer: "A detailed old offer.", Status: "completed", Result: "Delivered."}}
	b := arrangementBriefs(w)
	if b[0].Offer != "" || b[0].Result != "Delivered." || b[0].Status != "completed" {
		t.Fatal("brief lost outcome or copied old prose")
	}
	if w.Arrangements[0].Offer == "" || w.DirectorConnection().Offer == "" {
		t.Fatal("saved memory or callback damaged")
	}
	w.Log("Old task", "A quoted old offer.", "story")
	w.Log("Danger", "A public warning.", "danger")
	changes := recentWorldChanges(w)
	for _, r := range changes {
		if r.Kind == "story" {
			t.Fatal("duplicate story prose retained")
		}
	}
	if changes[len(changes)-1].Kind != "danger" {
		t.Fatal("world consequence omitted")
	}
}

func TestFocusedContextKeepsConnectionAndCorrection(t *testing.T) {
	w := core.New(27)
	for i := 0; i < 12; i++ {
		w.Arrangements = append(w.Arrangements, core.ArrangementMemory{Life: 1, Title: "Earlier task", Offer: "Old prose", Status: "completed", Result: "Completed"})
	}
	connection := w.DirectorConnection()
	c := focusedContext(w, "collection", connection, "Correct the speaker", []string{"", "bellandi", "russo"})
	if c["required_connection"] != connection || c["validation_feedback"] != "Correct the speaker" || c["required_operation"] != "collection" {
		t.Fatal("focused context lost required constraints")
	}
	if len(c["recent_arrangements"].([]core.ArrangementMemory)) != 8 || w.Arrangements[0].Offer != "Old prose" {
		t.Fatal("focused history bound or preservation failed")
	}
}
