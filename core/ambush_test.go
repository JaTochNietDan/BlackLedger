package core

import (
	"strings"
	"testing"
)

// "Can people only make attempts on your life while you're at home? They always
// seem to hit my home when I'm not there and they are coming after me."
//
// They came for the address rather than the person. A hit finds you where you
// are, and where that is decides how much of a chance you get.

func marked(seed uint32) *World {
	w := New(seed * 2654435761)
	w.Player.Health, w.Player.Contacts, w.Player.Security = 100, 0, 0
	// They have found the player. A family knows an address and not an evening,
	// so an attack away from home now goes to the house unless somebody of
	// theirs has laid eyes on you — which is the thing these tests are not
	// about. They are about what happens once you have been found.
	w.SeenBy = map[string]int{"bellandi": w.Minute}
	return w
}

// clearRoom empties an address, because the city already puts people in its
// bars and a test that does not say how many are standing there is measuring
// whatever the seed happened to do.
func clearRoom(w *World, id string) {
	for i := range w.NPCs {
		if w.NPCs[i].Location == id {
			w.NPCs[i].Location = "market"
		}
	}
}

// fill puts n people who are nobody's into a room.
func fill(w *World, id string, n int) int {
	put := 0
	for i := range w.NPCs {
		if put >= n {
			break
		}
		if !w.NPCs[i].Dead && w.NPCs[i].Location != id && w.NPCs[i].Faction != w.PlayerOrganizationID() {
			w.NPCs[i].Location = id
			put++
		}
	}
	return put
}

func hit(w *World) Plot {
	return Plot{ID: ID(), Kind: "hit", Life: w.Life, Due: w.Minute, Actor: "bellandi"}
}

func TestAHitFindsYouWhereYouAre(t *testing.T) {
	t.Parallel()
	const runs = 200
	untouched := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := marked(seed)
		w.Player.Location = "bar"
		before := w.Properties[w.Player.Home].Condition
		health := w.Player.Health
		w.Attack(hit(w))
		if w.Properties[w.Player.Home].Condition != before {
			t.Fatalf("they wrecked the house while the player was standing in a bar (seed %d)", seed)
		}
		if w.Player.Alive && w.Event == nil && w.Player.Health == health {
			untouched++
		}
	}
	if untouched > 0 {
		t.Fatalf("%d of %d attacks away from home did nothing to anybody", untouched, runs)
	}
}

func TestACrowdedRoomIsNotAnEmptyStreet(t *testing.T) {
	t.Parallel()
	const runs = 400
	// Deliberately below the number of people it takes for the room to warn
	// you. A scene raised is a player still standing, so counting survivors of
	// a warned attack measures the warning and not the room — the test would
	// pass whatever cover was worth.
	// Three arms, because two would not tell them apart: with cover deleted the
	// street is still worse than a bar, on the street penalty alone.
	room, empty, street := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		busy := marked(seed)
		busy.Player.Location = "bar"
		clearRoom(busy, "bar")
		put := fill(busy, "bar", SeenComingCrowd-1)
		if busy.Warned(hit(busy)) {
			t.Fatalf("a room of %d warned the player, so this measures nothing", put)
		}
		busy.Attack(hit(busy))
		if busy.Player.Alive {
			room++
		}

		bare := marked(seed)
		bare.Player.Location = "bar"
		clearRoom(bare, "bar")
		bare.Attack(hit(bare))
		if bare.Player.Alive {
			empty++
		}

		alone := marked(seed)
		alone.Player.Location = "transit"
		if alone.Warned(hit(alone)) {
			t.Fatal("the street warned the player")
		}
		alone.Attack(hit(alone))
		if alone.Player.Alive {
			street++
		}
	}
	t.Logf("of %d unwarned attacks: survived with two people in the bar %d, in an empty bar %d, on the street %d",
		runs, room, empty, street)
	if room <= empty {
		t.Fatalf("people in the room were worth nothing: %d against %d in an empty one", room, empty)
	}
	if empty <= street {
		t.Fatalf("four walls were worth nothing: %d against %d on the street", empty, street)
	}
}

