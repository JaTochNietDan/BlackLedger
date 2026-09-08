package core

import "testing"

func banker(t *testing.T) *World {
	t.Helper()
	w := New(503)
	w.Player.Location = "market"
	w.Player.Cash = 5000
	return w
}

func TestBankingMoneyCostsSomethingToDo(t *testing.T) {
	w := banker(t)
	cash := w.Player.Cash
	if err := w.Deposit(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-DepositLot {
		t.Fatalf("the lot was not taken: cash %d", w.Player.Cash)
	}
	if w.Offshore >= DepositLot {
		t.Fatalf("the arrangement took no cut: %d arrived of %d", w.Offshore, DepositLot)
	}
	if w.Offshore == 0 {
		t.Fatal("nothing arrived at all")
	}
	// It earns nothing sitting there.
	held := w.Offshore
	w.Advance(1440 * 10)
	if w.Offshore != held {
		t.Fatalf("the balance moved on its own: %d then %d", held, w.Offshore)
	}
}

func TestMoneySentOutSurvivesThePersonWhoSentIt(t *testing.T) {
	w := banker(t)
	for i := 0; i < 3; i++ {
		if err := w.Deposit(); err != nil {
			t.Fatal(err)
		}
	}
	banked := w.Offshore
	if banked == 0 {
		t.Fatal("nothing was banked")
	}
	if err := w.EstablishAccess(); err != nil {
		t.Fatal(err)
	}
	if !w.Player.Offshore {
		t.Fatal("access was not established")
	}

	w.Die("Test death")
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	if next.Offshore != banked {
		t.Fatalf("the balance did not survive: %d, expected %d", next.Offshore, banked)
	}
	// But nothing else is inherited, and the account does not answer to a stranger.
	if next.Player.Offshore {
		t.Fatal("a new person inherited access to the account")
	}
	if next.Player.Cash != 90 {
		t.Fatalf("a new person started with %d rather than the usual 90", next.Player.Cash)
	}
	if next.WithdrawReadiness() == "" {
		t.Fatal("a stranger could draw on it without establishing anything")
	}
}

func TestReachingItAsAStrangerCostsMoreThanAStrangerHas(t *testing.T) {
	w := New(509)
	w.Offshore = 4000
	w.Player.Location = "market"
	if w.Player.Cash >= AccessCost {
		t.Fatalf("a new arrival starts with %d, which already covers the %d access cost", w.Player.Cash, AccessCost)
	}
	if w.AccessReadiness() == "" {
		t.Fatal("a penniless new arrival could establish access")
	}
	if err := w.EstablishAccess(); err == nil {
		t.Fatal("the command ignored a requirement the action reports")
	}
	// Once they have earned enough, it is theirs.
	w.Player.Cash = AccessCost
	if reason := w.AccessReadiness(); reason != "" {
		t.Fatal("somebody who earned the fee was still refused:", reason)
	}
	if err := w.EstablishAccess(); err != nil {
		t.Fatal(err)
	}
	if err := w.Withdraw(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != 4000 {
		t.Fatalf("withdrawing brought back %d", w.Player.Cash)
	}
	if w.Offshore != 0 {
		t.Fatal("the account still holds money after being emptied")
	}
}

func TestAnEmptyAccountIsNotWorthReaching(t *testing.T) {
	w := banker(t)
	if w.AccessReadiness() == "" {
		t.Fatal("access was offered with nothing out there")
	}
	if w.WithdrawReadiness() == "" {
		t.Fatal("a withdrawal was offered from an empty account")
	}
	if err := w.Withdraw(); err == nil {
		t.Fatal("withdrew from an empty account")
	}
	// And none of it is arranged anywhere but where such things are arranged.
	w.Offshore = 1000
	w.Player.Location = "docks"
	for _, reason := range []string{w.DepositReadiness(), w.AccessReadiness(), w.WithdrawReadiness()} {
		if reason == "" {
			t.Fatal("banking was arranged at the waterfront")
		}
	}
}

func TestBankingIsALossNotASavingsPlan(t *testing.T) {
	// Sending money out and bringing it back must never be worth doing for its
	// own sake, or it becomes a way to launder value rather than a decision
	// about what to leave behind.
	w := banker(t)
	w.Player.Cash = 100000
	start := w.Player.Cash
	for i := 0; i < 20; i++ {
		if err := w.Deposit(); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.EstablishAccess(); err != nil {
		t.Fatal(err)
	}
	if err := w.Withdraw(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash >= start {
		t.Fatalf("a round trip through the account made money: %d from %d", w.Player.Cash, start)
	}
}
