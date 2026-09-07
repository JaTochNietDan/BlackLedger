package main

import (
	"blackledger/core"
	"errors"
	"fmt"
	"slices"
)

var errDirectorContextChanged = errors.New("director context changed during preparation")

// Called inside the save transaction, after generation and before queuing. Time,
// cash and ordinary work may advance while the model writes; ownership and the
// identity/availability of its speaker must still agree with the supplied facts.
// This guards stale input, not invented claims in otherwise current dialogue.
func validateDirectorFreshness(before, now *core.World, p core.Proposal, connection *core.ArrangementMemory) error {
	stale := func(reason string) error { return fmt.Errorf("%w: %s", errDirectorContextChanged, reason) }
	if len(before.Properties) != len(now.Properties) {
		return stale("property roster changed")
	}
	for id, property := range before.Properties {
		current, ok := now.Properties[id]
		if !ok || current.Owner != property.Owner {
			return stale("property ownership changed")
		}
	}
	original, current := before.NPC(p.Speaker), now.NPC(p.Speaker)
	if original == nil || current == nil || original.Name != current.Name || original.Role != current.Role {
		return stale("speaker identity changed")
	}
	if !slices.Equal(speakerBeneficiaries(before, p.Speaker), speakerBeneficiaries(now, p.Speaker)) || validateSpeakerBeneficiary(now, p.Speaker, p.Beneficiary) != nil {
		return stale("speaker affiliation changed")
	}
	if connection == nil {
		if !eligibleDirectorSpeakers(now)[p.Speaker] {
			return stale("contact is no longer available for new work")
		}
	} else {
		// Established continuations remain possible after goodwill changes, but
		// cannot refer to a removed or altered canonical completion.
		found := false
		for _, memory := range now.Arrangements {
			if memory.ID == connection.ID && memory.Life == now.Life && memory.Status == "completed" && memory.Result == connection.Result && memory.Speaker == connection.Speaker && memory.Beneficiary == connection.Beneficiary {
				found = true
				break
			}
		}
		if !found {
			return stale("completed connection changed")
		}
	}
	return nil
}
