package core

import "fmt"

// Two organizations that will not speak to each other will sometimes speak in
// front of somebody they both owe. That somebody can be the player, and it is
// the only thing in this game that can end a war without anybody losing one.
//
// It is also the most dangerous room in the city. One side may have come to
// settle and the other to finish it, and which is which is decided before
// anybody sits down, from what the two organizations have actually been doing
// to each other and from the temperament of the people leading them. The player
// finds out either from somebody who warned them, or when it starts.

const (
	// SitdownFee is what the room, the guarantees and the people on the doors
	// cost. It is paid whether or not anybody agrees to anything.
	SitdownFee = 220
	// SitdownMinutes is the evening it takes.
	SitdownMinutes = 180
	// SitdownStanding is the presence below which nobody would come because
	// you asked.
	SitdownStanding = 25
	// SitdownGround is the room this city has always used: somewhere neither
	// side holds, which meant the back of a bar and nowhere else.
	SitdownGround = "bar"
	// TrapHostility is the hostility above which somebody in the room is not
	// there to talk.
	TrapHostility = 62
	// SitdownWelcome is the standing below which an organization will not come
	// to a room the player arranged.
	SitdownWelcome = -25
)

// Quarrel is a pair of organizations that are not speaking, and what the player
// would be walking into.
type Quarrel struct {
	A, B *Faction
	// Trap is whether one of them came to finish it rather than settle it.
	Trap bool
	// Suspected is whether anybody warned the player.
	Suspected bool
}

// OpenQuarrel finds the worst quarrel in the city the player could stand
// between, and works out what kind of evening it would be.
func (w *World) OpenQuarrel() (Quarrel, bool) {
	var worst *Conflict
	for i := range w.Conflicts {
		c := &w.Conflicts[i]
		if c.State != "war" && c.State != "feud" {
			continue
		}
		if worst == nil || c.Hostility > worst.Hostility {
			worst = c
		}
	}
	if worst == nil {
		return Quarrel{}, false
	}
	a, b := w.faction(worst.A), w.faction(worst.B)
	if a == nil || b == nil {
		return Quarrel{}, false
	}

	// Somebody comes to finish it when the quarrel is past talking and the
	// person leading one side is not the waiting kind. Both halves have to be
	// true, so a hot-headed leader in a cooling quarrel is still just a man in
	// a bad mood.
	trap := worst.Hostility >= TrapHostility &&
		(w.leaderIs(a, "hot") || w.leaderIs(b, "hot") || w.leaderIs(a, "vain") || w.leaderIs(b, "vain"))

	return Quarrel{A: a, B: b, Trap: trap, Suspected: trap && (w.Reach() >= 2 || w.OwnGround())}, true
}

// OwnGround is whether the player is standing in a room of their own that two
// organizations would both agree to sit in.
//
// A restaurant is the other kind of neutral ground this city has. Not because
// it is neutral — it is the player's — but because a corner table with the
// plates still down is where this has always been done, and both sides know
// that a man with a dining room to lose has as much reason as they do for
// nobody to draw anything in it.
//
// It is also the restaurant's link past its own income, which it did not have.
// Two things come with holding the room: you are not renting it, so the fee is
// somebody else's problem, and your own people are on the door. Your staff
// notice who came heavy, which is the difference between walking into a trap
// and knowing about it — and information is the scarcest thing in this room.
func (w *World) OwnGround() bool {
	place, ok := PlaceByID(w.Player.Location)
	return ok && place.Kind == "restaurant" && w.Own(place.ID) && w.Staffed(place.ID)
}

// Staffed is whether a business of the player's has the hands its trade needs.
// A dining room with nobody in it is not a room either of them would sit in:
// the point of holding the meeting here is that the people on the door are
// yours, and an empty restaurant has none.
func (w *World) Staffed(id string) bool {
	trade, running := TradeOf(id)
	prop := w.Properties[id]
	return running && prop != nil && prop.Staff >= trade.Hands
}

// SitdownWhere reports whether a meeting can be held in this room at all.
func (w *World) SitdownWhere(id string) bool {
	if id == SitdownGround {
		return true
	}
	place, ok := PlaceByID(id)
	return ok && place.Kind == "restaurant" && w.Own(id) && w.Staffed(id)
}

// SitdownCost is what the room costs. Nothing, in a room of your own.
func (w *World) SitdownCost() int {
	if w.OwnGround() {
		return 0
	}
	return SitdownFee
}

// leaderIs reports whether whoever leads an organization is a particular kind
// of person.
func (w *World) leaderIs(f *Faction, temperament string) bool {
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Name == f.Leader {
			return TemperamentOf(n).ID == temperament
		}
	}
	return false
}

