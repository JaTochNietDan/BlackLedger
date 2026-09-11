package core

import "fmt"

// Saint Agnes, and hearing things.
//
// The last of the four rooms a family holds, and the busiest in the city:
// thirty people drink there of an evening, the fixer stands at one end of it,
// the envelope job starts there, there is a game behind it, and two families
// that will not speak to each other will speak in its back room. All of that
// was true whoever held the deed, and holding it did nothing whatsoever.
//
// What a bar is, that no other room in this city is, is somewhere everything
// gets said out loud. Reach is the game's word for how well the player hears
// things — it decides whether somebody warns you that one side came to a
// sitdown to finish it rather than settle it, and whether a word reaching Ward
// Street about what you have been moving can be traced to a name. Until now the
// only ways to raise it were buying the fixer coffee, five times, and putting a
// telephone in the hall.
//
// A bar of your own is the third. It is the same size as the telephone on
// purpose: a room is worth a line into the city, not two.

const (
	// BarReach is what a working bar of yours adds to how well you hear
	// things. The same as a telephone in the hall.
	BarReach = 1
	// SaloonKind is the trade a public house runs on.
	SaloonKind = "saloon"
)

// TheBar is a public house of the player's that is actually open, and nothing
// if they hold none. A room with nobody pulling and a room with something wrong
// in the cellar are both rooms where nothing gets said in front of you.
func (w *World) TheBar() string {
	for _, l := range Locations {
		if l.Kind != SaloonKind || !w.Own(l.ID) {
			continue
		}
		prop := w.Properties[l.ID]
		trade, runs := TradeOf(l.ID)
		if prop == nil || !runs || prop.Trouble || prop.Staff < trade.Hands {
			continue
		}
		return l.ID
	}
	return ""
}

// EarsInTheRoom is what holding a bar adds to Reach.
func (w *World) EarsInTheRoom() int {
	if w.TheBar() == "" {
		return 0
	}
	return BarReach
}

// WhatTheBarHears is the room saying what it is worth, for a card to carry.
func (w *World) WhatTheBarHears(id string) string {
	if w.TheBar() != id {
		return ""
	}
	return fmt.Sprintf("Everything in this district gets said out loud in here, and you own the room it is said in. Worth %d to how well you hear things, on top of anybody you know.", BarReach)
}
