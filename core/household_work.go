package core

import (
	"fmt"
	"sort"
)

const (
	HouseholdWorkPay     = 45
	HouseholdWorkMinutes = 90
)

// The building's notice board carries one small repair booking each day.
// Selection is stable on reads and rotates through residents, including tenants.
func (w *World) HouseholdWorkOffer(id string) (*NPC, string) {
	if !IsRentalHome(id) {
		return nil, "Household repair work is posted at residential buildings"
	}
	prop := w.Properties[id]
	if prop == nil {
		return nil, "This building is unavailable"
	}
	if prop.HouseholdWorkDay == w.Minute/1440+1 {
		return nil, "Today's household repair booking has already been taken"
	}
	hour := w.Minute % 1440
	if hour < 8*60 || hour+HouseholdWorkMinutes > 18*60 {
		return nil, "Repair bookings are accepted between 08:00 and 16:30"
	}
	if prop.Condition < 40 || prop.Trouble {
		return nil, "The building needs major repairs before residents can book small household jobs"
	}
	residents := w.Residents(id)
	sort.Slice(residents, func(i, j int) bool { return residents[i].ID < residents[j].ID })
	for i := range residents {
		n := residents[(i+w.Minute/1440)%len(residents)]
		if w.HouseholdWealth(n) >= HouseholdWorkPay+100 {
			return n, ""
		}
	}
	return nil, "No resident has money set aside for a repair booking today"
}

func (w *World) StartHouseholdWork(id string) (string, error) {
	n, reason := w.HouseholdWorkOffer(id)
	if reason != "" {
		return "", fmt.Errorf("%s", reason)
	}
	w.Properties[id].HouseholdWorkDay = w.Minute/1440 + 1
	return n.ID, nil
}

// The agreed resident must still be able to pay when work finishes. Do not
// silently replace the customer if the simulation changes their circumstances.
func (w *World) CompleteHouseholdWork(id, resident string) {
	n := w.NPC(resident)
	if n == nil || n.Dead || n.Home != id || w.HouseholdWealth(n) < HouseholdWorkPay+100 {
		w.Log("Household work cancelled", "The resident's circumstances changed during the booking. No fee was paid; today's booking remains taken.", "work")
		return
	}
	w.SpendHouseholdMoney(n, HouseholdWorkPay)
	w.Earn(HouseholdWorkPay)
	if w.Player.Respect < DockName {
		w.Player.Respect++
	}
	w.Log("Small repairs finished", fmt.Sprintf("You fixed sticking drawers and loose hinges in %s's home at %s. They paid you $%d from their household funds. Another booking can be taken tomorrow.", n.Name, placeName(id), HouseholdWorkPay), "work")
}
