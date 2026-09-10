package core

import "testing"

// "Let's verify, can you be attacked while traversing the map? Can your car
// affect this, whether it's armored, etc? This should all be the case."
//
// It is the case, and this is the proof rather than the assurance: the street
// is the most dangerous place in this city to be found, a plated car is the
// only cover out there, and what you are wearing helps wherever you are. Asked
// as numbers so it stays true.

func crossing(t *testing.T) *World {
	t.Helper()
	w := New(83)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 9000, 100
	w.Player.Location = "bar"
	return w
}

// onTheRoad puts the player between two addresses, which is what travelling is:
// a location the city has no address for.
func onTheRoad(w *World) { w.Player.Location = "transit" }

func TestYouCanBeCaughtOnTheRoad(t *testing.T) {
	w := crossing(t)
	if w.InTransit() {
		t.Fatal("a player standing in a bar is on the road")
	}
	onTheRoad(w)
	if !w.InTransit() {
		t.Fatal("a player between two addresses is not on the road")
	}
	if w.whereItHappens() != "the street" {
		t.Fatalf("an attack on the road happens at %q", w.whereItHappens())
	}
	// Nobody is watching for you out there, so it is the one place an attack
	// arrives without a warning first.
	w.Player.Security, w.Player.Contacts = 3, 0
	if w.Warned(Plot{Kind: "hit"}) {
		t.Fatal("somebody on the road was warned by the people who watch their door")
	}
}

func TestTheStreetIsTheWorstPlaceToBeFound(t *testing.T) {
	room := crossing(t)
	road := crossing(t)
	onTheRoad(road)
	inside, out := room.AmbushOddsHere(), road.AmbushOddsHere()
	t.Logf("an attack kills on %.0f%% in a bar and %.0f%% on the road", inside*100, out*100)
	if out <= inside {
		t.Fatalf("the open street is no worse than a room with people in it: %.2f against %.2f", out, inside)
	}
	// And the street is worse than an empty room, which is the part that is
	// about the street rather than about the people in the bar. Without this
	// the guard passed with the exposure set to nothing, because a bar with
	// drinkers in it is cover and a road is not.
	empty := crossing(t)
	empty.Player.Location = "estate"
	for i := range empty.NPCs {
		empty.NPCs[i].Location = "docks"
	}
	bare := empty.AmbushOddsHere()
	t.Logf("an empty room is %.0f%% against the road's %.0f%%", bare*100, out*100)
	if out <= bare {
		t.Fatalf("being out on the road costs nothing over standing in an empty room: %.2f against %.2f", out, bare)
	}
}

// A car is not cover. A plated one is, and that is the whole of what paying for
// the plate buys.
func TestAPlatedCarIsTheOnlyCoverOnTheRoad(t *testing.T) {
	bare := crossing(t)
	onTheRoad(bare)
	bare.Player.Car, bare.Player.CarWear, bare.Player.Fuel = 2, 100, 40

	plated := crossing(t)
	onTheRoad(plated)
	plated.Player.Car, plated.Player.CarWear, plated.Player.Fuel = 2, 100, 40
	plated.Player.Plate = PlateStages

	walking := crossing(t)
	onTheRoad(walking)
	walking.Player.Car = 0

	drove, armoured, walked := bare.AmbushOddsHere(), plated.AmbushOddsHere(), walking.AmbushOddsHere()
	t.Logf("on the road: %.0f%% on foot, %.0f%% in a car, %.0f%% in a plated one", walked*100, drove*100, armoured*100)
	if armoured >= drove {
		t.Fatalf("plating the car bought nothing: %.2f against %.2f", armoured, drove)
	}
	if drove != walked {
		t.Fatalf("an ordinary car is cover: %.2f against %.2f on foot", drove, walked)
	}
	if plated.PlateCoverHere() <= 0 {
		t.Fatal("a plated car is worth nothing on the road")
	}
	// And it is worth nothing anywhere else, because you are not in it.
	inside := crossing(t)
	inside.Player.Car, inside.Player.CarWear, inside.Player.Fuel = 2, 100, 40
	inside.Player.Plate = PlateStages
	if inside.PlateCoverHere() != 0 {
		t.Fatal("a car parked outside covered somebody standing in a bar")
	}
}

