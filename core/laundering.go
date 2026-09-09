package core

import "fmt"

// A laundry is a business that takes cash and gives back paperwork. That is why
// it is worth owning beyond its takings: it is the difference between money the
// police can ask about and money they cannot.
//
// It is not free. The cut is real, the premises take the strain, and a place
// that is falling apart cannot process much of anything.

// LaunderCooldown is how long the books need before they can absorb another
// round without the pattern becoming obvious.
const LaunderCooldown = 1440

// launderCapacity is how much police attention a laundry can clear in one pass,
// scaled by the state of the premises.
func (w *World) launderCapacity(id string) int {
	prop := w.Properties[id]
	trade, runs := TradeOf(id)
	if prop == nil || !runs {
		return 0
	}
	// What the trade is worth as a front, and how well the place is kept. Both
	// matter: a laundry with a caved-in roof explains nothing either.
	return max(0, trade.Cover*prop.Condition/100)
}

// LaunderFee is the cut taken to make the paperwork hold. Clearing more costs
// more, and a well-kept business does it more cheaply.
func (w *World) LaunderFee(id string) int {
	prop := w.Properties[id]
	if prop == nil {
		return 0
	}
	cleared := min(w.Player.Heat, w.launderCapacity(id))
	return 40 + cleared*22 + (100-prop.Condition)*2
}

// LaunderReadiness explains why the books cannot take it, or returns "".
func (w *World) LaunderReadiness(id string) string {
	// What is run here, not what kind of room it is. Asking the room meant a
	// casino could not put its own takings through its own books, which is
	// most of the reason anybody owns one.
	if trade, runs := TradeOf(id); !runs || trade.Cover <= 0 || !w.Own(id) {
		return "You need a business of your own that handles cash"
	}
	if w.Properties[id].Condition < 40 {
		return "The premises are in no state to explain anything"
	}
	if w.Player.Heat == 0 {
		return "Nobody is asking you about anything"
	}
	if w.Minute-w.Player.LastLaunder < LaunderCooldown && w.Player.LastLaunder > 0 {
		return "The books have already absorbed this much recently"
	}
	if w.Player.Cash < w.LaunderFee(id) {
		return "Not enough cash to cover the cut"
	}
	return ""
}

// Launder runs takings through an owned business, trading money and wear for
// police attention.
func (w *World) Launder(id string) error {
	if reason := w.LaunderReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	fee := w.LaunderFee(id)
	if err := w.Pay(fee); err != nil {
		return err
	}
	cleared := min(w.Player.Heat, w.launderCapacity(id))
	w.Player.Heat = max(0, w.Player.Heat-cleared)
	w.Player.LastLaunder = w.Minute
	prop := w.Properties[id]
	prop.Condition = max(0, prop.Condition-3)
	// The machines are always busy and the regulars go somewhere else. This is
	// what using a shop for something it is not costs the shop.
	w.ShiftCustom(id, "The books have been used for something else once too often", -CustomLaunderLoss)
	place, _ := PlaceByID(id)
	w.Log("The books absorb it", fmt.Sprintf("$%d through %s. Police attention falls by %d, to %d. The premises take a little more wear each time, and trade at %s is down to %d%%.",
		fee, place.Name, cleared, w.Player.Heat, place.Name, w.Custom(id)), "business")
	return nil
}
