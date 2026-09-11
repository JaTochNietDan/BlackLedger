package core

import "fmt"

// A residence was a rent line and a number of guards. Everything else a person
// might do with the place they live in — put a door on it that holds, a
// telephone in the hall, a safe behind the panelling, a dry cellar under it —
// did not exist.
//
// Comforts belong to the building rather than to the person. Moving out leaves
// them behind, which turns a housing tier from a strict upgrade into a decision
// about what you would be giving up. Each one costs money to fit, costs money
// every day, and does exactly one thing the rest of the game already cares
// about.

// Comfort is something fitted to a residence.
type Comfort struct {
	ID string
	// Label is what it is, in words.
	Label string
	// Detail is what it does, stated plainly enough to decide on.
	Detail string
	// Cost is what fitting it takes, and Upkeep what it costs a day.
	Cost, Upkeep int
	// Minutes is how long the work takes.
	Minutes int
}

var comforts = []Comfort{
	{ID: "door", Label: "A door that holds", Minutes: 90, Cost: 420, Upkeep: 2,
		Detail: "Steel behind the wood and bars on the ground-floor windows. Counts as another guard when somebody comes for you at home."},
	{ID: "telephone", Label: "A telephone in the hall", Minutes: 60, Cost: 300, Upkeep: 5,
		Detail: "People can reach you without walking here first. Counts as another contact for warnings and for hearing whose name is on something."},
	{ID: "safe", Label: "A safe behind the panelling", Minutes: 120, Cost: 560, Upkeep: 1,
		Detail: "Holds up to $900 that a fine cannot touch and a thief does not find."},
	{ID: "cellar", Label: "A dry cellar", Minutes: 150, Cost: 480, Upkeep: 2,
		Detail: "Hides 25 units of stock from the attention that holding it draws. A warrant served on your residence finds all of it."},
}

const (
	// SafeShelter is what a safe keeps out of reach.
	SafeShelter = 900
	// CellarHold is what a cellar hides from attention.
	CellarHold = 25
)

// ComfortByID is one of them, and whether it exists.
func ComfortByID(id string) (Comfort, bool) {
	for _, c := range comforts {
		if c.ID == id {
			return c, true
		}
	}
	return Comfort{}, false
}

// Comforts is what is fitted somewhere, in the order they are offered so the
// interface is stable.
func (w *World) Comforts(id string) []string {
	prop := w.Properties[id]
	if prop == nil {
		return nil
	}
	out := []string{}
	for _, c := range comforts {
		for _, fitted := range prop.Comforts {
			if fitted == c.ID {
				out = append(out, c.ID)
			}
		}
	}
	return out
}

// Fitted reports whether the place the player lives has something.
func (w *World) Fitted(comfort string) bool {
	for _, id := range w.Comforts(w.Player.Home) {
		if id == comfort {
			return true
		}
	}
	return false
}

// ComfortUpkeep is what everything fitted to the player's residence costs a
// day. Comforts somewhere they no longer live cost them nothing, because they
// belong to whoever is there now.
func (w *World) ComfortUpkeep() int {
	total := 0
	for _, id := range w.Comforts(w.Player.Home) {
		c, _ := ComfortByID(id)
		total += c.Upkeep
	}
	return total
}

// Reach is how well the player hears things: the people they know, a telephone
// if there is one in the hall, and a public house of their own, which is the
// one room in this city where everything gets said out loud in front of
// whoever owns it.
func (w *World) Reach() int {
	reach := w.Player.Contacts + w.EarsInTheRoom()
	if w.Fitted("telephone") {
		reach++
	}
	return reach
}

// Sheltered is the cash a fine cannot reach and a thief does not find.
func (w *World) Sheltered() int {
	if !w.Fitted("safe") {
		return 0
	}
	return min(SafeShelter, w.Player.Cash)
}

// Reachable is the money that can actually be taken from the player.
func (w *World) Reachable() int { return max(0, w.Player.Cash-w.Sheltered()) }

// FitReadiness explains why something cannot be fitted here, or returns "".
func (w *World) FitReadiness(place, comfort string) string {
	c, ok := ComfortByID(comfort)
	if !ok {
		return "There is no such thing"
	}
	l, known := PlaceByID(place)
	if !known || l.Type != "home" {
		return "Nobody lives here"
	}
	if place != w.Player.Home {
		return "You would be fitting out somebody else's home"
	}
	if w.Fitted(comfort) {
		return "There is already one here"
	}
	if w.Player.Cash < c.Cost {
		return "Not enough cash"
	}
	return ""
}

// Fit puts something into the place the player lives. It stays with the
// building.
func (w *World) Fit(place, comfort string) error {
	if reason := w.FitReadiness(place, comfort); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	c, _ := ComfortByID(comfort)
	if err := w.Pay(c.Cost); err != nil {
		return err
	}
	prop := w.Properties[place]
	prop.Comforts = append(prop.Comforts, comfort)
	l, _ := PlaceByID(place)
	w.Log(c.Label+" at "+l.Name, fmt.Sprintf("$%d and half a day's work. $%d a day from here, and it stays with the building if you ever move out.", c.Cost, c.Upkeep), "personal")
	return nil
}

// StripComforts is what a search does to a cellar: it finds it, and everything
// in it. Reported separately from the ordinary seizure because the point is
// that the hiding place was the thing that failed.
func (w *World) CellarFound() int {
	if !w.Fitted("cellar") {
		return 0
	}
	// Only what the cellar was holding. Anything under the floor of a car is
	// somewhere else entirely.
	held := min(w.Carrying(), CellarHold)
	if held == 0 {
		return 0
	}
	remaining := held
	for _, g := range w.Goods {
		if w.Player.Stock == nil {
			break
		}
		take := min(w.Player.Stock[g.ID], remaining)
		w.Player.Stock[g.ID] -= take
		remaining -= take
	}
	l, _ := PlaceByID(w.Player.Home)
	w.Log("They went through the cellar at "+l.Name, fmt.Sprintf("A dry cellar is the first place anybody looks twice at. %d units gone.", held), "danger")
	return held
}

// ResidenceDescription is what the place the player lives has in it, for the
// interface.
func (w *World) ResidenceDescription() map[string]any {
	fitted := []map[string]any{}
	for _, id := range w.Comforts(w.Player.Home) {
		c, _ := ComfortByID(id)
		fitted = append(fitted, map[string]any{"id": c.ID, "label": c.Label, "upkeep": c.Upkeep})
	}
	return map[string]any{
		"home":      w.Player.Home,
		"comforts":  fitted,
		"upkeep":    w.ComfortUpkeep(),
		"sheltered": w.Sheltered(),
		"reach":     w.Reach(),
	}
}

// doorProtection is what steel and bars are worth when nobody warned you.
// Substantial, because being unable to get through the door is most of what
// stops this kind of thing, and well short of safety.
func (w *World) doorProtection() float64 {
	if w.Fitted("door") {
		return .18
	}
	return 0
}
