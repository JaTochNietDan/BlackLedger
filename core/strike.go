package core

import "fmt"

// Going after somebody. Until now the only violence the player could aim at a
// person was a price on a name and a wait: you paid, and somebody you never met
// did or did not finish it. There was no version where you went yourself, and
// no version where you sent one of your own — which is most of what this city
// is actually about.
//
// The two are deliberately different bets. Going yourself is the best odds you
// can buy, because you bring what you are carrying, and it is the only one that
// can get you killed. Sending somebody puts a person between you and it: their
// face is the one that is seen, their loyalty is what decides it, and if they
// are taken alive they are yours and everybody knows it.

const (
	// StrikeMinutes is what an attempt costs in time, whoever makes it.
	StrikeMinutes = 60
	// StrikeHeat is what the police make of a killing in a room with people in
	// it, and StrikeQuiet what they make of one nobody saw.
	StrikeHeat  = 26
	StrikeQuiet = 12
	// SentHeat is what it costs when it was not your face at the scene.
	SentHeat = 9
	// StrikeRespect is what the street makes of somebody who does it in person.
	StrikeRespect = 9
	// HandCaught is how often somebody of yours who fails is taken alive rather
	// than getting away or being killed on the spot. Taken alive is the worst
	// of the three for you: they are known to be yours.
	HandCaught = .4
	// HandDies is how often a failure kills them outright.
	HandDies = .3
	// StrikeSore is what somebody who survives an attempt holds against whoever
	// sent it.
	StrikeSore = 60
)

// StrikeTarget is somebody the player could go after where they are standing:
// alive, here, and not one of their own.
func (w *World) StrikeTarget(id string) (*NPC, bool) {
	n := w.NPC(id)
	if n == nil || n.Dead || n.Location != w.Player.Location {
		return nil, false
	}
	if n.Faction != "" && n.Faction == w.PlayerOrganizationID() {
		return nil, false
	}
	for _, c := range w.Player.Crew {
		if c.ID == n.ID {
			return nil, false
		}
	}
	return n, true
}

// StrikeReadiness explains why the player cannot go after somebody themselves,
// or returns "".
func (w *World) StrikeReadiness(id string) string {
	if w.Player.HeldUntil > w.Minute {
		return "You are not going anywhere tonight"
	}
	n, ok := w.StrikeTarget(id)
	if !ok {
		if who := w.NPC(id); who != nil && who.Faction != "" && who.Faction == w.PlayerOrganizationID() {
			return "They are one of yours. Putting somebody out is what that is for"
		}
		return "They are not here"
	}
	if w.Player.Health < 40 {
		return "In this condition you would not finish it"
	}
	if n.Rank >= RankLeader && w.Presence() < 20 {
		return "Nobody who has never been anybody gets that close to " + n.Name
	}
	return ""
}

// SendReadiness explains why nobody of the player's can be sent after somebody,
// or returns "".
func (w *World) SendReadiness(id string) string {
	if _, ok := w.StrikeTarget(id); !ok {
		return "They are not here"
	}
	if _, ok := w.CrewHands(); !ok {
		return "You have nobody to send"
	}
	if reason := w.DelegateReadiness(); reason != "" {
		return reason
	}
	return ""
}

// strikeChance is how likely an attempt is to finish. Whoever makes it brings
// what they have — the player brings a weapon and a reputation, somebody sent
// brings how they feel about being sent — and whoever it is aimed at brings
// their standing and whatever is around them.
func (w *World) strikeChance(n *NPC, hand Hand) float64 {
	chance := .5 + w.HandEdge(hand)
	chance -= float64(n.Rank) / 260
	if f := w.faction(n.Faction); f != nil {
		chance -= float64(f.Power) / 500
	}
	// A room with people in it is a room with people watching.
	if crowd := len(w.PeopleHere(n.Location)) - 1; crowd > 0 {
		chance -= float64(min(crowd, 6)) * .02
	}
	return min64(.9, max64(.08, chance))
}

// Strike is an attempt on somebody's life, made by the player or by one of
// theirs. Nothing about it is decided anywhere else: it resolves here, before
// the clock moves, so a fatal one cannot also collect the hour it never
// survived.
func (w *World) Strike(id string, hand Hand) error {
	reason := w.StrikeReadiness(id)
	if hand.Crew {
		reason = w.SendReadiness(id)
	}
	if reason != "" {
		return fmt.Errorf("%s", reason)
	}
	n, _ := w.StrikeTarget(id)
	place, _ := PlaceByID(n.Location)
	crowd := len(w.PeopleHere(n.Location)) - 1
	family := w.faction(n.Faction)

	if w.Random() < w.strikeChance(n, hand) {
		w.finishThem(n, hand, place.Name, crowd, family)
		return nil
	}
	w.itWentWrong(n, hand, place.Name, family)
	return nil
}

