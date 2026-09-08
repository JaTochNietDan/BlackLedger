package core

import "fmt"

// People in this city rob each other, get passed over for promotion and belong
// to organizations that are at war. Until now none of it left a mark: nobody
// remembered who had done what to them, and nobody ever settled it.
//
// A grudge is one person's memory of one specific thing another person did.
// They accumulate from events the simulation actually committed, they decay,
// and when one gets heavy enough its holder does something about it. What
// follows is a killing the player usually reads about in the paper without
// knowing why, and an organization that now has a reason to hate another one.
//
// This is where private history becomes public war.

// Grudge is one person's reason to want another one dealt with.
type Grudge struct {
	// Holder is who is carrying it, Against is who it is about.
	Holder  string `json:"holder"`
	Against string `json:"against"`
	// Weight is how badly they want it settled; Because is what happened.
	Weight  int    `json:"weight"`
	Because string `json:"because"`
	// Since is when it was last added to, for decay.
	Since int `json:"since"`
}

const (
	// GrudgeActs is the weight at which somebody stops brooding and moves.
	GrudgeActs = 55
	// GrudgeCap is the most one person can hold against another, so a long
	// grudge does not become a certainty.
	GrudgeCap = 90
	// GrudgeFade is what a day takes off every grudge in the city. People let
	// things go, slowly.
	GrudgeFade = 1
	// MaxGrudges bounds the save. The oldest and lightest are forgotten first.
	MaxGrudges = 60
)

// Resent records that one person now has a reason to dislike another. Both must
// be people; nothing here works on the player, who has their own machinery.
func (w *World) Resent(holder, against string, weight int, because string) {
	if holder == "" || against == "" || holder == against || weight <= 0 {
		return
	}
	if w.NPC(holder) == nil || w.NPC(against) == nil {
		return
	}
	for i := range w.Grudges {
		g := &w.Grudges[i]
		if g.Holder == holder && g.Against == against {
			g.Weight = min(GrudgeCap, g.Weight+weight)
			g.Because, g.Since = because, w.Minute
			return
		}
	}
	w.Grudges = append(w.Grudges, Grudge{holder, against, min(GrudgeCap, weight), because, w.Minute})
	w.trimGrudges()
}

// trimGrudges keeps the save bounded by forgetting the lightest first, which is
// also the most human way to lose a grievance.
func (w *World) trimGrudges() {
	for len(w.Grudges) > MaxGrudges {
		lightest := 0
		for i := range w.Grudges {
			if w.Grudges[i].Weight < w.Grudges[lightest].Weight {
				lightest = i
			}
		}
		w.Grudges = append(w.Grudges[:lightest], w.Grudges[lightest+1:]...)
	}
}

// GrudgeDay lets everything cool a little, and drops what has been let go of
// entirely. Called before anybody acts, so a grudge settled today was still
// heavy this morning.
func (w *World) GrudgeDay() {
	kept := w.Grudges[:0]
	for _, g := range w.Grudges {
		if w.NPC(g.Holder) == nil || w.NPC(g.Against) == nil {
			continue
		}
		if holder := w.NPC(g.Holder); holder.Dead {
			continue // the dead carry nothing
		}
		if target := w.NPC(g.Against); target.Dead {
			continue // and nothing is owed by the dead either
		}
		g.Weight -= GrudgeFade
		if g.Weight <= 0 {
			continue
		}
		kept = append(kept, g)
	}
	w.Grudges = kept
}

// warGrudges is the standing reason people in organizations at war have to move
// against each other. It is computed rather than stored, because it stops being
// true the day the war ends.
func (w *World) warGrudge(holder, against *NPC) int {
	if holder.Faction == "" || against.Faction == "" || holder.Faction == against.Faction {
		return 0
	}
	c := w.Conflict(holder.Faction, against.Faction)
	if c == nil {
		return 0
	}
	switch c.State {
	case "war":
		return 35
	case "feud":
		return 15
	}
	return 0
}

// SettleGrudges gives the city's private quarrels a chance to become public
// ones. Ambition decides who is willing to act; skill decides whether it works.
func (w *World) SettleGrudges() {
	for i := range w.Grudges {
		g := w.Grudges[i]
		holder, target := w.NPC(g.Holder), w.NPC(g.Against)
		if holder == nil || target == nil || holder.Dead || target.Dead {
			continue
		}
		weight := g.Weight + w.warGrudge(holder, target)
		if weight < GrudgeActs {
			continue
		}
		// A loyal person does not move against their own people, whatever they
		// are owed, which is what makes a succession grudge inside a loyal
		// organization something that simply festers.
		if TemperamentOf(holder).Loyal && holder.Faction != "" && holder.Faction == target.Faction {
			continue
		}
		// Wanting it settled is not the same as being the kind of person who
		// settles things. Ambition and temperament are what turn a grievance
		// into a decision.
		if w.WorldRandom() >= w.Nerve(holder) {
			continue
		}
		w.settle(holder, target, g)
		return // one of these a day is plenty for one city
	}
}

