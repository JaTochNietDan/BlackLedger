package core

import (
	"encoding/json"
	"testing"
)

func TestIncendiaryDamagesOnlyTargetAndStartsPersistentFire(t *testing.T) {
	w := bomber(t)
	w.Player.Location = "club"
	p := w.Properties["club"]
	p.Condition = 100
	p.Supply = 40
	staff, bank, cash, heat := p.Staff, p.Bankroll, w.Player.Cash, w.Player.Heat
	other := w.Properties["laundry"].Condition
	if err := w.Incendiary("club"); err != nil {
		t.Fatal(err)
	}
	if p.Condition < 70 || p.Condition > 85 || p.Supply != 30 || p.Staff != staff || p.Bankroll != bank {
		t.Fatalf("unexpected fire damage: %+v", p)
	}
	if w.Player.Cash != cash-IncendiaryCost || w.Player.Heat != min(100, heat+IncendiaryHeat) || w.Properties["laundry"].Condition != other {
		t.Fatal("wrong cost, heat or collateral property damage")
	}
	fires := w.ActiveBuildingFires()
	if len(fires) != 1 || fires[0].Target != "club" {
		t.Fatal(fires)
	}
	last := w.VisualCues[len(w.VisualCues)-1]
	if last.Kind != "incendiary" || last.Target != "club" || last.Attacker == nil || last.Attacker.ID != "player" || last.Attacker.Name != w.Player.Name || last.Attacker.Weapon != 0 {
		t.Fatal("incendiary lost its target or unarmed attacker", last)
	}
	for _, cue := range w.VisualCues {
		if cue.Kind == "explosion" || cue.Kind == "killing" {
			t.Fatal("incendiary borrowed explosive casualties", cue)
		}
	}
	before := w.Player.Cash
	if w.Incendiary("club") == nil || w.Player.Cash != before {
		t.Fatal("paid for repeated attack on active fire")
	}
}

func TestIncendiaryRejectsWrongLocationOwnershipAndFunds(t *testing.T) {
	for _, mode := range []string{"away", "owned", "cash", "wrecked", "invalid"} {
		t.Run(mode, func(t *testing.T) {
			w := bomber(t)
			w.Player.Location = "club"
			id := "club"
			switch mode {
			case "away":
				w.Player.Location = "bar"
			case "owned":
				w.Properties[id].Owner = "player:1"
			case "cash":
				w.Player.Cash = 39
			case "wrecked":
				w.Properties[id].Condition = 0
			case "invalid":
				id = "missing"
			}
			cash := w.Player.Cash
			if w.Incendiary(id) == nil || w.Player.Cash != cash || len(w.ActiveBuildingFires()) != 0 {
				t.Fatal("invalid incendiary changed world")
			}
		})
	}
}

func TestIncendiaryCommandChargesOnceAndPersistsDamage(t *testing.T) {
	w := bomber(t)
	w.Player.Location = "laundry"
	w.Player.Heat = 0
	before, cash := w.Minute, w.Player.Cash
	act(t, &w, "incendiary", "laundry")
	if w.Minute != before+IncendiaryMinutes || w.Player.Cash != cash-IncendiaryCost {
		t.Fatal("wrong command duration or cost")
	}
	condition := w.Properties["laundry"].Condition
	raw, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var saved World
	if err = json.Unmarshal(raw, &saved); err != nil {
		t.Fatal(err)
	}
	if saved.Properties["laundry"].Condition != condition || len(saved.ActiveBuildingFires()) != 1 {
		t.Fatal("fire or damage lost on save")
	}
	saved.Minute = before + 180
	if saved.Properties["laundry"].Condition != condition {
		t.Fatal("extinguishing repaired damage")
	}
}
