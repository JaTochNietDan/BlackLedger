package core

import "fmt"

// Somebody talks.
//
// The last of the three risks layer 6 of docs/LIVING_WORLD.md asks the
// underground trade to carry: "seizure, informants, a rival who wants the
// route". Seizure was there, a rival went in last slice, and this is the one
// that is not about anybody wanting what you have. It is about somebody who
// wants you finished and has found a cheaper way to do it than a gun.
//
// It has to be a person. This game's rule about the people who rob the player —
// "named people with a place in the city, not anonymous thieves, and they can be
// answered" — holds here or the whole thing is a dice roll wearing a hat. So an
// informant is somebody already carrying a grudge the simulation committed,
// with a reason it recorded, and saying it costs them the grudge: telling the
// police is what they had to get off their chest.

const (
	// TalksAt is the weight of grievance at which somebody stops brooding and
	// picks up a telephone. Below the weight at which people move against
	// somebody themselves, because this is the coward's way and always was.
	TalksAt = 35
	// TalksChance is how often, per day, one of them does.
	TalksChance = .16
	// TalkHeat is what a word in the right ear is worth to the police.
	TalkHeat = 14
	// TalkSpent is how much of a grievance is discharged by saying it.
	TalkSpent = 25
	// TellableHeat is the attention at which somebody with a grievance has
	// something worth telling even if the player is carrying nothing: the
	// police already have a file, and a name in the right ear thickens it.
	TellableHeat = 25
)

// ConsiderInformant is somebody in this city deciding that what the player is
// carrying is worth telling the police about. Reports whether anybody did.
func (w *World) ConsiderInformant() bool {
	// There has to be something to tell them. Tying this to the trade alone was
	// too narrow and the measurement said so: a policy that ran the route made
	// no enemies, and the policy that made enemies never traded, so nobody in a
	// hundred campaigns ever picked up a telephone. What an informant needs is
	// somebody who hates them and something to point at — a trade being run,
	// crates in their hands, or a name the police already have a file on.
	if !w.Player.Alive {
		return false
	}
	if w.Player.Runs < RouteNotice/2 && w.Carrying() == 0 && w.Player.Heat < TellableHeat {
		return false
	}
	if w.WorldRandom() >= TalksChance {
		return false
	}
	var teller *NPC
	for _, n := range w.Aggrieved() {
		if n.Sore < TalksAt || n.Dead {
			continue
		}
		if n.Faction != "" && n.Faction == w.PlayerOrganizationID() {
			continue // your own do not, whatever they are owed
		}
		if len(w.Player.Crew) > 0 && w.Player.Crew[0].ID == n.ID {
			continue
		}
		teller = n
		break
	}
	if teller == nil {
		return false
	}

	about := teller.SoreAt
	if about == "" {
		about = "what you did"
	}
	teller.Sore = max(0, teller.Sore-TalkSpent)
	w.Player.Heat = min(100, w.Player.Heat+TalkHeat)

	// Whether the name comes back is the same question the city asks about a
	// robbery: somebody with contacts hears who, and somebody without hears
	// only that somebody did.
	if w.Reach() >= 2 {
		w.Log("Somebody talked", fmt.Sprintf("A word reached Ward Street about what you have been moving, and it takes a day and a few questions to find out whose. %s, over %s. Contacts are what put a name to it. Police interest is now %d.",
			teller.Name, about, w.Player.Heat), "danger")
		return true
	}
	w.Log("Somebody talked", fmt.Sprintf("Ward Street knows more about what you have been moving than anybody there worked out for themselves. You never saw a face and nobody will give you a name. Police interest is now %d.",
		w.Player.Heat), "danger")
	return true
}
