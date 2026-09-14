package main

import (
	"blackledger/core"
	"encoding/json"
	"testing"
)

func TestHTTPTournamentCommandsAndBracket(t *testing.T) {
	a := testApp(t)
	if err := a.s.Change(func(w *core.World) error {
		w.Player.Location = core.PoolPlace
		w.Player.Cash = 1000
		w.Event = nil
		w.Plots = nil
		w.NextPressure = 0
		for _, id := range []string{"leo", "mara", "elena"} {
			n := w.NPC(id)
			n.Location = core.PoolPlace
			n.Purse = 300
			n.Dead = false
			n.Heading = ""
			n.Sets, n.Arrives, n.Held = 0, 0, 0
		}
		return w.StartPoolTournament([]string{"leo", "mara", "elena"}, 25)
	}); err != nil {
		t.Fatal(err)
	}
	for _, body := range []string{
		`{"request_id":"tour-no-index","revision":0,"kind":"pool_tournament_place","pool":{"x":0.4,"y":0.4}}`,
		`{"request_id":"tour-other-table","revision":0,"kind":"pool_tournament_place","pool_game":1,"pool":{"x":0.4,"y":0.4}}`,
		`{"request_id":"tour-forged","revision":0,"kind":"pool_tournament_opponent","pool_game":1,"pool":{"speed":8}}`,
	} {
		r := request(a, "POST", "/api/action", body)
		if r.Code != 409 {
			t.Fatal(r.Code, r.Body.String())
		}
	}
	for _, body := range []string{
		`{"request_id":"tour-place","revision":0,"kind":"pool_tournament_place","pool_game":0,"pool":{"x":0.4,"y":0.4}}`,
		`{"request_id":"tour-shot","revision":1,"kind":"pool_tournament_shot","pool_game":0,"pool":{"angle":1.5707963267948966,"speed":7.5}}`,
	} {
		r := request(a, "POST", "/api/action", body)
		if r.Code != 200 {
			t.Fatal(r.Code, r.Body.String())
		}
		retry := request(a, "POST", "/api/action", body)
		if retry.Code != 200 || retry.Body.String() != r.Body.String() {
			t.Fatal("receipt changed")
		}
		var v any
		if err := json.Unmarshal(r.Body.Bytes(), &v); err != nil {
			t.Fatal(err)
		}
		missing := map[string]string{}
		nothings(v, "", missing)
		if len(missing) > 0 {
			t.Fatal(missing)
		}
	}
	r := request(a, "GET", "/api/state", "")
	var v struct {
		PoolTournament struct {
			Pot   int `json:"pot"`
			Games []struct {
				Index int `json:"index"`
				Table *struct {
					Shots  int    `json:"shots"`
					Replay string `json:"replay"`
				} `json:"table"`
			} `json:"games"`
		} `json:"pool_tournament"`
	}
	if err := json.Unmarshal(r.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	if v.PoolTournament.Pot != 100 || len(v.PoolTournament.Games) != 3 || v.PoolTournament.Games[0].Table.Shots != 1 || v.PoolTournament.Games[0].Table.Replay == "" || v.PoolTournament.Games[2].Table != nil {
		t.Fatal("incomplete public bracket")
	}
}
