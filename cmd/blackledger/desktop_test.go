package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDesktopPaths(t *testing.T) {
	root := t.TempDir()
	executable := filepath.Join(root, "Game with spaces", "blackledger")
	config := filepath.Join(root, "user settings")
	if _, _, err := desktopPaths(executable, config); err == nil {
		t.Fatal("missing distribution accepted")
	}
	web := filepath.Join(filepath.Dir(executable), "dist")
	if err := os.MkdirAll(web, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(web, "index.html"), []byte("game"), 0644); err != nil {
		t.Fatal(err)
	}
	gotWeb, gotSave, err := desktopPaths(executable, config)
	if err != nil {
		t.Fatal(err)
	}
	if gotWeb != web || gotSave != filepath.Join(config, "BlackLedger", "campaign.sqlite3") {
		t.Fatalf("unexpected paths: %s %s", gotWeb, gotSave)
	}
	if _, err := os.Stat(config); !os.IsNotExist(err) {
		t.Fatal("resolving paths must not create or mutate a save")
	}
}
