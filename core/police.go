package core

import "fmt"

// The police were a single number that only went up and did one thing. They are
// now a third force in the city, applying to organizations as much as to the
// player: attention accumulates from what people actually do, fades when they
// stop, and past a point it arrives at the door.
//
// Nothing here is a hidden roll the player cannot anticipate. Heat is public,
// the thresholds are stated in the actions that risk them, and a raid is
// announced by the same newspaper that reports everything else.

const (
	// RaidThreshold is where the police stop watching and start arriving.
	RaidThreshold = 45
	// ForfeitThreshold is where they stop taking money and start taking
	// premises.
	ForfeitThreshold = 80
	// BribeCeiling is the point past which a detective will not be seen with
	// you. Money stops working before the danger does.
	BribeCeiling = 70
)

// CoolOff is the attention that fades in a day when nothing new is added. It is
// deliberately slower than what a hard-run business generates, so choosing to
// skim still costs attention rather than being quietly absorbed.
const CoolOff = 1

// PoliceDay fades attention and decides whether anyone gets a visit. Runs once
// a game day alongside the other books.
func (w *World) PoliceDay() {
	if w.Player.Heat > 0 {
		w.Player.Heat = max(0, w.Player.Heat-CoolOff)
	}
	// Organizations draw attention too, and a war is the loudest thing in the
	// city. This is the same pressure the player feels, applied to them.
	for i := range w.Factions {
		f := &w.Factions[i]
		fighting := false
		for _, c := range w.Conflicts {
			if c.State == "war" && (c.A == f.ID || c.B == f.ID) {
				fighting = true
			}
		}
		if !fighting {
			continue
		}
		fine := min(f.Cash, 200+f.Power*8)
		if fine <= 0 {
			continue
		}
		f.Cash -= fine
		if w.WorldRandom() < .25 {
			w.Report("police", "POLICE PRESSURE ON "+upper(f.Name),
				fmt.Sprintf("Officers have raided premises connected to %s, citing the recent violence. Property was searched and records seized.", f.Name))
		}
	}
	w.considerRaid()
}

func upper(s string) string {
	out := []rune(s)
	for i, r := range out {
		if r >= 'a' && r <= 'z' {
			out[i] = r - 32
		}
	}
	return string(out)
}

// considerRaid decides whether the police come for the player today. The chance
// rises with attention, and having somewhere for them to search is what makes a
// raid worth their time.
func (w *World) considerRaid() {
	if w.Player.Heat < RaidThreshold || !w.Player.Alive {
		return
	}
	chance := float64(w.Player.Heat-RaidThreshold) / 120
	if w.WorldRandom() >= chance {
		return
	}
	w.Raid()
}

// Raid is the visit. They take what is carried, they fine what they can prove,
// and past a point they take the premises the money came through.
func (w *World) Raid() {
	// The premises they search is the one earning the most for the player.
	var target string
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || !w.Own(l.ID) || prop.Income <= 0 {
			continue
		}
		if target == "" || prop.Income > w.Properties[target].Income {
			target = l.ID
		}
	}

	seized := w.Seize("The police came with a warrant.")
	// A car is registered to somebody, and a warrant that turns up a false
	// floor takes the car with what is in it.
	if w.Concealed() > 0 && w.Carrying() > 0 && w.WorldRandom() < .35 {
		hidden := w.Carrying()
		for _, g := range w.Goods {
			if w.Player.Stock != nil {
				w.Player.Stock[g.ID] = 0
			}
		}
		seized += hidden
		w.LoseCar("They found the false floor and took the car with it.")
	}
	w.SeizeArms()
	w.Ruin(20) // being turned out against a wall is hard on good clothes
	still, foundStill := w.StillFound()
	fine := min(w.Player.Cash, 150+w.Player.Heat*12)
	if foundStill {
		// Finding a still is what turns a search into a case.
		fine = min(w.Player.Cash, fine*2+400)
	}
	w.Player.Cash -= fine
	if foundStill {
		w.Log("They found the still at "+still, fmt.Sprintf("Copper and pipe out through the front door in daylight. A fine of $%d and they will be back.", fine), "danger")
		w.Report("police", "STILL SEIZED AT "+upper(still),
			fmt.Sprintf("Officers dismantled an illegal still at %s. A prosecution is said to be likely.", still))
	}

	if target == "" {
		w.Player.Heat = max(0, w.Player.Heat-15)
		w.Log("They came to the door", fmt.Sprintf("A search and a fine of $%d. There was no business of yours for them to turn over.", fine), "danger")
		w.Report("police", "ARRESTS EXPECTED AFTER CITY SEARCHES", "Officers executed warrants across the district. A police spokesman said the operation was part of a continuing investigation.")
		return
	}

	place, _ := PlaceByID(target)
	prop := w.Properties[target]
	if w.Player.Heat >= ForfeitThreshold {
		prop.Owner = "independent"
		prop.Mode = ""
		w.Player.Heat = max(0, w.Player.Heat-30)
		w.Log("They took it", fmt.Sprintf("%s is forfeit. A fine of $%d, %d units of stock gone, and the business is no longer yours.", place.Name, fine, seized), "danger")
		w.Report("police", "AUTHORITIES SEIZE "+upper(place.Name),
			fmt.Sprintf("%s has been seized following an investigation into its accounts. The premises are closed pending proceedings.", place.Name))
		return
	}

	prop.Condition = max(0, prop.Condition-20)
	w.Player.Heat = max(0, w.Player.Heat-18)
	w.Log("Turned over at "+place.Name, fmt.Sprintf("A fine of $%d, %d units of stock gone, and %s was left in a state. They will be back if nothing changes.", fine, seized, place.Name), "danger")
	w.Report("police", "RAID AT "+upper(place.Name),
		fmt.Sprintf("Officers searched %s this morning. No charges have yet been brought.", place.Name))
}

// BribeCost is what a detective wants to lose the paperwork. It rises with what
// there is to lose.
func (w *World) BribeCost() int {
	return 120 + w.Player.Heat*28
}

// BribeReadiness explains why the arrangement cannot be made, or returns "".
func (w *World) BribeReadiness() string {
	if w.Player.Heat == 0 {
		return "Nobody is looking at you"
	}
	if w.Player.Heat >= BribeCeiling {
		return "You are too well known for anyone to be seen taking it"
	}
	if w.Player.Cash < w.BribeCost() {
		return "Not enough cash"
	}
	return ""
}

// Bribe buys the paperwork going missing. It works while the player is still
// somebody a detective can afford to be seen with.
func (w *World) Bribe() error {
	if reason := w.BribeReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	cost := w.BribeCost()
	if err := w.Pay(cost); err != nil {
		return err
	}
	cleared := min(w.Player.Heat, 12+w.Player.Contacts*3)
	w.Player.Heat = max(0, w.Player.Heat-cleared)
	if npc := w.NPC("harlow"); npc != nil {
		npc.Trust += 2
	}
	w.Log("An understanding with the detective", fmt.Sprintf("$%d, and a file goes to the bottom of a pile. Attention falls by %d, to %d. This does not buy the next one.", cost, cleared, w.Player.Heat), "personal")
	return nil
}
