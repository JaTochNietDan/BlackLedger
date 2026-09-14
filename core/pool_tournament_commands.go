package core

import (
	"blackledger/billiards"
	"fmt"
)

func (w *World) tournamentPlayerMatch(index int) (*billiards.BracketMatch, int, error) {
	t := w.PoolTournament
	if t == nil || t.Settled || t.Bracket == nil || index < 0 || index >= len(t.Bracket.Matches) {
		return nil, 0, fmt.Errorf("that tournament game is not active")
	}
	if t.Life != w.Life || !w.Player.Alive || w.Held() || w.Event != nil || w.Player.Location != PoolPlace || w.poolUnplayable() || t.Bracket.Withdrawn[t.PlayerID] {
		return nil, 0, fmt.Errorf("you are not available for this tournament game")
	}
	if w.Seated != "" || (w.Game != nil && !w.Game.Over) || (w.Hand != nil && !w.Hand.Done) {
		return nil, 0, fmt.Errorf("get up from the other game first")
	}
	cell := &t.Bracket.Matches[index]
	if cell.Rack == nil || cell.Resolved || cell.Rack.Winner >= 0 {
		return nil, 0, fmt.Errorf("that tournament game is not playing")
	}
	seat := -1
	for i, id := range cell.Players {
		if id == t.PlayerID {
			seat = i
		} else {
			n := w.NPC(id)
			if n == nil || n.Dead || n.Held > w.Minute || n.Location != PoolPlace || w.Travelling(n) || t.Bracket.Withdrawn[id] {
				return nil, 0, fmt.Errorf("the opponent is unavailable")
			}
		}
	}
	if seat < 0 {
		return nil, 0, fmt.Errorf("you are not playing at that table")
	}
	return cell, seat, nil
}
func (w *World) PlayPoolTournamentShot(index int, shot billiards.Shot, call billiards.Call) error {
	cell, seat, err := w.tournamentPlayerMatch(index)
	if err != nil {
		return err
	}
	next := *cell.Rack
	next.Balls = append([]billiards.Ball{}, cell.Rack.Balls...)
	result, err := next.Play(billiards.New(), seat, shot, call)
	if err != nil {
		return err
	}
	replay, err := billiards.EncodeReplay(result.Physics)
	if err != nil {
		return err
	}
	t := w.PoolTournament
	if t.Replays == nil {
		t.Replays = map[int]string{}
	}
	if t.Strokes == nil {
		t.Strokes = map[int]PoolStroke{}
	}
	t.Replays[index] = replay
	t.Strokes[index] = PoolStroke{Shooter: seat, Intent: billiards.Intent{Shot: shot, Call: call}}
	cell.Rack = &next
	w.ReconcilePoolTournament()
	return nil
}
func (w *World) poolTournamentCommand(c Command) (string, int, error) {
	t := w.PoolTournament
	if t == nil || t.Bracket == nil || t.Settled || t.Life != w.Life || !w.Player.Alive || w.Held() || w.Player.Location != PoolPlace {
		return "", 0, fmt.Errorf("there is no current tournament to play")
	}
	if c.Kind == "pool_tournament_withdraw" {
		err := t.Bracket.Withdraw(t.PlayerID)
		if err != nil {
			return "", 0, err
		}
		w.ReconcilePoolTournament()
		return "Withdraw from the tournament", 0, nil
	}
	if c.PoolGame == nil {
		return "", 0, fmt.Errorf("a tournament game index is required")
	}
	index := *c.PoolGame
	if c.Kind == "pool_tournament_opponent" {
		if w.Event != nil {
			return "", 0, fmt.Errorf("finish the current situation before watching a game")
		}
		return "Watch a tournament stroke", 2, w.PlayPoolTournamentBot(index)
	}
	cell, seat, err := w.tournamentPlayerMatch(index)
	if err != nil {
		return "", 0, err
	}
	switch c.Kind {
	case "pool_tournament_place":
		if c.Pool == nil {
			return "", 0, fmt.Errorf("a cue-ball position is required")
		}
		return "Place the tournament cue ball", 0, cell.Rack.PlaceCue(billiards.New(), seat, c.Pool.X, c.Pool.Y)
	case "pool_tournament_decide":
		return "Resolve the tournament break", 0, cell.Rack.DecideBreak(billiards.New(), seat, c.Choice)
	case "pool_tournament_shot":
		if c.Pool == nil {
			return "", 0, fmt.Errorf("cue input is required")
		}
		p := c.Pool
		return "Play a tournament shot", 2, w.PlayPoolTournamentShot(index, billiards.Shot{Angle: p.Angle, Speed: p.Speed, Top: p.Top, Side: p.Side}, billiards.Call{Ball: p.Ball, Pocket: p.Pocket, Safety: p.Safety})
	}
	return "", 0, fmt.Errorf("unknown tournament command")
}

func (w *World) PoolTournamentDescription() any {
	t := w.PoolTournament
	if t == nil || t.Bracket == nil || t.Life != w.Life || w.Player.Location != PoolPlace {
		return nil
	}
	names := map[string]string{}
	for _, id := range t.Bracket.Entrants {
		if id == t.PlayerID {
			names[id] = w.Player.Name
		} else if n := w.NPC(id); n != nil {
			names[id] = n.Name
		}
	}
	games := []map[string]any{}
	for i, cell := range t.Bracket.Matches {
		var table any
		seat := -1
		for j, id := range cell.Players {
			if id == t.PlayerID {
				seat = j
			}
		}
		if cell.Rack != nil {
			reason := ""
			if !cell.Resolved && !t.Settled && (w.Event != nil || !w.Player.Alive || w.Held() || w.poolUnplayable()) {
				reason = "the tournament is interrupted"
			}
			if seat >= 0 && !cell.Resolved && !t.Settled {
				if _, _, err := w.tournamentPlayerMatch(i); err != nil {
					reason = err.Error()
				}
			}
			g := &PoolGame{Place: PoolPlace, Life: t.Life, Opponent: cell.Players[1], OpponentName: names[cell.Players[1]], Stake: t.Fee, Escrow: t.Escrow, Settled: cell.Resolved || t.Settled, Voided: t.Voided, Match: cell.Rack, Replay: t.Replays[i]}
			if stroke, ok := t.Strokes[i]; ok {
				g.LastStroke = &stroke
			}
			table = poolDescription(g, reason)
		}
		games = append(games, map[string]any{"index": i, "round": cell.Round, "table_number": cell.Table, "players": cell.Players, "player_seat": seat, "resolved": cell.Resolved, "winner": cell.Winner, "table": table})
	}
	return map[string]any{"fee": t.Fee, "pot": t.Escrow, "settled": t.Settled, "voided": t.Voided, "finished": t.Bracket.Finished, "winner": t.Bracket.Winner, "player_id": t.PlayerID, "withdrawn": t.Bracket.Withdrawn[t.PlayerID], "names": names, "games": games}
}
