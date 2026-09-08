package main

import (
	"blackledger/core"
	"fmt"
	"regexp"
	"strings"
)

// Characters speak inside the city; they do not describe the cast list. The
// brief is direction for the writer, not lines for the speaker, and a model
// that dramatizes its own constraints produces dialogue like "Neither is you".
var outOfWorldVoice = []struct {
	name    string
	pattern *regexp.Regexp
}{
	{"reference to the player as a game role", regexp.MustCompile(`(?i)\bthe\s+player\b|\bthe\s+protagonist\b|\bplayer\s+character\b`)},
	{"reference to the cast", regexp.MustCompile(`(?i)\b(?:new\s+)?named\s+character\b|\bNPCs?\b|\bnon-?player\b`)},
	{"aside excluding the listener from the cast", regexp.MustCompile(`(?i)\bneither\s+(?:one\s+)?(?:is|of\s+(?:them|these|those|which)\s+is)\s+you\b|\bnone\s+of\s+(?:them|these|those)\s+(?:is|are)\s+you\b|\bneither\s+is\s+the\s+player\b`)},
	{"reference to the brief itself", regexp.MustCompile(`(?i)\bthis\s+(?:brief|proposal|schema|prompt)\b|\bapproach\s+label\b|\brequired\s+opening\b`)},
}

func validateInWorldVoice(p core.Proposal) error {
	fields := []struct{ name, text string }{{"title", p.Title}, {"body", p.Body}, {"outcome", p.Outcome}}
	for _, approach := range p.Approaches {
		fields = append(fields, struct{ name, text string }{"approach label", approach.Label})
	}
	for _, field := range fields {
		text := strings.ReplaceAll(field.text, "’", "'")
		for _, check := range outOfWorldVoice {
			if match := check.pattern.FindString(text); match != "" {
				return fmt.Errorf("%s breaks the scene with a %s %q; the brief is direction for you, not lines for the speaker. Say only what this character would say aloud in the city", field.name, check.name, match)
			}
		}
	}
	return nil
}

// The player runs their own operation and never leads one of the established
// families. Dialogue that hands them a rival's title, usually by attaching the
// speaker's own rank to the listener's name, reverses who is talking to whom.
var claimedFactionRank = regexp.MustCompile(`(?i),?\s+as\s+(?:the\s+)?(?:leader|head|boss|don|underboss)\s+of\b`)

func validatePlayerRole(w *core.World, p core.Proposal) error {
	text := strings.ReplaceAll(p.Body, "’", "'")
	subjects := []string{"you", strings.ToLower(w.Player.Name)}
	if first := strings.Fields(w.Player.Name); len(first) > 0 {
		subjects = append(subjects, strings.ToLower(first[0]))
	}
	lower := strings.ToLower(text)
	for _, subject := range subjects {
		if subject == "" {
			continue
		}
		for _, index := range regexp.MustCompile(`(?i)\b`+regexp.QuoteMeta(subject)+`\b`).FindAllStringIndex(lower, -1) {
			rest := lower[index[1]:]
			if location := claimedFactionRank.FindStringIndex(rest); location != nil && location[0] == 0 {
				return fmt.Errorf("body gives %s a family rank they do not hold (%q). The player leads no established family; if you mean your own rank, say it about yourself and address them by name", w.Player.Name, strings.TrimSpace(text[index[0]:index[1]+location[1]]))
			}
		}
	}
	return nil
}
