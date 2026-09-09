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
	return Hand{Crew: true, Name: w.Player.Crew[0].Name}, true
}

// DelegateReadiness explains why nobody can be sent, or returns "".
func (w *World) DelegateReadiness() string {
	if len(w.Player.Crew) == 0 {
		return "You have nobody to send"
	}
	if len(w.Tasks) > 0 {
		return w.Player.Crew[0].Name + " is already on assignment"
	}
	// The jobs this gates are aimed at a place, and the man who does them is
	// named nowhere in their ids, so the sweep that asks whether a subject can
	// be reached never saw them. Ask here, where all of them pass.
	if reason := w.OutOfReach(w.Player.Crew[0].ID); reason != "" {
		return reason
	}
	if w.Player.Crew[0].Loyalty < HandLoyalty {
		return fmt.Sprintf("%s will not do this below %d loyalty", w.Player.Crew[0].Name, HandLoyalty)
	}
	return ""
}

// HandEdge is what the person doing the work brings to it. The player brings
// what they are and what they are carrying; somebody sent brings how they feel
// about being sent, and no more than that.
func (w *World) HandEdge(hand Hand) float64 {
	if !hand.Crew {
		return float64(min(w.Presence(), 100))/400 + w.WeaponEdge()
	}
	return float64(w.Player.Crew[0].Loyalty)/500 - .08
}

// HandHurt is what a job going wrong costs, and to whom. The player takes it in
// health; somebody sent takes it in loyalty, and sometimes in everything.
func (w *World) HandHurt(hand Hand, injury int, what string) {
	if !hand.Crew {
		w.Ruin(30)
		w.Damage(15)
		w.Player.Health = max(0, w.Player.Health-w.Absorb(injury))
		return
	}
	if len(w.Player.Crew) == 0 {
		return
	}
	member := &w.Player.Crew[0]
	// A bad enough night and he does not come back from it.
	if injury >= 25 && w.Random() < .22 {
		name := member.Name
		w.Player.Crew = w.Player.Crew[:0]
		if n := w.NPC(member.ID); n != nil {
			w.KillBy(n.ID, nil, fmt.Sprintf("They had gone to %s for somebody else.", what))
		} else {
			w.Log(name+" did not come back", fmt.Sprintf("You sent them to %s and somebody was waiting. There is nobody to send now.", what), "danger")
		}
		return
	}
	before := member.Loyalty
	member.Loyalty = max(0, member.Loyalty-HandLoyaltyCost)
	w.Log(member.Name+" took it instead of you", fmt.Sprintf("They went out to %s and came back hurt. Loyalty falls from %d to %d, and they know whose idea it was.", what, before, member.Loyalty), "danger")
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
