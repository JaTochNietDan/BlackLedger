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
		income, condition, count := 0, 0, 0
		for _, id := range w.FamilyHoldings(f.ID) {
			prop := w.Properties[id]
			count++
			condition += prop.Condition
			income += prop.Income * prop.Condition / 100
			if prop.Condition < 100 {
				prop.Condition = min(100, prop.Condition+8)
			}
		}
		f.Cash = max(0, f.Cash+income*24)
		// Strength comes from holdings. An organization that holds nothing has
		// nothing to draw on and fades, rather than recovering to the strength
		// it had when it still owned half the waterfront.
		if count == 0 {
			f.Power = max(0, f.Power-4)
			continue
		}
		target := peak(f) * condition / (100 * count)
		switch {
		case f.Power < target:
			f.Power = min(target, f.Power+2)
		case f.Power > target:
			f.Power = max(target, f.Power-2)
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
func (w *World) sabotageChance(f *Faction) float64 {
	chance := .35 + float64(w.Player.Crew[0].Loyalty)/400 + float64(min(w.Player.Respect, 100))/500 - float64(f.Power)/300 + w.WeaponEdge()
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
func (w *World) Sabotage(id string) error {
	if reason := w.SabotageReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	f, _ := w.SabotageTarget(id)
	place, ok := PlaceByID(id)
	if !ok {
		return fmt.Errorf("unknown property")
	}
	prop := w.Properties[id]
	crew := w.Player.Crew[0].Name

	if w.Random() >= w.sabotageChance(f) {
		// Turned away. The family learns who came for them either way.
		injury := w.Absorb(12 + int(w.Random()*18))
		w.Ruin(30)
		w.Player.Health = max(0, w.Player.Health-injury)
		w.Player.Heat = min(100, w.Player.Heat+12)
		f.Goodwill = max(-100, f.Goodwill-20)
		w.Log("Turned away at "+place.Name, fmt.Sprintf("%s men were waiting. You and %s left without reaching anything, and you were hurt (-%d health). %s knows who came.", f.Name, crew, injury, f.Leader), "danger")
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
	f.Cash = max(0, f.Cash-damage*20)
	f.Goodwill = max(-100, f.Goodwill-30)
	w.Player.Respect += 4
	w.Player.Heat = min(100, w.Player.Heat+8)
	w.VisualCues = append(w.VisualCues, VisualCue{ID(), "attack", id, fmt.Sprintf("Your crew damaged %s. Condition is now %d%%.", place.Name, prop.Condition)})
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
	if w.Player.Contacts < 2 {
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
