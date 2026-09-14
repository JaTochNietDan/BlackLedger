package core

import "fmt"

// Opportunity is a suggestion based on public facts, never a disclosure of hidden plots.
type Opportunity struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Target string `json:"target"`
}

// Inspect the action as it would appear at the destination, without moving the
// real player. Actions is a read-only projection; travel can change availability
// before arrival, so guidance never executes the suggested work automatically.
func (w *World) opportunityAction(place, id string) (Action, bool) {
	view := *w
	view.Player.Location = place
	for _, a := range view.Actions(place) {
		if a.ID == id {
			return a, true
		}
	}
	return Action{}, false
}

func (w *World) actionOpportunity(title, place, id string) *Opportunity {
	a, ok := w.opportunityAction(place, id)
	if !ok || a.Disabled {
		return nil
	}
	terms := fmt.Sprintf("%d minutes", a.Minutes+a.Away)
	if fee := a.Cost + a.Asks; fee > 0 {
		terms += fmt.Sprintf(" · $%d", fee)
	}
	return &Opportunity{title, fmt.Sprintf("%s · %s. %s", a.Label, terms, a.Detail), place}
}

func (w *World) earningOpportunity(title, why string) *Opportunity {
	for _, job := range [][2]string{{"bar", "courier"}, {"docks", "dockwork"}} {
		if next := w.actionOpportunity(title, job[0], job[1]); next != nil {
			if why != "" {
				next.Detail = why + " " + next.Detail
			}
			return next
		}
	}
	return nil
}

func (w *World) NextOpportunity() *Opportunity {
	p := &w.Player
	if !p.Alive || w.Event != nil || w.Held() {
		return nil
	}
	for _, l := range Locations {
		if w.Own(l.ID) && w.Properties[l.ID].Income > 0 && w.Properties[l.ID].Condition < 70 {
			if next := w.actionOpportunity("Restore your income", l.ID, "repair"); next != nil {
				return next
			}
			if next := w.earningOpportunity("Put repair money aside", "Your business needs repairs. Paid work can cover the bill."); next != nil {
				return next
			}
		}
	}
	if p.JobCount == 0 {
		if next := w.earningOpportunity("Make your first connection", "A small job starts your reputation."); next != nil {
			return next
		}
	}
	if p.Respect < PremisesRespect {
		if next := w.earningOpportunity("Become a known face", fmt.Sprintf("%d of %d respect toward your first business.", p.Respect, PremisesRespect)); next != nil {
			return next
		}
	}
	owns := false
	for _, l := range Locations {
		if w.Own(l.ID) && w.Properties[l.ID].Income > 0 {
			owns = true
		}
	}
	if !owns {
		if next := w.actionOpportunity("Build a steady income", "laundry", "acquire"); next != nil {
			return next
		}
		// The first address can change hands. Offer another real, affordable
		// business rather than pointing forever at an unavailable laundry.
		for _, l := range Locations {
			if next := w.actionOpportunity("Build a steady income", l.ID, "acquire"); next != nil {
				return next
			}
		}
		why := "Build your capital before taking on a business and its daily bills."
		if a, ok := w.opportunityAction("laundry", "acquire"); ok && a.Reason == "Not enough cash" {
			why = fmt.Sprintf("Bluebird Laundry needs $%d; you have $%d. Keep money aside for its daily bills too.", a.Cost+a.Asks, p.Cash)
		}
		if next := w.earningOpportunity("Save for your first business", why); next != nil {
			return next
		}
	}
	if len(p.Crew) == 0 && len(w.OwnPeople()) == 0 {
		if next := w.actionOpportunity("Bring someone into the fold", "bar", "recruit"); next != nil {
			return next
		}
		if a, ok := w.opportunityAction("bar", "recruit"); ok && a.Reason == "Not enough cash" {
			if next := w.earningOpportunity("Set aside the hiring money", fmt.Sprintf("The introduction costs $%d, before daily wages.", a.Cost+a.Asks)); next != nil {
				return next
			}
		}
	}
	if p.Contacts < 2 {
		if next := w.actionOpportunity("Know who is asking about you", "bar", "contact"); next != nil {
			return next
		}
	}
	if w.District == 0 {
		if next := w.actionOpportunity("Reach beyond Old Harbor", "apartment", "expand"); next != nil {
			return next
		}
		if next := w.earningOpportunity("Build your connections fund", "Paid work builds the money and standing to reach the next district."); next != nil {
			return next
		}
	}
	if p.Home == "room" {
		if next := w.actionOpportunity("Find a safer address", "apartment", "move_home"); next != nil {
			return next
		}
	}
	for _, place := range []string{"garage", "casino"} {
		if next := w.actionOpportunity("Expand your organization", place, "acquire"); next != nil {
			return next
		}
	}
	return &Opportunity{"Keep the organization solvent", "Review income, upkeep and family relationships before your next expansion.", p.Home}
}
