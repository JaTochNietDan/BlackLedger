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
	Trace        []Step         `json:"trace,omitempty"`
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
