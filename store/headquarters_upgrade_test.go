package store

import (
	"blackledger/core"
	"encoding/json"
	"testing"
)

func TestCurrentVersionFamilyGainsMissingHeadquartersOnRead(t *testing.T) {
	s := testStore(t)
	w := core.New(59)
	w.Player.Respect = 45
	w.Properties["laundry"].Owner = w.PlayerOrganizationID()
	w.Incorporate()
	w.PlayerOrganization().Headquarters = ""
	raw, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.DB.Exec("UPDATE campaign SET state=? WHERE id=1", string(raw)); err != nil {
		t.Fatal(err)
	}
	got, err := s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if got.Headquarters(got.PlayerOrganizationID()) != "laundry" {
		t.Fatal("current-version family still lacks a base")
	}
	if got.Minute != w.Minute || got.Revision != w.Revision || got.Player.Cash != w.Player.Cash {
		t.Fatal("read advanced campaign")
	}
	var stored string
	if err = s.DB.QueryRow("SELECT state FROM campaign WHERE id=1").Scan(&stored); err != nil {
		t.Fatal(err)
	}
	if stored != string(raw) {
		t.Fatal("read wrote to database")
	}
}

func TestHeadquartersReadPreservesChoiceAndDoesNotFormIndependentOwner(t *testing.T) {
	for _, mode := range []string{"chosen", "lost", "independent"} {
		t.Run(mode, func(t *testing.T) {
			w := core.New(59)
			w.Player.Respect = 45
			w.Properties["laundry"].Owner = w.PlayerOrganizationID()
			w.Properties["garage"].Owner = w.PlayerOrganizationID()
			if mode != "independent" {
				w.Incorporate()
				w.PlayerOrganization().Headquarters = "laundry"
			}
			if mode == "lost" {
				w.Properties["laundry"].Owner = "independent"
			}
			raw, _ := json.Marshal(w)
			got, err := decode(string(raw))
			if err != nil {
				t.Fatal(err)
			}
			if mode == "independent" {
				if got.Incorporated() {
					t.Fatal("read auto-formed family")
				}
				return
			}
			if got.PlayerOrganization().Headquarters != "laundry" {
				t.Fatal("read replaced explicit base")
			}
			if mode == "lost" && got.Headquarters(got.PlayerOrganizationID()) != "" {
				t.Fatal("lost deed remained usable base")
			}
		})
	}
}
