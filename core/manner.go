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
	"mercercourt": {"shot beside the letter boxes before the morning shift", "strangled on a stair landing while a radio played behind a closed door", "followed onto the roof and found beneath the washing lines", "beaten in the basement beside the coal bins"},
	"herald": {"shot in the alley behind the loading bay, between editions",
		"found at the foot of the stairwell with a note nobody could read",
		"beaten in the composing room after the last shift went home",
		"taken from the front steps into a car, in front of two reporters who wrote nothing"},
	"precinct": {"found in a cell in the morning, hanged with a shirt the desk sergeant had signed for",
		"beaten in a corridor between two rooms nobody keeps a book for",
		"shot on the steps outside, walking out on bail with the paperwork still in hand",
		"taken ill in custody, according to a doctor who saw them once"},
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
		"crushed between a hull and the dockside, which happens to people who work there",
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
	"pawn": {"shot over the counter with the shutter half down",
		"found in the back room among other people's things",
		"beaten in the doorway under the three brass balls",
		"taken out through the yard door and not seen again"},
	"tailor": {"shot in the fitting room with the curtain drawn",
		"found under the cutting table with the shears still on it",
		"beaten in the workroom among half a dozen half-made suits",
		"taken out through the back where the deliveries come in"},
	"chapel": {"shot in the back yard between the hearse and the wall",
		"found in a casket that was not due to be opened again",
		"beaten in the preparation room where the taps run all day",
		"carried out through the side door in front of a family who said nothing"},
	"scrapyard": {"put through the press with what was left of a Ford",
		"found four cars down the stack a week later",
		"shot between the rows where nobody walks",
		"burned with the rest of it and nobody counted the ashes"},
	"filling": {"shot between the pumps with the nozzle still in the tank",
		"found in the inspection pit under a car that was never booked in",
		"beaten behind the rack of oil cans, out of sight of the road",
		"taken off the forecourt into a car that did not stop for petrol"},
	"pumps": {"shot under the arches where the light does not reach",
		"found at the back of the lot among the drums",
		"beaten against the pumps at four in the morning with nobody passing",
		"driven off from the forecourt and put out somewhere on the county road"},
	"dealer": {"found in the boot of a car nobody had sold yet",
		"shot on the forecourt between two rows of them",
		"crushed under a car that came off its ramp at the wrong moment",
		"driven off the lot in the passenger seat and not seen again"},
	"archway": {"beaten under the viaduct where the trains cover everything",
		"found in the inspection pit in the morning",
		"shot at the roller door with it half open",
		"burned in a car that was already going to be burned"},
	"goldenlily": {"shot in the cloakroom with the coat still over their arm",
		"found in the river stairs below the terrace at first light",
		"beaten in the counting room with the night's take on the table",
		"walked out between two of them and into a car that did not wait"},
	"steamworks": {"held under in a vat until the steam stopped mattering",
		"crushed in the mangle at the end of a long shift",
		"shot between the drying racks where nobody could see the door",
		"found in the boiler room, and the boiler still running"},
	"burlesque": {"shot in the wings while the band went on",
		"strangled in a dressing room with the mirror lights on",
		"beaten at the stage door and left in the alley",
		"taken up the back stairs to talk and carried down the front"},
	"cabstand": {"shot through the windscreen with the engine idling",
		"found in the pit under a car that nobody had raised",
		"beaten in the dispatch hut with the radio still on",
		"driven out in their own cab and brought back without them"},
	"restaurant": {"shot at the corner table with the plates still down",
		"held face down in the sink until the kitchen went quiet",
		"knifed in the yard by the bins and left against the wall",
		"taken through the back room and out to a car that was waiting"},
	"poolhall": {"beaten with a cue in front of six tables of people who saw nothing",
		"shot under the low lamp over the far table",
		"strangled in the passage by the payphone that rings every night",
		"put through the plate window and left in it"},
	"butcher": {"found in the cold room in the morning, on a hook meant for something else",
		"shot behind the counter with the block still wet",
		"beaten in the yard among the empty crates",
		"taken out in the back of the van and delivered nowhere"},
	"haulage": {"crushed under a trailer nobody had chocked",
		"shot in the cab with the engine running",
		"beaten in the yard between two trucks, out of the light",
		"put on a lorry going north and taken off it somewhere else"},
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
