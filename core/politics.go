package core

import "fmt"

// BusinessPressure is an authored political decision grounded in current ownership.
// The schedule is private; the player sees only the delivered demand.
func (w *World) BusinessPressure() {
	w.NextPressure = w.Minute + 720
	var target *Place
	for i := range Locations {
		l := &Locations[i]
		if w.Own(l.ID) && w.Properties[l.ID].Income > 0 && (target == nil || w.Properties[l.ID].Income > w.Properties[target.ID].Income) {
			target = l
		}
	}
	if target == nil {
		w.NextPressure = 0
		return
	}
	actor := 0
	if target.District > 0 {
		actor = 1
	}
	f := &w.Factions[actor]
	if f.Goodwill >= 25 {
		w.Log("An understanding holds", f.Name+" leaves your businesses alone for now.", "politics")
		return
	}
	speaker := "vittorio"
	if actor == 1 {
		speaker = "elena"
	}
	rival := w.Factions[1-actor]
	w.Event = &Scene{ID: ID(), Kind: "business_pressure", Source: "authored", Minute: w.Minute, Speaker: speaker, Actor: f.ID, Target: target.ID, Title: "A claim on your earnings", Body: fmt.Sprintf("“%s is doing business under your name now. My people expect a share. Pay for an understanding, or explain why I should tolerate a competitor.”", target.Name), Choices: []Choice{
		{ID: "pay", Label: "Pay $60 for an understanding", Cost: 60, Detail: "Improves this family's standing by 12. Postpones the next demand; existing personal threats remain."},
		{ID: "resist", Label: "Refuse their claim", Detail: "Gain 2 respect; lose 20 standing. The family may retaliate against your business."},
		{ID: "ally", Label: "Seek backing from " + rival.Name, Cost: 35, Detail: "Pay $35 for an introduction. Gain 12 standing with their rival; lose 12 with this family. Business retaliation remains possible."},
	}}
	w.Log("A family wants a share", f.Name+" has delivered a demand concerning "+target.Name+".", "politics")
}
func (w *World) ResolvePressure(e *Scene, choice string) error {
	if !w.Own(e.Target) {
		return fmt.Errorf("the business is no longer yours")
	}
	actor := -1
	for i := range w.Factions {
		if w.Factions[i].ID == e.Actor {
			actor = i
		}
	}
	if actor < 0 {
		return fmt.Errorf("unknown family")
	}
	f := &w.Factions[actor]
	switch choice {
	case "pay":
		f.Cash += 60
		f.Goodwill = min(100, f.Goodwill+12)
		w.NextPressure = w.Minute + 1440
		w.Log("A purchased understanding", "You pay $60 to "+f.Name+". The next business demand is postponed for a day; personal threats are unchanged.", "politics")
	case "resist", "ally":
		if choice == "resist" {
			w.Player.Respect += 2
			f.Goodwill = max(-100, f.Goodwill-20)
		} else {
			f.Goodwill = max(-100, f.Goodwill-12)
			other := &w.Factions[1-actor]
			other.Goodwill = min(100, other.Goodwill+12)
			other.Cash += 35
		}
		w.Log("Taking a side", "You reject "+f.Name+"'s terms. The relationship has worsened.", "politics")
		// A hidden decision is committed now, not invented at the moment of playback.
		if w.Random() < .75 {
			w.Plots = append(w.Plots, Plot{ID: ID(), Kind: "sabotage", Life: w.Life, Due: w.Minute + 120, Actor: f.ID, Target: e.Target, Strength: 35})
		}
	default:
		return fmt.Errorf("unknown business decision")
	}
	return nil
}
func (w *World) ResolveSabotage(p Plot) {
	if !w.Own(p.Target) {
		return
	}
	l, ok := PlaceByID(p.Target)
	if !ok {
		return
	}
	damage := p.Strength
	// Available loyal crew can protect businesses; guards protect the residence only.
	if len(w.Player.Crew) > 0 && len(w.Tasks) == 0 && w.Player.Crew[0].Loyalty >= 50 {
		damage = max(10, damage-20)
	}
	prop := w.Properties[p.Target]
	damage = min(damage, prop.Condition)
	prop.Condition -= damage
	w.Log("Broken glass at "+l.Name, fmt.Sprintf("An attack damaged the business by %d condition. Income is reduced until repairs are made. The attackers left no proof of who sent them.", damage), "danger")
}
