package core

import (
	"strings"
	"testing"
)

// The top bar grew to eight figures and never said what any of them meant. A
// new player reads "PRESENCE 253" and has no way to learn what it does short of
// dying of it.

func TestEveryFigureInTheTopBarExplainsItself(t *testing.T) {
	t.Parallel()
	w, _ := testator(t)
	stats := w.Dashboard()
	if len(stats) < 6 {
		t.Fatalf("the top bar carries %d figures", len(stats))
	}
	seen := map[string]bool{}
	for _, s := range stats {
		if s.ID == "" || s.Label == "" || s.Value == "" {
			t.Fatalf("a figure reads %q / %q / %q", s.ID, s.Label, s.Value)
		}
		if len(s.Meaning) < 40 {
			t.Fatalf("%q is explained as %q", s.Label, s.Meaning)
		}
		if seen[s.ID] {
			t.Fatalf("%q is listed twice", s.ID)
		}
		seen[s.ID] = true
	}
	if !seen["cash"] || !seen["heat"] || !seen["city"] {
		t.Fatalf("the top bar is missing something: %v", seen)
	}
}

// The explanations carry the real thresholds, so they cannot drift from the
// rules the way written copy does.
func TestTheExplanationsCarryTheRealNumbers(t *testing.T) {
	t.Parallel()
	w, _ := testator(t)
	for _, s := range w.Dashboard() {
		switch s.ID {
		case "heat":
			if !strings.Contains(s.Meaning, itoa(RaidThreshold)) || !strings.Contains(s.Meaning, itoa(ForfeitThreshold)) {
				t.Fatalf("attention is explained without its thresholds: %q", s.Meaning)
			}
		case "city":
			if !strings.Contains(s.Meaning, itoa(ScrutinyCrackdown)) {
				t.Fatalf("the city is explained without its threshold: %q", s.Meaning)
			}
		case "respect":
			if !strings.Contains(s.Meaning, itoa(OrganizationStanding)) {
				t.Fatalf("respect is explained without what it opens: %q", s.Meaning)
			}
		}
	}
}

func TestAFigureWorthWorryingAboutSaysSo(t *testing.T) {
	t.Parallel()
	w, _ := testator(t)
	w.Player.Heat, w.Player.Health = RaidThreshold+5, 20
	warned := map[string]bool{}
	for _, s := range w.Dashboard() {
		if s.Warn {
			warned[s.ID] = true
		}
	}
	if !warned["heat"] || !warned["health"] {
		t.Fatalf("at %d attention and %d health the bar warns about %v", w.Player.Heat, w.Player.Health, warned)
	}
	w.Player.Heat, w.Player.Health, w.Player.Cash = 0, 100, 100000
	for _, s := range w.Dashboard() {
		if s.Warn && s.ID != "city" {
			t.Fatalf("a healthy solvent player is warned about %q", s.ID)
		}
	}
}

func TestTheCityIsSaidTheSameWayEverywhere(t *testing.T) {
	t.Parallel()
	// The bar and the city description must not disagree about the word.
	w, _ := testator(t)
	for _, level := range []int{0, 40, 70, 95} {
		w.Attention = level
		word := w.ScrutinyDescription()["state"].(string)
		for _, s := range w.Dashboard() {
			if s.ID == "city" && s.Note != word {
				t.Fatalf("at %d the bar says %q and the city says %q", level, s.Note, word)
			}
		}
	}
}

func TestMoneyIsWrittenTheWayPeopleReadIt(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		n    int
		want string
	}{{0, "$0"}, {90, "$90"}, {999, "$999"}, {1000, "$1,000"}, {3778, "$3,778"},
		{1234567, "$1,234,567"}, {-450, "-$450"}} {
		if got := cash(c.n); got != c.want {
			t.Errorf("%d reads %q, expected %q", c.n, got, c.want)
		}
	}
}

// The top bar is the one thing always on screen, and it is where the player
// looks to know their situation. It said nothing at all about being locked in a
// cell: six figures, identical to a free man's, while the most important fact
// about the next five days was two clicks away on another screen.
func TestTheTopBarSaysWhenThePoliceHaveYou(t *testing.T) {
	t.Parallel()
	w := New(4)
	w.Player.Cash, w.Player.Respect = 6000, OrganizationStanding
	w.Player.Location = "bar"
	for _, s := range w.Dashboard() {
		if s.ID == "held" {
			t.Fatal("a free man is told he is in a cell")
		}
	}
	w.Confine(5, "a still in the back")
	var held Stat
	for _, s := range w.Dashboard() {
		if s.ID == "held" {
			held = s
		}
	}
	if held.ID == "" {
		t.Fatal("the police have him and the top bar does not mention it")
	}
	if !held.Warn {
		t.Fatal("being in a cell is not marked as anything to worry about")
	}
	if !contains(held.Value, itoa(w.DaysLeft())) {
		t.Fatalf("it does not say how long: %q", held.Value)
	}
	if !contains(held.Note, "Ward Street") {
		t.Fatalf("it does not say where: %q", held.Note)
	}
	// The meaning has to be something measured rather than asserted: the city
	// does not stop while you are inside.
	for _, want := range []string{"earning", "paid"} {
		if !contains(held.Meaning, want) {
			t.Fatalf("the meaning does not say what keeps running: %q", held.Meaning)
		}
	}
	// And it is the first thing, because it is the first thing that matters.
	if w.Dashboard()[0].ID != "held" {
		t.Fatalf("the top bar leads with %q", w.Dashboard()[0].ID)
	}
}
