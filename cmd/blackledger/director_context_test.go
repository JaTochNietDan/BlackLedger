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
	if len(c["avoid_recent_titles"].([]string)) != 8 || w.Arrangements[0].Offer != "Old prose" {
		t.Fatal("focused history bound or preservation failed")
	}
}

func TestFocusedFollowUpUsesRelevantContactAndCurrentOwnership(t *testing.T) {
	w := core.New(27)
	w.District = 1
	w.Properties["garage"].Owner = "player:1"
	connection := &core.ArrangementMemory{ID: "job", Life: 1, Speaker: "mara", Title: "A garage dispute", Offer: "Resolve the tool disagreement at Russo Motor Works.", Status: "completed"}
	c := focusedContext(w, "collection", connection, "", []string{"", "bellandi", "russo"})
	people := c["npcs"].([]map[string]any)
	places := c["places"].([]map[string]any)
	if len(people) != 1 || people[0]["id"] != "mara" {
		t.Fatal("unrelated contacts in follow-up brief")
	}
	if len(places) != 1 || places[0]["id"] != "garage" || places[0]["owner"] != "current player's organization" {
		t.Fatal("follow-up location or ownership wrong")
	}
	fresh := focusedContext(w, "mediation", nil, "", []string{"", "bellandi", "russo"})
	if len(fresh["npcs"].([]map[string]any)) != len(w.NPCs) || len(fresh["places"].([]map[string]any)) != len(accessibleJobLocations(w)) {
		t.Fatal("fresh requests cannot see the city")
	}
}
