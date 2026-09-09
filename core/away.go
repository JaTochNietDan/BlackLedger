package core

import (
	"fmt"
	"strings"
)

// Measured before this was written: over forty campaigns, four days away from
// Bellwether left the city different on the player's return every single time —
// the paper had printed something in all forty. Holdings changed hands, people
// died, families rose and fell.
//
// The player was told none of it. The record filed on their return read "4 days
// gone. The city did not wait", which is a claim with nothing in it. Everything
// that happened was sitting in the Herald, four screens away, filed under days
// the player had no reason to go back and read.
//
// Coming home is the moment that information is worth most. The city already
// knows what it did; this is it saying so.

// missedSince is what the paper carried while the player was elsewhere, biggest
// story first. Newsworthiness is the paper's own ordering, so what leads here is
// what led the front page.
func (w *World) missedSince(minute int) []Story {
	out := []Story{}
	for _, s := range w.News {
		if s.Minute >= minute && s.Life == w.Life {
			out = append(out, s)
		}
	}
	for a := range out {
		for b := a + 1; b < len(out); b++ {
			if Newsworthiness(out[b].Kind) > Newsworthiness(out[a].Kind) {
				out[a], out[b] = out[b], out[a]
			}
		}
	}
	return out
}

// MissedHeadlines is the number of stories worth carrying back, and how many
// there were in all.
const MissedHeadlines = 3

// WhatYouMissed is the sentence a returning player reads. The headlines are the
// paper's own, printed exactly as it printed them: sentence-casing them would
// lower-case the names, and paraphrasing them would be this file inventing
// news. Civic filler is left out — somebody who has been away for four days
// does not need to be told it rained.
func (w *World) WhatYouMissed(since int) string {
	// Two days of the same headline is one thing to tell somebody who has been
	// away. Read out of a real save: "POLICE PRESSURE ON RUSSO OUTFIT · POLICE
	// PRESSURE ON RUSSO OUTFIT" — the paper collapses repeats within a day and
	// this spans several.
	worth := []Story{}
	told := map[string]bool{}
	for _, s := range w.missedSince(since) {
		if Newsworthiness(s.Kind) <= Newsworthiness("civic") || told[s.Headline] {
			continue
		}
		told[s.Headline] = true
		worth = append(worth, s)
	}
	if len(worth) == 0 {
		return "Nothing in the paper you needed to be here for."
	}
	shown := worth
	if len(shown) > MissedHeadlines {
		shown = shown[:MissedHeadlines]
	}
	heads := make([]string, 0, len(shown))
	for _, s := range shown {
		heads = append(heads, s.Headline)
	}
	line := "While you were gone the paper carried: " + strings.Join(heads, " · ") + "."
	if rest := len(worth) - len(shown); rest > 0 {
		line += fmt.Sprintf(" And %d other %s.", rest, plainly(rest, "story", "stories"))
	}
	return line
}

// plainly picks a word to go with a count that is printed beside it.
func plainly(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
