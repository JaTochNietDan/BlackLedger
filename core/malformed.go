package core

import "regexp"

// What a sentence this city says must never look like.
//
// This started inside a test that read every card in every room. That reader
// only ever saw the cards, and a card is a small part of what the game says: the
// log and the paper carry the rest, and they are written while something is
// happening rather than while somebody is deciding. Three faults in three nights
// were sentences of that second kind — "0 crates of arms out through the front
// door in daylight" among them — so the reader lives here now and anything that
// can collect text can use it.
var malformed = []struct {
	name string
	re   *regexp.Regexp
}{
	{"a format verb that was never filled", regexp.MustCompile(`%!|%\[|%s|%d|%v`)},
	// The pronouns below are a pattern the reader matches against sentences the
	// game has already written, not words the game says. A city that never
	// assumes a gender still has to notice when something else did.
	{"an article doubled onto a name", regexp.MustCompile(`\b(your|their|his|her|its) an? `)}, // not prose

	{"two spaces", regexp.MustCompile(`  `)},
	{"a space before a full stop", regexp.MustCompile(` \.`)},
	{"an empty figure", regexp.MustCompile(`\$-`)},
	{"a sentence that starts in the middle", regexp.MustCompile(`^[a-z]`)},

	// Counting. A haul of nothing is not a haul, and one of a thing is not
	// several of it: the precinct used to sell a man back "for the 1 days still
	// on them", and a search of an empty room was reported as "0 crates of arms
	// out through the front door". Both were found by eye, one at a time.
	{"a count of none written as a quantity", regexp.MustCompile(`(^|[^.\d])0 [a-z]+s\b`)},
	{"one of something written as several", regexp.MustCompile(`(^|[^.\d])1 [a-z]+s\b`)},
}

// Malformed names what is wrong with a sentence, or returns "" if nothing is.
func Malformed(text string) string {
	for _, d := range malformed {
		if d.re.MatchString(text) {
			return d.name
		}
	}
	return ""
}

// MalformedKinds is every fault the reader can see, so a sweep that finds
// nothing can be told apart from a reader that cannot see anything.
func MalformedKinds() []string {
	out := []string{}
	for _, d := range malformed {
		out = append(out, d.name)
	}
	return out
}
