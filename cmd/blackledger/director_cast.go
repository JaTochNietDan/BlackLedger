package main

import "blackledger/core"

// New requests come through established contacts. Continuations preserve their
// saved speaker even when relationships have changed since the original offer.
// Among available contacts prefer the least recently heard, including queued
// requests, so preparing the next scene does not repeatedly select one person.
func directorSpeakers(w *core.World, connection *core.ArrangementMemory) []string {
	if connection != nil {
		return []string{connection.Speaker}
	}
	eligible := eligibleDirectorSpeakers(w)
	seen := map[string]int{}
	order := 0
	for _, memory := range w.Arrangements {
		if memory.Life == w.Life {
			order++
			seen[memory.Speaker] = order
		}
	}
	if w.Event != nil && w.Event.Kind == "proposal" {
		order++
		seen[w.Event.Speaker] = order
	}
	for _, offer := range w.Offers {
		if offer.Event != nil {
			order++
			seen[offer.Event.Speaker] = order
		}
	}
	result := []string{}
	oldest := order + 1
	for _, npc := range w.NPCs {
		if !eligible[npc.ID] {
			continue
		}
		if seen[npc.ID] < oldest {
			result = nil
			oldest = seen[npc.ID]
		}
		if seen[npc.ID] == oldest {
			result = append(result, npc.ID)
		}
	}
	return result
}

func eligibleDirectorSpeakers(w *core.World) map[string]bool {
	eligible := map[string]bool{"mara": true}
	for _, crew := range w.Player.Crew {
		if crew.Loyalty >= 30 {
			eligible[crew.ID] = true
		}
	}
	for _, faction := range w.Factions {
		if faction.Goodwill < 6 {
			continue
		}
		for _, npc := range w.NPCs {
			if npc.Name == faction.Leader {
				eligible[npc.ID] = true
			}
		}
	}
	return eligible
}
