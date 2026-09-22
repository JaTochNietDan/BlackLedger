package sim

import (
	"sort"
	"testing"
)

func TestWhatKillsThem(t *testing.T) {
	if testing.Short() {
		t.Skip("260-campaign death-cause diagnostic; run without -short to collect the report")
	}
	t.Parallel()
	causes := map[string]int{}
	dead, runs := 0, 0
	for _, who := range []string{
		"worker", "investor", "defiant", "reckless", "thief",
		"smuggler", "racketeer", "publican", "magpie", "distiller", "respectable", "diplomat", "soldier",
	} {
		for seed := uint32(1); seed <= 20; seed++ {
			runs++
			r := Run(seed*2654435761, who, "authored", 400, false)
			if r.Alive {
				continue
			}
			dead++
			causes[r.Died]++
		}
	}
	keys := []string{}
	for k := range causes {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(a, b int) bool { return causes[keys[a]] > causes[keys[b]] })
	t.Logf("%d of %d campaigns ended in a death, %d distinct causes", dead, runs, len(keys))
	for _, k := range keys {
		t.Logf("    %3d  %s", causes[k], k)
	}
}
