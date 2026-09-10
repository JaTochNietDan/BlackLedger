package core

import "fmt"

// A business was a number that produced money while you were elsewhere. This
// gives each one an inside: people who work there and have to be paid, stock or
// a float that runs down as it trades, and its own kind of trouble that stops
// it earning until somebody deals with it.
//
// The point is that a business rewards attention and punishes neglect, and that
// what goes wrong at a laundry is not what goes wrong at a casino.

// Trade describes how a particular kind of business runs.
type Trade struct {
	// Hands is the staff needed to run at full capacity.
	Hands int
	// Wage is what each of them costs a day.
	Wage int
	// Drain is how much stock or float a day of trading consumes.
	Drain int
	// Restock is what topping it back up costs, and how much it buys.
	Restock, RestockAmount int
	// Supplies is what the business runs on, in words.
	Supplies string
	// Trouble is what goes wrong here, and what fixing it is called.
	Trouble, Remedy, RemedyDetail string
	// RemedyCost is what putting it right costs.
	RemedyCost int
	// Cover is how much police attention a day's books can absorb. It is the
	// one thing a business is worth beyond its takings, and it is a fact about
	// the TRADE rather than the premises: a laundry is where the word comes
	// from, a casino handles more loose cash than anywhere else in the city,
	// and a yard full of trucks explains very little.
	//
	// Every business used to absorb the same fourteen points, and the gate
	// asked what kind of ROOM it was rather than what was run in it, so a
	// casino — the best front there is — could not launder a dollar.
	Cover int
	// Watched is how much this trade slows the city forgetting you. A revue bar
	// with a late licence and a gambling house are watched in a way a laundry
	// is not. It is called Watched rather than Notice because an operating mode
	// already has a Notice and it means something else — what a family notices
	// about your takings.
	//
	// It makes forgetting LESS FREQUENT rather than smaller, and both halves of
	// that are arithmetic rather than taste. Attention fades one point a day
	// and a hard-run business generates three, so anything that ADDS attention
	// climbs without limit — a player who bought a burlesque and did nothing
	// else would lose it inside a month with no way to stop it. And anything
	// that made the daily fade BIGGER would absorb skimming, which the fade is
	// deliberately too slow to do. Cooling every second or fourth night instead
	// leaves the rate right for a clean owner and slower for a watched one.
	Watched int
	// Hides is how many units of contraband this trade keeps out of sight. A
	// yard full of trucks and a cold room are places things sit without being
	// seen; a revue bar is not. It is the third shape of what a business is —
	// not a number it contributes to the books, but something the player can
	// only do because of what they own.
	Hides int
}

