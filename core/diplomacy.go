package core

import "fmt"

// AudienceReadiness explains why a family cannot see the player, or returns "".
// The audience is a seat across from whoever runs the family, so when nobody
// does there is no seat: the scene refuses to open, and the action has to say
// so rather than spending the player's evening on nothing.
// AudienceActor is the family you can sit down with in this room: whoever of
// them is standing in it with the authority to agree to anything.
//
// This used to be the address. The club meant Bellandi and the garage meant
// Russo, whatever either family was actually doing, so a sit-down at The
// Monarch seated the player across from a man who was two miles away at a
// filling station — and asking for one anywhere else was not possible even with
// their lead in the room.
//
// Preference goes to whoever holds the premises, because a family's own ground
// is where they receive people; failing that, anybody of theirs who is here.
func (w *World) AudienceActor(location string) string {
	if prop := w.Properties[location]; prop != nil && prop.Owner != "" {
		if f := w.faction(prop.Owner); f != nil && f.ID != w.PlayerOrganizationID() {
			if w.Speaker(f.ID, location) != nil {
				return f.ID
			}
		}
	}
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.ID == w.PlayerOrganizationID() {
			continue
		}
		if w.Speaker(f.ID, location) != nil {
			return f.ID
		}
	}
	return ""
}

func (w *World) AudienceReadiness(location string) string {
	actor := w.AudienceActor(location)
	if actor == "" {
		return w.nobodyToSitWith(location)
	}
	f := w.faction(actor)
	if f == nil {
		return "There is no such family in this city any more"
	}
	if w.Leader(actor) == nil {
		return "Nobody is left to speak for " + f.Name
	}
	return ""
}

// nobodyToSitWith says who you would have to find, and where they are, when
// there is nobody in this room to sit down with.
func (w *World) nobodyToSitWith(location string) string {
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.ID == w.PlayerOrganizationID() {
			continue
		}
		if lead := w.Leader(f.ID); lead != nil {
			return "Nobody here can settle anything for a family. " + w.whereabouts(lead)
		}
	}
	return "There is nobody here who speaks for anybody"
}

// Speaker is who can settle terms for a family in this room: the lead if they
// are standing here, otherwise the most senior of theirs who is. Nobody below a
// lieutenant is authorised to agree to anything.
func (w *World) Speaker(actor, location string) *NPC {
	if lead := w.Leader(actor); lead != nil && lead.Location == location && !w.Travelling(lead) {
		return lead
	}
	var best *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Faction != actor || n.Location != location || w.Travelling(n) {
			continue
		}
		if n.Rank < RankLieutenant {
			continue
		}
		if best == nil || n.Rank > best.Rank {
			best = n
		}
	}
	return best
}

// whereabouts says where somebody actually is, in the city's own words, so a
// refusal tells the player where to walk rather than only that they cannot.
func (w *World) whereabouts(n *NPC) string {
	if n == nil {
		return ""
	}
	if w.Travelling(n) {
		if to, ok := PlaceByID(n.Heading); ok {
			return n.Name + " is on the street, walking to " + to.Name + "."
		}
		return n.Name + " is somewhere on the street."
	}
	if place, ok := PlaceByID(n.Location); ok {
		return n.Name + " is at " + place.Name + "."
	}
	return "Nobody knows where " + n.Name + " is."
}

