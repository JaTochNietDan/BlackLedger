package store

import (
	"blackledger/billiards"
	"blackledger/core"
	"bytes"
	"math"
	"path/filepath"
	"sync"
	"testing"
)

func TestBilliardsCommandReceiptsSurviveRetriesAndReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pool-receipts.sqlite3")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { s.DB.Close() }()
	if err := s.Change(func(w *core.World) error {
		w.Player.Location = core.PoolPlace
		w.Player.Cash = 1000
		w.Event = nil
		w.Plots = nil
		w.NextPressure = 0
		n := w.NPC("leo")
		n.Location = core.PoolPlace
		n.Heading = ""
		n.Sets, n.Arrives, n.Held = 0, 0, 0
		n.Purse = 300
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	start := core.Command{RequestID: "pool-start-retry", Revision: 0, Kind: "pool_start", Target: "leo", Amount: 40}
	first, err := s.Command(start)
	if err != nil {
		t.Fatal(err)
	}
	// Concurrent retries contend for the same receipt, never a second stake.
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			again, e := s.Command(start)
			if e != nil || !bytes.Equal(first, again) {
				t.Error("start retry", e)
			}
		}()
	}
	wg.Wait()
	w, err := s.Read()
	if err != nil || w.Revision != 1 || w.Player.Cash != 960 || w.NPC("leo").Purse != 260 || w.Pool.Escrow != 80 {
		t.Fatal("stake duplicated", err)
	}
	// A late-rack position in this explicitly disposable save produces a real physical win.
	if err := s.Change(func(w *core.World) error {
		m := w.Pool.Match
		m.Breaking, m.InHand, m.HeadOnly = false, false, false
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
	shot := core.Command{RequestID: "pool-winning-shot", Revision: 1, Kind: "pool_shot", Pool: &core.PoolInput{Angle: math.Pi, Speed: .6, Ball: 8, Pocket: 4}}
	receipt, err := s.Command(shot)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.DB.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	again, err := s.Command(shot)
	if err != nil || !bytes.Equal(receipt, again) {
		t.Fatal("winning receipt changed after restart", err)
	}
	var before, after string
	if err = s.DB.QueryRow("SELECT state FROM campaign WHERE id=1").Scan(&before); err != nil {
		t.Fatal(err)
	}
	stale := shot
	stale.RequestID = "pool-stale-new-id"
	if _, err = s.Command(stale); err == nil {
		t.Fatal("accepted stale new request")
	}
	if err = s.DB.QueryRow("SELECT state FROM campaign WHERE id=1").Scan(&after); err != nil {
		t.Fatal(err)
	}
	if before != after {
		t.Fatal("stale shot changed persisted state")
	}
	var receipts int
	if err = s.DB.QueryRow("SELECT count(*) FROM receipts").Scan(&receipts); err != nil {
		t.Fatal(err)
	}
	if receipts != 2 {
		t.Fatal("unexpected receipts", receipts)
	}
	w, err = s.Read()
	if err != nil || w.Revision != 2 || w.Pool.Match.Shots != 1 || !w.Pool.Settled || w.Player.Cash != 1040 || w.NPC("leo").Purse != 260 {
		t.Fatal("shot or payment duplicated", err)
	}
	if _, err = billiards.DecodeReplay(w.Pool.Replay); err != nil {
		t.Fatal(err)
	}
}
