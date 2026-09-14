package core

import (
	"fmt"
	"strings"
	"testing"
)

// The Guide was prose written when this game had eight actions, and it rotted:
// it still told the player that autonomous family politics was future work, in
// a build where two families fight their own wars, hold ground, run coups and
// collapse. A guide that misdescribes the game is worse than no guide, so this
// one is the game reporting on itself.

func TestTheGuideAlwaysSaysSomethingTrue(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

// The rules are the one part of the guide that is written down once rather
// than answered by the game, so they are the one part that can rot silently.
// Every number in them is a constant somewhere; if the constant moves and the
// sentence does not, the guide is lying to a new player about the only
// thresholds this game promises are public.
func TestTheRulesQuoteNumbersTheGameStillUses(t *testing.T) {
	t.Parallel()
	joined := strings.Join(GuideRules(), " ")
	for _, quoted := range []struct {
		what  string
		value int
	}{
		{"the attention at which the police come to the door", RaidThreshold},
		{"the attention at which they take the premises", ForfeitThreshold},
	} {
		if !strings.Contains(joined, itoa(quoted.value)) {
			t.Errorf("the rules no longer mention %s (%d): %q", quoted.what, quoted.value, joined)
		}
	}
	// And no number appears in them that is not one of the game's own.
	real := map[string]bool{itoa(RaidThreshold): true, itoa(ForfeitThreshold): true, "90": true}
	for _, r := range GuideRules() {
		for i := 0; i < len(r); i++ {
			if r[i] < '0' || r[i] > '9' {
				continue
			}
			j := i
			for j < len(r) && r[j] >= '0' && r[j] <= '9' {
				j++
			}
			if !real[r[i:j]] {
				t.Errorf("the rules quote %q, which is not a threshold this game uses: %q", r[i:j], r)
			}
			i = j
		}
	}
}

// A new player has to be told that the city keeps its own hours. People walk
// between buildings now, and somebody who reads the rules and then travels
// across town to a room the screen named will find it empty.
func TestTheRulesSayThatPeopleMove(t *testing.T) {
	t.Parallel()
	joined := strings.ToLower(strings.Join(GuideRules(), " "))
	for _, want := range []string{"walk", "street"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("the rules never mention %q, in a game where the city walks around: %q", want, joined)
		}
	}
}

// The first thing a new player is told has to be what the first thing a new
// player does actually pays. The figure used to be a literal in four places —
// the payment, the log line, the button and this sentence — beside an
// unrelated 45 for how long the job takes.
func TestTheFirstStepQuotesTheFirstJob(t *testing.T) {
	t.Parallel()
	w := New(4)
	w.Player.Location = "bar"
	var button Action
	for _, a := range w.Actions("bar") {
		if a.ID == "courier" {
			button = a
		}
	}
	if button.ID == "" {
		t.Fatal("the first job is not offered at the bar")
	}
	opening := w.Guide()[0]
	for _, want := range []string{itoa(CourierPay), itoa(CourierRespect)} {
		if !strings.Contains(opening.What, want) {
			t.Errorf("the guide's opening step does not mention %q: %q", want, opening.What)
		}
		if !strings.Contains(button.Detail, want) {
			t.Errorf("the button does not mention %q: %q", want, button.Detail)
		}
	}
	// And doing it pays what both of them said.
	cash, respect := w.Player.Cash, w.Player.Respect
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "courier", Target: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Cash-cash != CourierPay || next.Player.Respect-respect != CourierRespect {
		t.Fatalf("the guide promises $%d and %d respect; the job paid $%d and %d",
			CourierPay, CourierRespect, next.Player.Cash-cash, next.Player.Respect-respect)
	}
}

