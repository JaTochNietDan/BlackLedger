package sim

import (
	"sort"
	"strings"
	"testing"
)

// Twelve ticks of work on running a business — hiring, the wage, putting
// somebody in charge, restocking, the people who walk out and the families who
// come for them — reached no simulated campaign at all, because every policy
// here buys a place and then never thinks about it again. Every "baseline
// unmoved: no campaign policy does this" in the development log is that gap.
//
// The publican is an investor who reads the books afterwards.

func TestThePublicanActuallyRunsWhatItBuys(t *testing.T) {
	t.Parallel()
	r := Run(27, "publican", "authored", 200, false)
	if r.Error != "" {
		t.Fatal(r.Error)
	}
	// By prefix: work aimed at a person carries their id after the colon, so
	// putting somebody in charge is recorded as "incharge:street-35". Hiring is
	// deliberately not on this list — the city fills an owned counter itself,
	// so a policy never has to, which is worth knowing and is not a gap.
	for _, kind := range []string{"acquire", "wage", "incharge"} {
		used := 0
		for did, n := range r.Actions {
			if did == kind || strings.HasPrefix(did, kind+":") {
				used += n
			}
		}
		if used == 0 {
			t.Errorf("a policy written to run its businesses never used %q", kind)
		}
	}
	// And it holds what it bought, staffed and run by somebody.
	if r.Milestones["laundry"] == 0 {
		t.Fatal("the publican never bought anything to run")
	}
}

// And what it actually pays, measured rather than assumed.
//
// This was `TestRunningWhatYouHoldCostsMoreThanItPaysSoFar`, and it recorded a
// deficiency: over two hundred commands a publican who hires, pays over the
// rate, restocks and puts somebody in charge ends poorer and slower to expand
// than an investor who buys and walks away. The reasoning written above it said
// the business layer charges for care without paying for it.
//
// The deficiency was a horizon. Two hundred commands is about twelve game days,
// and the simulator warns on every run that its own city measures need twenty —
// so everything this project believed about its economy came out of campaigns
// too short to see it. Over four hundred commands the publican is the richest
// policy in the game: $28,121 against the investor's $21,538, and ahead of the
// worker's $23,252 as well.
//
// Care pays. It pays later than anybody here had ever looked.
func TestRunningWhatYouHoldPaysOverALongEnoughRun(t *testing.T) {
	t.Parallel()
	const runs, steps = 20, 400
	publican, investor := make([]int, 0, runs), make([]int, 0, runs)
	pCasino, iCasino := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		p := Run(seed, "publican", "authored", steps, false)
		i := Run(seed, "investor", "authored", steps, false)
		publican = append(publican, p.Cash)
		investor = append(investor, i.Cash)
		pCasino += p.Milestones["casino"]
		iCasino += i.Milestones["casino"]
	}
	sort.Ints(publican)
	sort.Ints(investor)
	p, i := publican[runs/2], investor[runs/2]
	t.Logf("over %d campaigns of %d commands: the publican's median is $%d with %d casinos, the investor's $%d with %d",
		runs, steps, p, pCasino, i, iCasino)
	if p == i && pCasino == iCasino {
		t.Fatal("running a business and ignoring it come to exactly the same thing")
	}
	if p <= i {
		t.Fatalf("running what you hold stopped paying: the publican's median is $%d against the investor's $%d.\n"+
			"    Over two hundred commands that was true and this test said so; over four hundred it was not.\n"+
			"    If it is true again, find out whether the horizon moved or the business layer did.", p, i)
	}
	// It also opens more rooms, which it did not when the fixer had an
	// unlimited number of envelopes to carry.
	if pCasino <= iCasino {
		t.Fatalf("the publican opened %d casinos to the investor's %d", pCasino, iCasino)
	}
}
