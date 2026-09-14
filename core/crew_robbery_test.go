package core

import "testing"

func successfulCrewRobbery(t *testing.T) *World {
	t.Helper()
	for seed := uint32(1); seed < 100; seed++ {
		w := ordersFixture(t)
		w.RNG = scatter(seed)
		w.Player.Health = 20
		f := w.faction(w.Properties["club"].Owner)
		f.Cash = 37
		if err := w.StartCrewOrder("rob", "leo", "club"); err != nil {
			t.Fatal(err)
		}
		cash := w.Player.Cash
		orderStep(w)
		orderStep(w)
		if w.Player.Cash != cash || w.Player.Health != 20 || w.Player.Location != "laundry" {
			t.Fatal("remote robbery paid early or used player's body")
		}
		if w.CrewOrders[0].Loot > 0 {
			cue := w.VisualCues[len(w.VisualCues)-1]
			if cue.Kind != "robbery" || cue.Attacker == nil || cue.Attacker.ID != "leo" {
				t.Fatal("robbery scene missing its actual operative")
			}

			if w.CrewOrders[0].Loot != 37 || f.Cash != 0 {
				t.Fatal("robbery created money or debited twice")
			}
			return w
		}
	}
	t.Fatal("no successful robbery")
	return nil
}
func TestCrewRobberyCarriesFundedLootHomeAcrossReload(t *testing.T) {
	w := successfulCrewRobbery(t)
	w = w.Clone()
	cash := w.Player.Cash
	orderStep(w)
	w.SettleCrewOrders()
	if w.Player.Cash != cash+37 || w.CrewOrders[0].Loot != 0 || w.CrewOrders[0].Stage != "done" {
		t.Fatal("loot was lost or paid twice on return")
	}
}
func TestCrewRobberyCapturedLootNeverReachesPlayer(t *testing.T) {
	w := successfulCrewRobbery(t)
	cash := w.Player.Cash
	w.hold(w.NPC("leo"), 2)
	w.SettleCrewOrders()
	if w.Player.Cash != cash || w.CrewOrders[0].Loot != 0 || w.CrewOrders[0].active() {
		t.Fatal("captured loot reached player or stayed spendable")
	}
}
func TestCrewRobberyRechecksNewOwnership(t *testing.T) {
	w := ordersFixture(t)
	if err := w.StartCrewOrder("rob", "leo", "club"); err != nil {
		t.Fatal(err)
	}
	orderStep(w)
	w.Properties["club"].Owner = w.PlayerOrganizationID()
	condition := w.Properties["club"].Condition
	orderStep(w)
	if w.CrewOrders[0].Loot != 0 || w.Properties["club"].Condition != condition {
		t.Fatal("robbed the family's newly acquired business")
	}
}

func TestCrewRobberyDebitsPersonalProprietor(t *testing.T) {
	for seed := uint32(1); seed < 100; seed++ {
		w := ordersFixture(t)
		w.RNG = scatter(seed)
		w.Properties["club"].Owner = "mara"
		w.Properties["club"].Bankroll = 37
		before := w.businessFunds("club")
		if before <= 0 {
			t.Fatal("proprietor fixture has no takings")
		}
		if err := w.StartCrewOrder("rob", "leo", "club"); err != nil {
			t.Fatal(err)
		}
		orderStep(w)
		orderStep(w)
		if loot := w.CrewOrders[0].Loot; loot > 0 {
			if w.businessFunds("club") != before-loot || w.NPC("mara").Sore == 0 {
				t.Fatal("personal owner did not lose funds and remember robbery")
			}
			return
		}
	}
	t.Fatal("no successful personal-business robbery")
}
