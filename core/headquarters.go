package core

import "fmt"

// Headquarters deeds belong to the organization, not its current leader.
func (w *World) headquartersSite(owner, id string) bool {
	p := w.Properties[id]
	_, trade := TradeOf(id)
	return p != nil && p.Owner == owner && trade && p.Income > 0
}

func (w *World) Headquarters(owner string) string {
	f := w.faction(owner)
	if f != nil && w.headquartersSite(owner, f.Headquarters) {
		return f.Headquarters
	}
	return ""
}

// Existing families gain a base on upgrade. Player families whose recorded
// base was lost must choose another; NPC families relocate to a remaining deed.
func (w *World) SettleHeadquarters() {
	for i := range w.Factions {
		f := &w.Factions[i]
		if w.Headquarters(f.ID) != "" {
			continue
		}
		if f.ID == w.PlayerOrganizationID() && f.Headquarters != "" {
			continue
		}
		best := ""
		for _, p := range Locations {
			if w.headquartersSite(f.ID, p.ID) && (best == "" || w.Properties[p.ID].Income > w.Properties[best].Income) {
				best = p.ID
			}
		}
		if best != "" {
			f.Headquarters = best
		}
	}
}

func (w *World) HeadquartersReadiness(at string, forming bool) string {
	if !w.Player.Alive {
		return "This life has ended"
	}
	if !w.headquartersSite(w.PlayerOrganizationID(), at) {
		return "Choose a business your organization owns"
	}
	if w.Player.Location != at {
		return "Go to that business to establish the headquarters"
	}
	if w.Held() {
		return "You cannot establish a headquarters while in custody"
	}
	p := w.Properties[at]
	if p.Condition < 60 || p.Trouble {
		return "Put the premises in working order first"
	}
	if forming {
		if w.Incorporated() {
			return "Your family already exists"
		}
		if w.Player.Serves != "" {
			return "Leave your current family before forming your own"
		}
		if w.Player.Respect < OrganizationStanding {
			return fmt.Sprintf("Earn %d respect first", OrganizationStanding)
		}
	} else if !w.Incorporated() {
		return "Form your family first"
	}
	return ""
}

func (w *World) EstablishHeadquarters(at string, forming bool) error {
	if reason := w.HeadquartersReadiness(at, forming); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if forming {
		w.Incorporate()
	}
	f := w.PlayerOrganization()
	if f == nil {
		return fmt.Errorf("Your family could not be formed")
	}
	f.Headquarters = at
	place, _ := PlaceByID(at)
	w.Log("Headquarters established", fmt.Sprintf("%s operates from %s. The business remains a working concern; losing its deed or having it shut down interrupts headquarters operations.", f.Name, place.Name), "politics")
	return nil
}
