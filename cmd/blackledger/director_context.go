package main

import (
	"blackledger/core"
	"fmt"
	"strings"
)

type attributedArrangement struct {
	core.ArrangementMemory
	Participant     string `json:"participant"`
	OriginalRequest string `json:"original_request_claims_not_verified,omitempty"`
}

func historicalParticipant(w *core.World, life int) string {
	if life == w.Life {
		return w.Player.Name
	}
	for _, dead := range w.Dead {
		if dead.Life == life {
			return dead.Name
		}
	}
	return fmt.Sprintf("unknown former person from life %d", life)
}

func attributedMemory(w *core.World, m core.ArrangementMemory, includeRequest bool) attributedArrangement {
	claims := ""
	if includeRequest {
		claims = m.Offer
	}
	m.Offer = ""
	return attributedArrangement{ArrangementMemory: m, Participant: historicalParticipant(w, m.Life), OriginalRequest: claims}
}

// Separate the current protagonist's memories from city history. In particular,
// a prior life's second-person result must not read as the new person's work.
// Preserve original request prose only for the required callback, explicitly as
// claims: completing a job does not make every sentence in its offer true.
func attributeDirectorContext(c map[string]any, w *core.World, connection *core.ArrangementMemory) {
	current, former := []attributedArrangement{}, []attributedArrangement{}
	for _, m := range arrangementBriefs(w) {
		if m.Life == w.Life {
			current = append(current, attributedMemory(w, m, false))
		} else {
			former = append(former, attributedMemory(w, m, false))
		}
	}
	c["recent_arrangements"] = current
	c["previous_people_arrangements"] = former
	c["required_connection"] = nil
	if connection != nil {
		c["required_connection"] = attributedMemory(w, *connection, true)
	}
	currentHistory, formerHistory := []map[string]any{}, []map[string]any{}
	for _, r := range recentWorldChanges(w) {
		row := map[string]any{"participant": historicalParticipant(w, r.Life), "record": r}
		if r.Life == w.Life {
			currentHistory = append(currentHistory, row)
		} else {
			formerHistory = append(formerHistory, row)
		}
	}
	c["recent_history"] = currentHistory
	c["previous_people_history"] = formerHistory
	speakers := []map[string]string{}
	for _, id := range directorSpeakers(w, connection) {
		n := w.NPC(id)
		if n == nil {
			continue
		}
		relationship := "established contact requesting a job"
		for _, crew := range w.Player.Crew {
			if crew.ID == id {
				relationship = "the player's employee bringing their boss a lead"
			}
		}
		for _, f := range w.Factions {
			if f.Leader == n.Name {
				relationship = "leader of " + f.Name + ", requesting a favor from the player"
			}
		}
		speakers = append(speakers, map[string]string{"id": id, "name": n.Name, "relationship_to_addressee": relationship})
	}
	c["dialogue_participants"] = map[string]any{"addressee": w.Player.Name, "addressee_life": w.Life, "allowed_speakers": speakers}
}

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
		if l.District > w.District {
			continue
		}
		if connection != nil && connection.Location != l.ID && !strings.Contains(strings.ToLower(connection.Offer+" "+connection.Title), strings.ToLower(l.Name)) {
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

// The city's own quarrels are ordinary generation context. The director may
// narrate a war, a breakaway or a seizure the simulation actually committed; it
// may not invent one, and the existing guards still check every claim it makes
// about who owns what and who serves whom.
func cityConflicts(w *core.World) []map[string]any {
	out := []map[string]any{}
	for _, c := range w.PublicConflicts() {
		out = append(out, map[string]any{
			"between": c.Between, "state": c.State, "since_minute": c.Since,
		})
	}
	return out
}

// organizationHoldings reports what each organization currently holds, by the
// names a character would use aloud.
func organizationHoldings(w *core.World) map[string][]string {
	out := map[string][]string{}
	for _, f := range w.Factions {
		places := []string{}
		for _, id := range w.FamilyHoldings(f.ID) {
			if place, ok := core.PlaceByID(id); ok {
				places = append(places, place.Name)
			}
		}
		out[f.Name] = places
	}
	return out
}
