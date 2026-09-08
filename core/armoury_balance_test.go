package core

import (
	"fmt"
	"testing"
)

// Selling into a war has to be the most profitable thing in the game and the
// most dangerous, or it is just another business. This measures both against
// the alternative: sitting on the same money and running a laundry.

func TestSellingIntoAWarIsTheBestAndWorstMoneyInTheCity(t *testing.T) {
	const runs, days = 200, 45
	earned, raided, ruined := 0, 0, 0
	quiet := 0
	for seed := uint32(1); seed <= runs; seed++ {
		// A dealer: a room, stock, and a city at war to sell into.
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
		w.Player.Location, w.Player.Cash = "laundry", 30000
		if err := w.BuildArmoury("laundry"); err != nil {
			t.Fatal(err)
		}
		w.Properties["laundry"].Crates = ArmouryHold
		w.Conflicts[0].State, w.Conflicts[0].Hostility = "war", 80
		for i := range w.Factions {
			w.Factions[i].Cash = 40000
		}
		start, peak := w.Player.Earned, 0
		for day := 0; day < days; day++ {
			w.ArmouryDay()
			w.PoliceDay()
			peak = max(peak, w.Player.Heat)
		}
		earned += w.Player.Earned - start
		if _, ok := w.TheArmoury(); !ok {
			raided++
		}
		if peak >= RaidThreshold {
			ruined++
		}

		// The same laundry, no room, no crates.
		plain := New(seed)
		plain.MigrateLivingWorld()
		plain.WorldRNG = seed * 2654435761
		plain.Properties["laundry"].Owner = fmt.Sprintf("player:%d", plain.Life)
		plain.Player.Cash = 30000
		before := plain.Player.Earned
		plain.Advance(days * 1440)
		quiet += plain.Player.Earned - before
	}
	if earned <= quiet {
		t.Fatalf("selling into a war earned $%d against $%d running the same laundry quietly", earned/runs, quiet/runs)
	}
	if ruined == 0 {
		t.Fatal("nobody who did this ever reached the attention the police act on")
	}
	t.Logf("%d campaigns of %d days: selling into a war earns $%d a campaign against $%d running the laundry quietly; %d passed the attention the police act on and %d lost the room to a warrant",
		runs, days, earned/runs, quiet/runs, ruined, raided)
}
