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

// Record the scenario before Kill changes the victim's living/travelling status. Selection
// consumes the same one world draw previously used for manner-of-death wording;
// it neither rerolls the successful hit nor consumes combat RNG.
func (w *World) strikePresentation(victim *NPC, tier int) (CueStrike, string) {
	roll := w.WorldRandom()
	variant := "close-quarters"
	variants := []string{"killed in a close-quarters attack", "overpowered and killed", "beaten to death"}
	if tier > 0 {
		gun := []string{"", "a revolver", "a pump shotgun", "a Thompson"}[min(3, tier)]
		variant = "close-shot"
		variants = []string{"shot with " + gun, "shot at close range with " + gun, "killed by gunfire from " + gun}
		if tier == 3 {
			variant = "burst"
		}
		if tier == 1 && victim.Weapon == 0 && victim.Sore == 0 && !w.Travelling(victim) && roll < .65 {
			variant = "back-of-head"
			variants = []string{"caught unaware and shot once in the back of the head with a revolver"}
		}
	}
	how := variants[int(roll*float64(len(variants)))%len(variants)]
	place, _ := PlaceByID(victim.Location)
	setting := ""
	if w.residentAtHome(victim) {
		setting = "home"
	} else if !w.Travelling(victim) {
		setting = "interior"
	}
	return CueStrike{Setting: setting, Variant: variant, Victim: CueActor{ID: victim.ID, Name: victim.Name}},
		fmt.Sprintf("At %s, %s: %s was %s.", place.Name, hourOf(w.Minute), victim.Name, how)
}
