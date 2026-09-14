package core

import (
	"blackledger/billiards"
	"fmt"
	"hash/fnv"
)

type PoolTournament struct {
	HouseCutPercent int                `json:"house_cut_percent"`
	HouseCutPaid    int                `json:"house_cut_paid"`
	PrizePaid       int                `json:"prize_paid"`
	NextStrokes     map[int]int        `json:"next_strokes"`
	Replays         map[int]string     `json:"replays"`
	Strokes         map[int]PoolStroke `json:"strokes"`
	Life            int                `json:"life"`
	PlayerID        string             `json:"player_id"`
	Fee             int                `json:"fee"`
	Escrow          int                `json:"escrow"`
	Deposits        map[string]int     `json:"deposits"`
	Bracket         *billiards.Bracket `json:"bracket"`
	Settled         bool               `json:"settled"`
	Voided          bool               `json:"voided"`
}

// Shared funding entry point; scheduled admission calls this after checking its window.
// All participants and all money are validated before any wallet is touched.
func (w *World) StartPoolTournament(npcs []string, fee int) error {
	return w.startPoolTournament(npcs, fee, 0, true)
}
func (w *World) startPoolTournament(npcs []string, fee, cut int, enter bool) error {
	count := len(npcs)
	if enter {
		count++
	}
	if count != 4 && count != 8 {
		return fmt.Errorf("a hall tournament needs four or eight entrants")
	}
	if fee < PoolMinStake || fee > PoolMaxStake || cut < 0 || cut > 50 {
		return fmt.Errorf("entry must be $10–$500 and the house cut 0–50 percent")
	}
	if w.PoolTournament != nil && !w.PoolTournament.Settled {
		return fmt.Errorf("a tournament is already underway")
	}
	if !w.Player.Alive || w.Held() || w.Event != nil || w.Player.Location != PoolPlace || w.poolUnplayable() {
		return fmt.Errorf("the hall or host is unavailable")
	}
	if (w.Pool != nil && !w.Pool.Settled) || w.Seated != "" || (w.Game != nil && !w.Game.Over) || (w.Hand != nil && !w.Hand.Done) {
		return fmt.Errorf("finish the other game first")
	}
	if enter && w.Player.Cash < fee {
		return fmt.Errorf("you cannot cover your entry fee")
	}
	for _, id := range npcs {
		n := w.NPC(id)
		if n == nil || n.Dead || n.Held > w.Minute || n.Location != PoolPlace || w.Travelling(n) || w.Pockets(n) < fee {
			return fmt.Errorf("every entrant must be here and able to pay")
		}
	}
	player := fmt.Sprintf("player:%d", w.Life)
	entrants := append([]string{}, npcs...)
	if enter {
		entrants = append([]string{player}, npcs...)
	}
	bracket, err := billiards.NewBracket(entrants)
	if err != nil {
		return err
	}
	if enter {
		if err = w.Pay(fee); err != nil {
			return err
		}
	}
	deposits := map[string]int{}
	if enter {
		deposits[player] = fee
	}
	for _, id := range npcs {
		n := w.NPC(id)
		n.Purse -= fee
		n.Heading, n.Errand, n.Sets, n.Arrives = "", "", 0, 0
		deposits[id] = fee
	}
	w.PoolTournament = &PoolTournament{HouseCutPercent: cut, Life: w.Life, PlayerID: player, Fee: fee, Escrow: fee * len(entrants), Deposits: deposits, Bracket: bracket, Replays: map[int]string{}, Strokes: map[int]PoolStroke{}}
	w.schedulePoolTournament()
	w.Log("Entry money on the baize", fmt.Sprintf("%d entrants each pay $%d. The winner receives $%d after the posted %d%% house cut. Entrants forfeit their fee by leaving.", len(entrants), fee, fee*len(entrants)-fee*len(entrants)*cut/100, cut), "personal")
	return nil
}
func (w *World) tournamentParticipant(id string) bool {
	t := w.PoolTournament
	if t == nil || t.Settled || t.Bracket == nil || t.Bracket.Finished || t.Bracket.Withdrawn[id] || t.Deposits[id] == 0 {
		return false
	}
	for _, cell := range t.Bracket.Matches {
		if cell.Resolved && (cell.Players[0] == id || cell.Players[1] == id) && cell.Winner != id {
			return false
		}
	}
	return true
}
func (w *World) ReconcilePoolTournament() {
	t := w.PoolTournament
	if t == nil || t.Settled || t.Bracket == nil {
		return
	}
	// A completed rack is reconciled before a later departure can rewrite it.
	t.Bracket.Advance()
	if !t.Bracket.Finished && w.poolUnplayable() {
		w.refundPoolTournament()
		return
	}
	if !t.Bracket.Finished {
		missing := []string{}
		for id := range t.Deposits {
			if t.Bracket.Withdrawn[id] {
				continue
			}
			absent := false
			if id == t.PlayerID {
				absent = t.Life != w.Life || !w.Player.Alive || w.Held() || w.Player.Location != PoolPlace
			} else {
				n := w.NPC(id)
				absent = n == nil || n.Dead || n.Held > w.Minute || n.Location != PoolPlace || w.Travelling(n)
			}
			if absent {
				missing = append(missing, id)
			}
		}
		if len(missing) > 0 {
			_ = t.Bracket.WithdrawMany(missing)
		}
	}
	if t.Bracket.Finished && t.Bracket.Winner != "" {
		winner := t.Bracket.Winner
		property := w.Properties[PoolPlace]
		if property == nil {
			return
		}
		houseCut := t.Escrow * t.HouseCutPercent / 100
		prize := t.Escrow - houseCut
		winnerName := w.Player.Name
		if winner == t.PlayerID {
			if t.Life != w.Life {
				return
			}
			w.Player.Cash += prize
			w.Player.Earned += max(0, prize-t.Fee)
		} else {
			n := w.NPC(winner)
			if n == nil {
				return
			}
			n.Purse += prize
			winnerName = n.Name
		}
		if w.Own(PoolPlace) {
			w.Earn(houseCut)
		} else if !w.changeBusinessFunds(PoolPlace, houseCut) {
			property.Bankroll += houseCut
		}
		t.HouseCutPaid = houseCut
		w.Log("The tournament is decided", fmt.Sprintf("%s wins the tournament and receives the $%d prize. The hall takes $%d.", winnerName, prize, houseCut), "personal")
		t.PrizePaid = prize
		t.Escrow = 0
		t.Settled = true
		return
	}
	if w.poolUnplayable() || t.Bracket.Finished {
		w.refundPoolTournament()
	}
}
func (w *World) refundPoolTournament() {
	t := w.PoolTournament
	if t == nil || t.Settled {
		return
	}
	if w.Properties[PoolPlace] == nil {
		return
	}
	// Withdrawing forfeits an entry. A cancelled event returns the other entries;
	// forfeited money remains in the hall's till, never in a successor's pocket.
	for id := range t.Deposits {
		if t.Bracket.Withdrawn[id] {
			continue
		}
		if id == t.PlayerID {
			if t.Life != w.Life {
				return
			}
		} else if w.NPC(id) == nil {
			return
		}
	}
	forfeited := 0
	for id, fee := range t.Deposits {
		if t.Bracket.Withdrawn[id] {
			forfeited += fee
			continue
		}
		if id == t.PlayerID {
			w.Player.Cash += fee
		} else {
			w.NPC(id).Purse += fee
		}
	}
	w.Properties[PoolPlace].Bankroll += forfeited
	t.Escrow = 0
	t.Settled = true
	t.Voided = true
	w.Log("The tournament is called off", "Unforfeited entry fees are returned. Forfeited entries remain in the hall’s till.", "personal")
}

