package core

import "testing"

// "Businesses are run, not merely owned." Capacity says what a place loses for
// being short-handed, out of stock or in trouble, and there is no other half:
// a business kept properly earns exactly what a business scraping by earns, so
// every decision about the people behind the counter is a way to avoid losing
// money rather than a way to make any.
//
// Measured before this: the publican, a policy that puts somebody in charge and
// pays over the rate, ended a hundred campaigns with a median of $470 against
// an investor's $2,367.

func wellRun(t *testing.T) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	trade, _ := TradeOf("laundry")
	w.Properties["laundry"].Supply = trade.RestockAmount
	return w, "laundry"
}

func TestAPlaceRunBySomebodyEarnsMoreThanOneNobodyMinds(t *testing.T) {
	t.Parallel()
	minded, id := wellRun(t)
	nobody, _ := wellRun(t)
	who := minded.Properties[id].Hands[0]
	if err := minded.PutInCharge(id, who); err != nil {
		t.Fatal(err)
	}
	run, alone := minded.Capacity(id), nobody.Capacity(id)
	t.Logf("a laundry run by somebody works at %.2f, one nobody minds at %.2f", run, alone)
	if run <= alone {
		t.Fatalf("somebody running the place is worth nothing: %.2f against %.2f", run, alone)
	}
	// And it is a share rather than a doubling: a manager is worth having, not
	// worth more than the people doing the work.
	if run > alone*1.3 {
		t.Fatalf("a manager is worth %.0f%% of the whole place", (run/alone-1)*100)
	}
}

func TestPeopleWhoThinkWellOfYouWorkBetter(t *testing.T) {
	t.Parallel()
	liked, id := wellRun(t)
	strangers, _ := wellRun(t)
	for _, who := range liked.Properties[id].Hands {
		liked.NPC(who).Trust = 80
	}
	warm, cold := liked.Capacity(id), strangers.Capacity(id)
	t.Logf("a laundry whose people think well of you works at %.2f, one of strangers at %.2f", warm, cold)
	if warm <= cold {
		t.Fatalf("what the people behind the counter think of you is worth nothing: %.2f against %.2f", warm, cold)
	}
}

// And the ordinary case is unchanged: a place nobody has touched works exactly
// as well as it did, or this is a tax on everybody who has not read the manual.
func TestAnUntouchedBusinessWorksAsWellAsItEverDid(t *testing.T) {
	t.Parallel()
	w, id := wellRun(t)
	if got := w.Capacity(id); got != 1 {
		t.Fatalf("a fully staffed, fully stocked business with no trouble works at %.2f", got)
	}
}
