package core

import "fmt"

func (w *World) strikeAttacker(hand Hand) CueAttacker {
	if hand.Crew && len(w.Player.Crew) > 0 {
		member := w.Player.Crew[0]
		tier := 0
		if n := w.NPC(member.ID); n != nil {
			tier = n.Weapon
		}
		return CueAttacker{ID: member.ID, Name: member.Name, Weapon: min(3, max(0, tier))}
	}
	return CueAttacker{ID: "player", Name: w.Player.Name, Weapon: min(3, max(0, w.Player.Weapon))}
}

// A player-directed strike has known equipment; random unrelated manners must
// not contradict the recorded attack or its visual interpretation.
func (w *World) strikeManner(victim *NPC, tier int) string {
	variants := []string{"killed in a close-quarters attack", "overpowered and killed", "beaten to death"}
	if tier > 0 {
		gun := []string{"", "a revolver", "a pump shotgun", "a Thompson"}[min(3, tier)]
		variants = []string{"shot with " + gun, "shot at close range with " + gun, "killed by gunfire from " + gun}
	}
	how := variants[int(w.WorldRandom()*float64(len(variants)))%len(variants)]
	place, _ := PlaceByID(victim.Location)
	return fmt.Sprintf("At %s, %s: %s was %s.", place.Name, hourOf(w.Minute), victim.Name, how)
}
