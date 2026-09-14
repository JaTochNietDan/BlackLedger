package core

import "math"

// StreetSegment is an observed NPC journey during one committed travel command.
// Progress is the starting fraction; no future trip or private plan is exposed.
type StreetSegment struct {
	Journeying
	FromMinute  int     `json:"from_minute"`
	ToMinute    int     `json:"to_minute"`
	EndProgress float64 `json:"end_progress"`
}

func (w *World) recordStreetSegment(street []Journeying, from, to int) {
	for _, j := range street {
		end := min(to, from+j.Minutes)
		if end <= from {
			continue
		}
		progress := math.Min(1, j.Progress+(1-j.Progress)*float64(end-from)/float64(max(1, j.Minutes)))
		merged := false
		for i := len(w.streetTravel) - 1; i >= 0; i-- {
			old := &w.streetTravel[i]
			if old.ID == j.ID {
				if old.ToMinute == from && old.FromID == j.FromID && old.ToID == j.ToID && old.Vehicle == j.Vehicle && old.EndProgress < 1 {
					old.ToMinute = end
					old.EndProgress = progress
					merged = true
				}
				break
			}
		}
		if !merged {
			w.streetTravel = append(w.streetTravel, StreetSegment{Journeying: j, FromMinute: from, ToMinute: end, EndProgress: progress})
		}
	}
}
