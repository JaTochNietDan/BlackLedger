package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExplicitFrontendRootIsRespected(t *testing.T) {
	a := testApp(t)
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "index.html"), []byte("isolated-frontend"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BLACK_LEDGER_WEB", root)
	response := request(a, http.MethodGet, "/", "")
	if response.Code != 200 || !strings.Contains(response.Body.String(), "isolated-frontend") {
		t.Fatal("configured root ignored")
	}
	if request(a, http.MethodGet, "/art/missing.png", "").Code != 404 {
		t.Fatal("missing asset did not fail explicitly")
	}
}
