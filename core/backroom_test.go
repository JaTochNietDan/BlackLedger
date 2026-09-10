package core

import (
	"sort"
	"strings"
	"testing"
)

// A back room is the one game in this city where the money across the table
// belongs to somebody with a name, an address and a reason to remember losing
// it. Every other table in Black Ledger is the house: an edge, a settlement,
// and nobody on the other side of it.

func backroom(t *testing.T) (*World, []string) {
	t.Helper()
	w := New(71)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health, w.Player.Location = 4000, 100, BackRoom
	seated := []string{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || len(seated) >= 3 {
			continue
		}
		n.Location, n.Heading, n.Purse = BackRoom, "", 900
		seated = append(seated, n.ID)
	}
	if len(seated) < 3 {
		t.Fatal("the city could not put three people in a room")
	}
	return w, seated
}

func TestTheBackRoomSeatsPeopleWhoLiveHere(t *testing.T) {
	w, seated := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatalf("nobody could get a game: %v", err)
	}
	g := w.Game
	if g == nil {
		t.Fatal("a game was dealt and the world does not have one")
	}
	if len(g.Seats) < 2 {
		t.Fatalf("a card game with %d other players is patience", len(g.Seats))
	}
	for _, s := range g.Seats {
		n := w.NPC(s.Who)
		if n == nil {
			t.Fatalf("somebody nobody in this city has ever met sat down: %q", s.Who)
		}
		if n.Location != BackRoom {
			t.Fatalf("%s is at %s and is playing cards at the poolhall", n.Name, n.Location)
		}
		found := false
		for _, id := range seated {
			found = found || id == s.Who
		}
		if !found {
			t.Fatalf("%s was dealt in and was never in the room", n.Name)
		}
	}
	if len(g.Mine) != 5 {
		t.Fatalf("a hand of draw poker is five cards, not %d", len(g.Mine))
	}
}

func TestEveryCardOnTheTableIsADifferentCard(t *testing.T) {
	for seed := uint32(1); seed <= 200; seed++ {
		w, _ := backroom(t)
		w.RNG = seed * 2654435761
		if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
			t.Fatalf("nobody could get a game: %v", err)
		}
		seen := map[Card]bool{}
		all := append([]Card{}, w.Game.Mine...)
		for _, s := range w.Game.Seats {
			all = append(all, s.Cards...)
		}
		for _, c := range all {
			if seen[c] {
				t.Fatalf("the %s of %s was dealt to two people at the same table", c.Rank, c.Suit)
			}
			seen[c] = true
		}
	}
}

func TestTheBestHandTakesThePot(t *testing.T) {
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatalf("nobody could get a game: %v", err)
	}
	before := w.Player.Cash
	pot := w.Game.Pot
	if pot <= 0 {
		t.Fatal("nobody put anything in")
	}
	if err := w.ChangeCards(nil); err != nil {
		t.Fatalf("standing pat was refused: %v", err)
	}
	g := w.Game
	if !g.Done {
		t.Fatal("the hand never finished")
	}
	best, mine := BestAtTheTable(g), Rank(g.Mine)
	if best.Beats(mine) && w.Player.Cash != before {
		t.Fatalf("a losing hand took %d out of the pot", w.Player.Cash-before)
	}
	if mine.Beats(best) && w.Player.Cash != before+pot {
		t.Fatalf("the best hand at the table won %d of a %d pot", w.Player.Cash-before, pot)
	}
}

func hand(s ...string) []Card {
	out := make([]Card, 0, len(s))
	suitOf := map[byte]string{'s': "spades", 'h': "hearts", 'd': "diamonds", 'c': "clubs"}
	for _, t := range s {
		r := t[:len(t)-1]
		v := 0
		for _, k := range ranks {
			if k.Rank == r {
				v = k.Value
			}
		}
		out = append(out, Card{Rank: r, Suit: suitOf[t[len(t)-1]], Value: v})
	}
	return out
}

