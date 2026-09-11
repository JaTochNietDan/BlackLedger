package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// These suites read the view's own source and check that the two sides of the
// game still agree — that a screen is handed the work the core marks, that a
// figure in the top bar has a picture, that the room offers a way in. They are
// the cheapest guard there is against the interface quietly drifting from the
// core, and they cost nothing to run.
//
// They matched the source byte for byte, spacing included, which made every one
// of them a hostage to how the file happened to be laid out. Running a
// formatter over `src/` broke five of them at once while changing nothing the
// browser does: the bundle came out identical once you ignore where the spaces
// sit and how JSX splits a run of text into children.
//
// So they match on structure now. `source` reads a file with every run of
// whitespace collapsed to one space, and `holds` looks for a needle collapsed
// the same way, so "a={b} c={d}" and "a={b}\n  c={d}" are the same thing to a
// guard and neither is the thing being guarded.

var runsOfSpace = regexp.MustCompile(`\s+`)

// flat collapses every run of whitespace to a single space.
func flat(s string) string { return strings.TrimSpace(runsOfSpace.ReplaceAllString(s, " ")) }

// source reads one of the view's files, flattened. It skips the test rather
// than failing it where the interface is not beside this build.
func source(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile("../../" + path)
	if err != nil {
		t.Skip("no interface sources beside this build")
	}
	return flat(string(b))
}

// rawSource reads one of the view's files exactly as written. The stylesheet's
// guards count rules at the start of a line and the hook guard is about where a
// call sits relative to a return, and both of those are facts about lines.
func rawSource(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile("../../" + path)
	if err != nil {
		t.Skip("no interface sources beside this build")
	}
	return string(b)
}

// holds reports whether the flattened source contains the flattened needle.
func holds(src, needle string) bool { return strings.Contains(src, flat(needle)) }

// A guard over the guards.
//
// The brief's first fault shape is "a guard that cannot fail", and these text
// guards are the easiest place in the project to write one: a needle short
// enough or generic enough to appear in any file passes forever while the thing
// it was written about quietly goes away. Two were found by hand — `") : ("`
// matched every ternary in a file, and `"1440"` matched any arithmetic on a day
// — and finding them by hand is not a plan.
//
// So the needles are checked. Anything under this many characters is either a
// bare number, a fragment of punctuation, or a word common enough to appear
// somewhere by accident; a real needle is a phrase.
const ShortestNeedle = 10

func TestNoGuardIsWrittenWithANeedleThatCannotFail(t *testing.T) {
	t.Parallel()
	files, err := filepath.Glob("*_test.go")
	if err != nil || len(files) == 0 {
		t.Skip("no guards beside this build")
	}
	// Only what a guard demands is present. A short needle a guard demands is
	// ABSENT errs the safe way: it fails more often than it should rather than
	// less, and "the room must not contain the word backroom" is exactly as
	// strong as it sounds.
	needle := regexp.MustCompile(`!holds\([a-zA-Z]+, "((?:[^"\\]|\\.)*)"`)
	weak := []string{}
	for _, file := range files {
		body, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, found := range needle.FindAllStringSubmatch(string(body), -1) {
			if len(found[1]) < ShortestNeedle {
				weak = append(weak, fmt.Sprintf("%s: %q", filepath.Base(file), found[1]))
			}
		}
	}
	if len(weak) > 0 {
		t.Errorf("%d guard(s) ask for something short enough to be in any file, so they cannot fail: %v",
			len(weak), weak)
	}
}

// A rule that is written twice and quietly loses the argument.
//
// The stylesheet has grown by appending: a block gets rewritten, the new
// version goes on the end, and the old one stays where it was. Where both
// rules carry the same selector the later one simply wins, and the earlier
// declarations are dead — which is harmless until somebody reads the earlier
// block, believes it, and edits it.
//
// That is not hypothetical. It cost a tick twice on the same element. The
// drums did not spin because a `.drum` rule at the top of the slot machine
// block set `display:grid;place-items:center` and the newer rule further down
// never said otherwise, so the strip was centred instead of running; the fault
// was invisible because the newer block read correctly on its own.
//
// So the count is pinned. It does not have to be zero today — there are fifty
// of these and unpicking them is its own piece of work, to be done where the
// result can be looked at — but it must never go up, and the names are printed
// so the next person to open this file knows which blocks lie.
const shadowedDeclarations = 50

func TestNoRuleIsQuietlyOverriddenByALaterCopyOfItself(t *testing.T) {
	css := rawSource(t, "src/style.css")
	// Comments are blanked rather than cut so every offset still lines up with
	// the file as written.
	bare := regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllStringFunc(css, func(c string) string {
		return regexp.MustCompile(`[^\n]`).ReplaceAllString(c, " ")
	})
	type rule struct {
		selector string
		body     string
		line     int
	}
	rules := []rule{}
	depth, buf, selector, bodyAt, line, at := 0, "", "", 0, 1, 1
	for i, c := range bare {
		if c == '\n' {
			line++
		}
		switch c {
		case '{':
			if depth == 0 {
				selector, bodyAt, at = strings.TrimSpace(buf), i+1, line
			}
			depth++
			buf = ""
		case '}':
			depth--
			if depth == 0 {
				rules = append(rules, rule{selector, bare[bodyAt:i], at})
				buf = ""
			}
		default:
			buf += string(c)
		}
	}
	// Every rule a selector appears in, in the order they are written.
	where := map[string][]int{}
	for i, r := range rules {
		if strings.HasPrefix(r.selector, "@") {
			continue
		}
		for _, one := range strings.Split(r.selector, ",") {
			one = strings.TrimSpace(one)
			where[one] = append(where[one], i)
		}
	}
	declaration := regexp.MustCompile(`([a-z-]+)\s*:\s*([^;]*)`)
	declared := func(body string) map[string]string {
		out := map[string]string{}
		for _, m := range declaration.FindAllStringSubmatch(body, -1) {
			out[m[1]] = m[2]
		}
		return out
	}
	dead, said := 0, []string{}
	for _, idxs := range where {
		if len(idxs) < 2 {
			continue
		}
		for pos, idx := range idxs[:len(idxs)-1] {
			later := map[string]bool{}
			for _, j := range idxs[pos+1:] {
				for prop, val := range declared(rules[j].body) {
					// An earlier !important still wins, so the later copy is
					// not an override and the earlier one is not dead.
					if !strings.Contains(val, "!important") {
						later[prop] = true
					}
				}
			}
			for prop, val := range declared(rules[idx].body) {
				if later[prop] && !strings.Contains(val, "!important") {
					dead++
					said = append(said, fmt.Sprintf("    line %d  %s  %s", rules[idx].line, rules[idx].selector, prop))
				}
			}
		}
	}
	if dead > shadowedDeclarations {
		sort.Strings(said)
		t.Fatalf("%d declarations are overridden by a later copy of the same selector, up from %d:\n%s",
			dead, shadowedDeclarations, strings.Join(said, "\n"))
	}
	if dead < shadowedDeclarations {
		t.Fatalf("%d of these are left, down from %d — lower the pin in shadowedDeclarations so the ground that was won stays won",
			dead, shadowedDeclarations)
	}
}
