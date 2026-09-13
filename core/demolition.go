package core

import (
	"fmt"
	"strings"
)

// Sabotage is a crew breaking things in the dark and everybody pretending
// afterwards that they do not know who. A charge under the floor is not that.
// It is unmistakable, it is in the paper, it kills people who happened to be
// standing there, and it is the difference between a message and a declaration.
//
// It goes both ways, which is the whole of the design. Organizations at war
// use them on each other, an organization that hates the player enough uses
// them on the player, and none of it needs the player to be present.

const (
	// ChargeCost is what a charge costs off a boat. It is deliberately dearer
	// than the best gun in the city.
	ChargeCost = 850
	// ChargeMinutes is how long acquiring one takes.
	ChargeMinutes = 60
	// ChargeHeat is the attention each charge in your possession draws every
	// day. Nothing else in the game is this expensive to simply hold.
	ChargeHeat = 4
	// PlantMinutes is how long it takes to place one properly.
	PlantMinutes = 120
	// PlantStanding is what a room's worth of nerve costs: nobody with no name
	// gets near the underside of a rival's floor.
	PlantStanding = 20
	// BlastCasualty is the chance a charge kills somebody who worked there.
	BlastCasualty = .35
	// FactionChargeCost is what an organization spends to do the same thing.
	FactionChargeCost = 1200
)

// Charges is what the player is holding.
func (w *World) Charges() int { return w.Player.Charges }

// ChargeReadiness explains why one cannot be bought, or returns "".
func (w *World) ChargeReadiness() string {
	if !ArmsSource(w.Player.Location) {
		return "Nobody sells this here"
	}
	if w.Player.Charges >= 2 {
		return "You are carrying as much of this as anybody sane carries"
	}
	if w.Player.Cash < ChargeCost {
		return "Not enough cash"
	}
	return ""
}

// BuyCharge acquires one. Holding it is the most conspicuous thing a person in
// this city can do.
func (w *World) BuyCharge() error {
	if reason := w.ChargeReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(ChargeCost); err != nil {
		return err
	}
	w.Player.Charges++
	w.Player.Heat = min(100, w.Player.Heat+6)
	w.Log("A crate off a boat at Pier 14", fmt.Sprintf("$%d for something nobody in this city sells twice to the same customer. You are holding %d. It draws %d attention a day and a search that finds it ends everything.", ChargeCost, w.Player.Charges, ChargeHeat), "danger")
	return nil
}

// ChargeDay is what holding explosives costs while you decide what to do with
// them, which is the reason not to buy one until you know.
func (w *World) ChargeDay() {
	if w.Player.Charges <= 0 {
		return
	}
	w.Player.Heat = min(100, w.Player.Heat+ChargeHeat*w.Player.Charges)
}

// SeizeCharges is what a search costs somebody carrying explosives. Handled
// apart from arms because this is not a fine, it is a case.
func (w *World) SeizeCharges() bool {
	if w.Player.Charges <= 0 {
		return false
	}
	w.Player.Charges = 0
	w.Player.Heat = min(100, w.Player.Heat+25)
	w.Log("They found the crate", "What you were carrying is not something anybody talks their way out of. It is gone, and so is any doubt about what you are.", "danger")
	return true
}

// PlantReadiness explains why a charge cannot be placed here, or returns "".
func (w *World) PlantReadiness(id string) string {
	if w.Player.Charges <= 0 {
		return "You have nothing to place"
	}
	prop := w.Properties[id]
	if prop == nil || prop.Income <= 0 {
		return "There is nothing here worth the trouble"
	}
	if w.Own(id) {
		return "This is yours"
	}
	if w.Player.Location != id {
		return "You would have to be there"
	}
	if w.Presence() < PlantStanding {
		return fmt.Sprintf("Nobody who is nobody gets under that floor. You need %d presence", PlantStanding)
	}
	return ""
}

// Plant places a charge. It is resolved immediately and completely: what
// happens to the building, to whoever was in it, and to the player's standing
// with the people who owned it.
func (w *World) Plant(id string) error {
	if reason := w.PlantReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	prop := w.Properties[id]
	place, _ := PlaceByID(id)
	owner := w.faction(prop.Owner)
	w.Player.Charges--

	// Getting to the underside of somebody else's floor is the hard part.
	odds := .4 + float64(min(w.Presence(), 100))/300
	if len(w.Player.Crew) > 0 && w.Player.Crew[0].Loyalty >= 40 && len(w.Tasks) == 0 {
		odds += .12
	}
	if owner != nil {
		odds -= float64(owner.Power) / 280
	}
	odds = min64(.85, max64(.15, odds))

	if w.Random() >= odds {
		// It goes off in your hands, or somebody finds it with your face still
		// fresh in their memory.
		injury := w.Absorb(30 + int(w.Random()*35))
		w.Ruin(70)
		w.Damage(45)
		w.Player.Health = max(0, w.Player.Health-injury)
		w.Player.Heat = min(100, w.Player.Heat+30)
		if owner != nil {
			owner.Goodwill = max(-100, owner.Goodwill-45)
			w.RetaliationFrom(owner.ID)
		}
		w.Log("It went off early at "+place.Name, fmt.Sprintf("Something was wrong with it, or with the hour you chose. You are burned and cut (-%d health) and half the street saw somebody running.", injury), "danger")
		w.Report("attack", "EXPLOSION AT "+strings.ToUpper(place.Name),
			fmt.Sprintf("An explosion at %s is being treated as deliberate. Witnesses described somebody leaving on foot. Police say a prosecution is likely.", place.Name))
		w.Witness("explosion", id, "A charge went off early at "+place.Name+".", "EXPLOSION AT "+upper(place.Name))
		if w.Player.Health <= 0 {
			w.DieOf("a charge of your own", "A charge at "+place.Name+" went off with you still under it.")
		}
		return nil
	}

	w.detonate(id, fmt.Sprintf("A charge went off under %s.", place.Name))
	w.Player.Heat = min(100, w.Player.Heat+22)
	w.Player.Respect += 8
	if owner != nil {
		owner.Goodwill = max(-100, owner.Goodwill-50)
		w.RetaliationFrom(owner.ID)
		w.Log("There is no mistaking it", fmt.Sprintf("%s is wreckage. %s will not be wondering whether it was deliberate, and their standing with you is %+d.", place.Name, owner.Name, owner.Goodwill), "politics")
	}
	return nil
}

