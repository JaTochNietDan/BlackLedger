package core

// What a car is actually worth, said in minutes.
//
// The forecourt quoted a percentage: "journeys take 74% of the time they take on
// foot". That is true and it tells nobody anything, because nobody walks a
// percentage. What a car is worth is the road you keep walking and what it does
// to it — and the honest figure has to include what the car is carrying, since
// plate is weight and a plated car is slower than the table says.

// Worth is one road, timed both ways.
type Worth struct {
	To      string
	ToID    string
	Walking int
	Driving int
	// Plate is how much is on the car being quoted, so the figure can say why
	// it is not the number on the showroom card.
	Plate int
}

// CarWorth times the longest walk from where the player is standing, in a car of
// this kind. The longest, because that is where a car earns its money and where
// the difference is worth quoting.
func (w *World) CarWorth(tier int) Worth {
	from := w.Player.Location
	best := Worth{}
	for _, l := range Locations {
		if l.ID == from || l.District > w.District {
			continue
		}
		walk := TravelMinutes(from, l.ID)
		if walk <= best.Walking {
			continue
		}
		best = Worth{To: l.Name, ToID: l.ID, Walking: walk}
	}
	if best.ToID == "" {
		return best
	}
	best.Plate = min(PlateStages, w.Player.Plate)
	best.Driving = max(5, int(float64(best.Walking)*paceOf(tier, best.Plate)+.5))
	return best
}

// paceOf is what a car of this kind, carrying this much plate, does to a
// journey. In perfect order: the showroom figure, less what the weight costs.
// It is the same arithmetic World.Pace runs on the car the player owns, which
// is the point — a quote that does not agree with what happens is a lie told
// slowly.
func paceOf(tier, plate int) float64 {
	car := VehicleByTier(tier)
	return 1 - (1-car.Pace)*(1-float64(plate)*PlateWeight)
}
