package core

import (
	"strings"
	"testing"
)

// What a counter is worth talking to for.
//
// Asking the person behind one used to answer with the state of the premises:
// the wages, the trouble, whether it was short-handed, and how many came
// through the door. All of that is on "Review the books", which costs nothing
// and no time, so the whole of what the conversation added was a footfall
// count the game already had.
//
// A person who stands in a room all day sees three things a ledger never
// records: who is walking over, who has been asking after you, and the one
// thing this trade is in a position to notice. Those are what make the counter
// worth a visit.

func counterAt(t *testing.T, kind string) (*World, string) {
	t.Helper()
	id := ""
	for _, l := range Locations {
		if l.Kind == kind {
			id = l.ID
		}
	}
	if id == "" {
		t.Fatalf("the city has no %s", kind)
	}
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 90
	own(w, id)
	for d := 0; d < 8; d++ {
		w.Advance(1440)
		w.Event = nil
	}
	w.Player.Location = id
	return w, id
}

// answer is what the counter said, as one line.
func answer(t *testing.T, w *World, id string) string {
	t.Helper()
	prop := w.Properties[id]
	if len(prop.Hands) == 0 {
		t.Fatalf("nobody is behind the counter at %s", id)
	}
	n := w.NPC(prop.Hands[0])
	n.Location, n.Heading, n.Arrives = id, "", 0
	before := len(w.History)
	if err := w.AskTheCounter(n.ID); err != nil {
		t.Fatal(err)
	}
	said := []string{}
	for _, r := range w.History[before:] {
		said = append(said, r.Text)
	}
	return strings.Join(said, " ")
}

// everyTrade is each kind of business in the city, with an address that runs
// it. Written out of the city rather than listed here, because a list is how
// the first version of this asked five kinds out of seventeen and passed while
// two thirds of the city's counters said nothing at all.
func everyTrade() map[string]string {
	out := map[string]string{}
	for _, l := range Locations {
		if l.Kind != "" && out[l.Kind] == "" {
			out[l.Kind] = l.ID
		}
	}
	return out
}

func TestEachCounterSeesSomethingItsOwn(t *testing.T) {
	t.Parallel()
	// No two trades say the same thing, because the whole reason to hold more
	// than one kind of place is that each shows you a different part of the
	// same city.
	seen := map[string]string{}
	trades := everyTrade()
	if len(trades) < 15 {
		t.Fatalf("only %d kinds in this city, so this proves very little", len(trades))
	}
	for kind := range trades {
		w, id := counterAt(t, kind)
		line := w.fromBehindThisCounter(id)
		if line == "" {
			t.Errorf("%s: the counter has nothing of its own to say", kind)
			continue
		}
		if was, twice := seen[line]; twice {
			t.Errorf("%s and %s both say %q", was, kind, line)
		}
		seen[line] = kind
		// A count of none is a word, not a numeral. "0 people in this city
		// used to drive" is the shape the prose guards exist to catch, and a
		// counter that only ever reads on a busy city will print it the first
		// quiet week.
		if strings.HasPrefix(line, "0 ") || strings.Contains(line, " 0 ") {
			t.Errorf("%s counts in numerals where it should use a word: %q", kind, line)
		}
		// And it reaches the player, rather than being a function nothing calls.
		if !strings.Contains(answer(t, w, id), line) {
			t.Errorf("%s: the counter knows it and does not say it", kind)
		}
	}
}

func TestTheCounterSeesWhoIsWalkingOver(t *testing.T) {
	t.Parallel()
	w, id := counterAt(t, "butcher")
	// The fixture advances eight days to establish the business. Ordinary
	// commutes can already be underway; isolate the journey being tested.
	for i := range w.NPCs {
		w.NPCs[i].Heading = ""
		w.NPCs[i].Sets, w.NPCs[i].Arrives = 0, 0
	}
	if w.walkingOver(id) != "" {
		t.Fatal("somebody is walking over before anybody set off")
	}
	var coming *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Location != id {
			coming = n
			break
		}
	}
	coming.Heading, coming.Sets, coming.Arrives = id, 0, w.Minute+40
	coming.Errand, coming.Faction = "coming to buy something", ""
	plain := w.walkingOver(id)
	if !strings.Contains(plain, coming.Name) {
		t.Fatalf("the counter watched somebody walk over and did not name them: %q", plain)
	}
	// Whose man he is, when he is anybody's, because that is the part worth a
	// warning rather than a remark.
	coming.Faction = "bellandi"
	whose := w.walkingOver(id)
	if !strings.Contains(whose, "Bellandi") {
		t.Fatalf("a rival's man walked over as an ordinary customer: %q", whose)
	}
	if whose == plain {
		t.Fatal("it made no difference whose he was")
	}
	// Somebody already standing in the room is not walking over.
	coming.Heading, coming.Arrives = "", 0
	if w.walkingOver(id) != "" {
		t.Fatalf("nobody is on the street and the counter says: %q", w.walkingOver(id))
	}
}

func TestTheCounterHearsWhoHasBeenAsking(t *testing.T) {
	t.Parallel()
	w, id := counterAt(t, "butcher")
	if w.beenAsking(id) != "" {
		t.Fatalf("somebody has been asking before anybody was wronged: %q", w.beenAsking(id))
	}
	var near, far *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location == "" {
			continue
		}
		d := TravelMinutes(n.Location, id)
		if near == nil && d <= AskingDistance {
			near = n
		}
		if far == nil && d > AskingDistance {
			far = n
		}
	}
	if near == nil || far == nil {
		t.Skip("this city is all one distance from the butcher")
	}
	// Somebody brooding across town is not somebody this room has heard about.
	w.Aggrieve(far.ID, SoreActs+5, "a thing that happened")
	if said := w.beenAsking(id); said != "" {
		t.Fatalf("%s is %d minutes away and the counter has heard: %q",
			far.Name, TravelMinutes(far.Location, id), said)
	}
	// A little sore is not yet worth mentioning.
	w.Aggrieve(near.ID, SoreAsks-1, "a small thing")
	if said := w.beenAsking(id); said != "" {
		t.Fatalf("a grudge of %d was worth a warning: %q", near.Sore, said)
	}
	// Enough, and they say so; more, and they say what it sounded like.
	w.Aggrieve(near.ID, SoreAsks, "a larger thing")
	asking := w.beenAsking(id)
	if !strings.Contains(asking, near.Name) {
		t.Fatalf("somebody has been asking and the counter did not name them: %q", asking)
	}
	w.Aggrieve(near.ID, SoreActs, "a larger thing")
	worse := w.beenAsking(id)
	if worse == asking {
		t.Fatalf("somebody past acting on it sounds the same as somebody curious: %q", worse)
	}
}

func TestTheCounterSaysMoreThanTheBooksDo(t *testing.T) {
	t.Parallel()
	// The whole complaint. Everything the counter used to say was premises
	// state, and the books give that away for nothing.
	w, id := counterAt(t, "garage")
	var coming *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Location != id {
			coming = n
			coming.Heading, coming.Sets, coming.Arrives = id, 0, w.Minute+40
			coming.Faction = "bellandi"
			break
		}
	}
	said := answer(t, w, id)
	for _, mustSay := range []string{coming.Name, "Bellandi"} {
		if !strings.Contains(said, mustSay) {
			t.Fatalf("the counter never mentions %q: %s", mustSay, said)
		}
	}
	// None of which is anywhere in the books.
	books := ""
	if prop := w.Properties[id]; prop != nil {
		books = w.PlaceNote(id)
	}
	if strings.Contains(books, coming.Name) {
		t.Fatalf("the books already knew: %s", books)
	}
}
