package core

import (
	"fmt"
	"testing"
)

// Now that a room's takings depend on who is standing in it, an owner needs
// something to do about an empty room. Putting a night on is the oldest answer
// there is: pay for a band and a barrel, and the people who would have drunk
// somewhere else drink here instead.
//
// It has to move actual people. A number that goes up without anybody walking
// through the door would be the same abstraction the takings used to be.

func host(t *testing.T) *World {
	t.Helper()
	w := New(47)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 5000, 100
	own(w, "club")
	w.Player.Location = "club"
	// Let the city settle into its own habits first. Nobody has a place they
	// usually drink until they have had a day of standing somewhere, and a
	// night is a reason to go somewhere other than usual.
	for day := 0; day < 3; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	w.Player.Cash, w.Player.Location = 5000, "club"
	return w
}

func TestPuttingOnANightFillsTheRoom(t *testing.T) {
	t.Parallel()
	w := host(t)
	a := actionByID(w.Actions("club"), "night")
	if a == nil {
		t.Fatal("a room of yours cannot put a night on")
	}
	if a.Disabled {
		t.Fatalf("refused: %s", a.Reason)
	}
	cash := w.Player.Cash
	if err := w.PutOnANight("club"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash >= cash {
		t.Fatal("the band played for nothing")
	}
	// Booked, not playing. A night bought in the morning is for this evening,
	// and NightOn answers about the evening rather than the purchase.
	if w.NightReadiness("club") == "" {
		t.Fatal("the night was not booked")
	}

	// Against the same evening without one. Measuring the room before and after
	// the evening comes measures the evening: the bar, the club and the casino
	// fill up every night whether anybody paid for a band or not, and a test
	// that watched one room fill passed with the draw deleted.
	quiet := host(t)
	evening := func(x *World) int {
		for i := 0; i < 40 && !Evening(x.Minute); i++ {
			x.Advance(60)
		}
		x.Advance(180) // long enough for the walk
		return x.Footfall("club")
	}
	loud := evening(w)
	ordinary := evening(quiet)
	if loud <= ordinary {
		t.Fatalf("a night drew %d where an ordinary evening drew %d", loud, ordinary)
	}
}

func TestANightIsOneNight(t *testing.T) {
	t.Parallel()
	w := host(t)
	if err := w.PutOnANight("club"); err != nil {
		t.Fatal(err)
	}
	if w.NightReadiness("club") == "" {
		t.Fatal("two nights were put on at once")
	}
	w.Advance(NightLasts + 1440)
	if w.NightOn("club") {
		t.Fatal("the band never went home")
	}
	w.Player.Cash = 5000
	if reason := w.NightReadiness("club"); reason != "" {
		t.Fatalf("could not put another one on a week later: %s", reason)
	}
}

func TestYouCannotPutANightOnSomebodyElsesRoom(t *testing.T) {
	t.Parallel()
	w := host(t)
	w.Player.Location = "bar"
	if w.NightReadiness("bar") == "" {
		t.Fatal("paid for a band in a room that is not yours")
	}
	w.Player.Location = "laundry"
	w.Properties["laundry"].Owner = w.Properties["club"].Owner
	if w.NightReadiness("laundry") == "" {
		t.Fatal("a laundry put a night on")
	}
}

// A room the city has no habit of visiting is the only room worth paying for a
// band in, and it was the one room a night did nothing for. Two faults sat on
// top of each other. The draw was measured from the desk somebody was standing
// at when the question came up rather than from the room they would have drunk
// in, so a district with a band in it drew nobody whenever the day's work
// happened to be across town; and the word was held for a flat day from the
// minute it was bought, which expired on the exact minute of the next evening's
// departure and so covered no evening at all. The old guard missed both because
// it put its night on in a room the whole city already walked to.
func quiet(t *testing.T, id string, seed uint32) *World {
	t.Helper()
	w := New(seed)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 5000, 100
	own(w, id)
	for day := 0; day < 3; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	w.Player.Cash, w.Player.Location = 5000, id
	return w
}

// morning winds the clock to a minute of the day, so a night can be bought
// before the city goes out or after it has already gone.
func morning(w *World, at int) {
	for w.Minute%1440 != at {
		w.Advance(30)
		w.Event = nil
	}
}

// tonight walks the clock through one whole evening and says how many people
// stood in a room while it was on.
func tonight(w *World, id string) int {
	most := 0
	for i := 0; i < 48; i++ {
		w.Advance(30)
		w.Event = nil
		if Evening(w.Minute) {
			if n := w.Footfall(id); n > most {
				most = n
			}
		}
		if w.Minute%1440 == 0 {
			break
		}
	}
	return most
}

func TestANightFillsARoomNobodyWalksPast(t *testing.T) {
	t.Parallel()
	// Over several cities, because the fault this catches was invisible in a
	// city whose day's work happened to be near the room. One seed proves the
	// draw works where somebody already stands; it has to work everywhere.
	for _, seed := range []uint32{11, 29, 47, 61, 83} {
		w := quiet(t, "burlesque", seed)
		morning(w, 600)
		if err := w.PutOnANight("burlesque"); err != nil {
			t.Fatal(err)
		}
		loud := tonight(w, "burlesque")

		ordinary := quiet(t, "burlesque", seed)
		morning(ordinary, 600)
		usual := tonight(ordinary, "burlesque")
		if loud <= usual {
			t.Fatalf("city %d: a night at the burlesque drew %d where an ordinary evening drew %d",
				seed, loud, usual)
		}
	}
}

// The falloff itself, which no count of a room can show: a room holds people
// posted to it who never go anywhere, so no draw however total ever empties
// one, and a test that watched a room go dark could not fail.
func TestTheWalkDecidesWhoComes(t *testing.T) {
	t.Parallel()
	sample := func(minutes int) int {
		came := 0
		for i := 0; i < 500; i++ {
			if comes(fmt.Sprintf("person-%d", i), minutes) {
				came++
			}
		}
		return came
	}
	next, across := sample(10), sample(NightDraw)
	if next <= across {
		t.Fatalf("%d in 500 came from ten minutes away and %d from %d: the walk did not matter",
			next, across, NightDraw)
	}
	if across == 0 {
		t.Fatal("nobody at all crossed the far edge of the draw")
	}
	if next >= 500 {
		t.Fatal("a band two doors down took every single drinker: that is a switch, not a draw")
	}
}

func TestANightBoughtAfterTheCityHasGoneOutIsTomorrowsNight(t *testing.T) {
	t.Parallel()
	w := quiet(t, "burlesque", 47)
	morning(w, EveningFrom+120)
	if err := w.PutOnANight("burlesque"); err != nil {
		t.Fatal(err)
	}
	// Tonight was decided two hours ago and the band is not playing to a room
	// that filled up somewhere else.
	if w.NightOn("burlesque") {
		t.Fatal("a band booked at nine o'clock played to a crowd that had already gone out")
	}
	ordinary := quiet(t, "burlesque", 47)
	morning(ordinary, EveningFrom+120)
	usual := tonight(ordinary, "burlesque")
	if tonight(w, "burlesque") > usual {
		t.Fatal("people walked in on an evening nobody had been told about")
	}
	// But it is not money thrown away either: it buys the next one.
	morning(w, EveningFrom)
	if !w.NightOn("burlesque") {
		t.Fatal("the night that was bought never came")
	}
	morning(ordinary, EveningFrom)
	if came, would := tonight(w, "burlesque"), tonight(ordinary, "burlesque"); came <= would {
		t.Fatalf("the night that was paid for drew %d where an ordinary evening drew %d", came, would)
	}
}

func TestANightDrawsHardestOnTheRoomNextDoor(t *testing.T) {
	t.Parallel()
	near, far := "casino", "club"
	if TravelMinutes("burlesque", near) >= TravelMinutes("burlesque", far) {
		t.Fatalf("%s is not nearer the burlesque than %s", near, far)
	}
	// Both rooms read on the same pass of the same evening. Walking the clock
	// once per room read two different evenings and the answer moved on its own.
	evening := func(w *World) (int, int) {
		a, b := 0, 0
		for i := 0; i < 48; i++ {
			w.Advance(30)
			w.Event = nil
			if Evening(w.Minute) {
				a = max(a, w.Footfall(near))
				b = max(b, w.Footfall(far))
			}
			if w.Minute%1440 == 0 {
				break
			}
		}
		return a, b
	}
	ordinary := quiet(t, "burlesque", 47)
	morning(ordinary, 600)
	wasNear, wasFar := evening(ordinary)

	w := quiet(t, "burlesque", 47)
	morning(w, 600)
	if err := w.PutOnANight("burlesque"); err != nil {
		t.Fatal(err)
	}
	nowNear, nowFar := evening(w)
	if wasNear-nowNear <= wasFar-nowFar {
		t.Fatalf("the %s lost %d drinkers and the %s %d: the walk did not matter",
			near, wasNear-nowNear, far, wasFar-nowFar)
	}
}

// Two rooms can pay for a band on the same evening, and then the question is
// which one a drinker walks to. The nearer, which is the only answer that does
// not depend on the order the city's addresses happen to be listed in — and
// asked from every room in turn, because asked from one it does not: the list
// puts the right answer first often enough that taking the first match passes.
func TestBetweenTwoNightsADrinkerGoesToTheNearer(t *testing.T) {
	t.Parallel()
	hosts := []string{}
	for _, l := range Locations {
		if PlaysHost(l.ID) {
			hosts = append(hosts, l.ID)
		}
	}
	w := quiet(t, hosts[0], 47)
	morning(w, EveningFrom)
	for _, id := range hosts {
		w.Properties[id].Night = w.Minute
		if !w.NightOn(id) {
			t.Fatalf("%s was given a night and is not holding one", id)
		}
	}
	asked := 0
	for _, from := range hosts {
		nearest := ""
		for _, id := range hosts {
			if id == from || TravelMinutes(from, id) > NightDraw {
				continue
			}
			if nearest == "" || TravelMinutes(from, id) < TravelMinutes(from, nearest) {
				nearest = id
			}
		}
		if nearest == "" {
			continue
		}
		asked++
		for i := 0; i < 500; i++ {
			who := fmt.Sprintf("person-%d", i)
			to := w.theNight(from, who)
			if to == "" {
				continue
			}
			if to != nearest {
				t.Fatalf("out of the %s, %s walked %d minutes to the %s past the %s at %d",
					from, who, TravelMinutes(from, to), to, nearest, TravelMinutes(from, nearest))
			}
			break
		}
	}
	if asked < 2 {
		t.Fatalf("only %d rooms had a choice of nights: the question was never asked", asked)
	}
}
