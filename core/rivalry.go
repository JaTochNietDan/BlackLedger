package core

import (
	"fmt"
	"strings"
)

// Families hold property, and those holdings are the only thing a rival can
// actually take from them. Power used to be an inert display number; it is now
// a standing that damage reduces and time restores, so pressure runs in both
// directions instead of only from the families toward the player.

// FamilyHoldings lists the properties an established family owns, in the fixed
// order of Locations so reports and tests do not depend on map iteration.
func (w *World) FamilyHoldings(faction string) []string {
	out := []string{}
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Owner == faction {
			out = append(out, l.ID)
		}
	}
	return out
}

func (w *World) faction(id string) *Faction {
	for i := range w.Factions {
		if w.Factions[i].ID == id {
			return &w.Factions[i]
		}
	}
	return nil
}

// peak is the strength a family recovers toward once its holdings are whole.
// Saves written before families had holdings carry no peak; their current power
// is the right ceiling for them.
func peak(f *Faction) int {
	if f.Peak > 0 {
		return f.Peak
	}
	return f.Power
}

// FamilyDay runs once per game day. Families collect from their holdings,
// repair damage at their own pace, and drift toward the strength their
// remaining property can support.
func (w *World) FamilyDay() {
	for i := range w.Factions {
		f := &w.Factions[i]
		// The player's own organization is run by the player: its premises are
		// repaired when they pay for it, its money is their money, and its
		// strength is whatever their holdings and name are worth today. Letting
		// this loop have it would have quietly repaired their businesses for
		// free every morning.
		if f.ID == w.PlayerOrganizationID() {
			continue
		}
		income, condition, count := 0, 0, 0
		damaged := []*Property{}
		for _, id := range w.FamilyHoldings(f.ID) {
			prop := w.Properties[id]
			count++
			condition += prop.Condition
			income += prop.Income * prop.Condition / 100
			if prop.Condition < 100 {
				damaged = append(damaged, prop)
			}
		}
		// What the day costs them. Everybody who answers to a family is paid
		// what the game already says answering to somebody pays, and the
		// repairs that used to happen for free every morning are now bought.
		// Without this a family only ever got richer: no organization in the
		// city had an outgoing of any kind, so after a year the poorest of
		// them held sixty thousand dollars and none had ever been short.
		// Power is the game's own measure of how big an operation is, and a
		// family of ninety runs far more people than the handful who happen to
		// have names. Paying only the named ones would have made a family of
		// ninety and a family of ten cost the same to run.
		wages := w.FamilyBill(f)
		f.Cash += income * 24
		f.Cash -= wages
		if f.Cash < 0 {
			// The bill is met as far as it goes and the rest is owed to the
			// people, who notice. Money never goes below nothing, and the
			// count is of consecutive days, so a family that misses one week
			// reads differently from one that missed a payday once.
			f.Short++
			f.Cash = 0
		} else {
			f.Short = 0
		}
		// Repairs are what a family does with what is left, in the order the
		// holdings come, and a family with nothing spare watches its property
		// go down instead.
		for _, prop := range damaged {
			cost := min(FamilyRepair, 100-prop.Condition) * FamilyRepairCost
			if f.Cash < cost {
				continue
			}
			f.Cash -= cost
			prop.Condition = min(100, prop.Condition+FamilyRepair)
		}
		// Strength comes from holdings. An organization that holds nothing has
		// nothing to draw on and fades, rather than recovering to the strength
		// it had when it still owned half the waterfront.
		//
		// It used to fade on a timer, four points a day, and that ran faster
		// than its money did: a family stripped of every business shrank to
		// nothing without ever once failing to pay anybody. Starving a family
		// was not a way to beat it, because the fading was not caused by the
		// money at all. Now an empty family is worth nothing to work toward and
		// gets there the same way everybody else does — by running out.
		// A family with nothing left to earn from holds together for exactly as
		// long as its money does. That is the whole of the starvation story: it
		// is not the losing of the businesses that breaks them, it is the
		// payday they cannot meet afterwards.
		target := f.Power
		if count > 0 {
			target = peak(f) * condition / (100 * count)
		}
		// People who are not paid do not stay to be told why. A family that
		// cannot meet its wages shrinks toward what it can actually afford,
		// which is the floor its holdings will carry, and the longer it goes on
		// the faster they go.
		leaving := 2
		if f.Short > 0 {
			target = min(target, income*24/max(1, FamilyWage))
			leaving = 2 + min(f.Short, MostWhoLeaveAtOnce)
		}
		switch {
		case f.Power < target:
			f.Power = min(target, f.Power+2)
		case f.Power > target:
			f.Power = max(target, f.Power-leaving)
		}
	}
}

