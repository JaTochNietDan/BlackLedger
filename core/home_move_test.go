package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestOccupiedHouseMovePlansAndPreservesOtherState(t *testing.T) {
	w := New(7)
	w.Event = nil
	w.Plots = nil
	w.NextPressure = 0
	w.Minute = 600
	w.District = 2
	w.Player.Cash = 8000
	w.Player.Location = "estate"
	w.NPCs = []NPC{{ID: "resident", Name: "Rosa Varga", Home: "estate", Accommodation: "Private residence", Location: "bar", Post: "market", Purse: 600, Heading: "market", Arrives: 650}}
	before, _ := json.Marshal(w)
	changes, why := w.PlanHomeMove("estate")
	if why != "" || len(changes) != 1 || changes[0].To != "apartment" {
		t.Fatalf("plan %+v %s", changes, why)
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("reading move plan changed campaign")
	}
	found := false
	for _, a := range w.Actions("estate") {
		if a.ID == "move_home" {
			found = strings.Contains(a.Detail, "Rosa Varga will be rehoused at Ashbury Court")
		}
	}
	if !found {
		t.Fatal("rehousing absent from move offer")
	}
	next, err := Execute(w, Command{Kind: "move_home", Target: "estate", Revision: w.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	n := next.NPC("resident")
	if next.Player.Home != "estate" || n.Home != "apartment" || len(next.Residents("estate")) != 0 || n.Post != "market" || n.Location != "market" {
		t.Fatalf("move failed or teleported a journey: %+v", n)
	}
}

func TestNoRehousingWithoutAVacancy(t *testing.T) {
	w := New(7)
	w.Player.Home = "estate"
	w.Player.Location = "room"
	w.Properties["estate"].Owner = "player:1"
	w.NPCs = nil
	for i := 0; i < 136; i++ {
		home := "room"
		if i >= 24 {
			home = "apartment"
		}
		if i >= 88 {
			home = "mercercourt"
		}
		w.NPCs = append(w.NPCs, NPC{ID: fmt.Sprint(i), Name: fmt.Sprint(i), Home: home, Location: "bar", Post: "bar"})
	}
	changes, why := w.PlanHomeMove("room")
	if why == "" || len(changes) > 0 {
		t.Fatalf("overfilled move offered: %+v %q", changes, why)
	}
	before, _ := json.Marshal(w)
	if _, err := Execute(w, Command{Kind: "move_home", Target: "room", Revision: w.Revision, RequestID: ID()}); err == nil {
		t.Fatal("move with nowhere for resident accepted")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("rejected move mutated world")
	}
}

func TestHomeChangesDoNotTeleportOrEraseRentDebt(t *testing.T) {
	w := New(7)
	w.NPCs = []NPC{{ID: "resident", Name: "Resident", Home: "estate", Location: "bar", Post: "bar", Heading: "market", Sets: 800, Arrives: 830}}
	w.Properties["estate"].Rents = map[string]*RentAccount{"resident": {Arrears: 12}}
	changes, why := w.PlanHomeMove("estate")
	if why != "" {
		t.Fatal(why)
	}
	w.applyHomeChanges(changes)
	n := w.NPC("resident")
	if n.Home == "estate" || n.Location != "bar" || n.Post != "bar" || n.Heading != "market" || n.Sets != 800 || n.Arrives != 830 || w.Properties["estate"].Rents[n.ID].Arrears != 12 {
		t.Fatalf("rehousing altered unrelated state: %+v", n)
	}
}

func TestHousingChangeDuringPaperworkRefundsWithoutUndoingCity(t *testing.T) {
	w := New(7)
	w.Minute = 600
	w.Event = nil
	w.Plots = nil
	w.Tasks = nil
	w.NextPressure = 0
	w.WorldRNG = 1
	w.District = 2
	w.Player.Location = "estate"
	w.Player.Cash = 8000
	w.NPCs = []NPC{{ID: "resident", Name: "Resident", Home: "estate", Accommodation: "Private residence", Location: "bar", Post: "bar"}}
	w.Contracts = []Contract{{ID: "housing-interruption", Life: w.Life, Target: "resident", Tier: "specialist", Payer: "city", Due: 630}}
	next, err := Execute(w, Command{Kind: "move_home", Target: "estate", Revision: w.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	if !next.NPC("resident").Dead {
		t.Fatal("fixture did not change occupancy during paperwork")
	}
	if next.Player.Home != "room" || next.Own("estate") || next.Player.Cash != 8000 || next.Minute != 660 {
		t.Fatal("changed housing plan was committed, not refunded, or rolled back the city")
	}
	postponed := false
	for _, r := range next.History {
		if r.Title == "The move is postponed" {
			postponed = true
		}
	}
	if !postponed {
		t.Fatal("postponed move was not explained")
	}
}
