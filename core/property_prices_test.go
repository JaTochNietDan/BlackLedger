package core

import (
	"encoding/json"
	"testing"
)

func TestNeighborhoodPropertyPricesReactLocallyAndRecoverAcrossReload(t *testing.T) {
	w := propertyTrader()
	normalRoom, normalEstate := w.ResidencePrice("room"), w.ResidencePrice("estate")
	rng, worldRNG, minute := w.RNG, w.WorldRNG, w.Minute
	w.Witness("gunfight", "bar", "Recorded gunfire.", "")
	if w.ResidencePrice("room") >= normalRoom || w.ResidencePrice("estate") != normalEstate {
		t.Fatal("shock did not remain in the affected district")
	}
	if w.RNG != rng || w.WorldRNG != worldRNG || w.Minute != minute {
		t.Fatal("market recording consumed time/randomness")
	}
	if w.NeighborhoodPropertyIndex("room") != 94 {
		t.Fatal("wrong initial shock")
	}
	w.Minute += 720
	w.Witness("robbery", "bar", "Recorded robbery.", "")
	data, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var restored World
	if err = json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.NeighborhoodPropertyIndex("room") != 92 {
		t.Fatal("reload lost accumulated pressure")
	}
	restored.Minute += 720
	if restored.NeighborhoodPropertyIndex("room") != 93 {
		t.Fatal("partial-day decay lost when a new incident was recorded")
	}
	restored.Minute += 7 * 1440
	if restored.ResidencePrice("room") != normalRoom {
		t.Fatal("quiet neighborhood did not recover")
	}
}

func TestNeighborhoodPricesCannotBeMovedByReadingOrPrivatePlans(t *testing.T) {
	w := propertyTrader()
	w.Witness("politics", "bar", "A person left.", "")
	w.Witness("killing", "unknown", "No valid location.", "")
	if len(w.PropertyPressure) != 0 {
		t.Fatal("non-incident changed market")
	}
	w.Witness("explosion", "bar", "Explosion.", "")
	before, _ := json.Marshal(w)
	for i := 0; i < 5; i++ {
		w.PropertyMarket()
		w.Actions("room")
		w.PropertyOffer("room")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("reading prices mutated saved state")
	}
	for i := 0; i < 100; i++ {
		w.Witness("attack", "bar", "Attack.", "")
	}
	if w.NeighborhoodPropertyIndex("room") != 60 || w.ResidencePrice("room") <= w.PropertyOffer("room") {
		t.Fatal("price floor or broker spread broken")
	}
}

func TestEstateMarketActionAndPaymentAgreeAfterCrime(t *testing.T) {
	for _, kind := range []string{"buy_residence", "move_home"} {
		t.Run(kind, func(t *testing.T) {
			w := propertyTrader()
			w.Player.Location = "estate"
			w.District = 2
			w.Witness("attack", "estate", "Residence attacked.", "")
			price := w.ResidencePrice("estate")
			found := false
			for _, a := range w.Actions("estate") {
				if a.ID == kind {
					found = true
					if a.Cost != price {
						t.Fatalf("action=%d price=%d", a.Cost, price)
					}
				}
			}
			if !found {
				t.Fatal("missing purchase action")
			}
			for _, p := range w.PropertyMarket() {
				if p["id"] == "estate" && p["asking"] != price {
					t.Fatal("market quote differs")
				}
			}
			cash := w.Player.Cash
			next, err := Execute(w, Command{Kind: kind, Target: "estate", Revision: w.Revision, RequestID: ID()})
			if err != nil {
				t.Fatal(err)
			}
			if !next.Own("estate") || next.Player.Cash != cash-price {
				t.Fatal("purchase paid different price")
			}
			if kind == "buy_residence" && next.Player.Home != w.Player.Home {
				t.Fatal("deed purchase relocated player")
			}
		})
	}
}

func TestMarinerCrimeDiscountPreservesAcquisitionPremiumAndSaleSpread(t *testing.T) {
	w := propertyTrader()
	w.Witness("incendiary", "bar", "Fire.", "")
	price := AcquisitionCost(w, "room")
	if price >= MarinerFreehold || price != w.ResidencePrice("room") {
		t.Fatal("Mariner asking price ignored neighborhood")
	}
	cash := w.Player.Cash
	act(t, &w, "acquire", "room")
	if cash-w.Player.Cash != price {
		t.Fatal("Mariner charged wrong amount")
	}
	offer := w.PropertyOffer("room")
	cash = w.Player.Cash
	act(t, &w, "sell_property", "room")
	if w.Player.Cash-cash != offer || offer >= price {
		t.Fatal("sale created a cash loop")
	}
}

func TestNewLifeInheritsNeighborhoodMarketRatherThanResettingIt(t *testing.T) {
	w := propertyTrader()
	w.Witness("gunfight", "bar", "Gunfire.", "")
	w.Player.Alive = false
	next, err := Execute(w, Command{Kind: "new_life", Revision: w.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	if next.Life != w.Life+1 || next.NeighborhoodPropertyIndex("room") != 94 {
		t.Fatal("new arrival erased the neighborhood's recent violence")
	}
}
