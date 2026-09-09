package core

import "fmt"

// "Leo heads out" was the log line, and Leo did not head anywhere. He stood
// where he had been standing, the money arrived two hours later, and the city
// reported that he had come back from somewhere he never went.
//
// Everybody else in this city walks. The man the player pays was the last one
// who did not have to, and he is the one the player has most reason to watch.

// CollectionRound is where a round of collections is actually done: the
// player's best-earning premises, and failing that the corner the work comes
// from. Work with no address cannot be watched, and this game is about being
// able to see what your people are doing.
func (w *World) CollectionRound() string {
	best, income := "", 0
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || !w.Own(l.ID) || prop.Income <= 0 {
			continue
		}
		if worth := prop.Income * prop.Condition; worth > income {
			best, income = l.ID, worth
		}
	}
	if best != "" {
		return best
	}
	return "bar"
}

// SendOnCollections puts the crew member on the road to the round. He is out
// of whatever room he was in until he is back, which is what makes sending him
// a decision rather than a button that prints money.
func (w *World) SendOnCollections() {
	if len(w.Player.Crew) == 0 {
		return
	}
	n := w.NPC(w.Player.Crew[0].ID)
	if n == nil || n.Dead {
		return
	}
	round := w.CollectionRound()
	place, _ := PlaceByID(round)
	// Where he was standing, so the round can send him back to it.
	home := n.Location
	w.Tasks = append(w.Tasks, Task{ID(), n.Name + " · collections", w.Minute + CollectionMinutes})
	w.Homes = append(w.Homes, TaskHome{Task: w.Tasks[len(w.Tasks)-1].ID, Person: n.ID, Where: home})
	if round == n.Location {
		w.Log(n.Name+" starts the round", fmt.Sprintf("They are already at %s. Two hours of doors.", place.Name), "work")
		return
	}
	n.Heading = round
	n.Errand = "on collections at " + place.Name
	n.Arrives = w.Minute + TravelMinutes(n.Location, round)
	w.noticed(n, true)
	w.Log(n.Name+" heads out", fmt.Sprintf("%s, and %d minutes to walk it. They are not here while they are doing it.", place.Name, n.Arrives-w.Minute), "work")
}

// settleTasks pays for the work that has finished and sends whoever did it
// back where they set off from.
func (w *World) settleTasks() {
	for j := 0; j < len(w.Tasks); {
		if w.Tasks[j].Due > w.Minute {
			j++
			continue
		}
		done := w.Tasks[j]
		w.Earn(CollectionPay)
		w.sendHome(done.ID)
		w.Tasks = append(w.Tasks[:j], w.Tasks[j+1:]...)
	}
}

// sendHome walks whoever did a piece of work back to where they were standing
// when they were asked to do it.
func (w *World) sendHome(task string) {
	for i, home := range w.Homes {
		if home.Task != task {
			continue
		}
		w.Homes = append(w.Homes[:i], w.Homes[i+1:]...)
		n := w.NPC(home.Person)
		if n == nil || n.Dead || home.Where == "" {
			w.Log("The round is finished", fmt.Sprintf("$%d from collections.", CollectionPay), "business")
			return
		}
		if n.Location == home.Where || w.Travelling(n) {
			w.Log(n.Name+" returns", fmt.Sprintf("$%d from collections. They are available again.", CollectionPay), "business")
			return
		}
		place, _ := PlaceByID(home.Where)
		n.Heading = home.Where
		n.Errand = "coming back from the round"
		n.Arrives = w.Minute + TravelMinutes(n.Location, home.Where)
		w.noticed(n, true)
		w.Log(n.Name+" starts back", fmt.Sprintf("$%d from collections, and %d minutes back to %s.", CollectionPay, n.Arrives-w.Minute, place.Name), "business")
		return
	}
	w.Log("The round is finished", fmt.Sprintf("$%d from collections.", CollectionPay), "business")
}

const (
	// CollectionMinutes is how long a round of doors takes.
	CollectionMinutes = 120
	// CollectionPay is what it brings back.
	CollectionPay = 65
)

// onARound reports whether somebody is out doing a piece of work the player
// asked for, which is a different thing from whatever their job is called.
func (w *World) onARound(id string) bool {
	for _, home := range w.Homes {
		if home.Person == id {
			return true
		}
	}
	return false
}

const (
	// CourierPay and CourierRespect are what the first job in the game is
	// worth. The figure lived as a literal in four places — the payment, the
	// log line, the button and the guide — beside an unrelated 45 for how long
	// it takes, which is the shape of a number that eventually stops agreeing
	// with itself.
	CourierPay     = 45
	CourierRespect = 2
	// CourierMinutes is how long it takes, which is a different 45.
	CourierMinutes = 45
)
