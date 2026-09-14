package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func proprietorWorld() *World {
	w := New(7)
	w.NPCs = []NPC{{ID: "buyer", Name: "Rosa Finch", Location: "bar", Home: "riverside", Purse: 3000}}
	for id, p := range w.Properties {
		if id != "laundry" {
			p.Owner = "reserved"
		}
	}
	w.Properties["laundry"].Owner = "independent"
	return w
}
func TestProprietorBuysWithFundsAndCannotBeTakenAsUnheld(t *testing.T) {
	w := proprietorWorld()
	n := w.NPC("buyer")
	before := w.HouseholdWealth(n)
	price := AcquisitionCost(w, "laundry")
	w.ConsiderProprietors()
	if w.Properties["laundry"].Owner != n.ID || w.HouseholdWealth(n) != before-price {
		t.Fatal("purchase not funded")
	}
	if w.CanAcquire("laundry") || w.Unheld("laundry") || !strings.Contains(w.HolderName("laundry"), n.Name) {
		t.Fatal("private ownership not respected")
	}
	money := n.Purse
	w.ProprietorDay()
	if n.Purse != money {
		t.Fatal("purchase granted retroactive income")
	}
	w.Minute += 1440
	p := w.Properties["laundry"]
	trade, _ := TradeOf("laundry")
	expected := int(float64(p.Income*24*p.Condition)/100*w.Capacity("laundry")*w.TradeMultiplier("laundry")) - p.Staff*trade.Wage
	w.ProprietorDay()
	if n.Purse != money+expected {
		t.Fatalf("income/wages: got %d want %d", n.Purse, money+expected)
	}
	serialized, _ := json.Marshal(w)
	w.ProprietorDay()
	again, _ := json.Marshal(w)
	if string(serialized) != string(again) {
		t.Fatal("daily accounts settled twice")
	}
}
func TestProprietorDeathReleasesOnlyPersonalDeeds(t *testing.T) {
	w := New(7)
	n := w.NPC(w.Members(w.Factions[0].ID)[0].ID)
	faction := n.Faction
	w.Properties["laundry"].Owner = n.ID
	w.Properties["garage"].Owner = faction
	staff, condition := w.Properties["laundry"].Staff, w.Properties["laundry"].Condition
	if !w.Kill(n.ID, "Test fixture death") {
		t.Fatal("not killed")
	}
	if !w.CanAcquire("laundry") || w.Properties["garage"].Owner != faction || w.faction(faction) == nil {
		t.Fatal("wrong estate handling")
	}
	if w.Properties["laundry"].Staff != staff || w.Properties["laundry"].Condition != condition {
		t.Fatal("death reset premises")
	}
	before, _ := json.Marshal(w)
	if w.Kill(n.ID, "again") {
		t.Fatal("duplicate death")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("duplicate estate release")
	}
}
func TestProprietorCannotBuyWithoutWorkingReserve(t *testing.T) {
	w := proprietorWorld()
	n := w.NPC("buyer")
	n.Purse = AcquisitionCost(w, "laundry")
	w.ConsiderProprietors()
	if w.Properties["laundry"].Owner != "independent" {
		t.Fatal("purchase exhausted necessities")
	}
	n.Purse = 3000
	n.Dead = true
	w.ConsiderProprietors()
	if w.Properties["laundry"].Owner != "independent" {
		t.Fatal("dead buyer")
	}
}

func TestFamilyLeaderDeathKeepsFamilyBusinessesThroughSuccession(t *testing.T) {
	w := New(7)
	f := &w.Factions[0]
	family := f.ID
	previous := f.Leader
	var leader *NPC
	for i := range w.NPCs {
		if w.NPCs[i].Name == previous {
			leader = &w.NPCs[i]
			break
		}
	}
	if leader == nil {
		t.Fatal("missing leader fixture")
	}
	w.Properties["garage"].Owner = family
	if !w.Kill(leader.ID, "Succession fixture") {
		t.Fatal("death failed")
	}
	f = w.faction(family)
	if f == nil || f.Leader == previous || w.Properties["garage"].Owner != family || w.CanAcquire("garage") {
		t.Fatal("family deed did not survive succession")
	}
}
