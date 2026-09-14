package core

import (
	"fmt"
	"sort"
	"strings"
)

// Apartment deeds are separate from the freehold and its common parts. The
// building address remains the authoritative destination for each resident.
type ApartmentDeed struct {
	ID       string `json:"id"`
	Building string `json:"building"`
	Number   int    `json:"number"`
	Owner    string `json:"owner"`
	Resident string `json:"resident,omitempty"`
}

func (w *World) playerDeedID() string { return fmt.Sprintf("player:%d", w.Life) }
func (w *World) apartmentForResident(id string) *ApartmentDeed {
	for i := range w.Apartments {
		if w.Apartments[i].Resident == id {
			return &w.Apartments[i]
		}
	}
	return nil
}
func (w *World) apartment(id string) *ApartmentDeed {
	for i := range w.Apartments {
		if w.Apartments[i].ID == id {
			return &w.Apartments[i]
		}
	}
	return nil
}

// SettleApartments is called only after committed housing assignment, never
// from PlanHomeMove. Existing occupants retain their numbered flat.
func (w *World) SettleApartments() {
	existing := map[string]bool{}
	for _, u := range w.Apartments {
		existing[u.ID] = true
	}
	for _, building := range []string{"apartment", "mercercourt"} {
		for number := 1; number <= ResidentialCapacity[building]; number++ {
			id := fmt.Sprintf("%s-%02d", building, number)
			if !existing[id] {
				w.Apartments = append(w.Apartments, ApartmentDeed{ID: id, Building: building, Number: number, Owner: "independent"})
			}
		}
	}
	homes := map[string]string{}
	if w.Player.Alive {
		homes[w.playerDeedID()] = w.Player.Home
	}
	for _, n := range w.NPCs {
		if !n.Dead {
			homes[n.ID] = n.Home
		}
	}
	assigned := map[string]bool{}
	for i := range w.Apartments {
		u := &w.Apartments[i]
		if u.Resident == "" {
			continue
		}
		if homes[u.Resident] != u.Building || assigned[u.Resident] {
			u.Resident = ""
		} else {
			assigned[u.Resident] = true
		}
	}
	ids := make([]string, 0, len(homes))
	for id := range homes {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if assigned[id] {
			continue
		}
		// Prefer a vacant flat this resident already owns when returning home.
		var chosen *ApartmentDeed
		for i := range w.Apartments {
			u := &w.Apartments[i]
			if u.Building == homes[id] && u.Resident == "" {
				if chosen == nil || u.Owner == id {
					chosen = u
				}
				if u.Owner == id {
					break
				}
			}
		}
		if chosen != nil {
			chosen.Resident = id
		}
	}
}

func (w *World) ApartmentPrice(u *ApartmentDeed) int {
	if u == nil {
		return 0
	}
	base := 1800
	if u.Building == "mercercourt" {
		base = 1200
	}
	return base * w.NeighborhoodPropertyIndex(u.Building) / 100
}
func (w *World) ApartmentOwnerName(u *ApartmentDeed) string {
	if u.Owner == w.playerDeedID() {
		return "You"
	}
	if n := w.NPC(u.Owner); n != nil {
		return n.Name
	}
	return "Independent broker"
}
func (w *World) BuyApartmentReadiness(u *ApartmentDeed) string {
	if u == nil || u.Resident != w.playerDeedID() {
		return "You can buy the apartment you currently rent"
	}
	if u.Owner != "independent" {
		return "This owner is not offering the apartment"
	}
	if w.Player.Cash < w.ApartmentPrice(u) {
		return "Not enough cash"
	}
	return ""
}
func (w *World) BuyApartment(id string) error {
	u := w.apartment(id)
	if why := w.BuyApartmentReadiness(u); why != "" {
		return fmt.Errorf("%s", why)
	}
	price := w.ApartmentPrice(u)
	if err := w.Pay(price); err != nil {
		return err
	}
	u.Owner = w.playerDeedID()
	w.Log("An apartment in your name", fmt.Sprintf("You paid $%d for apartment %d at %s. Your rent ends; the deed remains yours if you move out.", price, u.Number, placeName(u.Building)), "business")
	return nil
}
func (w *World) SellApartment(id string) error {
	u := w.apartment(id)
	if u == nil || u.Owner != w.playerDeedID() {
		return fmt.Errorf("this apartment is not yours")
	}
	price := w.ApartmentPrice(u) * 65 / 100
	w.Player.Cash += price
	u.Owner = "independent"
	w.Log("An apartment sold", fmt.Sprintf("The broker paid $%d for apartment %d at %s. The resident keeps their home and resumes paying rent.", price, u.Number, placeName(u.Building)), "business")
	return nil
}

