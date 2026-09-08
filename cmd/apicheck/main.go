// Command apicheck plays a campaign against a running server using only the
// public HTTP API. cmd/simulate exercises the rules in process; this exercises
// the stack the browser actually uses -- command validation, the save
// transaction, event gating, revision conflicts and request idempotency.
//
// It chooses only from actions the API advertises, so it cannot drive the game
// through affordances a player does not have.
//
// Never point this at a real campaign. Give the server a fresh save:
//
//	BLACK_LEDGER_DB=.runtime/api-check.sqlite3 BLACK_LEDGER_PORT=8862 go run ./cmd/blackledger &
//	go run ./cmd/apicheck -base http://127.0.0.1:8862 -steps 120
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Progression first: these change what the campaign can do next. Earning is a
// fallback, so the policy follows the game's own suggestion instead of standing
// at one counter running errands forever.
var progression = []string{"repair", "delegate", "acquire", "recruit", "contact",
	"move_home", "security", "crew_bonus"}

// Earning actions, taken only when no progression step and no suggested
// destination is available.
var earning = []string{"courier", "dockwork"}

// Actions the core accepts without standing at the target location.
var remote = []string{"expand"}

// Passing time is a fallback, considered only after travel, so a campaign
// cannot settle in the starting room paying rent.
var idle = []string{"rest", "wait"}

var choicePreference = []string{"approach:careful", "accept", "pay", "escape", "acknowledge", "leave", "decline"}

type action struct {
	ID       string `json:"id"`
	Disabled bool   `json:"disabled"`
	Reason   string `json:"reason"`
	Target   string `json:"target"`
	Cost     int    `json:"cost"`
}

type place struct {
	ID      string   `json:"id"`
	Owned   bool     `json:"owned"`
	Actions []action `json:"actions"`
}

type choice struct {
	ID       string `json:"id"`
	Disabled bool   `json:"disabled"`
}

type snapshot struct {
	Revision int `json:"revision"`
	Life     int `json:"life"`
	Minute   int `json:"minute"`
	Player   struct {
		Name     string `json:"name"`
		Cash     int    `json:"cash"`
		Health   int    `json:"health"`
		Respect  int    `json:"respect"`
		Heat     int    `json:"heat"`
		Location string `json:"location"`
		Alive    bool   `json:"alive"`
	} `json:"player"`
	Locations []place `json:"locations"`
	Event     *struct {
		ID      string   `json:"id"`
		Kind    string   `json:"kind"`
		Title   string   `json:"title"`
		Choices []choice `json:"choices"`
	} `json:"event"`
	Opportunity *struct {
		Target string `json:"target"`
	} `json:"opportunity"`
	Factions []struct {
		ID       string `json:"id"`
		Goodwill int    `json:"goodwill"`
	} `json:"factions"`
}

type command struct {
	RequestID string `json:"request_id"`
	Revision  int    `json:"revision"`
	Kind      string `json:"kind"`
	Target    string `json:"target,omitempty"`
	Event     string `json:"event,omitempty"`
	Choice    string `json:"choice,omitempty"`
}

type client struct {
	base string
	http *http.Client
}

