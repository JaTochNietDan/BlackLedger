package core

import (
	"strings"
	"testing"
)

func operator(t *testing.T) *World {
	t.Helper()
	w := New(701)
	w.Properties["laundry"].Owner = "player:1"
	w.Player.Location = "laundry"
	w.Player.Cash = 5000
	return w
}

func TestABusinessComesAsAGoingConcern(t *testing.T) {
	t.Parallel()
	w := New(703)
	// Over the addresses rather than the trade table: a trade belongs to a kind
	// of business now, and "cabs" is not a place anybody can stand in.
	for _, l := range Locations {
		if l.Kind == "" {
			continue
		}
		id := l.ID
		trade, _ := TradeOf(id)
		prop := w.Properties[id]
		if prop.Staff != trade.Hands {
			t.Fatalf("%s starts with %d of %d positions filled", id, prop.Staff, trade.Hands)
		}
		if prop.Supply != trade.RestockAmount {
			t.Fatalf("%s starts with %d supplies", id, prop.Supply)
		}
		if w.Capacity(id) < 1 {
			t.Fatalf("%s starts working at %.2f of its potential", id, w.Capacity(id))
		}
	}
	// Somewhere that does not trade has no inside to manage.
	if _, running := TradeOf("room"); running {
		t.Fatal("a rented room is a business")
	}
	if w.Capacity("room") != 1 {
		t.Fatal("a place with no trade was given a capacity")
	}
}

func TestNeglectCostsCapacityWithoutClosingThePlace(t *testing.T) {
	t.Parallel()
	w := operator(t)
	full := w.Capacity("laundry")
	w.Properties["laundry"].Supply = 0
	empty := w.Capacity("laundry")
	if empty >= full {
		t.Fatal("running out of supplies cost nothing")
	}
	if empty <= 0 {
		t.Fatal("a business out of supplies stopped dead rather than limping")
	}
	w.Properties["laundry"].Supply = 40
	w.Properties["laundry"].Staff = 0
	if w.Capacity("laundry") >= full {
		t.Fatal("losing every worker cost nothing")
	}
	w.Properties["laundry"].Staff = trades["laundry"].Hands
	w.Properties["laundry"].Trouble = true
	if w.Capacity("laundry") >= full {
		t.Fatal("trouble on the premises cost nothing")
	}
}

func TestStaffCostWagesEveryDay(t *testing.T) {
	t.Parallel()
	w := operator(t)
	trade := trades["laundry"]
	if w.Wages() != trade.Hands*trade.Wage {
		t.Fatalf("wages are %d for %d hands at %d", w.Wages(), trade.Hands, trade.Wage)
	}
	before := w.DailyCost()
	if err := w.LayOff("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.DailyCost() >= before {
		t.Fatal("letting somebody go did not reduce what a day costs")
	}
	// Hiring costs a week up front and cannot exceed the places available.
	cash := w.Player.Cash
	if err := w.Hire("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-trade.Wage*7 {
		t.Fatal("hiring did not cost a week up front")
	}
	if w.HireReadiness("laundry") == "" {
		t.Fatal("a fully staffed business was still hiring")
	}
	// And only for a business the player owns.
	other := New(705)
	other.Player.Location = "laundry"
	if other.HireReadiness("laundry") == "" {
		t.Fatal("the player staffed a business they do not own")
	}
}

func TestTradingRunsSuppliesDownAndRestockingPutsThemBack(t *testing.T) {
	t.Parallel()
	w := operator(t)
	start := w.Properties["laundry"].Supply
	w.OperationsDay()
	if w.Properties["laundry"].Supply >= start {
		t.Fatal("a day of trading consumed nothing")
	}
	for day := 0; day < 40; day++ {
		w.OperationsDay()
	}
	if w.Properties["laundry"].Supply != 0 {
		t.Fatalf("supplies never ran out: %d left", w.Properties["laundry"].Supply)
	}
	cash := w.Player.Cash
	if err := w.Restock("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.Properties["laundry"].Supply != trades["laundry"].RestockAmount {
		t.Fatal("restocking did not fill it")
	}
	if w.Player.Cash != cash-trades["laundry"].Restock {
		t.Fatal("restocking was free")
	}
}

func TestTroubleFindsThePlacesNobodyIsWatching(t *testing.T) {
	t.Parallel()
	watched, neglected := 0, 0
	for i := uint32(1); i <= 300; i++ {
		good := New(i * 2654435761)
		good.Properties["laundry"].Owner = "player:1"
		for day := 0; day < 30; day++ {
			good.Properties["laundry"].Supply = 40
			good.Properties["laundry"].Condition = 100
			good.OperationsDay()
		}
		if good.Properties["laundry"].Trouble {
			watched++
		}

		bad := New(i * 2654435761)
		bad.Properties["laundry"].Owner = "player:1"
		bad.Properties["laundry"].Staff = 0
		for day := 0; day < 30; day++ {
			bad.Properties["laundry"].Condition = 40
			bad.OperationsDay()
		}
		if bad.Properties["laundry"].Trouble {
			neglected++
		}
	}
	t.Logf("of 300 businesses over 30 days: trouble at %d well-run, %d neglected", watched, neglected)
	if neglected <= watched {
		t.Fatal("neglect does not attract trouble")
	}
	if watched == 0 {
		t.Fatal("a well-run business never has a problem, which makes the remedy pointless")
	}
}

