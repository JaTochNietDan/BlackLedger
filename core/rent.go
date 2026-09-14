package core

// RentAccount belongs to the premises, including any outstanding tenant debt.
// A transfer of the premises transfers these receivables; a visit does not.
type RentAccount struct {
	Day       int `json:"day"`
	Due       int `json:"due"`
	Paid      int `json:"paid"`
	Arrears   int `json:"arrears"`
	Collected int `json:"collected"`
}

// The former $12 allowance included $8 lodging and $4 other necessities.
// Earnings now cover the actual accommodation instead of charging both costs.
const npcOtherLivingCost = 4

func (w *World) NPCRent(n *NPC) int {
	if n == nil || n.Dead || w.Properties[n.Home] == nil {
		return 0
	}
	switch n.Home {
	case "room":
		return 15
	case "apartment":
		if n.Accommodation == "Shared flat" {
			return 15
		}
		return 35
	}
	return 0
}

func (w *World) NPCLivingCost(n *NPC) int {
	if rent := w.NPCRent(n); rent > 0 {
		return npcOtherLivingCost + rent
	}
	return LivingCost
}

func (w *World) collectRent(n *NPC) {
	rent := w.NPCRent(n)
	if rent == 0 {
		return
	}
	p := w.Properties[n.Home]
	if p.Rents == nil {
		p.Rents = map[string]*RentAccount{}
	}
	account := p.Rents[n.ID]
	if account == nil {
		account = &RentAccount{}
		p.Rents[n.ID] = account
	}
	day := w.Minute/1440 + 1
	if account.Day == day {
		return
	}
	account.Day, account.Due = day, rent
	account.Arrears += rent
	account.Paid = min(max(0, n.Purse), account.Arrears)
	n.Purse -= account.Paid
	account.Arrears -= account.Paid
	account.Collected += account.Paid
	if w.Own(n.Home) {
		w.Earn(account.Paid)
	} else if f := w.faction(p.Owner); f != nil {
		f.Cash += account.Paid
	}
}

// RentalDaily is the contracted rate of living residents, regardless of where
// they happen to be standing. Cash is collected separately at daily settlement.
func (w *World) RentalDaily(id string) int {
	total := 0
	for _, n := range w.Residents(id) {
		total += w.NPCRent(n)
	}
	return total
}

func (w *World) RentRegister(id string) map[string]any {
	if id != "room" && id != "apartment" {
		return nil
	}
	tenants := []map[string]any{}
	for _, n := range w.Residents(id) {
		row := map[string]any{"id": n.ID, "name": n.Name, "daily": w.NPCRent(n), "accommodation": n.Accommodation}
		if w.Own(id) {
			if a := w.Properties[id].Rents[n.ID]; a != nil {
				row["account"] = *a
			}
		}
		tenants = append(tenants, row)
	}
	return map[string]any{"capacity": ResidentialCapacity[id], "occupied": len(tenants), "daily": w.RentalDaily(id), "tenants": tenants}
}

// HomeCost keeps the player's room rental separate from the deed. Cypress's
// existing upkeep remains a household cost even when its resident owns it.
func (w *World) HomeCost(id string) int {
	if (id == "room" || id == "apartment") && w.Own(id) {
		return 0
	}
	return HomeRent(id)
}
