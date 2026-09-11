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

// BackRoom is where the game started. A poolhall has a room behind the room.
const BackRoom = "poolhall"

// And it is not the only one. A game with no house in it belongs where there
// are people of an evening and nobody holding a float: a room that runs tables
// is already taking an edge off everybody in it, and a room that is empty after
// dark has nobody to deal to. That is the poolhall and the bar — and the bar
// has four times the poolhall's evening crowd, so the city's own game there
// gets up on nights the poolhall's does not.
//
// The club and the casino are left out on purpose. They run floats, the house
// edge is the point of them, and a no-house game in the same building would be
// two games competing for the same seats.
var backRooms = []string{BackRoom, "bar"}

// HasBackRoom reports whether this address has a room behind the room.
func HasBackRoom(id string) bool {
	for _, at := range backRooms {
		if at == id {
			return true
		}
	}
	return false
}

// BackRooms is every address with a game behind it, for the city's own night
// and for the guards.
func BackRooms() []string { return append([]string{}, backRooms...) }

// Ante bounds. Low enough that anybody in the room can sit down, high enough
// that losing is felt.
const (
	MinAnte = 10
	MaxAnte = 500
	// Players is how many other people sit down, at most.
	Players = 3
)

// A sitting, rather than a hand.
//
// The table dealt one hand and stopped. To play a second you got up and sat
// down again, which re-seated the room, re-read everybody's pockets and threw
// away everything the last hand had meant: "the game should continue until you
// stop playing, right now it just requires you to leave the table and rejoin.
// Realistically it feels like it should be more like actual poker, where you
// have a buy in and whatnot and you play until people go bust or you can
// leave."
//
// So there is money on the table now. Everybody who sits down puts a stake in
// front of them and plays out of it; the pot is chips, not pockets. What that
// buys is the two things a poker night is made of and neither of which a
// one-hand table can have: you can lose what you brought without losing what
// you own, and somebody can be cleaned out and leave.
const (
	// MinBuyIn and MaxBuyIn are what can be put on the table.
	MinBuyIn = 50
	MaxBuyIn = 10000
	// AntesInAStack is how many antes a full buy-in is worth. Twenty, because a
	// stake that cannot lose twenty hands is not a stake, it is one hand with
	// extra steps.
	AntesInAStack = 20
	// SitsOutUnder is the least anybody sits down with, in antes. Below this
	// they are not playing, they are waiting to be blinded off.
	SitsOutUnder = 4
	// LeastAnte is the smallest a hand is ever played for. The people in these
	// rooms carry fifty or seventy-five dollars, so a table pitched at what the
	// player brought would be a table nobody in the city could sit at.
	LeastAnte = 5
	// SitsDownWith is the least in somebody's pocket before they will take a
	// seat at all.
	SitsDownWith = 20
)

// TableAnte is the stake at a table: a twentieth of what the player put up, but
// never more than a quarter of what the shortest stack at the table can cover,
// because the game is played against the people in the room rather than against
// the player's bankroll. A rich man at a poor table plays for what the table
// plays for and sits there a long time, which is what a back room is.
func TableAnte(buyIn, shortest int) int {
	ante := buyIn / AntesInAStack
	if shortest > 0 {
		ante = min(ante, shortest/SitsOutUnder)
	}
	return max(LeastAnte, min(MaxAnte, ante))
}

