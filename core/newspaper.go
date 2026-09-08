package core

import "fmt"

// The Bellwether Herald prints what the city can see. It is assembled at the
// moment events are committed, not summarised afterwards, so a headline always
// refers to something that actually happened.
//
// The paper knows only public things. It never prints a plot that has not
// happened yet, never names who paid for a killing, and reports the player's
// own crimes the way a newspaper would: as an event with no suspect. Reading it
// is free and advances nothing.

type Story struct {
	ID       string `json:"id"`
	Minute   int    `json:"minute"`
	Life     int    `json:"life"`
	Headline string `json:"headline"`
	Body     string `json:"body"`
	Kind     string `json:"kind"`
}

// newsCapacity bounds the archive. A campaign runs for weeks of game time and
// the paper is a recent record, not a complete one.
const newsCapacity = 60

// Report files a story. Callers pass what the city could observe, which is why
// the wording differs from the private record written to the player's history.
func (w *World) Report(kind, headline, body string) {
	w.News = append(w.News, Story{
		ID: ID(), Minute: w.Minute, Life: w.Life,
		Headline: headline, Body: body, Kind: kind,
	})
	if len(w.News) > newsCapacity {
		w.News = w.News[len(w.News)-newsCapacity:]
	}
}

// Edition is the paper as it stands, most recent first, with the day each story
// ran. Only the current life's stories are printed: a new person picks up a
// paper that has been running the whole time, but the archive they can read
// starts when they arrive.
func (w *World) Edition() []map[string]any {
	out := []map[string]any{}
	for i := len(w.News) - 1; i >= 0; i-- {
		s := w.News[i]
		if s.Life != w.Life {
			continue
		}
		out = append(out, map[string]any{
			"id": s.ID, "headline": s.Headline, "body": s.Body,
			"kind": s.Kind, "minute": s.Minute,
			"day": s.Minute/1440 + 1, "time": fmt.Sprintf("%02d:%02d", s.Minute%1440/60, s.Minute%60),
		})
	}
	return out
}

// unattributed is how the paper describes a crime nobody has been charged with.
// The player reads their own work here with no name attached, which is exactly
// what the rest of the city sees.
func (w *World) unattributed(place string, what string) string {
	return fmt.Sprintf("%s Police have made no arrest and are appealing for anyone who saw the incident at %s.", what, place)
}
