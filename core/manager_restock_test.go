package core

import (
	"strings"
	"testing"
)

func TestManagerRestockUsesHaulagePriceAndRecordsThePurchase(t *testing.T) {
	w, id := toRun(t)
	who := w.Properties[id].Hands[0]
	if err := w.PutInCharge(id, who); err != nil {
		t.Fatal(err)
	}
	w.Properties["haulage"].Owner = w.PlayerOrganizationID()
	w.Properties[id].Supply = 0
	trade, _ := TradeOf(id)
	cost := w.RestockCost(id)
	if cost >= trade.Restock {
		t.Fatal("fixture has no haulage discount")
	}
	w.Player.Cash = cost
	w.TheyRunIt()
	if w.Player.Cash != 0 || w.Properties[id].Supply != trade.RestockAmount {
		t.Fatal("manager ignored the shared discounted price")
	}
	receipts := 0
	for _, r := range w.History {
		if strings.HasPrefix(r.Title, "Restocked at ") {
			receipts++
		}
	}
	if receipts != 1 {
		t.Fatal("purchase missing or duplicate in ledger", receipts)
	}
	w.TheyRunIt()
	if w.Player.Cash != 0 {
		t.Fatal("full stock charged again")
	}
}
func TestManagerInCustodyCannotRestockTheBusiness(t *testing.T) {
	w, id := toRun(t)
	who := w.Properties[id].Hands[0]
	if err := w.PutInCharge(id, who); err != nil {
		t.Fatal(err)
	}
	w.Properties[id].Supply = 0
	cash := w.Player.Cash
	w.hold(w.NPC(who), 2)
	w.TheyRunIt()
	if w.Player.Cash != cash || w.Properties[id].Supply != 0 {
		t.Fatal("held manager purchased supplies")
	}
}
