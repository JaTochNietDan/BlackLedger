package main

import (
	"blackledger/core"
	"strings"
	"testing"
)

func TestGeneratedDialogueStaysInsideTheCity(t *testing.T) {
	// Observed on qwen3.5:35b-a3b: the mediation brief said the disputing staff
	// were "neither the player nor a new named character", and Elena Russo
	// recited that restriction aloud to Alex as "Neither is you".
	for _, text := range []string{
		"At The Mariner, two staff members are locked in a disagreement. Neither is you, but the tension threatens our operations.",
		"Neither of them is you, so stay impartial.",
		"None of them are you; just settle it.",
		"The player should hear both sides.",
		"Do not introduce a new named character here.",
		"Another NPC will meet you at the door.",
		"Follow this brief and keep it quiet.",
		"Pick the careful approach label if you want it quiet.",
	} {
		t.Run(text, func(t *testing.T) {
			if validateInWorldVoice(core.Proposal{Body: text}) == nil {
				t.Fatal("out-of-world line accepted")
			}
		})
	}
	for _, text := range []string{
		"Two of the floor staff are deadlocking over who controls the rear loading bay.",
		"Neither wants to start a fight, but they won't back down.",
		"Hear both sides without taking a side.",
		"Neither of them will budge on the night shift.",
		"You are not the one who has to like it.",
	} {
		t.Run("clean/"+text, func(t *testing.T) {
			if err := validateInWorldVoice(core.Proposal{Body: text}); err != nil {
				t.Fatal("ordinary dialogue rejected:", err)
			}
		})
	}
	if validateInWorldVoice(core.Proposal{Title: "The player's dispute"}) == nil {
		t.Fatal("title meta-language accepted")
	}
	if validateInWorldVoice(core.Proposal{Outcome: "The player settled it."}) == nil {
		t.Fatal("outcome meta-language accepted")
	}
	if validateInWorldVoice(core.Proposal{Approaches: []core.Approach{{Method: "careful", Label: "Ask the NPC quietly"}}}) == nil {
		t.Fatal("approach label meta-language accepted")
	}
}

func roleWorld(name string) *core.World {
	w := &core.World{}
	w.Player.Name = name
	return w
}

func TestSpeakerDoesNotHandThePlayerAFamilyRank(t *testing.T) {
	// Observed on qwen3.5:35b-a3b: Elena Russo, who leads the Russo Outfit,
	// opened with "Alex, as leader of the Russo Outfit, I need you to step in",
	// attaching her own rank to the listener.
	w := roleWorld("Alex Varga")
	for _, text := range []string{
		"Alex, as leader of the Russo Outfit, I need you to step in.",
		"Alex Varga, as the leader of the Bellandi Family, you owe me this.",
		"You, as head of the Russo Outfit, should settle it.",
		"You as boss of the Bellandi Family can end this.",
	} {
		t.Run(text, func(t *testing.T) {
			err := validatePlayerRole(w, core.Proposal{Body: text})
			if err == nil {
				t.Fatal("misattributed family rank accepted")
			}
			if !strings.Contains(err.Error(), "Alex Varga") {
				t.Fatal("correction should name the player:", err)
			}
		})
	}
	for _, text := range []string{
		"As leader of the Russo Outfit, I need this settled quietly.",
		"Alex, I need you to step in at The Mariner.",
		"Elena Russo, as leader of the Russo Outfit, sent me.",
		"You work for yourself, and that is why I am asking.",
		"Alex, as far as the dock crews know, nothing happened.",
	} {
		t.Run("clean/"+text, func(t *testing.T) {
			if err := validatePlayerRole(w, core.Proposal{Body: text}); err != nil {
				t.Fatal("ordinary dialogue rejected:", err)
			}
		})
	}
}
