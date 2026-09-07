package main

import "blackledger/core"

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
	recent := arrangementBriefs(w)
	if len(recent) > 8 {
		recent = recent[len(recent)-8:]
	}
	return map[string]any{"required_operation": operation, "required_connection": connection,
		"validation_feedback": feedback, "allowed_beneficiary_ids": beneficiaries,
		"current_player": w.Player.Name, "current_life": w.Life, "current_location": w.Player.Location,
		"npcs": w.NPCs, "factions": w.Factions, "places": core.Locations, "properties": w.Properties,
		"recent_arrangements": recent}
}
