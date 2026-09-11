package core

import "fmt"

// Running a room of yours off your own still.
//
// What came off the back of a laundry had exactly one buyer in this city: the
// market, at whatever a crate was worth that day. So a player who owned a still
// and a bar was carrying crates past his own cellar to sell them to somebody
// else, and then paying cash to have drink delivered to the bar. The two halves
// were built a long way apart and nothing ever joined them.
//
// It is cheaper than buying the round in, and it is cheaper for a reason that
// is not a discount: the crates were already yours, and a crate on your hands
// is a crate the police can find. Moving it into a cellar behind a bar that
// sells drink is the one place in this city stock stops being contraband and
// starts being stock.

// OwnCellarReadiness explains why a room cannot be run off your own crates, or
// returns "".
func (w *World) OwnCellarReadiness(id string) string {
	trade, ok := TradeOf(id)
	if !ok || !w.Own(id) {
		return "This is not a business of yours"
	}
	if trade.Drink == 0 {
		return "Nobody drinks here"
	}
	if w.Player.Location != id {
		place, _ := PlaceByID(id)
		return "The crates would have to be carried into " + place.Name
	}
	if w.Properties[id].Supply >= trade.RestockAmount {
		return "It has all it needs"
	}
	if held := w.Holding("moonshine"); held < trade.Drink {
		return fmt.Sprintf("It takes %s and you are holding %d",
			plainly(trade.Drink, "one crate", fmt.Sprintf("%d crates", trade.Drink)), held)
	}
	return ""
}

// RunItOffYourOwn puts the room's supplies back up out of the player's crates.
func (w *World) RunItOffYourOwn(id string) error {
	if reason := w.OwnCellarReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	trade, _ := TradeOf(id)
	w.Player.Stock["moonshine"] -= trade.Drink
	prop := w.Properties[id]
	prop.Supply = trade.RestockAmount
	place, _ := PlaceByID(id)
	held := w.Holding("moonshine")
	w.Log("Out of your own cellar at "+place.Name,
		fmt.Sprintf("%s down the steps and nothing paid for it. %s is trading properly again and the $%d stayed in your pocket. You are holding %s.",
			plainly(trade.Drink, "One crate", fmt.Sprintf("%d crates", trade.Drink)),
			place.Name, w.RestockCost(id),
			plainly(held, "one crate", fmt.Sprintf("%d crates", held))), "business")
	return nil
}
