package core

import (
	"blackledger/billiards"
	"fmt"
)

const PoolPlace = "poolhall"
const PoolMinStake = 10
const PoolMaxStake = 500

type PoolGame struct {
	Place        string           `json:"place"`
	Life         int              `json:"life"`
	Opponent     string           `json:"opponent"`
	OpponentName string           `json:"opponent_name"`
	Stake        int              `json:"stake"`
	Escrow       int              `json:"escrow"`
	Settled      bool             `json:"settled"`
	Voided       bool             `json:"voided,omitempty"`
	Match        *billiards.Match `json:"match"`
	Replay       string           `json:"replay,omitempty"`
}

func (w *World) poolUnplayable() bool {
	if p := w.Properties[PoolPlace]; p == nil || p.Condition <= 0 {
		return true
	}
	for _, fire := range w.ActiveBuildingFires() {
		if fire.Target == PoolPlace && w.Minute < fire.ExtinguishedAt {
			return true
		}
	}
	return false
}

func (w *World) PoolReadiness(who string, stake int) string {
	if !w.Player.Alive || w.Held() || w.Event != nil {
		return "Finish the current situation before starting a rack."
	}
	if w.Player.Location != PoolPlace {
		return "You must be inside The Green Baize."
	}
	if w.poolUnplayable() {
		return "The hall is not fit to play in."
	}
	if w.Pool != nil && !w.Pool.Settled {
		return "There is already a rack with money on it."
	}
	if w.Seated != "" || (w.Game != nil && !w.Game.Over) || (w.Hand != nil && !w.Hand.Done) {
		return "Get up from the other game first."
	}
	if stake < PoolMinStake || stake > PoolMaxStake {
		return fmt.Sprintf("Each player must put up $%d to $%d.", PoolMinStake, PoolMaxStake)
	}
	if w.Player.Cash < stake {
		return "You cannot cover your half of the stake."
	}
	n := w.NPC(who)
	if n == nil || n.Dead || n.Location != PoolPlace || w.Travelling(n) || n.Held > w.Minute {
		return "That person is not here and available to play."
	}
	if w.Pockets(n) < stake {
		return "They cannot cover that stake."
	}
	return ""
}
func (w *World) StartPool(who string, stake int) error {
	if reason := w.PoolReadiness(who, stake); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	match, err := billiards.NewMatch(0)
	if err != nil {
		return err
	}
	n := w.NPC(who)
	if err = w.Pay(stake); err != nil {
		return err
	}
	n.Purse -= stake
	// An accepted rack replaces a planned departure; ordinary routine changes
	// do not let somebody walk away from a game they have agreed to finish.
	n.Heading, n.Errand, n.Sets, n.Arrives = "", "", 0, 0
	w.Pool = &PoolGame{Place: PoolPlace, Life: w.Life, Opponent: n.ID, OpponentName: n.Name, Stake: stake, Escrow: 2 * stake, Match: match}
	w.Log("A stake on the baize", fmt.Sprintf("You and %s each put up $%d at The Green Baize. The winner takes the $%d in the middle; leaving concedes the rack.", n.Name, stake, 2*stake), "personal")
	return nil
}
func (w *World) poolSession() (*PoolGame, error) {
	g := w.Pool
	if g == nil || g.Settled || g.Match == nil {
		return nil, fmt.Errorf("there is no active billiards match")
	}
	if g.Life != w.Life || !w.Player.Alive || w.Player.Location != g.Place {
		return nil, fmt.Errorf("you have left this rack")
	}
	n := w.NPC(g.Opponent)
	if n == nil || n.Dead || n.Location != g.Place || w.Travelling(n) || n.Held > w.Minute {
		return nil, fmt.Errorf("the other player is no longer at the table")
	}
	if w.poolUnplayable() {
		return nil, fmt.Errorf("the hall is no longer fit to play in")
	}
	if w.Event != nil {
		return nil, fmt.Errorf("finish the current situation before playing")
	}
	return g, nil
}
func (w *World) PlacePoolCue(x, y float64) error {
	g, err := w.poolSession()
	if err != nil {
		return err
	}
	return g.Match.PlaceCue(billiards.New(), 0, x, y)
}
func (w *World) DecidePoolBreak(choice string) error {
	g, err := w.poolSession()
	if err != nil {
		return err
	}
	return g.Match.DecideBreak(billiards.New(), 0, choice)
}
func (w *World) PlayPoolShot(shot billiards.Shot, call billiards.Call) error {
	g, err := w.poolSession()
	if err != nil {
		return err
	}
	// Match.Play is atomic. Encoding is also done before replacing saved state.
	next := *g.Match
	next.Balls = append([]billiards.Ball(nil), g.Match.Balls...)
	result, err := next.Play(billiards.New(), 0, shot, call)
	if err != nil {
		return err
	}
	replay, err := billiards.EncodeReplay(result.Physics)
	if err != nil {
		return err
	}
	g.Match = &next
	g.Replay = replay
	w.settlePool()
	return nil
}
func (w *World) ConcedePool() error {
	g := w.Pool
	if g == nil || g.Settled || g.Match == nil || g.Life != w.Life {
		return fmt.Errorf("there is no active rack to concede")
	}
	w.ReconcilePool()
	if g.Settled {
		return nil
	}
	if err := g.Match.Concede(0); err != nil {
		return err
	}
	w.settlePool()
	return nil
}
func (w *World) PoolOpponentPlaying(id string) bool {
	g := w.Pool
	return g != nil && !g.Settled && g.Match != nil && g.Match.Winner < 0 && g.Life == w.Life && w.Player.Alive && w.Player.Location == g.Place && g.Opponent == id
}

