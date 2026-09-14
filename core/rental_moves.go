package core

import (
	"fmt"
	"sort"
)

// One household reviews vacant rentals each day. Deeds are never transferred,
// tenants are never displaced, and an existing trip still reaches its endpoint.
func (w *World) RentalMoveDay() {
	people := []*NPC{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead && IsRentalHome(n.Home) {
			people = append(people, n)
		}
	}
	sort.Slice(people, func(i, j int) bool { return people[i].ID < people[j].ID })
	rank := map[string]int{"room": 0, "riverside": 1, "mercercourt": 2, "apartment": 3}
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
	for offset := range people {
		n := people[(offset+w.Minute/1440)%len(people)]
		old := w.apartmentForResident(n.ID)
		if old != nil && old.Owner == n.ID {
			continue
		}
		wealth := w.HouseholdWealth(n)
		// Prefer privacy, then consider the shared-flat leases already used
		// by the residential system. Their smaller reserve keeps affordable
		// housing reachable for ordinary room tenants.
		for _, accommodation := range []string{"Private apartment", "Shared flat"} {
			for _, u := range units {
				if u.Resident != "" || u.Building == n.Home {
					continue
				}
				prop := w.Properties[u.Building]
				if prop == nil || prop.Condition < 60 || prop.Trouble {
					continue
				}
				candidate := *n
				candidate.Home = u.Building
				candidate.Accommodation = accommodation
				rent := w.NPCRent(&candidate)
				reserveDays := 14
				if accommodation == "Shared flat" {
					reserveDays = 7
				}
				upgrading := rank[u.Building] > rank[n.Home] && wealth >= rent*reserveDays
				downsizing := rank[u.Building] < rank[n.Home] && wealth < w.NPCRent(n)*7 && rent < w.NPCRent(n)
				if rent <= 0 || (!upgrading && !downsizing) || wealth < rent*7 {
					continue
				}
				if old != nil {
					old.Resident = ""
				}
				u.Resident = n.ID
				n.Home = u.Building
				n.Accommodation = candidate.Accommodation
				w.Log("A new tenancy", fmt.Sprintf("%s moved into apartment %d at %s, renting for $%d a day. Their former apartment is now vacant.", n.Name, u.Number, placeName(u.Building), rent), "business")
				return
			}
		}
	}
}
