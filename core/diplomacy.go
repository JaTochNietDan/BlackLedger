package core

import "fmt"

// AudienceReadiness explains why a family cannot see the player, or returns "".
// The audience is a seat across from whoever runs the family, so when nobody
// does there is no seat: the scene refuses to open, and the action has to say
// so rather than spending the player's evening on nothing.
func (w *World) AudienceReadiness(location string) string {
	actor := "bellandi"
	if location == "garage" {
		actor = "russo"
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

// OpenAudience is reached only after the location's action completes without interruption.
func (w *World) OpenAudience(location string) {
	// Who sits across the table is a job, not a person. Every other authored
	// scene already asks the city who holds it — the fixer, the detective —
	// and this one named Vittorio and Elena, the two people who happened to
	// lead the two families on the first morning. Families change hands: the
	// leader is shot, the strongest survivor takes over, and the chair here
	// went on holding a corpse.
	actor, surname := "bellandi", "Bellandi"
	if location == "garage" {
		actor, surname = "russo", "Russo"
	}
	leader := w.Leader(actor)
	if leader == nil {
		// Nobody is left to grant one. The view falls back to the first person
		// in the city when it cannot find a speaker, so an empty chair here
		// would put a stranger's face and voice on a family's terms.
		return
	}
	speaker := leader.ID
	body := "“People mistake an open door for an invitation. Tell me you understand the difference.”"
	if actor == "russo" {
		body = "“You are doing business near my people. We can settle our differences, or you can show me that an arrangement with you is worth something.”"
	}
	w.Event = &Scene{ID: ID(), Title: "A seat across from " + surname, Body: body, Speaker: speaker, Actor: actor, Target: location, Kind: "audience", Source: "authored", Minute: w.Minute, Choices: []Choice{
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
