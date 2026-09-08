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
	// Situation is the committed state of the city around this request: a war
	// the speaker's organization is fighting, or ground it recently lost. It is
	// context to write from, not an outcome to decide.
	Situation string `json:"current_situation,omitempty"`
	// Because is the committed event this work exists on account of. It is fact,
	// not invention: the speaker may refer to it.
	Because string `json:"exists_because,omitempty"`
	// Constraints belong in their own field because the model dramatizes the
	// premise and task. Observed on qwen3.5:35b-a3b: a leader recited the cast
	// restriction aloud as "Neither is you, but the tension threatens our
	// operations". These bound what may be written; they are never spoken.
	Constraints []string `json:"constraints_never_spoken,omitempty"`
}

func jobBrief(w *core.World, operation string, connection *core.ArrangementMemory) narrativeBrief {
	place, _ := core.PlaceByID(w.Player.Location)
	if connection != nil {
		for _, l := range core.Locations {
			if l.District <= w.District && strings.Contains(strings.ToLower(connection.Offer+" "+connection.Title), strings.ToLower(l.Name)) {
				place = l
				break
			}
		}
	}
	if connection != nil && connection.Location != "" {
		if saved, ok := core.PlaceByID(connection.Location); ok && saved.District <= w.District {
			place = saved
		}
	}
	b := narrativeBrief{Location: place.Name}
	switch operation {
	case "collection":
		b.Premise = "A customer is ready to pay the establishment for its services."
		b.SourceRole = "the customer who owes payment"
		b.RecipientRole = "the establishment's manager who is owed payment"
		b.PlayerTask = "Collect the customer's payment and deliver it to the manager."
		b.Constraints = []string{"This is a new proposed task, not a consequence of an earlier job."}
	case "courier":
		b.Premise = "The establishment's manager needs a private message carried to the requesting contact."
		b.SourceRole = "the establishment's manager"
		b.RecipientRole = "the contact requesting this job"
		b.PlayerTask = "Carry the manager's sealed message to the requesting contact without exposing its contents."
	case "mediation":
		b.Premise = "Two staff members disagree about how to share access to a work area."
		b.PlayerTask = "Hear both staff members and negotiate shared access without violence."
		b.Constraints = []string{
			"The two staff members are neither the listener nor any newly named character.",
			"This job collects no money and delivers no package.",
		}
	case "escort":
		b.Premise = "Something of value has to cross the city while it is dangerous to do so."
		b.SourceRole = "the person sending it"
		b.RecipientRole = "the person expecting it at the other end"
		b.PlayerTask = "Travel with it and see that it arrives."
		b.Constraints = []string{"Do not decide whether the journey is attacked; only propose the work."}
	case "warning":
		b.Premise = "A message has to reach the other side of a quarrel, said to their face."
		b.SourceRole = "the organization sending the message"
		b.RecipientRole = "somebody on the other side of the quarrel"
		b.PlayerTask = "Deliver the message in person and leave without starting anything."
		b.Constraints = []string{
			"Do not threaten a specific act of violence or name a consequence the rules have not committed.",
			"Do not decide how the other side answers.",
		}
	case "recovery":
		b.Premise = "Something was left behind on ground that changed hands, and its owner wants it back."
		b.SourceRole = "the person who lost it"
		b.RecipientRole = "whoever holds the place now"
		b.PlayerTask = "Get it out without a confrontation."
		b.Constraints = []string{"Do not transfer the property itself or change who holds it."}
	case "settlement":
		b.Premise = "An arrangement made with somebody who is gone has to be settled with whoever replaced them."
		b.SourceRole = "the party owed the arrangement"
		b.RecipientRole = "the person who now leads the other side"
		b.PlayerTask = "Put the matter in front of the new leadership and come away with it settled."
		b.Constraints = []string{
			"Do not invent the terms of the original arrangement beyond what the supplied history shows.",
			"Do not decide the outcome of the negotiation.",
		}
	}
	// Conflict-derived work carries the committed fact that justifies it, so the
	// speaker can refer to a war or a seizure that actually happened.
	for _, situational := range w.SituationalOperations() {
		if situational.ID == operation {
			b.Because = situational.Because
		}
	}
	b.Situation = briefSituation(w, connection)
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

// briefSituation describes the pressure the requesting side is actually under,
// drawn from committed conflict state. An organization fighting a war has fewer
// people to spare, which is a reason for work to exist at all.
func briefSituation(w *core.World, connection *core.ArrangementMemory) string {
	organization := ""
	if connection != nil {
		organization = connection.Beneficiary
	}
	parts := []string{}
	for _, c := range w.PublicConflicts() {
		if c.State != "war" && c.State != "feud" {
			continue
		}
		state := "is at war with"
		if c.State == "feud" {
			state = "is feuding with"
		}
		parts = append(parts, fmt.Sprintf("%s %s %s", c.Between[0], state, c.Between[1]))
	}
	if len(parts) == 0 {
		return ""
	}
	prefix := "The city's standing quarrels: "
	if organization != "" {
		if f := w.FactionByID(organization); f != nil {
			prefix = f.Name + " is asking while the city stands like this: "
		}
	}
	return prefix + strings.Join(parts, "; ") + ". Use this only as the situation the request happens inside. Do not narrate its outcome or invent a new one."
}
