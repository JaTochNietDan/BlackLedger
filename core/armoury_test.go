package core

import (
	"fmt"
	"testing"
)

func dealer(t *testing.T) *World {
	t.Helper()
	w := New(71)
	w.MigrateLivingWorld()
	w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
	w.Player.Location, w.Player.Cash = "laundry", 20000
	return w
}

func TestARoomLikeThatGoesUnderABusinessAndOnlyOne(t *testing.T) {
	t.Parallel()
	w := dealer(t)
	if w.ArmouryReadiness("club") == "" || w.ArmouryReadiness("casino") == "" {
		t.Fatal("a casino floor took an armoury")
	}
	if w.ArmouryReadiness("laundry") != "" {
		t.Fatal("a laundry of your own would not:", w.ArmouryReadiness("laundry"))
	}
	if err := w.BuildArmoury("laundry"); err != nil {
		t.Fatal(err)
	}
	if _, ok := w.TheArmoury(); !ok {
		t.Fatal("there was no room afterwards")
	}
	w.Properties["garage"].Owner = fmt.Sprintf("player:%d", w.Life)
	if w.ArmouryReadiness("garage") == "" {
		t.Fatal("built a second one")
	}
}

func TestGunsComeOffABoatAndAreNotSoldOnAPublicFloor(t *testing.T) {
	t.Parallel()
	if !TradesAt("docks", "arms") {
		t.Fatal("the waterfront did not deal in arms")
	}
	if TradesAt("market", "arms") {
		t.Fatal("the exchange was selling guns across a public floor")
	}
	if !TradesAt("market", "moonshine") {
		t.Fatal("the exchange stopped dealing in everything else")
	}
}

func TestCratesGoIntoTheRoomAndOutOfYourHands(t *testing.T) {
	t.Parallel()
	w := dealer(t)
	w.BuildArmoury("laundry")
	if w.StockReadiness() == "" {
		t.Fatal("stocked a room with nothing in hand")
	}
	w.Player.Stock = map[string]int{"arms": 12}
	if err := w.StockArmoury(); err != nil {
		t.Fatal(err)
	}
	if w.Stocked() != 12 || w.Holding("arms") != 0 {
		t.Fatalf("%d in the room, %d in hand", w.Stocked(), w.Holding("arms"))
	}
	// The room takes what it takes and no more.
	w.Player.Stock["arms"] = ArmouryHold
	w.StockArmoury()
	if w.Stocked() != ArmouryHold {
		t.Fatalf("the room held %d", w.Stocked())
	}
	if w.StockReadiness() == "" {
		t.Fatal("a full room took more")
	}
	// And it has to be carried there.
	w.Player.Location = "bar"
	if w.StockReadiness() == "" {
		t.Fatal("crates moved themselves across the city")
	}
}

func TestTheCustomerIsAWar(t *testing.T) {
	t.Parallel()
	w := dealer(t)
	w.BuildArmoury("laundry")
	w.Properties["laundry"].Crates = 40
	for i := range w.Conflicts {
		w.Conflicts[i].State, w.Conflicts[i].Hostility = "cold", 10
	}
	cash := w.Player.Cash
	w.ArmouryDay()
	if w.Player.Cash != cash {
		t.Fatal("somebody bought guns in peacetime")
	}
	if w.Stocked() != 40 {
		t.Fatal("crates left the room with nobody fighting")
	}

	// A war is a customer.
	w.Conflicts[0].State, w.Conflicts[0].Hostility = "war", 80
	power, before := w.Factions[0].Power, w.Player.Cash
	w.ArmouryDay()
	if w.Player.Cash <= before {
		t.Fatal("a war bought nothing")
	}
	if w.Stocked() >= 40 {
		t.Fatal("the room did not empty")
	}
	if w.Factions[0].Power <= power && w.Factions[1].Power <= power {
		t.Fatal("what they bought made them no stronger")
	}
	if w.Player.Heat == 0 {
		t.Fatal("a room full of crates drew no attention")
	}
}

func TestNobodyBuysWhatTheyCannotAffordAndAnEmptyRoomSellsNothing(t *testing.T) {
	t.Parallel()
	w := dealer(t)
	w.BuildArmoury("laundry")
	w.Conflicts[0].State, w.Conflicts[0].Hostility = "war", 80
	for i := range w.Factions {
		w.Factions[i].Cash = 0
	}
	w.Properties["laundry"].Crates = 40
	cash := w.Player.Cash
	w.ArmouryDay()
	if w.Player.Cash != cash || w.Stocked() != 40 {
		t.Fatal("somebody bought guns with no money")
	}
	w.Properties["laundry"].Crates = 0
	heat := w.Player.Heat
	w.ArmouryDay()
	if w.Player.Heat != heat {
		t.Fatal("an empty room drew attention")
	}
}

func TestArmingOneSideIsSomethingTheOtherSideFindsOut(t *testing.T) {
	t.Parallel()
	const runs = 300
	noticed := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := dealer(t)
		w.WorldRNG = seed * 2654435761
		w.BuildArmoury("laundry")
		w.Properties["laundry"].Crates = 40
		w.Conflicts[0].State, w.Conflicts[0].Hostility = "war", 80
		for i := range w.Factions {
			w.Factions[i].Cash, w.Factions[i].Goodwill = 50000, 0
		}
		w.ArmouryDay()
		for i := range w.Factions {
			if w.Factions[i].Goodwill < 0 {
				noticed++
				break
			}
		}
	}
	if noticed == 0 {
		t.Fatal("arming one side was never noticed by the other")
	}
	if noticed == runs {
		t.Fatal("arming one side was always noticed, so there is no risk in it")
	}
	t.Logf("across %d days of selling into a war, the other side worked out where the guns came from %d times", runs, noticed)
}

func TestAWarrantThatFindsTheRoomTakesEverything(t *testing.T) {
	t.Parallel()
	w := dealer(t)
	w.BuildArmoury("laundry")
	w.Properties["laundry"].Crates = 40
	condition, heat := w.Properties["laundry"].Condition, w.Player.Heat
	place, crates, found := w.ArmouryFound()
	if !found || crates != 40 || place != "Bluebird Laundry" {
		t.Fatalf("a search found %d crates at %q", crates, place)
	}
	if w.Stocked() != 0 {
		t.Fatal("crates survived the search")
	}
	if _, ok := w.TheArmoury(); ok {
		t.Fatal("the room survived the search")
	}
	if w.Properties["laundry"].Condition >= condition || w.Player.Heat <= heat {
		t.Fatal("being found with it cost nothing")
	}
	if _, _, again := w.ArmouryFound(); again {
		t.Fatal("a second search found a room that was already gone")
	}
}

func TestTheInterfaceIsToldWhatIsUnderTheFloor(t *testing.T) {
	t.Parallel()
	w := dealer(t)
	if w.ArmouryDescription()["held"] != false {
		t.Fatal("a player with no room was told they had one")
	}
	w.BuildArmoury("laundry")
	w.Properties["laundry"].Crates = 15
	d := w.ArmouryDescription()
	if d["held"] != true || d["crates"] != 15 || d["place"] != "Bluebird Laundry" {
		t.Fatalf("the interface was told %v", d)
	}
	if d["price"].(int) <= w.Good("arms").Price {
		t.Fatal("a war paid no more than the waterfront")
	}
}
