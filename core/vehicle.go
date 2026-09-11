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
	// And there has to be something in it. A car with a dry tank is where you
	// left it, which is the whole reason a filling station is a business.
	return w.Player.Car > 0 && w.CarCondition() >= Wreck && w.Fuel() > 0
}

// Pace is the share of a walking journey the player's current transport takes.
// A car in poor order is worth part of what it was, and a wreck is worth
// nothing at all.
// CabPace is what a cab does against walking. Slower than anything the player
// could own, because it is somebody else's car on somebody else's route and it
// stops for other fares — but it is a ride, and the alternative is the
// streetcar.
const CabPace = .85

// RidingWithTheCabs reports whether the player is getting about in their own
// cabs: a yard of cars and drivers they hold, when whatever they own themselves
// cannot take them.
//
// The cabstand was an address that paid and did nothing else, which is true of
// nine of this city's thirteen trades. A cab yard is the one whose whole
// business is moving somebody across town, so it is the one where what it does
// for the rest of what you hold writes itself: you always have a ride.
func (w *World) RidingWithTheCabs() bool {
	return !w.Driving() && w.OwnsKind("cabs")
}

func (w *World) Pace() float64 {
	if !w.Driving() {
		if w.RidingWithTheCabs() {
			return CabPace
		}
		return 1
	}
	car := VehicleByTier(w.Player.Car)
	// Between a wreck and perfect order, the benefit scales with how well it is
	// running rather than switching off at a threshold.
	kept := float64(w.CarCondition()-Wreck) / float64(100-Wreck)
	// Plate is weight, and weight costs part of what the car was worth.
	return 1 - (1-car.Pace)*kept*(1-w.PlateDrag())
}

