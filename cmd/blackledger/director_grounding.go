package main

import (
	"blackledger/core"
	"fmt"
	"strings"
)

// Work that exists only because of a war, a seizure or a death at the top has
// to say so. Without this the model accepts the operation and then writes an
// errand: observed live on qwen3.5:35b-a3b, which took a war-derived warning
// job and produced a message delivery that never mentioned the war.
//
// This is a minimum lexical link to a committed fact, in the same spirit as the
// beneficiary mention check. It does not verify that the reference is apt.

// groundingSubjects are the names from the committed fact this work exists on
// account of: the organizations at war, the place that changed hands, whoever
// holds it now, or the person who just took over.
func groundingSubjects(w *core.World, operation string) []string {
	because := ""
	for _, situational := range w.SituationalOperations() {
		if situational.ID == operation {
			because = situational.Because
		}
	}
	if because == "" {
		return nil
	}
	subjects := []string{}
	seen := map[string]bool{}
	add := func(name string) {
		name = strings.TrimSpace(name)
		if name == "" || seen[name] || !strings.Contains(because, name) {
			return
		}
		seen[name] = true
		subjects = append(subjects, name)
	}
	for _, f := range w.Factions {
		add(f.Name)
		add(f.Leader)
	}
	for _, l := range core.Locations {
		add(l.Name)
	}
	return subjects
}

// validateSituationalGrounding requires the body to name at least one subject of
// the fact the work exists because of.
func validateSituationalGrounding(w *core.World, p core.Proposal) error {
	subjects := groundingSubjects(w, p.Operation)
	if len(subjects) == 0 {
		return nil
	}
	body := strings.ReplaceAll(p.Body, "’", "'")
	for _, subject := range subjects {
		if strings.Contains(body, subject) {
			return nil
		}
		// A family is as often called by its distinctive word as its full name.
		for _, word := range strings.Fields(subject) {
			if len(word) > 3 && !genericOrganizationWords[strings.ToLower(word)] && strings.Contains(body, word) {
				return nil
			}
		}
	}
	return fmt.Errorf("this work exists only because of a committed event, and the dialogue never refers to it; name at least one of %q so the request is grounded in what actually happened", subjects)
}
