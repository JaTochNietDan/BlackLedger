// Package sim exercises the production command boundary using public information only.
package sim

import (
	"blackledger/core"
	"encoding/json"
	"fmt"
)

type Place struct {
	ID        string        `json:"id"`
	Owned     bool          `json:"owned"`
	Locked    bool          `json:"locked"`
	Cost      int           `json:"cost"`
	Condition int           `json:"condition"`
	Income    int           `json:"income"`
	// How much longer this place is keeping its money somewhere else after
	// being robbed. The city knows it whether or not anybody is standing in the
	// room, which is what lets a policy walk to a different one instead.
	Shy int `json:"shy"`
	Actions   []core.Action `json:"actions"`
}
type Event struct {
	ID      string        `json:"id"`
	Kind    string        `json:"kind"`
	Choices []core.Choice `json:"choices"`
}

// View deliberately excludes plots, RNG, queued offers and private director memory.
type View struct {
	BusinessTruces map[string]int `json:"business_truces"`
	Revision       int            `json:"revision"`
	Minute         int            `json:"minute"`
	Player         core.Person    `json:"player"`
	Locations      []Place        `json:"locations"`
	Event          *Event         `json:"event"`
	District       int            `json:"district"`
}
type Step struct {
	Number  int          `json:"number"`
	Minute  int          `json:"minute"`
	Cash    int          `json:"cash"`
	Health  int          `json:"health"`
	Command core.Command `json:"command"`
}
type Report struct {
	Seed         uint32         `json:"seed"`
	Strategy     string         `json:"strategy"`
	Director     string         `json:"director"`
	Commands     int            `json:"commands"`
	Minutes      int            `json:"game_minutes"`
	Alive        bool           `json:"alive"`
	Cash         int            `json:"cash"`
	Respect      int            `json:"respect"`
	Milestones   map[string]int `json:"milestone_commands"`
	Actions      map[string]int `json:"action_counts"`
	Events       map[string]int `json:"event_counts"`
	Error        string         `json:"error,omitempty"`
	ReplayQueued int            `json:"replay_queued,omitempty"`
	// What the city did on its own while this run was happening.
	//
	// docs/LIVING_WORLD.md asks for exactly these and they did not exist: a
	// simulation reporting only deaths and cash cannot say whether the living
	// world is alive or whether every campaign plays out the same way. A system
	// is not finished because it runs; it is finished when a long simulation
	// shows it producing varied outcomes.
	World WorldMeasures `json:"world"`
	Trace []Step        `json:"trace,omitempty"`
}

// WorldMeasures counts the things the city did that the player did not do.
type WorldMeasures struct {
	// Conflicts that reached the state "war" during this run.
	WarsStarted int `json:"wars_started"`
	// Of those, the ones the player's organization was no party to.
	WarsElsewhere int `json:"wars_elsewhere"`
	// Premises whose owner changed. Between other organizations only: the
	// player buying a laundry is the player playing, not the city moving.
	HoldingsChangedHands int `json:"holdings_changed_hands"`
	// Organizations that came into existence and that stopped existing. Not
	// counting the player's own: naming your outfit is you playing.
	FactionsCreated   int `json:"factions_created"`
	FactionsDestroyed int `json:"factions_destroyed"`
	// Commands after which the player was worse off — health or a holding —
	// while a war they were no party to was running. This is the closest thing
	// to "affected by a conflict it had no part in" that can be counted
	// honestly: it does not prove the war caused the harm, only that the player
	// was taking damage while somebody else's war was on.
	HurtDuringOthersWar int `json:"hurt_during_others_war"`
}

// watcher remembers enough of the city to notice what changed between two
// commands. Nothing here reads the player's own actions: the question is what
// the city did while the player was doing something else.
type watcher struct {
	wars     map[string]bool
	owners   map[string]string
	factions map[string]bool
	health   int
	holdings int
}

