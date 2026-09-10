package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"blackledger/core"
)

// The interface is meant to stay functional so the game can be play-tested, and
// for several iterations it was not: a useEffect had been written below the
// early return that renders while the world is still loading. React counts
// hooks, the count changed between the first render and the second, and every
// boot ended in error #310 and a blank page. Every check in this repository
// passed the whole time, because none of them ever opened the page.
//
// This is the cheapest guard that would have caught it. It is a text check
// rather than a parser, and it is deliberately narrow: a hook below a top-level
// early return in a component is the specific mistake that cost days of play.

var (
	hookCall = regexp.MustCompile(`\buse(State|Effect|Memo|Callback|Ref|Reducer|Context|LayoutEffect)\s*\(`)
	// Only a return in the component's own body can change how many hooks run.
	// This first matched any depth, so `if (!from || !to) return [];` six
	// spaces deep inside a callback was read as an early return and every hook
	// after it reported. A guard that cries wolf gets reformatted around rather
	// than obeyed, so it is pinned to the one or two spaces of indentation a
	// component body uses in this codebase.
	earlyReturn  = regexp.MustCompile(`^ {1,2}if\s*\(.*\)\s*return\s`)
	componentTop = regexp.MustCompile(`^(export\s+)?function\s+[A-Z]`)
)

func TestNoHookIsWrittenBelowAnEarlyReturn(t *testing.T) {
	t.Parallel()
	files, err := filepath.Glob("../../src/*.tsx")
	if err != nil || len(files) == 0 {
		t.Skip("no interface sources beside this build")
	}
	for _, file := range files {
		body, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		returned, returnedAt := false, 0
		// Line by line and unflattened: this one is about where a call sits
		// relative to a return, which is a fact about lines.
		for i, line := range strings.Split(string(body), "\n") {
			switch {
			case componentTop.MatchString(line):
				returned = false // a new component starts its own count
			case earlyReturn.MatchString(line):
				returned, returnedAt = true, i+1
			case returned && hookCall.MatchString(line) && !strings.HasPrefix(strings.TrimSpace(line), "//"):
				t.Errorf("%s:%d calls a hook below the early return at line %d. React counts hooks: the count will change between the render before the data arrives and the one after, which is error #310 and a blank page.",
					filepath.Base(file), i+1, returnedAt)
			}
		}
	}
}

// A preference the game obeys and the player cannot reach is worse than no
// preference at all: the behaviour is there, somebody's browser is deciding it,
// and the only way to change it is to open developer tools. `black-ledger-motion`
// gated the theatre — a modal that holds the screen for several seconds after a
// killing — and nothing in the interface had ever written it.
//
// The rule: every stored preference the interface reads must also be written
// somewhere the player can click. Keys the game only uses as scratch storage
// (a pending request being retried, how much of the paper has been read) are
// not preferences and are exempt by name.
func TestEverySettingTheGameObeysCanBeChanged(t *testing.T) {
	t.Parallel()
	files, err := filepath.Glob("../../src/*.tsx")
	if err != nil || len(files) == 0 {
		t.Skip("no interface sources beside this build")
	}
	notAPreference := map[string]bool{
		"black-ledger-pending":   true, // a command being replayed after a refresh
		"black-ledger-news-seen": true, // how far through the paper the reader is
	}
	read, written := map[string]string{}, map[string]bool{}
	key := regexp.MustCompile(`localStorage\.(getItem|setItem|removeItem)\('(black-ledger-[a-z-]+)'`)
	for _, file := range files {
		body, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range key.FindAllStringSubmatch(string(body), -1) {
			if m[1] == "getItem" {
				read[m[2]] = filepath.Base(file)
				continue
			}
			written[m[2]] = true
		}
	}
	if len(read) == 0 {
		t.Fatal("no stored preferences were found at all, which means this guard is not reading the sources")
	}
	for name, where := range read {
		if notAPreference[name] || written[name] {
			continue
		}
		t.Errorf("%s obeys %q and nothing in the interface can set it: the player cannot change a setting the game is using", where, name)
	}
}

