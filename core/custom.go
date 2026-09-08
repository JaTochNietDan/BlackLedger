package core

import "fmt"

// A casino has a float, a nightly handle and somebody with real money at the
// table. A laundry has a number that goes up. Both were called businesses.
//
// Custom is the legitimate trade a place has built up: the people who bring
// their washing there rather than somewhere else, and the cars that come back.
// It is slow to build and quick to lose, it multiplies everything the place
// earns, and it is exactly what is spent when the books are used for something
// they were not meant for.

const (
	// CustomStart is the trade a place has when the player takes it on.
	CustomStart = 50
	// CustomCeiling is as good as a reputation gets, and CustomFloor as bad. A
	// stored zero means a save written before any of this existed, so the floor
	// is one rather than nothing: a shop nobody goes into any more still has a
	// door on it.
	CustomCeiling = 100
	CustomFloor   = 1
	// CustomGain is what a day of running properly is worth.
	CustomGain = 2
	// CustomTroubleLoss is what a day with something wrong costs.
	CustomTroubleLoss = 3
	// CustomLaunderLoss is what a round through the books costs. The machines
	// are always busy and the regulars go elsewhere.
	CustomLaunderLoss = 12
	// CustomRaidLoss is what a police search costs a shop's standing with the
	// people who used to walk past it.
	CustomRaidLoss = 20
	// OrderBonus is what a standing order pays a day while it is filled.
	OrderBonus = 26
	// OrderCapacity is the share of its potential a place must be working at
	// to keep one.
	OrderCapacity = .8
	// OrderCustom is the trade a place needs before anybody offers it one.
	OrderCustom = 65
	// OrderLoss is what losing one costs, on top of the money.
	OrderLoss = 15
)

// Custom is a place's legitimate trade. Absent in saves written before it
// existed, which reads as the trade a place has when it is taken on.
func (w *World) Custom(id string) int {
	prop := w.Properties[id]
	if prop == nil {
		return 0
	}
	if prop.Custom == 0 {
		return CustomStart
	}
	return max(CustomFloor, min(CustomCeiling, prop.Custom))
}

// Trade is what custom multiplies a place's earnings by: half at nothing, half
// again at the ceiling.
func (w *World) TradeMultiplier(id string) float64 {
	if _, running := TradeOf(id); !running {
		return 1
	}
	return float64(50+w.Custom(id)) / 100
}

// ShiftCustom moves a place's trade and keeps it in range.
func (w *World) ShiftCustom(id, why string, delta int) {
	prop := w.Properties[id]
	if prop == nil || delta == 0 {
		return
	}
	before := w.Custom(id)
	prop.Custom = max(CustomFloor, min(CustomCeiling, before+delta))
	if delta < 0 && before >= 50 && prop.Custom < 50 && why != "" {
		place, _ := PlaceByID(id)
		w.Log("The regulars are going elsewhere", fmt.Sprintf("%s at %s. Trade is at %d%% of what it was, and a business without custom earns like one.", why, place.Name, prop.Custom), "business")
	}
}

// CustomDay is a day of trade. A place run properly builds a reputation; one
// left with something wrong loses it faster than it gains.
func (w *World) CustomDay() {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || !w.Own(l.ID) {
			continue
		}
		if _, running := TradeOf(l.ID); !running {
			continue
		}
		switch {
		case prop.Trouble:
			w.ShiftCustom(l.ID, "Something has been wrong here for days", -CustomTroubleLoss)
		case prop.Condition < 50:
			w.ShiftCustom(l.ID, "The place is falling apart", -1)
		case w.Capacity(l.ID) >= OrderCapacity:
			w.ShiftCustom(l.ID, "", CustomGain)
		}
		if operatingMode(prop.Mode).ID == "hard" {
			w.ShiftCustom(l.ID, "It is being run for everything it will give", -2)
		}
	}
}

// OrderReadiness explains why a standing order cannot be taken on, or "".
func (w *World) OrderReadiness(id string) string {
	trade, running := TradeOf(id)
	if !running || !w.Own(id) {
		return "This is not a business of yours"
	}
	if w.Properties[id].Order {
		return "There is already one on the books"
	}
	if w.Custom(id) < OrderCustom {
		return fmt.Sprintf("Nobody offers those to a place with %d%% trade. It takes %d%%", w.Custom(id), OrderCustom)
	}
	if w.Capacity(id) < OrderCapacity {
		return "It could not fill one as it stands"
	}
	_ = trade
	return ""
}

// TakeOrder puts a standing order on the books: steady money for a place that
// can be relied on, and a promise that is expensive to break.
func (w *World) TakeOrder(id string) error {
	if reason := w.OrderReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	w.Properties[id].Order = true
	place, _ := PlaceByID(id)
	w.Log("A standing order at "+place.Name, fmt.Sprintf("Somebody respectable wants the same work every week. $%d a day while it is filled, and %s has to keep working at %d%% to fill it.", OrderBonus, place.Name, int(OrderCapacity*100)), "business")
	return nil
}

// OrderDay pays what the orders pay and takes the ones that cannot be filled.
func (w *World) OrderDay() {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || !prop.Order || !w.Own(l.ID) {
			continue
		}
		if w.Capacity(l.ID) >= OrderCapacity {
			w.Earn(OrderBonus)
			continue
		}
		prop.Order = false
		w.ShiftCustom(l.ID, "A standing order went unfilled", -OrderLoss)
		w.Log("The order is cancelled at "+l.Name, "They came for work that was not ready, and they will not come again. Trade is worse for it.", "business")
	}
}

// CustomDescription is a business's trade, for the interface.
func (w *World) CustomDescription(id string) map[string]any {
	prop := w.Properties[id]
	if prop == nil {
		return nil
	}
	if _, running := TradeOf(id); !running {
		return nil
	}
	return map[string]any{
		"custom": w.Custom(id), "multiplier": w.TradeMultiplier(id),
		"order": prop.Order, "order_pays": OrderBonus,
	}
}
