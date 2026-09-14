package main

import (
	"blackledger/billiards"
	"blackledger/core"
	"encoding/json"
	"testing"
)

func TestHTTPBilliardsIntentReplayAndReceipts(t *testing.T) {
	a := testApp(t)
	if err := a.s.Change(func(w *core.World) error {
		w.Player.Location = core.PoolPlace
		w.Player.Cash = 1000
		w.Event = nil
		w.Plots = nil
		w.NextPressure = 0
		n := w.NPC("leo")
		n.Location = core.PoolPlace
		n.Purse = 300
		n.Heading = ""
		n.Sets = 0
		n.Arrives = 0
		n.Held = 0
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	commands := []string{
		`{"request_id":"http-pool-start","revision":0,"kind":"pool_start","target":"leo","amount":40}`,
		`{"request_id":"http-pool-place","revision":1,"kind":"pool_place","pool":{"x":0.4,"y":0.4}}`,
		`{"request_id":"http-pool-shot","revision":2,"kind":"pool_shot","pool":{"angle":1.5707963267948966,"speed":7.5}}`,
	}
	for _, body := range commands {
		first := request(a, "POST", "/api/action", body)
		if first.Code != 200 {
			t.Fatal(first.Code, first.Body.String())
		}
		var payload any
		if err := json.Unmarshal(first.Body.Bytes(), &payload); err != nil {
			t.Fatal(err)
		}
		missing := map[string]string{}
		nothings(payload, "", missing)
		if len(missing) != 0 {
			t.Fatal("unexpected null public fields", missing)
		}
		retry := request(a, "POST", "/api/action", body)
		if retry.Code != 200 || first.Body.String() != retry.Body.String() {
			t.Fatal("retry changed receipt")
		}
	}
	out := request(a, "GET", "/api/state", "")
	var view struct {
		Revision int `json:"revision"`
		Pool     struct {
			Shots  int     `json:"shots"`
			Replay string  `json:"replay"`
			Width  float64 `json:"width"`
			Balls  []struct {
				ID       int        `json:"id"`
				Position [3]float64 `json:"position"`
				Pocket   int        `json:"pocket"`
			} `json:"balls"`
		} `json:"pool"`
	}
	if err := json.Unmarshal(out.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if view.Revision != 3 || view.Pool.Shots != 1 || len(view.Pool.Balls) != 16 || view.Pool.Width != billiards.Width {
		t.Fatal("incomplete public table")
	}
	if r, err := billiards.DecodeReplay(view.Pool.Replay); err != nil || len(r.Frames) < 2 {
		t.Fatal("invalid public replay", err)
	}
	before, err := a.s.Read()
	if err != nil {
		t.Fatal(err)
	}
	beforeJSON, _ := json.Marshal(before)
	for _, body := range []string{
		`{"request_id":"http-pool-stale","revision":2,"kind":"pool_shot","pool":{"angle":0,"speed":3}}`,
		`{"request_id":"http-pool-forged","revision":3,"kind":"pool_opponent","pool":{"angle":0,"speed":3}}`,
		`{"request_id":"http-pool-missing","revision":3,"kind":"pool_shot"}`,
	} {
		if r := request(a, "POST", "/api/action", body); r.Code != 409 {
			t.Fatal("invalid command accepted", r.Code, r.Body.String())
		}
	}
	if r := request(a, "POST", "/api/action", `{"request_id":"http-pool-malformed","revision":3,"kind":"pool_shot","pool":{"speed":"fast"}}`); r.Code != 400 {
		t.Fatal("malformed cue accepted", r.Code)
	}
	after, err := a.s.Read()
	if err != nil {
		t.Fatal(err)
	}
	afterJSON, _ := json.Marshal(after)
	if string(beforeJSON) != string(afterJSON) {
		t.Fatal("GET, stale or invalid requests changed save")
	}
}