// The guard above is a regular expression doing a parser's job, so its
// precision is worth testing directly: it has to keep catching the shape that
// broke the game and stop reporting the shape that cannot.
func TestTheHookGuardKnowsWhichReturnsMatter(t *testing.T) {
	t.Parallel()
	dangerous := []string{
		` if(!world)return <div className="loading">BLACK LEDGER</div>;`,
		`  if (!place) return null;`,
	}
	harmless := []string{
		`      if (!from || !to) return [];`,
		`    if (matchMedia('(prefers-reduced-motion: reduce)').matches) { setPaper(true); return }`,
		`        if (!ok) return errand{}, false`,
	}
	for _, line := range dangerous {
		if !earlyReturn.MatchString(line) {
			t.Errorf("the guard no longer sees a component's own early return: %q", line)
		}
	}
	for _, line := range harmless {
		if earlyReturn.MatchString(line) {
			t.Errorf("the guard reports a return that cannot change the hook count: %q", line)
		}
	}
}

// Every figure the core puts in the top bar is drawn with an icon looked up by
// its id, and an id with no path falls back to the city skyline — so a new stat
// silently wears the wrong picture. This is the cheapest check that the two
// sides still agree.
func TestEveryTopBarFigureHasAnIconOfItsOwn(t *testing.T) {
	t.Parallel()
	art := source(t, "src/art.ts")
	w := core.New(4)
	w.Confine(5, "a still in the back")
	for _, s := range w.Dashboard() {
		// Written `cash: 'M3 6h18…'` since the view was formatted, and `cash:'M…'`
		// before it. Flattening does not close a gap, so ask for the spelling
		// the file actually uses.
		if !holds(art, s.ID+": 'M") {
			t.Errorf("the top bar shows %q and nothing draws it, so it wears the city skyline", s.ID)
		}
	}
}

// Two rows of the same panel claimed the same thing about different numbers:
// "Working at 80%" beside "earns 100% of what it could", where the second was
// about how many regulars a place keeps and the first about everything. A
// player reading them together is reading a contradiction. Only one figure in
// this game answers "of what it could", and it is the one the core computes.
func TestOnlyOneFigureClaimsToBeWhatAPlaceCouldEarn(t *testing.T) {
	t.Parallel()
	body := source(t, "src/main.tsx")
	if n := strings.Count(string(body), "of what it could"); n != 0 {
		t.Errorf("the property panel claims %d times to say what a place could earn; that sentence belongs to the core's own note", n)
	}
}

// The result band reports how long a decision took. Sitting out a sentence the
// game itself calls "Do the 5 days" was reported as "120 hours" — true, and it
// tells the player nothing. This is a text guard only; that the band actually
// reads "5 days" was checked in a browser, because a regular expression cannot
// tell you what a player sees.
func TestTheResultBandCanCountInDays(t *testing.T) {
	t.Parallel()
	// The three things a band that counts in days has to do, rather than the
	// bare number 1440, which any arithmetic anywhere in the file satisfies.
	body := source(t, "src/Outcome.tsx")
	for _, part := range []string{"m >= 1440", "Math.floor(m / 1440)", "(m % 1440) / 60"} {
		if !holds(body, part) {
			t.Errorf("the result band cannot count in days: %q is not in it", part)
		}
	}
}

// The city without the map. Every address is reachable from a plain list of
// buttons behind the drawing, for somebody on a keyboard or without WebGL, and
// it was twenty-six names in whatever order the city held them: no distance, no
// sign of which were yours, nothing to choose on. The journeys in this city run
// from ten minutes to a hundred and twenty-five, so where a place is decides
// most of what going there costs, and that list is where it has to be legible.
func TestTheCityListSaysWhatEachJourneyCosts(t *testing.T) {
	t.Parallel()
	src := source(t, "src/CityIso.tsx")
	if !holds(src, "p.crossing.minutes") {
		t.Fatal("the address list names every place and says nothing about reaching any of them")
	}
	if !holds(src, "aria-label=\"City addresses\"") {
		t.Fatal("the list that is the city without the map is gone")
	}
	// Sorted, and by the journey rather than by anything else.
	if !holds(src, "a.crossing?.minutes ?? 0) - (b.crossing?.minutes ?? 0)") {
		t.Fatal("the addresses are not ordered by how long it takes to reach them")
	}
	// And what cannot be reached is not sorted into the middle by a journey
	// nobody can make.
	if !holds(src, "p.district > state.district ? 2 : 1") {
		t.Fatal("an address that is not open to the player is ordered as though it were")
	}
}
