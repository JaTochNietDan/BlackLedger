package core

import (
	"blackledger/billiards"
	"encoding/json"
	"math"
	"strings"
	"testing"
)

func poolFixture(t *testing.T) *World {
	t.Helper()
	w := New(7)
	w.Player.Location = PoolPlace
	w.Player.Cash = 1000
	w.Event = nil
	w.Plots = nil
	w.NextPressure = 0
	n := w.NPC("leo")
	if n == nil {
		t.Fatal("missing fixture opponent")
	}
	n.Dead = false
	n.Purse = 300
	n.Held = 0
	n.Location = PoolPlace
	n.Heading = ""
	n.Sets = 0
	n.Arrives = 0
	return w
}
func beginPool(t *testing.T) *World {
	t.Helper()
	w := poolFixture(t)
	if err := w.StartPool("leo", 40); err != nil {
		t.Fatal(err)
	}
	return w
}
func winningPoolPosition(w *World) {
	m := w.Pool.Match
	m.Breaking = false
	m.InHand = false
	m.HeadOnly = false
	m.Groups = [2]int{1, 2}
	for i := range m.Balls {
		b := &m.Balls[i]
		b.Pocketed = true
		switch b.ID {
		case 0:
			b.Pocketed = false
			b.Position = billiards.Vec{X: .34, Y: billiards.Length / 2, Z: billiards.Radius}
		case 8:
			b.Pocketed = false
			b.Position = billiards.Vec{X: .14, Y: billiards.Length / 2, Z: billiards.Radius}
		case 9:
			b.Pocketed = false
			b.Position = billiards.Vec{X: .9, Y: 2, Z: billiards.Radius}
		}
	}
}
func TestPoolStakeIsFundedByBothPlayers(t *testing.T) {
	w := beginPool(t)
	if w.Player.Cash != 960 || w.NPC("leo").Purse != 260 || w.Pool.Escrow != 80 {
		t.Fatal("stake was not reserved", w.Player.Cash, w.NPC("leo").Purse, w.Pool.Escrow)
	}
	if total := w.Player.Cash + w.NPC("leo").Purse + w.Pool.Escrow; total != 1300 {
		t.Fatal("money was created or lost", total)
	}
	before, _ := json.Marshal(w)
	if err := w.StartPool("leo", 40); err == nil {
		t.Fatal("started a second funded match")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("rejected second stake changed state")
	}
}
func TestPoolPhysicalWinSettlesOnceAndOnlyProfitIsEarned(t *testing.T) {
	w := beginPool(t)
	winningPoolPosition(w)
	earned := w.Player.Earned
	if err := w.PlayPoolShot(billiards.Shot{Angle: math.Pi, Speed: .6}, billiards.Call{Ball: 8, Pocket: 4}); err != nil {
		t.Fatal(err)
	}
	if !w.Pool.Settled || w.Pool.Escrow != 0 || w.Player.Cash != 1040 || w.NPC("leo").Purse != 260 || w.Player.Earned != earned+40 {
		t.Fatal("incorrect settlement", w.Pool, w.Player.Cash, w.Player.Earned)
	}
	if w.Pool.Replay == "" {
		t.Fatal("winning shot has no replay")
	}
	logCount := len(w.History)
	for i := 0; i < 4; i++ {
		w.ReconcilePool()
		w.settlePool()
	}
	if w.Player.Cash != 1040 || len(w.History) != logCount {
		t.Fatal("duplicate payout")
	}
	if err := w.ConcedePool(); err == nil {
		t.Fatal("settled winner could concede payout away")
	}
}
func TestPoolConcessionPaysOpponentAndCanBeRepeatedOnlyAsNoOpReconciliation(t *testing.T) {
	w := beginPool(t)
	if err := w.ConcedePool(); err != nil {
		t.Fatal(err)
	}
	w.ReconcilePool()
	if w.Player.Cash != 960 || w.NPC("leo").Purse != 340 || w.Pool.Escrow != 0 || !w.Pool.Settled {
		t.Fatal("opponent did not receive actual stake")
	}
	if err := w.ConcedePool(); err == nil {
		t.Fatal("second concession accepted")
	}
}
func TestPoolReadinessRejectsUnfundedUnavailableOrWrongAddress(t *testing.T) {
	for _, change := range []func(*World){
		func(w *World) { w.Player.Cash = 20 }, func(w *World) { w.NPC("leo").Purse = 20 }, func(w *World) { w.NPC("leo").Dead = true }, func(w *World) { w.NPC("leo").Location = "bar" }, func(w *World) { w.Player.Location = "bar" }, func(w *World) { w.NPC("leo").Held = w.Minute + 60 }, func(w *World) { w.Seated = PoolPlace }, func(w *World) { w.Properties[PoolPlace].Condition = 0 },
	} {
		w := poolFixture(t)
		change(w)
		before, _ := json.Marshal(w)
		if err := w.StartPool("leo", 40); err == nil {
			t.Fatal("accepted invalid match")
		}
		after, _ := json.Marshal(w)
		if string(before) != string(after) {
			t.Fatal("rejected stake changed world")
		}
	}
	for _, stake := range []int{-1, 0, 9, 501, 1 << 30} {
		w := poolFixture(t)
		if err := w.StartPool("leo", stake); err == nil {
			t.Fatal("accepted stake outside limits", stake)
		}
	}
}
func TestPoolSaveCloneKeepsStakeAndMatchIndependent(t *testing.T) {
	w := beginPool(t)
	copy := w.Clone()
	if copy.Pool.Escrow != 80 || copy.Pool.Match == nil {
		t.Fatal("clone lost match or escrow")
	}
	if err := copy.PlacePoolCue(.6, .4); err != nil {
		t.Fatal(err)
	}
	if !w.Pool.Match.InHand {
		t.Fatal("cue placement changed original save")
	}
	if err := copy.ConcedePool(); err != nil {
		t.Fatal(err)
	}
	if w.NPC("leo").Purse != 260 || w.Player.Cash != 960 || w.Pool.Settled {
		t.Fatal("settlement changed original save")
	}
}
func TestPoolFailedShotDoesNotChangeStakeOrMatch(t *testing.T) {
	w := beginPool(t)
	if err := w.PlacePoolCue(.6, .4); err != nil {
		t.Fatal(err)
	}
	before, _ := json.Marshal(w)
	if err := w.PlayPoolShot(billiards.Shot{Speed: math.NaN()}, billiards.Call{}); err == nil {
		t.Fatal("accepted non-finite shot")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("failed shot changed saved state")
	}
}
func TestPoolDepartureAndDeathDoNotAbandonEscrow(t *testing.T) {
	for _, what := range []string{"player-leaves", "opponent-leaves", "opponent-dies", "player-dies"} {
		t.Run(what, func(t *testing.T) {
			w := beginPool(t)
			switch what {
			case "player-leaves":
				w.Player.Location = "bar"
			case "opponent-leaves":
				w.NPC("leo").Location = "bar"
			case "opponent-dies":
				w.NPC("leo").Dead = true
			case "player-dies":
				w.DieOf("test", "isolated fixture")
			}
			w.Advance(0)
			if !w.Pool.Settled || w.Pool.Escrow != 0 {
				t.Fatal("interruption left escrow held", w.Pool)
			}
			if strings.HasPrefix(what, "player") {
				if w.NPC("leo").Purse != 340 {
					t.Fatal("opponent did not receive forfeit")
				}
			} else if w.Player.Cash != 1040 {
				t.Fatal("player did not receive forfeit")
			}
		})
	}
}
func TestPoolOpponentStaysForRackAndIsRetainedUntilPaid(t *testing.T) {
	w := beginPool(t)
	w.SetOut()
	if w.NPC("leo").Heading != "" {
		t.Fatal("routine sent opponent out mid-rack")
	}
	w.NPC("leo").Dead = true
	if !w.referenced("leo") {
		t.Fatal("unpaid opponent can be pruned")
	}
	w.ReconcilePool()
	if !w.Pool.Settled {
		t.Fatal("dead opponent did not settle")
	}
}
func TestPoolNewLifeNeverReceivesOldPrize(t *testing.T) {
	w := beginPool(t)
	w.Player.Alive = false
	baseline := w.Clone()
	baseline.ReconcilePool()
	baseline.Pool = nil
	command := Command{Kind: "new_life", Revision: w.Revision}
	got, err := Execute(w, command)
	if err != nil {
		t.Fatal(err)
	}
	want, err := Execute(baseline, command)
	if err != nil {
		t.Fatal(err)
	}
	if got.Player.Cash != want.Player.Cash || got.NPC("leo").Purse != want.NPC("leo").Purse || !got.Pool.Settled {
		t.Fatal("old escrow crossed into new life")
	}
}

