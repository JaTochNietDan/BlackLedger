package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
