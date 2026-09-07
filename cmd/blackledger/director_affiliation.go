package main

import (
	"blackledger/core"
	"fmt"
)

// Independent fixers and associates can bring work for either family. A family
// leader's ordinary job must serve their own organization; betrayal would need
// its own supported plot and consequences, which this operation does not have.
func speakerBeneficiaries(w *core.World, speaker string) []string {
	npc := w.NPC(speaker)
	if npc != nil {
		for _, f := range w.Factions {
			if f.Leader == npc.Name {
				return []string{f.ID}
			}
		}
	}
	ids := []string{""}
	for _, f := range w.Factions {
		ids = append(ids, f.ID)
	}
	return ids
}

func validateSpeakerBeneficiary(w *core.World, speaker, beneficiary string) error {
	allowed := speakerBeneficiaries(w, speaker)
	for _, id := range allowed {
		if id == beneficiary {
			return nil
		}
	}
	return fmt.Errorf("speaker %q must use exactly one of %q as beneficiary; ordinary jobs cannot invent a betrayal or change faction allegiance", speaker, allowed)
}

func directorConnection(w *core.World) *core.ArrangementMemory {
	connection := w.DirectorConnection()
	if connection != nil && validateSpeakerBeneficiary(w, connection.Speaker, connection.Beneficiary) != nil {
		// Keep historical saves intact, but do not force another contradictory job.
		return nil
	}
	return connection
}
