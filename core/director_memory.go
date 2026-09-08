package core

// ArrangementMemory preserves the actual offer separately from authoritative results.
// An offer's claims are not promoted into facts merely because the director wrote them.
type StoryConnection struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Result string `json:"result"`
}

type ArrangementMemory struct {
	Location    string `json:"location,omitempty"`
	ParentID    string `json:"parent_id,omitempty"`
	ID          string `json:"id"`
	Life        int    `json:"life"`
	Minute      int    `json:"minute"`
	Title       string `json:"title"`
	Offer       string `json:"offer"`
	Speaker     string `json:"speaker"`
	Operation   string `json:"operation"`
	Beneficiary string `json:"beneficiary"`
	Status      string `json:"status"`
	Result      string `json:"result,omitempty"`
}

func (w *World) RememberArrangement(scene *Scene, status string) {
	id := scene.ID
	if scene.JobID != "" {
		id = scene.JobID
	}
	for i := range w.Arrangements {
		if w.Arrangements[i].ID == id {
			w.Arrangements[i].Status = status
			if status == "completed" {
				w.Arrangements[i].Result = scene.Outcome
			}
			return
		}
	}
	// Old police scenes have no original offer. Do not invent one from the police dialogue.
	if scene.Kind != "proposal" {
		return
	}
	memory := ArrangementMemory{Location: scene.Target, ID: id, Life: w.Life, Minute: w.Minute, Title: scene.Title, Offer: scene.Body, Speaker: scene.Speaker, Operation: scene.Operation, Beneficiary: scene.Beneficiary, Status: status}
	if status == "completed" {
		memory.Result = scene.Outcome
	}
	if scene.Connection != nil {
		memory.ParentID = scene.Connection.ID
	}
	w.Arrangements = append(w.Arrangements, memory)
	if len(w.Arrangements) > 24 {
		w.Arrangements = w.Arrangements[len(w.Arrangements)-24:]
	}
}

// Select the least recently represented operation, including pending offers.
// Fiction remains model-authored; the core gives it a varied mechanical brief.
func (w *World) NextDirectorOperation() string {
	seen := map[string]int{}
	index := 1
	for _, m := range w.Arrangements {
		if m.Life == w.Life {
			seen[m.Operation] = index
			index++
		}
	}
	if w.Event != nil && w.Event.Operation != "" {
		seen[w.Event.Operation] = index
		index++
	}
	for _, offer := range w.Offers {
		if offer.Event != nil {
			seen[offer.Event.Operation] = index
			index++
		}
	}
	selected := "mediation"
	for _, operation := range []string{"collection", "courier"} {
		if seen[operation] < seen[selected] {
			selected = operation
		}
	}
	// Work that only exists because of something the city did takes precedence
	// over an errand, but only when it is the least recently told story. A war
	// should not produce nothing but escorts.
	// Situational work wins a tie. Without this the ordinary rotation always
	// holds the tie on a fresh campaign, so a city already at war could never
	// lead with a war story.
	for _, situational := range w.SituationalOperations() {
		if seen[situational.ID] <= seen[selected] {
			selected = situational.ID
		}
	}
	return selected
}

// Keep the next opportunity with a contact whose latest work actually completed.
// Declines and failures remain context, but cannot become successful callbacks.
func (w *World) DirectorConnection() *ArrangementMemory {
	for i := len(w.Arrangements) - 1; i >= 0; i-- {
		m := &w.Arrangements[i]
		if m.Life != w.Life {
			continue
		}
		if m.Status == "completed" {
			// A follow-up closes this short thread; do not force an endless same-contact chain.
			if m.ParentID != "" {
				return nil
			}
			for _, other := range w.Arrangements {
				if other.Life == w.Life && other.ParentID != "" && other.ParentID == m.ID {
					return nil
				}
			}
			if w.Event != nil && w.Event.Connection != nil && w.Event.Connection.ID == m.ID {
				return nil
			}
			for _, offer := range w.Offers {
				if offer.Event != nil && offer.Event.Connection != nil && offer.Event.Connection.ID == m.ID {
					return nil
				}
			}
			copy := *m
			return &copy
		}
		// A subsequent unresolved or refused offer should not force an older thread.
		return nil
	}
	return nil
}
