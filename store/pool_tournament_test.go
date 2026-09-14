package store

import (
	"blackledger/core"
	"path/filepath"
	"testing"
)

func TestPoolTournamentFeesAndWholePrizeSurviveReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tournament.sqlite3")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if s != nil {
			s.DB.Close()
		}
	}()
	if err = s.Change(func(w *core.World) error {
		w.Player.Location = core.PoolPlace
		w.Player.Cash = 1000
		w.Event = nil
		for _, id := range []string{"leo", "mara", "elena"} {
			n := w.NPC(id)
			n.Location = core.PoolPlace
			n.Purse = 300
			n.Held = 0
			n.Heading = ""
			n.Sets = 0
			n.Arrives = 0
		}
		return w.StartPoolTournament([]string{"leo", "mara", "elena"}, 25)
	}); err != nil {
		t.Fatal(err)
	}
	reopen := func() {
		t.Helper()
		if err = s.DB.Close(); err != nil {
			t.Fatal(err)
		}
		s, err = Open(path)
		if err != nil {
			t.Fatal(err)
		}
	}
	reopen()
	w, err := s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != 975 || w.PoolTournament.Escrow != 100 || w.NPC("leo").Purse != 275 {
		t.Fatal("deposits changed after restart")
	}
	// Actual concession decisions resolve this settlement fixture; core tests also
	// cover a physical called-eight final through the same settlement function.
	if err = s.Change(func(w *core.World) error {
		b := w.PoolTournament.Bracket
		if err := b.Matches[0].Rack.Concede(0); err != nil {
			return err
		}
		if err := b.Matches[1].Rack.Concede(1); err != nil {
			return err
		}
		w.ReconcilePoolTournament()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	reopen()
	if err = s.Change(func(w *core.World) error {
		if err := w.PoolTournament.Bracket.Matches[2].Rack.Concede(1); err != nil {
			return err
		}
		w.ReconcilePoolTournament()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	reopen()
	if err = s.Change(func(w *core.World) error { w.ReconcilePoolTournament(); w.ReconcilePool(); return nil }); err != nil {
		t.Fatal(err)
	}
	w, err = s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if !w.PoolTournament.Settled || w.PoolTournament.Escrow != 0 || w.NPC("leo").Purse != 375 || w.Player.Cash != 975 {
		t.Fatal("prize duplicated or lost across restart")
	}
}
