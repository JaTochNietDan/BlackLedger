package core

import (
	"fmt"
	"sort"
)

// Every table in this city is the house. Blackjack, the wheel, the dice and the
// machines all have an edge, a settlement and nobody on the other side of them:
// the player wins money from a room rather than from a person. The back room is
// the other thing. The money across the table belongs to somebody with a name,
// an address, a family and a reason to remember losing it, and the game has no
// edge at all because there is no house in it — five cards each, one draw, best
// hand takes what is on the table.
//
// It is also the cheapest way for the city's people to be people. A man who
// plays his hand the way somebody with his standing and his nerve would play it
// is characterised by the game rather than by a paragraph about him.

// BackRoom is where the game is. A poolhall has a room behind the room.
const BackRoom = "poolhall"

// Ante bounds. Low enough that anybody in the room can sit down, high enough
// that losing is felt.
const (
	MinAnte = 10
	MaxAnte = 500
	// Players is how many other people sit down, at most.
	Players = 3
)

// Seat is somebody at the table who is not the player.
type Seat struct {
	Who   string `json:"who"`
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
	// Threw is how many they changed, which is the only thing the player gets
	// to see about a hand before it is turned over.
	Threw int `json:"threw"`
}

// CardGame is a hand in progress in the back room. It lives on the world
// because money is on the table.
type CardGame struct {
	Place string `json:"place"`
	Ante  int    `json:"ante"`
	Pot   int    `json:"pot"`
	Mine  []Card `json:"mine"`
	Seats []Seat `json:"seats"`
	// Deck is what is left to draw from, so the draw cannot deal a card that is
	// already in somebody's hand.
	Deck []Card `json:"deck"`
	// Drawn is set once cards have been changed, so nobody draws twice.
	Drawn   bool   `json:"drawn,omitempty"`
	Done    bool   `json:"done,omitempty"`
	Outcome string `json:"outcome,omitempty"`
	Won     int    `json:"won,omitempty"`
}

// pokerRank is what a card is worth in a game where an ace is not eleven and a
// jack is not a ten. Blackjack's Value cannot answer this: it makes four
// different cards the same card.
func pokerRank(c Card) int {
	switch c.Rank {
	case "A":
		return 14
	case "K":
		return 13
	case "Q":
		return 12
	case "J":
		return 11
	default:
		return c.Value
	}
}

// deck is fifty-two distinct cards, shuffled. Blackjack draws with replacement,
// which is honest for a shoe and would put five aces in one hand here.
func (w *World) deck() []Card {
	cards := make([]Card, 0, 52)
	for _, r := range ranks {
		for _, s := range suits {
			cards = append(cards, Card{Rank: r.Rank, Suit: s, Value: r.Value})
		}
	}
	for i := len(cards) - 1; i > 0; i-- {
		j := min(i, int(w.Random()*float64(i+1)))
		cards[i], cards[j] = cards[j], cards[i]
	}
	return cards
}

// Hand categories, worst to best.
const (
	HighCard = iota
	Pair
	TwoPair
	Trips
	Straight
	Flush
	FullHouse
	Quads
	StraightFlush
)

// HandRank is what five cards are worth: the category, then the ranks that
// settle two hands in the same category, in the order they settle them.
type HandRank struct {
	Category int   `json:"category"`
	Order    []int `json:"order"`
}

// Beats reports whether this hand takes the pot from that one. Equal hands
// beat nothing, which is what a split is.
func (h HandRank) Beats(other HandRank) bool {
	if h.Category != other.Category {
		return h.Category > other.Category
	}
	for i := range h.Order {
		if i >= len(other.Order) {
			return false
		}
		if h.Order[i] != other.Order[i] {
			return h.Order[i] > other.Order[i]
		}
	}
	return false
}

// Name is what to call it at a showdown.
func (h HandRank) Name() string {
	switch h.Category {
	case StraightFlush:
		return "a straight flush"
	case Quads:
		return "four of a kind"
	case FullHouse:
		return "a full house"
	case Flush:
		return "a flush"
	case Straight:
		return "a straight"
	case Trips:
		return "three of a kind"
	case TwoPair:
		return "two pair"
	case Pair:
		return "a pair"
	}
	return "high card"
}

// Rank works out what a hand is. Counts first, because every category above a
// straight is a statement about how many of a rank there are.
func Rank(cards []Card) HandRank {
	count := map[int]int{}
	suit := map[string]int{}
	for _, c := range cards {
		count[pokerRank(c)]++
		suit[c.Suit]++
	}
	// Ranks sorted by how many of them there are, then by how high they are:
	// a full house is settled by the three before the two.
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

	flush := false
	for _, n := range suit {
		flush = flush || n == len(cards)
	}
	straight, high := isStraight(order, count)

	switch {
	case straight && flush:
		return HandRank{StraightFlush, []int{high}}
	case len(order) > 0 && count[order[0]] == 4:
		return HandRank{Quads, order}
	case len(order) == 2 && count[order[0]] == 3:
		return HandRank{FullHouse, order}
	case flush:
		return HandRank{Flush, order}
	case straight:
		return HandRank{Straight, []int{high}}
	case len(order) > 0 && count[order[0]] == 3:
		return HandRank{Trips, order}
	case len(order) == 3 && count[order[0]] == 2:
		return HandRank{TwoPair, order}
	case len(order) == 4:
		return HandRank{Pair, order}
	}
	return HandRank{HighCard, order}
}

