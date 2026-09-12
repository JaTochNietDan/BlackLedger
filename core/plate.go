package core

import "fmt"

// Plating a car.
//
// The top of the range has been called an armoured Packard since the day the
// list was written, and no rule in the game ever read that word: the plate was
// a sentence in a description and nothing else. And the one stretch of the city
// where nothing protects anybody is the street between two addresses — which is
// exactly where a car is.
//
// So a garage will fit it, in stages, to the car you actually have. It costs
// money, it costs time, and it costs speed, because plate is weight. What it
// buys is the only cover there is out there.

const (
	// PlateStages is how far a car can be built up. Two: the doors, and then
	// the glass. Past that it stops being a car.
	PlateStages = 2
	// PlateCost is what one stage costs at a garage.
	PlateCost = 900
	// PlateMinutes is how long the car is on the bench for it.
	PlateMinutes = 240
	// PlateCover is what one stage is worth against somebody who pulls level
	// with you on the street.
	PlateCover = .13
	// PlateWeight is the share of the car's advantage each stage costs. Plate
	// is weight, and weight is slower.
	PlateWeight = .09
)

// plateFitted is what a car of this kind carries before anybody works on it.
// The Packard is described as armoured, so it is.
// Plating is how much plate is on the car the player is driving.
//
// No car comes with any. The top of the range was called an armoured Packard
// and carried two stages nobody fitted, which made the only real way to get
// cover off a lot rather than off a bench — buying a car and arming a car were
// the same purchase, and there was nothing to decide. A Packard is a big quick
// car now and costs what one costs; plate goes on at a garage, to whatever the
// player actually drives, and the weight is the price of it.
func (w *World) Plating() int {
	if w.Player.Car == 0 {
		return 0
	}
	return min(PlateStages, w.Player.Plate)
}

// PlateReadiness explains why plate cannot be fitted here, or returns "".
func (w *World) PlateReadiness(id string) string {
	if !CarWorkshop(id) {
		return "Nobody works on cars here"
	}
	if w.Player.Car == 0 {
		return "There is no car of yours to plate"
	}
	if w.Plating() >= PlateStages {
		return "There is nothing more to put on it"
	}
	if w.Player.Cash < PlateCost {
		return fmt.Sprintf("Plating costs $%d", PlateCost)
	}
	return ""
}

// PlateLabel is what the next stage is called, so the button says what is being
// fitted rather than "upgrade armour".
func (w *World) PlateLabel() string {
	if w.Plating() == 0 {
		return "Plate the doors"
	}
	return "Fit glass that has stopped things"
}

// FitPlate puts one more stage on the car.
func (w *World) FitPlate(id string) error {
	if reason := w.PlateReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(PlateCost); err != nil {
		return err
	}
	w.Player.Plate++
	place, _ := PlaceByID(id)
	car := VehicleByTier(w.Player.Car)
	w.Log("On the bench at "+place.Name,
		fmt.Sprintf("%s. %s is heavier for it and slower with it, and there is something between you and the street now.",
			w.PlateLabel(), car.Label), "personal")
	return nil
}

// PlateCoverHere is what the plate is worth where the player is standing. On
// the street, in a car, it is the only cover there is. Anywhere else it is a
// car parked outside and worth nothing to somebody who has walked in the door.
func (w *World) PlateCoverHere() float64 {
	if !w.InTransit() || !w.Driving() {
		return 0
	}
	return float64(w.Plating()) * PlateCover
}

// PlateDrag is what the weight costs, as a share of what the car was worth.
func (w *World) PlateDrag() float64 {
	return float64(w.Plating()) * PlateWeight
}
