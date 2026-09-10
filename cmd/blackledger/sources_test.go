package main

import (
	"os"
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
