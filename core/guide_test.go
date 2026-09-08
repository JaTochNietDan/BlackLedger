package core

import "testing"

// The Guide was prose written when this game had eight actions, and it rotted:
// it still told the player that autonomous family politics was future work, in
// a build where two families fight their own wars, hold ground, run coups and
// collapse. A guide that misdescribes the game is worse than no guide, so this
// one is the game reporting on itself.

func TestTheGuideAlwaysSaysSomethingTrue(t *testing.T) {
	// On a brand new campaign every step is either open, done, or refused with
	// a reason. A blank is a screen that says nothing.
	w := New(1)
	steps := w.Guide()
	if len(steps) < 8 {
		t.Fatalf("the guide covers %d things", len(steps))
	}
	open := 0
	for _, s := range steps {
		if s.Title == "" || s.What == "" {
			t.Fatalf("a step reads %q / %q", s.Title, s.What)
		}
		if s.Open {
			open++
			if s.Reason != "" {
				t.Fatalf("%q is open and also refused: %q", s.Title, s.Reason)
			}
			continue
		}
		if !s.Done && s.Reason == "" {
			t.Fatalf("%q is closed and will not say why", s.Title)
		}
	}
	if open == 0 {
		t.Fatal("a new arrival is told there is nothing they can do")
	}
}

func TestTheGuideFollowsTheGameRatherThanDescribingIt(t *testing.T) {
	// The proof that it cannot rot: change the world and the guide changes,
	// because it asks the same readiness the buttons ask.
	w := New(2)
	before := map[string]Step{}
	for _, s := range w.Guide() {
		before[s.Title] = s
	}
	if before["Premises of your own"].Done {
		t.Fatal("a new arrival already owns premises")
	}

	rich, _ := testator(t)
	after := map[string]Step{}
	for _, s := range rich.Guide() {
		after[s.Title] = s
	}
	if !after["Premises of your own"].Done {
		t.Fatal("a man with two businesses is still told to buy one")
	}
	if !after["A name of your own"].Done {
		t.Fatal("an incorporated organization is still told to become one")
	}
	if !after["People who answer to you"].Done {
		t.Fatal("a man with people is still told to sign somebody on")
	}
	// And the one thing he has not done yet still says why.
	if step := after["Somebody on the door"]; step.Done {
		t.Fatal("nobody was put on a door and the guide says otherwise")
	}
	rich.Player.Location = "laundry"
	if err := rich.Post("laundry"); err != nil {
		t.Fatal(err)
	}
	done := false
	for _, s := range rich.Guide() {
		if s.Title == "Somebody on the door" {
			done = s.Done
		}
	}
	if !done {
		t.Fatal("somebody is standing on a door and the guide has not noticed")
	}
}

func TestTheRulesAreShortAndTrue(t *testing.T) {
	rules := GuideRules()
	if len(rules) < 4 {
		t.Fatalf("%d rules", len(rules))
	}
	for _, r := range rules {
		if len(r) < 30 {
			t.Fatalf("a rule reads %q", r)
		}
	}
}
