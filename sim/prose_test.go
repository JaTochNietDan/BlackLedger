package sim

import "testing"

// What a played city says, read back.
//
// The core sweeps forty campaigns that nobody touches, and it found two real
// faults the night it was written. What it cannot reach is everything a player
// provokes: a cell, a search, a crate going out of a room, a story kept out of
// the paper. Those sentences are written by commands, and the harness is the
// only thing in this project that issues thousands of them.
//
// So the harness reads itself back. Every policy, because they go to different
// places and the wrong sentence is usually in a room one of them never enters.
func TestNothingAPlayedCitySaysIsMalformed(t *testing.T) {
	t.Parallel()
	// Twenty-seven campaigns of seven hundred commands is a minute of wall
	// clock. The fast half of the gate runs three of them, which still reaches
	// a search and a cell; the full gate reads all nine policies.
	runs, floor := uint32(3), 27
	if testing.Short() {
		runs, floor = 1, 9
	}
	campaigns, bad := 0, 0
	for _, strategy := range []string{
		"worker", "investor", "defiant", "reckless", "thief",
		"smuggler", "racketeer", "publican", "magpie", "distiller",
	} {
		for seed := uint32(1); seed <= runs; seed++ {
			r := Run(seed*2654435761, strategy, "fixture", 700, false)
			campaigns++
			for _, line := range r.Malformed {
				bad++
				t.Errorf("%s seed %d: %s", strategy, seed, line)
			}
		}
	}
	if campaigns < floor {
		t.Fatalf("only %d campaigns were read, so this measures nothing", campaigns)
	}
	t.Logf("%d played campaigns read back, %d malformed sentences", campaigns, bad)
}
