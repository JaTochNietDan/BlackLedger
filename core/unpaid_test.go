package core

import (
	"strings"
	"testing"
)

// A day the player cannot cover costs them their security, their address and
// their crew's loyalty. The people behind their counters — whose wages are most
// of that bill — lose nothing and notice nothing: you can miss payroll for a
// month and every hand still turns up.
//
// The wage bill names them. Missing it should reach them.

func broke2(t *testing.T) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Location = "laundry"
	w.Player.Cash = 200000
	if err := w.apply(Command{Kind: "acquire", Target: "laundry", RequestID: "buyitbeforebroke"}); err != nil {
		t.Fatal(err)
	}
	for _, who := range w.Properties["laundry"].Hands {
		w.NPC(who).Trust = 60
	}
	return w, "laundry"
}

func TestMissingPayrollReachesThePeopleItNames(t *testing.T) {
	t.Parallel()
	w, id := broke2(t)
	before := 0
	for _, who := range w.Properties[id].Hands {
		before += w.NPC(who).Trust
	}
	short := w.NobodyGotPaid()
	if short != len(w.Properties[id].Hands) {
		t.Fatalf("%d hands went unpaid and the laundry has %d", short, len(w.Properties[id].Hands))
	}
	after := 0
	for _, who := range w.Properties[id].Hands {
		after += w.NPC(who).Trust
	}
	t.Logf("a day nobody covered: %d hands, thinking of the player at %d against %d", short, after, before)
	if after >= before {
		t.Fatalf("a day nobody covered cost their regard nothing: %d against %d", after, before)
	}
}

// And the night that does not clear says so, rather than listing the security
// and the address and leaving out the people whose wages were most of the bill.
func TestTheNightThatDoesNotClearSaysWhoWentUnpaid(t *testing.T) {
	t.Parallel()
	w, _ := broke2(t)
	// A bill nothing could cover: security is charged by the day and the
	// laundry does not earn that fast.
	w.Player.Security, w.Player.Cash = 40, 0
	for day := 0; day < 3; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	said := false
	for _, r := range w.History {
		if r.Title == "Your arrangements unravel" && strings.Contains(r.Text, "went unpaid") {
			said = true
		}
	}
	if !said {
		t.Fatal("three nights of unpaid bills and the report never mentions the people behind the counters")
	}
}

// And a player who pays their way is untouched, which is the ordinary case.
func TestPayingYourWayCostsYouNothingWithThem(t *testing.T) {
	t.Parallel()
	w, id := broke2(t)
	before := 0
	for _, who := range w.Properties[id].Hands {
		before += w.NPC(who).Trust
	}
	for day := 0; day < 5; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	after := 0
	for _, who := range w.Properties[id].Hands {
		after += w.NPC(who).Trust
	}
	if after < before {
		t.Fatalf("a player who covered every bill lost their people's regard anyway: %d against %d", after, before)
	}
}
