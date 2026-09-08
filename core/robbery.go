package core

import (
	"fmt"
	"strings"
)

// Taking something that is not yours, by hand, in the open. It is the crudest
// money in the city and the least deniable: a robbery is witnessed, and whoever
// owned what was taken knows it happened even when they do not know who.
//
// It runs both ways. The people who rob the player are named people with a
// place in the city, not anonymous thieves, and they can be answered.

// RobberyReadiness explains why premises cannot be robbed, or returns "".
func (w *World) RobberyReadiness(id string) string {
	prop := w.Properties[id]
	if prop == nil || prop.Income <= 0 {
		return "There is nothing here worth taking"
	}
	if w.Own(id) {
		return "You would be robbing yourself"
	}
	if w.Player.Health < 40 {
		return "You are in no condition for this"
	}
	return ""
}

// robberyOdds. A crew helps, a reputation helps, and premises belonging to a
// strong organization are watched.
func (w *World) robberyOdds(id string, hand Hand) float64 {
	odds := .5 + w.HandEdge(hand)
	if !hand.Crew && len(w.Player.Crew) > 0 && w.Player.Crew[0].Loyalty >= 40 && len(w.Tasks) == 0 {
		odds += .15 // somebody at your shoulder
	}
	if f := w.faction(w.Properties[id].Owner); f != nil {
		odds -= float64(f.Power) / 300
	}
	if odds < .15 {
		odds = .15
	}
	if odds > .85 {
		odds = .85
	}
	return odds
}

// Rob takes the day's cash out of premises the player does not own, by the
// player's own hands.
func (w *World) Rob(id string) error { return w.RobBy(id, w.OwnHands()) }

// RobBy is the same robbery whoever is standing there.
func (w *World) RobBy(id string, hand Hand) error {
	if reason := w.RobberyReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if hand.Crew {
		if reason := w.DelegateReadiness(); reason != "" {
			return fmt.Errorf("%s", reason)
		}
	}
	prop := w.Properties[id]
	place, ok := PlaceByID(id)
	if !ok {
		return fmt.Errorf("unknown premises")
	}
	owner := w.faction(prop.Owner)

	if w.Random() >= w.robberyOdds(id, hand) {
		injury := 10 + int(w.Random()*20)
		health := w.Player.Health
		w.HandHurt(hand, injury, "take the till at "+place.Name)
		w.Player.Heat = min(100, w.Player.Heat+w.HandHeat(hand, 15))
		if owner != nil {
			owner.Goodwill = max(-100, owner.Goodwill-15)
			w.RetaliationFrom(owner.ID)
		}
		if !hand.Crew {
			w.Log("It went wrong at "+place.Name, fmt.Sprintf("Somebody was waiting, or somebody shouted. You left with nothing and took a beating for it (-%d health).", health-w.Player.Health), "danger")
		} else {
			w.Log("It went wrong at "+place.Name, fmt.Sprintf("Somebody was waiting. %s left with nothing, and you were not there to be seen.", hand.Name), "danger")
		}
		if w.Player.Health <= 0 {
			w.Die("A robbery at " + place.Name + " went wrong.")
		}
		return nil
	}

	take := prop.Income*8 + int(w.Random()*float64(prop.Income*10))
	take = take * prop.Condition / 100
	w.Earn(take)
	w.Player.Heat = min(100, w.Player.Heat+w.HandHeat(hand, 10))
	w.Player.Respect += w.HandRespectFor(hand, 2)
	prop.Condition = max(0, prop.Condition-5)
	// A car outside is a thing witnesses describe, so driving to a robbery
	// makes it that much easier to work out who did it.
	if trail := w.CarTrail(); trail > 0 && !hand.Crew {
		w.Player.Heat = min(100, w.Player.Heat+trail*4)
		w.Log("Somebody described the car", fmt.Sprintf("A %s was parked where it had no business being. Attention is now %d.", lowerFirst(VehicleByTier(w.Player.Car).Label), w.Player.Heat), "danger")
	}
	if owner != nil {
		owner.Cash = max(0, owner.Cash-take)
		owner.Goodwill = max(-100, owner.Goodwill-25-w.CarTrail()*3)
		w.RetaliationFrom(owner.ID)
		w.Log("Taken from "+place.Name, fmt.Sprintf("$%d out of %s. %s will not need long to work out who would dare.", take, place.Name, owner.Name), "politics")
		w.Report("robbery", "ROBBERY AT "+strings.ToUpper(place.Name),
			w.unattributed(place.Name, fmt.Sprintf("A substantial sum was taken from %s, an establishment associated with %s.", place.Name, owner.Name)))
		return nil
	}
	w.Log("Taken from "+place.Name, fmt.Sprintf("$%d out of the till at %s. Nobody there answers to anyone who will come looking.", take, place.Name), "business")
	w.Report("robbery", "ROBBERY AT "+strings.ToUpper(place.Name),
		w.unattributed(place.Name, fmt.Sprintf("The day's takings were taken from %s.", place.Name)))
	return nil
}

