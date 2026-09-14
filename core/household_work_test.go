package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func householdWorkFixture() *World {
	w := New(81)
	w.District = 2
	w.Player.Location = "apartment"
	w.Minute = 8 * 60
	w.Event, w.Plots, w.Tasks = nil, nil, nil
	for i := range w.NPCs {
		w.NPCs[i].Purse = 0
	}
	w.HouseholdSavings = nil
	n := w.NPC("mara")
	n.Home, n.Purse = "apartment", 500
	w.SettleApartments()
	return w
}

func TestHouseholdWorkCommandTermsPaymentAndSavedQuota(t *testing.T) {
	w := householdWorkFixture()
	var offer Action
	before, _ := json.Marshal(w)
	for _, a := range w.Actions("apartment") {
		if a.ID == "householdwork" {
			offer = a
		}
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("reading offer mutated world")
	}
	if offer.ID == "" || offer.Disabled || offer.Group != "work" || offer.Subject != "mara" || offer.Minutes != 90 || !strings.Contains(offer.Detail, "$45") {
		t.Fatalf("bad offer: %+v", offer)
	}
	next, err := Execute(w, Command{Kind: "householdwork", Target: "apartment", RequestID: ID(), Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Cash != w.Player.Cash+45 || next.Player.Respect != w.Player.Respect+1 || next.Minute != w.Minute+90 {
		t.Fatalf("wrong reward or time: %+v", next.Player)
	}
	if _, err = Execute(next, Command{Kind: "householdwork", Target: "apartment", RequestID: ID(), Revision: next.Revision}); err == nil {
		t.Fatal("paid twice")
	}
	saved := next.Clone()
	saved.Life++
	if _, why := saved.HouseholdWorkOffer("apartment"); why == "" {
		t.Fatal("reload/new life reset booking")
	}
	saved.Minute += 1440
	if _, why := saved.HouseholdWorkOffer("apartment"); why != "" {
		t.Fatal("next day refused:", why)
	}
}

func TestHouseholdWorkTransfersRealFundsAndKeepsOriginalCustomer(t *testing.T) {
	for _, scenario := range []string{"paid", "dead", "moved", "poor"} {
		t.Run(scenario, func(t *testing.T) {
			w := householdWorkFixture()
			customer, err := w.StartHouseholdWork("apartment")
			if err != nil {
				t.Fatal(err)
			}
			n := w.NPC(customer)
			switch scenario {
			case "dead":
				n.Dead = true
			case "moved":
				n.Home = "room"
			case "poor":
				n.Purse = 144
			}
			cash, wealth := w.Player.Cash, w.HouseholdWealth(n)
			w.CompleteHouseholdWork("apartment", customer)
			if scenario == "paid" {
				if w.Player.Cash != cash+45 || w.HouseholdWealth(n) != wealth-45 {
					t.Fatal("payment not conserved")
				}
			} else if w.Player.Cash != cash || w.HouseholdWealth(n) != wealth {
				t.Fatal("invalid completion paid")
			}
		})
	}
}

func TestHouseholdWorkRejectsUnavailableBookingsWithoutMutation(t *testing.T) {
	for _, scenario := range []string{"early", "late", "damaged", "trouble", "poor", "wrong-place", "remote"} {
		t.Run(scenario, func(t *testing.T) {
			w := householdWorkFixture()
			target := "apartment"
			switch scenario {
			case "early":
				w.Minute = 479
			case "late":
				w.Minute = 991
			case "damaged":
				w.Properties[target].Condition = 39
			case "trouble":
				w.Properties[target].Trouble = true
			case "poor":
				w.NPC("mara").Purse = 144
			case "wrong-place":
				target = "bar"
				w.Player.Location = target
			case "remote":
				w.Player.Location = "bar"
			}
			before, _ := json.Marshal(w)
			if _, err := Execute(w, Command{Kind: "householdwork", Target: target, RequestID: ID(), Revision: w.Revision}); err == nil {
				t.Fatal("invalid job accepted")
			}
			after, _ := json.Marshal(w)
			if string(before) != string(after) {
				t.Fatal("rejection mutated input")
			}
		})
	}
	w := householdWorkFixture()
	w.Minute = 990
	if _, why := w.HouseholdWorkOffer("apartment"); why != "" {
		t.Fatal("last booking refused:", why)
	}
}

func TestHouseholdWorkInterruptedDoesNotPay(t *testing.T) {
	w := householdWorkFixture()
	w.Player.Contacts = 2
	w.Plots = []Plot{{ID: "repair-danger", Kind: "hit", Actor: "bellandi", Life: w.Life, Due: w.Minute + 120}}
	next, err := Execute(w, Command{Kind: "householdwork", Target: "apartment", RequestID: ID(), Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Event == nil || next.Event.Kind != "warning" {
		t.Fatal("fixture failed to interrupt")
	}
	if next.Player.Cash != w.Player.Cash || next.Properties["apartment"].HouseholdWorkDay == 0 {
		t.Fatal("interruption paid or released booking")
	}
}

func TestHouseholdWorkStartingDistrictAndSavings(t *testing.T) {
	w := householdWorkFixture()
	w.District, w.Player.Location = 0, "room"
	n := w.NPC("mara")
	n.Home, n.Purse = "room", 10
	w.HouseholdSavings = map[string]HouseholdAccount{"mara": {Cash: 200}}
	found := false
	for _, a := range w.Actions("room") {
		if a.ID == "householdwork" && !a.Disabled {
			found = true
		}
	}
	if !found {
		t.Fatal("starting district offers no household work")
	}
	customer, err := w.StartHouseholdWork("room")
	if err != nil {
		t.Fatal(err)
	}
	before := w.Player.Cash
	w.CompleteHouseholdWork("room", customer)
	if n.Purse != 0 || w.HouseholdSavings["mara"].Cash != 165 || w.Player.Cash != before+45 {
		t.Fatal("savings payment was not conserved")
	}
	if _, why := w.HouseholdWorkOffer("apartment"); why == "Today's household repair booking has already been taken" {
		t.Fatal("booking used another building's quota")
	}
}
