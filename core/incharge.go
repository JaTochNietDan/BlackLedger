package core

import "fmt"

// Putting somebody in charge.
//
// Somebody comes to run a business in this city already: one of the player's
// own people, when they have had enough of being badly paid, walks off with a
// holding and the paper prints it. The player could not do it on purpose. Every
// business they held was a place to walk to and restock by hand, which is the
// work a manager exists to take off somebody.
//
// A manager is one of the hands, not a new person: they stand behind the
// counter as well, they are paid the same wage, and they are somebody a rival
// can take or a family can put off coming in. What the player buys is not
// having to be there.

// RunsIt is whoever the player has put in charge here, and nothing where nobody
// is. Read off the role, which is where the city already keeps this.
func (w *World) RunsIt(id string) *NPC {
	place, ok := PlaceByID(id)
	if !ok {
		return nil
	}
	prop := w.Properties[id]
	if prop == nil {
		return nil
	}
	for _, who := range prop.Hands {
		if n := w.NPC(who); n != nil && !n.Dead && n.Role == "Runs "+place.Name {
			return n
		}
	}
	return nil
}

// InChargeReadiness explains why this person cannot be put in charge, or
// returns "".
func (w *World) InChargeReadiness(id, who string) string {
	n := w.NPC(who)
	if n == nil || n.Dead {
		return "There is nobody here by that name"
	}
	if !w.Own(id) {
		return "This is not a business of yours"
	}
	if w.EmployerOf(who) != id {
		return n.Name + " does not work here"
	}
	if running := w.RunsIt(id); running != nil {
		if running.ID == who {
			return n.Name + " already runs it"
		}
		return running.Name + " runs it"
	}
	return ""
}

// PutInCharge hands somebody the keys. It costs nothing today and a wage every
// day after, which they were already being paid.
func (w *World) PutInCharge(id, who string) error {
	if reason := w.InChargeReadiness(id, who); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	place, _ := PlaceByID(id)
	n := w.NPC(who)
	n.Role = "Runs " + place.Name
	n.Trust = min(100, n.Trust+InChargeTrust)
	w.Log(n.Name+" is running "+place.Name,
		fmt.Sprintf("They keep it stocked out of your money and you do not have to be there. It is a thing to be given, and %s knows it.",
			n.Name), "business")
	return nil
}

// InChargeTrust is what being handed the keys is worth to somebody. Being
// trusted with something is the cheapest loyalty there is.
const InChargeTrust = 20

// TheyRunIt is one day of the people the player has put in charge doing the
// thing they were put there for: keeping the place stocked out of the player's
// money, so a business is not a place to walk to every few days.
func (w *World) TheyRunIt() {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		trade, ok := TradeOf(l.ID)
		if prop == nil || !ok || !w.Own(l.ID) {
			continue
		}
		manager := w.RunsIt(l.ID)
		if manager == nil || w.Inside(manager) || w.CrewOrderFor(manager.ID) != nil || prop.Supply >= trade.RestockAmount {
			continue
		}
		// Use the same funded purchase, haulage discount and ledger receipt
		// as an order or a visit to the counter.
		if w.RestockReadiness(l.ID) == "" {
			_ = w.Restock(l.ID)
		}

	}
}

// RunsItName is who runs this place, for the room to say so. Empty where
// nobody does, and where the player has no business knowing.
func (w *World) RunsItName(id string) string {
	if !w.Own(id) {
		return ""
	}
	if n := w.RunsIt(id); n != nil {
		return n.Name
	}
	return ""
}
