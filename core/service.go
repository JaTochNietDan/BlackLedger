package core

import "fmt"

// Every path this game offers starts with a man buying premises. There is no
// version of it in which somebody arrives with ninety dollars, goes to work for
// one of the families, and comes up through it — which is the oldest story this
// city has.
//
// Answering to somebody is the other career. It pays a wage rather than costing
// one, it makes their quarrels yours, it protects you from them and from
// nobody else, and it rises: work done for them is standing inside them, and
// standing inside them is a share of what they take.

const (
	// ServiceGoodwill is the standing it takes before anybody has you.
	ServiceGoodwill = 20
	// ServiceRespect is the name it takes.
	ServiceRespect = 8
	// ServiceMinutes is the conversation.
	ServiceMinutes = 60
	// SoldierWage is what answering to somebody pays a day at the bottom.
	SoldierWage = 22
	// LieutenantShare is the share of a lieutenant's organization's daily take
	// that reaches them, as a percentage.
	LieutenantShare = 4
	// PromotionWork is how many pieces of work for them it takes to come up.
	PromotionWork = 3
	// LeavingCost is what walking out costs in standing with them.
	LeavingCost = 45
)

// Serving reports whose organization the player answers to, or "".
func (w *World) Serving() string { return w.Player.Serves }

// ServiceRank is the player's standing inside whoever they answer to.
func (w *World) ServiceRank() int {
	if w.Player.Serves == "" {
		return 0
	}
	if w.Player.Service >= PromotionWork*2 {
		return RankLieutenant
	}
	if w.Player.Service >= PromotionWork {
		return RankSoldier
	}
	return RankAssociate
}

// ServiceTitle is what they call the player inside it.
func (w *World) ServiceTitle() string {
	switch w.ServiceRank() {
	case RankLieutenant:
		return "Lieutenant"
	case RankSoldier:
		return "Soldier"
	}
	return "Associate"
}

// ServicePay is what answering to somebody is worth a day: a wage at the
// bottom, and a share of what they take once you are worth something to them.
func (w *World) ServicePay() int {
	f := w.faction(w.Player.Serves)
	if f == nil {
		return 0
	}
	pay := SoldierWage
	if w.ServiceRank() >= RankLieutenant {
		income := 0
		for _, id := range w.FamilyHoldings(f.ID) {
			prop := w.Properties[id]
			income += prop.Income * prop.Condition / 100 * 24
		}
		pay += income * LieutenantShare / 100
	}
	return pay
}

// ServeReadiness explains why the player cannot go to work for somebody, or
// returns "".
func (w *World) ServeReadiness(id string) string {
	f := w.faction(id)
	if f == nil || id == w.PlayerOrganizationID() {
		return "There is nobody of that description"
	}
	if w.Player.Serves != "" {
		return "You already answer to somebody"
	}
	if w.Incorporated() {
		return "You have your own thing. Nobody comes up through somebody else's while they have one"
	}
	if w.Player.Location != w.homeOf(id) {
		place, _ := PlaceByID(w.homeOf(id))
		return "That conversation happens at " + place.Name
	}
	if f.Goodwill < ServiceGoodwill {
		return fmt.Sprintf("They think of you at %+d. It takes %+d before anybody takes you on", f.Goodwill, ServiceGoodwill)
	}
	if w.Player.Respect < ServiceRespect {
		return fmt.Sprintf("Earn %d respect first", ServiceRespect)
	}
	return ""
}

// Serve puts the player to work for somebody.
func (w *World) Serve(id string) error {
	if reason := w.ServeReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	w.Player.Serves, w.Player.Service = id, 0
	f := w.faction(id)
	// Their quarrels are yours from this afternoon.
	for i := range w.Factions {
		other := &w.Factions[i]
		if other.ID == id {
			continue
		}
		if c := w.Conflict(id, other.ID); c != nil && c.State != "cold" {
			other.Goodwill = max(-100, other.Goodwill-15)
		}
	}
	w.Log("You answer to "+f.Name+" now", fmt.Sprintf("$%d a day and a place to stand. Their quarrels are yours from this afternoon, and coming up means doing what they ask.", SoldierWage), "politics")
	return nil
}

// LeaveService is walking out. Nobody takes it well.
func (w *World) LeaveService() error {
	if w.Player.Serves == "" {
		return fmt.Errorf("you answer to nobody")
	}
	f := w.faction(w.Player.Serves)
	w.Player.Serves, w.Player.Service = "", 0
	if f != nil {
		f.Goodwill = max(-100, f.Goodwill-LeavingCost)
		w.Log("You answer to nobody again", fmt.Sprintf("%s %s keep people who leave in their good books. Their standing with you is %+d.", f.Name, Agree(f.Name, "does not", "do not"), f.Goodwill), "politics")
		w.RetaliationFrom(f.ID)
	}
	return nil
}

// ServeAgainst credits work done against somebody whoever the player answers to
// is at odds with. Harm done to their enemies is work done for them, which is
// the only way most people in this business ever come up: an arrangement with a
// beneficiary is not something everybody is ever offered.
func (w *World) ServeAgainst(id string) {
	if w.Player.Serves == "" || id == "" || id == w.Player.Serves {
		return
	}
	c := w.Conflict(w.Player.Serves, id)
	if c == nil || c.State == "cold" {
		return
	}
	w.ServeWork(w.Player.Serves)
}

// ServeWork records a piece of work done for whoever the player answers to,
// which is the only way anybody comes up.
func (w *World) ServeWork(beneficiary string) {
	if w.Player.Serves == "" || beneficiary != w.Player.Serves {
		return
	}
	before := w.ServiceRank()
	w.Player.Service++
	if after := w.ServiceRank(); after != before {
		f := w.faction(w.Player.Serves)
		w.Log("They have moved you up", fmt.Sprintf("%s calls you a %s now. It pays $%d a day.", f.Name, lowerFirst(w.ServiceTitle()), w.ServicePay()), "politics")
	}
}

// ServiceDay pays the wage and holds the quarrel with whoever the player
// answers to at nothing, because nobody moves on their own people.
func (w *World) ServiceDay() {
	if w.Player.Serves == "" {
		return
	}
	f := w.faction(w.Player.Serves)
	if f == nil {
		// Whoever they answered to no longer exists.
		w.Player.Serves, w.Player.Service = "", 0
		return
	}
	pay := w.ServicePay()
	if f.Cash < pay {
		w.Log("Nothing came this week", f.Name+" did not have it. Nobody stays where the money stops.", "danger")
		return
	}
	f.Cash -= pay
	w.Earn(pay)
	f.Goodwill = min(100, f.Goodwill+1)
}

// ServiceDescription is who the player answers to, for the interface.
func (w *World) ServiceDescription() map[string]any {
	if w.Player.Serves == "" {
		return map[string]any{"serving": false}
	}
	f := w.faction(w.Player.Serves)
	name := "somebody"
	if f != nil {
		name = f.Name
	}
	return map[string]any{
		"serving": true, "name": name, "title": w.ServiceTitle(),
		"work": w.Player.Service, "pay": w.ServicePay(),
		"next": max(0, PromotionWork*2-w.Player.Service),
	}
}
