package core

import "fmt"

// A coup already existed: somebody below the leader of a failing organization
// took their place by force, the leader died, and the organization lost twelve
// strength. That was the whole of it. Nobody chose a side, nobody was counted
// afterwards, and the family that had just killed its own head went back to
// business the following morning.
//
// What was missing is the part everybody remembers: the weeks after. A family
// that has just eaten one of its own is divided, and the people who backed the
// wrong man have to live somewhere.

const (
	// CoupGrudge is the weight a standing grievance against your own leader
	// adds to the decision to move on them.
	CoupGrudge = 45
	// PurgeChance is how often somebody who backed the losing side is dealt
	// with rather than forgiven.
	PurgeChance = .45
	// WalkOutChance is how often the losers of a coup leave rather than stay
	// under the person who did it.
	WalkOutChance = .35
)

// coupNerve is how willing somebody is to move on the person above them. A
// loyal person will not, whatever they are owed; a grievance they have been
// carrying is worth more than any amount of ambition.
func (w *World) coupNerve(challenger, leader *NPC) float64 {
	if TemperamentOf(challenger).Loyal {
		return 0
	}
	nerve := float64(challenger.Ambition)/400 + .35
	for _, g := range w.Grudges {
		if g.Holder == challenger.ID && g.Against == leader.ID {
			nerve += float64(g.Weight) / 200
			break
		}
	}
	return nerve * TemperamentOf(challenger).Nerve
}

// CoupAftermath is what a family looks like the week after it ate one of its
// own. Somebody who backed the wrong side is dealt with, somebody else walks
// out, and everybody left has an opinion about it.
func (w *World) CoupAftermath(f *Faction, winner, loser *NPC) {
	members := w.Members(f.ID)
	if len(members) == 0 {
		return
	}

	// Everybody who is left now has a view on the person who did it. That is
	// the drama: a family under new leadership where somebody is still angry.
	for _, m := range members {
		if m.Dead || m.ID == winner.ID {
			continue
		}
		w.Resent(m.ID, winner.ID, 20, "what was done to "+loser.Name)
	}

	// The people who backed the wrong man are counted.
	if w.WorldRandom() < PurgeChance {
		for _, m := range members {
			if m.Dead || m.ID == winner.ID || m.Rank >= RankLeader {
				continue
			}
			w.KillBy(m.ID, winner, fmt.Sprintf("They had backed the wrong side when %s took %s.", winner.Name, f.Name))
			f.Power = max(10, f.Power-4)
			break
		}
	}

	// And somebody walks out rather than answer to the person who did it,
	// taking whatever they can hold with them.
	if w.WorldRandom() < WalkOutChance && w.Splinter(f) {
		w.Log("They would not answer to "+winner.Name,
			fmt.Sprintf("What is left of %s is not all of what it was. Some of them would rather start again than take orders from whoever did this.", f.Name), "politics")
		return
	}

	w.Report("politics", "UPHEAVAL IN "+upper(f.Name),
		fmt.Sprintf("%s is understood to be under new leadership following the death of %s. Associates describe the organization as divided.", f.Name, loser.Name))
}

// InternalMove is one person below the top deciding they should be at it. It
// replaces the older roll: who moves, how willing they are and what happens to
// everybody afterwards now all come from the same place the rest of the city's
// history does.
func (w *World) InternalMove(f *Faction) bool {
	// The player's organization has no leader inside it to move on: the leader
	// is the player. Somebody who has had enough of them defects instead, which
	// OwnPeopleDay handles.
	if f.ID == w.PlayerOrganizationID() {
		return false
	}
	members := w.Members(f.ID)
	if len(members) < 2 {
		return false
	}
	leader := members[0]

	// Whoever is most willing, rather than whoever happens to be second.
	var challenger *NPC
	best := 0.0
	for _, m := range members[1:] {
		if m.Dead || m.Rank < RankLieutenant {
			continue
		}
		if nerve := w.coupNerve(m, leader); nerve > best {
			challenger, best = m, nerve
		}
	}
	if challenger == nil {
		return false
	}

	// A leader who is still winning is a leader nobody moves on, unless
	// somebody has a specific reason that has been growing for weeks.
	owed := false
	for _, g := range w.Grudges {
		if g.Holder == challenger.ID && g.Against == leader.ID && g.Weight >= CoupGrudge {
			owed = true
			break
		}
	}
	if f.Power > peak(f)*3/5 && !owed {
		return false
	}
	if w.WorldRandom() >= best {
		return false
	}

	odds := .35 + float64(w.Poise(challenger)-w.Poise(leader))/200 + float64(challenger.Ambition)/400
	odds = min64(.8, max64(.15, odds))

	if w.WorldRandom() < odds {
		w.Log("A move against "+leader.Name,
			fmt.Sprintf("%s has taken %s's place at the head of %s by force. The organization is divided and bleeding.", challenger.Name, leader.Name, f.Name),
			"politics")
		w.KillBy(leader.ID, challenger, challenger.Name+" moved against them from inside "+f.Name+".")
		f.Power = max(10, f.Power-12)
		w.CoupAftermath(f, challenger, leader)
		return true
	}
	w.Log("A move that failed",
		fmt.Sprintf("%s tried to take %s from %s and did not survive it. The people who backed them are being counted.", challenger.Name, f.Name, leader.Name),
		"politics")
	w.KillBy(challenger.ID, leader, "They had moved against "+leader.Name+" and lost.")
	f.Power = max(10, f.Power-5)
	w.CoupAftermath(f, leader, challenger)
	return true
}

// movedThisDay is the same chance the city gives an organization to turn on its
// own leadership, exposed so a test can measure how often it happens over a
// season without reaching into the war machinery.
func (w *World) movedThisDay(f *Faction) bool {
	if w.WorldRandom() >= .03 {
		return false
	}
	return w.InternalMove(f)
}
