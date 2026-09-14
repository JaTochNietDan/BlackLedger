package core

import (
	"fmt"
	"strings"
)

const (
	BuildingDriveByCost    = 35
	BuildingDriveByMinutes = 10
	BuildingDriveByHeat    = 24
)

// BuildingDriveByReadiness keeps the driver distinct from the shooter. Like
// delegated jobs, this uses the first crew member; unlike a remote errand,
// driving together requires that person to be at the same address already.
func (w *World) BuildingDriveByReadiness(id string) string {
	if _, ok := PlaceByID(id); !ok {
		return "There is no business to attack here"
	}
	p := w.Properties[id]
	if p == nil || p.Income <= 0 {
		return "There is no business to attack here"
	}
	if !w.Player.Alive {
		return "This life is over"
	}
	if w.Held() {
		return "You are being held at Ward Street Station"
	}
	if w.Own(id) {
		return "This is yours"
	}
	if w.Player.Location != id {
		return "You would have to be there"
	}
	if p.Condition <= 0 {
		return "The premises are already wrecked"
	}
	for _, fire := range w.ActiveBuildingFires() {
		if fire.Target == id && w.Minute < fire.ExtinguishedAt {
			return "The premises are already burning"
		}
	}
	if w.Player.Weapon < 1 || w.Player.Weapon > 3 {
		return "You need a firearm equipped"
	}
	if !w.Driving() {
		return "You need your own running car with petrol"
	}
	if reason := w.DelegateReadiness(); reason != "" {
		return reason
	}
	driver := w.NPC(w.Player.Crew[0].ID)
	if driver.Location != id {
		return driver.Name + " must be here to drive"
	}
	if w.Player.Cash < BuildingDriveByCost {
		return "Not enough cash"
	}
	return ""
}

// BuildingDriveBy resolves a property attack. Command dispatch and its browser
// scene are deliberately not enabled yet; this method does not advance time.
func (w *World) BuildingDriveBy(id string) error {
	if reason := w.BuildingDriveByReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(BuildingDriveByCost); err != nil {
		return err
	}
	p := w.Properties[id]
	place, _ := PlaceByID(id)
	driver := w.NPC(w.Player.Crew[0].ID)
	car := VehicleByTier(w.Player.Car)
	before := p.Condition
	damage := min(before, 8+6*w.Player.Weapon+int(w.WorldRandom()*6))
	p.Condition -= damage
	p.Supply = max(0, p.Supply-damage/3)
	p.Trouble = true
	w.Burn(BuildingDriveByMinutes)
	w.Player.Heat = min(100, w.Player.Heat+BuildingDriveByHeat)
	if owner := w.faction(p.Owner); owner != nil {
		owner.Goodwill = max(-100, owner.Goodwill-25)
		w.RetaliationFrom(owner.ID)
	}
	body := fmt.Sprintf("%s drove your %s past %s while you fired from the passenger seat. Gunfire damaged the premises by %d points; condition is now %d%%.", driver.Name, strings.TrimPrefix(car.Label, "A "), place.Name, damage, p.Condition)
	w.Log("Drive-by at "+place.Name, body, "danger")
	headline := "GUNFIRE AT " + upper(place.Name)
	w.Report("attack", headline, fmt.Sprintf("Shots fired from a passing car damaged %s. The car left the scene.", place.Name))
	w.Witness("driveby-building", id, body, headline, driver.ID)
	cue := &w.VisualCues[len(w.VisualCues)-1]
	cue.Gravity = 7
	cue.Attacker = &CueAttacker{ID: "player", Name: w.Player.Name, Weapon: w.Player.Weapon}
	cue.DriveBy = &CueDriveBy{
		Driver:  CueActor{ID: driver.ID, Name: driver.Name},
		Vehicle: car.Label, VehicleTier: car.Tier,
		ConditionBefore: before, ConditionAfter: p.Condition,
	}
	return nil
}
