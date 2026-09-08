package main

import (
	"blackledger/core"
	"strings"
	"testing"
)

func familiarityWorld(life int, arrangements ...core.ArrangementMemory) *core.World {
	w := &core.World{Life: life, Arrangements: arrangements}
	w.Player.Name = "Nico Ward"
	return w
}

func TestStrangersDoNotClaimPriorDealingsWithThePlayer(t *testing.T) {
	// Observed on qwen3.5:35b-a3b: Mara greeted a brand-new life-two protagonist
	// with "I need you back at The Mariner", inheriting the dead Alex's standing.
	w := familiarityWorld(2, core.ArrangementMemory{Life: 1, Speaker: "mara", Status: "completed", Title: "Alex's old mediation"})
	for _, text := range []string{
		"I need you back at The Mariner. Two staff members are arguing over the work area.",
		"Welcome back. Two staff members are arguing over the work area.",
		"Good to see you again. Settle this quietly.",
		"You're back sooner than I expected.",
		"Handle it like last time and nobody gets hurt.",
		"Same as before: hear them both out.",
		"Last time you kept it quiet, so keep it quiet.",
		"You know the drill, so go and settle it.",
		"Our usual terms apply to this one.",
		"Settle it as always, without a scene.",
		"You've done this before, so it should be quick.",
		"You've handled worse than two arguing porters.",
	} {
		t.Run(text, func(t *testing.T) {
			err := validateEarnedFamiliarity(w, core.Proposal{Speaker: "mara", Body: text})
			if err == nil {
				t.Fatal("unearned familiarity accepted")
			}
			if !strings.Contains(err.Error(), "Nico Ward") {
				t.Fatal("correction should name the current protagonist:", err)
			}
		})
	}
}

func TestFirstMeetingProseAndOrdinaryDirectionsAreAccepted(t *testing.T) {
	w := familiarityWorld(2, core.ArrangementMemory{Life: 1, Speaker: "mara", Status: "completed"})
	for _, text := range []string{
		"Two staff members are arguing over who controls the rear loading bay.",
		"Collect the payment and bring it back to the manager.",
		"Report back once they have agreed.",
		"Get back to work once the dispute is settled.",
		"Do not come back empty-handed.",
		"Bring them back to the garage without a fight.",
		"You handled yourself well at the door just now.",
	} {
		t.Run(text, func(t *testing.T) {
			if err := validateEarnedFamiliarity(w, core.Proposal{Speaker: "mara", Body: text}); err != nil {
				t.Fatal("ordinary first-meeting prose rejected:", err)
			}
		})
	}
}

func TestArrangementsInThisLifeEarnFamiliarLanguage(t *testing.T) {
	body := "Welcome back. Two staff members are arguing over the work area."
	completed := familiarityWorld(1, core.ArrangementMemory{Life: 1, Speaker: "mara", Status: "completed"})
	if err := validateEarnedFamiliarity(completed, core.Proposal{Speaker: "mara", Body: body}); err != nil {
		t.Fatal("completed work with this speaker should earn familiarity:", err)
	}
	// Being turned down is still a shared past between these two people.
	declined := familiarityWorld(1, core.ArrangementMemory{Life: 1, Speaker: "mara", Status: "declined"})
	if err := validateEarnedFamiliarity(declined, core.Proposal{Speaker: "mara", Body: body}); err != nil {
		t.Fatal("a declined arrangement should still earn familiarity:", err)
	}
	// A different speaker in the same life has met nobody on the player's behalf.
	other := familiarityWorld(1, core.ArrangementMemory{Life: 1, Speaker: "leo", Status: "completed"})
	if validateEarnedFamiliarity(other, core.Proposal{Speaker: "mara", Body: body}) == nil {
		t.Fatal("another contact's history should not earn this speaker familiarity")
	}
	// The dead protagonist's record belongs to the city, not to the new person.
	previous := familiarityWorld(2, core.ArrangementMemory{Life: 1, Speaker: "mara", Status: "completed"})
	if validateEarnedFamiliarity(previous, core.Proposal{Speaker: "mara", Body: body}) == nil {
		t.Fatal("a previous life should not earn familiarity")
	}
}

func TestFamiliarityIsCheckedInTitlesAndApproachLabels(t *testing.T) {
	w := familiarityWorld(2)
	if validateEarnedFamiliarity(w, core.Proposal{Speaker: "mara", Title: "Welcome back to The Mariner"}) == nil {
		t.Fatal("title familiarity accepted")
	}
	if validateEarnedFamiliarity(w, core.Proposal{Speaker: "mara", Approaches: []core.Approach{{Method: "careful", Label: "Settle it like last time"}}}) == nil {
		t.Fatal("approach label familiarity accepted")
	}
	if validateEarnedFamiliarity(w, core.Proposal{Speaker: "mara", Body: "Hear both sides, then decide."}) != nil {
		t.Fatal("clean proposal rejected")
	}
}
