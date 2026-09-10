package core

import "fmt"

// Where people are, and how the player comes to know it.
//
// The city has always known where everybody is, and so has the player: the
// People screen listed every living soul with their address beside them, and
// going after somebody was a matter of reading it off the list. "It probably
// would make sense that we don't always know everyone's location and that's a
// service that we have to pay for to gain that information from... This would
// also act as a way of making it harder to make attempts on people's lives in
// the game as it would be a drawn out task more than anything."
//
// This is the foundation of that: the player knows where somebody is because
// they saw them there, and what they saw goes stale. Nothing is authored — a
// sighting is a fact about a room the player stood in.
//
// What follows it, in order, is asking other people (who will want something
// for it, and may think less of you for asking about their own), and the same
// limitation on the city's own killers, who should know where you live and not
// where you are.

// SightingLasts is how long a sighting is worth acting on. A day: long enough
// that a player who has been paying attention can use what they know, short
// enough that a name they saw last week is a name they have to find again.
const SightingLasts = 1440

// Seen is somebody the player laid eyes on somewhere, and when. Sighting is
// taken: that is the thing a journey passes in the street.
type Seen struct {
	Where string `json:"where"`
	When  int    `json:"when"`
}

// See records that the player is in the same room as somebody. Called wherever
// the room is drawn, because that is exactly the moment it becomes true.
func (w *World) See(n *NPC) {
	if n == nil || n.Dead || w.Travelling(n) {
		return
	}
	if w.Sightings == nil {
		w.Sightings = map[string]Seen{}
	}
	w.Sightings[n.ID] = Seen{Where: n.Location, When: w.Minute}
}

// SeeTheRoom records everybody standing where the player is. One call, so
// nothing has to remember to do it name by name.
func (w *World) SeeTheRoom() {
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Location == w.Player.Location {
			w.See(n)
		}
	}
}

// LastSeen is where the player last saw somebody and how long ago, and whether
// they ever did.
func (w *World) LastSeen(id string) (Seen, bool) {
	s, ok := w.Sightings[id]
	return s, ok
}

// KnowsWhere reports whether the player knows where somebody is well enough to
// go and find them. Their own people are always findable — they work for you —
// and so is anybody standing in front of them.
func (w *World) KnowsWhere(id string) bool {
	n := w.NPC(id)
	if n == nil || n.Dead {
		return false
	}
	if n.Location == w.Player.Location && !w.Travelling(n) {
		return true
	}
	if n.Faction != "" && n.Faction == w.PlayerOrganizationID() {
		return true
	}
	if w.EmployerOf(id) != "" && w.Own(w.EmployerOf(id)) {
		return true
	}
	s, ok := w.Sightings[id]
	return ok && w.Minute-s.When <= SightingLasts && s.Where == n.Location
}

// WhereNote is what the player can say about where somebody is, which is not
// always where they are.
func (w *World) WhereNote(id string) string {
	n := w.NPC(id)
	if n == nil || n.Dead {
		return ""
	}
	if w.KnowsWhere(id) {
		place, _ := PlaceByID(n.Location)
		return place.Name
	}
	s, ok := w.Sightings[id]
	if !ok {
		return "Nobody has told you where to find them"
	}
	place, _ := PlaceByID(s.Where)
	hours := (w.Minute - s.When) / 60
	if hours < 1 {
		return "Last seen at " + place.Name
	}
	return fmt.Sprintf("Last seen at %s, %s ago", place.Name, counted(hours, "hour", "hours"))
}
