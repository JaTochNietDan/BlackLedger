package core

import (
	"regexp"
	"strings"
)

// Letting the director write the city page, without letting it lie.
//
// The page is composed from facts — the weather, what a crate is fetching, how
// many people live here — and reads like it: correct, and a little flat. A
// language model can set the same fact in the register of a 1953 city paper,
// which is worth having.
//
// What it must not do is add anything. A model asked to make a paragraph more
// colourful will happily invent a councilman, a street, a number of arrests, a
// name. In a game whose whole premise is that the paper only prints what the
// city could see, an invented fact is not a blemish, it is a lie the player has
// no way to detect — they would read that eleven people were arrested and
// believe it.
//
// So the rewrite is checked against the original before it is allowed in, and
// the checks are about *addition* rather than about quality. Nothing here can
// tell whether the prose is any good. It can tell whether the prose is about
// the same facts.

// number finds anything the reader would take as a quantity.
var number = regexp.MustCompile(`\d+`)

// properName finds a capitalised word that is not at the start of a sentence —
// which in newspaper copy is very close to "a name the writer has introduced".
var properName = regexp.MustCompile(`([a-z,;]\s+)([A-Z][a-z]{2,})`)

// PolishLimit is how long a brief may be. A model asked for two sentences and
// given no limit will write eight.
const PolishLimit = 620

// AcceptPolish decides whether a rewritten brief may replace the original.
//
// It returns the text to use and whether the rewrite was taken. When it refuses
// it returns the original, so a caller that ignores the boolean still cannot
// print an invention.
func AcceptPolish(original, rewritten string) (string, bool) {
	clean := strings.TrimSpace(rewritten)
	// Models like to explain themselves. Anything before the copy is not copy.
	if cut := strings.LastIndex(clean, "\n\n"); cut > 0 && cut < len(clean)/2 {
		clean = strings.TrimSpace(clean[cut:])
	}
	clean = strings.Trim(clean, "\"")

	if len(clean) < 40 || len(clean) > PolishLimit {
		return original, false
	}
	// A rewrite that is barely a rewrite is not worth the risk of taking it.
	if clean == strings.TrimSpace(original) {
		return original, false
	}

	// Every number in the new copy has to have been in the old copy. This is
	// the one that matters: prices, counts and dates are what a reader will
	// take as fact and act on.
	had := map[string]bool{}
	for _, n := range number.FindAllString(original, -1) {
		had[n] = true
	}
	for _, n := range number.FindAllString(clean, -1) {
		if !had[n] {
			return original, false
		}
	}

	// And every name it uses has to have been there too. A model given a
	// paragraph about the weather will put a mayor in it.
	knew := map[string]bool{}
	for _, m := range properName.FindAllStringSubmatch(original, -1) {
		knew[m[2]] = true
	}
	for _, word := range strings.Fields(original) {
		knew[strings.Trim(word, ".,;:\"'")] = true
	}
	for _, m := range properName.FindAllStringSubmatch(clean, -1) {
		if !knew[m[2]] && !commonCapital[m[2]] {
			return original, false
		}
	}

	// The paper never addresses the reader and never knows who arranged
	// anything. A model writing "you" has stopped being a newspaper.
	lower := strings.ToLower(clean)
	for _, leak := range []string{" you ", "you ", " your ", "the player", "as an ai", "i have"} {
		if strings.HasPrefix(lower, strings.TrimSpace(leak)) || strings.Contains(lower, leak) {
			return original, false
		}
	}
	return clean, true
}

// commonCapital is the handful of capitalised words that are part of writing
// about this city rather than facts about it, so they may appear in a rewrite
// that did not have them in front of it.
var commonCapital = map[string]bool{
	"Bellwether": true, "Monday": true, "Tuesday": true, "Wednesday": true,
	"Thursday": true, "Friday": true, "Saturday": true, "Sunday": true,
	"January": true, "February": true, "March": true, "April": true, "May": true,
	"June": true, "July": true, "August": true, "September": true,
	"October": true, "November": true, "December": true, "The": true,
}

// NextPolish is the newest brief the director has not seen yet, if there is
// one. Only today's: yesterday's paper has gone out.
func (w *World) NextPolish() (Story, bool) {
	today := w.Minute / 1440
	for i := len(w.News) - 1; i >= 0; i-- {
		s := w.News[i]
		if s.Minute/1440 != today || s.Life != w.Life {
			continue
		}
		if s.Kind == "civic" && !s.Polished {
			return s, true
		}
	}
	return Story{}, false
}

// SetPolish writes an accepted rewrite back onto the story it belongs to, and
// marks the story seen either way. Matched by ID: the archive may have moved
// underneath the request while the model was thinking.
func (w *World) SetPolish(id, body string, took bool) bool {
	for i := range w.News {
		if w.News[i].ID != id {
			continue
		}
		w.News[i].Polished = true
		if took {
			w.News[i].Body = body
		}
		return true
	}
	return false
}
