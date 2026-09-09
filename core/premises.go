package core

import "fmt"

// The city view shows who is standing in each place and nothing about the place
// itself. A player looking at twelve addresses cannot see that one of them is
// out of soap, one has a press broken, one has a still running in the back and
// one has a man on the door — all of which the city knows and all of which
// change what they should do next.
//
// Six badges on every card would be a wall again. So this is the same
// discipline the people got: one line, the most important true thing about the
// place right now, and nothing a passer-by could not see unless the player owns
// it.

// PlaceNote is the one thing worth knowing about an address today.
func (w *World) PlaceNote(id string) string {
	prop := w.Properties[id]
	if prop == nil {
		return ""
	}
	l, _ := PlaceByID(id)
	trade, running := TradeOf(id)

	// What anybody walking past can see. A boarded window is not a secret.
	if !w.Own(id) {
		switch {
		case prop.Condition < 40:
			return "Boarded up and badly knocked about"
		case prop.Condition < 70:
			return "Visibly in poor repair"
		case l.Type == "home":
			return ""
		}
		return ""
	}

	// Inside your own premises, in the order a proprietor would worry.
	switch {
	case w.Player.Heat >= ForfeitThreshold && prop.Income > 0:
		return "At this much attention they can take it"
	case prop.Trouble && running:
		return trade.Trouble
	case running && prop.Supply <= 0:
		return "Out of " + trade.Supplies
	case running && prop.Staff < trade.Hands:
		return fmt.Sprintf("Short-handed: %d of %d", prop.Staff, trade.Hands)
	case prop.Condition < 60:
		return fmt.Sprintf("Wants repair at %d%%", prop.Condition)
	case prop.Still:
		return "A still running in the back"
	case w.PostedAt(id) != nil && w.Travelling(w.PostedAt(id)):
		n := w.PostedAt(id)
		return n.Name + " is on the way, " + counted(max(1, n.Arrives-w.Minute), "minute", "minutes") + " out"
	case w.PostedAt(id) != nil:
		return w.PostedAt(id).Name + " is on the door"
	case running && w.Custom(id) < 40:
		return "The regulars have gone elsewhere"
	case prop.Income > 0:
		return fmt.Sprintf("Trading at %d%% of an ordinary day", int(w.Trading(id)*100))
	}
	return ""
}

// PlaceWarn is whether that note is something to worry about rather than
// something to be pleased with, so the interface can colour it without
// re-deciding what any of it means.
func (w *World) PlaceWarn(id string) bool {
	prop := w.Properties[id]
	if prop == nil {
		return false
	}
	if !w.Own(id) {
		return prop.Condition < 70
	}
	trade, running := TradeOf(id)
	return (w.Player.Heat >= ForfeitThreshold && prop.Income > 0) ||
		(prop.Trouble && running) ||
		(running && prop.Supply <= 0) ||
		(running && prop.Staff < trade.Hands) ||
		prop.Condition < 60 ||
		(running && w.Custom(id) < 40)
}

// Trading is what a place is earning against an ordinary day at it:
// how well it is staffed and supplied, how much custom it keeps, and the
// condition of the building — which the clock uses to scale every dollar and
// which this figure used to leave out. A laundry knocked down to 60% earned
// 169 a day where a sound one earned 305, and told the player it was trading
// at everything it could.
//
// It is deliberately not a share of a maximum. A place with more regulars than
// usual earns more than an ordinary day, and the live campaign duly read
// "Trading at 102% of what it could" — a percentage of a ceiling, above the
// ceiling. The baseline is an ordinary day, and beating it is the point of
// keeping custom.
func (w *World) Trading(id string) float64 {
	prop := w.Properties[id]
	if prop == nil {
		return 0
	}
	return w.Capacity(id) * w.TradeMultiplier(id) * float64(prop.Condition) / 100
}
