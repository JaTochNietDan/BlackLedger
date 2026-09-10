package core

import "testing"

func player(t *testing.T) *World {
	t.Helper()
	w := New(211)
	w.MigrateLivingWorld()
	w.Player.Location, w.Player.Cash, w.Player.Respect = "club", 20000, HighTableStanding
	return w
}

func TestAHandIsOnTheTableUntilItIsSettled(t *testing.T) {
	t.Parallel()
	w := player(t)
	if w.HandDescription()["playing"] != false {
		t.Fatal("a hand existed before anybody sat down")
	}
	if err := w.DrawCard(); err == nil {
		t.Fatal("drew a card with no hand on the table")
	}
	if err := w.Stand(); err == nil {
		t.Fatal("stood on no hand")
	}
	cash := w.Player.Cash
	if err := w.Deal("club", 50); err != nil {
		t.Fatal(err)
	}
	stake, _ := tableStake("small")
	if w.Player.Cash != cash-stake.Amount {
		t.Fatalf("the stake was $%d", cash-w.Player.Cash)
	}
	d := w.HandDescription()
	if d["playing"] != true || d["player"].(int) < 4 || d["player"].(int) > 21 || d["dealer"].(int) < 2 {
		t.Fatalf("the table showed %v", d)
	}
	if err := w.Deal("club", 50); err == nil {
		t.Fatal("dealt a second hand over the first")
	}
	if err := w.Stand(); err != nil {
		t.Fatal(err)
	}
	if w.HandDescription()["playing"] != false {
		t.Fatal("a settled hand was still on the table")
	}
	if err := w.Stand(); err == nil {
		t.Fatal("stood on a hand that was already settled")
	}
}

func TestGoingOverEndsItImmediately(t *testing.T) {
	t.Parallel()
	w := player(t)
	found := false
	for seed := uint32(1); seed <= 600 && !found; seed++ {
		probe := player(t)
		probe.RNG = seed * 2654435761
		if err := probe.Deal("club", 50); err != nil {
			t.Fatal(err)
		}
		for probe.Hand != nil && !probe.Hand.Done {
			if err := probe.DrawCard(); err != nil {
				t.Fatal(err)
			}
		}
		if probe.Hand.Player > Bust {
			found = true
			if !probe.hasRecord("Gone at The Monarch") {
				t.Fatal("busting was not reported")
			}
		}
	}
	if !found {
		t.Fatal("drawing until the end never went over in 600 hands")
	}
	_ = w
}

// This test used to call `soften`, which brought a single ace down one at a
// time as the hand was built. The hand now keeps the cards it was dealt and
// `Total` adds them up, so the rule lives there instead: the same rule, asked
// of the thing that now owns it, and asked harder — soften could only ever
// rescue one ace, and a hand can hold four.
func TestAnAceComesDownRatherThanKillingYou(t *testing.T) {
	t.Parallel()
	ace := Card{Rank: "A", Suit: "spades", Value: 11}
	seven := Card{Rank: "7", Suit: "hearts", Value: 7}
	king := Card{Rank: "K", Suit: "clubs", Value: 10}
	four := Card{Rank: "4", Suit: "diamonds", Value: 4}

	if got := Total([]Card{ace, seven}); got != 18 {
		t.Errorf("an ace and a seven came to %d, and the ace should have stayed up", got)
	}
	if got := Total([]Card{ace, king, four}); got != 15 {
		t.Errorf("an ace, a king and a four came to %d rather than fifteen", got)
	}
	if got := Total([]Card{seven, king, four}); got != 21 {
		t.Errorf("a hand with no ace was softened to %d", got)
	}
	// Four aces is four, not forty-four: every one of them comes down, which
	// the old one-at-a-time rescue could not do.
	if got := Total([]Card{ace, ace, ace, ace}); got != 14 {
		t.Errorf("four aces came to %d", got)
	}
}

