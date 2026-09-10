package core

import "testing"

// The Ledger answered "what am I worth and what is this costing me" with two
// figures and sixty rows of history. DailyCost added nine things together and
// returned one number, so a player losing money had to guess which of the nine
// it was.

func TestTheBooksAddUpToWhatTheDayActuallyCosts(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	// A page of zeroes is noise. Somebody with no car, no crew and no
	// arrangements should not read lines about any of them.
	w := proprietor(t)
	for _, l := range w.Books()["lines"].([]Line) {
		if l.Amount == 0 {
			t.Fatalf("%q is listed at nothing", l.Label)
		}
	}
}

// The page whose whole job is "what am I worth and what is this costing me"
// said nothing about the one cost that has already gone wrong. The day's lines
// are what the player owes; being behind on the wages is where they have
// already failed to pay it, and it is the only one with people on the other
// side: a week of it and somebody stops coming in, and nobody new takes the job
// until it is paid.
func TestTheBooksSayWhichPremisesAreBehind(t *testing.T) {
	t.Parallel()
	w, id := unpaidRun(t)
	if behind := w.Books()["behind"].([]map[string]any); len(behind) != 0 {
		t.Fatalf("a laundry bought this morning is already behind on its wages: %v", behind)
	}
	for day := 0; day < 60; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	behind := w.Books()["behind"].([]map[string]any)
	if len(behind) != 1 {
		t.Fatalf("sixty days of paying nobody and the books name %d premises", len(behind))
	}
	row := behind[0]
	if row["id"] != id {
		t.Fatalf("the books are behind on %v rather than on the laundry", row["id"])
	}
	if row["nights"].(int) < PatienceRunsOut {
		t.Fatalf("the books say %v nights against a place nobody will work at", row["nights"])
	}
	if !row["shut"].(bool) {
		t.Fatal("the books do not say that nobody will take the job")
	}
	if row["positions"].(int) == 0 {
		t.Fatal("the books cannot say how many are missing, because they do not say how many it takes")
	}
}

// And a player who pays has nothing there, which is the ordinary case: an empty
// warning above the day's costs would be a page crying wolf every morning.
func TestTheBooksOfAPlayerWhoPaysAreQuiet(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Location, w.Player.Cash = "laundry", 200000
	if err := w.apply(Command{Kind: "acquire", Target: "laundry", RequestID: "buyitandpaythebooks"}); err != nil {
		t.Fatal(err)
	}
	for day := 0; day < 20; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	if behind := w.Books()["behind"].([]map[string]any); len(behind) != 0 {
		t.Fatalf("a player who pays every bill is told they are behind: %v", behind)
	}
}
