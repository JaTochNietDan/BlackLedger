package main

import (
	"blackledger/core"
	"blackledger/store"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func testApp(t *testing.T) *app {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "campaign.sqlite3"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.DB.Close() })
	return &app{s: s, port: "8791", client: http.DefaultClient}
}
func request(a *app, method, path, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	out := httptest.NewRecorder()
	a.ServeHTTP(out, r)
	return out
}
func TestHTTPCommandsSurviveDuplicateAndRejectStaleInput(t *testing.T) {
	t.Parallel()
	a := testApp(t)
	body := `{"kind":"travel","target":"bar","request_id":"first-journey","revision":0}`
	first := request(a, "POST", "/api/action", body)
	if first.Code != 200 {
		t.Fatal(first.Body.String())
	}
	duplicate := request(a, "POST", "/api/action", body)
	if duplicate.Code != 200 || duplicate.Body.String() != first.Body.String() {
		t.Fatal("duplicate did not return committed receipt")
	}
	stale := request(a, "POST", "/api/action", `{"kind":"courier","target":"bar","request_id":"stale-job","revision":0}`)
	if stale.Code != 409 {
		t.Fatal("stale action accepted")
	}
	w, _ := a.s.Read()
	if w.Player.Cash != 90 || w.Minute != 495 {
		t.Fatal("rejected/duplicate request advanced world")
	}
}
func TestHTTPPublicStateDoesNotDisclosePlot(t *testing.T) {
	t.Parallel()
	a := testApp(t)
	a.s.Change(func(w *core.World) error { w.Retaliation(); return nil })
	out := request(a, "GET", "/api/state", "")
	var public map[string]json.RawMessage
	if err := json.Unmarshal(out.Body.Bytes(), &public); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"plots", "rng", "offers"} {
		if _, ok := public[key]; ok {
			t.Fatal("private field disclosed", key)
		}
	}
	if out.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("state may be cached")
	}
}
func TestHTTPClosedSpeechAndForeignCommands(t *testing.T) {
	t.Parallel()
	a := testApp(t)
	if request(a, "POST", "/api/speech", `{"event":"already-closed"}`).Code != 409 {
		t.Fatal("closed conversation accepted")
	}
	r := httptest.NewRequest("POST", "/api/action", bytes.NewBufferString(`{}`))
	r.Header.Set("Origin", "https://unrelated.example")
	r.Header.Set("Content-Type", "application/json")
	out := httptest.NewRecorder()
	a.ServeHTTP(out, r)
	if out.Code != 403 {
		t.Fatal("foreign origin accepted")
	}
	if request(a, "POST", "/api/action", `{`).Code != 400 {
		t.Fatal("invalid JSON accepted")
	}
}
