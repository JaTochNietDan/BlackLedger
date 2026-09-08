package core

import "fmt"

// A casino is not only premises to own. It is a room you can walk into with
// money and walk out of without it. The house wins over time, which is exactly
// why owning one is worth more than playing at one, and why a player who wins
// heavily at somebody else's tables has taken it out of their pocket.

// Stakes are fixed so the interface can offer plain actions. Playing deep is
// how a run at the tables becomes worth an owner's attention.
type Stake struct {
	ID     string
	Label  string
	Amount int
}

var tableStakes = []Stake{
	{ID: "small", Label: "Play the small tables", Amount: 50},
	{ID: "high", Label: "Play the high tables", Amount: 300},
}

// HighTableStanding is what a room wants to see before it lets somebody sit
// down with real money: what they have done, or what they are wearing, or
// enough of both. A tailored suit alone is very nearly the price of entry,
// which is the point of owning one.
const HighTableStanding = 14

func tableStake(id string) (Stake, bool) {
	for _, s := range tableStakes {
		if s.ID == id {
			return s, true
		}
	}
	return Stake{}, false
}

// HasTables reports whether a place runs games at all.
func HasTables(id string) bool {
	place, ok := PlaceByID(id)
	return ok && place.Type == "casino"
}

// TableReadiness explains why the player cannot sit down, or returns "".
func (w *World) TableReadiness(id string, stake Stake) string {
	if !HasTables(id) {
		return "There are no tables here"
	}
	if w.Own(id) {
		return "You cannot win money from your own house"
	}
	if w.Properties[id] != nil && w.Properties[id].Condition < 30 {
		return "The room is in no state to run games"
	}
	if stake.ID == "high" && w.Presence() < HighTableStanding {
		return "The floor manager looks at you and seats you nowhere near that table"
	}
	if w.Player.Cash < stake.Amount {
		return "Not enough cash"
	}
	return ""
}

// Play resolves a session at the tables. The house holds an edge, so this is a
// way to lose money quickly and occasionally a way to make it.
func (w *World) Play(id, stakeID string) error {
	stake, ok := tableStake(stakeID)
	if !ok {
		return fmt.Errorf("no such game")
	}
	if reason := w.TableReadiness(id, stake); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	place, _ := PlaceByID(id)
	if err := w.Pay(stake.Amount); err != nil {
		return err
	}
	house := w.faction(w.Properties[id].Owner)

	// Expected return is under one: the house keeps roughly six percent of
	// everything staked, and most nights end down. Measured rather than
	// asserted, in TestTheHouseKeepsItsEdge.
	roll := w.Random()
	multiple := 0.0
	switch {
	case roll < .38:
		multiple = 0 // the stake is gone
	case roll < .72:
		multiple = .3 + w.Random()*.9 // most of these still leave you short
	case roll < .96:
		multiple = 1.5 + w.Random()
	default:
		multiple = 3 + w.Random()*4
	}
	returned := int(float64(stake.Amount) * multiple)
	net := returned - stake.Amount
	if returned > 0 {
		w.Earn(returned)
	}

	if house != nil {
		house.Cash = max(0, house.Cash-net)
	}
	switch {
	case net > stake.Amount*2:
		// Winning heavily at somebody's tables is noticed by the somebody.
		w.Player.Heat = min(100, w.Player.Heat+3)
		if house != nil {
			house.Goodwill = max(-100, house.Goodwill-8)
			w.Log("A good night at "+place.Name, fmt.Sprintf("You walk out $%d up. The floor manager watched you do it, and %s does not enjoy losing in its own room.", net, house.Name), "politics")
			return nil
		}
		w.Log("A good night at "+place.Name, fmt.Sprintf("You walk out $%d up. Nights like this are remembered.", net), "business")
	case net > 0:
		w.Log("A quiet win at "+place.Name, fmt.Sprintf("You leave $%d ahead. Nobody looks twice.", net), "business")
	case net == 0:
		w.Log("An even night at "+place.Name, "You leave with what you brought, which is more than most manage.", "business")
	default:
		w.Log("The house takes it at "+place.Name, fmt.Sprintf("You lose $%d. The tables do not care who you are.", -net), "business")
	}
	return nil
}