func (c *client) state() (*snapshot, error) {
	res, err := c.http.Get(c.base + "/api/state")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var s snapshot
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

// act returns the resulting snapshot, or the HTTP status and body when the
// server refuses the command.
func (c *client) act(cmd command) (*snapshot, int, string) {
	body, _ := json.Marshal(cmd)
	res, err := c.http.Post(c.base+"/api/action", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, 0, err.Error()
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode != 200 {
		return nil, res.StatusCode, string(raw)
	}
	var s snapshot
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, res.StatusCode, err.Error()
	}
	return &s, 200, ""
}

func requestID(step int) string {
	return fmt.Sprintf("apicheck-%d-%d", time.Now().UnixNano(), step)
}

func here(s *snapshot) map[string]action {
	out := map[string]action{}
	for _, p := range s.Locations {
		if p.ID != s.Player.Location {
			continue
		}
		for _, a := range p.Actions {
			if !a.Disabled {
				out[a.ID] = a
			}
		}
	}
	return out
}

// Work actions only appear once the player is at a place, so the reachable
// destination list is all the API gives us to plan with.
func reachable(s *snapshot) []string {
	out := []string{}
	for _, p := range s.Locations {
		if p.ID == s.Player.Location {
			continue
		}
		for _, a := range p.Actions {
			if a.ID == "travel" && !a.Disabled {
				out = append(out, p.ID)
			}
		}
	}
	return out
}

func destination(s *snapshot, visited map[string]int) string {
	places := reachable(s)
	if len(places) == 0 {
		return ""
	}
	if s.Opportunity != nil {
		for _, p := range places {
			if p == s.Opportunity.Target {
				return p
			}
		}
	}
	best := places[0]
	for _, p := range places {
		if visited[p] < visited[best] || (visited[p] == visited[best] && p < best) {
			best = p
		}
	}
	return best
}

func pick(s *snapshot, visited map[string]int) (command, bool) {
	if s.Event != nil {
		open := map[string]bool{}
		for _, c := range s.Event.Choices {
			if !c.Disabled {
				open[c.ID] = true
			}
		}
		for _, want := range choicePreference {
			if open[want] {
				return command{Kind: "choice", Event: s.Event.ID, Choice: want}, true
			}
		}
		for _, c := range s.Event.Choices {
			if !c.Disabled {
				return command{Kind: "choice", Event: s.Event.ID, Choice: c.ID}, true
			}
		}
		return command{}, false
	}
	available := here(s)
	take := func(names []string) (command, bool) {
		for _, want := range names {
			if a, ok := available[want]; ok {
				target := a.Target
				if target == "" {
					target = s.Player.Location
				}
				return command{Kind: want, Target: target}, true
			}
		}
		return command{}, false
	}
	if cmd, ok := take(progression); ok {
		return cmd, true
	}
	// Expansion is offered against a district the player is not standing in.
	for _, p := range s.Locations {
		for _, a := range p.Actions {
			if a.Disabled {
				continue
			}
			for _, want := range remote {
				if a.ID == want {
					return command{Kind: want, Target: p.ID}, true
				}
			}
		}
	}
	// Follow the suggested opportunity before settling for errands.
	if s.Opportunity != nil && s.Opportunity.Target != s.Player.Location {
		for _, p := range reachable(s) {
			if p == s.Opportunity.Target {
				return command{Kind: "travel", Target: p}, true
			}
		}
	}
	if cmd, ok := take(earning); ok {
		return cmd, true
	}
	if to := destination(s, visited); to != "" {
		return command{Kind: "travel", Target: to}, true
	}
	for _, want := range idle {
		if a, ok := available[want]; ok {
			target := a.Target
			if target == "" {
				target = s.Player.Location
			}
			return command{Kind: want, Target: target}, true
		}
	}
	return command{}, false
}

type failure struct {
	Step    int    `json:"step"`
	Problem string `json:"problem"`
}

type report struct {
	Base        string         `json:"base"`
	StepsTaken  int            `json:"steps_taken"`
	Idempotency string         `json:"idempotency"`
	Conflict    string         `json:"stale_revision"`
	Final       any            `json:"final"`
	Failures    []failure      `json:"invariant_failures"`
	Kinds       map[string]int `json:"command_counts"`
}

func main() {
	base := flag.String("base", "http://127.0.0.1:8862", "server base URL")
	steps := flag.Int("steps", 80, "maximum commands")
	out := flag.String("report", "", "write the JSON report here")
	verbose := flag.Bool("verbose", false, "log each committed command")
	flag.Parse()

	c := &client{base: *base, http: &http.Client{Timeout: 120 * time.Second}}
	s, err := c.state()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read state:", err)
		os.Exit(2)
	}

	rep := report{Base: *base, Kinds: map[string]int{}, Idempotency: "not reached", Conflict: "not reached"}
	fail := func(step int, format string, args ...any) {
		rep.Failures = append(rep.Failures, failure{Step: step, Problem: fmt.Sprintf(format, args...)})
	}
	visited := map[string]int{}
	checkedReplay, checkedConflict := false, false

	for step := 1; step <= *steps; step++ {
		cmd, ok := pick(s, visited)
		if !s.Player.Alive {
			cmd, ok = command{Kind: "new_life"}, true
		}
		if !ok {
			fail(step, "no action available at %s", s.Player.Location)
			break
		}
		cmd.RequestID, cmd.Revision = requestID(step), s.Revision
		before := s
		next, status, body := c.act(cmd)
		if status != 200 {
			// Refusing any command while an unresolved event is pending is
			// correct gating rather than a defect.
			if status == 409 && before.Event != nil {
				if s, err = c.state(); err != nil {
					fail(step, "state unreadable: %v", err)
					break
				}
				continue
			}
			fail(step, "HTTP %d on %s: %s", status, cmd.Kind, body)
			if s, err = c.state(); err != nil {
				break
			}
			continue
		}
		if next.Revision < before.Revision {
			fail(step, "revision went backwards: %d -> %d", before.Revision, next.Revision)
		}
		if next.Minute < before.Minute {
			fail(step, "clock went backwards: %d -> %d", before.Minute, next.Minute)
		}
		if next.Player.Cash < 0 {
			fail(step, "cash is negative: %d", next.Player.Cash)
		}
		if next.Player.Health < 0 || next.Player.Health > 100 {
			fail(step, "health out of range: %d", next.Player.Health)
		}
		if next.Life < before.Life {
			fail(step, "life count went backwards: %d -> %d", before.Life, next.Life)
		}

		// Replaying a committed request must not apply it a second time.
		if !checkedReplay && cmd.Kind != "new_life" {
			replay, replayStatus, _ := c.act(cmd)
			if replayStatus == 200 && replay != nil {
				if replay.Revision != next.Revision || replay.Player.Cash != next.Player.Cash || replay.Minute != next.Minute {
					fail(step, "replaying one request_id changed state: revision %d->%d cash %d->%d minute %d->%d",
						next.Revision, replay.Revision, next.Player.Cash, replay.Player.Cash, next.Minute, replay.Minute)
				} else {
					rep.Idempotency = "replayed request returned the same committed state"
				}
				next = replay
			} else {
				rep.Idempotency = fmt.Sprintf("replayed request refused with HTTP %d", replayStatus)
			}
			checkedReplay = true
		}

		// A command carrying a superseded revision must be refused, not applied.
		if !checkedConflict && cmd.Kind != "new_life" && next.Event == nil {
			stale := command{Kind: "wait", Target: next.Player.Location, RequestID: requestID(step * 1000), Revision: before.Revision - 1}
			_, staleStatus, staleBody := c.act(stale)
			if staleStatus == 200 {
				fail(step, "a command with a superseded revision was accepted")
			} else {
				rep.Conflict = fmt.Sprintf("HTTP %d: %s", staleStatus, staleBody)
			}
			if s, err = c.state(); err == nil {
				next = s
			}
			checkedConflict = true
		}

		if cmd.Kind == "travel" {
			visited[cmd.Target]++
		}
		rep.Kinds[cmd.Kind]++
		rep.StepsTaken++
		if *verbose {
			fmt.Fprintf(os.Stderr, "%3d %-12s day%d %02d:%02d cash=%d respect=%d heat=%d hp=%d\n",
				step, cmd.Kind, next.Minute/1440+1, next.Minute%1440/60, next.Minute%60,
				next.Player.Cash, next.Player.Respect, next.Player.Heat, next.Player.Health)
		}
		s = next
	}

	final, _ := c.state()
	owned := []string{}
	for _, p := range final.Locations {
		if p.Owned {
			owned = append(owned, p.ID)
		}
	}
	rep.Final = map[string]any{"life": final.Life, "minute": final.Minute, "cash": final.Player.Cash,
		"respect": final.Player.Respect, "heat": final.Player.Heat, "alive": final.Player.Alive, "owned": owned}

	if *out != "" {
		b, _ := json.MarshalIndent(rep, "", " ")
		if err := os.WriteFile(*out, append(b, '\n'), 0600); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	fmt.Printf("%d commands %v\n", rep.StepsTaken, rep.Final)
	fmt.Println("idempotency:", rep.Idempotency)
	fmt.Println("stale revision:", rep.Conflict)
	if len(rep.Failures) > 0 {
		fmt.Printf("%d invariant failure(s):\n", len(rep.Failures))
		for _, f := range rep.Failures {
			fmt.Println("  -", f.Step, f.Problem)
		}
		os.Exit(1)
	}
	fmt.Println("no invariant failures")
}