// ApartmentDay makes at most one funded transaction. A cash-strapped NPC may
// sell and remain a tenant; a resident with savings can buy from the broker or
// that seller. No preview, purse estimate or random valuation creates money.
func (w *World) ApartmentDay() {
	w.SettleApartments()
	for i := range w.Apartments {
		u := &w.Apartments[i]
		if u.Owner == w.playerDeedID() {
			continue
		}
		if strings.HasPrefix(u.Owner, "player:") {
			u.Owner = "independent"
		}
		seller := w.NPC(u.Owner)
		if seller != nil && seller.Dead {
			u.Owner = "independent"
			seller = nil
		}
		if u.Owner != "independent" && seller == nil {
			continue
		}
		if seller != nil && w.HouseholdWealth(seller) >= 100 {
			continue
		}
		buyer := w.NPC(u.Resident)
		if buyer == nil || buyer.Dead || buyer.ID == u.Owner {
			continue
		}
		price := w.ApartmentPrice(u)
		if w.HouseholdWealth(buyer) < price+200 {
			continue
		}
		w.SpendHouseholdMoney(buyer, price)
		if seller != nil {
			seller.Purse += price
		}
		previous := w.ApartmentOwnerName(u)
		u.Owner = buyer.ID
		w.Log("An apartment changes hands", fmt.Sprintf("%s bought apartment %d at %s from %s for $%d. The resident's address and current journey are unchanged.", buyer.Name, u.Number, placeName(u.Building), previous, price), "business")
		return
	}
	// Owners with little cash can sell to another NPC with savings, retaining
	// their tenancy. Stable identity order makes this independent of roster order.
	buyers := w.Residents("apartment")
	buyers = append(buyers, w.Residents("mercercourt")...)
	sort.Slice(buyers, func(i, j int) bool { return buyers[i].ID < buyers[j].ID })
	for i := range w.Apartments {
		u := &w.Apartments[i]
		seller := w.NPC(u.Owner)
		if seller == nil || seller.Dead || w.HouseholdWealth(seller) >= 100 {
			continue
		}
		price := w.ApartmentPrice(u)
		for _, buyer := range buyers {
			if buyer.ID == seller.ID || w.HouseholdWealth(buyer) < price+500 {
				continue
			}
			w.SpendHouseholdMoney(buyer, price)
			seller.Purse += price
			u.Owner = buyer.ID
			w.Log("A private apartment sale", fmt.Sprintf("%s paid %s $%d for apartment %d at %s. Existing residents keep their tenancy.", buyer.Name, seller.Name, price, u.Number, placeName(u.Building)), "business")
			return
		}
	}
}

func (w *World) ApartmentMarket() []map[string]any {
	out := []map[string]any{}
	for i := range w.Apartments {
		u := &w.Apartments[i]
		owned := u.Owner == w.playerDeedID()
		home := u.Resident == w.playerDeedID()
		if !owned && !home {
			continue
		}
		resident := "Vacant"
		rent := 0
		if home {
			resident = "You"
		} else if n := w.NPC(u.Resident); n != nil && !n.Dead {
			resident = n.Name
			rent = w.NPCRent(n)
		}
		out = append(out, map[string]any{"id": u.ID, "building": u.Building, "number": u.Number, "address": placeName(u.Building), "owned": owned, "home": home, "available": u.Owner == "independent", "owner": w.ApartmentOwnerName(u), "resident": resident, "asking": w.ApartmentPrice(u), "offer": w.ApartmentPrice(u) * 65 / 100, "daily_rent": rent})
	}
	return out
}

func (w *World) ApartmentRentIncome() int {
	total := 0
	for _, n := range w.NPCs {
		if u := w.apartmentForResident(n.ID); u != nil && u.Building == n.Home && u.Owner == w.playerDeedID() {
			total += w.NPCRent(&n)
		}
	}
	return total
}
