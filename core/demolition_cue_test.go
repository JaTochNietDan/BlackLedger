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
