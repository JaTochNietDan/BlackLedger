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

// everything is what the table is worth in total: pockets, and the chips in
// front of anybody. The money used to live only in pockets, so a hand could be
// weighed by reading them before and after. It lives on the table during a
// sitting now, and reading the pockets alone says the ante vanished.
func everything(w *World, _ []string) int {
	total := w.Player.Cash
	// Every pocket in the city, not only the three this test sat down: the
	// room fills its own seats, so a stranger's chips would otherwise look like
	// money appearing from nowhere.
	for i := range w.NPCs {
		total += w.NPCs[i].Purse
	}
	if g := w.Game; g != nil {
		total += g.Stack
		for _, s := range g.Seats {
			total += s.Stack
		}
		// The pot only while it is still a pot. A hand that has been settled
		// has already paid it into the stacks, and the figure is kept so the
		// screen can say what was played for — counting both is counting the
		// same money twice.
		if !g.Done {
			total += g.Pot
		}
	}
	return total
}

// theirs is what the room has, in pockets and in front of them.
func theirs(w *World) int {
	total := 0
	for i := range w.NPCs {
		total += w.NPCs[i].Purse
	}
	if g := w.Game; g != nil {
		for _, s := range g.Seats {
			total += s.Stack
		}
		// A pot in play was put in by both sides, so it belongs to neither
		// until it is settled; once it is, it is already in the stacks.
		if !g.Done {
			total += g.Pot - g.MyBet - g.Ante
		}
	}
	return total
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
	t.Parallel()
	w, seated := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
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
	t.Parallel()
	for seed := uint32(1); seed <= 200; seed++ {
		w, _ := backroom(t)
		w.RNG = seed * 2654435761
		if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
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
	t.Parallel()
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
		t.Fatalf("nobody could get a game: %v", err)
	}
	pot := w.Game.Pot
	if pot <= 0 {
		t.Fatal("nobody put anything in")
	}
	playOut(t, w)
	g := w.Game
	if !g.Done {
		t.Fatal("the hand never finished")
	}
	// What the hand paid, which is money on the table rather than money in a
	// pocket. Sitting down puts a stake in front of you and the pot is settled
	// against it; `Player.Cash` only moves when the money is picked up again.
	//
	// This asked whether the cash in hand had risen by the whole pot, which is
	// the wrong pocket and also double-counts the player's own ante. It passed
	// for as long as the one deal it ran on never gave the player the best
	// hand — and the moment the seeds were spread it dealt one and reported
	// that the best hand at the table had won nothing.
	best, mine := BestAtTheTable(g), BestOfSeven(g.Mine, g.Board)
	switch {
	case best.Beats(mine) && g.Won > 0:
		t.Fatalf("a losing hand took $%d off the table", g.Won)
	case mine.Beats(best) && g.Won <= 0:
		t.Fatalf("the best hand at the table came away %d from a $%d pot", g.Won, pot)
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
	t.Parallel()
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
	t.Parallel()
	const hands = 3000
	table, mine := 0, 0
	for seed := uint32(1); seed <= hands; seed++ {
		w, seated := backroom(t)
		w.RNG = seed * 2654435761
		// Counted before the ante, or the antes are missing from one side of
		// the sum and the pot looks like money out of nowhere.
		cash := w.Player.Cash
		was, theirsWas := everything(w, seated), theirs(w)
		if err := w.SitInTheBackRoom(BackRoom, 800); err != nil {
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
		// What the hand did to the player, whether it is in their pocket or in
		// front of them.
		// The buy-in has already left the pocket, so the position is what is
		// carried plus what is in front of them against what they walked in with.
		mine += w.Player.Cash + w.Game.Stack - cash
		table += theirs(w) - theirsWas
		if now := everything(w, seated); now != was {
			t.Fatalf("$%d appeared at the table out of nowhere", now-was)
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
	t.Parallel()
	// Everything here is measured in antes rather than in dollars. It used to
	// bet a flat $100 whatever the game was played for, which is two antes at
	// one stake and five at another — so the same rule read as holding at $50
	// and inverting at $40, and the alarm was the measurement rather than the
	// game. Three stakes now, because a property that is only true at one of
	// them is not the property this claims to guard.
	run := func(buyIn int, fold, bet bool) int {
		total := 0
		for seed := uint32(1); seed <= 1000; seed++ {
			w, _ := backroom(t)
			w.RNG = seed * 2654435761
			cash := w.Player.Cash
			if err := w.SitInTheBackRoom(BackRoom, buyIn); err != nil {
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
					put = min(w.Game.Ante*2, w.Game.Stack)
				}
				if err := w.PlaceBet(put); err != nil {
					t.Fatalf("betting was refused: %v", err)
				}
			}
			// What the hand did, whether it ended in the pocket or in front of
			// them. The buy-in has already left the pocket by here.
			total += w.Player.Cash + w.Game.Stack - cash
		}
		return total
	}
	for _, buyIn := range []int{400, 1000, 2000} {
		calls, folds, plays := run(buyIn, false, false), run(buyIn, true, false), run(buyIn, true, true)
		t.Logf("1000 hands on a $%d buy-in: calling everything is $%d, folding what is beaten is $%d, folding and betting the good ones is $%d",
			buyIn, calls, folds, plays)
		if folds <= calls {
			t.Fatalf("at $%d: throwing beaten hands in cost more than paying for them: %d against %d",
				buyIn, folds, calls)
		}
		if plays <= folds {
			t.Fatalf("at $%d: putting money on a made hand earned nothing: %d against %d",
				buyIn, plays, folds)
		}
	}
}

// Random hands almost never tie — none did in four thousand — so the split has
// to be put on the table by hand or it is a branch nobody has ever run.
func TestATiedPotIsSplitRatherThanGivenAway(t *testing.T) {
	t.Parallel()
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
		t.Fatalf("no game: %v", err)
	}
	g := w.Game
	// Into the chips in front of them rather than into a pocket: the money is
	// on the table until somebody picks it up.
	before := g.Stack
	theirs := g.Seats[0].Stack
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
	if got := g.Stack - before; got != half {
		t.Fatalf("half of a $%d pot came to $%d", g.Pot, got)
	}
	if got := g.Seats[0].Stack - theirs; got != half {
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
	t.Parallel()
	w, _ := backroom(t)
	offered := func(kind string) *Action {
		for _, a := range w.Actions(w.Player.Location) {
			if a.ID == kind {
				return &a
			}
		}
		return nil
	}
	// Through the door first: the table is only offered to somebody sitting at
	// it, because the room offering both a way in and a way to buy in was the
	// same door twice.
	if err := w.Sit(BackRoom, Backroom); err != nil {
		t.Fatal(err)
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
	if err := w.apply(Command{Kind: "cards", Amount: 400, RequestID: "backroomsitdown1"}); err != nil {
		t.Fatalf("sitting down was refused: %v", err)
	}
	// The figure the player types is what they put on the table, not the ante.
	// The ante follows from it and from what the room can play for.
	if w.Game == nil || w.Game.BuyIn != 400 {
		t.Fatalf("a player who typed $400 put $%d on the table", w.Game.BuyIn)
	}
	if w.Player.Cash != before-400 {
		t.Fatalf("buying in took $%d", before-w.Player.Cash)
	}
	if w.Game.Ante <= 0 || w.Game.Ante > 400/AntesInAStack {
		t.Fatalf("a $400 buy-in is playing for $%d a hand", w.Game.Ante)
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
	t.Parallel()
	w, seated := backroom(t)
	was := everything(w, seated)
	if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
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
	// Weighed with the chips in it, because most of this money is on the table
	// rather than in anybody's pocket until the sitting ends.
	if d := everything(w, seated) - was; d != 0 {
		t.Fatalf("$%d appeared at the table out of nowhere", d)
	}
	if w.Game.Folded && w.Game.Won != -(w.Game.Ante+w.Game.MyBet) {
		t.Fatalf("a folded hand cost $%d against $%d put in", -w.Game.Won, w.Game.Ante+w.Game.MyBet)
	}
}

// The table has to talk, or the player is reading a spreadsheet. Every seat
// says what it did, and nobody's cards are visible until the hand is over.
func TestTheTableSaysWhatItDid(t *testing.T) {
	t.Parallel()
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
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
	t.Parallel()
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
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
	t.Parallel()
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
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

// A game needs people in the room, and a back room is an evening.
//
// Measured before the poolhall was put on the list of places the city goes of an
// evening: the most anybody ever stood in it in a week was one person, at nine
// in the morning, so the back room was a feature nobody could ever have used.
//
// The first version of this counted HOURS in which two or more people were in
// the room, on one seed, and asked that the evening have more of them than the
// day. That measure saturates: two people is nearly always true, so both
// columns sit near the eighty-four hours a week has of each, and across twenty
// seeds twelve of them come out an exact tie. It passed on its one seed by a
// coincidence, and adding a single address to the city — three more people with
// somewhere else to be — flipped it. The property was never in doubt; the
// instrument could not see it.
//
// Heads, not hours, and twenty cities rather than one. The evening carries
// half again as many people as the day, which is the thing worth guarding.
func TestTheBackRoomIsAnEveningRatherThanAnAfternoon(t *testing.T) {
	t.Parallel()
	night, day, hours, most := 0, 0, 0, 0
	for seed := uint32(400); seed < 420; seed++ {
		w := New(seed)
		for step := 0; step < 24*7; step++ {
			w.Advance(60)
			n := w.InTheRoom(BackRoom)
			if n > most {
				most = n
			}
			if Evening(w.Minute) {
				night += n
			} else {
				day += n
			}
			hours++
		}
	}
	t.Logf("twenty weeks: %d in the room of an evening and %d in the daytime, at most %d at once",
		night, day, most)
	if hours < 3000 {
		t.Fatalf("only %d hours were watched, so this measures nothing", hours)
	}
	if most < 2 {
		t.Fatalf("the most anybody ever saw in the back room was %d, so there is no game", most)
	}
	// Half again as many. A poolhall with people in it at noon is a poolhall;
	// a city whose evenings are no busier than its afternoons has stopped
	// going to work.
	if night*100 < day*125 {
		t.Fatalf("%d in the room of an evening against %d in the daytime, which is not an evening trade",
			night, day)
	}
}

// The whole reason for a game with no house in it is that the money belongs to
// somebody. A man who loses a night's money to you across a table has a reason
// to remember you, and a man you paid has a reason to like you. Until this,
// both walked away with nothing on their mind.
func TestTakingSomebodysMoneyAtCardsIsSomethingTheyRemember(t *testing.T) {
	t.Parallel()
	w, _ := backroom(t)
	// Enough on the table that everybody puts everything they have on it, so
	// what is in front of them is what they own. The stake used to be named
	// here as an ante of $300 against a $900 pocket, which is not a bet
	// anybody makes: a night's money is lost by pushing chips in over a hand,
	// not by an ante set to a third of what somebody is carrying.
	if err := w.SitInTheBackRoom(BackRoom, 2000); err != nil {
		t.Fatalf("no game: %v", err)
	}
	g := w.Game
	loser := w.NPC(g.Seats[0].Who)
	if loser.Sore != 0 {
		t.Fatalf("%s was already sore before a card was turned over", loser.Name)
	}
	// They push most of what is in front of them in, and it is beaten.
	put := g.Seats[0].Stack * 3 / 4
	g.Seats[0].Stack -= put
	g.Seats[0].In += put
	g.Pot += put
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
	t.Parallel()
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
	t.Parallel()
	w, _ := backroom(t)
	// The smallest table there is, against pockets of $900: what anybody loses
	// here is small against what they are carrying, which is the whole point.
	// This used to name a $10 ante directly; the figure is the buy-in now and
	// the ante follows from it.
	if err := w.SitInTheBackRoom(BackRoom, MinBuyIn); err != nil {
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
	t.Parallel()
	w, seated := backroom(t)
	// One sitting, hand after hand, which is what this is now: the test used to
	// get up and sit down again between every hand, because that was the only
	// way to play a second one.
	if err := w.SitInTheBackRoom(BackRoom, 2000); err != nil {
		t.Fatalf("no game: %v", err)
	}
	g := w.Game
	hands, sore := 0, 0
	for hands < 60 && !g.Over {
		// The player wins every hand, which is the fastest honest way to the
		// end of the room's money.
		g.Mine = hand("Ah", "As", "Ad", "Ac", "Kh")
		for i := range g.Seats {
			g.Seats[i].Cards = hand("2h", "7s", "9c", "Jc", "4h")
			// And everybody puts most of what is in front of them in, or a
			// night of antes takes longer than anybody would sit for.
			put := g.Seats[i].Stack / 2
			g.Seats[i].Stack -= put
			g.Seats[i].In += put
			g.Pot += put
		}
		g.Street = River
		if err := w.showdown(); err != nil {
			t.Fatal(err)
		}
		hands++
		if err := w.DealAgain(); err != nil {
			break
		}
	}
	for _, id := range seated {
		if w.NPC(id).Sore > 0 {
			sore++
		}
	}
	t.Logf("after %d winning hands: %d of the three are sore, the table is over=%v (%q), and %d are still sitting",
		hands, sore, g.Over, g.Ended, len(g.Seats))
	if hands >= 60 {
		t.Fatal("sixty winning hands and the room still had money in it")
	}
	if !g.Over {
		t.Fatal("the room was cleaned out and the game is still going")
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
	t.Parallel()
	round := func(sore int) (int, int) {
		called, raised := 0, 0
		for seed := uint32(1); seed <= 1200; seed++ {
			w, seated := backroom(t)
			w.RNG = seed * 2654435761
			for _, id := range seated {
				w.NPC(id).Sore = sore
			}
			if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
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
	t.Parallel()
	run := func(sore int) (int, int) {
		total, swing := 0, 0
		for seed := uint32(1); seed <= 2000; seed++ {
			w, seated := backroom(t)
			w.RNG = seed * 2654435761
			for _, id := range seated {
				w.NPC(id).Sore = sore
			}
			cash := w.Player.Cash
			if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
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
			// What the hand did, counting the chips: the buy-in has left the
			// pocket and the winnings have not come back to it yet.
			moved := w.Player.Cash + w.Game.Stack - cash
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
	t.Parallel()
	w, seated := backroom(t)
	w.NPC(seated[1]).Sore = SoreAtCards
	if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
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

// The sentence a losing hand is described with. A hand names itself with the
// article in front of it — "a pair", "a flush" — and the table's own line puts
// "your" in front of that: "Leo Carver had a pair against your a pair", which is
// what a playtest read back off the felt.
func TestALosingHandIsDescribedInEnglish(t *testing.T) {
	t.Parallel()
	for _, category := range []int{HighCard, Pair, TwoPair, Trips, Straight, Flush, FullHouse, Quads, StraightFlush} {
		h := HandRank{Category: category}
		if h.Name() == "" || h.Bare() == "" {
			t.Fatalf("a hand of category %d names itself as nothing", category)
		}
		line := "Leo had " + h.Name() + " against your " + h.Bare() + "."
		if strings.Contains(line, "your a ") || strings.Contains(line, "your an ") {
			t.Fatalf("reads wrong: %q", line)
		}
		// And the article is only dropped where there was one.
		if !strings.HasPrefix(h.Name(), "a ") && h.Bare() != h.Name() {
			t.Fatalf("%q lost something other than its article: %q", h.Name(), h.Bare())
		}
	}

	// And the line the felt actually builds, off a hand the player loses.
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 1000); err != nil {
		t.Fatal(err)
	}
	g := w.Game
	// The player has to hold a hand that names itself with an article, or the
	// sentence reads the same either way and this measures nothing: "high card"
	// against "high card" is fine however it is joined.
	g.Board = hand("2h", "7s", "9c", "Jc", "4h")
	g.Mine = hand("Ah", "Ac")
	g.Seats[0].Cards = hand("2s", "2d")
	for i := range g.Seats {
		g.Seats[i].Folded = i != 0
	}
	g.Street = River
	if err := w.showdown(); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", g.Outcome)
	if !strings.Contains(g.Outcome, "pair") {
		t.Fatalf("the player was meant to lose with a hand that has an article: %q", g.Outcome)
	}
	if strings.Contains(g.Outcome, "your a ") {
		t.Fatalf("the felt says %q", g.Outcome)
	}
}
