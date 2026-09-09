package core

import "fmt"

// The player can change what an organization is worth in exactly one way: by
// hurting it. There is no version of this city's politics in which two people
// who both have a problem with a third stand together, which is most of what
// politics is.
//
// A pact is that. It costs money every day, it stops the two of you moving on
// each other, it brings somebody to the door when you are raided, and it puts
// you in the middle of every quarrel they have — which is the price, and it is
// not a small one.

const (
	// PactOpening is what it takes to put one in place.
	PactOpening = 700
	// PactTribute is what holding it costs a day. Miss it and it lapses.
	PactTribute = 30
	// PactMinutes is the evening it takes to agree.
	PactMinutes = 120
	// PactGoodwill is the standing it takes before anybody will consider one.
	PactGoodwill = 15
	// PactAnswer is how often an ally turns up when the player is raided.
	PactAnswer = .45
	// PactDrag is what being known to stand with somebody costs the player
	// each day with everybody that somebody is fighting.
	PactDrag = 2
)

// Pact is an understanding between the player and an organization.
type Pact struct {
	// With is the organization, Since when, and Life which protagonist agreed
	// it: nobody inherits somebody else's friends.
	With  string `json:"with"`
	Since int    `json:"since"`
	Life  int    `json:"life"`
}

// Allied reports whether the player has an understanding with an organization.
func (w *World) Allied(id string) bool {
	for _, p := range w.Pacts {
		if p.With == id && p.Life == w.Life {
			return true
		}
	}
	return false
}

// Allies is everybody the player stands with.
func (w *World) Allies() []*Faction {
	out := []*Faction{}
	for _, p := range w.Pacts {
		if p.Life != w.Life {
			continue
		}
		if f := w.faction(p.With); f != nil {
			out = append(out, f)
		}
	}
	return out
}

// PactCost is what the understandings cost a day, which joins the rest of the
// bill.
func (w *World) PactCost() int { return len(w.Allies()) * PactTribute }

// commonEnemy is somebody both the player and an organization are at odds with,
// which is the only reason anybody agrees to one of these.
func (w *World) commonEnemy(id string) (*Faction, bool) {
	for i := range w.Factions {
		third := &w.Factions[i]
		if third.ID == id || third.ID == w.PlayerOrganizationID() {
			continue
		}
		theirs := w.Conflict(id, third.ID)
		mine := w.Conflict(w.PlayerOrganizationID(), third.ID)
		if theirs == nil || mine == nil {
			continue
		}
		if theirs.State == "cold" || mine.State == "cold" {
			continue
		}
		return third, true
	}
	return nil, false
}

// PactReadiness explains why an understanding cannot be reached, or returns "".
func (w *World) PactReadiness(id string) string {
	if !w.Incorporated() {
		return "Nobody makes an arrangement like this with one person. It takes an organization"
	}
	f := w.faction(id)
	if f == nil || f.ID == w.PlayerOrganizationID() {
		return "There is nobody of that description"
	}
	if w.Allied(id) {
		return "You already stand with them"
	}
	if f.Goodwill < PactGoodwill {
		return fmt.Sprintf("They think of you at %+d. It takes %+d before anybody would consider it", f.Goodwill, PactGoodwill)
	}
	if _, ok := w.commonEnemy(id); !ok {
		return "Neither of you has a problem the other one shares"
	}
	if w.Intelligence(id) < 2 {
		return "You do not know enough about them to be sure what you would be agreeing to"
	}
	if w.Player.Cash < PactOpening {
		return fmt.Sprintf("It takes $%d to put one in place", PactOpening)
	}
	return ""
}

// MakePact agrees an understanding.
func (w *World) MakePact(id string) error {
	if reason := w.PactReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(PactOpening); err != nil {
		return err
	}
	f := w.faction(id)
	enemy, _ := w.commonEnemy(id)
	w.Pacts = append(w.Pacts, Pact{With: id, Since: w.Minute, Life: w.Life})
	// Standing together is not standing quietly. Whoever they are fighting
	// learns that you are in it now.
	w.Antagonize(enemy.ID, w.PlayerOrganizationID(), 15)
	w.Log("An understanding with "+f.Name, fmt.Sprintf("$%d to open it and $%d a day to keep it. Neither of you moves on the other, they may answer when somebody comes for you, and %s now has a problem with you that they did not have this morning.", PactOpening, PactTribute, enemy.Name), "politics")
	w.Report("politics", upper(f.Name)+" AND "+upper(w.Player.Name),
		fmt.Sprintf("Interests associated with %s and %s are understood to have reached an accommodation. Neither would describe its terms.", f.Name, w.Player.Name))
	return nil
}

// BreakPact ends one, from either side.
func (w *World) BreakPact(id, why string) {
	if !w.Allied(id) {
		return
	}
	kept := w.Pacts[:0]
	for _, p := range w.Pacts {
		if p.With == id && p.Life == w.Life {
			continue
		}
		kept = append(kept, p)
	}
	w.Pacts = kept
	if f := w.faction(id); f != nil {
		w.Log("The understanding with "+f.Name+" is over", why, "politics")
	}
}

// PactDay is what standing with somebody costs: the tribute, and being in every
// quarrel they are in.
func (w *World) PactDay() {
	for _, f := range w.Allies() {
		if w.Player.Cash < w.DailyCost() {
			w.BreakPact(f.ID, "The tribute stopped arriving. Nobody stands with somebody who cannot pay for it.")
			continue
		}
		// Their quarrels are yours now.
		for i := range w.Factions {
			third := &w.Factions[i]
			if third.ID == f.ID || third.ID == w.PlayerOrganizationID() {
				continue
			}
			if c := w.Conflict(f.ID, third.ID); c != nil && c.State != "cold" {
				w.Antagonize(third.ID, w.PlayerOrganizationID(), PactDrag)
			}
		}
		// And neither of you moves on the other while it holds.
		if c := w.Conflict(f.ID, w.PlayerOrganizationID()); c != nil {
			c.Hostility = max(0, c.Hostility-8)
			c.State = classify(c)
		}
	}
}

// AllyAnswers reports whether somebody who stands with the player turns up when
// a raid comes, and drives it off. Their strength is what decides it.
func (w *World) AllyAnswers(attacker *Faction) (*Faction, bool) {
	for _, f := range w.Allies() {
		if f.ID == attacker.ID {
			continue
		}
		if w.WorldRandom() >= PactAnswer*float64(f.Power)/100 {
			continue
		}
		return f, true
	}
	return nil, false
}

// PactDescription is who the player stands with, for the interface.
func (w *World) PactDescription() []map[string]any {
	out := []map[string]any{}
	for _, f := range w.Allies() {
		out = append(out, map[string]any{
			"id": f.ID, "name": f.Name, "tribute": PactTribute,
			"strength": w.Intelligence(f.ID) >= 2,
			"power":    f.Power,
		})
	}
	return out
}
