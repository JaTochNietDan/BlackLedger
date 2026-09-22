package store

import (
	"blackledger/core"
	"encoding/json"
	"testing"
)

func TestExistingRentReceiptRetainsUnknownRecipientOnRead(t *testing.T) {
	s := testStore(t)
	w := core.New(59)
	w.SettleApartments()
	var unit *core.ApartmentDeed
	for i := range w.Apartments {
		u := &w.Apartments[i]
		if w.NPC(u.Resident) != nil {
			unit = u
			break
		}
	}
	if unit == nil {
		t.Fatal("fixture has no tenant")
	}
	unit.Owner = "player:1"
	w.Properties[unit.Building].Rents = map[string]*core.RentAccount{unit.Resident: {Day: w.Minute/1440 + 1, Paid: 12, Collected: 120, Arrears: 7}}
	raw, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.DB.Exec("UPDATE campaign SET state=? WHERE id=1", string(raw)); err != nil {
		t.Fatal(err)
	}
	got, err := s.Read()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, row := range got.ApartmentMarket() {
		if row["id"] == unit.ID {
			found = true
			if row["rent_paid_today"] != nil || row["rent_arrears"] != 7 {
				t.Fatal("legacy receipt inferred or debt changed", row)
			}
		}
	}
	if !found {
		t.Fatal("owned unit missing")
	}
	var stored string
	if err = s.DB.QueryRow("SELECT state FROM campaign WHERE id=1").Scan(&stored); err != nil {
		t.Fatal(err)
	}
	if stored != string(raw) || got.Minute != w.Minute || got.Player.Cash != w.Player.Cash {
		t.Fatal("read changed saved finances")
	}
}
