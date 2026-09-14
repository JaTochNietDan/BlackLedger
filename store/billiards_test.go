package store

import (
	"blackledger/billiards"
	"blackledger/core"
	"math"
	"path/filepath"
	"testing"
)

func TestBilliardsEscrowAndPhysicalSettlementSurviveDatabaseReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pool.sqlite3")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Change(func(w *core.World) error {
		w.Player.Location = core.PoolPlace
		w.Player.Cash = 1000
		w.Event = nil
		n := w.NPC("leo")
		n.Location = core.PoolPlace
		n.Heading = ""
		n.Purse = 300
		n.Held = 0
		if err := w.StartPool("leo", 40); err != nil {
			return err
		}
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
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err = s.DB.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Pool == nil || loaded.Pool.Escrow != 80 || loaded.Player.Cash != 960 || loaded.NPC("leo").Purse != 260 {
		t.Fatal("held stake changed on reopen")
	}
	if err = s.Change(func(w *core.World) error {
		return w.PlayPoolShot(billiards.Shot{Angle: math.Pi, Speed: .6}, billiards.Call{Ball: 8, Pocket: 4})
	}); err != nil {
		t.Fatal(err)
	}
	if err = s.DB.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.DB.Close()
	if err = s.Change(func(w *core.World) error { w.ReconcilePool(); return nil }); err != nil {
		t.Fatal(err)
	}
	loaded, err = s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.Pool.Settled || loaded.Pool.Escrow != 0 || loaded.Player.Cash != 1040 || loaded.NPC("leo").Purse != 260 {
		t.Fatal("winner paid incorrectly across reopen")
	}
	replay, err := billiards.DecodeReplay(loaded.Pool.Replay)
	if err != nil || len(replay.Frames) < 2 {
		t.Fatal("saved replay lost", err)
	}
}
