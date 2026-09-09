package core

import "fmt"

// Coming up stops at lieutenant. There is nothing above it, which makes the
// second career a wage with a ceiling rather than a story — and this city has
// exactly one thing above lieutenant, which is the man in the chair.
//
// Everything needed for this already exists: the city runs coups on itself
// every season, weighing a challenger's poise against a leader's and taking the
// organization apart afterwards. This is the same move with the player as the
// challenger, and what it wins is not a promotion. It is the organization.

const (
	// TakeoverMinutes is the evening it takes.
	TakeoverMinutes = 90
	// TakeoverStanding is the presence it takes before anybody would follow
	// you rather than him.
	TakeoverStanding = 40
	// TakeoverWeak is the share of its peak an organization must have fallen to
	// before its people would consider it, unless the player has a reason of
	// their own that has been growing.
	TakeoverWeak = 3
	// TakeoverKeep is the share of its people who stay for the man who did it.
	TakeoverKeep = .6
)

// Leader is whoever is at the top of an organization, as a person.
func (w *World) Leader(id string) *NPC {
	f := w.faction(id)
	if f == nil {
		return nil
	}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Name == f.Leader {
			return n
		}
	}
	return nil
}

// TakeoverReadiness explains why the player cannot move on the man above them,
// or returns "".
func (w *World) TakeoverReadiness() string {
	if w.Player.Serves == "" {
		return "You answer to nobody"
	}
	f := w.faction(w.Player.Serves)
	if f == nil {
		return "There is nothing left to take"
	}
	if w.ServiceRank() < RankLieutenant {
		return "Nobody below a lieutenant gets close enough to try"
	}
	leader := w.Leader(f.ID)
	if leader == nil {
		return "Nobody is in the chair"
	}
	if w.Player.Location != leader.Location {
		place, _ := PlaceByID(leader.Location)
		return "They are at " + place.Name
	}
	if w.Presence() < TakeoverStanding {
		return fmt.Sprintf("Nobody would follow you rather than them. You need %d presence", TakeoverStanding)
	}
	if w.Player.Health < 50 {
		return "You are in no condition for this"
	}
	if f.Power > peak(f)*TakeoverWeak/5 {
		return f.Name + " is doing too well for anybody to move on the man running it"
	}
	return ""
}

// takeoverOdds is the player against the man in the chair: what they are worth
// standing there, against his own competence and what the organization can put
// behind him.
func (w *World) takeoverOdds(leader *NPC, f *Faction) float64 {
	mine := float64(min(w.Presence(), 100))/3 + float64(w.Player.Service)*4 + w.WeaponEdge()*100
	if len(w.Player.Crew) > 0 && w.Player.Crew[0].Loyalty >= 40 {
		mine += 10
	}
	his := float64(w.Poise(leader)) + float64(f.Power)/2
	return min64(.85, max64(.15, mine/(mine+his)))
}

// TakeOver is the move. It is resolved before the clock advances, because a man
// who does not survive it does not collect the evening.
func (w *World) TakeOver() error {
	if reason := w.TakeoverReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	f := w.faction(w.Player.Serves)
	leader := w.Leader(f.ID)
	name := f.Name

	if w.Random() >= w.takeoverOdds(leader, f) {
		// He was ready, or somebody told him.
		injury := w.Absorb(40 + int(w.Random()*45))
		w.Ruin(80)
		w.Player.Health = max(0, w.Player.Health-injury)
		w.Player.Serves, w.Player.Service = "", 0
		f.Goodwill = max(-100, f.Goodwill-70)
		w.Log("They were expecting it", fmt.Sprintf("%s knew before you were through the door. You are not one of theirs any more and %s is not a name you can use.", leader.Name, name), "danger")
		w.RetaliationFrom(f.ID)
		if w.Player.Health <= 0 {
			w.Die("A move on " + leader.Name + " that they saw coming.")
		}
		return nil
	}

	w.KillBy(leader.ID, nil, fmt.Sprintf("They had run %s, and somebody who worked for them had come up far enough to want it.", name))
	w.Player.Serves, w.Player.Service = "", 0
	w.Player.Respect += 25
	w.Player.Heat = min(100, w.Player.Heat+20)

	// What was theirs is yours: the ground, the people who stay, and every
	// quarrel the name was in.
	me := w.PlayerOrganizationID()
	held := w.FamilyHoldings(f.ID)
	for _, id := range held {
		w.Properties[id].Owner = me
	}
	kept, left := 0, 0
	for _, n := range w.Members(f.ID) {
		if w.WorldRandom() < TakeoverKeep {
			n.Faction, n.Rank, n.Role, n.Trust = me, RankSoldier, "Yours", 30
			kept++
			continue
		}
		n.Faction, n.Rank, n.Role, n.Trust = "", RankAssociate, "Out of work", 0
		left++
	}
	for i := range w.Conflicts {
		c := &w.Conflicts[i]
		if c.A == f.ID {
			c.A = me
		}
		if c.B == f.ID {
			c.B = me
		}
	}
	w.Dissolve(f.ID)
	w.OrganizationDay()

	w.Log("It is yours", fmt.Sprintf("%s is dead and %s is a name nobody uses now. You hold %d of its premises, %d of its people stayed and %d would not, and every quarrel it was in is yours.", leader.Name, name, len(held), kept, left), "politics")
	w.Report("politics", upper(name)+" IS FINISHED",
		fmt.Sprintf("%s is dead and the organization they ran is understood to have passed to somebody who worked for them. Police have not commented.", leader.Name))
	return nil
}
