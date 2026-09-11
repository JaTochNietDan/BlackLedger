package core

import "fmt"

// The Monarch, and the three rooms like it.
//
// Every address in this city that can be bought is a business: it has hands to
// hire, supplies to keep up, a kind of trouble of its own, a front to launder
// behind, and somebody behind a counter worth asking what they have seen. Every
// address that can only be taken by force is a number.
//
// That is exactly backwards. The four addresses a family holds — Saint Agnes,
// Pier 14, The Monarch and the Mercer Exchange — are the busiest rooms in the
// game. Thirty people drink at the bar of an evening and twenty-eight at the
// club. You fight a war for a family's seat, win it, and what you have won is
// an income figure and nothing to do about it: no staff, no counter, no
// supplies, no cover, no trouble. The best room in the city goes inert the
// moment it is yours.
//
// The club is the first of the four to be given a trade. The other three are
// the same gap and are written down as next.
//
// Its own reach is standing rather than money, which nothing else in the game
// pays. A room where the city drinks, with your name over the door, is where a
// reputation is actually made — and standing is what gates buying anything at
// all, so a club earns you the thing that lets you own the rest.

const (
	// ClubPerHead is how many people have to drink in a room of yours before
	// one of them says your name somewhere it matters.
	ClubPerHead = 9
	// ClubNightly caps what one night in one room is worth. A club builds a
	// reputation; it does not hand one over.
	ClubNightly = 3
	// ClubKind is the trade the city's night rooms run on.
	ClubKind = "club"
)

// ClubFloor reports whether this address is a room the city drinks in that
// earns its holder standing.
func ClubFloor(id string) bool {
	place, ok := PlaceByID(id)
	return ok && place.Kind == ClubKind
}

// TheEveningCrowd is whose evening this room is, rather than who is standing in
// it when the day turns over.
//
// The clock settles the city at midnight, by which time everybody has gone
// home: measured at that moment, the busiest room in the game held eight
// people, and a rule pitched at nine paid nothing almost every night. The
// city's own card game had already hit this and says so in `core/citygame.go`.
// Asking whose evening a room is gives the twenty-eight who actually drink
// there.
func (w *World) TheEveningCrowd(id string) int {
	crowd := 0
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || haunt(n.ID) != id {
			continue
		}
		// Your own people are staff and hangers-on. A room full of nobody but
		// your own is not a room the city is talking about.
		if n.Faction != "" && n.Faction == w.PlayerOrganizationID() {
			continue
		}
		crowd++
	}
	return crowd
}

// StandingFromTheDoor is what a night at a room of yours is worth in
// reputation: who drinks there, over how many it takes to be talked about.
// Nothing if it is short-handed or in trouble — a room with a queue and no
// staff, or a room with something wrong in the back, is a room people talk
// about for the wrong reason.
func (w *World) StandingFromTheDoor(id string) int {
	if !ClubFloor(id) || !w.Own(id) {
		return 0
	}
	prop := w.Properties[id]
	trade, runs := TradeOf(id)
	if prop == nil || !runs || prop.Trouble || prop.Staff < trade.Hands {
		return 0
	}
	return min(ClubNightly, w.TheEveningCrowd(id)/ClubPerHead)
}

// ClubNight pays it, once a day from the clock.
func (w *World) ClubNight() {
	for _, l := range Locations {
		gained := w.StandingFromTheDoor(l.ID)
		if gained <= 0 {
			continue
		}
		w.Player.Respect += gained
		w.Log("A good night at "+l.Name,
			fmt.Sprintf("%s drink here, and your name is over the door. %s standing.",
				counted(w.TheEveningCrowd(l.ID), "person", "people"), plainly(gained, "+1", fmt.Sprintf("+%d", gained))),
			"business")
	}
}

// SeatOf reports whether this address is one of a family's own holdings —
// somewhere with no price that changes hands only when somebody takes it.
func (w *World) SeatOf(owner, id string) bool {
	for _, held := range w.FamilyHoldings(owner) {
		if held == id {
			return true
		}
	}
	return false
}
