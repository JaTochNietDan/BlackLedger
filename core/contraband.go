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

// bulkGoods is how each good reads after a count of units. The market lists
// "Crated arms", which is right on a price board and wrong in a sentence: the
// ledger read "5 crates of Crated arms for $1100". Like pluralGoods this is a
// fact about the word and is deliberately not stored in the save — the first
// version of it was a field on Good, and every campaign begun before tonight
// went on reading "5 crates of crated arms" because the saved struct had no
// such field. That is the second time this exact mistake has been made.
var bulkGoods = map[string]string{
	"moonshine":  "moonshine",
	"cigarettes": "untaxed cigarettes",
	"arms":       "arms",
}

// InBulk is how the good reads after a count of units: "5 crates of arms".
func (g Good) InBulk() string {
	if b, ok := bulkGoods[g.ID]; ok {
		return b
	}
	return lowerFirst(g.Name)
}

// Agrees picks the verb form that goes with this good's name.
func (g Good) Agrees(singular, plural string) string {
	if pluralGoods[g.ID] {
		return plural
	}
	return singular
}

// Lot is what a trade moves when nobody names a number. It was how much changed
// hands full stop, because the interface had no field to type into — which made
// the decision the whole underground trade is built around, how much do you
// dare carry, a constant somebody else chose. A card can carry a figure now.
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

// The route.
//
// There was one price for the whole city, so buying at the docks and selling at
// the market was the same transaction done twice. Measured: a policy that
// trades made 152 purchases and 26 sales over sixty campaigns, because there
// was nothing to carry anything towards. Layer 6 of docs/LIVING_WORLD.md asks
// for "sources, routes, and buyers"; what existed was a price that moved and a
// pocket that held.
//
// A floor is dearer or cheaper than the city's price by a fact about the floor.
// The waterfront is where it comes ashore, so it is cheap there and dear on the
// exchange floor where the buyers are — and the difference is paid for by the
// walk between them, which is the one stretch of this city where nothing covers
// anybody and the reason a car is worth plating.
var floorSpread = map[string]map[string]int{
	"docks": {
		// Off a boat, in quantity, from somebody who wants it gone tonight.
		"moonshine": 78, "arms": 82,
	},
	"market": {
		// The floor where the buyers are, and where a case of anything has
		// already changed hands twice before it gets there.
		"moonshine": 118, "cigarettes": 104,
	},
}

// PriceAt is what one unit costs where the player is standing: the city's own
// price, moved by what this floor is. Zero where nothing is traded.
func (w *World) PriceAt(location, good string) int {
	if !TradesAt(location, good) {
		return 0
	}
	g := w.Good(good)
	if g == nil {
		return 0
	}
	spread := 100
	if floor, ok := floorSpread[location]; ok {
		if at, ok := floor[good]; ok {
			spread = at
		}
	}
	return max(1, g.Price*spread/100)
}

// otherFloor is the other address in this city that deals in the same good, so
// a card can quote the price at the far end of the route rather than leaving
// the player to walk over and find out.
func (w *World) otherFloor(from, good string) string {
	for _, l := range Locations {
		if l.ID != from && TradesAt(l.ID, good) {
			return l.ID
		}
	}
	return ""
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
func (w *World) TradeReadiness(good, side string, units int) string {
	g := w.Good(good)
	if g == nil {
		return "Nobody deals in that here"
	}
	if !TradesAt(w.Player.Location, good) {
		return "There is no market for this here"
	}
	if side == "buy" {
		units = w.BuyUnits(good, units)
		if units <= 0 {
			return "Name how many"
		}
		if room := w.CarryLimit() - w.Carrying(); units > room {
			return fmt.Sprintf("You can carry %s more", counted(max(0, room), g.Unit, g.Unit+"s"))
		}
		if price := w.PriceAt(w.Player.Location, good); w.Player.Cash < price*units {
			return fmt.Sprintf("That is $%d", price*units)
		}
		return ""
	}
	held := w.Holding(good)
	if held == 0 {
		return "You are not carrying any"
	}
	if units > held {
		return fmt.Sprintf("You are carrying %s", counted(held, g.Unit, g.Unit+"s"))
	}
	return ""
}

// BuyUnits is how many a purchase actually moves: what was typed, or the lot
// when nothing was, and never more than there is room or money for.
func (w *World) BuyUnits(good string, units int) int {
	if units > 0 {
		return units
	}
	price := w.PriceAt(w.Player.Location, good)
	if price <= 0 {
		return 0
	}
	room := max(0, w.CarryLimit()-w.Carrying())
	return max(1, min(Lot, min(room, w.Player.Cash/price)))
}

// Buy takes a lot at the current price.
func (w *World) Buy(good string, units int) error {
	if reason := w.TradeReadiness(good, "buy", units); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	g := w.Good(good)
	units = w.BuyUnits(good, units)
	price := w.PriceAt(w.Player.Location, good)
	cost := price * units
	if err := w.Pay(cost); err != nil {
		return err
	}
	if w.Player.Stock == nil {
		w.Player.Stock = map[string]int{}
	}
	w.Player.Stock[good] += units
	w.Player.Heat = min(100, w.Player.Heat+1)
	w.Log("A quiet purchase", fmt.Sprintf("%s of %s for $%d, at $%d each. Holding stock draws attention until it is sold.",
		counted(units, g.Unit, g.Unit+"s"), g.InBulk(), cost, price), "business")
	return nil
}

// Sell moves everything the player is carrying of one good at the current
// price. Selling is where the profit is realised and where the risk ends.
func (w *World) Sell(good string, units int) error {
	if reason := w.TradeReadiness(good, "sell", units); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	g := w.Good(good)
	held := w.Holding(good)
	// Nothing named sells all of it, which is what selling has always meant
	// here: the risk ends when the last of it is gone.
	if units <= 0 || units > held {
		units = held
	}
	price := w.PriceAt(w.Player.Location, good)
	takings := price * units
	w.Player.Stock[good] = held - units
	w.Earn(takings)
	// A load that crossed the city is what gets a man noticed. Selling is the
	// half that shows: buying is somebody with money, selling is somebody with
	// a trade.
	w.RouteRun(units)
	w.Log("The goods move on", fmt.Sprintf("%s of %s sold for $%d, at $%d each.",
		counted(units, g.Unit, g.Unit+"s"), g.InBulk(), takings, price), "business")
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
		w.Log("The goods are gone", fmt.Sprintf("%s You lose %s of stock.", reason,
			plainly(lost, "one unit", fmt.Sprintf("%d units", lost))), "danger")
	}
	return lost
}
