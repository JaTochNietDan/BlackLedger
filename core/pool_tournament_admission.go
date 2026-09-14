package core

import (
	"fmt"
	"sort"
)

const PoolTournamentFee = 25
const poolTournamentPeriod = 3 * 1440
const poolTournamentHour = 18 * 60
const poolTournamentWindow = 2 * 60

// The notice is a fixed city schedule. Reading it never rolls entrants or money.
func (w *World) poolTournamentWindow() (opens, closes int) {
	opens = (w.Minute/poolTournamentPeriod)*poolTournamentPeriod + poolTournamentHour
	if w.Minute >= opens+poolTournamentWindow {
		opens += poolTournamentPeriod
	}
	return opens, opens + poolTournamentWindow
}
func (w *World) poolTournamentEntrants() []string {
	ids := []string{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != PoolPlace || n.Held > w.Minute || w.Travelling(n) || w.Pockets(n) < PoolTournamentFee || w.PoolOpponentPlaying(n.ID) {
			continue
		}
		ids = append(ids, n.ID)
	}
	sort.Strings(ids)
	// A full eight-player draw takes priority; otherwise use four. Never invent
	// participants or waive the fee to fill a bracket.
	if len(ids) >= 7 {
		return ids[:7]
	}
	if len(ids) >= 3 {
		return ids[:3]
	}
	return ids
}
func (w *World) poolTournamentAdmission() (int, int, []string, string) {
	opens, closes := w.poolTournamentWindow()
	ids := w.poolTournamentEntrants()
	reason := ""
	switch {
	case w.PoolTournament != nil && !w.PoolTournament.Settled:
		reason = "A tournament is already underway."
	case w.Minute < opens:
		reason = "Entry opens at the posted time."
	case w.LastPoolTournamentSlot == opens:
		reason = "This evening's tournament has already been held."
	case len(ids) < 3:
		reason = "At least three other funded players must be here to enter."
	default:
		reason = w.PoolReadiness(ids[0], PoolTournamentFee)
	}
	return opens, closes, ids, reason
}
func (w *World) EnterPoolTournament() error {
	opens, _, ids, reason := w.poolTournamentAdmission()
	if reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.StartPoolTournament(ids, PoolTournamentFee); err != nil {
		return err
	}
	w.LastPoolTournamentSlot = opens
	return nil
}
func (w *World) PoolTournamentNotice() any {
	if w.Player.Location != PoolPlace {
		return nil
	}
	opens, closes, ids, reason := w.poolTournamentAdmission()
	entrants := []map[string]any{}
	for _, id := range ids {
		entrants = append(entrants, map[string]any{"id": id, "name": w.NPC(id).Name})
	}
	return map[string]any{"can_host": w.Own(PoolPlace), "opens": opens, "closes": closes, "fee": PoolTournamentFee, "entrants": entrants, "pot": PoolTournamentFee * (len(ids) + 1), "can_enter": reason == "", "unavailable": reason}
}
