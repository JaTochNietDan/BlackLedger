package core

import "fmt"

// A protagonist dies and their premises pass to their estate, which is a string
// in an ownership field and nothing else. Their people keep a faction id that no
// longer resolves to anything: orphaned, unfindable, still counted as somebody's
// men by a lookup that returns nil.
//
// What actually happens when a man who ran something is killed is that somebody
// who worked for him is running it by the end of the week. The next protagonist
// arrives into a city that contains their predecessor's organization, under a
// name they never chose, holding premises they used to own — and can deal with
// it, work for it, or take it back.

// InheritCost is the strength an organization loses when the person it was
// built around is killed. It is a great deal: most of what it was worth was him.
const InheritCost = 30

// Inherit is what the player's organization becomes when the player is gone. It
// reports the name of whoever took it over, or "" when there was nobody to.
func (w *World) Inherit() string {
	me := w.PlayerOrganizationID()
	people := w.Members(me)
	if len(people) == 0 {
		// Nobody was left to carry it. Their people were the organization.
		w.orphan(me)
		return ""
	}
	successor := people[0]
	holdings := w.FamilyHoldings(me)

	// A new organization under a new name, with the old one's quarrels.
	id, name := "estate:"+successor.ID, successor.Name+"'s people"
	power := max(10, w.PlayerStrength()-InheritCost)
	w.Factions = append(w.Factions, Faction{
		ID: id, Name: name, Leader: successor.Name,
		Power: power, Peak: power, Cash: max(0, w.Player.Cash/3),
	})
	for _, n := range people {
		n.Faction = id
	}
	successor.Rank, successor.Role = RankLeader, "Head of "+name
	for _, held := range holdings {
		w.Properties[held].Owner = id
	}
	// The city's opinion of the man carries to the thing he left.
	for i := range w.Conflicts {
		c := &w.Conflicts[i]
		if c.A == me {
			c.A = id
		}
		if c.B == me {
			c.B = id
		}
	}
	w.Log(successor.Name+" is running it now", fmt.Sprintf("Everything that answered to %s answers to %s by the end of the week, under a name %s never chose. They hold %d of the premises and rather less of the strength.", w.Player.Name, successor.Name, w.Player.Name, len(holdings)), "politics")
	w.Report("politics", upper(name)+" TAKE OVER WHAT IS LEFT",
		fmt.Sprintf("Interests formerly associated with %s are understood to have passed to %s. Associates describe the transition as orderly.", w.Player.Name, successor.Name))
	// The id, not the name: whatever asks what became of this needs something
	// it can look the organization up by.
	return id
}

// orphan releases people whose organization no longer exists, so nobody is left
// answering to an id that resolves to nothing.
func (w *World) orphan(id string) {
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Faction == id {
			n.Faction, n.Rank, n.Role, n.Trust = "", RankAssociate, "Out of work", 0
		}
	}
}

// EstateName is what the city calls whatever a dead protagonist left, for the
// interface and for anything that needs to name it.
func (w *World) EstateName(id string) string {
	if f := w.faction(id); f != nil {
		return f.Name
	}
	return ""
}
