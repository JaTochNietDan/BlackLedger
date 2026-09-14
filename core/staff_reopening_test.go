package core

import "testing"

func TestClosedHiringDoesNotInventStaffWithoutPeople(t *testing.T) {
	for _, unpaid := range []bool{false, true} {
		w := New(53)
		p := w.Properties["laundry"]
		p.Owner = w.PlayerOrganizationID()
		p.Staff = 0
		p.Hands = nil
		if unpaid {
			p.Unpaid = PatienceRunsOut
			p.Staff = 3
		} else {
			p.Shorthanded = w.Minute + 1440
		}
		w.EmptyChairs()
		if p.Staff != 0 || len(p.Hands) != 0 {
			t.Fatal("blocked hiring restored phantom staff")
		}
		p.Unpaid = 0
		p.Shorthanded = 0
		w.EmptyChairs()
		if p.Staff == 0 || p.Staff != len(p.Hands) {
			t.Fatal("reopening did not hire named workers")
		}
	}
}
