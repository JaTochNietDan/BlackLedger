package core

import "fmt"

// Bellwether is not the only place on the map, and the things it cannot give
// you are worth days of your life to go and get. A journey out of town is a
// single decision with a long price: the city runs without you for days, your
// businesses go unwatched, work you promised somebody runs down toward its
// deadline, and whatever a rival had planned happens to an empty house.
//
// That last part is not a bug. Being somewhere else is a real way to survive a
// week you were not going to survive here, and it costs exactly as much as
// being away costs.

// Destination is somewhere the player can go, and the one thing it is for.
type Destination struct {
	ID, Name string
	// Blurb is the place, and Purpose the reason anybody makes the trip.
	Blurb, Purpose string
	// Days is how long the whole journey takes, and Fare what the ticket costs.
	Days int
	Fare int
	// Relief is the attention that falls off each day spent there. Somewhere
	// nobody has heard of you is worth more of it than a port two hours away.
	Relief int
}

var destinations = []Destination{
	{ID: "rockridge", Name: "Rockridge", Days: 2, Fare: 60, Relief: 4,
		Blurb:   "A mill town upriver where the stills outnumber the churches.",
		Purpose: "Buy moonshine at country prices and bring back everything you can carry."},
	{ID: "kingsport", Name: "Kingsport", Days: 3, Fare: 110, Relief: 5,
		Blurb:   "A port city with banks that ask fewer questions than the ones here.",
		Purpose: "Arrange the account abroad in person, which costs a fraction of arranging it by wire."},
	{ID: "halloway", Name: "Halloway", Days: 4, Fare: 90, Relief: 12,
		Blurb:   "Far enough inland that nobody there has heard of you, which is the point.",
		Purpose: "Be nowhere anybody is looking. Attention falls hard, and you come back knowing somebody new."},
}

const (
	// CountryPrice is what a crate costs where it is made, as a percentage of
	// what it costs in the city.
	CountryPrice = 55
	// PocketLoad is what a person can carry without a hiding place: enough to
	// make the trip worth the fare and not enough to make a car pointless.
	PocketLoad = 20
	// KingsportAccess is what establishing the account costs in person.
	KingsportAccess = 150
	// TripDepart is where journeys are booked.
	TripDepart = "market"
)

// DestinationByID is one of them, and whether it exists.
func DestinationByID(id string) (Destination, bool) {
	for _, d := range destinations {
		if d.ID == id {
			return d, true
		}
	}
	return Destination{}, false
}

// Destinations is everywhere the player could go, for the interface.
func Destinations() []Destination { return destinations }

// CarryLimit is everything the player could bring back: what fits in their
// pockets, under the floor of a car, and in a cellar at home.
func (w *World) CarryLimit() int { return PocketLoad + w.Concealed() }

// TripCost is the fare plus, for a buying trip, what the stock will cost. The
// player is told the whole number before they commit to four days away.
func (w *World) TripCost(id string) int {
	d, ok := DestinationByID(id)
	if !ok {
		return 0
	}
	switch id {
	case "rockridge":
		return d.Fare + w.countryPrice()*w.CarryLimit()
	case "kingsport":
		if !w.Player.Offshore && w.Offshore > 0 {
			return d.Fare + KingsportAccess
		}
	}
	return d.Fare
}

func (w *World) countryPrice() int {
	g := w.Good("moonshine")
	if g == nil {
		return 0
	}
	return max(1, g.Price*CountryPrice/100)
}

// TripReadiness explains why a journey cannot be made, or returns "".
func (w *World) TripReadiness(id string) string {
	d, ok := DestinationByID(id)
	if !ok {
		return "There is nowhere of that name"
	}
	if w.Player.Location != TripDepart {
		return "Journeys are booked at the exchange"
	}
	if w.Player.Health < 35 {
		return "You are in no condition to travel"
	}
	if id == "kingsport" && w.Player.Offshore && w.Offshore == 0 {
		return "There is nothing out there and nothing to arrange"
	}
	if id == "rockridge" && w.CarryLimit() == 0 {
		return "There would be nothing to carry it back in"
	}
	if w.Player.Cash < w.TripCost(id) {
		return fmt.Sprintf("The trip costs $%d all in", w.TripCost(id))
	}
	_ = d
	return ""
}

// Trip is the whole journey: the fare, the days, what happened in the city
// while the player was not in it, and what they came back with.
func (w *World) Trip(id string) error {
	if reason := w.TripReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	d, _ := DestinationByID(id)
	cost := w.TripCost(id)
	if err := w.Pay(cost); err != nil {
		return err
	}

	// What the trip was for, settled on arrival rather than on return, because
	// the days that follow are the city's rather than the player's.
	arrival := ""
	switch id {
	case "rockridge":
		load := w.CarryLimit()
		if w.Player.Stock == nil {
			w.Player.Stock = map[string]int{}
		}
		w.Player.Stock["moonshine"] += load
		arrival = fmt.Sprintf("%d crates at $%d apiece, which is what they cost where they are made. It all has to come back with you.", load, w.countryPrice())
	case "kingsport":
		if !w.Player.Offshore {
			w.Player.Offshore = true
			arrival = fmt.Sprintf("A morning in a panelled room and the account answers to you. It cost $%d rather than the $%d arranging it from here would have.", KingsportAccess, AccessCost)
		} else {
			arrival = "A morning in a panelled room, and the people there now know your face rather than your signature."
		}
	case "halloway":
		before := w.Player.Contacts
		w.Player.Contacts = min(5, w.Player.Contacts+1)
		arrival = "Four days where nobody knows your name is four days nobody is looking for it."
		if w.Player.Contacts > before {
			arrival += " You came back knowing somebody who owes you a conversation."
		}
	}
	w.Log("Out of the city: "+d.Name, fmt.Sprintf("$%d and %d days. %s", cost, d.Days, arrival), "travel")

	// The city runs. Anything that came for the player found an empty house,
	// which is the whole of what leaving buys and the reason it is worth doing.
	// Out of the city means out of the city: nothing that asks where the player
	// is standing can be satisfied while they are two hundred miles from it.
	w.Player.Location = "transit"
	dodged := 0
	relief := 0
	for day := 0; day < d.Days && w.Player.Alive; day++ {
		w.Advance(1440)
		if w.Event != nil {
			w.Event = nil
			dodged++
		}
		fell := min(w.Player.Heat, d.Relief)
		w.Player.Heat -= fell
		relief += fell
	}
	if dodged > 0 {
		w.Log("Somebody called while you were away", fmt.Sprintf("Whatever was arranged for you happened %d times to a locked door. Nobody who goes to that trouble takes it well.", dodged), "danger")
	}
	if relief > 0 {
		w.Log("Off the books for a while", fmt.Sprintf("Attention fell by %d while you were nowhere anybody could find you. It is now %d.", relief, w.Player.Heat), "personal")
	}
	if w.Player.Alive {
		w.Player.Location = TripDepart
		w.Log("Back in Bellwether", fmt.Sprintf("%d days gone. The city did not wait.", d.Days), "travel")
	}
	return nil
}
