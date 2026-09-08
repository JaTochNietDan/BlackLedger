package core

import "fmt"

// Owning premises was passive: the money arrived whether you thought about it
// or not. How a business is run is now a standing decision with a cost on both
// sides. Skimming hard pays for the crew and draws the wrong attention; running
// clean survives scrutiny and barely pays; the middle is unremarkable, which is
// its own kind of safety.
//
// Nothing here is a one-off gamble. A mode persists until it is changed, and
// its consequences accrue every day the business is open.

type OperatingMode struct {
	ID     string
	Label  string
	Detail string
	// Take scales income. Heat is police attention added per day. Wear is the
	// condition lost per day. Notice is how much more likely a family is to
	// take an interest in the earnings.
	Take   float64
	Heat   int
	Wear   int
	Notice int
}

var operatingModes = []OperatingMode{
	{ID: "clean", Label: "Run it clean",
		Detail: "Books that survive an audit. Earns well short of what the premises could, and nobody looks twice.",
		Take:   .7, Heat: 0, Wear: 0, Notice: -1},
	{ID: "standard", Label: "Run it as usual",
		Detail: "An ordinary business doing ordinary trade. No attention beyond the usual.",
		Take:   1, Heat: 0, Wear: 0, Notice: 0},
	{ID: "hard", Label: "Skim what it will bear",
		Detail: "Half again the money, police attention every day, premises that need regular repair, and earnings a family may decide are worth a share.",
		Take:   1.5, Heat: 1, Wear: 1, Notice: 2},
}

func operatingMode(id string) OperatingMode {
	for _, m := range operatingModes {
		if m.ID == id {
			return m
		}
	}
	return operatingModes[1] // saves written before modes existed run as usual
}

// Mode reports how a property is being run.
func (w *World) Mode(id string) OperatingMode {
	if prop := w.Properties[id]; prop != nil {
		return operatingMode(prop.Mode)
	}
	return operatingModes[1]
}

// SetMode changes how an owned business is run. It commits immediately and
// costs no time: it is a decision about the future, not an errand.
func (w *World) SetMode(id, mode string) error {
	if !w.Own(id) {
		return fmt.Errorf("you do not own this business")
	}
	prop := w.Properties[id]
	if prop.Income <= 0 {
		return fmt.Errorf("this property does not trade")
	}
	found := false
	for _, m := range operatingModes {
		if m.ID == mode {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("unknown way of running a business")
	}
	if prop.Mode == mode || (prop.Mode == "" && mode == "standard") {
		return fmt.Errorf("it is already run that way")
	}
	prop.Mode = mode
	place, _ := PlaceByID(id)
	w.Log("A change at "+place.Name, operatingMode(mode).Label+". "+operatingMode(mode).Detail, "business")
	return nil
}

// BusinessDay applies the standing consequences of how each owned business is
// run. Called once per game day, alongside the player's bills.
func (w *World) BusinessDay() {
	heat, worn := 0, []string{}
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || !w.Own(l.ID) || prop.Income <= 0 {
			continue
		}
		mode := operatingMode(prop.Mode)
		heat += mode.Heat
		if mode.Wear > 0 && prop.Condition > 0 {
			prop.Condition = max(0, prop.Condition-mode.Wear)
			if prop.Condition < 50 {
				worn = append(worn, l.Name)
			}
		}
	}
	if heat > 0 {
		w.Player.Heat = min(100, w.Player.Heat+heat)
	}
	if len(worn) > 0 {
		w.Log("The premises are showing it", fmt.Sprintf("%s need work. A business in poor condition earns less of what it could.", joinNames(worn)), "business")
	}
}

func joinNames(names []string) string {
	switch len(names) {
	case 0:
		return ""
	case 1:
		return names[0]
	case 2:
		return names[0] + " and " + names[1]
	}
	out := ""
	for i, n := range names[:len(names)-1] {
		if i > 0 {
			out += ", "
		}
		out += n
	}
	return out + " and " + names[len(names)-1]
}

// SkimNotice is how much extra interest the player's businesses attract from
// families deciding whether to demand a share. Running everything clean is a
// reason to be left alone.
func (w *World) SkimNotice() int {
	notice := 0
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && w.Own(l.ID) && prop.Income > 0 {
			notice += operatingMode(prop.Mode).Notice
		}
	}
	return notice
}
