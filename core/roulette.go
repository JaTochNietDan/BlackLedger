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

// Chip is one bet on the cloth: what it backs and what is on it. A real table
// takes as many of these as somebody can reach, all settled against the same
// pocket, which is most of what makes roulette a game rather than a coin toss
// with extra numbers.
type Chip struct {
	Bet    string `json:"bet"`
	Amount int    `json:"amount"`
	// Won and Back are filled in once the ball has dropped, so the cloth can be
	// drawn afterwards showing which chips came in.
	Won  bool `json:"won,omitempty"`
	Back int  `json:"back,omitempty"`
}

// Spin is the wheel turning, kept on the world for the same reason a hand is:
// the money is down and the result is real.
type Spin struct {
	// Place is the room. Stake was one of two fixed lots by name and is kept so
	// a spin in an older save still reads; Down is what was actually on the
	// cloth, because a player names the number now.
	Place, Stake string
	Down         int `json:"down,omitempty"`
	// Bet is what was backed, Pocket what came up, and Done whether it has
	// been settled — so nothing settles twice.
	Bet    string `json:"bet"`
	Pocket int    `json:"pocket"`
	Won    bool   `json:"won"`
	Done   bool   `json:"done"`
	// Every chip that was on the cloth. A spin from before the table took more
	// than one bet has none, and reads by Bet and Down as it always did.
	Chips []Chip `json:"chips"`
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

// ChipsReadiness explains why this cloth cannot be spun, or returns "".
//
// The house limit is per bet, the way a real table's is: two chips at the limit
// are two bets and are fine. What is not fine is more on the cloth than the
// player has, however it is spread about.
func (w *World) ChipsReadiness(id string, chips []Chip) string {
	if len(chips) == 0 {
		return "There is nothing on the cloth"
	}
	total := 0
	for _, c := range chips {
		if _, ok := RouletteBetByID(c.Bet); !ok {
			return "That is not a bet this room takes"
		}
		if reason := w.TableReadiness(id, Stake{Amount: c.Amount}); reason != "" {
			return reason
		}
		total += c.Amount
	}
	if w.Player.Cash < total {
		return fmt.Sprintf("There is $%d on the cloth and you have $%d", total, w.Player.Cash)
	}
	if w.Hand != nil && !w.Hand.Done {
		return "There is a hand of cards on the table already"
	}
	if w.Dice != nil && !w.Dice.Done {
		return "There is a point on at the dice table"
	}
	return ""
}

// SpinChips puts a whole cloth down and turns the wheel once. Every chip is
// settled against the same pocket, which is the only honest way to do it: one
// number comes up and the table pays everybody who had it.
func (w *World) SpinChips(id string, chips []Chip) error {
	if reason := w.ChipsReadiness(id, chips); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	total := 0
	for _, c := range chips {
		total += c.Amount
	}
	if err := w.Pay(total); err != nil {
		return err
	}
	// Chips on the cloth is sitting down at it.
	w.takeASeatFor(id, Floor)
	// The player is at this table, so this is the player's stream.
	pocket := int(w.Random() * Pockets)
	if pocket >= Pockets {
		pocket = Pockets - 1
	}
	returned, won := 0, false
	settled := make([]Chip, 0, len(chips))
	for _, c := range chips {
		bet, _ := RouletteBetByID(c.Bet)
		if bet.Wins(pocket) {
			c.Won, c.Back = true, c.Amount*(bet.Pays+1)
			returned += c.Back
			won = true
		}
		settled = append(settled, c)
	}
	// Bet and Down keep reading for one chip, so a save written before the
	// table took a cloth still describes itself, and so does this one.
	lead := settled[0].Bet
	w.Spin = &Spin{Place: id, Down: total, Bet: lead, Pocket: pocket, Won: won, Done: true, Chips: settled}
	if returned > 0 {
		w.Earn(returned)
	}

	place, _ := PlaceByID(id)
	house := w.faction(w.Properties[id].Owner)
	net := returned - total
	w.changeBusinessFunds(id, -net)
	w.tableAftermath(place.Name, house, net, total)
	dropped := fmt.Sprintf("%s, %s.", NumberWord(pocket), PocketColour(pocket))
	what := chipWords(settled)
	if returned > 0 {
		w.Log("The wheel at "+place.Name, fmt.Sprintf("%s on %s. %s $%d comes back across the cloth.",
			what, place.Name, dropped, returned), "business")
	} else {
		w.Log("The wheel at "+place.Name, fmt.Sprintf("%s on %s. %s The $%d stays where it is.",
			what, place.Name, dropped, total), "business")
	}
	return nil
}

// chipWords describes what was on the cloth without listing thirty of them.
func chipWords(chips []Chip) string {
	if len(chips) == 1 {
		bet, _ := RouletteBetByID(chips[0].Bet)
		return bet.Label
	}
	return fmt.Sprintf("%d chips on the cloth", len(chips))
}

// PlayWheel is one chip on one bet, which is what the table took before it took
// a cloth. Kept because a single bet is still the commonest thing anybody does.
func (w *World) PlayWheel(id, betID string, amount int) error {
	if reason := w.SpinReadiness(id, betID, Stake{Amount: amount}); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	return w.SpinChips(id, []Chip{{Bet: betID, Amount: amount}})
}

// WheelDescription is the last spin, for the interface to show.
func (w *World) WheelDescription() map[string]any {
	if w.Spin == nil {
		return map[string]any{"spun": false}
	}
	place, _ := PlaceByID(w.Spin.Place)
	bet, _ := RouletteBetByID(w.Spin.Bet)
	stake := Stake{Amount: w.spinDown()}
	// The whole cloth, so the table can show which chips came in rather than
	// describing one of them.
	chips := []map[string]any{}
	for _, c := range w.Spin.Chips {
		on, _ := RouletteBetByID(c.Bet)
		chips = append(chips, map[string]any{"bet": c.Bet, "label": on.Label, "amount": c.Amount, "won": c.Won, "back": c.Back})
	}
	back := 0
	for _, c := range w.Spin.Chips {
		back += c.Back
	}
	return map[string]any{
		"spun": true, "place": place.Name, "stake": stake.Amount,
		"bet": bet.Label, "pocket": w.Spin.Pocket, "colour": PocketColour(w.Spin.Pocket),
		"won": w.Spin.Won, "pays": bet.Pays, "chips": chips, "back": back,
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

// spinDown is what was on the cloth. A spin recorded before players could name
// a number carries one of the two old lots by name instead.
func (w *World) spinDown() int {
	if w.Spin == nil {
		return 0
	}
	if w.Spin.Down > 0 {
		return w.Spin.Down
	}
	if stake, ok := tableStake(w.Spin.Stake); ok {
		return stake.Amount
	}
	return 0
}
