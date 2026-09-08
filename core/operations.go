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
}

var trades = map[string]Trade{
	"laundry": {
		Hands: 3, Wage: 6, Drain: 5, Restock: 90, RestockAmount: 40, Supplies: "soap and coal",
		Trouble: "A press has broken and the back room is standing idle.",
		Remedy:  "Repair the press", RemedyDetail: "Gets the back room working again.", RemedyCost: 120,
	},
	"garage": {
		Hands: 3, Wage: 8, Drain: 6, Restock: 130, RestockAmount: 40, Supplies: "parts",
		Trouble: "Parts are walking out of the store faster than they are booked in.",
		Remedy:  "Find out who is taking the parts", RemedyDetail: "Stops the losses and puts somebody out of a job.", RemedyCost: 90,
	},
	"casino": {
		Hands: 5, Wage: 11, Drain: 9, Restock: 240, RestockAmount: 45, Supplies: "a float at the tables",
		Trouble: "A dealer is working with somebody on the floor and the tables are losing.",
		Remedy:  "Deal with the dealer", RemedyDetail: "Ends the arrangement, one way or another.", RemedyCost: 160,
	},
}

// TradeOf reports how a place runs, and whether it runs at all.
func TradeOf(id string) (Trade, bool) {
	t, ok := trades[id]
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
