package core

import (
	"fmt"
	"strings"
)

// Organizations have opinions about each other, not only about the player. A
// conflict is the recorded state between two of them: it builds, it breaks into
// open war, holdings change hands, and it settles or ends in a conquest. None of
// it needs the player's participation, and all of it is committed by the same
// rules the player is judged by.

// Conflict is stored with its participants in sorted order so a pair has exactly
// one record however it is looked up.
type Conflict struct {
	A         string `json:"a"`
	B         string `json:"b"`
	Hostility int    `json:"hostility"`
	State     string `json:"state"` // "cold", "feud" or "war"
	Since     int    `json:"since"`
}

const (
	// The level an undisturbed quarrel between established families settles at.
	settledHostility = 52
	feudAt           = 40
	feudEndsAt       = 30
	warAt            = 70
	warEndsAt        = 45
)

func pairKey(a, b string) (string, string) {
	if a > b {
		return b, a
	}
	return a, b
}

// Conflict returns the recorded state between two organizations, creating a cold
// record the first time a pair is considered.
func (w *World) Conflict(a, b string) *Conflict {
	first, second := pairKey(a, b)
	if first == second {
		return nil
	}
	for i := range w.Conflicts {
		if w.Conflicts[i].A == first && w.Conflicts[i].B == second {
			return &w.Conflicts[i]
		}
	}
	w.Conflicts = append(w.Conflicts, Conflict{A: first, B: second, State: "cold"})
	return &w.Conflicts[len(w.Conflicts)-1]
}

// Rival is the strongest organization other than the given one. It replaces the
// former assumption that the city holds exactly two families.
func (w *World) Rival(id string) *Faction {
	var best *Faction
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.ID == id {
			continue
		}
		if best == nil || f.Power > best.Power {
			best = f
		}
	}
	return best
}

// FactionByID finds an organization by id, or nil when the holder is the player,
// an independent owner or a dead protagonist's estate.
func (w *World) FactionByID(id string) *Faction {
	if id == "" {
		return nil
	}
	return w.faction(id)
}

// PropertyHolder names whoever holds a property, for reports and narration.
func (w *World) PropertyHolder(id string) string {
	prop := w.Properties[id]
	if prop == nil {
		return ""
	}
	return prop.Owner
}

// contest resolves one attempt by an attacking organization against a defending
// one. The attacker goes after the defender's weakest holding, which is how a
// war eventually strips an organization rather than grinding forever.
func (w *World) contest(attacker, defender *Faction) {
	holdings := w.FamilyHoldings(defender.ID)
	if len(holdings) == 0 {
		return
	}
	weakest := holdings[0]
	for _, id := range holdings {
		if w.Properties[id].Condition < w.Properties[weakest].Condition {
			weakest = id
		}
	}
	prop := w.Properties[weakest]
	place, ok := PlaceByID(weakest)
	if !ok {
		return
	}
	// Strength decides the odds; a weakened defender loses ground faster.
	odds := .35 + float64(attacker.Power-defender.Power)/200
	if odds < .1 {
		odds = .1
	}
	if odds > .85 {
		odds = .85
	}
	if w.WorldRandom() >= odds {
		attacker.Power = max(10, attacker.Power-2)
		w.Log("A raid repelled at "+place.Name, fmt.Sprintf("%s moved against %s and was driven off. %s holds the ground.", attacker.Name, place.Name, defender.Name), "politics")
		return
	}
	damage := min(prop.Condition, 20+int(w.WorldRandom()*25))
	prop.Condition -= damage
	defender.Power = max(10, defender.Power-max(2, damage/5))
	defender.Cash = max(0, defender.Cash-damage*15)
	// A raid reaches people, not only premises.
	if w.WorldRandom() < .18 {
		if victim := w.casualty(defender.ID); victim != nil {
			w.Kill(victim.ID, fmt.Sprintf("Killed at %s when %s came for it.", place.Name, attacker.Name))
		}
	}

	// A holding that can no longer be defended changes hands.
	if prop.Condition <= 15 && attacker.Power > defender.Power {
		prop.Owner = attacker.ID
		prop.Condition = max(prop.Condition, 25)
		attacker.Power = min(100, attacker.Power+6)
		attacker.Peak = max(peak(attacker), attacker.Power)
		if c := w.Conflict(attacker.ID, defender.ID); c != nil {
			c.Hostility = min(100, c.Hostility+10)
		}
		w.Log("A holding changes hands", fmt.Sprintf("%s has taken %s from %s. The city notices who could not hold it.", attacker.Name, place.Name, defender.Name), "politics")
		return
	}
	w.Log("Trouble at "+place.Name, fmt.Sprintf("%s struck %s, a holding of %s. Its condition is now %d%%.", attacker.Name, place.Name, defender.Name, prop.Condition), "politics")
}

// Antagonize records that two organizations have been set against each other,
// whether by their own actions or by someone playing one against the other.
func (w *World) Antagonize(a, b string, amount int) {
	if c := w.Conflict(a, b); c != nil {
		c.Hostility = min(100, max(0, c.Hostility+amount))
	}
}

// considerSplinters gives weakened or embattled organizations a chance to lose
// people who would rather run their own operation. Checked before relations are
// advanced so a new organization takes part in the same turn it is created.
func (w *World) considerSplinters() {
	for i := 0; i < len(w.Factions); i++ {
		// A failing organization may lose its leader to the person below them
		// before it ever loses anyone to a rival.
		if w.WorldRandom() < 0.03 && w.ConsiderInternalMove(&w.Factions[i]) {
			return
		}
		if w.WorldRandom() < 0.05 && w.Splinter(&w.Factions[i]) {
			return // one upheaval at a time
		}
	}
}

