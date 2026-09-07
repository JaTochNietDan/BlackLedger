package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAttackPresentationContainsOnlyCommittedPublicResult(t *testing.T) {
	w := New(27)
	w.Player.Location = "laundry"
	w.Properties["laundry"].Owner = "player:1"
	w.Plots = []Plot{{ID: "private-plot", Kind: "sabotage", Life: 1, Due: 495, Actor: "bellandi", Target: "laundry", Strength: 35}}
	act(t, &w, "wait", "laundry")
	if w.Properties["laundry"].Condition != 65 || len(w.LastResult.Cues) != 1 {
		t.Fatal("missing committed attack cue")
	}
	cue := w.LastResult.Cues[0]
	if cue.Target != "laundry" || cue.Kind != "attack" || !strings.Contains(cue.Caption, "65%") {
		t.Fatal("wrong cue", cue)
	}
	data, _ := json.Marshal(cue)
	if strings.Contains(string(data), "bellandi") || strings.Contains(string(data), "private-plot") {
		t.Fatal("hidden plot leaked")
	}
	cloned := w.Clone()
	if len(cloned.LastResult.Cues) != 1 {
		t.Fatal("saved result lost presentation")
	}
	act(t, &w, "wait", "laundry")
	if len(w.LastResult.Cues) != 0 {
		t.Fatal("old attack replayed as new outcome")
	}
}

func TestUnresolvedAttackDoesNotEmitOutcomeScene(t *testing.T) {
	w := New(27)
	w.Player.Security = 3
	w.Retaliation()
	act(t, &w, "rest", "room")
	if w.Event == nil || w.Event.Kind != "attack" || len(w.LastResult.Cues) != 0 {
		t.Fatal("unresolved outcome was presented")
	}
	choice(t, &w, "defend")
	if len(w.LastResult.Cues) != 1 || w.LastResult.Cues[0].Target != "room" {
		t.Fatal("resolved attack lacks scene")
	}
}

func TestPeacefulRecoveryDoesNotInventAnAttack(t *testing.T) {
	w := New(27)
	w.Player.Heat = 20
	w.Player.Location = "market"
	act(t, &w, "lie_low", "market")
	if len(w.LastResult.Cues) != 0 {
		t.Fatal("peaceful action produced an attack scene")
	}
}
