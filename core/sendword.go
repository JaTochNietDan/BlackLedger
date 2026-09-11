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

// TheExpected names the family that said it would have somebody at this address
// today, if any did. A seat arranged by telephone outranks whoever the room
// happened to put in front of you.
func (w *World) TheExpected(location string) string {
	for id, until := range w.Expected {
		if until > w.Minute && w.homeOf(id) == location && w.Leader(id) != nil {
			return id
		}
	}
	return ""
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
	w.Log("Word sent to "+f.Name,
		fmt.Sprintf("%s and somebody who can agree to something will be at %s for the rest of the day. $%d, and they know you want to talk before you are through the door.",
			how, place.Name, WordCost), "politics")
	return nil
}
