package core

// What a war does to the people who are not in it.
//
// Measured over a hundred campaigns: the investor was party to no war, watched
// seven or nine break out elsewhere, and was never touched by one — lowest
// health a hundred, nothing taken, takings unchanged. A city could be at war
// around somebody who owned half of it and cost them nothing at all. Layer 2 of
// docs/LIVING_WORLD.md says the player may be a bystander, a beneficiary, or
// collateral, and only the first two were true.
//
// What a war does to everybody else is keep them at home. Nobody walks to a bar
// on the night two families are shooting at each other, and a room's takings are
// who is standing in it. So the cost reaches an owner who has done nothing and
// offended nobody, which is what being collateral means.

// StayIn is one in how many of the rest of the city — the people whose family
// is not the one doing the shooting — stays home anyway. Not most of it: a city
// at war is still a city, and the people with least to lose are the ones still
// on the street.
//
// The first version kept two people in three at home whoever they were, and the
// habit guard caught it: a moving face was where it usually is on only 63% of
// person-hours, against a promise of ninety. That promise is the point of the
// routine — a player learns somebody drinks at The Blue Hour and finds them
// there. So the curfew now falls first on the people with a reason to keep off
// the street, which is the two families in the fight, and only lightly on
// everybody else.
const StayIn = 6

// warring returns the families currently shooting at each other.
func (w *World) warring() []string {
	var sides []string
	for _, c := range w.Conflicts {
		if c.State == "war" {
			sides = append(sides, c.A, c.B)
		}
	}
	return sides
}

// CityAtWar reports whether any two organizations are actually fighting.
func (w *World) CityAtWar() bool { return len(w.warring()) > 0 }

// staysIn reports whether this person is one of the ones who does not go out
// tonight. Their own family being in it keeps them in; otherwise it is derived
// from their id, so the same people are the cautious ones every night of the
// same war rather than a different handful each evening.
func (w *World) staysIn(n *NPC) bool {
	if n.Faction != "" {
		for _, side := range w.warring() {
			if side == n.Faction {
				return true
			}
		}
	}
	sum := 0
	for i := 0; i < len(n.ID); i++ {
		sum = sum*31 + int(n.ID[i])
	}
	if sum < 0 {
		sum = -sum
	}
	return sum%StayIn == 0
}
