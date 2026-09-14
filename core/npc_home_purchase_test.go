package core

import (
	"encoding/json"
	"testing"
)

func homeBuyerWorld() *World {
	w := New(7)
	w.NPCs = []NPC{{ID: "buyer", Name: "Rosa", Home: "room", Accommodation: "Rented room", Location: "garage", Post: "cabstand", Heading: "market", Sets: 610, Arrives: 650, Purse: 1500}}
	w.HouseholdSavings = map[string]HouseholdAccount{}
	a := w.HouseholdSavings["buyer"]
	a.Cash = 1000
	w.HouseholdSavings["buyer"] = a
	w.SettleApartments()
	return w
}
func TestNPCHomePurchaseUsesExactVacantDeedAndSavings(t *testing.T) {
	w := homeBuyerWorld()
	n := w.NPC("buyer")
	w.Properties["room"].Rents = map[string]*RentAccount{"buyer": {Arrears: 30}}
	before := *n
	wealth := w.HouseholdWealth(n)
	w.ApartmentDay()
	u := w.apartmentForResident(n.ID)
	if u == nil || u.Building != "apartment" || u.Owner != n.ID || n.Home != u.Building {
		t.Fatalf("missing purchased home: %+v %+v", n, u)
	}
	if w.HouseholdWealth(n) != wealth-w.ApartmentPrice(u) || w.NPCRent(n) != 0 {
		t.Fatal("purchase was not funded or owner charged rent")
	}
	if n.Location != before.Location || n.Post != before.Post || n.Heading != before.Heading || n.Sets != before.Sets || n.Arrives != before.Arrives || w.Properties["room"].Rents[n.ID].Arrears != 30 {
		t.Fatal("move altered trip/work/debt")
	}
	owned := 0
	for _, u := range w.Apartments {
		if u.Owner == n.ID {
			owned++
		}
	}
	if owned != 1 {
		t.Fatal("more than one purchase")
	}
	beforeJSON, _ := json.Marshal(w)
	w.ApartmentDay()
	afterJSON, _ := json.Marshal(w)
	if string(beforeJSON) != string(afterJSON) {
		t.Fatal("owner occupant bought or moved again")
	}
}
func TestNPCHomePurchasePreservesTenantsAndRejectsUnfundedMove(t *testing.T) {
	w := homeBuyerWorld()
	n := w.NPC("buyer")
	// None of the more expensive flats may be appropriated from their owners.
	for i := range w.Apartments {
		if w.Apartments[i].Building != "riverside" {
			w.Apartments[i].Owner = w.playerDeedID()
		}
	}
	n.Purse = 499
	a := w.HouseholdSavings[n.ID]
	a.Cash = 1000
	w.HouseholdSavings[n.ID] = a
	before, _ := json.Marshal(w)
	if w.purchaseNPCHome() {
		t.Fatal("spent the required reserve")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("failed purchase mutated world")
	}
	n.Purse++
	if !w.purchaseNPCHome() || n.Home != "riverside" || w.HouseholdWealth(n) != 500 {
		t.Fatal("affordable vacant home not purchased")
	}
	for _, u := range w.Apartments {
		if u.Building != "riverside" && u.Owner != w.playerDeedID() {
			t.Fatal("private deed taken")
		}
	}
}
func TestNPCHomePurchaseReleasesOnlyBuyersOldTenancy(t *testing.T) {
	w := homeBuyerWorld()
	n := w.NPC("buyer")
	n.Home = "riverside"
	w.NPCs = append(w.NPCs, NPC{ID: "tenant", Name: "Tenant", Home: "apartment", Purse: 0})
	w.SettleApartments()
	old := w.apartmentForResident(n.ID)
	old.Owner = w.playerDeedID()
	occupied := w.apartmentForResident("tenant")
	occupiedBefore := *occupied
	if !w.purchaseNPCHome() {
		t.Fatal("no move")
	}
	if old.Resident != "" || old.Owner != w.playerDeedID() || *occupied != occupiedBefore {
		t.Fatal("another deed or tenancy altered")
	}
	if w.apartmentForResident(n.ID).ID == old.ID {
		t.Fatal("old unit retained")
	}
}
