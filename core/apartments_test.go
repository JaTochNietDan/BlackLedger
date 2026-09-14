package core

import (
	"encoding/json"
	"fmt"
	"testing"
)

func apartmentWorld() *World {
	w := propertyTrader()
	w.Player.Home = "mercercourt"
	w.Player.Location = "mercercourt"
	w.NPCs = []NPC{{ID: "tenant", Name: "Tenant", Home: "mercercourt", Accommodation: "Private apartment", Purse: 3000, Location: "bar"}, {ID: "buyer", Name: "Buyer", Home: "apartment", Accommodation: "Private apartment", Purse: 5000, Location: "garage"}}
	w.SettleApartments()
	return w
}

func TestApartmentRegistryPreservesHomesAndIsReadOnlyDuringMovePlanning(t *testing.T) {
	w := apartmentWorld()
	u := *w.apartmentForResident("tenant")
	before, _ := json.Marshal(w)
	w.PlanHomeMove("apartment")
	w.ApartmentMarket()
	w.Actions("mercercourt")
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("read mutated apartment registry")
	}
	w.SettleApartments()
	if *w.apartmentForResident("tenant") != u || len(w.Apartments) != 384 {
		t.Fatal("stable apartment assignment changed")
	}
	for _, n := range w.NPCs {
		if n.Location == n.Home {
			t.Fatal("assignment teleported NPC")
		}
	}
}

func TestPlayerApartmentPurchaseSaleAndMovePreserveDeed(t *testing.T) {
	w := apartmentWorld()
	id := w.apartmentForResident(w.playerDeedID()).ID
	cash := w.Player.Cash
	act(t, &w, "buy_apartment:"+id, "mercercourt")
	if cash-w.Player.Cash != 1200 || w.HomeCost("mercercourt") != 0 || w.Own("mercercourt") {
		t.Fatal("purchase cost, rent or building ownership wrong")
	}
	w.Player.Location = "room"
	act(t, &w, "move_home", "room")
	if w.apartment(id).Owner != w.playerDeedID() || w.apartment(id).Resident != "" {
		t.Fatal("moving lost deed or retained occupant")
	}
	w.Player.Location = "mercercourt"
	cash = w.Player.Cash
	earned := w.Player.Earned
	act(t, &w, "sell_apartment:"+id, "mercercourt")
	if w.Player.Cash-cash != 780 || w.Player.Earned != earned {
		t.Fatal("sale creates earnings or wrong payment")
	}
	if _, err := Execute(w, Command{Kind: "sell_apartment:" + id, Target: "mercercourt", Revision: w.Revision, RequestID: ID()}); err == nil {
		t.Fatal("sold same deed twice")
	}
}

func TestNPCResidentBuysHomeThenTradesDeedWithoutTeleportOrMoneyCreation(t *testing.T) {
	w := apartmentWorld()
	w.apartmentForResident("buyer").Owner = "buyer"
	tenant := w.NPC("tenant")
	u := w.apartmentForResident(tenant.ID)
	id := u.ID
	w.ApartmentDay()
	if w.apartment(id).Owner != tenant.ID || tenant.Purse != 1800 || w.NPCRent(tenant) != 0 {
		t.Fatal("resident failed to purchase own flat")
	}
	// Stop the other resident buying their broker flat first: an already-owned
	// home still leaves them able to buy a second deed from a cash-poor seller.
	w.apartmentForResident("buyer").Owner = "buyer"
	tenant.Purse = 50
	buyer := w.NPC("buyer")
	total := tenant.Purse + buyer.Purse
	oldHome, oldPlace := tenant.Home, tenant.Location
	w.ApartmentDay()
	if w.apartment(id).Owner != "buyer" || tenant.Purse+buyer.Purse != total {
		t.Fatal("private sale failed to conserve money")
	}
	if tenant.Home != oldHome || tenant.Location != oldPlace || w.NPCRent(tenant) != 25 {
		t.Fatal("sale changed address/journey or omitted rent")
	}
	before := buyer.Purse
	w.collectRent(tenant)
	if buyer.Purse-before != 25 {
		t.Fatal("rent not paid to new owner")
	}
}