// detonate is the whole of what a charge does to a place, whoever set it. The
// same function for the player and for an organization, so a bombing is a
// bombing whichever end of it you are on.
func (w *World) detonate(id, cause string) {
	prop := w.Properties[id]
	place, _ := PlaceByID(id)
	if prop == nil {
		return
	}
	damage := min(prop.Condition, 45+int(w.WorldRandom()*30))
	prop.Condition -= damage
	// A wrecked business does not trade. Stock, staff and the money behind the
	// tables go with the building.
	//
	// And the position that goes is a person who goes. This took the count down
	// and left the name on the books, so a bombed club named five people behind
	// a counter it said held four — which is the exact thing `Property.Hands`
	// was introduced to stop, a staff that is a number rather than people. A
	// counter that loses a position loses whoever was standing at it, and the
	// room says so, because somebody not coming back to work is not a statistic
	// the player should have to infer from a figure moving.
	prop.Supply = 0
	if prop.Staff > 0 {
		prop.Staff--
		if gone := w.lastHand(id); gone != "" {
			w.letGo(id)
			if n := w.NPC(gone); n != nil {
				w.Log("One of them is not coming back", fmt.Sprintf("%s was at %s when it went up, and will not be standing behind that counter again.", n.Name, place.Name), "danger")
			}
		}
	}
	prop.Trouble = true
	if prop.Bankroll > 0 {
		prop.Bankroll = prop.Bankroll / 3
	}
	if prop.Still {
		prop.Still = false
	}
	owner := w.faction(prop.Owner)
	if owner != nil {
		owner.Power = max(10, owner.Power-max(4, damage/4))
		owner.Cash = max(0, owner.Cash-damage*30)
	}

	killed := ""
	if w.WorldRandom() < BlastCasualty {
		// A blast can only kill somebody physically inside these premises.
		// Choosing from the entire owning family killed people across town and
		// sent the city camera to an unrelated address while describing this one.
		present := []*NPC{}
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if !n.Dead && n.Location == id && !w.Travelling(n) {
				present = append(present, n)
			}
		}
		if len(present) > 0 {
			victim := present[int(w.WorldRandom()*float64(len(present)))]
			killed = victim.Name
			w.Kill(victim.ID, fmt.Sprintf("%s was inside %s when a charge went off under it, %s.", victim.Name, place.Name, hourOf(w.Minute)))
		}
	}

	body := fmt.Sprintf("%s The premises are wreckage and will not trade for some time.", cause)
	if killed != "" {
		body = fmt.Sprintf("%s %s was inside. The premises are wreckage.", cause, killed)
	}
	w.Log("An explosion at "+place.Name, body, "danger")
	headline := "EXPLOSION DESTROYS " + strings.ToUpper(place.Name)
	if killed != "" {
		headline = "ONE DEAD IN EXPLOSION AT " + strings.ToUpper(place.Name)
	}
	w.Report("attack", headline,
		fmt.Sprintf("An explosion at %s is being treated as deliberate. %s Police have appealed for witnesses and say they expect none.", place.Name, body))
	w.Witness("explosion", id, fmt.Sprintf("An explosion wrecked %s. Condition is now %d%%.", place.Name, prop.Condition), "EXPLOSION AT "+upper(place.Name))
}

func max64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// DemolitionDay is the city using the same thing the player can. An
// organization at war with money to spend will occasionally stop raiding a
// rival's premises and simply destroy them, and one that hates the player
// enough will do it to theirs.
func (w *World) DemolitionDay() {
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.Cash < FactionChargeCost {
			continue
		}
		target, reason := w.demolitionTarget(f)
		if target == "" {
			continue
		}
		// Rare on purpose. Anything more often than this and wreckage stops
		// being an event and becomes weather.
		if w.WorldRandom() >= .012 {
			continue
		}
		f.Cash -= FactionChargeCost
		w.detonate(target, reason)
	}
}

// demolitionTarget is what an organization would destroy, and why. A war is the
// usual reason; the player having made an enemy of them is the other.
func (w *World) demolitionTarget(f *Faction) (string, string) {
	// Somebody they are at war with, and the best thing that side holds.
	if rival := w.fighting(f.ID); rival != nil {
		best := ""
		for _, id := range w.FamilyHoldings(rival.ID) {
			if best == "" || w.Properties[id].Income > w.Properties[best].Income {
				best = id
			}
		}
		if best != "" {
			return best, fmt.Sprintf("A charge went off under it during the war between %s and %s.", f.Name, rival.Name)
		}
	}
	// Or the player, if they have given this organization enough reason.
	if f.Goodwill <= -60 {
		best := ""
		for _, l := range Locations {
			if prop := w.Properties[l.ID]; prop != nil && w.Own(l.ID) && prop.Income > 0 {
				if best == "" || prop.Income > w.Properties[best].Income {
					best = l.ID
				}
			}
		}
		if best != "" {
			return best, fmt.Sprintf("A charge went off under it. %s did not send anybody to explain why.", f.Name)
		}
	}
	return "", ""
}
