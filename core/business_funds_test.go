package core

import "testing"

func TestBusinessFundsFollowPersonalDeedAndCasinoFloat(t *testing.T) {
	w := New(7)
	n := w.NPC("mara")
	n.Purse = 100
	w.Properties["poolhall"].Owner = n.ID
	w.changeBusinessFunds("poolhall", 25)
	if n.Purse != 125 {
		t.Fatal("hall fee missed owner")
	}
	w.changeBusinessFunds("poolhall", -30)
	if n.Purse != 95 {
		t.Fatal("owner loss not charged")
	}
	w.Properties["casino"].Owner = n.ID
	w.Properties["casino"].Bankroll = 500
	w.changeBusinessFunds("casino", -80)
	if n.Purse != 95 || w.businessFunds("casino") != 420 {
		t.Fatal("gambling used household instead of float")
	}
	w.changeBusinessFunds("casino", 20)
	if w.Properties["casino"].Bankroll != 440 {
		t.Fatal("house win missed float")
	}
	n.Dead = true
	if w.changeBusinessFunds("casino", 10) || w.businessFunds("casino") != 0 {
		t.Fatal("dead proprietor traded")
	}
}
func TestProprietorCasinoPurchaseFundsFloatAndPreservesReserve(t *testing.T) {
	w := proprietorWorld()
	w.Properties["laundry"].Owner = "reserved"
	w.Properties["casino"].Owner = "independent"
	n := w.NPC("buyer")
	n.Purse = 10000
	before := w.HouseholdWealth(n)
	price := AcquisitionCost(w, "casino")
	float := w.Properties["casino"].Bankroll
	w.ConsiderProprietors()
	if w.Properties["casino"].Owner != n.ID || w.HouseholdWealth(n) != before-price-BankrollLot || w.Properties["casino"].Bankroll != float+BankrollLot {
		t.Fatal("purchase created unfunded float")
	}
	if w.NightHandleAt("casino") == 0 {
		t.Fatal("personal casino cannot trade")
	}
}
func TestPersonalCasinoFailureDoesNotPenalizePlayer(t *testing.T) {
	w := New(7)
	n := w.NPC("mara")
	p := w.Properties["casino"]
	p.Owner = n.ID
	p.Bankroll = BankrollFull
	respect := w.Player.Respect
	// Deterministic range exercises both winning and losing nightly outcomes.
	for i := 0; i < 100; i++ {
		p.Bankroll = 200
		w.CasinoDay()
	}
	if w.Player.Respect != respect {
		t.Fatal("NPC casino changed player reputation")
	}
}
func TestAllPricedBusinessTradesHavePersonalPurchasePath(t *testing.T) {
	for _, l := range Locations {
		if l.Cost > 0 && l.Type != "home" {
			if _, ok := TradeOf(l.ID); ok && !personalBusiness(l.ID) {
				t.Fatal(l.ID)
			}
		}
	}
}

func TestFuelPurchasePaysIndividualStationOwner(t *testing.T) {
	w := New(7)
	n := w.NPC("mara")
	w.Properties["filling"].Owner = n.ID
	w.Player.Car = 1
	w.Player.Fuel = 0
	w.Player.Fuelled = 1
	w.Player.Cash = 1000
	fee := w.FuelFee("filling")
	before := n.Purse
	if err := w.FillUp("filling"); err != nil {
		t.Fatal(err)
	}
	if n.Purse != before+fee || w.Player.Cash != 1000-fee {
		t.Fatal("fuel money missed proprietor")
	}
}
func TestNPCHallOwnerReceivesTournamentCutExactlyOnce(t *testing.T) {
	w := hostFixture(t)
	w.NPCs = append(w.NPCs, NPC{ID: "hall-owner", Name: "Hall Owner", Purse: 1000})
	w.Properties[PoolPlace].Owner = "hall-owner"
	if err := w.startPoolTournament([]string{"leo", "mara", "elena", "vittorio"}, 50, 20, false); err != nil {
		t.Fatal(err)
	}
	tour := w.PoolTournament
	for i := 0; i < 2; i++ {
		tour.Bracket.Matches[i].Rack.Concede(1)
	}
	w.ReconcilePoolTournament()
	tour.Bracket.Matches[2].Rack.Concede(1)
	w.ReconcilePoolTournament()
	if !tour.Settled || tour.HouseCutPaid != 40 || tour.PrizePaid != 160 || w.NPC("hall-owner").Purse != 1040 {
		t.Fatal("cut did not reach owner")
	}
	w.ReconcilePoolTournament()
	if w.NPC("hall-owner").Purse != 1040 {
		t.Fatal("cut paid twice")
	}
}

func TestPersonalMarinerReceivesTenantPaymentsWithoutPhantomRent(t *testing.T) {
	w := New(7)
	owner := w.NPC("mara")
	tenant := w.NPC("leo")
	p := w.Properties["room"]
	p.Owner = owner.ID
	p.Rents = map[string]*RentAccount{}
	tenant.Home = "room"
	tenant.Purse = 100
	before := owner.Purse
	rent := w.NPCRent(tenant)
	w.collectRent(tenant)
	if owner.Purse != before+rent || tenant.Purse != 100-rent {
		t.Fatal("rent did not follow lodging deed")
	}
	if !personalBusiness("room") {
		t.Fatal("lodging freehold excluded")
	}
	owner.Purse = 1000
	p.Income = 10000
	p.Condition = 100
	p.Supply = 10000
	p.ProprietorDay = 0
	trade, _ := TradeOf("room")
	w.ProprietorDay()
	if owner.Purse != 1000-p.Staff*trade.Wage {
		t.Fatal("lodging gained abstract rent")
	}
}
