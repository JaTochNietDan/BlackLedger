package main

import (
	"blackledger/core"
	"encoding/json"
	"testing"
)

func TestHeadquartersOrderHTTPRetryDoesNotReserveOrDispatchTwice(t *testing.T) {
	a := testApp(t)
	if err := a.s.Change(func(w *core.World) error {
		w.Event = nil
		w.Plots = nil
		w.Tasks = nil
		w.Homes = nil
		w.Player.Location = "laundry"
		w.Player.Respect = 80
		w.Player.Cash = 5000
		w.Properties["laundry"].Owner = w.PlayerOrganizationID()
		w.Properties["garage"].Owner = w.PlayerOrganizationID()
		w.Properties["garage"].Supply = 0
		w.Player.Crew = []core.Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 80}}
		n := w.NPC("leo")
		n.Location = "laundry"
		n.Heading = ""
		n.Arrives = 0
		n.Held = 0
		w.Incorporate()
		return w.EstablishHeadquarters("laundry", false)
	}); err != nil {
		t.Fatal(err)
	}
	w, err := a.s.Read()
	if err != nil {
		t.Fatal(err)
	}
	body, _ := json.Marshal(core.Command{RequestID: "headquarters-retry", Revision: w.Revision, Kind: "crew_order:restock", Choice: "leo", Target: "garage"})
	first := request(a, "POST", "/api/action", string(body))
	if first.Code != 200 {
		t.Fatal(first.Body.String())
	}
	committed, err := a.s.Read()
	if err != nil {
		t.Fatal(err)
	}
	second := request(a, "POST", "/api/action", string(body))
	if second.Code != 200 || second.Body.String() != first.Body.String() {
		t.Fatal("retry changed committed response")
	}
	again, err := a.s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if len(again.CrewOrders) != 1 || again.Player.Cash != committed.Player.Cash || again.Revision != committed.Revision {
		t.Fatal("retry repeated order or payment")
	}
}
