package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// Test-only provider recorder. It retains final response content and request
// fingerprints, not the model's reasoning text. Cancellation and response bytes
// pass through to the normal generation path, including rejected attempts.
func captureLiveProvider(t *testing.T, target string) func() []map[string]any {
	t.Helper()
	var mu sync.Mutex
	var active sync.WaitGroup
	attempts := []map[string]any{}
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		active.Add(1)
		defer active.Done()
		start := time.Now()
		requestBody, err := io.ReadAll(io.LimitReader(r.Body, 2<<20))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		digest := sha256.Sum256(requestBody)
		row := map[string]any{"request_sha256": hex.EncodeToString(digest[:])}
		record := func() {
			row["seconds"] = time.Since(start).Seconds()
			mu.Lock()
			attempts = append(attempts, row)
			mu.Unlock()
		}
		req, err := http.NewRequestWithContext(r.Context(), r.Method, strings.TrimRight(target, "/")+r.URL.RequestURI(), bytes.NewReader(requestBody))
		if err != nil {
			row["error"] = err.Error()
			record()
			http.Error(w, err.Error(), 502)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			row["error"] = err.Error()
			record()
			http.Error(w, err.Error(), 502)
			return
		}
		defer response.Body.Close()
		body, readErr := io.ReadAll(io.LimitReader(response.Body, 2<<20))
		row["http_status"] = response.StatusCode
		if readErr != nil {
			row["read_error"] = readErr.Error()
		}
		var decoded map[string]any
		if json.Unmarshal(body, &decoded) == nil {
			if message, ok := decoded["message"].(map[string]any); ok {
				row["content"] = message["content"]
			}
			for _, key := range []string{"done_reason", "error", "eval_count", "prompt_eval_count", "eval_duration", "load_duration", "total_duration"} {
				if value, ok := decoded[key]; ok {
					row[key] = value
				}
			}
		}
		record()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		_, _ = w.Write(body)
	}))
	t.Cleanup(proxy.Close)
	t.Setenv("BLACK_LEDGER_OLLAMA", proxy.URL)
	return func() []map[string]any {
		active.Wait()
		mu.Lock()
		defer mu.Unlock()
		out := make([]map[string]any, len(attempts))
		copy(out, attempts)
		return out
	}
}

func TestLiveProviderRecorderPreservesFailuresAndOmitsReasoning(t *testing.T) {
	const output = `{"message":{"content":"final dialogue","thinking":"private reasoning"},"eval_count":12}`
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(503); _, _ = io.WriteString(w, output) }))
	defer target.Close()
	records := captureLiveProvider(t, target.URL)
	response, err := http.Post(env("BLACK_LEDGER_OLLAMA", ""), "application/json", strings.NewReader(`{"model":"test"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil || response.StatusCode != 503 || string(body) != output {
		t.Fatal("recorder changed provider response", err)
	}
	rows := records()
	if len(rows) != 1 || rows[0]["content"] != "final dialogue" || rows[0]["http_status"] != 503 {
		t.Fatal("failed provider attempt missing", rows)
	}
	encoded, _ := json.Marshal(rows)
	if strings.Contains(string(encoded), "private reasoning") || len(rows[0]["request_sha256"].(string)) != 64 {
		t.Fatal("invalid report provenance or leaked reasoning")
	}
}

func TestLiveProviderRecorderWaitsForCancelledAttempt(t *testing.T) {
	started := make(chan struct{})
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		close(started)
		<-r.Context().Done()
	}))
	defer target.Close()
	records := captureLiveProvider(t, target.URL)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "POST", env("BLACK_LEDGER_OLLAMA", ""), strings.NewReader(`{"model":"test"}`))
	finished := make(chan error, 1)
	go func() {
		response, err := http.DefaultClient.Do(req)
		if response != nil {
			response.Body.Close()
		}
		finished <- err
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("request did not reach provider")
	}
	cancel()
	if err := <-finished; err == nil {
		t.Fatal("cancelled request unexpectedly succeeded")
	}
	rows := records()
	if len(rows) != 1 || rows[0]["error"] == nil {
		t.Fatal("cancelled provider attempt missing", rows)
	}
}
