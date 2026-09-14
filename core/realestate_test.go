package core

import "testing"

func propertyTrader() *World {
	w := New(7)
	w.Minute = 600
	w.Event = nil
	w.Plots = nil
	w.Tasks = nil
	w.Contracts = nil
	w.NextPressure = 0
	w.Player.Cash = 20000
	w.Player.Respect = 30
	w.Player.Location = "room"
	return w
}

func TestPropertySaleKeepsHomeTenantsAndAccounts(t *testing.T) {
	w := propertyTrader()
	p := w.Properties["room"]
	p.Owner = "player:1"
	p.Posted = "guard"
	p.Rents = map[string]*RentAccount{"tenant": {Arrears: 9, Collected: 6}}
	homes := len(w.Residents("room"))
	cash, earned := w.Player.Cash, w.Player.Earned
	next, err := Execute(w, Command{Kind: "sell_property", Target: "room", Revision: w.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	if next.Own("room") || next.Properties["room"].Owner != "independent" || next.Player.Cash-cash != 2340 || next.Player.Earned != earned {
		t.Fatal("sale money or ownership incorrect")
	}
	if next.Player.Home != "room" || next.Player.Location != "room" || next.HomeCost("room") != 15 || len(next.Residents("room")) != homes {
		t.Fatal("sale displaced residents")
	}
	if next.Properties["room"].Rents["tenant"].Arrears != 9 || next.Properties["room"].Posted != "" || next.Wages() != 0 {
		t.Fatal("sale lost accounts or retained player business obligations")
	}
	if _, err = Execute(next, Command{Kind: "sell_property", Target: "room", Revision: next.Revision, RequestID: ID()}); err == nil {
		t.Fatal("same deed sold twice")
	}
}

func TestBuySellBuyCannotFarmCashOrRespect(t *testing.T) {
	w := propertyTrader()
	cash := w.Player.Cash
	act(t, &w, "acquire", "room")
	respect := w.Player.Respect
	act(t, &w, "sell_property", "room")
	act(t, &w, "acquire", "room")
	if w.Player.Respect != respect || w.Player.Cash >= cash || w.Properties["room"].BoughtLife != w.Life {
		t.Fatal("deed loop generated progress")
	}
}

func TestBuyingDeedDoesNotMoveAnybody(t *testing.T) {
	w := propertyTrader()
	w.Player.Location = "estate"
	w.District = 2
	playerHome := w.Player.Home
	residents := len(w.Residents("estate"))
	cash := w.Player.Cash
	next, err := Execute(w, Command{Kind: "buy_residence", Target: "estate", Revision: w.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	if !next.Own("estate") || next.Player.Cash != cash-3500 || next.Player.Home != playerHome || len(next.Residents("estate")) != residents {
		t.Fatal("deed purchase moved residents or charged incorrectly")
	}
	next.Player.Home = "estate"
	next.Player.Location = "estate"
	act(t, &next, "sell_property", "estate")
	if next.Player.Home != "estate" || next.Own("estate") || next.HomeCost("estate") != 90 {
		t.Fatal("sale and rent-back did not retain home")
	}
}

func TestBrokerOfferReflectsDamageAndRejectsInvalidDeeds(t *testing.T) {
	w := propertyTrader()
	w.Properties["estate"].Owner = "player:1"
	full := w.PropertyOffer("estate")
	w.Properties["estate"].Condition = 50
	if w.PropertyOffer("estate") >= full || w.PropertyOffer("estate") <= 0 {
		t.Fatal("damage not reflected")
	}
	w.Properties["estate"].Condition = 0
	if w.SellPropertyReadiness("estate") == "" {
		t.Fatal("worthless deed offered")
	}
	if w.SellPropertyReadiness("laundry") == "" || w.SellPropertyReadiness("unknown") == "" {
		t.Fatal("invalid deed sale allowed")
	}
}