// Seat is somebody at the table who is not the player.
type Seat struct {
	Who   string `json:"who"`
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
	// Threw is left from draw poker, where it was how many cards somebody
	// bought and the only thing you could read off them before the showdown. In
	// hold'em nobody changes a card: what you read is what they do with the
	// money. Kept at zero so a save written before the game changed still
	// loads.
	Threw int `json:"threw"`
	// Had is what was in their pocket before the ante, so what a night at the
	// table did to somebody can be told from what they were carrying. Measuring
	// it at the showdown instead measures nothing: by then the ante and every
	// bet have already gone, and a seat that simply is not paid has lost
	// nothing between one line and the next.
	Had int `json:"had,omitempty"`
	// In is what this seat has put in beyond the ante, and Folded is whether
	// they threw their hand in rather than pay to see yours.
	In     int  `json:"in,omitempty"`
	Folded bool `json:"folded,omitempty"`
	// Said is what they did when the money went round, kept so the interface
	// can show a table talking rather than a row of totals.
	Said string `json:"said,omitempty"`
	// Stack is what this seat has in front of them. Everything at this table is
	// played out of it: the ante, every bet, and what a pot pays. When it will
	// not cover the ante they are cleaned out and they leave, which is the
	// thing a table of one hand could never do.
	Stack int `json:"stack"`
	// Out is set on the hand somebody is cleaned out, so the table can say so
	// once before the seat goes.
	Out bool `json:"out,omitempty"`
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
	// Board is what is face up in the middle: three, then one, then one. Never
	// omitted: the screen reads its length, and a board that arrives as null
	// before the flop is a blank screen the moment anybody sits down.
	Board []Card `json:"board"`
	// Street is how far the hand has got — preflop, flop, turn, river, shown.
	Street string `json:"street,omitempty"`
	// Bet is what it costs to stay in beyond the ante, and Mine is what the
	// player has already put toward it. Raised is set once somebody has put
	// the money up a second time, because this room allows one raise: a table
	// that can raise for ever is a table nobody can write a decision for.
	Bet    int  `json:"bet,omitempty"`
	MyBet  int  `json:"my_bet,omitempty"`
	Raised bool `json:"raised,omitempty"`
	// Facing is set when the money has come back at the player and they owe
	// something to see the hands.
	Facing bool `json:"facing,omitempty"`
	// Folded is the player having thrown their hand in.
	Folded  bool   `json:"folded,omitempty"`
	Done    bool   `json:"done,omitempty"`
	Outcome string `json:"outcome,omitempty"`
	Won     int    `json:"won,omitempty"`
	// Stack is the player's chips, BuyIn is what they put on the table when
	// they sat down, and Hands is how many have been dealt since. Over is the
	// sitting finished rather than the hand: everybody else cleaned out, or the
	// player with nothing left to ante.
	Stack int  `json:"stack"`
	BuyIn int  `json:"buy_in,omitempty"`
	Hands int  `json:"hands,omitempty"`
	Over  bool `json:"over,omitempty"`
	// Ended says why, in words, when it is over.
	Ended string `json:"ended,omitempty"`
	// Left is everybody who has been cleaned out at this table tonight. They do
	// not come back to it: being cleaned out is going home, not sitting out a
	// hand.
	Left []string `json:"left"`
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
		if r := BestOfSeven(s.Cards, g.Board); best.Category < 0 || r.Beats(best) {
			best = r
		}
	}
	return best
}

// SitInTheBackRoom deals a hand against whoever is actually in the room and can
// cover the ante. Nobody is invented for this: an empty room has no game, which
// is why the game is worth walking somewhere for.
func (w *World) SitInTheBackRoom(place string, buyIn int) error {
	if !HasBackRoom(place) {
		return fmt.Errorf("there is no game here")
	}
	if w.Game != nil && !w.Game.Done {
		return fmt.Errorf("there is a hand on the table already")
	}
	if buyIn < MinBuyIn || buyIn > MaxBuyIn {
		return fmt.Errorf("the table takes between $%d and $%d on it", MinBuyIn, MaxBuyIn)
	}
	if w.Player.Cash < buyIn {
		return fmt.Errorf("you cannot put that on the table")
	}
	// Everybody at the table plays out of what they put in front of them: what
	// they can stand to, up to what the player has put up. The stake follows
	// from the shortest of those, so the game is the room's game.
	var seats []Seat
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != place || w.Travelling(n) || len(seats) >= Players {
			continue
		}
		if w.Pockets(n) < SitsDownWith {
			continue
		}
		seats = append(seats, Seat{Who: n.ID, Name: n.Name, Had: w.Pockets(n),
			Stack: min(w.Pockets(n), buyIn)})
	}
	if len(seats) < 2 {
		return fmt.Errorf("there is nobody in the back room with money to lose")
	}
	shortest := seats[0].Stack
	for _, s := range seats {
		shortest = min(shortest, s.Stack)
	}
	ante := TableAnte(buyIn, shortest)
	for i := range seats {
		if n := w.NPC(seats[i].Who); n != nil {
			n.Purse -= seats[i].Stack
		}
	}
	if err := w.Pay(buyIn); err != nil {
		return err
	}
	g := &CardGame{Place: place, Ante: ante, Seats: seats, BuyIn: buyIn, Stack: buyIn}
	w.Game = g
	// And the player is at the table. Buying in from the room's own list left
	// them standing in the room with a hand of cards going on somewhere the
	// screen could not show: "I set the amount I want to buy in but it just ran
	// a simulation instead of letting me play the game." Putting money on a
	// table is sitting down at it, whichever button was pressed.
	w.Seated, w.SeatedTo = place, Backroom
	w.Log("A seat in the back room",
		fmt.Sprintf("$%d on the table at $%d a hand, against %s. You play out of what is in front of you and take home what is left.%s",
			buyIn, ante, listNames(g.Seats), w.tableRemembers(g)), "personal")
	return w.DealAgain()
}

