package core

import "testing"

func TestPublicJourneyVehicleReflectsOperationalCar(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name      string
		car       int
		dry, hurt bool
		want      string
	}{
		{"pedestrian", 0, false, false, ""}, {"Ford", 1, false, false, VehicleByTier(1).Label},
		{"Hudson", 2, false, false, VehicleByTier(2).Label}, {"Packard", 3, false, false, VehicleByTier(3).Label},
		{"dry", 2, true, false, ""}, {"damaged", 3, false, true, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := New(4)
			n := w.NPC("mara")
			n.Location = "bar"
			n.Heading = "laundry"
			n.Sets = 0
			n.Arrives = w.Minute + 5
			n.Car = tt.car
			n.Dry = tt.dry
			n.Hurt = tt.hurt
			minute, revision := w.Minute, w.Revision
			for _, j := range w.OnTheStreet() {
				if j.ID == n.ID {
					if j.Vehicle != tt.want {
						t.Fatalf("vehicle %q, want %q", j.Vehicle, tt.want)
					}
					if w.Minute != minute || w.Revision != revision {
						t.Fatal("reading presentation changed simulation")
					}
					return
				}
			}
			t.Fatal("journey missing")
		})
	}
}

func TestFailedChargeStillPublishesExplosionCue(t *testing.T) {
	t.Parallel()
	for seed := uint32(1); seed < 100; seed++ {
		w := New(seed)
		w.Player.Location = "club"
		w.Player.Charges = 1
		w.Player.Respect = 60
		if err := w.Plant("club"); err != nil {
			t.Fatal(err)
		}
		if w.Player.Health == 100 {
			continue
		}
		count := 0
		for _, cue := range w.VisualCues {
			if cue.Kind == "explosion" && cue.Target == "club" {
				count++
			}
		}
		if count != 1 {
			t.Fatalf("failed charge produced %d explosions", count)
		}
		return
	}
	t.Fatal("fixture did not exercise failed charge")
}

func TestBlastCannotKillSomebodyAtAnotherAddressOrOnTheStreet(t *testing.T) {
	t.Parallel()
	sawCasualty := false
	for seed := uint32(1); seed <= 100; seed++ {
		w := New(seed)
		for i := range w.NPCs {
			w.NPCs[i].Location = "room"
			w.NPCs[i].Heading = ""
			w.NPCs[i].Sets = 0
		}
		present := w.NPC("mara")
		present.Location = "club"
		travelling := w.NPC("leo")
		travelling.Location = "club"
		travelling.Heading = "laundry"
		travelling.Arrives = w.Minute + 20
		w.WorldRNG = seed * 2654435761
		w.detonate("club", "An isolated blast.")
		for _, n := range w.NPCs {
			if n.Dead {
				if n.ID != present.ID {
					t.Fatalf("blast killed %s at %s, heading %s", n.ID, n.Location, n.Heading)
				}
				sawCasualty = true
			}
		}
		for _, cue := range w.VisualCues {
			if cue.Kind == "killing" && cue.Target != "club" {
				t.Fatalf("blast death was staged at %s", cue.Target)
			}
		}
	}
	if !sawCasualty {
		t.Fatal("fixture never exercised casualty selection")
	}
}
