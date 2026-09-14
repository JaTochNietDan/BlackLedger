package core

import "testing"

func TestResidentWalksHomeAndKeepsWorkAddress(t *testing.T) {
	w := New(7)
	w.Minute = 1440
	w.NPCs = []NPC{{ID: "resident", Name: "Resident", Location: "club", Post: "laundry", Home: "room"}}
	n := &w.NPCs[0]
	w.SetOut()
	if n.Heading != "room" || n.Location != "club" || n.Sets <= w.Minute || n.Sets > w.Minute+300 {
		t.Fatalf("home journey: %+v", n)
	}
	walksOut(w, n)
	street := w.OnTheStreet()
	if len(street) != 1 || street[0].ToID != "room" || street[0].Progress != 0 {
		t.Fatalf("street: %+v", street)
	}
	w.Minute = n.Arrives
	w.Arrivals()
	if n.Location != "room" || n.Post != "laundry" || n.Heading != "" {
		t.Fatalf("home arrival lost workplace: %+v", n)
	}
	w.SetOut()
	if n.Heading != "" {
		t.Fatal("resident leaves home before morning")
	}
	w.Minute = 1440 + HomeUntil
	w.SetOut()
	if n.Heading != "laundry" {
		t.Fatalf("morning destination: %+v", n)
	}
}

func TestHomeRoutinePreservesDutyAndCustody(t *testing.T) {
	w := New(7)
	w.Minute = 1440
	w.NPCs = []NPC{
		{ID: "manager", Location: "bar", Post: "laundry", Home: "room", Role: "Runs Bluebird Laundry"},
		{ID: "held", Location: "police", Post: "laundry", Home: "room", Held: 2000},
		{ID: "dead", Location: "bar", Post: "laundry", Home: "room", Dead: true},
	}
	w.SetOut()
	if w.NPCs[0].Heading != "laundry" || w.NPCs[1].Heading != "" || w.NPCs[2].Heading != "" {
		t.Fatalf("duties/custody: %+v", w.NPCs)
	}
}

func TestAdvanceDispatchesMorningCommuteAtBoundary(t *testing.T) {
	w := New(7)
	w.Minute = 1440 + HomeUntil - 1
	w.Event = nil
	w.NextPressure = 10000
	w.Plots = nil
	w.NPCs = []NPC{{ID: "resident", Name: "Resident", Location: "room", Post: "laundry", Home: "room"}}
	w.Advance(1)
	if w.Minute != 1440+HomeUntil || w.NPCs[0].Heading != "laundry" {
		t.Fatalf("morning clock did not dispatch: minute=%d npc=%+v", w.Minute, w.NPCs[0])
	}
}

func TestHomeJourneyIsRecordedDuringPlayerTravel(t *testing.T) {
	w := New(7)
	w.Minute = 1440
	w.District = 2
	w.Player.Location = "bar"
	w.Player.Cash = 10000
	w.Event = nil
	w.Plots = nil
	w.Tasks = nil
	w.Contracts = nil
	w.NextPressure = 0
	w.NPCs = []NPC{{ID: "resident", Name: "Resident", Location: "club", Post: "laundry", Home: "room"}}
	w.SetOut()
	departure := w.NPCs[0].Sets
	w.Minute = departure - 1
	next, err := Execute(w, Command{Revision: w.Revision, RequestID: ID(), Kind: "travel", Target: "dealer"})
	if err != nil {
		t.Fatal(err)
	}
	legs := next.LastResult.StreetTravel
	if len(legs) != 1 || legs[0].ToID != "room" || legs[0].FromMinute != departure || legs[0].Progress != 0 {
		t.Fatalf("home trip missing from player travel: %+v", legs)
	}
}

func TestReplacementOfficialKeepsOfficeDuty(t *testing.T) {
	w := New(7)
	w.Minute = 1440
	for _, office := range Officials() {
		n := &NPC{ID: "successor", Role: office.Role, Location: office.Place(), Post: office.Place(), Home: "room"}
		if _, ok := w.routine(n); ok {
			t.Fatalf("replacement %s abandoned duty for routine", office.Role)
		}
	}
}

func TestHomePurchaseDuringWalkContinuesFromOldAddress(t *testing.T) {
	w := homeBuyerWorld()
	w.Minute = 1440
	n := w.NPC("buyer")
	n.Location = "club"
	n.Heading = ""
	n.Sets = 0
	n.Arrives = 0
	w.SetOut()
	if n.Heading != "room" {
		t.Fatalf("initial trip %+v", n)
	}
	walksOut(w, n)
	arrives := n.Arrives
	if !w.purchaseNPCHome() {
		t.Fatal("purchase failed")
	}
	if n.Heading != "room" || n.Arrives != arrives || n.Location != "club" {
		t.Fatal("purchase changed in-flight trip")
	}
	w.Minute = arrives
	w.Arrivals()
	if n.Location != "room" || n.Heading != "apartment" || n.Post != "cabstand" || n.Sets <= w.Minute {
		t.Fatalf("no onward journey: %+v", n)
	}
	walksOut(w, n)
	w.Minute = n.Arrives
	w.Arrivals()
	if n.Location != "apartment" || n.Heading != "" || n.Post != "cabstand" {
		t.Fatalf("new home arrival: %+v", n)
	}
}

func TestOldHomeArrivalAfterDawnReturnsToWork(t *testing.T) {
	w := homeBuyerWorld()
	w.Minute = 1440 + HomeUntil
	n := w.NPC("buyer")
	n.Home = "apartment"
	n.Location = "club"
	n.Heading = "room"
	n.Errand = "heading home to The Mariner"
	n.Sets = 0
	n.Arrives = w.Minute
	w.Arrivals()
	if n.Location != "room" || n.Post != "cabstand" || n.Heading != "cabstand" {
		t.Fatalf("dawn detour lost work: %+v", n)
	}
}
