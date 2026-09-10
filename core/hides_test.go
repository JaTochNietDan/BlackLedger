package core

import "testing"

// The third shape of what a business is: not a number it contributes, but
// something the player can only do because of what they own. Holding contraband
// draws police attention every day for every unit that is not out of sight, and
// out of sight meant a false floor in a car or a cellar under the house. A yard
// full of trucks and a butcher's cold room are places things sit without being
// seen, and that is what those trades are for.

func TestSomeTradesKeepThingsOutOfSight(t *testing.T) {
	t.Parallel()
	hiding := 0
	for kind, trade := range trades {
		if trade.Hides < 0 {
			t.Errorf("a %s hides %d units, which is not a quantity", kind, trade.Hides)
		}
		if trade.Hides > 0 {
			hiding++
		}
	}
	if hiding == 0 {
		t.Fatal("no trade in the city can keep anything out of sight")
	}
	if hiding == len(trades) {
		t.Error("every trade hides things, so there is nothing to choose between")
	}
	haulage, _ := TradeOfKind("haulage")
	burlesque, _ := TradeOfKind("burlesque")
	if haulage.Hides <= burlesque.Hides {
		t.Errorf("a haulage yard hides %d and a revue bar %d", haulage.Hides, burlesque.Hides)
	}
}

// The wiring: what a business hides has to reach the attention holding stock
// draws, or it is a number in a table nobody reads.
func TestAYardFullOfTrucksHidesWhatACarCannot(t *testing.T) {
	t.Parallel()
	stocked := func(own string) (concealed, exposed, heat int) {
		w := New(53)
		w.District = 2
		w.Player.Cash, w.Player.Respect = 40000, 200
		// Moonshine, and only five crates. Attention from held stock is capped
		// at six a day, and eight crates of arms is twenty-four before the cap
		// — so every case measured the cap rather than the hiding place, and
		// a cab yard and no yard at all both came back as six.
		w.Player.Stock = map[string]int{"moonshine": 5}
		if own != "" {
			w.Properties[own].Owner = "player:1"
		}
		before := w.Player.Heat
		w.ContrabandDay()
		return w.Concealed(), w.Exposed(), w.Player.Heat - before
	}
	bare, _, bareHeat := stocked("")
	yard, yardExposed, yardHeat := stocked("cabstand")
	haul, _, haulHeat := stocked("haulage")

	if yard <= bare || haul <= bare {
		t.Errorf("owning nothing hides %d, a cab yard %d, a haulage yard %d", bare, yard, haul)
	}
	if haul <= yard {
		t.Errorf("a haulage yard hides %d and a cab yard %d", haul, yard)
	}
	if yardHeat >= bareHeat {
		t.Errorf("five crates drew %d attention with a yard and %d without", yardHeat, bareHeat)
	}
	if haulHeat > yardHeat {
		t.Errorf("the bigger yard drew more attention: haulage %d, cabs %d", haulHeat, yardHeat)
	}
	t.Logf("five crates: nothing hides %d and draws %d; a cab yard hides %d (%d still out) and draws %d; haulage hides %d and draws %d",
		bare, bareHeat, yard, yardExposed, yardHeat, haul, haulHeat)
}

// And the address-named ownership check. A car costs half to keep when the
// player owns a garage, and that rule named the ONE garage in the city. It asks
// the kind now.
//
// Being straight about what this proves: there is still only one garage, so the
// garage rule itself cannot be caught failing — swapping the fix back out leaves
// this passing. What can be proved is the thing the fix rests on, over a kind
// the city genuinely has two of.
func TestOwningAnyOneOfAKindCountsAsOwningThatKind(t *testing.T) {
	t.Parallel()
	byKind := map[string][]string{}
	for _, l := range Locations {
		if l.Kind != "" {
			byKind[l.Kind] = append(byKind[l.Kind], l.ID)
		}
	}
	proved := 0
	for kind, ids := range byKind {
		if len(ids) < 2 {
			continue
		}
		proved++
		for _, id := range ids {
			w := New(53)
			if w.OwnsKind(kind) {
				t.Fatalf("a player who owns nothing already owns a %s", kind)
			}
			w.Properties[id].Owner = "player:1"
			if !w.OwnsKind(kind) {
				t.Errorf("owning %s does not count as owning a %s", id, kind)
			}
		}
	}
	if proved == 0 {
		t.Fatal("no kind has two addresses, so this proves nothing")
	}
	// The garage rule itself. Last slice this could not be caught failing,
	// because the city had one garage and "the garage" and "any garage" were
	// the same thing. There are two now, so it is checked on EACH of them
	// ALONE — owning only the second one has to be enough.
	garages := 0
	for _, l := range Locations {
		if l.Kind != "garage" {
			continue
		}
		garages++
		w := New(53)
		w.District = 2
		w.Player.Cash, w.Player.Respect, w.Player.Car = 40000, 200, 1
		full := w.CarUpkeep()
		w.Properties[l.ID].Owner = "player:1"
		if w.CarUpkeep() >= full {
			t.Errorf("owning %s alone costs %d a day against %d owning none", l.ID, w.CarUpkeep(), full)
		}
	}
	if garages < 2 {
		t.Errorf("the city has %d garages, so this is still not provable", garages)
	}
}
