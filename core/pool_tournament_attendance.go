package core

import "hash/fnv"

// Eight regulars plan a tournament evening. Selection does not depend on their
// current journey, so starting to walk cannot reshuffle the guest list.
func (w *World) poolTournamentVisit(n *NPC) bool {
	opens, closes := w.poolTournamentWindow()
	if w.Minute < opens-60 || w.Minute >= closes || w.LastPoolTournamentSlot == opens || (w.PoolTournament != nil && !w.PoolTournament.Settled) || w.poolUnplayable() {
		return false
	}
	war := w.CityAtWar()
	eligible := func(p *NPC) bool {
		return !p.Dead && !w.keepsPost(p) && p.Post != "" && w.Pockets(p) >= PoolTournamentFee && !(war && w.staysIn(p))
	}
	if !eligible(n) {
		return false
	}
	score := func(id string) uint64 {
		h := fnv.New64a()
		h.Write([]byte(id))
		return h.Sum64() ^ uint64(opens/poolTournamentPeriod)*0x9e3779b97f4a7c15
	}
	target := score(n.ID)
	ahead := 0
	for i := range w.NPCs {
		p := &w.NPCs[i]
		if p.ID == n.ID || !eligible(p) {
			continue
		}
		value := score(p.ID)
		if value < target || (value == target && p.ID < n.ID) {
			ahead++
			if ahead >= 8 {
				return false
			}
		}
	}
	return true
}