// dissolve removes organizations that hold nothing and have no strength left to
// take anything back. Their quarrels go with them; the people who led them
// remain in the city as ordinary names.
func (w *World) dissolve() {
	kept := w.Factions[:0]
	gone := map[string]bool{}
	for _, f := range w.Factions {
		if len(w.FamilyHoldings(f.ID)) == 0 && f.Power <= 15 && len(w.Factions)-len(gone) > 2 {
			gone[f.ID] = true
			w.Log("An organization ends", f.Name+" no longer holds anything worth defending. What remains of it answers to someone else now.", "politics")
			continue
		}
		kept = append(kept, f)
	}
	if len(gone) == 0 {
		return
	}
	w.Factions = kept
	conflicts := w.Conflicts[:0]
	for _, c := range w.Conflicts {
		if !gone[c.A] && !gone[c.B] {
			conflicts = append(conflicts, c)
		}
	}
	w.Conflicts = conflicts
	plots := w.Plots[:0]
	for _, p := range w.Plots {
		if !gone[p.Actor] {
			plots = append(plots, p)
		}
	}
	w.Plots = plots
}

// FactionTurn advances relations between organizations once per day. Ambition,
// weakness and proximity move hostility; open war produces raids and seizures.
func (w *World) FactionTurn() {
	w.dissolve()
	w.considerSplinters()
	for i := range w.Factions {
		for j := i + 1; j < len(w.Factions); j++ {
			a, b := &w.Factions[i], &w.Factions[j]
			c := w.Conflict(a.ID, b.ID)
			if c == nil {
				continue
			}
			// Relations drift without direction. What actually starts wars is
			// events: a visible weakness, a seizure, or someone playing one
			// organization against another. A city at peace can stay at peace.
			drift := int(w.WorldRandom()*7) - 3
			switch {
			case len(w.FamilyHoldings(a.ID)) == 0 || len(w.FamilyHoldings(b.ID)) == 0:
				drift = -3 // nothing left to fight over
			case c.State == "war":
				drift = 0 // exhaustion is applied after the raids below
			default:
				// Relations pull back toward an uneasy normal, so quarrels cool
				// as readily as they build. What breaks that equilibrium is
				// events: a seizure, a betrayal, or someone with an interest in
				// seeing these two fight. A badly weakened neighbour is also a
				// standing temptation.
				drift += (settledHostility - c.Hostility) / 35
				if gap := a.Power - b.Power; gap > 40 || gap < -40 {
					drift++
				}
			}
			c.Hostility = min(100, max(0, c.Hostility+drift))

			// Hysteresis: a war is easier to stay in than to enter. Without it a
			// quarrel crossing the threshold flickers back the same day and
			// never costs anyone ground.
			previous := c.State
			switch {
			case c.Hostility >= warAt, previous == "war" && c.Hostility >= warEndsAt:
				c.State = "war"
			case c.Hostility >= feudAt, previous == "feud" && c.Hostility >= feudEndsAt:
				c.State = "feud"
			default:
				c.State = "cold"
			}
			if c.State != previous {
				c.Since = w.Minute
				switch c.State {
				case "war":
					w.Log("Open war in the city", fmt.Sprintf("%s and %s are now at war. Their quarrel is not yours, but the city will feel it.", a.Name, b.Name), "politics")
				case "feud":
					w.Log("A quarrel hardens", fmt.Sprintf("%s and %s are no longer on speaking terms.", a.Name, b.Name), "politics")
				case "cold":
					if previous == "war" {
						w.Log("A war burns out", fmt.Sprintf("%s and %s have stopped short of destroying each other.", a.Name, b.Name), "politics")
					}
				}
			}
			if c.State != "war" {
				continue
			}
			// Each side gets one attempt per day, strongest moving first.
			attacker, defender := a, b
			if b.Power > a.Power {
				attacker, defender = b, a
			}
			// Roughly one raid a day between them. More than that strips an
			// organization before anyone can react to the war starting.
			if w.WorldRandom() < .35 {
				w.contest(attacker, defender)
			}
			if w.WorldRandom() < .2 {
				w.contest(defender, attacker)
			}
			// Losing everything ends the quarrel; so does mutual exhaustion.
			if len(w.FamilyHoldings(defender.ID)) == 0 || len(w.FamilyHoldings(attacker.ID)) == 0 {
				c.Hostility = 0
				c.State = "cold"
				c.Since = w.Minute
			} else {
				// Open war is expensive for both sides and burns itself out.
				c.Hostility = max(0, c.Hostility-3)
			}
		}
	}
}

// PublicConflicts reports the quarrels the city can see. Hostility is a private
// number; what the street knows is who is feuding and who is at war, and since
// when. Cold relations are not news and are omitted.
type PublicConflict struct {
	Between []string `json:"between"`
	State   string   `json:"state"`
	Since   int      `json:"since"`
}

func (w *World) PublicConflicts() []PublicConflict {
	out := []PublicConflict{}
	for _, c := range w.Conflicts {
		if c.State == "cold" {
			continue
		}
		a, b := w.factionName(c.A), w.factionName(c.B)
		out = append(out, PublicConflict{Between: []string{a, b}, State: c.State, Since: c.Since})
	}
	return out
}

// HolderName describes who holds a property in words the interface can show.
func (w *World) HolderName(id string) string {
	prop := w.Properties[id]
	if prop == nil {
		return ""
	}
	switch {
	case w.Own(id):
		return "Your organization"
	case prop.Owner == "independent":
		return "Independent"
	case strings.HasPrefix(prop.Owner, "former:"):
		return strings.TrimPrefix(prop.Owner, "former:") + " (deceased)"
	case strings.HasPrefix(prop.Owner, "player:"):
		return "A previous life"
	}
	if f := w.faction(prop.Owner); f != nil {
		return f.Name
	}
	return prop.Owner
}