func TestApartmentLandlordIncomeComesFromTenantAndSurvivesSave(t *testing.T) {
	w := apartmentWorld()
	u := w.apartmentForResident("tenant")
	u.Owner = w.playerDeedID()
	w.NPC("tenant").Purse = 10
	cash := w.Player.Cash
	w.collectRent(w.NPC("tenant"))
	if w.Player.Cash-cash != 10 || w.NPC("tenant").Purse != 0 {
		t.Fatal("created rent not held by tenant")
	}
	if w.ApartmentRentIncome() != 25 {
		t.Fatal("contracted income missing")
	}
	w.collectRent(w.NPC("tenant"))
	if w.Player.Cash-cash != 10 {
		t.Fatal("double collected rent")
	}
	data, _ := json.Marshal(w)
	var restored World
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.apartment(u.ID).Owner != w.playerDeedID() {
		t.Fatal("reload lost deed")
	}
	restored.Player.Alive = false
	next, err := Execute(&restored, Command{Kind: "new_life", Revision: restored.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	if next.apartment(u.ID).Owner == next.playerDeedID() {
		t.Fatal("new protagonist inherited prior life deed")
	}
}

func TestApartmentNPCOwnerCannotBeDisplacedByPlayerMove(t *testing.T) {
	w := apartmentWorld()
	w.NPCs = nil
	for i := 0; i < ResidentialCapacity["mercercourt"]; i++ {
		id := ID()
		w.NPCs = append(w.NPCs, NPC{ID: id, Name: "Resident", Home: "mercercourt"})
	}
	w.Player.Home = "room"
	w.SettleApartments()
	for i := range w.Apartments {
		u := &w.Apartments[i]
		if u.Resident != "" {
			u.Owner = u.Resident
		}
	}
	if _, why := w.PlanHomeMove("mercercourt"); why == "" {
		t.Fatal("owner evicted to make room for player")
	}
}

func TestBrokerInvestmentPurchaseKeepsTenantAndCollectsOnlyRealCash(t *testing.T) {
	w := apartmentWorld()
	w.District = 2
	u := w.apartmentForResident("tenant")
	id, home, location := u.ID, w.Player.Home, w.NPC("tenant").Location
	cash, earned := w.Player.Cash, w.Player.Earned
	act(t, &w, "buy_apartment:"+id, "mercercourt")
	if w.Player.Cash != cash-1200 || w.Player.Earned != earned {
		t.Fatal("investment charged wrong principal or invented earnings")
	}
	if w.Player.Home != home || w.NPC("tenant").Home != "mercercourt" || w.NPC("tenant").Location != location || w.apartment(id).Resident != "tenant" {
		t.Fatal("deed transfer moved or evicted somebody")
	}
	w.NPC("tenant").Purse = 7
	cash = w.Player.Cash
	w.collectRent(w.NPC("tenant"))
	if w.Player.Cash != cash+7 || w.NPC("tenant").Purse != 0 {
		t.Fatal("investment rent was not funded by actual tenant cash")
	}
	act(t, &w, "sell_apartment:"+id, "mercercourt")
	if w.apartment(id).Resident != "tenant" {
		t.Fatal("sale evicted tenant")
	}
}

func TestBrokerListingsAreBoundedStableAndRotateAfterPurchase(t *testing.T) {
	w := apartmentWorld()
	for i := 0; i < 12; i++ {
		w.NPCs = append(w.NPCs, NPC{ID: fmt.Sprintf("listing-%02d", i), Name: "Tenant", Home: "mercercourt", Purse: 50})
	}
	w.SettleApartments()
	before, _ := json.Marshal(w)
	listed := w.ApartmentListings()
	count := 0
	for _, u := range w.Apartments {
		if listed[u.ID] && u.Building == "mercercourt" {
			count++
		}
	}
	if count != 5 {
		t.Fatalf("expected home, three tenants and one vacancy, got %d", count)
	}
	w.ApartmentMarket()
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("broker preview mutated world")
	}
	var bought string
	for _, u := range w.Apartments {
		if listed[u.ID] && u.Building == "mercercourt" && u.Resident != "" && u.Resident != w.playerDeedID() {
			bought = u.ID
			break
		}
	}
	if err := w.BuyApartment(bought); err != nil {
		t.Fatal(err)
	}
	next := w.ApartmentListings()
	if !next[bought] || len(next) != len(listed)+1 {
		t.Fatal("owned holding missing or broker failed to replenish")
	}
	w.apartment(bought).Owner = "tenant"
	if w.BuyApartmentReadiness(w.apartment(bought)) == "" {
		t.Fatal("private owner forced to sell")
	}
}

func TestApartmentMarketExplainsCurrentNeighborhoodDiscount(t *testing.T) {
	w := apartmentWorld()
	w.recordPropertyIncident("explosion", "mercercourt")
	before, _ := json.Marshal(w)
	found := false
	for _, row := range w.ApartmentMarket() {
		building := row["building"].(string)
		if row["neighborhood_index"] != w.NeighborhoodPropertyIndex(building) {
			t.Fatal("listing shows another neighborhood's pressure", row)
		}
		if building == "mercercourt" {
			found = true
			if row["neighborhood_index"] != 90 {
				t.Fatal("discount omitted", row)
			}
		}
	}
	if !found {
		t.Fatal("missing local listing")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("reading discount mutated campaign")
	}
}
