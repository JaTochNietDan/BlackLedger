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
	w.contestAt(attacker, defender, weakest)
}

// contestAt is the same raid against a named holding. The city always goes for
// the weakest thing somebody holds; the player picks, which is the only
// difference between what they can do and what is done to them.
func (w *World) contestAt(attacker, defender *Faction, weakest string) {
	prop := w.Properties[weakest]
	if prop == nil || prop.Owner != defender.ID {
		return
	}
	// Nobody moves on somebody they stand with, in either direction.
	if w.Allied(attacker.ID) && defender.ID == w.PlayerOrganizationID() {
		return
	}
	if w.Allied(defender.ID) && attacker.ID == w.PlayerOrganizationID() {
		return
	}
	// And somebody who stands with the player may turn up at the door.
	if defender.ID == w.PlayerOrganizationID() {
		if ally, answered := w.AllyAnswers(attacker); answered {
			place, _ := PlaceByID(weakest)
			attacker.Power = max(10, attacker.Power-3)
			w.Antagonize(attacker.ID, ally.ID, 6)
			w.Log(ally.Name+" was already there", fmt.Sprintf("%s came for %s and found people who were not yours waiting in it. They left without it.", attacker.Name, place.Name), "politics")
			return
		}
	}
	place, ok := PlaceByID(weakest)
	if !ok {
		return
	}
	// Strength decides the odds; a weakened defender loses ground faster.
	// A man standing on the door is the difference between a place that is
	// walked into and a place that has to be taken.
	odds := .35 + float64(attacker.Power-defender.Power-w.PostingDefenceAt(weakest))/200
	if odds < .1 {
		odds = .1
	}
	if odds > .85 {
		odds = .85
	}
	if w.WorldRandom() >= odds {
		attacker.Power = max(10, attacker.Power-2)
		w.Log("A raid repelled at "+place.Name, fmt.Sprintf("%s moved against %s and %s driven off. %s %s the ground.",
			Leads(attacker.Name), place.Name, Agree(attacker.Name, "was", "were"),
			Leads(defender.Name), Agree(defender.Name, "holds", "hold")), "politics")
		w.Witness("gunfight", weakest, fmt.Sprintf("%s came for %s and %s driven off. %s still %s it.",
			Leads(attacker.Name), place.Name, Agree(attacker.Name, "was", "were"),
			Leads(defender.Name), Agree(defender.Name, "holds", "hold")), "")
		return
	}
	damage := min(prop.Condition, 20+int(w.WorldRandom()*25))
	prop.Condition -= damage
	defender.Power = max(10, defender.Power-max(2, damage/5))
	defender.Cash = max(0, defender.Cash-damage*15)
	// The player's organization is the player: money taken off it comes out of
	// their pocket, not out of a number that is rewritten every morning.
	if defender.ID == w.PlayerOrganizationID() {
		w.Player.Cash = max(0, w.Player.Cash-damage*15)
		w.ShiftCustom(weakest, "Somebody came through the front of it", -6)
	}
	// A raid reaches people, not only premises.
	if w.WorldRandom() < .18 {
		// Whoever is on the door is the one standing in it: the first thing a
		// raid reaches, which is the price of what they are worth.
		if victim := w.StoodInIt(weakest); victim != nil && victim.Faction == defender.ID {
			w.Properties[weakest].Posted = ""
			w.KillBy(victim.ID, nil, fmt.Sprintf("%s had come for %s and %s was on the door.", attacker.Name, place.Name, victim.Name))
		} else if victim := w.casualty(defender.ID); victim != nil {
			w.KillBy(victim.ID, nil, fmt.Sprintf("%s had come for %s.", attacker.Name, place.Name))
		} else if defender.ID == w.PlayerOrganizationID() && len(w.Player.Crew) > 0 {
			// Whoever stands with the player is who a raid reaches, because
			// they have no soldiers of their own to lose.
			member := w.Player.Crew[0]
			w.Player.Crew = w.Player.Crew[:0]
			if n := w.NPC(member.ID); n != nil {
				w.KillBy(n.ID, nil, fmt.Sprintf("%s had come for %s.", attacker.Name, place.Name))
			} else {
				w.Log(member.Name+" did not come out of it", fmt.Sprintf("%s came for %s and %s was standing in it.", attacker.Name, place.Name, member.Name), "danger")
			}
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
		w.Log("A holding changes hands", fmt.Sprintf("%s %s taken %s from %s. The city notices who could not hold it.",
			Leads(attacker.Name), Agree(attacker.Name, "has", "have"), place.Name, defender.Name), "politics")
		w.Report("seizure", strings.ToUpper(place.Name)+" CHANGES HANDS",
			fmt.Sprintf("%s now controls %s, previously held by %s. Neither organization would comment.",
				Leads(attacker.Name), place.Name, defender.Name))
		w.Witness("seizure", weakest, fmt.Sprintf("%s took %s off %s. Their people were in it by the evening.", attacker.Name, place.Name, defender.Name),
			strings.ToUpper(place.Name)+" CHANGES HANDS")
		return
	}
	w.Log("Trouble at "+place.Name, fmt.Sprintf("%s struck %s, a holding of %s. Its condition is now %d%%.", attacker.Name, place.Name, defender.Name, prop.Condition), "politics")
}

// classify names the state a quarrel is in. It is harder to enter a war than to
// leave one, so a conflict that has crossed the line does not flicker back the
// moment hostility dips.
func classify(c *Conflict) string {
	switch {
	case c.Hostility >= warAt, c.State == "war" && c.Hostility >= warEndsAt:
		return "war"
	case c.Hostility >= feudAt, c.State == "feud" && c.Hostility >= feudEndsAt:
		return "feud"
	}
	return "cold"
}

// Antagonize records that two organizations have been set against each other,
// whether by their own actions or by someone playing one against the other. The
// recorded state follows immediately, so hostility and state never disagree
// between one turn and the next.
func (w *World) Antagonize(a, b string, amount int) {
	c := w.Conflict(a, b)
	if c == nil {
		return
	}
	c.Hostility = min(100, max(0, c.Hostility+amount))
	if state := classify(c); state != c.State {
		c.State = state
		c.Since = w.Minute
	}
}

// considerSplinters gives weakened or embattled organizations a chance to lose
// people who would rather run their own operation. Checked before relations are
// advanced so a new organization takes part in the same turn it is created.
func (w *World) considerSplinters() {
	for i := 0; i < len(w.Factions); i++ {
		// A failing organization may lose its leader to the person below them
		// before it ever loses anyone to a rival.
		if w.movedThisDay(&w.Factions[i]) {
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
		// The player's organization ends when the player does, not when a bad
		// month leaves them holding nothing.
		if f.ID == w.PlayerOrganizationID() {
			kept = append(kept, f)
			continue
		}
		if len(w.FamilyHoldings(f.ID)) == 0 && f.Power <= 15 && len(w.Factions)-len(gone) > 2 {
			gone[f.ID] = true
			w.Log("An organization ends", fmt.Sprintf("%s no longer %s anything worth defending. What remains of it answers to someone else now.",
				Leads(f.Name), Agree(f.Name, "holds", "hold")), "politics")
			w.Report("collapse", "END OF "+strings.ToUpper(f.Name),
				fmt.Sprintf("%s %s ceased to operate as an organization. Its remaining interests have been absorbed by others.",
					Leads(f.Name), Agree(f.Name, "has", "have")))
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

// considerReestablish gives an organization that still has people but no ground
// a way back: it moves onto premises nobody is holding. Without this a city can
// settle permanently into one family owning everything, with the loser present
// in name and incapable of ever acting again.
func (w *World) considerReestablish() {
	for i := range w.Factions {
		f := &w.Factions[i]
		// People are what makes a comeback possible, not strength: an
		// organization reduced to nothing still has members who want somewhere
		// to work. Requiring strength first left beaten families permanently
		// inert, which locked the city into one family owning everything.
		if len(w.FamilyHoldings(f.ID)) > 0 || len(w.Members(f.ID)) == 0 {
			continue
		}
		if w.WorldRandom() >= .12 {
			continue
		}
		for _, l := range Locations {
			prop := w.Properties[l.ID]
			// Only premises that trade and that no living person holds: never
			// the player's, and never another organization's.
			if prop == nil || prop.Income <= 0 || w.Own(l.ID) || w.faction(prop.Owner) != nil {
				continue
			}
			if !strings.HasPrefix(prop.Owner, "former:") && prop.Owner != "independent" {
				continue
			}
			prop.Owner = f.ID
			prop.Condition = max(prop.Condition, 40)
			f.Power = min(peak(f), f.Power+8)
			w.Log("They are back on their feet", fmt.Sprintf("%s has taken over %s. An organization with nothing left has found somewhere to start again.", f.Name, l.Name), "politics")
			w.Report("recovery", upper(f.Name)+" MOVES INTO "+upper(l.Name),
				fmt.Sprintf("%s %s taken over the running of %s, which had been standing without an owner.",
					Leads(f.Name), Agree(f.Name, "has", "have"), l.Name))
			return
		}
	}
}

// FactionTurn advances relations between organizations once per day. Ambition,
// weakness and proximity move hostility; open war produces raids and seizures.
func (w *World) FactionTurn() {
	w.dissolve()
	w.considerReestablish()
	w.considerSplinters()
	w.ConsiderFactionContracts()
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

			previous := c.State
			c.State = classify(c)
			if c.State != previous {
				c.Since = w.Minute
				switch c.State {
				case "war":
					w.Log("Open war in the city", fmt.Sprintf("%s and %s are now at war. Their quarrel is not yours, but the city will feel it.", a.Name, b.Name), "politics")
					// Named. Every war in the city used to be reported under the
					// same headline, so a player who read the paper twice could
					// not tell that the second one was a different war.
					w.Report("war", upper(a.Name)+" AND "+upper(b.Name)+" AT WAR",
						fmt.Sprintf("Violence between %s and %s has escalated beyond the usual. Businesses in the affected districts are advised that the police cannot guarantee protection.", a.Name, b.Name))
				case "feud":
					// A quarrel reaching this state means two opposite things
					// depending on where it came from. Coming up from cold, two
					// families have stopped speaking. Coming down from war, the
					// shooting has stopped — and the first version of this
					// printed "bad blood between them" the day a war ended,
					// which is the wrong story told at the wrong moment.
					if previous == "war" {
						w.Log("A war burns out", fmt.Sprintf("%s and %s have stopped short of destroying each other.", a.Name, b.Name), "politics")
						w.Report("politics", "THE FIGHTING STOPS BETWEEN "+upper(a.Name)+" AND "+upper(b.Name),
							w.howItEnded(a, b))
						break
					}
					w.Log("A quarrel hardens", fmt.Sprintf("%s and %s are no longer on speaking terms.", a.Name, b.Name), "politics")
					// A quarrel hardening was written only into the player's own
					// record, so the city could be two moves from a war and the
					// paper had never mentioned it.
					w.Report("civic", "BAD BLOOD BETWEEN "+upper(a.Name)+" AND "+upper(b.Name),
						fmt.Sprintf("%s and %s are no longer on speaking terms, by the account of people who deal with both. Nothing has been said openly and nothing needs to be.", a.Name, b.Name))
				case "cold":
					if previous == "war" {
						w.Log("A war burns out", fmt.Sprintf("%s and %s have stopped short of destroying each other.", a.Name, b.Name), "politics")
						// And the end of a war was never reported at all: the
						// paper announced every war and never once said one was
						// over, so as far as a reader could tell they were all
						// still running.
						w.Report("politics", "THE FIGHTING STOPS BETWEEN "+upper(a.Name)+" AND "+upper(b.Name),
							w.howItEnded(a, b))
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

// howItEnded says why the shooting stopped, from what is left on the ground.
// A war that burned out and a war that finished somebody are the same
// transition in the model and are not the same story in the city.
func (w *World) howItEnded(a, b *Faction) string {
	beaten := ""
	switch {
	case len(w.FamilyHoldings(a.ID)) == 0:
		beaten = a.Name
	case len(w.FamilyHoldings(b.ID)) == 0:
		beaten = b.Name
	}
	if beaten != "" {
		return fmt.Sprintf("%s is not holding anything in the district any more. Whether that is the end of them is a question nobody is asking out loud.", beaten)
	}
	return fmt.Sprintf("%s and %s have stopped short of destroying each other. Both are smaller than they were, and both are still here.", a.Name, b.Name)
}
