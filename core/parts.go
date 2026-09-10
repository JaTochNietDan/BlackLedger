package core

import "fmt"

// Cars belong to people now, which makes a car a thing that can be taken off
// somebody. A stripped car is four links in one act: whoever drove it is
// walking, the parts have a buyer, the buyer's trade picks up because there is
// more work about, and the forecourt sells that person another one in a week.
// That is what the city being joined up means.

const (
	// PartsValue is what a night's work off one car is worth, before what the
	// car was is taken into account.
	PartsValue = 95
	// PartsHeat is the attention it draws. Somebody's car going to pieces in
	// the street is not a quiet crime.
	PartsHeat = 5
	// PartsSore is what the owner holds against whoever did it.
	PartsSore = 35
	// PartsTrade is how much better a garage does for every car going to
	// pieces in this city. Trade is capped, so this cannot run away.
	PartsTrade = 3
	// WreckTrade is what a burned-out shell is worth to the yard that ends up
	// with it. It is worth nothing at all to a garage: a bench lives on parts
	// and glass, and there are no parts worth having on a car that went up in
	// the street.
	WreckTrade = 2
)

// StripTarget finds a car worth taking apart: somebody standing here, who the
// player knows, who drives, and who is not one of their own.
func (w *World) StripTarget(location string) (*NPC, bool) {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != location || n.Car == 0 || !w.Known(n) {
			continue
		}
		if n.Faction != "" && n.Faction == w.PlayerOrganizationID() {
			continue // not your own people's cars
		}
		if len(w.Player.Crew) > 0 && w.Player.Crew[0].ID == n.ID {
			continue
		}
		return n, true
	}
	return nil, false
}

// StripReadiness explains why there is nothing to take apart, or returns "".
func (w *World) StripReadiness(location string) string {
	if w.Player.HeldUntil > w.Minute {
		return "You are not going anywhere tonight"
	}
	if _, ok := w.StripTarget(location); !ok {
		return "Nothing parked here belongs to anybody you know"
	}
	if left := w.Curtains(location); left > 0 {
		place, _ := PlaceByID(location)
		return fmt.Sprintf("The street at %s is still being watched from the windows since the last time. Worth walking down again %s", place.Name, soon(left))
	}
	return ""
}

// PartsWorth is what one car is worth in pieces: what it was, and whether the
// player has a garage of their own to take them to rather than selling them on
// at whatever is offered.
func (w *World) PartsWorth(n *NPC) int {
	if n == nil || n.Car == 0 {
		return 0
	}
	worth := PartsValue * n.Car
	if w.OwnsKind("garage") {
		worth = worth * 3 / 2 // your own bench, and nobody taking a cut
	}
	return worth
}

// StripCar takes a car apart for what is on it. The owner is walking, the
// parts go to the trade that buys them, and every garage in the city has a
// little more work than it had yesterday.
func (w *World) StripCar(location string) error {
	if reason := w.StripReadiness(location); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	mark, _ := w.StripTarget(location)
	worth := w.PartsWorth(mark)
	// The street remembers. Nobody in a row that lost a car last night parks in
	// it again without looking, and somebody is at the window.
	if prop := w.Properties[location]; prop != nil {
		prop.Stripped = max(1, w.Minute)
	}
	label := VehicleByTier(mark.Car).Label
	mark.Car = 0

	w.Earn(worth)
	w.Player.Heat = min(100, w.Player.Heat+PartsHeat)
	// Against the PLAYER, not against whoever else was standing there. Resent
	// records a grudge between two people in the city; Aggrieve is the one that
	// reaches the protagonist, and it is what makes a man who lost his car the
	// one who comes looking for you later.
	w.Aggrieve(mark.ID, PartsSore, "what was left of their car")

	// The buyers. A garage is where parts go, so a city with more of them
	// going around is a city where a garage has more work — which is the whole
	// of why anybody would hold one.
	w.PartsAbout(PartsTrade)

	// Nobody takes one car apart in a quiet street and leaves the rest of the
	// row untouched. What is still there in the morning is still there with the
	// glass out of it, and that is a garage's actual trade.
	glass := w.BreakGlass(location, mark.ID)

	place, _ := PlaceByID(location)
	w.Log(label+" in pieces at "+place.Name,
		fmt.Sprintf("%s's car went for parts. $%d for the night's work, and %s will know by morning that it was somebody.",
			mark.Name, worth, mark.Name), "danger")
	if glass > 0 {
		w.Log("Glass in the road at "+place.Name,
			fmt.Sprintf("%s in the street were gone through and left where they stood. Somebody will be paid to put them right.", plural(glass, "more car", "more cars")), "danger")
	}
	w.Report("theft", "CAR STRIPPED IN "+upper(place.Name),
		w.unattributed(place.Name, fmt.Sprintf("A car was taken apart in the street at %s overnight. Garages in the district report no shortage of work.", place.Name)))
	return nil
}

// PartsAbout is what more cars going to pieces does to the trades that live off
// it. A garage takes the parts and a scrapyard takes what is left, so both do
// better when there is more of it happening — which is the whole of why anybody
// would hold one.
//
// Trade is capped, so this cannot run away however much of the city ends up on
// bricks.
func (w *World) PartsAbout(delta int) {
	for _, l := range Locations {
		switch l.Kind {
		case "garage":
			w.ShiftCustom(l.ID, "parts coming in off the street", delta)
		case "scrapyard":
			w.ShiftCustom(l.ID, "more of the city arriving on a low-loader", delta)
		}
	}
}

// WreckArrives is a car that is not coming back: burned in the street, or what
// is left of one after a raid. Somebody has to take it away, and the only trade
// in this city that wants it is the yard. A raid used to lift every garage and
// every yard by the same one point, which says a shell is worth what a stripped
// car is worth — and it is not, to anybody.
func (w *World) WreckArrives() {
	for _, l := range Locations {
		if l.Kind == "scrapyard" {
			w.ShiftCustom(l.ID, "another one in off a low-loader", WreckTrade)
		}
	}
}
