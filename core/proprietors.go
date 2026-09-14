package core

import (
	"fmt"
	"sort"
)

// Personal deeds remain distinct from faction deeds even if an owner later
// joins a family. Only the person's own properties return to the broker.
func (w *World) releaseProprietor(id string) {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || prop.Owner != id {
			continue
		}
		prop.Owner = "independent"
		w.Log("A business returns to the market", fmt.Sprintf("%s is offered for purchase after the death of its individual proprietor. Its staff and condition remain as they were.", l.Name), "business")
	}
}

// Priced trading premises can be acquired by individual proprietors.
func personalBusiness(id string) bool {
	l, ok := PlaceByID(id)
	if !ok || ((l.Cost <= 0 || IsRentalHome(id) || l.Type == "home") && id != "room") {
		return false
	}
	_, trading := TradeOf(id)
	return trading
}

// At most one funded purchase per day; existing NPCs can enter the market later
// in a campaign, rather than ownership being assigned only at world creation.
func (w *World) ConsiderProprietors() {
	buyers := []*NPC{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead && n.Faction == "" && !IsOfficial(n.ID) && len(w.FamilyHoldings(n.ID)) == 0 {
			buyers = append(buyers, n)
		}
	}
	sort.Slice(buyers, func(i, j int) bool { return buyers[i].ID < buyers[j].ID })
	for offset := range buyers {
		n := buyers[(offset+w.Minute/1440)%len(buyers)]
		for _, l := range Locations {
			prop := w.Properties[l.ID]
			if !personalBusiness(l.ID) || prop == nil || prop.Owner != "independent" || prop.Income <= 0 || prop.Condition < 60 || prop.Trouble {
				continue
			}
			trade, ok := TradeOf(l.ID)
			if !ok {
				continue
			}
			price := AcquisitionCost(w, l.ID)
			float := 0
			if HasBankroll(l.ID) {
				float = BankrollLot
			}
			reserve := max(500, 7*(trade.Hands*trade.Wage+w.NPCLivingCost(n))+trade.Restock)
			if price <= 0 || w.HouseholdWealth(n) < price+reserve+float || !w.SpendHouseholdMoney(n, price+float) {
				continue
			}
			prop.Owner = n.ID
			prop.Bankroll += float
			prop.ProprietorDay = w.Minute/1440 + 1 // No full day's takings on purchase morning.
			w.Log("An independent proprietor", fmt.Sprintf("%s bought %s for $%d from their household funds. The premises remain individually owned.", n.Name, l.Name, price), "business")
			return
		}
	}
}

// Personal business accounts use their own capacity/custom, with real wages,
// stock and repair bills. Player licences and collection bonuses never apply.
func (w *World) ProprietorDay() {
	day := w.Minute/1440 + 1
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil {
			continue
		}
		n := w.NPC(prop.Owner)
		if n == nil {
			continue
		}
		if n.Dead {
			w.releaseProprietor(n.ID)
			continue
		}
		if prop.ProprietorDay == day {
			continue
		}
		trade, ok := TradeOf(l.ID)
		if !ok || (IsRentalHome(l.ID) && l.ID != "room") {
			continue
		}
		prop.ProprietorDay = day
		if HasBankroll(l.ID) {
			take := max(0, prop.Bankroll-BankrollFull)
			prop.Bankroll -= take
			n.Purse += take
			funding := min(max(0, BankrollLot-prop.Bankroll), max(0, w.HouseholdWealth(n)-500))
			if funding > 0 && w.SpendHouseholdMoney(n, funding) {
				prop.Bankroll += funding
			}
		}
		income := int(float64(prop.Income*24*prop.Condition) / 100 * w.Capacity(l.ID) * w.TradeMultiplier(l.ID))
		if l.ID != "room" {
			n.Purse += max(0, income)
		} // Lodging earns actual tenant payments.
		wages := prop.Staff * trade.Wage
		if w.SpendHouseholdMoney(n, wages) {
			prop.Unpaid = 0
		} else {
			prop.Unpaid++
			prop.Condition = max(0, prop.Condition-1)
		}
		prop.Supply = max(0, prop.Supply-trade.Drain)
		if prop.Supply < trade.Drain && w.SpendHouseholdMoney(n, trade.Restock) {
			prop.Supply += trade.RestockAmount
		}
		repair := min(2, 100-prop.Condition)
		if repair > 0 && w.SpendHouseholdMoney(n, repair*FamilyRepairCost) {
			prop.Condition += repair
		}
	}
}
