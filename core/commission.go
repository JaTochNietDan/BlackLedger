package core

import "fmt"

// Everything the player could do until now was either an errand that resolved
// the moment they committed to it, or a business decision that paid out by the
// hour. Nothing in the game ever asked them for something and waited.
//
// A commission is a standing piece of work an organization wants done, with a
// deadline, a reward and a cost for failing. It is deliberately not a script:
// the objective is a condition the world can check for itself, so it can be met
// any way the player can think of. "Take two thousand out of Bellandi" does not
// say how, and the game does not care how.

// Objective kinds. Each one is something the world can look at and answer, so
// nothing here needs to watch the player or guess at their intent.
const (
	// Deliver: be carrying a quantity of a good at a place.
	ObjectiveDeliver = "deliver"
	// Damage: put a named property below a condition.
	ObjectiveDamage = "damage"
	// Drain: take an amount out of an organization's money, from wherever it
	// stood when the work was accepted.
	ObjectiveDrain = "drain"
	// Remove: a named person is no longer alive.
	ObjectiveRemove = "remove"
	// Standing: be worth a certain amount to the city.
	ObjectiveStanding = "standing"
)

// Commission is work an organization wants done.
type Commission struct {
	ID string `json:"id"`
	// Patron is the organization that wants it, and Giver the person who asked.
	PatronID  string `json:"patron"`
	GiverName string `json:"giver"`
	// Kind is one of the objectives above; Target names what it is about and
	// Amount how much of it is wanted.
	Kind   string `json:"kind"`
	Target string `json:"target,omitempty"`
	Amount int    `json:"amount"`
	// GoodID is what a delivery is of, held separately from Target, which is
	// where it has to be.
	GoodID string `json:"good,omitempty"`
	// Baseline is what the measured thing stood at when the work was accepted,
	// so a drain is what the player took rather than what somebody else did.
	Baseline int `json:"baseline,omitempty"`
	// Brief is the work in the words it was asked in.
	Brief string `json:"brief"`
	// Due is the minute it must be done by; Life ties it to one protagonist.
	Due  int `json:"due"`
	Life int `json:"life"`
	// Pay is the money, Respect the standing, Goodwill what it is worth with
	// the patron, and Penalty what failing costs with them.
	Pay      int `json:"pay"`
	Respect  int `json:"respect"`
	Goodwill int `json:"goodwill"`
	Penalty  int `json:"penalty"`
	// Done and Failed are set once, so a commission is settled exactly once.
	Done   bool `json:"done,omitempty"`
	Failed bool `json:"failed,omitempty"`
}

// CommissionWindow is how long an organization will wait. Three days is long
// enough to plan and short enough that accepting two at once is a decision.
const CommissionWindow = 4320

// Live is the work the current protagonist has outstanding.
func (w *World) Live() []Commission {
	out := []Commission{}
	for _, c := range w.Commissions {
		if c.Life == w.Life && !c.Done && !c.Failed {
			out = append(out, c)
		}
	}
	return out
}

// Met reports whether a commission's objective currently holds. This is the
// whole of the design: the world answers a question about itself, and how the
// player got it there is their business.
func (w *World) Met(c Commission) bool {
	switch c.Kind {
	case ObjectiveDeliver:
		return w.Player.Location == c.Target && w.Holding(c.GoodID) >= c.Amount
	case ObjectiveDamage:
		prop := w.Properties[c.Target]
		return prop != nil && prop.Condition <= c.Amount
	case ObjectiveDrain:
		f := w.faction(c.Target)
		return f != nil && c.Baseline-f.Cash >= c.Amount
	case ObjectiveRemove:
		for i := range w.NPCs {
			if w.NPCs[i].ID == c.Target {
				return w.NPCs[i].Dead
			}
		}
		return true // somebody who is not in the city any more is not a problem
	case ObjectiveStanding:
		return w.Presence() >= c.Amount
	}
	return false
}

// SettleCommissions pays for what is done and charges for what was not. Called
// from the clock, so a deadline passes whether or not the player is looking.
func (w *World) SettleCommissions() {
	for i := range w.Commissions {
		c := &w.Commissions[i]
		if c.Life != w.Life || c.Done || c.Failed {
			continue
		}
		if w.Met(*c) {
			c.Done = true
			// A delivery is a delivery: the crates change hands. What was paid
			// for them is the player's problem, which is what makes the market
			// price part of whether the work was worth taking.
			if c.Kind == ObjectiveDeliver && w.Player.Stock != nil {
				w.Player.Stock[c.GoodID] -= c.Amount
			}
			w.Earn(c.Pay)
			w.Player.Respect += c.Respect
			w.ServeWork(c.PatronID)
			if f := w.faction(c.PatronID); f != nil {
				f.Goodwill = min(100, f.Goodwill+c.Goodwill)
				w.Log("Settled with "+f.Name, fmt.Sprintf("%s asked and it was done. $%d, %d respect, and %s thinks better of you by %d.", c.GiverName, c.Pay, c.Respect, f.Name, c.Goodwill), "politics")
			}
			continue
		}
		if w.Minute >= c.Due {
			c.Failed = true
			if f := w.faction(c.PatronID); f != nil {
				f.Goodwill = max(-100, f.Goodwill-c.Penalty)
				w.Log("Nothing came of it", fmt.Sprintf("%s asked for something and the time for it has passed. %s thinks worse of you by %d.", c.GiverName, f.Name, c.Penalty), "danger")
			}
		}
	}
}

