package core

func (w *World) unattendedPoolGames() []int {
	out := []int{}
	t := w.PoolTournament
	if t == nil || t.Settled || t.Bracket == nil {
		return out
	}
	for i, c := range t.Bracket.Matches {
		if !c.Resolved && c.Rack != nil && c.Rack.Winner < 0 && c.Players[0] != t.PlayerID && c.Players[1] != t.PlayerID {
			out = append(out, i)
		}
	}
	return out
}
func (w *World) schedulePoolTournament() {
	for _, i := range w.unattendedPoolGames() {
		t := w.PoolTournament
		if t.NextStrokes == nil {
			t.NextStrokes = map[int]int{}
		}
		if t.NextStrokes[i] == 0 {
			t.NextStrokes[i] = w.Minute + 2
		}
	}
}
func (w *World) nextPoolTournamentStroke() int {
	next := 0
	for _, i := range w.unattendedPoolGames() {
		due := max(w.Minute+1, w.PoolTournament.NextStrokes[i])
		if next == 0 || due < next {
			next = due
		}
	}
	return next
}
func (w *World) advancePoolTournament() {
	w.ReconcilePoolTournament()
	if w.Event != nil || !w.Player.Alive {
		return
	}
	// Snapshot active matches: a qualifier gets its own two-minute setup period,
	// never another stroke on a newly opened rack in this same clock boundary.
	for _, i := range w.unattendedPoolGames() {
		t := w.PoolTournament
		if t.NextStrokes[i] > 0 && t.NextStrokes[i] <= w.Minute {
			t.NextStrokes[i] = w.Minute + 2
			if err := w.PlayPoolTournamentBot(i); err != nil {
				w.Log("A pause at the baize", err.Error(), "personal")
			}
		}
	}
	w.schedulePoolTournament()
}
