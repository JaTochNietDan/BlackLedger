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
	if err := w.Deal("club", "small"); err != nil {
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
	if err := w.Deal("club", "small"); err == nil {
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
	w := player(t)
	found := false
	for seed := uint32(1); seed <= 600 && !found; seed++ {
		probe := player(t)
		probe.RNG = seed * 2654435761
		if err := probe.Deal("club", "small"); err != nil {
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

func TestAnAceComesDownRatherThanKillingYou(t *testing.T) {
	if soften(22, 11) != 12 {
		t.Fatalf("an ace on twenty-two left %d", soften(22, 11))
	}
	if soften(22, 7) != 22 {
		t.Fatal("a seven behaved like an ace")
	}
	if soften(18, 11) != 18 {
		t.Fatal("an ace came down when it did not have to")
	}
}

// The edge has to land near the six percent the rest of the economy is built
// on, or the game the player sits down to and the game an owner's books model
// are different games.
func TestTheEdgeAtTheTableMatchesTheEdgeInTheBooks(t *testing.T) {
	// A player standing on seventeen or better, which is what most people do.
	const hands = 20000
	staked, returned := 0, 0
	w := player(t)
	w.Player.Cash = 1 << 30
	stake, _ := tableStake("small")
	for i := 0; i < hands; i++ {
		before := w.Player.Cash
		if err := w.Deal("club", "small"); err != nil {
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
	for seed := uint32(1); seed <= 400; seed++ {
		w := player(t)
		w.RNG = seed * 2654435761
		if err := w.Deal("club", "small"); err != nil {
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
	w := player(t)
	w.Player.Cash = 10
	if err := w.Deal("club", "small"); err == nil {
		t.Fatal("sat down with ten dollars")
	}
	w.Player.Cash = 20000
	w.Properties["club"].Owner = "player:1"
	w.Life = 1
	if err := w.Deal("club", "small"); err == nil {
		t.Fatal("played at their own house")
	}
}

func TestATieGivesTheMoneyBack(t *testing.T) {
	pushed := false
	for seed := uint32(1); seed <= 2000 && !pushed; seed++ {
		w := player(t)
		w.RNG = seed * 2654435761
		cash := w.Player.Cash
		if err := w.Deal("club", "small"); err != nil {
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
