package core

import "testing"

func TestCrewBombReservesChargeAndRecallReturnsItOnce(t *testing.T) {
	w := ordersFixture(t)
	w.Player.Charges = 1
	w.NPC("leo").Skill = 90
	if err := w.StartCrewOrder("bomb", "leo", "club"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Charges != 0 || w.CrewOrders[0].Charges != 1 || w.CrewOrders[0].Reserved != 0 {
		t.Fatal("charge was not reserved exactly")
	}
	w = w.Clone() // Reserved equipment survives the same JSON path as a reload.
	if err := w.RecallCrewOrder(w.CrewOrders[0].ID); err != nil {
		t.Fatal(err)
	}
	for w.CrewOrders[0].active() {
		orderStep(w)
	}
	w.SettleCrewOrders()
	if w.Player.Charges != 1 || w.CrewOrders[0].Charges != 0 {
		t.Fatal("unused charge not returned exactly once")
	}
}
func TestCrewBombRechecksDeedAndDoesNotBombNewFamilyProperty(t *testing.T) {
	w := ordersFixture(t)
	w.Player.Charges = 1
	w.NPC("leo").Skill = 90
	if err := w.StartCrewOrder("bomb", "leo", "club"); err != nil {
		t.Fatal(err)
	}
	orderStep(w)
	condition := w.Properties["club"].Condition
	w.Properties["club"].Owner = w.PlayerOrganizationID()
	for w.CrewOrders[0].active() {
		orderStep(w)
	}
	if w.Properties["club"].Condition != condition || w.Player.Charges != 1 {
		t.Fatal("bombed newly owned property or consumed unused charge")
	}
}
func TestCrewBombBothOutcomesConsumeOnlyTheReservedCharge(t *testing.T) {
	blasts, accidents := 0, 0
	for seed := uint32(1); seed <= 50; seed++ {
		w := ordersFixture(t)
		w.RNG = scatter(seed)
		w.Player.Charges = 2
		w.NPC("leo").Skill = 90
		if err := w.StartCrewOrder("bomb", "leo", "club"); err != nil {
			t.Fatal(err)
		}
		orderStep(w)
		if w.CrewOrders[0].Due-w.Minute != PlantMinutes {
			t.Fatal("incorrect planting duration")
		}
		orderStep(w)
		if w.Player.Health != 100 || w.Player.Location != "laundry" || w.Player.Charges != 1 || w.CrewOrders[0].Charges != 0 {
			t.Fatal("remote bombing charged player body or extra charge")
		}
		found := false
		for _, cue := range w.VisualCues {
			if cue.Kind == "explosion" && cue.Attacker != nil && cue.Attacker.ID == "leo" {
				found = true
				if cue.Detonation == "premature" {
					accidents++
				} else {
					blasts++
					if w.NPC("leo").Dead {
						t.Fatal("successful planter was selected as collateral casualty")
					}
				}
			}
		}
		if !found {
			t.Fatal("bomb cue lacks named operative")
		}
	}
	if blasts == 0 || accidents == 0 {
		t.Fatal("both outcomes not exercised", blasts, accidents)
	}
}
func TestCapturedOperativeDoesNotReturnExplosives(t *testing.T) {
	w := ordersFixture(t)
	w.Player.Charges = 1
	w.NPC("leo").Skill = 90
	if err := w.StartCrewOrder("bomb", "leo", "club"); err != nil {
		t.Fatal(err)
	}
	w.hold(w.NPC("leo"), 2)
	w.SettleCrewOrders()
	if w.Player.Charges != 0 || w.CrewOrders[0].Charges != 0 || w.CrewOrders[0].active() {
		t.Fatal("captured explosives returned or job survived custody")
	}
}
