package core

import (
	"blackledger/billiards"
	"encoding/json"
	"math"
	"testing"
)

func tournamentFixture(t *testing.T) *World {
	w := poolFixture(t)
	for _, id := range []string{"leo", "mara", "elena"} {
		n := w.NPC(id)
		n.Dead = false
		n.Location = PoolPlace
		n.Purse = 300
		n.Heading = ""
		n.Sets = 0
		n.Arrives = 0
		n.Held = 0
	}
	return w
}
func fundedTournament(t *testing.T) *World {
	t.Helper()
	w := tournamentFixture(t)
	if err := w.StartPoolTournament([]string{"leo", "mara", "elena"}, 25); err != nil {
		t.Fatal(err)
	}
	return w
}
func tournamentTotal(w *World) int {
	n := w.Player.Cash + w.Properties[PoolPlace].Bankroll
	for _, id := range []string{"leo", "mara", "elena"} {
		n += w.NPC(id).Purse
	}
	if w.PoolTournament != nil {
		n += w.PoolTournament.Escrow
	}
	return n
}
func tournamentWinningPosition(w *World, index int) {
	saved := w.Pool
	w.Pool = &PoolGame{Match: w.PoolTournament.Bracket.Matches[index].Rack}
	winningPoolPosition(w)
	w.Pool = saved
}
func TestPoolTournamentEntryIsFullyFundedAndAtomic(t *testing.T) {
	w := tournamentFixture(t)
	total := tournamentTotal(w)
	for _, ids := range [][]string{{"leo", "leo", "mara"}, {"leo", "mara", "missing"}, {"leo"}} {
		before, _ := json.Marshal(w)
		if err := w.StartPoolTournament(ids, 25); err == nil {
			t.Fatal("invalid entrants accepted")
		}
		after, _ := json.Marshal(w)
		if string(before) != string(after) {
			t.Fatal("partial debit")
		}
	}
	w.NPC("elena").Purse = 24
	before, _ := json.Marshal(w)
	if err := w.StartPoolTournament([]string{"leo", "mara", "elena"}, 25); err == nil {
		t.Fatal("unfunded entry accepted")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("partial stake")
	}
	w.NPC("elena").Purse = 300
	if err := w.StartPoolTournament([]string{"leo", "mara", "elena"}, 25); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != 975 || w.PoolTournament.Escrow != 100 || tournamentTotal(w) != total {
		t.Fatal("entry created/lost money")
	}
	if err := w.StartPool("leo", 20); err == nil {
		t.Fatal("casual game collided with tournament")
	}
	if err := w.StartPoolTournament([]string{"leo", "mara", "elena"}, 25); err == nil {
		t.Fatal("double entry")
	}
}
func TestPoolTournamentPlayerTakesWholePhysicalPrizeOnce(t *testing.T) {
	w := fundedTournament(t)
	total := tournamentTotal(w)
	earned := w.Player.Earned
	b := w.PoolTournament.Bracket
	_ = b.Matches[0].Rack.Concede(1)
	_ = b.Matches[1].Rack.Concede(1)
	w.ReconcilePoolTournament()
	tournamentWinningPosition(w, 2)
	if _, err := b.Matches[2].Rack.Play(billiards.New(), 0, billiards.Shot{Angle: math.Pi, Speed: .6}, billiards.Call{Ball: 8, Pocket: 4}); err != nil {
		t.Fatal(err)
	}
	w.ReconcilePoolTournament()
	w = w.Clone()
	for i := 0; i < 4; i++ {
		w.ReconcilePoolTournament()
	}
	if w.Player.Cash != 1075 || w.Player.Earned != earned+75 || !w.PoolTournament.Settled || w.PoolTournament.Escrow != 0 || tournamentTotal(w) != total {
		t.Fatal("wrong whole-pool payout")
	}
}
func TestPoolTournamentNPCPhysicalWinPaysEntrant(t *testing.T) {
	w := fundedTournament(t)
	total := tournamentTotal(w)
	b := w.PoolTournament.Bracket
	_ = b.Matches[0].Rack.Concede(0)
	_ = b.Matches[1].Rack.Concede(1)
	w.ReconcilePoolTournament()
	tournamentWinningPosition(w, 2)
	if err := w.PlayPoolTournamentBot(2); err != nil {
		t.Fatal(err)
	}
	if !w.PoolTournament.Settled || w.NPC("leo").Purse != 375 || w.Player.Cash != 975 || tournamentTotal(w) != total {
		t.Fatal("NPC prize incorrect")
	}
	if w.PoolTournament.Strokes[2].Shooter != 0 {
		t.Fatal("missing physical shooter")
	}
	if _, err := billiards.DecodeReplay(w.PoolTournament.Replays[2]); err != nil {
		t.Fatal(err)
	}
	before, _ := json.Marshal(w)
	if err := w.PlayPoolTournamentBot(2); err == nil {
		t.Fatal("replayed finished game")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("rejected replay changed payout")
	}
}
func TestPoolTournamentClosureRefundsAndWithdrawalsForfeit(t *testing.T) {
	for _, left := range []bool{false, true} {
		w := fundedTournament(t)
		total := tournamentTotal(w)
		bank := w.Properties[PoolPlace].Bankroll
		if left {
			w.Player.Location = "bar"
			w.ReconcilePool()
		}
		w.Properties[PoolPlace].Condition = 0
		w.ReconcilePool()
		w.ReconcilePool()
		want := 1000
		if left {
			want = 975
			bank += 25
		}
		if !w.PoolTournament.Voided || !w.PoolTournament.Settled || w.Player.Cash != want || w.Properties[PoolPlace].Bankroll != bank || tournamentTotal(w) != total {
			t.Fatal("refund/forfeiture lost money", left)
		}
		for _, id := range []string{"leo", "mara", "elena"} {
			if w.NPC(id).Purse != 300 {
				t.Fatal("entrant not refunded", id)
			}
		}
	}
}
func TestPoolTournamentPinsEligibleEntrantsAndRetainsUnpaidIdentities(t *testing.T) {
	w := fundedTournament(t)
	if !w.PoolOpponentPlaying("leo") || !w.referenced("leo") {
		t.Fatal("unpaid entrant not retained")
	}
	_ = w.PoolTournament.Bracket.Matches[0].Rack.Concede(1)
	w.ReconcilePoolTournament()
	if w.PoolOpponentPlaying("leo") {
		t.Fatal("eliminated entrant pinned")
	}
	if !w.PoolOpponentPlaying("mara") || !w.referenced("leo") {
		t.Fatal("active player or unpaid deposit forgotten")
	}
	before, _ := json.Marshal(w)
	if err := w.PlayPoolTournamentBot(99); err == nil {
		t.Fatal("unknown table accepted")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("invalid game changed state")
	}
}
func TestPoolTournamentNewLifeCannotReceiveOldFees(t *testing.T) {
	w := fundedTournament(t)
	w.DieOf("test", "An isolated tournament death fixture.")
	next, err := Execute(w, Command{Kind: "new_life", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	w = next
	w.Properties[PoolPlace].Condition = 0
	w.ReconcilePool()
	if w.Player.Cash != 90 || !w.PoolTournament.Settled {
		t.Fatal("old fees went to new protagonist")
	}
}

func TestPoolTournamentSimultaneousLossAndBotAuthority(t *testing.T) {
	w := fundedTournament(t)
	total := tournamentTotal(w)
	bank := w.Properties[PoolPlace].Bankroll
	before, _ := json.Marshal(w)
	if err := w.PlayPoolTournamentBot(0); err == nil {
		t.Fatal("bot took player shot")
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("rejected player automation mutated table")
	}
	w.Player.Alive = false
	for _, id := range []string{"leo", "mara", "elena"} {
		w.NPC(id).Dead = true
	}
	w.ReconcilePool()
	if !w.PoolTournament.Settled || !w.PoolTournament.Voided || w.PoolTournament.Bracket.Winner != "" || w.Properties[PoolPlace].Bankroll != bank+100 || tournamentTotal(w) != total {
		t.Fatal("simultaneous departures invented a champion or lost escrow")
	}
}
