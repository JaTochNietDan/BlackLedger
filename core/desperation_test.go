package core

import "testing"

// The inbox asked for this in as many words: a family with a lot of money is
// stable but a target, and a family with little is struggling, aggressive and
// unpredictable. The four states existed and nothing in the world read them —
// only the director's context did, and the director demonstrably ignores it.
//
// The measurement is a control, because comparing a family to its own past only
// says time passed. Two cities from the same seed, identical in every way
// except how the money sits, run through the same number of days.

// The measurement is the pressure itself, not a month of the city. Money
// changes which branches a family takes — whether it can afford a killing,
// whether a breakaway is worth anything — and every branch changes what is
// drawn from the world's random stream. So two cities that differ only in money
// are not running the same run of luck, and a month of simulation cannot tell
// the two apart. An earlier version of this test measured exactly that and
// reported a clear result in both directions before a line of the mechanism
// existed.

func TestADesperateFamilyPressesAQuarrelHarder(t *testing.T) {
	w := New(41)
	w.District = 2
	a, b := w.faction("bellandi"), w.faction("russo")
	if a == nil || b == nil {
		t.Fatal("the seeded families are missing")
	}
	a.Power, a.Peak, b.Power, b.Peak = 55, 55, 55, 55

	// Everything else held still; only how the money sits changes.
	settle := func(f *Faction, cash, short int) { f.Cash, f.Short = cash, short }

	settle(a, 60000, 0)
	settle(b, 60000, 0)
	if w.HowTheyArePlaced(a) != "comfortable" {
		t.Fatalf("a family with $60000 is %q", w.HowTheyArePlaced(a))
	}
	easy := w.MoneyPressure(a, b)

	settle(a, 0, 4)
	hard := w.MoneyPressure(a, b)
	if hard <= easy {
		t.Errorf("a family that cannot pay anybody presses at %d and a comfortable one at %d", hard, easy)
	}

	settle(a, 900, 0)
	for _, id := range w.FamilyHoldings(a.ID) {
		w.Properties[id].Owner = ""
	}
	if placed := w.HowTheyArePlaced(a); placed != "struggling" {
		t.Fatalf("a family with $900 and no income is %q", placed)
	}
	pressed := w.MoneyPressure(a, b)
	if pressed <= easy || pressed >= hard {
		t.Errorf("struggling presses at %d, between comfortable %d and cannot-pay %d", pressed, easy, hard)
	}
	t.Logf("pressure on a quarrel: comfortable %d, struggling %d, cannot pay %d", easy, pressed, hard)
}

// The other half of the same fact. A neighbour with money is worth moving on —
// but only if there is ground to take, because money nobody can reach is not a
// temptation.
func TestARichNeighbourIsWorthMovingOn(t *testing.T) {
	w := New(41)
	w.District = 2
	a, b := w.faction("bellandi"), w.faction("russo")
	a.Power, a.Peak, b.Power, b.Peak = 55, 55, 55, 55
	a.Cash, a.Short = 6000, 0
	b.Cash, b.Short = 300, 0
	for _, id := range w.FamilyHoldings(b.ID) {
		w.Properties[id].Owner = ""
	}
	poor := w.MoneyPressure(a, b)

	b.Cash = 60000
	if w.HowTheyArePlaced(b) != "comfortable" {
		t.Fatalf("a family with $60000 and no bill to speak of is %q", w.HowTheyArePlaced(b))
	}
	richWithNothing := w.MoneyPressure(a, b)
	if richWithNothing != poor {
		t.Errorf("a rich neighbour holding no ground drew %d against %d — money nobody can reach is not a temptation", richWithNothing, poor)
	}

	w.Properties["market"].Owner = b.ID
	rich := w.MoneyPressure(a, b)
	if rich <= poor {
		t.Errorf("a rich neighbour with premises drew %d and a poor one %d", rich, poor)
	}
	t.Logf("pressure toward a neighbour: poor %d, rich holding nothing %d, rich with premises %d", poor, richWithNothing, rich)
}

// And the wiring: the pressure has to reach the quarrel, not sit in a function
// nobody calls. My first attempt at this called MoneyPressure directly and
// passed with the line in FactionTurn commented out, which is the exact fault
// this project has been caught by before.
//
// Two cities, identical, with the SAME cash on both sides — so nothing that
// reads cash takes a different branch and the random stream is consumed the
// same way — and the world's stream reset to the same value before every turn.
// The only difference is that one pair has missed payday.
func TestTheQuarrelActuallyReadsTheMoney(t *testing.T) {
	run := func(short int) int {
		total := 0
		for _, seed := range []uint32{13, 41, 77, 109, 233, 311} {
			w := New(seed)
			w.District = 2
			a, b := w.faction("bellandi"), w.faction("russo")
			if a == nil || b == nil {
				t.Fatal("the seeded families are missing")
			}
			w.Antagonize(a.ID, b.ID, 30)
			for day := 0; day < 20; day++ {
				a, b = w.faction("bellandi"), w.faction("russo")
				if a == nil || b == nil {
					break
				}
				a.Power, a.Peak, b.Power, b.Peak = 55, 55, 55, 55
				a.Cash, b.Cash = 40000, 40000
				a.Short, b.Short = short, short
				// The same run of luck in both cities, every turn.
				w.WorldRNG = seed*2654435761 + uint32(day)
				w.FactionTurn()
			}
			if c := w.Conflict("bellandi", "russo"); c != nil {
				total += c.Hostility
			}
		}
		return total
	}
	paid, unpaid := run(0), run(5)
	if unpaid <= paid {
		t.Errorf("across six cities with the same money and the same luck, families that had missed payday reached %d hostility and families that had not reached %d — the pressure never reaches the quarrel",
			unpaid, paid)
	}
	t.Logf("same cash, same luck, twenty days: paid %d, missed payday %d", paid, unpaid)
}
