package core

import "testing"

// Money into and out of a business moves in whatever figure the player types,
// not in fixed lots. The lot survives only as what the field starts on.

func TestTheFloatTakesTheFigureYouType(t *testing.T) {
	w := houseKeeper(t)
	cash := w.Player.Cash
	if err := w.Bankroll("casino", 700); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-700 || w.Properties["casino"].Bankroll != 700 {
		t.Fatalf("cash %d, float %d", w.Player.Cash, w.Properties["casino"].Bankroll)
	}
	if err := w.Draw("casino", 425); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-275 || w.Properties["casino"].Bankroll != 275 {
		t.Fatalf("after drawing 425: cash %d, float %d", w.Player.Cash, w.Properties["casino"].Bankroll)
	}
}

func TestAFigureNobodyTypedIsStillTheLot(t *testing.T) {
	w := houseKeeper(t)
	cash := w.Player.Cash
	if err := w.Bankroll("casino", 0); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-BankrollLot || w.Properties["casino"].Bankroll != BankrollLot {
		t.Fatalf("cash %d, float %d", w.Player.Cash, w.Properties["casino"].Bankroll)
	}
}

func TestTheFloatRefusesWhatIsNotThere(t *testing.T) {
	w := houseKeeper(t)
	w.Player.Cash = 400
	if w.BankrollReadiness("casino", 900) == "" {
		t.Fatal("put more behind the tables than the player had")
	}
	if w.BankrollReadiness("casino", 1) == "" {
		t.Fatal("a dollar is not a float")
	}
	if w.BankrollReadiness("casino", 400) != "" {
		t.Fatalf("refused every dollar the player had: %s", w.BankrollReadiness("casino", 400))
	}
	if err := w.Bankroll("casino", 400); err != nil {
		t.Fatal(err)
	}
	if w.DrawReadiness("casino", 401) == "" {
		t.Fatal("drew more than the float held")
	}
}

func TestMoneyLeavesAndComesHomeInTypedFigures(t *testing.T) {
	w := banker(t)
	cash := w.Player.Cash
	if err := w.Deposit(1200); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-1200 {
		t.Fatalf("the typed figure was not taken: cash %d", w.Player.Cash)
	}
	if w.Offshore != 1200*(100-DepositCut)/100 {
		t.Fatalf("the cut was not taken on the typed figure: %d", w.Offshore)
	}
	if w.DepositReadiness(9000) == "" {
		t.Fatal("wired out more than the player had")
	}
	w.Player.Offshore = true
	held := w.Offshore
	if err := w.Withdraw(300); err != nil {
		t.Fatal(err)
	}
	if w.Offshore != held-300 {
		t.Fatalf("bringing part home emptied the account: %d", w.Offshore)
	}
	if err := w.Withdraw(0); err != nil {
		t.Fatal(err)
	}
	if w.Offshore != 0 {
		t.Fatalf("a figure nobody typed left %d out there", w.Offshore)
	}
}

func TestTheInterfaceOffersAFieldForEachOfThem(t *testing.T) {
	w := houseKeeper(t)
	w.Event, w.District = nil, 9
	w.Properties["casino"].Bankroll = 800
	for _, id := range []string{"bankroll", "draw"} {
		a := actionByID(w.Actions("casino"), id)
		if a == nil {
			t.Fatalf("%s is not offered", id)
		}
		if a.Sum == nil {
			t.Fatalf("%s takes no typed figure", id)
		}
		if a.Sum.Least <= 0 || a.Sum.Most < a.Sum.Least || a.Sum.Preset < a.Sum.Least || a.Sum.Preset > a.Sum.Most {
			t.Fatalf("%s field is bounded %d..%d starting at %d", id, a.Sum.Least, a.Sum.Most, a.Sum.Preset)
		}
	}
	if a := actionByID(w.Actions("casino"), "draw"); a.Sum.Most != 800 {
		t.Fatalf("the draw field lets you take %d of an 800 float", a.Sum.Most)
	}
	b := banker(t)
	b.Event, b.District = nil, 9
	b.Offshore, b.Player.Offshore = 900, true
	for _, id := range []string{"deposit", "withdraw"} {
		a := actionByID(b.Actions("market"), id)
		if a == nil || a.Sum == nil {
			t.Fatalf("%s takes no typed figure", id)
		}
	}
	if a := actionByID(b.Actions("market"), "withdraw"); a.Sum.Most != 900 {
		t.Fatalf("the withdraw field caps at %d of 900 out there", a.Sum.Most)
	}
}

func actionByID(actions []Action, id string) *Action {
	for i := range actions {
		if actions[i].ID == id {
			return &actions[i]
		}
	}
	return nil
}
