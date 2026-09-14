package core

import (
	"encoding/json"
	"testing"
)

func TestPoolTournamentRegularsTravelAndKeepTheirPlace(t *testing.T) {
	w := tournamentFixture(t)
	w.Minute = 1020
	// An isolated free evening: retain the real city's people and addresses,
	// remove duty/war/grudge claims so ordinary leisure can be observed.
	w.Grudges = nil
	for i := range w.NPCs {
		n := &w.NPCs[i]
		n.Rank = 0
		n.Role = ""
		n.Faction = ""
		n.Held = 0
		n.Purse = 100
		n.Post = "bar"
		n.Location = "bar"
		n.Heading = ""
		n.Sets = 0
		n.Arrives = 0
		n.Car = 0
		n.Drove = 0
	}
	ids := []string{}
	for i := range w.NPCs {
		if w.poolTournamentVisit(&w.NPCs[i]) {
			ids = append(ids, w.NPCs[i].ID)
		}
	}
	if len(ids) != 8 {
		t.Fatalf("expected eight regulars, got %d", len(ids))
	}
	raw, _ := json.Marshal(w)
	var clockWorld World
	if err := json.Unmarshal(raw, &clockWorld); err != nil {
		t.Fatal(err)
	}
	clockWorld.Minute = 1019
	clockWorld.Advance(61)
	for _, id := range ids {
		if clockWorld.NPC(id).Location != PoolPlace {
			t.Fatal("clock skipped tournament invitation/arrival", id, clockWorld.Minute)
		}
	}
	w.SetOut()
	for _, id := range ids {
		n := w.NPC(id)
		if n.Location != "bar" || n.Heading != PoolPlace || n.Arrives <= w.Minute {
			t.Fatal("attendance skipped travel", id, n.Heading)
		}
		if !w.poolTournamentVisit(n) {
			t.Fatal("trip reshuffled guests")
		}
	}
	w.Minute = 1080
	w.Arrivals()
	w.SetOut()
	for _, id := range ids {
		n := w.NPC(id)
		if n.Location != PoolPlace || n.Heading != "" {
			t.Fatal("regular did not wait for entry", id, n.Location, n.Heading, n.Arrives)
		}
	}
	if err := w.EnterPoolTournament(); err != nil {
		t.Fatal(err)
	}
	w.SetOut()
	for id := range w.PoolTournament.Deposits {
		if id == w.PoolTournament.PlayerID {
			continue
		}
		if w.NPC(id).Heading != "" {
			t.Fatal("entrant left funded event")
		}
	}
}

func TestPoolTournamentAttendanceRespectsScheduleAndDuty(t *testing.T) {
	w := tournamentFixture(t)
	w.Minute = 1020
	var selected *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		n.Purse = 100
		n.Post = "bar"
		if w.poolTournamentVisit(n) {
			selected = n
			break
		}
	}
	if selected == nil {
		t.Fatal("no available regular")
	}
	selected.Held = w.Minute + 60
	if w.poolTournamentVisit(selected) {
		t.Fatal("custody ignored")
	}
	selected.Held = 0
	for _, minute := range []int{1019, 1200, 1440 + 1080} {
		w.Minute = minute
		if w.poolTournamentVisit(selected) {
			t.Fatal("wrong schedule", minute)
		}
	}
	w.Minute = 1020
	w.Properties[PoolPlace].Condition = 0
	if w.poolTournamentVisit(selected) {
		t.Fatal("visit to unusable hall")
	}
}
