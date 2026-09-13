package core

import (
	"encoding/json"
	"testing"
)

func TestRaidPresenceSurvivesResultReplacementAndReload(t *testing.T) {
	w := proprietor(t)
	start := w.Minute
	w.Witness("raid", "bar", "Police search the premises.", "")
	w.VisualCues = nil
	w.LastResult = nil
	data, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var restored World
	if err = json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	got := restored.Public()["police_presence"].([]PolicePresence)
	if len(got) != 1 || got[0].Target != "bar" || got[0].CleanupAt != start+120 {
		t.Fatalf("presence lost: %+v", got)
	}
	got[0].Target = "changed"
	if restored.ActivePolicePresence()[0].Target != "bar" {
		t.Fatal("projection aliases save")
	}
	restored.Minute = start + 119
	if len(restored.ActivePolicePresence()) != 1 {
		t.Fatal("cleaned early")
	}
	restored.Minute = start + 120
	if len(restored.ActivePolicePresence()) != 0 {
		t.Fatal("not cleaned at deadline")
	}
	if len(restored.PolicePresence) != 1 {
		t.Fatal("GET mutated save")
	}
}

func TestRaidPresenceMergesAddressWithoutInventingOtherScenes(t *testing.T) {
	w := proprietor(t)
	w.Witness("arrest", "bar", "Arrest.", "")
	w.Witness("raid", "unknown-place", "Invalid.", "")
	if len(w.ActivePolicePresence()) != 0 {
		t.Fatal("invented cordon")
	}
	w.Witness("raid", "bar", "Raid.", "")
	first := w.ActivePolicePresence()[0]
	w.recordPolicePresence(w.VisualCues[len(w.VisualCues)-1])
	if len(w.ActivePolicePresence()) != 1 || w.ActivePolicePresence()[0] != first {
		t.Fatal("duplicate restarted cordon")
	}
	w.Minute += 60
	w.Witness("raid", "bar", "Second raid.", "")
	got := w.ActivePolicePresence()
	if len(got) != 1 || got[0].ID != first.ID || got[0].CleanupAt != w.Minute+120 {
		t.Fatalf("wrong merged cordon: %+v", got)
	}
}
