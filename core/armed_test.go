package core

import "testing"

func TestArmedResistanceDoesNotRiskAnAbsentPlayer(t *testing.T) {
	w := New(31)
	before := w.RNG
	if w.armedResistance(1, 100, Hand{Crew: true}) || w.RNG != before {
		t.Fatal("delegated work exposed the player or consumed their risk roll")
	}
	w.Player.Armour = 3
	if w.armedResistance(TillGun, 0, w.OwnHands()) || w.RNG != before {
		t.Fatal("armour did not remove the unbacked shop's armed risk")
	}
}

func TestFailedRobberiesCanEndAHealthyPlayersLife(t *testing.T) {
	for _, kind := range []string{"business", "street"} {
		t.Run(kind, func(t *testing.T) {
			for seed := uint32(1); seed <= 1000; seed++ {
				var w *World
				if kind == "business" {
					w = robberWorld(t)
				} else {
					w, _ = mugger(t)
				}
				w.RNG = seed * 2654435761
				w.Player.Health = 100
				cash := w.Player.Cash
				var err error
				if kind == "business" {
					err = w.Rob("club")
				} else {
					err = w.Mug("club", w.OwnHands())
				}
				if err != nil {
					t.Fatal(err)
				}
				if !w.Player.Alive {
					if w.Player.Cash != cash || len(w.Dead) != 1 {
						t.Fatal("fatal failed robbery paid out or recorded death incorrectly")
					}
					return
				}
			}
			t.Fatal("no healthy player faced lethal resistance across the seeded attempts")
		})
	}
}
