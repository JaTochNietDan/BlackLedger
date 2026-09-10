package core

import "fmt"

// The back room was a game only the player had ever played in. The people
// standing in it have money, a reason to gamble it and nothing else to do of an
// evening, and the room has an owner who ought to be earning from the table
// whether or not the player is at it. A business the player is not standing in
// is still a business, which is the same principle as the rest of this city:
// every actor runs on the same rules, and what happens off-screen is a
// transaction the simulation actually performed.

const (
	// TableCharge is what the room takes off each player for the seat. Not a
	// rake: the pot stays exactly what everybody put in, which is the whole
	// point of a game with no house in it. The house is renting a table.
	TableCharge = 8
	// CityAnte is what the city plays for among themselves. Small — these are
	// people's wages, not a family's money.
	CityAnte = 25
	// SoreAtEachOther is what a bad night at the table is worth between two
	// people who are not the player. It goes through the city's own grudges,
	// which is the machinery that already exists for two people falling out.
	SoreAtEachOther = 9
)

// BackRoomNight plays one hand among whoever spent the evening in the back
// room. Called once a day from the clock, like every other thing the city does
// to itself — which is at midnight, by which time everybody has gone home. So
// this asks whose evening the poolhall is rather than who is standing in it at
// the moment the day turns over: the first version asked the second question
// and played one hand in thirty days.
func (w *World) BackRoomNight() {
	var players []*NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead && !w.Travelling(n) && haunt(n.ID) == BackRoom && w.Pockets(n) >= CityAnte+TableCharge {
			players = append(players, n)
		}
		if len(players) >= Players+1 {
			break
		}
	}
	if len(players) < 3 {
		return
	}
	pot, deck := 0, w.worldDeck()
	hands := make([][]Card, len(players))
	for i, n := range players {
		n.Purse -= CityAnte
		pot += CityAnte
		hands[i] = deck[i*5 : i*5+5]
		w.chargeForTheTable(n)
	}
	// Everybody plays their own hand the way the room plays a hand. Nobody
	// bets: what is interesting off-screen is who won and who is short, and a
	// betting round nobody watches is arithmetic for its own sake.
	at := len(players) * 5
	for i, n := range players {
		kept, threw := keepFor(hands[i], n.Skill)
		hands[i] = append(kept, deck[at:at+threw]...)
		at += threw
	}
	best, who := HandRank{Category: -1}, 0
	for i := range players {
		if r := Rank(hands[i]); r.Beats(best) {
			best, who = r, i
		}
	}
	players[who].Purse += pot
	w.BackRoomHands++
	// The player hears about it if it is their room or they were in it. A game
	// nobody told them about is still a game, and the city is full of them.
	if w.Own(BackRoom) || w.Player.Location == BackRoom {
		w.Log("The back room", w.BackRoomNote(players[who], pot), "personal")
	}
	// A bad night between two people who are not the player is the city's own
	// business, and the city already has somewhere to put it. One person a
	// night at most, and only somebody who cannot sit down again: falling out
	// with every loser put nine of this city's sixteen grudges down to a card
	// game inside two months, which is a city whose quarrels are all about
	// poker.
	worst := -1
	for i, n := range players {
		if i == who || w.Pockets(n) >= CityAnte {
			continue
		}
		if worst < 0 || n.Purse < players[worst].Purse {
			worst = i
		}
	}
	if worst >= 0 {
		w.Resent(players[worst].ID, players[who].ID, SoreAtEachOther, "a night at cards in the back room")
	}
}

// keepFor plays a hand the way somebody of that skill plays it, and reports
// what they kept and how many they bought.
func keepFor(cards []Card, skill int) ([]Card, int) {
	if Rank(cards).Category >= Straight {
		return append([]Card{}, cards...), 0
	}
	count := map[int]int{}
	for _, c := range cards {
		count[pokerRank(c)]++
	}
	best := 0
	for _, c := range cards {
		if pokerRank(c) > best {
			best = pokerRank(c)
		}
	}
	kept := make([]Card, 0, 5)
	for _, c := range cards {
		// Somebody who cannot read a hand keeps their highest card and hopes,
		// which is how a poor player loses money to a good one at the same
		// table on the same cards.
		if count[pokerRank(c)] > 1 || (skill < ReadsIt && pokerRank(c) == best) {
			kept = append(kept, c)
		}
	}
	if len(kept) == 0 {
		kept = append(kept, cards[0])
	}
	if len(kept) > 5 {
		kept = kept[:5]
	}
	return kept, 5 - len(kept)
}

// chargeForTheTable takes the seat money and gives it to whoever holds the
// room. An address nobody holds keeps it, which is what an unowned back room
// does with it.
func (w *World) chargeForTheTable(n *NPC) {
	if w.Pockets(n) < TableCharge {
		return
	}
	n.Purse -= TableCharge
	w.BackRoomTake += TableCharge
	prop := w.Properties[BackRoom]
	if prop == nil {
		return
	}
	// The room's money, not the holder's pocket. Somebody managing the Green
	// Baize should be able to look at what the table has taken, add to it and
	// draw it out, which they cannot do with money that has already gone into
	// their own hands.
	if w.Own(BackRoom) {
		prop.Bankroll += TableCharge
		return
	}
	if f := w.faction(prop.Owner); f != nil {
		f.Cash += TableCharge
	}
}

// worldDeck is a shuffled deck for a hand nobody watches. It runs on a stream
// of its own, derived from the world's seed and the day, rather than off either
// of the city's two.
//
// The first version shuffled off WorldRNG, which is fifty-one draws a night
// taken out from under every other thing the city does off-screen. Nothing was
// wrong with the cards; everything downstream of them moved. It showed up in a
// balance test on a knife edge — a car netted $11,482 against $11,618 on foot,
// where it had been ahead — and that is a real thing to know about the car,
// but it is not a thing a card game in a poolhall should be able to do.
func (w *World) worldDeck() []Card {
	cards := make([]Card, 0, 52)
	for _, r := range ranks {
		for _, s := range suits {
			cards = append(cards, Card{Rank: r.Rank, Suit: s, Value: r.Value})
		}
	}
	seed := w.RNG ^ 0x517cc1b7 ^ uint32(w.Minute/1440)*2654435761
	next := func() float64 {
		seed = 1664525*seed + 1013904223
		return float64(seed) / 4294967296
	}
	for i := len(cards) - 1; i > 0; i-- {
		j := min(i, int(next()*float64(i+1)))
		cards[i], cards[j] = cards[j], cards[i]
	}
	return cards
}

// BackRoomNote is what the player hears about a game they were not in, when
// they are somewhere they would hear it.
func (w *World) BackRoomNote(winner *NPC, pot int) string {
	return fmt.Sprintf("%s had the best of it in the back room at the poolhall last night, and walked out with $%d of it.", winner.Name, pot)
}
