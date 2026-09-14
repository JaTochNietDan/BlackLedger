package core

import "fmt"

const BurglaryMinutes = 45

type HouseholdAccount struct {
	Cash int `json:"cash"`
	Day  int `json:"day"`
}

// Household cash is private saved state, funded from actual purses, not a
// loot roll or a public estimate. Residents retain it when they move home.
func (w *World) HouseholdSavingsDay() {
	if w.HouseholdSavings == nil {
		w.HouseholdSavings = map[string]HouseholdAccount{}
	}
	day := w.Minute/1440 + 1
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || ResidentialCapacity[n.Home] == 0 {
			continue
		}
		a := w.HouseholdSavings[n.ID]
		if a.Day == day {
			continue
		}
		a.Day = day
		reserve := max(50, 3*w.NPCLivingCost(n))
		// Keep enough household savings to reach the local deed price and a
		// modest reserve. A fixed $1000 ceiling made $1800 flats impossible
		// to buy through ordinary earnings, even over an entire campaign.
		limit := 1000
		if unit := w.apartmentForResident(n.ID); unit != nil {
			limit = max(limit, w.ApartmentPrice(unit)+500)
		}
		deposit := min(60, min(max(0, (n.Purse-reserve)/5), max(0, limit-a.Cash)))
		n.Purse -= deposit
		a.Cash += deposit
		w.HouseholdSavings[n.ID] = a
	}
}

func (w *World) BurglaryReadiness(id string) string {
	n := w.NPC(id)
	if n == nil || n.Dead || !w.Known(n) {
		return "You do not know a living resident by that name"
	}
	if ResidentialCapacity[n.Home] == 0 {
		return "They have no known residence"
	}
	if w.Player.Location != n.Home {
		return "You must be at their home address"
	}
	if w.Player.HeldUntil > w.Minute {
		return "You are in custody"
	}
	if w.Player.Health < 25 {
		return "You need at least 25 health to risk a break-in"
	}
	if n.Faction == w.PlayerOrganizationID() {
		return "They are one of your own people"
	}
	for _, c := range w.Player.Crew {
		if c.ID == id {
			return "They are one of your crew"
		}
	}
	return ""
}

func (w *World) residentAtHome(n *NPC) bool {
	return n != nil && !n.Dead && n.Home != "" && n.Location == n.Home && !w.Travelling(n) && n.Held <= w.Minute
}

func (w *World) burglaryChance(n *NPC) float64 {
	chance := .76 + min64(.12, float64(w.Presence())/500) - float64(n.Rank)/700
	if w.residentAtHome(n) {
		chance -= .30
	}
	if n.Home == "estate" {
		chance -= .12
	}
	return min64(.88, max64(.18, chance))
}

func (w *World) Burgle(id string) error {
	if why := w.BurglaryReadiness(id); why != "" {
		return fmt.Errorf("%s", why)
	}
	n := w.NPC(id)
	home := n.Home
	present := w.residentAtHome(n)
	success := w.Random() < w.burglaryChance(n)
	presentation := CueBurglary{Intruder: CueActor{ID: "player", Name: w.Player.Name}, Resident: CueActor{ID: n.ID, Name: n.Name}, ResidentPresent: present, Success: success}
	headline := "BREAK-IN AT " + upper(placeName(home))
	var detail string
	if success {
		a := w.HouseholdSavings[id]
		loot := max(0, a.Cash)
		presentation.Taken = loot
		a.Cash = 0
		if w.HouseholdSavings != nil {
			w.HouseholdSavings[id] = a
		}
		if loot > 0 {
			w.Earn(loot)
			detail = fmt.Sprintf("You took $%d from the cash %s kept at home. Their wallet and the other residents' belongings were untouched.", loot, n.Name)
		} else {
			detail = "You searched " + n.Name + "'s home and found no cash to take."
		}
		w.Player.Heat = min(100, w.Player.Heat+7)
	} else {
		injury := 8 + int(w.Random()*15)
		if present {
			injury += 8
		}
		before := w.Player.Health
		w.Player.Health = max(0, w.Player.Health-w.Absorb(injury))
		presentation.HealthLost = before - w.Player.Health
		w.Ruin(15)
		w.Player.Heat = min(100, w.Player.Heat+14)
		detail = fmt.Sprintf("The break-in at %s's home went wrong. You escaped with nothing and lost %d health.", n.Name, before-w.Player.Health)
	}
	// A resident who sees the intruder can name them; otherwise witnesses may
	// identify a failed intruder. There is no automatic knowledge from absence.
	identified := present || (!success && w.Random() < .45)
	presentation.Identified = identified
	if identified {
		w.Aggrieve(n.ID, 35, "the break-in at their home")
		w.answerFor(n, w.OwnHands(), 25)
		detail += " You were identified; " + n.Name + " knows who came through the door."
	}
	if !w.Player.Alive || w.Player.Health <= 0 {
		presentation.Fatal = true
		detail = "The break-in at " + n.Name + "'s home ended in a fatal confrontation."
		w.DieOf("a burglary", detail)
	}
	w.Log("Burglary at "+placeName(home), detail, "danger")
	w.Report("robbery", headline, "A resident's home at "+placeName(home)+" was broken into. Police are asking for witnesses.")
	w.Witness("robbery", home, detail, headline)
	w.VisualCues[len(w.VisualCues)-1].Burglary = &presentation
	return nil
}

// Cash at home is available for necessities before the day's wages and rent.
func (w *World) HouseholdBills() {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		a, ok := w.HouseholdSavings[n.ID]
		if !ok || n.Dead {
			continue
		}
		take := min(a.Cash, max(0, w.NPCLivingCost(n)-n.Purse))
		n.Purse += take
		a.Cash -= take
		w.HouseholdSavings[n.ID] = a
	}
}

func (w *World) HouseholdWealth(n *NPC) int {
	if n == nil {
		return 0
	}
	return max(0, n.Purse) + max(0, w.HouseholdSavings[n.ID].Cash)
}

// Used only after an affordability check. It returns false without mutation
// if a caller's quote is no longer funded.
func (w *World) SpendHouseholdMoney(n *NPC, amount int) bool {
	if n == nil || amount < 0 || w.HouseholdWealth(n) < amount {
		return false
	}
	pocket := min(max(0, n.Purse), amount)
	n.Purse -= pocket
	if left := amount - pocket; left > 0 {
		a := w.HouseholdSavings[n.ID]
		a.Cash -= left
		w.HouseholdSavings[n.ID] = a
	}
	return true
}
