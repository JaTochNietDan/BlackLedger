package core

import "fmt"

// Sitting down at the tables was one roll and a paragraph. The house kept its
// edge, the money moved, and the player made no decision after the one to sit
// down — which is the only decision that does not belong to a gambler.
//
// A hand is a decision the player actually makes, several times, with the same
// arithmetic anybody at a table faces: what they are showing, what the dealer is
// showing, and whether one more card is worth it. The house's edge comes from
// the rules of the room rather than from a weighted table: the dealer draws to
// sixteen, stands on seventeen, and takes ties.

const (
	// Bust is the total above which a hand is dead.
	Bust = 21
	// DealerStands is the total the dealer draws to and no further.
	DealerStands = 17
	// NaturalPays is what a two-card twenty-one returns, as a percentage of the
	// stake, on top of the stake itself.
	NaturalPays = 150
)

// TableHand is a session at the tables that has not finished. It lives on the world
// because it is real: money has been staked and the cards are on the table.
type TableHand struct {
	// Place is the room. Stake was one of two fixed lots by name and is kept
	// only so a hand already on the table in an older save can be read; Down is
	// what is actually on the cloth, in dollars, because a player picks the
	// number now.
	Place, Stake string
	Down         int `json:"down,omitempty"`
	// Player and Dealer are the totals showing.
	Player, Dealer int
	// Cards is how many the player has taken, so a natural can be told from a
	// twenty-one reached the hard way.
	Cards int
	// Done is set when the hand has been settled, so nothing settles twice.
	Done bool
	// What was actually dealt. The hand used to be two totals and a count, so
	// there was no such thing in this game as the seven of clubs — which meant
	// a table drawing playing cards would have been drawing cards nobody had
	// dealt. Absent in saves written before the deck had faces, which read as a
	// hand whose cards are not known; the totals are still true of it.
	Mine   []Card `json:"mine"`
	Theirs []Card `json:"theirs"`
	// What happened, kept on the hand after it is settled. A finished hand used
	// to stop being described at all: the felt vanished the instant the dealer
	// turned their card over, so the one moment the player was waiting for —
	// what the dealer actually had — was the one thing the game never showed
	// them. It stays on the table until the next hand is dealt.
	Won     bool   `json:"won,omitempty"`
	Outcome string `json:"outcome,omitempty"`
}

// Card is one card off the deck: what it says and what suit it is. Value is
// carried with it so that what a hand is worth and what is lying on the table
// can never disagree — the total is added up from these and from nowhere else.
type Card struct {
	Rank  string `json:"rank"`
	Suit  string `json:"suit"`
	Value int    `json:"value"`
}

// ranks are what a card can say, in order, with what each is worth. An ace is
// eleven here and softened later if the hand would otherwise be dead.
var ranks = []struct {
	Rank  string
	Value int
}{
	{"A", 11}, {"2", 2}, {"3", 3}, {"4", 4}, {"5", 5}, {"6", 6}, {"7", 7},
	{"8", 8}, {"9", 9}, {"10", 10}, {"J", 10}, {"Q", 10}, {"K", 10},
}

// suits are the four of them, named as the interface will draw them.
var suits = []string{"spades", "hearts", "diamonds", "clubs"}

// card is what comes off the deck. Face cards are ten and an ace is eleven
// unless that busts, which is the one piece of arithmetic the house does for
// you.
//
// The suit is drawn from the same stream as the rank because it is part of the
// same card. It changes nothing about the odds — every suit is worth the same —
// but a card without one is not a card.
func (w *World) card() Card {
	r := ranks[min(len(ranks)-1, int(w.Random()*float64(len(ranks))))]
	suit := suits[min(len(suits)-1, int(w.Random()*float64(len(suits))))]
	return Card{Rank: r.Rank, Suit: suit, Value: r.Value}
}

// Total is what a row of cards is worth, with aces brought down one at a time
// while the hand would otherwise be dead. It is the only place a hand's value
// is worked out, so the number on the screen is always the sum of the cards
// beside it.
func Total(cards []Card) int {
	total, aces := 0, 0
	for _, c := range cards {
		total += c.Value
		if c.Rank == "A" {
			aces++
		}
	}
	for total > Bust && aces > 0 {
		total -= 10
		aces--
	}
	return total
}

// Deal starts a hand. The stake is taken now, because it is on the table now.
func (w *World) Deal(id string, amount int) error {
	stake := Stake{ID: "typed", Label: "A hand", Amount: amount}
	if reason := w.TableReadiness(id, stake); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if w.Hand != nil && !w.Hand.Done {
		return fmt.Errorf("there is a hand on the table already")
	}
	if err := w.Pay(stake.Amount); err != nil {
		return err
	}
	mine := []Card{w.card(), w.card()}
	theirs := []Card{w.card()}
	w.Hand = &TableHand{Place: id, Down: amount, Mine: mine, Theirs: theirs,
		Player: Total(mine), Dealer: Total(theirs), Cards: len(mine)}
	place, _ := PlaceByID(id)
	w.Log("A hand at "+place.Name, fmt.Sprintf("$%d down. You are showing %d and the dealer is showing %d.", stake.Amount, w.Hand.Player, w.Hand.Dealer), "personal")
	return nil
}

