package main

import (
	"blackledger/core"
	"fmt"
	"regexp"
	"strings"
)

// These lexical checks catch explicit unsupplied terms, not every possible
// semantic promise. The authoritative choice details own price and duration.
// Do not rewrite the prose: ask for the one bounded correction instead.
const quantity = `(?:\d[\d,.]*|a|an|one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve|thirteen|fourteen|fifteen|sixteen|seventeen|eighteen|nineteen|twenty|thirty|forty|fifty|sixty|seventy|eighty|ninety|hundred|thousand|million)`

var spokenMoney = regexp.MustCompile(`(?i)(?:[\p{Sc}]\s*` + quantity + `|\b(?:USD|GBP|EUR)\s*\d|\b` + quantity + `(?:[ -]+` + quantity + `)*[ -]+(?:dollars?|bucks?|cents?|pounds?|euros?|grand|USD|GBP|EUR)\b)`)
var spokenDuration = regexp.MustCompile(`(?i)\b` + quantity + `(?:[ -]+` + quantity + `)*[ -]+(?:minutes?|hours?|days?|weeks?|months?)\b`)
var spokenDeadline = regexp.MustCompile(`(?i)\b(?:before|by|until|at)\s+(?:(?:the|next|this)\s+)?(?:noon|midnight|dawn|dusk|sunrise|sunset|sundown|nightfall|lunch|lunchtime|closing(?:\s+time)?|tomorrow|tonight)\b|\b(?:before|by|until)\s+(?:the\s+)?(?:market|bar|laundry|casino|club|shop|day|week|month|shift|night|business)\s+(?:closes?|opens?|ends?)\b|\b(?:before|by|until|at)\s+\d{1,2}(?::\d{2}|\s*(?:a\.?m\.?|p\.?m\.?|o'clock))\b|\b(?:before|by|until|at)\s+` + quantity + `\s+o'clock\b|\b(?:before|by|until)\s+(?:the\s+)?(?:end|close)\s+of\s+(?:(?:the|this)\s+)?(?:day|week|month|shift|night|business)\b`)

func validateSpokenTerms(p core.Proposal) error {
	fields := []struct{ name, text string }{{"title", p.Title}, {"body", p.Body}}
	for _, approach := range p.Approaches {
		fields = append(fields, struct{ name, text string }{"approach label", approach.Label})
	}
	for _, field := range fields {
		text := strings.ReplaceAll(field.text, "’", "'")
		for _, check := range []struct {
			name    string
			pattern *regexp.Regexp
		}{
			{"monetary amount", spokenMoney}, {"duration", spokenDuration}, {"deadline", spokenDeadline},
		} {
			if match := check.pattern.FindString(text); match != "" {
				return fmt.Errorf("%s includes an unsupported %s %q; omit amounts, durations and deadlines from spoken terms and labels, because authoritative choice details supply them. Keep the task and its subject", field.name, check.name, match)
			}
		}
	}
	return nil
}
