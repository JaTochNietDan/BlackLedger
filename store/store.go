package store

import (
	"blackledger/core"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

type Store struct{ DB *sql.DB }

func Open(path string) (*Store, error) {
	if path != ":memory:" {
		if e := os.MkdirAll(filepath.Dir(path), 0755); e != nil {
			return nil, e
		}
	}
	db, e := sql.Open("sqlite", path)
	if e != nil {
		return nil, e
	}
	db.SetMaxOpenConns(1)
	s := &Store{db}
	for _, q := range []string{"PRAGMA journal_mode=WAL", "PRAGMA busy_timeout=20000", "CREATE TABLE IF NOT EXISTS campaign (id INTEGER PRIMARY KEY, state TEXT NOT NULL)", "CREATE TABLE IF NOT EXISTS receipts (id TEXT PRIMARY KEY, result TEXT NOT NULL)"} {
		if _, e = db.Exec(q); e != nil {
			return nil, e
		}
	}
	b, _ := json.Marshal(core.New(27))
	_, e = db.Exec("INSERT OR IGNORE INTO campaign VALUES (1,?)", string(b))
	return s, e
}
func decode(data string) (*core.World, error) {
	var w core.World
	e := json.Unmarshal([]byte(data), &w)
	return &w, e
}
func (s *Store) Read() (*core.World, error) {
	var data string
	e := s.DB.QueryRow("SELECT state FROM campaign WHERE id=1").Scan(&data)
	if e != nil {
		return nil, e
	}
	return decode(data)
}
func (s *Store) Change(fn func(*core.World) error) error {
	tx, e := s.DB.Begin()
	if e != nil {
		return e
	}
	defer tx.Rollback()
	var data string
	if e = tx.QueryRow("SELECT state FROM campaign WHERE id=1").Scan(&data); e != nil {
		return e
	}
	w, e := decode(data)
	if e != nil {
		return e
	}
	if e = fn(w); e != nil {
		return e
	}
	b, e := json.Marshal(w)
	if e != nil {
		return e
	}
	if _, e = tx.Exec("UPDATE campaign SET state=? WHERE id=1", string(b)); e != nil {
		return e
	}
	return tx.Commit()
}
func (s *Store) Command(c core.Command) (json.RawMessage, error) {
	if len(c.RequestID) < 8 || len(c.RequestID) > 100 {
		return nil, fmt.Errorf("a request ID is required")
	}
	tx, e := s.DB.Begin()
	if e != nil {
		return nil, e
	}
	defer tx.Rollback()
	var receipt string
	e = tx.QueryRow("SELECT result FROM receipts WHERE id=?", c.RequestID).Scan(&receipt)
	if e == nil {
		return json.RawMessage(receipt), nil
	}
	if e != sql.ErrNoRows {
		return nil, e
	}
	var data string
	if e = tx.QueryRow("SELECT state FROM campaign WHERE id=1").Scan(&data); e != nil {
		return nil, e
	}
	w, e := decode(data)
	if e != nil {
		return nil, e
	}
	w, e = core.Execute(w, c)
	if e != nil {
		return nil, e
	}
	raw, e := json.Marshal(w)
	if e != nil {
		return nil, e
	}
	public, e := json.Marshal(w.Public())
	if e != nil {
		return nil, e
	}
	if _, e = tx.Exec("UPDATE campaign SET state=? WHERE id=1", string(raw)); e != nil {
		return nil, e
	}
	if _, e = tx.Exec("INSERT INTO receipts VALUES (?,?)", c.RequestID, string(public)); e != nil {
		return nil, e
	}
	if e = tx.Commit(); e != nil {
		return nil, e
	}
	return public, nil
}
