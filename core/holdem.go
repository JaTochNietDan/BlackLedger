package core

// Texas hold'em, because the user asked for it: "I'd also prefer if this game
// was the Texas Hold Em version as it's better to play so we can fix that
// maybe." Draw poker gives the player one decision — which cards to throw — and
// then turns everything over. Hold'em gives them four rounds of the only
// decision that matters at a card table, against a board everybody can see, so
// what the people across from you do with the same three cards is the game.
//
// Two cards each, three on the table, then one, then one. The best five of the
// seven a player can see wins. HandRank and Beats are unchanged: what five
// cards are worth was already the one place that arithmetic lived.

// BestOfSeven is what somebody holding two cards makes of the board. Twenty-one
// hands of five, and the best of them stands. Small enough to do honestly.
func BestOfSeven(hole, board []Card) HandRank {
	seven := make([]Card, 0, 7)
	seven = append(seven, hole...)
	seven = append(seven, board...)
	if len(seven) < 5 {
		// Before the flop there is no five-card hand yet, so a pair of aces is
		// worth what a pair of aces is worth and nothing else has a name.
		return Rank(seven)
	}
	best := HandRank{Category: -1}
	five := make([]Card, 5)
	for a := 0; a < len(seven)-4; a++ {
		for b := a + 1; b < len(seven)-3; b++ {
			for c := b + 1; c < len(seven)-2; c++ {
				for d := c + 1; d < len(seven)-1; d++ {
					for e := d + 1; e < len(seven); e++ {
						five[0], five[1], five[2], five[3], five[4] = seven[a], seven[b], seven[c], seven[d], seven[e]
						if r := Rank(five); best.Category < 0 || r.Beats(best) {
							best = r
						}
					}
				}
			}
		}
	}
	return best
}

// The streets, in the order they come.
const (
	Preflop = "preflop"
	Flop    = "flop"
	Turn    = "turn"
	River   = "river"
	Shown   = "shown"
)

// nextStreet is what follows this one, and how many cards it puts on the table.
func nextStreet(street string) (string, int) {
	switch street {
	case Preflop:
		return Flop, 3
	case Flop:
		return Turn, 1
	case Turn:
		return River, 1
	}
	return Shown, 0
}

// StreetName is what to call the street in a sentence.
func StreetName(street string) string {
	switch street {
	case Preflop:
		return "before the flop"
	case Flop:
		return "on the flop"
	case Turn:
		return "on the turn"
	case River:
		return "on the river"
	}
	return "at the showdown"
}
