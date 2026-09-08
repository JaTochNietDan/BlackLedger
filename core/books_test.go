package core

import "testing"

// The Ledger answered "what am I worth and what is this costing me" with two
// figures and sixty rows of history. DailyCost added nine things together and
// returned one number, so a player losing money had to guess which of the nine
// it was.

func TestTheBooksAddUpToWhatTheDayActuallyCosts(t *testing.T) {
	w, _ := testator(t)
	w.Player.Cash, w.Player.Security = 40000, 2
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 70}}
	w.ensureOfficials()
	w.Player.Location = CityHall
	w.Player.Respect = 90
	if w.RetainerReadiness("commissioner") == "" {
		if err := w.Retain("commissioner"); err != nil {
			t.Fatal(err)
		}
	}

	books := w.Books()
	lines := books["lines"].([]Line)
	if len(lines) < 4 {
		t.Fatalf("the books break the day down into %d lines", len(lines))
	}
	sum := 0
	for _, l := range lines {
		if l.Label == "" || l.Amount == 0 {
			t.Fatalf("a line reads %q at %d", l.Label, l.Amount)
		}
		sum += l.Amount
	}
	if sum != books["costs"].(int) {
		t.Fatalf("the lines add to %d and the books say %d", sum, books["costs"])
	}
	// And the total must agree with the number the clock actually charges.
	if sum != w.DailyCost() {
		t.Fatalf("the books say %d a day and the clock charges %d", sum, w.DailyCost())
	}
	if books["net"].(int) != books["income"].(int)-books["costs"].(int) {
		t.Fatal("the net is not the difference")
	}
}

func TestTheBooksSayWhereTheMoneyIs(t *testing.T) {
	w, _ := testator(t)
	w.Player.Cash, w.District, w.Player.Respect = 40000, 2, 90
	// Money out on the street belongs in the books at what it comes back at.
	lent := 0
	for _, n := range w.Civilians() {
		w.Player.Location = n.Location
		if w.LendReadiness(n.ID) != "" {
			continue
		}
		if err := w.Lend(n.ID); err == nil {
			lent++
		}
		if lent >= 2 {
			break
		}
	}
	if lent == 0 {
		t.Skip("nobody would borrow")
	}
	books := w.Books()
	if books["lent"].(int) <= 0 {
		t.Fatal("money is out on the street and the books do not show it")
	}
	if books["owed"].(int) <= books["lent"].(int) {
		t.Fatalf("the books say $%d went out and $%d comes back", books["lent"], books["owed"])
	}
	if books["cash"].(int) != w.Player.Cash {
		t.Fatal("the books disagree with the pocket")
	}
	if books["holdings"].(int) == 0 {
		t.Fatal("a man with two businesses is shown holding nothing")
	}
}

func TestNothingIsListedThatCostsNothing(t *testing.T) {
	// A page of zeroes is noise. Somebody with no car, no crew and no
	// arrangements should not read lines about any of them.
	w := proprietor(t)
	for _, l := range w.Books()["lines"].([]Line) {
		if l.Amount == 0 {
			t.Fatalf("%q is listed at nothing", l.Label)
		}
	}
}
