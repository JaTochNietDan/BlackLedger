package core

import "testing"

func TestCrewBudgetLostWithCarrierAcrossReload(t *testing.T) {
	for _, stage := range []string{"outbound", "working", "returning"} {
		for _, disposition := range []string{"dead", "captured", "missing"} {
			t.Run(stage+"/"+disposition, func(t *testing.T) {
				w := ordersFixture(t)
				before := w.Player.Cash
				cost := w.RestockCost("garage")
				if err := w.StartCrewOrder("restock", "leo", "garage"); err != nil {
					t.Fatal(err)
				}
				if stage != "outbound" {
					orderStep(w)
				}
				if stage == "returning" {
					// Recall leaves the unused budget with the returning carrier.
					if err := w.RecallCrewOrder(w.CrewOrders[0].ID); err != nil {
						t.Fatal(err)
					}
				}
				actor := w.NPC("leo")
				purse := actor.Purse
				switch disposition {
				case "dead":
					actor.Dead = true
				case "captured":
					actor.Held = w.Minute + 1000
				case "missing":
					actor.ID = "removed-operative"
				}
				w = w.Clone()
				w.SettleCrewOrders()
				w = w.Clone()
				w.SettleCrewOrders()
				o := w.CrewOrders[0]
				if o.Stage != "cancelled" || o.Reserved != 0 || w.Player.Cash != before-cost || w.Properties["garage"].Supply != 0 {
					t.Fatal("lost carrier refunded or completed work", o, w.Player.Cash)
				}
				if n := w.NPC("leo"); n != nil && n.Purse != purse {
					t.Fatal("funds transferred to unavailable carrier")
				}
			})
		}
	}
}
