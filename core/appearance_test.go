package core

import (
	"fmt"
	"testing"
)

func dressed(t *testing.T) *World {
	t.Helper()
	w := New(41)
	w.Event, w.District = nil, 9
	w.Player.Location = "tailor"
	w.Player.Cash = 5000
	return w
}

// Renamed with the rule it guards. Clothes were bought at the exchange, which
// sold everything else and so sold this too — the one purchase in this game
// about how a man is read, made at a counter between the fish and the cloth.
func TestClothesAreSoldAtTheTailorsAndNowhereElse(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	if w.DressReadiness(1) != "" {
		t.Fatal("the tailor refused to sell:", w.DressReadiness(1))
	}
	for _, elsewhere := range []string{"bar", "docks", "club", "room", "market"} {
		w.Player.Location = elsewhere
		if w.DressReadiness(1) == "" {
			t.Fatalf("%s was selling suits", elsewhere)
		}
		if err := w.BuyAttire(1); err == nil {
			t.Fatalf("bought a suit at %s", elsewhere)
		}
	}
}

func TestDressWorksUpwardAndRunsOut(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Cash = 100000
	for tier := 1; tier < len(attires); tier++ {
		if err := w.BuyAttire(tier); err != nil {
			t.Fatal(err)
		}
		if w.Player.Dress != tier || w.DressCondition() != 100 {
			t.Fatalf("tier %d at %d condition", w.Player.Dress, w.DressCondition())
		}
	}
	if w.DressReadiness(len(attires)) == "" {
		t.Fatal("the rail was selling something that does not exist")
	}
}

func TestPresenceIsWhatYouHaveDonePlusWhatYouAreWearing(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Respect = 10
	if w.Presence() != 10 {
		t.Fatalf("working clothes were worth %d", w.Presence()-10)
	}
	if err := w.BuyAttire(w.Player.Dress + 1); err != nil {
		t.Fatal(err)
	}
	if w.Presence() != 10+attires[1].Presence {
		t.Fatalf("a pressed suit was worth %d", w.Presence()-10)
	}
}

func TestWearTakesStandingAwayAndRuinRemovesItEntirely(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Cash = 100000
	w.BuyAttire(w.Player.Dress + 1)
	w.BuyAttire(w.Player.Dress + 1) // tailored
	full := w.Standing()
	if full == 0 {
		t.Fatal("a tailored suit was worth nothing")
	}
	w.Player.DressWear = 60
	if partial := w.Standing(); partial >= full || partial == 0 {
		t.Fatalf("a worn suit was worth %d against %d", partial, full)
	}
	w.Ruin(30) // now below Shabby
	if w.Standing() != 0 {
		t.Fatalf("a ruined suit was still worth %d", w.Standing())
	}
	if !w.hasRecord("The suit is finished") {
		t.Fatal("nobody was told the suit was finished")
	}
}

func TestABeatingCostsTheSuitAsWellAsTheHealth(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Cash = 100000
	w.BuyAttire(w.Player.Dress + 1)
	w.BuyAttire(w.Player.Dress + 1)
	before := w.DressCondition()
	w.Ruin(30)
	if w.DressCondition() >= before {
		t.Fatalf("condition went %d to %d", before, w.DressCondition())
	}
}

func TestWorkingClothesCannotBeRuined(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Ruin(80)
	if w.DressCondition() != 100 || w.Standing() != 0 {
		t.Fatalf("working clothes were damaged: %d", w.DressCondition())
	}
}

func TestGoodClothesWearOutAndAreNoticed(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Cash = 100000
	w.BuyAttire(w.Player.Dress + 1)
	w.BuyAttire(w.Player.Dress + 1)
	w.BuyAttire(w.Player.Dress + 1) // bespoke: noticed
	heat, condition := w.Player.Heat, w.DressCondition()
	w.DressDay()
	if w.DressCondition() != condition-DressUpkeep {
		t.Fatalf("a day of wear took %d", condition-w.DressCondition())
	}
	if w.Player.Heat != heat+attires[3].Notice {
		t.Fatalf("bespoke drew %d attention", w.Player.Heat-heat)
	}
}

func TestShabbyClothesDrawNoAttention(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Cash = 100000
	w.BuyAttire(w.Player.Dress + 1)
	w.BuyAttire(w.Player.Dress + 1)
	w.BuyAttire(w.Player.Dress + 1)
	w.Player.DressWear = Shabby - 1
	heat := w.Player.Heat
	w.DressDay()
	if w.Player.Heat != heat {
		t.Fatal("a ruined suit still made somebody curious")
	}
}

func TestPressingIsFreeAtALaundryOfYourOwn(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Cash = 100000
	w.BuyAttire(w.Player.Dress + 1)
	w.Player.DressWear = 50
	w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
	w.Player.Location = "laundry"
	if fee := w.PressFee("laundry"); fee != 0 {
		t.Fatalf("your own laundry charged $%d", fee)
	}
	cash := w.Player.Cash
	if err := w.Press("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash || w.DressCondition() != 95 {
		t.Fatalf("cash %d, condition %d", w.Player.Cash, w.DressCondition())
	}
}

func TestPressingCostsMoneyAtHomeAndNeverFullyRestoresARuinedSuit(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Cash = 100000
	w.BuyAttire(w.Player.Dress + 1)
	w.Player.DressWear = 10
	w.Player.Location = w.Player.Home
	cash := w.Player.Cash
	if err := w.Press(w.Player.Home); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-PressCost {
		t.Fatalf("pressing cost $%d", cash-w.Player.Cash)
	}
	if w.DressCondition() != 55 {
		t.Fatalf("condition %d", w.DressCondition())
	}
	w.Player.DressWear = 90
	w.Press(w.Player.Home)
	if w.DressCondition() != 100 {
		t.Fatalf("a suit still in one piece did not come back to 100: %d", w.DressCondition())
	}
}

func TestTheHighTablesJudgeYouAtTheDoor(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Cash = 100000
	w.Player.Location = "club"
	w.Player.Respect = 0
	high, _ := tableStake("high")
	small, _ := tableStake("small")
	if w.TableReadiness("club", high) == "" {
		t.Fatal("a man in working clothes with no name sat at the high tables")
	}
	if w.TableReadiness("club", small) != "" {
		t.Fatal("the small tables turned him away too:", w.TableReadiness("club", small))
	}
	w.Player.Location = "tailor"
	w.BuyAttire(w.Player.Dress + 1)
	w.BuyAttire(w.Player.Dress + 1)
	w.BuyAttire(w.Player.Dress + 1)
	w.Player.Location = "club"
	if w.TableReadiness("club", high) != "" {
		t.Fatal("bespoke did not open the room:", w.TableReadiness("club", high))
	}
	w.Player.DressWear = 20
	if w.TableReadiness("club", high) == "" {
		t.Fatal("a ruined suit still opened the room")
	}
}

func TestClothesDoNotSurviveTheirOwner(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Cash = 100000
	w.BuyAttire(w.Player.Dress + 1)
	w.Player.Alive = false
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	w = next
	if w.Player.Dress != 0 || w.Standing() != 0 {
		t.Fatalf("the next arrival inherited a suit: tier %d", w.Player.Dress)
	}
}

func TestAnOlderSaveIsWearingItsSuitInGoodOrder(t *testing.T) {
	t.Parallel()
	w := dressed(t)
	w.Player.Dress, w.Player.DressWear = 2, 0
	w.MigrateLivingWorld()
	if w.DressCondition() != 100 {
		t.Fatalf("migration left it at %d", w.DressCondition())
	}
}
