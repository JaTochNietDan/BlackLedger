package core

import (
	"strings"
	"testing"
)

func TestMarinerPurchaseKeepsTenantsAndStartsBusinessObligations(t *testing.T) {
	w := New(7)
	w.Event = nil
	w.Plots = nil
	w.NextPressure = 0
	w.Player.Cash = 10000
	w.Player.Respect = 30
	w.Player.Location = "room"
	home, _ := PlaceByID("room")
	if home.Cost != 0 || AcquisitionCost(w, "room") != MarinerFreehold {
		t.Fatal("freehold price replaced room rental")
	}
	residents := len(w.Residents("room"))
	before := w.Player.Cash
	next, err := Execute(w, Command{Kind: "acquire", Target: "room", Revision: w.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	if !next.Own("room") || before-next.Player.Cash != MarinerFreehold || next.Holdings() != 1 {
		t.Fatalf("purchase did not transfer at quoted price: cash%d", next.Player.Cash)
	}
	if len(next.Residents("room")) != residents {
		t.Fatal("purchase changed tenancies")
	}
	p := next.Properties["room"]
	if p.Staff != 3 || len(p.Hands) != 3 || p.Supply != 40 || next.Wages() != 18 || next.HomeCost("room") != 0 {
		t.Fatalf("business obligations: %+v wages%d", p, next.Wages())
	}
	found := map[string]bool{}
	for _, a := range next.Actions("room") {
		found[a.ID] = true
	}
	for _, id := range []string{"inspect", "repair", "hire", "layoff", "restock", "operate:clean"} {
		if !found[id] {
			t.Errorf("missing %s", id)
		}
	}
	if found["order"] || next.OrderReadiness("room") == "" {
		t.Fatal("lodging invents an unoccupied standing-order income")
	}
}

func TestMarinerRentRespondsToServiceAndDamage(t *testing.T) {
	w := New(7)
	n := &NPC{ID: "tenant", Home: "room"}
	p := w.Properties["room"]
	full := w.NPCRent(n)
	if full != 15 {
		t.Fatalf("full rent %d", full)
	}
	p.Supply = 0
	if w.NPCRent(n) >= full {
		t.Fatal("no service credit when supplies exhausted")
	}
	p.Supply = 40
	p.Trouble = true
	if w.NPCRent(n) >= full {
		t.Fatal("broken boiler has no rent consequence")
	}
	p.Trouble = false
	p.Staff = 0
	if w.NPCRent(n) >= full {
		t.Fatal("no staff has no rent consequence")
	}
	p.Staff = 3
	p.Condition = 0
	if w.NPCRent(n) != 0 {
		t.Fatal("destroyed building still charges rent")
	}
	p.Condition = 100
	p.Owner = "player:1"
	supply := p.Supply
	w.OperationsDay()
	if p.Supply >= supply {
		t.Fatal("lodging supplies never run down")
	}
}

func TestMarinerMigrationPreservesExistingRentalAccounts(t *testing.T) {
	w := New(7)
	w.Version = 16
	p := w.Properties["room"]
	p.Income = 0
	p.Staff = 0
	p.Supply = 0
	p.Owner = "player:1"
	p.Rents = map[string]*RentAccount{"tenant": {Day: 2, Arrears: 9, Collected: 6}}
	cash := w.Player.Cash
	w.MigrateLivingWorld()
	if p.Income != 15 || p.Staff != 3 || p.Supply != 40 || p.Rents["tenant"].Arrears != 9 || p.Owner != "player:1" || w.Player.Cash != cash {
		t.Fatal("migration lost ownership/accounts or failed to establish lodging")
	}
}

func TestMarinerCopyDescribesRentNotHourlyCash(t *testing.T) {
	w := New(7)
	w.Player.Cash = 10000
	w.Player.Respect = 30
	for _, a := range w.Actions("room") {
		if a.ID == "acquire" && (strings.Contains(a.Detail, "an hour") || !strings.Contains(a.Detail, "living tenants")) {
			t.Fatalf("misleading acquisition: %s", a.Detail)
		}
	}
}

func TestLodgingUpgradeDoesNotRefillEstablishedBusinesses(t *testing.T) {
	w := New(7)
	w.Version = 16
	w.Properties["room"].Income = 0
	w.Properties["laundry"].Staff = 0
	w.Properties["laundry"].Supply = 0
	w.MigrateLivingWorld()
	if w.Properties["laundry"].Staff != 0 || w.Properties["laundry"].Supply != 0 {
		t.Fatal("lodging upgrade gifted unrelated staff or supplies")
	}
	w.Version = SaveVersion
	w.Properties["room"].Supply = 0
	w.Properties["room"].Staff = 0
	w.MigrateLivingWorld()
	if w.Properties["room"].Supply != 0 || w.Properties["room"].Staff != 0 {
		t.Fatal("repeated migration refilled lodging")
	}
}
