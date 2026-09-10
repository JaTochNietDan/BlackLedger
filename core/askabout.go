package core

import "fmt"

// Asking after somebody. The player knows where people are because they saw
// them; this is the other way of finding out, and the one the note was really
// about: "maybe only some people know where someone is and we would have to be
// able to convince us to give them their location, especially if it's a rival
// family lead... We don't want it to be a situation where we can just talk to
// one singular character and find out people's locations by paying a fee,
// that's too simplistic. We want a dynamic knowledge situation and bargainning
// situation. Some people may not want to give us that information and even
// asking them could affect our reputation with them."
//
// So there is no broker and no fee. You ask whoever is standing in front of
// you, they answer if they know and if they think enough of you, and asking
// somebody to give up their own family is a thing they remember.

const (
	// TellsYou is the trust at which somebody will tell you where a stranger
	// is. Below it they say they have not seen them, whether or not that is
	// true.
	TellsYou = 25
	// GivesUpTheirOwn is what it takes to give up somebody from their own
	// family, which is a different thing from gossip about a stranger.
	GivesUpTheirOwn = 70
	// AskingGalls is what somebody thinks of being asked to give up one of
	// their own and refusing. Being asked is the insult; the refusal is not.
	AskingGalls = 6
	// Secondhand is how old the answer is. What somebody tells you is where
	// they last saw them, not where they are now, so it is worth acting on and
	// not the same as having seen it yourself.
	Secondhand = 240
)

// Knows reports whether this person could say where the target is. Standing in
// the same room, working in the same building, or belonging to the same family
// — all of which are facts the city already holds.
func (w *World) Knows(who, target string) bool {
	n, mark := w.NPC(who), w.NPC(target)
	if n == nil || mark == nil || n.Dead || mark.Dead || who == target {
		return false
	}
	if n.Location == mark.Location {
		return true
	}
	if mark.Faction != "" && n.Faction == mark.Faction {
		return true
	}
	if at := w.EmployerOf(target); at != "" && at == w.EmployerOf(who) {
		return true
	}
	return false
}

// AskAboutReadiness explains why there is no point asking, or returns "".
func (w *World) AskAboutReadiness(who, target string) string {
	n, mark := w.NPC(who), w.NPC(target)
	if n == nil || n.Dead {
		return "There is nobody here by that name"
	}
	if mark == nil || mark.Dead {
		return "There is nobody to ask about"
	}
	if n.Location != w.Player.Location || w.Travelling(n) {
		return n.Name + " is not here"
	}
	if w.KnowsWhere(target) {
		return "You know where to find " + mark.Name
	}
	return ""
}

// AskAbout is the question and the answer. Whether they know is one thing and
// whether they will say is another.
func (w *World) AskAbout(who, target string) error {
	if reason := w.AskAboutReadiness(who, target); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	n, mark := w.NPC(who), w.NPC(target)
	theirs := mark.Faction != "" && n.Faction == mark.Faction
	switch {
	case !w.Knows(who, target):
		w.Log("You ask "+n.Name+" about "+mark.Name,
			fmt.Sprintf("%s has not seen %s and does not pretend otherwise.", n.Name, mark.Name), "personal")
	case theirs && n.Trust < GivesUpTheirOwn:
		// Asking somebody to give up one of their own is the insult, whatever
		// they answer.
		w.Aggrieve(n.ID, AskingGalls, "the day you asked them to give up one of their own")
		n.Trust = max(0, n.Trust-AskingGalls)
		w.Log("You ask "+n.Name+" about "+mark.Name,
			fmt.Sprintf("%s is one of theirs, and %s looks at you for a moment before saying they could not tell you.", mark.Name, n.Name), "personal")
	case n.Trust < TellsYou:
		w.Log("You ask "+n.Name+" about "+mark.Name,
			fmt.Sprintf("%s says they have not seen %s. They say it the way people say it when they would rather not be asked again.", n.Name, mark.Name), "personal")
	default:
		// They tell you where they last saw them, which is not the same as
		// where they are.
		if w.Sightings == nil {
			w.Sightings = map[string]Seen{}
		}
		w.Sightings[target] = Seen{Where: mark.Location, When: w.Minute - Secondhand}
		place, _ := PlaceByID(mark.Location)
		n.Trust = min(100, n.Trust+1)
		w.Log("You ask "+n.Name+" about "+mark.Name,
			fmt.Sprintf("%s saw %s at %s, and says so without making anything of it.", n.Name, mark.Name, place.Name), "personal")
	}
	return nil
}