// Progress is how far along a commission is, in the terms it was set in, so the
// interface can show something honest without knowing anything about the rules.
func (w *World) Progress(c Commission) string {
	switch c.Kind {
	case ObjectiveDeliver:
		return fmt.Sprintf("%d of %d units in hand", w.Holding(c.GoodID), c.Amount)
	case ObjectiveDamage:
		if prop := w.Properties[c.Target]; prop != nil {
			place, _ := PlaceByID(c.Target)
			return fmt.Sprintf("%s is at %d%%, wanted at %d%% or worse", place.Name, prop.Condition, c.Amount)
		}
	case ObjectiveDrain:
		if f := w.faction(c.Target); f != nil {
			return fmt.Sprintf("$%d of $%d taken out of %s", max(0, c.Baseline-f.Cash), c.Amount, f.Name)
		}
	case ObjectiveRemove:
		return "Still walking around"
	case ObjectiveStanding:
		return fmt.Sprintf("%d presence of %d", w.Presence(), c.Amount)
	}
	return ""
}

// PublicCommissions is the outstanding work, for the interface, with how far
// along each one is and how long is left.
func (w *World) PublicCommissions() []map[string]any {
	out := []map[string]any{}
	for _, c := range w.Live() {
		out = append(out, map[string]any{
			"id": c.ID, "brief": c.Brief, "giver": c.GiverName,
			"patron":   w.factionName(c.PatronID),
			"progress": w.Progress(c), "met": w.Met(c),
			"pay": c.Pay, "respect": c.Respect,
			"due": c.Due, "minutes_left": max(0, c.Due-w.Minute),
		})
	}
	return out
}

// MaxCommissions is how many pieces of work the player can carry at once. More
// than this and a deadline stops meaning anything.
const MaxCommissions = 3

// commissionGiver finds somebody at this place who could ask for something on
// behalf of an organization that would deal with the player at all.
func (w *World) commissionGiver(location string) *NPC {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != location || n.Faction == "" || n.Rank < RankLieutenant {
			continue
		}
		f := w.faction(n.Faction)
		if f == nil || f.Goodwill < -20 {
			continue // people who want you dead do not hand you work
		}
		if w.hasLiveCommissionFrom(n.Faction) {
			continue
		}
		return n
	}
	return nil
}

func (w *World) hasLiveCommissionFrom(faction string) bool {
	for _, c := range w.Live() {
		if c.PatronID == faction {
			return true
		}
	}
	return false
}

