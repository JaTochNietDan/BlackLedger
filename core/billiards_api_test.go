package core

import (
	"blackledger/billiards"
	"encoding/json"
	"math"
	"testing"
)

func poolExecute(t *testing.T, w *World, c Command) *World {
	t.Helper()
	c.Revision = w.Revision
	next, err := Execute(w, c)
	if err != nil {
		t.Fatal(c.Kind, err)
	}
	return next
}
func poolReject(t *testing.T, w *World, c Command) {
	t.Helper()
	before, _ := json.Marshal(w)
	c.Revision = w.Revision
	if _, err := Execute(w, c); err == nil {
		t.Fatal("accepted", c.Kind)
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("rejection mutated original")
	}
}
func TestPoolCommandsValidateIntentAndPreserveStakes(t *testing.T) {
	w := poolFixture(t)
	minute := w.Minute
	w = poolExecute(t, w, Command{Kind: "pool_start", Target: "leo", Amount: 40})
	if w.Revision != 1 || w.Minute != minute+2 || w.Pool.Escrow != 80 || w.Player.Cash != 960 {
		t.Fatal("start did not reserve exactly one stake and two minutes")
	}
	for _, c := range []Command{
		{Kind: "pool_place"}, {Kind: "pool_shot"}, {Kind: "pool_opponent"},
		{Kind: "pool_opponent", Pool: &PoolInput{Speed: 8}}, {Kind: "pool_close"},
		{Kind: "pool_place", Target: "casino", Pool: &PoolInput{X: .4, Y: .4}},
		{Kind: "pool_start", Target: "leo", Amount: 40}, {Kind: "pool_forged"},
		{Kind: "pool_shot", Pool: &PoolInput{Angle: math.Pi / 2, Speed: 7}},
		{Kind: "travel", Target: "bar"}, {Kind: "rest", Target: PoolPlace},
	} {
		poolReject(t, w, c)
	}
	w = poolExecute(t, w, Command{Kind: "face", Choice: "1"})
	if w.Pool.Settled || w.Pool.Escrow != 80 {
		t.Fatal("appearance change ended rack")
	}
	w = poolExecute(t, w, Command{Kind: "pool_place", Pool: &PoolInput{X: .4, Y: .4}})
	if w.Minute != minute+2 || w.Pool.Match.InHand {
		t.Fatal("placement advanced time or failed")
	}
	w = poolExecute(t, w, Command{Kind: "pool_shot", Pool: &PoolInput{Angle: math.Pi / 2, Speed: 7.5}})
	if w.Minute != minute+4 || w.Pool.Match.Shots != 1 || w.Pool.Replay == "" {
		t.Fatal("physical shot did not commit")
	}
	w = poolExecute(t, w, Command{Kind: "pool_concede"})
	if !w.Pool.Settled || w.Pool.Match.Winner != 1 || w.NPC("leo").Purse != 340 {
		t.Fatal("concession failed")
	}
	poolReject(t, w, Command{Kind: "pool_concede"})
	w = poolExecute(t, w, Command{Kind: "pool_close"})
	if w.Pool != nil {
		t.Fatal("closed rack remains")
	}
}
func TestPoolCommandPhysicalNPCWinAndEventPause(t *testing.T) {
	w := beginPool(t)
	winningPoolPosition(w)
	w.Pool.Match.Turn = 1
	w.Pool.Match.Groups = [2]int{2, 1}
	poolReject(t, w, Command{Kind: "pool_shot", Pool: &PoolInput{Angle: math.Pi, Speed: .6, Ball: 8, Pocket: 4}})
	w = poolExecute(t, w, Command{Kind: "pool_opponent"})
	if !w.Pool.Settled || w.Pool.Match.Winner != 1 || w.Pool.LastStroke.Shooter != 1 || w.Pool.Match.Shots != 1 {
		t.Fatal("NPC command did not physically win")
	}
	if _, err := billiards.DecodeReplay(w.Pool.Replay); err != nil {
		t.Fatal(err)
	}
	w = beginPool(t)
	w.Event = &Scene{ID: "pool-pause", Title: "An interruption"}
	poolReject(t, w, Command{Kind: "pool_place", Pool: &PoolInput{X: .4, Y: .4}})
	if w.PoolDescription().(map[string]any)["unavailable"] == "" {
		t.Fatal("paused table appears playable")
	}
	w = poolExecute(t, w, Command{Kind: "pool_concede"})
	w = poolExecute(t, w, Command{Kind: "pool_close"})
	if w.Event == nil || w.Event.ID != "pool-pause" {
		t.Fatal("conceding erased interruption")
	}
}
func TestPoolProjectionIsReadOnlyAndLocal(t *testing.T) {
	w := beginPool(t)
	winningPoolPosition(w)
	w = poolExecute(t, w, Command{Kind: "pool_shot", Pool: &PoolInput{Angle: math.Pi, Speed: .6, Ball: 8, Pocket: 4}})
	before, _ := json.Marshal(w)
	view := w.PoolDescription().(map[string]any)
	if view["winner"] != 0 || view["pot"] != 0 || view["settled"] != true || view["width"] != billiards.Width {
		t.Fatal("incorrect public outcome")
	}
	balls := view["balls"].([]map[string]any)
	for _, b := range balls {
		if b["id"] == 8 && b["pocket"] != 4 {
			t.Fatal("eight not in called pocket")
		}
		if b["id"] == 0 && b["pocket"] != -1 {
			t.Fatal("live cue marked pocketed")
		}
	}
	opponents := w.PoolOpponents()
	found := false
	for _, n := range opponents {
		if n["id"] == "leo" {
			found = true
		}
	}
	if !found {
		t.Fatal("local opponent omitted")
	}
	// DTO values must not alias mutable match data.
	balls[0]["position"] = [3]float64{99, 99, 99}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("reading/mutating projection changed world")
	}
	w.Player.Location = "bar"
	if w.PoolDescription() != nil || len(w.PoolOpponents()) != 0 {
		t.Fatal("remote table exposed")
	}
	w.Player.Location = PoolPlace
	w.Life++
	if w.PoolDescription() != nil {
		t.Fatal("old-life table exposed")
	}
}

func TestPoolCommandMillimetreDrawAndOffsetRejection(t *testing.T) {
	w := beginPool(t)
	winningPoolPosition(w)
	for i := range w.Pool.Match.Balls {
		b := &w.Pool.Match.Balls[i]
		if b.ID == 0 {
			b.Position = billiards.Vec{X: .4, Y: 1, Z: billiards.Radius}
		}
		if b.ID == 8 {
			b.Position = billiards.Vec{X: .4 + 2*billiards.Radius + .015, Y: 1, Z: billiards.Radius}
		}
	}
	poolReject(t, w, Command{Kind: "pool_shot", Pool: &PoolInput{Speed: 1.2, Top: .02, Safety: true}})
	w = poolExecute(t, w, Command{Kind: "pool_shot", Pool: &PoolInput{Speed: 1.2, Top: -.012, Safety: true}})
	if w.Pool.LastStroke.Intent.Shot.Top != -.012 {
		t.Fatal("saved offset changed units")
	}
	for _, b := range w.Pool.Match.Balls {
		if b.ID == 0 && b.Position.X >= .4 {
			t.Fatal("12mm draw did not reverse the cue ball", b.Position)
		}
	}
	if _, err := billiards.DecodeReplay(w.Pool.Replay); err != nil {
		t.Fatal(err)
	}
}