func TestABusyRoomSeesThemComing(t *testing.T) {
	t.Parallel()
	w := marked(7)
	w.Player.Location = "bar"
	clearRoom(w, "bar")
	if w.Warned(hit(w)) {
		t.Fatal("an empty bar warned the player")
	}
	put := fill(w, "bar", SeenComingCrowd)
	if !w.Warned(hit(w)) {
		t.Fatalf("a room of %d saw nothing", put)
	}
}

func TestYouStillGetTheChoiceAwayFromHome(t *testing.T) {
	t.Parallel()
	w := marked(3)
	w.Player.Location, w.Player.Contacts = "bar", 3
	w.Player.Security = 2
	w.Attack(hit(w))
	if w.Event == nil {
		t.Fatal("a warned player away from home got no scene at all")
	}
	if w.Event.Kind != "attack" {
		t.Fatalf("the scene was a %s", w.Event.Kind)
	}
	ids := map[string]bool{}
	for _, c := range w.Event.Choices {
		ids[c.ID] = true
	}
	for _, want := range []string{"escape", "defend", "bargain"} {
		if !ids[want] {
			t.Fatalf("no way to %s away from home", want)
		}
	}
	place, _ := PlaceByID("bar")
	if !strings.Contains(w.Event.Body, place.Name) {
		t.Fatalf("the scene never said where it was happening: %q", w.Event.Body)
	}
}

func TestTheyTakeItOutOnTheHouseOnlyWhenTheyCannotReachYou(t *testing.T) {
	t.Parallel()
	w := marked(11)
	w.Player.HeldUntil = w.Minute + 2*1440
	w.Player.Location = "precinct"
	before := w.Properties[w.Player.Home].Condition
	health := w.Player.Health
	w.Attack(hit(w))
	if w.Properties[w.Player.Home].Condition >= before {
		t.Fatal("a player nobody could reach lost nothing at all")
	}
	if w.Player.Health != health || !w.Player.Alive {
		t.Fatal("they got to somebody in a police cell")
	}
}

// The street is not a room.
//
// `Warned` had a branch saying "the people who watch your door are at your
// door", returning false for somebody caught between two addresses — and it
// could never run, because a flat check on the player's reach sat above it and
// answered first. Two cups of coffee bought a warning on an empty street where
// by that same function's account there is nobody to give one.
//
// The order is the fix, and this is the guard: in a room, contacts are enough;
// on the street, they are not, and it takes everything the city offers.
func TestContactsDoNotFollowYouIntoTheStreet(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash = 100, 20000
	w.Player.Contacts = RoomWarning
	if w.Reach() < RoomWarning {
		t.Fatalf("a reach of %d with %d contacts", w.Reach(), w.Player.Contacts)
	}
	// Standing in a room, that is enough to be told.
	w.Player.Location = "bar"
	if w.InTransit() {
		t.Fatal("standing still and in transit")
	}
	if !w.Warned(Plot{}) {
		t.Fatalf("a reach of %d in a room and nobody said anything", w.Reach())
	}
	// On the street between two addresses, it is not.
	// Between two addresses is a location that is not one.
	w.Player.Location = "on the way to Pier 14"
	if !w.InTransit() {
		t.Skip("this build has no journey to be caught on")
	}
	if w.Warned(Plot{}) {
		t.Fatalf("a reach of %d told somebody on an empty street", w.Reach())
	}
	// And enough of it is enough anywhere.
	w.Player.Contacts = StreetWarning
	if w.Reach() < StreetWarning {
		t.Skipf("a reach of %d is the most this city offers", w.Reach())
	}
	if !w.Warned(Plot{}) {
		t.Fatalf("every contact in the city, a reach of %d, and still nothing", w.Reach())
	}
}
