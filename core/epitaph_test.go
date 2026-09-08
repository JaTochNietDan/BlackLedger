package core

import (
	"strings"
	"testing"
)

// The death screen told the player that their properties pass to "your former
// organization" and that another life begins without their money, rank or
// authority. What a protagonist leaves has passed to the strongest of their own
// people for a long time now, as an organization under a name they never chose.
// The most dramatic screen in the game had rotted the way the Guide had.

func TestTheDeadSayWhatBecameOfWhatTheyBuilt(t *testing.T) {
	w, member := testator(t)
	w.Player.Respect, w.Player.Earned = 180, 44000
	name := w.Player.Name
	w.Die("Shot on the way home.")

	if len(w.Dead) == 0 {
		t.Fatal("nobody died")
	}
	last := w.Dead[len(w.Dead)-1]
	if last.Estate == "" {
		t.Fatalf("%s left people behind and the record says nothing became of it", name)
	}
	if !strings.Contains(last.Estate, member.Name) {
		t.Fatalf("the estate is called %q and the man who took it was %s", last.Estate, member.Name)
	}

	e := w.Epitaph()
	if e == nil {
		t.Fatal("a dead protagonist has no epitaph")
	}
	if e["name"] != name || e["cause"] != "Shot on the way home." {
		t.Fatalf("the epitaph reads %v", e)
	}
	if !strings.Contains(e["became"].(string), member.Name) {
		t.Fatalf("what became of it reads %q", e["became"])
	}
	if e["respect"].(int) != 180 || e["earned"].(int) != 44000 {
		t.Fatal("the epitaph disagrees with the life")
	}
	if e["inherits"].(string) == "" {
		t.Fatal("the next person is told nothing about what they get")
	}
}

func TestSomebodyWhoHadNobodyLeavesNothingBehind(t *testing.T) {
	w := proprietor(t)
	own(w, "laundry")
	w.Die("Shot on the way home.")
	e := w.Epitaph()
	if e == nil {
		t.Fatal("no epitaph")
	}
	if e["estate"] != "" {
		t.Fatalf("a man with nobody left %q", e["estate"])
	}
	if !strings.Contains(e["became"].(string), "nobody") {
		t.Fatalf("what became of it reads %q", e["became"])
	}
	// But the premises are still standing, in his name.
	standing := e["standing"].([]string)
	if len(standing) == 0 {
		t.Fatal("a man who owned a laundry left no trace of it")
	}
}

func TestWhatSurvivesTheCityIsSaidPlainly(t *testing.T) {
	w, _ := testator(t)
	w.Offshore = 0
	w.Die("Shot on the way home.")
	if got := w.Epitaph()["inherits"].(string); !strings.Contains(got, "Nothing but the city") {
		t.Fatalf("with nothing abroad the next life is told %q", got)
	}

	w2, _ := testator(t)
	w2.Offshore = 5200
	w2.Die("Shot on the way home.")
	got := w2.Epitaph()["inherits"].(string)
	if !strings.Contains(got, "5200") {
		t.Fatalf("with $5200 abroad the next life is told %q", got)
	}
}
