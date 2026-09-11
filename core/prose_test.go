package core

import (
	"strings"
	"testing"
)

// Every sentence the city can say, read back. A playtest found "against your a
// pair" by eye; this looks for the rest of that family without anybody having
// to read six hundred cards.
// The list of faults and the reader itself live in core/malformed.go, so the
// log and the paper can be read by the same eye that reads the cards.
func TestTheReaderCanSeeEachKindOfDefect(t *testing.T) {
	t.Parallel()
	samples := map[string]string{
		"a format verb that was never filled":   "Paid %!d(MISSING) for it.",
		"an article doubled onto a name":        "Leo had a pair against your a pair.",
		"two spaces":                            "The room  is empty.",
		"a space before a full stop":            "Nobody is here .",
		"an empty figure":                       "It cost $-40 to settle.",
		"a sentence that starts in the middle":  "nobody is here.",
		"a count of none written as a quantity": "They took 0 crates and left.",
		"one of something written as several":   "$40 for the 1 days still on them.",
	}
	for name, text := range samples {
		if found := Malformed(text); found != name {
			t.Errorf("%q should read as %q and reads as %q", text, name, found)
		}
	}
	// And it does not cry wolf over ordinary prose.
	for _, fine := range []string{
		"Buy a used Ford", "$1,750 and two days on the bench.",
		"Leo Carver had three of a kind against your pair.",
		"They are not behind the counter at Bluebird Laundry this morning.",
	} {
		if found := Malformed(fine); found != "" {
			t.Errorf("%q is ordinary prose and reads as %q", fine, found)
		}
	}
	// Every fault the reader knows about has a sample above it, or the sweeps
	// below are trusting an eye that has never been tested.
	for _, kind := range MalformedKinds() {
		if _, ok := samples[kind]; !ok {
			t.Errorf("nothing shows the reader what %q looks like", kind)
		}
	}
}

func TestNothingTheCitySaysIsMalformed(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 40, 50000
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["poolhall"].Owner = "player:1"

	said, bad := 0, 0
	check := func(where, text string) {
		if strings.TrimSpace(text) == "" {
			return
		}
		said++
		if found := Malformed(text); found != "" {
			bad++
			t.Errorf("%s: %s — %q", where, found, text)
		}
	}
	for _, l := range Locations {
		if l.District > w.District {
			continue
		}
		w.Player.Location, w.Event = l.ID, nil
		for _, a := range w.Actions(l.ID) {
			check(l.ID+" "+a.ID+" label", a.Label)
			check(l.ID+" "+a.ID+" detail", a.Detail)
			check(l.ID+" "+a.ID+" reason", a.Reason)
		}
		check(l.ID+" note", w.PlaceNote(l.ID))
		check(l.ID+" room", w.RoomNote(l.ID))
	}
	t.Logf("%d sentences read, %d malformed", said, bad)
	if said < 500 {
		t.Fatalf("only %d sentences read across forty campaigns, so this measures nothing", said)
	}
}

// And the same eye over what the city says while it is running.
//
// The sweep above reads the cards, which is what the game says while somebody is
// deciding. The log and the paper are what it says while something is happening,
// and nothing had ever read those: "0 crates of arms out through the front door
// in daylight" was in the log, not on a card, and it was found by walking a
// chain rather than by any reader.
//
// Passive campaigns across a spread of seeds, because a city left alone still
// takes premises off people, starts wars, fills desks and buries the dead, and
// every one of those writes a sentence.
func TestNothingTheCitySaysWhileItRunsIsMalformed(t *testing.T) {
	t.Parallel()
	said, bad := 0, 0
	check := func(where, text string) {
		if strings.TrimSpace(text) == "" {
			return
		}
		said++
		if found := Malformed(text); found != "" {
			bad++
			t.Errorf("%s: %s — %q", where, found, text)
		}
	}
	for seed := uint32(1); seed <= 40; seed++ {
		w := New(seed)
		w.Event, w.District = nil, 9
		w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 120, 500000
		for _, id := range []string{"laundry", "poolhall", "garage"} {
			if prop := w.Properties[id]; prop != nil {
				prop.Owner = "player:1"
			}
		}
		// A quiet city says quiet things. Half of these campaigns are run hot —
		// a room under a floor with crates in it, a still, stock on the player
		// and the police already interested — because the sentences that were
		// wrong were written by raids, seizures and arrests rather than by the
		// weather.
		if seed%2 == 0 {
			w.Player.Heat = 70
			if prop := w.Properties["laundry"]; prop != nil {
				prop.Armoury, prop.Crates, prop.Still = true, 3, true
			}
			w.Player.Stock = map[string]int{}
			for _, g := range w.Goods {
				w.Player.Stock[g.ID] = 6
			}
			w.Player.Car = 2
		}
		w.Advance(1440 * 60)
		for _, r := range w.History {
			check("record title", r.Title)
			check("record text", r.Text)
		}
		for _, s := range w.News {
			check("headline", s.Headline)
			check("story", s.Body)
		}
		// A scene the city raised is a sentence too, and it is the one the
		// player is made to read.
		if w.Event != nil {
			check("scene title", w.Event.Title)
			check("scene body", w.Event.Body)
			for _, c := range w.Event.Choices {
				check("scene choice", c.Label)
				check("scene detail", c.Detail)
			}
		}
	}
	t.Logf("%d sentences read across forty campaigns, %d malformed", said, bad)
	if said < 10000 {
		t.Fatalf("only %d sentences read across forty campaigns, so this measures nothing", said)
	}
}
