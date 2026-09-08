package core

import "fmt"

// A location was a set of premises with a list of verbs attached. But half of
// what a player does anywhere is done to somebody who happens to be standing
// there — lending them money, collecting it, putting them on, putting them out,
// taking what they are carrying — and rendering "Lend Perla Mraz money" as one
// more card in a column of verbs throws away the thing the decision is actually
// about, which is Perla Mraz.
//
// This is who is in the room. The core already knows all of it; it has simply
// never been asked the question in this shape.

// Presence is one person where the player is standing, and everything the
// player has earned the right to know about them.
type Presence struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Role    string `json:"role,omitempty"`
	Faction string `json:"faction,omitempty"`
	// Standing is the one line that says who this person is to the player.
	Standing string `json:"standing"`
	// Trust is what they think of the player, and Sore what they hold against
	// them. Both are only shown for people the player has reason to know.
	Trust int `json:"trust,omitempty"`
	Sore  int `json:"sore,omitempty"`
	// Owes is what they owe, if anything, and Overdue whether it is late.
	Owes    int  `json:"owes,omitempty"`
	Overdue bool `json:"overdue,omitempty"`
	// Yours is whether they answer to the player, Known whether the player has
	// any reason to know their name at all.
	Yours bool `json:"yours,omitempty"`
	Known bool `json:"known,omitempty"`
	// Temperament and Manner are what the player has learned of their
	// character; empty for a stranger.
	Temperament string `json:"temperament,omitempty"`
}

// standingOf is the single line under somebody's name: the most important true
// thing about them from where the player is standing.
func (w *World) standingOf(n *NPC) string {
	switch {
	case w.Inside(n):
		return "In a cell at Ward Street Station"
	case n.Faction == w.PlayerOrganizationID():
		if posted := w.postedWhere(n.ID); posted != "" {
			return "Yours · on the door at " + posted
		}
		return "Yours · " + fmt.Sprintf("$%d a day", MemberWage)
	case len(w.Player.Crew) > 0 && w.Player.Crew[0].ID == n.ID:
		return fmt.Sprintf("Your crew · %d loyalty", w.Player.Crew[0].Loyalty)
	case IsOfficial(n.ID):
		if w.Retained(n.ID) {
			o, _ := OfficialByID(n.ID)
			return fmt.Sprintf("%s · paid $%d a day", n.Role, o.Retainer)
		}
		return n.Role
	case w.isRoleHolder(n):
		return n.Role
	case n.Faction != "":
		if n.Rank >= RankLeader {
			return "Head of " + w.factionName(n.Faction)
		}
		return w.factionName(n.Faction)
	case n.Role != "":
		return n.Role
	}
	return "Nobody in particular, yet"
}

// postedWhere is the place somebody of the player's is standing on the door of.
func (w *World) postedWhere(id string) string {
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Posted == id {
			return l.Name
		}
	}
	return ""
}

// PeopleHere is everybody standing where the player is, in the order a person
// would notice them: their own first, then the people they know, then whoever
// else is in the room.
func (w *World) PeopleHere(id string) []Presence {
	out := []Presence{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != id {
			continue
		}
		known := w.Known(n)
		p := Presence{
			ID: n.ID, Name: n.Name, Role: n.Role,
			Standing: w.standingOf(n),
			Yours:    n.Faction == w.PlayerOrganizationID(),
			Known:    known,
		}
		if known {
			p.Faction = w.factionName(n.Faction)
			p.Trust, p.Sore = n.Trust, n.Sore
			p.Temperament = TemperamentOf(n).Label
		}
		if l := w.LoanTo(n.ID); l != nil {
			p.Owes, p.Overdue = l.Owed, l.Missed > 0
		}
		out = append(out, p)
	}
	// Yours first, then anybody who owes you, then people you know.
	rank := func(p Presence) int {
		switch {
		case p.Yours:
			return 0
		case p.Owes > 0:
			return 1
		case p.Known:
			return 2
		}
		return 3
	}
	for i := range out {
		for j := i + 1; j < len(out); j++ {
			if rank(out[j]) < rank(out[i]) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