// DealAgain puts the next hand out. It is the same work whether it is the first hand
// of a sitting or the ninth, which is the whole of what "the game should
// continue until you stop playing" asks for: the table between hands is a table
// with people and money still at it.
func (w *World) DealAgain() error {
	g := w.Game
	if g == nil {
		return fmt.Errorf("you are not at the table")
	}
	if g.Over {
		return fmt.Errorf("the game is over")
	}
	if g.Hands > 0 && !g.Done {
		return fmt.Errorf("there is a hand on the table already")
	}
	// Whoever cannot cover the ante is cleaned out. They take what is left in
	// front of them back to their pocket and go.
	kept := make([]Seat, 0, len(g.Seats))
	var gone []string
	for _, s := range g.Seats {
		if s.Stack >= g.Ante {
			kept = append(kept, s)
			continue
		}
		if n := w.NPC(s.Who); n != nil {
			n.Purse += s.Stack
		}
		gone = append(gone, s.Name)
		g.Left = append(g.Left, s.Who)
	}
	g.Seats = kept
	if len(gone) > 0 {
		w.Log("Cleaned out at the table",
			fmt.Sprintf("%s %s nothing left in front of %s and left the room.",
				joinNames(gone), was(len(gone)), them(len(gone))), "personal")
	}
	if g.Stack < g.Ante {
		return w.endSitting("You have nothing left in front of you.")
	}
	// And whoever else is in the room takes the empty chair. A back room does
	// not close because one man went home; it closes when there is nobody left
	// in it. Without this a table at the bar, where seven people spend their
	// evening, ended after nine hands because the first three were cleaned out.
	w.fillTheTable(g)
	if len(g.Seats) < 2 {
		return w.endSitting("There is nobody left in the room with money to play for.")
	}
	// A fresh hand: the felt is cleared and the antes go in.
	g.Mine, g.Board, g.Deck = nil, []Card{}, w.deck()
	g.Street, g.Bet, g.MyBet, g.Raised = Preflop, 0, 0, false
	g.Facing, g.Folded, g.Done, g.Outcome, g.Won = false, false, false, "", 0
	g.Pot = g.Ante
	g.Stack -= g.Ante
	g.Mine = g.take(2)
	for i := range g.Seats {
		s := &g.Seats[i]
		s.Cards, s.In, s.Folded, s.Said, s.Out = g.take(2), 0, false, "", false
		s.Stack -= g.Ante
		g.Pot += g.Ante
	}
	g.Hands++
	return nil
}

// endSitting closes the table and gives everybody back what is in front of
// them. Nobody walks away from a back room leaving their chips on the baize.
func (w *World) endSitting(why string) error {
	g := w.Game
	if g == nil || g.Over {
		return nil
	}
	took := g.Stack
	w.Player.Cash += took
	g.Stack = 0
	for i := range g.Seats {
		if n := w.NPC(g.Seats[i].Who); n != nil {
			n.Purse += g.Seats[i].Stack
		}
		g.Seats[i].Stack = 0
	}
	g.Over, g.Done, g.Ended = true, true, why
	up := took - g.BuyIn
	how := fmt.Sprintf("You put $%d on the table and picked $%d up", g.BuyIn, took)
	switch {
	case up > 0:
		how = fmt.Sprintf("You put $%d on the table and picked $%d up, $%d to the good", g.BuyIn, took, up)
	case up < 0:
		how = fmt.Sprintf("You put $%d on the table and picked $%d up, $%d of it gone", g.BuyIn, took, -up)
	}
	w.Log("Up from the table",
		fmt.Sprintf("%s. %s %s", why, how, plural(g.Hands, "hand", "hands")+" played."), "personal")
	return nil
}

