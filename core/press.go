package core

import "strings"

// A headline with nothing beside it is a line of text. The paper the inbox
// asked for has a picture under every story: a coarse black and white block of
// whatever the story is about.
//
// Nothing here generates an image, and nothing here needs to. What the core
// owes the presentation is the truth of what the picture should be of, and that
// can be read out of the story the city already filed: the premises named in the
// headline, or the person, or the city itself when it is neither.

// Subject is what a story's picture is of.
type Subject struct {
	// Kind is "place", "person" or "city".
	Kind string `json:"kind"`
	// ID is the place or person; empty for the city.
	ID string `json:"id,omitempty"`
	// Name is what to caption it with.
	Name string `json:"name"`
}

// SubjectOf works out what a story is about from the story itself. Headlines
// name premises and people in capitals, which is what makes this readable
// rather than a second record that could disagree with the first.
func (w *World) SubjectOf(s Story) Subject {
	// Whoever filed it said who it was about, so there is nothing to work out.
	// Everything below this is for stories written before a story could say.
	if s.About != "" {
		if n := w.NPC(s.About); n != nil {
			return Subject{Kind: "person", ID: n.ID, Name: n.Name}
		}
		if place, ok := PlaceByID(s.About); ok {
			return Subject{Kind: "place", ID: place.ID, Name: place.Name}
		}
	}
	headline := strings.ToUpper(s.Headline)

	// Premises first: a place is the more specific thing when a headline names
	// both, because "KILLED AT THE MONARCH" is a picture of The Monarch.
	best := Subject{}
	for _, l := range Locations {
		name := strings.ToUpper(l.Name)
		if !strings.Contains(headline, name) {
			continue
		}
		if best.ID == "" || len(name) > len(best.Name) {
			best = Subject{Kind: "place", ID: l.ID, Name: l.Name}
		}
	}
	if best.ID != "" {
		best.Name = placeName(best.ID)
		return best
	}

	// Then anybody the city has a name for, living or not: a story about
	// somebody's death is a picture of them.
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Name != "" && strings.Contains(headline, strings.ToUpper(n.Name)) {
			return Subject{Kind: "person", ID: n.ID, Name: n.Name}
		}
	}

	// And otherwise the city, which is what a newspaper prints when it has no
	// photograph of anything in particular.
	return Subject{Kind: "city", Name: "Bellwether"}
}

func placeName(id string) string {
	if place, ok := PlaceByID(id); ok {
		return place.Name
	}
	return ""
}
