package main

import (
	"blackledger/core"
	"fmt"
	"strings"
	"unicode"
)

func normalizedOffer(s string) string {
	return strings.Join(strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) }), " ")
}

// Five-word overlap catches recycled passages with substituted names/locations.
// Short dialogue and shared genre vocabulary alone must not trigger rejection.
func recycledPassage(a, b string) bool {
	shingles := func(s string) map[string]bool {
		words := strings.Fields(normalizedOffer(s))
		out := map[string]bool{}
		for i := 0; i+5 <= len(words); i++ {
			out[strings.Join(words[i:i+5], " ")] = true
		}
		return out
	}
	x, y := shingles(a), shingles(b)
	if len(x) < 20 || len(y) < 20 {
		return false
	}
	common := 0
	for phrase := range x {
		if y[phrase] {
			common++
		}
	}
	denominator := len(x)
	if len(y) < denominator {
		denominator = len(y)
	}
	return float64(common)/float64(denominator) >= .65
}
func repeatedProposal(w *core.World, p core.Proposal) error {
	check := func(title, body string) error {
		if normalizedOffer(title) == normalizedOffer(p.Title) || normalizedOffer(body) == normalizedOffer(p.Body) || recycledPassage(body, p.Body) {
			return fmt.Errorf("proposal repeats a recent offer %q; write a distinct task and title rather than swapping names or locations, preserving any required contact connection", title)
		}
		return nil
	}
	for _, m := range w.Arrangements {
		if m.Life == w.Life {
			if err := check(m.Title, m.Offer); err != nil {
				return err
			}
		}
	}
	for _, o := range w.Offers {
		if o.Event != nil {
			if err := check(o.Event.Title, o.Event.Body); err != nil {
				return err
			}
		}
	}
	if w.Event != nil && w.Event.Kind == "proposal" {
		return check(w.Event.Title, w.Event.Body)
	}
	return nil
}
