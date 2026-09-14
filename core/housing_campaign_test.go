package core

import (
	"fmt"
	"testing"
)

// Unattended economy simulation, not a played campaign: the observer is funded
// to avoid personal rent bankruptcy and dismisses player-only event prompts.
// NPC wages, trade, housing, recruitment and conflict use the real clock.
func TestHousingEconomyFundsNPCDeedsAndHomesForLaterResidents(t *testing.T) {
	for _, seed := range []uint32{7, 27, 61} {
		t.Run(fmt.Sprint(seed), func(t *testing.T) {
			w := New(seed)
			w.Player.Cash = 100000
			w.Plots, w.Tasks, w.Contracts = nil, nil, nil
			start := w.Minute
			for turn := 0; w.Player.Alive && w.Minute < start+90*1440 && turn < 1000; turn++ {
				w.Event = nil
				w.Advance(240)
			}
			if w.Minute < start+90*1440 {
				t.Fatal("economy simulation ended early")
			}
			if w.HousingShortage() != 0 {
				t.Fatalf("living residents have no homes: %d", w.HousingShortage())
			}
			deeds := 0
			for _, u := range w.Apartments {
				if n := w.NPC(u.Owner); n != nil && !n.Dead {
					deeds++
				}
			}
			if deeds == 0 {
				t.Fatal("ordinary earnings never funded a home purchase")
			}
			for id, capacity := range ResidentialCapacity {
				occupied := len(w.Residents(id))
				if w.Player.Alive && w.Player.Home == id {
					occupied++
				}
				if occupied > capacity {
					t.Fatalf("%s: %d occupants in %d places", id, occupied, capacity)
				}
			}
			t.Logf("90 days: %d living NPC-owned deeds, zero shortage", deeds)
		})
	}
}

func TestCommittedCommandAssignsHomesToNewResidents(t *testing.T) {
	w := New(7)
	w.Event, w.Plots, w.Tasks = nil, nil, nil
	n := w.AddCivilian()
	if n == nil {
		t.Fatal("new civilian missing")
	}
	id := n.ID
	act(t, &w, "wait", w.Player.Location)
	n = w.NPC(id)
	if n.Home == "" {
		t.Fatal("committed world has an unhoused newcomer despite vacancies")
	}
}
