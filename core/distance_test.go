package core

import (
	"strings"
	"testing"
)

// The city view showed twelve addresses and gave no sense of how far any of
// them was: a place across town and one on the next corner looked identical,
// and the journey time only appeared after the player had already decided where
// they were going.

func TestEveryAddressSaysHowFarItIs(t *testing.T) {
	w := proprietor(t)
	w.District = 2
	w.Player.Location = "bar"
	if w.Away("bar") != 0 || w.TravelNote("bar") != "You are here" {
		t.Fatalf("standing at the bar reads %q", w.TravelNote("bar"))
	}
	said := map[string]bool{}
	for _, l := range Locations {
		note := w.TravelNote(l.ID)
		if note == "" {
			t.Fatalf("%s does not say how far it is", l.Name)
		}
		if l.ID != "bar" && w.Away(l.ID) <= 0 {
			t.Fatalf("%s is %d minutes away", l.Name, w.Away(l.ID))
		}
		said[w.ReachWord(l.ID)] = true
	}
	// The words have to distinguish somewhere, or they say nothing.
	if len(said) < 3 {
		t.Fatalf("twelve addresses are described in %d ways: %v", len(said), said)
	}
	t.Logf("from Saint Agnes the city reads as %d different distances", len(said))
}

func TestDrivingIsNeverSlowerThanWalking(t *testing.T) {
	w := proprietor(t)
	w.District = 2
	w.Player.Location = "room"
	onFoot := map[string]int{}
	for _, l := range Locations {
		onFoot[l.ID] = w.Away(l.ID)
	}
	// Give them a car and the same journeys must not get longer.
	// CarWear is the car's condition rather than its wear, despite the name:
	// a hundred is a car that runs.
	w.Player.Car, w.Player.CarWear = 1, 100
	if !w.Driving() {
		t.Skip("a car of that tier does not run")
	}
	shorter := 0
	for _, l := range Locations {
		if driving := w.Away(l.ID); driving > onFoot[l.ID] {
			t.Fatalf("%s takes %d minutes driving and %d on foot", l.Name, driving, onFoot[l.ID])
		} else if driving < onFoot[l.ID] {
			shorter++
		}
	}
	if shorter == 0 {
		t.Fatal("a car made no journey in this city any shorter")
	}
	// And the note says both, so the player can see what the car is buying.
	for _, l := range Locations {
		if l.ID == w.Player.Location {
			continue
		}
		if note := w.TravelNote(l.ID); !strings.Contains(note, "on foot") {
			t.Fatalf("driving to %s reads %q and never mentions walking", l.Name, note)
		}
		break
	}
}

func TestDistanceFollowsWhereYouAreStanding(t *testing.T) {
	w := proprietor(t)
	w.District = 2
	w.Player.Location = "room"
	far := w.Away("estate")
	w.Player.Location = "estate"
	if w.Away("estate") != 0 {
		t.Fatal("standing somewhere is not being there")
	}
	if w.Away("room") == 0 {
		t.Fatal("the room is nowhere now")
	}
	// Roughly symmetric: the city is not one-way.
	if back := w.Away("room"); back != far {
		t.Fatalf("estate to room is %d and room to estate is %d", back, far)
	}
}
