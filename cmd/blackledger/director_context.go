package main

import "blackledger/core"

// Keep outcomes and identity without repeatedly quoting the same old offer prose.
// DirectorConnection separately supplies the full one-job callback context.
func arrangementBriefs(w *core.World) []core.ArrangementMemory {
	out := make([]core.ArrangementMemory, len(w.Arrangements))
	copy(out, w.Arrangements)
	for i := range out {
		out[i].Offer = ""
	}
	return out
}
func recentWorldChanges(w *core.World) []core.Record {
	out := []core.Record{}
	for _, r := range w.History {
		if r.Kind != "story" {
			out = append(out, r)
		}
	}
	if len(out) > 12 {
		out = out[len(out)-12:]
	}
	return out
}