// was and them are the grammar of a line that names one person or several.
func was(n int) string {
	if n == 1 {
		return "has"
	}
	return "have"
}
func them(n int) string {
	if n == 1 {
		return "them"
	}
	return "them"
}

// tableRemembers is what the player is told when they sit down against people
// they have played before. The city already knows who is carrying something
// against them; this is the one table where they find out before the money
// goes in rather than afterwards.
func (w *World) tableRemembers(g *CardGame) string {
	var holding []string
	for _, s := range g.Seats {
		if n := w.NPC(s.Who); grudging(n) {
			holding = append(holding, n.Name)
		}
	}
	if len(holding) == 0 {
		return ""
	}
	if len(holding) == 1 {
		return " " + holding[0] + " has not forgotten the last time."
	}
	return " " + joinNames(holding) + " have not forgotten the last time."
}

// cards is a row of cards that is never nil, because the screen counts it and
// a nil slice reaches the browser as null rather than as an empty row.
func cards(in []Card) []Card {
	if in == nil {
		return []Card{}
	}
	return in
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

// Streets. A betting round ends when nobody is short of the bet; then the next
// cards go face up in the middle and the money goes round again. Four rounds,
// which is four times the decision draw poker gave the player.

// nextRound puts the next street on the table and clears what everybody has in
// front of them, or calls the hand if there is nothing left to come.
func (w *World) nextRound() error {
	g := w.Game
	street, cards := nextStreet(g.Street)
	if street == Shown {
		return w.showdown()
	}
	g.Street = street
	g.Board = append(g.Board, g.take(cards)...)
	// The money in front of people goes into the pot at the end of a street;
	// what is owed on the next one starts at nothing.
	g.Bet, g.MyBet, g.Raised, g.Facing = 0, 0, false, false
	for i := range g.Seats {
		g.Seats[i].In, g.Seats[i].Said = 0, ""
	}
	// Everybody but the player may bet into the new street before it comes back
	// to them.
	w.roundOfBetting()
	if g.Bet > g.MyBet {
		g.Facing = true
	}
	return nil
}

// stillIn counts the hands that have not been thrown away, the player's
// included. One left is a hand that is over without anybody showing anything.
func (w *World) stillIn() int {
	g := w.Game
	n := 0
	if !g.Folded {
		n++
	}
	for _, s := range g.Seats {
		if !s.Folded {
			n++
		}
	}
	return n
}

// showdown turns over everybody still in and moves the money. A pot goes to the
// best hand, or is split between every hand that good — the first version
// handed a tie to the room and gave the player their ante back, which is not a
// split and cost the player five and a half percent of the ante a hand over
// three thousand of them.
func (w *World) showdown() error {
	g := w.Game
	mine := BestOfSeven(g.Mine, g.Board)
	// Everybody who paid to be here.
	var winners []int
	best := HandRank{Category: -1}
	if !g.Folded {
		best = mine
		winners = []int{-1}
	}
	for i, s := range g.Seats {
		if s.Folded {
			continue
		}
		r := BestOfSeven(s.Cards, g.Board)
		switch {
		case r.Beats(best):
			best, winners = r, []int{i}
		case !best.Beats(r):
			winners = append(winners, i)
		}
	}
	if len(winners) == 0 {
		// Everybody threw their hand in, which cannot happen while the player
		// is in it and can if the player folded to a table that then folded to
		// nobody. The money sits with whoever put in most.
		winners = []int{0}
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
			g.Stack += took
			g.Won = took - g.Ante - g.MyBet
			names = append(names, "you")
			continue
		}
		g.Seats[seat].Stack += took
		names = append(names, g.Seats[seat].Name)
	}
	if g.Won == 0 && (g.Folded || winners[0] >= 0) {
		g.Won = -g.Ante - g.MyBet
	}
	switch {
	case g.Folded:
		g.Outcome = fmt.Sprintf("You threw in %s and it cost you $%d.", mine.Name(), -g.Won)
	case len(winners) == 1 && winners[0] < 0:
		g.Outcome = fmt.Sprintf("You had %s and took $%d off the table.", mine.Name(), g.Won)
	case len(winners) == 1:
		g.Outcome = fmt.Sprintf("%s had %s against your %s, and the table went to %s.", names[0], best.Name(), mine.Name(), names[0])
	default:
		g.Outcome = fmt.Sprintf("%s all had %s, and $%d went each way.", joinNames(names), best.Name(), share)
	}
	g.Done = true
	w.remember()
	w.Log("The hand in the back room", g.Outcome, "personal")
	return nil
}