var trades = map[string]Trade{
	"laundry": {
		Hands: 3, Wage: 6, Drain: 5, Restock: 90, RestockAmount: 40, Supplies: "soap and coal",
		Trouble: "A press has broken and the back room is standing idle.",
		Remedy:  "Repair the press", RemedyDetail: "Gets the back room working again.", RemedyCost: 120, Cover: 18, Watched: 0, Hides: 2,
	},
	"garage": {
		Hands: 3, Wage: 8, Drain: 6, Restock: 130, RestockAmount: 40, Supplies: "parts",
		Trouble: "Parts are walking out of the store faster than they are booked in.",
		Remedy:  "Find out who is taking the parts", RemedyDetail: "Stops the losses and puts somebody out of a job.", RemedyCost: 90, Cover: 9, Watched: 0, Hides: 3,
	},
	"casino": {
		Hands: 5, Wage: 11, Drain: 9, Restock: 240, RestockAmount: 45, Supplies: "a float at the tables",
		Trouble: "A dealer is working with somebody on the floor and the tables are losing.",
		Remedy:  "Deal with the dealer", RemedyDetail: "Ends the arrangement, one way or another.", RemedyCost: 160, Cover: 20, Watched: 3, Hides: 0,
	},
	// What goes wrong at a laundry is not what goes wrong at a casino, and it
	// is not what goes wrong at a restaurant either. Each of these has its own
	// trouble, because a business the player has learned to run is only
	// interesting while there is a kind of trouble they have not met.
	"restaurant": {
		Hands: 4, Wage: 7, Drain: 7, Restock: 150, RestockAmount: 45, Supplies: "the week's food order",
		Trouble: "The kitchen failed an inspection and the dining room is half empty.",
		Remedy:  "Put the kitchen right", RemedyDetail: "New fittings and a word with the inspector.", RemedyCost: 170, Cover: 16, Watched: 1, Hides: 1,
	},
	"poolhall": {
		Hands: 2, Wage: 5, Drain: 4, Restock: 70, RestockAmount: 40, Supplies: "cloth, chalk and drink",
		Trouble: "Somebody is running their own book out of the back and taking the room's money with it.",
		Remedy:  "Put the outside book out", RemedyDetail: "The room takes its own bets again.", RemedyCost: 110, Cover: 11, Watched: 2, Hides: 0,
	},
	"butcher": {
		Hands: 3, Wage: 9, Drain: 8, Restock: 200, RestockAmount: 45, Supplies: "stock and ice",
		Trouble: "The cold room failed overnight and a week of stock went with it.",
		Remedy:  "Get the cold room running", RemedyDetail: "An engineer, and the stock replaced.", RemedyCost: 210, Cover: 8, Watched: 0, Hides: 4,
	},
	"haulage": {
		Hands: 6, Wage: 12, Drain: 11, Restock: 320, RestockAmount: 50, Supplies: "fuel and parts",
		Trouble: "A driver has been talking to somebody at Ward Street and the yard knows it.",
		Remedy:  "Find out which driver", RemedyDetail: "One driver off the books, and the runs are quiet again.", RemedyCost: 240, Cover: 6, Watched: 1, Hides: 7,
	},
	"burlesque": {
		Hands: 7, Wage: 10, Drain: 9, Restock: 260, RestockAmount: 45, Supplies: "the bar and the wardrobe",
		Trouble:      "Somebody from outside is leaning on the dancers for a cut of what they take.",
		Remedy:       "Have a word with whoever is standing at the stage door",
		RemedyDetail: "The cut stops and the room keeps its own money.", RemedyCost: 190, Cover: 14, Watched: 3, Hides: 0,
	},
	"scrapyard": {
		Hands: 3, Wage: 8, Drain: 5, Restock: 150, RestockAmount: 40, Supplies: "the torch and the crane",
		Trouble: "A car came in that somebody is still looking for, and it is halfway down the stack.",
		Remedy:  "Make that one disappear properly", RemedyDetail: "Cut up, weighed in, and off the books.", RemedyCost: 180,
		Cover: 10, Watched: 1, Hides: 6,
	},
	"pawn": {
		Hands: 2, Wage: 9, Drain: 6, Restock: 180, RestockAmount: 40, Supplies: "the float behind the counter",
		Trouble: "Somebody came in looking for a watch they say was theirs, and it is in the window.",
		Remedy:  "Settle it over the counter", RemedyDetail: "The watch goes back and nobody writes anything down.", RemedyCost: 170,
		Cover: 11, Watched: 2, Hides: 5,
	},
	"dealer": {
		Hands: 4, Wage: 11, Drain: 6, Restock: 340, RestockAmount: 45, Supplies: "cars on the lot",
		Trouble: "Two cars on the forecourt turn out to have come off a boat, and somebody official has noticed.",
		Remedy:  "Get the paperwork straight", RemedyDetail: "New documents, and the pair of them off the lot.", RemedyCost: 260,
		Cover: 12, Watched: 1, Hides: 5,
	},
	"filling": {
		Hands: 3, Wage: 8, Drain: 7, Restock: 210, RestockAmount: 45, Supplies: "petrol and the rack behind the counter",
		Trouble: "The tanker did not come, and the pumps are running on what is in the ground.",
		Remedy:  "Pay somebody to bring a load out of hours", RemedyDetail: "A tanker at four in the morning, and nobody writes it down.", RemedyCost: 190,
		Cover: 9, Watched: 1, Hides: 4,
	},
	"cabs": {
		Hands: 8, Wage: 9, Drain: 10, Restock: 280, RestockAmount: 50, Supplies: "fuel and tyres",
		Trouble: "Two cars are off the road and the dispatcher is turning work away.",
		Remedy:  "Get the cars back on the road", RemedyDetail: "Both back out by the evening shift.", RemedyCost: 200, Cover: 7, Watched: 1, Hides: 4,
	},
}

