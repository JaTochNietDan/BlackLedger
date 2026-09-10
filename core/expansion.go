package core

import "fmt"

// A family moving in on a shop that answers to nobody.
//
// Measured before this was written: counting only owners that resolve to a live
// organization, the families held nineteen percent of the city after two months
// and twenty-five percent after sixteen. Three quarters of it was "independent"
// — the placeholder for a shop nobody holds — and stayed that way for ever.
//
// Nothing in the rules ever took one. Property moved between families in a war,
// a splinter took one holding with it, and a family that had lost everything
// could start again on unheld ground. A family in good order never grew.
//
// Layer 1 of docs/LIVING_WORLD.md says organizations own income-earning
// property. Layer 4 says war changes what actors do and creates the openings a
// player exploits. Neither is true of a map nobody wants, and a city the player
// is competing for has to be a city somebody else is competing for.

const (
	// MoveInStrength is the power a family needs before it is worth their while
	// leaning on somebody who is not bothering anybody.
	MoveInStrength = 45
	// MoveInPrice is what it costs them: the payment, the persuasion, and the
	// week of somebody's time. It comes out of the family's own money, so a
	// family that has been bled cannot buy its way across the city.
	MoveInPrice = 320
	// MoveInPower is what holding one more address is worth to them.
	MoveInPower = 5
	// MoveInChance is how often, per family, per half-day, one of them tries.
	MoveInChance = .045
	// MoveInShare is the most of the city's earning addresses one organization
	// will hold by moving in. Past that they are not expanding, they are
	// winning, and a city with a winner is a city with nothing left to play.
	MoveInShare = 26
)

// earningAddresses is how many places in this city make money at all, which is
// the denominator for anybody's share of it.
func (w *World) earningAddresses() int {
	n := 0
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Income > 0 {
			n++
		}
	}
	return n
}

// Unheld reports whether this address earns and answers to nobody: not the
// player's, not a family's, and not a dead protagonist's estate.
func (w *World) Unheld(id string) bool {
	prop := w.Properties[id]
	if prop == nil || prop.Income <= 0 || w.Own(id) {
		return false
	}
	if w.faction(prop.Owner) != nil {
		return false
	}
	if prop.Owner == "independent" || prop.Owner == "" {
		return true
	}
	// An id that no longer resolves to anybody: a family that has been buried
	// still has its name on the doors it held, and those doors are standing
	// open. Not a dead protagonist's estate, which has its own meaning and its
	// own way of changing hands.
	return !hasPrefix(prop.Owner, "player:") && !hasPrefix(prop.Owner, "former:")
}

// ConsiderExpansion is a family in good order taking somewhere that answers to
// nobody. Called from the faction turn, off the world's own stream, because the
// player is not party to it.
func (w *World) ConsiderExpansion() {
	total := w.earningAddresses()
	if total == 0 {
		return
	}
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.ID == w.PlayerOrganizationID() {
			continue
		}
		held := len(w.FamilyHoldings(f.ID))
		// A family with nothing is a different case and has its own rule:
		// considerReestablish gets them back on their feet for nothing, because
		// they have nothing to pay with.
		if held == 0 || f.Power < MoveInStrength || f.Cash < MoveInPrice {
			continue
		}
		if held*100/total >= MoveInShare {
			continue
		}
		if w.WorldRandom() >= MoveInChance {
			continue
		}
		// The nearest thing to their own ground, so a family grows outward from
		// where it already is rather than appearing across the city.
		seat := w.homeOf(f.ID)
		best, closest := "", 1<<30
		for _, l := range Locations {
			if !w.Unheld(l.ID) {
				continue
			}
			if away := TravelMinutes(seat, l.ID); away < closest {
				best, closest = l.ID, away
			}
		}
		if best == "" {
			return // nothing left in this city that answers to nobody
		}
		prop := w.Properties[best]
		f.Cash -= MoveInPrice
		prop.Owner = f.ID
		f.Power = min(peak(f), f.Power+MoveInPower)
		place, _ := PlaceByID(best)
		w.Log("Somebody else's name over the door",
			fmt.Sprintf("%s %s taken over %s. Whoever was running it is still behind the counter and is not the one being paid now.",
				Leads(f.Name), Agree(f.Name, "has", "have"), place.Name), "politics")
		// The headline agrees with the name too: "Vera Kohl's people" is a
		// plural, and the guard that reads every passage in the game caught
		// TAKES where it should have been TAKE.
		w.Report("business", upper(f.Name)+" "+upper(Agree(f.Name, "takes", "take"))+" OVER "+upper(place.Name),
			fmt.Sprintf("%s %s acquired the running of %s. The premises had been trading independently.",
				Leads(f.Name), Agree(f.Name, "has", "have"), place.Name))
		return // one a turn, so the city changes hands at a pace anybody can watch
	}
}
