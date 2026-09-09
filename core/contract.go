package core

import (
	"fmt"
	"strings"
)

// A name and a price. Anyone alive in this city can be killed, and the same
// arrangement works against the player: a contract is a plot like any other,
// and the rules do not care who paid for it.
//
// Nothing is decided when the money changes hands. The contract is committed
// with a due time, and it resolves later, which is what makes a hit something
// you can be warned about, run from, or answer.

type Tier struct {
	ID    string
	Label string
	// Noun is the tier in the middle of a sentence, with whatever article it
	// needs. Label is a button and reads "A professional"; splicing that into
	// prose printed "The a professional you paid $2501 to reach Elena Russo".
	Noun    string
	Detail  string
	Skill   int     // how good they are
	Rate    float64 // multiplier on the target's price
	Capture float64 // chance of being taken alive after a failure
}

var contractTiers = []Tier{
	{ID: "cheap", Label: "Someone who needs the money", Noun: "the one who needed the money",
		Detail: "Cheap and unreliable. If it goes wrong they are very likely to be caught, and they will not hold their tongue.",
		Skill:  35, Rate: 1, Capture: .5},
	{ID: "professional", Label: "A professional", Noun: "the professional",
		Detail: "Costs more and usually finishes. If it goes wrong there is a fair chance they are taken alive.",
		Skill:  55, Rate: 2.2, Capture: .25},
	{ID: "specialist", Label: "A specialist", Noun: "the specialist",
		Detail: "Expensive, quiet, and rarely caught. Ask what that costs before you decide it is worth it.",
		Skill:  78, Rate: 4.5, Capture: .08},
}

func contractTier(id string) (Tier, bool) {
	for _, t := range contractTiers {
		if t.ID == id {
			return t, true
		}
	}
	return Tier{}, false
}

// Contract is a commissioned killing awaiting its moment. It is private: the
// public projection never exposes who is coming for whom.
type Contract struct {
	ID     string `json:"id"`
	Life   int    `json:"life"`
	Target string `json:"target"`
	Tier   string `json:"tier"`
	Payer  string `json:"payer"`
	Due    int    `json:"due"`
	Fee    int    `json:"fee"`
}

// ContractPrice is what it costs to have someone killed. Standing, competence
// and the strength of whoever protects them all raise the price.
func (w *World) ContractPrice(target string, tier Tier) int {
	person := w.NPC(target)
	if person == nil {
		return 0
	}
	base := 120 + person.Rank*6 + person.Skill*3
	if f := w.faction(person.Faction); f != nil {
		base += f.Power * 4
	}
	return int(float64(base) * tier.Rate)
}

