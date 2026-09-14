package core

import "fmt"

func saleableResidence(id string) bool { return id == "room" || id == "estate" }

// Broker offers exclude acquisition premiums and retain a spread, so selling
// and buying back the same deed cannot manufacture cash.
func (w *World) PropertyOffer(id string) int {
	if !saleableResidence(id) || w.Properties[id] == nil {
		return 0
	}
	base := MarinerFreehold
	if id == "estate" {
		l, _ := PlaceByID(id)
		base = l.Cost
	}
	return base * 65 * max(0, min(100, w.Properties[id].Condition)) / 10000
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
	l, _ := PlaceByID(id)
	if w.Player.Cash < l.Cost {
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
			asking = l.Cost
		}
		out = append(out, map[string]any{"id": id, "name": l.Name, "owned": w.Own(id), "available": w.CanAcquire(id), "holder": w.HolderName(id), "asking": asking, "offer": w.PropertyOffer(id), "condition": p.Condition, "residents": len(w.Residents(id)), "home": w.Player.Home == id, "locked": l.District > w.District})
	}
	return out
}