// SoreAtCards is the most somebody holds against a player who cleaned them out
// in one hand, and BearsIt is the share of what they had that a loss has to be
// before they carry it at all. A tenner off nine hundred is a tenner; a night's
// money is a reason to remember a face.
const (
	SoreAtCards = 34
	BearsIt     = .2
)

// remember is what the table takes away from the hand. The point of a game with
// no house in it is that the money belongs to somebody: a man who loses a
// night's money to you across a table has a reason to remember you, and a man
// you paid has a reason to like you. Before this, both walked away with nothing
// on their mind, which made the back room a slot machine with faces on it.
func (w *World) remember() {
	g := w.Game
	for _, s := range g.Seats {
		n := w.NPC(s.Who)
		if n == nil || n.Dead {
			continue
		}
		// What the night did to them, not what is in their pocket. The money
		// is in front of them while they are sitting at the table, so reading
		// the pocket alone says a man with four hundred in chips was cleaned
		// out and a man who won a pot went home with nothing.
		had, moved := s.Had, n.Purse+s.Stack-s.Had
		switch {
		case moved > 0:
			// Money you handed over is goodwill, and more of it is more of it.
			n.Trust = min(100, n.Trust+1+moved/(g.Ante*4))
		case moved < 0 && had > 0:
			share := float64(-moved) / float64(had)
			if share < BearsIt {
				continue
			}
			if share > 1 {
				share = 1
			}
			weight := int(float64(SoreAtCards) * share)
			because := "the night you took them at cards"
			if n.Purse+s.Stack == 0 {
				because = "the night you cleaned them out at cards"
			}
			w.Aggrieve(n.ID, weight, because)
		}
	}
}

// BackRoomAnte is what the player is staking: what they typed, or the smallest
// the game runs for.
func (w *World) BackRoomAnte(amount int) int {
	if amount <= 0 {
		return MinAnte
	}
	return amount
}

// BackRoomBuyIn is what the player is putting on the table: what they named, or
// a sensible default they can afford. A twentieth of it is the ante, so this
// one figure sets the size of the game.
func (w *World) BackRoomBuyIn(id string, amount int) int {
	if amount > 0 {
		return amount
	}
	want := MinBuyIn * 2
	if w.Player.Cash < want {
		want = MinBuyIn
	}
	return min(want, max(MinBuyIn, w.Player.Cash))
}

// DealReadiness explains why the next hand cannot go out, or returns "".
func (w *World) DealReadiness() string {
	g := w.Game
	switch {
	case g == nil:
		return "You are not at the table"
	case g.Over:
		return "The game is over"
	case !g.Done:
		return "There is a hand on the table already"
	case g.Stack < g.Ante:
		return "You have nothing left in front of you"
	case len(g.Seats) < 2:
		return "There is nobody left at the table"
	}
	return ""
}