// isStraight reports five in a row, counting the wheel: an ace plays low at the
// bottom of it and the hand is a five high, not an ace high.
func isStraight(order []int, count map[int]int) (bool, int) {
	if len(order) != 5 {
		return false, 0
	}
	sorted := append([]int{}, order...)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	run := true
	for i := 1; i < len(sorted); i++ {
		run = run && sorted[i] == sorted[i-1]-1
	}
	if run {
		return true, sorted[0]
	}
	if sorted[0] == 14 && sorted[1] == 5 && sorted[2] == 4 && sorted[3] == 3 && sorted[4] == 2 {
		return true, 5
	}
	return false, 0
}

// BestAtTheTable is the best hand anybody else is holding.
func BestAtTheTable(g *CardGame) HandRank {
	best := HandRank{Category: -1}
	for _, s := range g.Seats {
		if r := Rank(s.Cards); best.Category < 0 || r.Beats(best) {
			best = r
		}
	}
	return best
}

// SitInTheBackRoom deals a hand against whoever is actually in the room and can
// cover the ante. Nobody is invented for this: an empty room has no game, which
// is why the game is worth walking somewhere for.
func (w *World) SitInTheBackRoom(place string, ante int) error {
	if place != BackRoom {
		return fmt.Errorf("there is no game here")
	}
	if w.Game != nil && !w.Game.Done {
		return fmt.Errorf("there is a hand on the table already")
	}
	if ante < MinAnte || ante > MaxAnte {
		return fmt.Errorf("the game runs between $%d and $%d a hand", MinAnte, MaxAnte)
	}
	if w.Player.Cash < ante {
		return fmt.Errorf("you cannot cover the ante")
	}
	var seats []Seat
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != place || w.Travelling(n) || len(seats) >= Players {
			continue
		}
		if w.Pockets(n) < ante {
			continue
		}
		seats = append(seats, Seat{Who: n.ID, Name: n.Name})
	}
	if len(seats) < 2 {
		return fmt.Errorf("there is nobody in the back room with money to lose")
	}
	if err := w.Pay(ante); err != nil {
		return err
	}
	g := &CardGame{Place: place, Ante: ante, Pot: ante, Seats: seats, Deck: w.deck()}
	g.Mine = g.take(5)
	for i := range g.Seats {
		g.Seats[i].Cards = g.take(5)
		if n := w.NPC(g.Seats[i].Who); n != nil {
			n.Purse -= ante
			g.Pot += ante
		}
	}
	w.Game = g
	w.Log("A game in the back room", fmt.Sprintf("$%d a head with %s. Five cards, one draw.", ante, listNames(g.Seats)), "personal")
	return nil
}

// take deals off the top of what is left.
func (g *CardGame) take(n int) []Card {
	out := g.Deck[:n]
	g.Deck = g.Deck[n:]
	return out
}

func listNames(seats []Seat) string {
	names := make([]string, 0, len(seats))
	for _, s := range seats {
		names = append(names, s.Name)
	}
	return joinNames(names)
}

// Draw changes the cards the player named, lets everybody else change theirs,
// and turns them over. Passing nothing is standing pat.
func (w *World) ChangeCards(discards []int) error {
	g := w.Game
	if g == nil || g.Done {
		return fmt.Errorf("there is no hand on the table")
	}
	if g.Drawn {
		return fmt.Errorf("the cards have been changed already")
	}
	if len(discards) > 3 {
		return fmt.Errorf("the house lets you change three")
	}
	throw := map[int]bool{}
	for _, i := range discards {
		if i < 0 || i >= len(g.Mine) {
			return fmt.Errorf("you are not holding a card there")
		}
		throw[i] = true
	}
	kept := make([]Card, 0, 5)
	for i, c := range g.Mine {
		if !throw[i] {
			kept = append(kept, c)
		}
	}
	g.Mine = append(kept, g.take(len(throw))...)
	for i := range g.Seats {
		g.Seats[i].Cards, g.Seats[i].Threw = g.change(w, g.Seats[i].Cards)
	}
	g.Drawn = true
	return w.showdown()
}

