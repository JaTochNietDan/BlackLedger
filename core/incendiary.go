package core

import "fmt"

const (
	IncendiaryCost    = 40
	IncendiaryMinutes = 15
	IncendiaryHeat    = 18
)

func (w *World) IncendiaryReadiness(id string) string {
	prop := w.Properties[id]
	if prop == nil || prop.Income <= 0 {
		return "There is no business to attack here"
	}
	if w.Own(id) {
		return "This is yours"
	}
	if w.Player.Location != id {
		return "You would have to be there"
	}
	if prop.Condition <= 0 {
		return "The premises are already wrecked"
	}
	for _, fire := range w.ActiveBuildingFires() {
		if fire.Target == id && w.Minute < fire.ExtinguishedAt {
			return "The premises are already burning"
		}
	}
	if w.Player.Cash < IncendiaryCost {
		return "Not enough cash"
	}
	return ""
}

// Incendiary attacks damage premises without borrowing the explosive charge's
// casualties, staff removal or bankroll destruction. Fire response is shared.
func (w *World) Incendiary(id string) error {
	if reason := w.IncendiaryReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(IncendiaryCost); err != nil {
		return err
	}
	prop := w.Properties[id]
	place, _ := PlaceByID(id)
	damage := min(prop.Condition, 15+int(w.WorldRandom()*16))
	prop.Condition -= damage
	prop.Supply = max(0, prop.Supply-10)
	prop.Trouble = true
	w.Player.Heat = min(100, w.Player.Heat+IncendiaryHeat)
	if owner := w.faction(prop.Owner); owner != nil {
		owner.Goodwill = max(-100, owner.Goodwill-25)
		w.RetaliationFrom(owner.ID)
	}
	w.igniteBuilding(id)
	body := fmt.Sprintf("You threw an incendiary bottle at %s and fled the frontage. Fire damaged the premises by %d points; condition is now %d%%. Supplies were lost, and the fire brigade has been called.", place.Name, damage, prop.Condition)
	w.Log("Fire at "+place.Name, body, "danger")
	headline := "ARSON AT " + upper(place.Name)
	w.Report("attack", headline, fmt.Sprintf("An incendiary attack damaged %s. Witnesses saw the attacker flee. Police are treating the fire as deliberate.", place.Name))
	w.Witness("incendiary", id, body, headline)
	w.VisualCues[len(w.VisualCues)-1].Attacker = &CueAttacker{ID: "player", Name: w.Player.Name, Weapon: 0}
	return nil
}
