package main

import (
	"blackledger/core"
	"encoding/json"
	"os"
	"testing"
	"time"
)

// This explicit local-provider experiment is skipped by ordinary test runs.
// Passing proves delivery/structural invariants, not narrative correctness.
// The persisted prose still needs review against the accompanying facts.
func TestLiveDirectorScenarios(t *testing.T) {
	if os.Getenv("BLACK_LEDGER_LIVE_DIRECTOR_EVAL") != "1" {
		t.Skip("local-model evaluation requires explicit opt-in")
	}
	path := os.Getenv("BLACK_LEDGER_LIVE_EVAL_REPORT")
	if path == "" {
		t.Fatal("set BLACK_LEDGER_LIVE_EVAL_REPORT to a new report path")
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	rows := []map[string]any{}
	report := map[string]any{"model": env("BLACK_LEDGER_MODEL", "qwen3:14b"), "brief_mode": env("BLACK_LEDGER_DIRECTOR_BRIEF", "full"), "thinking": env("BLACK_LEDGER_DIRECTOR_THINK", "0") == "1", "started_utc": time.Now().UTC(), "method": "Selected scenarios use isolated temporary saves, actual local provider and production generation/validation path. No user campaign access. Prose requires manual assessment; test success is not semantic acceptance."}
	writeReport := func() {
		report["cases"] = rows
		b, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(b, '\n'), 0600); err != nil {
			t.Fatal(err)
		}
	}
	cases := []struct {
		name    string
		prepare func(*core.World)
	}{
		{"crew_after_declined_delivery", func(w *core.World) {
			w.Player.Crew = []core.Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 65}}
			w.Player.Respect = 6
			w.Arrangements = []core.ArrangementMemory{{ID: core.ID(), Life: w.Life, Speaker: "mara", Operation: "courier", Status: "declined", Title: "A delivery declined", Offer: "Carry a sealed envelope."}}
		}},
		{"leader_new_request", func(w *core.World) {
			w.Factions[1].Goodwill = 9
			w.Arrangements = []core.ArrangementMemory{{ID: core.ID(), Life: w.Life, Speaker: "mara", Operation: "courier", Status: "declined", Title: "A refused introduction"}}
		}},
		{"completed_tool_dispute", func(w *core.World) {
			w.District = 1
			w.Player.Location = "garage"
			w.Properties["garage"].Owner = "player:1"
			w.Arrangements = []core.ArrangementMemory{{ID: core.ID(), Life: w.Life, Speaker: "mara", Operation: "mediation", Location: "garage", Status: "completed", Title: "A tool shared", Offer: "Two mechanics at Russo Motor Works disagree about access to a tool. Help them agree on how to share it.", Result: "You completed the requested mediation without violence."}}
		}},
		{"new_person_after_death", func(w *core.World) {
			w.Properties["laundry"].Owner = "player:1"
			w.Arrangements = []core.ArrangementMemory{{ID: core.ID(), Life: w.Life, Speaker: "mara", Operation: "mediation", Status: "completed", Title: "Alex's old mediation", Result: "You completed the requested mediation without violence."}}
			w.Die("An attack ended the previous life.")
			next, err := core.Execute(w, core.Command{Kind: "new_life", Revision: w.Revision})
			if err != nil {
				t.Fatal(err)
			}
			*w = *next
		}},
	}
	writeReport()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := testApp(t)
			attempts := captureLiveProvider(t, env("BLACK_LEDGER_OLLAMA", "http://127.0.0.1:11435"))
			if err := a.s.Change(func(w *core.World) error { tc.prepare(w); return nil }); err != nil {
				t.Fatal(err)
			}
			before, err := a.s.Read()
			if err != nil {
				t.Fatal(err)
			}
			row := map[string]any{"name": tc.name, "status": "writing", "player": before.Player, "life": before.Life, "properties": before.Properties, "factions": before.Factions, "previous_arrangements": before.Arrangements, "operation": before.NextDirectorOperation(), "allowed_speakers": directorSpeakers(before, directorConnection(before))}
			rows = append(rows, row)
			writeReport()
			start := time.Now()
			err = a.generate(before)
			row["seconds"] = time.Since(start).Seconds()
			row["attempts"] = attempts()
			after, readErr := a.s.Read()
			if readErr != nil {
				t.Fatal(readErr)
			}
			row["offers"] = after.Offers
			if err != nil {
				row["status"], row["error"] = "failed", err.Error()
				t.Errorf("generation: %v", err)
			} else if len(after.Offers) != 1 || after.Minute != before.Minute || after.Player.Cash != before.Player.Cash {
				row["status"] = "invariant_failed"
				t.Error("expected one queued offer and unchanged gameplay clock/cash")
			} else {
				row["status"] = "ready"
			}
			writeReport()
			t.Logf("%s: %v, %.1fs", tc.name, row["status"], row["seconds"])
		})
	}
	report["finished_utc"] = time.Now().UTC()
	writeReport()
}
