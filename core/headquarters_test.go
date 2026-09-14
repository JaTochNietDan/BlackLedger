package core

import (
	"encoding/json"
	"testing"
)

func TestFamilyFormationRequiresAChosenOwnedBusiness(t *testing.T) {
	w := New(71)
	w.Event = nil
	w.Player.Respect = 80
	if w.HeadquartersReadiness("laundry", true) == "" {
		t.Fatal("formed without a business")
	}
	w.Properties["laundry"].Owner = w.PlayerOrganizationID()
	w.OrganizationDay()
	if w.Incorporated() {
		t.Fatal("formed automatically")
	}
	w.Player.Location = "laundry"
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "form_family", Target: "laundry"})
	if err != nil {
		t.Fatal(err)
	}
	if !next.Incorporated() || next.Headquarters(next.PlayerOrganizationID()) != "laundry" {
		t.Fatal("chosen base was not recorded")
	}
	if w.Incorporated() {
		t.Fatal("execute mutated original")
	}
	raw, err := json.Marshal(next)
	if err != nil {
		t.Fatal(err)
	}
	var loaded World
	if err = json.Unmarshal(raw, &loaded); err != nil {
		t.Fatal(err)
	}
	loaded.MigrateLivingWorld()
	if loaded.Headquarters(loaded.PlayerOrganizationID()) != "laundry" {
		t.Fatal("reload changed selected base")
	}
}

func TestHeadquartersCannotBeRemoteUnownedOrRuined(t *testing.T) {
	w := New(71)
	w.Event = nil
	w.Player.Respect = 80
	w.Properties["laundry"].Owner = w.PlayerOrganizationID()
	if w.HeadquartersReadiness("laundry", true) == "" {
		t.Fatal("formed remotely")
	}
	w.Player.Location = "laundry"
	w.Properties["laundry"].Condition = 59
	if w.HeadquartersReadiness("laundry", true) == "" {
		t.Fatal("formed in ruins")
	}
	w.Properties["laundry"].Condition = 100
	w.Player.Serves = "bellandi"
	if w.HeadquartersReadiness("laundry", true) == "" {
		t.Fatal("formed while serving another family")
	}
}

func TestFamilyBaseSurvivesSuccessionAndRelocatesOnlyWhenNeeded(t *testing.T) {
	w := New(71)
	f := w.faction("bellandi")
	base := w.Headquarters(f.ID)
	if base == "" {
		t.Fatal("NPC family has no base")
	}
	w.Properties["laundry"].Owner = f.ID
	w.Properties["laundry"].Income = 200
	w.SettleHeadquarters()
	if w.Headquarters(f.ID) != base {
		t.Fatal("richer acquisition moved established base")
	}
	w.Kill(w.Leader(f.ID).ID, "Died at home.")
	if w.Headquarters(f.ID) != base {
		t.Fatal("succession lost family base")
	}
	w.Properties[base].Owner = "independent"
	w.SettleHeadquarters()
	if w.Headquarters(f.ID) == base || w.Headquarters(f.ID) == "" {
		t.Fatal("NPC family did not relocate after losing base")
	}
}

func TestPlayerMustChooseReplacementAfterLosingHeadquarters(t *testing.T) {
	w := New(71)
	w.Event = nil
	w.Player.Respect = 80
	w.Player.Location = "laundry"
	for _, id := range []string{"laundry", "garage"} {
		w.Properties[id].Owner = w.PlayerOrganizationID()
	}
	if err := w.EstablishHeadquarters("laundry", true); err != nil {
		t.Fatal(err)
	}
	w.Properties["laundry"].Owner = "independent"
	w.OrganizationDay()
	if !w.Incorporated() || w.Headquarters(w.PlayerOrganizationID()) != "" {
		t.Fatal("lost base silently moved or dissolved family")
	}
	w.Player.Location = "garage"
	if err := w.EstablishHeadquarters("garage", false); err != nil {
		t.Fatal(err)
	}
	if w.Headquarters(w.PlayerOrganizationID()) != "garage" {
		t.Fatal("replacement not selected")
	}
}

func TestPlayerSuccessionKeepsTheChosenBase(t *testing.T) {
	w, member := testator(t)
	w.Player.Location = "laundry"
	if err := w.EstablishHeadquarters("laundry", false); err != nil {
		t.Fatal(err)
	}
	w.Die("Shot.")
	if got := w.Headquarters("estate:" + member.ID); got != "laundry" {
		t.Fatalf("successor moved the chosen base: %s", got)
	}
}
