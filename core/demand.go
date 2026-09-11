package core

import "fmt"

// Leaning on somebody else's business.
//
// A family that holds ground on your streets sends for a share of what your
// places take, and pays for an understanding or comes for the deed. The player
// could do none of that. Everything aimed at a family was either violence — go
// in with your crew, plant a charge, put a name up for a contract — or politics
// at a distance: point one of them at another, buy an official, smear a name in
// the paper. There was no version of the thing the whole genre is about, which
// is walking into a room that is not yours and saying that it is now.
//
// This is that, and it is deliberately the mirror of what they do to you. The
// same share of the same takings, refused or paid for the same reason: whether
// the people in the room think you can make it stick.

const (
	// DemandMinutes is the afternoon it takes.
	DemandMinutes = 60
	// DemandStanding is the presence below which nobody takes you seriously
	// enough to open the till.
	DemandStanding = 30
	// DemandGrudge is what being leaned on costs their opinion of you, whether
	// or not they pay. Asking is the insult; paying is only the money.
	DemandGrudge = 18
	// DemandRespect is what the street makes of somebody who walked in and was
	// paid. Nothing at all if they were refused: it is the being paid that is
	// worth saying.
	DemandRespect = 3
)

// Leaning reports whether this address is somebody else's business worth
// leaning on: a family's, earning, and not the player's.
func Leaning(w *World, id string) bool {
	prop := w.Properties[id]
	if prop == nil || prop.Income <= 0 || w.Own(id) {
		return false
	}
	return w.faction(prop.Owner) != nil
}

// WouldPay reports whether the family holding this place opens the till, and is
// decided from what each side is rather than from a roll. A family stronger
// than your name in the city tells you to leave; one weaker than it pays and
// remembers.
func (w *World) WouldPay(id string) bool {
	prop := w.Properties[id]
	f := w.faction(prop.Owner)
	return f != nil && w.Presence() > f.Power
}

// DemandReadiness explains why a share cannot be asked for here, or "".
func (w *World) DemandReadiness(id string) string {
	if w.Player.Location != id {
		return "This is said in the room"
	}
	if !Leaning(w, id) {
		return "There is nobody here to ask"
	}
	if w.Presence() < DemandStanding {
		return fmt.Sprintf("Nobody here has heard of you. You need %d presence", DemandStanding)
	}
	return ""
}

// DemandAShare walks in and says it.
func (w *World) DemandAShare(id string) error {
	if reason := w.DemandReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	prop := w.Properties[id]
	f := w.faction(prop.Owner)
	place, _ := PlaceByID(id)
	share := w.TheirShare(id)
	// Asking is the insult, whatever comes of it.
	f.Goodwill = max(-100, f.Goodwill-DemandGrudge)
	if !w.WouldPay(id) {
		w.Log("They told you to leave",
			fmt.Sprintf("%s at %s. %s %s what you are worth against what %s %s, and %s %s not worried. Their standing with you is %+d.",
				f.Name, place.Name, Leads(f.Name), Agree(f.Name, "weighed", "weighed"), Agree(f.Name, "has", "have"),
				"behind them", Leads(f.Name), Agree(f.Name, "is", "are"), f.Goodwill), "politics")
		w.RetaliationFrom(f.ID)
		return nil
	}
	paid := min(share, f.Cash)
	f.Cash = max(0, f.Cash-paid)
	w.Earn(paid)
	w.Player.Respect += DemandRespect
	w.Log("A share of "+place.Name,
		fmt.Sprintf("$%d out of %s's till, which is a week of what the place takes. Nobody argued in front of the room. Their standing with you is %+d.",
			paid, f.Name, f.Goodwill), "politics")
	return nil
}