// settle is one person moving against another. It is the same arithmetic the
// player faces: skill against skill, with whoever the target answers to
// counting for something.
func (w *World) settle(holder, target *NPC, g Grudge) {
	defence := w.Poise(target) + 10
	if f := w.faction(target.Faction); f != nil {
		defence += f.Power / 4
	}
	attack := w.Poise(holder)
	odds := float64(attack) / float64(attack+defence)

	// Clear the memory either way: it has been acted on now, and what happens
	// next is a new thing rather than the old one.
	w.forget(holder.ID, target.ID)

	if w.WorldRandom() >= odds {
		// It went wrong, which is how the wrong person ends up dead.
		if w.WorldRandom() < .3 {
			w.KillBy(holder.ID, target, fmt.Sprintf("They had gone after %s over %s.", target.Name, g.Because))
			w.hostility(target, holder, 6)
			return
		}
		w.Log("It came to nothing", fmt.Sprintf("%s moved against %s over %s. It did not go their way, and now %s knows.", holder.Name, target.Name, g.Because, target.Name), "politics")
		w.Resent(target.ID, holder.ID, 30, "being moved against")
		w.hostility(holder, target, 8)
		return
	}

	w.KillBy(target.ID, holder, fmt.Sprintf("It was over %s.", g.Because))
	// Whoever the dead answered to now has a reason of their own, and if the
	// two belong to different organizations the city has a new quarrel.
	w.hostility(holder, target, 14)
	for _, peer := range w.Members(target.Faction) {
		if !peer.Dead && peer.ID != holder.ID {
			w.Resent(peer.ID, holder.ID, 25, "what happened to "+target.Name)
			break
		}
	}
}

// forget removes one person's grievance against another once it has been acted
// on.
func (w *World) forget(holder, against string) {
	kept := w.Grudges[:0]
	for _, g := range w.Grudges {
		if g.Holder == holder && g.Against == against {
			continue
		}
		kept = append(kept, g)
	}
	w.Grudges = kept
}

// hostility is how a private quarrel between two people becomes a public one
// between the organizations they answer to. This is the whole point of the
// system: a war the player is not part of, started by something nobody told
// them about.
func (w *World) hostility(actor, victim *NPC, amount int) {
	if actor.Faction == "" || victim.Faction == "" || actor.Faction == victim.Faction {
		return
	}
	before := ""
	if c := w.Conflict(actor.Faction, victim.Faction); c != nil {
		before = c.State
	}
	w.Antagonize(actor.Faction, victim.Faction, amount)
	if c := w.Conflict(actor.Faction, victim.Faction); c != nil && c.State != before {
		a, b := w.factionName(actor.Faction), w.factionName(victim.Faction)
		switch c.State {
		case "war":
			w.Log("It is a war now", fmt.Sprintf("%s and %s are past talking. Nobody outside either organization was told why.", a, b), "politics")
			w.Report("politics", "VIOLENCE BETWEEN TWO FAMILIES ESCALATES",
				fmt.Sprintf("Police attribute a series of incidents to a dispute between interests associated with %s and %s. Neither would comment.", a, b))
		case "feud":
			w.Log("Bad blood", fmt.Sprintf("%s and %s have stopped being civil about something.", a, b), "politics")
		}
	}
}

// GrudgeSummary is what the player can find out about who wants who dealt with.
// Only what somebody would tell them: it takes contacts to hear this at all.
func (w *World) GrudgeSummary() []map[string]any {
	if w.Reach() < 2 {
		return []map[string]any{}
	}
	out := []map[string]any{}
	for _, g := range w.Grudges {
		// Deliberately below the weight at which anybody acts: the useful thing
		// to hear is that two people have a problem, before one of them
		// settles it.
		if g.Weight < 15 {
			continue
		}
		holder, target := w.NPC(g.Holder), w.NPC(g.Against)
		if holder == nil || target == nil || holder.Dead || target.Dead {
			continue
		}
		out = append(out, map[string]any{
			"holder": holder.Name, "against": target.Name, "because": g.Because,
		})
	}
	return out
}
