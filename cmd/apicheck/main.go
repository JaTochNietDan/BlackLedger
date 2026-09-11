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
	"regexp"
	"strings"
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

// Work that names the person it is done to. These cannot be listed by id, so
// the harness matches them by prefix.
var named = []string{"lend:", "lean:", "extend:", "forgive:", "bail:", "sign:", "share:", "smear:"}

// Actions the core accepts without standing at the target location.
var remote = []string{"expand"}

// Passing time is a fallback, considered only after travel, so a campaign
// cannot settle in the starting room paying rent.
var idle = []string{"rest", "wait"}

// Ventures are the risky, optional systems: they are never required to make
// progress, so a policy that only climbs the ladder never touches them. Mixed
// in deliberately, because an untried system is an unverified one.
var ventures = []string{
	"launder", "bribe", "rob", "sabotage", "move", "incite", "contract",
	"enquire:bellandi", "enquire:russo", "pact:bellandi", "pact:russo", "serve:bellandi", "serve:russo", "leave_service", "takeover",
	"deposit", "offshore_access", "withdraw",
	"hire", "restock", "remedy", "layoff", "still", "dismantle",
	"play:small", "play:high", "buy:moonshine", "buy:cigarettes",
	"sell:moonshine", "sell:cigarettes", "arms:weapon", "arms:armour",
	"operate:hard", "operate:clean", "operate:standard", "inspect", "investigate", "lie_low",
	"dress", "press", "bankroll", "draw", "car", "service",
	"fit:door", "fit:telephone", "fit:safe", "fit:cellar", "commission",
	"trip:rockridge", "trip:kingsport", "trip:halloway", "charge", "plant", "sitdown", "retain:commissioner", "retain:mayor", "retain:editor", "spike", "puff", "rob:crew", "sabotage:crew", "armoury", "stock_arms", "buy:arms", "sell:arms", "mug", "mug:crew", "hit", "stand", "order", "post", "unpost", "sit_out", "lawyer", "talk",
	// Signing somebody on and lending them money both name them, so the
	// harness cannot list those by id.
}

var choicePreference = []string{"approach:careful", "accept", "pay", "escape", "acknowledge", "listen", "leave", "decline"}

// Multi-step events need their own preference or the policy always takes the
// exit: "leave" is offered at every step of arranging a contract.
func eventChoice(s *snapshot, turn int) (string, bool) {
	open := []string{}
	for _, c := range s.Event.Choices {
		if !c.Disabled {
			open = append(open, c.ID)
		}
	}
	if len(open) == 0 {
		return "", false
	}
	committing := []string{}
	for _, id := range open {
		if strings.HasPrefix(id, "mark:") || strings.HasPrefix(id, "hire:") {
			committing = append(committing, id)
		}
	}
	// Follow an arrangement through often enough to exercise it, but not always.
	if len(committing) > 0 && turn%3 != 0 {
		return committing[turn%len(committing)], true
	}
	for _, want := range choicePreference {
		for _, id := range open {
			if id == want {
				return id, true
			}
		}
	}
	return open[0], true
}

type action struct {
	ID       string `json:"id"`
	Disabled bool   `json:"disabled"`
	Reason   string `json:"reason"`
	Target   string `json:"target"`
	Cost     int    `json:"cost"`
	// What the card carries besides its id. Asking somebody where to find
	// somebody else publishes the mark here as well as in the id, for a client
	// that sends the old one-name form — and a harness that drops it is not
	// sending what the browser sends. Without it this refused four times in a
	// run with "There is nobody to ask about".
	Choice string `json:"choice,omitempty"`
}

