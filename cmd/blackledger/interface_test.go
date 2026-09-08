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
	hookCall     = regexp.MustCompile(`\buse(State|Effect|Memo|Callback|Ref|Reducer|Context|LayoutEffect)\s*\(`)
	earlyReturn  = regexp.MustCompile(`^\s*if\s*\(.*\)\s*return\s`)
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
