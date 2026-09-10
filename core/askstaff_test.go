package core

import (
	"strings"
	"testing"
)

// The people behind your counters have names, a wage and a view of the street,
// and nothing to say. They notice a car parked across the road on their own
// schedule; a player standing in front of one should be able to ask.
//
// Everything the answer contains is a fact the city already holds: who has been
// in, what the trade is doing, and whether somebody has been looking the place
// over. Nothing is invented for it.

func employer(t *testing.T) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	w.Player.Location = "laundry"
	return w, "laundry"
}

func TestYouCanAskYourOwnPeopleWhatTheyHaveSeen(t *testing.T) {
	t.Parallel()
	w, id := employer(t)
	hand := w.Properties[id].Hands[0]
	var ask *Action
	for _, a := range w.Actions(id) {
		if a.ID == "ask:"+hand {
			ask = &a
		}
	}
	if ask == nil {
		t.Fatalf("%s works behind this counter and there is no way to ask them anything", w.NPC(hand).Name)
	}
	if ask.Disabled {
		t.Fatalf("asking your own employee is refused: %s", ask.Reason)
	}
	if ask.Subject != hand {
		t.Fatalf("a question for %s is filed under %q", w.NPC(hand).Name, ask.Subject)
	}
	if err := w.apply(Command{Kind: "ask:" + hand, Target: id, RequestID: "askthecounter1"}); err != nil {
		t.Fatalf("asking was refused: %v", err)
	}
	last := w.History[len(w.History)-1]
	if !strings.Contains(last.Title, w.NPC(hand).Name) {
		t.Fatalf("the answer is filed under %q", last.Title)
	}
	if last.Text == "" {
		t.Fatal("they were asked and said nothing at all")
	}
}

// Somebody who works for a rival is not going to tell you anything, and nobody
// who works for nobody has a counter to watch.
func TestYouCanOnlyAskThePeopleWhoWorkForYou(t *testing.T) {
	t.Parallel()
	w, id := employer(t)
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	w.Properties["butcher"].Owner = rival
	w.EmptyChairs()
	theirs := w.Properties["butcher"].Hands[0]
	// Standing in their shop, in front of their employee.
	w.Player.Location = "butcher"
	for _, a := range w.Actions("butcher") {
		if a.ID == "ask:"+theirs && !a.Disabled {
			t.Fatalf("%s works for a rival and answers your questions", w.NPC(theirs).Name)
		}
	}
	// And refused if it is asked for anyway, which is the half the panel does
	// not cover: the question is only offered inside a room you hold, so
	// removing the check on who they work for failed nothing at all.
	if err := w.AskTheCounter(theirs); err == nil {
		t.Fatalf("%s works for a rival and answered", w.NPC(theirs).Name)
	}
	// And the answer comes from the room they actually stand in.
	w.Player.Location = id
	mine := w.Properties[id].Hands[0]
	if err := w.AskTheCounter(mine); err != nil {
		t.Fatalf("asking your own was refused: %v", err)
	}
	place, _ := PlaceByID(id)
	if !strings.Contains(w.History[len(w.History)-1].Text, place.Name) {
		t.Fatalf("the answer does not mention the room it came from: %q", w.History[len(w.History)-1].Text)
	}
}

// And what they say is drawn from what the city holds rather than invented: a
// place somebody is casing reads differently from a quiet one.
func TestWhatTheySayComesFromWhatIsActuallyHappening(t *testing.T) {
	t.Parallel()
	quiet, id := employer(t)
	watched, _ := employer(t)
	rival := watched.Factions[0].ID
	if rival == watched.PlayerOrganizationID() {
		rival = watched.Factions[1].ID
	}
	watched.Plots = append(watched.Plots, Plot{ID: ID(), Kind: "sabotage", Life: watched.Life,
		Due: watched.Minute + 4320, Actor: rival, Target: id, Strength: 35, Known: true})

	if err := quiet.AskTheCounter(quiet.Properties[id].Hands[0]); err != nil {
		t.Fatal(err)
	}
	if err := watched.AskTheCounter(watched.Properties[id].Hands[0]); err != nil {
		t.Fatal(err)
	}
	calm := quiet.History[len(quiet.History)-1].Text
	uneasy := watched.History[len(watched.History)-1].Text
	if calm == uneasy {
		t.Fatalf("a laundry somebody is looking over reads exactly like a quiet one: %q", calm)
	}
	if !strings.Contains(uneasy, "car") {
		t.Fatalf("somebody casing the place is not mentioned: %q", uneasy)
	}
}
