package core

import "testing"

// Waging war on a family's money rather than its people: take or wreck what it
// earns from, and it cannot pay anybody. This asks whether the city actually
// carries that through, end to end, using only what a player can do.

func TestTakingAFamilysBusinessesStarvesIt(t *testing.T) {
	t.Parallel()
	w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
	f := w.faction("vasco")
	f.Cash, f.Power, f.Peak = 6000, 70, 70
	held := w.FamilyHoldings("vasco")
	if len(held) < 2 {
		t.Fatalf("the family holds %d businesses, so there is nothing to take", len(held))
	}

	// A fortnight left alone: this is the family at its own pace.
	for day := 0; day < 14; day++ {
		w.FamilyDay()
	}
	settled, settledPower := f.Cash, f.Power

	// Now everything they earn from is taken off them, one at a time, the way
	// a player takes a business over.
	for _, id := range held {
		w.Properties[id].Owner = "player:1"
	}
	worst := 0
	for day := 0; day < 45; day++ {
		w.FamilyDay()
		worst = max(worst, f.Short)
	}

	if f.Cash >= settled {
		t.Errorf("a family with everything taken off it holds $%d, up from $%d", f.Cash, settled)
	}
	if worst == 0 {
		t.Error("a family with no income was never short of its wages")
	}
	if f.Power >= settledPower {
		t.Errorf("nobody was paid for six weeks and the operation is %d strong, up from %d", f.Power, settledPower)
	}
	t.Logf("left alone: $%d at %d strong. Everything taken: $%d at %d strong, short on %d days running",
		settled, settledPower, f.Cash, f.Power, worst)
}

// Wrecking what they earn from, rather than taking it, is a slower kind of
// attack and a single act of it is absorbed: the family buys the damage back
// over about a fortnight and is earning normally again. So it is measured
// against the same family over the same weeks, left alone — the control is the
// only thing that says what the wrecking actually cost them.
func TestWreckingAFamilysBusinessesCostsItMoney(t *testing.T) {
	t.Parallel()
	settle := func(wreck bool) (int, int, int) {
		w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
		f := w.faction("vasco")
		f.Cash, f.Power, f.Peak = 6000, 70, 70
		for day := 0; day < 14; day++ {
			w.FamilyDay()
		}
		if wreck {
			for _, id := range w.FamilyHoldings("vasco") {
				w.Properties[id].Condition = 15
			}
		}
		worst := 0
		for day := 0; day < 30; day++ {
			w.FamilyDay()
			worst = max(worst, f.Short)
		}
		return f.Cash, f.Power, worst
	}
	alone, alonePower, _ := settle(false)
	hit, hitPower, worst := settle(true)

	if hit >= alone {
		t.Errorf("wrecking both businesses left the family with $%d against $%d for the same weeks left alone", hit, alone)
	}
	t.Logf("left alone $%d at %d strong; wrecked to 15 once, $%d at %d strong, worst run of %d days short",
		alone, alonePower, hit, hitPower, worst)
}

// And the attack only bites if it is kept up. A family wrecked every week
// cannot both repair and pay, and that is what a campaign against their money
// actually looks like.
func TestWreckingThemEveryWeekStarvesThem(t *testing.T) {
	t.Parallel()
	w := withFamily(41, "vasco", "Vasco Company", "Ilse Vasco")
	f := w.faction("vasco")
	f.Cash, f.Power, f.Peak = 6000, 70, 70
	worst := 0
	for day := 0; day < 60; day++ {
		if day%7 == 0 {
			for _, id := range w.FamilyHoldings("vasco") {
				w.Properties[id].Condition = 15
			}
		}
		w.FamilyDay()
		worst = max(worst, f.Short)
	}
	if worst == 0 {
		t.Error("a family wrecked every week for two months never missed a payday")
	}
	if f.Power >= 70 {
		t.Errorf("a family wrecked every week for two months is still %d strong", f.Power)
	}
	t.Logf("wrecked weekly: $%d at %d strong, worst run of %d days short", f.Cash, f.Power, worst)
}
