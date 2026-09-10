package core

import (
	"fmt"
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

// playOut checks the hand through every street: no bet from the player, and pay
// whatever comes back at them. The cheapest way to reach a showdown from a test.
func playOut(t *testing.T, w *World) {
	t.Helper()
	for i := 0; i < 12 && !w.Game.Done; i++ {
		if w.Game.Facing {
			if err := w.CallBet(); err != nil {
				t.Fatalf("calling was refused: %v", err)
			}
			continue
		}
		if err := w.PlaceBet(0); err != nil {
			t.Fatalf("checking was refused: %v", err)
		}
	}
	if !w.Game.Done {
		t.Fatal("a hand of hold'em never reached a showdown")
	}
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
	if len(g.Mine) != 2 {
		t.Fatalf("a hand of hold'em is two cards, not %d", len(g.Mine))
	}
	if g.Street != Preflop {
		t.Fatalf("a hand was dealt and the street is %q", g.Street)
	}
	if len(g.Board) != 0 {
		t.Fatalf("%d cards were face up before anybody had bet", len(g.Board))
	}
}

func TestEveryCardOnTheTableIsADifferentCard(t *testing.T) {
	for seed := uint32(1); seed <= 200; seed++ {
		w, _ := backroom(t)
		w.RNG = seed * 2654435761
		if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
			t.Fatalf("nobody could get a game: %v", err)
		}
		playOut(t, w)
		seen := map[Card]bool{}
		all := append([]Card{}, w.Game.Mine...)
		all = append(all, w.Game.Board...)
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
	playOut(t, w)
	g := w.Game
	if !g.Done {
		t.Fatal("the hand never finished")
	}
	best, mine := BestAtTheTable(g), BestOfSeven(g.Mine, g.Board)
	if best.Beats(mine) && w.Player.Cash > before {
		t.Fatalf("a losing hand took %d out of the pot", w.Player.Cash-before)
	}
	if mine.Beats(best) && w.Player.Cash < before+pot {
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
		playOut(t, w)
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
}

// The draw alone is close to an equal share of the pots, which is what having
// no house in the game means. Once there is money to answer, calling every bet
// with any hand is the way to lose it, and folding what is beaten is the way to
// stop. Both numbers are here because the first is the game being fair and the
// second is the game being a game.
func TestFoldingIsWorthMoreThanTheCardsAre(t *testing.T) {
	run := func(fold, bet bool) int {
		total := 0
		for seed := uint32(1); seed <= 3000; seed++ {
			w, _ := backroom(t)
			w.RNG = seed * 2654435761
			cash := w.Player.Cash
			if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
				t.Fatalf("no game: %v", err)
			}
			// Four streets of the same decision, which is the whole of what
			// hold'em asks and four times what the draw asked.
			for i := 0; i < 12 && !w.Game.Done; i++ {
				// What "beaten" means depends on the street. Before the flop
				// nobody has a hand yet, so folding everything under three of
				// a kind there is folding every hand in the game — which is
				// what the first version of this test did, and it made folding
				// look like the losing policy.
				rank := BestOfSeven(w.Game.Mine, w.Game.Board)
				made := rank.Category >= Trips
				worthPaying := rank.Category >= Pair || w.Game.Street == Preflop
				if w.Game.Facing {
					// Somebody put money in. Pay to see it, or believe them.
					var err error
					if fold && !worthPaying {
						err = w.FoldHand()
					} else {
						err = w.CallBet()
					}
					if err != nil {
						t.Fatalf("answering the bet was refused: %v", err)
					}
					continue
				}
				put := 0
				if bet && made {
					put = 100
				}
				if err := w.PlaceBet(put); err != nil {
					t.Fatalf("betting was refused: %v", err)
				}
			}
			total += w.Player.Cash - cash
		}
		return total
	}
	calls, folds, plays := run(false, false), run(true, false), run(true, true)
	t.Logf("over 3000 hands at $50: calling everything is $%d, folding what is beaten is $%d, folding and betting the good ones is $%d", calls, folds, plays)
	if folds <= calls {
		t.Fatalf("throwing beaten hands in cost more than paying for them: %d against %d", folds, calls)
	}
	if plays <= folds {
		t.Fatalf("putting money on a made hand earned nothing: %d against %d", plays, folds)
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
	// The same hand for two people, which hold'em produces far more often than
	// draw poker did: the board is most of everybody's hand.
	g.Board = hand("Kh", "Kd", "9s", "5h", "3s")
	g.Mine = hand("Ah", "Qc")
	g.Seats[0].Cards = hand("Ad", "Qs")
	g.Seats[1].Cards = hand("2h", "7s")
	g.Seats[2].Cards = hand("2s", "7h")
	g.Street = River
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
	if a := offered("bet"); a == nil || a.Disabled {
		t.Fatal("a hand was dealt and there is no way to put money on it")
	}
	// Four streets, checked and called through the same path the interface
	// uses. Each command needs an id of its own or the store refuses the repeat.
	for i := 0; i < 12 && !w.Game.Done; i++ {
		id := fmt.Sprintf("backroomstreet%02d", i)
		var err error
		if w.Game.Facing {
			err = w.apply(Command{Kind: "call", RequestID: id})
		} else {
			err = w.apply(Command{Kind: "bet", Amount: 0, RequestID: id})
		}
		if err != nil {
			t.Fatalf("playing the hand out was refused: %v", err)
		}
	}
	if !w.Game.Done {
		t.Fatal("the hand never came to a showdown")
	}
	if w.CardsDescription()["outcome"] == "" {
		t.Fatal("the hand finished and the interface is told nothing about it")
	}
}

// Money the player never sees again has to leave the table properly: what
// somebody folds stays in the pot, and the pot is always exactly what everybody
// put in it.
func TestWhatIsFoldedStaysInThePot(t *testing.T) {
	w, seated := backroom(t)
	cash, purses := w.Player.Cash, 0
	for _, id := range seated {
		purses += w.NPC(id).Purse
	}
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatalf("no game: %v", err)
	}
	if err := w.PlaceBet(100); err != nil {
		t.Fatalf("betting was refused: %v", err)
	}
	// A bet before the flop can win the hand where it stands: everybody else
	// throws theirs in and there is nothing left to fold.
	if !w.Game.Done {
		if err := w.FoldHand(); err != nil {
			t.Fatal(err)
		}
	}
	if !w.Game.Done {
		t.Fatal("the hand never finished")
	}
	after := 0
	for _, id := range seated {
		after += w.NPC(id).Purse
	}
	if d := w.Player.Cash - cash + after - purses; d != 0 {
		t.Fatalf("$%d appeared at the table out of nowhere", d)
	}
	if w.Game.Folded && w.Game.Won != -(w.Game.Ante+w.Game.MyBet) {
		t.Fatalf("a folded hand cost $%d against $%d put in", -w.Game.Won, w.Game.Ante+w.Game.MyBet)
	}
}

// The table has to talk, or the player is reading a spreadsheet. Every seat
// says what it did, and nobody's cards are visible until the hand is over.
func TestTheTableSaysWhatItDid(t *testing.T) {
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatalf("no game: %v", err)
	}
	for _, seat := range w.CardsDescription()["seats"].([]map[string]any) {
		if _, shown := seat["cards"]; shown {
			t.Fatalf("%v's hand is on the screen before it is turned over", seat["name"])
		}
	}
	playOut(t, w)
	said := 0
	for _, s := range w.Game.Seats {
		if s.Said != "" {
			said++
		}
	}
	if said == 0 {
		t.Fatal("the money went round the table and nobody said anything")
	}
	for _, seat := range w.CardsDescription()["seats"].([]map[string]any) {
		if _, shown := seat["cards"]; !shown {
			t.Fatalf("the hand is over and %v's cards are still face down", seat["name"])
		}
	}
}