// SabotageTarget reports whether a property is a rival holding the player could
// attack, and which family owns it.
func (w *World) SabotageTarget(id string) (*Faction, bool) {
	prop := w.Properties[id]
	if prop == nil || w.Own(id) {
		return nil, false
	}
	f := w.faction(prop.Owner)
	return f, f != nil
}

// SabotageReadiness explains why an attack cannot be attempted, or returns "".
// The same reasons gate the offered action and the committed command, so the
// interface can never present an attack the rules would refuse.
func (w *World) SabotageReadiness(id string) string {
	if _, ok := w.SabotageTarget(id); !ok {
		return "No rival organization holds this property"
	}
	if w.Player.Respect < 12 {
		return "Earn 12 respect first"
	}
	if len(w.Player.Crew) == 0 {
		return "Recruit crew before moving against a family"
	}
	if w.Player.Crew[0].Loyalty < 40 {
		return w.Player.Crew[0].Name + " is not loyal enough for this"
	}
	if len(w.Tasks) > 0 {
		return w.Player.Crew[0].Name + " is already on assignment"
	}
	return ""
}

// sabotageChance is the probability the attack lands. A stronger family is
// harder to reach; a loyal crew and a known name help.
func (w *World) sabotageChance(f *Faction, hand Hand) float64 {
	chance := .35 + w.HandEdge(hand) - float64(f.Power)/300
	if !hand.Crew {
		chance += float64(w.Player.Crew[0].Loyalty) / 400
	}
	if chance > .85 {
		chance = .85
	}
	if chance < .10 {
		chance = .10
	}
	return chance
}

// Sabotage resolves a player attack on a rival holding. It commits a result
// immediately; nothing here is left for presentation to decide.
func (w *World) Sabotage(id string) error { return w.SabotageBy(id, w.OwnHands()) }

// SabotageBy is the same attack whoever carries it out.
func (w *World) SabotageBy(id string, hand Hand) error {
	if reason := w.SabotageReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if hand.Crew {
		if reason := w.DelegateReadiness(); reason != "" {
			return fmt.Errorf("%s", reason)
		}
	}
	f, _ := w.SabotageTarget(id)
	place, ok := PlaceByID(id)
	if !ok {
		return fmt.Errorf("unknown property")
	}
	prop := w.Properties[id]
	crew := w.Player.Crew[0].Name

	if w.Random() >= w.sabotageChance(f, hand) {
		// Turned away. The family learns who came for them either way.
		injury := 12 + int(w.Random()*18)
		health := w.Player.Health
		w.HandHurt(hand, injury, "move against "+place.Name)
		w.Player.Heat = min(100, w.Player.Heat+w.HandHeat(hand, 12))
		f.Goodwill = max(-100, f.Goodwill-20)
		if hand.Crew {
			w.Log("Turned away at "+place.Name, fmt.Sprintf("%s had people waiting. %s went in without you and came back with nothing. %s knows who sent them.", f.Name, hand.Name, f.Leader), "danger")
		} else {
			w.Log("Turned away at "+place.Name, fmt.Sprintf("%s had people waiting. You and %s left without reaching anything, and you were hurt (-%d health). %s knows who came.", f.Name, crew, health-w.Player.Health, f.Leader), "danger")
		}
		w.RetaliationFrom(f.ID)
		if w.Player.Health <= 0 {
			w.Die("An attack on " + place.Name + " went wrong.")
			return nil
		}
		return nil
	}

	damage := min(prop.Condition, 25+int(w.Random()*20))
	prop.Condition -= damage
	lostPower := max(2, damage/5)
	f.Power = max(10, f.Power-lostPower)
	f.Cash = max(0, f.Cash-damage*RaidTakes)
	f.Goodwill = max(-100, f.Goodwill-30)
	w.Player.Respect += w.HandRespectFor(hand, 4)
	// Harm done to somebody's enemy is work done for them.
	w.ServeAgainst(f.ID)
	w.Player.Heat = min(100, w.Player.Heat+w.HandHeat(hand, 8))
	w.Witness("attack", id, fmt.Sprintf("Your crew damaged %s. Condition is now %d%%.", place.Name, prop.Condition), "")
	w.Report("attack", "DAMAGE AT "+strings.ToUpper(place.Name),
		w.unattributed(place.Name, fmt.Sprintf("%s, an establishment associated with %s, was attacked overnight.", place.Name, f.Name)))
	w.Log("A message to "+f.Name, fmt.Sprintf("You and %s damaged %s by %d condition. %s standing falls to %+d and their strength to %d. They will answer this.", crew, place.Name, damage, f.Name, f.Goodwill, f.Power), "politics")
	w.RetaliationFrom(f.ID)
	return nil
}

