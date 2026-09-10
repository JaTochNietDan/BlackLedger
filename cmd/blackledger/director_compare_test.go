package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"blackledger/core"
)

// A controlled comparison, run only when BLACK_LEDGER_COMPARE is set, because
// it needs a model and takes minutes. Shipping the money picture to the
// director says nothing about whether the model writes differently for a family
// that cannot pay its people. This asks.
//
// One fixed city, one fixed family driven into each of two states, and the same
// number of requests each way — with organization_money and without. Everything
// else is identical, including the seed, so the only difference between the two
// runs is whether the model was told how the family is placed.

func askTheModel(t *testing.T, system string, body []byte) (string, bool) {
	t.Helper()
	payload, _ := json.Marshal(map[string]any{
		"model": env("BLACK_LEDGER_MODEL", "qwen3:14b"), "stream": false, "think": false,
		"messages": []map[string]string{{"role": "system", "content": system}, {"role": "user", "content": string(body)}},
		"options":  map[string]any{"temperature": .8, "num_predict": 700},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "POST", env("BLACK_LEDGER_OLLAMA", "http://127.0.0.1:11435")+"/api/chat", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Logf("the model did not answer: %v", err)
		return "", false
	}
	defer res.Body.Close()
	var out struct {
		Message struct{ Content string } `json:"message"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Logf("could not read the answer: %v", err)
		return "", false
	}
	return out.Message.Content, true
}

// starving drives one family into the state named, on a fixed seed.
func starving(t *testing.T, broke bool) (*core.World, *core.Faction) {
	t.Helper()
	w := core.New(41)
	w.Player.Cash, w.Player.Respect, w.Player.Health = 3000, 60, 100
	var f *core.Faction
	for i := range w.Factions {
		if w.Factions[i].ID == "bellandi" {
			f = &w.Factions[i]
		}
	}
	if f == nil {
		t.Fatal("no bellandi")
	}
	if broke {
		for _, id := range w.FamilyHoldings("bellandi") {
			w.Properties[id].Owner = ""
		}
		f.Cash = 0
		for day := 0; day < 20; day++ {
			w.FamilyDay()
			w.PayTheCity()
		}
	}
	return w, f
}

func TestTheModelWritesDifferentlyForAFamilyThatCannotPay(t *testing.T) {
	if os.Getenv("BLACK_LEDGER_COMPARE") == "" {
		t.Skip("set BLACK_LEDGER_COMPARE=1 to run the model comparison")
	}
	runs := 8
	// Words a speaker would only reach for if they knew the family was in
	// trouble. Written out rather than derived from the prompt, so the check
	// cannot agree with the prompt by construction.
	pressed := []string{"struggl", "cannot pay", "can't pay", "unpaid", "short",
		"desperat", "broke", "owe", "wages", "payday", "money is tight",
		"lean", "hungry", "no money", "out of money"}

	measure := func(withMoney, broke bool) (mentions, total int) {
		w, f := starving(t, broke)
		operation := w.NextDirectorOperation()
		connection := directorConnection(w)
		beneficiaries := []string{""}
		for _, faction := range w.Factions {
			beneficiaries = append(beneficiaries, faction.ID)
		}
		for i := 0; i < runs; i++ {
			c := directorContext(w, operation, connection, "", beneficiaries, false)
			if !withMoney {
				delete(c, "organization_money")
			}
			b, _ := json.Marshal(c)
			answer, ok := askTheModel(t, prompt, b)
			if !ok {
				continue
			}
			total++
			low := strings.ToLower(answer)
			for _, word := range pressed {
				if strings.Contains(low, word) {
					mentions++
					break
				}
			}
		}
		_ = f
		return mentions, total
	}

	type row struct {
		label           string
		mentions, total int
	}
	rows := []row{}
	for _, c := range []struct {
		label            string
		withMoney, broke bool
	}{
		{"told, family cannot pay", true, true},
		{"not told, family cannot pay", false, true},
		{"told, family comfortable", true, false},
		{"not told, family comfortable", false, false},
	} {
		m, n := measure(c.withMoney, c.broke)
		rows = append(rows, row{c.label, m, n})
		t.Logf("%-30s %d of %d scenes reach for a word about money", c.label, m, n)
	}
	if rows[0].total == 0 {
		t.Skip("the model answered nothing, so this measured nothing")
	}
	fmt.Fprintln(os.Stderr, "comparison complete")
}
