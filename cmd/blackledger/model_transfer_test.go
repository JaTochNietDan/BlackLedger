package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLargeModelLosslessTransferAndHTTPVariants(t *testing.T) {
	root := t.TempDir()
	data := bytes.Repeat([]byte("glTF-model-buffer-"), 350000)
	if err := os.WriteFile(filepath.Join(root, "room.glb"), data, 0600); err != nil {
		t.Fatal(err)
	}
	fetch := func(encoding, rangeValue, modified string) *httptest.ResponseRecorder {
		r := httptest.NewRequest("GET", "/room.glb", nil)
		r.Header.Set("Accept-Encoding", encoding)
		r.Header.Set("Range", rangeValue)
		r.Header.Set("If-Modified-Since", modified)
		out := httptest.NewRecorder()
		serveStatic(out, r, root)
		return out
	}
	out := fetch("br, gzip", "", "")
	if out.Code != 200 || out.Header().Get("Content-Encoding") != "gzip" || out.Header().Get("Content-Length") != "" || out.Header().Get("Vary") != "Accept-Encoding" {
		t.Fatal("invalid compressed headers", out.Code, out.Header())
	}
	reader, err := gzip.NewReader(out.Body)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := io.ReadAll(reader)
	reader.Close()
	if err != nil || !bytes.Equal(data, decoded) {
		t.Fatal("model bytes changed", err)
	}
	for _, accept := range []string{"", "gzip;q=0", "gzip;q=invalid", "br"} {
		plain := fetch(accept, "", "")
		if plain.Header().Get("Content-Encoding") != "" || !bytes.Equal(plain.Body.Bytes(), data) {
			t.Fatal("unsupported encoding", accept)
		}
	}
	ranged := fetch("gzip", "bytes=10-29", "")
	if ranged.Code != 206 || ranged.Header().Get("Content-Encoding") != "" || !bytes.Equal(ranged.Body.Bytes(), data[10:30]) {
		t.Fatal("range corrupted")
	}
	cached := fetch("gzip", "", out.Header().Get("Last-Modified"))
	if cached.Code != 304 || cached.Body.Len() != 0 || cached.Header().Get("Content-Encoding") != "" {
		t.Fatal("conditional response corrupted")
	}
	missing := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/missing.glb", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	serveStatic(missing, r, root)
	if missing.Code != 404 || missing.Header().Get("Content-Encoding") != "" {
		t.Fatal("missing asset encoded as model")
	}
}
