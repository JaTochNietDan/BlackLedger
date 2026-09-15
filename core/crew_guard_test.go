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
	w.Player.Crew = nil
	w.SettleCrewOrders()
	if w.CrewOrders[0].active() || w.Properties["garage"].Posted != "" {
		t.Fatal("departed member retained guard order")
	}
}

func TestHiredAssociateCanGuardWithoutJoiningFamily(t *testing.T) {
	w := ordersFixture(t)
	n := w.NPC("leo")
	n.Faction = ""
	n.Trust = 0
	if err := w.StartCrewOrder("guard", n.ID, "garage"); err != nil {
		t.Fatal(err)
	}
	orderStep(w)
	orderStep(w)
	if w.PostedAt("garage") != n || w.PostingDefenceAt("garage") <= PostingDefence {
		t.Fatal("associate did not defend business")
	}
	if got := w.PostingDescription("garage")["trust"]; got != 80 {
		t.Fatal("hired guard ignored loyalty", got)
	}
	if n.Faction != "" {
		t.Fatal("guard order forced family membership")
	}
	if w.CrewOrderReadiness("restock", n.ID, "laundry") == "" {
		t.Fatal("guard double booked")
	}
	w.Die("Test employee contract ends")
	w.SettleCrewOrders()
	if w.CrewOrders[0].active() || w.PostedAt("garage") != nil {
		t.Fatal("personal employment survived employer death")
	}
}

func TestHeadquartersGuardDoesNotWalkOutOfTheSameBuilding(t *testing.T) {
	w := ordersFixture(t)
	n := w.NPC("leo")
	if err := w.StartCrewOrder("guard", n.ID, "laundry"); err != nil {
		t.Fatal(err)
	}
	if w.Travelling(n) || n.Heading != "" {
		t.Fatal("same-address dispatch created a street journey")
	}
	orderStep(w)
	orderStep(w)
	if err := w.RecallCrewOrder(w.CrewOrders[0].ID); err != nil {
		t.Fatal(err)
	}
	if w.Travelling(n) || n.Heading != "" || n.Location != "laundry" {
		t.Fatal("same-address report-back left the building")
	}
	orderStep(w)
	if w.CrewOrders[0].Stage != "done" {
		t.Fatal("same-address report not settled")
	}
}

func TestLegacyGuardPostBlocksAllDelegation(t *testing.T) {
	w := ordersFixture(t)
	n := w.NPC("leo")
	n.Faction = w.PlayerOrganizationID()
	postAndArrive(t, w, "garage")
	hand, _ := w.NamedHands(n.ID)
	if w.HandReadiness(hand) == "" || w.DelegateReadiness() == "" {
		t.Fatal("legacy guard sent on a second job")
	}
	if err := w.Unpost("garage"); err != nil {
		t.Fatal(err)
	}
	if w.HandReadiness(hand) != "" {
		t.Fatal("relieved guard stayed unavailable", w.HandReadiness(hand))
	}
}
func TestLegacyPostingSkipsTravellersAndNamedErrands(t *testing.T) {
	for _, reason := range []string{"travel", "named-task"} {
		t.Run(reason, func(t *testing.T) {
			w := ordersFixture(t)
			n := w.NPC("leo")
			n.Faction = w.PlayerOrganizationID()
			if reason == "travel" {
				n.Heading = "garage"
				n.Sets = 0
				n.Arrives = w.Minute + 30
			} else {
				w.Homes = []TaskHome{{Task: "errand", Person: n.ID, Where: "laundry"}}
			}
			for _, free := range w.Unposted() {
				if free.ID == n.ID {
					t.Fatal("busy operative offered for guard duty")
				}
			}
		})
	}
}

func TestLegacyGuardCancelsAnUnstartedRoutine(t *testing.T) {
	w := ordersFixture(t)
	n := w.NPC("leo")
	n.Faction = w.PlayerOrganizationID()
	n.Heading, n.Sets, n.Arrives = "bar", w.Minute+100, w.Minute+140
	if err := w.Post("laundry"); err != nil {
		t.Fatal(err)
	}
	if n.Heading != "" || n.Sets != 0 || n.Arrives != 0 {
		t.Fatal("guard kept a pending departure")
	}
}
