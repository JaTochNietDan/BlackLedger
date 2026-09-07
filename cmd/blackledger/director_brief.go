package main

import (
	"blackledger/core"
	"fmt"
	"strings"
)

// This is a proposed job's cast and direction, not a historical fact or an
// off-screen transfer of money. The usual operation rules still own rewards.
type narrativeBrief struct {
	Location      string `json:"location"`
	Premise       string `json:"proposed_premise"`
	PlayerTask    string `json:"player_task"`
	SourceRole    string `json:"source_role,omitempty"`
	RecipientRole string `json:"recipient_role,omitempty"`
	Opening       string `json:"required_opening,omitempty"`
}

func jobBrief(w *core.World, operation string, connection *core.ArrangementMemory) narrativeBrief {
	place, _ := core.PlaceByID(w.Player.Location)
	if connection != nil {
		for _, l := range core.Locations {
			if strings.Contains(strings.ToLower(connection.Offer+" "+connection.Title), strings.ToLower(l.Name)) {
				place = l
				break
			}
		}
	}
	b := narrativeBrief{Location: place.Name}
	switch operation {
	case "collection":
		b.Premise = "A customer is ready to pay the establishment for its services. This is a new proposed task, not a consequence of an earlier job."
		b.SourceRole = "the customer who owes payment"
		b.RecipientRole = "the establishment's manager who is owed payment"
		b.PlayerTask = "Collect the customer's payment and deliver it to the manager."
	case "courier":
		b.Premise = "The establishment's manager needs a private message carried to the requesting contact."
		b.SourceRole = "the establishment's manager"
		b.RecipientRole = "the NPC requesting this job"
		b.PlayerTask = "Carry the manager's sealed message to the requesting contact without exposing its contents."
	case "mediation":
		b.Premise = "Two staff members disagree about how to share access to a work area. Neither is the player or a new named character."
		b.PlayerTask = "Hear both staff members and negotiate shared access without violence. Do not collect money or deliver a package."
	}
	if connection != nil {
		// A factual spoken acknowledgement that does not promote the old offer's
		// unverified claims into accomplished events.
		switch connection.Operation {
		case "courier":
			b.Opening = "You completed our last delivery."
		case "mediation":
			b.Opening = "You handled our last mediation without violence."
		case "collection":
			b.Opening = "You collected the last agreed payment."
		}
	}
	return b
}
func validateBriefOpening(body string, brief narrativeBrief) error {
	opening := strings.TrimSpace(brief.Opening)
	text := strings.TrimLeft(strings.TrimSpace(body), "\"“”")
	if opening != "" && !strings.HasPrefix(text, opening) {
		return fmt.Errorf("body must begin with this exact factual acknowledgement: %s Then describe the new task using the supplied source and recipient roles", opening)
	}
	return nil
}
