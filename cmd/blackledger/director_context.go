package main

import (
	"blackledger/core"
	"strings"
)

// Keep outcomes and identity without repeatedly quoting the same old offer prose.
// DirectorConnection separately supplies the full one-job callback context.
func arrangementBriefs(w *core.World) []core.ArrangementMemory {
	out := make([]core.ArrangementMemory, len(w.Arrangements))
	copy(out, w.Arrangements)
	for i := range out {
		out[i].Offer = ""
	}
	return out
}
func recentWorldChanges(w *core.World) []core.Record {
	out := []core.Record{}
	for _, r := range w.History {
		if r.Kind != "story" {
			out = append(out, r)
		}
	}
	if len(out) > 12 {
		out = out[len(out)-12:]
	}
	return out
}

// Focused context is an experimental brief; validation remains identical to full mode.
func focusedContext(w *core.World, operation string, connection *core.ArrangementMemory, feedback string, beneficiaries []string) map[string]any {
	people := []map[string]any{}
	for _, n := range w.NPCs {
		if connection != nil && n.ID != connection.Speaker {
			continue
		}
		people = append(people, map[string]any{"id": n.ID, "name": n.Name, "role": n.Role, "trust": n.Trust})
	}
	places := []map[string]any{}
	for _, l := range core.Locations {
		if connection != nil && !strings.Contains(strings.ToLower(connection.Offer+" "+connection.Title), strings.ToLower(l.Name)) {
			continue
		}
		p := w.Properties[l.ID]
		owner := p.Owner
		if w.Own(l.ID) {
			owner = "current player's organization"
		}
		places = append(places, map[string]any{"id": l.ID, "name": l.Name, "owner": owner, "condition": p.Condition})
	}
	recentTitles := []string{}
	for _, m := range w.Arrangements {
		if m.Life == w.Life {
			recentTitles = append(recentTitles, m.Title)
		}
	}
	if len(recentTitles) > 8 {
		recentTitles = recentTitles[len(recentTitles)-8:]
	}
	return map[string]any{"job_brief": jobBrief(w, operation, connection), "required_operation": operation, "required_connection": connection,
		"validation_feedback": feedback, "allowed_beneficiary_ids": beneficiaries,
		"current_player": w.Player.Name, "current_life": w.Life,
		"npcs": people, "factions": w.Factions, "places": places, "avoid_recent_titles": recentTitles}
}
