package core

import "testing"

func TestApartmentPurchaseDoesNotClaimFormerOwnersRent(t *testing.T) {
	w := homeBuyerWorld()
	n := w.NPC("buyer")
	n.Home, n.Accommodation, n.Purse = "apartment", "Private apartment", 500
	w.SettleApartments()
	u := w.apartmentForResident(n.ID)
	u.Owner = "independent"
	w.collectRent(n)
	w.Player.Location, w.Player.Cash, w.District = "apartment", 10000, 9
	if err := w.BuyApartment(u.ID); err != nil {
		t.Fatal(err)
	}
	id := u.ID
	paid := func() any {
		for _, row := range w.ApartmentMarket() {
			if row["id"] == id {
				return row["rent_paid_today"]
			}
		}
		t.Fatal("missing owned apartment")
		return nil
	}
	if paid() != 0 {
		t.Fatal("claimed seller's receipt", paid())
	}
	w = w.Clone()
	n = w.NPC("buyer")
	w.Minute += 1440
	before := w.Player.Cash
	w.collectRent(n)
	if paid() != 35 || w.Player.Cash != before+35 {
		t.Fatal("new landlord receipt missing", paid())
	}
	account := w.Properties[n.Home].Rents[n.ID]
	account.PaidTo = ""
	if paid() != nil {
		t.Fatal("legacy recipient inferred", paid())
	}
}