// SitdownReadiness explains why a meeting cannot be called, or returns "".
func (w *World) SitdownReadiness() string {
	if !w.SitdownWhere(w.Player.Location) {
		return "This is not somewhere either of them would come"
	}
	q, ok := w.OpenQuarrel()
	if !ok {
		return "Nobody in this city is quarrelling badly enough to need you"
	}
	if w.Presence() < SitdownStanding {
		return fmt.Sprintf("Neither of them would cross the street for you. You need %d presence", SitdownStanding)
	}
	// Naming who, and by how much. "One of them" cannot tell the player which
	// door is shut, and when this refusal turned up on an action the game had
	// just offered as available, it could not tell me either.
	if q.A.Goodwill < SitdownWelcome {
		return fmt.Sprintf("%s would not sit in a room you arranged. They think of you at %+d and it takes %+d", q.A.Name, q.A.Goodwill, SitdownWelcome)
	}
	if q.B.Goodwill < SitdownWelcome {
		return fmt.Sprintf("%s would not sit in a room you arranged. They think of you at %+d and it takes %+d", q.B.Name, q.B.Goodwill, SitdownWelcome)
	}
	if w.Player.Cash < w.SitdownCost() {
		return fmt.Sprintf("The room and the guarantees cost $%d", w.SitdownCost())
	}
	return ""
}

// CallSitdown arranges the meeting and puts the player in the room. What
// happens next is a decision rather than an outcome, which is why this opens a
// scene rather than resolving one.
func (w *World) CallSitdown() error {
	if reason := w.SitdownReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	q, _ := w.OpenQuarrel()
	return w.openSitdown(q)
}

// CallSitdownAs opens the meeting that was agreed to, whatever the three hours
// since did to anybody's opinion.
//
// The room takes an evening, and the clock is advanced before it opens so that
// the evening costs the evening. That left three hours between the check that
// enabled the button and the check that ran the command, and an organization
// whose standing drifted one point in between was refused after the player had
// committed. Found by cmd/apicheck, twice in twenty-two runs: "Doyle Crew would
// not sit in a room you arranged. They think of you at -26 and it takes -25."
//
// Three guesses at the cause failed to reproduce it before the refusal was made
// to name the family and the figure, at which point it said so itself. The room
// was arranged when they were willing; they came.
func (w *World) CallSitdownAs(q Quarrel) error {
	if q.A == nil || q.B == nil {
		return fmt.Errorf("there is nobody left to bring")
	}
	return w.openSitdown(q)
}

func (w *World) openSitdown(q Quarrel) error {
	if err := w.Pay(w.SitdownCost()); err != nil {
		return err
	}

	body := fmt.Sprintf("“%s and %s, in the same room, because you asked. Nobody has said anything yet.”", q.A.Name, q.B.Name)
	if q.Suspected {
		// Whoever actually told you. It said "Mara caught your eye" in every
		// campaign, including the ones in which Mara had been dead a month —
		// the same fault as the coffee bought for a fixer who no longer exists.
		who := w.RoleName("fixer")
		if w.OwnGround() {
			who = "one of your own people on the door"
		}
		body = fmt.Sprintf("“%s and %s, in the same room, because you asked. One of them brought more people than the room needs, and %s caught your eye on the way in.”", q.A.Name, q.B.Name, who)
	}
	speaker := w.HolderID("fixer")
	w.Event = &Scene{
		ID: ID(), Title: "A room nobody wanted to be in", Body: body,
		Speaker: speaker, Actor: q.A.ID, Target: q.B.ID, Kind: "sitdown",
		Source: "authored", Minute: w.Minute,
		Choices: []Choice{
			{ID: "press", Label: "Press them both to settle it tonight",
				Detail: "Ends the quarrel outright if they came to talk. If one of them did not, you are standing in the middle of it."},
			{ID: "listen", Label: "Let them talk and say nothing",
				Detail: "Cools it a little. Costs nothing and settles nothing. Safe unless somebody had already decided."},
			{ID: "side", Label: "Come down on " + q.A.Name + "'s side",
				Detail: fmt.Sprintf("Standing with %s, and %s will not forget it. Hardens the quarrel.", q.A.Name, q.B.Name)},
			{ID: "leave", Label: "Say your piece and get out",
				Detail: "You paid for the room. Both of them think less of you for calling a meeting you would not sit through."},
		},
	}
	return nil
}

