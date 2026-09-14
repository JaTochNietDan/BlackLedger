package core

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"
)

func TestStartingCityHasHomesWithoutMovingWorkOrPeople(t *testing.T) {
	w := New(7)
	used := map[string]int{}
	for _, n := range w.NPCs {
		if !n.Dead {
			if n.Home == "" {
				t.Fatalf("%s has no home", n.ID)
			}
			used[n.Home]++
		}
	}
	used[w.Player.Home]++
	for id, count := range used {
		if count > ResidentialCapacity[id] {
			t.Fatalf("%s over capacity: %d", id, count)
		}
	}
	before, _ := json.Marshal(w.NPCs)
	w.SettleHousing()
	after, _ := json.Marshal(w.NPCs)
	if string(before) != string(after) {
		t.Fatal("settling housing twice changed the city")
	}
	clone := w.Clone()
	for i, n := range w.NPCs {
		if clone.NPCs[i].Home != n.Home || clone.NPCs[i].Accommodation != n.Accommodation {
			t.Fatal("home lost in save round trip")
		}
	}
}

func TestHousingTracksStandingAndReleasesDeadResidents(t *testing.T) {
	w := New(1)
	w.NPCs = []NPC{{ID: "wealthy", Purse: 1000, Location: "bar", Post: "docks"}, {ID: "comfortable", Purse: 120, Location: "market"}, {ID: "boarder", Purse: 25, Location: "bar"}}
	w.SettleHousing()
	if w.NPC("wealthy").Home != "estate" || w.NPC("comfortable").Home != "apartment" || w.NPC("boarder").Home != "room" {
		t.Fatal("housing ignores wealth")
	}
	if w.NPC("wealthy").Location != "bar" || w.NPC("wealthy").Post != "docks" {
		t.Fatal("moving home teleported a person or changed their work")
	}
	w.NPC("boarder").Dead = true
	if len(w.Residents("room")) != 0 {
		t.Fatal("dead tenant counts as an occupant")
	}
	w.SettleHousing()
	if w.NPC("boarder").Home != "" {
		t.Fatal("dead resident retains an assigned room")
	}
}

func TestHousingMigrationIsIdempotentAndDoesNotTakePlayerEstate(t *testing.T) {
	w := New(1)
	w.Properties["estate"].Owner = "player:1"
	w.NPCs = []NPC{{ID: "rich", Purse: 10000}}
	w.SettleHousing()
	if w.NPCs[0].Home == "estate" {
		t.Fatal("stranger assigned to player's estate")
	}
	w.NPCs[0].Purse = 0
	before := w.NPCs[0].Home
	w.SettleHousing()
	if w.NPCs[0].Home != before {
		t.Fatal("temporary cash loss silently evicted tenant")
	}
}

func TestHousingCapacityShortageIsExplicitAndExistingResidentsStay(t *testing.T) {
	w := New(1)
	w.NPCs = nil
	capacity := -1 // Reserve the player’s room; low-standing NPCs cannot use Cypress.
	for id, n := range ResidentialCapacity {
		if IsRentalHome(id) {
			capacity += n
		}
	}
	for i := 0; i < capacity+15; i++ {
		w.NPCs = append(w.NPCs, NPC{ID: fmt.Sprintf("resident-%03d", i), Purse: 30})
	}
	w.SettleHousing()
	if got := w.HousingShortage(); got != 15 {
		t.Fatalf("expected 15 without homes after %d affordable places, got %d", capacity, got)
	}
	before := map[string]string{}
	for _, n := range w.NPCs {
		before[n.ID] = n.Home
	}
	slices.Reverse(w.NPCs)
	w.SettleHousing()
	for _, n := range w.NPCs {
		if n.Home != before[n.ID] {
			t.Fatal("roster order moved an established resident")
		}
	}
}

func TestOldSaveReceivesHomesWithoutChangingMoneyOrOwnership(t *testing.T) {
	w := New(1)
	w.Version = 14
	for i := range w.NPCs {
		w.NPCs[i].Home = ""
		w.NPCs[i].Accommodation = ""
	}
	w.Properties["laundry"].Owner = "player:1"
	cash := w.Player.Cash
	w.MigrateLivingWorld()
	if w.HousingShortage() != 0 {
		t.Fatal("migration left the starting population without accommodation")
	}
	if !w.Own("laundry") || w.Player.Cash != cash {
		t.Fatal("home assignment changed property or player money")
	}
	before, _ := json.Marshal(w.NPCs)
	w.SettleHousing()
	after, _ := json.Marshal(w.NPCs)
	if string(before) != string(after) {
		t.Fatal("migrated homes are not stable")
	}
}
