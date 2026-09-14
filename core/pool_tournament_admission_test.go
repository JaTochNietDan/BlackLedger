package core

import (
	"encoding/json"
	"testing"
)

func TestPoolTournamentScheduledEntryAndOneDrawPerEvening(t *testing.T) {
	w := tournamentFixture(t)
	w.Minute = poolTournamentHour - 1
	poolReject(t, w, Command{Kind: "pool_tournament_enter"})
	w.Minute++
	before, _ := json.Marshal(w)
	notice := w.PoolTournamentNotice().(map[string]any)
	after, _ := json.Marshal(w)
	if string(before) != string(after) || notice["can_enter"] != true {
		t.Fatal("notice changed state or entry unavailable", notice)
	}
	w = poolExecute(t, w, Command{Kind: "pool_tournament_enter"})
	if w.LastPoolTournamentSlot != poolTournamentHour || w.PoolTournament.Escrow != 100 || w.Player.Cash != 975 {
		t.Fatal("entry not fully funded")
	}
	poolReject(t, w, Command{Kind: "pool_tournament_enter"})
	w.Properties[PoolPlace].Condition = 0
	w.ReconcilePoolTournament()
	w.Properties[PoolPlace].Condition = 100
	poolReject(t, w, Command{Kind: "pool_tournament_enter"})
	w.Minute = poolTournamentHour + poolTournamentPeriod
	w = poolExecute(t, w, Command{Kind: "pool_tournament_enter"})
	if w.LastPoolTournamentSlot != w.Minute {
		t.Fatal("later evening unavailable")
	}
}
func TestPoolTournamentAdmissionDeadlineAndFunds(t *testing.T) {
	w := tournamentFixture(t)
	w.Minute = poolTournamentHour + poolTournamentWindow
	opens, _, _, _ := w.poolTournamentAdmission()
	if opens != poolTournamentHour+poolTournamentPeriod {
		t.Fatal("deadline did not advance notice")
	}
	poolReject(t, w, Command{Kind: "pool_tournament_enter"})
	w.Minute = poolTournamentHour
	for i := range w.NPCs {
		if w.NPCs[i].ID != "leo" && w.NPCs[i].ID != "mara" && w.NPCs[i].ID != "elena" {
			w.NPCs[i].Location = "bar"
		}
	}
	w.NPC("elena").Purse = 24
	poolReject(t, w, Command{Kind: "pool_tournament_enter"})
	w.NPC("elena").Purse = 300
	w.Player.Cash = 24
	poolReject(t, w, Command{Kind: "pool_tournament_enter"})
	w.Player.Location = "bar"
	if w.PoolTournamentNotice() != nil {
		t.Fatal("remote notice")
	}
}