// robber picks who comes for the player: someone real, with a place in the
// city and a reason. Low standing and high ambition make a person likelier to
// try it, but anyone might.
func (w *World) robber() *NPC {
	candidates := []*NPC{}
	for _, n := range w.People() {
		own := false
		for _, c := range w.Player.Crew {
			if c.ID == n.ID {
				own = true
			}
		}
		if own || n.Rank >= RankLeader {
			continue // the person at the top does not do this personally
		}
		candidates = append(candidates, n)
	}
	if len(candidates) == 0 {
		return nil
	}
	weights, total := make([]int, len(candidates)), 0
	for i, n := range candidates {
		// Somebody who has a reason is likelier than somebody who merely has
		// an opportunity. This is how a collection in the street comes back.
		weights[i] = max(1, n.Ambition+RankLeader-n.Rank+n.Sore*3)
		total += weights[i]
	}
	roll := int(w.WorldRandom() * float64(total))
	for i, weight := range weights {
		roll -= weight
		if roll < 0 {
			return candidates[i]
		}
	}
	return candidates[len(candidates)-1]
}

// robberyExposure is how much the player looks worth robbing. Carrying goods is
// the clearest signal; so is walking around with a great deal of cash.
func (w *World) robberyExposure() float64 {
	exposure := 0.0
	// Only what is visibly on you. Stock under the floor of a car or down a
	// cellar is not a reason for anybody in the street to pick you out.
	if carrying := w.Exposed(); carrying > 0 {
		exposure += .04 + float64(min(carrying, 40))/400
	}
	if w.Player.Cash > 1500 {
		exposure += .03
	}
	// Somebody at your shoulder is a reason to pick a different mark.
	if len(w.Player.Crew) > 0 && w.Player.Crew[0].Loyalty >= 40 && len(w.Tasks) == 0 {
		exposure /= 2
	}
	if w.Player.Security > 0 && w.Player.Location == w.Player.Home {
		exposure /= 2
	}
	return exposure
}

// ConsiderRobbery gives the city a chance to take something from the player.
// Whether they learn who did it depends on whether they built anyone who would
// tell them.
func (w *World) ConsiderRobbery() {
	// Nobody is robbed in the street while the police are holding them. The one
	// thing a cell is good for.
	if w.Held() {
		return
	}
	exposure := w.robberyExposure()
	if exposure <= 0 || w.WorldRandom() >= exposure {
		return
	}
	thief := w.robber()
	if thief == nil {
		return
	}
	attribution := "You never saw a face worth describing."
	if w.Reach() >= 2 {
		attribution = fmt.Sprintf("It takes a day and a few questions, but the name that comes back is %s.", thief.Name)
	}

	if carrying := w.Carrying(); carrying > 0 {
		w.Seize("You were jumped and your stock was taken. " + attribution)
		w.Player.Respect = max(0, w.Player.Respect-2)
		return
	}
	// What is behind the panelling at home is not in a pocket in the street.
	loss := min(w.Reachable(), 60+int(w.WorldRandom()*float64(w.Reachable())/4))
	if loss <= 0 {
		return
	}
	w.Player.Cash -= loss
	w.Player.Respect = max(0, w.Player.Respect-2)
	w.Log("Robbed in the street", fmt.Sprintf("$%d taken off you. %s", loss, attribution), "danger")
}
