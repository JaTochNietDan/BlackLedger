package core

import (
	"encoding/json"
	"fmt"
	"testing"
)

func buildingDriveByFixture(t *testing.T, seed uint32) *World {
	t.Helper()
	w := New(seed)
	w.MigrateLivingWorld()
	w.Player.Location = "club"
	w.Player.Cash = 20000
	w.Player.Weapon = 3
	w.Player.Car, w.Player.CarWear = 1, 100
	w.Player.Fuel, w.Player.Fuelled = 30, 1
	w.Player.Crew = []Crew{{ID: "leo", Name: "Outdated crew name", Loyalty: 90}}
	n := w.NPC("leo")
	n.Location, n.Heading, n.Arrives, n.Sets = "club", "", 0, 0
	n.Hurt, n.Dead = false, false
	w.Properties["club"].Condition = 100
	w.Properties["club"].Supply = 40
	w.Tasks = nil
	return w
}

func TestBuildingDriveByResolvedDamageAndCapturedCast(t *testing.T) {
	for tier := 1; tier <= 3; tier++ {
		for seed := uint32(1); seed <= 32; seed++ {
			t.Run(fmt.Sprintf("weapon%d/seed%d", tier, seed), func(t *testing.T) {
				w := buildingDriveByFixture(t, seed)
				w.Player.Weapon, w.Player.Car = tier, tier
				p := w.Properties["club"]
				staff, bank, cash, heat, minute := p.Staff, p.Bankroll, w.Player.Cash, w.Player.Heat, w.Minute
				respect, health, charges := w.Player.Respect, w.Player.Health, w.Player.Charges
				people, _ := json.Marshal(w.NPCs)
				other := w.Properties["laundry"].Condition
				if err := w.BuildingDriveBy("club"); err != nil {
					t.Fatal(err)
				}
				damage := 100 - p.Condition
				if damage < 8+6*tier || damage > 13+6*tier || p.Supply != 40-damage/3 || !p.Trouble {
					t.Fatalf("unexpected damage: %+v", p)
				}
				if p.Staff != staff || p.Bankroll != bank || w.Properties["laundry"].Condition != other || len(w.ActiveBuildingFires()) != 0 {
					t.Fatal("invented collateral loss or fire")
				}
				if w.Player.Cash != cash-BuildingDriveByCost || w.Player.Heat != min(100, heat+BuildingDriveByHeat) || w.Fuel() != 29 || w.Minute != minute {
					t.Fatal("wrong cost, fuel, heat or premature clock advancement")
				}
				if w.Player.Respect != respect || w.Player.Health != health || w.Player.Charges != charges {
					t.Fatal("unrelated player consequences")
				}
				afterPeople, _ := json.Marshal(w.NPCs)
				if string(people) != string(afterPeople) {
					t.Fatal("property attack changed occupants")
				}
				if len(w.VisualCues) != 1 {
					t.Fatal(w.VisualCues)
				}
				cue := w.VisualCues[0]
				if cue.Kind != "driveby-building" || cue.Target != "club" || cue.Minute != minute || cue.Gravity != 7 || cue.Headline == "" {
					t.Fatal(cue)
				}
				if cue.Attacker == nil || cue.Attacker.ID != "player" || cue.Attacker.Name != w.Player.Name || cue.Attacker.Weapon != tier {
					t.Fatal(cue.Attacker)
				}
				d := cue.DriveBy
				if d == nil || d.Driver.ID != "leo" || d.Driver.Name != w.NPC("leo").Name || d.VehicleTier != tier || d.Vehicle != VehicleByTier(tier).Label || d.ConditionBefore != 100 || d.ConditionAfter != p.Condition {
					t.Fatal(d)
				}
				if len(cue.Actors) != 1 || cue.Actors[0] != d.Driver {
					t.Fatal("driver missing from cast", cue.Actors)
				}
				// A later load and equipment/name changes cannot rewrite the event.
				// Commands persist presentation cues in LastResult, not VisualCues.
				w.LastResult = &Result{Cues: w.VisualCues}
				saved := w.Clone()
				saved.Player.Weapon, saved.Player.Car = 0, 0
				saved.NPC("leo").Name = "Changed after the event"
				saved.Properties["club"].Condition = 100
				a, _ := json.Marshal(cue)
				b, _ := json.Marshal(saved.LastResult.Cues[0])
				if string(a) != string(b) {
					t.Fatal("saved scene facts changed")
				}
			})
		}
	}
}

