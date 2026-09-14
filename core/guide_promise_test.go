package core

import (
	"strings"
	"testing"
)

// The guide page says of itself: "This page is not written down anywhere. It
// asks the game the same question the buttons ask, so it cannot tell you
// something the rules do not."
//
// That is a promise, and it is testable. When the guide says a step is open,
// the work it describes has to be something the player can actually press
// somewhere in the city. When it gives a reason instead, nothing of that kind
// should be live.

// what each step points at, by action id or prefix.
var guideMeans = map[string][]string{
	"Somewhere to start":            {"courier"},
	"Somebody who knows people":     {"contact"},
	"Premises of your own":          {"acquire"},
	"A name of your own":            {}, // becomes true on its own
	"People who answer to you":      {"sign:", "recruit"},
	"Somebody on the door":          {"post"},
	"Money on the street":           {"lend:"},
	"Something of your own to sell": {"still"},
	"Somebody in the building":      {"retain:"},
	"An understanding":              {"pact:"},
	"Somebody else's ladder":        {"serve:"},
	"The chair":                     {"takeover"},
}

func liveSomewhere(w *World, prefixes []string) (bool, string) {
	best := ""
	for i := range Locations {
		id := Locations[i].ID
		w.Player.Location = id
		for _, a := range w.Actions(id) {
			for _, want := range prefixes {
				if a.ID != want && !strings.HasPrefix(a.ID, want) {
					continue
				}
				if !a.Disabled {
					return true, ""
				}
				if best == "" {
					best = a.ID + " at " + id + ": " + a.Reason
				}
			}
		}
	}
	return false, best
}

func TestTheGuideDoesNotOfferWhatTheButtonsRefuse(t *testing.T) {
	t.Parallel()
	// A player who has made a start and has almost nothing left, which is the
	// ordinary state of a campaign at the end of a bad week.
	for _, cash := range []int{51, 400, 4000} {
		w := New(19)
		w.District = 2
		w.Player.Cash, w.Player.Respect = cash, 33
		w.Player.Location = "room"
		checked := 0
		for _, step := range w.Guide() {
			if !step.Open {
				continue
			}
			prefixes, known := guideMeans[step.Title]
			if !known {
				t.Errorf("the guide offers %q and this test does not know what it points at", step.Title)
				continue
			}
			if len(prefixes) == 0 {
				continue
			}
			checked++
			if live, refused := liveSomewhere(w, prefixes); !live {
				t.Errorf("with $%d the guide says %q can be done now, and every button for it is refused — %s",
					cash, step.Title, refused)
			}
		}
		t.Logf("$%d: checked %d open steps", cash, checked)
	}
}

// The other half. The guide describes a set of candidates — every person, every
// premises, every family — and when none of them will do it reports the
// briefest refusal they gave. Each of those searches used to be seeded with an
// invented sentence for the empty case, and those sentences are short, so they
// beat every real reason: a player with no organization was told "There is
// nobody in this city for that yet" while the buttons in front of them said
// "Nobody signs on with one person. They sign on with something that has a
// name."
//
// So: the made-up sentence may only appear when there was genuinely nothing to
// ask about.
func TestTheGuideOnlyInventsAReasonWhenThereIsNothingToAsk(t *testing.T) {
	t.Parallel()
	invented := map[string][]string{
		"There is nobody in this city for that yet":         nil,
		"There is nowhere in this city for that yet":        nil,
		"There is nobody to reach an understanding with":    nil,
		"Nobody is taking anybody on":                       nil,
		"Nobody in that building will take a call from you": nil,
	}
	for _, cash := range []int{51, 400, 4000} {
		w := New(19)
		w.District = 2
		w.Player.Cash, w.Player.Respect = cash, 33
		people, families := len(w.People()), len(w.Factions)
		for _, step := range w.Guide() {
			if _, made := invented[step.Reason]; !made {
				continue
			}
			switch step.Reason {
			case "There is nobody in this city for that yet":
				if people > 0 {
					t.Errorf("with $%d the guide says %q for %q, and there are %d people to ask about",
						cash, step.Reason, step.Title, people)
				}
			case "There is nobody to reach an understanding with", "Nobody is taking anybody on":
				if families > 0 {
					t.Errorf("with $%d the guide says %q for %q, and there are %d families",
						cash, step.Reason, step.Title, families)
				}
			case "There is nowhere in this city for that yet":
				t.Errorf("with $%d the guide says %q for %q, and the city has %d places",
					cash, step.Reason, step.Title, len(Locations))
			case "Nobody in that building will take a call from you":
				t.Errorf("with $%d the guide says %q for %q, and the city hall has officials in it",
					cash, step.Reason, step.Title)
			}
		}
	}
}

// And the helper itself, both ways.
func TestTheGuideReportsTheBriefestRealRefusal(t *testing.T) {
	t.Parallel()
	if got := shortest("nothing to ask", []string{"a long refusal indeed", "short one"}); got != "short one" {
		t.Errorf("got %q, want the briefest real refusal", got)
	}
	if got := shortest("nothing to ask", []string{"a refusal", "", "another"}); got != "" {
		t.Errorf("one candidate was willing and it reported %q", got)
	}
	if got := shortest("nothing to ask", nil); got != "nothing to ask" {
		t.Errorf("with nothing to ask it reported %q", got)
	}
	// The fault: a fallback shorter than every real reason must not win.
	if got := shortest("no", []string{"a real and rather long reason"}); got == "no" {
		t.Error("the invented sentence beat a real refusal")
	}
}
