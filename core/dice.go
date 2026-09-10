package core

import "fmt"

// The dice, which every room like this had and this one did not.
//
// The three games in here are three different shapes of decision. A hand of
// cards is a decision on every card. The wheel is one decision and then nothing
// at all. This is one decision followed by a run of rolls you cannot affect but
// have to sit through — the point is set, and then it is the point or the seven
// and nobody at the table knows which.
//
// As with the wheel, the edge is not a number applied to the odds. Every payout
// here is even money and the whole of the advantage is in the arithmetic of two
// dice: the pass line wins 244 times in 495 and loses 251, which is 1.41% and
// cannot be tuned without making the game a lie.

const (
	// PassPays and DontPays are even money, which is what these bets pay.
	PassPays = 1
	// DontBar is the total that pushes on the don't pass rather than winning.
	// Without it the wrong-way bettor would have the advantage, which is why
	// every table in the world bars one number.
	DontBar = 12
)

// Craps is the game in front of the shooter. Kept on the world for the same
// reason a hand is: the money is down and the point is real.
type Craps struct {
	Place string `json:"place"`
	// Bet is what was backed: the pass line, the don't, or the field.
	Bet  string `json:"bet"`
	Down int    `json:"down"`
	// Point is the number that has to be made, or zero on the come-out.
	Point int    `json:"point"`
	Dice  [2]int `json:"dice"`
	Total int    `json:"total"`
	Rolls int    `json:"rolls"`
	Won   bool   `json:"won"`
	Push  bool   `json:"push"`
	Done  bool   `json:"done"`
	// Outcome is what the table would say out loud.
	Outcome string `json:"outcome"`
}

// DiceBet is one thing the player can back on the dice.
type DiceBet struct {
	ID, Label, Detail string
}

// DiceBets is what this table takes. Three bets, all of them real, and the two
// line bets are the only ones in the house worth playing.
func DiceBets() []DiceBet {
	return []DiceBet{
		{"pass", "The pass line", "Seven or eleven wins it now, two three or twelve loses it now, anything else is the point and has to come again before a seven. The house keeps about 1.4 in every hundred."},
		{"dont", "Don't pass", "The other side of the same game: two or three wins now, seven or eleven loses now, twelve is a standoff, and after that the seven is what you want. About 1.4 in every hundred, the other way round."},
		{"field", "The field", "One roll. Two three four nine ten eleven or twelve, and the two pays double and the twelve triple. Everything else loses. It looks like most of the numbers and it is not: the house keeps about 2.8 in every hundred."},
	}
}

// diceBetLabel is what to call a bet in a sentence.
func diceBetLabel(id string) string {
	if b, ok := DiceBetByID(id); ok {
		return b.Label
	}
	return "the line"
}

// DiceBetByID finds a bet by name.
func DiceBetByID(id string) (DiceBet, bool) {
	for _, b := range DiceBets() {
		if b.ID == id {
			return b, true
		}
	}
	return DiceBet{}, false
}

// roll is two dice out of the player's own stream, because the player is
// standing at this table.
//
// A note on measuring this, because I got it wrong first. Over 200,000
// decisions the pass line came out at 1.75% against a true 1.414%, and that
// looked like the world's linear congruential generator failing on consecutive
// draws — the classic lattice problem, and craps is the only game in the house
// that resolves over a sequence of rolls rather than one. So the draw was
// mixed. Then the same measurement at five million decisions gave 1.402%
// unmixed and 1.431% mixed, both inside one standard error of the truth.
//
// The standard error at 200,000 decisions is 0.22%. The number that looked like
// a fault was 1.5 of those from the true value, which is nothing at all. The
// generator was never the problem and the mixing is gone; what the test does
// instead is take a sample big enough for its own claim.
func (w *World) roll() (int, int) {
	return 1 + int(w.Random()*6), 1 + int(w.Random()*6)
}

