package main

import (
	"blackledger/core"
	"encoding/json"
	"os"
	"testing"
)

// Replays proposals captured from a real local-model run against the current
// guards. This needs no provider: the recorded prose is the fixture. It proves
// the guards act on actual model output rather than only on hand-written
// strings, and it fails if a later change stops catching a known defect.
type recordedSuite struct {
	Model string `json:"model"`
	Cases []struct {
		Name   string `json:"name"`
		Life   int    `json:"life"`
		Player struct {
			Name string `json:"name"`
		} `json:"player"`
		PreviousArrangements []core.ArrangementMemory `json:"previous_arrangements"`
		Offers               []struct {
			Event struct {
				Title   string `json:"title"`
				Body    string `json:"body"`
				Speaker string `json:"speaker"`
				Outcome string `json:"outcome"`
				Choices []struct {
					ID    string `json:"id"`
					Label string `json:"label"`
				} `json:"choices"`
			} `json:"event"`
		} `json:"offers"`
	} `json:"cases"`
}

func TestRecordedProposalsAreJudgedByCurrentGuards(t *testing.T) {
	t.Parallel()
	const path = "../../docs/director-qwen35-moe-focused-scenarios.json"
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Skip("recorded suite not present:", err)
	}
	var suite recordedSuite
	if err := json.Unmarshal(raw, &suite); err != nil {
		t.Fatal(err)
	}
	// What manual review of that run found, case by case.
	expected := map[string]string{
		"crew_after_declined_delivery": "",
		"leader_new_request":           "in-world voice",
		"completed_tool_dispute":       "",
		"new_person_after_death":       "familiarity or title",
	}
	if len(suite.Cases) != len(expected) {
		t.Fatalf("recorded suite has %d cases, expected %d", len(suite.Cases), len(expected))
	}
	for _, c := range suite.Cases {
		t.Run(c.Name, func(t *testing.T) {
			if len(c.Offers) == 0 {
				t.Fatal("recorded case has no offer")
			}
			event := c.Offers[0].Event
			w := &core.World{Life: c.Life, Arrangements: c.PreviousArrangements}
			w.Player.Name = c.Player.Name
			proposal := core.Proposal{Title: event.Title, Body: event.Body, Speaker: event.Speaker, Outcome: event.Outcome}
			for _, choice := range event.Choices {
				// Only generated approaches carry model wording; accept/decline are ours.
				if choice.ID != "accept" && choice.ID != "decline" {
					proposal.Approaches = append(proposal.Approaches, core.Approach{Method: choice.ID, Label: choice.Label})
				}
			}
			verdicts := map[string]error{
				"in-world voice": validateInWorldVoice(proposal),
				"player role":    validatePlayerRole(w, proposal),
				"familiarity":    validateEarnedFamiliarity(w, proposal),
				"scene title":    validateSceneTitle(proposal),
				"spoken terms":   validateSpokenTerms(proposal),
			}
			caught := []string{}
			for name, err := range verdicts {
				if err != nil {
					caught = append(caught, name+": "+err.Error())
				}
			}
			for _, line := range caught {
				t.Log("rejected by", line)
			}
			if expected[c.Name] == "" && len(caught) > 0 {
				t.Fatal("previously acceptable proposal is now rejected; check for a false positive")
			}
			if expected[c.Name] != "" && len(caught) == 0 {
				t.Fatal("known defective proposal passed every guard")
			}
		})
	}
}
