package core

import "testing"

// People are the first thing in this game that costs money every day and can
// walk out. This measures whether keeping them is worth the bill.

func TestPeopleAreWorthWhatTheyCost(t *testing.T) {
	heavy(t)
	const runs, days = 200, 60
	measure := func(hire int, pay bool) (survived, holdings, lost int) {
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(seed)
			w.MigrateLivingWorld()
			w.WorldRNG = seed * 2654435761
			w.Player.Cash, w.Player.Health, w.Player.Contacts = 40000, 100, 3
			own(w, "laundry", "garage")
			w.Player.Respect = OrganizationStanding
			w.OrganizationDay()
			taken := 0
			for _, n := range w.Civilians() {
				if taken >= hire {
					break
				}
				if IsOfficial(n.ID) || w.isRoleHolder(n) {
					continue
				}
				w.Player.Location = n.Location
				n.Trust = 1
				if w.SignOn(n.ID) == nil {
					taken++
				}
			}
			// Somebody with a reason, so there is something to survive.
			for i := range w.Factions {
				if w.Factions[i].ID == w.PlayerOrganizationID() {
					continue
				}
				if c := w.Conflict(w.Factions[i].ID, w.PlayerOrganizationID()); c != nil {
					c.Hostility, c.State = 75, "war"
				}
			}
			for day := 0; day < days; day++ {
				if !pay {
					w.Player.Cash = 0
				}
				w.FactionTurn()
				w.FamilyDay()
				w.OrganizationDay()
				w.OwnPeopleDay()
				w.Minute += 1440
			}
			held := len(w.FamilyHoldings(w.PlayerOrganizationID()))
			holdings += held
			if held > 0 {
				survived++
			}
			lost += hire - len(w.OwnPeople())
		}
		return
	}

	aloneHeld := 0
	_, aloneHeld, _ = measure(0, true)
	_, backedHeld, _ := measure(3, true)
	_, unpaidHeld, unpaidLost := measure(3, false)

	if backedHeld <= aloneHeld {
		t.Fatalf("three people kept %d holdings across %d campaigns against %d alone", backedHeld, runs, aloneHeld)
	}
	if unpaidLost == 0 {
		t.Fatal("nobody ever left an organization that stopped paying them")
	}
	t.Logf("%d campaigns of %d days at war: alone keeps %d holdings, three people kept and paid keep %d, three people not paid keep %d and %d of them walk out",
		runs, days, aloneHeld, backedHeld, unpaidHeld, unpaidLost)
}
