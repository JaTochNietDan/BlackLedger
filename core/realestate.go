package core

import "fmt"

// DistrictPropertyPressure keeps sub-point decay so frequent incidents cannot
// postpone recovery by repeatedly discarding partial days. It survives lives.
type DistrictPropertyPressure struct {
	Units  int `json:"units"`
	Minute int `json:"minute"`
}

// Each 1,440 pressure units reduces prices by one percentage point. Pressure
// recovers one unit per game minute, bounded at a forty-percent discount.
func (w *World) propertyPressure(district int) int {
	p := w.PropertyPressure[district]
	return max(0, min(40*1440, p.Units)-max(0, w.Minute-p.Minute))
}

func (w *World) recordPropertyIncident(kind, id string) {
	points := map[string]int{"killing": 8, "explosion": 10, "incendiary": 7, "gunfight": 6, "driveby-building": 6, "attack": 3, "robbery": 2}[kind]
	place, ok := PlaceByID(id)
	if !ok || points == 0 {
		return
	}
	units := min(40*1440, w.propertyPressure(place.District)+points*1440)
	if w.PropertyPressure == nil {
		w.PropertyPressure = map[int]DistrictPropertyPressure{}
	}
	w.PropertyPressure[place.District] = DistrictPropertyPressure{Units: units, Minute: w.Minute}
}

func (w *World) NeighborhoodPropertyIndex(id string) int {
	place, ok := PlaceByID(id)
	if !ok {
		return 100
	}
	return 100 - (w.propertyPressure(place.District)+1439)/1440
}

// ResidencePrice is shared by standalone deed and buy-and-move commands.
func (w *World) ResidencePrice(id string) int {
	place, ok := PlaceByID(id)
	if !ok {
		return 0
	}
	base := place.Cost
	if id == "room" {
		base = MarinerFreehold
	}
	return base * w.NeighborhoodPropertyIndex(id) / 100
}

func saleableResidence(id string) bool { return id == "room" || id == "estate" }

// Broker offers exclude acquisition premiums and retain a spread, so selling
// and buying back the same deed cannot manufacture cash.
func (w *World) PropertyOffer(id string) int {
	if !saleableResidence(id) || w.Properties[id] == nil {
		return 0
	}
	return w.ResidencePrice(id) * 65 * max(0, min(100, w.Properties[id].Condition)) / 10000
}

func (w *World) SellPropertyReadiness(id string) string {
	if !saleableResidence(id) || !w.Own(id) {
		return "This is not a residence you own"
	}
	if w.PropertyOffer(id) <= 0 {
		return "Repair the building before seeking an offer"
	}
	return ""
}

func (w *World) SellProperty(id string) error {
	if why := w.SellPropertyReadiness(id); why != "" {
		return fmt.Errorf("%s", why)
	}
	price := w.PropertyOffer(id)
	p := w.Properties[id]
	p.Owner = "independent"
	p.Posted = ""
	w.Player.Cash += price // Asset sale, not new earnings or job progress.
	detail := fmt.Sprintf("The broker paid $%d for %s. The deed, fixtures, staff and tenant accounts transfer to the buyer.", price, placeName(id))
	if w.Player.Home == id {
		detail += fmt.Sprintf(" You remain in your home as a renter at $%d a day.", w.HomeCost(id))
	}
	w.Log("A deed changes hands", detail, "business")
	return nil
}

func (w *World) BuyResidenceReadiness(id string) string {
	if id != "estate" {
		return "This residence is not offered separately"
	}
	if w.Own(id) {
		return "This deed is already yours"
	}
	if w.Properties[id] == nil || !w.CanAcquire(id) {
		return "The owner is not offering this residence"
	}
	if w.Player.Cash < w.ResidencePrice(id) {
		return "Not enough cash"
	}
	return ""
}

func (w *World) PropertyMarket() []map[string]any {
	out := []map[string]any{}
	for _, id := range []string{"room", "estate"} {
		l, ok := PlaceByID(id)
		p := w.Properties[id]
		if !ok || p == nil {
			continue
		}
		asking := AcquisitionCost(w, id)
		if id == "estate" {
			asking = w.ResidencePrice(id)
		}
		out = append(out, map[string]any{"id": id, "name": l.Name, "owned": w.Own(id), "available": w.CanAcquire(id), "holder": w.HolderName(id), "asking": asking, "offer": w.PropertyOffer(id), "condition": p.Condition, "neighborhood_index": w.NeighborhoodPropertyIndex(id), "residents": len(w.Residents(id)), "home": w.Player.Home == id, "locked": l.District > w.District})
	}
	return out
}
