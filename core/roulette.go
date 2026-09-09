package core

import "fmt"

// A wheel, which the room has and the game did not. Blackjack is a hand the
// player keeps making decisions inside; roulette is one decision made before
// anything happens and then nothing you can do about it, which is a different
// thing to offer a person and the reason a room has both.
//
// The house edge here is not a fudge factor applied to the odds. Every payout
// is the true one — a single number pays 35 to 1 against 36 other pockets, red
// pays even against eighteen others — and the whole of the advantage is the
// zero, which belongs to nobody. That is how a real wheel works and it means
// the edge can be stated rather than tuned: one pocket in thirty-seven.

const (
	// Pockets is the wheel: nought to thirty-six, single zero.
	Pockets = 37
	// Zero is the pocket that belongs to the house. It is not red, not black,
	// not odd, not even, and not in any dozen, so it takes every outside bet
	// on the table and that is the entire edge.
	Zero = 0
)

// reds are the red pockets on a single-zero wheel, which are not every other
// number and cannot be worked out arithmetically. This is the actual layout.
var reds = map[int]bool{
	1: true, 3: true, 5: true, 7: true, 9: true, 12: true, 14: true, 16: true,
	18: true, 19: true, 21: true, 23: true, 25: true, 27: true, 30: true,
	32: true, 34: true, 36: true,
}

// Red reports whether a pocket is red. The zero is neither colour.
func Red(pocket int) bool { return reds[pocket] }

// Black reports whether a pocket is black.
func Black(pocket int) bool { return pocket != Zero && !reds[pocket] }

// RouletteBet is one thing the player can back, what it pays, and how to tell
// whether it came in. Pays is to one: a single number pays 35 to 1, so a
// winning dollar comes back as thirty-six.
type RouletteBet struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Pays  int    `json:"pays"`
	// Wins reports whether this bet comes in on that pocket.
	Wins func(pocket int) bool `json:"-"`
}

// RouletteBets is everything on the cloth. Straight numbers are generated
// rather than written out, because thirty-seven hand-written entries is
// thirty-seven chances to write one of them wrong.
func RouletteBets() []RouletteBet {
	bets := []RouletteBet{
		{ID: "red", Label: "Red", Pays: 1, Wins: Red},
		{ID: "black", Label: "Black", Pays: 1, Wins: Black},
		{ID: "odd", Label: "Odd", Pays: 1, Wins: func(p int) bool { return p != Zero && p%2 == 1 }},
		{ID: "even", Label: "Even", Pays: 1, Wins: func(p int) bool { return p != Zero && p%2 == 0 }},
		{ID: "low", Label: "One to eighteen", Pays: 1, Wins: func(p int) bool { return p >= 1 && p <= 18 }},
		{ID: "high", Label: "Nineteen to thirty-six", Pays: 1, Wins: func(p int) bool { return p >= 19 && p <= 36 }},
		{ID: "dozen1", Label: "First dozen", Pays: 2, Wins: func(p int) bool { return p >= 1 && p <= 12 }},
		{ID: "dozen2", Label: "Second dozen", Pays: 2, Wins: func(p int) bool { return p >= 13 && p <= 24 }},
		{ID: "dozen3", Label: "Third dozen", Pays: 2, Wins: func(p int) bool { return p >= 25 && p <= 36 }},
	}
	for n := 0; n <= 36; n++ {
		number := n
		bets = append(bets, RouletteBet{
			ID:    fmt.Sprintf("number:%d", number),
			Label: fmt.Sprintf("Straight up on %s", NumberWord(number)),
			Pays:  35,
			Wins:  func(p int) bool { return p == number },
		})
	}
	return bets
}

// RouletteBetByID finds a bet, and whether there is one.
func RouletteBetByID(id string) (RouletteBet, bool) {
	for _, b := range RouletteBets() {
		if b.ID == id {
			return b, true
		}
	}
	return RouletteBet{}, false
}

// NumberWord is how the room says a number aloud. Nobody at a table says
// "zero"; they say the house number.
func NumberWord(n int) string {
	if n == Zero {
		return "the nought"
	}
	return fmt.Sprintf("%d", n)
}

