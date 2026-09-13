package core

// Aftermath is the observable scene left by a public killing. Its lifetime uses
// the game clock, so a paused decision, replay or browser reload cannot clean it.
type Aftermath struct {
	ID        string   `json:"id"`
	Target    string   `json:"target"`
	Victim    CueActor `json:"victim"`
	Minute    int      `json:"minute"`
	PoliceAt  int      `json:"police_at"`
	CleanupAt int      `json:"cleanup_at"`
}

func (w *World) recordAftermath(cue VisualCue) {
	if cue.Kind != "killing" {
		return
	}
	active := w.ActiveAftermath()
	for _, actor := range cue.Actors {
		person := w.NPC(actor.ID)
		if person == nil || !person.Dead || person.DiedAt != cue.Minute {
			continue
		}
		duplicate := false
		for _, held := range active {
			if held.Victim.ID == actor.ID {
				duplicate = true
				break
			}
		}
		if !duplicate {
			active = append(active, Aftermath{ID: cue.ID + ":" + actor.ID, Target: cue.Target,
				Victim: actor, Minute: cue.Minute, PoliceAt: cue.Minute + 5, CleanupAt: cue.Minute + 180})
		}
	}
	w.Aftermath = active
}

// ActiveAftermath returns a detached public projection and does not mutate a
// save on GET. Expired records are discarded on the next witnessed killing.
func (w *World) ActiveAftermath() []Aftermath {
	out := []Aftermath{}
	for _, scene := range w.Aftermath {
		if scene.Minute <= w.Minute && w.Minute < scene.CleanupAt {
			out = append(out, scene)
		}
	}
	return out
}