// InciteReadiness explains why a rivalry cannot be stoked, or returns "".
// Working two organizations against each other needs contacts who will carry a
// story and a rival worth pointing at.
func (w *World) InciteReadiness(id string) string {
	f, ok := w.SabotageTarget(id)
	if !ok {
		return "No rival organization holds this property"
	}
	if w.Rival(f.ID) == nil {
		return "There is no other organization to point them at"
	}
	if w.Reach() < 2 {
		return "Build an information network first"
	}
	if w.Player.Respect < 8 {
		return "Earn 8 respect first"
	}
	return ""
}

// Incite spends money and standing to make one organization believe another
// moved against it. It commits a real change to their quarrel; whether that
// becomes a war is decided by the same rules that govern every other feud.
func (w *World) Incite(id string) error {
	if reason := w.InciteReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	f, _ := w.SabotageTarget(id)
	other := w.Rival(f.ID)
	place, _ := PlaceByID(id)

	// A story that does not hold up comes back to the person who told it.
	if w.Random() < .25 {
		f.Goodwill = max(-100, f.Goodwill-15)
		w.Player.Heat = min(100, w.Player.Heat+6)
		w.Log("A story that did not hold", fmt.Sprintf("Your word against %s did not survive scrutiny at %s. %s knows where it came from.", other.Name, place.Name, f.Name), "danger")
		return nil
	}
	w.Antagonize(f.ID, other.ID, 18)
	w.Player.Heat = min(100, w.Player.Heat+3)
	w.Player.Respect++
	c := w.Conflict(f.ID, other.ID)
	w.Log("A word in the right ear", fmt.Sprintf("You leave %s believing %s moved against them. Their quarrel is now %s.", f.Name, other.Name, c.State), "politics")
	return nil
}

const (
	// FamilyRepair is how much condition a family buys back on one holding in
	// a day, and FamilyRepairCost is what each point of it costs them. Both
	// used to be free: the morning put eight points back on every damaged
	// property in the city and took nothing for it.
	FamilyRepair     = 8
	FamilyRepairCost = 9
	// FamilyWage is what one point of an organization's strength costs to keep
	// on the street for a day.
	FamilyWage = 8
	// MostWhoLeaveAtOnce caps how fast an unpaid organization comes apart, so
	// a long bad run is a decline and not a disappearance.
	MostWhoLeaveAtOnce = 6
)

// FamilyBill is what a day costs an organization. Power covers the anonymous
// mass of it; the named people are paid on top, and they are the ones whose
// pockets the rest of the game can see. Billing on strength alone meant a
// family that had collapsed to nothing owed nothing, stopped being short, and
// started paying its five remaining men again out of an empty safe.
//
// It is a function rather than a line inside the day because two things ask
// it: the morning that charges it, and everything that wants to say how a
// family is placed. Those must not be able to drift apart.
func (w *World) FamilyBill(f *Faction) int {
	if f == nil {
		return 0
	}
	return f.Power*FamilyWage + len(w.Members(f.ID))*SoldierWage
}

// FamilyIncome is what an organization's holdings bring in over a day.
func (w *World) FamilyIncome(f *Faction) int {
	if f == nil {
		return 0
	}
	income := 0
	for _, id := range w.FamilyHoldings(f.ID) {
		prop := w.Properties[id]
		income += prop.Income * prop.Condition / 100
	}
	return income * 24
}

// DaysOfCover is how long an organization could go on paying everybody if the
// money stopped coming in tomorrow. It is the honest measure of how a family is
// placed, and cash alone is not: a family holding five thousand dollars against
// a bill of eight hundred a day is in more trouble than one holding two
// thousand against a bill of ninety.
//
// A family whose income already covers its bill is not counting days at all.
func (w *World) DaysOfCover(f *Faction) int {
	if f == nil {
		return 0
	}
	shortfall := w.FamilyBill(f) - w.FamilyIncome(f)
	if shortfall <= 0 {
		return WellCovered
	}
	return min(WellCovered, f.Cash/shortfall)
}

// WellCovered is the point past which counting days stops meaning anything: an
// organization that is living within its income is not running out of money on
// any particular morning.
const WellCovered = 90

// HowTheyArePlaced describes an organization's finances in the terms anything
// reasoning about them should use. The words are the core's, not a screen's or
// a prompt's, because what it means to be struggling is a fact about the world.
func (w *World) HowTheyArePlaced(f *Faction) string {
	if f == nil {
		return "nobody"
	}
	switch {
	case f.Short > 0:
		return "cannot pay its people"
	case w.DaysOfCover(f) < 7:
		return "struggling"
	case w.DaysOfCover(f) < 30:
		return "getting by"
	}
	return "comfortable"
}
