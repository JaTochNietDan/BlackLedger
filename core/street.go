package core

import "fmt"

// Travel was minutes passing. The city was doing things during those minutes —
// settling grudges, raiding premises, robbing people — and the player never saw
// any of it because none of it happened where they were.
//
// A journey now sometimes puts them at the scene. Nothing here invents an
// incident: it shows one the city already had a reason for, and whether the
// player understands what they are looking at depends on whether anybody would
// tell them. A man shot in the street for a reason nobody explains is the point
// rather than a failure.

const (
	// StreetChance is how often a journey crosses something. Deliberately low:
	// a city where every walk is an incident is a shooting gallery.
	StreetChance = .16
	// StrayInjury is what being too close to a war costs somebody who was only
	// passing.
	StrayInjury = 18
)

// Sighting is something the player passed, and whether they understood it.
type Sighting struct {
	// Title and Body are what they saw.
	Title, Body string
	// Explained is whether anybody told them why.
	Explained bool
	// Stray is the injury they took for being there, which is usually none.
	Stray int
}

// OnTheWay is what the player passes between two places. Everything it can
// return is drawn from something the city is currently doing, in the order of
// how hard it would be to miss.
func (w *World) OnTheWay(from, to string) (Sighting, bool) {
	if w.WorldRandom() >= StreetChance {
		return Sighting{}, false
	}
	where, _ := PlaceByID(to)
	explained := w.Reach() >= 2

	// A war in progress is the loudest thing on any street in this city, and
	// the one thing that can reach somebody who was only passing.
	for _, c := range w.Conflicts {
		if c.State != "war" {
			continue
		}
		a, b := w.factionName(c.A), w.factionName(c.B)
		body := fmt.Sprintf("Two cars and a lot of shouting outside %s, and then everybody was somewhere else. Somebody was face down in the road when the street filled in again.", where.Name)
		if explained {
			body = fmt.Sprintf("Two cars outside %s and people out of both of them. %s and %s, and no doubt about which was which. Somebody was face down in the road when it was over.", where.Name, a, b)
		}
		stray := 0
		if w.WorldRandom() < .3 {
			stray = StrayInjury
		}
		return Sighting{Title: "Shooting on the way to " + where.Name, Body: body, Explained: explained, Stray: stray}, true
	}

	// Somebody with a heavy grievance, on the day it is heavy enough to matter.
	for _, g := range w.Grudges {
		if g.Weight < GrudgeActs-15 {
			continue
		}
		holder, target := w.NPC(g.Holder), w.NPC(g.Against)
		if holder == nil || target == nil || holder.Dead || target.Dead {
			continue
		}
		body := fmt.Sprintf("Two of them in a doorway near %s, one doing all the talking and neither enjoying it. Neither looked at you.", where.Name)
		if explained || w.Known(holder) {
			body = fmt.Sprintf("%s has %s against a wall in a doorway near %s. It is about %s, and it is not finished.", holder.Name, target.Name, where.Name, g.Because)
		}
		return Sighting{Title: "Something in a doorway", Body: body, Explained: explained || w.Known(holder)}, true
	}

	// Or the police, where the police have a reason to be.
	if w.Player.Heat >= RaidThreshold/2 {
		body := fmt.Sprintf("A car that was not there when you set out is there when you pass %s, and it is still there when you look back.", where.Name)
		if explained {
			body = fmt.Sprintf("Harlow's people, parked where they can see the door of %s. Not hiding it, which is the message.", where.Name)
		}
		return Sighting{Title: "Somebody is watching the street", Body: body, Explained: explained}, true
	}

	// And in a quiet city, the city being quiet.
	return Sighting{}, false
}

// PassThrough shows the player what they went past, and charges them for it if
// they were too close. Called on a completed journey, so an interrupted one
// never collects a scene the player did not finish walking into.
func (w *World) PassThrough(from, to string) {
	sighting, ok := w.OnTheWay(from, to)
	if !ok {
		return
	}
	kind := "politics"
	body := sighting.Body
	if sighting.Stray > 0 {
		injury := w.Absorb(sighting.Stray)
		w.Player.Health = max(0, w.Player.Health-injury)
		w.Ruin(25)
		kind = "danger"
		body += fmt.Sprintf(" You were close enough to it that something came off a wall and found you (-%d health).", injury)
	}
	if !sighting.Explained {
		body += " Nobody you know can tell you what it was about."
	}
	w.Log(sighting.Title, body, kind)
	if w.Player.Health <= 0 {
		w.Die("Caught on the street by somebody else's war.")
	}
}
