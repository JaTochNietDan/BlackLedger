package core

import (
	"strings"
	"testing"
)

// The day's bill was all or nothing. A player a dollar short lost their
// security, their address and their crew's loyalty in a single night — and then
// the bill was $33 instead of $2,033, so it cleared every night after and
// nothing was ever unpaid again. Being briefly short was a cliff, and being
// persistently short was impossible.
//
// It gives things up in order now, one at a time, and each only while what is
// left of the bill is still out of reach.
//
// These press on Settle itself rather than on a day of the world. A day pays
// the player as it goes, so a world set a dollar short at dawn is comfortable
// by midnight and proves nothing about the order things are given up in.

func settling(t *testing.T) *World {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Location = "laundry"
	w.Player.Cash = 200000
	if err := w.apply(Command{Kind: "acquire", Target: "laundry", RequestID: "buyitthensettle"}); err != nil {
		t.Fatal(err)
	}
	w.Player.Security, w.Player.Home = 20, "apartment"
	if w.Wages() == 0 {
		t.Fatal("nobody stands behind this counter, so the wages cannot be the last thing to go")
	}
	return w
}

// payrollOnly is the day with everything given up that can be: no security, no
// crew, and a rented room. What is left is the payroll and the standing costs
// nobody can hand back at midnight.
func payrollOnly(w *World) int {
	return w.DailyCost() - 10*w.Player.Security - 12*len(w.Player.Crew) -
		(HomeRent(w.Player.Home) - HomeRent("room"))
}

func TestBeingALittleShortCostsTheLeastThingFirst(t *testing.T) {
	t.Parallel()
	w := settling(t)
	w.Player.Cash = w.DailyCost() - 1
	w.Settle()
	if w.Player.Security > 0 {
		t.Fatal("a player who could not cover the day kept their security")
	}
	if w.Player.Home != "apartment" {
		t.Fatalf("being a dollar short cost them their address as well: %q", w.Player.Home)
	}
	if w.Properties["laundry"].Unpaid > 0 {
		t.Fatal("being a dollar short left the wages unpaid too")
	}
	if w.Player.Cash <= 0 {
		t.Fatal("the security was given up and the day was still taken in full")
	}
}

// And the wages are the last thing to go: a player with nothing but the payroll
// still makes the payroll, having given up everything else to do it.
func TestTheWagesAreTheLastThingToGoUnpaid(t *testing.T) {
	t.Parallel()
	w := settling(t)
	w.Player.Cash = payrollOnly(w)
	w.Settle()
	if w.Player.Security > 0 || w.Player.Home != "room" {
		t.Fatalf("a player down to the payroll kept their security or their address: %d, %q",
			w.Player.Security, w.Player.Home)
	}
	if w.Properties["laundry"].Unpaid != 0 {
		t.Fatal("the wages went unpaid while there was still the money for them")
	}

	// A dollar less and they do not.
	x := settling(t)
	x.Player.Cash = payrollOnly(x) - 1
	x.Settle()
	if x.Properties["laundry"].Unpaid == 0 {
		t.Fatal("a player who could not cover the payroll paid it anyway")
	}
	said := false
	for _, r := range x.History {
		said = said || strings.Contains(r.Text, "went unpaid")
	}
	if !said {
		t.Fatal("nobody was told the counters went unpaid")
	}
}

// A player who covers everything is untouched, which is the ordinary case.
func TestAPlayerWhoCoversTheDayLosesNothing(t *testing.T) {
	t.Parallel()
	w := settling(t)
	home, guard := w.Player.Home, w.Player.Security
	for day := 0; day < 5; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if w.Player.Home != home || w.Player.Security != guard {
		t.Fatalf("a player who paid every bill lost their address or their security: %q, %d",
			w.Player.Home, w.Player.Security)
	}
	if w.Properties["laundry"].Unpaid != 0 {
		t.Fatal("a player who paid every bill has unpaid nights against their laundry")
	}
}