func TestBuildingDriveByRefusalsDoNotMutateWorld(t *testing.T) {
	cases := map[string]func(*World){
		"away":              func(w *World) { w.Player.Location = "bar" },
		"owned":             func(w *World) { w.Properties["club"].Owner = fmt.Sprintf("player:%d", w.Life) },
		"cash":              func(w *World) { w.Player.Cash = BuildingDriveByCost - 1 },
		"wrecked property":  func(w *World) { w.Properties["club"].Condition = 0 },
		"no business":       func(w *World) { delete(w.Properties, "club") },
		"no income":         func(w *World) { w.Properties["club"].Income = 0 },
		"burning":           func(w *World) { w.igniteBuilding("club") },
		"unarmed":           func(w *World) { w.Player.Weapon = 0 },
		"invalid firearm":   func(w *World) { w.Player.Weapon = 4 },
		"no car":            func(w *World) { w.Player.Car = 0 },
		"wrecked car":       func(w *World) { w.Player.CarWear = Wreck - 1 },
		"dry car":           func(w *World) { w.Player.Fuel = 0 },
		"no crew":           func(w *World) { w.Player.Crew = nil },
		"missing driver":    func(w *World) { w.Player.Crew[0].ID = "missing" },
		"busy driver":       func(w *World) { w.Tasks = []Task{{}} },
		"disloyal driver":   func(w *World) { w.Player.Crew[0].Loyalty = HandLoyalty - 1 },
		"dead driver":       func(w *World) { w.NPC("leo").Dead = true },
		"held driver":       func(w *World) { w.NPC("leo").Held = w.Minute + 30 },
		"away driver":       func(w *World) { w.NPC("leo").Location = "bar" },
		"travelling driver": func(w *World) { n := w.NPC("leo"); n.Heading = "bar"; n.Arrives = w.Minute + 30 },
		"dead player":       func(w *World) { w.Player.Alive = false },
		"held player":       func(w *World) { w.Player.HeldUntil = w.Minute + 30 },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			w := buildingDriveByFixture(t, 61)
			change(w)
			before, _ := json.Marshal(w)
			if w.BuildingDriveByReadiness("club") == "" || w.BuildingDriveBy("club") == nil {
				t.Fatal("invalid attack accepted")
			}
			after, _ := json.Marshal(w)
			if string(before) != string(after) {
				t.Fatal("refused attack mutated world")
			}
		})
	}
}

func TestBuildingDriveByBoundsRemainingDamage(t *testing.T) {
	w := buildingDriveByFixture(t, 61)
	w.Properties["club"].Condition, w.Properties["club"].Supply = 2, 0
	w.Player.Heat = 99
	w.Player.Fuel = 1
	if err := w.BuildingDriveBy("club"); err != nil {
		t.Fatal(err)
	}
	if w.Properties["club"].Condition != 0 || w.Properties["club"].Supply != 0 || w.Player.Heat != 100 || w.Fuel() != 0 {
		t.Fatal("unbounded damage, fuel or heat")
	}
	if w.VisualCues[0].DriveBy.ConditionBefore != 2 || w.VisualCues[0].DriveBy.ConditionAfter != 0 {
		t.Fatal("wrong actual damage in scene")
	}
}

func TestBuildingDriveByUsesPlayerCarNotDriversCar(t *testing.T) {
	w := buildingDriveByFixture(t, 61)
	// Hurt and Dry describe the NPC's own car, not personal injury. They must
	// not prevent the crew member driving the player's operational car.
	n := w.NPC("leo")
	n.Car, n.Hurt, n.Dry = 3, true, true
	if err := w.BuildingDriveBy("club"); err != nil {
		t.Fatal(err)
	}
	if w.VisualCues[0].DriveBy.VehicleTier != 1 || !n.Hurt || !n.Dry {
		t.Fatal("borrowed or repaired the driver's own car")
	}
}

func TestBuildingDriveByUnknownAddressIsRejected(t *testing.T) {
	w := buildingDriveByFixture(t, 61)
	before, _ := json.Marshal(w)
	if w.BuildingDriveBy("missing") == nil {
		t.Fatal("unknown address accepted")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("unknown address mutated world")
	}
}

func TestBuildingDriveByProvokesActualOwner(t *testing.T) {
	w := buildingDriveByFixture(t, 61)
	w.Properties["club"].Owner = "russo"
	w.Plots = nil
	w.faction("russo").Goodwill = -90
	other := w.faction("bellandi").Goodwill
	if err := w.BuildingDriveBy("club"); err != nil {
		t.Fatal(err)
	}
	if w.faction("russo").Goodwill != -100 || w.faction("bellandi").Goodwill != other {
		t.Fatal("wrong family or unbounded goodwill")
	}
	if len(w.Plots) != 1 || w.Plots[0].Actor != "russo" || w.Plots[0].Kind != "hit" || w.Plots[0].Due <= w.Minute {
		t.Fatal("missing owner retaliation", w.Plots)
	}
}