// Folding has to lose. A hand thrown in wins nothing however good it was, or
// the button is decoration: with the fold not recorded, a player could throw in
// the best hand at the table and still be paid for it.
func TestAHandThrownInWinsNothingHoweverGoodItWas(t *testing.T) {
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatalf("no game: %v", err)
	}
	g := w.Game
	g.Board = hand("Ah", "Ad", "Ac", "Kh", "7d")
	g.Mine = hand("As", "Ks")
	g.Seats[0].Cards = hand("2h", "7s")
	g.Seats[1].Cards = hand("2s", "8h")
	g.Seats[2].Cards = hand("3s", "6h")
	g.Street = River
	before := w.Player.Cash
	if err := w.FoldHand(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != before {
		t.Fatalf("four aces were thrown in and paid $%d", w.Player.Cash-before)
	}
	if !strings.Contains(g.Outcome, "threw in") {
		t.Fatalf("a folded hand was described as %q", g.Outcome)
	}
}

// The view invents nothing, which only holds if the core publishes everything
// the view reads. These are the names src/Tables.tsx takes off the table: a
// rename here that is not a rename there is a blank screen at a table with
// money on it, and nothing else in the build would notice.
func TestTheTablePublishesEverythingTheScreenReads(t *testing.T) {
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatalf("no game: %v", err)
	}
	table := w.CardsDescription()
	for _, key := range []string{"place", "ante", "pot", "mine", "hand", "seats",
		"board", "street", "street_name", "bet", "my_bet", "facing", "folded",
		"done", "outcome", "won"} {
		if _, ok := table[key]; !ok {
			t.Errorf("the screen reads %q off the table and the core does not send it", key)
		}
	}
	seats, ok := table["seats"].([]map[string]any)
	if !ok || len(seats) == 0 {
		t.Fatal("the table has no seats on it")
	}
	for _, key := range []string{"who", "name", "threw", "in", "folded", "said", "sore"} {
		if _, ok := seats[0][key]; !ok {
			t.Errorf("the screen reads %q off a seat and the core does not send it", key)
		}
	}
	// And nothing at all when nobody is playing, because the screen tests for
	// the table's presence to decide whether to draw one.
	w.Game = nil
	if w.CardsDescription() != nil {
		t.Error("an empty back room published a table")
	}
}