// BackRoomReadiness explains why there is no game to sit in on, or returns "".
func (w *World) BackRoomReadiness(id string, buyIn int) string {
	if !HasBackRoom(id) {
		return "There is no game here"
	}
	if w.Game != nil && !w.Game.Done {
		return "You are in the middle of a hand"
	}
	if buyIn < MinBuyIn || buyIn > MaxBuyIn {
		return fmt.Sprintf("The table takes between $%d and $%d on it", MinBuyIn, MaxBuyIn)
	}
	if w.Player.Cash < buyIn {
		return "You cannot put that on the table"
	}
	if w.seatable(id, SitsDownWith) < 2 {
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
		seat := map[string]any{"who": s.Who, "name": s.Name, "threw": s.Threw,
			"in": s.In, "folded": s.Folded, "said": s.Said, "stack": s.Stack}
		// What this hand did to them, and whether they were already carrying
		// something about the last one. A table where the player has taken
		// money off somebody twice should look different from one where
		// everybody has just sat down.
		if n := w.NPC(s.Who); n != nil {
			seat["sore"] = n.Sore
			if g.Done {
				// What the night has done to them, not the hand: their chips
				// plus whatever is left in their pocket, against what they
				// walked in with. Reading the pocket alone says a man with
				// four hundred in front of him is broke.
				seat["moved"] = n.Purse + s.Stack - s.Had
			}
		}
		// Nobody sees a hand before it is turned over.
		if g.Done {
			seat["cards"] = s.Cards
			seat["hand"] = Rank(s.Cards).Name()
		}
		seats = append(seats, seat)
	}
	return map[string]any{
		"place": g.Place, "ante": g.Ante, "pot": g.Pot, "mine": g.Mine,
		"hand": BestOfSeven(g.Mine, g.Board).Name(), "seats": seats,
		"board": cards(g.Board), "street": g.Street, "street_name": StreetName(g.Street),
		"bet": g.Bet, "my_bet": g.MyBet, "facing": g.Facing, "folded": g.Folded,
		"done": g.Done, "outcome": g.Outcome, "won": g.Won,
		// The sitting, as opposed to the hand: what is in front of the player,
		// what they brought, how many hands they have played, and whether the
		// night is finished.
		"stack": g.Stack, "buy_in": g.BuyIn, "hands": g.Hands,
		"over": g.Over, "ended": g.Ended, "up": g.Stack - g.BuyIn,
	}
}

// Betting. The draw is arithmetic anybody can do; the money is where the people
// at the table stop being a distribution and start being people. A man who
// raises on nothing twice a night is somebody the player learns to call, and a
// man who never puts a dollar in without the hand to back it is somebody they
// learn to believe.

const (
	// OneRaise is the house rule. A table that can raise for ever is a table
	// nobody can write a decision for, and this room plays a single raise.
	OneRaise = true
	// BluffNerve is the ambition above which somebody in this room will put
	// money on a hand that cannot win.
	BluffNerve = 55
	// ReadsIt is the skill above which somebody folds a hand that is behind
	// rather than paying to be shown it.
	ReadsIt = 45
	// TakesItPersonally is how sore somebody has to be at the player before it
	// changes how they sit down against them. A memory that stays in a field is
	// worth nothing: a man you took a night's money off pays to see your hand
	// rather than believing you, because what he wants back is what you took.
	TakesItPersonally = 15
)

// grudging reports whether this person is playing the player rather than the
// cards. It is the same number everything else in the city reads to decide who
// would move against them.
func grudging(n *NPC) bool { return n != nil && n.Sore >= TakesItPersonally }

// drawNote is what the player is told when the cards have been changed: what
// they are holding and what everybody else bought, which is all the reading
// anybody gets before the money goes round.
func (w *World) drawNote() string {
	g := w.Game
	said := make([]string, 0, len(g.Seats))
	for _, s := range g.Seats {
		switch s.Threw {
		case 0:
			said = append(said, s.Name+" stood pat")
		case 1:
			said = append(said, s.Name+" took one")
		default:
			said = append(said, fmt.Sprintf("%s took %d", s.Name, s.Threw))
		}
	}
	return fmt.Sprintf("You are holding %s. %s.", Rank(g.Mine).Name(), joinNames(said))
}

// strength is how good a hand looks to the person holding it. Skill is how well
// they read it: somebody who cannot tell a pair of threes from a pair of
// queens plays them the same way.
func seatStrength(cards []Card, skill int) int {
	return handStrength(Rank(cards), skill)
}

// tableStrength is what somebody makes of their two cards and the board.
func tableStrength(hole, board []Card, skill int) int {
	return handStrength(BestOfSeven(hole, board), skill)
}

// handStrength is how good a made hand looks to the person holding it.
func handStrength(r HandRank, skill int) int {
	s := r.Category * 2
	if r.Category == Pair && len(r.Order) > 0 && r.Order[0] >= 11 && skill >= ReadsIt {
		s++ // a pair worth playing, if you can tell
	}
	if r.Category == HighCard && len(r.Order) > 0 && r.Order[0] == 14 && skill >= ReadsIt {
		s++
	}
	return s
}

