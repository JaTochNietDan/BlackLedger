package core

import "fmt"

// Burying one of your own.
//
// The city kills people constantly and the player's own crew are among them.
// What happened to a man who died working for you was a line in the log and
// nothing else: his name left the books, the wage stopped, and everybody else
// on the payroll went on as though the week had been ordinary. That is the one
// thing about this business nobody in it would let pass.
//
// So there is a funeral, and it is the undertaker's other half. The trade came
// in living on the city's dead and it earns on strangers; this is the part that
// costs the player money and buys the thing money does not usually buy, which
// is the confidence of people who can see what happens when it is their turn.

const (
	// FuneralWindow is how long a family waits. A man buried nine days after he
	// died was not buried by anybody who thought much of him.
	FuneralWindow = 6 * 1440
	// FuneralCost is what a decent one costs, and FuneralOwn what is left of
	// that when the parlour is yours: the cars and the box are your own, and
	// what remains is the plot and the notices.
	FuneralCost = 260
	FuneralOwn  = 90
	// FuneralTrust is what everybody still on the payroll makes of it, and
	// FuneralRespect what the street does.
	FuneralTrust   = 12
	FuneralRespect = 3
	// FuneralMinutes is a morning.
	FuneralMinutes = 180
)

// PauperTrust is what everybody still on the books makes of a man of theirs
// going into the ground on the parish's money.
//
// Without this the funeral was a card with no decision in it: pay and be
// thought better of, or do nothing and lose nothing. A thing that is free to
// skip is not a choice, and this file's own rule is to ask whether the new
// thing is strictly better than the old. It is the same six days from the other
// side — the morning after the last one is the morning everybody notices.
const PauperTrust = 8

// TheUnburied runs once a day and closes the window on anybody of the player's
// nobody arranged anything for.
func (w *World) TheUnburied() {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead || n.Buried || n.DiedAt <= 0 {
			continue
		}
		if n.Faction != w.PlayerOrganizationID() {
			continue
		}
		if w.Minute-n.DiedAt <= FuneralWindow {
			continue
		}
		// Marked either way. A man is only buried badly once.
		n.Buried = true
		left := w.OwnPeople()
		for _, other := range left {
			other.Trust = max(0, other.Trust-PauperTrust)
		}
		if len(left) == 0 {
			continue
		}
		w.Log("Nobody arranged anything for "+n.Name,
			fmt.Sprintf("They went into the ground on the parish's money, with nobody there from the firm they died working for. %s who still works for you knows it.",
				plainly(len(left), "The one person", fmt.Sprintf("Every one of the %d people", len(left)))), "danger")
	}
}

// FuneralFee is what burying one of yours costs here.
func (w *World) FuneralFee(at string) int {
	if w.Own(at) {
		return FuneralOwn
	}
	return FuneralCost
}

// Unburied is everybody of the player's who died recently and has not been
// buried, oldest first so the list does not shuffle between refreshes.
func (w *World) Unburied() []*NPC {
	out := []*NPC{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead || n.Buried || n.DiedAt <= 0 {
			continue
		}
		if n.Faction != w.PlayerOrganizationID() {
			continue
		}
		if w.Minute-n.DiedAt > FuneralWindow {
			continue
		}
		out = append(out, n)
	}
	return out
}

// FuneralReadiness explains why somebody cannot be buried, or returns "".
func (w *World) FuneralReadiness(at, id string) string {
	place, ok := PlaceByID(at)
	if !ok || place.Kind != "undertaker" {
		return "This is not a funeral director's"
	}
	if w.Player.Location != at {
		return "The arrangements are made at " + place.Name
	}
	n := w.NPC(id)
	if n == nil || !n.Dead || n.Buried {
		return "There is nobody of that description to bury"
	}
	if n.Faction != w.PlayerOrganizationID() {
		return "They were not one of yours"
	}
	if w.Minute-n.DiedAt > FuneralWindow {
		return "It has been too long. Whatever was going to be said has been said"
	}
	if w.Player.Cash < w.FuneralFee(at) {
		return "Not enough cash"
	}
	return ""
}

// BuryYourOwn pays for it. Everybody still on the payroll was watching.
func (w *World) BuryYourOwn(at, id string) error {
	if reason := w.FuneralReadiness(at, id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	fee := w.FuneralFee(at)
	if err := w.Pay(fee); err != nil {
		return err
	}
	n := w.NPC(id)
	n.Buried = true
	if !w.Own(at) {
		w.funeralProceeds(at, n.Name, fee)
	}
	// The trade is somebody else's unless it is yours, the same as every other
	// business in this city.
	w.ShiftCustom(at, "", BurialTrade)

	lifted := 0
	for _, other := range w.OwnPeople() {
		if other.ID == id {
			continue
		}
		other.Trust = min(100, other.Trust+FuneralTrust)
		lifted++
	}
	w.Player.Respect += FuneralRespect
	place, _ := PlaceByID(at)
	watched := "Nobody was left to see it."
	if lifted > 0 {
		watched = fmt.Sprintf("%s who still works for you knows what you do when it is one of them.",
			plainly(lifted, "The one person", "Every one of the "+fmt.Sprintf("%d people", lifted)))
	}
	w.Log("Buried out of "+place.Name,
		fmt.Sprintf("$%d for the cars, the box and the notices. %s", fee, watched), "politics")
	return nil
}
