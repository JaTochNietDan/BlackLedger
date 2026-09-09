package core

import "testing"

// Adding four businesses to the city crashed every existing save. The repair
// that gives a campaign a record for an address added since it began was gated
// behind `Version < SaveVersion`, and those saves were already at the current
// version — nothing had bumped the number, because adding a place does not
// change the shape of a save, only its contents.
//
// So the repair is not a version migration and must not be gated like one.

func TestASaveAtTheCurrentVersionStillGetsANewAddress(t *testing.T) {
	w := New(31)
	w.Version = SaveVersion
	// A campaign that has never heard of the haulage yard, at the version this
	// build writes.
	delete(w.Properties, "haulage")
	delete(w.Properties, "restaurant")

	w.SettleNewPlaces()

	for _, id := range []string{"haulage", "restaurant"} {
		prop := w.Properties[id]
		if prop == nil {
			t.Fatalf("%s is still missing from a save at the current version", id)
		}
		if prop.Income != PlaceIncome[id] {
			t.Errorf("%s earns %d in an old campaign and %d in a new one", id, prop.Income, PlaceIncome[id])
		}
		if prop.Condition != 100 || prop.Owner != "independent" {
			t.Errorf("%s arrives owned by %q at %d condition", id, prop.Owner, prop.Condition)
		}
		trade, running := TradeOf(id)
		if running && (prop.Staff != trade.Hands || prop.Supply != trade.RestockAmount) {
			t.Errorf("%s arrives with %d of %d hands and %d of %d supplies — a business is a going concern before anybody buys it",
				id, prop.Staff, trade.Hands, prop.Supply, trade.RestockAmount)
		}
	}
}

// The crash itself: reading the world walks every address and would find
// nothing where a property should be.
func TestReadingTheWorldNeverFindsAnAddressWithNoRecord(t *testing.T) {
	w := New(31)
	w.Version = SaveVersion
	for _, l := range Locations {
		delete(w.Properties, l.ID)
	}
	w.SettleNewPlaces()
	// Public walks every location and reads its property. Before the repair ran
	// on every load, this is where the live game died.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("reading the world with a campaign that predates the city: %v", r)
		}
	}()
	if w.Public() == nil {
		t.Fatal("the world read back as nothing")
	}
}

// And what a place earns is one table, so a business added later is not worth
// nothing forever in a campaign that was already running.
func TestEveryEarningPlaceHasItsIncomeInOneTable(t *testing.T) {
	w := New(31)
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop.Income != PlaceIncome[l.ID] {
			t.Errorf("%s earns %d in a new city and the table says %d", l.ID, prop.Income, PlaceIncome[l.ID])
		}
	}
	for id := range trades {
		if PlaceIncome[id] <= 0 {
			t.Errorf("%s is a trading business that earns nothing", id)
		}
	}
}
