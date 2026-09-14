package core

import "fmt"

const (
	RushOrderPay     = 80
	RushOrderMinutes = 60
	RushOrderSupply  = 2
)

// One external customer's linen order per day. Its reservation belongs to the
// premises, so a new protagonist or a change of owner cannot create it again.
func (w *World) RushOrderReadiness(id string) string {
	if id != "laundry" {
		return "The hotel leaves its linen orders at Bluebird Laundry"
	}
	prop := w.Properties[id]
	if prop == nil {
		return "The laundry is unavailable"
	}
	if w.Own(id) {
		return "You own the laundry; its work already earns through your business accounts"
	}
	hour := w.Minute % 1440
	if hour < 8*60 || hour+RushOrderMinutes > 18*60 {
		return "Rush orders are taken between 08:00 and 17:00"
	}
	if prop.RushOrderDay == w.Minute/1440+1 {
		return "Today's hotel order has already been taken"
	}
	if prop.Condition < 40 || prop.Trouble {
		return "The laundry needs repairs before it can take a rush order"
	}
	if prop.Staff < 1 {
		return "The laundry has nobody to supervise the order"
	}
	if prop.Supply < RushOrderSupply {
		return "The laundry needs soap and coal before it can take the order"
	}
	return ""
}

func (w *World) StartRushOrder(id string) error {
	if reason := w.RushOrderReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	prop := w.Properties[id]
	prop.RushOrderDay = w.Minute/1440 + 1
	prop.Supply -= RushOrderSupply
	return nil
}

func (w *World) CompleteRushOrder() {
	w.Earn(RushOrderPay)
	if w.Player.Respect < DockName {
		w.Player.Respect++
	}
	w.Log("Hotel linen ready", fmt.Sprintf("The hotel's rush order is folded and packed. You receive $%d for helping at Bluebird Laundry. Another order can be taken tomorrow.", RushOrderPay), "work")
}
