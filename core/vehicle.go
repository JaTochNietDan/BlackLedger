package core

import "fmt"

// Everybody in this city walks or takes a streetcar, and the clock is the one
// thing nobody can buy more of: every journey is time in which businesses earn,
// rivals move and arrangements come due. A car buys that time back.
//
// It is also the most visible thing a person can own. A car is registered, it
// is described by witnesses, it sits outside places it has no business being,
// and somebody else can take it or burn it. So it is time and a hiding place
// against a standing cost, a trail, and something that can be taken.

// Vehicle is what the player drives.
type Vehicle struct {
	Tier int
	// Label is the car, in words.
	Label string
	// Detail is why anybody would buy it.
	Detail string
	// Cost is the price of moving up to it.
	Cost int
	// Pace is the share of the walking time a journey takes at full condition.
	Pace float64
	// Compartment is how many units of stock a search will not find.
	Compartment int
	// Upkeep is what it costs a day in fuel, a space and somebody's attention.
	Upkeep int
	// Trail is how much more likely it is that somebody works out who did
	// something, because a car parked outside is a thing people describe.
	Trail int
}

var vehicles = []Vehicle{
	{Tier: 0, Label: "On foot and by streetcar", Pace: 1},
	{Tier: 1, Label: "A used Ford", Detail: "Rattles, starts most mornings, and gets you across town before the day is gone.", Cost: 620, Pace: .74, Upkeep: 4, Trail: 1},
	{Tier: 2, Label: "A Hudson with a false floor", Detail: "Respectable from the outside, and there is room under it for things that should not be in the boot.", Cost: 1750, Pace: .6, Compartment: 20, Upkeep: 7, Trail: 2},
	{Tier: 3, Label: "An armoured Packard", Detail: "Plated doors, glass that has stopped things before, and a driver's seat people have walked away from.", Cost: 4400, Pace: .5, Compartment: 35, Upkeep: 12, Trail: 3},
}

const (
	// CarWear is the condition a car loses on a day it is driven at all.
	CarWear = 2
	// CarService is what putting one right costs anywhere but a garage of your
	// own, where the people who work for you do it.
	CarService = 90
	// CarServiceMinutes is how long that takes.
	CarServiceMinutes = 90
	// Wreck is the condition below which a car is not worth the space it takes
	// up: it makes no journey faster and hides nothing.
	Wreck = 30
)

// VehicleByTier is what the player drives, clamped so an unknown save is safe.
func VehicleByTier(tier int) Vehicle { return vehicles[max(0, min(tier, len(vehicles)-1))] }

func nextVehicle(tier int) (Vehicle, bool) {
	if tier+1 >= len(vehicles) {
		return Vehicle{}, false
	}
	return vehicles[tier+1], true
}

// CarCondition is the state of what the player drives. Saves written before
// cars existed carry a zero and no car, which is nothing to keep running.
func (w *World) CarCondition() int {
	if w.Player.Car == 0 {
		return 0
	}
	return max(0, min(100, w.Player.CarWear))
}

// Driving reports whether the player has something that actually runs.
func (w *World) Driving() bool {
	return w.Player.Car > 0 && w.CarCondition() >= Wreck
}

// Pace is the share of a walking journey the player's current transport takes.
// A car in poor order is worth part of what it was, and a wreck is worth
// nothing at all.
func (w *World) Pace() float64 {
	if !w.Driving() {
		return 1
	}
	car := VehicleByTier(w.Player.Car)
	// Between a wreck and perfect order, the benefit scales with how well it is
	// running rather than switching off at a threshold.
	kept := float64(w.CarCondition()-Wreck) / float64(100-Wreck)
	return 1 - (1-car.Pace)*kept
}

// Journey is how long it takes the player to get somewhere, which is the whole
// reason to own a car.
func (w *World) Journey(from, to string) int {
	return max(5, int(float64(TravelMinutes(from, to))*w.Pace()+.5))
}

// Concealed is how much stock a search will not find, because it is under the
// floor of something parked two streets away.
func (w *World) Concealed() int {
	hidden := 0
	if w.Driving() {
		hidden += VehicleByTier(w.Player.Car).Compartment
	}
	// A dry cellar under the place you live hides the same way a false floor
	// does, right up until somebody serves a warrant on the place you live.
	if w.Fitted("cellar") {
		hidden += CellarHold
	}
	return hidden
}

// Exposed is the stock that is not hidden: what draws attention every day and
// what a search actually takes.
func (w *World) Exposed() int { return max(0, w.Carrying()-w.Concealed()) }