// PlaceBet is the player putting money in, or checking with nothing. Everybody
// else answers it, in the order they are sitting.
func (w *World) PlaceBet(amount int) error {
	g := w.Game
	if g == nil || g.Done {
		return fmt.Errorf("there is no hand on the table")
	}
	if g.Facing {
		return fmt.Errorf("the bet is back with you: call it or throw the hand in")
	}
	if g.Bet > 0 {
		return fmt.Errorf("the betting is finished")
	}
	if amount < 0 || amount > MaxAnte {
		return fmt.Errorf("the room takes up to $%d on one bet", MaxAnte)
	}
	// Out of what is in front of you. Nobody at a table reaches into their
	// coat: what you can bet is what you brought, which is the whole point of
	// buying in.
	if amount > g.Stack {
		return fmt.Errorf("you have only $%d in front of you", g.Stack)
	}
	if amount > 0 {
		g.Stack -= amount
		g.Bet, g.MyBet, g.Pot = amount, amount, g.Pot+amount
	}
	w.roundOfBetting()
	if g.Bet > g.MyBet {
		g.Facing = true
		return nil
	}
	if w.stillIn() < 2 {
		return w.showdown()
	}
	return w.nextRound()
}

// roundOfBetting walks the table until nobody is short. Anybody short of the
// bet calls it, folds or, once in a hand, puts it up again — and somebody
// checked to may bet, which is why one pass is not enough: the first version
// skipped every seat whose stake already equalled a bet of nothing, so nobody
// in this room ever bet into a checked pot and the betting round moved no money
// at all. Three passes is the most a single raise can need.
func (w *World) roundOfBetting() {
	for pass := 0; pass < 3; pass++ {
		if !w.bettingPass() {
			return
		}
	}
}

// bettingPass walks the table once and reports whether anybody did anything.
func (w *World) bettingPass() bool {
	g := w.Game
	moved := false
	for i := range g.Seats {
		s := &g.Seats[i]
		if s.Folded || (s.In >= g.Bet && (g.Bet > 0 || s.Said != "")) {
			continue
		}
		moved = true
		n := w.NPC(s.Who)
		if n == nil {
			s.Folded, s.Said = true, "is not at the table any more"
			continue
		}
		if s.In >= g.Bet && g.Bet > 0 {
			continue
		}
		strength := tableStrength(s.Cards, g.Board, n.Skill)
		owed := g.Bet - s.In
		// Nothing to call: a hand worth showing bets it, and so does somebody
		// with the nerve to represent one they have not got.
		if owed == 0 {
			bluff := n.Ambition >= BluffNerve && w.Random() < .18
			if grudging(n) {
				bluff = bluff || w.Random() < .12
			}
			s.Said = "checks"
			if strength >= 4 || bluff {
				put := min(g.Ante*2, s.Stack)
				if put <= 0 {
					continue
				}
				s.Stack -= put
				s.In += put
				g.Bet, g.Pot = s.In, g.Pot+put
				s.Said = "bets $" + fmt.Sprint(put)
				if bluff && strength < 4 {
					s.Said = "bets $" + fmt.Sprint(put) + " without blinking"
				}
			}
			continue
		}
		if s.Stack < owed {
			s.Folded, s.Said = true, "has not got it and throws the hand in"
			continue
		}
		// A raise, once, and only from somebody holding something and willing
		// to say so.
		// Somebody with something against the player puts it up on less, and
		// does not need the nerve for it: wanting their money back is the
		// nerve.
		up := strength >= 6 && n.Ambition >= BluffNerve
		if grudging(n) {
			up = up || strength >= 4
		}
		if !g.Raised && up {
			put := min(owed+g.Ante*2, s.Stack)
			s.Stack -= put
			s.In += put
			g.Bet, g.Pot, g.Raised = s.In, g.Pot+put, true
			s.Said = "puts it up $" + fmt.Sprint(s.In-(g.Bet-put))
			s.Said = fmt.Sprintf("sees it and puts it up to $%d", s.In)
			continue
		}
		// Otherwise: pay to see it, or read that it is beaten and stop paying.
		// What it takes to make somebody put the hand down. The first version
		// folded anything under two pair to any bet worth more than the ante,
		// which made the whole table fold to a bet and turned betting a made
		// hand into a way of buying the antes: measured over three thousand
		// hands, betting the good ones was $1,800 worse than checking them.
		// A pair somebody can read is worth a call.
		folds := strength <= 1 || (strength <= 2 && owed > g.Ante*2) ||
			(strength <= 3 && n.Skill < ReadsIt && owed > g.Ante*3)
		// A man who wants his money back pays to see the hand. He is wrong more
		// often for it, which is what makes him worth betting into.
		if grudging(n) && strength >= 1 {
			folds = false
		}
		if folds {
			s.Folded, s.Said = true, "throws the hand in"
			continue
		}
		s.Stack -= owed
		s.In += owed
		g.Pot += owed
		s.Said = "calls"
	}
	return moved
}

