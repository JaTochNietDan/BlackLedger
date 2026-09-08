package core

import "fmt"

// The city has always told the player what happened in prose and never once
// shown them. Six places emitted a visual cue and all six called it an
// "attack"; a killing, a police raid, a seizure of ground and an arrest — the
// four loudest things that happen in this game — produced nothing to look at
// at all. So the most dramatic moments of a campaign arrive as a paragraph in
// a list, and the paper reports them the next morning to a player who never saw
// them.
//
// A cue is the core saying what a moment looked like and where it happened. It
// is not choreography — the interface decides how a thing is played — but the
// truth of it belongs here: what kind of event, which building, who was in it,
// and the headline the Herald will carry afterwards, so the paper can arrive
// after the scene rather than instead of it.

// gravity is how much a moment is worth stopping for, so the interface never
// has to guess which of five things that happened in one command is the one to
// show. It is deliberately the same ordering the city uses for its own
// temperature: what a city notices.
var gravity = map[string]int{
	"killing": 9, "explosion": 8, "gunfight": 7, "raid": 6,
	"seizure": 5, "arrest": 5, "attack": 4, "robbery": 2,
}

// Gravity is what a kind of moment is worth stopping for.
func Gravity(kind string) int { return gravity[kind] }

// Hold is how long the theatre should stay on a moment, in milliseconds. A
// robbery and a killing were played for exactly the same two and a half
// seconds; the city already knew one was worth more than the other.
func Hold(kind string) int {
	return 1800 + Gravity(kind)*320
}

// Witness records that something worth seeing happened somewhere. Callers pass
// the place it happened in, the words for it, and whoever was standing in it.
func (w *World) Witness(kind, place, caption, headline string, actors ...string) {
	if _, ok := PlaceByID(place); !ok {
		// A moment with nowhere to happen is a moment nothing can show. Better
		// to drop it than to send the interface somewhere that does not exist.
		return
	}
	named := []string{}
	for _, id := range actors {
		if n := w.NPC(id); n != nil {
			named = append(named, n.Name)
		}
	}
	w.VisualCues = append(w.VisualCues, VisualCue{
		ID: ID(), Kind: kind, Target: place, Caption: caption,
		Headline: headline, Actors: named, Gravity: Gravity(kind), Minute: w.Minute,
	})
}

// Worth is the moment out of everything that happened worth taking the player
// to. Nothing is more distracting than five scenes in a row for one decision.
func (w *World) Worth() (VisualCue, bool) {
	best, found := VisualCue{}, false
	for _, cue := range w.VisualCues {
		if !found || cue.Gravity > best.Gravity {
			best, found = cue, true
		}
	}
	return best, found
}

// whereTheyWere is the place somebody was standing when it happened, falling
// back to the street outside if the city has lost track of them.
func (w *World) whereTheyWere(n *NPC) string {
	if n == nil {
		return ""
	}
	if _, ok := PlaceByID(n.Location); ok {
		return n.Location
	}
	return ""
}

// witnessKilling is the cue for a death, built from the record the city already
// wrote so the scene and the paper cannot disagree.
func (w *World) witnessKilling(victim *NPC, cause, headline string) {
	place := w.whereTheyWere(victim)
	if place == "" {
		return
	}
	l, _ := PlaceByID(place)
	w.Witness("killing", place,
		fmt.Sprintf("%s at %s. %s", victim.Name, l.Name, cause),
		headline, victim.ID)
}
