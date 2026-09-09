package core

import "fmt"

// The city could rob premises and the player could rob premises, and nobody in
// it could rob a person. That is the crudest thing anybody does here and it was
// the one thing missing: walking up to somebody you know and taking what they
// are carrying.
//
// It is personal in a way nothing else is. There is no unattributed headline
// about a till: a man knows who took his money, and the better known the face
// the more certain he is. A famous man is better at every part of this except
// getting away with it.

const (
	// MuggingMinutes is how long it takes to do and get clear.
	MuggingMinutes = 45
	// MuggingHeat is the attention it draws.
	MuggingHeat = 9
	// RecognisedAt is the presence above which a face is a face people know.
	// Everything else in this game rewards standing; this is the one place it
	// costs.
	RecognisedAt = 35
)

// Pockets is what somebody is carrying, derived from what they are: standing
// inside an organization, and how well that organization is doing. Nobody in
// this city walks around with nothing.
func (w *World) Pockets(n *NPC) int {
	if n == nil || n.Dead {
		return 0
	}
	purse := 30 + n.Rank*2
	if f := w.faction(n.Faction); f != nil {
		purse += min(400, f.Cash/40) * max(1, n.Rank) / RankLeader
	}
	if IsOfficial(n.ID) {
		purse *= 3 // a man with a salary and an arrangement carries both
	}
	if TemperamentOf(n).ID == "vain" {
		purse = purse * 3 / 2 // it is on him, and it shows
	}
	return purse
}

// MuggingTarget finds who is here and worth taking something off. Only somebody
// the player knows: you do not walk up to a stranger in this city.
func (w *World) MuggingTarget(location string) (*NPC, bool) {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != location || !w.Known(n) {
			continue
		}
		if len(w.Player.Crew) > 0 && w.Player.Crew[0].ID == n.ID {
			continue // not your own people
		}
		return n, true
	}
	return nil, false
}

// MuggingReadiness explains why nobody can be taken for what they have, or
// returns "".
func (w *World) MuggingReadiness(location string) string {
	if w.Player.Location != location {
		return "You would have to be there"
	}
	n, ok := w.MuggingTarget(location)
	if !ok {
		return "There is nobody here you know well enough to walk up to"
	}
	if w.Player.Health < 40 {
		return "You are in no condition for this"
	}
	if w.Pockets(n) < 40 {
		return n.Name + " has nothing on them worth the trouble"
	}
	return ""
}

// Recognised reports whether the person robbed can name who did it. A known
// face is a known face, and a crewman sent in the player's place is the one
// being described instead.
func (w *World) Recognised(hand Hand) bool {
	if hand.Crew {
		return false
	}
	return w.Presence() >= RecognisedAt
}

// Mug takes what somebody is carrying. The same rules whoever does it.
func (w *World) Mug(location string, hand Hand) error {
	if reason := w.MuggingReadiness(location); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if hand.Crew {
		if reason := w.DelegateReadiness(); reason != "" {
			return fmt.Errorf("%s", reason)
		}
	}
	mark, _ := w.MuggingTarget(location)
	purse := w.Pockets(mark)
	place, _ := PlaceByID(location)

	// A man standing next to a strong organization is a different proposition
	// from a man standing on his own, which is most of what decides who is
	// worth walking up to.
	defence := w.Poise(mark) + 15
	if f := w.faction(mark.Faction); f != nil {
		defence += f.Power / 3
	}
	odds := .74 + w.HandEdge(hand) - float64(defence)/250
	odds = min64(.85, max64(.12, odds))

	if w.Random() >= odds {
		injury := 12 + int(w.Random()*20)
		health := w.Player.Health
		w.HandHurt(hand, injury, "take something off "+mark.Name)
		w.Player.Heat = min(100, w.Player.Heat+w.HandHeat(hand, MuggingHeat))
		// He knows either way: somebody put hands on him.
		w.Resent(mark.ID, w.crewID(hand), 30, "what was tried at "+place.Name)
		w.answerFor(mark, hand, 12)
		if hand.Crew {
			w.Log("It went wrong at "+place.Name, fmt.Sprintf("%s went at %s and came off worse. %s knows who sent them.", hand.Name, mark.Name, mark.Name), "danger")
		} else {
			w.Log("It went wrong at "+place.Name, fmt.Sprintf("%s was not as easy as they looked. You came away with nothing and %d less health, and they saw all of it.", mark.Name, health-w.Player.Health), "danger")
		}
		if w.Player.Health <= 0 {
			w.Die("A robbery at " + place.Name + " that should have been simple.")
		}
		return nil
	}

	w.Earn(purse)
	w.Player.Heat = min(100, w.Player.Heat+w.HandHeat(hand, MuggingHeat))
	w.Player.Respect += w.HandRespectFor(hand, 3)
	if f := w.faction(mark.Faction); f != nil {
		f.Cash = max(0, f.Cash-purse/2)
	}

	// Whether he can put a name to it is the whole of the risk.
	if w.Recognised(hand) {
		w.Resent(mark.ID, w.crewID(hand), 45, "being robbed at "+place.Name)
		w.answerFor(mark, hand, 25)
		w.Log("They know your face", fmt.Sprintf("$%d off %s, and somebody with your standing is not a person anybody has to describe twice.", purse, mark.Name), "danger")
	} else {
		w.Resent(mark.ID, w.crewID(hand), 20, "being robbed at "+place.Name)
		w.Log("Taken off "+mark.Name, fmt.Sprintf("$%d, and nobody who could put a name to it. %s will be asking, though.", purse, mark.Name), "politics")
	}
	w.Report("robbery", "ROBBERY IN "+upper(place.Name),
		w.unattributed(place.Name, fmt.Sprintf("Somebody was robbed at %s. Police have asked anybody who saw it to come forward.", place.Name)))

	// Robbing a man with a title is not robbing a man.
	if IsOfficial(mark.ID) {
		w.Player.Heat = min(100, w.Player.Heat+20)
		w.EndRetainerQuietly(mark.ID)
		w.Log("Of all the people in this city", fmt.Sprintf("%s is not somebody anybody robs. Whatever you had with that building, you do not have it now.", mark.Name), "danger")
	}
	return nil
}

// crewID is who the mark holds a grudge against. The player is not a person in
// the city's own records, so when they do it themselves the mark's answer is
// their organization's standing rather than a private grievance — blaming the
// crewman for something the player did was wrong and read as a bug.
func (w *World) crewID(hand Hand) string {
	if hand.Crew && len(w.Player.Crew) > 0 {
		return w.Player.Crew[0].ID
	}
	return ""
}

// answerFor is the mark's organization taking an interest, which is the part
// that outlives the money.
func (w *World) answerFor(mark *NPC, hand Hand, goodwill int) {
	f := w.faction(mark.Faction)
	if f == nil {
		return
	}
	if !w.Recognised(hand) && goodwill < 20 {
		return // nobody to be angry at yet
	}
	f.Goodwill = max(-100, f.Goodwill-goodwill)
	if goodwill >= 20 {
		w.RetaliationFrom(f.ID)
	}
}
