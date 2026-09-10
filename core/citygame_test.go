package core

import "testing"

// The back room is a game only the player has ever played in. The people
// standing in it have money, a reason to gamble it, and nothing else to do of
// an evening — and the room has an owner who ought to be earning from the
// table whether or not the player is at it. A business the player is not
// standing in should still be a business.

func TestTheCityPlaysCardsWithoutThePlayer(t *testing.T) {
	w := New(404)
	w.Event = nil
	moved, played := 0, 0
	before := map[string]int{}
	for i := range w.NPCs {
		before[w.NPCs[i].ID] = w.NPCs[i].Purse
	}
	for day := 0; day < 30; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Purse != before[n.ID] {
			moved++
		}
	}
	played = w.BackRoomHands
	t.Logf("over thirty days the city played %d hands in the back room and %d people's pockets are a different size", played, moved)
	if played == 0 {
		t.Fatal("in a month nobody in this city ever played a hand of cards without the player dealing themselves in")
	}
	if moved == 0 {
		t.Fatal("the city played cards and nobody's money moved")
	}
	// Nobody plays with money they have not got.
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Purse < 0 {
			t.Fatalf("%s is carrying $%d", n.Name, n.Purse)
		}
	}
}

// And the room takes something for the table, which is what a back room is for
// from the owner's side.
func TestTheHouseTakesSomethingForTheTable(t *testing.T) {
	w := New(404)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, BackRoom)
	earned := w.Player.Earned
	for day := 0; day < 30; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	t.Logf("thirty days: %d hands, $%d taken for the table, and the holder is $%d better off all told",
		w.BackRoomHands, w.BackRoomTake, w.Player.Earned-earned)
	if w.BackRoomHands == 0 {
		t.Fatal("nobody played a hand in a month")
	}
	// The poolhall earns whether or not anybody plays cards, so the question is
	// what the table itself is worth rather than what the room made.
	if w.BackRoomTake == 0 {
		t.Fatal("a month of card games in a room you own took nothing for the table")
	}
	if w.Player.Earned-earned < w.BackRoomTake {
		t.Fatalf("the table took $%d and the holder is only $%d better off", w.BackRoomTake, w.Player.Earned-earned)
	}
}

// A hand nobody watches must not move the cards the player is about to be
// dealt. Two streams, and the city's own games run on the city's.
func TestTheCitysOwnCardsDoNotTouchThePlayersDeal(t *testing.T) {
	deal := func(cityPlays bool) []Card {
		w := New(404)
		w.Event, w.District = nil, 9
		w.Player.Cash, w.Player.Health, w.Player.Location = 4000, 100, BackRoom
		w.RNG, w.WorldRNG = 99, 99
		seated := 0
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if n.Dead || seated >= 3 {
				continue
			}
			n.Location, n.Heading, n.Purse = BackRoom, "", 900
			seated++
		}
		if cityPlays {
			w.BackRoomNight()
		}
		if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
			t.Fatalf("no game: %v", err)
		}
		return w.Game.Mine
	}
	quiet, busy := deal(false), deal(true)
	for i := range quiet {
		if quiet[i] != busy[i] {
			t.Fatalf("a game the player was not in changed what they were dealt: %s of %s became %s of %s",
				quiet[i].Rank, quiet[i].Suit, busy[i].Rank, busy[i].Suit)
		}
	}
}

// And somebody who has a bad night against somebody else falls out with them,
// through the city's own machinery for two people falling out rather than a
// second one written for cards.
func TestABadNightAtCardsIsSomethingTwoPeopleFallOutOver(t *testing.T) {
	w := New(404)
	w.Event = nil
	for day := 0; day < 60; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	cards := 0
	for _, g := range w.Grudges {
		if g.Because == "a night at cards in the back room" {
			cards++
		}
	}
	t.Logf("after sixty days, %d of %d grudges in this city are about a card game", cards, len(w.Grudges))
	if cards == 0 {
		t.Fatal("sixty days of card games and nobody fell out with anybody over one")
	}
	// And not a city whose only quarrels are about poker. Asked as "there are
	// other reasons people fall out here" rather than as a share: the city
	// holds seven to sixteen grudges at a time, and a ratio on a denominator
	// that small says more about the sample than about the game.
	other := len(w.Grudges) - cards
	if other == 0 {
		t.Fatalf("every one of this city's %d grudges is about a card game", len(w.Grudges))
	}
}
