package store

import (
	"blackledger/core"
	"encoding/json"
	"path/filepath"
	"testing"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	s, e := Open(filepath.Join(t.TempDir(), "test.sqlite3"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { s.DB.Close() })
	return s
}
func TestIdempotentCommand(t *testing.T) {
	s := testStore(t)
	c := core.Command{RequestID: "duplicate-key", Revision: 0, Kind: "travel", Target: "bar"}
	a, e := s.Command(c)
	if e != nil {
		t.Fatal(e)
	}
	b, e := s.Command(c)
	if e != nil || string(a) != string(b) {
		t.Fatal("duplicate changed result")
	}
	w, _ := s.Read()
	if w.Revision != 1 {
		t.Fatal("double commit")
	}
}
func TestStaleRevision(t *testing.T) {
	s := testStore(t)
	_, _ = s.Command(core.Command{RequestID: "first-command", Revision: 0, Kind: "travel", Target: "bar"})
	_, e := s.Command(core.Command{RequestID: "second-command", Revision: 0, Kind: "courier", Target: "bar"})
	w, _ := s.Read()
	if e == nil || w.Player.Cash != 90 {
		t.Fatal("stale command changed cash")
	}
}
func TestInterruptedIncidentReload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.sqlite3")
	s, e := Open(path)
	if e != nil {
		t.Fatal(e)
	}
	e = s.Change(func(w *core.World) error { w.Player.Security = 1; w.Retaliation(); w.Advance(240); return nil })
	if e != nil {
		t.Fatal(e)
	}
	before, _ := s.Read()
	s.DB.Close()
	s, e = Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer s.DB.Close()
	after, _ := s.Read()
	if before.Event.ID != after.Event.ID || after.Minute != 720 {
		t.Fatal("incident rerolled on restart")
	}
}
func TestRejectedPayment(t *testing.T) {
	s := testStore(t)
	_, e := s.Command(core.Command{RequestID: "unaffordable", Revision: 0, Kind: "security", Target: "room"})
	w, _ := s.Read()
	if e == nil || w.Player.Cash != 90 || w.Revision != 0 {
		t.Fatal("rejected command not atomic")
	}
}

func TestLegacyEstatePurchaseMigration(t *testing.T) {
	w := core.New(27)
	w.Version = 1
	w.Player.Home = "estate"
	w.Player.Location = "estate"
	bytes, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	migrated, err := decode(string(bytes))
	if err != nil {
		t.Fatal(err)
	}
	if !migrated.Own("estate") || migrated.Player.BestHome != 2 || migrated.Version != 2 {
		t.Fatal("legacy purchase was not restored")
	}
	if migrated.Player.Cash != w.Player.Cash || migrated.Revision != w.Revision || migrated.Minute != w.Minute {
		t.Fatal("migration charged or advanced the player")
	}
	w.Player.Alive = false
	bytes, _ = json.Marshal(w)
	migrated, err = decode(string(bytes))
	if err != nil {
		t.Fatal(err)
	}
	if migrated.Own("estate") {
		t.Fatal("migration invented ownership for a dead player")
	}
}
