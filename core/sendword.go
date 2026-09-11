package core

import "fmt"

const (
	WordMinutes  = 30
	WordCost     = 25
	WordHolds    = 480
	WordContacts = 2
)

func (w *World) Expecting(id string) bool {
	return w.Expected != nil && w.Expected[id] > w.Minute
}

func (w *World) SendWordReadiness(id string) string {
	f := w.faction(id)
	if f == nil || id == w.PlayerOrganizationID() {
		return "There is nobody of that description"
	}
	if w.Expecting(id) {
		return f.Name + " is already expecting you"
	}
	if w.Leader(id) == nil {
		return "Nobody is left to speak for " + f.Name
	}
	if len(w.FamilyHoldings(id)) == 0 {
		return f.Name + " has nowhere left to receive anybody"
	}
	if !w.Fitted("telephone") && w.Player.Contacts < WordContacts {
		return fmt.Sprintf("Nobody would carry a message for you. It takes a telephone or %d people who know your name", WordContacts)
	}
	if f.Goodwill < 0 {
		return f.Name + " would not take the call"
	}
	if w.Player.Cash < WordCost {
		return "Not enough cash"
	}
	return ""
}

func (w *World) SendWord(id string) error {
	if reason := w.SendWordReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(WordCost); err != nil {
		return err
	}
	if w.Expected == nil {
		w.Expected = map[string]int{}
	}
	w.Expected[id] = w.Minute + WordHolds
	f := w.faction(id)
	place, _ := PlaceByID(w.homeOf(id))
	how := "Somebody carries the message for you"
	if w.Fitted("telephone") {
		how = "You telephone"
	}
	w.sendFor(id)
	w.Log("Word sent to "+f.Name,
		fmt.Sprintf("%s. Somebody who can agree to something is making their way to %s and will be there for the rest of the day. $%d, and they know you want to talk before you are through the door.",
			how, place.Name, WordCost), "politics")
	return nil
}

// sendFor starts whoever can agree to something walking to their own hall.
//
// The seat is not conjured into the room: they go on their own feet, take the
// time it takes, and can be seen going, like everybody else on the street. It
// has to be done here rather than left to the next time the city reconsiders
// where people should be, because that happens twice a day and a message sent
// at nine in the morning is about this morning.
func (w *World) sendFor(id string) {
	lead, hall := w.Leader(id), w.homeOf(id)
	if lead == nil || hall == "" || lead.Location == hall {
		return
	}
	place, ok := PlaceByID(hall)
	if !ok {
		return
	}
	lead.Heading = hall
	lead.Errand = "expected at " + place.Name
	lead.Sets = w.Minute
	lead.Arrives = w.Minute + TravelMinutes(lead.Location, hall)
}