// TradeOf reports how a place runs, and whether it runs at all. It asks the
// address what kind of business it is and looks the rules up by that, so a city
// with three laundries in it has three laundries and not three sets of rules
// that happen to agree today.
func TradeOf(id string) (Trade, bool) {
	place, ok := PlaceByID(id)
	if !ok || place.Kind == "" {
		return Trade{}, false
	}
	t, ok := trades[place.Kind]
	return t, ok
}

// TradeOfKind reports the rules for a kind of business directly, for anything
// asking about the kind rather than about an address.
func TradeOfKind(kind string) (Trade, bool) {
	t, ok := trades[kind]
	return t, ok
}

// Capacity is the share of its potential a business is currently working at.
// Understaffed, out of supplies or in trouble, it earns less; none of these
// stops it dead, because a business limping is more interesting than one shut.
func (w *World) Capacity(id string) float64 {
	trade, ok := TradeOf(id)
	prop := w.Properties[id]
	if !ok || prop == nil {
		return 1
	}
	capacity := 1.0
	if trade.Hands > 0 {
		staffed := float64(min(prop.Staff, trade.Hands)) / float64(trade.Hands)
		capacity *= .35 + .65*staffed
	}
	if prop.Supply <= 0 {
		capacity *= .45
	}
	if prop.Trouble {
		capacity *= .55
	}
	return capacity
}

// Wages is what the people working the player's businesses cost each day.
func (w *World) Wages() int {
	total := 0
	for id, prop := range w.Properties {
		if !w.Own(id) {
			continue
		}
		if trade, ok := TradeOf(id); ok {
			total += prop.Staff * trade.Wage
		}
	}
	return total
}

// OperationsDay runs each owned business for a day: stock and float run down
// with trading, and trouble finds the places that are not being watched.
func (w *World) OperationsDay() {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		trade, ok := TradeOf(l.ID)
		if !ok || prop == nil || !w.Own(l.ID) {
			continue
		}
		if prop.Supply > 0 {
			drain := max(1, trade.Drain*max(1, prop.Staff)/max(1, trade.Hands))
			prop.Supply = max(0, prop.Supply-drain)
			if prop.Supply == 0 {
				w.Log("Out of "+trade.Supplies+" at "+l.Name,
					fmt.Sprintf("%s has run out of %s and is barely trading. It needs restocking.", l.Name, trade.Supplies), "business")
			}
		}
		// Trouble arrives where nobody is paying attention: a short-handed or
		// run-down business is the one it happens to.
		if prop.Trouble {
			continue
		}
		risk := .02
		if prop.Staff < trade.Hands {
			risk += .04
		}
		if prop.Condition < 60 {
			risk += .03
		}
		if operatingMode(prop.Mode).ID == "hard" {
			risk += .03
		}
		if w.WorldRandom() < risk {
			prop.Trouble = true
			w.Log("Trouble at "+l.Name, trade.Trouble+" It will keep earning less until somebody deals with it.", "business")
		}
	}
}

// HireReadiness explains why nobody else can be taken on, or returns "".
func (w *World) HireReadiness(id string) string {
	trade, ok := TradeOf(id)
	if !ok || !w.Own(id) {
		return "This is not a business of yours"
	}
	if w.Properties[id].Staff >= trade.Hands {
		return "It is fully staffed"
	}
	if w.Player.Cash < trade.Wage*7 {
		return "Not enough cash to cover their first week"
	}
	return ""
}