// A game needs people in the room. Measured before the poolhall was put on the
// list of places the city goes of an evening: the most anybody ever stood in it
// in a week was one person, at nine in the morning, so the back room was a
// feature nobody could ever have used.
func TestThereIsSomebodyInTheBackRoomToPlayAgainst(t *testing.T) {
	w := New(404)
	night, day, most := 0, 0, 0
	for step := 0; step < 24*7; step++ {
		w.Advance(60)
		n := w.InTheRoom(BackRoom)
		if n > most {
			most = n
		}
		if n < 2 {
			continue
		}
		if Evening(w.Minute) {
			night++
		} else {
			day++
		}
	}
	t.Logf("in a week: at most %d in the back room, a game on %d evening hours and %d daytime hours", most, night, day)
	if night == 0 {
		t.Fatal("there was never an evening hour in a week when two people were in the back room")
	}
	// A back room is an evening. Some daytime hours are honest — a poolhall
	// with two people in it at noon is a poolhall — but if the day is as busy
	// as the night then the city has stopped going to work.
	if day >= night {
		t.Fatalf("a game was on for %d daytime hours against %d evening ones", day, night)
	}
}

// The whole reason for a game with no house in it is that the money belongs to
// somebody. A man who loses a night's money to you across a table has a reason
// to remember you, and a man you paid has a reason to like you. Until this,
// both walked away with nothing on their mind.
func TestTakingSomebodysMoneyAtCardsIsSomethingTheyRemember(t *testing.T) {
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 300); err != nil {
		t.Fatalf("no game: %v", err)
	}
	g := w.Game
	loser := w.NPC(g.Seats[0].Who)
	if loser.Sore != 0 {
		t.Fatalf("%s was already sore before a card was turned over", loser.Name)
	}
	g.Mine = hand("Ah", "As", "Ad", "Ac", "Kh")
	g.Seats[0].Cards = hand("2h", "7s", "9c", "Jc", "4h")
	g.Seats[1].Cards = hand("2s", "7h", "8c", "Jh", "4s")
	g.Seats[2].Cards = hand("3s", "6h", "8d", "Qh", "5s")
	g.Street = River
	if err := w.showdown(); err != nil {
		t.Fatal(err)
	}
	if loser.Sore == 0 {
		t.Fatalf("%s lost a third of everything they had and holds nothing against anybody", loser.Name)
	}
	if loser.SoreAt == "" {
		t.Fatal("somebody is sore about nothing in particular")
	}
}

