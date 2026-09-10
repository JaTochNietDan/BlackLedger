package core

import (
	"regexp"
	"testing"
)

// Reading every action's description side by side found three that quote terms
// for a transaction that cannot happen. A loan the player has no room left to
// make reads "$0 out, $0 back inside 7 days". Bringing money home with no
// account reads "Brings $0 back into the city". Servicing a car the player does
// not own reads "Restores up to 55 condition, currently 0 of 100".
//
// Nobody is misled about what they can do: all three are disabled, and their
// refusals already say why. The fault is narrower and still real — a
// description exists to tell the player what an action would do, and quoting
// zero dollars tells them nothing while looking like a figure.

var zeroMoney = regexp.MustCompile(`\$0\b`)

func TestNoDescriptionQuotesTermsOfZero(t *testing.T) {
	t.Parallel()
	w := New(7)
	w.District = 2
	w.Player.Location = "market"
	// A player with nothing: no cash, no car, no account, no room on the book.
	w.Player.Cash = 0
	w.Offshore = 0
	w.Player.Car = 0
	checked := 0
	for i := range Locations {
		id := Locations[i].ID
		w.Player.Location = id
		for _, a := range w.Actions(id) {
			checked++
			if m := zeroMoney.FindString(a.Detail); m != "" {
				t.Errorf("%s at %s quotes %q:\n  %s", a.ID, id, m, a.Detail)
			}
		}
	}
	t.Logf("read %d actions", checked)
}

// The car is the same shape without the dollar sign: a condition out of a
// hundred, for a car that is not there.
func TestServicingNoCarDoesNotQuoteItsCondition(t *testing.T) {
	t.Parallel()
	w := New(7)
	w.District = 2
	w.Player.Cash = 500
	w.Player.Car = 0
	w.Player.Location = "garage"
	for _, a := range w.Actions("garage") {
		if a.ID != "service" {
			continue
		}
		if !a.Disabled {
			t.Fatal("a car that does not exist can be worked on")
		}
		if contains(a.Detail, "of 100") {
			t.Fatalf("the description states the condition of a car nobody owns: %q", a.Detail)
		}
		return
	}
	t.Fatal("the garage does not offer to work on a car at all")
}

// And the other half, so this is not just a rule about suppressing text: when
// there IS something to state, the figure is still stated.
func TestRealTermsAreStillQuoted(t *testing.T) {
	t.Parallel()
	w := New(7)
	w.District = 2
	w.Player.Cash = 4000
	w.Offshore = 1200
	w.Player.Location = "market"
	sawWithdraw, sawLend := false, false
	for _, a := range w.Actions("market") {
		if a.ID == "withdraw" {
			sawWithdraw = true
			if !contains(a.Detail, "$1200") {
				t.Errorf("an account with money in it does not say how much: %q", a.Detail)
			}
		}
		if len(a.ID) > 5 && a.ID[:5] == "lend:" {
			sawLend = true
			if !contains(a.Detail, "$") || contains(a.Detail, "when there is room on your book") {
				t.Errorf("a loan the player can make does not quote its terms: %q", a.Detail)
			}
		}
	}
	if !sawWithdraw {
		t.Fatal("the market did not offer to bring the account home")
	}
	if !sawLend {
		t.Fatal("the market offered nobody a loan")
	}
	// A car in the yard still reports its condition.
	w.Player.Car = 1
	w.Player.CarWear = 40
	w.Player.Location = "garage"
	for _, a := range w.Actions("garage") {
		if a.ID == "service" {
			if !contains(a.Detail, "of 100") {
				t.Fatalf("a car the player owns does not report its condition: %q", a.Detail)
			}
			return
		}
	}
	t.Fatal("the garage did not offer to work on the car")
}