func TestPoolFireVoidsTheContestWithoutCreatingProfit(t *testing.T) {
	w := beginPool(t)
	earned := w.Player.Earned
	w.igniteBuilding(PoolPlace)
	w.Advance(0)
	if !w.Pool.Settled || !w.Pool.Voided || w.Pool.Escrow != 0 || w.Player.Cash != 1000 || w.NPC("leo").Purse != 300 || w.Player.Earned != earned {
		t.Fatal("fire did not return original stakes")
	}
	if err := w.StartPool("leo", 40); err == nil {
		t.Fatal("started a rack in a burning hall")
	}
	w.ReconcilePool()
	if w.Player.Cash != 1000 {
		t.Fatal("fire refund repeated")
	}
}

func TestPoolCannotShootThroughAnUnreconciledFire(t *testing.T) {
	w := beginPool(t)
	winningPoolPosition(w)
	w.igniteBuilding(PoolPlace)
	before, _ := json.Marshal(w)
	if err := w.PlayPoolShot(billiards.Shot{Angle: math.Pi, Speed: .6}, billiards.Call{Ball: 8, Pocket: 4}); err == nil {
		t.Fatal("played on a burning table")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("rejected fire shot altered money")
	}
	if err := w.ConcedePool(); err != nil {
		t.Fatal(err)
	}
	if !w.Pool.Voided || w.Player.Cash != 1000 || w.NPC("leo").Purse != 300 {
		t.Fatal("leaving burning table forfeited instead of refunding")
	}
}
