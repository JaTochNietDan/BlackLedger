package core

import (
	"encoding/json"
	"testing"
)

func ordersFixture(t *testing.T) *World {
	t.Helper()
	w := withCrew(t, 80)
	w.Event = nil
	w.Plots = nil
	w.Tasks = nil
	w.Homes = nil
	w.Player.Location = "laundry"
	w.Properties["laundry"].Owner = w.PlayerOrganizationID()
	w.Incorporate()
	w.PlayerOrganization().Headquarters = "laundry"
	w.Properties["garage"].Owner = w.PlayerOrganizationID()
	w.Properties["garage"].Supply = 0
	n := w.NPC("leo")
	n.Location = "laundry"
	n.Heading = ""
	n.Arrives = 0
	n.Held = 0
	return w
}
func orderStep(w *World) { w.Minute = w.CrewOrders[0].Due; w.Arrivals(); w.SettleCrewOrders() }
func TestCrewOrderReservesOnceTravelsWorksAndReturnsAcrossReload(t *testing.T) {
	w := ordersFixture(t)
	cash := w.Player.Cash
	cost := w.RestockCost("garage")
	if err := w.StartCrewOrder("restock", "leo", "garage"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-cost || w.NPC("leo").Location != "laundry" || w.Player.Location != "laundry" {
		t.Fatal("dispatch did not reserve exactly or teleported")
	}
	if w.CrewOrderReadiness("restock", "leo", "garage") == "" {
		t.Fatal("double booked")
	}
	raw, _ := json.Marshal(w)
	var loaded World
	if err := json.Unmarshal(raw, &loaded); err != nil {
		t.Fatal(err)
	}
	w = &loaded
	orderStep(w)
	if w.CrewOrders[0].Stage != "working" || w.NPC("leo").Location != "garage" {
		t.Fatal("missing arrival")
	}
	w.SetOut()
	if w.NPC("leo").Heading != "" {
		t.Fatal("routine stole working operative")
	}
	orderStep(w)
	trade, _ := TradeOf("garage")
	if w.Properties["garage"].Supply != trade.RestockAmount || w.CrewOrders[0].Stage != "returning" {
		t.Fatal("did not restock")
	}
	orderStep(w)
	w.SettleCrewOrders()
	if w.CrewOrders[0].Stage != "done" || w.Player.Cash != cash-cost || w.NPC("leo").Location != "laundry" {
		t.Fatal("return or exactly-once settlement failed")
	}
}
func TestCrewOrderRecallAndChangedDeedRefundUnusedBudget(t *testing.T) {
	for _, recall := range []bool{true, false} {
		w := ordersFixture(t)
		cash := w.Player.Cash
		if err := w.StartCrewOrder("restock", "leo", "garage"); err != nil {
			t.Fatal(err)
		}
		if recall {
			if err := w.RecallCrewOrder(w.CrewOrders[0].ID); err != nil {
				t.Fatal(err)
			}
			if w.NPC("leo").Heading != "garage" {
				t.Fatal("recall teleported")
			}
		} else {
			w.Properties["garage"].Owner = "bellandi"
		}
		for w.CrewOrders[0].active() {
			orderStep(w)
		}
		if w.Player.Cash != cash || w.Properties["garage"].Supply != 0 {
			t.Fatal("unused funds lost or another owner’s business restocked")
		}
	}
}
func TestCrewOrderRechecksVictimAndKeepsPlayerAtHeadquarters(t *testing.T) {
	w := ordersFixture(t)
	victim := w.NPC("mara")
	victim.Faction = ""
	victim.Location = "garage"
	victim.Heading = ""
	victim.Arrives = 0
	w.MeetPerson(victim.ID)
	w.See(victim)
	if err := w.StartCrewOrder("assassinate", "leo", victim.ID); err != nil {
		t.Fatal(err)
	}
	orderStep(w)
	victim.Location = "bar"
	orderStep(w)
	if victim.Dead || w.Player.Location != "laundry" || w.Player.Health != 100 {
		t.Fatal("moved victim attacked or player used as proxy")
	}
}
func TestCrewOrderBlocksGuardAndNewLifeInheritsNoBudget(t *testing.T) {
	w := ordersFixture(t)
	w.Properties["garage"].Posted = "leo"
	if w.CrewOrderReadiness("restock", "leo", "garage") == "" {
		t.Fatal("guard double booked")
	}
	w.Properties["garage"].Posted = ""
	if err := w.StartCrewOrder("restock", "leo", "garage"); err != nil {
		t.Fatal(err)
	}
	w.Life++
	cash := w.Player.Cash
	w.SettleCrewOrders()
	if w.Player.Cash != cash || w.CrewOrders[0].active() {
		t.Fatal("new life received former budget or order")
	}
	if w.NPC("leo").Heading != "garage" {
		t.Fatal("cancellation teleported a living operative")
	}
}

func TestCrewOrdersRunConcurrentlyAndExecuteDoesNotMutateRejectedState(t *testing.T) {
	w := ordersFixture(t)
	mara := w.NPC("mara")
	mara.Faction = w.PlayerOrganizationID()
	mara.Trust = 80
	mara.Location = "laundry"
	mara.Heading = ""
	mara.Arrives = 0
	w.Properties["laundry"].Supply = 0
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "crew_order:restock", Choice: "leo", Target: "garage"})
	if err != nil {
		t.Fatal(err)
	}
	if len(w.CrewOrders) != 0 {
		t.Fatal("original mutated")
	}
	if _, err = Execute(next, Command{Revision: next.Revision, Kind: "crew_order:restock", Choice: "leo", Target: "garage"}); err == nil {
		t.Fatal("duplicate assignment accepted")
	}
	next.Event = nil
	if err = next.StartCrewOrder("restock", "mara", "laundry"); err != nil {
		t.Fatal(err)
	}
	if len(next.CrewOrders) != 2 || next.CrewOrderFor("leo") == nil || next.CrewOrderFor("mara") == nil {
		t.Fatal("missing concurrent assignments")
	}
	for i := 0; i < 20; i++ {
		next.Event = nil
		next.Advance(10)
	}
	for _, o := range next.CrewOrders {
		if o.active() {
			t.Fatal("real advance clock failed to settle order", o)
		}
	}
	if next.Player.Location != "laundry" {
		t.Fatal("player moved")
	}
}

