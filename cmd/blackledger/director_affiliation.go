package main

import (
	"blackledger/core"
	"fmt"
	"strings"
	"unicode"
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

// A minimum observable link between spoken terms and the faction receiving
// credit. Mentioning a family is not proof of semantic consistency, but an
// offer naming only its rival cannot silently award credit to this beneficiary.
func validateBeneficiaryMention(w *core.World, p core.Proposal) error {
	if p.Beneficiary == "" {
		return nil
	}
	words := func(s string) string {
		return " " + strings.Join(strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }), " ") + " "
	}
	for _, f := range w.Factions {
		if f.ID != p.Beneficiary {
			continue
		}
		body := words(p.Body)
		if strings.Contains(body, words(f.Name)) || strings.Contains(body, words(f.ID)) {
			return nil
		}
		return fmt.Errorf("body must name beneficiary %q and explain how this job serves that family; naming only a rival contradicts the displayed faction reward", f.Name)
	}
	return fmt.Errorf("unknown beneficiary %q", p.Beneficiary)
}

func directorConnection(w *core.World) *core.ArrangementMemory {
	connection := w.DirectorConnection()
	if connection != nil && validateSpeakerBeneficiary(w, connection.Speaker, connection.Beneficiary) != nil {
		// Keep historical saves intact, but do not force another contradictory job.
		return nil
	}
	return connection
}