// AvailableCommission is the work on offer at a place, derived entirely from
// what the organization asking is currently living through. A family at war
// wants its rival hurt; one that is short wants stock moved; one that has lost
// ground wants the premises that took it damaged.
func (w *World) AvailableCommission(location string) (Commission, bool) {
	giver := w.commissionGiver(location)
	if giver == nil || len(w.Live()) >= MaxCommissions {
		return Commission{}, false
	}
	f := w.faction(giver.Faction)
	c := Commission{
		ID: ID(), PatronID: f.ID, GiverName: giver.Name,
		Due: w.Minute + CommissionWindow, Life: w.Life,
	}

	// Somebody they are actually fighting is the first thing on their mind. A
	// family with no quarrel has other problems, which is what stops every
	// organization in the city asking for the same thing.
	if rival := w.fighting(f.ID); rival != nil {
		// What they ask for depends on where the rival is strong. A rival with
		// more money than them is a rival who can buy people; one whose
		// strength is in the ground they hold is one whose ground should go.
		// At war, and only at war, they will ask for a person. It is the
		// heaviest thing anybody asks the player for and it pays like it.
		if w.atWar(f.ID, rival.ID) {
			if mark := w.expendable(rival.ID); mark != nil {
				c.Kind, c.Target = ObjectiveRemove, mark.ID
				c.Brief = fmt.Sprintf("%s of %s is not to see the end of the week. %s did not say this to you and will deny it if asked.", mark.Name, rival.Name, Leads(f.Name))
				c.Pay, c.Respect, c.Goodwill, c.Penalty = 900, 10, 22, 14
				return c, true
			}
		}
		if rival.Cash > f.Cash {
			c.Kind, c.Target, c.Amount, c.Baseline = ObjectiveDrain, rival.ID, 900, rival.Cash
			c.Brief = fmt.Sprintf("Take nine hundred dollars out of %s. %s %s care how it leaves their hands, only that it does.", rival.Name, Leads(f.Name), Agree(f.Name, "does not", "do not"))
			c.Pay, c.Respect, c.Goodwill, c.Penalty = 500, 6, 16, 10
			return c, true
		}
		// Ground the rival holds and could lose is the most useful thing to
		// take off them.
		if holdings := w.FamilyHoldings(rival.ID); len(holdings) > 0 {
			target := holdings[0]
			place, _ := PlaceByID(target)
			prop := w.Properties[target]
			if prop.Condition > 45 {
				c.Kind, c.Target, c.Amount = ObjectiveDamage, target, 45
				c.Brief = fmt.Sprintf("Put %s out of the state it is in. %s %s it at forty-five percent or worse, and %s want to be seen doing it.", place.Name, Leads(f.Name), Agree(f.Name, "wants", "want"), Agree(f.Name, "does not", "do not"))
				c.Pay, c.Respect, c.Goodwill, c.Penalty = 420, 5, 14, 8
				return c, true
			}
		}
		c.Kind, c.Target, c.Amount, c.Baseline = ObjectiveDrain, rival.ID, 900, rival.Cash
		c.Brief = fmt.Sprintf("Take nine hundred dollars out of %s. %s %s care how it leaves their hands, only that it does.", rival.Name, Leads(f.Name), Agree(f.Name, "does not", "do not"))
		c.Pay, c.Respect, c.Goodwill, c.Penalty = 500, 6, 16, 10
		return c, true
	}

	// An organization with less money than its neighbours wants stock it can
	// sell, brought to a place it holds.
	if richest := w.richest(f.ID); richest != nil && f.Cash < richest.Cash {
		where := w.homeOf(f.ID)
		if holdings := w.FamilyHoldings(f.ID); len(holdings) > 0 {
			where = holdings[0]
		}
		place, _ := PlaceByID(where)
		c.Kind, c.Target, c.GoodID, c.Amount = ObjectiveDeliver, where, "moonshine", 25
		c.Brief = fmt.Sprintf("Bring twenty-five crates of moonshine to %s and hand them over. %s %s short and would rather owe you than owe anybody else.", place.Name, Leads(f.Name), Agree(f.Name, "is", "are"))
		c.Pay, c.Respect, c.Goodwill, c.Penalty = 1150, 4, 12, 8
		return c, true
	}

	// A quiet, solvent organization wants somebody worth being seen with.
	target := ((w.Presence()/10)+2)*10 + 10
	c.Kind, c.Amount = ObjectiveStanding, target
	c.Brief = fmt.Sprintf("Be worth %d to this city. %s will not be seen doing business with somebody nobody has heard of, and would like to do business with you.", target, Leads(f.Name))
	c.Pay, c.Respect, c.Goodwill, c.Penalty = 300, 0, 18, 6
	return c, true
}

// CommissionReadiness explains why work cannot be taken on here, or returns "".
func (w *World) CommissionReadiness(location string) string {
	if len(w.Live()) >= MaxCommissions {
		return fmt.Sprintf("You are already carrying %d pieces of work", MaxCommissions)
	}
	if _, ok := w.AvailableCommission(location); !ok {
		return "Nobody here is asking for anything"
	}
	return ""
}

// TakeCommission accepts the work on offer here. The baseline is recorded now,
// so what somebody else takes out of a rival's pocket is not credited to the
// player.
func (w *World) TakeCommission(location string) error {
	if reason := w.CommissionReadiness(location); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	c, _ := w.AvailableCommission(location)
	if giver := w.commissionGiver(location); giver != nil {
		w.MeetPerson(giver.ID)
	}
	w.Commissions = append(w.Commissions, c)
	w.Log(c.GiverName+" asks for something", fmt.Sprintf("%s $%d and %d standing with %s if it is done inside three days. Failing costs %d.", c.Brief, c.Pay, c.Goodwill, w.factionName(c.PatronID), c.Penalty), "politics")
	return nil
}

// fighting is the organization this one is at war or in a feud with, rather
// than merely the strongest of its neighbours.
func (w *World) fighting(id string) *Faction {
	for _, c := range w.Conflicts {
		if c.State != "war" && c.State != "feud" {
			continue
		}
		other := ""
		switch id {
		case c.A:
			other = c.B
		case c.B:
			other = c.A
		default:
			continue
		}
		if f := w.faction(other); f != nil && len(w.FamilyHoldings(f.ID)) > 0 {
			return f
		}
	}
	return nil
}

// richest is the organization other than this one with the most money, which is
// who a family compares itself to when it is deciding whether it is short.
func (w *World) richest(id string) *Faction {
	var best *Faction
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.ID == id {
			continue
		}
		if best == nil || f.Cash > best.Cash {
			best = f
		}
	}
	return best
}

// atWar reports whether two organizations are actually at war rather than
// merely quarrelling.
func (w *World) atWar(a, b string) bool {
	c := w.Conflict(a, b)
	return c != nil && c.State == "war"
}

// expendable is somebody in an organization whose death would be felt without
// ending it: not the person at the top, and somebody who is actually there.
func (w *World) expendable(faction string) *NPC {
	for _, n := range w.Members(faction) {
		if !n.Dead && n.Rank < RankLeader {
			return n
		}
	}
	return nil
}