func watch(w *core.World) *watcher {
	m := &watcher{wars: map[string]bool{}, owners: map[string]string{}, factions: map[string]bool{}}
	m.note(w)
	m.health = w.Player.Health
	m.holdings = owned(w)
	return m
}

func (m *watcher) note(w *core.World) {
	for _, c := range w.Conflicts {
		if c.State == "war" {
			m.wars[c.A+"|"+c.B] = true
		}
	}
	for _, l := range core.Locations {
		if prop := w.Properties[l.ID]; prop != nil {
			m.owners[l.ID] = prop.Owner
		}
	}
	for i := range w.Factions {
		m.factions[w.Factions[i].ID] = true
	}
}

func owned(w *core.World) int {
	n := 0
	for _, l := range core.Locations {
		if w.Own(l.ID) {
			n++
		}
	}
	return n
}

// changed folds one command's worth of city movement into the measures.
func (m *watcher) changed(w *core.World, into *WorldMeasures) {
	player := w.PlayerOrganizationID()
	elsewhere := false
	for _, c := range w.Conflicts {
		if c.State != "war" {
			continue
		}
		key := c.A + "|" + c.B
		if !m.wars[key] {
			into.WarsStarted++
			if c.A != player && c.B != player {
				into.WarsElsewhere++
			}
		}
		if c.A != player && c.B != player {
			elsewhere = true
		}
	}
	for _, l := range core.Locations {
		prop := w.Properties[l.ID]
		if prop == nil {
			continue
		}
		was, had := m.owners[l.ID]
		// Only movement between other organizations counts. The player taking
		// premises is the player playing.
		if had && was != prop.Owner && was != player && prop.Owner != player {
			into.HoldingsChangedHands++
		}
	}
	// The player naming their own organization is the player playing, not the
	// city making a new family. Counting it made this read as one new
	// organization per campaign for the strategies that form one, and almost
	// none for the strategies that do not — which says something about the
	// strategies and nothing at all about the city.
	live := map[string]bool{}
	for i := range w.Factions {
		id := w.Factions[i].ID
		live[id] = true
		if !m.factions[id] && id != player {
			into.FactionsCreated++
		}
	}
	for id := range m.factions {
		if !live[id] && id != player {
			into.FactionsDestroyed++
		}
	}
	if elsewhere && (w.Player.Health < m.health || owned(w) < m.holdings) {
		into.HurtDuringOthersWar++
	}
	m.wars = map[string]bool{}
	m.owners = map[string]string{}
	m.factions = map[string]bool{}
	m.note(w)
	m.health = w.Player.Health
	m.holdings = owned(w)
}

