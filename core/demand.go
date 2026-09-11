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

// And when asking is no longer the point.
//
// A family can come for your deed once you have refused them long enough. The
// player could not do the same to them: a takeover reaches only a family you
// serve, and everything else aimed at a rival was a war. So a weak family that
// already hates you could go on holding a shop on your own street for ever,
// paying you a share every time you asked and losing nothing else.
//
// This is the other half of the claim, and it is deliberately harder than
// theirs. They need you to be weak on standing; you need them to be weak
// outright — beaten down below what your name is worth by a wide margin, at the
// bottom of their opinion of you, with somewhere else to go. Nobody is wiped out
// of this city in an afternoon.

const (
	// PushMinutes is how long it takes to say and be believed.
	PushMinutes = 90
	// PushMargin is how far your name has to be above what a family is before
	// they will walk out of a room rather than fight for it.
	PushMargin = 40
	// PushLeft is how many other places they have to hold. A family with one
	// address left is a family with nothing to lose, and taking the last thing
	// a family has is a different kind of trouble.
	PushLeft = 1
	// PushRespect and PushHeat are what the street and Ward Street make of it.
	PushRespect = 12
	PushHeat    = 14
	// PushPower is what it takes out of them.
	PushPower = 12
)

// Pushable reports whether this is a room worth even asking about — somebody
// else's, and held by a family that has stopped being on terms with the player.
//
// The card is drawn on that rather than on the whole rule, because a card that
// is refused in every room it ever appears in is a card nobody can press, and
// this game has a guard for exactly that. Below nothing they are on the path to
// walking out, and the refusal tells the player how far along it they are.
func Pushable(w *World, id string) bool {
	if !Leaning(w, id) {
		return false
	}
	f := w.faction(w.Properties[id].Owner)
	return f != nil && f.Goodwill < 0
}

// PushReadiness explains why a family will not walk out of this room, or "".
func (w *World) PushReadiness(id string) string {
	if w.Player.Location != id {
		return "This is said in the room"
	}
	if !Leaning(w, id) {
		return "There is nobody here to take it from"
	}
	f := w.faction(w.Properties[id].Owner)
	if f.Goodwill > DemandFloor {
		// upper1, because a family in this city can be called "the Duarte
		// Brothers" and a sentence that starts with one starts in lower case.
		return fmt.Sprintf("%s would fight you for it. They think of you at %+d and it takes %+d", upper1(f.Name), f.Goodwill, DemandFloor)
	}
	if w.Presence() < f.Power+PushMargin {
		return fmt.Sprintf("%s is worth more in this city than you are. You need %d presence against their %d",
			upper1(f.Name), f.Power+PushMargin, f.Power)
	}
	if len(w.FamilyHoldings(f.ID)) <= PushLeft {
		return upper1(f.Name) + " has nowhere else to go, and somebody who takes the last thing a family has is in a different kind of trouble"
	}
	return ""
}

// DemandFloor is the standing at or below which a family has stopped pretending
// to be on terms with the player, which is what it takes before they will give
// up a room rather than fight over it.
const DemandFloor = -60

// TakeItFromThem walks in and takes the deed.
func (w *World) TakeItFromThem(id string) error {
	if reason := w.PushReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	prop := w.Properties[id]
	f := w.faction(prop.Owner)
	place, _ := PlaceByID(id)
	prop.Owner = fmt.Sprintf("player:%d", w.Life)
	f.Power = max(0, f.Power-PushPower)
	f.Goodwill = -100
	w.Player.Respect += PushRespect
	w.Player.Heat = min(100, w.Player.Heat+PushHeat)
	// Nobody else in the city is pleased about it either. A man who takes a
	// room off one family is a man who might take one off another.
	for i := range w.Factions {
		if w.Factions[i].ID != f.ID && w.Factions[i].ID != w.PlayerOrganizationID() {
			w.Factions[i].Goodwill = max(-100, w.Factions[i].Goodwill-PushOthers)
		}
	}
	w.Log("They walked out of "+place.Name,
		fmt.Sprintf("%s is yours. Nobody raised a hand, because %s %s what it would cost and %s %s somewhere else to be. Every family in this city heard about it by the evening.",
			place.Name, Leads(f.Name), Agree(f.Name, "counted", "counted"), Leads(f.Name), Agree(f.Name, "has", "have")), "politics")
	w.RetaliationFrom(f.ID)
	return nil
}

// PushOthers is what the rest of the city makes of watching it happen.
const PushOthers = 10
