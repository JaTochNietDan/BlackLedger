package core

import (
	"encoding/json"
	"testing"
)

func rentalWorld() *World {
	w := New(7)
	w.Minute = 1440
	w.NPCs = []NPC{{ID: "tenant", Name: "Tenant", Home: "room", Accommodation: "Rented room", Location: "bar", Post: "bar", Purse: 30}}
	w.Properties["room"].Owner = "player:1"
	return w
}
func TestRentReplacesBundledExpenseAndTransfersOnce(t *testing.T) {
	w := rentalWorld()
	before := w.Player.Cash
	w.PayTheCity()
	n := &w.NPCs[0]
	a := w.Properties["room"].Rents[n.ID]
	if n.Purse != 30 || w.Player.Cash-before != 15 || a.Paid != 15 || a.Arrears != 0 {
		t.Fatalf("rent accounting: purse=%d cash=%d account=%+v", n.Purse, w.Player.Cash-before, a)
	}
	w.PayTheCity()
	if n.Purse != 30 || w.Player.Cash-before != 15 || a.Collected != 15 {
		t.Fatal("duplicate settlement")
	}
	if w.HourlyIncome("room") != 0 || w.Books()["income"].(int) != 15 {
		t.Fatal("rent forecast and cash accrual duplicated or missing")
	}
}
func TestUnpaidRentIsDebtNotInventedIncome(t *testing.T) {
	w := rentalWorld()
	n := &w.NPCs[0]
	n.Faction = "bellandi"
	n.Purse = 10
	w.faction("bellandi").Short = 2
	before := w.Player.Cash
	w.PayTheCity()
	a := w.Properties["room"].Rents[n.ID]
	if n.Purse != 0 || a.Paid != 6 || a.Arrears != 9 || w.Player.Cash-before != 6 {
		t.Fatalf("short payday: %+v purse%d", a, n.Purse)
	}
	w.Minute += 1440
	w.faction("bellandi").Short = 0
	n.Purse = 50
	w.PayTheCity()
	if a.Arrears != 0 || a.Paid != 24 || a.Collected != 30 {
		t.Fatalf("arrears collection: %+v", a)
	}
}
func TestAbsentAndDeadTenantsAndTransfer(t *testing.T) {
	w := rentalWorld()
	n := &w.NPCs[0]
	n.Heading = "market"
	n.Arrives = w.Minute + 20
	w.PayTheCity()
	a := w.Properties["room"].Rents[n.ID]
	if a.Collected != 15 {
		t.Fatal("travelling resident lost tenancy")
	}
	n.Dead = true
	w.Minute += 1440
	before := w.Player.Cash
	w.PayTheCity()
	if w.RentalDaily("room") != 0 || w.Player.Cash != before {
		t.Fatal("dead resident paid rent")
	}
	n.Dead = false
	w.Properties["room"].Owner = "russo"
	cash := w.faction("russo").Cash
	w.PayTheCity()
	if w.Player.Cash != before || w.faction("russo").Cash-cash != 15 {
		t.Fatal("rent did not follow ownership")
	}
}
func TestRentalSettlementSurvivesReload(t *testing.T) {
	w := rentalWorld()
	w.PayTheCity()
	data, e := json.Marshal(w)
	if e != nil {
		t.Fatal(e)
	}
	var restored World
	if e = json.Unmarshal(data, &restored); e != nil {
		t.Fatal(e)
	}
	before := restored.Player.Cash
	restored.PayTheCity()
	if restored.Player.Cash != before || restored.Properties["room"].Rents["tenant"].Collected != 15 {
		t.Fatal("reload charged twice")
	}
}

func TestOwnerDoesNotPayRoomRentToThemselves(t *testing.T) {
	w := rentalWorld()
	w.Player.Home = "room"
	if w.HomeCost("room") != 0 || w.DailyCost() != 0 || w.Books()["costs"].(int) != 0 {
		t.Fatal("owner charged their own rent")
	}
	w.Properties["room"].Owner = "independent"
	if w.HomeCost("room") != 15 || w.DailyCost() != 15 {
		t.Fatal("ordinary rental charge lost")
	}
}

func TestRentRegisterHidesPrivateAccountsAndPublicIncomeMatches(t *testing.T) {
	w := rentalWorld()
	w.PayTheCity()
	register := w.RentRegister("room")
	rows := register["tenants"].([]map[string]any)
	if len(rows) != 1 || rows[0]["account"] == nil || register["daily"] != 15 {
		t.Fatalf("register: %+v", register)
	}
	if w.Public()["income"].(float64) != 15.0/24 {
		t.Fatal("public hourly forecast excludes rent")
	}
	w.Properties["room"].Owner = "independent"
	rows = w.RentRegister("room")["tenants"].([]map[string]any)
	if _, exists := rows[0]["account"]; exists {
		t.Fatal("private tenant balance exposed to nonowner")
	}
	if w.RentRegister("bar") != nil {
		t.Fatal("bar has a residential register")
	}
}
