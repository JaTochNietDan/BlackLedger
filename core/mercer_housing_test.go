package core

import (
	"fmt"
	"testing"
)

func TestMercerCourtAbsorbsGrowthWithoutMovingEstablishedHomes(t *testing.T) {
	w := New(7)
	old := map[string]string{}
	for _, n := range w.NPCs {
		old[n.ID] = n.Home
	}
	for i := 0; i < 30; i++ {
		w.NPCs = append(w.NPCs, NPC{ID: fmt.Sprintf("arrival-%02d", i), Name: "New arrival", Purse: 30, Location: "docks", Post: "docks"})
	}
	w.SettleHousing()
	if w.HousingShortage() != 0 || len(w.Residents("mercercourt")) == 0 {
		t.Fatal("new housing did not accommodate growth")
	}
	for _, n := range w.NPCs {
		if home, ok := old[n.ID]; ok && n.Home != home {
			t.Fatal("growth moved an existing tenant")
		}
	}
	for _, n := range w.Residents("mercercourt") {
		if _, established := old[n.ID]; !established && n.Location != "docks" {
			t.Fatal("lease teleported a resident")
		}
	}
	if len(w.Residents("mercercourt")) > 48 {
		t.Fatal("building exceeds capacity")
	}
}

func TestMercerCourtRentalAndPlayerMove(t *testing.T) {
	w := New(7)
	w.Event = nil
	w.Plots = nil
	w.Tasks = nil
	w.Contracts = nil
	w.Player.Location = "mercercourt"
	w.Player.Cash = 1000
	w.NPCs = []NPC{{ID: "tenant", Name: "Tenant", Home: "mercercourt", Accommodation: "Shared flat", Purse: 100, Location: "docks", Post: "docks"}}
	if w.NPCRent(&w.NPCs[0]) != 12 {
		t.Fatal("shared flat rent")
	}
	w.collectRent(&w.NPCs[0])
	if w.NPCs[0].Purse != 88 || w.Properties["mercercourt"].Rents["tenant"].Paid != 12 {
		t.Fatal("rent not paid from tenant purse")
	}
	next, err := Execute(w, Command{Kind: "move_home", Target: "mercercourt", Revision: w.Revision, RequestID: ID()})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Home != "mercercourt" || next.Player.Cash != 880 || next.HomeCost(next.Player.Home) != 25 || HomeRank(next.Player.Home) != 1 {
		t.Fatal("player lease price or housing tier incorrect")
	}
	if next.Own("mercercourt") {
		t.Fatal("renting gave away a building deed")
	}
	if next.Watchers() != 1 || next.Guard() != 1 {
		t.Fatal("apartment protection not applied")
	}
	if next.RentRegister("mercercourt") == nil {
		t.Fatal("missing tenant register")
	}
}

func TestOldCampaignGainsMercerCourtWithoutChangingExistingLeases(t *testing.T) {
	w := New(7)
	delete(w.Properties, "mercercourt")
	w.NPCs = []NPC{{ID: "old", Home: "apartment", Accommodation: "Private apartment", Location: "bar"}, {ID: "new", Purse: 30, Location: "docks"}}
	cash, minute := w.Player.Cash, w.Minute
	w.MigrateLivingWorld()
	if w.Properties["mercercourt"] == nil || w.Properties["mercercourt"].Owner != "independent" {
		t.Fatal("missing independently held new address")
	}
	if w.NPC("old").Home != "apartment" || w.NPC("old").Location != "bar" || w.Player.Cash != cash || w.Minute != minute {
		t.Fatal("migration changed existing campaign")
	}
}
