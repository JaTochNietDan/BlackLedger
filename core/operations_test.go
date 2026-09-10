package core

import "testing"

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