// One physical NPC stroke shared by spectator commands and scheduled play.
// It cannot choose a shot on behalf of the player.
func (w *World) PlayPoolTournamentBot(index int) error {
	t := w.PoolTournament
	if t == nil || t.Settled || t.Bracket == nil || index < 0 || index >= len(t.Bracket.Matches) {
		return fmt.Errorf("there is no such tournament game")
	}
	cell := &t.Bracket.Matches[index]
	if cell.Resolved || cell.Rack == nil || cell.Rack.Winner >= 0 {
		return fmt.Errorf("that tournament game is not playing")
	}
	shooter := cell.Players[cell.Rack.Turn]
	if shooter == t.PlayerID {
		return fmt.Errorf("the player must choose their own shot")
	}
	if w.poolUnplayable() {
		return fmt.Errorf("the hall is not fit to play in")
	}
	for _, id := range cell.Players {
		if id == t.PlayerID {
			if t.Life != w.Life || !w.Player.Alive || w.Held() || w.Player.Location != PoolPlace || w.Event != nil {
				return fmt.Errorf("the player's match is interrupted")
			}
		} else {
			n := w.NPC(id)
			if n == nil || n.Dead || n.Held > w.Minute || n.Location != PoolPlace || w.Travelling(n) || t.Bracket.Withdrawn[id] {
				return fmt.Errorf("an entrant is no longer available")
			}
		}
	}
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(shooter))
	identity := hash.Sum64()
	turn, err := billiards.Opponent(billiards.New(), cell.Rack, cell.Rack.Turn, .65+float64(identity%26)/100, identity^uint64(t.Life)*0x9e3779b97f4a7c15^uint64(index)*193^uint64(cell.Rack.Shots)*7919)
	if err != nil {
		return err
	}
	replay, err := billiards.EncodeReplay(turn.Result.Physics)
	if err != nil {
		return err
	}
	if t.Replays == nil {
		t.Replays = map[int]string{}
	}
	if t.Strokes == nil {
		t.Strokes = map[int]PoolStroke{}
	}
	t.Replays[index] = replay
	t.Strokes[index] = PoolStroke{Shooter: cell.Rack.Turn, Intent: turn.Intent, Placement: turn.Placement, Decision: turn.Decision}
	if t.NextStrokes == nil {
		t.NextStrokes = map[int]int{}
	}
	t.NextStrokes[index] = w.Minute + 2
	cell.Rack = &turn.Match
	w.ReconcilePoolTournament()
	return nil
}
