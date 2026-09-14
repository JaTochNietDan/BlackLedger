package core

import (
	"encoding/json"
	"testing"
)

func hostFixture(t *testing.T) *World {
	w := tournamentFixture(t)
	w.Properties[PoolPlace].Owner = "player:1"
	for _, id := range []string{"leo", "mara", "elena", "vittorio"} {
		n := w.NPC(id)
		if n == nil {
			continue
		}
		n.Location = PoolPlace
		n.Dead = false
		n.Held = 0
		n.Purse = 300
		n.Heading = ""
		n.Sets = 0
		n.Arrives = 0
	}
	return w
}
func TestPoolOwnerCanHostWithoutEnteringAndCollectCut(t *testing.T) {
	w := hostFixture(t)
	w.Player.Cash = 0
	w = poolExecute(t, w, Command{Kind: "pool_tournament_host", PoolHost: &PoolHostInput{Fee: 50, CutPercent: 20}})
	tour := w.PoolTournament
	if tour.Deposits[tour.PlayerID] != 0 || tour.Escrow != 200 || w.Player.Cash != 0 {
		t.Fatal("nonplaying owner charged")
	}
	for i := 0; i < 2; i++ {
		tour.Bracket.Matches[i].Rack.Concede(1)
	}
	w.ReconcilePoolTournament()
	tournamentWinningPosition(w, 2)
	w.Player.Location = "bar"
	w.Advance(2)
	if !tour.Settled || tour.PrizePaid != 160 || tour.HouseCutPaid != 40 || w.Player.Cash != 40 {
		t.Fatal("split not settled", tour.PrizePaid, tour.HouseCutPaid, w.Player.Cash)
	}
	w.ReconcilePoolTournament()
	if w.Player.Cash != 40 {
		t.Fatal("cut paid twice")
	}
}
func TestPoolOwnerEntryAndInvalidSettingsAreAtomic(t *testing.T) {
	w := hostFixture(t)
	for _, settings := range []*PoolHostInput{nil, {Fee: 9}, {Fee: 50, CutPercent: 51}, {Fee: 50, CutPercent: -1}} {
		poolReject(t, w, Command{Kind: "pool_tournament_host", PoolHost: settings})
	}
	w.Properties[PoolPlace].Owner = ""
	poolReject(t, w, Command{Kind: "pool_tournament_host", PoolHost: &PoolHostInput{Fee: 50}})
	w.Properties[PoolPlace].Owner = "player:1"
	w = poolExecute(t, w, Command{Kind: "pool_tournament_host", PoolHost: &PoolHostInput{Fee: 50, CutPercent: 20, Enter: true}})
	if w.Player.Cash != 950 || w.PoolTournament.Deposits[w.PoolTournament.PlayerID] != 50 {
		t.Fatal("owner did not pay same fee")
	}
	before, _ := json.Marshal(w)
	poolReject(t, w, Command{Kind: "pool_tournament_host", PoolHost: &PoolHostInput{Fee: 10}})
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("active terms changed")
	}
	w.Properties[PoolPlace].Condition = 0
	w.ReconcilePoolTournament()
	if w.Player.Cash != 1000 || w.PoolTournament.HouseCutPaid != 0 {
		t.Fatal("cancelled event took a cut")
	}
}

func TestPoolHostOffersAndWaitCommand(t *testing.T) {
	w := hostFixture(t)
	before, _ := json.Marshal(w)
	notice := w.PoolTournamentNotice().(map[string]any)
	offers := notice["host_offers"].([]map[string]any)
	if len(offers) < 4 {
		t.Fatal("owner did not see funded candidates")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("offers changed save")
	}
	w = poolExecute(t, w, Command{Kind: "pool_tournament_host", PoolHost: &PoolHostInput{Fee: 50, CutPercent: 20}})
	minute := w.Minute
	w = poolExecute(t, w, Command{Kind: "pool_tournament_wait"})
	if w.Minute != minute+10 {
		t.Fatal("wait duration")
	}
	for i := 0; i < 2; i++ {
		if w.PoolTournament.Bracket.Matches[i].Rack.Shots != 5 {
			t.Fatal("unattended table did not progress", i)
		}
	}
	w.Event = &Scene{ID: "pause", Title: "Interrupted"}
	poolReject(t, w, Command{Kind: "pool_tournament_wait"})
}
