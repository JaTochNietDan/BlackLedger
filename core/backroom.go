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
	Drawn bool `json:"drawn,omitempty"`
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
		seats = append(seats, Seat{Who: n.ID, Name: n.Name, Had: w.Pockets(n)})
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
	w.Log("A game in the back room", fmt.Sprintf("$%d a head with %s. Five cards, one draw.%s",
		ante, listNames(g.Seats), w.tableRemembers(g)), "personal")
	return nil
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
	// The hands are made. What they are worth is now a question of who will
	// pay to see them, which is the half of this game the cards do not decide.
	w.Log("The draw", w.drawNote(), "personal")
	return nil
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

// showdown turns over everybody still in and moves the money. A pot goes to the
// best hand, or is split between every hand that good — the first version
// handed a tie to the room and gave the player their ante back, which is not a
// split and cost the player five and a half percent of the ante a hand over
// three thousand of them.
func (w *World) showdown() error {
	g := w.Game
	mine := Rank(g.Mine)
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
		r := Rank(s.Cards)
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
			w.Player.Cash += took
			g.Won = took - g.Ante - g.MyBet
			names = append(names, "you")
			continue
		}
		if n := w.NPC(g.Seats[seat].Who); n != nil {
			n.Purse += took
		}
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
		had, moved := s.Had, n.Purse-s.Had
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
			if n.Purse == 0 {
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
		seat := map[string]any{"who": s.Who, "name": s.Name, "threw": s.Threw,
			"in": s.In, "folded": s.Folded, "said": s.Said}
		// What this hand did to them, and whether they were already carrying
		// something about the last one. A table where the player has taken
		// money off somebody twice should look different from one where
		// everybody has just sat down.
		if n := w.NPC(s.Who); n != nil {
			seat["sore"] = n.Sore
			if g.Done {
				seat["moved"] = n.Purse - s.Had
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
		"hand": Rank(g.Mine).Name(), "seats": seats, "drawn": g.Drawn,
		"bet": g.Bet, "my_bet": g.MyBet, "facing": g.Facing, "folded": g.Folded,
		"done": g.Done, "outcome": g.Outcome, "won": g.Won,
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
	r := Rank(cards)
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
	if !g.Drawn {
		return fmt.Errorf("the cards have not been changed yet")
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
	if amount > w.Player.Cash {
		return fmt.Errorf("you cannot cover that")
	}
	if amount > 0 {
		if err := w.Pay(amount); err != nil {
			return err
		}
		g.Bet, g.MyBet, g.Pot = amount, amount, g.Pot+amount
	}
	w.roundOfBetting()
	if g.Bet > g.MyBet {
		g.Facing = true
		return nil
	}
	return w.showdown()
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
		strength := seatStrength(s.Cards, n.Skill)
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
				put := min(g.Ante*2, w.Pockets(n))
				if put <= 0 {
					continue
				}
				n.Purse -= put
				s.In += put
				g.Bet, g.Pot = s.In, g.Pot+put
				s.Said = "bets $" + fmt.Sprint(put)
				if bluff && strength < 4 {
					s.Said = "bets $" + fmt.Sprint(put) + " without blinking"
				}
			}
			continue
		}
		if w.Pockets(n) < owed {
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
			put := min(owed+g.Ante*2, w.Pockets(n))
			n.Purse -= put
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
		n.Purse -= owed
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
	if owed > w.Player.Cash {
		return fmt.Errorf("you cannot cover that")
	}
	if err := w.Pay(owed); err != nil {
		return err
	}
	g.MyBet, g.Pot, g.Facing = g.Bet, g.Pot+owed, false
	// Whoever called the smaller figure has to match the bigger one or get out.
	w.roundOfBetting()
	return w.showdown()
}

// FoldHand is the player throwing it in rather than paying. What is already in
// the pot stays in the pot, which is the whole of what folding costs.
func (w *World) FoldHand() error {
	g := w.Game
	if g == nil || g.Done {
		return fmt.Errorf("there is no hand on the table")
	}
	if !g.Drawn {
		return fmt.Errorf("you have not seen your hand yet")
	}
	g.Folded, g.Facing = true, false
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
