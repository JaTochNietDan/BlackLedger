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

// And what it actually costs, measured rather than assumed. Running what you
// hold is not a free lunch and it is not even, on these numbers, a better one
// inside two hundred commands: paying over the rate and restocking out of your
// own pocket leaves a publican poorer and slower to expand than an investor who
// buys and walks away.
//
// That is worth having in front of us rather than argued about. It says the
// business layer currently charges for care without paying for it — a manager
// saves the player walking to the shop, which costs a policy nothing, and the
// wage over the rate buys loyalty against pressures a 200-command campaign
// rarely lives long enough to feel.
func TestRunningWhatYouHoldCostsMoreThanItPaysSoFar(t *testing.T) {
	t.Parallel()
	const runs = 20
	publican, investor := make([]int, 0, runs), make([]int, 0, runs)
	pCasino, iCasino := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		p := Run(seed, "publican", "authored", 200, false)
		i := Run(seed, "investor", "authored", 200, false)
		publican = append(publican, p.Cash)
		investor = append(investor, i.Cash)
		pCasino += p.Milestones["casino"]
		iCasino += i.Milestones["casino"]
	}
	sort.Ints(publican)
	sort.Ints(investor)
	p, i := publican[runs/2], investor[runs/2]
	t.Logf("over %d campaigns: the publican's median is $%d with %d casinos, the investor's $%d with %d",
		runs, p, pCasino, i, iCasino)
	// The claim is only that the two play differently, which is what makes the
	// measurement worth taking every time the balance moves.
	if p == i && pCasino == iCasino {
		t.Fatal("running a business and ignoring it come to exactly the same thing")
	}
	if pCasino > iCasino {
		t.Fatalf("the publican expanded faster than the investor: %d casinos against %d", pCasino, iCasino)
	}
}
