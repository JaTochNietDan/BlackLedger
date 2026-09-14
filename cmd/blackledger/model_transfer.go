package main

import (
	"compress/gzip"
	"net/http"
	"strconv"
	"strings"
)

func acceptsModelGzip(value string) bool {
	for _, entry := range strings.Split(value, ",") {
		parts := strings.Split(entry, ";")
		if !strings.EqualFold(strings.TrimSpace(parts[0]), "gzip") {
			continue
		}
		quality := 1.0
		for _, p := range parts[1:] {
			key, v, ok := strings.Cut(strings.TrimSpace(p), "=")
			if ok && strings.EqualFold(key, "q") {
				n, err := strconv.ParseFloat(v, 64)
				if err != nil || n < 0 || n > 1 {
					return false
				}
				quality = n
			}
		}
		return quality > 0
	}
	return false
}

// GLB buffers compress losslessly. Keep range/error/conditional responses in
// their normal representation; never advertise the uncompressed byte length
// for a compressed transfer.
type modelGzipWriter struct {
	http.ResponseWriter
	gzip  *gzip.Writer
	wrote bool
}

func (w *modelGzipWriter) WriteHeader(status int) {
	if w.wrote {
		return
	}
	w.wrote = true
	if status == http.StatusOK {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
		w.Header().Del("Accept-Ranges")
		w.gzip, _ = gzip.NewWriterLevel(w.ResponseWriter, gzip.BestSpeed)
	}
	w.ResponseWriter.WriteHeader(status)
}
func (w *modelGzipWriter) Write(b []byte) (int, error) {
	if !w.wrote {
		w.WriteHeader(http.StatusOK)
	}
	if w.gzip != nil {
		return w.gzip.Write(b)
	}
	return w.ResponseWriter.Write(b)
}
func serveStatic(w http.ResponseWriter, r *http.Request, root string) {
	files := http.FileServer(http.Dir(root))
	if strings.HasSuffix(strings.ToLower(r.URL.Path), ".glb") {
		w.Header().Add("Vary", "Accept-Encoding")
		if r.Header.Get("Range") == "" && acceptsModelGzip(r.Header.Get("Accept-Encoding")) {
			compressed := &modelGzipWriter{ResponseWriter: w}
			files.ServeHTTP(compressed, r)
			if compressed.gzip != nil {
				_ = compressed.gzip.Close()
			}
			return
		}
	}
	files.ServeHTTP(w, r)
}