// Journey is how long it takes the player to get somewhere, which is the whole
// reason to own a car.
func (w *World) Journey(from, to string) int {
	// Standing where you already are is no journey. Everything else is at
	// least five minutes, because two addresses on the same street are still
	// two addresses — but the floor used to apply to the trip nobody takes, so
	// the city reported ten minutes to reach the room the player was in.
	if from == to {
		return 0
	}
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
	// And the trades whose whole business is things sitting somewhere without
	// being looked at. A yard full of trucks and a cold room are that; a revue
	// bar is not.
	for _, l := range Locations {
		if !w.Own(l.ID) {
			continue
		}
		if trade, runs := TradeOf(l.ID); runs {
			hidden += trade.Hides
		}
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
	if w.OwnsKind("garage") {
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

// CarSource is where cars change hands: a forecourt. It used to be the motor
// works, and the comment above it said "the one place in this city that has
// any" — which was true of a city with one garage and nowhere to buy a car.
// A garage repairs what you already have; a dealer is where it came from.
func CarSource(id string) bool {
	place, ok := PlaceByID(id)
	return ok && place.Kind == "dealer"
}

// CarWorkshop is where a car is worked on: a garage. This used to be CarSource,
// the same function that answered where cars were SOLD — one question doing two
// jobs, which nobody noticed while the answer to both was "the garage". Moving
// sales to a forecourt moved servicing with it, and a motor works could no
// longer touch a car.
func CarWorkshop(id string) bool {
	place, ok := PlaceByID(id)
	return ok && place.Kind == "garage"
}

// DealerMargin is the share of a car's price that stays with the forecourt
// rather than going wherever cars come from. It is what a lot is worth holding.
const DealerMargin = 25

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
	// A car is sold by somebody. The forecourt keeps its margin, and if the
	// player holds the lot they are buying from themselves — the margin never
	// leaves their pocket, so the car costs them less.
	lot := w.Properties[w.Player.Location]
	margin := next.Cost * DealerMargin / 100
	price := next.Cost
	if w.Own(w.Player.Location) {
		price -= margin
	}
	if err := w.Pay(price); err != nil {
		return err
	}
	if lot != nil && !w.Own(w.Player.Location) {
		if house := w.faction(lot.Owner); house != nil {
			house.Cash += margin
		}
	}
	// Plate is fitted to a car, not to a person. What you had on the last one
	// is on the last one.
	w.Player.Car, w.Player.CarWear, w.Player.Plate = next.Tier, 100, 0
	// A car off the lot comes with a tank in it.
	w.Player.Fuel, w.Player.Fuelled = FuelFull, max(1, w.Minute)
	w.Log("Off the lot at Russo Motor Works", fmt.Sprintf("%s, $%d. $%d a day to keep on the road. %s", next.Label, next.Cost, w.CarUpkeep(), next.Detail), "personal")
	return nil
}

// ServiceReadiness explains why a car cannot be worked on, or returns "".
func (w *World) ServiceReadiness(id string) string {
	if !CarWorkshop(id) {
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
	if w.OwnsKind("garage") {
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
		"condition": w.CarCondition(), "fuel": w.Fuel(), "tank": FuelFull,
		"pace":      w.Pace(),
		"concealed": w.Concealed(),
		"exposed":   w.Exposed(),
		"upkeep":    w.CarUpkeep(),
		"running":   w.Driving(),
		"plate":     w.Plating(), "plate_max": PlateStages,
	}
}

// DrivesAt is the standing at which somebody in this city keeps a car. Below it
// they walk and take the streetcar like everybody else.
const DrivesAt = RankSoldier

// WouldDrive reports whether this person is somebody who would keep a car: what
// they are worth, and what they are. A family lieutenant drives. Somebody
// mending nets at the docks does not, however long they save.
func (w *World) WouldDrive(n *NPC) bool {
	if n == nil || n.Dead {
		return false
	}
	if IsOfficial(n.ID) {
		return true // a man with a title and an arrangement has a car
	}
	return n.Rank >= DrivesAt
}

// SettleCars puts a car under the people who would have one. It runs when a
// city is made and again when one is loaded, so a save written before the city
// drove reads as a city that always did — and it only ever settles somebody who
// has never had one, so it cannot quietly replace a car that was taken.
func (w *World) SettleCars() {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Car != 0 || n.Drove != 0 || !w.WouldDrive(n) {
			continue
		}
		n.Car, n.Drove = 1+n.Rank/RankLieutenant, max(1, w.Minute)
	}
}

// theForecourt is where cars change hands: the first lot in the list, so
// reports and errands do not depend on map iteration order. It is the same
// shape as theGarage, and for the same reason — the city needs one answer to
// "where would somebody go for this".
func (w *World) theForecourt() string {
	for _, l := range Locations {
		if l.Kind == "dealer" && w.Properties[l.ID] != nil {
			return l.ID
		}
	}
	return ""
}

// CarTrade is the city buying cars. Somebody who would drive, has none and can
// afford one goes to a forecourt and buys, and whoever holds that forecourt
// takes the margin — the same margin the player pays, because it is the same
// transaction seen from the other side.
//
// It used to be one sale a day at most, on the grounds that a city where
// everybody replaces a car on the same morning is a city where nothing was ever
// taken from anybody. The walk is what makes that true now: a person has to
// cross Bellwether to the lot and be standing on it, and eight people who lost
// cars in a raid arrive over days rather than together. Rationing it on top of
// that left somebody queueing at a forecourt for a month.
func (w *World) CarTrade() {
	lot := w.theForecourt()
	if lot == "" {
		return
	}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Location == lot {
			w.sellCarTo(n)
		}
	}
}

// sellCarTo is one sale, to somebody standing on the forecourt. Called when
// they walk on and again when the day turns over: a sweep at midnight finds a
// lot everybody left hours ago, which is how a man with four thousand dollars
// in his pocket spent a month walking past the place every afternoon.
func (w *World) sellCarTo(n *NPC) {
	lot := w.theForecourt()
	price := VehicleByTier(1).Cost
	if lot == "" || n == nil || n.Dead || n.Car != 0 || n.Location != lot {
		return
	}
	if !w.WouldDrive(n) || n.Purse < price {
		return
	}
	n.Purse -= price
	n.Car, n.Drove = 1, max(1, w.Minute)
	if house := w.faction(w.Properties[lot].Owner); house != nil {
		house.Cash += price * DealerMargin / 100
	}
	if w.Own(lot) {
		w.Earn(price * DealerMargin / 100)
		place, _ := PlaceByID(lot)
		w.Log("A car sold at "+place.Name, fmt.Sprintf("%s bought one off the lot. $%d of it is yours.",
			n.Name, price*DealerMargin/100), "business")
	}
}
