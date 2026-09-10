package core

import "fmt"

// "It seems like you can rob places or take from people's cars multiple times
// in a row, that should probably be tracked and time limited etc."
//
// You could. Nothing about a robbery was remembered by the place it happened
// to, so the same till could be emptied every forty-five minutes for as long as
// the player felt like standing there, and the same street could lose a car a
// night for a week without anybody in it noticing.
//
// What stops it is not a cooldown bolted onto a button. It is that the place
// changes. A shop that was robbed on Tuesday does not leave Wednesday's takings
// in the same drawer, and a row of cars that lost one to a man with a bag of
// tools is a row that gets looked at from the window after that. Both wear off,
// because a place that can never be robbed again has been deleted rather than
// defended.

const (
	// TillRest is how long a robbed premises keeps its money somewhere else and
	// somebody near the door.
	TillRest = 2 * 1440
	// StreetRest is how long a street that lost a car is worth nothing to
	// somebody carrying a bag of tools.
	StreetRest = 1440
)

// Shy reports how much longer a place is still careful after being robbed, in
// minutes, and zero once it has stopped being.
func (w *World) Shy(id string) int {
	prop := w.Properties[id]
	if prop == nil || prop.Robbed == 0 {
		return 0
	}
	return max(0, prop.Robbed+TillRest-w.Minute)
}

// Curtains reports how much longer a street is being looked at after a car was
// taken apart in it. Watched is taken: that is the police looking at the
// player, which is a different kind of attention entirely.
func (w *World) Curtains(id string) int {
	prop := w.Properties[id]
	if prop == nil || prop.Stripped == 0 {
		return 0
	}
	return max(0, prop.Stripped+StreetRest-w.Minute)
}

// soon puts a number of minutes into the words somebody would actually use. A
// refusal that says "1,432 minutes" is a refusal written by a clock.
func soon(minutes int) string {
	switch {
	case minutes <= 0:
		return "now"
	case minutes < 90:
		return "within the hour"
	case minutes < 1440:
		return fmt.Sprintf("in about %d hours", (minutes+59)/60)
	case minutes < 2880:
		return "tomorrow"
	}
	return fmt.Sprintf("in about %d days", (minutes+1439)/1440)
}