// DiceReadiness explains why the dice cannot be played, or returns "".
func (w *World) DiceReadiness(id, betID string, stake Stake) string {
	if reason := w.TableReadiness(id, stake); reason != "" {
		return reason
	}
	if _, ok := DiceBetByID(betID); !ok {
		return "That is not a bet this table takes"
	}
	if w.Hand != nil && !w.Hand.Done {
		return "There is a hand of cards on the table already"
	}
	if w.Dice != nil && !w.Dice.Done {
		return "There is a point on already"
	}
	return ""
}

// PlayDice puts money on the line and throws the come-out.
func (w *World) PlayDice(id, betID string, amount int) error {
	stake := Stake{ID: "typed", Label: "The dice", Amount: amount}
	if reason := w.DiceReadiness(id, betID, stake); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(stake.Amount); err != nil {
		return err
	}
	w.Dice = &Craps{Place: id, Bet: betID, Down: stake.Amount}
	w.throw()
	return nil
}

// RollReadiness explains why there is nothing to throw, or returns "".
func (w *World) RollReadiness() string {
	if w.Dice == nil || w.Dice.Done {
		return "There is nothing on the line"
	}
	if w.Player.Location != w.Dice.Place {
		return "The table is not here"
	}
	return ""
}

// RollAgain throws again for a point that is still out there. No new money:
// the bet is already down, which is the whole shape of the game.
func (w *World) RollAgain() error {
	if reason := w.RollReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	w.throw()
	return nil
}

// throw rolls the dice and works out what it did to the bet on the table.
func (w *World) throw() {
	g := w.Dice
	a, b := w.roll()
	g.Dice, g.Total, g.Rolls = [2]int{a, b}, a+b, g.Rolls+1

	switch g.Bet {
	case "field":
		// One roll and it is over, whatever it is.
		pays := 0
		switch g.Total {
		case 2:
			pays = 2
		case 12:
			pays = 3
		case 3, 4, 9, 10, 11:
			pays = 1
		}
		g.Won, g.Done = pays > 0, true
		if pays > 0 {
			g.Outcome = fmt.Sprintf("%s. The field pays %d to 1.", diceWords(a, b), pays)
			w.settleDice(g.Down * (pays + 1))
			return
		}
		g.Outcome = diceWords(a, b) + ". The field does not take that number."
		w.settleDice(0)
		return
	case "dont":
		if g.Point == 0 {
			switch g.Total {
			case 2, 3:
				g.Won, g.Done, g.Outcome = true, true, comeOutWords(a, b)+". The wrong way is the right way."
				w.settleDice(g.Down * (PassPays + 1))
			case DontBar:
				// Barred: the money stays where it is and the come-out is
				// thrown again. Nothing is decided and nothing changes hands.
				g.Outcome = comeOutWords(a, b) + ". Barred. Nothing is decided; throw again."
			case 7, 11:
				g.Won, g.Done, g.Outcome = false, true, comeOutWords(a, b)+". A natural, and the wrong way pays for it."
				w.settleDice(0)
			default:
				g.Point = g.Total
				g.Outcome = fmt.Sprintf("%s. The point is %d, and now the seven is yours.", comeOutWords(a, b), g.Point)
			}
			return
		}
		switch g.Total {
		case 7:
			g.Won, g.Done, g.Outcome = true, true, diceWords(a, b)+". Seven out, and it pays."
			w.settleDice(g.Down * (PassPays + 1))
		case g.Point:
			g.Won, g.Done, g.Outcome = false, true, fmt.Sprintf("%s. The point is made and the wrong way pays for it.", diceWords(a, b))
			w.settleDice(0)
		default:
			g.Outcome = fmt.Sprintf("%s. Still %d to beat.", diceWords(a, b), g.Point)
		}
		return
	}

	// The pass line.
	if g.Point == 0 {
		switch g.Total {
		case 7, 11:
			g.Won, g.Done, g.Outcome = true, true, comeOutWords(a, b)+". A natural, and it pays."
			w.settleDice(g.Down * (PassPays + 1))
		case 2, 3, 12:
			g.Won, g.Done, g.Outcome = false, true, comeOutWords(a, b)+". The line is taken."
			w.settleDice(0)
		default:
			g.Point = g.Total
			g.Outcome = fmt.Sprintf("%s. The point is %d, and it has to come again before a seven.", comeOutWords(a, b), g.Point)
		}
		return
	}
	switch g.Total {
	case g.Point:
		g.Won, g.Done, g.Outcome = true, true, fmt.Sprintf("%s. The point, made, and it pays.", diceWords(a, b))
		w.settleDice(g.Down * (PassPays + 1))
	case 7:
		g.Won, g.Done, g.Outcome = false, true, diceWords(a, b)+". Seven out. The line is taken."
		w.settleDice(0)
	default:
		g.Outcome = fmt.Sprintf("%s. Still looking for %d.", diceWords(a, b), g.Point)
	}
}

