package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
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
