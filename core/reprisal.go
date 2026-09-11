package core

import "fmt"

// What a family does when it comes for you and you are not there.
//
// It wrecked the house, every time. Held in a cell, out of the city, warned and
// standing somewhere else, or simply not seen lately — four different ways of
// being out of reach and one answer to all of them, which is a family that
// cannot think of anything to do about a man they cannot find except break his
// furniture.
//
// They can find the people who work for him. That is the whole point of having
// people, from their side of it: somebody who answers to you is somebody who
// can be reached instead of you, and until this the only thing that ever killed
// one of the player's own was the player walking them into a rival's holding.
// Signing somebody on had no risk in it at all, and the funeral built for them
// could not be reached by anything the city does on its own.

// ReprisalChance is how often a family that cannot reach the player takes it out
// on somebody who works for them instead. The rest of the time it is still the
// house: a third is often enough to make people a commitment and rare enough
// that a campaign is not a funeral every week.
const ReprisalChance = .34

// ReprisalTrust is what the rest of them make of it. Somebody died for working
// for you, which is a different thing from a wage arriving late.
const ReprisalTrust = 10

// tookOneOfYours reports whether the reprisal fell on a person rather than a
// building, and does it if so.
func (w *World) tookOneOfYours(plot Plot) bool {
	if !w.Incorporated() || w.Random() >= ReprisalChance {
		return false
	}
	victim := w.casualty(w.PlayerOrganizationID())
	if victim == nil {
		return false
	}
	name, family := victim.Name, w.factionName(plot.Actor)
	w.KillBy(victim.ID, nil, fmt.Sprintf("%s could not find you.", family))
	for _, other := range w.OwnPeople() {
		other.Trust = max(0, other.Trust-ReprisalTrust)
	}
	w.Antagonize(plot.Actor, w.PlayerOrganizationID(), 8)
	w.Log("They took it out on "+name,
		fmt.Sprintf("%s went looking for you and found somebody who works for you instead. Everybody else on your books knows which of the two they are.", family),
		"danger")
	return true
}

// outOfReach is what happens when they come and you are not there to be found.
func (w *World) outOfReach(plot Plot) {
	if w.tookOneOfYours(plot) {
		return
	}
	w.wreckTheHouse()
}
