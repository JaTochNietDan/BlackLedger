package core

import "testing"

func TestFuneralBillsFollowTheDeedAndAreFundedOnce(t *testing.T) {
	for _, owner := range []string{"player", "family", "person"} {
		t.Run(owner, func(t *testing.T) {
			w, at := theParlour(t)
			victim := w.AddCivilian()
			proprietor := w.AddCivilian()
			if victim == nil || proprietor == nil {
				t.Fatal("missing fixture people")
			}
			victimID := victim.ID
			victim = w.NPC(victimID)
			victim.Faction, victim.Purse = "", 300
			delete(w.HouseholdSavings, victimID)
			switch owner {
			case "player":
				w.Properties[at].Owner = w.PlayerOrganizationID()
			case "family":
				w.Properties[at].Owner = w.Factions[0].ID
			case "person":
				w.Properties[at].Owner = proprietor.ID
			}
			funds := func() int {
				if owner == "player" {
					return w.Player.Cash
				}
				return w.businessFunds(at)
			}
			before := funds()
			if !w.Kill(victimID, "Died at home.") {
				t.Fatal("death failed")
			}
			if got := funds() - before; got != FuneralCost-FuneralOwn {
				t.Fatalf("owner margin %d", got)
			}
			if w.NPC(victimID).Purse != 300-FuneralCost || !w.NPC(victimID).Buried {
				t.Fatal("estate bill or burial missing")
			}
			w.Bury(w.NPC(victimID))
			w.Kill(victimID, "Duplicate.")
			if funds()-before != FuneralCost-FuneralOwn {
				t.Fatal("duplicate funeral payment")
			}
		})
	}
}

func TestFuneralFamilyCoversOnlyTheEstateShortfall(t *testing.T) {
	w, at := theParlour(t)
	n := w.AddCivilian()
	n.Faction, n.Purse = w.Factions[0].ID, 40
	delete(w.HouseholdSavings, n.ID)
	w.Properties[at].Owner = w.PlayerOrganizationID()
	familyBefore, cashBefore := w.Factions[0].Cash, w.Player.Cash
	w.Kill(n.ID, "Died at home.")
	if w.Factions[0].Cash != familyBefore-(FuneralCost-40) || n.Purse != 0 {
		t.Fatal("wrong estate/family split")
	}
	if w.Player.Cash-cashBefore != FuneralCost-FuneralOwn {
		t.Fatal("wrong owner proceeds")
	}
}

func TestPoorEstateBuysOnlyAnAffordableFuneral(t *testing.T) {
	for _, available := range []int{0, 24, 25, 80} {
		w, at := theParlour(t)
		n := w.AddCivilian()
		n.Faction, n.Purse = "", available
		delete(w.HouseholdSavings, n.ID)
		w.Properties[at].Owner = w.PlayerOrganizationID()
		before := w.Player.Cash
		w.Kill(n.ID, "Died at home.")
		if available < BurialPurse {
			if n.Purse != available || w.Player.Cash != before {
				t.Fatal("unfunded burial created a bill")
			}
		} else {
			cost := (available*FuneralOwn + FuneralCost - 1) / FuneralCost
			if n.Purse != 0 || w.Player.Cash-before != available-cost {
				t.Fatal("partial funeral was not funded")
			}
		}
	}
}

func TestCrewFuneralWaitsForPlayerAndPaysTheProprietor(t *testing.T) {
	w, at, dead := aDeadHandOfYours(t)
	if w.NPC(dead).Buried {
		t.Fatal("crew was buried before the player's decision")
	}
	proprietor := w.AddCivilian()
	w.Properties[at].Owner = proprietor.ID
	before := w.HouseholdWealth(proprietor)
	w.Player.Location = at
	if err := w.BuryYourOwn(at, dead); err != nil {
		t.Fatal(err)
	}
	if w.HouseholdWealth(proprietor)-before != FuneralCost-FuneralOwn {
		t.Fatal("crew payment never reached the proprietor")
	}
	if err := w.BuryYourOwn(at, dead); err == nil {
		t.Fatal("funeral paid twice")
	}
}

func TestFuneralUsesSavingsAndNeverOverdrawsAFamily(t *testing.T) {
	w, at := theParlour(t)
	n := w.AddCivilian()
	n.Faction, n.Purse = w.Factions[0].ID, 10
	account := w.HouseholdSavings[n.ID]
	account.Cash = 20
	w.HouseholdSavings[n.ID] = account
	w.Factions[0].Cash = 30
	w.Properties[at].Owner = w.PlayerOrganizationID()
	before := w.Player.Cash
	w.Bury(n)
	if w.Player.Cash != before || n.Purse != 10 {
		t.Fatal("billed a living person")
	}
	w.Kill(n.ID, "Died at home.")
	if w.HouseholdWealth(n) != 0 || w.Factions[0].Cash != 0 {
		t.Fatal("available estate/family funds were not settled")
	}
	cost := (60*FuneralOwn + FuneralCost - 1) / FuneralCost
	if w.Player.Cash-before != 60-cost {
		t.Fatal("payment exceeded the available funds")
	}
}