// What you are wearing is yours wherever you are, which is the difference
// between armour and cover.
func TestWhatYouAreWearingHelpsOnTheRoadAndInTheRoom(t *testing.T) {
	for _, where := range []string{"a bar", "the road"} {
		bare, worn := crossing(t), crossing(t)
		if where == "the road" {
			onTheRoad(bare)
			onTheRoad(worn)
		}
		worn.Player.Armour = 3
		if worn.AmbushOddsHere() >= bare.AmbushOddsHere() {
			t.Errorf("armour was worth nothing in %s: %.2f against %.2f", where, worn.AmbushOddsHere(), bare.AmbushOddsHere())
		}
	}
}

// And the crossing says so before the player sets off, which is the half of
// this that is worth anything to somebody deciding whether to go.
func TestTheCrossingSaysWhatItIsWalkingInto(t *testing.T) {
	w := crossing(t)
	w.Plots = append(w.Plots, Plot{ID: ID(), Kind: "hit", Life: w.Life,
		Due: w.Minute + 600, Actor: w.Factions[0].ID, Known: true, Strength: 40})
	quiet := crossing(t)
	told := w.Crossing("bar", "docks")
	said := quiet.Crossing("bar", "docks")
	if told["warned"] != true {
		t.Fatal("somebody has a price on the player's head and the crossing does not say so")
	}
	if said["warned"] != false {
		t.Fatal("a crossing with nobody after you says somebody is")
	}
	if told["note"] == "" {
		t.Fatal("the crossing is warned and says nothing about it")
	}
	// Driving changes what it says, because it changes what is true.
	w.Player.Car, w.Player.CarWear, w.Player.Fuel = 2, 100, 40
	w.Player.Plate = PlateStages
	if driven := w.Crossing("bar", "docks"); driven["note"] == told["note"] {
		t.Fatal("the crossing reads the same on foot as in a plated car")
	}
}

// The other half of what a car is worth on the road: how long you are on it. A
// journey a car makes in ten minutes is ten minutes of being findable, and the
// same journey on foot is longer. Nothing has to be added for this to be true —
// it falls out of the clock — but it is worth a number, because "can your car
// affect this" is a fair question to want an answer to.
func TestACarShortensHowLongYouAreFindable(t *testing.T) {
	walking := crossing(t)
	driving := crossing(t)
	driving.Player.Car, driving.Player.CarWear, driving.Player.Fuel = 2, 100, 40

	onFoot, byCar := 0, 0
	for _, to := range []string{"docks", "club", "market", "casino", "archway"} {
		onFoot += walking.Journey("bar", to)
		byCar += driving.Journey("bar", to)
	}
	t.Logf("five crossings out of the bar: %d minutes on foot, %d driving", onFoot, byCar)
	if byCar >= onFoot {
		t.Fatalf("a car saved no time at all: %d minutes against %d", byCar, onFoot)
	}
	// And a wreck is not a car. Somebody whose car will not start is walking,
	// and is exposed for exactly as long as somebody who never had one.
	wreck := crossing(t)
	wreck.Player.Car, wreck.Player.CarWear, wreck.Player.Fuel = 2, 0, 40
	if wreck.Journey("bar", "docks") != walking.Journey("bar", "docks") {
		t.Fatal("a car that will not start still got somebody across the city faster")
	}
	// A dry tank is the same thing by another route. Fuelled has to be set for
	// the tank to read empty: a car nobody has ever put petrol into is a car
	// that came with a tank, which is what keeps every save written before
	// petrol existed from standing at the kerb.
	dry := crossing(t)
	dry.Player.Car, dry.Player.CarWear = 2, 100
	dry.Player.Fuel, dry.Player.Fuelled = 0, 1
	if dry.Journey("bar", "docks") != walking.Journey("bar", "docks") {
		t.Fatal("a car with nothing in the tank still got somebody across the city faster")
	}
}