// The number on the screen has to be the sum of the cards beside it. Keeping a
// total and a list of cards is keeping the same fact twice, and this is the
// check that they cannot drift.
func TestAHandIsWorthExactlyWhatIsLyingOnTheTable(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 20000, 200, 100
	w.Player.Dress, w.Player.Location = 1, "casino"
	if err := w.Deal("casino", 50); err != nil {
		t.Fatal(err)
	}
	for step := 0; step < 8 && w.Hand != nil && !w.Hand.Done; step++ {
		if w.Hand.Player != Total(w.Hand.Mine) {
			t.Fatalf("the hand says %d and the cards on the table come to %d", w.Hand.Player, Total(w.Hand.Mine))
		}
		if w.Hand.Cards != len(w.Hand.Mine) {
			t.Fatalf("the hand counts %d cards and %d are on the table", w.Hand.Cards, len(w.Hand.Mine))
		}
		if w.Hand.Dealer != Total(w.Hand.Theirs) {
			t.Fatalf("the dealer shows %d and their cards come to %d", w.Hand.Dealer, Total(w.Hand.Theirs))
		}
		if err := w.DrawCard(); err != nil {
			t.Fatal(err)
		}
	}
}

// And every card is a real card: something a deck actually contains.
func TestEveryCardDealtIsACardThatExists(t *testing.T) {
	t.Parallel()
	w := New(11)
	seenRanks, seenSuits := map[string]bool{}, map[string]bool{}
	for i := 0; i < 4000; i++ {
		c := w.card()
		known := false
		for _, r := range ranks {
			if r.Rank == c.Rank && r.Value == c.Value {
				known = true
			}
		}
		if !known {
			t.Fatalf("the deck produced a %s worth %d", c.Rank, c.Value)
		}
		right := false
		for _, s := range suits {
			if s == c.Suit {
				right = true
			}
		}
		if !right {
			t.Fatalf("the deck produced a card of %q", c.Suit)
		}
		seenRanks[c.Rank] = true
		seenSuits[c.Suit] = true
	}
	if len(seenRanks) != len(ranks) {
		t.Errorf("4000 cards showed %d of %d ranks", len(seenRanks), len(ranks))
	}
	if len(seenSuits) != len(suits) {
		t.Errorf("4000 cards showed %d of %d suits", len(seenSuits), len(suits))
	}
}

// The edge has to land near the six percent the rest of the economy is built
// on, or the game the player sits down to and the game an owner's books model
// are different games.
func TestTheEdgeAtTheTableMatchesTheEdgeInTheBooks(t *testing.T) {
	t.Parallel()
	// A player standing on seventeen or better, which is what most people do.
	const hands = 20000
	staked, returned := 0, 0
	w := player(t)
	w.Player.Cash = 1 << 30
	stake, _ := tableStake("small")
	for i := 0; i < hands; i++ {
		before := w.Player.Cash
		if err := w.Deal("club", 50); err != nil {
			t.Fatal(err)
		}
		for w.Hand != nil && !w.Hand.Done && w.Hand.Player < DealerStands {
			if err := w.DrawCard(); err != nil {
				t.Fatal(err)
			}
		}
		if w.Hand != nil && !w.Hand.Done {
			if err := w.Stand(); err != nil {
				t.Fatal(err)
			}
		}
		staked += stake.Amount
		returned += w.Player.Cash - before + stake.Amount
	}
	edge := float64(staked-returned) / float64(staked) * 100
	if edge < HouseEdge-3 || edge > HouseEdge+3 {
		t.Fatalf("the house keeps %.1f%% at the table against %d%% in the books", edge, HouseEdge)
	}
	t.Logf("over %d hands played to seventeen, the house kept %.1f%% of everything staked", hands, edge)
}

func TestTheDealerDrawsToSixteenAndStandsOnSeventeen(t *testing.T) {
	t.Parallel()
	for seed := uint32(1); seed <= 400; seed++ {
		w := player(t)
		w.RNG = seed * 2654435761
		if err := w.Deal("club", 50); err != nil {
			t.Fatal(err)
		}
		if err := w.Stand(); err != nil {
			t.Fatal(err)
		}
		if w.Hand.Dealer < DealerStands && w.Hand.Dealer <= Bust {
			t.Fatalf("the dealer stopped on %d", w.Hand.Dealer)
		}
	}
}

func TestYouCannotSitAtYourOwnTablesOrOnesYouCannotAfford(t *testing.T) {
	t.Parallel()
	w := player(t)
	w.Player.Cash = 10
	if err := w.Deal("club", 50); err == nil {
		t.Fatal("sat down with ten dollars")
	}
	w.Player.Cash = 20000
	w.Properties["club"].Owner = "player:1"
	w.Life = 1
	if err := w.Deal("club", 50); err == nil {
		t.Fatal("played at their own house")
	}
}