// The guide's whole promise is that it asks the game the same question the
// buttons ask, so it cannot tell the player something the rules do not. It
// asked the readiness functions directly, and those answer about a rule rather
// than about the player's situation — so a man locked in a cell at Ward Street
// Station was told he could carry envelopes at Saint Agnes and buy premises,
// while the action list correctly offered him three things: sit it out, pay a
// lawyer, or name somebody.
func TestTheGuideKnowsWhenYouAreInACell(t *testing.T) {
	t.Parallel()
	w := New(4)
	w.Player.Cash, w.Player.Respect = 6000, OrganizationStanding
	w.Player.Location = "bar"
	before := 0
	for _, s := range w.Guide() {
		if s.Open {
			before++
		}
	}
	if before == 0 {
		t.Fatal("nothing was open to a free man with money")
	}
	w.Confine(5, "a still in the back")
	if !w.Held() {
		t.Fatal("he was not confined")
	}
	for _, s := range w.Guide() {
		if s.Open {
			t.Errorf("%q is offered to a man in a cell", s.Title)
		}
		if !s.Done && s.Reason == "" {
			t.Errorf("%q is closed and says nothing about why", s.Title)
		}
	}
	// And it says the one thing that is true, in the words the game uses.
	closed := w.Guide()[0].Reason
	if !contains(closed, "Ward Street") {
		t.Fatalf("the guide does not say where he is: %q", closed)
	}
	// Let him out and put him back where he was — released at the station, a
	// man cannot lend money to somebody who is at the bar, which is the rules
	// working rather than the guide failing.
	w.Player.HeldUntil = 0
	w.Player.Location = "bar"
	open := 0
	for _, s := range w.Guide() {
		if s.Open {
			open++
		}
	}
	if open != before {
		t.Fatalf("released, %d things are open where %d were before", open, before)
	}
}

// And a guide for somebody who is dead is a guide to nothing.
func TestTheGuideKnowsWhenYouAreDead(t *testing.T) {
	t.Parallel()
	w := New(4)
	w.Player.Cash = 6000
	w.Die("Shot on the steps of the Monarch.")
	for _, s := range w.Guide() {
		if s.Open {
			t.Errorf("%q is offered to somebody who is dead", s.Title)
		}
	}
}

// The guide's own header promises it "cannot tell you something the rules do
// not". It could: the standing it takes to buy premises was written out three
// times as a bare 6 — in the rule, in the opportunity that suggests it, and in
// the guide line itself — so one of the three could change and the page would
// go on saying six. Every figure the guide states now comes from the constant
// the rule uses.
func TestTheGuideStatesNoFigureOfItsOwn(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	w.Player.Respect = PremisesRespect - 1
	found := false
	for _, s := range w.Guide() {
		if s.Title != "Premises of your own" {
			continue
		}
		found = true
		if s.Reason != fmt.Sprintf("Earn %d respect first", PremisesRespect) {
			t.Fatalf("the guide says %q while the rule uses %d", s.Reason, PremisesRespect)
		}
	}
	if !found {
		t.Fatal("the guide does not mention premises")
	}
	// And the rule itself agrees, at the one place the player meets it.
	w.Player.Respect = PremisesRespect
	for _, s := range w.Guide() {
		if s.Title == "Premises of your own" && s.Reason != "" {
			t.Fatalf("at exactly %d respect the guide still says %q", PremisesRespect, s.Reason)
		}
	}
}

func TestGuideRecognizesFirstAssociateBeforeIncorporation(t *testing.T) {
	w := New(81)
	w.Event = nil
	w.Player.Respect = PremisesRespect
	w.Player.Cash = 90
	w.Player.Location = "room"
	var step Step
	for _, s := range w.Guide() {
		if s.Title == "People who answer to you" {
			step = s
		}
	}
	if !step.Open || step.Done || step.Reason != "" {
		t.Fatalf("available first driver blocked by organization rule: %+v", step)
	}
	if w.Incorporated() {
		t.Fatal("fixture is already incorporated")
	}
	if w.Player.Location != "room" {
		t.Fatal("guide moved player to hiring address")
	}
	w.Player.Location = "bar"
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "recruit", Target: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if len(next.Player.Crew) != 1 {
		t.Fatal("recruitment did not occur")
	}
	for _, s := range next.Guide() {
		if s.Title == "People who answer to you" && (!s.Done || s.Open) {
			t.Fatalf("hired driver did not complete milestone: %+v", s)
		}
	}
}

func TestGuideHiringUsesActualCashAndAvailableCandidates(t *testing.T) {
	w := New(81)
	w.Event = nil
	w.Player.Respect = PremisesRespect
	w.Player.Cash = 89
	for _, s := range w.Guide() {
		if s.Title == "People who answer to you" && (s.Open || s.Reason != "Not enough cash") {
			t.Fatalf("unaffordable first hire misdescribed: %+v", s)
		}
	}
	w.Player.Cash = 90
	if driver := w.Holder("driver"); driver != nil {
		driver.Dead = true
	}
	for _, s := range w.Guide() {
		if s.Title == "People who answer to you" && s.Open {
			t.Fatalf("guide offered dead driver: %+v", s)
		}
	}
}
