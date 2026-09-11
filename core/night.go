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
	// NightLasts is the evening itself: half a day, from the hour the crowd
	// sets off to midnight. One night, and the next evening the city is back
	// where it was.
	NightLasts = 720
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

// nightFrom is the hour the crowd sets off for the evening a night arranged
// at this minute pays for. Word that goes round before the crowd has left is
// word for tonight; word that goes round after it has is word for tomorrow,
// because tonight is already decided and the room would be paying for an
// evening that had begun without it.
func nightFrom(minute int) int {
	day := minute - minute%1440
	if minute%1440 <= EveningFrom {
		return day + EveningFrom
	}
	return day + 1440 + EveningFrom
}

// NightOn reports whether a place is holding one this evening. It is a window
// and not a deadline: a deadline cannot say which evening was bought, and the
// two differ whenever the afternoon spent arranging one runs past the hour the
// city goes out.
func (w *World) NightOn(id string) bool {
	prop := w.Properties[id]
	return prop != nil && prop.Night <= w.Minute && w.Minute < prop.Night+NightLasts
}

// nightBooked reports whether one is on or paid for and still to come.
func (w *World) nightBooked(id string) bool {
	prop := w.Properties[id]
	return prop != nil && prop.Night > 0 && w.Minute < prop.Night+NightLasts
}

// NightReadiness explains why a night cannot be put on, or returns "".
func (w *World) NightReadiness(id string) string {
	if !w.Own(id) {
		return "This is not a room of yours"
	}
	if !PlaysHost(id) {
		return "Nobody comes out for an evening at a place like this"
	}
	if w.nightBooked(id) {
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
	prop.Night = nightFrom(w.Minute + NightMinutes)
	place, _ := PlaceByID(id)
	when := "tonight"
	if prop.Night >= w.Minute+1440 {
		when = "tomorrow evening"
	}
	w.Log("A night at "+place.Name,
		fmt.Sprintf("$%d on a band, a barrel and the word going round. People who drink elsewhere will drink here %s, and what a room takes is who is standing in it.", NightCost, when), "business")
	w.Report("business", "MUSIC AT "+upper(place.Name),
		fmt.Sprintf("%s is holding an evening %s. Word has gone round the district.", place.Name, when))
	return nil
}

// NightPull is how much of a room two doors down a band takes, in people per
// hundred. The rest of the hundred is lost across the walk: somebody at the far
// edge of the draw mostly stays where they are.
const NightPull = 80

// comes decides whether one person walks the extra distance, from their id and
// nothing else, so the same night draws the same faces in a replayed game.
// Nearness is the whole of it. A band took every single drinker out of both
// rooms inside the radius and left them dark, which is not a night out, it is a
// switch; a draw that falls off with the walk leaves the near room thin and the
// far room barely touched, and makes the radius a pull rather than a cliff.
func comes(id string, minutes int) bool {
	share := NightPull - NightPull*minutes/(2*NightDraw)
	sum := 0
	for i := 0; i < len(id); i++ {
		sum = sum*37 + int(id[i])
	}
	if sum < 0 {
		sum = -sum
	}
	return sum%100 < share
}

// theNight is where a night is on, near enough to walk to, nearest of them if
// more than one room has paid for a band, and worth the walk to this person.
//
// It is asked from the room somebody would have drunk in, not from the desk
// they are standing at when the question comes up. A night takes trade off the
// room that would have had it, so the distance that matters is between the two
// rooms. Measured from the desk instead, a district with a band in it drew
// nobody at all whenever the day's work happened to be across town.
func (w *World) theNight(from, who string) string {
	best, near := "", 0
	for _, l := range Locations {
		if !w.NightOn(l.ID) || l.ID == from {
			continue
		}
		if d := TravelMinutes(from, l.ID); d <= NightDraw && (best == "" || d < near) {
			best, near = l.ID, d
		}
	}
	if best == "" || !comes(who, near) {
		return ""
	}
	return best
}
