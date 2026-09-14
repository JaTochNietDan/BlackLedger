package core

// Non-player takings follow the deed. Casino proprietors keep the exposed
// gambling float separate from household funds, like player-owned casinos.
func (w *World) businessFunds(id string) int {
	p := w.Properties[id]
	if p == nil {
		return 0
	}
	if f := w.faction(p.Owner); f != nil {
		return f.Cash
	}
	if n := w.NPC(p.Owner); n != nil && !n.Dead {
		if HasBankroll(id) {
			return p.Bankroll
		}
		return w.HouseholdWealth(n)
	}
	return 0
}
func (w *World) changeBusinessFunds(id string, amount int) bool {
	p := w.Properties[id]
	if p == nil || w.Own(id) {
		return false
	}
	if f := w.faction(p.Owner); f != nil {
		f.Cash = max(0, f.Cash+amount)
		return true
	}
	if n := w.NPC(p.Owner); n != nil && !n.Dead {
		if HasBankroll(id) {
			p.Bankroll = max(0, p.Bankroll+amount)
		} else if amount >= 0 {
			n.Purse += amount
		} else {
			w.SpendHouseholdMoney(n, min(-amount, w.HouseholdWealth(n)))
		}
		return true
	}
	return false
}
func (w *World) personalCasino(id string) bool {
	p := w.Properties[id]
	if p == nil || !HasBankroll(id) {
		return false
	}
	n := w.NPC(p.Owner)
	return n != nil && !n.Dead
}
