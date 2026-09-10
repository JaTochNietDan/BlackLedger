package core

import "fmt"

// Petrol. A car cost money every day and needed nothing: it drove for ever on
// an upkeep figure, which meant a filling station would have been a shop
// selling something no vehicle in this city consumed.
//
// A tank now goes down with the driving and has to be filled somewhere, which
// is the whole of what makes a forecourt with two pumps on it a business. The
// city's own drivers run dry as well, and when they do they go and buy a tank
// off whoever holds the station — the same shape as a broken windscreen taking
// somebody to a garage.

const (
	// FuelFull is a full tank, in the units the gauge is read in.
	FuelFull = 100
	// FuelPerHour is what an hour behind the wheel burns. A journey across
	// Bellwether is under an hour, so a tank is most of a week of ordinary
	// running and a bad day of chasing people around the city.
	FuelPerHour = 4
	// FuelPrice is what a full tank costs at the pumps, and PumpTrade is what
	// one sale is worth to the station's custom.
	FuelPrice = 26
	PumpTrade = 2
	// FillMinutes is how long it takes to put one in.
	FillMinutes = 15
)

// Fuel is what is left in the player's tank. Nothing without a car, because a
// gauge on a car nobody owns is a number about nothing.
func (w *World) Fuel() int {
	if w.Player.Car == 0 {
		return 0
	}
	// A tank nobody has ever touched is a tank that came with the car. That is
	// what Fuelled is for: without it, every save written before petrol existed
	// — and every car handed to somebody by anything other than a forecourt —
	// would be standing empty at the kerb the moment this shipped.
	if w.Player.Fuelled == 0 {
		return FuelFull
	}
	return max(0, min(FuelFull, w.Player.Fuel))
}

// SettleFuel fills the tank of a car that has never had one filled. It runs
// when a car is bought and again on every load, so a save written before petrol
// existed reads as a car somebody had been keeping on the road all along.
//
// Fuelled is when the tank was last touched, and it is the whole reason this
// cannot quietly undo running dry: a tank at nothing with a minute on it was
// emptied, and a tank at nothing with no minute never existed.
func (w *World) SettleFuel() {
	if w.Player.Car == 0 || w.Player.Fuelled != 0 {
		return
	}
	w.Player.Fuel, w.Player.Fuelled = FuelFull, max(1, w.Minute)
}

// Burn takes what a stretch of driving costs out of the tank.
func (w *World) Burn(minutes int) {
	if w.Player.Car == 0 || minutes <= 0 {
		return
	}
	// Read through Fuel(), not off the field: a tank nobody has touched reads
	// full and holds nothing, so taking the journey off the raw number emptied
	// every car in the city the first time it was driven.
	was := w.Fuel()
	spent := (minutes*FuelPerHour + 59) / 60
	w.Player.Fuel = max(0, was-spent)
	w.Player.Fuelled = max(1, w.Minute)
	if was > 0 && w.Player.Fuel == 0 {
		w.Log("Out of petrol", fmt.Sprintf("%s is standing wherever it stopped. It is a walk from here until somebody puts a tank in it.",
			VehicleByTier(w.Player.Car).Label), "personal")
	}
}

// Pumps reports whether a place sells petrol.
func Pumps(id string) bool {
	place, ok := PlaceByID(id)
	return ok && place.Kind == "filling"
}

// FuelFee is what a tank costs here. Your own pumps are your own petrol, so it
// is only what the petrol cost the station.
func (w *World) FuelFee(id string) int {
	short := FuelFull - w.Fuel()
	price := FuelPrice * short / FuelFull
	if w.Own(id) {
		price = price / 3
	}
	return price
}

// FillReadiness explains why the tank cannot be filled, or returns "".
func (w *World) FillReadiness(id string) string {
	if !Pumps(id) {
		return "There are no pumps here"
	}
	if w.Player.Car == 0 {
		return "There is nothing of yours to put it in"
	}
	if w.Fuel() >= FuelFull {
		return "The tank is full"
	}
	if w.Player.Cash < w.FuelFee(id) {
		return "Not enough cash"
	}
	return ""
}

// FillUp buys a tank.
func (w *World) FillUp(id string) error {
	if reason := w.FillReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	fee := w.FuelFee(id)
	if err := w.Pay(fee); err != nil {
		return err
	}
	w.Player.Fuel, w.Player.Fuelled = FuelFull, max(1, w.Minute)
	place, _ := PlaceByID(id)
	if house := w.faction(w.Properties[id].Owner); house != nil {
		house.Cash += fee
	}
	w.ShiftCustom(id, "cars in off the road", PumpTrade)
	who := fmt.Sprintf("$%d at %s", fee, place.Name)
	if w.Own(id) {
		who = "Your own pumps, at what the petrol cost"
	}
	w.Log("A full tank", fmt.Sprintf("%s. %s is good for another week of it.", who, VehicleByTier(w.Player.Car).Label), "personal")
	return nil
}

// theStation is where the city buys its petrol: the player's own pumps if they
// hold any, and otherwise the first on the list, so errands and reports do not
// depend on map iteration order. The same shape as theGarage and theForecourt,
// for the same reason.
func (w *World) theStation() string {
	first := ""
	for _, l := range Locations {
		if l.Kind != "filling" || w.Properties[l.ID] == nil {
			continue
		}
		if w.Own(l.ID) {
			return l.ID
		}
		if first == "" {
			first = l.ID
		}
	}
	return first
}

// DryDay is the city's own cars running low. One driver a day finds the needle
// on the pin, which is enough to keep a forecourt busy without every car in
// Bellwether queueing at the same pump.
func (w *World) DryDay() {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Car == 0 || n.Dry || n.Purse < FuelPrice {
			continue
		}
		if w.WorldRandom() < .5 {
			continue
		}
		n.Dry = true
		return
	}
}

// fillFor is one tank sold to somebody standing on the forecourt, called when
// they walk on and again when the day turns over — the same rule as the bench
// and the lot, and for the same reason: at midnight the pumps are shut and
// everybody who came that day has gone home.
func (w *World) fillFor(n *NPC) {
	station := w.theStation()
	if station == "" || n == nil || n.Dead || !n.Dry || n.Car == 0 || n.Location != station {
		return
	}
	if n.Purse < FuelPrice {
		return
	}
	n.Purse -= FuelPrice
	n.Dry = false
	if house := w.faction(w.Properties[station].Owner); house != nil {
		house.Cash += FuelPrice
	}
	w.ShiftCustom(station, "cars in off the road", PumpTrade)
	if w.Own(station) {
		w.Earn(FuelPrice)
		place, _ := PlaceByID(station)
		w.Log("A tank sold at "+place.Name, fmt.Sprintf("%s filled up on the way through. $%d on the counter.", n.Name, FuelPrice), "business")
	}
}

// PumpDay is the sweep for anybody left standing on a forecourt at the turn of
// the day.
func (w *World) PumpDay() {
	station := w.theStation()
	if station == "" {
		return
	}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Location == station {
			w.fillFor(n)
		}
	}
}