func TestATieGivesTheMoneyBack(t *testing.T) {
	t.Parallel()
	pushed := false
	for seed := uint32(1); seed <= 2000 && !pushed; seed++ {
		w := player(t)
		w.RNG = seed * 2654435761
		cash := w.Player.Cash
		if err := w.Deal("club", 50); err != nil {
			t.Fatal(err)
		}
		for w.Hand != nil && !w.Hand.Done && w.Hand.Player < DealerStands {
			w.DrawCard()
		}
		if w.Hand != nil && !w.Hand.Done {
			w.Stand()
		}
		if w.Hand.Player <= Bust && w.Hand.Player == w.Hand.Dealer {
			pushed = true
			if w.Player.Cash != cash {
				t.Fatalf("a stand-off cost $%d", cash-w.Player.Cash)
			}
			if !w.hasRecord("A stand-off at The Monarch") {
				t.Fatal("a stand-off was not reported")
			}
		}
	}
	if !pushed {
		t.Fatal("no hand in 2000 ever tied")
	}
}

// The rules bug the cards uncovered. Softening used to look only at the card
// just drawn, so an ace already in the hand could not come down later: an ace,
// a five and a ten is sixteen at any table in the world, and this game called
// it twenty-six and took the money. Keeping the cards makes the hand's value a
// question about all of them, which is the only way the answer is right.
func TestAnAceAlreadyInTheHandStillComesDownLater(t *testing.T) {
	t.Parallel()
	ace := Card{Rank: "A", Suit: "spades", Value: 11}
	five := Card{Rank: "5", Suit: "hearts", Value: 5}
	ten := Card{Rank: "10", Suit: "clubs", Value: 10}
	if got := Total([]Card{ace, five, ten}); got != 16 {
		t.Errorf("an ace, a five and a ten came to %d rather than sixteen", got)
	}
	if got := Total([]Card{ace, five}); got != 16 {
		t.Errorf("an ace and a five came to %d rather than sixteen", got)
	}
	// And a hand that is genuinely over stays over.
	if got := Total([]Card{ten, ten, five}); got != 25 {
		t.Errorf("two tens and a five came to %d", got)
	}
}

// A hand that is over is the moment the player sat down for: the dealer turns
// their card over and you find out. The table used to stop describing a hand
// the instant it settled, so the felt vanished before it had said anything and
// the only account of it was a line in the ledger.
func TestASettledHandStaysOnTheTableAndSaysWhatHappened(t *testing.T) {
	t.Parallel()
	w := player(t)
	if err := w.Deal("club", 50); err != nil {
		t.Fatal(err)
	}
	// Stand without going over: a player who busts is finished before the
	// dealer plays at all, so there is no second card to turn over and the
	// table is right to show none.
	for w.Hand != nil && !w.Hand.Done && w.Hand.Player < 12 {
		if err := w.DrawCard(); err != nil {
			t.Fatal(err)
		}
	}
	if w.Hand == nil || w.Hand.Done {
		t.Fatal("the hand ended before anybody stood on it")
	}
	if err := w.Stand(); err != nil {
		t.Fatal(err)
	}
	d := w.HandDescription()
	if d["playing"] != false {
		t.Fatal("a finished hand is still being played")
	}
	if d["settled"] != true {
		t.Fatal("a finished hand is not on the table at all")
	}
	if d["outcome"] == "" || d["outcome"] == nil {
		t.Error("the hand is over and the table does not say what happened")
	}
	if len(d["theirs"].([]Card)) < 2 {
		t.Errorf("the dealer never turned their card over: %v", d["theirs"])
	}
	if d["where"] != "club" {
		t.Errorf("the table cannot say which room it is in: %v", d["where"])
	}
	// And it is gone the moment the next one is dealt over it.
	if err := w.Deal("club", 50); err != nil {
		t.Fatal(err)
	}
	if next := w.HandDescription(); next["playing"] != true || next["settled"] != false {
		t.Errorf("the old hand is still on the table: %v", next)
	}
}