type place struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Room    string   `json:"room"`
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
	History []struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"history"`
	Editions []struct {
		Day      int    `json:"day"`
		Dateline string `json:"dateline"`
		Stories  []struct {
			Headline string `json:"headline"`
			Body     string `json:"body"`
		} `json:"stories"`
	} `json:"editions"`
}

// The prose is assembled at runtime from names, roles and prices, and no test
// sees the sentence that comes out. Every copy fault found so far was found by
// reading a real save's output: a splinter's name carrying its own lower-case
// article into the start of a sentence, a plural organization name taking a
// singular verb, a role with no article, the Sunday page under a Monday
// dateline. These are the patterns, checked against everything the player and
// the reader actually see.
var badCopy = []struct{ pattern, why string }{
	{"arms is ", "a plural good took a singular verb"},
	{"cigarettes is ", "a plural good took a singular verb"},
	{"They were Lieutenant", "a role was used without an article"},
	{"They were Soldier", "a role was used without an article"},
	{"They were  ", "a role was empty and left a hole"},
	{"come forward. Police", "the police line was printed twice"},
	// A label written for a button, spliced into a sentence: the contract
	// tiers read "A professional", and the ledger printed "The a professional
	// you paid $2501 to reach Elena Russo".
	{"The a ", "a label with its own article was given another one"},
	{"The an ", "a label with its own article was given another one"},
	{"the a ", "a label with its own article was given another one"},
	{"The someone ", "a label written for a button was spliced into a sentence"},
	{"The The ", "an article was added to a phrase that had one"},
	{" a a ", "an article was doubled"},
	// "5 crates of Crated arms for $1100." The market lists a good by a name
	// fit for a price board, and the ledger put it after a count of units.
	{"crates of Crated", "a good's price-board name was used after a unit count"},
	{"cases of Untaxed", "a good's price-board name was used after a unit count"},
	{"crates of Moonshine", "a good's price-board name was used after a unit count"},
	// Both found by reading a campaign's paper end to end a second time, and
	// both introduced by earlier fixes on this same list.
	{"a head of", "a titled office was given an article"},
	{"an head of", "a titled office was given an article"},
	{"twice such", "a frequency was used where a count belongs"},
	// "You know Mayor Ellis Crane now: They answered to nobody."
	{"a police commissioner", "a unique office took an indefinite article"},
	{"a mayor of", "a unique office took an indefinite article"},
	{"an editor of", "a unique office took an indefinite article"},
	{"three times such", "a frequency was used where a count belongs"},
}

// singularOne catches "Whatever was arranged for you happened 1 times to a
// locked door". The boundary matters: without it "11 people in here" was
// reported as a fault, and a check that cries wolf is worse than no check.
// pluralSubject catches "Franca Sabbatini's people has people asking where you
// sleep". It is anchored at a sentence start and refuses to cross a comma or
// the words that introduce a second party, because "Violence between Brenner
// Company and Franca Sabbatini's people has escalated" is correct English and
// a plain substring reported it as a fault four times in fifteen runs.
//
// Every verb listed is one that ONLY exists in the singular. Adding "moved"
// and "came" to widen the net reported "Nico Ward's people moved against Saint
// Agnes and were driven off", which is correct — a past tense is the same for
// both. That was the third time this scan cried wolf. A check that misses
// something is worth more than one that has to be ignored.
var pluralSubject = regexp.MustCompile(`(^|\. )[A-Z][A-Za-z']*(\s[A-Z][A-Za-z']*)*'s people (has|is|was|holds|controls|takes|wants|makes|thinks) `)

var singularOne = regexp.MustCompile(`(^|[^0-9])1 (times|days|people|stories|others|minutes|crates|men)\b`)

// readable checks everything the city has written down.
func readable(s *snapshot, fail func(int, string, ...any)) {
	check := func(where, text string) {
		if text != "" && strings.ToLower(text[:1]) == text[:1] && strings.ToUpper(text[:1]) != text[:1] {
			fail(0, "%s begins in lower case: %q", where, first(text, 60))
		}
		for _, bad := range badCopy {
			if strings.Contains(text, bad.pattern) {
				fail(0, "%s: %s (%q)", where, bad.why, first(text, 80))
			}
		}
		if pluralSubject.MatchString(text) {
			fail(0, "%s: a plural organization name took a singular verb (%q)", where, first(text, 80))
		}
		if singularOne.MatchString(text) {
			fail(0, "%s: a count of one took a plural noun (%q)", where, first(text, 80))
		}
	}
	for _, r := range s.History {
		check("a ledger record", r.Title)
		check("a ledger record", r.Text)
	}
	// What the room says about itself is prose too, assembled from a headcount
	// and the hour.
	for _, l := range s.Locations {
		check("the note in "+l.Name, l.Room)
	}
	days := map[int]int{}
	for _, e := range s.Editions {
		days[e.Day]++
		if days[e.Day] > 1 {
			fail(0, "the paper printed day %d twice, both %s", e.Day, e.Dateline)
		}
		for _, story := range e.Stories {
			check("a story", story.Body)
			if strings.Contains(story.Headline, "SUNDAY") && !strings.Contains(e.Dateline, "Sunday") {
				fail(0, "the Sunday page ran under %q", e.Dateline)
			}
		}
	}
}

