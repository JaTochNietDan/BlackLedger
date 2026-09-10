package core

import "fmt"

// Putting a night on.
//
// A room's takings depend on who is standing in it now, which leaves an owner
// with an empty room and nothing to do about it. This is the oldest answer
// there is: pay for a band and a barrel, and the people who would have drunk
// somewhere else drink here instead.
//
// It moves actual people. A number that went up without anybody walking through
// the door would be the same abstraction the takings were before — and the
// promise the routine makes, that a person drinks in the same place every
// evening for the rest of their life, is kept: a night is one night, and
// afterwards they go back to where they always go.

const (
	// NightCost is the band, the barrel and the word going round.
	NightCost = 260
	// NightMinutes is the afternoon it takes to arrange.
	NightMinutes = 90
	// NightLasts is how long the word holds. One night, and the next evening
	// the city is back where it was.
	NightLasts = 1440
	// NightDraw is how far somebody will come for it. A person drinks near
	// where they are; a band is a reason to walk a bit further and not a reason
	// to cross the city.
	NightDraw = 40
)

// PlaysHost reports whether a place is somewhere a night can be put on: it has
// to be somewhere people drink.
func PlaysHost(id string) bool {
	for _, h := range haunts {
		if h == id {
			return true
		}
	}
	place, ok := PlaceByID(id)
	return ok && place.Type == "burlesque"
}

// NightOn reports whether a place is holding one tonight.
func (w *World) NightOn(id string) bool {
	prop := w.Properties[id]
	return prop != nil && prop.Night > w.Minute
}

// NightReadiness explains why a night cannot be put on, or returns "".
func (w *World) NightReadiness(id string) string {
	if !w.Own(id) {
		return "This is not a room of yours"
	}
	if !PlaysHost(id) {
		return "Nobody comes out for an evening at a place like this"
	}
	if w.NightOn(id) {
		return "There is one on already"
	}
	if w.Player.Cash < NightCost {
		return fmt.Sprintf("A band and a barrel costs $%d", NightCost)
	}
	return ""
}

// PutOnANight pays for it and puts the word about.
func (w *World) PutOnANight(id string) error {
	if reason := w.NightReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(NightCost); err != nil {
		return err
	}
	prop := w.Properties[id]
	prop.Night = w.Minute + NightLasts
	place, _ := PlaceByID(id)
	w.Log("A night at "+place.Name,
		fmt.Sprintf("$%d on a band, a barrel and the word going round. People who drink elsewhere will drink here tonight, and what a room takes is who is standing in it.", NightCost), "business")
	w.Report("business", "MUSIC AT "+upper(place.Name),
		fmt.Sprintf("%s is holding an evening tonight. Word has gone round the district.", place.Name))
	return nil
}

// theNight is where a night is on, and near enough to walk to from here.
func (w *World) theNight(from string) string {
	for _, l := range Locations {
		if !w.NightOn(l.ID) || l.ID == from {
			continue
		}
		if TravelMinutes(from, l.ID) <= NightDraw {
			return l.ID
		}
	}
	return ""
}
