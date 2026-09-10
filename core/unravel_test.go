package core

import (
	"strings"
	"testing"
)

// Running a city for four hundred game days turned 154 of its 180 ledger
// entries into one headline. Most of that is the harness — I kept a broke
// player alive far past where a campaign would have ended — but reading the
// entry showed something the harness did not cause.
//
// Every midnight the player cannot cover the bills, the ledger says "Security
// leaves; your residence is now a rented room. Unpaid crew lose loyalty."
// After the first such night there is no security to leave, the residence is
// already a rented room, and the crew's loyalty is already nothing. The
// sentence goes on reporting three losses that already happened.

func brokeAtMidnight(t *testing.T) *World {
	t.Helper()
	w := New(13)
	w.Player.Cash = 0
	w.Player.Home = "apartment"
	w.Player.Security = 2
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 60}}
	return w
}

func midnight(t *testing.T, w *World) Record {
	t.Helper()
	before := len(w.History)
	for i := 0; i < 96 && len(w.History) == before; i++ {
		w.Advance(30)
	}
	if len(w.History) == before {
		t.Fatal("a day passed without the books being done")
	}
	return w.History[len(w.History)-1]
}

func TestTheBillsOnlyReportWhatActuallyWentWrong(t *testing.T) {
	t.Parallel()
	w := brokeAtMidnight(t)
	first := midnight(t, w)
	if !contains(first.Text, "Security leaves") {
		t.Fatalf("the first night does not report losing security: %q", first.Text)
	}
	if w.Player.Security != 0 || w.Player.Home != "room" {
		t.Fatalf("the first night did not actually take them: security %d, home %q",
			w.Player.Security, w.Player.Home)
	}
	// Nothing left to lose. The second night must not claim otherwise.
	w.Player.Crew[0].Loyalty = 0
	second := midnight(t, w)
	for _, claim := range []string{"Security leaves", "now a rented room", "lose loyalty"} {
		if strings.Contains(second.Text, claim) {
			t.Errorf("a second night with nothing left to lose still claims %q:\n  %s", claim, second.Text)
		}
	}
	if strings.TrimSpace(second.Text) == "" {
		t.Fatal("the second night says nothing at all")
	}
	t.Logf("second night reads: %s", second.Text)
}

// Sixth of the class: a refusal that cannot tell two cases apart. Lending is
// capped at a share of everything the player has, out and in hand, so a full
// book and an empty pocket both fall below the floor. The message spoke only of
// the book, which for a player holding nothing read "You have $0 out already,
// which is as much as you can afford to be owed".
func TestBeingBrokeIsRefusedDifferentlyFromHavingLentTooMuch(t *testing.T) {
	t.Parallel()
	setup := func(cash int, out bool) string {
		w := New(17)
		w.Player.Location = "bar"
		w.Player.Respect = LendStanding + 20
		w.Player.Cash = cash
		// Somebody with no title: the fixer and the driver are refused earlier,
		// for a different and correct reason.
		w.NPCs = append(w.NPCs, NPC{ID: "borrower", Name: "Ilona Kovac", Location: "bar"})
		if out {
			w.NPCs = append(w.NPCs, NPC{ID: "debtor", Name: "Otto Reiss", Location: "docks"})
			w.Loans = []Loan{{ID: ID(), Debtor: "debtor", Life: w.Life,
				Principal: cash * 4, Owed: cash * 5, Due: w.Minute + 1440}}
		}
		reason := w.LendReadiness("borrower")
		if reason == "" {
			t.Fatalf("cash %d, out %v: lending was not refused at all", cash, out)
		}
		return reason
	}
	broke := setup(0, false)
	if contains(broke, "$0 out already") {
		t.Errorf("a player with nothing is told about their book: %q", broke)
	}
	if !contains(broke, "money") {
		t.Errorf("the refusal does not say the money is the problem: %q", broke)
	}
	if !contains(broke, "$") {
		t.Errorf("the refusal does not say how much it would take: %q", broke)
	}
	full := setup(400, true)
	if !contains(full, "out already") {
		t.Errorf("a full book is not described as a full book: %q", full)
	}
	if broke == full {
		t.Errorf("both cases are refused with the same words: %q", broke)
	}
	t.Logf("broke: %s", broke)
	t.Logf("full:  %s", full)
}