func TestLosingToSomebodyAtCardsIsAlsoSomethingTheyRemember(t *testing.T) {
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 300); err != nil {
		t.Fatalf("no game: %v", err)
	}
	g := w.Game
	winner := w.NPC(g.Seats[0].Who)
	trust := winner.Trust
	g.Mine = hand("2h", "7s", "9c", "Jc", "4h")
	g.Seats[0].Cards = hand("Ah", "As", "Ad", "Ac", "Kh")
	g.Seats[1].Cards = hand("2s", "7h", "8c", "Jh", "4s")
	g.Seats[2].Cards = hand("3s", "6h", "8d", "Qh", "5s")
	g.Street = River
	if err := w.showdown(); err != nil {
		t.Fatal(err)
	}
	if winner.Trust <= trust {
		t.Fatalf("%s took a night's money off you and thinks no better of you for it", winner.Name)
	}
	if winner.Sore != 0 {
		t.Fatalf("%s won and is sore about it", winner.Name)
	}
}

// A small pot is a small thing. Somebody who drops a tenth of their pocket is
// not carrying it around a week later, or the whole city ends up sore at a
// player who plays cards.
func TestASmallLossIsNotHeldAgainstAnybody(t *testing.T) {
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 10); err != nil {
		t.Fatalf("no game: %v", err)
	}
	g := w.Game
	g.Mine = hand("Ah", "As", "Ad", "Ac", "Kh")
	g.Seats[0].Cards = hand("2h", "7s", "9c", "Jc", "4h")
	g.Seats[1].Cards = hand("2s", "7h", "8c", "Jh", "4s")
	g.Seats[2].Cards = hand("3s", "6h", "8d", "Qh", "5s")
	g.Street = River
	if err := w.showdown(); err != nil {
		t.Fatal(err)
	}
	for _, s := range g.Seats {
		if n := w.NPC(s.Who); n.Sore != 0 {
			t.Fatalf("%s lost $10 of $900 and holds it against you", n.Name)
		}
	}
}

// The money is real, so a room can be emptied. A player who keeps winning runs
// out of people willing and able to sit down, which is the natural end of a
// game with no house behind it: the house never runs out, and these people do.
func TestARoomYouHaveCleanedOutHasNoGameLeftInIt(t *testing.T) {
	w, seated := backroom(t)
	hands, sore := 0, 0
	for hands < 40 {
		if reason := w.BackRoomReadiness(BackRoom, 300); reason != "" {
			break
		}
		if err := w.SitInTheBackRoom(BackRoom, 300); err != nil {
			t.Fatalf("no game: %v", err)
		}
		g := w.Game
		// The player wins every hand, which is the fastest honest way to the
		// end of the room's money.
		g.Mine = hand("Ah", "As", "Ad", "Ac", "Kh")
		g.Seats[0].Cards = hand("2h", "7s", "9c", "Jc", "4h")
		g.Seats[1].Cards = hand("2s", "7h", "8c", "Jh", "4s")
		g.Seats[2].Cards = hand("3s", "6h", "8d", "Qh", "5s")
		g.Street = River
		if err := w.showdown(); err != nil {
			t.Fatal(err)
		}
		hands++
	}
	for _, id := range seated {
		if w.NPC(id).Sore > 0 {
			sore++
		}
	}
	left := w.BackRoomReadiness(BackRoom, 300)
	t.Logf("after %d winning hands: %d of the three are sore, and the room says %q", hands, sore, left)
	if hands >= 40 {
		t.Fatal("forty winning hands at $300 and the room still had money in it")
	}
	if left == "" {
		t.Fatal("the game stopped and the room says there is nothing wrong")
	}
	if sore == 0 {
		t.Fatal("a room was emptied and nobody in it minded")
	}
}

