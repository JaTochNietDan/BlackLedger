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

	// A business short of what it runs on needs somebody to fetch it, and a
	// business in trouble needs somebody who can deal with that kind of thing.
	if place, need, ok := w.businessNeed(); ok {
		out = append(out, SituationalOperation{
			ID:      "supply",
			Effect:  Effect{Reward: 95, Respect: 4, Heat: 3, Minutes: 60},
			Because: fmt.Sprintf("%s is %s, and somebody has to go and get what it needs.", place, need),
		})
	}

	// Stock that has to move is the most ordinary reason in this city for
	// somebody to be asked to carry something.
	if held := w.Carrying(); held >= 10 {
		out = append(out, SituationalOperation{
			ID:      "distribution",
			Effect:  Effect{Reward: 175, Respect: 5, Heat: 7, Minutes: 75},
			Because: fmt.Sprintf("There are %d units of stock sitting where they should not be sitting, and every day they sit there is a day somebody could find them.", held),
		})
	}

	// A room full of crates and a war outside is a situation nobody has to be
	// told about twice.
	if _, ok := w.TheArmoury(); ok && w.Stocked() > 0 {
		if buyers := w.buyers(); len(buyers) > 0 {
			out = append(out, SituationalOperation{
				ID:     "consignment",
				Effect: Effect{Reward: 240, Respect: 4, Heat: 9, Minutes: 90},
				Because: fmt.Sprintf("%s is fighting and you are sitting on %d crates. They would rather not send anybody to a laundry in daylight.",
					buyers[0].Name, w.Stocked()),
			})
		}
	}

	// Two people the player knows, one of whom has been carrying something
	// against the other for weeks. Somebody would rather it did not go the way
	// it is going.
	if holder, target, because, ok := w.liveGrievance(); ok {
		out = append(out, SituationalOperation{
			ID:      "grievance",
			Effect:  Effect{Reward: 130, Respect: 7, Heat: 2, Minutes: 75},
			Because: fmt.Sprintf("%s has not forgiven %s for %s, and it is far enough along that somebody is going to get hurt over it.", holder, target, because),
		})
	}

	// Work the player has already promised somebody, which is the most
	// pressing thing in their life and the easiest thing for a contact to
	// mention.
	if live := w.Live(); len(live) > 0 {
		c := live[0]
		out = append(out, SituationalOperation{
			ID:     "obligation",
			Effect: Effect{Reward: 150, Respect: 3, Heat: 4, Minutes: 60},
			Because: fmt.Sprintf("%s asked you for something on behalf of %s and the time for it is running out. %s",
				c.GiverName, w.factionName(c.PatronID), w.Progress(c)),
		})
	}

	// An organization that has stopped talking about the player and started
	// deciding about them.
	if f, ok := w.dangerousEnemy(); ok {
		out = append(out, SituationalOperation{
			ID:      "warning_off",
			Effect:  Effect{Reward: 175, Respect: 6, Heat: 5, Minutes: 75},
			Because: fmt.Sprintf("%s %s stopped complaining about you, which is worse than complaining. Somebody who talks to them thinks it is worth one conversation.", Leads(f.Name), Agree(f.Name, "has", "have")),
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
		"escort":       {Reward: 165, Respect: 6, Heat: 4, Minutes: 90},
		"warning":      {Reward: 110, Respect: 7, Heat: 8, Minutes: 60},
		"recovery":     {Reward: 150, Respect: 5, Heat: 6, Minutes: 75},
		"settlement":   {Reward: 140, Respect: 6, Heat: 3, Minutes: 75},
		"supply":       {Reward: 95, Respect: 4, Heat: 3, Minutes: 60},
		"distribution": {Reward: 175, Respect: 5, Heat: 7, Minutes: 75},
		"consignment":  {Reward: 240, Respect: 4, Heat: 9, Minutes: 90},
		"grievance":    {Reward: 130, Respect: 7, Heat: 2, Minutes: 75},
		"obligation":   {Reward: 150, Respect: 3, Heat: 4, Minutes: 60},
		"warning_off":  {Reward: 175, Respect: 6, Heat: 5, Minutes: 75},
	}
}

// businessNeed finds a business of the player's that is short of something,
// naming the place and what is wrong with it in words a character could say.
func (w *World) businessNeed() (string, string, bool) {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		trade, running := TradeOf(l.ID)
		if !running || prop == nil || !w.Own(l.ID) {
			continue
		}
		switch {
		case prop.Supply == 0:
			return l.Name, "out of " + trade.Supplies, true
		case prop.Trouble:
			return l.Name, "not running properly", true
		case prop.Staff < trade.Hands:
			return l.Name, "short-handed", true
		}
	}
	return "", "", false
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

// liveGrievance finds a quarrel between two people the player has reason to
// know about, which is what makes it something a contact would mention.
func (w *World) liveGrievance() (string, string, string, bool) {
	if w.Reach() < 2 {
		return "", "", "", false
	}
	for _, g := range w.Grudges {
		if g.Weight < GrudgeActs-15 {
			continue
		}
		holder, target := w.NPC(g.Holder), w.NPC(g.Against)
		if holder == nil || target == nil || holder.Dead || target.Dead {
			continue
		}
		return holder.Name, target.Name, g.Because, true
	}
	return "", "", "", false
}

// dangerousEnemy is an organization whose standing with the player has passed
// the point where they stop sending messages.
func (w *World) dangerousEnemy() (*Faction, bool) {
	for i := range w.Factions {
		if f := &w.Factions[i]; f.Goodwill <= -55 {
			return f, true
		}
	}
	return nil, false
}

// PublicFactionName is an organization's name, for anything outside the core
// that needs to say who somebody answers to.
func (w *World) PublicFactionName(id string) string { return w.factionName(id) }
