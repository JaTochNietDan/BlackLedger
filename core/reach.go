package core

// One question, asked once, about everything aimed at a person: is he here to
// be dealt with?
//
// The sweep this replaces asked only whether somebody was out walking, because
// walking was what had just been built. Every other way of being unreachable
// went on being ignored, so a player could send a man on collections from the
// far side of the city, order one out of a police cell, and pay a bonus to a
// man who had been shot the day before.

// OutOfReach explains why somebody cannot be dealt with, in the city's own
// words, or returns "" when they can. Every action with a subject is measured
// against it after the list is built, so an action about a person cannot be
// written without this being asked.
func (w *World) OutOfReach(id string) string {
	n := w.NPC(id)
	// An office is a desk, not a man. Work aimed at one carries the office's
	// own id, and asking for that by name found whoever held it on the first
	// morning — so a campaign that had buried the mayor was told "Mayor Ellis
	// Crane is dead" for ever, while somebody else sat at his desk.
	if IsOfficial(id) {
		n = w.OfficeHolder(id)
		if n == nil {
			o, _ := OfficialByID(id)
			return "There is nobody behind the " + lowerFirst(o.Role) + "'s desk this week"
		}
	}
	if n == nil {
		return "There is nobody by that name in this city"
	}
	if n.Dead {
		return n.Name + " is dead"
	}
	if w.Inside(n) {
		return n.Name + " is being held at Ward Street Station"
	}
	if w.Travelling(n) {
		out := n.Name + " is out on the street"
		if to, ok := PlaceByID(n.Heading); ok {
			out += ", walking to " + to.Name + " — " + counted(max(1, n.Arrives-w.Minute), "minute", "minutes") + " out"
		}
		return out
	}
	if w.CrewOrderFor(n.ID) != nil {
		return n.Name + " is on a headquarters assignment"
	}
	// Being in another building is deliberately not a reason. It was tried
	// here and taken out again: requiring the player to stand in the same room
	// is a change to how the game plays rather than a correction of something
	// it was claiming falsely, and it broke the opening — Leo starts at the bar
	// and the player at the Mariner. The line that does hold is physical: a man
	// at an address can be reached, by going to him or sending word; a man
	// between two addresses is nowhere, and that is what the street means.
	return ""
}
