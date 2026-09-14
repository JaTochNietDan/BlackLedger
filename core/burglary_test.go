package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func burglaryWorld() *World {
	w := propertyTrader()
	w.Player.Location = "mercercourt"
	w.Player.Respect = 40
	n := w.NPC("mara")
	n.Home = "mercercourt"
	n.Location = "bar"
	n.Heading = ""
	n.Faction = ""
	n.Purse = 90
	w.MeetPerson(n.ID)
	w.HouseholdSavings = map[string]HouseholdAccount{n.ID: {Cash: 180, Day: 1}}
	w.RNG = 1
	return w
}

func TestBurglaryStealsFiniteHomeCashWithoutTakingWalletOrOtherResidents(t *testing.T) {
	w := burglaryWorld()
	cash := w.Player.Cash
	purse := w.NPC("mara").Purse
	w.HouseholdSavings["leo"] = HouseholdAccount{Cash: 80, Day: 1}
	next, err := Execute(w, Command{Kind: "burgle:mara", Target: "mercercourt", Revision: w.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Cash-cash != 180 || next.HouseholdSavings["mara"].Cash != 0 || next.NPC("mara").Purse != purse || next.HouseholdSavings["leo"].Cash != 80 {
		t.Fatal("burglary did not debit only the targeted home account")
	}
	if next.Minute-w.Minute != BurglaryMinutes || next.Player.Heat != 7 || next.NeighborhoodPropertyIndex("mercercourt") >= 100 {
		t.Fatal("missing time, heat or local market consequence")
	}
	next.RNG = 1
	cash = next.Player.Cash
	if err := next.Burgle("mara"); err != nil {
		t.Fatal(err)
	}
	if next.Player.Cash != cash {
		t.Fatal("empty home regenerated loot")
	}
}

func TestHouseholdSavingsAreFundedBoundedAndAvailableForBills(t *testing.T) {
	w := burglaryWorld()
	w.HouseholdSavings = nil
	n := w.NPC("mara")
	n.Purse = 1000
	before := n.Purse
	w.HouseholdSavingsDay()
	a := w.HouseholdSavings[n.ID]
	if a.Cash <= 0 || a.Cash > 60 || a.Cash+n.Purse != before {
		t.Fatal("unfunded savings")
	}
	w.HouseholdSavingsDay()
	if w.HouseholdSavings[n.ID] != a {
		t.Fatal("saved twice in one day")
	}
	n.Purse = 0
	cash := a.Cash
	w.HouseholdBills()
	if n.Purse+w.HouseholdSavings[n.ID].Cash != cash || n.Purse != w.NPCLivingCost(n) {
		t.Fatal("household bills created money or ignored savings")
	}
	a = w.HouseholdSavings[n.ID]
	a.Cash = 999
	a.Day = 0
	w.HouseholdSavings[n.ID] = a
	n.Purse = 10000
	w.HouseholdSavingsDay()
	if w.HouseholdSavings[n.ID].Cash != 1000 {
		t.Fatal("savings cap exceeded")
	}
}

func TestHouseholdSavingsStayPrivateAndSurviveReload(t *testing.T) {
	w := burglaryWorld()
	saved, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var loaded World
	if err = json.Unmarshal(saved, &loaded); err != nil {
		t.Fatal(err)
	}
	if loaded.HouseholdSavings["mara"].Cash != 180 {
		t.Fatal("save lost cash")
	}
	snapshot, _ := json.Marshal(w.Public())
	if strings.Contains(string(snapshot), "household_savings") {
		t.Fatal("private balances leaked")
	}
	n := w.NPC("mara")
	n.Home = "apartment"
	if w.BurglaryReadiness("mara") == "" {
		t.Fatal("old home still burgleable")
	}
	w.Player.Location = n.Home
	if w.HouseholdSavings[n.ID].Cash != 180 {
		t.Fatal("moving lost household cash")
	}
}

func TestBurglaryFailureAndOccupiedHomeRetaliation(t *testing.T) {
	w := burglaryWorld()
	n := w.NPC("mara")
	empty := w.burglaryChance(n)
	n.Location = n.Home
	if w.burglaryChance(n) >= empty {
		t.Fatal("occupied home no harder")
	}
	cash := w.Player.Cash
	health := w.Player.Health
	// Choose the combat stream directly for a reproducible failed attempt.
	for seed := uint32(1); ; seed++ {
		trial := w.Clone()
		trial.RNG = seed
		if trial.Random() >= w.burglaryChance(n) {
			w.RNG = seed
			break
		}
	}
	if err := w.Burgle(n.ID); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash || w.HouseholdSavings[n.ID].Cash != 180 || w.Player.Health >= health || w.Player.Heat != 14 || n.Sore == 0 {
		t.Fatal("failure lacks injury, retained cash or recognized retaliation")
	}
}

func TestHomeAttackRequiresPhysicalPresenceAndCapturesHomeSetting(t *testing.T) {
	w := burglaryWorld()
	n := w.NPC("mara")
	if w.HomeStrikeReadiness(n.ID) == "" {
		t.Fatal("absent target attackable at home")
	}
	n.Location = n.Home
	n.Heading = "bar"
	n.Sets = 0
	n.Arrives = w.Minute + 30
	if w.HomeStrikeReadiness(n.ID) == "" {
		t.Fatal("departed target attackable at home")
	}
	n.Heading = ""
	n.Arrives = 0
	if why := w.HomeStrikeReadiness(n.ID); why != "" {
		t.Fatal(why)
	}
	w.RNG = 1
	w.Player.Weapon = 1
	next, err := Execute(w, Command{Kind: "home_strike:mara", Target: n.Home, Revision: w.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	captured := false
	for _, cue := range next.LastResult.Cues {
		if cue.Strike != nil && cue.Strike.Setting == "home" {
			captured = true
		}
	}
	if !next.NPC(n.ID).Dead || !captured {
		t.Fatal("home hit did not preserve its setting")
	}
	if _, err := Execute(next, Command{Kind: "home_strike:mara", Target: n.Home, Revision: next.Revision, RequestID: ID()}); err == nil {
		t.Fatal("dead resident attackable")
	}
}

func TestSavedHouseholdCashCanFundApartmentInsteadOfStrandingSavings(t *testing.T) {
	w := apartmentWorld()
	n := w.NPC("tenant")
	n.Purse = 500
	w.HouseholdSavings = map[string]HouseholdAccount{n.ID: {Cash: 1000, Day: 1}}
	// The other resident already owns their home, leaving this buyer's offer.
	w.apartmentForResident("buyer").Owner = "buyer"
	w.ApartmentDay()
	if w.apartmentForResident(n.ID).Owner != n.ID || w.HouseholdWealth(n) != 300 {
		t.Fatal("apartment purchase ignored savings or charged incorrectly")
	}
	before := w.HouseholdWealth(n)
	if w.SpendHouseholdMoney(n, 1000) || w.HouseholdWealth(n) != before {
		t.Fatal("unfunded household spending mutated cash")
	}
}