func Public(w *core.World) View {
	b, _ := json.Marshal(w.Public())
	var v View
	if err := json.Unmarshal(b, &v); err != nil {
		panic(err)
	}
	return v
}
func (v View) place(id string) Place {
	for _, p := range v.Locations {
		if p.ID == id {
			return p
		}
	}
	return Place{}
}
func (v View) action(target, kind string) (core.Command, bool) {
	for _, a := range v.place(target).Actions {
		if a.ID == kind && !a.Disabled {
			return core.Command{Kind: kind, Target: target, Revision: v.Revision}, true
		}
	}
	return core.Command{}, false
}
func (v View) at(target, kind string) (core.Command, bool) {
	if v.Player.Location != target {
		return v.action(target, "travel")
	}
	return v.action(target, kind)
}
func Choose(v View, strategy string) (core.Command, error) {
	if v.Event != nil {
		priorities := []string{"approach:careful", "accept", "pay", "escape", "acknowledge", "leave", "decline"}
		if strategy == "defiant" && v.Event.Kind == "business_pressure" {
			priorities = []string{"resist"}
		}
		if strategy == "diplomat" && v.Event.Kind == "audience" {
			priorities = []string{"business_truce", "leave"}
		}
		if strategy == "reckless" {
			priorities = []string{"approach:press", "accept", "resist", "defend", "leave", "decline"}
		}
		for _, id := range priorities {
			for _, c := range v.Event.Choices {
				if c.ID == id && !c.Disabled {
					return core.Command{Kind: "choice", Choice: id, Event: v.Event.ID, Revision: v.Revision}, nil
				}
			}
		}
		for _, c := range v.Event.Choices {
			if !c.Disabled {
				return core.Command{Kind: "choice", Choice: c.ID, Event: v.Event.ID, Revision: v.Revision}, nil
			}
		}
		return core.Command{}, fmt.Errorf("no affordable event choice")
	}
	if strategy == "reckless" {
		// One explicit early challenge, followed by staying home without security.
		// Respect is a public result of that challenge; no hidden schedule is read.
		target, kind := "club", "provoke"
		if v.Player.Respect > 0 {
			target, kind = v.Player.Home, "rest"
		}
		if c, ok := v.at(target, kind); ok {
			return c, nil
		}
	}

	if v.Player.Health < 65 {
		if c, ok := v.at(v.Player.Home, "rest"); ok {
			return c, nil
		}
	}
	// The thief builds the ordinary way first — a name, a crew, somebody on the
	// door — and only then starts taking tills. A policy that robs from the
	// first minute is dead inside six commands and measures nothing.
	if strategy == "thief" && v.Player.Respect >= 6 && v.Player.Security >= 1 {
		// A robbery that goes wrong is a beating, and two in a row is a death,
		// so this one waits until it is whole before trying another.
		if v.Player.Health < 85 {
			if c, ok := v.at(v.Player.Home, "rest"); ok {
				return c, nil
			}
		}
		// Takes tills, and nothing else while there is one to take. No other
		// policy here ever robs anything, so a rule about robbing the same
		// place twice could be added or deleted and every number in the report
		// would stay where it was.
		for _, l := range v.Locations {
			if l.Owned || l.Locked || l.Income <= 0 || l.Shy > 0 {
				continue
			}
			if c, ok := v.at(l.ID, "rob"); ok {
				return c, nil
			}
		}
	}
	if strategy == "worker" {
		if c, ok := v.at("docks", "dockwork"); ok {
			return c, nil
		}
	}
	if strategy == "diplomat" && v.Player.Cash >= 200 {
		for _, p := range v.Locations {
			if !p.Owned || p.Income <= 0 {
				continue
			}
			actor, venue := "bellandi", "club"
			if p.ID == "garage" || p.ID == "casino" {
				actor, venue = "russo", "garage"
			}
			if v.BusinessTruces[actor] <= v.Minute {
				if c, ok := v.at(venue, "audience"); ok {
					return c, nil
				}
			}
		}
	}
	if len(v.Player.Crew) > 0 && v.Player.Crew[0].Loyalty < 50 {
		if c, ok := v.action(v.Player.Location, "crew_bonus"); ok {
			return c, nil
		}
	}
	if c, ok := v.action(v.Player.Location, "delegate"); ok {
		return c, nil
	}
	target, kind, cost := "bar", "courier", 0
	switch {
	case v.Player.Respect < 6:
	case !v.place("laundry").Owned:
		target, kind, cost = "laundry", "acquire", v.place("laundry").Cost
	case len(v.Player.Crew) == 0:
		target, kind, cost = "bar", "recruit", 90
	case v.Player.Contacts < 2:
		target, kind, cost = "bar", "contact", 10
	case v.District == 0:
		target, kind, cost = "apartment", "expand", 100
	case v.Player.Home == "room":
		target, kind, cost = "apartment", "move_home", 180
	case v.Player.Security < 1:
		target, kind, cost = v.Player.Home, "security", 100
	case !v.place("garage").Owned:
		target, kind, cost = "garage", "acquire", v.place("garage").Cost
	case !v.place("casino").Owned:
		target, kind, cost = "casino", "acquire", v.place("casino").Cost
	default:
		target, kind = v.Player.Home, "rest"
	}
	for _, p := range v.Locations {
		if p.Owned && p.Income > 0 && p.Condition < 70 {
			target, kind, cost = p.ID, "repair", 50
			break
		}
	}
	if v.Player.Cash < cost {
		target, kind = "bar", "courier"
	}
	if kind == "expand" {
		if c, ok := v.action(target, kind); ok {
			return c, nil
		}
	}
	if c, ok := v.at(target, kind); ok {
		return c, nil
	}
	if c, ok := v.at("bar", "courier"); ok {
		return c, nil
	}
	return core.Command{}, fmt.Errorf("no policy action at %s", v.Player.Location)
}
func Run(seed uint32, strategy, director string, limit int, trace bool) Report {
	return RunRecorded(seed, strategy, director, limit, trace, nil)
}

