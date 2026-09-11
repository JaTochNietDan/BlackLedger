package core

import (
	"strings"
	"testing"
)

// Answering to a family that no longer exists.
//
// The city announces an organization ending and counts the people it put on the
// street — "nine people are out of work tonight". A player who answered to that
// name was one of them and was never told: the wage simply stopped, the rank
// went, and the only sign was a card that had quietly changed.
//
// Reachable now in a way it was not. A city left alone buries three
// organizations in a hundred long campaigns, where the sample this project
// measured through buried none at all.
func TestAFamilyEndingIsSaidToTheManWhoWorkedForIt(t *testing.T) {
	t.Parallel()
	w := New(spread(3))
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 40000, 200

	family := ""
	for i := range w.Factions {
		if f := &w.Factions[i]; f.ID != w.PlayerOrganizationID() {
			family = f.ID
			break
		}
	}
	if family == "" {
		t.Skip("this city has nobody to answer to")
	}
	// Signed on the way the game signs somebody on, then the name goes.
	w.Player.Serves, w.Player.Service = family, PromotionWork
	rank := w.ServiceTitle()
	if rank == "" {
		t.Fatal("somebody who has done the work has no title, so nothing is lost")
	}
	kept := w.Factions[:0]
	for _, f := range w.Factions {
		if f.ID != family {
			kept = append(kept, f)
		}
	}
	w.Factions = kept

	before := len(w.History)
	w.ServiceDay()
	if w.Player.Serves != "" {
		t.Fatalf("the family is gone and the player still answers to %q", w.Player.Serves)
	}
	if w.Player.Service != 0 {
		t.Fatalf("the family is gone and %d pieces of work still count", w.Player.Service)
	}
	told := false
	for _, r := range w.History[min(before, len(w.History)):] {
		if strings.Contains(r.Title, "nobody to answer to") {
			told = true
			if !strings.Contains(r.Text, lowerFirst(rank)) {
				t.Errorf("the player is not told what they were: %q", r.Text)
			}
		}
	}
	if !told {
		t.Fatal("the name they worked under ended and the city said nothing to them")
	}
}
