package core

import "testing"

// Two cards each and five on the table. What a player has is the best five of
// the seven they can see, which is a different question from what five cards
// are worth — and the one the whole game is played on.

func TestTheBestFiveOfSevenIsWhatYouHave(t *testing.T) {
	t.Parallel()
	cases := []struct {
		what  string
		cat   int
		hole  []Card
		board []Card
	}{
		{"a flush made with one card from your hand", Flush,
			hand("Ah", "2s"), hand("9h", "8h", "7h", "6h", "Kc")},
		{"the board's own straight, which everybody has", Straight,
			hand("2s", "3c"), hand("9h", "8s", "7d", "6c", "5h")},
		{"a full house made across the two", FullHouse,
			hand("9c", "5s"), hand("9h", "9d", "5h", "2c", "Kd")},
		{"four of a kind on the board", Quads,
			hand("2s", "3c"), hand("9h", "9d", "9c", "9s", "Kd")},
		{"two pair, where three pair is not a hand", TwoPair,
			hand("Ah", "As"), hand("Kh", "Kd", "2c", "2s", "7d")},
		{"nothing at all", HighCard,
			hand("Ah", "3s"), hand("Kh", "9d", "7c", "5s", "2d")},
	}
	for _, c := range cases {
		if got := BestOfSeven(c.hole, c.board); got.Category != c.cat {
			t.Errorf("%s came out as %s", c.what, got.Name())
		}
	}
	// Two pair where the third pair is the one that counts: aces and kings play
	// the aces and the kings, not the twos.
	made := BestOfSeven(hand("Ah", "As"), hand("Kh", "Kd", "2c", "2s", "7d"))
	if len(made.Order) == 0 || made.Order[0] != 14 || made.Order[1] != 13 {
		t.Errorf("aces and kings and twos played %v", made.Order)
	}
}

func TestWhatYouHoldStillDecidesIt(t *testing.T) {
	t.Parallel()
	board := hand("9h", "9d", "5h", "2c", "Kd")
	mine := BestOfSeven(hand("9c", "5s"), board)
	theirs := BestOfSeven(hand("Ah", "Qs"), board)
	if !mine.Beats(theirs) {
		t.Fatalf("nines full of fives lost to ace high on the same board: %s against %s", mine.Name(), theirs.Name())
	}
	// And the same two cards on a board that gives everybody the same hand is a
	// split, which is a thing hold'em does and draw poker almost never did.
	flat := hand("Ah", "Ad", "Ac", "As", "Kd")
	if BestOfSeven(hand("2s", "3c"), flat).Beats(BestOfSeven(hand("4d", "5c"), flat)) {
		t.Fatal("four aces and a king on the table was not the same hand for both of them")
	}
}

func TestTheStreetsComeInOrder(t *testing.T) {
	t.Parallel()
	street, cards := Preflop, 0
	got := []string{street}
	for street != Shown {
		street, cards = nextStreet(street)
		got = append(got, street)
		switch street {
		case Flop:
			if cards != 3 {
				t.Errorf("the flop put %d cards on the table", cards)
			}
		case Turn, River:
			if cards != 1 {
				t.Errorf("the %s put %d cards on the table", street, cards)
			}
		}
	}
	if len(got) != 5 {
		t.Fatalf("a hand of hold'em went %v", got)
	}
}

// A hand of hold'em is four rounds of the same decision against a board that
// grows. This walks one from the seat to the showdown and says what the table
// looked like at each step.
func TestAHandGoesPreflopFlopTurnRiver(t *testing.T) {
	t.Parallel()
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatalf("no game: %v", err)
	}
	seen := map[string]int{}
	for i := 0; i < 12 && !w.Game.Done; i++ {
		seen[w.Game.Street] = len(w.Game.Board)
		if w.Game.Facing {
			if err := w.CallBet(); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if err := w.PlaceBet(0); err != nil {
			t.Fatal(err)
		}
	}
	if !w.Game.Done {
		t.Fatal("checking every street never reached a showdown")
	}
	for street, cards := range map[string]int{Preflop: 0, Flop: 3, Turn: 4, River: 5} {
		if got, reached := seen[street]; !reached {
			t.Errorf("the hand never went through the %s", street)
		} else if got != cards {
			t.Errorf("%d cards were face up %s, not %d", got, StreetName(street), cards)
		}
	}
	if len(w.Game.Board) != 5 {
		t.Fatalf("a hand reached the showdown with %d cards on the table", len(w.Game.Board))
	}
	// Nobody sees a hole card until it is over, and then everybody does.
	for _, s := range w.Game.Seats {
		if len(s.Cards) != 2 {
			t.Fatalf("%s is holding %d cards", s.Name, len(s.Cards))
		}
	}
}

// A hand nobody is left contesting is over where it stands. There is no reason
// to deal a river to somebody playing against nobody.
func TestAHandEverybodyElseThrewInIsOverWhereItStands(t *testing.T) {
	t.Parallel()
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatalf("no game: %v", err)
	}
	g := w.Game
	before := w.Player.Cash
	for i := range g.Seats {
		g.Seats[i].Folded = true
	}
	if err := w.PlaceBet(0); err != nil {
		t.Fatal(err)
	}
	if !g.Done {
		t.Fatal("the player was left playing against nobody and the hand went on")
	}
	if w.Player.Cash != before+g.Pot {
		t.Fatalf("everybody threw their hand in and the pot of $%d came to $%d", g.Pot, w.Player.Cash-before)
	}
}

// Two cards are not a hand yet, and the screen says what you have from the
// moment they are dealt. Rank answers the question it is asked — five cards —
// so asked about two of a suit it called them a flush, and a player holding the
// queen and three of diamonds was told they had one.
func TestTwoCardsAreNotAFlush(t *testing.T) {
	t.Parallel()
	if got := BestOfSeven(hand("Qd", "3d"), nil); got.Category != HighCard {
		t.Errorf("the queen and three of diamonds is %s", got.Name())
	}
	if got := BestOfSeven(hand("Qd", "Qs"), nil); got.Category != Pair {
		t.Errorf("two queens is %s", got.Name())
	}
	if got := BestOfSeven(hand("Ah", "Kh"), hand("2c", "7d", "9s")); got.Category != HighCard {
		t.Errorf("ace king of hearts on a rainbow flop is %s", got.Name())
	}
}

// And the board reaches the screen as a row rather than as nothing at all. A
// nil slice arrives in the browser as null, the screen counts its length, and
// the whole page goes blank the moment anybody sits down.
func TestTheBoardIsAlwaysARow(t *testing.T) {
	t.Parallel()
	w, _ := backroom(t)
	if err := w.SitInTheBackRoom(BackRoom, 50); err != nil {
		t.Fatal(err)
	}
	table := w.CardsDescription()
	board, ok := table["board"].([]Card)
	if !ok {
		t.Fatalf("the board is published as %T before the flop", table["board"])
	}
	if board == nil {
		t.Fatal("the board is published as null before the flop, which is a blank screen")
	}
	if len(board) != 0 {
		t.Fatalf("%d cards are face up before anybody has bet", len(board))
	}
}
