package core

import "fmt"

// The city now keeps hours: people are at their posts through the morning and
// in the bars and clubs after midday. None of that reached the player. Standing
// in Saint Agnes at ten at night with fourteen people in it read exactly like
// standing in it at nine in the morning with two — the only trace of the
// difference was a number in the corner of a filter, and a number with nothing
// to compare it to says nothing.
//
// A room that is full should say so, in the words somebody standing in the
// doorway would use. This is a description of a real count and nothing else:
// no threshold here changes what anybody can do.

// crowdBands are read as a share of everybody still living in the city, so the
// description keeps meaning something as the population rises and falls.
const (
	roomFull  = .20
	roomBusy  = .12
	roomQuiet = .04
)

// InTheRoom counts the people actually in a room. Somebody out on the street
// between two addresses is in neither of them.
func (w *World) InTheRoom(id string) int {
	n := 0
	for i := range w.NPCs {
		if p := &w.NPCs[i]; !p.Dead && p.Location == id && !w.Travelling(p) {
			n++
		}
	}
	return n
}

// living is everybody still on their feet, which is what a room's headcount is
// worth comparing against.
func (w *World) living() int {
	n := 0
	for i := range w.NPCs {
		if !w.NPCs[i].Dead {
			n++
		}
	}
	return n
}

// RoomNote says how busy a room is, or nothing when it is neither. A room with
// an unremarkable number of people in it is not worth a sentence.
func (w *World) RoomNote(id string) string {
	here, all := w.InTheRoom(id), w.living()
	if all == 0 {
		return ""
	}
	place, ok := PlaceByID(id)
	if !ok {
		return ""
	}
	evening := Evening(w.Minute)
	share := float64(here) / float64(all)
	switch {
	case here == 0:
		return "There is nobody else in here."
	case share >= roomFull && evening:
		return fmt.Sprintf("%s is full tonight. %d people, and the ones by the door are watching who comes in.", place.Name, here)
	case share >= roomFull:
		return fmt.Sprintf("%d people in here, which for this hour is a crowd.", here)
	case share >= roomBusy && evening:
		return fmt.Sprintf("A good crowd in tonight, %d of them, and the talk stops when a stranger passes.", here)
	case share >= roomBusy:
		return fmt.Sprintf("Busy: %d people, and somebody is always within earshot.", here)
	case share <= roomQuiet && evening:
		return fmt.Sprintf("Nearly empty at this hour: %s, and nobody else worth counting.", heads(here))
	case share <= roomQuiet:
		return fmt.Sprintf("Quiet. %s, and anything said here is heard.", heads(here))
	}
	return ""
}

// heads counts people the way somebody looking round a room would.
func heads(n int) string {
	if n == 1 {
		return "one other person"
	}
	return fmt.Sprintf("%d other people", n)
}
