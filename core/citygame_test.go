package core

import (
	"strings"
	"testing"
)

// The back room is a game only the player has ever played in. The people
// standing in it have money, a reason to gamble it, and nothing else to do of
// an evening — and the room has an owner who ought to be earning from the
// table whether or not the player is at it. A business the player is not
// standing in should still be a business.

func TestTheCityPlaysCardsWithoutThePlayer(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	w := New(404)
	w.Event = nil
	for day := 0; day < 60; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	cards := 0
	for _, g := range w.Grudges {
		// The reason names the room now that there is more than one of them:
		// "a night at cards behind The Mariner". The wording changed with the
		// code, not the rule being guarded.
		if strings.HasPrefix(g.Because, "a night at cards behind ") {
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

// The game behind the poolhall was the only one in the city. A room with no
// house in it belongs wherever there are people of an evening and nobody
// holding a float, and the bar has four times the poolhall's evening crowd —
// thirty-two people against six — so a three-handed game gets up there on
// nights the poolhall's does not.
func TestTheCityPlaysInEveryBackRoom(t *testing.T) {
	t.Parallel()
	if len(BackRooms()) < 2 {
		t.Fatal("there is still only one room in this city with a game behind it")
	}
	// Every one of them has to be a room people are actually in after dark, or
	// it is a table nobody ever sits at.
	w := New(61)
	w.Event, w.District = nil, 9
	for _, at := range BackRooms() {
		if _, ok := PlaceByID(at); !ok {
			t.Fatalf("%s is not an address in this city", at)
		}
		if HasBankroll(at) {
			t.Fatalf("%s runs a float, so a game with no house in it is two games for the same seats", at)
		}
		evenings := 0
		for i := range w.NPCs {
			if !w.NPCs[i].Dead && haunt(w.NPCs[i].ID) == at {
				evenings++
			}
		}
		if evenings < 3 {
			t.Fatalf("%s has %d people in it of an evening, which will not make a three-handed game",
				at, evenings)
		}
	}
	for day := 0; day < 60; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	t.Logf("sixty days: %d hands in the city's back rooms, $%d taken for the seats",
		w.BackRoomHands, w.BackRoomTake)
	if w.BackRoomHands < 60 {
		t.Fatalf("two rooms and only %d hands in sixty days, which is one room's worth", w.BackRoomHands)
	}
}

// And the seat money stays in the room it was taken in. The player holds the
// poolhall here and not the bar, which the Russos have and which does not change
// hands — so one table fills and the other pays a rival, and neither pays the
// other. A hand played across town used to be worth money to the poolhall
// because the poolhall was the only room the code could name.
func TestTheSeatMoneyStaysInTheRoomItWasTakenIn(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Location, w.Player.Cash = BackRoom, 200000
	if err := w.apply(Command{Kind: "acquire", Target: BackRoom, RequestID: "holdthehall"}); err != nil {
		t.Fatal(err)
	}
	rival := w.Properties["bar"].Owner
	if w.Own("bar") {
		t.Fatal("the player holds the bar, so there is nobody else to pay")
	}
	theirs := 0
	if f := w.faction(rival); f != nil {
		theirs = f.Cash
	}
	for day := 0; day < 60; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	after := 0
	if f := w.faction(rival); f != nil {
		after = f.Cash
	}
	t.Logf("sixty days: the poolhall's table has taken $%d; the bar's holder is $%d different",
		w.Properties[BackRoom].Bankroll, after-theirs)
	if w.Properties[BackRoom].Bankroll == 0 {
		t.Fatal("the player holds the poolhall and its table has taken nothing")
	}
	if w.Properties["bar"].Bankroll != 0 {
		t.Fatal("a room the player does not hold is filling a table of its own rather than paying its holder")
	}
}

// And the city says where the games are. While there was one of them nobody
// needed telling which room it was; with two, a player who has always played
// behind the poolhall has no reason to walk into a bar and find out there is a
// table there as well.
func TestTheCitySaysWhichRoomsHaveAGameBehindThem(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	said := map[string]bool{}
	for _, l := range w.Public()["locations"].([]map[string]any) {
		id := l["id"].(string)
		if l["back_room"].(bool) != HasBackRoom(id) {
			t.Fatalf("%s: the payload says back room %v and the core says %v",
				id, l["back_room"], HasBackRoom(id))
		}
		if strings.Contains(l["note"].(string), "game in the back") {
			said[id] = true
		}
	}
	for _, at := range BackRooms() {
		if !said[at] {
			place, _ := PlaceByID(at)
			t.Fatalf("%s has a game behind it and nothing walking past says so", place.Name)
		}
	}
	if len(said) != len(BackRooms()) {
		t.Fatalf("%d rooms claim a game and %d have one", len(said), len(BackRooms()))
	}
}

// And a room the player holds says it too, because holding the Green Baize
// should not hide what the Green Baize is.
func TestARoomYouHoldStillSaysItHasAGameBehindIt(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Location, w.Player.Cash = BackRoom, 200000
	if err := w.apply(Command{Kind: "acquire", Target: BackRoom, RequestID: "holdandsay"}); err != nil {
		t.Fatal(err)
	}
	if note := w.PlaceNote(BackRoom); !strings.Contains(note, "game in the back") {
		t.Fatalf("the player holds the room with the game in it and the room says %q", note)
	}
}
