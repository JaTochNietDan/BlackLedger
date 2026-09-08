package core

import "fmt"

// A killing was a sentence the code wrote: killed at a place, by an
// organization, over something. It said what happened and never how, so twelve
// different deaths in a campaign read as the same death twelve times.
//
// The manner of it is composed from what the world already knows: where the
// person was standing, what hour it was, and the kind of person who came for
// them. Nothing here is invented at the moment of telling — it is all state the
// simulation committed before anybody died.

// mannerByPlace is how it is done, by where the person was. The detail is the
// point: a hit that is only a number is not a story.
// mannerByPlace is how it is done, by where the person was. The place and the
// hour are stated first, so an entry here is free to end on a clause without
// the sentence falling over — the first build appended the place after them and
// produced "walked twenty feet before anybody caught him at The Monarch".
//
// Nobody is a he or a she in these. The city does not know and it does not
// need to.
var mannerByPlace = map[string][]string{
	"bar": {"shot twice at the counter, in front of everyone and nobody",
		"stabbed in the passage behind the building and left where the bins are",
		"held down over a table until the room emptied",
		"followed out and put down against the wall by the door"},
	"club": {"shot in the doorway while the band kept playing",
		"beaten to death in the office upstairs",
		"knifed in the crowd, and got twenty feet before anybody noticed",
		"taken out through the kitchen and never brought back"},
	"market": {"shot at the loading doors, early, before the traders came",
		"found under a tarpaulin with a wire still around the throat",
		"run down by a lorry that neither stopped nor slowed",
		"beaten in a storeroom with the door shut on it"},
	"docks": {"put into the water with pockets full of chain",
		"shot on the quay and rolled off it",
		"crushed between a hull and the dockside, which happens to men who work there",
		"taken onto a boat that came back without them"},
	"laundry": {"held under in a press until it stopped",
		"shot through the window from a car that did not stop",
		"found in the back room with the machines still running",
		"strangled with something off the line and left sitting up"},
	"garage": {"crushed under a car that was not on a jack by accident",
		"shot in the pit, twice, close",
		"burned in a car nobody had reported stolen",
		"beaten with a wrench in full view of the road"},
	"casino": {"shot on the steps, walking out with a good night behind them",
		"found in a service corridor with no marks worth reporting",
		"knifed at the tables while the room looked at the wheel",
		"taken upstairs to talk and carried back down"},
	"apartment": {"shot on their own landing",
		"strangled in the stairwell",
		"put through a window from four floors up",
		"met in their own kitchen by somebody who was already sitting in it"},
	"estate": {"shot on the drive, getting out of the car",
		"found in the grounds two days later",
		"drowned in the pond, which the police called an accident twice",
		"shot through a window from the lawn"},
	"room": {"shot through the door of their own room",
		"smothered in bed",
		"pushed down the stairs by somebody waiting on the landing",
		"found in the bath with the water still warm"},
}

// hourOf is the time of day in the words a newspaper would use.
func hourOf(minute int) string {
	switch hour := (minute % 1440) / 60; {
	case hour < 5:
		return "in the small hours"
	case hour < 9:
		return "before the streets were properly awake"
	case hour < 12:
		return "in the middle of the morning"
	case hour < 15:
		return "at lunchtime, in daylight"
	case hour < 18:
		return "in the late afternoon"
	case hour < 22:
		return "after dark"
	}
	return "late, with the city quiet"
}

// signature is what the kind of person who did it leaves behind. It is the one
// thing about a killing that says who to look for.
func signature(killer *NPC) string {
	if killer == nil {
		return ""
	}
	switch TemperamentOf(killer).ID {
	case "hot":
		return " It was not quiet and it was not quick."
	case "careful":
		return " Nobody heard it and nobody has said anything since."
	case "greedy":
		return " Their pockets were empty when they were found."
	case "vain":
		return " Whoever did it wanted it seen."
	case "loyal":
		return " It was done properly, the way these things are supposed to be."
	}
	return ""
}

// Manner is how somebody died, in one sentence, composed from where they were,
// when it was and who came for them.
func (w *World) Manner(victim, killer *NPC) string {
	if victim == nil {
		return ""
	}
	place, known := PlaceByID(victim.Location)
	options := mannerByPlace[victim.Location]
	if !known || len(options) == 0 {
		return fmt.Sprintf("In the street, %s: %s was killed, and the street has said nothing about it.", hourOf(w.Minute), victim.Name) + signature(killer)
	}
	how := options[int(w.WorldRandom()*float64(len(options)))%len(options)]
	return fmt.Sprintf("At %s, %s: %s was %s.", place.Name, hourOf(w.Minute), victim.Name, how) + signature(killer)
}

// KillBy ends a life and says how. The manner comes first because it is what
// anybody in the city would hear; the reason follows for whoever knows it.
func (w *World) KillBy(id string, killer *NPC, why string) bool {
	victim := w.NPC(id)
	if victim == nil || victim.Dead {
		return false
	}
	cause := w.Manner(victim, killer)
	if why != "" {
		cause += " " + why
	}
	return w.Kill(id, cause)
}