// change is how somebody who is not the player plays their hand: keep what is
// worth keeping and buy the rest. A pair draws three, two pair draws one,
// anything made stands pat, and nothing at all keeps its two highest cards.
//
// The house rule is three cards, and it is three for everybody. The first
// version let the room throw four while the player could throw three, which is
// not a house edge — there is no house — but it is the same thing wearing a
// different coat, and the money measured it at ten percent of the ante a hand.
func (g *CardGame) change(w *World, hand []Card) ([]Card, int) {
	r := Rank(hand)
	if r.Category >= Straight {
		return hand, 0
	}
	count := map[int]int{}
	for _, c := range hand {
		count[pokerRank(c)]++
	}
	// On nothing at all, keep the two highest: throwing the other three is the
	// most the room allows anybody.
	high := map[int]bool{}
	if r.Category == HighCard {
		ranked := make([]int, 0, len(hand))
		for _, c := range hand {
			ranked = append(ranked, pokerRank(c))
		}
		sort.Sort(sort.Reverse(sort.IntSlice(ranked)))
		high[ranked[0]], high[ranked[1]] = true, true
	}
	keep := make([]Card, 0, 5)
	for _, c := range hand {
		if count[pokerRank(c)] > 1 || high[pokerRank(c)] {
			keep = append(keep, c)
		}
	}
	threw := len(hand) - len(keep)
	return append(keep, g.take(threw)...), threw
}

// showdown turns everything over and moves the money. A pot goes to the best
// hand, or is split between every hand that good — the first version handed a
// tie to the room and gave the player their ante back, which is not a split and
// cost the player five and a half percent of the ante a hand over three
// thousand of them.
func (w *World) showdown() error {
	g := w.Game
	mine := Rank(g.Mine)
	best := mine
	for _, s := range g.Seats {
		if r := Rank(s.Cards); r.Beats(best) {
			best = r
		}
	}
	// Everybody holding the best hand, the player counted as seat -1.
	var winners []int
	if !best.Beats(mine) {
		winners = append(winners, -1)
	}
	for i, s := range g.Seats {
		if !best.Beats(Rank(s.Cards)) {
			winners = append(winners, i)
		}
	}
	share := g.Pot / len(winners)
	// A pot that will not divide leaves a dollar or two on the table. It goes
	// to the first winner clockwise, which is how it is settled in a room.
	rest := g.Pot - share*len(winners)
	names := make([]string, 0, len(winners))
	for k, seat := range winners {
		took := share
		if k == 0 {
			took += rest
		}
		if seat < 0 {
			w.Player.Cash += took
			g.Won = took - g.Ante
			names = append(names, "you")
			continue
		}
		if n := w.NPC(g.Seats[seat].Who); n != nil {
			n.Purse += took
		}
		names = append(names, g.Seats[seat].Name)
	}
	if len(winners) == 1 && winners[0] < 0 {
		g.Outcome = fmt.Sprintf("You had %s and took $%d off the table.", mine.Name(), g.Won)
	} else if len(winners) == 1 {
		g.Won = -g.Ante
		g.Outcome = fmt.Sprintf("%s had %s against your %s, and the table went to %s.", names[0], best.Name(), mine.Name(), names[0])
	} else {
		if g.Won == 0 {
			g.Won = -g.Ante
		}
		g.Outcome = fmt.Sprintf("%s all had %s, and $%d went each way.", joinNames(names), best.Name(), share)
	}
	g.Done = true
	w.Log("The hand in the back room", g.Outcome, "personal")
	return nil
}

// BackRoomAnte is what the player is staking: what they typed, or the smallest
// the game runs for.
func (w *World) BackRoomAnte(amount int) int {
	if amount <= 0 {
		return MinAnte
	}
	return amount
}

// BackRoomReadiness explains why there is no game to sit in on, or returns "".
func (w *World) BackRoomReadiness(id string, ante int) string {
	if id != BackRoom {
		return "There is no game here"
	}
	if w.Game != nil && !w.Game.Done {
		return "You are in the middle of a hand"
	}
	if ante < MinAnte || ante > MaxAnte {
		return fmt.Sprintf("The game runs between $%d and $%d a hand", MinAnte, MaxAnte)
	}
	if w.Player.Cash < ante {
		return "You cannot cover the ante"
	}
	if w.seatable(id, ante) < 2 {
		return "There is nobody in the back room with money to lose"
	}
	return ""
}

// seatable counts who could sit down for that ante.
func (w *World) seatable(id string, ante int) int {
	n := 0
	for i := range w.NPCs {
		p := &w.NPCs[i]
		if !p.Dead && p.Location == id && !w.Travelling(p) && w.Pockets(p) >= ante {
			n++
		}
	}
	return n
}

// CardsDescription is the game as the interface draws it, and nothing when
// nobody is playing.
func (w *World) CardsDescription() map[string]any {
	g := w.Game
	if g == nil {
		return nil
	}
	seats := make([]map[string]any, 0, len(g.Seats))
	for _, s := range g.Seats {
		seat := map[string]any{"who": s.Who, "name": s.Name, "threw": s.Threw}
		// Nobody sees a hand before it is turned over.
		if g.Done {
			seat["cards"] = s.Cards
			seat["hand"] = Rank(s.Cards).Name()
		}
		seats = append(seats, seat)
	}
	return map[string]any{
		"place": g.Place, "ante": g.Ante, "pot": g.Pot, "mine": g.Mine,
		"hand": Rank(g.Mine).Name(), "seats": seats, "drawn": g.Drawn,
		"done": g.Done, "outcome": g.Outcome, "won": g.Won,
	}
}
