package core

import "sort"

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
		// Before the flop there is no five-card hand yet. Rank answers the
		// question it is asked — five cards — and asked about two it called
		// any two of a suit a flush, so the screen told a player holding the
		// queen and three of diamonds that they had one.
		return beforeTheFlop(seven)
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

// beforeTheFlop is what a hand is worth before there are five cards to make one of. A
// pair is a pair; two of a suit is two of a suit and nothing more.
func beforeTheFlop(cards []Card) HandRank {
	count := map[int]int{}
	for _, c := range cards {
		count[pokerRank(c)]++
	}
	order := make([]int, 0, len(count))
	for r := range count {
		order = append(order, r)
	}
	sort.Slice(order, func(i, j int) bool {
		if count[order[i]] != count[order[j]] {
			return count[order[i]] > count[order[j]]
		}
		return order[i] > order[j]
	})
	if len(order) > 0 && count[order[0]] >= 2 {
		return HandRank{Pair, order}
	}
	return HandRank{HighCard, order}
}
