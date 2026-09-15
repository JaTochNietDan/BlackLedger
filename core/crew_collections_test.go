package core

import "testing"

func TestNamedCollectionRoundTravelsAndPaysOnlyOnReturn(t *testing.T) {
	w := ordersFixture(t)
	before := w.Player.Cash
	if err := w.StartCrewOrder("collections", "leo", "garage"); err != nil {
		t.Fatal(err)
	}
	if len(w.Tasks) != 0 || w.Player.Cash != before {
		t.Fatal("legacy task or early payment")
	}
	if w.DelegateReadiness() == "" || w.crewCollectionTarget("garage", "") == "" {
		t.Fatal("duplicate round allowed")
	}
	orderStep(w)
	if w.CrewOrders[0].Due != w.Minute+CollectionMinutes {
		t.Fatal("wrong work time")
	}
	w = w.Clone()
	orderStep(w)
	if w.CrewOrders[0].Loot != CollectionPay || w.Player.Cash != before || w.Player.Location != "laundry" {
		t.Fatal("collection custody or player location")
	}
	w = w.Clone()
	orderStep(w)
	w.SettleCrewOrders()
	if w.Player.Cash != before+CollectionPay || w.CrewOrders[0].Loot != 0 || w.CrewOrders[0].Stage != "done" {
		t.Fatal("incorrect return settlement")
	}
}

func TestNamedCollectionInvalidationAndCarrierLoss(t *testing.T) {
	for _, cause := range []string{"recall", "sold", "ruined", "trouble", "captured", "dead"} {
		t.Run(cause, func(t *testing.T) {
			w := ordersFixture(t)
			before := w.Player.Cash
			if err := w.StartCrewOrder("collections", "leo", "garage"); err != nil {
				t.Fatal(err)
			}
			orderStep(w)
			switch cause {
			case "recall":
				if err := w.RecallCrewOrder(w.CrewOrders[0].ID); err != nil {
					t.Fatal(err)
				}
			case "sold":
				w.Properties["garage"].Owner = "bellandi"
			case "ruined":
				w.Properties["garage"].Condition = 0
			case "trouble":
				w.Properties["garage"].Trouble = true
			case "captured", "dead":
				orderStep(w) // Work committed, money still carried on the return leg.
				if cause == "captured" {
					w.NPC("leo").Held = w.Minute + 1000
				} else {
					w.NPC("leo").Dead = true
				}
			}
			for i := 0; i < 4 && w.CrewOrders[0].active(); i++ {
				orderStep(w)
			}
			w.SettleCrewOrders()
			if w.Player.Cash != before || w.CrewOrders[0].active() || w.CrewOrders[0].Loot != 0 {
				t.Fatal("invalid round paid", w.CrewOrders[0])
			}
		})
	}
}

func TestCollectionOffersAndLegacyExclusion(t *testing.T) {
	w := ordersFixture(t)
	found := false
	for _, offer := range w.CrewOrderOffers() {
		if offer.Kind == "collections" && offer.Actor == "leo" && offer.Target == "garage" {
			found = true
			if offer.Reason != "" || offer.Minutes != 2*TravelMinutes("laundry", "garage")+CollectionMinutes || offer.Cost != 0 {
				t.Fatal("wrong quote", offer)
			}
		}
	}
	if !found {
		t.Fatal("missing offer")
	}
	w.SendOnCollections()
	if w.crewCollectionTarget("garage", "") == "" {
		t.Fatal("legacy round overlaps")
	}
}

func TestCollectionCommandUsesNamedOperative(t *testing.T) {
	w := ordersFixture(t)
	mara := w.NPC("mara")
	mara.Faction, mara.Trust, mara.Location = w.PlayerOrganizationID(), 80, "laundry"
	mara.Heading, mara.Arrives, mara.Held = "", 0, 0
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "crew_order:collections", Choice: mara.ID, Target: "garage"})
	if err != nil {
		t.Fatal(err)
	}
	if len(w.CrewOrders) != 0 || len(next.CrewOrders) != 1 || next.CrewOrders[0].Actor != mara.ID || next.Player.Location != "laundry" {
		t.Fatal("command used wrong actor or mutated source")
	}
}
