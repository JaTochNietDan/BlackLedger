package core

import "fmt"

// BusinessPressure is an authored political decision grounded in current ownership.
// The schedule is private; the player sees only the delivered demand.
func (w *World) BusinessPressure() {
	// A business visibly skimming invites attention sooner; one run clean and
	// quiet buys time. Notice is a standing consequence of how the player runs
	// what they own, not a hidden roll.
	interval := 720 - 90*w.SkimNotice()
	w.NextPressure = w.Minute + max(240, min(1440, interval))
	var target *Place
	actor := 0
	ownsBusiness := false
	for i := range Locations {
		l := &Locations[i]
		if !w.Own(l.ID) || w.Properties[l.ID].Income <= 0 {
			continue
		}
		ownsBusiness = true
		family := 0
		if l.District > 0 {
			family = 1
		}
		// A local understanding does not suppress another family's territorial claim.
		if w.Factions[family].Goodwill >= 25 || w.BusinessTruces[w.Factions[family].ID] > w.Minute {
			continue
		}
		// The premises worth demanding a share of are the ones visibly earning.
		weight := int(float64(w.Properties[l.ID].Income) * operatingMode(w.Properties[l.ID].Mode).Take)
		best := 0
		if target != nil {
			best = int(float64(w.Properties[target.ID].Income) * operatingMode(w.Properties[target.ID].Mode).Take)
		}
		if target == nil || weight > best {
			target, actor = l, family
		}
	}
	if target == nil {
		if !ownsBusiness {
			w.NextPressure = 0
		} else {
			w.Log("An understanding holds", "Your relationships with the families in your business districts keep their demands at bay for now.", "politics")
		}
		return
	}
	f := &w.Factions[actor]
	speaker := "vittorio"
	if actor == 1 {
		speaker = "elena"
	}
	rival := w.Rival(f.ID)
	if rival == nil {
		rival = f
	}
	w.Event = &Scene{ID: ID(), Kind: "business_pressure", Source: "authored", Minute: w.Minute, Speaker: speaker, Actor: f.ID, Target: target.ID, Title: "A claim on your earnings", Body: fmt.Sprintf("“%s is doing business under your name now. My people expect a share. Pay for an understanding, or explain why I should tolerate a competitor.”", target.Name), Choices: []Choice{
		{ID: "pay", Label: "Pay $60 for an understanding", Cost: 60, Detail: "Improves this family's standing by 12. Postpones the next demand; existing personal threats remain."},
		{ID: "resist", Label: "Refuse their claim", Detail: "Gain 2 respect; lose 20 standing. The family may retaliate against your business. Repeated defiance can put your life at risk."},
		{ID: "ally", Label: "Seek backing from " + rival.Name, Cost: 35, Detail: "Pay $35 for an introduction. Gain 12 standing with their rival; lose 12 with this family. Business retaliation remains possible; a deepening feud can also put your life at risk."},
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
			other := w.Rival(f.ID)
			if other == nil {
				other = f
			}
			other.Goodwill = min(100, other.Goodwill+12)
			other.Cash += 35
			// Being seen to take one family's side against another is exactly
			// the kind of thing that turns a standing quarrel into a war.
			w.Antagonize(f.ID, other.ID, 12)
			w.Log("A rival introduction", fmt.Sprintf("You pay $35 for an introduction to %s. Their standing is now %+d; %s standing is now %+d. This does not guarantee protection, and %s will hear who you went to.", other.Name, other.Goodwill, f.Name, f.Goodwill, f.Name), "politics")
		}
		if choice == "resist" {
			w.Log("Taking a side", "You reject "+f.Name+"'s terms. The relationship has worsened.", "politics")
		}
		if f.Goodwill <= -40 {
			w.RetaliationFrom(f.ID)
		}
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
	if w.BusinessTruces[p.Actor] > w.Minute {
		return
	}
	if !w.Own(p.Target) {
		return
	}
	l, ok := PlaceByID(p.Target)
	if !ok {
		return
	}
	damage := max(0, p.Strength)
	unguarded := min(damage, w.Properties[p.Target].Condition)
	// Available loyal crew can protect businesses; guards protect the residence only.
	if len(w.Player.Crew) > 0 && len(w.Tasks) == 0 && w.Player.Crew[0].Loyalty >= 50 {
		damage = min(damage, max(10, damage-20))
	}
	prop := w.Properties[p.Target]
	damage = min(damage, prop.Condition)
	prop.Condition -= damage
	defense := ""
	if prevented := unguarded - damage; prevented > 0 {
		defense = fmt.Sprintf(" %s helped defend the business, preventing %d additional condition loss.", w.Player.Crew[0].Name, prevented)
	}
	evidence := "The attackers left no proof of who sent them."
	if p.Known {
		evidence = "This matches the operation your sources uncovered."
	}
	w.VisualCues = append(w.VisualCues, VisualCue{ID(), "attack", p.Target, fmt.Sprintf("The attack damaged %s. Condition is now %d%%.", l.Name, prop.Condition)})
	w.Log("Broken glass at "+l.Name, fmt.Sprintf("An attack damaged the business by %d condition. Income is reduced until repairs are made. %s", damage, evidence)+defense, "danger")
}

// Investigation reveals only actual, still-active plans and identifies their real target.
func (w *World) Investigate() {
	found := false
	for i := range w.Plots {
		p := &w.Plots[i]
		if p.Life != w.Life {
			continue
		}
		found = true
		p.Known = true
		actor := "An unknown organization"
		for _, f := range w.Factions {
			if f.ID == p.Actor {
				actor = f.Name
			}
		}
		detail := actor + " has commissioned an attack against you. Avoid home, negotiate, or arrange residential protection."
		if p.Kind == "sabotage" {
			l, ok := PlaceByID(p.Target)
			if !ok {
				continue
			}
			detail = actor + " has commissioned an attack on " + l.Name + ". Available loyal crew can limit business damage; home security cannot."
		}
		w.Log("Word on the street", detail, "intel")
	}
	if !found {
		w.Log("Word on the street", "Your sources have no evidence of an active operation against you. This is not a guarantee of safety.", "intel")
	}
}

// Sensitive arrangements expose the player to a bounded police decision once heat is high.
func (w *World) PoliceStop(job *Scene) {
	// Whoever has the caseload this week. FillRoles guarantees somebody does.
	w.FillRoles()
	w.RememberArrangement(job, "awaiting_police")
	w.Event = &Scene{JobID: job.ID, Operation: job.Operation, ID: ID(), Kind: "police_stop", Source: "authored", Minute: w.Minute, Speaker: w.HolderID("detective"), Actor: job.Speaker, Beneficiary: job.Beneficiary, Title: "Too familiar a face", Body: "“Your name keeps coming up in the same places. Before you finish this arrangement, we are going to have a conversation. You can settle this inconvenience, or leave the business unfinished.”", Effect: job.Effect, Outcome: job.Outcome, Choices: []Choice{
		{ID: "pay", Label: "Pay $40 and finish the job", Cost: 40, Detail: fmt.Sprintf("Receive the agreed $%d reward and reputation afterward. Police attention decreases before this job's heat is applied.", job.Effect.Reward)},
		{ID: "abandon", Label: "Abandon the arrangement", Detail: "No payment or reward. Lose 6 heat; the work and time already spent are lost."},
	}}
	w.Log("Police attention catches up", "A detective interrupts the arrangement before its reward is paid.", "danger")
}
func (w *World) CompleteArrangement(job *Scene) {
	w.RememberArrangement(job, "completed")
	w.Earn(job.Effect.Reward)
	w.Player.Respect += job.Effect.Respect
	w.Player.Heat = min(100, w.Player.Heat+job.Effect.Heat)
	speaker := job.Speaker
	if job.Kind == "police_stop" {
		speaker = job.Actor
	}
	if npc := w.NPC(speaker); npc != nil {
		npc.Trust += 3
	}
	w.ResolveBeneficiary(job.Beneficiary)
	// Work done for whoever the player answers to is the only way anybody
	// comes up inside an organization.
	w.ServeWork(job.Beneficiary)
	title := job.Title
	if job.Kind == "police_stop" {
		// The detective's scene is an interruption, not the arrangement's name.
		// Older saves may lack a linked memory; avoid inventing a title then.
		title = "An arrangement completed"
		for _, memory := range w.Arrangements {
			if memory.ID == job.JobID && memory.Life == w.Life && memory.Title != "" {
				title = memory.Title
				break
			}
		}
	}
	w.Log(title, job.Outcome+fmt.Sprintf(" ($%d, respect +%d)", job.Effect.Reward, job.Effect.Respect), "story")
}

func (w *World) ResolveBeneficiary(id string) {
	if id == "" {
		return
	}
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.ID == id {
			f.Goodwill = min(100, f.Goodwill+6)
			w.Log("An alliance takes shape", f.Name+" appreciates the completed arrangement. Standing improves by 6.", "politics")
		} else {
			f.Goodwill = max(-100, f.Goodwill-3)
		}
	}
	// Supporting rivals repeatedly can create a genuine feud, even without direct provocation.
	for _, f := range w.Factions {
		if f.Goodwill <= -30 {
			w.RetaliationFrom(f.ID)
		}
	}
}
