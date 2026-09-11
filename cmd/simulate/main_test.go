package main

import (
	"sync/atomic"
	"testing"
)

// The campaigns run at once now. Every one of them builds its own world and
// draws from its own streams; there is no package-level state in the core that
// anything writes after `init`, no global randomness, and no clock or
// environment read anywhere in the campaign path — so the only thing running
// them in sequence bought was the order the reports landed in, and an index
// buys that more cheaply. Eight hundred campaigns went from 414 seconds to 81
// on a machine with sixteen cores, with the output identical to the byte.
//
// What could go wrong is quiet: a piece of work skipped, a piece done twice, or
// a report landing in somebody else's slot. None of those would fail the build
// and all three would move the baseline, so they are checked here.

func TestEveryPieceOfWorkIsDoneExactlyOnce(t *testing.T) {
	t.Parallel()
	for _, workers := range []int{0, 1, 3, 8, 64} {
		const n = 40
		done := make([]int32, n)
		var total int32
		inParallel(workers, n, func(k int) {
			atomic.AddInt32(&done[k], 1)
			atomic.AddInt32(&total, 1)
		})
		if int(total) != n {
			t.Fatalf("%d workers: %d pieces of work for %d slots", workers, total, n)
		}
		for k, times := range done {
			if times != 1 {
				t.Fatalf("%d workers: slot %d was done %d times", workers, k, times)
			}
		}
	}
	// And nothing at all is not a deadlock.
	inParallel(4, 0, func(int) { t.Fatal("work was done where there was none") })
}

// And the reports are the same whether they are run at once or one after
// another. This is the whole claim: the speed is free because the campaigns
// never touch each other.
func TestRunningThemAtOnceChangesNothing(t *testing.T) {
	if testing.Short() {
		t.Skip("campaigns: run without -short to compare them")
	}
	t.Parallel()
	const runs = 3
	strategies := []string{"worker", "publican"}
	one := runCampaigns(1, strategies, runs, 1, "authored", 120, false, nil)
	many := runCampaigns(8, strategies, runs, 1, "authored", 120, false, nil)
	for k := range one {
		if one[k].Seed != many[k].Seed || one[k].Strategy != many[k].Strategy {
			t.Fatalf("slot %d holds %s/%d run one at a time and %s/%d run at once",
				k, one[k].Strategy, one[k].Seed, many[k].Strategy, many[k].Seed)
		}
		if one[k].Cash != many[k].Cash || one[k].Minutes != many[k].Minutes ||
			one[k].Alive != many[k].Alive || one[k].Respect != many[k].Respect {
			t.Fatalf("%s on seed %d: $%d over %d minutes one at a time, $%d over %d at once",
				one[k].Strategy, one[k].Seed, one[k].Cash, one[k].Minutes, many[k].Cash, many[k].Minutes)
		}
	}
	// And they are in the order they were asked for, not the order they
	// finished. Scrambling the slots leaves every report present and every
	// figure in it right, so nothing above would notice.
	for k, r := range many {
		if want := strategies[k/runs]; r.Strategy != want {
			t.Fatalf("slot %d holds a %s and should hold a %s", k, r.Strategy, want)
		}
		if want := uint32(1) + uint32(k%runs)*0x9e3779b9; r.Seed != want {
			t.Fatalf("slot %d holds seed %d and should hold %d", k, r.Seed, want)
		}
	}
	t.Logf("%d campaigns, identical run at once and one after another", len(one))
}