// ContractTargets are the people the player could put a price on: everyone
// alive except themselves and their own crew.
// ContractTargets is who the player could put a name on. It used to be
// everybody alive except their own crew, which in a city of fifty rendered as
// fifty identical rows in a single scene — a wall rather than a decision, and
// a wall that included a laundress the player has never heard of.
//
// A name is something you have a reason to say. So it is the people this
// protagonist actually knows, plus anybody who has given them a reason
// whether they know them or not: somebody carrying a grudge against them,
// somebody who owes them and has stopped paying. Their own people are never on
// it — that is what dismissing somebody is for.
//
// Ordered by how much reason there is, so the name the player is most likely
// to be thinking of is the first one they read.
func (w *World) ContractTargets() []*NPC {
	reason := func(n *NPC) int {
		switch {
		case n.Sore >= 30:
			return 0 // they have made themselves a problem
		case w.LoanTo(n.ID) != nil && w.LoanTo(n.ID).Missed > 0:
			return 1 // they owe you and have stopped answering
		case n.Rank >= RankLeader:
			return 2 // everybody knows who runs things
		case IsOfficial(n.ID) || w.isRoleHolder(n):
			return 3
		case w.Known(n):
			return 4
		}
		return -1
	}
	out := []*NPC{}
	for _, n := range w.People() {
		if n.Faction == w.PlayerOrganizationID() {
			continue // not your own people
		}
		own := false
		for _, c := range w.Player.Crew {
			if c.ID == n.ID {
				own = true
			}
		}
		if !own && reason(n) >= 0 {
			out = append(out, n)
		}
	}
	for i := range out {
		for j := i + 1; j < len(out); j++ {
			if reason(out[j]) < reason(out[i]) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

// OpenContract offers the names available. Choosing one asks what class of
// person should do it, so the price is seen before anything is committed.
func (w *World) OpenContract() {
	choices := []Choice{}
	for _, n := range w.ContractTargets() {
		where := "somewhere in the city"
		if place, ok := PlaceByID(n.Location); ok {
			where = "usually found at " + place.Name
		}
		role := n.Role
		if f := w.faction(n.Faction); f != nil {
			role = n.Role + ", " + f.Name
		}
		choices = append(choices, Choice{
			ID:     "mark:" + n.ID,
			Label:  n.Name,
			Detail: role + ", " + where + ".",
		})
	}
	choices = append(choices, Choice{ID: "leave", Label: "Say nothing", Detail: "No name is given. No money changes hands."})
	w.Event = &Scene{
		ID: ID(), Kind: "contract", Source: "authored", Minute: w.Minute,
		Title:   "A name, and what it costs",
		Body:    "“Everyone in this city is a price. Say a name and I will tell you what it takes. After that it is not your arrangement any more, and it is not mine either.”",
		Speaker: w.HolderID("fixer"), Choices: choices,
	}
}

// openContractTerms is the second step: the mark is chosen, and now the class
// of person doing the work decides both the price and the risk.
func (w *World) openContractTerms(target string) error {
	person := w.NPC(target)
	if person == nil || person.Dead {
		return fmt.Errorf("that person is beyond anyone's reach now")
	}
	choices := []Choice{}
	for _, tier := range contractTiers {
		fee := w.ContractPrice(target, tier)
		choices = append(choices, Choice{
			ID: "hire:" + tier.ID, Label: tier.Label, Cost: fee,
			Detail: fmt.Sprintf("$%d. %s", fee, tier.Detail),
		})
	}
	choices = append(choices, Choice{ID: "leave", Label: "Think about it", Detail: "No money changes hands and no name is passed on."})
	w.Event = &Scene{
		ID: ID(), Kind: "contract_terms", Source: "authored", Minute: w.Minute,
		Target: target, Speaker: w.HolderID("fixer"),
		Title:   "What it takes to reach " + person.Name,
		Body:    fmt.Sprintf("“%s. That can be arranged. Who does it decides what it costs, and what happens if it goes wrong.”", person.Name),
		Choices: choices,
	}
	return nil
}

// Commission books the killing. The money is gone from this moment whether or
// not it works.
func (w *World) Commission(target, tier string) error {
	person := w.NPC(target)
	if person == nil || person.Dead {
		return fmt.Errorf("that person is beyond anyone's reach now")
	}
	t, ok := contractTier(tier)
	if !ok {
		return fmt.Errorf("nobody like that is available")
	}
	// The fee is charged by the choice that reached here, which carries it as
	// its cost. Charging again here would take it twice.
	fee := w.ContractPrice(target, t)
	w.Contracts = append(w.Contracts, Contract{
		ID: ID(), Life: w.Life, Target: target, Tier: tier, Payer: "player",
		Due: w.Minute + 180 + int(w.Random()*600), Fee: fee,
	})
	w.Player.Heat = min(100, w.Player.Heat+2)
	w.Log("A name passed on", fmt.Sprintf("$%d for %s. It happens when it happens, and not because you are watching.", fee, person.Name), "danger")
	return nil
}

// killingMethod is how a commissioned killing was done. The manner comes from
// the same place every other death in this city gets it; what the tier adds is
// how much of a mess was left behind, which is most of what the fee buys.
func (w *World) killingMethod(person *NPC, tier Tier) string {
	manner := w.Manner(person, nil)
	switch tier.ID {
	case "cheap":
		return manner + " Whoever did it left a great deal behind them."
	case "professional":
		return manner + " Whoever did it had done it before."
	case "specialist":
		return manner + " There is nothing to find and there was never going to be."
	}
	return manner
}

// ResolveContracts settles every commissioned killing whose moment has come.
// Called from the clock, so a hit lands while the player is doing something
// else, which is how it should feel.
func (w *World) ResolveContracts() {
	remaining := w.Contracts[:0]
	for _, c := range w.Contracts {
		if c.Life != w.Life || c.Due > w.Minute {
			remaining = append(remaining, c)
			continue
		}
		w.resolveContract(c)
	}
	w.Contracts = remaining
}

func (w *World) resolveContract(c Contract) {
	person := w.NPC(c.Target)
	tier, _ := contractTier(c.Tier)
	if person == nil || person.Dead {
		return // somebody else got there first
	}
	// Standing and competence protect a person; so does the strength of whoever
	// they answer to.
	defence := person.Skill/2 + person.Rank/4
	if f := w.faction(person.Faction); f != nil {
		defence += f.Power / 5
	}
	odds := .3 + float64(tier.Skill-defence)/120
	if odds < .08 {
		odds = .08
	}
	if odds > .92 {
		odds = .92
	}

	if w.WorldRandom() < odds {
		method := w.killingMethod(person, tier)
		w.Kill(c.Target, method)
		if c.Payer == "player" {
			w.Player.Heat = min(100, w.Player.Heat+6)
			w.Player.Respect += 3
			w.Log("Your arrangement is settled",
				fmt.Sprintf("%s you paid $%d for reached %s. %s Nobody has connected it to you yet.",
					upper1(tier.Noun), c.Fee, person.Name, method), "danger")
		}
		return
	}

	// A shot that misses is still a shooting. The city reports the attempt
	// without knowing who arranged it.
	w.Report("attempt", "ATTEMPT ON THE LIFE OF "+strings.ToUpper(person.Name),
		fmt.Sprintf("%s survived an attack near %s. Police describe the assault as targeted and say no arrest has been made.",
			person.Name, w.placeName(person.Location)))

	// It failed. Whether it comes back to whoever paid is the real risk.
	if c.Payer == "player" {
		w.Log("Your arrangement failed",
			fmt.Sprintf("%s you paid $%d to reach %s did not finish it. %s is alive and now knows somebody wants them dead.",
				upper1(tier.Noun), c.Fee, person.Name, person.Name), "danger")
	} else {
		w.Log("An attempt that failed", fmt.Sprintf("Someone went for %s and did not finish it. %s knows they are worth killing now.", person.Name, person.Name), "danger")
	}
	if w.WorldRandom() >= tier.Capture {
		return // they got away, and said nothing
	}
	if c.Payer != "player" {
		return
	}
	// Taken alive, and questioned.
	w.Player.Heat = min(100, w.Player.Heat+25)
	w.Report("arrest", "ARREST AFTER ATTACK ON "+strings.ToUpper(person.Name),
		"Somebody taken at the scene is said to be assisting police with their enquiries. Sources suggest they have given a name.")
	if f := w.faction(person.Faction); f != nil {
		f.Goodwill = max(-100, f.Goodwill-45)
		w.RetaliationFrom(f.ID)
		w.Log("They gave up a name", fmt.Sprintf("The one who went for %s was taken alive and questioned. %s %s who paid, and so do the police.", person.Name, f.Name, Agree(f.Name, "knows", "know")), "danger")
		return
	}
	w.Log("They gave up a name", fmt.Sprintf("The one who went for %s was taken alive and questioned. Your name came out of it.", person.Name), "danger")
}

// ConsiderFactionContracts lets organizations buy the same service the player
// can. A family losing a war, or one whose standing with a rival has collapsed,
// will pay to remove the person at the top of the other side. The player is not
// exempt: this is the same arrangement pointed the other way.
func (w *World) ConsiderFactionContracts() {
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.Cash < 400 || w.WorldRandom() >= .04 {
			continue
		}
		for _, c := range w.Conflicts {
			if c.State != "war" || (c.A != f.ID && c.B != f.ID) {
				continue
			}
			enemy := c.A
			if enemy == f.ID {
				enemy = c.B
			}
			rival := w.faction(enemy)
			if rival == nil {
				continue
			}
			// Go for the person at the top of the other side.
			members := w.Members(rival.ID)
			if len(members) == 0 {
				continue
			}
			target := members[0]
			tier := contractTiers[0]
			if f.Cash > 3000 {
				tier = contractTiers[1]
			}
			fee := w.ContractPrice(target.ID, tier)
			if fee > f.Cash {
				continue
			}
			f.Cash -= fee
			w.Contracts = append(w.Contracts, Contract{
				ID: ID(), Life: w.Life, Target: target.ID, Tier: tier.ID, Payer: f.ID,
				Due: w.Minute + 240 + int(w.WorldRandom()*900), Fee: fee,
			})
			return // one arrangement at a time
		}
	}
}

// PendingArrangements is what the player has paid for and not yet seen the end
// of. They know the name; they do not know when, and they never see anyone
// else's arrangements.
func (w *World) PendingArrangements() []map[string]any {
	out := []map[string]any{}
	for _, c := range w.Contracts {
		if c.Payer != "player" || c.Life != w.Life {
			continue
		}
		person := w.NPC(c.Target)
		if person == nil {
			continue
		}
		tier, _ := contractTier(c.Tier)
		out = append(out, map[string]any{
			"target": person.Name,
			"hired":  tier.Label,
			"paid":   c.Fee,
			"status": "Arranged. It happens when it happens.",
		})
	}
	return out
}

func (w *World) placeName(id string) string {
	if place, ok := PlaceByID(id); ok {
		return place.Name
	}
	return "the waterfront"
}
