package core

import "sort"

// Residential capacity is distinct from public footfall and workplace staffing.
// The apartment building includes shared flats as well as private tenancies.
var ResidentialCapacity = map[string]int{"room": 24, "apartment": 64, "estate": 1, "mercercourt": 48}

func IsRentalHome(id string) bool {
	return id == "room" || id == "apartment" || id == "mercercourt"
}

func (w *World) HousingStanding(n *NPC) int {
	if n == nil || n.Dead {
		return 0
	}
	return max(n.Purse, w.StandingPurse(n))
}

// SettleHousing assigns only missing/invalid homes. A visit, robbery, temporary
// absence or roster reorder cannot silently move an established resident.
func (w *World) SettleHousing() {
	used := map[string]int{}
	people := []*NPC{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead {
			n.Home = ""
			n.Accommodation = ""
			continue
		}
		people = append(people, n)
	}
	sort.Slice(people, func(i, j int) bool {
		a, b := people[i], people[j]
		if w.HousingStanding(a) != w.HousingStanding(b) {
			return w.HousingStanding(a) > w.HousingStanding(b)
		}
		return a.ID < b.ID
	})
	// Existing residents have priority over new arrivals, within actual capacity.
	for _, n := range people {
		reserved := 0
		if w.Player.Home == n.Home {
			reserved = 1
		}
		if ResidentialCapacity[n.Home] > used[n.Home]+reserved {
			used[n.Home]++
			continue
		}
		n.Home = ""
		n.Accommodation = ""
	}
	for _, n := range people {
		if n.Home != "" {
			continue
		}
		wealth := w.HousingStanding(n)
		choices := []string{"room", "mercercourt", "apartment"}
		if wealth >= 80 {
			choices = []string{"apartment", "mercercourt", "room"}
		}
		if wealth >= 500 {
			choices = []string{"estate", "apartment", "mercercourt", "room"}
		}
		for _, home := range choices {
			// Cypress is a private residence, never a stranger's rented bed in an
			// already-owned estate. Its holder's own household may live there.
			if home == "estate" {
				p := w.Properties[home]
				if p == nil || (p.Owner != "independent" && p.Owner != n.Faction) {
					continue
				}
				if w.Player.Home == home {
					continue
				}
			}
			reserved := 0
			if w.Player.Home == home {
				reserved = 1
			}
			if used[home]+reserved >= ResidentialCapacity[home] {
				continue
			}
			n.Home = home
			switch home {
			case "room":
				n.Accommodation = "Rented room"
			case "estate":
				n.Accommodation = "Private residence"
			default:
				if wealth >= 80 {
					n.Accommodation = "Private apartment"
				} else {
					n.Accommodation = "Shared flat"
				}
			}
			used[home]++
			break
		}
	}
}

func (w *World) Residents(id string) []*NPC {
	out := []*NPC{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead && n.Home == id {
			out = append(out, n)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (w *World) HousingShortage() int {
	count := 0
	for _, n := range w.NPCs {
		if !n.Dead && n.Home == "" {
			count++
		}
	}
	return count
}
