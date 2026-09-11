package core

import "testing"

func TestPausedArrangementKeepsTermsAndOnlyRemainingTime(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		method           string
		duration, reward int
	}{{"accept", 45, 75}, {"approach:careful", 75, 60}, {"approach:press", 30, 95}} {
		t.Run(tc.method, func(t *testing.T) {
			w := pressureWorld()
			start := w.Minute
			e, err := w.ValidateProposal(Proposal{Location: "bar", Title: "A sealed message", Body: "Deliver this message at Saint Agnes.", Speaker: "mara", Operation: "courier", Outcome: "Done.", Approaches: []Approach{{Method: "careful", Label: "Check the route"}, {Method: "press", Label: "Take the shortcut"}}})
			if err != nil {
				t.Fatal(err)
			}
			w.Event = e
			choice(t, &w, tc.method)
			if w.SuspendedJob == nil || w.SuspendedJob.Remaining != tc.duration-30 || w.Minute != start+30 || w.Event.Kind != "business_pressure" {
				t.Fatal("job progress not saved")
			}
			if w.Arrangements[0].Status != "paused" {
				t.Fatal("paused job considered complete")
			}
			w = w.Clone()
			paid := w.TheirShare(w.Event.Target)
			choice(t, &w, "pay")
			if w.Event == nil || w.Event.Kind != "resume_job" || w.Minute != start+30 {
				t.Fatal("missing resume decision")
			}
			id := w.Event.ID
			choice(t, &w, "resume")
			if w.Minute != start+tc.duration || w.SuspendedJob != nil || w.Event != nil || w.Arrangements[0].Status != "completed" {
				t.Fatal("incorrect resumed completion")
			}
			// The demand is a week of what the place takes rather than a flat
			// sixty, so the arithmetic asks the world what was paid instead of
			// restating a number that has stopped being one.
			want := 200 - paid + tc.reward + 14*tc.duration/60
			if w.Player.Cash != want {
				t.Fatalf("cash %d, want %d", w.Player.Cash, want)
			}
			if _, err := Execute(w, Command{Revision: w.Revision, Kind: "choice", Event: id, Choice: "resume"}); err == nil {
				t.Fatal("resumed job replayed")
			}
		})
	}
}

func TestPausedJobAbandonmentAndPoliceCompletion(t *testing.T) {
	t.Parallel()
	for _, abandon := range []bool{false, true} {
		w := pressureWorld()
		w.Player.Heat = 14
		e, _ := w.ValidateProposal(Proposal{Title: "A delivery", Body: "Carry a message.", Speaker: "mara", Operation: "courier", Outcome: "Done."})
		w.Event = e
		choice(t, &w, "accept")
		choice(t, &w, "pay")
		if abandon {
			before := w.Player.Cash
			choice(t, &w, "abandon")
			if w.Player.Cash != before || w.SuspendedJob != nil || w.Arrangements[0].Status != "abandoned" {
				t.Fatal("abandonment paid reward")
			}
		} else {
			choice(t, &w, "resume")
			if w.Event == nil || w.Event.Kind != "police_stop" || w.Arrangements[0].Status != "awaiting_police" {
				t.Fatal("resume bypassed police")
			}
			before := w.Player.Cash
			choice(t, &w, "pay")
			if w.Player.Cash != before-40+75 || w.Arrangements[0].Status != "completed" {
				t.Fatal("police lost reward")
			}
		}
	}
}

func TestUrgentWarningDoesNotSuspendJobAsRoutineConversation(t *testing.T) {
	t.Parallel()
	w := New(27)
	w.Player.Contacts = 2
	w.Plots = []Plot{{ID: "danger", Kind: "hit", Actor: "bellandi", Life: 1, Due: w.Minute + 120}}
	e, _ := w.ValidateProposal(Proposal{Title: "A delivery", Body: "Carry a message.", Speaker: "mara", Operation: "courier", Outcome: "Done."})
	w.Event = e
	choice(t, &w, "accept")
	if w.Event == nil || w.Event.Kind != "warning" || w.SuspendedJob != nil || w.Arrangements[0].Status != "interrupted" {
		t.Fatal("warning treated as routine")
	}
}

func TestRepeatedDemandsPreserveRemainingWork(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	e, _ := w.ValidateProposal(Proposal{Title: "A collection", Body: "Collect the agreed payment.", Speaker: "mara", Operation: "collection", Outcome: "Done."})
	w.Event = e
	choice(t, &w, "accept")
	choice(t, &w, "pay")
	w.NextPressure = w.Minute + 10
	choice(t, &w, "resume")
	if w.SuspendedJob == nil || w.SuspendedJob.Remaining != 50 {
		t.Fatal("second pause lost elapsed work")
	}
	choice(t, &w, "pay")
	choice(t, &w, "resume")
	if w.Minute != 570 || w.SuspendedJob != nil || w.Arrangements[0].Status != "completed" {
		t.Fatal("nested pause duplicated duration")
	}
}
