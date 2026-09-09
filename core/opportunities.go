package core

// Opportunity is a suggestion based on public facts, never a disclosure of hidden plots.
type Opportunity struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Target string `json:"target"`
}

func (w *World) NextOpportunity() *Opportunity {
	p := &w.Player
	if !p.Alive || w.Event != nil {
		return nil
	}
	for _, l := range Locations {
		if w.Own(l.ID) && w.Properties[l.ID].Income > 0 && w.Properties[l.ID].Condition < 70 {
			return &Opportunity{"Restore your income", "Repair the damaged business to recover its earning capacity.", l.ID}
		}
	}
	if p.JobCount == 0 {
		return &Opportunity{"Make your first connection", "Mara has paid work at Saint Agnes. A small job starts your reputation.", "bar"}
	}
	if p.Respect < PremisesRespect {
		return &Opportunity{"Become a known face", "Earn 6 respect to recruit an associate or establish your first business.", "bar"}
	}
	owns := false
	for _, l := range Locations {
		if w.Own(l.ID) && w.Properties[l.ID].Income > 0 {
			owns = true
		}
	}
	if !owns && w.CanAcquire("laundry") {
		return &Opportunity{"Build a steady income", "Bluebird Laundry needs 6 respect and capital; buying out a former organization costs more. Its income grows with game time.", "laundry"}
	}
	if len(p.Crew) == 0 {
		// Whoever drives, the same person the recruit button offers. Nobody
		// holds a job in this city forever.
		hand := "A driver"
		if n := w.Holder("driver"); n != nil {
			hand = n.Name
		}
		return &Opportunity{"Bring someone into the fold", hand + " costs $90 to recruit and $12 a day. They can collect money or protect businesses.", "bar"}
	}
	if p.Contacts < 2 {
		return &Opportunity{"Know who is asking about you", "Develop your information network through Mara. Good contacts can warn of personal danger.", "bar"}
	}
	if w.District == 0 {
		return &Opportunity{"Reach beyond Old Harbor", "With 10 respect and $100, establish contacts in Ashbury.", "apartment"}
	}
	if p.Home == "room" {
		return &Opportunity{"Find a safer address", "Ashbury Court costs $180 to rent and $35 a day, with a base level of residential protection.", "apartment"}
	}
	if !w.Own("garage") && w.CanAcquire("garage") {
		return &Opportunity{"Expand your organization", "Russo Motor Works needs 10 respect and capital; former operators require a buyout. Taking it will affect your standing with Russo.", "garage"}
	}
	if !w.Own("casino") && w.CanAcquire("casino") {
		return &Opportunity{"Open your own casino", "The Blue Hour requires 20 respect and capital; former operators require a buyout. Build the capital and connections to run it.", "casino"}
	}
	return &Opportunity{"Keep the organization solvent", "Review income, upkeep and family relationships before your next expansion.", p.Home}
}
