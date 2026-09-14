package core

import "fmt"

// A car taken apart in the street is one link; the row it was parked in is
// another. Somebody working a street at night goes through what is there, and
// most of what they touch they leave behind with the glass out of it. Those
// cars are not gone. They are work, and work is the thing a garage actually
// lives on — parts coming in pay it something, but a person standing at the
// counter paying to be put back on the road is the trade itself.
//
// The money comes out of the owner's own purse, which is why a city where
// people are broke is a city of broken cars: somebody who cannot find the fee
// keeps driving it as it is, and the garage never sees them.

const (
	// GlassCost is what a garage charges to put right what a thief left behind:
	// glass, a lock, a steering column somebody went at with a screwdriver.
	GlassCost = 40
	// RepairTrade is what one of those jobs is worth to the trade doing it.
	// Custom is capped, so a bad month for the city cannot run a garage away.
	RepairTrade = 2
)

// BreakGlass is what a night's thieving does to everything else parked in the
// street. It returns how many cars were left needing a garage, so the act that
// caused it can say so.
func (w *World) BreakGlass(location, except string) int {
	touched := 0
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != location || n.Car == 0 || n.ID == except || n.Hurt {
			continue
		}
		n.Hurt = true
		touched++
	}
	return touched
}

// theGarage is where the city takes that work. Every address in this city has
// a holding behind it, so this is not a filter on whether a garage is open: it
// is the order the work goes in. A garage of the player's own is where their
// own city's business ends up, and otherwise it is the first bench in the list,
// which keeps reports and tests off map iteration order.
func (w *World) theGarage() string {
	first := ""
	for _, l := range Locations {
		if l.Kind != "garage" || w.Properties[l.ID] == nil {
			continue
		}
		if w.Own(l.ID) {
			return l.ID
		}
		if first == "" {
			first = l.ID
		}
	}
	return first
}

// RepairsDay is one person taking a broken car in. It is one a day at most, for
// the same reason a forecourt sells one a day: a city where every window is
// replaced on the same morning is a city where nothing was ever done to
// anybody.
func (w *World) RepairsDay() {
	garage := w.theGarage()
	if garage == "" {
		return
	}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Location == garage {
			w.putRight(n)
		}
	}
}

// putRight is one car being worked on, for somebody standing at the bench. It
// is called when they walk in and again when the day turns over, because a
// sweep at midnight finds a garage everybody has already gone home from: they
// bring the car in during the afternoon and are somewhere else by the small
// hours.
func (w *World) putRight(n *NPC) {
	garage := w.theGarage()
	if garage == "" || n == nil || n.Dead || !n.Hurt || n.Car == 0 || n.Location != garage {
		return
	}
	// Somebody who cannot find the fee drives it broken. That is the link
	// running the other way: a garage in a poor district has less work than the
	// same garage in a district with money in it, off the same crimes.
	if n.Purse < GlassCost {
		return
	}
	n.Purse -= GlassCost
	n.Hurt = false
	w.changeBusinessFunds(garage, GlassCost)
	w.ShiftCustom(garage, "glass and locks after a night's thieving", RepairTrade)
	if w.Own(garage) {
		w.Earn(GlassCost)
		place, _ := PlaceByID(garage)
		w.Log("Work in at "+place.Name,
			fmt.Sprintf("%s brought one in with the glass out of it. $%d for the job.", n.Name, GlassCost), "business")
	}
}
