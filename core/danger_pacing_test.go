package core

import "testing"

func TestRoutineOffersWaitForKnownDangerWithoutRevealingHiddenPlans(t *testing.T) {
	for _, tc := range []struct {
		name       string
		known      bool
		life       int
		deferOffer bool
	}{{"discovered", true, 1, true}, {"undiscovered", false, 1, false}, {"previous life", true, 0, false}} {
		t.Run(tc.name, func(t *testing.T) {
			w := New(27)
			scene, err := w.ValidateProposal(Proposal{Title: "A quiet job", Body: "Deliver a message.", Speaker: "mara", Operation: "courier", Outcome: "Delivered."})
			if err != nil {
				t.Fatal(err)
			}
			w.Offers = []Offer{{Ready: w.Minute, Event: scene}}
			w.Plots = []Plot{{ID: "threat", Life: tc.life, Kind: "hit", Actor: "bellandi", Known: tc.known, Due: w.Minute + 120}}
			before := w.Minute
			w.OfferIfReady()
			if (w.Event == nil) != tc.deferOffer || w.Minute != before {
				t.Fatal("wrong delivery or changed time")
			}
			if tc.deferOffer {
				if len(w.Offers) != 1 || len(w.Arrangements) != 0 {
					t.Fatal("postponed request consumed or remembered as offered")
				}
				w = w.Clone()
				w.Plots = nil
				w.OfferIfReady()
				if w.Event == nil || w.Event.ID != scene.ID || len(w.Offers) != 0 {
					t.Fatal("saved offer did not resume after threat cleared")
				}
			}
		})
	}
}

func TestKnownThreatDoesNotBlockChosenNegotiation(t *testing.T) {
	w := New(27)
	w.Player.Location = "club"
	w.Plots = []Plot{{ID: "threat", Life: 1, Kind: "hit", Actor: "bellandi", Known: true, Due: w.Minute + 240}}
	act(t, &w, "audience", "club")
	if w.Event == nil || w.Event.Kind != "audience" {
		t.Fatal("explicit negotiation blocked")
	}
}
