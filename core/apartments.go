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
	for _, building := range []string{"apartment", "mercercourt", "riverside"} {
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
	} else if u.Building == "riverside" {
		base = 1000
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
	if u == nil {
		return "This apartment is not in the registry"
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
	buyers := []*NPC{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead && n.Home != "" {
			buyers = append(buyers, n)
		}
	}
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

// A small broker board rotates as deeds sell: three occupied investments and
// one vacant flat per building, plus every player holding/current home. Reads
// never allocate units, evict tenants or change the order of the saved registry.
func (w *World) ApartmentListings() map[string]bool {
	listed := map[string]bool{}
	units := make([]*ApartmentDeed, 0, len(w.Apartments))
	for i := range w.Apartments {
		units = append(units, &w.Apartments[i])
	}
	sort.Slice(units, func(i, j int) bool { return units[i].ID < units[j].ID })
	occupied, vacant := map[string]int{}, map[string]int{}
	for _, u := range units {
		if u.Owner == w.playerDeedID() || u.Resident == w.playerDeedID() {
			listed[u.ID] = true
			continue
		}
		if u.Owner != "independent" {
			continue
		}
		if u.Resident == "" && vacant[u.Building] < 1 {
			listed[u.ID] = true
			vacant[u.Building]++
		} else if n := w.NPC(u.Resident); n != nil && !n.Dead && occupied[u.Building] < 3 {
			listed[u.ID] = true
			occupied[u.Building]++
		}
	}
	return listed
}

func (w *World) ApartmentMarket() []map[string]any {
	listed := w.ApartmentListings()
	out := []map[string]any{}
	for i := range w.Apartments {
		u := &w.Apartments[i]
		owned := u.Owner == w.playerDeedID()
		home := u.Resident == w.playerDeedID()
		if !listed[u.ID] {
			continue
		}
		location, ok := PlaceByID(u.Building)
		if !ok {
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
		out = append(out, map[string]any{"id": u.ID, "building": u.Building, "number": u.Number, "address": placeName(u.Building), "owned": owned, "home": home, "available": u.Owner == "independent", "owner": w.ApartmentOwnerName(u), "resident": resident, "asking": w.ApartmentPrice(u), "offer": w.ApartmentPrice(u) * 65 / 100, "daily_rent": rent, "neighborhood_index": w.NeighborhoodPropertyIndex(u.Building), "locked": location.District > w.District})
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
