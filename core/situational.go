package core

import "fmt"

// Courier, mediation and collection exist in every city on every day. They are
// the reason the director could be handed a war and still write an errand:
// nothing in the operation set was about a conflict.
//
// These operations exist only because of something the simulation committed.
// They cannot be selected in a quiet city, which is what makes them worth
// telling: when one arrives, the world caused it.

// SituationalOperation is work that only exists under a condition, with the
// committed fact that justifies it.
type SituationalOperation struct {
	ID      string
	Effect  Effect
	Because string
}

// SituationalOperations reports the conflict-derived work the city currently
// supports, in a fixed order so selection is deterministic.
func (w *World) SituationalOperations() []SituationalOperation {
	out := []SituationalOperation{}

	// A war between two organizations means people who need moving safely and
	// messages nobody will carry openly.
	for _, c := range w.Conflicts {
		if c.State != "war" {
			continue
		}
		a, b := w.factionName(c.A), w.factionName(c.B)
		out = append(out, SituationalOperation{
			ID:      "escort",
			Effect:  Effect{Reward: 165, Respect: 6, Heat: 4, Minutes: 90},
			Because: fmt.Sprintf("%s and %s are at war, so moving anything of value across the city needs somebody watching it.", a, b),
		})
		out = append(out, SituationalOperation{
			ID:      "warning",
			Effect:  Effect{Reward: 110, Respect: 7, Heat: 8, Minutes: 60},
			Because: fmt.Sprintf("%s and %s are at war, and a message delivered in person carries further than one that is sent.", a, b),
		})
		break
	}

	// Ground that changed hands recently leaves things behind that somebody
	// wants back.
	if place, holder, ok := w.recentSeizure(); ok {
		out = append(out, SituationalOperation{
			ID:      "recovery",
			Effect:  Effect{Reward: 150, Respect: 5, Heat: 6, Minutes: 75},
			Because: fmt.Sprintf("%s changed hands recently and is now held by %s. Not everything inside it belonged to the people who lost it.", place, holder),
		})
	}

	// A death at the top of an organization leaves arrangements that were only
	// ever held together by the person who is gone.
	if name, organization, ok := w.recentLeadershipChange(); ok {
		out = append(out, SituationalOperation{
			ID:      "settlement",
			Effect:  Effect{Reward: 140, Respect: 6, Heat: 3, Minutes: 75},
			Because: fmt.Sprintf("%s now leads %s. Arrangements made with the person before them are not settled.", name, organization),
		})
	}
	return out
}

// SituationalEffects is the reward table for the conflict-derived operations,
// merged into the ordinary catalog when a proposal is validated.
func SituationalEffects() map[string]Effect {
	return map[string]Effect{
		"escort":     {Reward: 165, Respect: 6, Heat: 4, Minutes: 90},
		"warning":    {Reward: 110, Respect: 7, Heat: 8, Minutes: 60},
		"recovery":   {Reward: 150, Respect: 5, Heat: 6, Minutes: 75},
		"settlement": {Reward: 140, Respect: 6, Heat: 3, Minutes: 75},
	}
}

// recentSeizure finds ground that changed hands within the last few days, from
// the city's own news rather than from any hidden record.
func (w *World) recentSeizure() (string, string, bool) {
	for i := len(w.News) - 1; i >= 0; i-- {
		s := w.News[i]
		if s.Life != w.Life || s.Kind != "seizure" || w.Minute-s.Minute > 4320 {
			continue
		}
		for _, l := range Locations {
			if prop := w.Properties[l.ID]; prop != nil && upper(l.Name)+" CHANGES HANDS" == s.Headline {
				return l.Name, w.HolderName(l.ID), true
			}
		}
	}
	return "", "", false
}

// recentLeadershipChange finds an organization whose head changed within the
// last few days and is still standing.
func (w *World) recentLeadershipChange() (string, string, bool) {
	for i := len(w.History) - 1; i >= 0; i-- {
		r := w.History[i]
		if r.Life != w.Life || r.Kind != "politics" || w.Minute-r.Minute > 4320 {
			continue
		}
		for _, f := range w.Factions {
			if r.Title == f.Leader+" takes over "+f.Name && len(w.FamilyHoldings(f.ID)) > 0 {
				return f.Leader, f.Name, true
			}
		}
	}
	return "", "", false
}
