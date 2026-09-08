package core

import "fmt"

// The city view shows twelve addresses and gives no sense of how far any of
// them is. A place across town and one on the next corner look identical, and
// the journey time only appears after the place has been selected and the
// travel action read — which is to say, after the player has already decided
// where they were going.
//
// The clock is the scarcest thing in this game. How long it takes to get
// somewhere is a cost, and a cost belongs on the thing it is a cost of.

// Away is how long it takes the player to reach a place from where they are
// standing, by whatever they actually travel by.
func (w *World) Away(id string) int {
	if id == w.Player.Location {
		return 0
	}
	return w.Journey(w.Player.Location, id)
}

// ReachWord describes that distance the way somebody who lives here would.
func (w *World) ReachWord(id string) string {
	switch minutes := w.Away(id); {
	case minutes == 0:
		return "You are here"
	case minutes <= 15:
		return "Round the corner"
	case minutes <= 30:
		return "A short walk"
	case minutes <= 60:
		return "The other side of the district"
	case minutes <= 90:
		return "Across town"
	}
	return "The far side of the city"
}

// TravelNote is what it costs to go, said once: the words and the minutes.
func (w *World) TravelNote(id string) string {
	minutes := w.Away(id)
	if minutes == 0 {
		return "You are here"
	}
	if w.Driving() {
		walk := TravelMinutes(w.Player.Location, id)
		if walk > minutes {
			return fmt.Sprintf("%s · %d min driving, %d on foot", w.ReachWord(id), minutes, walk)
		}
	}
	return fmt.Sprintf("%s · %d min", w.ReachWord(id), minutes)
}