func TestCrewOrderAssassinationUsesRemoteActorAndInteriorCue(t *testing.T) {
	killed := 0
	for seed := uint32(1); seed <= 30; seed++ {
		w := ordersFixture(t)
		w.RNG = scatter(seed)
		victim := w.NPC("mara")
		victim.Faction = ""
		victim.Location = "garage"
		victim.Heading = ""
		victim.Arrives = 0
		w.MeetPerson(victim.ID)
		w.See(victim)
		w.NPC("leo").Weapon = 1
		if err := w.StartCrewOrder("assassinate", "leo", victim.ID); err != nil {
			t.Fatal(err)
		}
		orderStep(w)
		orderStep(w)
		if w.Player.Location != "laundry" || w.Player.Health != 100 {
			t.Fatal("remote order moved or wounded player")
		}
		if w.NPC(victim.ID).Dead {
			killed++
			found := false
			for _, cue := range w.VisualCues {
				if cue.Attacker != nil && cue.Attacker.ID == "leo" && cue.Strike != nil {
					found = true
				}
			}
			if !found {
				t.Fatal("successful remote strike has no actual operative cue")
			}
		}
	}
	if killed == 0 {
		t.Fatal("never exercised successful strike")
	}
}

func TestCrewOrderCannotTrackAnUnseenTarget(t *testing.T) {
	w := ordersFixture(t)
	n := w.NPC("mara")
	n.Faction = ""
	n.Location = "garage"
	n.Heading = ""
	n.Arrives = 0
	w.MeetPerson(n.ID)
	delete(w.Sightings, n.ID)
	if w.CrewOrderReadiness("assassinate", "leo", n.ID) == "" {
		t.Fatal("known name disclosed hidden address")
	}
	w.See(n)
	n.Location = "casino"
	if err := w.StartCrewOrder("assassinate", "leo", n.ID); err != nil {
		t.Fatal(err)
	}
	if w.CrewOrders[0].Place != "garage" {
		t.Fatal("order tracked hidden movement")
	}
	orderStep(w)
	if w.CrewOrders[0].Stage != "returning" || n.Dead {
		t.Fatal("did not abandon stale address")
	}
}
