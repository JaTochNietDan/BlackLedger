package core

import "fmt"

// A casino used to be a room the player played in and a number that ran itself.
// The city's own people never sat down: the money behind the tables came from
// nowhere anybody lived in, and nobody in Bellwether was ever poorer on a
// Tuesday because of a Monday night.
//
// They play now. One person a night at each room that runs a float puts down
// what they can stand to, the house holds the same edge it holds against the
// player, and what the city loses is exactly what the house takes — no more,
// and out of pockets the rest of this city already reads from. Somebody cleaned
// out at the tables cannot find a garage's fee for the glass, and cannot put
// money down on a car at the forecourt. That is the point of it.

const (
	// TableLimit is the most anybody in this city will put down in one night.
	TableLimit = 120
	// TablePart is the share of what somebody has on them that they are willing
	// to lose in an evening, as a divisor.
	TablePart = 4
	// TableFloor is the least somebody will bother sitting down with.
	TableFloor = 40
)

// OnTheFloor is who is standing in a room that runs tables. The city has always
// known where everybody is; a gaming floor with nobody on it is a room, not a
// house.
func (w *World) OnTheFloor(id string) []*NPC {
	out := []*NPC{}
	if !HasBankroll(id) {
		return out
	}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Location == id {
			out = append(out, n)
		}
	}
	return out
}

// TableStakeFor is what this person will put down tonight, and nothing if they
// are not sitting down at all.
func (w *World) TableStakeFor(n *NPC) int {
	if n == nil || n.Dead {
		return 0
	}
	stake := min(TableLimit, n.Purse/TablePart)
	if stake < TableFloor {
		return 0
	}
	return stake
}

// TableNight is one person a night at every room that runs a float. It is one
// because the same reason holds here as at the forecourt: a city where everyone
// gambles on the same evening is a city where nothing means anything.
func (w *World) TableNight() {
	for _, l := range Locations {
		if !HasBankroll(l.ID) {
			continue
		}
		prop := w.Properties[l.ID]
		if prop == nil {
			continue
		}
		for _, n := range w.OnTheFloor(l.ID) {
			stake := w.TableStakeFor(n)
			if stake == 0 {
				continue
			}
			// The house has to be able to pay them. A room that cannot cover
			// what is put down in front of it is not running a game.
			mine := w.Own(l.ID)
			switch {
			case mine && prop.Bankroll < stake:
				continue
			case !mine && w.businessFunds(l.ID) < stake:
				continue
			}
			// An even-money bet with the house's edge on it: the same edge the
			// player faces across the same felt, seen from the other side.
			won := w.WorldRandom() < float64(100-HouseEdge)/200
			if won {
				n.Purse += stake
			} else {
				n.Purse -= stake
			}
			take := stake
			if won {
				take = -stake
			}
			if mine {
				// Money won and lost behind your own tables stays behind them,
				// which is the rule the nightly float already runs on.
				prop.Bankroll = max(0, prop.Bankroll+take)
				place, _ := PlaceByID(l.ID)
				if won {
					w.Log("A winner at "+place.Name, fmt.Sprintf("%s took $%d off the table.", n.Name, stake), "business")
				}
			} else {
				w.changeBusinessFunds(l.ID, take)
			}
			break // one a night, per room
		}
	}
}
