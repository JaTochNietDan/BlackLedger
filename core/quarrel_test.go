package core

import (
	"strings"
	"testing"
)

// A quarrel has an arc — it hardens, it becomes a war, it stops — and the paper
// used to carry only the middle of it. Every war was announced under the same
// headline and no war was ever reported as over, so a reader could not tell one
// war from another and had no reason to think any of them had ended.

func headlinesOfKind(w *World, kind string) []string {
	out := []string{}
	for _, s := range w.News {
		if s.Kind == kind {
			out = append(out, s.Headline)
		}
	}
	return out
}

func TestTwoWarsAreTwoDifferentHeadlines(t *testing.T) {
	w := New(121)
	a := &w.Factions[0]
	b := &w.Factions[1]
	c := Faction{ID: "third", Name: "The Alcaraz People", Power: 50, Peak: 50}
	w.Factions = append(w.Factions, c)

	w.Report("war", upper(a.Name)+" AND "+upper(b.Name)+" AT WAR", "x")
	w.Report("war", upper(a.Name)+" AND "+upper(c.Name)+" AT WAR", "y")
	seen := headlinesOfKind(w, "war")
	if len(seen) != 2 || seen[0] == seen[1] {
		t.Errorf("two different wars produced %v", seen)
	}
}

func TestAWarThatFinishedSomebodyReadsDifferently(t *testing.T) {
	w := New(122)
	a, b := &w.Factions[0], &w.Factions[1]

	// Nobody is beaten: both still hold something.
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil {
			prop.Owner = a.ID
			break
		}
	}
	held := false
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Owner != a.ID {
			prop.Owner = b.ID
			held = true
			break
		}
	}
	if !held {
		t.Skip("not enough premises in this world")
	}
	both := w.howItEnded(a, b)
	if !strings.Contains(both, "stopped short") {
		t.Errorf("a war both sides survived was reported as %q", both)
	}

	// Now strip one side of everything.
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Owner == b.ID {
			prop.Owner = "independent"
		}
	}
	beaten := w.howItEnded(a, b)
	if !strings.Contains(beaten, b.Name) {
		t.Errorf("a war that finished %s did not name them: %q", b.Name, beaten)
	}
	if beaten == both {
		t.Error("burning out and being finished read the same way")
	}
}

func TestTheEndOfAWarIsQuieterThanTheStart(t *testing.T) {
	// The paper says a war has started loudly and says it is over calmly. If
	// the ending carried the same weight, a city would look hardest at the
	// moment the shooting stopped.
	if scrutinyWeight["politics"] >= scrutinyWeight["war"] {
		t.Errorf("the end of a war carries %d against the war's %d",
			scrutinyWeight["politics"], scrutinyWeight["war"])
	}
	// And a quarrel merely hardening costs nothing at all.
	if scrutinyWeight["civic"] != 0 {
		t.Errorf("bad blood between two families raises the city's temperature by %d", scrutinyWeight["civic"])
	}
}

func TestAWarCoolingIntoAFeudIsReportedAsTheFightingStopping(t *testing.T) {
	// A war does not only end by going cold. It usually decays into a feud
	// first, and the first version of this printed "bad blood between them" on
	// the day the shooting stopped — the wrong story at the wrong moment. Only
	// the API run showed it; the tests were all looking at war-to-cold.
	w := New(123)
	a, b := &w.Factions[0], &w.Factions[1]
	w.Conflicts = []Conflict{{A: a.ID, B: b.ID, State: "war", Hostility: warEndsAt + 1}}
	for i := 0; i < 3 && len(headlinesOfKind(w, "politics")) == 0; i++ {
		w.Minute += 1440
		w.FactionTurn()
	}
	stopped := headlinesOfKind(w, "politics")
	if len(stopped) == 0 {
		t.Fatalf("a war decayed out of open fighting and the paper said nothing: state now %q", w.Conflicts[0].State)
	}
	if !strings.Contains(stopped[0], "FIGHTING STOPS") {
		t.Errorf("the end of the fighting was reported as %q", stopped[0])
	}
	for _, h := range headlinesOfKind(w, "civic") {
		if strings.Contains(h, "BAD BLOOD") {
			t.Errorf("a war ending was reported as a quarrel hardening: %q", h)
		}
	}
}
