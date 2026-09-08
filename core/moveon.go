package core

import "fmt"

// The city can take the player's ground. The player cannot take anybody's.
// Sabotage damages a holding and a charge wrecks one, and neither of them has
// ever moved a name on a deed — which meant a war with the player was something
// they could only survive or end, never win.
//
// A move is the same raid the city runs, with the player's organization as the
// attacker. The one difference is that the city always goes for the weakest
// thing somebody holds and the player picks.

const (
	// MoveMinutes is what committing your people to it takes.
	MoveMinutes = 120
	// MoveHeat is what a daylight move on somebody else's premises draws.
	MoveHeat = 14
	// MoveStanding is the strength below which nobody would follow you into it.
	MoveStanding = 30
)

// MoveTarget reports the organization holding a property, and whether the
// player could move on it: they are an organization, and the two of them are
// past being civil.
func (w *World) MoveTarget(id string) (*Faction, bool) {
	if !w.Incorporated() {
		return nil, false
	}
	prop := w.Properties[id]
	if prop == nil || prop.Owner == w.PlayerOrganizationID() {
		return nil, false
	}
	holder := w.faction(prop.Owner)
	if holder == nil {
		return nil, false
	}
	c := w.Conflict(holder.ID, w.PlayerOrganizationID())
	if c == nil || (c.State != "war" && c.State != "feud") {
		return nil, false
	}
	return holder, true
}

// MoveOnReadiness explains why a move cannot be made, or returns "".
func (w *World) MoveOnReadiness(id string) string {
	if !w.Incorporated() {
		return "A man does not take ground. An organization does"
	}
	holder, ok := w.MoveTarget(id)
	if !ok {
		return "There is nobody here you are at odds with"
	}
	if w.Player.Location != id {
		return "You would have to be there"
	}
	if w.PlayerStrength() < MoveStanding {
		return fmt.Sprintf("Nobody would follow you into that. You need %d strength", MoveStanding)
	}
	if w.Player.Health < 40 {
		return "You are in no condition for this"
	}
	if w.Allied(holder.ID) {
		return "You stand with them. Break that first"
	}
	_ = holder
	return ""
}

// MoveOn commits the player's organization against a rival's holding. It is
// resolved by the same function the city uses on them, so nothing here is a
// separate set of rules for the player.
func (w *World) MoveOn(id string) error {
	if reason := w.MoveOnReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	holder, _ := w.MoveTarget(id)
	me := w.PlayerOrganization()
	place, _ := PlaceByID(id)

	before := len(w.FamilyHoldings(w.PlayerOrganizationID()))
	strength := me.Power
	w.contestAt(me, holder, id)
	w.Player.Heat = min(100, w.Player.Heat+MoveHeat)
	holder.Goodwill = max(-100, holder.Goodwill-25)
	w.ServeAgainst(holder.ID)
	if c := w.Conflict(holder.ID, w.PlayerOrganizationID()); c != nil {
		w.Antagonize(holder.ID, w.PlayerOrganizationID(), 12)
	}

	// A move that is driven off costs the people who went, not a number.
	if me.Power < strength {
		if victim := w.casualty(w.PlayerOrganizationID()); victim != nil {
			w.KillBy(victim.ID, nil, fmt.Sprintf("They had gone with you to %s.", place.Name))
		} else if len(w.Player.Crew) > 0 && w.Random() < .3 {
			member := w.Player.Crew[0]
			w.Player.Crew = w.Player.Crew[:0]
			if n := w.NPC(member.ID); n != nil {
				w.KillBy(n.ID, nil, fmt.Sprintf("They had gone with you to %s.", place.Name))
			}
		} else {
			injury := w.Absorb(14 + int(w.Random()*20))
			w.Ruin(30)
			w.Player.Health = max(0, w.Player.Health-injury)
			w.Log("You went in with them", fmt.Sprintf("It did not go your way at %s and you came out of it with %d less health.", place.Name, injury), "danger")
			if w.Player.Health <= 0 {
				w.Die("A move on " + place.Name + " that should not have been made.")
			}
		}
	}

	if len(w.FamilyHoldings(w.PlayerOrganizationID())) > before {
		w.Player.Respect += 10
		w.Log(place.Name+" is yours", fmt.Sprintf("%s could not hold it and you could. This is how ground changes hands in this city, and it is how it will be taken back.", holder.Name), "politics")
	}
	return nil
}
