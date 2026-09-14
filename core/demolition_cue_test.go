package core

import (
	"encoding/json"
	"testing"
)

func TestPlantCuePreservesOutcomeAndPlanter(t *testing.T) {
	seen := map[string]bool{}
	for seed := uint32(1); seed <= 64 && len(seen) < 2; seed++ {
		w := bomber(t)
		w.Player.Location, w.Player.Charges = "club", 1
		w.RNG = seed * 2654435761
		name, minute := w.Player.Name, w.Minute
		act(t, &w, "plant", "club")
		if w.LastResult == nil {
			t.Fatal("missing command result")
		}
		var cue *VisualCue
		for i := range w.LastResult.Cues {
			c := &w.LastResult.Cues[i]
			if c.Kind == "explosion" && c.Target == "club" && c.Minute == minute {
				cue = c
				break
			}
		}
		if cue == nil {
			t.Fatal("missing explosion cue")
		}
		want := "premature"
		if w.Properties["club"].Condition < 100 {
			want = "planted"
		}
		if cue.Detonation != want || cue.Attacker == nil || cue.Attacker.ID != "player" || cue.Attacker.Name != name || cue.Attacker.Weapon != 0 {
			t.Fatalf("wrong causal cast: %+v", cue)
		}
		seen[want] = true
		raw, err := json.Marshal(w)
		if err != nil {
			t.Fatal(err)
		}
		var restored World
		if err := json.Unmarshal(raw, &restored); err != nil {
			t.Fatal(err)
		}
		found := false
		for _, c := range restored.LastResult.Cues {
			if c.ID == cue.ID {
				found = c.Detonation == want && c.Attacker != nil && c.Attacker.Name == name
			}
		}
		if !found {
			t.Fatal("save lost demolition outcome or planter")
		}
	}
	if len(seen) != 2 {
		t.Fatalf("did not exercise both outcomes: %v", seen)
	}
}

func TestFactionDetonationDoesNotInventAnIndividualPlanter(t *testing.T) {
	w := bomber(t)
	w.detonate("club", "A faction charge.")
	cue := w.VisualCues[len(w.VisualCues)-1]
	if cue.Kind != "explosion" || cue.Detonation != "planted" || cue.Attacker != nil {
		t.Fatalf("invented faction planter: %+v", cue)
	}
	var legacy VisualCue
	if err := json.Unmarshal([]byte(`{"kind":"explosion","target":"club"}`), &legacy); err != nil {
		t.Fatal(err)
	}
	if legacy.Detonation != "" || legacy.Attacker != nil {
		t.Fatal("legacy cue invented a preamble")
	}
}

func TestFactionPlanterLeavesBeforeTheBlastAndSurvivesSave(t *testing.T) {
	w := bomber(t)
	w.NPCs = []NPC{{ID: "planter", Name: "Avery Holt", Faction: "bellandi", Location: "club", Home: "room", Skill: 70}}
	minute := w.Minute
	w.factionDetonation("bellandi", "club", "A faction charge.")
	n := w.NPC("planter")
	cue := w.VisualCues[len(w.VisualCues)-1]
	if cue.Attacker == nil || cue.Attacker.ID != n.ID || cue.Attacker.Name != n.Name || cue.Attacker.Weapon != 0 || cue.Detonation != "planted" {
		t.Fatalf("wrong planter: %+v", cue)
	}
	if n.Dead || !w.Travelling(n) || n.Location != "club" || n.Heading != "room" || w.Minute != minute {
		t.Fatalf("wrong departure: %+v", n)
	}
	raw, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var restored World
	if err = json.Unmarshal(raw, &restored); err != nil {
		t.Fatal(err)
	}
	if !restored.Travelling(restored.NPC("planter")) {
		t.Fatal("save lost departure")
	}
	restored.Minute = n.Arrives
	restored.Arrivals()
	if restored.NPC("planter").Location != "room" || restored.Travelling(restored.NPC("planter")) {
		t.Fatal("escape did not arrive")
	}
}

func TestFactionPlanterDoesNotUseAbsentOrUnavailablePeople(t *testing.T) {
	for _, change := range []func(*NPC){
		func(n *NPC) { n.Dead = true }, func(n *NPC) { n.Faction = "fassano" }, func(n *NPC) { n.Location = "bar" },
		func(n *NPC) { n.Held = 999999 }, func(n *NPC) { n.Heading = "room"; n.Arrives = 999999 },
		func(n *NPC) { n.Heading = "room"; n.Sets = 999999 }, func(n *NPC) { n.Home = "club" }, func(n *NPC) { n.Home = "missing" },
	} {
		w := bomber(t)
		n := NPC{ID: "unavailable", Name: "Avery Holt", Faction: "bellandi", Location: "club", Home: "room"}
		change(&n)
		w.NPCs = []NPC{n}
		w.factionDetonation("bellandi", "club", "A faction charge.")
		if cue := w.VisualCues[len(w.VisualCues)-1]; cue.Attacker != nil {
			t.Fatalf("borrowed unavailable planter: %+v", n)
		}
		got := w.NPC(n.ID)
		if got.Heading != n.Heading || got.Arrives != n.Arrives || got.Sets != n.Sets {
			t.Fatal("interrupted an existing schedule")
		}
	}
}

func TestFactionPlanterEscapesEvenWhenAnOccupantDies(t *testing.T) {
	casualties := 0
	for seed := uint32(1); seed <= 64; seed++ {
		w := bomber(t)
		w.WorldRNG = seed * 2654435761
		w.NPCs = []NPC{
			{ID: "planter", Name: "Avery Holt", Faction: "bellandi", Location: "club", Home: "room", Skill: 70},
			{ID: "occupant", Name: "Morgan Dale", Location: "club", Home: "room"},
		}
		w.factionDetonation("bellandi", "club", "A faction charge.")
		cue := w.VisualCues[len(w.VisualCues)-1]
		linked := false
		for _, story := range w.News {
			if story.Headline == cue.Headline && story.Minute == cue.Minute {
				linked = true
			}
		}
		if !linked {
			t.Fatal("blast cue does not identify its newspaper story")
		}
		if w.NPC("planter").Dead {
			t.Fatal("departed planter became an indoor casualty")
		}
		if w.NPC("occupant").Dead {
			casualties++
		}
	}
	if casualties == 0 {
		t.Fatal("did not exercise a fatal blast")
	}
}
