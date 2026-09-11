package core

import "testing"

// A campaign two jobs in whose fixer is dead crashed every command after it.
// The second job is offered by Mara, who is named by role rather than by name;
// the role is filled by whoever holds it, and a city that has buried them holds
// nobody. The proposal did not stand up, the error was thrown away, and the
// next line read a field off the nil scene — so the request panicked, the
// connection died under the browser, and the page said it had failed to fetch
// and could not reconnect, every time, for ever.
func TestTwoJobsInWithNobodyToBringTheThirdDoesNotCrash(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.JobCount = 2
	// Bury whoever holds the role, which is the only way this proposal stops
	// standing up.
	buried := 0
	for i := range w.NPCs {
		if w.NPCs[i].Role == "Fixer" || w.NPCs[i].ID == w.HolderID("fixer") {
			w.NPCs[i].Dead = true
			buried++
		}
	}
	if buried == 0 {
		t.Fatal("no fixer to bury, so this measures nothing")
	}
	if w.HolderID("fixer") != "" {
		t.Fatalf("the role refilled itself with %q, so the crash cannot be reached this way",
			w.HolderID("fixer"))
	}
	// This panicked.
	w.OfferIfReady()
	if w.Event != nil {
		t.Fatalf("a city with nobody to bring the work brought some: %q", w.Event.Title)
	}
	// And it is still quiet tomorrow rather than crashing then instead.
	for day := 0; day < 3; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
		w.OfferIfReady()
	}
}

// And a city that does have a fixer still gets the offer, which is the
// ordinary case and the thing the crash was in the way of.
func TestTwoJobsInWithAFixerStillGetsTheOffer(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.JobCount = 2
	if w.HolderID("fixer") == "" {
		t.Fatal("this city has no fixer to begin with")
	}
	w.OfferIfReady()
	if w.Event == nil || w.Event.Title != "A favor with a price" {
		t.Fatalf("the second job was not offered: %v", w.Event)
	}
}
