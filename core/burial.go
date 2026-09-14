package core

import "fmt"

const (
	// BurialTrade is the additional custom from a paid funeral.
	BurialTrade  = 2
	BurialPauper = 0
)

// TheParlour is where this city's funerals go. A parlour of the player's own
// first, then the first one on the list, which keeps this off map iteration
// order the same way the garage does.
func (w *World) TheParlour() string {
	first := ""
	for _, l := range Locations {
		if l.Kind != "undertaker" || w.Properties[l.ID] == nil {
			continue
		}
		if w.Own(l.ID) {
			return l.ID
		}
		if first == "" {
			first = l.ID
		}
	}
	return first
}

// Bury settles an ordinary funeral once. Player crew retain their explicit
// arrangement window; their money is not spent before that decision.
func (w *World) Bury(person *NPC) {
	if person == nil || !person.Dead || person.Buried || person.Faction == w.PlayerOrganizationID() {
		return
	}
	id := w.TheParlour()
	if id == "" {
		return
	}
	person.Buried = true
	estate := w.HouseholdWealth(person)
	family := w.faction(person.Faction)
	available := estate
	if family != nil {
		available += max(0, family.Cash)
	}
	// An unfunded burial still takes place, without invented parish income.
	if available < BurialPurse {
		return
	}
	paid := min(FuneralCost, available)
	fromEstate := min(estate, paid)
	if !w.SpendHouseholdMoney(person, fromEstate) {
		return
	}
	if family != nil {
		family.Cash -= paid - fromEstate
	}
	w.funeralProceeds(id, person.Name, paid)
	w.deathServicePayments(person, (paid*FuneralOwn+FuneralCost-1)/FuneralCost)
	w.ShiftCustom(id, "nobody in this district dying", BurialTrade)
}

// The plot, transport and notices are outside costs. Reduced means buy a
// smaller service at the same cost ratio; only the funded margin reaches the
// proprietor. An independent unowned business keeps its takings off-ledger.
func (w *World) funeralProceeds(id, name string, paid int) {
	if paid <= 0 {
		return
	}
	cost := (paid*FuneralOwn + FuneralCost - 1) / FuneralCost
	margin := max(0, paid-cost)
	if w.Own(id) {
		w.Earn(margin)
		place, _ := PlaceByID(id)
		w.Log("Funeral accounts at "+place.Name,
			fmt.Sprintf("%s's funeral brought in $%d; $%d covered the plot, transport and notices. The remaining $%d went into your accounts.", name, paid, cost, margin), "business")
	} else {
		w.changeBusinessFunds(id, margin)
	}
}

// BurialPurse is the minimum available estate/family money for a paid service.
const BurialPurse = 25