// CallBet is the player paying to see the hands after somebody put it up.
func (w *World) CallBet() error {
	g := w.Game
	if g == nil || g.Done || !g.Facing {
		return fmt.Errorf("there is nothing to call")
	}
	owed := g.Bet - g.MyBet
	if owed > g.Stack {
		return fmt.Errorf("you have only $%d in front of you", g.Stack)
	}
	g.Stack -= owed
	g.MyBet, g.Pot, g.Facing = g.Bet, g.Pot+owed, false
	// Whoever called the smaller figure has to match the bigger one or get out.
	w.roundOfBetting()
	if g.Bet > g.MyBet {
		g.Facing = true
		return nil
	}
	if w.stillIn() < 2 {
		return w.showdown()
	}
	return w.nextRound()
}

// FoldHand is the player throwing it in rather than paying. What is already in
// the pot stays in the pot, which is the whole of what folding costs.
func (w *World) FoldHand() error {
	g := w.Game
	if g == nil || g.Done {
		return fmt.Errorf("there is no hand on the table")
	}
	g.Folded, g.Facing = true, false
	// A hand nobody is left contesting is over: the pot goes to whoever is
	// still in it without anybody having to show anything.
	return w.showdown()
}

// CallReadiness explains why the player cannot pay to see the hands.
func (w *World) CallReadiness() string {
	g := w.Game
	if g == nil || g.Done || !g.Facing {
		return "There is nothing to call"
	}
	if owed := g.Bet - g.MyBet; owed > w.Player.Cash {
		return fmt.Sprintf("You are $%d short of calling it", owed-w.Player.Cash)
	}
	return ""
}

// TableAnteAt is what a hand would cost at this room for this buy-in, read off
// who is actually in there. The offer says it before the player commits,
// because "a twentieth of what you put up" is not the whole truth at a table
// where somebody is carrying fifty dollars.
func (w *World) TableAnteAt(id string, buyIn int) int {
	shortest := 0
	seats := 0
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != id || w.Travelling(n) || seats >= Players {
			continue
		}
		if w.Pockets(n) < SitsDownWith {
			continue
		}
		seats++
		stake := min(w.Pockets(n), buyIn)
		if shortest == 0 || stake < shortest {
			shortest = stake
		}
	}
	return TableAnte(buyIn, shortest)
}

// fillTheTable sits down whoever else is in the room, up to the table's size.
// Somebody who has already been cleaned out here tonight does not come back:
// they went home, which is what being cleaned out means.
func (w *World) fillTheTable(g *CardGame) {
	if len(g.Seats) >= Players {
		return
	}
	seated := map[string]bool{}
	for _, s := range g.Seats {
		seated[s.Who] = true
	}
	for _, who := range g.Left {
		seated[who] = true
	}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if len(g.Seats) >= Players {
			return
		}
		if n.Dead || seated[n.ID] || n.Location != g.Place || w.Travelling(n) {
			continue
		}
		stake := min(w.Pockets(n), g.BuyIn)
		if stake < g.Ante*SitsOutUnder {
			continue
		}
		n.Purse -= stake
		g.Seats = append(g.Seats, Seat{Who: n.ID, Name: n.Name, Had: w.Pockets(n) + stake, Stack: stake})
		place, _ := PlaceByID(g.Place)
		w.Log("Somebody takes the empty chair",
			fmt.Sprintf("%s sits down at %s with $%d in front of them.", n.Name, place.Name, stake), "personal")
	}
}
