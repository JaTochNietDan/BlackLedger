package store

import (
	"blackledger/core"
	"encoding/json"
	"reflect"
	"testing"
)

func TestCurrentVersionCrowdedSaveGainsRiversideHomes(t *testing.T) {
	w := core.New(7)
	for len(w.NPCs) < core.MaxPeople {
		if w.AddCivilian() == nil {
			t.Fatal("population ceiling")
		}
	}
	// A pre-Riverside save has the current schema, but no new property or deeds.
	delete(w.Properties, "riverside")
	units := w.Apartments[:0]
	for _, u := range w.Apartments {
		if u.Building != "riverside" {
			units = append(units, u)
		}
	}
	w.Apartments = units
	oldHomes := map[string]string{}
	for _, n := range w.NPCs {
		oldHomes[n.ID] = n.Home
	}
	data, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	next, err := decode(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if next.HousingShortage() != 0 || len(next.Apartments) != 384 {
		t.Fatal("current-version save did not gain homes and deeds")
	}
	if next.Minute != w.Minute || !reflect.DeepEqual(next.Player, w.Player) || next.Revision != w.Revision {
		t.Fatal("housing migration advanced campaign")
	}
	for _, n := range next.NPCs {
		if oldHomes[n.ID] != "" && n.Home != oldHomes[n.ID] {
			t.Fatal("existing tenancy moved")
		}
	}
	data, err = json.Marshal(next)
	if err != nil {
		t.Fatal(err)
	}
	again, err := decode(string(data))
	if err != nil {
		t.Fatal(err)
	}
	repeated, _ := json.Marshal(again)
	if string(repeated) != string(data) {
		t.Fatal("housing migration is not idempotent")
	}
}
