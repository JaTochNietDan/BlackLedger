package core

import (
	"fmt"
	"sort"
)

// A voluntary move follows a funded vacant-flat purchase. Existing owner-
// occupants keep their home; no tenant is displaced and no trip is teleported.
func (w *World) purchaseNPCHome() bool {
	rank := map[string]int{"room": 0, "riverside": 1, "mercercourt": 2, "apartment": 3}
	buyers := []*NPC{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if _, rental := rank[n.Home]; !n.Dead && rental {
			buyers = append(buyers, n)
		}
	}
	sort.Slice(buyers, func(i, j int) bool { return buyers[i].ID < buyers[j].ID })
	units := []*ApartmentDeed{}
	for i := range w.Apartments {
		units = append(units, &w.Apartments[i])
	}
	sort.Slice(units, func(i, j int) bool {
		if rank[units[i].Building] != rank[units[j].Building] {
			return rank[units[i].Building] > rank[units[j].Building]
		}
		return units[i].ID < units[j].ID
	})
	for _, n := range buyers {
		old := w.apartmentForResident(n.ID)
		if old != nil && old.Owner == n.ID {
			continue
		}
		for _, u := range units {
			if u.Owner != "independent" || u.Resident != "" || rank[u.Building] <= rank[n.Home] {
				continue
			}
			if p := w.Properties[u.Building]; p == nil || p.Condition <= 0 {
				continue
			}
			price := w.ApartmentPrice(u)
			if w.HouseholdWealth(n) < price+500 {
				continue
			}
			if !w.SpendHouseholdMoney(n, price) {
				continue
			}
			from := n.Home
			if old != nil {
				old.Resident = ""
			}
			u.Owner, u.Resident = n.ID, n.ID
			n.Home, n.Accommodation = u.Building, "Private apartment"
			w.Log("A home of their own", fmt.Sprintf("%s bought vacant apartment %d at %s for $%d and gave up their tenancy at %s. Their current journey and any old rent debt remain unchanged.", n.Name, u.Number, placeName(u.Building), price, placeName(from)), "business")
			return true
		}
	}
	return false
}
