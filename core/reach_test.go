package core

import "testing"

// The action sweep added when people started walking asked one question — is
// this person out on the street? — and every other way of being unreachable
// went on being ignored. A player could send a man on collections from the far
// side of the city, order one out of a police cell, and pay a bonus to a man
// who had been shot the day before: the crew list kept him, the button offered
// him work, and the command was accepted.

func TestYouCannotSendAManWhoIsDead(t *testing.T) {
	w, leo := crewman(t)
	w.Player.Cash = 3000
	w.Kill(leo.ID, "Shot at the counter.")
	if len(w.Player.Crew) != 0 {
		t.Fatalf("a dead man is still on the books: %+v", w.Player.Crew)
	}
	for _, a := range w.Actions(w.Player.Location) {
		if a.ID == "delegate" || a.ID == "crew_bonus" {
			t.Fatalf("%q is offered for somebody who is dead", a.Label)
		}
	}
	if _, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "delegate", Target: w.Player.Location}); err == nil {
		t.Fatal("a dead man was sent on collections")
	}
}

func TestYouCannotSendAManOutOfACell(t *testing.T) {
	w, leo := crewman(t)
	w.Player.Cash = 3000
	leo.Held = w.Minute + 2880
	for _, a := range w.Actions(w.Player.Location) {
		if (a.ID == "delegate" || a.ID == "crew_bonus") && !a.Disabled {
			t.Fatalf("%q is offered for a man the police are holding", a.Label)
		}
		if a.ID == "delegate" && !contains(a.Reason, "held") {
			t.Fatalf("the refusal does not say where he is: %q", a.Reason)
		}
	}
}

// Being in a different building is deliberately not out of reach. A man at an
// address can be reached; a man between two addresses is nowhere. This is here
// so that the distinction is a decision on the record rather than an accident,
// and so that reinstating the stricter rule has to be done deliberately.
func TestBeingInAnotherBuildingIsNotOutOfReach(t *testing.T) {
	w, leo := crewman(t)
	w.Player.Cash = 3000
	leo.Location = "club"
	w.Player.Location = "bar"
	if reason := w.OutOfReach(leo.ID); reason != "" {
		t.Fatalf("a man standing at an address across the city is unreachable: %q", reason)
	}
	// But once he steps out of it he is nowhere.
	leo.Heading, leo.Arrives, leo.Errand = "bar", w.Minute+30, "on his way"
	if reason := w.OutOfReach(leo.ID); reason == "" {
		t.Fatal("a man on the street can still be dealt with")
	}
}

// Taking the dead off the books uncovered an older bug it had been hiding:
// "Recruit Leo Carver" was refused only because he was already in the crew, so
// the moment death removed him the button offered to hire him again.
func TestYouCannotRecruitAManWhoIsDead(t *testing.T) {
	w, leo := crewman(t)
	w.Player.Cash = 3000
	w.Kill(leo.ID, "Shot twice outside the Mariner.")
	for _, a := range w.Actions("bar") {
		if a.ID != "recruit" {
			continue
		}
		if !a.Disabled {
			t.Fatal("the game offered to recruit a dead man")
		}
		if !contains(a.Reason, "dead") {
			t.Fatalf("it refuses for the wrong reason: %q", a.Reason)
		}
		return
	}
	t.Fatal("recruiting is not offered at the bar at all")
}
