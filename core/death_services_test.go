package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDeathServicesSplitTheExistingBillWithoutCreatingMoney(t *testing.T) {
	w, chapel := theParlour(t)
	person := w.AddCivilian()
	person.Faction, person.Purse = "", FuneralCost
	delete(w.HouseholdSavings, person.ID)
	for _, id := range []string{chapel, "mortuary", "cemetery", "crematorium"} {
		w.Properties[id].Owner = w.PlayerOrganizationID()
	}
	before := w.Player.Cash
	disposition := deathDisposition(person.ID)
	custom := w.Custom(disposition)
	other := "cemetery"
	if disposition == other {
		other = "crematorium"
	}
	otherCustom := w.Custom(other)
	w.Kill(person.ID, "Died at home.")
	// 170 parlour margin plus15 care and30 disposition, with45 external costs.
	if w.Player.Cash-before != 215 || person.Purse != 0 {
		t.Fatal("service payments failed to conserve the260 bill")
	}
	if w.Custom(disposition) != custom+BurialTrade || w.Custom(other) != otherCustom {
		t.Fatal("funeral paid for both burial and cremation")
	}
	w.Bury(person)
	if w.Player.Cash-before != 215 {
		t.Fatal("duplicate service settlement")
	}
}

func TestDeathServicesPayFamilyAndPersonalOwners(t *testing.T) {
	w, _ := theParlour(t)
	person := w.AddCivilian()
	personID := person.ID
	proprietor := w.AddCivilian()
	person = w.NPC(personID)
	person.Faction, person.Purse = "", FuneralCost
	delete(w.HouseholdSavings, personID)
	w.Properties["mortuary"].Owner = w.Factions[0].ID
	w.Properties[deathDisposition(personID)].Owner = proprietor.ID
	familyBefore, personalBefore := w.Factions[0].Cash, w.HouseholdWealth(proprietor)
	w.Kill(personID, "Died at home.")
	if w.Factions[0].Cash-familyBefore != 15 || w.HouseholdWealth(proprietor)-personalBefore != 30 {
		t.Fatal("service income did not follow the deeds")
	}
}

func TestClosedDeathServiceDoesNotTakeACommission(t *testing.T) {
	for _, reason := range []string{"damage", "trouble", "staff", "supplies"} {
		w, _ := theParlour(t)
		person := w.AddCivilian()
		p := w.Properties["mortuary"]
		p.Owner = w.PlayerOrganizationID()
		switch reason {
		case "damage":
			p.Condition = 59
		case "trouble":
			p.Trouble = true
		case "staff":
			p.Staff = 0
		case "supplies":
			p.Supply = 0
		}
		before := w.Player.Cash
		w.deathServicePayments(person, FuneralOwn)
		if w.Player.Cash != before {
			t.Fatalf("%s provider took money", reason)
		}
	}
}

func TestNewDeathBusinessesHaveDeedsTradesAndSaveRepair(t *testing.T) {
	w := New(53)
	for _, id := range []string{"mortuary", "cemetery", "crematorium"} {
		if !personalBusiness(id) {
			t.Fatalf("%s cannot be personally acquired", id)
		}
		delete(w.Properties, id)
	}
	w.SettleNewPlaces()
	for _, id := range []string{"mortuary", "cemetery", "crematorium"} {
		p := w.Properties[id]
		if p == nil || p.Owner != "independent" || p.Staff == 0 || p.Supply == 0 || p.Income != PlaceIncome[id] {
			t.Fatalf("%s did not enter the saved city as a going concern", id)
		}
	}
}

func TestServiceBookReportsPaidWorkAcrossReloads(t *testing.T) {
	w, chapel := theParlour(t)
	person := w.AddCivilian()
	person.Faction, person.Purse = "", FuneralCost
	delete(w.HouseholdSavings, person.ID)
	w.Kill(person.ID, "Died at home.")
	if len(w.DeathServiceReceipts) != 3 {
		t.Fatal("missing service receipts")
	}
	raw, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var restored World
	if err = json.Unmarshal(raw, &restored); err != nil {
		t.Fatal(err)
	}
	restored.Properties["mortuary"].Owner = restored.PlayerOrganizationID()
	line := restored.fromBehindThisCounter("mortuary")
	if !strings.Contains(line, "$30") || !strings.Contains(line, "$15") || !strings.Contains(line, "1 paid arrangement") {
		t.Fatalf("counter lost its actual receipts: %s", line)
	}
	if strings.Contains(restored.fromBehindThisCounter(chapel), "receiving and preparation") {
		t.Fatal("receipts attributed to the wrong venue")
	}
	before := len(restored.DeathServiceReceipts)
	restored.Bury(restored.NPC(person.ID))
	if len(restored.DeathServiceReceipts) != before {
		t.Fatal("reload repeated service receipts")
	}
	restored.Minute += 8 * 1440
	if !strings.HasPrefix(restored.fromBehindThisCounter("mortuary"), "No paid work") {
		t.Fatal("old receipt reported as recent")
	}
}
