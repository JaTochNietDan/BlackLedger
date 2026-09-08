package core

import "fmt"

// Until now the moonshine on the market came from nowhere: it could be bought
// and sold but never made. A still turns a business the player already owns
// into a source, which is the difference between trading on somebody else's
// supply and having one of your own.
//
// It is the most profitable thing a business can do and the most dangerous. It
// produces stock that has to be moved, and stock draws attention every day it
// is held. A raid that finds a still costs far more than a raid that does not,
// and the people who rob premises would rather have crates than a till.

const (
	// StillCost is what setting one up takes.
	StillCost = 450
	// StillMinutes is how long that takes.
	StillMinutes = 180
	// StillHeat is the attention a working still draws each day, on top of
	// whatever the stock it produces draws.
	StillHeat = 2
	// StillSupplyDraw is the extra stock or fuel a still consumes each day.
	StillSupplyDraw = 4
	// StillHoard is the point at which there is nowhere left to put the output
	// and production stops until some of it moves.
	StillHoard = 60
)

// StillSites are the premises a still can be hidden in. A casino floor is not
// one of them: too many people, too much attention already.
func StillSite(id string) bool { return id == "laundry" || id == "garage" }

// StillReadiness explains why one cannot be set up, or returns "".
func (w *World) StillReadiness(id string) string {
	if !StillSite(id) || !w.Own(id) {
		return "There is nowhere here to hide one"
	}
	if w.Properties[id].Still {
		return "There is already one running here"
	}
	if w.Properties[id].Staff == 0 {
		return "Somebody has to work it"
	}
	if w.Player.Cash < StillCost {
		return "Not enough cash"
	}
	return ""
}

// BuildStill puts one in.
func (w *World) BuildStill(id string) error {
	if reason := w.StillReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(StillCost); err != nil {
		return err
	}
	w.Properties[id].Still = true
	place, _ := PlaceByID(id)
	w.Log("A still at "+place.Name, fmt.Sprintf("$%d of copper and pipe in the back of %s. It makes its own stock now, and stock has to be moved.", StillCost, place.Name), "business")
	return nil
}

// DismantleReadiness explains why one cannot be taken out, or returns "".
func (w *World) DismantleReadiness(id string) string {
	if !w.Own(id) || w.Properties[id] == nil || !w.Properties[id].Still {
		return "There is nothing here to take out"
	}
	return ""
}

// Dismantle removes a still, which is what somebody does when the police are
// getting close and the profit is no longer worth the risk.
func (w *World) Dismantle(id string) error {
	if reason := w.DismantleReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	w.Properties[id].Still = false
	w.Player.Heat = max(0, w.Player.Heat-4)
	place, _ := PlaceByID(id)
	w.Log("The still comes out of "+place.Name, "Copper and pipe out through the back door. There is less to find here now.", "business")
	return nil
}

// StillOutput is what a still makes in a day, which depends on the business
// running well enough to hide it.
func (w *World) StillOutput(id string) int {
	prop := w.Properties[id]
	if prop == nil || !prop.Still {
		return 0
	}
	return max(1, int(float64(6)*w.Capacity(id)))
}

// StillDay runs every still the player owns. Output arrives as stock, which is
// where the risk starts rather than ends.
func (w *World) StillDay() {
	heat := 0
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || !prop.Still || !w.Own(l.ID) {
			continue
		}
		heat += StillHeat
		if prop.Supply <= 0 || prop.Staff == 0 {
			w.Log("The still is cold at "+l.Name, "Nothing to run it with and nobody to work it. It makes nothing today.", "business")
			continue
		}
		if w.Holding("moonshine") >= StillHoard {
			w.Log("Nowhere to put it at "+l.Name, fmt.Sprintf("There are already %d crates on hand. The still stands idle until some of it moves.", w.Holding("moonshine")), "business")
			continue
		}
		prop.Supply = max(0, prop.Supply-StillSupplyDraw)
		made := w.StillOutput(l.ID)
		if w.Player.Stock == nil {
			w.Player.Stock = map[string]int{}
		}
		w.Player.Stock["moonshine"] += made
		w.Log("Run off at "+l.Name, fmt.Sprintf("%d crates out of the back of %s. You are holding %d, and holding is what gets noticed.", made, l.Name, w.Holding("moonshine")), "business")
	}
	if heat > 0 {
		w.Player.Heat = min(100, w.Player.Heat+heat)
	}
}

// StillFound is what a search costs when there is one to find. Separate from
// seizing stock, because the still is the thing that makes them come back.
func (w *World) StillFound() (string, bool) {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop != nil && prop.Still && w.Own(l.ID) {
			prop.Still = false
			prop.Condition = max(0, prop.Condition-25)
			w.Player.Heat = min(100, w.Player.Heat+10)
			return l.Name, true
		}
	}
	return "", false
}
