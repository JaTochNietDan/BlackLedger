package core

import "testing"

// The boat at Pier 14.
//
// The waterfront is cheap for anybody standing on it — that is a fact about the
// floor. Holding the wharf does not buy a better price, it buys knowing when a
// boat is in: a quantity ashore tonight from somebody who has to be at sea
// before it is light, at a price no floor in this city offers, gone in the
// morning. What is left is the decision the whole underground trade is built
// on, which is how much of it you dare carry.

const thePier = "docks"

func wharfing(t *testing.T, hold bool) *World { return wharfingIn(t, hold, 53) }

func wharfingIn(t *testing.T, hold bool, seed uint32) *World {
	t.Helper()
	w := New(seed)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 90
	if hold {
		own(w, thePier)
	}
	w.Player.Location = thePier
	w.Event = nil
	return w
}

// waitForABoat runs nights until one comes in, and says how many it took.
func waitForABoat(w *World, nights int) int {
	for i := 1; i <= nights; i++ {
		w.Advance(1440)
		w.Event = nil
		if w.BoatIsIn() {
			return i
		}
	}
	return 0
}

func TestABoatComesInAtAPierOfYoursAndNowhereElse(t *testing.T) {
	t.Parallel()
	theirs := wharfing(t, false)
	if waitForABoat(theirs, 60) != 0 {
		t.Fatal("a boat unloaded for somebody who does not hold the wharf")
	}
	mine := wharfing(t, true)
	after := waitForABoat(mine, 60)
	if after == 0 {
		t.Fatal("sixty nights on your own pier and no boat came in")
	}
	t.Logf("a boat on night %d: %d %s at $%d", after, mine.Landed.Units, mine.Landed.Good, mine.Landed.Price)
	if mine.Landed.Where != thePier {
		t.Fatalf("the boat came in at %s", mine.Landed.Where)
	}
	if mine.Landed.Units <= 0 {
		t.Fatal("a boat came in empty")
	}
}

func TestTheBoatIsCheaperThanTheFloorItLandsOn(t *testing.T) {
	t.Parallel()
	w := wharfing(t, true)
	if waitForABoat(w, 60) == 0 {
		t.Skip("no boat in this city")
	}
	floor := w.PriceAt(thePier, w.Landed.Good)
	exchange := w.PriceAt("market", w.Landed.Good)
	if w.Landed.Price >= floor {
		t.Fatalf("ashore at $%d where the floor is $%d, so the pier buys nothing",
			w.Landed.Price, floor)
	}
	if exchange > 0 && floor >= exchange {
		t.Fatalf("the waterfront at $%d is dearer than the exchange at $%d", floor, exchange)
	}
}

func TestTheBoatIsGoneInTheMorning(t *testing.T) {
	t.Parallel()
	w := wharfing(t, true)
	if waitForABoat(w, 60) == 0 {
		t.Skip("no boat in this city")
	}
	if !w.BoatIsIn() {
		t.Fatal("the boat is not in")
	}
	w.Advance(LandingLasts + 60)
	w.Event = nil
	if w.BoatIsIn() {
		t.Fatal("they were still tied up a day later")
	}
	if w.LandingReadiness(1) == "" {
		t.Fatal("you can still unload a boat that has sailed")
	}
}

func TestNobodyUnloadsAtAPierThatIsNotWorking(t *testing.T) {
	t.Parallel()
	trade, _ := TradeOf(thePier)
	for _, wrong := range []string{"trouble", "short-handed"} {
		w := wharfing(t, true)
		prop := w.Properties[thePier]
		for night := 0; night < 60; night++ {
			if wrong == "trouble" {
				prop.Trouble = true
			} else {
				prop.Staff = trade.Hands - 1
			}
			w.Advance(1440)
			w.Event = nil
			if w.BoatIsIn() {
				t.Fatalf("a boat unloaded at a pier that is %s", wrong)
			}
		}
	}
}

func TestWhatComesOffTheBoatIsYoursAndIsOnYou(t *testing.T) {
	t.Parallel()
	w := wharfing(t, true)
	if waitForABoat(w, 60) == 0 {
		t.Skip("no boat in this city")
	}
	good, units, price := w.Landed.Good, w.Landed.Units, w.Landed.Price
	held, cash := w.Holding(good), w.Player.Cash
	// Part of it, because naming a number is the decision this is for.
	take := units / 2
	if take < 1 {
		take = 1
	}
	if err := w.TakeTheLanding(take); err != nil {
		t.Fatal(err)
	}
	if w.Holding(good) != held+take {
		t.Fatalf("took %d and are holding %d more", take, w.Holding(good)-held)
	}
	if w.Player.Cash != cash-price*take {
		t.Fatalf("$%d for %d at $%d each", cash-w.Player.Cash, take, price)
	}
	if w.Landed.Units != units-take {
		t.Fatalf("%d left on the boat where %d should be", w.Landed.Units, units-take)
	}
	// And it is on you. Not necessarily where anybody can find it — a bonded
	// shed is the best hiding place in the city and the pier is yours, which is
	// why the first version of this line asked for exposure and failed against
	// correct behaviour. What it costs you is that you are carrying it.
	if w.Carrying() < take {
		t.Fatalf("a lot came ashore and you are carrying %d", w.Carrying())
	}
	// The rest of it, by asking for more than is there.
	if err := w.TakeTheLanding(units * 10); err != nil {
		t.Fatal(err)
	}
	if w.Landed.Units != 0 {
		t.Fatalf("%d still on the boat after taking the lot", w.Landed.Units)
	}
	if w.BoatIsIn() {
		t.Fatal("an empty boat is still in")
	}
}

func TestTheSameCityLandsTheSameBoat(t *testing.T) {
	t.Parallel()
	one, two := wharfing(t, true), wharfing(t, true)
	a, b := waitForABoat(one, 60), waitForABoat(two, 60)
	if a == 0 {
		t.Skip("no boat in this city")
	}
	if a != b {
		t.Fatalf("the same city saw a boat on night %d and night %d", a, b)
	}
	if one.Landed != two.Landed {
		t.Fatalf("the same city landed %+v and %+v", one.Landed, two.Landed)
	}
	// A different city, which must not land the identical boat on the identical
	// night. The first version of this compared two worlds of the same seed and
	// skipped itself.
	other := wharfingIn(t, true, 83)
	if waitForABoat(other, 60) == a && other.Landed == one.Landed {
		t.Fatalf("two different cities both landed %+v on night %d", one.Landed, a)
	}
}
