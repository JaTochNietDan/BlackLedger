package core

import "testing"

// "When gambling you should be able to actually choose how much to gamble, not
// use set amounts, up to a max limit. The max limit should be defined by the
// casino owner dynamically, whether by you the owner by or by someone else who
// owns it."
//
// Two fixed lots were the whole of it: fifty dollars or three hundred, and a
// nickel or a dollar at a machine.

func punter(t *testing.T) *World {
	t.Helper()
	w := New(43)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 5000, 60, 100
	w.Player.Location = "casino"
	return w
}

func TestARoomHasAHouseLimitAndItIsTheHoldersToSet(t *testing.T) {
	t.Parallel()
	w := punter(t)
	// Nobody has said otherwise, so the room runs on what it is worth.
	base := w.TableLimit("casino")
	if base <= 0 {
		t.Fatalf("a gambling house with no limit at all: %d", base)
	}
	// A family holds it, and the limit is theirs, not the player's.
	w.Properties["casino"].Owner = "bellandi"
	if reason := w.LimitReadiness("casino"); reason == "" {
		t.Error("you can set the limit in somebody else's house")
	}
	w.Properties["casino"].Owner = "player:1"
	if reason := w.LimitReadiness("casino"); reason != "" {
		t.Fatalf("in your own house: %q", reason)
	}
	if err := w.SetLimit("casino", 750); err != nil {
		t.Fatal(err)
	}
	if got := w.TableLimit("casino"); got != 750 {
		t.Errorf("the limit reads %d", got)
	}
	// And it survives being read back through what the interface is told.
	if desc := w.LimitDescription("casino"); desc["limit"] != 750 {
		t.Errorf("the room says %v", desc["limit"])
	}
}

func TestYouChooseWhatYouPutDown(t *testing.T) {
	t.Parallel()
	w := punter(t)
	w.Properties["casino"].Owner = "bellandi"
	limit := w.TableLimit("casino")
	// Any amount up to the limit, not two of them.
	for _, amount := range []int{5, 37, 120, limit} {
		run := punter(t)
		run.Properties["casino"].Owner = "bellandi"
		cash := run.Player.Cash
		if err := run.Deal("casino", amount); err != nil {
			t.Fatalf("$%d refused: %v", amount, err)
		}
		if run.Player.Cash != cash-amount {
			t.Errorf("put down $%d and paid $%d", amount, cash-run.Player.Cash)
		}
		if run.Hand == nil || run.HandDescription()["stake"] != amount {
			t.Errorf("the table says the stake is %v", run.HandDescription()["stake"])
		}
	}
	// And nothing above it, nothing below a dollar, and nothing you do not have.
	for _, bad := range []int{0, -50, limit + 1, 99999} {
		if err := w.Deal("casino", bad); err == nil {
			t.Errorf("$%d was accepted", bad)
		}
	}
}

func TestTheWheelAndTheMachinesTakeWhatYouChooseToo(t *testing.T) {
	t.Parallel()
	w := punter(t)
	w.Properties["casino"].Owner = "bellandi"
	if f := w.faction("bellandi"); f != nil {
		f.Cash = 50000
	}
	cash := w.Player.Cash
	if err := w.PlayWheel("casino", "red", 175); err != nil {
		t.Fatal(err)
	}
	if w.Spin == nil || w.WheelDescription()["stake"] != 175 {
		t.Errorf("the wheel says %v", w.WheelDescription()["stake"])
	}
	if w.Player.Cash > cash-175 && !w.Spin.Won {
		t.Errorf("a losing spin of $175 cost $%d", cash-w.Player.Cash)
	}
	if err := w.PlayWheel("casino", "red", w.TableLimit("casino")+1); err == nil {
		t.Error("the wheel took more than the house allows")
	}

	machine := punter(t)
	machine.Player.Location = "bar"
	machine.Properties["bar"].Owner = "bellandi"
	before := machine.Player.Cash
	if err := machine.PullHandle("bar", 12); err != nil {
		t.Fatal(err)
	}
	if machine.Reels == nil || machine.MachineDescription()["stake"] != 12 {
		t.Errorf("the machine says %v", machine.MachineDescription()["stake"])
	}
	if machine.Player.Cash > before-12 && machine.Reels.Pays == 0 {
		t.Error("a losing pull cost nothing")
	}
	if err := machine.PullHandle("bar", machine.MachineLimit("bar")+1); err == nil {
		t.Error("a nickel machine took more than the house allows")
	}
}

// The wiring: the holder's own room offers it, and a typed figure comes through
// the command layer rather than being a lot the room picked.
func TestTheHolderSetsTheLimitThroughTheRoom(t *testing.T) {
	t.Parallel()
	w := punter(t)
	w.Properties["casino"].Owner = "player:1"
	var limit *Action
	for i, a := range w.Actions("casino") {
		if a.ID == "limit" {
			limit = &w.Actions("casino")[i]
			_ = i
		}
	}
	if limit == nil {
		t.Fatal("your own gambling house offers no limit to set")
	}
	if limit.Disabled {
		t.Fatalf("in your own house: %q", limit.Reason)
	}
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "limit", Target: "casino", Amount: 900})
	if err != nil {
		t.Fatal(err)
	}
	if got := next.TableLimit("casino"); got != 900 {
		t.Errorf("the room takes %d", got)
	}
	// A figure outside what a house can run is refused rather than clamped
	// quietly: a limit nobody agreed to is not a limit.
	if _, err := Execute(next, Command{Revision: next.Revision, Kind: "limit", Target: "casino", Amount: 1}); err == nil {
		t.Error("a house limit of a dollar was accepted")
	}
	// And a room somebody else holds is not yours to run.
	other := punter(t)
	other.Properties["casino"].Owner = "bellandi"
	if _, err := Execute(other, Command{Revision: other.Revision, Kind: "limit", Target: "casino", Amount: 900}); err == nil {
		t.Error("you set the limit in Bellandi's house")
	}
}

// And a typed stake goes through the command layer, which is what the interface
// actually presses.
func TestATypedStakeReachesTheTable(t *testing.T) {
	t.Parallel()
	w := punter(t)
	w.Properties["casino"].Owner = "bellandi"
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "play", Target: "casino", Amount: 143})
	if err != nil {
		t.Fatal(err)
	}
	if next.HandDescription()["stake"] != 143 {
		t.Errorf("the table took %v", next.HandDescription()["stake"])
	}
	if next.Player.Cash != w.Player.Cash-143 {
		t.Errorf("it cost $%d", w.Player.Cash-next.Player.Cash)
	}
	// Nothing named still deals a hand, because a button offered as available
	// has to work when it is pressed.
	bare, err := Execute(w, Command{Revision: w.Revision, Kind: "play", Target: "casino"})
	if err != nil {
		t.Fatalf("a hand with no figure named: %v", err)
	}
	if bare.HandDescription()["stake"] == 0 {
		t.Error("a hand with nothing on it")
	}
	t.Logf("nothing named puts down $%v", bare.HandDescription()["stake"])
}