// ReconcilePool resolves forced departures and death even when the current
// command is unrelated to billiards. A different life never receives an old
// protagonist's winnings. The unpaid opponent remains referenced for settlement.
func (w *World) ReconcilePool() {
	g := w.Pool
	if g == nil || g.Settled || g.Match == nil {
		return
	}
	n := w.NPC(g.Opponent)
	if g.Life != w.Life || !w.Player.Alive || w.Player.Location != g.Place {
		if g.Match.Winner < 0 {
			_ = g.Match.Concede(0)
		}
	} else if w.poolUnplayable() && g.Match.Winner < 0 {
		// An unusable hall voids the contest; neither player profits from a fire.
		if n == nil {
			return
		}
		w.Player.Cash += g.Stake
		n.Purse += g.Stake
		g.Escrow = 0
		g.Settled = true
		g.Voided = true
		w.Log("The rack is called off", "The hall is no longer fit to play in. Both stakes are returned.", "personal")
		return
	} else if n == nil || n.Dead || n.Location != g.Place || w.Travelling(n) || n.Held > w.Minute {
		if g.Match.Winner < 0 {
			_ = g.Match.Concede(1)
		}
	}
	w.settlePool()
}
func (w *World) settlePool() {
	g := w.Pool
	if g == nil || g.Settled || g.Match == nil || g.Match.Winner < 0 {
		return
	}
	// A damaged old-life save cannot redirect its escrow into a new life.
	// Normally DieOf / new_life settle it before the life number changes.
	if g.Match.Winner == 0 && g.Life != w.Life {
		return
	}
	if g.Match.Winner == 1 && w.NPC(g.Opponent) == nil {
		return
	}
	amount := g.Escrow
	if g.Match.Winner == 0 {
		w.Player.Cash += amount
		w.Player.Earned += max(0, amount-g.Stake)
		w.Log("The rack is yours", fmt.Sprintf("You take the $%d staked with %s at The Green Baize, a $%d profit.", amount, g.OpponentName, max(0, amount-g.Stake)), "personal")
	} else {
		w.NPC(g.Opponent).Purse += amount
		w.Log("The other side of the baize", fmt.Sprintf("%s takes the $%d in the middle. Your stake cost you $%d.", g.OpponentName, amount, g.Stake), "personal")
	}
	g.Escrow = 0
	g.Settled = true
}
