package main

import (
	"blackledger/core"
	"fmt"
	"regexp"
	"strings"
)

// Permanent death is the point of this game: a new person arrives owning nothing
// and knowing nobody, while the city keeps the dead protagonist's record. Prose
// that greets the player as a returning acquaintance invents a shared past the
// world never recorded, and it reads worst exactly where it matters most, in the
// first scene of a new life.
//
// Any earlier arrangement with the same speaker in the current life earns this
// language, whether it was completed or declined; being turned down is still a
// shared past. A previous protagonist's record never earns it, however
// faithfully the model attributes that record elsewhere.
func spokenAcquaintance(w *core.World, speaker string) bool {
	for _, m := range w.Arrangements {
		if m.Life == w.Life && m.Speaker == speaker {
			return true
		}
	}
	return false
}

// These are second-person claims of prior dealings, not every warm phrase. The
// authoritative saved arrangements own what has actually happened between these
// two people. Do not rewrite the prose: ask for the one bounded correction.
var unearnedFamiliarity = []struct {
	name    string
	pattern *regexp.Regexp
}{
	{"returning greeting", regexp.MustCompile(`(?i)\b(?:welcome\s+back|you(?:'|’)?re\s+back|you\s+are\s+back|good\s+to\s+(?:see|have)\s+you\s+(?:back|again)|see\s+you\s+again)\b`)},
	{"summons back to a place", regexp.MustCompile(`(?i)\b(?:need|want|bring|get)\s+you\s+back\b|\byou\s+back\s+(?:at|in|on|to)\b`)},
	{"appeal to a previous occasion", regexp.MustCompile(`(?i)\b(?:like|as)\s+(?:last\s+time|before|always|usual)\b|\bsame\s+as\s+(?:last\s+time|before|always)\b|\blast\s+time\s+you\b|\bthe\s+way\s+you\s+(?:always|usually)\b`)},
	{"assumed routine", regexp.MustCompile(`(?i)\byou\s+know\s+the\s+drill\b|\bour\s+usual\b|\bthe\s+usual\s+arrangement\b|\bas\s+usual\b|\bas\s+always\b`)},
	{"claimed track record", regexp.MustCompile(`(?i)\byou(?:'|’)?ve\s+done\s+this\s+before\b|\byou\s+did\s+(?:this|it)\s+before\b|\byou(?:'|’)?ve\s+(?:handled|worked)\b`)},
}

func validateEarnedFamiliarity(w *core.World, p core.Proposal) error {
	if spokenAcquaintance(w, p.Speaker) {
		return nil
	}
	fields := []struct{ name, text string }{{"title", p.Title}, {"body", p.Body}}
	for _, approach := range p.Approaches {
		fields = append(fields, struct{ name, text string }{"approach label", approach.Label})
	}
	for _, field := range fields {
		text := strings.ReplaceAll(field.text, "’", "'")
		for _, check := range unearnedFamiliarity {
			if match := check.pattern.FindString(text); match != "" {
				return fmt.Errorf("%s claims prior dealings with this person through a %s %q, but %s has no earlier arrangement with this speaker. Address them as someone they are meeting for the first time and describe the task on its own terms", field.name, check.name, match, w.Player.Name)
			}
		}
	}
	return nil
}
