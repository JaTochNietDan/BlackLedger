package core

import (
	"regexp"
	"strings"
	"testing"
)

// Every sentence the city can say, read back. A playtest found "against your a
// pair" by eye; this looks for the rest of that family without anybody having
// to read six hundred cards.
var defects = []struct {
	name string
	re   *regexp.Regexp
}{
	{"a format verb that was never filled", regexp.MustCompile(`%!|%\[|%s|%d|%v`)},
	{"an article doubled onto a name", regexp.MustCompile(`\b(your|their|his|her|its) an? `)},

	{"two spaces", regexp.MustCompile(`  `)},
	{"a space before a full stop", regexp.MustCompile(` \.`)},
	{"an empty figure", regexp.MustCompile(`\$-`)},
	{"a sentence that starts in the middle", regexp.MustCompile(`^[a-z]`)},
}

// And the reader is worth no more than what it can see, so it is shown one of
// each first. A sweep that finds nothing is either a clean city or a broken
// reader, and there is no way to tell them apart from the outside.
func TestTheReaderCanSeeEachKindOfDefect(t *testing.T) {
	t.Parallel()
	samples := map[string]string{
		"a format verb that was never filled":  "Paid %!d(MISSING) for it.",
		"an article doubled onto a name":       "Leo had a pair against your a pair.",
		"two spaces":                           "The room  is empty.",
		"a space before a full stop":           "Nobody is here .",
		"an empty figure":                      "It cost $-40 to settle.",
		"a sentence that starts in the middle": "nobody is here.",
	}
	for name, text := range samples {
		found := ""
		for _, d := range defects {
			if d.re.MatchString(text) {
				found = d.name
				break
			}
		}
		if found != name {
			t.Errorf("%q should read as %q and reads as %q", text, name, found)
		}
	}
	// And it does not cry wolf over ordinary prose.
	for _, fine := range []string{
		"Buy a used Ford", "$1,750 and two days on the bench.",
		"Leo Carver had three of a kind against your pair.",
		"They are not behind the counter at Bluebird Laundry this morning.",
	} {
		for _, d := range defects {
			if d.re.MatchString(fine) {
				t.Errorf("%q is ordinary prose and reads as %q", fine, d.name)
			}
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
		for _, d := range defects {
			if d.re.MatchString(text) {
				bad++
				t.Errorf("%s: %s — %q", where, d.name, text)
				break
			}
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
		t.Fatalf("only %d sentences read, so this measures nothing", said)
	}
}
