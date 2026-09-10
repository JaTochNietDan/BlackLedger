package store

import (
	"blackledger/core"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
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
	t.Parallel()
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
	t.Parallel()
	s := testStore(t)
	_, _ = s.Command(core.Command{RequestID: "first-command", Revision: 0, Kind: "travel", Target: "bar"})
	_, e := s.Command(core.Command{RequestID: "second-command", Revision: 0, Kind: "courier", Target: "bar"})
	w, _ := s.Read()
	if e == nil || w.Player.Cash != 90 {
		t.Fatal("stale command changed cash")
	}
}
func TestInterruptedIncidentReload(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	s := testStore(t)
	_, e := s.Command(core.Command{RequestID: "unaffordable", Revision: 0, Kind: "security", Target: "room"})
	w, _ := s.Read()
	if e == nil || w.Player.Cash != 90 || w.Revision != 0 {
		t.Fatal("rejected command not atomic")
	}
}

func TestLegacyEstatePurchaseMigration(t *testing.T) {
	t.Parallel()
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
	if !migrated.Own("estate") || migrated.Player.BestHome != 2 || migrated.Version != core.SaveVersion {
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

// Adding four businesses to the city crashed every save already at the current
// version: the repair that gives a campaign a record for a new address was
// gated behind a version bump, and adding a place does not bump the version.
// This loads a campaign that has never heard of an address and asks whether it
// comes back with one.
func TestASaveThatNeverHeardOfAnAddressGetsOneOnLoad(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "test.sqlite3")
	s, e := Open(path)
	if e != nil {
		t.Fatal(e)
	}
	// A campaign at the version this build writes, from which one address has
	// been struck out entirely — which is exactly what an old save looks like
	// the morning after a business is added to the city.
	missing := "haulage"
	e = s.Change(func(w *core.World) error {
		w.Version = core.SaveVersion
		delete(w.Properties, missing)
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	s.DB.Close()

	s, e = Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer s.DB.Close()
	w, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	prop := w.Properties[missing]
	if prop == nil {
		t.Fatalf("%s is still missing after a load, so reading the world would panic on it", missing)
	}
	if prop.Income != core.PlaceIncome[missing] {
		t.Errorf("%s came back earning %d against %d in a new city", missing, prop.Income, core.PlaceIncome[missing])
	}
	// The crash itself was here.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("reading a loaded campaign: %v", r)
		}
	}()
	if w.Public() == nil {
		t.Fatal("the world read back as nothing")
	}
}

// A save written before a field existed has nothing under its name, and Go
// reads nothing as a nil slice, which goes down the wire as `null` for the view
// to count and blank on. Loading fills them, so an old campaign is no more
// dangerous to open than a new one.
func TestAnOldSaveComesBackWithEveryListFilled(t *testing.T) {
	t.Parallel()
	w, err := decode(`{"version":1,"life":1,"seed":7,"properties":{},"npcs":[]}`)
	if err != nil {
		t.Fatal(err)
	}
	empty := []string{}
	findNil(reflect.ValueOf(w), "", &empty, map[uintptr]bool{})
	if len(empty) > 0 {
		t.Fatalf("an old save came back with %d lists that are nothing rather than empty: %s",
			len(empty), strings.Join(empty[:min(6, len(empty))], ", "))
	}
}

func findNil(v reflect.Value, path string, out *[]string, seen map[uintptr]bool) {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return
		}
		if v.Kind() == reflect.Pointer {
			if seen[v.Pointer()] {
				return
			}
			seen[v.Pointer()] = true
		}
		findNil(v.Elem(), path, out, seen)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).IsExported() {
				findNil(v.Field(i), path+"/"+v.Type().Field(i).Name, out, seen)
			}
		}
	case reflect.Slice:
		if v.IsNil() {
			*out = append(*out, path)
			return
		}
		for i := 0; i < v.Len() && i < 3; i++ {
			findNil(v.Index(i), path+"/*", out, seen)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			findNil(v.MapIndex(key), path+"/*", out, seen)
		}
	}
}
