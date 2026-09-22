package core

import (
	"strconv"
	"testing"
)

func TestDepositOfferAcceptsAffordableTypedAmountBelowLegacyLot(t *testing.T) {
	for _, cash := range []int{DepositLeast, 250, DepositLot - 1, DepositLot} {
		t.Run(strconv.Itoa(cash), func(t *testing.T) {
			w := banker(t)
			w.Event = nil
			w.Player.Cash = cash
			offer := actionByID(w.Actions("market"), "deposit")
			if offer == nil || offer.Disabled || offer.Sum == nil || offer.Sum.Preset > cash {
				t.Fatalf("affordable wire not offered at cash %d: %+v", cash, offer)
			}
			next, err := Execute(w, Command{Kind: "deposit", Target: "market", Amount: DepositLeast, Revision: w.Revision, RequestID: ID()})
			if err != nil {
				t.Fatal(err)
			}
			if next.Player.Cash != cash-DepositLeast || next.Offshore != DepositLeast*(100-DepositCut)/100 {
				t.Fatalf("wrong balances: cash=%d offshore=%d", next.Player.Cash, next.Offshore)
			}
		})
	}
}

func TestDepositOfferDoesNotRelaxCommandLimitsOrLegacyDefault(t *testing.T) {
	w := banker(t)
	w.Event = nil
	w.Player.Cash = 250
	for _, amount := range []int{0, DepositLeast - 1, 251} {
		if _, err := Execute(w, Command{Kind: "deposit", Target: "market", Amount: amount, Revision: w.Revision, RequestID: ID()}); err == nil {
			t.Fatalf("accepted invalid/legacy unaffordable amount %d", amount)
		}
	}
	w.Player.Cash = DepositLeast - 1
	offer := actionByID(w.Actions("market"), "deposit")
	if offer == nil || !offer.Disabled {
		t.Fatalf("offered wire below minimum: %+v", offer)
	}
}
