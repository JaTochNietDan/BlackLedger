package core

import "fmt"

// Everything in a saved world that points at a name, checked against whether
// the name is still there.
//
// Two nights running, a thing survived the organization it belonged to. The
// player went on answering to a family that had ended, and the books listed an
// understanding with "An unidentified family" because the name resolved to
// nothing and the ledger printed it anyway. Neither showed up as a number going
// wrong; both showed up as a sentence, and only because somebody happened to
// walk that path.
//
// A city churns. Organizations fall and form, people die, businesses change
// hands, and every one of those leaves an id somewhere that used to mean
// something. This asks the question once, of everything, rather than one field
// at a time after somebody notices.
//
// It is deliberately a list of what is wrong rather than a bool: a caller that
// finds three wants to say which three.
func (w *World) Dangling() []string {
	family := map[string]bool{}
	for _, f := range w.Factions {
		family[f.ID] = true
	}
	person := map[string]bool{}
	for i := range w.NPCs {
		person[w.NPCs[i].ID] = true
	}
	place := map[string]bool{}
	for _, l := range Locations {
		place[l.ID] = true
	}

	out := []string{}
	note := func(what, id string) { out = append(out, fmt.Sprintf("%s points at %q", what, id)) }

	if w.Player.Serves != "" && !family[w.Player.Serves] {
		note("the player's employer", w.Player.Serves)
	}
	for _, p := range w.Pacts {
		if p.Life == w.Life && !family[p.With] {
			note("an understanding", p.With)
		}
	}
	for id := range w.BusinessTruces {
		if !family[id] {
			note("a ceasefire", id)
		}
	}
	// The player's own name is always good, incorporated or not: the city keeps
	// quarrels with it from the first time anybody takes offence, and the
	// organization is filed later.
	mine := w.PlayerOrganizationID()
	for _, c := range w.Conflicts {
		if c.A != mine && !family[c.A] {
			note("a quarrel", c.A)
		}
		if c.B != mine && !family[c.B] {
			note("a quarrel", c.B)
		}
	}
	for _, p := range w.Plots {
		if p.Life != w.Life {
			continue
		}
		if !family[p.Actor] {
			note("a plan", p.Actor)
		}
		if p.Target != "" && p.Kind == "sabotage" && !place[p.Target] {
			note("a plan's target", p.Target)
		}
	}
	for _, c := range w.Commissions {
		if c.Life != w.Life || c.Done || c.Failed {
			continue
		}
		if c.PatronID != "" && !family[c.PatronID] {
			note("a commission's patron", c.PatronID)
		}
	}
	for _, l := range w.Loans {
		if l.Life == w.Life && !person[l.Debtor] {
			note("a debt", l.Debtor)
		}
	}
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil {
			continue
		}
		// An owner is a family, the player, or nobody in particular.
		if prop.Owner != "" && prop.Owner != "independent" &&
			prop.Owner != mine && !family[prop.Owner] {
			note("the deeds to "+l.ID, prop.Owner)
		}
		if prop.Posted != "" && !person[prop.Posted] {
			note("somebody on the door at "+l.ID, prop.Posted)
		}
		for _, who := range prop.Hands {
			if !person[who] {
				note("somebody behind the counter at "+l.ID, who)
			}
		}
	}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Faction != "" && n.Faction != mine && !family[n.Faction] {
			note("somebody answering to", n.Faction)
		}
	}
	return out
}
