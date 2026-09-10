package core

// Who is actually in the room.
//
// The city has a real rhythm and none of it reached the ledger. Measured over
// twenty days of one city: at ten in the morning the market holds twelve people
// and the bar seven; at nine at night the bar, the club and the casino hold
// twenty-five each and the docks and the butcher are empty. People go to their
// posts in the morning and to the places that sell a drink in the evening, and
// a butcher earned exactly the same at three in the morning as at noon.
//
// This is the join between the people and the money. Everything that moves
// somebody — a post, a haunt, an errand, a lead holding court, a family taking
// over the shop across the road — now shows up in what a business takes.

const (
	// TypicalRoom is what an ordinary address holds. The city runs about
	// eighty-six people over twenty-six addresses, so three is the middle and
	// the number every room is measured against.
	TypicalRoom = 3
	// PerHead is what one more or one fewer body is worth to a day's takings.
	PerHead = .02
	// QuietRoom and BusyRoom bound it. A shop with nobody in it still has a
	// door and a standing order; a room with thirty people in it is a good
	// night and not a different business.
	QuietRoom = .75
	BusyRoom  = 1.45
)

// Footfall is how many people are in a place who are there as custom: not the
// player, not somebody of theirs standing on the door, and not anybody who is
// out on the street between two addresses.
func (w *World) Footfall(id string) int {
	n := 0
	for i := range w.NPCs {
		who := &w.NPCs[i]
		if who.Dead || who.Location != id || w.Travelling(who) {
			continue
		}
		// Your own people are staff. A man you pay to stand on the door is not
		// somebody who came in and spent anything.
		if who.Faction != "" && who.Faction == w.PlayerOrganizationID() {
			continue
		}
		n++
	}
	return n
}

// RoomTrade is what the room itself is worth to a day's takings today. It
// modulates what custom and capacity already decide rather than replacing them:
// a full room is worth more, not everything.
func (w *World) RoomTrade(id string) float64 {
	if _, running := TradeOf(id); !running {
		return 1
	}
	factor := 1 + float64(w.Footfall(id)-TypicalRoom)*PerHead
	if factor < QuietRoom {
		return QuietRoom
	}
	if factor > BusyRoom {
		return BusyRoom
	}
	return factor
}