// RunRecorded consumes each recorded proposal once; exhaustion falls back to authored play.
func RunRecorded(seed uint32, strategy, director string, limit int, trace bool, corpus []core.Proposal) Report {
	r := Report{Seed: seed, Strategy: strategy, Director: director, Milestones: map[string]int{}, Actions: map[string]int{}, Events: map[string]int{}}
	w := core.New(seed)
	eyes := watch(w)
	start := w.Minute
	nextOffer := start + 240
	cursor := 0
	for i := 0; i < limit && w.Player.Alive; i++ {
		// This is a deterministic test provider, not the real AI director. It uses the same validator.
		if (director == "fixture" || (director == "replay" && cursor < len(corpus))) && w.Minute >= nextOffer && w.Event == nil && len(w.Offers) == 0 {
			operation := w.NextDirectorOperation()
			p := core.Proposal{Title: "Simulation arrangement", Body: "I need help with a discreet neighborhood job.", Speaker: "mara", Operation: operation, Outcome: "The agreed job is complete.", Approaches: []core.Approach{{Method: "careful", Label: "Prepare carefully"}, {Method: "press", Label: "Push the schedule"}}}
			if director == "replay" {
				p = corpus[cursor]
			}
			scene, err := w.ValidateProposal(p)
			if err != nil {
				r.Error = err.Error()
				break
			}
			if director == "replay" {
				cursor++
			}
			w.Offers = append(w.Offers, core.Offer{Ready: w.Minute + 30, Event: scene})
			nextOffer = w.Minute + 240
		}
		v := Public(w)
		c, err := Choose(v, strategy)
		if err != nil {
			r.Error = err.Error()
			break
		}
		if trace {
			r.Trace = append(r.Trace, Step{i + 1, w.Minute, w.Player.Cash, w.Player.Health, c})
		}
		n, err := core.Execute(w, c)
		if err != nil {
			r.Error = fmt.Sprintf("command %d %s: %v", i+1, c.Kind, err)
			break
		}
		if n.Minute < w.Minute || n.Revision != w.Revision+1 || n.Player.Cash < 0 || n.Player.Health < 0 || n.Player.Health > 100 {
			r.Error = "command invariant failed"
			break
		}
		if v.Event != nil {
			r.Events[v.Event.Kind]++
		}
		r.Actions[c.Kind]++
		w = n
		eyes.changed(w, &r.World)
		r.Commands++
		marks := map[string]bool{"crew": len(w.Player.Crew) > 0, "laundry": w.Own("laundry"), "garage": w.Own("garage"), "casino": w.Own("casino"), "housing": w.Player.Home != "room", "security": w.Player.Security > 0, "death": !w.Player.Alive}
		for name, done := range marks {
			if done && r.Milestones[name] == 0 {
				r.Milestones[name] = r.Commands
			}
		}
	}
	r.ReplayQueued = cursor
	r.Minutes = w.Minute - start
	r.Cash = w.Player.Cash
	r.Respect = w.Player.Respect
	r.Alive = w.Player.Alive
	return r
}