// PocketColour is how a pocket is described once the ball has dropped.
func PocketColour(pocket int) string {
	switch {
	case pocket == Zero:
		return "green"
	case Red(pocket):
		return "red"
	}
	return "black"
}

// Spin is the wheel turning, kept on the world for the same reason a hand is:
// the money is down and the result is real.
type Spin struct {
	// Place is the room and Stake what is on the cloth.
	Place, Stake string
	// Bet is what was backed, Pocket what came up, and Done whether it has
	// been settled — so nothing settles twice.
	Bet    string `json:"bet"`
	Pocket int    `json:"pocket"`
	Won    bool   `json:"won"`
	Done   bool   `json:"done"`
}

// SpinReadiness explains why the wheel cannot be played, or returns "".
func (w *World) SpinReadiness(id, betID string, stake Stake) string {
	if reason := w.TableReadiness(id, stake); reason != "" {
		return reason
	}
	if _, ok := RouletteBetByID(betID); !ok {
		return "That is not a bet this room takes"
	}
	if w.Hand != nil && !w.Hand.Done {
		return "There is a hand of cards on the table already"
	}
	return ""
}

// PlayWheel puts money on the cloth and turns the wheel. Unlike a hand of
// cards there is nothing to decide afterwards, so it settles in one go.
func (w *World) PlayWheel(id, betID, stakeID string) error {
	stake, ok := tableStake(stakeID)
	if !ok {
		return fmt.Errorf("no such game")
	}
	if reason := w.SpinReadiness(id, betID, stake); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	bet, _ := RouletteBetByID(betID)
	if err := w.Pay(stake.Amount); err != nil {
		return err
	}
	// The player is at this table, so this is the player's stream.
	pocket := int(w.Random() * Pockets)
	if pocket >= Pockets {
		pocket = Pockets - 1
	}
	won := bet.Wins(pocket)
	w.Spin = &Spin{Place: id, Stake: stakeID, Bet: betID, Pocket: pocket, Won: won, Done: true}

	place, _ := PlaceByID(id)
	house := w.faction(w.Properties[id].Owner)
	returned := 0
	if won {
		// The stake comes back with it: a dollar on a number that comes in is
		// thirty-six dollars, not thirty-five.
		returned = stake.Amount * (bet.Pays + 1)
		w.Earn(returned)
	}
	net := returned - stake.Amount
	if house != nil {
		house.Cash = max(0, house.Cash-net)
	}
	w.tableAftermath(place.Name, house, net, stake.Amount)
	dropped := fmt.Sprintf("%s, %s.", NumberWord(pocket), PocketColour(pocket))
	if won {
		w.Log("The wheel at "+place.Name, fmt.Sprintf("%s on %s. %s $%d comes back across the cloth.",
			bet.Label, place.Name, dropped, returned), "business")
	} else {
		w.Log("The wheel at "+place.Name, fmt.Sprintf("%s on %s. %s The $%d stays where it is.",
			bet.Label, place.Name, dropped, stake.Amount), "business")
	}
	return nil
}

// WheelDescription is the last spin, for the interface to show.
func (w *World) WheelDescription() map[string]any {
	if w.Spin == nil {
		return map[string]any{"spun": false}
	}
	place, _ := PlaceByID(w.Spin.Place)
	bet, _ := RouletteBetByID(w.Spin.Bet)
	stake, _ := tableStake(w.Spin.Stake)
	return map[string]any{
		"spun": true, "place": place.Name, "stake": stake.Amount,
		"bet": bet.Label, "pocket": w.Spin.Pocket, "colour": PocketColour(w.Spin.Pocket),
		"won": w.Spin.Won, "pays": bet.Pays,
	}
}

// stakeSuffix distinguishes the two wheels in a list of buttons without
// repeating the price, which the button already carries.
func stakeSuffix(stake Stake) string {
	if stake.ID == "high" {
		return " at the high tables"
	}
	return ""
}
