package core

import (
	"fmt"
	"strings"
)

// An organization is only as loyal as its prospects. A family that has been
// beaten down, or that is bleeding into a war it is not winning, gives the
// people below its leader a reason to take what they are standing on and keep
// it. Splinters are how the city gains organizations nobody authored.

var splinterSurnames = []string{"Falcone", "Serra", "Marchetti", "Delano", "Kovac",
	"Brenner", "Alcaraz", "Vance", "Doyle", "Rizzo", "Lindqvist", "Amato"}

// No article. The seeded families are "Bellandi Family" and "Russo Outfit", and
// a splinter that carried "the " in its own name started every sentence it
// began in lower case — in the paper, in the ledger, and in the warning that
// somebody is asking where you sleep. Leads still repairs the names that older
// saves already hold.
var splinterForms = []string{"%s Crew", "%s Combine", "%s Brothers",
	"%s Syndicate", "%s Company"}

var splinterFirstNames = []string{"Sal", "Renata", "Tomas", "Ida", "Bruno", "Vera",
	"Otto", "Cleo", "Marco", "Dita"}

// maxOrganizations keeps the city legible. A splinter cannot form once the board
// is this crowded; existing organizations must fall first.
const maxOrganizations = 5

func (w *World) nameTaken(name string) bool {
	for _, f := range w.Factions {
		if f.Name == name || f.Leader == name {
			return true
		}
	}
	for _, n := range w.NPCs {
		if n.Name == name {
			return true
		}
	}
	return false
}

// newOrganizationName draws a name that no living organization or person is
// already using. It is drawn from the reserved world stream, so the city
// reshaping itself never shifts the odds of a player's decision.
func (w *World) newOrganizationName() (string, string, string, bool) {
	for attempt := 0; attempt < 40; attempt++ {
		surname := splinterSurnames[int(w.WorldRandom()*float64(len(splinterSurnames)))%len(splinterSurnames)]
		form := splinterForms[int(w.WorldRandom()*float64(len(splinterForms)))%len(splinterForms)]
		first := splinterFirstNames[int(w.WorldRandom()*float64(len(splinterFirstNames)))%len(splinterFirstNames)]
		name := fmt.Sprintf(form, surname)
		leader := first + " " + surname
		if w.nameTaken(name) || w.nameTaken(leader) {
			continue
		}
		id := fmt.Sprintf("splinter-%s", surname)
		if w.faction(id) != nil {
			continue
		}
		return id, name, leader, true
	}
	return "", "", "", false
}

// splinterReady reports whether a family is weak enough, and holds enough, for
// someone inside it to break away.
func (w *World) splinterReady(f *Faction) bool {
	if len(w.Factions) >= maxOrganizations {
		return false
	}
	// Somebody walking out on the player takes a business rather than founding
	// a family, which is one path rather than two doing the same thing.
	if f.ID == w.PlayerOrganizationID() {
		return false
	}
	holdings := w.FamilyHoldings(f.ID)
	if len(holdings) < 2 {
		return false // nobody walks away with the last thing the family owns
	}
	weakened := f.Power <= peak(f)*3/4
	fighting := false
	for _, c := range w.Conflicts {
		if (c.A == f.ID || c.B == f.ID) && c.State == "war" {
			fighting = true
		}
	}
	// An organization with nobody left to fight fractures from the inside.
	// Without this a city that ends in one family owning everything stays that
	// way, and nothing further can happen in it.
	unopposed := true
	for _, other := range w.Factions {
		if other.ID != f.ID && len(w.FamilyHoldings(other.ID)) > 0 {
			unopposed = false
		}
	}
	return weakened || fighting || unopposed
}

// Splinter breaks a new organization out of an existing one, taking a holding
// with it. The parent loses the ground, the strength and the money that went
// with it, and the two are enemies from the first day.
func (w *World) Splinter(parent *Faction) bool {
	if !w.splinterReady(parent) {
		return false
	}
	id, name, leader, ok := w.newOrganizationName()
	if !ok {
		return false
	}
	holdings := w.FamilyHoldings(parent.ID)
	taken := holdings[0]
	for _, candidate := range holdings {
		if w.Properties[candidate].Condition < w.Properties[taken].Condition {
			taken = candidate
		}
	}
	place, exists := PlaceByID(taken)
	if !exists {
		return false
	}

	strength := max(20, parent.Power/3)
	money := parent.Cash / 6
	parent.Power = max(10, parent.Power-strength/2)
	parent.Cash = max(0, parent.Cash-money)
	w.Properties[taken].Owner = id

	w.Factions = append(w.Factions, Faction{ID: id, Name: name, Leader: leader,
		Power: strength, Goodwill: 0, Cash: money, Peak: strength})
	// The new leader is a person the city can deal with and the director can speak through.
	w.NPCs = append(w.NPCs, NPC{ID: id, Name: leader, Role: "Head of " + name, Trust: 0,
		Voice: w.voiceFor(leader), Color: "#8d7f6a", Faction: id,
		Location: taken, Rank: RankLeader, Ambition: 60 + int(w.WorldRandom()*40), Skill: 40 + int(w.WorldRandom()*50)})
	// A breakaway takes people with it, not only ground.
	w.AddMember(id, "Lieutenant", RankLieutenant, taken)

	// A breakaway is a betrayal and starts hotter than any ordinary quarrel,
	// but not already at war: the parent has to decide whether it can afford
	// to take the ground back, and sometimes it cannot.
	w.Antagonize(parent.ID, id, 62)
	if c := w.Conflict(parent.ID, id); c != nil {
		c.State = "feud"
		c.Since = w.Minute
	}
	w.Log("A family splits", fmt.Sprintf("%s %s broken away from %s, taking %s. %s leads them, and %s wants it back.",
		Leads(name), Agree(name, "has", "have"), parent.Name, place.Name, leader, parent.Name), "politics")
	w.Report("split", "SPLIT IN "+strings.ToUpper(parent.Name),
		fmt.Sprintf("A faction led by %s has broken from %s and taken control of %s. Observers expect the dispute to be settled outside the courts.", leader, parent.Name, place.Name))
	return true
}
