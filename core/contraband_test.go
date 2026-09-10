package core

import "testing"

func trader(t *testing.T) *World {
	t.Helper()
	w := New(9)
	w.Player.Location = "market"
	w.Player.Cash = 1000
	return w
}

func TestBuyingAndSellingMovesMoneyAndStock(t *testing.T) {
	w := trader(t)
	price := w.Good("moonshine").Price
	cash := w.Player.Cash
	if err := w.Buy("moonshine", 0); err != nil {
		t.Fatal(err)
	}
	if w.Holding("moonshine") != Lot {
		t.Fatalf("bought a lot but hold %d", w.Holding("moonshine"))
	}
	if w.Player.Cash != cash-price*Lot {
		t.Fatalf("cash is %d, expected %d", w.Player.Cash, cash-price*Lot)
	}
	// Selling into a higher price is where the profit comes from.
	w.Good("moonshine").Price = price * 2
	if err := w.Sell("moonshine", 0); err != nil {
		t.Fatal(err)
	}
	if w.Holding("moonshine") != 0 {
		t.Fatal("stock survived the sale")
	}
	if w.Player.Cash != cash-price*Lot+price*2*Lot {
		t.Fatalf("takings are wrong: cash %d", w.Player.Cash)
	}
}

func TestTradeNeedsAMarketAndTheMeans(t *testing.T) {
	w := trader(t)
	w.Player.Cash = 10
	if w.TradeReadiness("moonshine", "buy", 0) == "" {
		t.Fatal("bought without the cash")
	}
	if err := w.Buy("moonshine", 0); err == nil {
		t.Fatal("the command ignored a requirement the action reports")
	}
	w.Player.Cash = 1000
	if err := w.Sell("moonshine", 0); err == nil {
		t.Fatal("sold goods that were never held")
	}
	// The waterfront deals in what comes off a boat, not in everything.
	w.Player.Location = "docks"
	if w.TradeReadiness("cigarettes", "buy", 0) == "" {
		t.Fatal("the docks traded a good it does not deal in")
	}
	if w.TradeReadiness("moonshine", "buy", 0) != "" {
		t.Fatal("the docks refused the one good it does deal in")
	}
	// Nowhere else is a market at all.
	w.Player.Location = "room"
	if w.TradeReadiness("moonshine", "buy", 0) == "" {
		t.Fatal("a rented room traded contraband")
	}
	if w.Good("nonsense") != nil {
		t.Fatal("an unknown good exists")
	}
}

func TestCarryingGoodsDrawsAttention(t *testing.T) {
	w := trader(t)
	if err := w.Buy("moonshine", 0); err != nil {
		t.Fatal(err)
	}
	heat := w.Player.Heat
	w.ContrabandDay()
	if w.Player.Heat <= heat {
		t.Fatal("carrying contraband drew no attention")
	}
	// Attention stops once the goods are gone.
	if err := w.Sell("moonshine", 0); err != nil {
		t.Fatal(err)
	}
	settled := w.Player.Heat
	w.ContrabandDay()
	if w.Player.Heat != settled {
		t.Fatal("attention accrued with nothing being carried")
	}
	// Cigarettes are the quiet trade.
	quiet := trader(t)
	if err := quiet.Buy("cigarettes", 0); err != nil {
		t.Fatal(err)
	}
	before := quiet.Player.Heat
	quiet.ContrabandDay()
	if quiet.Player.Heat != before {
		t.Fatal("untaxed cigarettes drew police attention on their own")
	}
}

func TestASearchTakesWhatIsBeingCarried(t *testing.T) {
	w := trader(t)
	if err := w.Buy("moonshine", 0); err != nil {
		t.Fatal(err)
	}
	if lost := w.Seize("A search."); lost != Lot {
		t.Fatalf("a search took %d units, expected %d", lost, Lot)
	}
	if w.Carrying() != 0 {
		t.Fatal("stock survived a search")
	}
	if w.Seize("Nothing to find.") != 0 {
		t.Fatal("a search of an empty pocket took something")
	}
}

func TestPricesMoveAndStayWithinReach(t *testing.T) {
	w := New(11)
	moved := false
	start := w.Good("moonshine").Price
	for i := 0; i < 400; i++ {
		w.MarketPrices()
		p := w.Good("moonshine").Price
		if p != start {
			moved = true
		}
		base := w.Good("moonshine").Base
		if p < base/3 || p > base*3 {
			t.Fatalf("price left its bounds: %d against a base of %d", p, base)
		}
	}
	if !moved {
		t.Fatal("the market never moved")
	}
}

func TestWarMakesGoodsDearer(t *testing.T) {
	peace, war := New(12), New(12)
	war.Antagonize("bellandi", "russo", 100)
	if war.Conflict("bellandi", "russo").State != "war" {
		t.Fatal("the test world is not at war")
	}
	peaceTotal, warTotal := 0, 0
	for i := 0; i < 200; i++ {
		peace.MarketPrices()
		war.MarketPrices()
		peaceTotal += peace.Good("moonshine").Price
		warTotal += war.Good("moonshine").Price
	}
	if warTotal <= peaceTotal {
		t.Fatalf("war did not make goods scarcer: %d against %d", warTotal, peaceTotal)
	}
}

func TestOlderSavesCarryNothing(t *testing.T) {
	w := New(13)
	w.Player.Stock = nil
	if w.Carrying() != 0 || w.Holding("moonshine") != 0 {
		t.Fatal("a save with no recorded stock was carrying something")
	}
	if err := w.Sell("moonshine", 0); err == nil {
		t.Fatal("sold from an absent stock record")
	}
}
