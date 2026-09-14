package core

import (
	"strings"
	"testing"
)

// An agreement does not outlive the organization it was with.
//
// When a family ends, its quarrels and its plans are cleared and its people are
// put on the street. What the player had arranged with it was left where it
// was: the books listed an understanding with "An unidentified family", because
// the name resolved to nothing and the ledger asked for it anyway, and a
// business ceasefire with nobody sat in the save for the rest of the campaign.
//
// Reachable now in a way it was not — a city left alone buries three
// organizations in a hundred long campaigns where the sample this project used
// to measure through buried none. Fourteen of these were left dangling across
// forty cities of a hundred and twenty days.
func TestAnAgreementDoesNotOutliveTheFamily(t *testing.T) {
	t.Parallel()
	kept, dropped := 0, 0
	told := 0
	for n := uint32(1); n <= 40; n++ {
		w := New(spread(n))
		w.Event = nil
		w.Player.Cash, w.Player.Respect = 40000, 200
		w.BusinessTruces = map[string]int{}
		for _, f := range w.Factions {
			if f.ID == w.PlayerOrganizationID() {
				continue
			}
			w.BusinessTruces[f.ID] = w.Minute + 200*1440
			w.Pacts = append(w.Pacts, Pact{With: f.ID, Since: w.Minute, Life: w.Life})
		}
		before := len(w.Pacts)
		// History retains only 180 entries. Observe notifications while time
		// passes, rather than requiring an early notice to survive 120 days.
		end := w.Minute + 120*1440
		notices := map[string]bool{}
		for w.Minute < end && w.Player.Alive && w.Event == nil {
			w.Advance(min(1440, end-w.Minute))
			for _, record := range w.History {
				if strings.Contains(record.Title, "understanding with nobody") {
					notices[record.ID] = true
				}
			}
		}
		told += len(notices)

		alive := map[string]bool{}
		for _, f := range w.Factions {
			alive[f.ID] = true
		}
		for _, p := range w.Pacts {
			if !alive[p.With] {
				t.Errorf("an understanding with %q outlived the organization", p.With)
			}
		}
		for id := range w.BusinessTruces {
			if !alive[id] {
				t.Errorf("a ceasefire with %q outlived the organization", id)
			}
		}
		kept += len(w.Pacts)
		dropped += before - len(w.Pacts)

	}
	t.Logf("across forty cities: %d understandings kept, %d ended with the family, %d said so",
		kept, dropped, told)
	if dropped == 0 {
		t.Fatal("no organization ended in forty cities of a hundred and twenty days, so this measures nothing")
	}
	if told == 0 {
		t.Fatal("agreements ended and the player was never told")
	}
}
