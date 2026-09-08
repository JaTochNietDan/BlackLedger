package core

import (
	"strings"
	"testing"
)

func TestJobLocationIsValidatedDisplayedAndRemembered(t *testing.T) {
	w := New(27)
	p := Proposal{Location: "garage", Title: "A local delivery", Body: "Carry the sealed message.", Speaker: "mara", Operation: "courier", Outcome: "Delivered."}
	if _, err := w.ValidateProposal(p); err == nil {
		t.Fatal("locked district accepted")
	}
	p.Location = "missing"
	if _, err := w.ValidateProposal(p); err == nil {
		t.Fatal("unknown place accepted")
	}
	p.Location = "laundry"
	p.Approaches = []Approach{{Method: "careful", Label: "Check the route"}}
	scene, err := w.ValidateProposal(p)
	if err != nil {
		t.Fatal(err)
	}
	// The venue is a condition of the job rather than of any one way of doing
	// it, so it is stated once above the choices instead of on each of them.
	if !strings.HasPrefix(scene.Conditions, "Bluebird Laundry · ") {
		t.Fatalf("missing displayed venue: %q", scene.Conditions)
	}
	w.Event = scene
	w.RememberArrangement(scene, "offered")
	choice(t, &w, "accept")
	if w.Player.Location != "room" || w.Minute != 525 || w.Player.Cash != 165 {
		t.Fatal("abstract job changed base or charged unquoted travel")
	}
	memory := w.Clone().Arrangements[0]
	if memory.Location != "laundry" || memory.Status != "completed" {
		t.Fatal("job venue lost on completion/save")
	}
	p.Location = "" // Saved/authored legacy proposals remain compatible.
	if _, err := w.ValidateProposal(p); err != nil {
		t.Fatal(err)
	}
}
