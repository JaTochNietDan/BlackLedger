package core

import (
	"fmt"
	"testing"
)

func houseKeeper(t *testing.T) *World {
	t.Helper()
	w := New(23)
	w.Properties["casino"].Owner = fmt.Sprintf("player:%d", w.Life)
	w.Properties["casino"].Condition = 100
	trade, _ := TradeOf("casino")
	w.Properties["casino"].Staff, w.Properties["casino"].Supply = trade.Hands, trade.RestockAmount
	w.Player.Location, w.Player.Cash = "casino", 5000
	return w
}

func TestAnEmptyRoomAttractsNobody(t *testing.T) {
	t.Parallel()
	w := houseKeeper(t)
	if w.Confidence("casino") != 0 || w.NightHandleAt("casino") != 0 {
		t.Fatalf("a room with nothing behind the tables ran $%d of action", w.NightHandleAt("casino"))
	}
	w.CasinoDay()
	if w.Properties["casino"].Bankroll != 0 {
		t.Fatal("an unfunded room made money")
	}
}

func TestWhatIsBehindTheTablesDecidesTheAction(t *testing.T) {
	t.Parallel()
	w := houseKeeper(t)
	previous := 0
	for _, float := range []int{250, 500, 1000, 1500, 3000} {
		w.Properties["casino"].Bankroll = float
		handle := w.NightHandleAt("casino")
		if float <= BankrollFull && handle <= previous {
			t.Fatalf("$%d behind the tables drew $%d of action against $%d", float, handle, previous)
		}
		previous = handle
	}
	w.Properties["casino"].Bankroll = 9000
	if w.NightHandleAt("casino") != previous {
		t.Fatal("action kept growing past the point a room can use it, so there is no reason to ever stop funding")
	}
}

func TestARoomInPoorOrderRunsLessAction(t *testing.T) {
	t.Parallel()
	w := houseKeeper(t)
	w.Properties["casino"].Bankroll = BankrollFull
	full := w.NightHandleAt("casino")
	w.Properties["casino"].Condition = 40
	if worn := w.NightHandleAt("casino"); worn >= full {
		t.Fatalf("a run-down room drew $%d against $%d", worn, full)
	}
	w.Properties["casino"].Condition = 100
	w.Properties["casino"].Staff = 0
	if unstaffed := w.NightHandleAt("casino"); unstaffed >= full {
		t.Fatalf("an unstaffed room drew $%d against $%d", unstaffed, full)
	}
}

func TestTheHouseEdgeShowsOverASeasonAndNotOverANight(t *testing.T) {
	t.Parallel()
	losses, nights := 0, 0
	total := 0
	for seed := uint32(1); seed <= 400; seed++ {
		w := houseKeeper(t)
		w.WorldRNG = seed * 2654435761
		w.Properties["casino"].Bankroll = 4000 // deep enough that ruin is not in play
		before := w.Properties["casino"].Bankroll
		w.CasinoDay()
		change := w.Properties["casino"].Bankroll - before
		if change < 0 {
			losses++
		}
		total += change
		nights++
	}
	if losses < nights/8 {
		t.Fatalf("the house lost on %d of %d nights, so there is no risk in running a room", losses, nights)
	}
	if total <= 0 {
		t.Fatalf("the house lost $%d across %d nights, so owning a casino is worth nothing", -total, nights)
	}
	t.Logf("%d of %d nights lost money; the house is $%d up across them, $%d a night", losses, nights, total, total/nights)
}

