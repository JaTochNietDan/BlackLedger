package core

import (
	"encoding/json"
	"testing"
)

func TestCrewPropertyWorkReservesAndSurvivesReload(t *testing.T) {
	for _, kind := range []string{"repair", "remedy"} {
		t.Run(kind, func(t *testing.T) {
			w := ordersFixture(t)
			p := w.Properties["garage"]
			p.Condition, p.Trouble = 35, true
			cost, _ := w.crewPropertyWork(kind, "garage")
			w.Player.Cash = cost // Work must succeed with the wallet empty after reservation.
			if err := w.StartCrewOrder(kind, "leo", "garage"); err != nil {
				t.Fatal(err)
			}
			if w.Player.Cash != 0 || p.Condition != 35 || !p.Trouble {
				t.Fatal("premature work or wrong reservation")
			}
			orderStep(w)
			if w.CrewOrders[0].Stage != "working" || w.CrewOrders[0].Due-w.Minute != PropertyWorkMinutes {
				t.Fatal("wrong work duration")
			}
			raw, err := json.Marshal(w)
			if err != nil {
				t.Fatal(err)
			}
			var loaded World
			if err := json.Unmarshal(raw, &loaded); err != nil {
				t.Fatal(err)
			}
			w = &loaded
			orderStep(w)
			p = w.Properties["garage"]
			if kind == "repair" && (p.Condition != 75 || !p.Trouble) {
				t.Fatal("wrong repair effect")
			}
			if kind == "remedy" && (p.Trouble || p.Condition != 35) {
				t.Fatal("wrong remedy effect")
			}
			orderStep(w)
			w.SettleCrewOrders()
			if w.Player.Cash != 0 || w.CrewOrders[0].Reserved != 0 || w.CrewOrders[0].Stage != "done" {
				t.Fatal("double payment or unfinished work")
			}
		})
	}
}

func TestCrewPropertyWorkInvalidationReturnsBudget(t *testing.T) {
	for _, kind := range []string{"repair", "remedy"} {
		for _, change := range []string{"deed-outbound", "deed-working", "already-done", "recall"} {
			t.Run(kind+"/"+change, func(t *testing.T) {
				w := ordersFixture(t)
				p := w.Properties["garage"]
				p.Condition, p.Trouble = 35, true
				cash := w.Player.Cash
				if err := w.StartCrewOrder(kind, "leo", "garage"); err != nil {
					t.Fatal(err)
				}
				if change != "deed-outbound" {
					orderStep(w)
				}
				switch change {
				case "deed-outbound", "deed-working":
					p.Owner = "bellandi"
				case "already-done":
					p.Condition, p.Trouble = 100, false
				case "recall":
					if err := w.RecallCrewOrder(w.CrewOrders[0].ID); err != nil {
						t.Fatal(err)
					}
				}
				for i := 0; i < 4 && w.CrewOrders[0].active(); i++ {
					orderStep(w)
				}
				if w.CrewOrders[0].active() || w.Player.Cash != cash {
					t.Fatal("budget not returned")
				}
				if change != "already-done" && (p.Condition != 35 || !p.Trouble) {
					t.Fatal("invalidated job changed property")
				}
			})
		}
	}
}
