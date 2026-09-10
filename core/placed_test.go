package core

import "testing"

// How a family is placed has to be a fact about the world rather than a phrase
// a screen invents, and it has to tell the four states apart. Cash alone cannot:
// a family holding five thousand against a bill of eight hundred a day is in
// more trouble than one holding two thousand against a bill of ninety.

func TestHowAFamilyIsPlacedTellsTheFourStatesApart(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range []struct {
		name        string
		cash, power int
		holdings    bool
		short       int
		want        string
	}{
		{"paying nobody", 0, 60, false, 3, "cannot pay its people"},
		{"a few days left", 900, 60, false, 0, "struggling"},
		{"a few weeks left", 6000, 60, false, 0, "getting by"},
		{"living within its income", 6000, 20, true, 0, "comfortable"},
	} {
		w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
		f := w.faction("vasco")
		f.Cash, f.Power, f.Peak, f.Short = c.cash, c.power, c.power, c.short
		if !c.holdings {
			for _, id := range w.FamilyHoldings("vasco") {
				w.Properties[id].Owner = ""
			}
		}
		got := w.HowTheyArePlaced(f)
		if got != c.want {
			t.Errorf("%s: %q, wanted %q (bill $%d a day, income $%d, %d days of cover)",
				c.name, got, c.want, w.FamilyBill(f), w.FamilyIncome(f), w.DaysOfCover(f))
		}
		seen[got] = true
	}
	if len(seen) != 4 {
		t.Errorf("four different situations produced %d descriptions: %v", len(seen), seen)
	}
}

// The two figures behind it must be the ones the morning actually charges,
// or a family can be described as comfortable on a bill it is not paying.
func TestTheBillDescribedIsTheBillCharged(t *testing.T) {
	w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
	f := w.faction("vasco")
	f.Cash, f.Power, f.Peak = 50000, 40, 40
	for _, id := range w.FamilyHoldings("vasco") {
		w.Properties[id].Owner = ""
	}
	bill := w.FamilyBill(f)
	before := f.Cash
	w.FamilyDay()
	if spent := before - f.Cash; spent != bill {
		t.Errorf("the bill is described as $%d and the morning took $%d", bill, spent)
	}
}

// And a family living within its income is not counting days at all.
func TestAFamilyLivingWithinItsIncomeIsNotCountingDays(t *testing.T) {
	w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
	f := w.faction("vasco")
	f.Cash, f.Power, f.Peak = 100, 8, 8
	if w.FamilyIncome(f) <= w.FamilyBill(f) {
		t.Skipf("this family's income of $%d does not cover its bill of $%d", w.FamilyIncome(f), w.FamilyBill(f))
	}
	if got := w.DaysOfCover(f); got != WellCovered {
		t.Errorf("a family earning more than it spends has %d days of cover on $%d", got, f.Cash)
	}
}
