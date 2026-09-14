package core

import (
	"math"
	"testing"
)

func TestAccountsUseTheSameCurrentIncomeAsTheClock(t *testing.T) {
	w := New(2)
	for id, p := range w.Properties {
		if w.Own(id) {
			p.Owner = "independent"
		}
	}
	p := w.Properties["laundry"]
	p.Owner = "player:1"
	p.Condition = 73
	w.NPCs = nil
	expected := float64(p.Income*p.Condition) / 100 * operatingMode(p.Mode).Take * w.Capacity("laundry") * w.TradeMultiplier("laundry") * w.RoomTrade("laundry") * (1 + w.LicenceTake()) * w.CollectionShare("laundry")
	if math.Abs(w.HourlyIncome("laundry")-expected) > 1e-9 {
		t.Fatal("clock rate lost an operating factor")
	}
	if got := w.Books()["income"]; got != int(expected*24) {
		t.Fatalf("accounts %v, actual daily rate %d", got, int(expected*24))
	}
	if w.HourlyIncome("missing") != 0 || w.HourlyIncome("bar") != 0 {
		t.Fatal("income credited for an unowned address")
	}
}
