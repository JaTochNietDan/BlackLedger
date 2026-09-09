package core

import "fmt"

// The underground trade is the fastest money in the city and the least
// forgiving. A price that moves is only half of it: goods have to be carried,
// and carrying them is what draws attention, invites a robbery and turns an
// ordinary police stop into a disaster. Nothing here is a dice roll the player
// cannot see coming; the price, the stock and the attention are all public.

type Good struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Unit  string `json:"unit"`
	Base  int    `json:"base"`
	Price int    `json:"price"`
	// Heat is the attention a single unit draws per day while it is held.
	Heat int `json:"heat"`
}

// pluralGoods are the goods whose names take a plural verb. This is a fact
// about the word, not about the market, so it is not stored in the save: a
// campaign begun before anybody read the paper would otherwise keep printing
// "Crated arms is fetching more than it did" forever.
var pluralGoods = map[string]bool{"cigarettes": true, "arms": true}

// Agrees picks the verb form that goes with this good's name.
func (g Good) Agrees(singular, plural string) string {
	if pluralGoods[g.ID] {
		return plural
	}
	return singular
}

// Lot is how much changes hands in one transaction. The interface offers plain
// actions rather than a quantity field, so trade happens in fixed lots.
const Lot = 5

func newGoods() []Good {
	return []Good{
		{ID: "moonshine", Name: "Moonshine", Unit: "crate", Base: 40, Price: 40, Heat: 1},
		{ID: "cigarettes", Name: "Untaxed cigarettes", Unit: "case", Base: 22, Price: 22, Heat: 0},
		// Guns are the dearest thing that comes off a boat and the worst thing
		// to be found holding. They also have a customer nothing else has: an
		// organization at war.
		{ID: "arms", Name: "Crated arms", Unit: "crate", Base: 165, Price: 165, Heat: 3},
	}
}

// Good finds a tradeable good by id.
func (w *World) Good(id string) *Good {
	for i := range w.Goods {
		if w.Goods[i].ID == id {
			return &w.Goods[i]
		}
	}
	return nil
}

// Holding reports how much of a good the player is carrying.
func (w *World) Holding(id string) int {
	if w.Player.Stock == nil {
		return 0
	}
	return w.Player.Stock[id]
}

// Carrying is everything the player is holding, for attention and risk.
func (w *World) Carrying() int {
	total := 0
	for _, g := range w.Goods {
		total += w.Holding(g.ID)
	}
	return total
}

// TradesAt reports whether a place has a market for a good. The exchange buys
// and sells everything; the waterfront deals in what comes off a boat.
func TradesAt(location, good string) bool {
	switch location {
	case "market":
		// Everything but guns, which nobody moves across a public floor.
		return good != "arms"
	case "docks":
		return good == "moonshine" || good == "arms"
	}
	return false
}

// MarketPrices moves prices toward their base with noise, so a run of good
// prices is temporary and a crash recovers. War disrupts supply: while any two
// organizations are fighting, scarcity pushes prices up. Drawn from the world
// stream, because the market moves whether or not the player is looking.
func (w *World) MarketPrices() {
	if len(w.Goods) == 0 {
		w.Goods = newGoods()
	}
	fighting := false
	for _, c := range w.Conflicts {
		if c.State == "war" {
			fighting = true
		}
	}
	for i := range w.Goods {
		g := &w.Goods[i]
		drift := (g.Base - g.Price) / 4
		swing := int(w.WorldRandom()*float64(g.Base)/2) - g.Base/4
		scarcity := 0
		if fighting {
			scarcity = g.Base / 10
		}
		g.Price = max(g.Base/3, min(g.Base*3, g.Price+drift+swing+scarcity))
	}
}

// ContrabandDay applies what carrying goods costs. Attention accrues for as
// long as the stock is held, which is what makes moving it quickly the point.
func (w *World) ContrabandDay() {
	// Only what is not under the floor of a car draws attention. A false floor
	// is the difference between holding stock and being seen to hold it.
	exposed := w.Exposed()
	if exposed == 0 {
		return
	}
	heat, counted := 0, 0
	for _, g := range w.Goods {
		held := min(w.Holding(g.ID), exposed-counted)
		if held <= 0 {
			continue
		}
		counted += held
		heat += g.Heat * held
	}
	if heat == 0 {
		return
	}
	w.Player.Heat = min(100, w.Player.Heat+min(6, heat))
	w.Log("Goods under the floorboards", fmt.Sprintf("What you are holding is drawing attention. Police interest is now %d.", w.Player.Heat), "danger")
}

// TradeReadiness explains why a trade cannot happen, or returns "".
func (w *World) TradeReadiness(good, side string) string {
	g := w.Good(good)
	if g == nil {
		return "Nobody deals in that here"
	}
	if !TradesAt(w.Player.Location, good) {
		return "There is no market for this here"
	}
	if side == "buy" {
		if w.Player.Cash < g.Price*Lot {
			return "Not enough cash"
		}
		return ""
	}
	if w.Holding(good) == 0 {
		return "You are not carrying any"
	}
	return ""
}

// Buy takes a lot at the current price.
func (w *World) Buy(good string) error {
	if reason := w.TradeReadiness(good, "buy"); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	g := w.Good(good)
	cost := g.Price * Lot
	if err := w.Pay(cost); err != nil {
		return err
	}
	if w.Player.Stock == nil {
		w.Player.Stock = map[string]int{}
	}
	w.Player.Stock[good] += Lot
	w.Player.Heat = min(100, w.Player.Heat+1)
	w.Log("A quiet purchase", fmt.Sprintf("%d %ss of %s for $%d, at $%d each. Holding stock draws attention until it is sold.",
		Lot, g.Unit, g.Name, cost, g.Price), "business")
	return nil
}

// Sell moves everything the player is carrying of one good at the current
// price. Selling is where the profit is realised and where the risk ends.
func (w *World) Sell(good string) error {
	if reason := w.TradeReadiness(good, "sell"); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	g := w.Good(good)
	held := w.Holding(good)
	takings := g.Price * held
	w.Player.Stock[good] = 0
	w.Earn(takings)
	w.Log("The goods move on", fmt.Sprintf("%d %ss of %s sold for $%d, at $%d each.",
		held, g.Unit, g.Name, takings, g.Price), "business")
	return nil
}

// Seize takes contraband out of the player's hands, for a stated reason. The
// same call serves a police search and a robbery.
func (w *World) Seize(reason string) int {
	// A search turns out the premises and the person. What is under the floor
	// of a car parked two streets away is not there to be found.
	remaining := w.Exposed()
	lost := 0
	for _, g := range w.Goods {
		take := min(w.Holding(g.ID), remaining)
		if take <= 0 {
			continue
		}
		remaining -= take
		lost += take
		if w.Player.Stock != nil {
			w.Player.Stock[g.ID] -= take
		}
	}
	if lost > 0 {
		w.Log("The goods are gone", fmt.Sprintf("%s You lose %d units of stock.", reason, lost), "danger")
	}
	return lost
}
