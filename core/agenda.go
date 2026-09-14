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

// cityPace scales each person's chance by how many people there are, so that
// filling the streets makes the city deeper rather than four times busier. A
// city of eight and a city of thirty-five produce roughly the same number of
// incidents a week; what changes is who they happen to.
func (w *World) cityPace() float64 {
	living := len(w.People())
	if living <= 8 {
		return 1
	}
	return 8 / float64(living)
}

// PeopleDay gives every living person one chance a day to pursue what they
// want. Called from the clock, so it happens whether or not anyone is watching.
func (w *World) PeopleDay() {
	w.PayTheCity()
	// Somebody ambitious inside an organization spends every day looking at the
	// person above them. Two a day against one a day of forgetting means it
	// takes about a month and a half to become a reason, which is the only
	// source of grievance in this city that does not need somebody to have died
	// first.
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Faction == "" || n.Rank < RankLieutenant || n.Rank >= RankLeader || n.Ambition < 60 || IsOfficial(n.ID) {
			continue
		}
		if members := w.Members(n.Faction); len(members) > 0 && members[0].ID != n.ID {
			w.Resent(n.ID, members[0].ID, 2, "standing in the way of what they want")
		}
	}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		// keepsPost covers the top of an organization, the officials, anybody
		// the police are holding, and — the reason this changed — anybody
		// holding one of the city's standing jobs. The paper carried
		// "DETECTIVE HARLOW TAKES OVER THE BLUE HOUR": the city detective had
		// walked off his beat and seized a casino.
		if n.Dead || w.keepsPost(n) {
			continue
		}
		if w.WorldRandom() >= dailyAmbition(n)*w.cityPace() {
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
	// Or somebody standing next to them, which is the only thing in this city
	// that ever reaches the people who are not in an organization and do not
	// run premises — the fixer, the driver, the detective. Without it they were
	// immortal by omission rather than by design.
	if w.WorldRandom() < .3 && w.takeFromAPerson(n) {
		return
	}
	w.takeFromSomebody(n)
}

// takeFromAPerson is one of the city's own people robbing another in the
// street, by the same arithmetic the player faces. It reports whether anything
// happened.
func (w *World) takeFromAPerson(n *NPC) bool {
	var mark *NPC
	for i := range w.NPCs {
		other := &w.NPCs[i]
		if other.Dead || other.ID == n.ID || other.Location != n.Location {
			continue
		}
		if other.Faction != "" && other.Faction == n.Faction {
			continue // not your own people
		}
		if IsOfficial(other.ID) {
			continue // a man with a title is not robbed in the street
		}
		mark = other
		break
	}
	if mark == nil {
		return false
	}
	purse := w.Pockets(mark)
	if purse < 40 {
		return false
	}
	place, _ := PlaceByID(n.Location)
	defence := w.Poise(mark) + 15
	if f := w.faction(mark.Faction); f != nil {
		defence += f.Power / 3
	}
	if w.Poise(n)*100/(w.Poise(n)+defence) < int(w.WorldRandom()*100) {
		w.Resent(mark.ID, n.ID, 25, "what was tried at "+place.Name)
		w.Log("Somebody tried it at "+place.Name, fmt.Sprintf("%s went at %s in the street and came off worse.", n.Name, mark.Name), "politics")
		return true
	}
	if f := w.faction(mark.Faction); f != nil {
		f.Cash = max(0, f.Cash-purse/2)
	}
	w.Resent(mark.ID, n.ID, 40, "being robbed at "+place.Name)
	w.Log("Robbed in the street", fmt.Sprintf("%s took $%d off %s near %s.", n.Name, purse, mark.Name, place.Name), "politics")
	w.Report("robbery", "ROBBERY IN "+upper(place.Name),
		w.unattributed(place.Name, fmt.Sprintf("Somebody was robbed near %s. Police have asked anybody who saw it to come forward.", place.Name)))
	return true
}

