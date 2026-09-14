package core

import (
	"encoding/json"
	"testing"
)

func TestFamilyRecallsInheritedOrdersAndReceivesTheirBudget(t *testing.T) {
	for _, stage := range []string{"outbound", "working", "returning"} {
		t.Run(stage, func(t *testing.T) {
			w := ordersFixture(t)
			actor := w.NPC("leo")
			actor.Faction = w.PlayerOrganizationID()
			if err := w.StartCrewOrder("restock", actor.ID, "garage"); err != nil {
				t.Fatal(err)
			}
			cost := w.CrewOrders[0].Reserved
			if stage != "outbound" {
				orderStep(w)
			}
			if stage == "returning" {
				orderStep(w)
			}
			w.Die("Test succession")
			estate := "estate:" + actor.ID
			f := w.faction(estate)
			if f == nil || w.CrewOrders[0].Estate != estate {
				t.Fatal("order did not pass to family")
			}
			funds := f.Cash
			w.Dissolve(w.PlayerOrganizationID())
			w.Life++
			w.Player = newPerson(w.Life)
			cash := w.Player.Cash
			raw, _ := json.Marshal(w)
			var loaded World
			if err := json.Unmarshal(raw, &loaded); err != nil {
				t.Fatal(err)
			}
			w = &loaded
			for i := 0; i < 4 && w.CrewOrders[0].active(); i++ {
				orderStep(w)
			}
			w.SettleCrewOrders()
			expected := funds + cost
			if stage == "returning" {
				expected = funds
			}
			if w.CrewOrders[0].Stage != "done" || w.faction(estate).Cash != expected || w.Player.Cash != cash {
				t.Fatal("wrong successor settlement", w.CrewOrders[0], w.faction(estate).Cash, expected)
			}
			if stage != "returning" && w.Properties["garage"].Supply != 0 {
				t.Fatal("recalled work was executed")
			}
			if w.NPC("leo").Location != "laundry" {
				t.Fatal("operative did not return to inherited base")
			}
		})
	}
}

func TestInheritedProceedsStayWithFamilyUnlessOperativeIsLost(t *testing.T) {
	for _, disposition := range []string{"returns", "captured", "dead", "family-gone"} {
		t.Run(disposition, func(t *testing.T) {
			w := ordersFixture(t)
			actor := w.NPC("leo")
			actor.Faction = w.PlayerOrganizationID()
			// A committed robbery is already on its return leg when the leader dies.
			w.CrewOrders = []CrewOrder{{ID: ID(), Life: w.Life, Actor: actor.ID, Name: actor.Name, Kind: "rob", Place: "garage", Base: "laundry", Stage: "returning", Loot: 125}}
			actor.Location = "garage"
			w.crewOrderJourney(&w.CrewOrders[0], "laundry")
			w.Die("Test succession")
			estate := "estate:" + actor.ID
			f := w.faction(estate)
			before, purse := f.Cash, actor.Purse
			switch disposition {
			case "captured":
				actor.Held = w.Minute + 1000
			case "dead":
				actor.Dead = true
			case "family-gone":
				w.Dissolve(estate)
			}
			for i := 0; i < 3 && w.CrewOrders[0].active(); i++ {
				orderStep(w)
			}
			w.SettleCrewOrders()
			if w.CrewOrders[0].Loot != 0 {
				t.Fatal("unsettled proceeds")
			}
			if disposition == "returns" && f.Cash != before+125 {
				t.Fatal("family did not receive proceeds")
			}
			if (disposition == "captured" || disposition == "dead") && f.Cash != before {
				t.Fatal("lost proceeds were recovered")
			}
			if disposition == "family-gone" && actor.Purse != purse+125 {
				t.Fatal("survivor lost funds after family dissolution")
			}
		})
	}
}

func TestSuccessionRecallsAnUncommittedAttack(t *testing.T) {
	w := ordersFixture(t)
	actor := w.NPC("leo")
	actor.Faction = w.PlayerOrganizationID()
	victim := w.NPC("mara")
	victim.Faction, victim.Location, victim.Heading, victim.Arrives = "", "garage", "", 0
	w.MeetPerson(victim.ID)
	w.See(victim)
	if err := w.StartCrewOrder("assassinate", actor.ID, victim.ID); err != nil {
		t.Fatal(err)
	}
	orderStep(w)
	w.Die("Test succession")
	w.Life++
	w.Player = newPerson(w.Life)
	for i := 0; i < 4 && w.CrewOrders[0].active(); i++ {
		orderStep(w)
	}
	if victim.Dead || w.CrewOrders[0].Stage != "done" || w.CrewOrders[0].Estate == "" {
		t.Fatal("attack survived leadership change")
	}
}
