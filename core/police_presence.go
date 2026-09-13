package core

// PolicePresence is the visible cordon remaining after a witnessed raid.
// It uses game time and is independent of the latest command result.
type PolicePresence struct {
	ID        string `json:"id"`
	Target    string `json:"target"`
	Minute    int    `json:"minute"`
	CleanupAt int    `json:"cleanup_at"`
}

func (w *World) recordPolicePresence(cue VisualCue) {
	if cue.Kind != "raid" {
		return
	}
	active := w.ActivePolicePresence()
	// One cordon per address. A later real raid can extend its attendance,
	// while duplicate delivery of the same cue cannot restart its clock.
	for i := range active {
		if active[i].Target == cue.Target {
			if cue.Minute > active[i].Minute {
				active[i].Minute = cue.Minute
				active[i].CleanupAt = cue.Minute + 120
			}
			w.PolicePresence = active
			return
		}
	}
	w.PolicePresence = append(active, PolicePresence{ID: cue.ID, Target: cue.Target,
		Minute: cue.Minute, CleanupAt: cue.Minute + 120})
}

func (w *World) ActivePolicePresence() []PolicePresence {
	out := []PolicePresence{}
	for _, scene := range w.PolicePresence {
		if scene.Minute <= w.Minute && w.Minute < scene.CleanupAt {
			out = append(out, scene)
		}
	}
	return out
}
