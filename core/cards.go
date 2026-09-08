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
	// Place is the room, and Stake what is down.
	Place, Stake string
	// Player and Dealer are the totals showing.
	Player, Dealer int
	// Cards is how many the player has taken, so a natural can be told from a
	// twenty-one reached the hard way.
	Cards int
	// Done is set when the hand has been settled, so nothing settles twice.
	Done bool
}

// card is what comes off the deck. Face cards are ten and an ace is eleven
// unless that busts, which is the one piece of arithmetic the house does for
// you.
func (w *World) card() int {
	value := 1 + int(w.Random()*13)
	if value > 10 {
		return 10
	}
	if value == 1 {
		return 11
	}
	return value
}

// soften brings an eleven down to a one when the hand would otherwise be dead.
func soften(total, drawn int) int {
	if total > Bust && drawn == 11 {
		return total - 10
	}
	return total
}

// Deal starts a hand. The stake is taken now, because it is on the table now.
func (w *World) Deal(id, stakeID string) error {
	stake, ok := tableStake(stakeID)
	if !ok {
		return fmt.Errorf("no such game")
	}
	if reason := w.TableReadiness(id, stake); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if w.Hand != nil && !w.Hand.Done {
		return fmt.Errorf("there is a hand on the table already")
	}
	if err := w.Pay(stake.Amount); err != nil {
		return err
	}
	first, second := w.card(), w.card()
	total := first + second
	if total > Bust {
		total -= 10 // two aces
	}
	w.Hand = &TableHand{Place: id, Stake: stakeID, Player: total, Dealer: w.card(), Cards: 2}
	place, _ := PlaceByID(id)
	w.Log("A hand at "+place.Name, fmt.Sprintf("$%d down. You are showing %d and the dealer is showing %d.", stake.Amount, total, w.Hand.Dealer), "personal")
	return nil
}

// DrawCard takes another card. Going over is the end of it.
func (w *World) DrawCard() error {
	if w.Hand == nil || w.Hand.Done {
		return fmt.Errorf("there is no hand on the table")
	}
	drawn := w.card()
	w.Hand.Player = soften(w.Hand.Player+drawn, drawn)
	w.Hand.Cards++
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
		drawn := w.card()
		w.Hand.Dealer = soften(w.Hand.Dealer+drawn, drawn)
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
	stake, _ := tableStake(hand.Stake)
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
	if house != nil {
		house.Cash = max(0, house.Cash-net)
	}
	w.tableAftermath(place.Name, house, net, stake.Amount)
	if won {
		w.Log("Paid out at "+place.Name, fmt.Sprintf("%s $%d comes back across the table.", why, returned), "business")
	} else {
		w.Log("Gone at "+place.Name, fmt.Sprintf("%s The $%d stays where it is.", why, stake.Amount), "business")
	}
	return nil
}

// push returns the stake on a tie and leaves the room none the wiser.
func (w *World) push(why string) error {
	hand := w.Hand
	hand.Done = true
	stake, _ := tableStake(hand.Stake)
	place, _ := PlaceByID(hand.Place)
	w.Earn(stake.Amount)
	w.Log("A stand-off at "+place.Name, why+" Your money comes back and nobody is any better off.", "business")
	return nil
}

// HandDescription is what is on the table, for the interface.
func (w *World) HandDescription() map[string]any {
	if w.Hand == nil || w.Hand.Done {
		return map[string]any{"playing": false}
	}
	place, _ := PlaceByID(w.Hand.Place)
	stake, _ := tableStake(w.Hand.Stake)
	return map[string]any{
		"playing": true, "place": place.Name, "stake": stake.Amount,
		"player": w.Hand.Player, "dealer": w.Hand.Dealer, "cards": w.Hand.Cards,
	}
}
