package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestExhaustedModelBudgetDoesNotRetryOrQueue(t *testing.T) {
	for _, thinking := range []string{"0", "1"} {
		t.Run(thinking, func(t *testing.T) {
			a := testApp(t)
			t.Setenv("BLACK_LEDGER_DIRECTOR_THINK", thinking)
			var calls atomic.Int32
			model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				json.NewEncoder(w).Encode(map[string]any{"done_reason": "length", "message": map[string]string{"content": ""}})
			}))
			defer model.Close()
			t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
			before, _ := a.s.Read()
			err := a.generate(before)
			after, _ := a.s.Read()
			if err == nil || !strings.Contains(err.Error(), "response budget") || calls.Load() != 1 {
				t.Fatal("token exhaustion incorrectly retried", err, calls.Load())
			}
			if len(after.Offers) != 0 || after.Minute != before.Minute || after.Player.Cash != before.Player.Cash {
				t.Fatal("budget failure changed campaign")
			}
		})
	}
}
