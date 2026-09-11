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

// Every table in the game makes a noise except the one behind the poolhall. The
// casino floor has had a hum under it and a card on every deal since the tables
// took the screen; the back room, which takes the screen the same way, had
// silence — and silence at a card table reads as a thing that is broken.
func TestTheBackRoomIsNotSilent(t *testing.T) {
	t.Parallel()
	src := source(t, "src/BackRoomScene.tsx")
	if !holds(src, "roomTone(true)") || !holds(src, "return () => roomTone(false)") {
		t.Fatal("the back room has no room under it, or leaves it running after the player walks out")
	}
	if !holds(src, "playTable('card')") {
		t.Fatal("cards land on the table in silence")
	}
	if !holds(src, "playTable('chips'") {
		t.Fatal("money goes into the middle in silence")
	}
	// Off what the core sent, never off a clock of the interface's own: the
	// board is the count of cards the core has dealt.
	if !holds(src, "cards?.board?.length ?? 0") {
		t.Fatal("the cards are counted by something other than what the core dealt")
	}
	// And the first look is not a thing that just happened, the way sitting
	// down used to deal a card the moment the room opened.
	if !holds(src, "seenBoard.current === null") || !holds(src, "seenPot.current === null") {
		t.Fatal("walking into the room plays the hand that was already on the table")
	}
	// A new hand clears the board, and cards coming off the table are not cards
	// landing on it.
	if !holds(src, "board < seenBoard.current") {
		t.Fatal("the board being cleared for a new hand is heard as cards being dealt")
	}
	if !holds(source(t, "src/sound.ts"), "case 'chips':") {
		t.Fatal("nothing knows what chips sound like")
	}
}

// And the ledger draws it. The core can hold a fact the interface never puts on
// a page, which is the same as not holding it.
func TestTheLedgerDrawsWhatIsBehind(t *testing.T) {
	t.Parallel()
	src := source(t, "src/LedgerScreen.tsx")
	if !holds(src, "b.behind?.length") {
		t.Fatal("the ledger never asks whether anything is behind")
	}
	if !holds(src, "b.behind.map") {
		t.Fatal("the ledger knows something is behind and does not say which premises")
	}
	if !holds(src, "nobody will take the job until it is paid") {
		t.Fatal("a place nobody will work at reads the same as one that is a night late")
	}
	// Above the breakdown of the day's costs, because what is owed and what has
	// already gone wrong are not the same thing and the second one is the one
	// to act on.
	behind := strings.Index(src, "books-behind")
	lines := strings.Index(src, "books-lines")
	if behind < 0 || lines < 0 || behind > lines {
		t.Fatal("what is already unpaid is drawn below the breakdown of what is owed")
	}
	if !holds(source(t, "src/style.css"), ".books-behind li.shut") {
		t.Fatal("a place nobody will work at is drawn exactly like one a night behind")
	}
}

// The address list behind the map says which rooms have a game in them, the
// same way it says what each journey costs and which places are yours. The
// list is the city for somebody on a keyboard, so a fact that only reaches the
// drawing does not reach them.
func TestTheCityListSaysWhereTheGamesAre(t *testing.T) {
	t.Parallel()
	src := source(t, "src/CityIso.tsx")
	if !holds(src, "p.back_room ? ' · a game in the back' : ''") {
		t.Fatal("the address list says nothing about which rooms have a game behind them")
	}
}

// And the screen follows what the player sat down to rather than what the room
// holds. Picking it from the address is what put somebody who asked for the
// machines at Saint Agnes into a hand of cards.
func TestTheScreenFollowsWhatWasSatDownTo(t *testing.T) {
	t.Parallel()
	src := source(t, "src/main.tsx")
	if !holds(src, "world!.seated_to === 'back'") {
		t.Fatal("the screen is chosen by something other than what the player sat down to")
	}
	if holds(src, "l.id === world!.seated)?.back_room") {
		t.Fatal("the screen is still chosen by what the room holds")
	}
}