func TestAHandIsWorthWhatItIsWorth(t *testing.T) {
	cases := []struct {
		what string
		cat  int
		hand []Card
	}{
		{"a straight flush", StraightFlush, hand("9h", "8h", "7h", "6h", "5h")},
		{"four of a kind", Quads, hand("9h", "9s", "9d", "9c", "5h")},
		{"a full house", FullHouse, hand("9h", "9s", "9d", "5c", "5h")},
		{"a flush", Flush, hand("Ah", "8h", "7h", "6h", "2h")},
		{"a straight", Straight, hand("9h", "8s", "7h", "6h", "5h")},
		{"the wheel", Straight, hand("Ah", "2s", "3h", "4h", "5h")},
		{"three of a kind", Trips, hand("9h", "9s", "9d", "4c", "5h")},
		{"two pair", TwoPair, hand("9h", "9s", "4d", "4c", "5h")},
		{"a pair", Pair, hand("9h", "9s", "J", "4c", "5h")[:0]},
		{"a pair", Pair, hand("9h", "9s", "Jd", "4c", "5h")},
		{"nothing at all", HighCard, hand("Kh", "9s", "Jd", "4c", "5h")},
	}
	for _, c := range cases {
		if len(c.hand) == 0 {
			continue
		}
		if got := Rank(c.hand); got.Category != c.cat {
			t.Errorf("%s came out as %s", c.what, got.Name())
		}
	}
	// The wheel is a five high, not an ace high, so it loses to a six high.
	wheel, six := Rank(hand("Ah", "2s", "3h", "4h", "5h")), Rank(hand("6h", "2s", "3h", "4h", "5c"))
	if !six.Beats(wheel) || wheel.Beats(six) {
		t.Error("an ace-low straight beat a six high")
	}
	// Same category, higher card.
	if !Rank(hand("Kh", "Ks", "2d", "3c", "4h")).Beats(Rank(hand("Qh", "Qs", "2d", "3c", "4h"))) {
		t.Error("kings did not beat queens")
	}
	// A full house is settled by the three, not by the two.
	if !Rank(hand("9h", "9s", "9d", "2c", "2h")).Beats(Rank(hand("8h", "8s", "8d", "Ac", "Ah"))) {
		t.Error("nines full of twos lost to eights full of aces")
	}
	// Nothing beats itself, which is what a split pot is.
	if Rank(hand("Kh", "Ks", "2d", "3c", "4h")).Beats(Rank(hand("Kd", "Kc", "2s", "3h", "4s"))) {
		t.Error("a hand beat the same hand")
	}
}

// A back room has no house in it. Over enough hands a player who plays the same
// way as everybody else at the table should come out level, because the only
// thing moving the money is the cards.
func TestTheBackRoomTakesNoRake(t *testing.T) {
	const hands = 3000
	table, mine := 0, 0
	for seed := uint32(1); seed <= hands; seed++ {
		w, seated := backroom(t)
		w.RNG = seed * 2654435761
		// Counted before the ante, or the antes are missing from one side of
		// the sum and the pot looks like money out of nowhere.
		cash, purses := w.Player.Cash, 0
		for _, id := range seated {
			purses += w.NPC(id).Purse
		}
		if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
			t.Fatalf("no game: %v", err)
		}
		// Play the hand the way the room plays it, so nothing but the cards is
		// different between the seats.
		keep := map[int]bool{}
		count := map[int]int{}
		for _, c := range w.Game.Mine {
			count[pokerRank(c)]++
		}
		made := Rank(w.Game.Mine).Category >= Straight
		high := map[int]bool{}
		if Rank(w.Game.Mine).Category == HighCard {
			var r []int
			for _, c := range w.Game.Mine {
				r = append(r, pokerRank(c))
			}
			sort.Sort(sort.Reverse(sort.IntSlice(r)))
			high[r[0]], high[r[1]] = true, true
		}
		var throw []int
		for i, c := range w.Game.Mine {
			held := made || count[pokerRank(c)] > 1 || high[pokerRank(c)]
			keep[i] = held
			if !held {
				throw = append(throw, i)
			}
		}
		if err := w.ChangeCards(throw); err != nil {
			t.Fatalf("the draw was refused: %v", err)
		}
		mine += w.Player.Cash - cash
		after := 0
		for _, id := range seated {
			after += w.NPC(id).Purse
		}
		table += after - purses
		if w.Player.Cash-cash+after-purses != 0 {
			t.Fatalf("$%d appeared at the table out of nowhere", w.Player.Cash-cash+after-purses)
		}
	}
	t.Logf("over %d hands the player is $%d up and the room is $%d up", hands, mine, table)
	if mine+table != 0 {
		t.Fatalf("the money at the table does not add up: %d against %d", mine, table)
	}
	// A quarter of the ante per hand either way is noise; a rake is not.
	if edge := float64(mine) / float64(hands*50); edge < -.06 || edge > .06 {
		t.Fatalf("a game with no house in it returned %.1f%% of the ante per hand", edge*100)
	}
}

