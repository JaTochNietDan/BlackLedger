package core

import (
	"blackledger/billiards"
	"encoding/json"
	"math"
	"testing"
)

func TestPoolTournamentCommandsValidateTableAndIntent(t *testing.T) {
	w := fundedTournament(t)
	zero, one, future := 0, 1, 2
	for _, c := range []Command{
		{Kind: "pool_tournament_place", Pool: &PoolInput{X: .4, Y: .4}},
		{Kind: "pool_tournament_place", PoolGame: &one, Pool: &PoolInput{X: .4, Y: .4}},
		{Kind: "pool_tournament_place", PoolGame: &future, Pool: &PoolInput{X: .4, Y: .4}},
		{Kind: "pool_tournament_shot", PoolGame: &zero},
		{Kind: "pool_tournament_opponent", PoolGame: &zero},
		{Kind: "pool_tournament_opponent", PoolGame: &one, Pool: &PoolInput{Speed: 8}},
		{Kind: "pool_tournament_place", PoolGame: &zero, Target: "casino", Pool: &PoolInput{X: .4, Y: .4}},
	} {
		poolReject(t, w, c)
	}
	minute := w.Minute
	w = poolExecute(t, w, Command{Kind: "pool_tournament_place", PoolGame: &zero, Pool: &PoolInput{X: .4, Y: .4}})
	w = poolExecute(t, w, Command{Kind: "pool_tournament_shot", PoolGame: &zero, Pool: &PoolInput{Angle: math.Pi / 2, Speed: 7.5}})
	if w.Minute != minute+2 || w.PoolTournament.Bracket.Matches[0].Rack.Shots != 1 || w.PoolTournament.Escrow != 100 {
		t.Fatal("stroke/time/escrow incorrect")
	}
	if _, err := billiards.DecodeReplay(w.PoolTournament.Replays[0]); err != nil {
		t.Fatal(err)
	}
	w.Event = &Scene{ID: "pause", Title: "An interruption"}
	poolReject(t, w, Command{Kind: "pool_tournament_opponent", PoolGame: &one})
	w = poolExecute(t, w, Command{Kind: "pool_tournament_withdraw"})
	if !w.PoolTournament.Bracket.Withdrawn[w.PoolTournament.PlayerID] || w.Event == nil {
		t.Fatal("withdrawal lost interruption")
	}
}

func TestPoolTournamentFinalCommandAndPublicProjection(t *testing.T) {
	w := fundedTournament(t)
	before, _ := json.Marshal(w)
	view := w.PoolTournamentDescription().(map[string]any)
	games := view["games"].([]map[string]any)
	if len(games) != 3 || games[2]["table"] != nil || games[0]["player_seat"] != 0 || games[1]["player_seat"] != -1 {
		t.Fatal("incorrect bracket view")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("projection changed save")
	}
	w.Player.Location = "bar"
	if w.PoolTournamentDescription() != nil {
		t.Fatal("remote bracket exposed")
	}
	w.Player.Location = PoolPlace
	w.PoolTournament.Bracket.Matches[0].Rack.Concede(1)
	w.PoolTournament.Bracket.Matches[1].Rack.Concede(1)
	w.ReconcilePoolTournament()
	tournamentWinningPosition(w, 2)
	index := 2
	w = poolExecute(t, w, Command{Kind: "pool_tournament_shot", PoolGame: &index, Pool: &PoolInput{Angle: math.Pi, Speed: .6, Ball: 8, Pocket: 4}})
	if !w.PoolTournament.Settled || w.Player.Cash != 1075 || w.PoolTournament.Escrow != 0 {
		t.Fatal("physical championship did not pay")
	}
	poolReject(t, w, Command{Kind: "pool_tournament_shot", PoolGame: &index, Pool: &PoolInput{Angle: math.Pi, Speed: .6, Ball: 8, Pocket: 4}})
	view = w.PoolTournamentDescription().(map[string]any)
	if view["pot"] != 0 || view["winner"] != w.PoolTournament.PlayerID {
		t.Fatal("championship missing from view")
	}
}
