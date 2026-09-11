package core

// What an undertaker lives on.
//
// The trade went into the city with a counter that counts the week's dead and a
// line in the working brief calling it "the one trade in this city whose custom
// is made entirely by everybody else's work". The counter said it and nothing
// did it: the room earned its hourly figure whether the district buried two
// people that week or twenty. That is a sentence describing a consequence with
// no code behind it, which is the fault this log has now recorded five times,
// and the fifth one was written by me the night before.
//
// It is the same shape as a garage and the city's broken glass. A car goes to
// pieces in the street, every bench in the city has more work, and the one the
// player holds is where their own city's business ends up. A funeral is that
// with nobody to argue about the bill.

const (
	// BurialTrade is what one funeral does for the room that takes it. A city
	// alone kills about twenty-eight people in a long campaign, so a parlour
	// left to the city drifts up rather than jumps.
	BurialTrade = 2
	// BurialPauper is the exception. Somebody with no name in this city and
	// nobody to send a card to is buried out of the parish's money, which is
	// not money. The room does the work and is not better off for it.
	BurialPauper = 0
)

// TheParlour is where this city's funerals go. A parlour of the player's own
// first, then the first one on the list, which keeps this off map iteration
// order the same way the garage does.
func (w *World) TheParlour() string {
	first := ""
	for _, l := range Locations {
		if l.Kind != "undertaker" || w.Properties[l.ID] == nil {
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

// Bury is a funeral, and what it is worth to the room that takes it. Called
// wherever somebody dies, so a quiet week is a quiet week for the parlour too.
func (w *World) Bury(person *NPC) {
	id := w.TheParlour()
	if id == "" || person == nil {
		return
	}
	// Somebody nobody knew and nobody is paying for.
	worth := BurialTrade
	if person.Faction == "" && person.Rank < RankSoldier && person.Purse < BurialPurse {
		worth = BurialPauper
	}
	if worth == 0 {
		return
	}
	w.ShiftCustom(id, "nobody in this district dying", worth)
}

// BurialPurse is what somebody has to be carrying before anybody expects to be
// paid for burying them.
const BurialPurse = 25
