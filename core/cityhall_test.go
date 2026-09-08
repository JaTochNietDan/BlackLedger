package core

import (
	"fmt"
	"testing"
)

func petitionerAtHall(t *testing.T) *World {
	t.Helper()
	w := New(37)
	w.MigrateLivingWorld()
	w.Player.Location, w.Player.Cash, w.Player.Respect = CityHall, 20000, 40
	return w
}

func TestThereIsSomebodyInThatBuilding(t *testing.T) {
	w := petitionerAtHall(t)
	for _, o := range Officials() {
		n := w.NPC(o.ID)
		if n == nil || n.Dead || n.Name != o.Name {
			t.Fatalf("%s was not a person in the city", o.ID)
		}
		if n.Voice == "" {
			t.Fatalf("%s had no voice", o.Name)
		}
	}
	// And they are people rather than machinery: they do not go out robbing.
	for i := 0; i < 200; i++ {
		w.PeopleDay()
	}
	for _, o := range Officials() {
		for _, r := range w.History {
			if containsName(r.Title, o.Name) && containsName(r.Title, "Taken from") {
				t.Fatalf("%s was out taking tills: %s", o.Name, r.Title)
			}
		}
	}
}

func TestArrangementsAreMadeAtTheExchangeAndNowhereElse(t *testing.T) {
	w := petitionerAtHall(t)
	if w.RetainerReadiness("mayor") != "" {
		t.Fatal("could not reach the mayor:", w.RetainerReadiness("mayor"))
	}
	for _, elsewhere := range []string{"bar", "docks", "club", "room"} {
		w.Player.Location = elsewhere
		if w.RetainerReadiness("mayor") == "" {
			t.Fatalf("arranged with the mayor at %s", elsewhere)
		}
	}
}

func TestNobodyInThatBuildingIsSeenWithSomebodyTooHot(t *testing.T) {
	w := petitionerAtHall(t)
	o, _ := OfficialByID("mayor")
	w.Player.Heat = o.Ceiling + 1
	if w.RetainerReadiness("mayor") == "" {
		t.Fatal("a mayor was seen with somebody the whole city was watching")
	}
	w.Player.Heat = o.Ceiling
	if w.RetainerReadiness("mayor") != "" {
		t.Fatal("refused at exactly the stated ceiling:", w.RetainerReadiness("mayor"))
	}
	// And a nobody cannot get a call returned at any price.
	w.Player.Respect = 0
	if w.RetainerReadiness("mayor") == "" {
		t.Fatal("a nobody had the mayor on a retainer")
	}
}

func TestACommissionerBuysRoomAndAMayorBuysMoney(t *testing.T) {
	w := petitionerAtHall(t)
	if w.RaidRelief() != 0 || w.LicenceTake() != 0 {
		t.Fatal("somebody was getting the benefit without paying for it")
	}
	if err := w.Retain("commissioner"); err != nil {
		t.Fatal(err)
	}
	if w.RaidRelief() != RetainerRelief {
		t.Fatalf("a commissioner absorbed %d attention", w.RaidRelief())
	}
	if err := w.Retain("mayor"); err != nil {
		t.Fatal(err)
	}
	if w.LicenceTake() != MayorTake {
		t.Fatalf("a mayor was worth %.2f", w.LicenceTake())
	}
	// Both join the daily bill.
	o1, _ := OfficialByID("commissioner")
	o2, _ := OfficialByID("mayor")
	if w.RetainerCost() != o1.Retainer+o2.Retainer {
		t.Fatalf("the arrangements cost $%d a day", w.RetainerCost())
	}
}

func TestSomebodyRicherCanOutbidYouWithoutSayingSo(t *testing.T) {
	w := petitionerAtHall(t)
	w.Retain("commissioner")
	for i := range w.Factions {
		w.Factions[i].Cash = 0
	}
	if w.Outbid("commissioner") || w.RaidRelief() == 0 {
		t.Fatal("an arrangement was worthless with nobody competing for it")
	}
	w.Factions[0].Cash = w.Player.Cash + OutbidBy + 1
	if w.Outbid("commissioner") {
		t.Fatal("a rich organization with no quarrel with the player bought the commissioner out from under them")
	}
	w.Factions[0].Goodwill = -40
	if !w.Outbid("commissioner") {
		t.Fatal("an organization far richer than the player bought nothing")
	}
	if w.RaidRelief() != 0 {
		t.Fatal("an outbid arrangement was still worth something")
	}
	// And the player is still paying for it, which is the point.
	if w.RetainerCost() == 0 {
		t.Fatal("being outbid stopped the bill")
	}
}

