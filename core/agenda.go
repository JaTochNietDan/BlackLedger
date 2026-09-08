package core

import "fmt"

// People in this city want things. Until now they had ambition and a standing
// and did nothing with either: everything that happened to them was done by the
// player or by their organization. This gives each of them a day of their own.
//
// They act through the same rules the player does. When somebody robs a
// business it is the same money leaving the same books, and the newspaper
// reports it the same way, which is what makes the city produce news nobody
// wrote.

// dailyAmbition is the chance a person pursues something on a given day.
// Ambition drives it, so the people who want more are the ones the city hears
// about, and it stays low enough that a day is mostly quiet.
func dailyAmbition(n *NPC) float64 {
	return float64(n.Ambition) / 1400
}

// PeopleDay gives every living person one chance a day to pursue what they
// want. Called from the clock, so it happens whether or not anyone is watching.
func (w *World) PeopleDay() {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Rank >= RankLeader {
			continue // whoever is at the top has people for this
		}
		if w.WorldRandom() >= dailyAmbition(n) {
			continue
		}
		w.pursue(n)
	}
}

// pursue is what one person does with one opportunity. An unaffiliated person
// wants somewhere of their own; somebody inside an organization wants to be
// worth more inside it.
func (w *World) pursue(n *NPC) {
	if n.Faction == "" {
		if w.claimPremises(n) {
			return
		}
	}
	w.takeFromSomebody(n)
}

// claimPremises is how an unaffiliated person stops being nobody: they take
// over premises no organization holds. Never the player's, and never a rival's.
func (w *World) claimPremises(n *NPC) bool {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || prop.Income <= 0 || w.Own(l.ID) || w.faction(prop.Owner) != nil {
			continue
		}
		if prop.Owner != "independent" && !hasPrefix(prop.Owner, "former:") {
			continue
		}
		if n.Skill+n.Ambition < 90 {
			return false // not yet somebody who could hold it
		}
		n.Rank = max(n.Rank, RankSoldier)
		n.Location = l.ID
		n.Role = "Runs " + l.Name
		w.Log(n.Name+" takes over "+l.Name, fmt.Sprintf("%s has put themselves in charge of %s. Nobody stopped them.", n.Name, l.Name), "politics")
		w.Report("business", upper(n.Name)+" TAKES OVER "+upper(l.Name),
			fmt.Sprintf("%s is now running %s, which had been standing without anyone to answer for it.", n.Name, l.Name))
		return true
	}
	return false
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// takeFromSomebody is a person helping themselves to somebody else's takings.
// The player's business is a target like any other, which is the half of this
// the player actually feels.
func (w *World) takeFromSomebody(n *NPC) {
	// Weighted by what a place actually takes, rather than always the richest.
	// Always picking the top earner meant the player's business was never a
	// target unless it was the best in the city, which is not how this works.
	candidates, weights, total := []string{}, []int{}, 0
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || prop.Income <= 0 || prop.Owner == n.Faction {
			continue // not from your own people
		}
		weight := max(1, prop.Income*prop.Condition/100)
		candidates = append(candidates, l.ID)
		weights = append(weights, weight)
		total += weight
	}
	if total == 0 {
		return
	}
	target := candidates[len(candidates)-1]
	roll := int(w.WorldRandom() * float64(total))
	for i, weight := range weights {
		roll -= weight
		if roll < 0 {
			target = candidates[i]
			break
		}
	}
	place, _ := PlaceByID(target)
	prop := w.Properties[target]

	// Skill against whatever protects the place.
	defence := 30
	if f := w.faction(prop.Owner); f != nil {
		defence += f.Power / 3
	}
	if w.Own(target) {
		defence += w.Guard() * 8
		if len(w.Player.Crew) > 0 && w.Player.Crew[0].Loyalty >= 40 {
			defence += 10
		}
	}
	if n.Skill*100/(n.Skill+defence) < int(w.WorldRandom()*100) {
		w.Log("Somebody tried it at "+place.Name, fmt.Sprintf("%s went at %s and was turned away.", n.Name, place.Name), "politics")
		return
	}

	take := prop.Income*6 + int(w.WorldRandom()*float64(prop.Income*8))
	take = take * prop.Condition / 100
	prop.Condition = max(0, prop.Condition-4)
	n.Rank = min(RankLieutenant, n.Rank+2)

	if w.Own(target) {
		// Somewhere with a still has something better than a till, and the
		// people who rob premises know what it is worth.
		if prop.Still && w.Holding("moonshine") > 0 {
			stolen := min(w.Holding("moonshine"), 3+int(w.WorldRandom()*8))
			w.Player.Stock["moonshine"] -= stolen
			attribution := "Nobody will say who."
			if w.Reach() >= 2 {
				attribution = "The name that comes back is " + n.Name + "."
			}
			w.Log("Crates gone from "+place.Name, fmt.Sprintf("%d crates went out of the back of %s. %s", stolen, place.Name, attribution), "danger")
			w.Report("robbery", "THEFT AT "+upper(place.Name),
				w.unattributed(place.Name, fmt.Sprintf("Goods were taken from %s overnight.", place.Name)))
			return
		}
		// The player loses the takings themselves, and hears about it.
		w.Player.Cash = max(0, w.Player.Cash-take)
		attribution := "Nobody will say who."
		if w.Reach() >= 2 {
			attribution = "The name that comes back is " + n.Name + "."
		}
		w.Log("Taken from "+place.Name, fmt.Sprintf("$%d went out of %s while you were elsewhere. %s", take, place.Name, attribution), "danger")
		w.Report("robbery", "ROBBERY AT "+upper(place.Name),
			w.unattributed(place.Name, fmt.Sprintf("The day's takings were taken from %s.", place.Name)))
		return
	}

	if f := w.faction(prop.Owner); f != nil {
		f.Cash = max(0, f.Cash-take)
		w.Log("Taken from "+place.Name, fmt.Sprintf("%s helped themselves to $%d of %s money at %s.", n.Name, take, f.Name, place.Name), "politics")
		w.Report("robbery", "ROBBERY AT "+upper(place.Name),
			w.unattributed(place.Name, fmt.Sprintf("A sum was taken from %s, an establishment associated with %s.", place.Name, f.Name)))
		return
	}
	w.Log("Taken from "+place.Name, fmt.Sprintf("%s helped themselves to $%d at %s.", n.Name, take, place.Name), "politics")
	w.Report("robbery", "ROBBERY AT "+upper(place.Name),
		w.unattributed(place.Name, fmt.Sprintf("The day's takings were taken from %s.", place.Name)))
}