// DrawCard takes another card. Going over is the end of it.
func (w *World) DrawCard() error {
	if w.Hand == nil || w.Hand.Done {
		return fmt.Errorf("there is no hand on the table")
	}
	w.Hand.Mine = append(w.Hand.Mine, w.card())
	w.Hand.Player = Total(w.Hand.Mine)
	w.Hand.Cards = len(w.Hand.Mine)
	if w.Hand.Player > Bust {
		return w.settleHand(false, "You went over.")
	}
	return nil
}

// Stand leaves it with the dealer, who draws to sixteen and stands on
// seventeen, and who takes ties. That last rule is where the house's edge lives.
func (w *World) Stand() error {
	if w.Hand == nil || w.Hand.Done {
		return fmt.Errorf("there is no hand on the table")
	}
	for w.Hand.Dealer < DealerStands {
		w.Hand.Theirs = append(w.Hand.Theirs, w.card())
		w.Hand.Dealer = Total(w.Hand.Theirs)
	}
	if w.Hand.Dealer > Bust {
		return w.settleHand(true, fmt.Sprintf("The dealer went over on %d.", w.Hand.Dealer))
	}
	if w.Hand.Player > w.Hand.Dealer {
		return w.settleHand(true, fmt.Sprintf("%d against %d.", w.Hand.Player, w.Hand.Dealer))
	}
	if w.Hand.Player == w.Hand.Dealer {
		return w.push(fmt.Sprintf("%d each. Nobody's night.", w.Hand.Player))
	}
	return w.settleHand(false, fmt.Sprintf("%d against %d.", w.Hand.Player, w.Hand.Dealer))
}

// settleHand pays or does not, and moves the money through the same books a
// session at the tables always moved it through.
func (w *World) settleHand(won bool, why string) error {
	hand := w.Hand
	hand.Done = true
	stake := Stake{Amount: w.handDown(hand)}
	place, _ := PlaceByID(hand.Place)
	house := w.faction(w.Properties[hand.Place].Owner)

	returned := 0
	if won {
		returned = stake.Amount * 2
		if hand.Player == Bust && hand.Cards == 2 {
			returned = stake.Amount + stake.Amount*NaturalPays/100
		}
		w.Earn(returned)
	}
	net := returned - stake.Amount
	w.changeBusinessFunds(hand.Place, -net)
	w.tableAftermath(place.Name, house, net, stake.Amount)
	hand.Won = won
	if won {
		hand.Outcome = fmt.Sprintf("%s $%d comes back across the table.", why, returned)
		w.Log("Paid out at "+place.Name, hand.Outcome, "business")
	} else {
		hand.Outcome = fmt.Sprintf("%s The $%d stays where it is.", why, stake.Amount)
		w.Log("Gone at "+place.Name, hand.Outcome, "business")
	}
	return nil
}

// push returns the stake on a tie and leaves the room none the wiser.
func (w *World) push(why string) error {
	hand := w.Hand
	hand.Done = true
	stake := Stake{Amount: w.handDown(hand)}
	place, _ := PlaceByID(hand.Place)
	w.Earn(stake.Amount)
	w.Log("A stand-off at "+place.Name, why+" Your money comes back and nobody is any better off.", "business")
	return nil
}

// HandDescription is what is on the table, for the interface.
func (w *World) HandDescription() map[string]any {
	if w.Hand == nil {
		return map[string]any{"playing": false}
	}
	place, _ := PlaceByID(w.Hand.Place)
	stake := Stake{Amount: w.handDown(w.Hand)}
	// The cards themselves, so a table can draw what was dealt rather than a
	// number. Never nil: a row the interface has to guard is a row it will
	// eventually forget to guard.
	mine, theirs := w.Hand.Mine, w.Hand.Theirs
	if mine == nil {
		mine = []Card{}
	}
	if theirs == nil {
		theirs = []Card{}
	}
	return map[string]any{
		"playing": !w.Hand.Done, "settled": w.Hand.Done, "won": w.Hand.Won,
		"outcome": w.Hand.Outcome, "place": place.Name, "stake": stake.Amount,
		"player": w.Hand.Player, "dealer": w.Hand.Dealer, "cards": w.Hand.Cards,
		"mine": mine, "theirs": theirs, "where": w.Hand.Place,
	}
}

// handDown is what is on the cloth. A hand written before players could name a
// number carries one of the two old lots by name instead, and reads as that.
func (w *World) handDown(hand *TableHand) int {
	if hand == nil {
		return 0
	}
	if hand.Down > 0 {
		return hand.Down
	}
	if stake, ok := tableStake(hand.Stake); ok {
		return stake.Amount
	}
	return 0
}