// settleDice moves the money once a decision is reached, and tells the house
// about it the same way every other table does.
func (w *World) settleDice(returned int) {
	g := w.Dice
	if returned > 0 {
		w.Earn(returned)
	}
	place, _ := PlaceByID(g.Place)
	house := w.faction(w.Properties[g.Place].Owner)
	net := returned - g.Down
	if house != nil {
		house.Cash = max(0, house.Cash-net)
	}
	w.tableAftermath(place.Name, house, net, g.Down)
	bet, _ := DiceBetByID(g.Bet)
	w.Log("The dice at "+place.Name, fmt.Sprintf("%s. %s", bet.Label, g.Outcome), "business")
}

// diceWords is what the stickman calls it. A two, a three or a twelve is only
// "craps" on the come-out; once a point is on, a three is a three and the only
// numbers anybody in the room cares about are the point and the seven.
func diceWords(a, b int) string { return call(a, b, false) }

func comeOutWords(a, b int) string { return call(a, b, true) }

func call(a, b int, comeOut bool) string {
	total := a + b
	name := map[int]string{2: "snake eyes", 7: "seven", 11: "yo, eleven", 12: "boxcars"}[total]
	if total == 3 {
		name = "three"
	}
	if name == "" {
		name = spokenNumber(total)
	}
	if comeOut && (total == 2 || total == 3 || total == 12) {
		name += ", craps"
	}
	return fmt.Sprintf("%d and %d — %s", a, b, name)
}

// spokenNumber is the total said rather than printed. "8, 8" is a table read
// out by a machine.
func spokenNumber(n int) string {
	words := map[int]string{4: "four", 5: "five", 6: "six", 8: "eight", 9: "nine", 10: "ten"}
	if word, ok := words[n]; ok {
		return word
	}
	return fmt.Sprintf("%d", n)
}

// DiceDescription is the table as it stands, for the interface to draw.
func (w *World) DiceDescription() map[string]any {
	out := map[string]any{"playing": false, "settled": false, "bets": diceBetList()}
	if w.Dice == nil {
		return out
	}
	place, _ := PlaceByID(w.Dice.Place)
	bet, _ := DiceBetByID(w.Dice.Bet)
	out["playing"] = !w.Dice.Done
	out["settled"] = w.Dice.Done
	out["place"] = place.Name
	out["where"] = w.Dice.Place
	out["bet"] = w.Dice.Bet
	out["bet_label"] = bet.Label
	out["stake"] = w.Dice.Down
	out["point"] = w.Dice.Point
	out["dice"] = []int{w.Dice.Dice[0], w.Dice.Dice[1]}
	out["total"] = w.Dice.Total
	out["rolls"] = w.Dice.Rolls
	out["won"] = w.Dice.Won
	out["outcome"] = w.Dice.Outcome
	return out
}

func diceBetList() []map[string]any {
	out := []map[string]any{}
	for _, b := range DiceBets() {
		out = append(out, map[string]any{"id": b.ID, "label": b.Label, "detail": b.Detail})
	}
	return out
}
