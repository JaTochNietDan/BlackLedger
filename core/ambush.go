package core

import "fmt"

// A hit came for an address rather than for a person. If the player was not
// standing in their own home when it landed, the attackers wrecked the house
// and went away — which is why somebody who spends their day out working could
// never be shot at, and kept coming back to a broken door instead.
//
// They come for you where you are. What that place is worth is the whole
// decision: a room with people in it makes them careful and gives somebody the
// chance to see them come in, your own people there make them slower, and the
// street between two addresses has none of that.

const (
	// AmbushOdds is the chance an unwarned attack kills, before anything about
	// where you are standing or what you are wearing.
	AmbushOdds = .78
	// RoomWarning is the reach at which somebody tells you in time, in a room.
	// StreetWarning is what it takes between two addresses, where there is
	// nobody to tell you and the word has to have reached you before you set
	// out: every contact this city offers, and a telephone or a room of your
	// own on top of them.
	RoomWarning   = 2
	StreetWarning = 6
	// StreetExposure is what being between two addresses adds to that. No
	// walls, no door, and nobody who knows you to shout.
	StreetExposure = .10
	// WitnessCover is what one person in the room is worth, and the most that
	// a crowd can be worth however big it gets. Somebody watching does not stop
	// a killing; it makes people hurry, and hurrying is most of what saves you.
	WitnessCover    = .05
	WitnessCoverMax = .25
	// OwnCover is what one of your own standing there is worth, and its
	// ceiling. Worth more per head than a stranger, because they are the ones
	// who get in the way.
	OwnCover    = .10
	OwnCoverMax = .20
	// CoverMax is what any room can be worth in total. A crowd is not armour.
	CoverMax = .35
	// SeenComingCrowd is how many people it takes for the room itself to warn
	// you: enough that a stranger walking in with his hand in his coat is
	// something somebody notices out loud.
	SeenComingCrowd = 3
)

// InTransit reports whether the player is between two addresses. OnTheStreet
// is taken: it is the list of everybody the city has walking about.
func (w *World) InTransit() bool {
	_, ok := PlaceByID(w.Player.Location)
	return !ok
}

// Cover is what the room the player is standing in is worth against somebody
// who came to kill them. Home has a door instead; the street has nothing.
func (w *World) Cover() float64 {
	if w.Player.Location == w.Player.Home {
		return w.doorProtection()
	}
	if w.InTransit() {
		// Nothing out here but what you are driving, which is why anybody
		// bothers plating one.
		return w.PlateCoverHere()
	}
	strangers, own := 0., 0.
	for _, who := range w.PeopleHere(w.Player.Location) {
		if who.Yours {
			own += OwnCover
		} else {
			strangers += WitnessCover
		}
	}
	return min64(CoverMax, min64(WitnessCoverMax, strangers)+min64(OwnCoverMax, own))
}

// Warned reports whether the player gets a moment to decide rather than a roll
// of the dice. Somebody told you, somebody of yours was watching, your contacts
// reach far enough — or the room is busy enough that a man walking in with his
// hand in his coat is something people say out loud.
func (w *World) Warned(plot Plot) bool {
	if w.InTransit() {
		// The people who watch your door are at your door. This line was
		// already written and could never run: the reach check sat above it and
		// returned true first, so two cups of coffee bought a warning on an
		// empty street between two addresses, where by this file's own account
		// there is nobody to give one.
		//
		// It takes a great deal more than two. Measured with a policy that
		// mugs, strikes and robs its way across the city: at a reach of two it
		// survived every one of twelve campaigns past seventy days, and without
		// the coffee it died in every one of twelve at twenty-five. Fifty
		// dollars was the difference between certain death and no death at all,
		// while doing the same amount of harm and walking further.
		return w.Reach() >= StreetWarning
	}
	if w.Reach() >= RoomWarning {
		return true
	}
	if w.Watchers() > 0 && w.Player.Location == w.Player.Home {
		return true
	}
	if w.Player.Security > 0 {
		return true
	}
	return len(w.PeopleHere(w.Player.Location)) >= SeenComingCrowd
}

// AmbushOddsHere is the chance an unwarned attack kills the player where they
// are standing, all of it in one place so the scene and the roll cannot come to
// different answers about the same room.
func (w *World) AmbushOddsHere() float64 {
	odds := AmbushOdds - float64(w.Player.Armour)*.09 - w.Cover()
	if w.InTransit() {
		odds += StreetExposure
	}
	return odds
}

// whereItHappens names the place for the telling. A journey has no address, so
// it is described rather than named.
func (w *World) whereItHappens() string {
	if w.InTransit() {
		return "the street"
	}
	place, _ := PlaceByID(w.Player.Location)
	return place.Name
}

// wreckTheHouse is what they do when the player is somewhere they cannot be
// reached. A cell is the only such place: the city can get at you anywhere else
// it can walk.
func (w *World) wreckTheHouse() {
	w.Properties[w.Player.Home].Condition = max(10, w.Properties[w.Player.Home].Condition-45)
	w.Witness("attack", w.Player.Home, "Somebody came armed and damaged your residence while you were away.", "")
	w.Log("Someone came looking", "They could not get to you, so they went to where you sleep. Armed, and gone before anyone could identify them.", "danger")
}

// ambushScene is the moment before it happens, wherever it is happening.
func (w *World) ambushScene(plot Plot) {
	where := w.whereItHappens()
	var body string
	switch {
	case w.InTransit():
		body = "A car pulls level with you and slows. "
		if w.Driving() {
			body += "There is nothing between you and it but your own doors."
		} else {
			body += "There is nowhere on this street to be."
		}
	case w.Player.Location == w.Player.Home:
		body = "A car stops outside " + where + ". "
		if w.Guard() > 0 {
			body += "Your security raises the alarm."
		} else {
			body += "A contact calls: leave by the back, now."
		}
	default:
		body = "Two of them come through the door of " + where + " and do not look at the bar. "
		if w.Guard() > 0 {
			body += "One of yours is already moving."
		} else if len(w.PeopleHere(w.Player.Location)) >= SeenComingCrowd {
			body += "The room has gone quiet around you."
		} else {
			body += "There is nobody in here but you and them."
		}
	}
	w.Event = &Scene{ID: "attack-" + plot.ID, Title: "They came for you", Body: body + " You have moments to act.",
		Speaker: w.HolderID("fixer"), Kind: "attack", Source: "authored", Minute: w.Minute,
		Choices: []Choice{
			{ID: "escape", Label: "Get out", Detail: "A chance to escape. Security and contacts help; injuries reduce your odds."},
			{ID: "defend", Label: "Stand your ground", Detail: "Rely on what you are carrying and who is with you. Injuries and a weak defence can be fatal."},
			{ID: "bargain", Label: "Offer $180 to stand down", Cost: 180, Detail: "Money may settle this incident, but your standing suffers."},
		}}
}

// survived is the aftermath of an attack the player lived through.
func (w *World) survived() {
	p := &w.Player
	p.Health = max(1, p.Health-w.Absorb(65))
	w.Ruin(55)
	w.Damage(40)
	worn := "Nobody warned you."
	if p.Armour > 0 {
		worn = "What you were wearing took the worst of it."
	}
	w.Log("You survived by inches", fmt.Sprintf("They found you at %s and left you wounded. %s You need rest and protection.", w.whereItHappens(), worn), "danger")
}
