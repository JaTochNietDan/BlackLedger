package main

import (
	"blackledger/core"
	"testing"
)

func TestHTTPOwnerTournamentTermsAndEntryReceipt(t *testing.T) {
	a := testApp(t)
	setup := func(owned bool) error {
		return a.s.Change(func(w *core.World) error {
			w.Player.Location = core.PoolPlace
			w.Player.Cash = 1000
			w.Event = nil
			for i := range w.NPCs {
				w.NPCs[i].Location = "bar"
			}
			for _, id := range []string{"leo", "mara", "elena", "vittorio"} {
				n := w.NPC(id)
				n.Location = core.PoolPlace
				n.Purse = 300
				n.Dead = false
				n.Held = 0
				n.Heading = ""
				n.Sets = 0
				n.Arrives = 0
			}
			if owned {
				w.Properties[core.PoolPlace].Owner = "player:1"
			}
			return nil
		})
	}
	if err := setup(false); err != nil {
		t.Fatal(err)
	}
	body := `{"kind":"pool_tournament_host","request_id":"host-once","revision":0,"pool_host":{"fee":50,"cut_percent":20,"enter":false}}`
	if r := request(a, "POST", "/api/action", body); r.Code != 409 {
		t.Fatal("non-owner arranged event", r.Code)
	}
	if err := setup(true); err != nil {
		t.Fatal(err)
	}
	first := request(a, "POST", "/api/action", body)
	if first.Code != 200 {
		t.Fatal(first.Code, first.Body.String())
	}
	again := request(a, "POST", "/api/action", body)
	if again.Code != 200 || again.Body.String() != first.Body.String() {
		t.Fatal("owner entry receipt changed")
	}
	w, err := a.s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != 1000 || w.PoolTournament.Escrow != 200 || w.PoolTournament.HouseCutPercent != 20 || len(w.PoolTournament.Deposits) != 4 {
		t.Fatal("terms or funding changed")
	}
}
