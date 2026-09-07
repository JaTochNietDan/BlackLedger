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

func TestDirectorAttributesFormerLifeWithoutGivingNewPersonItsAchievements(t *testing.T) {
	w := core.New(27)
	w.Arrangements = []core.ArrangementMemory{{Life: 1, Speaker: "mara", Operation: "mediation", Status: "completed", Offer: "An unverified allegation.", Result: "You completed the mediation."}}
	w.Log("Work finished", "You completed the mediation.", "result")
	w.Die("Previous life ended.")
	next, err := core.Execute(w, core.Command{Kind: "new_life", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	c := map[string]any{}
	attributeDirectorContext(c, next, nil)
	if len(c["recent_arrangements"].([]attributedArrangement)) != 0 || c["required_connection"] != nil {
		t.Fatal("old work attributed to newcomer")
	}
	old := c["previous_people_arrangements"].([]attributedArrangement)
	if len(old) != 1 || old[0].Participant != w.Player.Name || old[0].Result != "You completed the mediation." || old[0].Offer != "" || old[0].OriginalRequest != "" {
		t.Fatal("former identity or canonical result lost", old)
	}
	if len(c["previous_people_history"].([]map[string]any)) == 0 {
		t.Fatal("former history omitted")
	}
	for _, row := range c["previous_people_history"].([]map[string]any) {
		if row["participant"] != w.Player.Name {
			t.Fatal("history participant wrong")
		}
	}
	participants := c["dialogue_participants"].(map[string]any)
	if participants["addressee"] != next.Player.Name || participants["addressee"] == w.Player.Name {
		t.Fatal("new addressee missing")
	}
	if next.Arrangements[0].Offer != "An unverified allegation." {
		t.Fatal("save memory mutated")
	}
}

func TestDirectorCallbackSeparatesClaimsAndNamesSpeakerHierarchy(t *testing.T) {
	w := core.New(27)
	w.Player.Crew = []core.Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 65}}
	m := core.ArrangementMemory{Life: w.Life, Speaker: "leo", Operation: "courier", Status: "completed", Offer: "The supplier secretly owns this building.", Result: "You completed the requested delivery."}
	c := map[string]any{}
	attributeDirectorContext(c, w, &m)
	connection := c["required_connection"].(attributedArrangement)
	if connection.Offer != "" || connection.OriginalRequest != m.Offer || connection.Result != m.Result || connection.Participant != w.Player.Name {
		t.Fatal("callback claims promoted or attribution lost")
	}
	participants := c["dialogue_participants"].(map[string]any)
	speakers := participants["allowed_speakers"].([]map[string]string)
	if len(speakers) != 1 || speakers[0]["id"] != "leo" || speakers[0]["relationship_to_addressee"] != "the player's employee bringing their boss a lead" {
		t.Fatal("crew hierarchy missing", speakers)
	}
}