// Hire takes somebody on. They cost a week up front and a wage every day after.
func (w *World) Hire(id string) error {
	if reason := w.HireReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	trade, _ := TradeOf(id)
	if err := w.Pay(trade.Wage * 7); err != nil {
		return err
	}
	prop := w.Properties[id]
	prop.Staff++
	place, _ := PlaceByID(id)
	w.Log("Another pair of hands at "+place.Name,
		fmt.Sprintf("%d of %d positions filled. Wages are now $%d a day across your businesses.", prop.Staff, trade.Hands, w.Wages()), "business")
	return nil
}

// LayOffReadiness explains why nobody can be let go, or returns "".
func (w *World) LayOffReadiness(id string) string {
	if _, ok := TradeOf(id); !ok || !w.Own(id) {
		return "This is not a business of yours"
	}
	if w.Properties[id].Staff == 0 {
		return "There is nobody to let go"
	}
	return ""
}

// LayOff cuts the wage bill at the cost of what the place can handle.
func (w *World) LayOff(id string) error {
	if reason := w.LayOffReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	trade, _ := TradeOf(id)
	prop := w.Properties[id]
	prop.Staff--
	place, _ := PlaceByID(id)
	w.Log("Somebody is let go at "+place.Name,
		fmt.Sprintf("%d of %d positions filled. It handles less, and wages are now $%d a day.", prop.Staff, trade.Hands, w.Wages()), "business")
	return nil
}

// RestockReadiness explains why it cannot be topped up, or returns "".
func (w *World) RestockReadiness(id string) string {
	trade, ok := TradeOf(id)
	if !ok || !w.Own(id) {
		return "This is not a business of yours"
	}
	if w.Properties[id].Supply >= trade.RestockAmount {
		return "It has all it needs"
	}
	if w.Player.Cash < trade.Restock {
		return "Not enough cash"
	}
	return ""
}

// Restock puts stock or a float back in.
func (w *World) Restock(id string) error {
	if reason := w.RestockReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	trade, _ := TradeOf(id)
	if err := w.Pay(trade.Restock); err != nil {
		return err
	}
	prop := w.Properties[id]
	prop.Supply = trade.RestockAmount
	place, _ := PlaceByID(id)
	w.Log("Restocked at "+place.Name, fmt.Sprintf("$%d on %s. %s is trading properly again.", trade.Restock, trade.Supplies, place.Name), "business")
	return nil
}

// RemedyReadiness explains why the trouble cannot be dealt with, or returns "".
func (w *World) RemedyReadiness(id string) string {
	trade, ok := TradeOf(id)
	if !ok || !w.Own(id) {
		return "This is not a business of yours"
	}
	if !w.Properties[id].Trouble {
		return "Nothing is wrong here"
	}
	if w.Player.Cash < trade.RemedyCost {
		return "Not enough cash"
	}
	return ""
}

// Remedy deals with whatever has gone wrong.
func (w *World) Remedy(id string) error {
	if reason := w.RemedyReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	trade, _ := TradeOf(id)
	if err := w.Pay(trade.RemedyCost); err != nil {
		return err
	}
	w.Properties[id].Trouble = false
	place, _ := PlaceByID(id)
	w.Log("Sorted at "+place.Name, fmt.Sprintf("$%d and it is dealt with. %s is earning what it should again.", trade.RemedyCost, place.Name), "business")
	return nil
}

// OwnsKind reports whether the player owns any business of this kind. Rules
// about what owning a garage does for you were written when the city could hold
// exactly one, so they named the address; a second garage made them wrong.
func (w *World) OwnsKind(kind string) bool {
	for _, l := range Locations {
		if l.Kind == kind && w.Own(l.ID) {
			return true
		}
	}
	return false
}
