package core

import "fmt"

// A car for one of your own, and plate on it.
//
// "We could probably also extend to be able to provide cars, and armor the cars
// for people in our family to keep them more protected from attacks."
//
// What makes this worth building rather than a number on somebody else's sheet
// is where it lands. Sending one of your own after somebody is the one decision
// in the game whose whole point is that it is not you taking the risk — and
// when it goes wrong there are three ways it ends: they are killed, they are
// taken alive and your name comes out of it, or they get out with nothing.
// The difference between the second and the third is whether there was
// something running at the kerb.

const (
	// TheirCarCost is what putting one of your own on the road costs. It is the
	// forecourt price: the lot does not do favours, though holding the lot
	// yourself keeps the margin in your pocket.
	TheirCarCost = 620
	// TheirCarMinutes is the afternoon it takes to do it properly.
	TheirCarMinutes = 60
	// DrivesOut is how much a car of their own shifts a job that went wrong
	// away from the two endings you do not want.
	DrivesOut = .14
	// PlatedOut is what each stage of plate on it adds to that.
	PlatedOut = .07
)

// yours reports whether this person is one of the player's, alive, and somebody
// the player could actually be standing with.
func (w *World) yours(id string) (*NPC, bool) {
	n := w.NPC(id)
	if n == nil || n.Dead {
		return nil, false
	}
	if n.Faction == "" || n.Faction != w.PlayerOrganizationID() {
		return nil, false
	}
	return n, true
}

// TheirCarReadiness explains why one of your own cannot be put in a car, or
// returns "".
func (w *World) TheirCarReadiness(id string) string {
	n, ok := w.yours(id)
	if !ok {
		return "They are not one of yours"
	}
	if !CarSource(w.Player.Location) {
		return "Cars are sold on a forecourt"
	}
	if n.Location != w.Player.Location {
		return n.Name + " would have to be here to be put in one"
	}
	if n.Car > 0 {
		return n.Name + " already has something to drive"
	}
	if w.Player.Cash < w.TheirCarPrice() {
		return fmt.Sprintf("A car for somebody costs $%d", w.TheirCarPrice())
	}
	return ""
}

// TheirCarPrice is what the lot charges for it, less the margin if the lot is
// yours.
func (w *World) TheirCarPrice() int {
	price := TheirCarCost
	if w.Own(w.Player.Location) {
		price -= TheirCarCost * DealerMargin / 100
	}
	return price
}

// BuyCarFor puts one of your own on the road.
func (w *World) BuyCarFor(id string) error {
	if reason := w.TheirCarReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	n, _ := w.yours(id)
	price := w.TheirCarPrice()
	if err := w.Pay(price); err != nil {
		return err
	}
	lot := w.Properties[w.Player.Location]
	if lot != nil && !w.Own(w.Player.Location) {
		if house := w.faction(lot.Owner); house != nil {
			house.Cash += TheirCarCost * DealerMargin / 100
		}
	}
	n.Car, n.Drove = 1, max(1, w.Minute)
	n.Trust = min(100, n.Trust+4)
	place, _ := PlaceByID(w.Player.Location)
	w.Log(n.Name+" has something to drive",
		fmt.Sprintf("$%d off the lot at %s. They will get where they are sent faster, and away from it quicker.", price, place.Name), "work")
	return nil
}

// TheirPlating is how much plate is on one of your own people's car.
func (w *World) TheirPlating(id string) int {
	n, ok := w.yours(id)
	if !ok || n.Car == 0 {
		return 0
	}
	return min(PlateStages, plateFitted(n.Car)+n.Plate)
}

// TheirPlateReadiness explains why their car cannot be plated, or returns "".
func (w *World) TheirPlateReadiness(id string) string {
	n, ok := w.yours(id)
	if !ok {
		return "They are not one of yours"
	}
	if !CarWorkshop(w.Player.Location) {
		return "Nobody works on cars here"
	}
	if n.Car == 0 {
		return n.Name + " has no car to put anything on"
	}
	if n.Location != w.Player.Location {
		return n.Name + " would have to bring it in"
	}
	if w.TheirPlating(id) >= PlateStages {
		return "There is nothing more to put on it"
	}
	if w.Player.Cash < PlateCost {
		return fmt.Sprintf("Plating costs $%d", PlateCost)
	}
	return ""
}

// PlateTheirCar puts one more stage on one of your own people's car.
func (w *World) PlateTheirCar(id string) error {
	if reason := w.TheirPlateReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	n, _ := w.yours(id)
	if err := w.Pay(PlateCost); err != nil {
		return err
	}
	n.Plate++
	n.Trust = min(100, n.Trust+6)
	place, _ := PlaceByID(w.Player.Location)
	w.Log("On the bench at "+place.Name,
		fmt.Sprintf("%s's car goes out heavier than it came in. Somebody paying for that is not a thing people forget.", n.Name), "work")
	return nil
}

// GetsOut is how much better a job that went wrong ends for whoever you sent,
// because of what was waiting at the kerb.
func (w *World) GetsOut(id string) float64 {
	n, ok := w.yours(id)
	if !ok || n.Car == 0 {
		return 0
	}
	return DrivesOut + float64(w.TheirPlating(id))*PlatedOut
}
