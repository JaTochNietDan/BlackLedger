package core

import (
	"encoding/json"
	"testing"
)

func TestDetonationLeavesSavedFireWithoutRepairingDamage(t *testing.T) {
	w := proprietor(t)
	start := w.Minute
	w.detonate("laundry", "Test charge.")
	fires := w.Public()["building_fires"].([]BuildingFire)
	if len(fires) != 1 || fires[0].Target != "laundry" || fires[0].BrigadeAt != start+10 || fires[0].ExtinguishedAt != start+45 || fires[0].CleanupAt != start+90 {
		t.Fatalf("bad fire: %+v", fires)
	}
	condition := w.Properties["laundry"].Condition
	if condition >= 100 {
		t.Fatal("detonation did not damage building")
	}
	w.LastResult = nil
	w.VisualCues = nil
	raw, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var restored World
	if err = json.Unmarshal(raw, &restored); err != nil {
		t.Fatal(err)
	}
	restored.Minute = start + 45
	if len(restored.ActiveBuildingFires()) != 1 {
		t.Fatal("brigade disappeared at extinguishing")
	}
	if restored.Properties["laundry"].Condition != condition {
		t.Fatal("extinguishing repaired building")
	}
	restored.Minute = start + 90
	if len(restored.ActiveBuildingFires()) != 0 {
		t.Fatal("cleanup did not expire")
	}
	if len(restored.BuildingFires) != 1 {
		t.Fatal("public read mutated save")
	}
}
func TestRepeatedDetonationRenewsFireButNotAnArrivedBrigade(t *testing.T) {
	w := proprietor(t)
	w.Witness("explosion", "bar", "An early charge.", "")
	if len(w.ActiveBuildingFires()) != 0 {
		t.Fatal("cue alone invented a building fire")
	}
	w.detonate("laundry", "First charge.")
	first := w.ActiveBuildingFires()[0]
	copy := w.ActiveBuildingFires()
	copy[0].Target = "changed"
	if w.ActiveBuildingFires()[0].Target != "laundry" {
		t.Fatal("projection aliases save")
	}
	w.igniteBuilding("laundry")
	if w.ActiveBuildingFires()[0] != first {
		t.Fatal("duplicate ignition changed deadlines")
	}
	w.Minute += 20
	w.detonate("laundry", "Second charge.")
	got := w.ActiveBuildingFires()
	if len(got) != 1 || got[0].ID != first.ID || got[0].BrigadeAt != first.BrigadeAt || got[0].ExtinguishedAt != w.Minute+45 {
		t.Fatalf("bad renewed fire: %+v", got)
	}
}
