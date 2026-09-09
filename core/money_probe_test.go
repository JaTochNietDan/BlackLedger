package core

import (
	"sort"
	"testing"
)

// A probe, not an assertion: what does a family's money actually do over a
// year, and is it a number worth making decisions from?
func TestProbeFamilyMoney(t *testing.T) {
	for _, days := range []int{30, 120, 400} {
		var all []int
		for _, seed := range []uint32{7, 31, 88, 149, 219} {
			w := New(seed)
			w.District = 2
			live(t, w, days*48)
			for _, f := range w.Factions {
				if f.ID == w.PlayerOrganizationID() {
					continue
				}
				all = append(all, f.Cash)
				income := 0
				for _, id := range w.FamilyHoldings(f.ID) {
					prop := w.Properties[id]
					income += prop.Income * prop.Condition / 100
				}
				t.Logf("  %-18s cash %8d  power %3d  holdings %d  people %2d  day's income %6d",
					f.Name, f.Cash, f.Power, len(w.FamilyHoldings(f.ID)), len(w.Members(f.ID)), income*24)
			}
		}
		sort.Ints(all)
		zero := 0
		for _, v := range all {
			if v == 0 {
				zero++
			}
		}
		t.Logf("day %3d: %2d families  min %8d  median %8d  max %9d  broke %d",
			days, len(all), all[0], all[len(all)/2], all[len(all)-1], zero)
	}
}
