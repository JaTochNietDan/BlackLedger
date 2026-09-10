package core

import (
	"sort"
	"strings"
	"testing"
)

// The existing scan gives one invented family a name carrying its own article
// and waits for the city to write about it. That surfaces a site only when the
// world happens to do the thing that writes it, so the same fault has come back
// four times in one night, one sentence at a time.
//
// This renames every organization in the city to carry a lower case article, so
// every passage that starts a sentence with an organization's name is wrong at
// once and all of them can be fixed together.
func TestNoFamilyNameEverBeginsASentenceInLowerCase(t *testing.T) {
	t.Parallel()
	faults := map[string]int{}
	total := 0
	for _, seed := range []uint32{31, 47, 88, 103, 219} {
		w := New(seed)
		w.District = 2
		for i := range w.Factions {
			f := &w.Factions[i]
			if !strings.HasPrefix(f.Name, "the ") {
				f.Name = "the " + f.Name
			}
		}
		live(t, w, 4000)
		// Anything the world invented after the rename gets the same treatment
		// on the way in, so splinters and successors are covered too.
		for i := range w.Factions {
			if !strings.HasPrefix(w.Factions[i].Name, "the ") {
				w.Factions[i].Name = "the " + w.Factions[i].Name
			}
		}
		live(t, w, 2000)
		for _, text := range writings(w) {
			total++
			for _, sentence := range strings.Split(text, ". ") {
				if strings.HasPrefix(sentence, "the ") {
					head := sentence
					if len(head) > 70 {
						head = head[:70]
					}
					faults[head]++
				}
			}
		}
	}
	if len(faults) > 0 {
		keys := make([]string, 0, len(faults))
		for k := range faults {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			t.Errorf("a sentence begins in lower case (%d times): %q", faults[k], k)
		}
	}
	t.Logf("read %d passages across five cities with every organization named lower case", total)
}
