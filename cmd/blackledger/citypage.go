package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"blackledger/core"
)

// Letting the director set the city page.
//
// The page itself is composed from facts and is always there. This asks the
// local model to set one of those briefs in the register of a 1953 city paper,
// and core.AcceptPolish decides whether what comes back may be printed — the
// test of a rewrite is not whether it reads well but whether it added anything.
//
// Everything about this is arranged so that the page does not depend on it. The
// model may be absent, slow, or wrong; in all three cases the brief keeps the
// words the facts gave it and the reader never knows a request was made.

// polishPrompt is deliberately narrow. A model given room to be creative about
// a newspaper will invent an arrest.
const polishPrompt = `You are a rewrite man on a city newspaper in 1953.

You will be given one short paragraph of filler copy. Set it in the voice of a
working American city paper of that year: plain, unsentimental, a little dry,
short sentences.

Absolute rules:
- Add no fact that is not in the paragraph. No numbers that are not there. No
  names of people, streets, departments or organisations that are not there.
- Do not address the reader. Never write "you".
- Two or three sentences. Never more.
- Reply with the rewritten paragraph and nothing else. No preamble, no quotes,
  no explanation.`

// polish asks for one brief to be set better, in the background, at most one at
// a time. It shares the director's lock so the two cannot queue behind each
// other on a machine with one model loaded.
func (a *app) polish() bool {
	if !a.ai.TryLock() {
		return false
	}
	snapshot, e := a.s.Read()
	if e != nil {
		a.ai.Unlock()
		return false
	}
	brief, ok := snapshot.NextPolish()
	if !ok {
		a.ai.Unlock()
		return false
	}
	go func() {
		defer a.ai.Unlock()
		better, e := a.rewrite(brief.Body)
		if e != nil {
			// Not worth a status on the interface: the page is already correct,
			// and the reader has no idea a rewrite was ever attempted.
			log.Println("City page:", e)
			return
		}
		body, took := core.AcceptPolish(brief.Body, better)
		if !took {
			log.Printf("City page: rewrite refused for %q — %s", brief.Headline, core.PolishRefusal(brief.Body, better))
		}
		_ = a.s.Change(func(w *core.World) error {
			if w.ID != snapshot.ID || w.Life != snapshot.Life {
				return nil // a different campaign is loaded; the request is stale
			}
			w.SetPolish(brief.ID, body, took)
			return nil
		})
	}()
	return true
}

func (a *app) rewrite(original string) (string, error) {
	payload, _ := json.Marshal(map[string]any{
		"model": env("BLACK_LEDGER_MODEL", "qwen3:14b"), "stream": false, "think": false,
		"messages": []map[string]string{
			{"role": "system", "content": polishPrompt},
			{"role": "user", "content": original},
		},
		// Short and cool: this is a rewrite, not an invention.
		"options": map[string]any{"temperature": .55, "num_predict": 220},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "POST",
		env("BLACK_LEDGER_OLLAMA", "http://127.0.0.1:11435")+"/api/chat", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res, e := a.client.Do(req)
	if e != nil {
		return "", e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("rewrite unavailable: %s", res.Status)
	}
	var reply struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if e := json.NewDecoder(io.LimitReader(res.Body, 20000)).Decode(&reply); e != nil {
		return "", e
	}
	return reply.Message.Content, nil
}
