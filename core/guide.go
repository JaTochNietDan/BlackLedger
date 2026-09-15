package core

import (
	"fmt"
	"strings"
)

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

// Only eligible, owned premises can explain the next business improvement.
// Refusals from unrelated addresses must not hide the remaining prerequisite.
func (w *World) ownedPremisesReason(check func(string) string, eligible func(string) bool, missing string) string {
	reasons := []string{}
	for _, place := range Locations {
		prop := w.Properties[place.ID]
		if place.District > w.District || !w.Own(place.ID) || prop == nil || prop.Income <= 0 || (eligible != nil && !eligible(place.ID)) {
			continue
		}
		reasons = append(reasons, check(place.ID))
	}
	return shortest(missing, reasons)
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

	// Before any rule is asked: can this player act at all? The readiness
	// functions answer about a rule — whether a still can be built, whether
	// somebody would take a loan — and know nothing about where the player is
	// or whether they are alive. So a man locked in a cell was told he could
	// carry envelopes across town, while the buttons beside him offered three
	// things: sit it out, pay a lawyer, or name somebody. The guide's whole
	// promise is that it cannot say what the rules do not, and this is the one
	// question the rules ask first.
	stopped := ""
	switch {
	case !p.Alive:
		stopped = "This life is over"
	case w.Held():
		stopped = fmt.Sprintf("You are being held at Ward Street Station, %s", plural(w.DaysLeft(), "day to go", "days to go"))
	}

	step := func(title, what, reason string, done bool) Step {
		if stopped != "" {
			return Step{Title: title, What: what, Open: false, Reason: stopped, Done: done}
		}
		return Step{Title: title, What: what, Open: reason == "" && !done, Reason: reason, Done: done}
	}

	steps := []Step{
		step("Somewhere to start", fmt.Sprintf("Carry envelopes at Saint Agnes for $%d and %d respect. It takes %d minutes and no upfront cash.", CourierPay, CourierRespect, CourierMinutes),
			w.CourierReadiness(), p.JobCount > 0),
		step("Somebody who knows people", "Buy "+w.RoleName("fixer")+" a coffee. Contacts are how you hear that somebody is coming before they arrive.",
			need(p.Contacts >= 5, "Your information network is fully developed"), p.Contacts > 1),
		step("Premises of your own", "A business can earn while you are elsewhere. Staff, stock, rent and repairs determine what you keep.",
			w.acquisitionReason(), owned > 0),
		step("A name of your own", fmt.Sprintf("Own a business and earn %d respect, then form your family there and choose it as headquarters.", OrganizationStanding),
			w.incorporationReason(), w.Incorporated()),
		step("People who answer to you", "Hire a driver to start. Once your organization has a name, you can sign on more people. Your people add to what you are worth in a fight.",
			w.guideHiringReason(), len(p.Crew) > 0 || len(w.OwnPeople()) > 0),
		step("Somebody on the door", "One of your people, standing at a business. Harder to rob, harder to take, and they are the one standing in it when somebody comes. They have to walk there first, and the door is worth nothing until they arrive.",
			w.ownedPremisesReason(w.PostReadiness, nil, "First acquire a business of your own"), w.anyPosted()),
		step("Money on the street", fmt.Sprintf("Lend at %d%% over %d days. It is the oldest business this trade has.", int(LoanRate*100), LoanTermDays),
			firstOpenPerson(w, w.LendReadiness), len(w.Book()) > 0),
		step("Something of your own to sell", "A still turns a business into a source. It is the most profitable thing premises can do and the most dangerous.",
			w.ownedPremisesReason(w.StillReadiness, StillSite, "First acquire Bluebird Laundry or Russo Motor Works"), w.anyStill()),
		step("Somebody in the building", "An official on a retainer. Files go to the bottom of piles, or licences arrive, or the paper does not run it.",
			w.anyRetainerReason(), len(p.Retainers) > 0),
		step("An understanding", "Stand with an organization. Protection bought with money and paid for in enemies.",
			w.anyPactReason(), len(w.Pacts) > 0),
		step("Somebody else's ladder", "Go to work for a family instead of building your own. A slower living that needs no capital.",
			w.anyServiceReason(), p.Serves != ""),
		step("The chair", "Move on whoever runs the organization you answer to. What it wins is not a promotion; it is the organization.",
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
	reasons := []string{}
	for _, n := range w.People() {
		reasons = append(reasons, check(n.ID))
	}
	return shortest("There is nobody in this city for that yet", reasons)
}

func (w *World) anyPosted() bool {
	for _, l := range Locations {
		if w.Own(l.ID) && w.PostedAt(l.ID) != nil {
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
	reasons := []string{}
	for _, o := range officials {
		reasons = append(reasons, w.RetainerReadiness(o.ID))
	}
	return shortest("Nobody in that building will take a call from you", reasons)
}

func (w *World) anyPactReason() string {
	reasons := []string{}
	for i := range w.Factions {
		reasons = append(reasons, w.PactReadiness(w.Factions[i].ID))
	}
	return shortest("There is nobody to reach an understanding with", reasons)
}

func (w *World) anyServiceReason() string {
	reasons := []string{}
	for i := range w.Factions {
		reasons = append(reasons, w.ServeReadiness(w.Factions[i].ID))
	}
	return shortest("Nobody is taking anybody on", reasons)
}

// shortest is how the guide reports on a set of candidates: the briefest real
// refusal any of them gave, and the fallback only when there was nothing to
// ask about at all.
//
// Each of these used to seed the search with its fallback sentence and then
// keep the shortest string. The fallbacks are short, so they beat every real
// reason: the guide told a player with no organization "There is nobody in
// this city for that yet" while the buttons in front of them said "Nobody
// signs on with one person. They sign on with something that has a name."
// That is exactly what the page promises it cannot do.
func shortest(nothing string, reasons []string) string {
	best := ""
	for _, r := range reasons {
		if r == "" {
			return ""
		}
		if best == "" || len(r) < len(best) {
			best = r
		}
	}
	if best == "" {
		return nothing
	}
	return best
}

// acquisitionReason is the premises step's refusal. It asks the same question
// the button asks, but only of the places that are businesses at all: a rented
// room and the police station both refuse, and "There is nothing here to take
// over" is a sentence about a place, which the guide is not standing in.
func (w *World) acquisitionReason() string {
	reasons := []string{}
	for i := range Locations {
		id := Locations[i].ID
		if prop := w.Properties[id]; prop == nil || prop.Income <= 0 {
			continue
		}
		reasons = append(reasons, w.AcquireReadiness(id))
	}
	return shortest("There is nothing in this city to take over yet", reasons)
}

// Both the first associate and later organization members are real hiring
// routes. Inspect actual destination actions so a travelling/unavailable
// candidate, locked district or lack of cash cannot become an open promise.
func (w *World) guideHiringReason() string {
	reasons := []string{}
	view := *w
	for _, place := range Locations {
		if place.District > w.District {
			continue
		}
		view.Player.Location = place.ID
		for _, a := range view.Actions(place.ID) {
			if a.ID != "recruit" && !strings.HasPrefix(a.ID, "sign:") {
				continue
			}
			if !a.Disabled {
				return ""
			}
			reasons = append(reasons, a.Reason)
		}
	}
	return shortest("There is nobody available to hire yet", reasons)
}
