package core

import "fmt"

// OpenAudience is reached only after the location's action completes without interruption.
func (w *World) OpenAudience(location string) {
	actor, speaker, surname := "bellandi", "vittorio", "Bellandi"
	if location == "garage" {
		actor, speaker, surname = "russo", "elena", "Russo"
	}
	body := "“People mistake an open door for an invitation. Tell me you understand the difference.”"
	if actor == "russo" {
		body = "“You are doing business near my people. We can settle our differences, or you can show me that an arrangement with you is worth something.”"
	}
	w.Event = &Scene{ID: ID(), Title: "A seat across from " + surname, Body: body, Speaker: speaker, Actor: actor, Target: location, Kind: "audience", Source: "authored", Minute: w.Minute, Choices: []Choice{
		{ID: "tribute", Label: "Offer $150 in tribute", Cost: 150, Detail: "Gain 8 standing. Cancels this family's current operations against you; other families and future demands are unchanged."},
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
	case "tribute":
		faction.Goodwill = min(100, faction.Goodwill+8)
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
