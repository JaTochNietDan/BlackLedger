package core

import "fmt"

// A sitting is a thing the world knows about.
//
// It used to be a screen the interface opened on its own: a local flag said
// "show the tables" and nothing in the city changed. That is fine until you
// notice what the tables are made of. The last hand, the last spin and where
// the three drums stopped are all saved state, kept so a game can be drawn
// rather than described — so walking up to a table showed you the game the
// last person played, cards face up and a ball in a pocket. Sitting down
// looked exactly like a game that had started without you.
//
// So sitting down is a decision, and taking a seat clears the table in front
// of you. Getting up is a decision too, and you cannot make it in the middle
// of a hand: the money is already down.

// Playable reports whether there is anything to sit down to here. The room
// behind the poolhall counts: a hand of cards against people who live here is a
// game you sit down to and leave, not a box above the action list.
func Playable(id string) bool { return HasTables(id) || HasMachines(id) || HasBackRoom(id) }

// SitReadiness explains why no seat can be taken here, or returns "".
func (w *World) SitReadiness(id string) string {
	if !Playable(id) {
		return "There is nothing to play in here"
	}
	if w.Player.Location != id {
		return "You are not in the room"
	}
	return ""
}

// Sit takes a seat, and clears whatever the last player left on the table.
func (w *World) Sit(id string) error {
	if reason := w.SitReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	w.Seated = id
	w.clearTable()
	return nil
}

// RiseReadiness explains why the player cannot get up, or returns "".
func (w *World) RiseReadiness() string {
	if w.Hand != nil && !w.Hand.Done {
		return "Finish the hand first"
	}
	if w.Game != nil && !w.Game.Done {
		return "There is money of yours in the middle of the table"
	}
	return ""
}

// Rise gets up and leaves, taking the table's memory with it. Nothing is left
// for whoever sits down next, including the player themselves.
func (w *World) Rise() error {
	if reason := w.RiseReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	w.Seated = ""
	w.clearTable()
	return nil
}

// clearTable forgets the last hand, spin and pull. A settled hand is a picture
// of something that is over; a live one cannot be cleared and is not, because
// RiseReadiness refuses before this is ever reached.
func (w *World) clearTable() {
	if w.Hand != nil && w.Hand.Done {
		w.Hand = nil
	}
	if w.Game != nil && w.Game.Done {
		w.Game = nil
	}
	w.Spin, w.Reels = nil, nil
}

// LeaveTable is what walking out of the room does. A seat is taken in one room
// and the city is full of others; nobody is at the tables in a place they are
// not standing in.
func (w *World) LeaveTable() {
	if w.Seated == "" || w.Seated == w.Player.Location {
		return
	}
	w.Seated = ""
	w.clearTable()
}
