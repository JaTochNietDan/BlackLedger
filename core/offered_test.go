package core

import (
	"strings"
	"testing"
)

// An action the game offers is an action the game accepts.
//
// cmd/apicheck found one violation by accident — a sitdown listed as available
// and then refused with a 409, because three hours pass between the check that
// enables the button and the check that runs the command. That fault is not
// special to the sitdown: every action whose work is done after the clock
// advances can refuse for a reason that became true during those hours. This
// states the rule for all of them at once.
func TestEveryOfferedActionIsAccepted(t *testing.T) {
	tried, offered := 0, map[string]bool{}
	for seed := 0; seed < 40; seed++ {
		w := New(uint32(seed)*2654435761 + 3)
		w.MigrateLivingWorld()
		// Varied states: money, standing, the clock in different places, and
		// enough elapsed for the city to have made its own arrangements.
		w.Player.Cash = 200 + seed*400
		w.Player.Respect = seed * 3
		w.Player.Health = 40 + seed
		w.Advance(seed * 137)
		for step := 0; step < 6; step++ {
			if !w.Player.Alive || w.Event != nil {
				break
			}
			here := w.Player.Location
			for _, a := range w.Actions(here) {
				if a.Disabled || a.ID == "travel" || a.ID == "new_life" {
					continue
				}
				target := a.Target
				if target == "" {
					target = here
				}
				tried++
				offered[a.ID] = true
				if _, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: a.ID, Target: target}); err != nil {
					// Money spent by an earlier action in this loop is not the
					// fault under test; each attempt runs from the same world.
					t.Fatalf("seed %d: %q was offered as available and refused with %q",
						seed, a.ID, err)
				}
			}
			// Move the world on so the next pass sees a different city.
			next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "rest", Target: here})
			if err != nil {
				break
			}
			w = next
		}
	}
	// Random states almost never sit a value exactly on a threshold, and that
	// is precisely where this class of fault lives: the sitdown was refused
	// because a family stood at -26 against a limit of -25. Reverting the fix
	// and running the pass above did not catch it. So the second pass puts the
	// world ON the boundaries and starts the clock where the hours an action
	// takes will cross the point at which organizations reconsider each other.
	for seed := 0; seed < 60; seed++ {
		w := New(uint32(seed)*2654435761 + 17)
		w.MigrateLivingWorld()
		w.Player.Cash, w.Player.Respect, w.Player.Health = SitdownFee, SitdownStanding, 100
		w.Player.Location = SitdownGround
		w.Minute = 720*3 - SitdownMinutes/2
		for i := range w.Factions {
			w.Factions[i].Goodwill = SitdownWelcome
		}
		for i := range w.Conflicts {
			w.Conflicts[i].State = "war"
			w.Conflicts[i].Hostility = 40 + seed%50
		}
		for _, a := range w.Actions(w.Player.Location) {
			if a.Disabled || a.ID == "travel" || a.ID == "new_life" {
				continue
			}
			target := a.Target
			if target == "" {
				target = w.Player.Location
			}
			tried++
			offered[a.ID] = true
			if _, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: a.ID, Target: target}); err != nil {
				t.Fatalf("on the boundary, seed %d: %q was offered as available and refused with %q", seed, a.ID, err)
			}
		}
	}

	t.Logf("%d offers across %d distinct actions", tried, len(offered))
	if len(offered) < 12 {
		t.Fatalf("only %d distinct actions were ever offered, so this proves little", len(offered))
	}
	_ = strings.TrimSpace
}
