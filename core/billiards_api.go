package core

import (
	"blackledger/billiards"
	"fmt"
	"strings"
)

// PoolInput is cue intent, never a physical outcome or an NPC's chosen input.
type PoolInput struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Angle  float64 `json:"angle"`
	Speed  float64 `json:"speed"`
	Top    float64 `json:"top"`
	Side   float64 `json:"side"`
	Ball   int     `json:"ball"`
	Pocket int     `json:"pocket"`
	Safety bool    `json:"safety"`
}

func (w *World) poolCommand(c Command) (string, int, error) {
	if c.Kind != "pool_start" && c.Target != "" && c.Target != PoolPlace {
		return "", 0, fmt.Errorf("that is not the billiards table")
	}
	if c.Pool != nil && c.Kind != "pool_place" && c.Kind != "pool_shot" && c.Kind != "pool_tournament_place" && c.Kind != "pool_tournament_shot" {
		return "", 0, fmt.Errorf("this command does not accept cue input")
	}
	if strings.HasPrefix(c.Kind, "pool_tournament_") {
		return w.poolTournamentCommand(c)
	}
	switch c.Kind {
	case "pool_start":
		stake := c.Amount
		if stake == 0 {
			stake = 20
		}
		return "Stake a billiards rack", 2, w.StartPool(c.Target, stake)
	case "pool_place":
		if c.Pool == nil {
			return "", 0, fmt.Errorf("a cue-ball position is required")
		}
		return "Place the cue ball", 0, w.PlacePoolCue(c.Pool.X, c.Pool.Y)
	case "pool_shot":
		if c.Pool == nil {
			return "", 0, fmt.Errorf("cue input is required")
		}
		p := c.Pool
		return "Play a billiards shot", 2, w.PlayPoolShot(billiards.Shot{Angle: p.Angle, Speed: p.Speed, Top: p.Top, Side: p.Side}, billiards.Call{Ball: p.Ball, Pocket: p.Pocket, Safety: p.Safety})
	case "pool_opponent":
		return "The other player shoots", 2, w.PlayPoolOpponent()
	case "pool_decide":
		return "Resolve the break", 0, w.DecidePoolBreak(c.Choice)
	case "pool_concede":
		return "Concede the billiards rack", 0, w.ConcedePool()
	case "pool_close":
		if w.Pool == nil || w.Pool.Life != w.Life {
			return "", 0, fmt.Errorf("there is no billiards table to leave")
		}
		if !w.Pool.Settled {
			return "", 0, fmt.Errorf("finish or concede the rack before leaving")
		}
		w.Pool = nil
		return "Leave the billiards table", 0, nil
	default:
		return "", 0, fmt.Errorf("unknown billiards command")
	}
}
func poolIntentView(intent billiards.Intent) map[string]any {
	s, c := intent.Shot, intent.Call
	return map[string]any{"angle": s.Angle, "speed": s.Speed, "top": s.Top, "side": s.Side, "ball": c.Ball, "pocket": c.Pocket, "safety": c.Safety}
}
func (w *World) PoolDescription() any {
	g := w.Pool
	if g == nil || g.Match == nil || g.Life != w.Life || g.Place != w.Player.Location {
		return nil
	}
	reason := ""
	if !g.Settled {
		if _, err := w.poolSession(); err != nil {
			reason = err.Error()
		}
	}
	return poolDescription(g, reason)
}
func poolDescription(g *PoolGame, reason string) any {
	m := g.Match
	balls := []map[string]any{}
	legal := []int{}
	for _, b := range m.Balls {
		pocket := -1
		if b.Pocketed {
			pocket = b.Pocket
		}
		balls = append(balls, map[string]any{"id": b.ID, "position": [3]float64{b.Position.X, b.Position.Y, b.Position.Z}, "rotation": b.Orientation, "pocket": pocket})
		if !b.Pocketed && m.LegalFirst(m.Turn, b.ID) {
			legal = append(legal, b.ID)
		}
	}
	choices := []map[string]string{}
	labels := map[string]string{"accept": "Accept the table", "spot-eight": "Spot the eight", "hand-behind": "Take cue ball behind the head string", "rerack-self": "Re-rack and break", "rerack-opponent": "Have the other player break again"}
	for _, choice := range m.BreakChoices() {
		choices = append(choices, map[string]string{"id": choice, "label": labels[choice]})
	}
	var stroke any
	if s := g.LastStroke; s != nil {
		item := map[string]any{"shooter": s.Shooter, "intent": poolIntentView(s.Intent), "decision": s.Decision}
		if s.Placement != nil {
			item["placement"] = [3]float64{s.Placement.X, s.Placement.Y, s.Placement.Z}
		}
		stroke = item
	}
	return map[string]any{"place": g.Place, "opponent": g.Opponent, "opponent_name": g.OpponentName, "stake": g.Stake, "pot": g.Escrow, "settled": g.Settled, "voided": g.Voided, "turn": m.Turn, "groups": m.Groups, "breaking": m.Breaking, "ball_in_hand": m.InHand, "behind_head_string": m.HeadOnly, "winner": m.Winner, "shots": m.Shots, "balls": balls, "legal_balls": legal, "break_choices": choices, "outcome": m.Last.Reason, "foul": m.Last.Foul, "replay": g.Replay, "stroke": stroke, "unavailable": reason, "width": billiards.Width, "length": billiards.Length, "radius": billiards.Radius}
}
func (w *World) PoolOpponents() []map[string]any {
	out := []map[string]any{}
	if !w.Player.Alive || w.Player.Location != PoolPlace {
		return out
	}
	for _, n := range w.NPCs {
		if n.Dead || n.Location != PoolPlace || w.Travelling(&n) || n.Held > w.Minute {
			continue
		}
		maximum := min(PoolMaxStake, w.Pockets(&n))
		reason := w.PoolReadiness(n.ID, PoolMinStake)
		out = append(out, map[string]any{"id": n.ID, "name": n.Name, "max_stake": maximum, "unavailable": reason})
	}
	return out
}
