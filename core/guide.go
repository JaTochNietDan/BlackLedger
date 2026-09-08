package core

import "fmt"

// The Guide was prose written when this game had eight actions, and it rotted.
// It still told the player that "broader autonomous family politics" was future
// work, in a build where two families fight their own wars, hold ground, run
// coups and collapse; and it described death wrongly, since what a protagonist
// leaves behind became an estate under one of their own people. A guide that
// misdescribes the game is worse than no guide.
//
// So this one is not prose about the game. It is the game reporting on itself:
// the things a player can be doing, each answered by the same readiness
// function that answers the button. It cannot go stale, because if the rule
// changes the guide changes with it.

// Step is one thing a player might be doing, and whether they can yet.
type Step struct {
	Title string `json:"title"`
	// What it is, in a sentence.
	What string `json:"what"`
	// Open is whether it is available now; Reason is why not, in the game's own
	// words, when it is not.
	Open   bool   `json:"open"`
	Reason string `json:"reason,omitempty"`
	// Done is whether the player is already past it.
	Done bool `json:"done,omitempty"`
}

// need is the same shape the action list uses: a reason when something is in
// the way, and nothing when it is not.
func need(blocked bool, reason string) string {
	if blocked {
		return reason
	}
	return ""
}

// firstOpen runs a readiness function against every place in the city and
// returns the friendliest answer: empty if it can be done somewhere, otherwise
// the reason given where it came closest to being possible.
func firstOpen(check func(string) string) string {
	reason := "There is nowhere in this city for that yet"
	for _, l := range Locations {
		if r := check(l.ID); r == "" {
			return ""
		} else if len(r) < len(reason) {
			reason = r
		}
	}
	return reason
}

// Guide is where the player stands, in the order these things usually happen.
func (w *World) Guide() []Step {
	p := &w.Player
	owned := 0
	for _, l := range Locations {
		if w.Own(l.ID) {
			owned++
		}
	}

	step := func(title, what, reason string, done bool) Step {
		return Step{Title: title, What: what, Open: reason == "" && !done, Reason: reason, Done: done}
	}

	steps := []Step{
		step("Somewhere to start", fmt.Sprintf("Carry envelopes at Saint Agnes for $%d and %d respect. It costs nothing but the hour.", CourierPay, CourierRespect),
			"", p.JobCount > 0),
		step("Somebody who knows people", "Buy Mara a coffee. Contacts are how you hear that somebody is coming before they arrive.",
			need(p.Contacts >= 5, "Your information network is fully developed"), p.Contacts > 1),
		step("Premises of your own", "A business earns while you are elsewhere, and is the only income that does.",
			need(owned == 0 && p.Respect < 6, "Earn 6 respect first"), owned > 0),
		step("A name of your own", fmt.Sprintf("Two premises and %d respect and the city files you with the families.", OrganizationStanding),
			w.incorporationReason(), w.Incorporated()),
		step("People who answer to you", "Sign somebody on. They add to what you are worth in a fight and stand in front of what comes at you.",
			firstOpenPerson(w, w.SignOnReadiness), len(w.OwnPeople()) > 0),
		step("Somebody on the door", "One of your people, standing at a business. Harder to rob, harder to take, and they are the one standing in it when somebody comes. They have to walk there first, and the door is worth nothing until they arrive.",
			firstOpen(w.PostReadiness), w.anyPosted()),
		step("Money on the street", fmt.Sprintf("Lend at %d%% over %d days. It is the oldest business this trade has.", int(LoanRate*100), LoanTermDays),
			firstOpenPerson(w, w.LendReadiness), len(w.Book()) > 0),
		step("Something of your own to sell", "A still turns a business into a source. It is the most profitable thing premises can do and the most dangerous.",
			firstOpen(w.StillReadiness), w.anyStill()),
		step("Somebody in the building", "An official on a retainer. Files go to the bottom of piles, or licences arrive, or the paper does not run it.",
			w.anyRetainerReason(), len(p.Retainers) > 0),
		step("An understanding", "Stand with an organization. Protection bought with money and paid for in enemies.",
			w.anyPactReason(), len(w.Pacts) > 0),
		step("Somebody else's ladder", "Go to work for a family instead of building your own. A slower living that needs no capital.",
			w.anyServiceReason(), p.Serves != ""),
		step("The chair", "Move on the man running the organization you answer to. What it wins is not a promotion; it is the organization.",
			w.TakeoverReadiness(), false),
	}
	return steps
}

// The rules that do not change, which are the only part of a guide that can
// safely be written down once.
func GuideRules() []string {
	return []string{
		"The clock only moves when you commit to something. Reading, inspecting and choosing cost nothing.",
		fmt.Sprintf("Attention is public and so are the thresholds. Past %d the police come to the door; past %d they take the premises.", RaidThreshold, ForfeitThreshold),
		"The city does not scale to you. An organization at ninety strength will kill you on your first day if you give it a reason.",
		"Death is permanent. What you built passes to the strongest of your own people and becomes an organization you can deal with, or fight.",
		"Everything anybody in this city does is done by the same rules you use. There is no separate arithmetic for them.",
		"People are not furniture. They walk between buildings on their own errands, and anybody out on the street is not at the address they left and cannot be dealt with until they arrive — including your own, when you send them somewhere.",
	}
}

func (w *World) incorporationReason() string {
	if w.Incorporated() {
		return ""
	}
	needs := w.PlayerOrganizationDescription()
	if needs == nil {
		return ""
	}
	if list, ok := needs["needs"].([]string); ok && len(list) > 0 {
		return "Still short of " + join(list)
	}
	return ""
}

func firstOpenPerson(w *World, check func(string) string) string {
	reason := "There is nobody in this city for that yet"
	for _, n := range w.People() {
		if r := check(n.ID); r == "" {
			return ""
		} else if len(r) < len(reason) {
			reason = r
		}
	}
	return reason
}

func (w *World) anyPosted() bool {
	for _, l := range Locations {
		if w.PostedAt(l.ID) != nil {
			return true
		}
	}
	return false
}

func (w *World) anyStill() bool {
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Still && w.Own(l.ID) {
			return true
		}
	}
	return false
}

func (w *World) anyRetainerReason() string {
	reason := "Nobody in that building will take a call from you"
	for _, o := range officials {
		if r := w.RetainerReadiness(o.ID); r == "" {
			return ""
		} else if len(r) < len(reason) {
			reason = r
		}
	}
	return reason
}

func (w *World) anyPactReason() string {
	reason := "There is nobody to reach an understanding with"
	for i := range w.Factions {
		if r := w.PactReadiness(w.Factions[i].ID); r == "" {
			return ""
		} else if len(r) < len(reason) {
			reason = r
		}
	}
	return reason
}

func (w *World) anyServiceReason() string {
	reason := "Nobody is taking anybody on"
	for i := range w.Factions {
		if r := w.ServeReadiness(w.Factions[i].ID); r == "" {
			return ""
		} else if len(r) < len(reason) {
			reason = r
		}
	}
	return reason
}
