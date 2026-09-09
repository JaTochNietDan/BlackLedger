package core

import "testing"

// Families had income and no outgoings of any kind. Over a year the poorest
// organization in the city held sixty thousand dollars and not one had ever
// been short, so "a family with little money" was a state the world could not
// reach and nothing could be decided from a family's wealth.
//
// A family now pays its people every day and buys its own repairs.

// A family that runs an operation and holds nothing to pay for it goes broke,
// and shrinks toward what it can actually afford.
func TestAFamilyWithNoHoldingsRunsOutOfMoney(t *testing.T) {
	w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
	f := w.faction("vasco")
	for _, id := range w.FamilyHoldings("vasco") {
		w.Properties[id].Owner = ""
	}
	f.Cash, f.Power, f.Peak = 900, 60, 60
	before, worst := f.Power, 0
	for day := 0; day < 30; day++ {
		w.FamilyDay()
		worst = max(worst, f.Short)
	}
	if f.Cash != 0 {
		t.Errorf("a family with no holdings and sixty people still has $%d after a month", f.Cash)
	}
	// Short is read at the end of a month only for what it says now, and an
	// operation that has shrunk to nobody costs nothing to run, so it is not
	// short any more. What matters is that it went through it.
	if worst == 0 {
		t.Error("a family that could pay nobody was never once short")
	}
	if f.Power >= before {
		t.Errorf("nobody was paid for a month and the operation is still %d strong, up from %d", f.Power, before)
	}
}

// The other side of it: a family whose holdings cover the bill is not short and
// does not shrink for want of money.
func TestAFamilyThatCoversItsBillIsNotShort(t *testing.T) {
	w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
	f := w.faction("vasco")
	f.Cash, f.Power, f.Peak = 5000, 20, 20
	for day := 0; day < 30; day++ {
		w.FamilyDay()
	}
	if f.Short != 0 {
		t.Errorf("a small family on two businesses was short %d days running", f.Short)
	}
	if f.Cash <= 5000 {
		t.Errorf("two businesses did not cover twenty people: $%d, down from $5000", f.Cash)
	}
}

// Money is a quantity, not a sign. Nothing may drive a family below nothing.
func TestAFamilyNeverHoldsLessThanNothing(t *testing.T) {
	w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
	f := w.faction("vasco")
	for _, id := range w.FamilyHoldings("vasco") {
		w.Properties[id].Owner = ""
	}
	f.Cash, f.Power, f.Peak = 10, 90, 90
	for day := 0; day < 60; day++ {
		w.FamilyDay()
		if f.Cash < 0 {
			t.Fatalf("day %d: the family holds $%d", day, f.Cash)
		}
	}
}

// Repairs used to happen every morning for free. A family with nothing spare
// watches its property go down instead, which is how a bad year compounds.
func TestABrokeFamilyCannotRepairItsHoldings(t *testing.T) {
	w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
	f := w.faction("vasco")
	held := w.FamilyHoldings("vasco")
	if len(held) == 0 {
		t.Fatal("the family holds nothing to damage")
	}
	prop := w.Properties[held[0]]
	prop.Condition = 30
	// Broke, and running an operation far larger than one damaged business
	// will pay for, so nothing is ever spare.
	f.Cash, f.Power, f.Peak = 0, 90, 90
	for _, id := range held[1:] {
		w.Properties[id].Owner = ""
	}
	w.FamilyDay()
	if prop.Condition != 30 {
		t.Errorf("a family with no money repaired a business from 30 to %d", prop.Condition)
	}

	// The same damage, the same family, with money behind it.
	f.Cash = 20000
	w.FamilyDay()
	if prop.Condition <= 30 {
		t.Errorf("a family with $20000 left its business at %d", prop.Condition)
	}
}

// The player's organization is run by the player. This loop must never pay
// their wages or repair their businesses for them.
func TestThePlayersOwnOrganizationPaysItsOwnWay(t *testing.T) {
	w := New(41)
	w.Player.Cash, w.Player.Respect = 20000, 200
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	w.Incorporate()
	own := w.PlayerOrganization()
	if own == nil {
		t.Fatal("the player has no organization")
	}
	w.Properties["laundry"].Condition = 40
	own.Cash, own.Power = 500, 40
	w.FamilyDay()
	if w.Properties["laundry"].Condition != 40 {
		t.Errorf("the morning repaired the player's business to %d for nothing", w.Properties["laundry"].Condition)
	}
	if own.Cash != 500 || own.Short != 0 {
		t.Errorf("the loop billed the player's own organization: $%d, short %d", own.Cash, own.Short)
	}
}
