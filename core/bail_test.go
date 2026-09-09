package core

import (
	"strings"
	"testing"
)

// Auditing the actions that had never had their label, gate and effect read
// together. Bail is offered at the precinct for anybody of the player's who is
// inside.
//
// Two things. `BailReadiness` exists, refuses when the player cannot afford it,
// and is called by `Bail` — but the button was added with an empty readiness,
// so it was never shown as unavailable. A player with no money saw a live
// button and got an error back instead of a reason on the button, which is the
// opposite of how every other action in the game behaves.
//
// And the description counts days while the ledger entry beside it does not:
// the effect already uses `plainly` to write "for the day still on them", and
// the description said "for the 1 days still on them".

func heldMan(t *testing.T, days int) (*World, *NPC) {
	t.Helper()
	w := New(23)
	w.Player.Location = "precinct"
	w.Player.Cash = 5000
	w.Player.Respect = 120
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	w.Incorporate()
	if w.PlayerOrganization() == nil {
		t.Fatal("the player has no organization, so nobody answers to them")
	}
	w.NPCs = append(w.NPCs, NPC{ID: "held", Name: "Otto Reiss", Location: "precinct",
		Faction: w.PlayerOrganizationID(), Rank: RankSoldier,
		Held: w.Minute + days*1440})
	n := w.NPC("held")
	if !w.Inside(n) {
		t.Fatal("he is not inside")
	}
	return w, n
}

func bailAction(t *testing.T, w *World) Action {
	t.Helper()
	for _, a := range w.Actions("precinct") {
		if strings.HasPrefix(a.ID, "bail:") {
			return a
		}
	}
	t.Fatal("the precinct does not offer to bail anybody out")
	return Action{}
}

func TestBailIsRefusedOnTheButtonWhenItCannotBePaid(t *testing.T) {
	w, _ := heldMan(t, 2)
	if a := bailAction(t, w); a.Disabled {
		t.Fatalf("a player who can afford it is refused: %q", a.Reason)
	}
	w.Player.Cash = 0
	a := bailAction(t, w)
	if !a.Disabled {
		t.Fatal("a player with no money is offered a live bail button")
	}
	if !contains(a.Reason, "cash") && !contains(a.Reason, "money") {
		t.Fatalf("the refusal does not say the money is the problem: %q", a.Reason)
	}
	// And the command still refuses, as it always did.
	if err := w.Bail("held"); err == nil {
		t.Fatal("bail was paid with no money")
	}
}

func TestBailCountsOneDayAsOneDay(t *testing.T) {
	w, _ := heldMan(t, 1)
	a := bailAction(t, w)
	if contains(a.Detail, "1 days") {
		t.Errorf("the description counts one day as many: %q", a.Detail)
	}
	if _, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: a.ID, Target: "precinct"}); err != nil {
		t.Fatalf("bail failed: %v", err)
	}
	// The ledger entry has always got this right; the description should agree.
	last := w.History[len(w.History)-1]
	if contains(last.Text, "1 days") {
		t.Errorf("the ledger counts one day as many: %q", last.Text)
	}
}

// The other half, and the more important one: the exception must be exactly one
// action wide. Everything else aimed at somebody in a cell stays refused, which
// is the rule the sweep was written for.
func TestOnlyBailReachesIntoACell(t *testing.T) {
	w, n := heldMan(t, 2)
	w.Player.Crew = []Crew{{n.ID, n.Name, 60}}
	w.Player.Location = "bar"
	n.Location = "bar"
	offered, refused := 0, 0
	for _, a := range w.Actions("bar") {
		if a.Subject != n.ID {
			continue
		}
		offered++
		if a.Disabled && contains(a.Reason, "held at Ward Street Station") {
			refused++
			continue
		}
		t.Errorf("%q is offered for a man the police are holding: disabled=%v %q",
			a.Label, a.Disabled, a.Reason)
	}
	if offered == 0 {
		t.Fatal("nothing was aimed at him at all, so this proved nothing")
	}
	t.Logf("%d of %d actions aimed at him are refused because he is inside", refused, offered)
	// And bail itself, at the precinct, is not.
	w.Player.Location = "precinct"
	if a := bailAction(t, w); a.Disabled {
		t.Fatalf("bail is still refused: %q", a.Reason)
	}
}