// The memory is worth nothing if it stays in a field. A man you took a night's
// money off does not sit down against you the same way twice: he pays to see
// your hand rather than believing you, because what he wants back is what you
// took. That has to be visible in the money, not in a sentence.
func TestSomebodyYouTookMoneyOffPlaysYouHarder(t *testing.T) {
	round := func(sore int) (int, int) {
		called, raised := 0, 0
		for seed := uint32(1); seed <= 1200; seed++ {
			w, seated := backroom(t)
			w.RNG = seed * 2654435761
			for _, id := range seated {
				w.NPC(id).Sore = sore
			}
			if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
				t.Fatalf("no game: %v", err)
			}
			// The player bets into them every time, so what comes back is the
			// only thing that differs between the two runs.
			if err := w.PlaceBet(150); err != nil {
				t.Fatal(err)
			}
			for _, s := range w.Game.Seats {
				if !s.Folded && s.In > 0 {
					called++
				}
			}
			if w.Game.Raised {
				raised++
			}
			if w.Game.Facing {
				if err := w.CallBet(); err != nil {
					t.Fatal(err)
				}
			}
		}
		return called, raised
	}
	freshCalls, freshRaises := round(0)
	soreCalls, soreRaises := round(SoreAtCards)
	t.Logf("against a $150 bet: a fresh table calls %d times and puts it up %d, a table you have taken money off calls %d and puts it up %d",
		freshCalls, freshRaises, soreCalls, soreRaises)
	if soreCalls <= freshCalls {
		t.Fatalf("people you cleaned out play you exactly as they did the first night: %d calls against %d", soreCalls, freshCalls)
	}
	if soreRaises <= freshRaises {
		t.Fatalf("nobody who is sore at you ever puts it up: %d against %d", soreRaises, freshRaises)
	}
}

// And what that is worth to the player, which is the point of it being in the
// game rather than in a paragraph: a table with a reason to want your money
// costs you.
//
// Under draw poker the same crude policy came out ahead against a grudge —
// $23,100 against $13,400 — because a table that calls light pays off a made
// hand. Under hold'em it comes out badly behind, because four streets against
// people who raise on less punish a policy that folds everything under three of
// a kind. Measured: $-54,800 against a grudge and $-7,800 against a fresh
// table, with less money crossing the felt rather than more, because the hands
// end earlier. The first version of this asked about the money moving and had
// to be corrected — that was true of the draw and is not true here.
func TestATableWithAGrudgeCostsYou(t *testing.T) {
	run := func(sore int) (int, int) {
		total, swing := 0, 0
		for seed := uint32(1); seed <= 2000; seed++ {
			w, seated := backroom(t)
			w.RNG = seed * 2654435761
			for _, id := range seated {
				w.NPC(id).Sore = sore
			}
			cash := w.Player.Cash
			if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
				t.Fatalf("no game: %v", err)
			}
			for i := 0; i < 12 && !w.Game.Done; i++ {
				made := BestOfSeven(w.Game.Mine, w.Game.Board).Category >= Trips
				if w.Game.Facing {
					var err error
					if made {
						err = w.CallBet()
					} else {
						err = w.FoldHand()
					}
					if err != nil {
						t.Fatal(err)
					}
					continue
				}
				put := 0
				if made {
					put = 100
				}
				if err := w.PlaceBet(put); err != nil {
					t.Fatal(err)
				}
			}
			moved := w.Player.Cash - cash
			total += moved
			if moved < 0 {
				swing -= moved
			} else {
				swing += moved
			}
		}
		return total, swing
	}
	fresh, freshSwing := run(0)
	sore, soreSwing := run(SoreAtCards)
	t.Logf("over 2000 hands: a fresh table leaves the player $%d with $%d changing hands, a table with a grudge leaves them $%d with $%d changing hands",
		fresh, freshSwing, sore, soreSwing)
	if sore >= fresh {
		t.Fatalf("a table that wants your money back cost you nothing: $%d against $%d", sore, fresh)
	}
	if freshSwing == 0 || soreSwing == 0 {
		t.Fatal("no money crossed the felt either way, so this proves nothing")
	}
}

// A player sitting down against people they have cleaned out before should be
// told so before the money goes in, not after.
func TestYouAreToldWhoAtTheTableRemembersYou(t *testing.T) {
	w, seated := backroom(t)
	w.NPC(seated[1]).Sore = SoreAtCards
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatalf("no game: %v", err)
	}
	last := w.History[len(w.History)-1]
	if !strings.Contains(last.Text, w.NPC(seated[1]).Name) ||
		!strings.Contains(last.Text, "not forgotten") {
		t.Fatalf("sitting down against somebody who is sore at you read: %q", last.Text)
	}
	// And a fresh table says nothing of the kind.
	fresh, _ := backroom(t)
	if err := fresh.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatal(err)
	}
	if body := fresh.History[len(fresh.History)-1].Text; strings.Contains(body, "not forgotten") {
		t.Fatalf("a table nobody had played before remembered something: %q", body)
	}
}
