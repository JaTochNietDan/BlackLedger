package core

import "fmt"

// Every violent thing the player could do, they had to do standing there. A man
// with a crew does not do that, and the reason he has a crew is so he does not
// have to.
//
// Sending somebody is worse at the job and better for the life expectancy. Your
// own hands bring your standing, your gun and your armour to it and put all of
// the risk on you. Leo brings his loyalty and takes the risk himself, which
// means the risk is now that you lose Leo.

const (
	// HandLoyalty is the loyalty below which nobody goes out on this kind of
	// errand for anybody.
	HandLoyalty = 40
	// HandRespect is the share of the standing the player earns when it was
	// somebody else's hands. A man who sends people is respected less than a
	// man who goes, which is most of why anybody goes.
	HandRespect = 2
	// HandHeatRelief is the police attention that lands on the man who was
	// actually seen rather than on the player.
	HandHeatRelief = 6
	// HandLoyaltyCost is what a job that goes wrong costs the person who was
	// sent to do it, in how they feel about being sent.
	HandLoyaltyCost = 25
)

// Hand is who does a piece of work.
type Hand struct {
	// ID binds the work to a particular hired person, even if the roster changes.
	ID string
	// Crew is true when somebody else's hands were on it.
	Crew bool
	// Name is who did it, for the record.
	Name string
}

// OwnHands is the player doing it themselves.
func (w *World) OwnHands() Hand { return Hand{} }

// CrewHands is the player sending somebody, and whether there is anybody to
// send.
func (w *World) CrewHands() (Hand, bool) {
	if len(w.Player.Crew) == 0 {
		return Hand{}, false
	}
	return Hand{ID: w.Player.Crew[0].ID, Crew: true, Name: w.Player.Crew[0].Name}, true
}

// DelegateReadiness explains why nobody can be sent, or returns "".
func (w *World) DelegateReadiness() string {
	for _, o := range w.CrewOrders {
		if o.Kind == "collections" && o.active() {
			return "The headquarters collection round must finish first"
		}
	}
	if len(w.Player.Crew) == 0 {
		return "You have nobody to send"
	}
	if w.CrewOrderFor(w.Player.Crew[0].ID) != nil {
		return "They are on a headquarters assignment"
	}
	if len(w.Tasks) > 0 {
		return w.Player.Crew[0].Name + " is already on assignment"
	}
	hand, _ := w.CrewHands()
	return w.HandReadiness(hand)
}

// NamedHands resolves either an original associate or a signed family member.
func (w *World) NamedHands(id string) (Hand, bool) {
	for _, c := range w.Player.Crew {
		if c.ID == id {
			return Hand{ID: id, Crew: true, Name: c.Name}, true
		}
	}
	for _, n := range w.OwnPeople() {
		if n.ID == id {
			return Hand{ID: id, Crew: true, Name: n.Name}, true
		}
	}
	return Hand{}, false
}

func (w *World) handMember(hand Hand) (Crew, bool) {
	if !hand.Crew {
		return Crew{}, false
	}
	id := hand.ID
	// Older synchronous callers have no ID; keep their original default.
	if id == "" && len(w.Player.Crew) > 0 {
		id = w.Player.Crew[0].ID
	}
	for _, c := range w.Player.Crew {
		if c.ID == id {
			return c, true
		}
	}
	for _, n := range w.OwnPeople() {
		if n.ID == id {
			return Crew{ID: n.ID, Name: n.Name, Loyalty: n.Trust}, true
		}
	}
	return Crew{}, false
}

func (w *World) HandReadiness(hand Hand) string {
	if w.CrewOrderFor(hand.ID) != nil {
		return "They are on a headquarters assignment"
	}
	c, ok := w.handMember(hand)
	if !ok {
		return "They do not work for you"
	}
	if why := w.OutOfReach(c.ID); why != "" {
		return why
	}
	if w.personHasLegacyTask(c.ID) {
		return c.Name + " is already on assignment"
	}
	for _, place := range Locations {
		if n := w.PostedAt(place.ID); n != nil && n.ID == c.ID {
			return "Relieve " + c.Name + " from guard duty first"
		}
	}
	if c.Loyalty < HandLoyalty {
		return fmt.Sprintf("%s will not do this below %d loyalty", c.Name, HandLoyalty)
	}
	return ""
}

func (w *World) HandEdge(hand Hand) float64 {
	if !hand.Crew {
		return float64(min(w.Presence(), 100))/400 + w.WeaponEdge()
	}
	c, ok := w.handMember(hand)
	if !ok {
		return 0
	}
	edge := float64(c.Loyalty)/500 - .08
	if n := w.NPC(c.ID); n != nil {
		edge += float64(n.Weapon) * WeaponWorth
	}
	return edge
}

func (w *World) removeAssociate(id string) {
	for i, c := range w.Player.Crew {
		if c.ID == id {
			w.Player.Crew = append(w.Player.Crew[:i], w.Player.Crew[i+1:]...)
			return
		}
	}
}

// HandHurt charges the selected person's trust, never another associate's life.
func (w *World) HandHurt(hand Hand, injury int, what string) {
	if !hand.Crew {
		w.Ruin(30)
		w.Damage(15)
		w.Player.Health = max(0, w.Player.Health-w.Absorb(injury))
		return
	}
	c, ok := w.handMember(hand)
	if !ok {
		return
	}
	if injury >= 25 && w.Random() < .22 {
		if n := w.NPC(c.ID); n != nil {
			w.KillBy(n.ID, nil, fmt.Sprintf("They had gone to %s for somebody else.", what))
		} else {
			w.Log(c.Name+" did not come back", fmt.Sprintf("You sent them to %s and somebody was waiting.", what), "danger")
		}
		w.removeAssociate(c.ID)
		return
	}
	after := max(0, c.Loyalty-HandLoyaltyCost)
	associate := false
	for i := range w.Player.Crew {
		if w.Player.Crew[i].ID == c.ID {
			w.Player.Crew[i].Loyalty = after
			associate = true
			break
		}
	}
	if !associate {
		if n := w.NPC(c.ID); n != nil {
			n.Trust = after
		}
	}
	w.Log(c.Name+" took it instead of you", fmt.Sprintf("They went out to %s and came back hurt. Loyalty falls from %d to %d, and they know whose idea it was.", what, c.Loyalty, after), "danger")
}

// HandRespectFor is the standing a piece of work is worth, which is less when
// somebody else's hands were on it.
func (w *World) HandRespectFor(hand Hand, full int) int {
	if !hand.Crew {
		return full
	}
	return max(1, full*HandRespect/5)
}

// HandHeat is the attention a piece of work draws to the player. Somebody else
// standing there is somebody else being described.
func (w *World) HandHeat(hand Hand, full int) int {
	if !hand.Crew {
		return full
	}
	return max(0, full-HandHeatRelief)
}

// Older tasks carry their actor in Homes; pre-named saves used the first hire.
func (w *World) personHasLegacyTask(id string) bool {
	for _, task := range w.Homes {
		if task.Person == id {
			return true
		}
	}
	return len(w.Tasks) > 0 && len(w.Homes) == 0 && len(w.Player.Crew) > 0 && w.Player.Crew[0].ID == id
}
