package core

import (
	"encoding/json"
	"testing"
)

func TestAftermathSurvivesResultsAndReloadUntilGameClockCleanup(t *testing.T) {
	w := proprietor(t)
	w.NPC("mara").Location = "bar"
	start := w.Minute
	if !w.Kill("mara", "Shot at Saint Agnes.") {
		t.Fatal("fixture death failed")
	}
	scenes := w.ActiveAftermath()
	if len(scenes) != 1 {
		t.Fatalf("aftermath: %+v", scenes)
	}
	scene := scenes[0]
	if scene.Target != "bar" || scene.Victim.ID != "mara" || scene.PoliceAt != start+5 || scene.CleanupAt != start+180 {
		t.Fatalf("wrong scene: %+v", scene)
	}
	w.VisualCues = nil
	w.LastResult = nil
	encoded, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var restored World
	if err = json.Unmarshal(encoded, &restored); err != nil {
		t.Fatal(err)
	}
	if len(restored.ActiveAftermath()) != 1 {
		t.Fatal("reload lost scene")
	}
	restored.Minute = start + 179
	if len(restored.ActiveAftermath()) != 1 {
		t.Fatal("cleaned early")
	}
	restored.Minute = start + 180
	if len(restored.ActiveAftermath()) != 0 {
		t.Fatal("scene survived cleanup")
	}
	if len(restored.Aftermath) != 1 {
		t.Fatal("public read mutated save")
	}
}

func TestAftermathDoesNotInventVictimsOrDuplicateBodies(t *testing.T) {
	w := proprietor(t)
	w.Witness("killing", "bar", "preview", "", "mara")
	if len(w.ActiveAftermath()) != 0 {
		t.Fatal("living person became a body")
	}
	w.NPC("mara").Location = "bar"
	w.Kill("mara", "Shot.")
	w.Witness("killing", "bar", "repeat", "", "mara")
	if len(w.ActiveAftermath()) != 1 {
		t.Fatal("duplicate body")
	}
	scenes := w.ActiveAftermath()
	scenes[0].Target = "elsewhere"
	if w.ActiveAftermath()[0].Target != "bar" {
		t.Fatal("projection aliases save")
	}
}

func TestOldDeathCannotRestartCleanup(t *testing.T) {
	w := proprietor(t)
	w.NPC("mara").Location = "bar"
	w.Kill("mara", "Shot.")
	w.Minute += 200
	w.Witness("killing", "bar", "old report", "", "mara")
	if len(w.ActiveAftermath()) != 0 {
		t.Fatal("old report recreated a cleaned scene")
	}
}

func TestOfficialKillingAlsoLeavesPublicAftermath(t *testing.T) {
	w := proprietor(t)
	mayor := w.NPC("mayor")
	if mayor == nil {
		t.Fatal("missing mayor")
	}
	mayor.Location = "bar"
	w.Kill(mayor.ID, "Shot at Saint Agnes.")
	scenes := w.Public()["aftermath"].([]Aftermath)
	if len(scenes) != 1 || scenes[0].Victim.ID != mayor.ID {
		t.Fatalf("official missing: %+v", scenes)
	}
}