// CarUpkeep is what the car costs a day. A garage of your own is people who
// work on it anyway, so it costs half.
func (w *World) CarUpkeep() int {
	if w.Player.Car == 0 {
		return 0
	}
	upkeep := VehicleByTier(w.Player.Car).Upkeep
	if w.Own("garage") {
		upkeep = (upkeep + 1) / 2
	}
	return upkeep
}

// CarTrail is how much more likely somebody is to work out who did something,
// because a car outside is a thing witnesses describe.
func (w *World) CarTrail() int {
	if !w.Driving() {
		return 0
	}
	return VehicleByTier(w.Player.Car).Trail
}

// Damage is what a bad night does to a car: a chase, a beating in the street,
// somebody taking an interest in it while it is parked.
func (w *World) Damage(amount int) {
	if w.Player.Car == 0 || amount <= 0 {
		return
	}
	before := w.CarCondition()
	w.Player.CarWear = max(0, before-amount)
	if before >= Wreck && w.Player.CarWear < Wreck {
		w.Log("It will not start", fmt.Sprintf("%s is finished until somebody works on it. You are walking.", VehicleByTier(w.Player.Car).Label), "personal")
	}
}

// CarDay is the wear of driving a car around this city, and what it costs to
// keep on the road.
func (w *World) CarDay() {
	if w.Player.Car == 0 {
		return
	}
	w.Damage(CarWear)
}

// CarSource is where cars change hands: a motor works, which is the one place
// in this city that has any.
func CarSource(id string) bool { return id == "garage" }

// CarReadiness explains why the next car cannot be bought, or returns "".
func (w *World) CarReadiness() string {
	if !CarSource(w.Player.Location) {
		return "Nobody sells cars here"
	}
	next, ok := nextVehicle(w.Player.Car)
	if !ok {
		return "There is nothing better on the lot"
	}
	if w.Player.Cash < next.Cost {
		return "Not enough cash"
	}
	return ""
}

// BuyVehicle takes the next step up. The old one goes toward it, which is why
// each step costs what it costs rather than the difference.
func (w *World) BuyVehicle() error {
	if reason := w.CarReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	next, _ := nextVehicle(w.Player.Car)
	if err := w.Pay(next.Cost); err != nil {
		return err
	}
	w.Player.Car, w.Player.CarWear = next.Tier, 100
	w.Log("Off the lot at Russo Motor Works", fmt.Sprintf("%s, $%d. $%d a day to keep on the road. %s", next.Label, next.Cost, w.CarUpkeep(), next.Detail), "personal")
	return nil
}

// ServiceReadiness explains why a car cannot be worked on, or returns "".
func (w *World) ServiceReadiness(id string) string {
	if !CarSource(id) {
		return "Nobody works on cars here"
	}
	if w.Player.Car == 0 {
		return "There is nothing of yours to work on"
	}
	if w.CarCondition() >= 100 {
		return "It is running as well as it ever will"
	}
	if w.Player.Cash < w.ServiceFee() {
		return "Not enough cash"
	}
	return ""
}

// ServiceFee is nothing at a motor works of your own: the people there work on
// cars all day and one of them is yours.
func (w *World) ServiceFee() int {
	if w.Own("garage") {
		return 0
	}
	return CarService
}

// Service puts a car back on the road.
func (w *World) Service(id string) error {
	if reason := w.ServiceReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	fee := w.ServiceFee()
	if err := w.Pay(fee); err != nil {
		return err
	}
	w.Player.CarWear = min(100, w.CarCondition()+55)
	who := fmt.Sprintf("$%d at Russo Motor Works", fee)
	if fee == 0 {
		who = "Your own people at Russo Motor Works"
	}
	w.Log("Back on the road", fmt.Sprintf("%s. %s is at %d of 100.", who, VehicleByTier(w.Player.Car).Label, w.Player.CarWear), "personal")
	return nil
}

// LoseCar is what happens when somebody takes it or burns it. The player is
// walking again, and everything under the floor goes with it.
func (w *World) LoseCar(reason string) bool {
	if w.Player.Car == 0 {
		return false
	}
	label := VehicleByTier(w.Player.Car).Label
	w.Player.Car, w.Player.CarWear = 0, 0
	w.Log(label+" is gone", reason+" You are walking again, and anything that was under the floor of it went too.", "danger")
	return true
}

// VehicleDescription is what the player drives, for the interface.
func (w *World) VehicleDescription() map[string]any {
	return map[string]any{
		"car":       VehicleByTier(w.Player.Car).Label,
		"condition": w.CarCondition(),
		"pace":      w.Pace(),
		"concealed": w.Concealed(),
		"exposed":   w.Exposed(),
		"upkeep":    w.CarUpkeep(),
		"running":   w.Driving(),
	}
}
