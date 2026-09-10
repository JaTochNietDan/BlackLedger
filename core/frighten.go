package core

import "fmt"

// Putting the wind up somebody else's counter.
//
// Moving against a business meant damaging the property: go in with the crew,
// take the condition off it, and wait for the answer. That is the loud option
// and it was the only one. A business is people now — named, standing in a
// room, with a wage and a view of the street — and people can be put off coming
// in. It costs the holder the same trade for a couple of days without leaving a
// broken window for anybody to point at.
//
// Nothing here is free. The person you frightened remembers your face, the
// family it belongs to thinks less of you, and standing in their shop doing it
// is not something anybody misses.

const (
	// FrightenMinutes is how long it takes to make the point.
	FrightenMinutes = 25
	// FrightenGalls is what the family thinks of somebody who does this to
	// their people. Less than a broken window, more than nothing.
	FrightenGalls = 14
	// FrightenedOff is what the person carries away from it.
	FrightenedOff = 30
)

// FrightenReadiness explains why there is nobody here to lean on, or returns "".
func (w *World) FrightenReadiness(id string) string {
	prop := w.Properties[id]
	if prop == nil {
		return "There is nothing here to lean on"
	}
	if w.Player.Location != id {
		return "You are not in the room"
	}
	if w.Own(id) {
		return "These are your own people"
	}
	if prop.Owner == "" || prop.Owner == "independent" {
		return "Nobody holds this, so there is nobody it would cost"
	}
	if len(prop.Hands) == 0 {
		return "There is nobody behind the counter to put off"
	}
	if w.Player.Health < 40 {
		return "In this condition nobody would take you seriously"
	}
	// Somebody of theirs standing in it. A quiet word with the staff is a thing
	// you do when nobody is watching the door, and a family minding its own
	// holding is the oldest answer there is to it.
	if who := w.theirsHere(id, prop.Owner); who != nil {
		return who.Name + " is standing in here, and this is not the kind of thing you do in front of somebody"
	}
	return ""
}

// theirsHere is somebody of that organization standing in the room, and nothing
// when nobody is watching the door.
func (w *World) theirsHere(id, faction string) *NPC {
	if faction == "" {
		return nil
	}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead && n.Faction == faction && n.Location == id && !w.Travelling(n) {
			return n
		}
	}
	return nil
}

// Frighten puts one of them off coming in. Which one is whoever is standing
// there — the city already knows, and picking the most useful of them would be
// the player reading a list rather than walking into a shop.
func (w *World) Frighten(id string) error {
	if reason := w.FrightenReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	prop := w.Properties[id]
	who := prop.Hands[0]
	n := w.NPC(who)
	w.walkOut(who, id, "")
	w.Aggrieve(who, FrightenedOff, "the afternoon you came into where they work")
	holder := w.faction(prop.Owner)
	if holder != nil {
		holder.Goodwill = max(-100, holder.Goodwill-FrightenGalls)
	}
	place, _ := PlaceByID(id)
	// Rooms talk. Doing this in front of people is the difference between a
	// quiet word and a thing the street knows about by the evening.
	w.Witness("politics", id, n.Name+" walked out of "+place.Name+" and did not come back.", "")
	holderName := w.HolderName(id)
	w.Log("A word at "+place.Name,
		fmt.Sprintf("%s will not be behind that counter tomorrow, and %s %s short-handed at %d. Nobody broke anything, which is the point of doing it this way.",
			n.Name, holderName, Agree(holderName, "is", "are"), prop.Staff), "danger")
	return nil
}

// And the same move, made against the player.
//
// Every actor in the city runs on the same rules: that is the principle the
// whole living world is built on. The player can put the wind up somebody
// else's counter, and until this the only thing a family could do to a business
// of theirs was break it. A family that has fallen out with them sends
// somebody round instead, which costs the same trade and leaves nothing to
// point at — the same bargain from the other side.

const (
	// TheyLean is the day's chance a family that has fallen out with the player
	// puts somebody off one of their counters. Low: over two months it is about
	// one pair of hands, which is a nuisance rather than a siege.
	TheyLean = .03
	// FallenOut is the goodwill below which a family will do it at all. Nobody
	// leans on the counter of somebody they have no quarrel with.
	FallenOut = -40
)

// TheyFrighten is one day of the families the player has fallen out with
// deciding whether to come for their people. Called every business day.
func (w *World) TheyFrighten() {
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.ID == w.PlayerOrganizationID() || f.Goodwill > FallenOut {
			continue
		}
		if w.WorldRandom() >= TheyLean {
			continue
		}
		for _, l := range Locations {
			prop := w.Properties[l.ID]
			if prop == nil || !w.Own(l.ID) || len(prop.Hands) == 0 {
				continue
			}
			// Not past somebody of the player's on the door, which is the
			// answer to this and the reason to post anybody anywhere.
			if w.PostedAt(l.ID) != nil {
				continue
			}
			who := prop.Hands[0]
			n := w.NPC(who)
			w.walkOut(who, l.ID, "")
			place, _ := PlaceByID(l.ID)
			w.Witness("politics", l.ID, "Somebody came into "+place.Name+" and did not buy anything.", "")
			w.Log("A word at "+place.Name,
				fmt.Sprintf("%s will not be behind that counter tomorrow. They will not say who asked them not to be, and %s %s no reason to hide it.",
					n.Name, f.Name, Agree(f.Name, "has", "have")), "danger")
			break // one address a day, or a family empties you in a week
		}
	}
}
