package core

import (
	"encoding/json"
	"testing"
)

func TestCrewGuardTravelsTakesPostAndReturnsOnRecall(t *testing.T) {
	w := ordersFixture(t)
	n := w.NPC("leo")
	n.Faction = w.PlayerOrganizationID()
	if err := w.StartCrewOrder("guard", n.ID, "garage"); err != nil {
		t.Fatal(err)
	}
	if w.PostingDefenceAt("garage") != 0 {
		t.Fatal("remote defence")
	}
	orderStep(w)
	if w.CrewOrders[0].Due-w.Minute != PostingMinutes {
		t.Fatal("wrong setup time")
	}
	orderStep(w)
	if w.CrewOrders[0].Stage != "guarding" || w.PostingDefenceAt("garage") == 0 {
		t.Fatal("guard not established")
	}
	w.SetOut()
	if n.Heading != "" {
		t.Fatal("routine stole guard")
	}
	raw, _ := json.Marshal(w)
	var saved World
	if err := json.Unmarshal(raw, &saved); err != nil {
		t.Fatal(err)
	}
	w = &saved
	if err := w.RecallCrewOrder(w.CrewOrders[0].ID); err != nil {
		t.Fatal(err)
	}
	if w.PostingDefenceAt("garage") != 0 || w.CrewOrders[0].Stage != "returning" {
		t.Fatal("recall did not release post")
	}
	orderStep(w)
	if w.CrewOrders[0].Stage != "done" || w.NPC("leo").Location != "laundry" {
		t.Fatal("guard failed to return")
	}
}
func TestCrewGuardPostPersistsThroughSuccessionAndEndsOnLostDeed(t *testing.T) {
	w := ordersFixture(t)
	n := w.NPC("leo")
	n.Faction = w.PlayerOrganizationID()
	if err := w.StartCrewOrder("guard", n.ID, "garage"); err != nil {
		t.Fatal(err)
	}
	orderStep(w)
	orderStep(w)
	w.Die("Test succession")
	w.Life++
	w.Player = newPerson(w.Life)
	w.SettleCrewOrders()
	if w.CrewOrders[0].Stage != "guarding" || w.PostingDefenceAt("garage") == 0 {
		t.Fatal("inherited guard abandoned post")
	}
	w.Properties["garage"].Owner = "independent"
	w.SettleCrewOrders()
	if w.CrewOrders[0].Stage != "returning" || w.Properties["garage"].Posted != "" {
		t.Fatal("lost deed retained guard")
	}
}
func TestCrewGuardRechecksPostBeforeCommitting(t *testing.T) {
	w := ordersFixture(t)
	n := w.NPC("leo")
	n.Faction = w.PlayerOrganizationID()
	if err := w.StartCrewOrder("guard", n.ID, "garage"); err != nil {
		t.Fatal(err)
	}
	orderStep(w)
	w.Properties["garage"].Owner = "independent"
	orderStep(w)
	if w.CrewOrders[0].Stage != "returning" || w.Properties["garage"].Posted != "" {
		t.Fatal("guard posted after lost deed")
	}
}

func TestCrewGuardCannotRemainOnDutyAfterLeavingFamily(t *testing.T) {
	w := ordersFixture(t)
	n := w.NPC("leo")
	n.Faction = w.PlayerOrganizationID()
	if err := w.StartCrewOrder("guard", n.ID, "garage"); err != nil {
		t.Fatal(err)
	}
	orderStep(w)
	orderStep(w)
	n.Faction = ""
	w.SettleCrewOrders()
	if w.CrewOrders[0].active() || w.Properties["garage"].Posted != "" {
		t.Fatal("departed member retained guard order")
	}
}