// claimPremises is how an unaffiliated person stops being nobody: they take
// over premises no organization holds. Never the player's, and never a rival's.
func (w *World) claimPremises(n *NPC) bool {
	if n == nil || n.Dead || n.Faction != "" || w.Inside(n) || len(w.FamilyHoldings(n.ID)) > 0 {
		return false
	}
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if !personalBusiness(l.ID) || prop == nil || prop.Income <= 0 || w.Own(l.ID) || w.faction(prop.Owner) != nil {
			continue
		}
		if prop.Owner != "independent" && !hasPrefix(prop.Owner, "former:") {
			continue
		}
		if n.Skill+n.Ambition < 90 {
			return false // not yet somebody who could hold it
		}
		// A takeover must claim the deed, not leave the same business open for
		// every subsequent ambitious resident to run at once.
		prop.Owner = n.ID
		prop.ProprietorDay = w.Minute/1440 + 1
		n.Heading, n.Arrives, n.Sets, n.Errand = "", 0, 0, ""
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
		defence += w.Guard()*8 + w.PostingDefenceAt(target)*2
		if len(w.Player.Crew) > 0 && w.Player.Crew[0].Loyalty >= 40 {
			defence += 10
		}
	}
	if w.Poise(n)*100/(w.Poise(n)+defence) < int(w.WorldRandom()*100) {
		w.Log("Somebody tried it at "+place.Name, fmt.Sprintf("%s went at %s and was turned away.", n.Name, place.Name), "politics")
		return
	}

	// Whoever answers for the place remembers who came for it. This is where
	// most of the city's private history starts.
	for i := range w.NPCs {
		if keeper := &w.NPCs[i]; !keeper.Dead && keeper.ID != n.ID && keeper.Location == target && keeper.Rank >= RankSoldier {
			w.Resent(keeper.ID, n.ID, 22, "what happened at "+place.Name)
			break
		}
	}

	take := prop.Income*6 + int(w.WorldRandom()*float64(prop.Income*8))
	take = int(float64(take*prop.Condition/100) * TemperamentOf(n).Greed)
	prop.Condition = max(0, prop.Condition-4)
	n.Rank = min(RankLieutenant, n.Rank+2)

	if w.Own(target) {
		// Somewhere with a still has something better than a till, and the
		// people who rob premises know what it is worth.
		if prop.Still && w.Holding("moonshine") > 0 {
			stolen := min(w.Holding("moonshine"), 3+int(w.WorldRandom()*8))
			w.Player.Stock["moonshine"] -= stolen
			attribution := "Nobody will say who."
			if w.Reach() >= 2 || (TemperamentOf(n).ID == "vain" && w.Reach() >= 1) {
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
		// Somebody who wants to be seen doing it is how anybody ever finds out.
		if w.Reach() >= 2 || (TemperamentOf(n).ID == "vain" && w.Reach() >= 1) {
			attribution = "The name that comes back is " + n.Name + "."
		}
		w.Log("Taken from "+place.Name, fmt.Sprintf("$%d went out of %s while you were elsewhere. %s", take, place.Name, attribution), "danger")
		w.Report("robbery", "ROBBERY AT "+upper(place.Name),
			w.unattributed(place.Name, fmt.Sprintf("The day's takings were taken from %s.", place.Name)))
		return
	}

	if f := w.faction(prop.Owner); f != nil {
		f.Cash = max(0, f.Cash-take)
		w.Log("Taken from "+place.Name, theftLine(n.Name, take, f.Name, place.Name), "politics")
		w.Report("robbery", "ROBBERY AT "+upper(place.Name),
			w.unattributed(place.Name, fmt.Sprintf("A sum was taken from %s, an establishment associated with %s.", place.Name, f.Name)))
		return
	}
	w.Log("Taken from "+place.Name, fmt.Sprintf("%s helped themselves to $%d at %s.", n.Name, take, place.Name), "politics")
	w.Report("robbery", "ROBBERY AT "+upper(place.Name),
		w.unattributed(place.Name, fmt.Sprintf("The day's takings were taken from %s.", place.Name)))
}

// theftLine says whose money was taken without putting a possessive after an
// organization's name. The seeded families read fine either way — "$133 of
// Bellandi Family money" is ordinary attributive English — but the living
// world names the families it makes after the person who broke away, and
// "$133 of Franca Sabbatini's people money" stacks one possessive on another.
// Read on the Ledger screen in a browser. The verb follows the name too.
func theftLine(who string, take int, family, place string) string {
	return fmt.Sprintf("%s helped themselves to $%d of what %s %s at %s.",
		who, take, family, Agree(family, "keeps", "keep"), place)
}
