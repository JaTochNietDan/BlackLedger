package core

import (
	"fmt"
	"strings"
	"testing"
)

// The table dealt one hand and stopped. To play a second you got up and sat
// down again, which re-seated the room, re-read everybody's pockets and threw
// away everything the last hand had meant: "the game should continue until you
// stop playing, right now it just requires you to leave the table and rejoin.
// Realistically it feels like it should be more like actual poker, where you
// have a buy in and whatnot and you play until people go bust or you can
// leave."

func attheTable(t *testing.T, at string, buyIn int) *World {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Cash = 5000
	w.Player.Location = at
	// An evening, which is when anybody is in the room.
	for w.Minute%1440 < 1200 {
		w.Event = nil
		w.Advance(60)
		w.Event = nil
	}
	if err := w.Sit(at, Backroom); err != nil {
		t.Fatal(err)
	}
	if err := w.SitInTheBackRoom(at, buyIn); err != nil {
		t.Fatal(err)
	}
	return w
}

// checkThrough checks every street through to the showdown.
func checkThrough(t *testing.T, w *World) {
	t.Helper()
	for g := w.Game; !g.Done; {
		var err error
		if g.Facing {
			err = w.CallBet()
		} else {
			err = w.PlaceBet(0)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestTheMoneyOnTheTableIsWhatYouPutThere(t *testing.T) {
	t.Parallel()
	w := attheTable(t, "bar", 1000)
	g := w.Game
	if g.BuyIn != 1000 {
		t.Fatalf("brought $1000 and the table says $%d", g.BuyIn)
	}
	// The buy-in left the pocket and is on the table.
	if w.Player.Cash != 4000 {
		t.Fatalf("bought in for $1000 out of $5000 and still carry $%d", w.Player.Cash)
	}
	if g.Stack+g.Ante != 1000 {
		t.Fatalf("$%d in front of you and $%d anted out of a $1000 buy-in", g.Stack, g.Ante)
	}
	// Everybody else is playing out of chips too, and the stake is the room's
	// rather than the player's: these people carry fifty dollars.
	for _, s := range g.Seats {
		if s.Stack <= 0 {
			t.Fatalf("%s sat down with nothing in front of them", s.Name)
		}
		// Their ante is already in the pot, so what they sat down with is what
		// is in front of them plus it.
		if s.Stack+g.Ante < g.Ante*SitsOutUnder {
			t.Fatalf("%s sat down with $%d at $%d a hand, which is not sitting down",
				s.Name, s.Stack+g.Ante, g.Ante)
		}
	}
	t.Logf("$%d on the table at $%d a hand against %d", g.BuyIn, g.Ante, len(g.Seats))

}

// And what the player can bet is what is in front of them, not what they carry.
// This is the whole point of a buy-in, and it takes a short stack to see: a
// thousand dollars on the table is more than the room takes on one bet anyway,
// so the ceiling on a bet hid the rule that matters.
func TestYouCanOnlyBetWhatIsInFrontOfYou(t *testing.T) {
	t.Parallel()
	w := attheTable(t, "bar", MinBuyIn)
	g := w.Game
	if g.Stack >= MaxAnte {
		t.Fatalf("$%d in front of you is more than the room takes on a bet, so this measures nothing",
			g.Stack)
	}
	over := g.Stack + 1
	if over > MaxAnte || over > w.Player.Cash {
		t.Fatalf("$%d is not a bet the pocket could have covered, so this measures nothing", over)
	}
	if err := w.PlaceBet(over); err == nil {
		t.Fatalf("bet $%d with $%d in front of you and $%d in your pocket",
			over, g.Stack, w.Player.Cash)
	}
	// And what is in front of you can be bet.
	if err := w.PlaceBet(g.Stack); err != nil {
		t.Fatalf("could not push in the $%d in front of you: %v", g.Stack, err)
	}
}

func TestTheGameGoesOnUntilYouStopPlaying(t *testing.T) {
	t.Parallel()
	w := attheTable(t, "bar", 1000)
	g := w.Game
	// Ten hands without ever leaving the table.
	for hand := 0; hand < 10 && !g.Over; hand++ {
		checkThrough(t, w)
		if g.Over {
			break
		}
		if err := w.DealAgain(); err != nil {
			t.Fatalf("hand %d would not go out: %v", hand+2, err)
		}
	}
	t.Logf("%d hands played, $%d in front of you, %d still at the table",
		g.Hands, g.Stack, len(g.Seats))
	if g.Hands < 2 {
		t.Fatalf("only %d hand was dealt, so the table still stops after one", g.Hands)
	}
	// Nobody got up and sat down again: it is the same sitting throughout.
	if w.Seated != "bar" || w.SeatedTo != Backroom {
		t.Fatal("the player left the table at some point")
	}
}

// Somebody who cannot cover the ante is cleaned out, takes what is left back to
// their pocket, and goes.
func TestSomebodyCleanedOutLeavesTheTable(t *testing.T) {
	t.Parallel()
	w := attheTable(t, "bar", 1000)
	g := w.Game
	checkThrough(t, w)
	// Take the shortest stack down to nothing, the way a hand would.
	short := &g.Seats[0]
	who, name, left := short.Who, short.Name, short.Stack
	short.Stack = 0
	pocket := 0
	if n := w.NPC(who); n != nil {
		pocket = n.Purse
	}
	if err := w.DealAgain(); err != nil {
		t.Fatal(err)
	}
	for _, s := range g.Seats {
		if s.Who == who {
			t.Fatalf("%s has nothing and is still in the hand", name)
		}
	}
	if n := w.NPC(who); n == nil || n.Purse != pocket {
		t.Fatalf("%s left with $%d rather than the nothing in front of them", name, left)
	}
	told := false
	for _, r := range w.History {
		told = told || strings.Contains(r.Title, "Cleaned out at the table")
	}
	if !told {
		t.Fatal("somebody was cleaned out and nobody said so")
	}
}

// And you can pick your money up. What is in front of you goes back in your
// pocket, and what is in front of everybody else goes back in theirs.
func TestPickingYourMoneyUpTakesItWithYou(t *testing.T) {
	t.Parallel()
	w := attheTable(t, "bar", 1000)
	g := w.Game
	checkThrough(t, w)
	stack, cash := g.Stack, w.Player.Cash
	theirs := map[string]int{}
	for _, s := range g.Seats {
		if n := w.NPC(s.Who); n != nil {
			theirs[s.Who] = n.Purse + s.Stack
		}
	}
	if err := w.endSitting("You pick your money up off the table."); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash+stack {
		t.Fatalf("picked up $%d off a $%d stack", w.Player.Cash-cash, stack)
	}
	if !g.Over {
		t.Fatal("the money came off the table and the game is still going")
	}
	for who, had := range theirs {
		if n := w.NPC(who); n == nil || n.Purse != had {
			t.Fatalf("somebody's chips were left on the baize")
		}
	}
	// And nothing is conjured: the table's money is the buy-ins, no more.
	if g.Stack != 0 {
		t.Fatalf("$%d is still in front of a player who got up", g.Stack)
	}
}

// The whole night, through the same path a player uses, with nothing reached
// into directly.
func TestASittingPlayedThroughTheCommands(t *testing.T) {
	t.Parallel()
	w := attheTable(t, "bar", 600)
	hands := 0
	for i := 0; i < 40 && !w.Game.Over; i++ {
		w.Event = nil
		if w.Game.Done {
			if err := w.apply(Command{Kind: "deal", Target: "bar",
				RequestID: "deal-" + fmt.Sprint(i)}); err != nil {
				break
			}
			hands++
			continue
		}
		kind := "bet"
		if w.Game.Facing {
			kind = "call"
		}
		if err := w.apply(Command{Kind: kind, Target: "bar", Amount: 0,
			RequestID: "play-" + fmt.Sprint(i)}); err != nil {
			t.Fatalf("%s at hand %d: %v", kind, w.Game.Hands, err)
		}
	}
	t.Logf("%d hands through the command path, over=%v, $%d in front of you",
		w.Game.Hands, w.Game.Over, w.Game.Stack)
	if w.Game.Hands < 2 {
		t.Fatalf("only %d hand played through the commands", w.Game.Hands)
	}
}

// Putting money on a table is sitting down at it, whichever button was pressed.
// The buy-in is offered in the room's own list as well as behind the door, and
// taking it there left the player standing in the room with a hand of cards
// going on somewhere the screen could not draw: "I set the amount I want to buy
// in but it just ran a simulation instead of letting me play the game."
func TestBuyingInSitsYouDown(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 30
	w.Player.Cash = 5000
	w.Player.Location = "bar"
	for w.Minute%1440 < 1200 {
		w.Event = nil
		w.Advance(60)
		w.Event = nil
	}
	// Straight from the room, with no seat taken first.
	if w.Seated != "" {
		t.Fatal("the player is already at a table")
	}
	if err := w.apply(Command{Kind: "cards", Target: "bar", Amount: 600,
		RequestID: "buyinfromtheroom"}); err != nil {
		t.Fatal(err)
	}
	if w.Game == nil {
		t.Fatal("no hand was dealt")
	}
	if w.Seated != "bar" {
		t.Fatalf("money is on the table at the bar and the player is seated at %q", w.Seated)
	}
	if w.SeatedTo != Backroom {
		t.Fatalf("the player bought into a card game and is sitting at the %q", w.SeatedTo)
	}
}
