package core

import (
	"fmt"
	"testing"
)

func TestRentalMoveFillsPlayerVacancyAndPaysLandlord(t *testing.T) {
	w := homeBuyerWorld()
	n := w.NPC("buyer")
	n.Purse = 550
	w.HouseholdSavings = nil
	old := w.apartmentForResident(n.ID)
	var target *ApartmentDeed
	for i := range w.Apartments {
		u := &w.Apartments[i]
		if u.Building == "apartment" && u.Resident == "" {
			u.Resident = "reserved"
			if target == nil {
				target = u
			}
		}
	}
	target.Resident = ""
	target.Owner = w.playerDeedID()
	before := *n
	w.RentalMoveDay()
	if target.Resident != n.ID || target.Owner != w.playerDeedID() || (old != nil && old.Resident != "") {
		t.Fatal("vacancy not rented with deed preserved")
	}
	if n.Location != before.Location || n.Heading != before.Heading || n.Arrives != before.Arrives || n.Post != before.Post {
		t.Fatal("move changed trip/work")
	}
	cash := w.Player.Cash
	w.collectRent(n)
	if w.Player.Cash != cash+35 {
		t.Fatal("landlord not paid")
	}
	for _, row := range w.ApartmentMarket() {
		if row["id"] == target.ID {
			if row["rent_paid_today"] != 35 || row["vacant"] != false {
				t.Fatalf("bad accounts: %+v", row)
			}
			return
		}
	}
	t.Fatal("owned apartment missing")
}
func TestRentalMoveDoesNotDisplaceOrMoveOwnerOccupants(t *testing.T) {
	w := homeBuyerWorld()
	n := w.NPC("buyer")
	n.Home = "riverside"
	w.SettleApartments()
	old := w.apartmentForResident(n.ID)
	old.Owner = n.ID
	w.RentalMoveDay()
	if old.Resident != n.ID || n.Home != old.Building {
		t.Fatal("owner moved")
	}
	old.Owner = "independent"
	n.Purse = 0
	w.HouseholdSavings = nil
	w.RentalMoveDay()
	if old.Resident != n.ID {
		t.Fatal("unfunded move")
	}
}

func TestRentalUpgradeCanReachEveryAffordableTier(t *testing.T) {
	for _, tc := range []struct {
		wealth int
		home   string
	}{{69, "room"}, {70, "riverside"}, {84, "mercercourt"}, {105, "apartment"}, {280, "riverside"}, {349, "riverside"}, {350, "mercercourt"}, {489, "mercercourt"}, {490, "apartment"}} {
		t.Run(tc.home+fmt.Sprint(tc.wealth), func(t *testing.T) {
			w := homeBuyerWorld()
			n := w.NPC("buyer")
			n.Purse = tc.wealth
			w.HouseholdSavings = nil
			for _, id := range []string{"riverside", "mercercourt", "apartment"} {
				w.Properties[id].Condition = 100
				w.Properties[id].Trouble = false
			}
			w.RentalMoveDay()
			if n.Home != tc.home {
				t.Fatalf("$%d chose %s, want %s", tc.wealth, n.Home, tc.home)
			}
		})
	}
}
func TestRentalDowngradeChoosesBestAffordableVacancy(t *testing.T) {
	w := homeBuyerWorld()
	n := w.NPC("buyer")
	n.Home = "apartment"
	n.Accommodation = "Private apartment"
	n.Purse = 200
	w.HouseholdSavings = nil
	w.SettleApartments()
	w.RentalMoveDay()
	if n.Home != "mercercourt" {
		t.Fatalf("expected affordable middle tier, got %s", n.Home)
	}
}
