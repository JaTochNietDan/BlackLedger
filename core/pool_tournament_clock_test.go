package core

import (
	"blackledger/billiards"
	"encoding/json"
	"testing"
)

func TestPoolTournamentClockUsesPhysicalStrokesAndSurvivesSave(t *testing.T) {
	w := fundedTournament(t)
	minute := w.Minute
	w.Advance(1)
	if w.PoolTournament.Bracket.Matches[1].Rack.Shots != 0 {
		t.Fatal("early stroke")
	}
	data, _ := json.Marshal(w)
	var restored World
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	restored.Advance(1)
	if restored.Minute != minute+2 || restored.PoolTournament.Bracket.Matches[1].Rack.Shots != 1 || restored.PoolTournament.Bracket.Matches[0].Rack.Shots != 0 {
		t.Fatal("clock failed or played for human")
	}
	if _, err := billiards.DecodeReplay(restored.PoolTournament.Replays[1]); err != nil {
		t.Fatal(err)
	}
	w.Advance(1)
	a, _ := json.Marshal(w.PoolTournament)
	b, _ := json.Marshal(restored.PoolTournament)
	if string(a) != string(b) {
		t.Fatal("save changed scheduled shot")
	}
}
func TestPoolTournamentWatchedStrokeIsNotImmediatelyRepeated(t *testing.T) {
	w := fundedTournament(t)
	index := 1
	w = poolExecute(t, w, Command{Kind: "pool_tournament_opponent", PoolGame: &index})
	if w.PoolTournament.Bracket.Matches[1].Rack.Shots != 1 {
		t.Fatal("watched stroke doubled")
	}
	w.Advance(2)
	if w.PoolTournament.Bracket.Matches[1].Rack.Shots != 2 {
		t.Fatal("automatic next stroke missing")
	}
}
func TestPoolTournamentClockPhysicallySettlesAbsentPlayersFinal(t *testing.T) {
	w := fundedTournament(t)
	w.PoolTournament.Bracket.Withdraw(w.PoolTournament.PlayerID)
	w.PoolTournament.Bracket.Matches[1].Rack.Concede(1)
	w.ReconcilePoolTournament()
	tournamentWinningPosition(w, 2)
	w.Advance(2)
	if !w.PoolTournament.Settled || w.PoolTournament.Escrow != 0 || w.NPC("leo").Purse != 375 {
		t.Fatal("automatic physical final did not pay", w.PoolTournament)
	}
	cash := w.NPC("leo").Purse
	w.Advance(4)
	if w.NPC("leo").Purse != cash {
		t.Fatal("repeated payout")
	}
}

func TestPoolTournamentUnattendedDrawCompletesRealRacks(t *testing.T) {
	w := fundedTournament(t)
	if err := w.PoolTournament.Bracket.Withdraw(w.PoolTournament.PlayerID); err != nil {
		t.Fatal(err)
	}
	w.ReconcilePoolTournament()
	for i := 0; i < 120 && !w.PoolTournament.Settled; i++ {
		w.Advance(2)
	}
	tour := w.PoolTournament
	if !tour.Settled || tour.Voided || tour.Bracket.Winner == "" {
		t.Fatal("unattended tournament failed to finish", w.Minute)
	}
	for _, index := range []int{1, 2} {
		rack := tour.Bracket.Matches[index].Rack
		if rack == nil || rack.Shots < 2 || rack.Winner < 0 || tour.Replays[index] == "" {
			t.Fatal("round skipped physical play", index)
		}
		t.Logf("round %d: %d physical strokes", index, rack.Shots)
	}
	if tour.Escrow != 0 {
		t.Fatal("prize not settled")
	}
}

func TestPoolTournamentClockPausesAndReconcilesClosure(t *testing.T) {
	w := fundedTournament(t)
	minute := w.Minute
	w.Event = &Scene{ID: "interruption", Title: "Interrupted"}
	w.Advance(10)
	if w.Minute != minute || w.PoolTournament.Bracket.Matches[1].Rack.Shots != 0 {
		t.Fatal("played through paused time")
	}
	w.Event = nil
	w.Properties[PoolPlace].Condition = 0
	w.Advance(2)
	if !w.PoolTournament.Voided || w.PoolTournament.Bracket.Matches[1].Rack.Shots != 0 || w.Player.Cash != 1000 {
		t.Fatal("closure did not precede scheduled play")
	}
}
