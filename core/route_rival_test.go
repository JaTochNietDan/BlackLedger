package core

import (
	"strings"
	"testing"
)

// "Profit comes from risk taken knowingly: seizure, informants, a rival who
// wants the route."
//
// Measured last slice: of the three, one existed. A smuggler over a hundred
// campaigns carried thirty-seven attention and lost the goods fourteen times,
// and was never once hurt, because nobody in the city had any opinion about a
// man carrying crates through it. Now that the two ends of the route disagree
// about what a crate is worth, running it is a thing worth taking off somebody.

func carrier(t *testing.T) *World {
	t.Helper()
	w := New(83)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 8000, 100
	w.Player.Location = "docks"
	return w
}

func TestRunningTheRouteGetsYouNoticed(t *testing.T) {
	w := carrier(t)
	if w.Player.Runs != 0 {
		t.Fatal("somebody who has never traded is already known for it")
	}
	if err := w.Buy("moonshine", 6); err != nil {
		t.Fatal(err)
	}
	w.Player.Location = "market"
	if err := w.Sell("moonshine", 6); err != nil {
		t.Fatal(err)
	}
	if w.Player.Runs < 6 {
		t.Fatalf("six crates across the city and the word is %d", w.Player.Runs)
	}
}

func TestSomebodyElseWantsTheRoute(t *testing.T) {
	w := carrier(t)
	w.Player.Runs = RouteNotice * 2
	for i := 0; i < 400 && !w.routePlotted(); i++ {
		w.WorldRNG = uint32(i*2654435761 + 7)
		w.ConsiderRoute()
	}
	if !w.routePlotted() {
		t.Fatal("nobody in this city ever wanted a trade somebody else was running")
	}
	// It is aimed at the goods rather than at the man, which is the difference
	// between wanting the route and wanting him dead.
	plot := Plot{}
	for _, p := range w.Plots {
		if p.Kind == "route" {
			plot = p
		}
	}
	if plot.Actor == "" {
		t.Fatal("a plot with nobody behind it")
	}
	w.Player.Location, w.Player.Stock = "transit", map[string]int{"moonshine": 6}
	w.TakeTheRoute(plot)
	if w.Holding("moonshine") != 0 {
		t.Fatalf("they came for the load and left %d of it", w.Holding("moonshine"))
	}
	if !w.Player.Alive {
		t.Fatal("wanting somebody's trade is not wanting them dead")
	}
}

func TestTheyWarnYouOffWhenYouAreCarryingNothing(t *testing.T) {
	w := carrier(t)
	w.Player.Stock = map[string]int{}
	before := len(w.History)
	w.TakeTheRoute(Plot{ID: ID(), Kind: "route", Life: w.Life, Actor: "bellandi", Due: w.Minute})
	if len(w.History) == before {
		t.Fatal("they came for a load that was not there and said nothing about it")
	}
	said := w.History[len(w.History)-1]
	if !strings.Contains(strings.ToLower(said.Title+said.Text), "trade") &&
		!strings.Contains(strings.ToLower(said.Title+said.Text), "carry") {
		t.Fatalf("the warning does not say what it is about: %q %q", said.Title, said.Text)
	}
}