// OpenAudience is reached only after the location's action completes without interruption.
func (w *World) OpenAudience(location string) {
	// Who sits across the table is a job, not a person. Every other authored
	// scene already asks the city who holds it — the fixer, the detective —
	// and this one named Vittorio and Elena, the two people who happened to
	// lead the two families on the first morning. Families change hands: the
	// leader is shot, the strongest survivor takes over, and the chair here
	// went on holding a corpse.
	actor := w.AudienceActor(location)
	if actor == "" {
		return
	}
	f := w.faction(actor)
	if f == nil {
		return
	}
	surname := f.Name
	// Whoever of theirs is actually in this room, which is the lead when the
	// lead is here. The view falls back to the first person in the city when it
	// cannot find a speaker, so an empty chair here would put a stranger's face
	// and voice on a family's terms.
	across := w.Speaker(actor, location)
	if across == nil {
		return
	}
	speaker := across.ID
	body := "“People mistake an open door for an invitation. Tell me you understand the difference.”"
	if actor == "russo" {
		body = "“You are doing business near my people. We can settle our differences, or you can show me that an arrangement with you is worth something.”"
	}
	// A lieutenant settling terms in his lead's absence says so. He can agree to
	// the same things — the money is the family's either way — but the player
	// should know whose word they are taking.
	title := "A seat across from " + surname
	if lead := w.Leader(actor); lead == nil || across.ID != lead.ID {
		title = "A seat across from " + across.Name
		body = "“" + surname + " is not here. I can speak for the family, and what I agree to holds.”"
	}
	w.Event = &Scene{ID: ID(), Title: title, Body: body, Speaker: speaker, Actor: actor, Target: location, Kind: "audience", Source: "authored", Minute: w.Minute, Choices: []Choice{
		{ID: "tribute", Label: "Offer $150 in tribute", Cost: 150, Detail: fmt.Sprintf("Gain %d standing; how you are dressed is part of what they are willing to grant. Cancels this family's current operations against you; other families and future demands are unchanged.", 8+w.Standing()/4)},
		{ID: "business_truce", Label: "Pay $100 for a business ceasefire", Cost: 100, Detail: "For 24 game hours, this family suspends business demands and sabotage. Cancels its pending sabotage; personal threats, other families and ownership remain unchanged. Renews from now, not cumulatively."},
		{ID: "work", Label: "Offer to do a favor", Detail: "Hear a paid courier offer for this family. Completion improves their standing and worsens their rival's. Existing threats remain until resolved separately."},
		{ID: "leave", Label: "Leave without an agreement", Detail: "No payment. Existing threats remain."},
	}}
}
func (w *World) ResolveAudience(e *Scene, choice string) error {
	actor := e.Actor
	// Saves from before actor-specific audiences contained only Bellandi audiences.
	if actor == "" {
		actor = "bellandi"
	}
	var faction *Faction
	for i := range w.Factions {
		if w.Factions[i].ID == actor {
			faction = &w.Factions[i]
			break
		}
	}
	if faction == nil {
		return fmt.Errorf("unknown audience family")
	}
	switch choice {
	case "business_truce":
		if w.BusinessTruces == nil {
			w.BusinessTruces = map[string]int{}
		}
		w.BusinessTruces[actor] = w.Minute + 1440
		faction.Cash += 100
		remaining := []Plot{}
		for _, plot := range w.Plots {
			if plot.Actor != actor || plot.Kind != "sabotage" {
				remaining = append(remaining, plot)
			}
		}
		w.Plots = remaining
		w.Log("A business ceasefire", fmt.Sprintf("%s accepts $100. Business demands and sabotage pause until Day %d %02d:%02d. Personal threats and other families remain unaffected.", faction.Name, (w.Minute+1440)/1440+1, w.Minute%1440/60, w.Minute%60), "politics")
	case "tribute":
		// The same money, taken more seriously from somebody who looks like
		// they have more of it.
		gained := 8 + w.Standing()/4
		faction.Goodwill = min(100, faction.Goodwill+gained)
		faction.Cash += 150
		remaining := []Plot{}
		for _, p := range w.Plots {
			if p.Actor != actor {
				remaining = append(remaining, p)
			}
		}
		w.Plots = remaining
		w.Log("A temporary understanding", fmt.Sprintf("%s accepts $150 and calls off its current operations against you. Standing is now %+d. Other families and future demands are unchanged.", faction.Name, faction.Goodwill), "politics")
	case "work":
		speaker := e.Speaker
		if speaker == "" {
			speaker = "vittorio"
		}
		scene, err := w.ValidateProposal(Proposal{Title: "A message for " + faction.Name, Body: fmt.Sprintf("“Carry a sealed message for %s. Keep it private and report back when it is delivered. Finish this properly and we will have a reason to speak again.”", faction.Name), Speaker: speaker, Operation: "courier", Outcome: "The sealed message was delivered.", Beneficiary: actor, Approaches: []Approach{{Method: "careful", Label: "Check the route first"}, {Method: "press", Label: "Take the direct route"}}})
		if err != nil {
			return err
		}
		scene.Source = "authored"
		w.Event = scene
		w.RememberArrangement(scene, "offered")
		w.Log(scene.Title, scene.Body, "story")
	case "leave":
		w.Log("No agreement reached", "You leave the meeting with "+faction.Name+". Existing threats remain.", "personal")
	default:
		return fmt.Errorf("unknown audience decision")
	}
	return nil
}

// Only current, publicly agreed terms are exposed; private plans stay private.
func (w *World) ActiveBusinessTruces() map[string]int {
	active := map[string]int{}
	if !w.Player.Alive {
		return active
	}
	for actor, until := range w.BusinessTruces {
		if until > w.Minute {
			active[actor] = until
		}
	}
	return active
}
