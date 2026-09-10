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
	return ""
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