// finishThem is an attempt that works.
func (w *World) finishThem(n *NPC, hand Hand, where string, crowd int, family *Faction) {
	name := n.Name
	attacker := w.strikeAttacker(hand)
	w.Kill(n.ID, w.strikeManner(n, attacker.Weapon))
	heat := StrikeQuiet
	if crowd > 0 {
		heat = StrikeHeat
	}
	if hand.Crew {
		heat = SentHeat
	}
	w.Player.Heat = min(100, w.Player.Heat+heat)
	if !hand.Crew {
		w.Player.Respect += StrikeRespect
	}
	kind := "gunfight"
	if attacker.Weapon == 0 {
		kind = "attack"
	}
	w.Witness(kind, n.Location, name+" was killed at "+where+".", "")
	if len(w.VisualCues) > 0 {
		w.VisualCues[len(w.VisualCues)-1].Attacker = &attacker
	}
	w.ReportAbout("killing", "KILLING AT "+upper(where),
		w.unattributed(where, fmt.Sprintf("%s was found dead at %s. Police have no arrest and describe the killing as targeted.", name, where)), n.ID)
	if family != nil {
		family.Goodwill = max(-100, family.Goodwill-40)
		w.RetaliationFrom(family.ID)
	}
	if hand.Crew {
		w.Log(w.Player.Crew[0].Name+" finished it", fmt.Sprintf("%s is dead at %s. It was not your face anybody saw.", name, where), "danger")
		return
	}
	w.Log("You finished it", fmt.Sprintf("%s is dead at %s. %s", name, where,
		map[bool]string{true: "There were people in the room, and rooms talk.", false: "Nobody was close enough to say who did it."}[crowd > 0]), "danger")
}

// itWentWrong is an attempt that does not work, which is where the difference
// between going and sending actually lives.
func (w *World) itWentWrong(n *NPC, hand Hand, where string, family *Faction) {
	w.Aggrieve(n.ID, StrikeSore, "the night somebody came for them")
	if family != nil {
		family.Goodwill = max(-100, family.Goodwill-25)
		w.RetaliationFrom(family.ID)
	}
	w.ReportAbout("attempt", "ATTEMPT ON THE LIFE OF "+upper(n.Name),
		w.unattributed(where, fmt.Sprintf("%s survived an attack at %s. Police say the assault was targeted and no arrest has been made.", n.Name, where)), n.ID)

	if !hand.Crew {
		// You were there, and it can be the end of you.
		w.Player.Heat = min(100, w.Player.Heat+StrikeHeat)
		if w.Random() < .22 {
			w.DieOf("an attempt of your own", "An attempt on "+n.Name+" at "+where+" that went the other way.")
			return
		}
		w.HandHurt(hand, 55, "")
		w.Log("It went wrong", fmt.Sprintf("%s is alive, knows your face, and you came out of %s worse than you went in.", n.Name, where), "danger")
		return
	}

	// Somebody of yours. Three ways it ends for them, and the worst one for you
	// is the one where they are still breathing.
	who := w.Player.Crew[0]
	w.Player.Heat = min(100, w.Player.Heat+SentHeat)
	// What was waiting at the kerb. A man with something running gets off the
	// street; a man on foot is still on it when the doors open. This is the
	// whole of what buying one of your own a car is for.
	roll := w.Random() + w.GetsOut(who.ID)
	switch {
	case roll < HandDies:
		w.Kill(who.ID, fmt.Sprintf("Shot at %s, going for %s on somebody else's word.", where, n.Name))
		w.Player.Crew = nil
		w.Log("They did not come back", fmt.Sprintf("%s went for %s at %s and did not walk out of it.", who.Name, n.Name, where), "danger")
	case roll < HandDies+HandCaught:
		// Taken alive. Their face is known and so is whose face it is.
		w.Player.Heat = min(100, w.Player.Heat+25)
		w.hold(w.NPC(who.ID), 3)
		w.Player.Crew = nil
		w.ReportAbout("arrest", "ARREST AFTER ATTACK ON "+upper(n.Name),
			"Somebody taken at the scene is assisting police with their enquiries. Sources suggest a name has been given.", n.ID)
		if family != nil {
			family.Goodwill = max(-100, family.Goodwill-45)
			w.RetaliationFrom(family.ID)
			w.Log("They gave up a name", fmt.Sprintf("%s was taken alive at %s. %s %s who sent them, and so do the police.",
				who.Name, where, Leads(family.Name), Agree(family.Name, "knows", "know")), "danger")
			return
		}
		w.Log("They gave up a name", fmt.Sprintf("%s was taken alive at %s and questioned. Your name came out of it.", who.Name, where), "danger")
	default:
		w.HandHurt(hand, 40, "")
		w.Log("They got away with nothing", fmt.Sprintf("%s went for %s at %s, did not finish it, and got out. %s is alive and looking.",
			who.Name, n.Name, where, n.Name), "danger")
	}
}

// factionOrStreet names who answers for somebody: their organization if they
// have one, and the street if they do not — because somebody with nobody behind
// them is not somebody nobody misses.
func (w *World) factionOrStreet(n *NPC) string {
	if f := w.faction(n.Faction); f != nil {
		return f.Name
	}
	return "the street"
}

// theirOrTheir keeps the sentence right whether the name has a family or not.
func theirOrTheir(n *NPC) string { return n.Name }
