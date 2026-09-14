package core

import "testing"

func TestRiversideCoversPopulationCeilingWithoutMovingResidents(t *testing.T) {
	w := New(7)
	old := map[string]NPC{}
	for _, n := range w.NPCs {
		old[n.ID] = n
	}
	deeds := append([]ApartmentDeed(nil), w.Apartments...)
	for len(w.NPCs) < MaxPeople {
		if w.AddCivilian() == nil {
			t.Fatal("population stopped before ceiling")
		}
	}
	w.SettleHousing()
	w.SettleApartments()
	if w.HousingShortage() != 0 {
		t.Fatalf("%d people have no home", w.HousingShortage())
	}
	if len(w.Residents("riverside")) == 0 {
		t.Fatal("new homes unused")
	}
	for _, n := range w.NPCs {
		if prev, ok := old[n.ID]; ok && (n.Home != prev.Home || n.Location != prev.Location) {
			t.Fatalf("established resident moved: %s", n.ID)
		}
		if n.Home != "room" && n.Home != "estate" && w.apartmentForResident(n.ID) == nil {
			t.Fatalf("missing numbered apartment: %s", n.ID)
		}
	}
	for _, d := range deeds {
		current := w.apartment(d.ID)
		if current == nil || current.Owner != d.Owner || (d.Resident != "" && current.Resident != d.Resident) {
			t.Fatalf("existing deed changed: %s", d.ID)
		}
	}
	for id, capacity := range ResidentialCapacity {
		occupied := len(w.Residents(id))
		if w.Player.Home == id {
			occupied++
		}
		if occupied > capacity {
			t.Fatalf("%s over capacity", id)
		}
	}
}

func TestRiversideLeaseDeedRentAndMigration(t *testing.T) {
	w := New(7)
	w.Event = nil
	w.Plots = nil
	w.Tasks = nil
	w.Contracts = nil
	delete(w.Properties, "riverside")
	kept := w.Apartments[:0]
	for _, u := range w.Apartments {
		if u.Building != "riverside" {
			kept = append(kept, u)
		}
	}
	w.Apartments = kept
	cash, minute := w.Player.Cash, w.Minute
	w.MigrateLivingWorld()
	w.SettleHousing()
	w.SettleApartments()
	if w.Properties["riverside"] == nil || len(w.Apartments) != 384 || w.Player.Cash != cash || w.Minute != minute {
		t.Fatal("additive migration failed")
	}
	w.District = 2
	w.Player.Location = "riverside"
	w.Player.Cash = 2000
	next, err := Execute(w, Command{Kind: "move_home", Target: "riverside", Revision: w.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	w = next
	if w.Player.Home != "riverside" || w.Player.Cash != 1900 || w.HomeCost("riverside") != 20 || HomeRank("riverside") != 1 {
		t.Fatal("riverside lease terms incorrect")
	}
	u := w.apartmentForResident(w.playerDeedID())
	if u == nil {
		t.Fatal("player has no flat")
	}
	price := w.ApartmentPrice(u)
	if price != 1000 {
		t.Fatalf("initial deed price: %d", price)
	}
	if err := w.BuyApartment(u.ID); err != nil {
		t.Fatal(err)
	}
	if w.HomeCost("riverside") != 0 || w.Player.Cash != 900 {
		t.Fatal("ownership failed to end rent")
	}
	if err := w.SellApartment(u.ID); err != nil {
		t.Fatal(err)
	}
	if w.HomeCost("riverside") != 20 || w.Player.Cash != 1550 || u.Resident != w.playerDeedID() {
		t.Fatal("sale changed tenancy or proceeds")
	}
	n := NPC{ID: "riverside-tenant", Home: "riverside", Accommodation: "Shared flat", Purse: 100}
	if w.NPCRent(&n) != 10 {
		t.Fatal("shared rent incorrect")
	}
	w.collectRent(&n)
	if n.Purse != 90 || w.Properties["riverside"].Rents[n.ID].Paid != 10 {
		t.Fatal("rent did not debit resident")
	}
}

func TestRiversidePrivateSaleConservesMoneyAndRespondsToNeighborhoodCrime(t *testing.T) {
	w := New(7)
	w.NPCs = []NPC{{ID: "seller", Name: "Seller", Home: "riverside", Location: "bar", Purse: 50}, {ID: "buyer", Name: "Buyer", Home: "riverside", Location: "docks", Purse: 3000}}
	w.SettleApartments()
	unit := w.apartmentForResident("seller")
	unit.Owner = "seller"
	w.apartmentForResident("buyer").Owner = "buyer"
	w.Witness("gunfight", "garage", "Recorded gunfire.", "")
	price := w.ApartmentPrice(unit)
	if price >= 1000 || w.ApartmentPrice(w.apartment("mercercourt-01")) != 1200 {
		t.Fatal("crime price did not stay local")
	}
	w.ApartmentDay()
	if unit.Owner != "buyer" || w.NPC("buyer").Purse != 3000-price || w.NPC("seller").Purse != 50+price {
		t.Fatal("private sale failed to transfer actual funds")
	}
	if unit.Resident != "seller" || w.NPC("seller").Home != "riverside" || w.NPC("seller").Location != "bar" {
		t.Fatal("sale displaced resident")
	}
}
