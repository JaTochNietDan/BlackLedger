package core

import "testing"

func TestArrangementMemoryTracksPoliceOutcome(t *testing.T) {
	for _, decision := range []string{"pay", "abandon"} {
		w := New(27)
		w.Player.Heat = 14
		scene, _ := w.ValidateProposal(approachProposal())
		w.Event = scene
		choice(t, &w, "accept")
		if len(w.Arrangements) != 1 || w.Arrangements[0].Status != "awaiting_police" {
			t.Fatal("premature completion")
		}
		w = w.Clone()
		choice(t, &w, decision)
		m := w.Arrangements[0]
		expected := "completed"
		if decision == "abandon" {
			expected = "abandoned"
		}
		if len(w.Arrangements) != 1 || m.Status != expected || m.ID != scene.ID || m.Offer != scene.Body || m.Speaker != "mara" {
			t.Fatal("lost original job memory", m)
		}
		if (m.Result != "") != (decision == "pay") {
			t.Fatal("unearned result remembered")
		}
	}
}

func TestDirectorVarietyIncludesPendingAndDeclinedWork(t *testing.T) {
	w := New(27)
	if w.NextDirectorOperation() != "mediation" {
		t.Fatal("unexpected opening brief")
	}
	p := approachProposal()
	p.Operation = "mediation"
	w.Event, _ = w.ValidateProposal(p)
	choice(t, &w, "decline")
	if w.Arrangements[0].Status != "declined" || w.Arrangements[0].Result != "" {
		t.Fatal("decline became a completion")
	}
	if w.NextDirectorOperation() != "collection" {
		t.Fatal("repeated declined operation")
	}
	p.Operation = "collection"
	e, _ := w.ValidateProposal(p)
	w.Offers = append(w.Offers, Offer{Ready: w.Minute + 30, Event: e})
	if w.Clone().NextDirectorOperation() != "courier" {
		t.Fatal("pending operation ignored")
	}
	w.Die("test")
	act(t, &w, "new_life", "")
	if len(w.Arrangements) != 1 || w.Arrangements[0].Life != 1 || w.NextDirectorOperation() != "mediation" {
		t.Fatal("cross-life memory ownership wrong")
	}
}

func TestInterruptedArrangementHasNoRememberedSuccess(t *testing.T) {
	w := New(27)
	w.Player.Contacts = 2
	w.Retaliation()
	w.Advance(120)
	w.Event, _ = w.ValidateProposal(approachProposal())
	choice(t, &w, "accept")
	if w.Event.Kind != "warning" || w.Arrangements[0].Status != "interrupted" || w.Arrangements[0].Result != "" {
		t.Fatal("interrupted work remembered as finished")
	}
}

func TestDirectorConnectionRequiresLatestSuccessfulCurrentLifeWork(t *testing.T) {
	w := New(27)
	scene, _ := w.ValidateProposal(approachProposal())
	w.Event = scene
	choice(t, &w, "accept")
	connection := w.Clone().DirectorConnection()
	if connection == nil || connection.ID != scene.ID || connection.Status != "completed" {
		t.Fatal("completed work not available for followup")
	}
	scene, _ = w.ValidateProposal(approachProposal())
	w.Event = scene
	choice(t, &w, "decline")
	if w.DirectorConnection() != nil {
		t.Fatal("refusal forced an older successful thread")
	}
	w.Die("test")
	act(t, &w, "new_life", "")
	if w.DirectorConnection() != nil {
		t.Fatal("previous person's work assigned to new life")
	}
}
