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
	employed := n != nil && n.Faction != "" && n.Faction == prop.Owner
	if w.Player.Alive && w.Own(id) {
		if _, hired := w.NamedHands(prop.Posted); hired {
			employed = true
		}
	}
	if n == nil || n.Dead || w.Inside(n) || !employed {
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
		if !posted[n.ID] && w.OutOfReach(n.ID) == "" && !w.personHasLegacyTask(n.ID) {
			out = append(out, n)
		}
	}
	return out
}

// PostingDefenceAt is what somebody on the door adds to holding a place: what
// they are worth, and how much they mean it.
func (w *World) PostingDefenceAt(id string) int {
	n := w.PostedAt(id)
	// Somebody still on their way to the door is not on the door.
	if n == nil || n.Location != id || w.Travelling(n) {
		return 0
	}
	return PostingDefence + w.Poise(n)/4 + w.guardReliability(id, n)/10
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
	if len(w.OwnPeople()) == 0 {
		return "Sign someone into your organization before assigning a guard"
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
	// Replace a scheduled routine before it starts; active travel was excluded.
	best.Heading, best.Arrives, best.Sets, best.Errand = "", 0, 0, ""
	w.Properties[id].Posted = best.ID
	place, _ := PlaceByID(id)
	if best.Location == id {
		w.Log(best.Name+" is on the door at "+place.Name, fmt.Sprintf("Somewhere to stand and something to do. %s is harder to walk into now, and they are the one standing in it when somebody comes.", place.Name), "business")
		return nil
	}
	// Everybody else in this city walks. The player's own people used to be the
	// only ones who could be in two places in the same minute, and the door was
	// defended from the moment of the decision rather than from the moment
	// somebody was standing in it.
	best.Heading = id
	best.Sets = 0
	best.Errand = "sent to stand on the door at " + place.Name
	best.Arrives = w.Minute + TravelMinutes(best.Location, id)
	w.noticed(best, true)
	w.Log(best.Name+" is sent to the door at "+place.Name,
		fmt.Sprintf("It is %d minutes across the city. Until they get there the door is exactly as easy to walk into as it was.", best.Arrives-w.Minute), "business")
	return nil
}

// Unpost takes somebody off a door.
func (w *World) Unpost(id string) error {
	if !w.Own(id) {
		return fmt.Errorf("this is not a business of yours")
	}
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
	out := map[string]any{
		"id": n.ID, "name": n.Name, "trust": w.guardReliability(id, n),
		"worth": w.PostingDefenceAt(id),
	}
	// Naming somebody on a door they have not reached tells the player the
	// place is held when it is not, which is exactly the moment they would
	// stop worrying about it.
	if w.Travelling(n) && n.Heading == id {
		out["coming"] = true
		out["minutes"] = max(1, n.Arrives-w.Minute)
	} else if n.Location != id || w.Travelling(n) {
		out["away"] = true
	}
	return out
}

// StoodInIt is whoever a raid on a place reaches first. Somebody on the door
// is the reason the place held, and the reason it costs somebody when it does
// not: they are standing in the doorway either way. Without one, the raid reaches
// whoever the organization can least afford to lose, as before.
func (w *World) StoodInIt(id string) *NPC {
	n := w.PostedAt(id)
	if n == nil || n.Location != id || w.Travelling(n) {
		return nil
	}
	return n
}

// Signed members answer through trust; hired associates through paid loyalty.
func (w *World) guardReliability(id string, n *NPC) int {
	if w.Own(id) && w.Player.Alive {
		for _, c := range w.Player.Crew {
			if c.ID == n.ID {
				return c.Loyalty
			}
		}
	}
	return n.Trust
}