// ResolveSitdown settles what happened in the room. Every branch is decided
// from state the world already holds, so nothing here is a coin toss dressed
// up as a decision.
func (w *World) ResolveSitdown(e *Scene, choice string) error {
	a, b := w.faction(e.Actor), w.faction(e.Target)
	if a == nil || b == nil {
		return fmt.Errorf("one of them is no longer in the city")
	}
	c := w.Conflict(a.ID, b.ID)
	trap := c != nil && c.Hostility >= TrapHostility &&
		(w.leaderIs(a, "hot") || w.leaderIs(b, "hot") || w.leaderIs(a, "vain") || w.leaderIs(b, "vain"))

	switch choice {
	case "leave":
		a.Goodwill = max(-100, a.Goodwill-8)
		b.Goodwill = max(-100, b.Goodwill-8)
		w.Log("You did not stay for it", "You said what you came to say and were on the street before anybody answered. Both of them noticed.", "politics")
		return nil

	case "listen":
		if c != nil {
			w.Antagonize(a.ID, b.ID, -6)
		}
		w.Player.Respect++
		w.Log("They talked", fmt.Sprintf("Nothing was agreed and nobody drew anything. %s and %s are a little further from the edge than they were.", a.Name, b.Name), "politics")
		return nil

	case "side":
		a.Goodwill = min(100, a.Goodwill+18)
		b.Goodwill = max(-100, b.Goodwill-30)
		if c != nil {
			w.Antagonize(a.ID, b.ID, 12)
		}
		w.RetaliationFrom(b.ID)
		w.Log("You took a side", fmt.Sprintf("You said it plainly and %s heard every word. %s will remember that you were in the room and whose side you were on.", b.Name, b.Leader), "politics")
		w.Report("politics", "TALKS COLLAPSE BETWEEN TWO FAMILIES",
			fmt.Sprintf("A meeting between interests associated with %s and %s is understood to have ended without agreement.", a.Name, b.Name))
		return nil

	case "press":
		if !trap {
			// Both of them came to talk, and somebody finally made them.
			if c != nil {
				c.Hostility = max(0, c.Hostility-45)
				previous := c.State
				c.State = classify(c)
				if c.State != previous {
					c.Since = w.Minute
				}
			}
			a.Goodwill = min(100, a.Goodwill+22)
			b.Goodwill = min(100, b.Goodwill+22)
			w.Player.Respect += 12
			w.Log("It is settled", fmt.Sprintf("%s and %s shook on it in front of you, which means it holds as long as you do. Both of them owe you an evening they will not talk about.", a.Name, b.Name), "politics")
			w.Report("politics", "TRUCE BETWEEN TWO FAMILIES",
				fmt.Sprintf("Sources say the dispute between interests associated with %s and %s has been settled. Neither would say who arranged it.", a.Name, b.Name))
			return nil
		}
		w.bloodbath(a, b)
		return nil
	}
	return fmt.Errorf("nothing was said")
}

// attend puts somebody in the room the meeting is in. Whoever dies at a
// sitdown came to it.
func (w *World) attend(n *NPC) {
	n.Location = SitdownGround
	n.Heading, n.Arrives, n.Errand, n.Sets = "", 0, "", 0
}

// bloodbath is the evening going the way one side always intended. The player
// is standing between them when it does.
func (w *World) bloodbath(a, b *Faction) {
	if c := w.Conflict(a.ID, b.ID); c != nil {
		w.Antagonize(a.ID, b.ID, 25)
	}
	// Both sides lose people. Whoever planned it loses fewer.
	//
	// They are put in the room before they are killed in it. A casualty is
	// drawn from anywhere in the organization, and the paper reports a death
	// where the person was standing — so the city carried "At The Monarch:
	// Ennio Zanetti was knifed in the crowd. A meeting between Bellandi Family
	// and Russo Outfit went the way somebody had already decided", naming a
	// room three streets from the one the player was sitting in.
	if victim := w.casualty(b.ID); victim != nil {
		w.attend(victim)
		w.KillBy(victim.ID, nil, fmt.Sprintf("A meeting between %s and %s went the way somebody had already decided.", a.Name, b.Name))
	}
	if w.WorldRandom() < .5 {
		if victim := w.casualty(a.ID); victim != nil {
			w.attend(victim)
			w.KillBy(victim.ID, nil, "The same room, the same evening.")
		}
	}

	// The player is in the middle of it. Armour and a weapon are the difference
	// between a bad night and the last one.
	injury := w.Absorb(35 + int(w.Random()*40))
	w.Ruin(80)
	w.Player.Health = max(0, w.Player.Health-injury)
	w.Player.Heat = min(100, w.Player.Heat+18)
	a.Goodwill = max(-100, a.Goodwill-15)
	b.Goodwill = max(-100, b.Goodwill-15)

	// So is whoever came with you.
	lost := ""
	if len(w.Player.Crew) > 0 && w.Random() < .3 {
		lost = w.Player.Crew[0].Name
		w.Player.Crew = w.Player.Crew[:0]
	}

	body := fmt.Sprintf("Somebody stood up before anybody had finished a sentence. You went out through the kitchen with %d less health than you came in with.", injury)
	if lost != "" {
		body = fmt.Sprintf("Somebody stood up before anybody had finished a sentence. You went out through the kitchen with %d less health than you came in with. %s did not come out at all.", injury, lost)
	}
	w.Log("It was never a meeting", body, "danger")
	w.Report("killing", "SHOOTING AT SAINT AGNES",
		fmt.Sprintf("Several people were shot at a bar on the waterfront. Police believe a meeting between interests associated with %s and %s was the occasion.", a.Name, b.Name))
	w.Witness("gunfight", SitdownGround, "A meeting at Saint Agnes ended in gunfire.", "")
	if w.Player.Health <= 0 {
		w.Die("A meeting at Saint Agnes that was never a meeting.")
	}
}