func first(text string, n int) string {
	if len(text) <= n {
		return text
	}
	return text[:n] + "…"
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

// asCommand turns a card into the command a client would post for it, carrying
// whatever the card publishes beside its id.
func asCommand(a action, s *snapshot) command {
	return command{Kind: a.ID, Target: targetOr(a, s), Choice: a.Choice}
}

func targetOr(a action, s *snapshot) string {
	if a.Target != "" {
		return a.Target
	}
	return s.Player.Location
}

// Staying on your feet. Measured over fifteen runs of eighty commands: the
// harness reached life fourteen, and twelve of its thirteen deaths were the
// same line — "Caught on the street by somebody else's war". It travelled
// relentlessly looking for systems it had not tried, at any health, through
// any war, and died before it ever accumulated the money or the standing that
// half the game is gated behind. A dead harness verifies nothing, and a player
// who is hurt goes home rather than walking across town.
const (
	restBelow   = 55
	lieLowAbove = 85
)

func pick(s *snapshot, visited map[string]int, tried map[string]int, turn int) (command, bool) {
	if s.Event != nil {
		choice, ok := eventChoice(s, turn)
		if !ok {
			return command{}, false
		}
		return command{Kind: "choice", Event: s.Event.ID, Choice: choice}, true
	}
	available := here(s)
	// Before anything else, do not die. This is not timidity: it is the only
	// way the expensive half of the game is ever reached, and it costs nothing
	// in coverage because resting is itself a command.
	if s.Player.Health < restBelow {
		if a, ok := available["rest"]; ok {
			return command{Kind: "rest", Target: targetOr(a, s)}, true
		}
	}
	if s.Player.Heat > lieLowAbove {
		if a, ok := available["lie_low"]; ok {
			return command{Kind: "lie_low", Target: targetOr(a, s)}, true
		}
	}
	// Anything at all this run has not taken yet, before the written list.
	//
	// `ventures` below is a hand-written list of the risky optional systems,
	// and a hand-written list only ever covers what somebody thought of. The
	// line that reported "1 crates of arms" off a boat at the pier was found by
	// accident, because taking a lot off a boat is not on that list and never
	// was. So the first question each turn is the one the exploring policy in
	// the simulator asks: of everything this room is actually offering, what
	// has this run taken least? It needs no list and it cannot go stale.
	if turn%2 == 0 {
		want, fewest := "", 0
		for id := range available {
			if id == "new_life" || id == "rest" || id == "wait" || id == "lie_low" {
				continue
			}
			if n := tried[id]; want == "" || n < fewest {
				want, fewest = id, n
			}
		}
		if want != "" && fewest == 0 {
			return asCommand(available[want], s), true
		}
	}
	// Prefer a system this run has not exercised yet: an untried system is an
	// unverified one, and coverage is the point of the harness. Fall back to
	// rotating through the rest so behaviour is still varied.
	if turn%2 == 0 {
		// Start the scan at a different place each turn. Reading the coverage
		// report showed ten ventures that were offered, enabled, and never
		// taken in fifteen runs — including all three trips out of the city.
		// The scan began at the top of the list every time, so an untried
		// venture near the front won every turn and anything late in the list
		// was starved no matter how often the game offered it.
		start := turn % len(ventures)
		for offset := 0; offset < len(ventures); offset++ {
			want := ventures[(start+offset)%len(ventures)]
			if tried[want] > 0 {
				continue
			}
			if a, ok := available[want]; ok {
				return command{Kind: want, Target: targetOr(a, s)}, true
			}
		}
		for offset := 0; offset < len(ventures); offset++ {
			want := ventures[(turn/2+offset)%len(ventures)]
			if a, ok := available[want]; ok {
				return command{Kind: want, Target: targetOr(a, s)}, true
			}
		}
	}
	// Some work names the person it is done to, so the harness cannot list it
	// by id. Take the first available action carrying one of these prefixes,
	// which is how lending, collecting and bailing get exercised at all.
	if turn%3 == 0 {
		for _, prefix := range named {
			for id, a := range available {
				if strings.HasPrefix(id, prefix) {
					return command{Kind: id, Target: targetOr(a, s)}, true
				}
			}
		}
	}
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
	// What the game was showing when it refused. A 409 on an action the game
	// had just listed as available is the only way the timing-drift class has
	// ever been caught, and the report used to carry the message and nothing
	// else — the state that produced it was gone by the time anyone read it,
	// and finding the cause took three failed hypotheses and seven hundred
	// attempts. Now the refusal brings its own evidence.
	Refused *refusal `json:"refused,omitempty"`
}

// refusal is the action as the game offered it, and the world it was offered
// in.
type refusal struct {
	Kind     string `json:"kind"`
	Target   string `json:"target"`
	Offered  bool   `json:"offered_as_available"`
	Reason   string `json:"reason_shown,omitempty"`
	Minute   int    `json:"minute"`
	Location string `json:"location"`
	Cash     int    `json:"cash"`
	Health   int    `json:"health"`
	Heat     int    `json:"heat"`
	Respect  int    `json:"respect"`
	Standing string `json:"standing,omitempty"`
}

// asOffered describes the action the way the player saw it a moment earlier.
func asOffered(s *snapshot, cmd command) *refusal {
	r := &refusal{Kind: cmd.Kind, Target: cmd.Target, Minute: s.Minute,
		Location: s.Player.Location, Cash: s.Player.Cash, Health: s.Player.Health,
		Heat: s.Player.Heat, Respect: s.Player.Respect}
	for _, p := range s.Locations {
		if p.ID != s.Player.Location {
			continue
		}
		for _, a := range p.Actions {
			if a.ID == cmd.Kind {
				r.Offered = !a.Disabled
				r.Reason = a.Reason
			}
		}
	}
	// Standing with every organization, because that is what moved under the
	// sitdown and nothing in the report would have shown it.
	parts := []string{}
	for _, f := range s.Factions {
		parts = append(parts, fmt.Sprintf("%s %+d", f.ID, f.Goodwill))
	}
	r.Standing = strings.Join(parts, ", ")
	return r
}

// blocked is why a venture was never exercised. A coverage list that says only
// "never tried" cannot tell the difference between a system the harness chose
// not to reach and one the game never offered it — and those want opposite
// fixes. This separates them: a venture the player was never even shown is a
// reachability problem, one shown but always greyed out is a wealth or state
// problem, and one offered and enabled but not taken is the harness's fault.
type blocked struct {
	Venture string `json:"venture"`
	Why     string `json:"why"`
	// Reason is the game's own explanation, when it gave one.
	Reason string `json:"reason,omitempty"`
}

type report struct {
	Base        string         `json:"base"`
	StepsTaken  int            `json:"steps_taken"`
	Idempotency string         `json:"idempotency"`
	Conflict    string         `json:"stale_revision"`
	Final       any            `json:"final"`
	Failures    []failure      `json:"invariant_failures"`
	Kinds       map[string]int `json:"command_counts"`
	Untried     []string       `json:"never_tried,omitempty"`
	Blocked     []blocked      `json:"why_never_tried,omitempty"`
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
	// What the game put in front of the player, and whether it was greyed out.
	offered, enabled, greyed := map[string]bool{}, map[string]bool{}, map[string]string{}
	note := func(s *snapshot) {
		if s == nil {
			return
		}
		for _, p := range s.Locations {
			for _, a := range p.Actions {
				offered[a.ID] = true
				if a.Disabled {
					if _, seen := greyed[a.ID]; !seen {
						greyed[a.ID] = a.Reason
					}
				} else {
					enabled[a.ID] = true
				}
			}
		}
	}
	note(s)
	checkedReplay, checkedConflict := false, false

	for step := 1; step <= *steps; step++ {
		cmd, ok := pick(s, visited, rep.Kinds, step)
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
			if len(rep.Failures) > 0 {
				rep.Failures[len(rep.Failures)-1].Refused = asOffered(before, cmd)
			}
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
		note(next)
		s = next
	}

	final, _ := c.state()
	if final != nil {
		readable(final, fail)
	}
	owned := []string{}
	for _, p := range final.Locations {
		if p.Owned {
			owned = append(owned, p.ID)
		}
	}
	rep.Final = map[string]any{"life": final.Life, "minute": final.Minute, "cash": final.Player.Cash,
		"respect": final.Player.Respect, "heat": final.Player.Heat, "health": final.Player.Health,
		"alive": final.Player.Alive, "owned": owned}

	// Coverage is part of the result. A clean run that never tried a system has
	// not tested it, and saying so is the difference between evidence and noise.
	// This has to happen before the report is written: it was computed after,
	// so every JSON report ever produced carried an empty coverage list while
	// the terminal showed the real one.
	untried := []string{}
	for _, venture := range ventures {
		if rep.Kinds[venture] == 0 {
			untried = append(untried, venture)
		}
	}
	rep.Untried = untried
	for _, venture := range untried {
		switch {
		case !offered[venture]:
			rep.Blocked = append(rep.Blocked, blocked{Venture: venture, Why: "never appeared in any location's action list"})
		case !enabled[venture]:
			rep.Blocked = append(rep.Blocked, blocked{Venture: venture, Why: "offered but always out of reach", Reason: greyed[venture]})
		default:
			rep.Blocked = append(rep.Blocked, blocked{Venture: venture, Why: "offered and available, and the harness never took it"})
		}
	}

	if *out != "" {
		b, _ := json.MarshalIndent(rep, "", " ")
		if err := os.WriteFile(*out, append(b, '\n'), 0600); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	fmt.Printf("%d commands %v\n", rep.StepsTaken, rep.Final)
	fmt.Println("idempotency:", rep.Idempotency)
	fmt.Println("stale revision:", rep.Conflict)
	fmt.Printf("exercised %d kinds of command: %v\n", len(rep.Kinds), rep.Kinds)
	if len(untried) > 0 {
		fmt.Printf("never tried (%d): %v\n", len(untried), untried)
	}
	if len(rep.Failures) > 0 {
		fmt.Printf("%d invariant failure(s):\n", len(rep.Failures))
		for _, f := range rep.Failures {
			fmt.Println("  -", f.Step, f.Problem)
		}
		os.Exit(1)
	}
	fmt.Println("no invariant failures")
}
