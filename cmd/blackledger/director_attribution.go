package main

import (
	"blackledger/core"
	"fmt"
	"regexp"
	"strings"
)

// The people and premises in a job belong to whoever the rules say they belong
// to. Observed on qwen3.5:35b-a3b: a mechanically neutral mediation at the
// independently owned Saint Agnes described "two Bellandi staff", which invites
// the player to expect Bellandi standing to move when no goodwill is at stake.
//
// Mentioning a family remains legal; this rejects only attributing the job's own
// people or premises to a family that neither receives the credit nor owns the
// place the work happens in.
var attributedParty = `(?:staff|men|crew|people|workers?|employees?|soldiers?|associates?|guys|boys)`

// Organizations formed during play are named things like "the Falcone Crew", so
// the distinctive part is neither reliably the first word nor the last. Articles
// and the words every organization shares carry no identity, and matching on
// them would reject ordinary lines like "the men are arguing".
var genericOrganizationWords = map[string]bool{
	"the": true, "a": true, "an": true, "of": true,
	"family": true, "outfit": true, "crew": true, "combine": true,
	"syndicate": true, "company": true, "brothers": true, "sons": true,
}

func factionWords(f core.Faction) []string {
	words := []string{}
	for _, part := range strings.Fields(f.Name) {
		if !genericOrganizationWords[strings.ToLower(part)] {
			words = append(words, part)
		}
	}
	// A leader stands for their family, by either name: "Russo's men" and
	// "Vittorio's people" both hand the job's people to the Bellandi.
	words = append(words, strings.Fields(f.Leader)...)
	seen, out := map[string]bool{}, []string{}
	for _, w := range words {
		if w != "" && !seen[strings.ToLower(w)] {
			seen[strings.ToLower(w)] = true
			out = append(out, w)
		}
	}
	return out
}

func validateFactionAttribution(w *core.World, p core.Proposal) error {
	location := p.Location
	if location == "" {
		location = w.Player.Location
	}
	owner := ""
	if property := w.Properties[location]; property != nil {
		owner = property.Owner
	}
	body := strings.ReplaceAll(p.Body, "’", "'")
	for _, f := range w.Factions {
		if f.ID == p.Beneficiary || f.ID == owner {
			continue
		}
		for _, word := range factionWords(f) {
			pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(word) + `(?:'s)?\s+` + attributedParty + `\b`)
			if match := pattern.FindString(body); match != "" {
				place, _ := core.PlaceByID(location)
				return fmt.Errorf("body assigns this job's people to a family through %q, but the work is %s and %s is owned by %q. Describe the people involved without giving them a family, or make that family the beneficiary", match, neutralOrFor(w, p.Beneficiary), place.Name, owner)
			}
		}
	}
	return nil
}

func neutralOrFor(w *core.World, beneficiary string) string {
	if beneficiary == "" {
		return "neutral, earning no family standing"
	}
	for _, f := range w.Factions {
		if f.ID == beneficiary {
			return "credited to " + f.Name
		}
	}
	return "credited to " + beneficiary
}
