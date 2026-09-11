package core

import (
	"fmt"
	"sort"
)

const (
	// PressureFloor is the least a family will take for an understanding about
	// a business, and what the demand used to be for every business in the
	// city whatever it earned.
	PressureFloor = 60
	// PressureCeiling is the most. A claim on one shop is a claim on one shop.
	PressureCeiling = 220
)

// BusinessPressure is an authored political decision grounded in current ownership.
// The schedule is private; the player sees only the delivered demand.
func (w *World) BusinessPressure() {
	// A business visibly skimming invites attention sooner; one run clean and
	// quiet buys time. Notice is a standing consequence of how the player runs
	// what they own, not a hidden roll.
	interval := 720 - 90*w.SkimNotice()
	w.NextPressure = w.Minute + max(240, min(1440, interval))
	var target *Place
	var actor *Faction
	ownsBusiness := false
	claims := w.Claimants()
	for i := range Locations {
		l := &Locations[i]
		if !w.Own(l.ID) || w.Properties[l.ID].Income <= 0 {
			continue
		}
		ownsBusiness = true
		family := claims[l.District]
		if family == nil {
			// Nobody is left with a claim on this street.
			continue
		}
		// A local understanding does not suppress another family's territorial claim.
		if family.Goodwill >= 25 || w.BusinessTruces[family.ID] > w.Minute {
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
	f := actor
	// The demand comes from whoever runs the family, for the same reason the
	// audience does. This named Vittorio and Elena by position in the faction
	// list, so a successor's demand arrived in a dead predecessor's voice.
	leader := w.Leader(f.ID)
	if leader == nil {
		// A family with nobody left to lead it is not collecting from anybody.
		return
	}
	speaker := leader.ID
	rival := w.Rival(f.ID)
	if rival == nil {
		rival = f
	}
	// Enough of the district in one pair of hands, and the demand stops being
	// a share and becomes the place.
	if buyout, wanted := w.TheirClaim(f.ID, target.ID); wanted {
		w.Event = &Scene{ID: ID(), Kind: "business_pressure", Source: "authored", Minute: w.Minute,
			Speaker: speaker, Actor: f.ID, Target: target.ID, Title: "They want the place",
			Body: fmt.Sprintf("“You have taken half this district under your own name and gone on sending me an envelope for it. I am not asking for a share of %s any more. I am asking for %s.”", target.Name, target.Name),
			Choices: []Choice{
				{ID: "hand", Label: "Hand them " + target.Name,
					Detail: fmt.Sprintf("The deed goes to %s and the people behind the counter stay where they are. Their standing with you improves by 30, and they leave your other places alone for five days.", f.Name)},
				{ID: "pay", Label: fmt.Sprintf("Buy them off for $%d", buyout), Cost: buyout,
					Detail: fmt.Sprintf("Six times an ordinary understanding, because this is the price of not deciding. Improves their standing by 12 and postpones the next demand by a day.")},
				{ID: "resist", Label: "Tell them no",
					Detail: "Gain 2 respect; lose 20 standing. A family that came for the place and was refused does not send another envelope. They may come for the premises, and repeated defiance can put your life at risk."},
			}}
		w.Log("A family wants the place", f.Name+" is no longer asking for a share of "+target.Name+".", "politics")
		return
	}
	w.Event = &Scene{ID: ID(), Kind: "business_pressure", Source: "authored", Minute: w.Minute, Speaker: speaker, Actor: f.ID, Target: target.ID, Title: "A claim on your earnings", Body: fmt.Sprintf("“%s is doing business under your name now. My people expect a share. Pay for an understanding, or explain why I should tolerate a competitor.”", target.Name), Choices: []Choice{
		{ID: "pay", Label: fmt.Sprintf("Pay $%d for an understanding", w.TheirShare(target.ID)), Cost: w.TheirShare(target.ID),
			Detail: fmt.Sprintf("A week of what %s takes. Improves this family's standing by 12. Postpones the next demand; existing personal threats remain.", target.Name)},
		{ID: "resist", Label: "Refuse their claim", Detail: "Gain 2 respect; lose 20 standing. The family may retaliate against your business. Repeated defiance can put your life at risk."},
		{ID: "ally", Label: "Seek backing from " + rival.Name, Cost: 35, Detail: "Pay $35 for an introduction. Gain 12 standing with their rival; lose 12 with this family. Business retaliation remains possible; a deepening feud can also put your life at risk."},
	}}
	w.Log("A family wants a share", f.Name+" has delivered a demand concerning "+target.Name+".", "politics")
}

// TheirShare is what a family wants for an understanding about one business.
//
// It was sixty dollars, flat, whether the place was a two-room laundry or the
// best casino in the city, and whether you held one address or six. The scene
// has always said "my people expect a share" and the number was a fee — the
// same shape of fault as a card that names a dead man: the words describe a
// rule the code does not have.
//
// A week of what the place takes is a share. It is floored at what the demand
// used to be, so nothing about this makes a family cheaper to deal with, and it
// is capped, because a claim on one shop is a claim on one shop: a family that
// could ask for a month of a casino's takings would be running the business
// rather than leaning on it.
func (w *World) TheirShare(id string) int {
	prop := w.Properties[id]
	if prop == nil {
		return PressureFloor
	}
	return min(PressureCeiling, max(PressureFloor, prop.Income*7))
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
	case "hand":
		return w.HandOver(e.Target, f.ID)
	case "pay":
		share := w.TheirShare(e.Target)
		if buyout, wanted := w.TheirClaim(f.ID, e.Target); wanted {
			share = buyout
		}
		f.Cash += share
		f.Goodwill = min(100, f.Goodwill+12)
		w.NextPressure = w.Minute + 1440
		what := "a week of what the place takes"
		if _, wanted := w.TheirClaim(f.ID, e.Target); wanted {
			what = "six times an understanding, because they had come for the place"
		}
		w.Log("A purchased understanding", fmt.Sprintf("You pay $%d to %s, which is %s. The next business demand is postponed for a day; personal threats are unchanged.", share, f.Name, what), "politics")
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
	// A room that saw them coming is a room that is ready for them: the
	// shutters are down, the stock is out the back, and there are people
	// standing in it. Knowing used to be worth a sentence in the report — "this
	// matches the operation your sources uncovered" — and nothing else at all,
	// which made both investigating and the word from your own counter flavour.
	if p.Known {
		if hands := len(w.Properties[p.Target].Hands); hands > 0 {
			damage = max(damage/3, damage-ReadyFor*hands)
		}
	}
	prop := w.Properties[p.Target]
	damage = min(damage, prop.Condition)
	prop.Condition -= damage
	// Who turned it away, and by how much. This used to name the player's first
	// crew member without checking there was one, because crew were the only
	// thing that could reduce the damage. The people behind the counter can
	// too, and the first business that defended itself without a crew panicked
	// the game.
	defense := ""
	if prevented := unguarded - damage; prevented > 0 {
		switch {
		case len(w.Player.Crew) > 0 && len(w.Tasks) == 0 && w.Player.Crew[0].Loyalty >= 50:
			defense = fmt.Sprintf(" %s helped defend the business, preventing %d additional condition loss.", w.Player.Crew[0].Name, prevented)
		case len(w.Properties[p.Target].Hands) > 0:
			defense = fmt.Sprintf(" The place was expecting them and %s was standing in it, which cost the attack %d condition.",
				w.HandNames(p.Target), prevented)
		default:
			defense = fmt.Sprintf(" Something turned %d of it away.", prevented)
		}
	}
	evidence := "The attackers left no proof of who sent them."
	if p.Known {
		evidence = "This matches the operation your sources uncovered."
	}
	w.Witness("attack", p.Target, fmt.Sprintf("The attack damaged %s. Condition is now %d%%.", l.Name, prop.Condition), "")
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

// Claimants decides which family collects in each district of the city.
//
// This was two slice indices: family zero for the player's first district,
// family one for anything beyond it. Dissolve removes a family from that
// slice, and destroying families is most of what the living world does —
// seventeen of them across twenty long campaigns. So the indices stopped
// meaning the families they were written for as soon as the city did its job,
// and with one family left the second index was not an index into anything: a
// player who owned earning premises outside their first district and had seen
// one family fall crashed the city on the next collection.
//
// Two things have to hold. Ground is the honest answer to whose street this
// is, so a family holding premises in a district collects there. And separate
// districts need separate claimants, or buying one family off would silence
// every demand in the city — which is the property the truce tests state. So a
// district nobody holds goes to the strongest family not already collecting
// from the player, and only falls back to a family that is when there is
// nobody else. When there are no families at all, nobody collects.
func (w *World) Claimants() map[int]*Faction {
	mine := w.PlayerOrganizationID()
	candidates := []*Faction{}
	for i := range w.Factions {
		if w.Factions[i].ID != mine {
			candidates = append(candidates, &w.Factions[i])
		}
	}
	sort.SliceStable(candidates, func(a, b int) bool {
		if candidates[a].Power != candidates[b].Power {
			return candidates[a].Power > candidates[b].Power
		}
		return candidates[a].ID < candidates[b].ID
	})
	claims := map[int]*Faction{}
	if len(candidates) == 0 {
		return claims
	}
	districts := map[int]bool{}
	for i := range Locations {
		districts[Locations[i].District] = true
	}
	ordered := []int{}
	for d := range districts {
		ordered = append(ordered, d)
	}
	sort.Ints(ordered)
	taken := map[string]bool{}
	unheld := []int{}
	for _, d := range ordered {
		var best *Faction
		bestHeld := 0
		for _, f := range candidates {
			held := 0
			for i := range Locations {
				if Locations[i].District != d {
					continue
				}
				if prop := w.Properties[Locations[i].ID]; prop != nil && prop.Owner == f.ID {
					held++
				}
			}
			// candidates is already strongest-first, so a plain > keeps the
			// stronger family on a tie without a second comparison here.
			if held > bestHeld {
				best, bestHeld = f, held
			}
		}
		if best != nil {
			claims[d], taken[best.ID] = best, true
			continue
		}
		unheld = append(unheld, d)
	}
	for _, d := range unheld {
		for _, f := range candidates {
			if !taken[f.ID] {
				claims[d], taken[f.ID] = f, true
				break
			}
		}
		if claims[d] == nil {
			claims[d] = candidates[0]
		}
	}
	return claims
}