func TestNobodyInThatBuildingGoesDownWithYou(t *testing.T) {
	w := petitionerAtHall(t)
	w.Retain("mayor")
	o, _ := OfficialByID("mayor")
	w.Player.Heat = o.Ceiling + 5
	w.CityHallDay()
	if w.Retained("mayor") {
		t.Fatal("a mayor stayed bought while the city was watching")
	}
	if !w.hasRecord(o.Name + " is not taking calls") {
		t.Fatal("nobody told the player it had ended")
	}
	if w.RetainerCost() != 0 {
		t.Fatal("still paying somebody who had stopped answering")
	}
}

func TestAnArrangementCanBeEndedAndCostsTheOpeningAgain(t *testing.T) {
	w := petitionerAtHall(t)
	o, _ := OfficialByID("commissioner")
	cash := w.Player.Cash
	w.Retain("commissioner")
	if err := w.EndRetainer("commissioner"); err != nil {
		t.Fatal(err)
	}
	if w.Retained("commissioner") || w.RetainerCost() != 0 {
		t.Fatal("stopping did not stop it")
	}
	if err := w.EndRetainer("commissioner"); err == nil {
		t.Fatal("ended an arrangement that did not exist")
	}
	w.Retain("commissioner")
	if w.Player.Cash != cash-o.Opening*2 {
		t.Fatalf("reopening cost $%d rather than the opening twice", cash-w.Player.Cash)
	}
}

func TestKillingAManWithATitleIsTheLoudestThingInTheCity(t *testing.T) {
	w := petitionerAtHall(t)
	w.Retain("commissioner")
	heat := w.Player.Heat
	power, cash := w.Factions[0].Power, w.Factions[0].Cash
	if !w.Kill("commissioner", "Shot outside his own house.") {
		t.Fatal("a commissioner could not be killed")
	}
	if w.Player.Heat <= heat {
		t.Fatalf("attention went %d to %d", heat, w.Player.Heat)
	}
	if w.Factions[0].Power >= power || w.Factions[0].Cash >= cash {
		t.Fatal("the organizations that had nothing to do with it paid nothing")
	}
	if w.Retained("commissioner") || w.RetainerCost() != 0 {
		t.Fatal("still paying a dead man")
	}
	if !w.hasNewsKind("police") {
		t.Fatal("the paper did not carry it")
	}
	// And nothing they were worth survives them.
	if w.RaidRelief() != 0 {
		t.Fatal("a dead commissioner was still losing files")
	}
}

func TestNobodyInheritsAnArrangement(t *testing.T) {
	w := petitionerAtHall(t)
	w.Retain("mayor")
	w.Player.Alive = false
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	if next.Retained("mayor") || next.RetainerCost() != 0 {
		t.Fatal("the next arrival had the mayor on a retainer")
	}
}

func TestAnOlderSaveGetsAPeopleItNeverHad(t *testing.T) {
	w := New(11)
	w.NPCs = w.NPCs[:0]
	w.MigrateLivingWorld()
	for _, o := range Officials() {
		if w.NPC(o.ID) == nil {
			t.Fatalf("migration left %s out of the city", o.Name)
		}
	}
	if w.RetainerCost() != 0 {
		t.Fatalf("migration handed out arrangements costing $%d", w.RetainerCost())
	}
}

func TestALicenceIsWorthWhatItSays(t *testing.T) {
	const days = 30
	// Income is compared per elapsed minute rather than per campaign, because
	// an encounter pauses the clock and the two runs do not stop in the same
	// place. What is being measured is the rate, not the wall time.
	rate := func(retained bool) float64 {
		w := petitionerAtHall(t)
		w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
		if retained {
			if err := w.Retain("mayor"); err != nil {
				t.Fatal(err)
			}
		}
		// Both keep enough to pay every bill: an unpaid bill has consequences
		// of its own, and the retainer makes the licensed run's bill larger.
		w.Player.Cash, w.Player.Earned = 50000, 0
		start := w.Minute
		w.Advance(days * 1440)
		elapsed := w.Minute - start
		if elapsed == 0 {
			t.Fatal("no time passed at all")
		}
		return float64(w.Player.Earned) / float64(elapsed)
	}

	without, with := rate(false), rate(true)
	if with <= without {
		t.Fatalf("a licensed laundry earned $%.3f a minute against $%.3f unlicensed", with, without)
	}
	t.Logf("a laundry earns $%.3f a minute unlicensed and $%.3f with the mayor paid, which is %.0f%% more", without, with, (with/without-1)*100)
}