// Random hands almost never tie — none did in four thousand — so the split has
// to be put on the table by hand or it is a branch nobody has ever run.
func TestATiedPotIsSplitRatherThanGivenAway(t *testing.T) {
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatalf("no game: %v", err)
	}
	g := w.Game
	before := w.Player.Cash
	purse := w.NPC(g.Seats[0].Who).Purse
	g.Mine = hand("Kh", "Ks", "9d", "5c", "3h")
	g.Seats[0].Cards = hand("Kd", "Kc", "9s", "5h", "3s")
	g.Seats[1].Cards = hand("2h", "7s", "9c", "Jc", "4h")
	g.Seats[2].Cards = hand("2s", "7h", "8c", "Jh", "4s")
	g.Drawn = true
	if err := w.showdown(); err != nil {
		t.Fatal(err)
	}
	half := g.Pot / 2
	if got := w.Player.Cash - before; got != half {
		t.Fatalf("half of a $%d pot came to $%d", g.Pot, got)
	}
	if got := w.NPC(g.Seats[0].Who).Purse - purse; got != half {
		t.Fatalf("the other half of a $%d pot came to $%d", g.Pot, got)
	}
	if !strings.Contains(g.Outcome, "each way") {
		t.Fatalf("a split pot was described as %q", g.Outcome)
	}
}

// The engine is not the game. An action nobody can press is a feature in a
// file, so this walks the same path the interface walks: the action is offered
// where the game is, the ante the player typed is the ante taken, and the two
// buttons of a hand both work.
func TestTheGameCanBePlayedThroughTheSamePathAsEverythingElse(t *testing.T) {
	w, _ := backroom(t)
	offered := func(kind string) *Action {
		for _, a := range w.Actions(w.Player.Location) {
			if a.ID == kind {
				return &a
			}
		}
		return nil
	}
	sit := offered("cards")
	if sit == nil {
		t.Fatal("the back room offers no game to sit in on")
	}
	if sit.Disabled {
		t.Fatalf("the game is refused in the room it is played in: %s", sit.Reason)
	}
	if sit.Group != "tables" {
		t.Fatalf("a card game is filed under %q", sit.Group)
	}
	before := w.Player.Cash
	if err := w.apply(Command{Kind: "cards", Amount: 120, RequestID: "backroomsitdown1"}); err != nil {
		t.Fatalf("sitting down was refused: %v", err)
	}
	if w.Game == nil || w.Game.Ante != 120 {
		t.Fatalf("a player who typed $120 is playing for $%d", w.Game.Ante)
	}
	if w.Player.Cash != before-120 {
		t.Fatalf("the ante took $%d", before-w.Player.Cash)
	}
	if offered("cards") != nil {
		t.Fatal("a second game was offered while a hand was on the table")
	}
	if a := offered("change"); a == nil || a.Disabled {
		t.Fatal("a hand was dealt and there is no way to change a card")
	}
	if err := w.apply(Command{Kind: "change", Choice: "0,2", RequestID: "backroomdrawtwo1"}); err != nil {
		t.Fatalf("changing two cards was refused: %v", err)
	}
	if !w.Game.Done {
		t.Fatal("the hand never came to a showdown")
	}
	if w.CardsDescription()["outcome"] == "" {
		t.Fatal("the hand finished and the interface is told nothing about it")
	}
}
