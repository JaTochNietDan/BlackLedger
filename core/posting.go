package core

import "fmt"

// Somebody who answers to the player stands wherever they happen to be. They
// add to what the organization is worth in a fight, in the abstract, and there
// is no way to put one of them anywhere — so a business being taken apart is
// something the player can only watch and pay for.
//
// Somebody on the door is the oldest answer to that. They make the place
// harder to walk into, for the police as much as for anybody, and they are the
// one standing in it when somebody comes.

const (
	// PostingMinutes is the morning it takes to arrange.
	PostingMinutes = 45
	// PostingDefence is what somebody on the door is worth against anybody
	// coming for the premises, before what they are worth themselves.
	PostingDefence = 12
	// PostingSteadies is the share of an ordinary robbery he turns away.
	PostingSteadies = .45
)

// PostedAt is whoever is standing at a place, or nil.
func (w *World) PostedAt(id string) *NPC {
	prop := w.Properties[id]
	if prop == nil || prop.Posted == "" {
		return nil
	}
	n := w.NPC(prop.Posted)
	if n == nil || n.Dead || w.Inside(n) || n.Faction != w.PlayerOrganizationID() {
		prop.Posted = ""
		return nil
	}
	return n
}

// Unposted is everybody who answers to the player and is not standing anywhere
// in particular.
func (w *World) Unposted() []*NPC {
	posted := map[string]bool{}
	for _, l := range Locations {
		if n := w.PostedAt(l.ID); n != nil {
			posted[n.ID] = true
		}
	}
	out := []*NPC{}
	for _, n := range w.OwnPeople() {
		// Somebody the police are holding is not standing anywhere.
		if !posted[n.ID] && !w.Inside(n) {
			out = append(out, n)
		}
	}
	return out
}

// PostingDefenceAt is what somebody on the door adds to holding a place: what
// they are worth, and how much they mean it.
func (w *World) PostingDefenceAt(id string) int {
	n := w.PostedAt(id)
	if n == nil {
		return 0
	}
	return PostingDefence + w.Poise(n)/4 + n.Trust/10
}

// PostReadiness explains why nobody can be put on the door here, or returns "".
func (w *World) PostReadiness(id string) string {
	if !w.Own(id) {
		return "This is not a business of yours"
	}
	if prop := w.Properties[id]; prop == nil || prop.Income <= 0 {
		return "There is nothing here to stand over"
	}
	if w.PostedAt(id) != nil {
		return w.PostedAt(id).Name + " is already on the door here"
	}
	if len(w.Unposted()) == 0 {
		return "Everybody who answers to you is standing somewhere already"
	}
	return ""
}

// Post puts the most reliable of the player's people on the door.
func (w *World) Post(id string) error {
	if reason := w.PostReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	free := w.Unposted()
	best := free[0]
	for _, n := range free {
		if n.Trust > best.Trust {
			best = n
		}
	}
	w.Properties[id].Posted = best.ID
	best.Location = id
	place, _ := PlaceByID(id)
	w.Log(best.Name+" is on the door at "+place.Name, fmt.Sprintf("Somewhere to stand and something to do. %s is harder to walk into now, and they are the one standing in it when somebody comes.", place.Name), "business")
	return nil
}

// Unpost takes somebody off a door.
func (w *World) Unpost(id string) error {
	n := w.PostedAt(id)
	if n == nil {
		return fmt.Errorf("nobody is on the door there")
	}
	w.Properties[id].Posted = ""
	place, _ := PlaceByID(id)
	w.Log(n.Name+" comes off the door at "+place.Name, "They are free to be put somewhere else now.", "business")
	return nil
}

// PostingDescription is who is standing where, for the interface.
func (w *World) PostingDescription(id string) map[string]any {
	n := w.PostedAt(id)
	if n == nil {
		return nil
	}
	return map[string]any{
		"id": n.ID, "name": n.Name, "trust": n.Trust,
		"worth": w.PostingDefenceAt(id),
	}
}

// StoodInIt is whoever a raid on a place reaches first. Somebody on the door
// is the reason the place held, and the reason it costs somebody when it does
// not: they are standing in the doorway either way. Without one, the raid reaches
// whoever the organization can least afford to lose, as before.
func (w *World) StoodInIt(id string) *NPC {
	return w.PostedAt(id)
}
