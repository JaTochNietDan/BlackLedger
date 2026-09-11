package main

import (
	"blackledger/core"
	"encoding/json"
	"fmt"
	"testing"
)

// Full command/API/store route; no injected money, reputation or death state.
func TestHTTPRiseFallAndNewLife(t *testing.T) {
	t.Parallel()
	a := testApp(t)
	count := 0
	state := func() *core.World {
		w, e := a.s.Read()
		if e != nil {
			t.Fatal(e)
		}
		return w
	}
	command := func(kind, target, decision string) {
		w := state()
		count++
		c := core.Command{Kind: kind, Target: target, Choice: decision, Revision: w.Revision, RequestID: fmt.Sprintf("campaign-step-%04d", count)}
		if w.Event != nil {
			c.Event = w.Event.ID
		}
		data, _ := json.Marshal(c)
		r := request(a, "POST", "/api/action", string(data))
		if r.Code != 200 {
			t.Fatalf("step %d %s/%s: %s", count, kind, target, r.Body.String())
		}
	}
	settle := func() {
		for i := 0; i < 4 && state().Event != nil; i++ {
			e := state().Event
			decision := "accept"
			if e.Kind == "warning" {
				decision = "acknowledge"
			}
			if e.Kind == "business_pressure" {
				decision = "pay"
				if state().Player.Cash < 60 {
					decision = "resist"
				}
			}
			if e.Kind == "attack" {
				t.Fatal("unexpected personal attack during peaceful progression")
			}
			command("choice", "", decision)
		}
	}
	travel := func(target string) {
		for tries := 0; state().Player.Location != target && tries < 4; tries++ {
			settle()
			command("travel", target, "")
			settle()
		}
		if state().Player.Location != target {
			t.Fatal("could not reach", target)
		}
	}
	// Envelopes first, and the pier for the rest of it. The fixer has three
	// favours a day; a campaign that wanted sixteen used to take them all from
	// him and now works for most of its money, which is what this city looks
	// like from the inside.
	carried := 0
	earn := func(jobs int) {
		for i := 0; i < jobs; i++ {
			if carried < core.CourierADay {
				travel("bar")
				command("courier", "bar", "")
				settle()
				carried++
				continue
			}
			travel("docks")
			command("dockwork", "docks", "")
			settle()
		}
	}
	earn(2)
	command("recruit", "bar", "")
	settle()
	// The keys to a laundry cost four times what they did, so the walk to it is
	// longer: a freehold is meant to be a thing a player builds up to rather
	// than an afternoon's courier work.
	earn(14)
	travel("laundry")
	command("acquire", "laundry", "")
	settle()
	if !state().Own("laundry") || len(state().Player.Crew) != 1 {
		t.Fatal("first organization not established")
	}
	earn(10)
	command("expand", "apartment", "")
	settle()
	travel("apartment")
	command("move_home", "apartment", "")
	settle()
	command("security", "apartment", "")
	settle()
	if state().Player.Home != "apartment" || state().Guard() != 2 {
		t.Fatal("housing/security progression failed")
	}
	travel("bar")
	command("contact", "bar", "")
	settle()
	command("contact", "bar", "")
	settle()
	// Persistent IDs and properties survive server-side reads throughout the run.
	worldID := state().ID
	for round := 0; round < 8 && state().Player.Alive; round++ {
		travel("club")
		command("provoke", "club", "")
		travel("apartment")
		for tick := 0; tick < 6 && state().Player.Alive && state().Event == nil; tick++ {
			command("rest", "apartment", "")
			if state().Event != nil && state().Event.Kind == "warning" {
				command("choice", "", "acknowledge")
			}
		}
		if state().Event != nil {
			if state().Event.Kind == "attack" {
				command("choice", "", "defend")
			} else {
				settle()
			}
		}
	}
	if state().Player.Alive {
		t.Fatal("seeded repeated exposed defense did not reach death")
	}
	oldName := state().Player.Name
	oldMinute := state().Minute
	command("new_life", "", "")
	w := state()
	if w.ID != worldID || w.Life != 2 || w.Player.Name == oldName || w.Player.Cash != 90 || w.Player.Respect != 0 || len(w.Player.Crew) != 0 || w.Own("laundry") || w.Properties["laundry"].Owner != "former:"+oldName || w.Minute <= oldMinute {
		t.Fatal("new life did not preserve city and reset personal authority")
	}
	t.Logf("Completed %d committed HTTP commands; death at minute %d, city %s persisted", count, oldMinute, worldID)
}
