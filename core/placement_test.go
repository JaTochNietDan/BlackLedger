package core

import (
	"strings"
	"testing"
)

// Three complaints from the inbox, all of them about work being labelled or
// placed wrong rather than about what it does.

func shopper(t *testing.T) *World {
	t.Helper()
	w := New(19)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	return w
}

// "Is buying a business called 'establish protection'? That's not super clear
// what that means, we should fix the wording on that to explain what that
// actually entails."
func TestBuyingABusinessSaysThatIsWhatItIs(t *testing.T) {
	w := shopper(t)
	w.Player.Location = "laundry"
	a := actionByID(w.Actions("laundry"), "acquire")
	if a == nil {
		t.Fatal("the laundry is not for sale")
	}
	if strings.Contains(strings.ToLower(a.Label), "protection") {
		t.Fatalf("buying premises is still called %q", a.Label)
	}
	place, _ := PlaceByID("laundry")
	if !strings.Contains(a.Label, place.Name) && !strings.Contains(a.Label, "laundry") {
		t.Fatalf("the label does not say what is being bought: %q", a.Label)
	}
	// And it says what you are taking on, not only what you will earn.
	for _, want := range []string{"yours", "a day"} {
		if !strings.Contains(strings.ToLower(a.Detail), want) {
			t.Fatalf("the description never mentions %q: %q", want, a.Detail)
		}
	}
}

// "Why does it seem like you can send Leo Carver on collections in practically
// every single building's action menu?"
func TestOrdersToYourOwnPeopleFollowYou(t *testing.T) {
	w := shopper(t)
	w.Player.Location = "bar"
	w.Player.Crew = []Crew{{"street-24", "Dita Toth", 60}}
	for _, id := range []string{"delegate", "crew_bonus"} {
		a := actionByID(w.Actions("bar"), id)
		if a == nil {
			t.Fatalf("%s is not offered at all", id)
		}
		if !a.Anywhere {
			t.Fatalf("%s is filed as the bar's own work, so it is in every building's list", id)
		}
	}
}

// "I don't think 'moving against X business yourself' should required respect,
// that doesn't make sense."
func TestGoingInYourselfAsksNothingAboutYourName(t *testing.T) {
	w := shopper(t)
	w.Player.Respect = 0
	w.Player.Crew = []Crew{{"street-24", "Dita Toth", 60}}
	held := ""
	for _, l := range Locations {
		if _, ok := w.SabotageTarget(l.ID); ok {
			held = l.ID
			break
		}
	}
	if held == "" {
		t.Skip("no rival holding in this city")
	}
	if reason := w.SabotageReadiness(held); strings.Contains(strings.ToLower(reason), "respect") {
		t.Fatalf("going in yourself still asks for a reputation: %s", reason)
	}
	// Asking somebody else to do it on your account still does.
	if reason := w.SendAgainstReadiness(held); !strings.Contains(strings.ToLower(reason), "respect") {
		t.Fatalf("sending your crew asks nothing about your name: %q", reason)
	}
	w.Player.Respect = 40
	if reason := w.SendAgainstReadiness(held); strings.Contains(strings.ToLower(reason), "respect") {
		t.Fatalf("a known name was still refused: %s", reason)
	}
}
