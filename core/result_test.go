package core

import "testing"

// The moment after a player commits to something, what happened is the most
// important thing on the screen. The result used to carry the records the
// command wrote and nothing about the command itself, so the interface could
// show the city's minutes and never the player's decision.

func TestTheResultSaysWhatWasDoneAndWhatItCost(t *testing.T) {
	t.Parallel()
	w := New(1)
	w.Player.Location = "bar"
	next, err := Execute(w, Command{Revision: w.Revision, RequestID: ID(), Kind: "courier", Target: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	r := next.LastResult
	if r == nil {
		t.Fatal("nothing was reported")
	}
	if r.Action != "Carry a discreet envelope" {
		t.Fatalf("the result calls it %q", r.Action)
	}
	if r.Kind != "courier" {
		t.Fatalf("the result files it under %q", r.Kind)
	}
	if r.Elapsed <= 0 {
		t.Fatalf("%d minutes passed", r.Elapsed)
	}
	if r.Cash != next.Player.Cash-w.Player.Cash {
		t.Fatalf("the result says $%d and the books say $%d", r.Cash, next.Player.Cash-w.Player.Cash)
	}
	if r.Cash <= 0 || r.Respect <= 0 {
		t.Fatalf("carrying an envelope paid $%d and %d respect", r.Cash, r.Respect)
	}
	if len(r.Records) == 0 {
		t.Fatal("nothing was written down")
	}
}

func TestWhatItCostIsMeasuredAcrossTheWholeCommand(t *testing.T) {
	t.Parallel()
	// Not just the price on the button: whatever the day charged while the
	// clock moved is part of what the decision cost, and the player should be
	// told once rather than left to diff two screens.
	w := New(2)
	w.Player.Location = "room"
	w.Player.Cash = 5000
	before := w.Player.Cash
	next, err := Execute(w, Command{Revision: w.Revision, RequestID: ID(), Kind: "rest", Target: "room"})
	if err != nil {
		t.Fatal(err)
	}
	r := next.LastResult
	if r.Cash != next.Player.Cash-before {
		t.Fatalf("the result says $%d and the books moved $%d", r.Cash, next.Player.Cash-before)
	}
	if r.Health != next.Player.Health-w.Player.Health {
		t.Fatal("resting reported the wrong recovery")
	}
	if r.Action == "" {
		t.Fatal("the player is not told what they chose")
	}
}

func TestEveryOrdinaryActionReportsItself(t *testing.T) {
	t.Parallel()
	// Whatever the player presses, the result names it. A blank here is a
	// screen that says something happened and will not say what.
	w := New(3)
	w.Properties["laundry"].Owner = "player:1"
	w.Player.Cash, w.Player.Respect, w.District = 9000, 40, 2
	w.Properties["laundry"].Supply = 0
	w.Properties["laundry"].Staff = 0
	w.Properties["laundry"].Condition = 40
	tried := 0
	for _, id := range []string{"wait", "courier", "contact", "repair", "restock", "hire", "travel"} {
		for _, l := range Locations {
			w.Player.Location = l.ID
			available := false
			for _, a := range w.Actions(l.ID) {
				if a.ID == id && !a.Disabled {
					available = true
				}
			}
			if !available {
				continue
			}
			next, err := Execute(w, Command{Revision: w.Revision, RequestID: ID(), Kind: id, Target: l.ID})
			if err != nil {
				continue
			}
			tried++
			if next.LastResult == nil || next.LastResult.Action == "" {
				t.Errorf("%s at %s reported no action", id, l.ID)
			}
			if next.LastResult.Kind != id {
				t.Errorf("%s at %s came back as %q", id, l.ID, next.LastResult.Kind)
			}
			w = next
			break
		}
	}
	if tried < 4 {
		t.Fatalf("only %d actions could be taken; this is not measuring what it claims to", tried)
	}
	t.Logf("%d ordinary actions taken, every one of them named in its own result", tried)
}

// A new life is not a decision with a price. The deltas are computed across a
// change of protagonist, so a rich man dying and a pauper waking up read as
// though the pauper had just lost eight thousand dollars and every ounce of
// standing they had. Nobody lost anything: they are two different people.
func TestBeginningAgainIsNotReportedAsALoss(t *testing.T) {
	t.Parallel()
	w := New(1)
	w.Player.Cash, w.Player.Respect = 9000, 70
	w.Die("Shot on the steps of the Monarch.")
	next, err := Execute(w, Command{Revision: w.Revision, RequestID: ID(), Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	r := next.LastResult
	if r == nil {
		t.Fatal("nothing was reported")
	}
	if r.Cash != 0 || r.Respect != 0 || r.Heat != 0 || r.Health != 0 {
		t.Fatalf("waking up as somebody else reported $%d, %d respect, %d attention and %d health",
			r.Cash, r.Respect, r.Heat, r.Health)
	}
	if r.Kind != "new_life" || len(r.Records) == 0 {
		t.Fatalf("the arrival was filed as %q with %d records", r.Kind, len(r.Records))
	}
}