func TestTroubleCanBeDealtWithAndEachTradeHasItsOwn(t *testing.T) {
	t.Parallel()
	w := operator(t)
	if w.RemedyReadiness("laundry") == "" {
		t.Fatal("a remedy was offered with nothing wrong")
	}
	w.Properties["laundry"].Trouble = true
	cash := w.Player.Cash
	if err := w.Remedy("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.Properties["laundry"].Trouble {
		t.Fatal("the trouble was not dealt with")
	}
	if w.Player.Cash >= cash {
		t.Fatal("dealing with it was free")
	}
	// What goes wrong at a laundry is not what goes wrong at a casino.
	seen := map[string]bool{}
	for id, trade := range trades {
		if trade.Trouble == "" || trade.Remedy == "" || trade.Supplies == "" {
			t.Fatalf("%s has no trouble, remedy or supplies of its own", id)
		}
		if seen[trade.Trouble] {
			t.Fatalf("%s shares its trouble with another trade", id)
		}
		seen[trade.Trouble] = true
	}
}

// Businesses that link together. The garage was the only address in the city
// that did anything for the rest of what you hold — half off the car's upkeep
// and its repairs — and the user's standing words ask for businesses that link
// to each other rather than eight separate incomes. Every business in the city
// buys its stock from somebody and pays somebody to bring it; a player who
// holds the haulier is paying themselves for the second half of that.
func yardTest(t *testing.T, buyYard bool) *World {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Cash = 500000
	for _, id := range []string{"laundry", "butcher"} {
		w.Event, w.Player.Location = nil, id
		if err := w.apply(Command{Kind: "acquire", Target: id, RequestID: "yard" + id}); err != nil {
			t.Fatal(id, err)
		}
	}
	if buyYard {
		w.Event, w.Player.Location = nil, "haulage"
		if err := w.apply(Command{Kind: "acquire", Target: "haulage", RequestID: "yardyard"}); err != nil {
			t.Fatal(err)
		}
	}
	w.Event = nil
	return w
}

func TestAYardOfYourOwnCarriesYourOwnStock(t *testing.T) {
	t.Parallel()
	without, with := yardTest(t, false), yardTest(t, true)
	for _, id := range []string{"laundry", "butcher"} {
		trade, _ := TradeOf(id)
		full, carried := without.RestockCost(id), with.RestockCost(id)
		t.Logf("%s: $%d the rate, $%d without a yard, $%d with one", id, trade.Restock, full, carried)
		if full != trade.Restock {
			t.Fatalf("%s costs $%d to stock with no yard and the trade says $%d", id, full, trade.Restock)
		}
		if carried >= full {
			t.Fatalf("%s costs $%d with a yard of your own and $%d without", id, carried, full)
		}
	}
	// And the yard cannot carry its own fuel for nothing.
	trade, _ := TradeOf("haulage")
	if with.RestockCost("haulage") != trade.Restock {
		t.Fatalf("the yard stocks itself at $%d instead of $%d",
			with.RestockCost("haulage"), trade.Restock)
	}
}

// What the player pays is what the button says and what the log says, which is
// where a discount like this usually comes apart.
func TestTheYardsSavingIsWhatIsActuallyPaid(t *testing.T) {
	t.Parallel()
	w := yardTest(t, true)
	w.Properties["butcher"].Supply = 0
	before, want := w.Player.Cash, w.RestockCost("butcher")
	trade, _ := TradeOf("butcher")
	if err := w.Restock("butcher"); err != nil {
		t.Fatal(err)
	}
	paid := before - w.Player.Cash
	if paid != want {
		t.Fatalf("the button said $%d and the till took $%d", want, paid)
	}
	if paid >= trade.Restock {
		t.Fatalf("the yard saved nothing: $%d against the rate of $%d", paid, trade.Restock)
	}
	said := w.History[len(w.History)-1].Text
	if !strings.Contains(said, "your own yard") && !strings.Contains(said, "Your own yard") {
		t.Fatalf("the saving is taken and nobody is told where it came from: %q", said)
	}
}

// A player with no yard is told nothing about one, which is the ordinary case.
func TestAPlayerWithNoYardHearsNothingAboutOne(t *testing.T) {
	t.Parallel()
	w := yardTest(t, false)
	w.Properties["butcher"].Supply = 0
	if err := w.Restock("butcher"); err != nil {
		t.Fatal(err)
	}
	if said := w.History[len(w.History)-1].Text; strings.Contains(said, "yard") {
		t.Fatalf("a player who holds no yard is told one carried their stock: %q", said)
	}
}
