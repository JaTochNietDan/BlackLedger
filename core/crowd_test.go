package core

import (
	"strings"
	"testing"
)

// The city was given hours and none of it reached the player. Standing in
// Saint Agnes at ten at night with fourteen people in it read exactly like
// standing in it at nine in the morning with two.

func TestAFullRoomSaysSo(t *testing.T) {
	w := New(21)
	// Put most of the city in one room.
	living := 0
	for i := range w.NPCs {
		if !w.NPCs[i].Dead {
			living++
		}
	}
	for i := range w.NPCs {
		if i*3 < living*2 {
			w.NPCs[i].Location, w.NPCs[i].Heading = "bar", ""
		}
	}
	note := w.RoomNote("bar")
	if note == "" {
		t.Fatal("two thirds of the city is in one bar and the room says nothing")
	}
	if !strings.Contains(note, "Saint Agnes") && !strings.Contains(note, "crowd") {
		t.Fatalf("a full room reads %q", note)
	}
}

func TestAnEmptyRoomSaysSo(t *testing.T) {
	w := New(22)
	for i := range w.NPCs {
		w.NPCs[i].Location, w.NPCs[i].Heading = "market", ""
	}
	if note := w.RoomNote("bar"); !strings.Contains(note, "nobody else") {
		t.Fatalf("an empty room reads %q", note)
	}
}

// An unremarkable room is not worth a sentence. A note on every room in the
// city is a note the player stops reading.
func TestAnOrdinaryRoomIsNotWorthASentence(t *testing.T) {
	w := New(23)
	living := w.living()
	spread := 0
	for i := range w.NPCs {
		if w.NPCs[i].Dead {
			continue
		}
		// About a fifteenth of the city in the bar: neither full nor empty.
		if spread < living/15 {
			w.NPCs[i].Location = "bar"
		} else {
			w.NPCs[i].Location = "market"
		}
		w.NPCs[i].Heading = ""
		spread++
	}
	if note := w.RoomNote("bar"); note != "" {
		t.Fatalf("an ordinary room says %q", note)
	}
}

// Somebody out on the street is in neither building, and must not be counted
// into a room they have left.
func TestSomebodyWalkingIsNotInTheRoom(t *testing.T) {
	w := New(24)
	for i := range w.NPCs {
		w.NPCs[i].Location, w.NPCs[i].Heading = "bar", ""
	}
	full := w.InTheRoom("bar")
	n := &w.NPCs[0]
	n.Heading, n.Sets, n.Arrives = "market", 0, w.Minute+30
	if w.InTheRoom("bar") != full-1 {
		t.Fatal("somebody on the street is still being counted into the room they left")
	}
}

// The description is a description. It must never be the thing that decides
// what can be done in a room.
func TestTheRoomNoteChangesNothing(t *testing.T) {
	w := New(25)
	before := len(w.Actions("bar"))
	for i := range w.NPCs {
		w.NPCs[i].Location, w.NPCs[i].Heading = "bar", ""
	}
	if w.RoomNote("bar") == "" {
		t.Fatal("the whole city is in the bar and it is not remarked on")
	}
	crowded := 0
	for _, a := range w.Actions("bar") {
		if a.Subject == "" {
			crowded++
		}
	}
	impersonal := 0
	w2 := New(25)
	for _, a := range w2.Actions("bar") {
		if a.Subject == "" {
			impersonal++
		}
	}
	if crowded != impersonal {
		t.Fatalf("filling the room changed what can be done in it: %d against %d", crowded, impersonal)
	}
	_ = before
}
