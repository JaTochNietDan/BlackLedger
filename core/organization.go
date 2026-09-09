package core

import "fmt"

// Everything the city does to an organization — declaring war on it, raiding
// its holdings, weighing its strength against a neighbour's, deciding it is
// weak enough to move on — it could not do to the player, because the player
// was not an organization. They were a man with businesses, standing outside
// the machinery that runs everybody else.
//
// Past a certain point that stops being true. A man holding two premises with a
// name people know is not a man any more, and the city starts filing him where
// it files the others: in the same list, with the same strength, in the same
// conflicts, raided by the same function.

const (
	// OrganizationHoldings is the ground it takes before anybody thinks of you
	// as a thing rather than a person.
	OrganizationHoldings = 2
	// PremisesRespect is the standing it takes before anybody will sell you
	// premises. It was written out three times as a bare 6 — in the rule, in
	// the opportunity that suggests it, and in the guide page whose own header
	// promises it "cannot tell you something the rules do not". One of the
	// three could have changed and the guide would have gone on saying six.
	PremisesRespect = 6
	// OrganizationStanding is the name it takes.
	OrganizationStanding = 25
)

// PlayerOrganizationID is how the city files the player's holdings, which is
// already how their ownership is written. Nothing here invents a new identity:
// it names the one the properties already carry.
func (w *World) PlayerOrganizationID() string { return fmt.Sprintf("player:%d", w.Life) }

// PlayerOrganization is the player's own organization, once there is one.
func (w *World) PlayerOrganization() *Faction {
	return w.faction(w.PlayerOrganizationID())
}

// Incorporated reports whether the city treats the player as an organization.
func (w *World) Incorporated() bool { return w.PlayerOrganization() != nil }

// OrganizationReady reports whether the player is big enough to be filed with
// the others.
func (w *World) OrganizationReady() bool {
	// Nobody starts their own thing while they are still somebody's soldier.
	if w.Player.Serves != "" {
		return false
	}
	return len(w.FamilyHoldings(w.PlayerOrganizationID())) >= OrganizationHoldings &&
		w.Player.Respect >= OrganizationStanding
}

// PlayerStrength is what the player's organization is worth in a fight: the
// ground they hold, the name they have, and whoever stands with them.
func (w *World) PlayerStrength() int {
	strength := 10
	for _, id := range w.FamilyHoldings(w.PlayerOrganizationID()) {
		prop := w.Properties[id]
		strength += 8 + prop.Condition/10
	}
	strength += min(30, w.Player.Respect/3)
	strength += w.Guard() * 3
	for _, c := range w.Player.Crew {
		strength += c.Loyalty / 8
	}
	// And whoever answers to you, weighted by whether they mean it.
	for _, n := range w.OwnPeople() {
		strength += 4 + n.Trust/20
	}
	return min(100, strength)
}

// Incorporate files the player with the others. Idempotent, and it does not
// undo itself: an organization that has been seen is an organization, whatever
// happens to it afterwards.
func (w *World) Incorporate() {
	if w.Incorporated() || !w.OrganizationReady() {
		return
	}
	name := w.Player.Name + "'s people"
	w.Factions = append(w.Factions, Faction{
		ID: w.PlayerOrganizationID(), Name: name, Leader: w.Player.Name,
		Power: w.PlayerStrength(), Peak: w.PlayerStrength(), Cash: w.Player.Cash,
	})
	// Everybody already in the city has a view on a new organization, which is
	// the same cold start any splinter gets.
	for i := range w.Factions {
		other := &w.Factions[i]
		if other.ID == w.PlayerOrganizationID() {
			continue
		}
		// Conflict creates the relationship if there is not one, at nothing.
		// Where they start is what they already thought of the player.
		c := w.Conflict(other.ID, w.PlayerOrganizationID())
		if c == nil {
			continue
		}
		c.Hostility = min(100, max(0, 40-other.Goodwill/2))
		c.State, c.Since = classify(c), w.Minute
	}
	w.Log("They have started calling you something", fmt.Sprintf("%s. Two premises and a name is the point at which this city stops thinking of you as a person and starts thinking of you as a thing it has to deal with. You are in the same book as the others now, and the same things are done to what is in that book.", name), "politics")
	w.Report("politics", "A NEW NAME ON THE WATERFRONT",
		fmt.Sprintf("Interests associated with %s are now spoken of as an organization rather than a proprietor. Rivals are said to have noticed.", w.Player.Name))
}

// OrganizationDay keeps the player's entry true: their strength is whatever
// their holdings, name and people are worth today, and their money is theirs.
func (w *World) OrganizationDay() {
	w.Incorporate()
	f := w.PlayerOrganization()
	if f == nil {
		return
	}
	f.Power = w.PlayerStrength()
	f.Peak = max(f.Peak, f.Power)
	f.Cash = w.Player.Cash
	f.Leader = w.Player.Name
}

// Dissolve removes the player's organization when the person it belonged to is
// gone. Their premises pass to their estate the way they always did, and the
// name goes with them.
func (w *World) Dissolve(id string) {
	kept := w.Factions[:0]
	for _, f := range w.Factions {
		if f.ID != id {
			kept = append(kept, f)
		}
	}
	w.Factions = kept
	conflicts := w.Conflicts[:0]
	for _, c := range w.Conflicts {
		if c.A != id && c.B != id {
			conflicts = append(conflicts, c)
		}
	}
	w.Conflicts = conflicts
}

// PlayerOrganizationDescription is what the city calls the player, for the
// interface.
func (w *World) PlayerOrganizationDescription() map[string]any {
	f := w.PlayerOrganization()
	if f == nil {
		needs := []string{}
		if held := len(w.FamilyHoldings(w.PlayerOrganizationID())); held < OrganizationHoldings {
			needs = append(needs, fmt.Sprintf("%d more premises", OrganizationHoldings-held))
		}
		if w.Player.Respect < OrganizationStanding {
			needs = append(needs, fmt.Sprintf("%d more respect", OrganizationStanding-w.Player.Respect))
		}
		return map[string]any{"named": false, "needs": needs}
	}
	return map[string]any{
		"named": true, "name": f.Name, "power": f.Power, "peak": f.Peak,
	}
}
