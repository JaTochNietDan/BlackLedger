package core

import "fmt"

type HomeChange struct {
	ID, Name, From, To, Accommodation string
}

// PlanHomeMove works on copied residence records. Reading an action must never
// evict a tenant, move an actor or change the campaign's money.
func (w *World) PlanHomeMove(target string) ([]HomeChange, string) {
	changes := []HomeChange{}
	if ResidentialCapacity[target] <= 0 {
		return changes, "This address has no residential accommodation"
	}
	trial := *w
	trial.NPCs = append([]NPC(nil), w.NPCs...)
	trial.Player.Home = target
	trial.SettleHousing()
	for i, n := range w.NPCs {
		if n.Dead {
			continue
		}
		after := trial.NPCs[i]
		if n.Home != "" && after.Home == "" {
			return nil, "There is no vacant accommodation for " + n.Name + "; this move cannot be arranged"
		}
		if n.Home != after.Home {
			changes = append(changes, HomeChange{n.ID, n.Name, n.Home, after.Home, after.Accommodation})
		}
	}
	return changes, ""
}

func homePlansMatch(a, b []HomeChange) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func homeMoveDetail(changes []HomeChange) string {
	out := ""
	for _, change := range changes {
		if change.To == "" {
			continue
		}
		out += fmt.Sprintf(" %s will be rehoused at %s.", change.Name, placeName(change.To))
	}
	return out
}

func (w *World) applyHomeChanges(changes []HomeChange) {
	for _, change := range changes {
		if n := w.NPC(change.ID); n != nil && !n.Dead {
			n.Home, n.Accommodation = change.To, change.Accommodation
			w.Log("A change of address", n.Name+" now lives at "+placeName(change.To)+". Their workplace and current journey are unchanged.", "personal")
		}
	}
}
