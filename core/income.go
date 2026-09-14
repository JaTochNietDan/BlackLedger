package core

// HourlyIncome is the owner's current rate, shared by the clock and accounts.
// It includes actual footfall and the share collected through posted crew.
func (w *World) HourlyIncome(id string) float64 {
	prop := w.Properties[id]
	if prop == nil || !w.Own(id) || id == "room" || id == "apartment" {
		return 0
	}
	return float64(prop.Income*prop.Condition) / 100 * operatingMode(prop.Mode).Take *
		w.Capacity(id) * w.TradeMultiplier(id) * w.RoomTrade(id) *
		(1 + w.LicenceTake()) * w.CollectionShare(id)
}