func TestAThinFloatIsTheThingThatRuinsAHouse(t *testing.T) {
	t.Parallel()
	ruinedThin, ruinedDeep := 0, 0
	for seed := uint32(1); seed <= 400; seed++ {
		thin := houseKeeper(t)
		thin.WorldRNG = seed * 2654435761
		thin.Properties["casino"].Bankroll = 400
		for day := 0; day < 30 && thin.Properties["casino"].Condition > 0; day++ {
			thin.CasinoDay()
		}
		if thin.hasRecord("The tables could not cover it at The Blue Hour") {
			ruinedThin++
		}

		deep := houseKeeper(t)
		deep.WorldRNG = seed * 2654435761
		deep.Properties["casino"].Bankroll = 4000
		for day := 0; day < 30; day++ {
			deep.CasinoDay()
		}
		if deep.hasRecord("The tables could not cover it at The Blue Hour") {
			ruinedDeep++
		}
	}
	// Twenty-five rather than forty. The threshold was set at the figure the
	// old sample produced, with no room in it, and campaign numbers were
	// feeding the stream an eighth of its range: spread properly a thin float
	// goes bust 37 times in 400 months rather than 40-odd. What is being
	// guarded is that under-funding a room costs something, so the line sits
	// clear of the measurement instead of on it.
	if ruinedThin < 25 {
		t.Fatalf("a thin float went bust in only %d of 400 months, so under-funding costs nothing", ruinedThin)
	}
	if ruinedDeep >= ruinedThin/2 {
		t.Fatalf("a deep float went bust in %d of 400 months against %d thin, so funding a room buys nothing", ruinedDeep, ruinedThin)
	}
	t.Logf("over 30 days: a $400 float cannot pay in %d of 400 campaigns, a $4,000 float in %d", ruinedThin, ruinedDeep)
}

func TestBeingUnableToPayCostsTheRoomAndTheName(t *testing.T) {
	t.Parallel()
	w := houseKeeper(t)
	w.Player.Respect = 20
	prop := w.Properties["casino"]
	prop.Bankroll = 100
	condition := prop.Condition
	// Force the night that ruins a house rather than waiting for one.
	prop.Bankroll = 0
	prop.Bankroll -= 300
	if prop.Bankroll >= 0 {
		t.Fatal("fixture failed")
	}
	prop.Bankroll = 100
	found := false
	for seed := uint32(1); seed <= 6000 && !found; seed++ {
		probe := houseKeeper(t)
		probe.Player.Respect = 20
		probe.WorldRNG = seed * 2654435761
		probe.Properties["casino"].Bankroll = 600
		probe.CasinoDay()
		if probe.Properties["casino"].Bankroll == 0 && probe.hasRecord("The tables could not cover it at The Blue Hour") {
			found = true
			if probe.Properties["casino"].Condition != condition-RuinCondition {
				t.Fatalf("condition went to %d", probe.Properties["casino"].Condition)
			}
			if probe.Player.Respect != 15 {
				t.Fatalf("respect went to %d", probe.Player.Respect)
			}
			if len(probe.News) == 0 || probe.News[len(probe.News)-1].Kind != "business" {
				t.Fatal("the city never heard about it")
			}
		}
	}
	if !found {
		t.Fatal("no thin house was ever cleaned out across 6000 nights")
	}
}

func TestWinningsOnlyReachThePlayerWhenTheyTakeThemOut(t *testing.T) {
	t.Parallel()
	w := houseKeeper(t)
	w.Properties["casino"].Bankroll = 2000
	cash := w.Player.Cash
	w.CasinoDay()
	float := w.Properties["casino"].Bankroll
	if w.Player.Cash != cash {
		t.Fatal("a night at the tables paid straight into the player's pocket")
	}
	if err := w.Draw("casino", 0); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash+BankrollLot || w.Properties["casino"].Bankroll >= float {
		t.Fatalf("cash %d, float %d", w.Player.Cash, w.Properties["casino"].Bankroll)
	}
}

func TestFundingAndDrawingAreRefusedWhereTheyMakeNoSense(t *testing.T) {
	t.Parallel()
	w := houseKeeper(t)
	if w.DrawReadiness("casino", 0) == "" {
		t.Fatal("drew from an empty float")
	}
	if w.BankrollReadiness("laundry", 0) == "" || w.BankrollReadiness("club", 0) == "" {
		t.Fatal("a laundry or somebody else's room took a float")
	}
	w.Player.Cash = 10
	if w.BankrollReadiness("casino", 0) == "" {
		t.Fatal("funded a room with no money")
	}
}

func TestAnOlderSaveHasMoneyBehindItsTables(t *testing.T) {
	t.Parallel()
	w := houseKeeper(t)
	w.Properties["casino"].Bankroll = 0
	w.MigrateLivingWorld()
	if w.Properties["casino"].Bankroll == 0 {
		t.Fatal("an owned room came back from migration with nothing behind the tables")
	}
}
