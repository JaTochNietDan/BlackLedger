package core

import (
	"encoding/json"
	"testing"
)

func burglaryCue(t *testing.T, w *World) CueBurglary {
	t.Helper()
	for i := len(w.VisualCues) - 1; i >= 0; i-- {
		if w.VisualCues[i].Burglary != nil {
			return *w.VisualCues[i].Burglary
		}
	}
	if w.LastResult != nil {
		for _, cue := range w.LastResult.Cues {
			if cue.Burglary != nil {
				return *cue.Burglary
			}
		}
	}
	t.Fatal("missing completed burglary cue")
	return CueBurglary{}
}

func TestBurglaryPresentationCapturesActualCashOnceAndSurvivesSave(t *testing.T) {
	w := burglaryWorld()
	before := w.Player.Cash
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "burgle:mara", Target: "mercercourt"})
	if err != nil {
		t.Fatal(err)
	}
	w = next
	cue := burglaryCue(t, w)
	if !cue.Success || cue.Taken != 180 || w.Player.Cash-before != cue.Taken || cue.ResidentPresent || cue.Fatal || cue.HealthLost != 0 {
		t.Fatalf("wrong committed outcome: %+v", cue)
	}
	raw, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var loaded World
	if err = json.Unmarshal(raw, &loaded); err != nil {
		t.Fatal(err)
	}
	if burglaryCue(t, &loaded) != cue {
		t.Fatal("save changed the recorded break-in")
	}
	w.HouseholdSavings["mara"] = HouseholdAccount{Cash: 913, Day: 1}
	if burglaryCue(t, w).Taken != 180 {
		t.Fatal("old cue changed with the current household balance")
	}
}

func TestEmptySuccessfulBurglaryIsNotAFailedAttempt(t *testing.T) {
	w := burglaryWorld()
	w.HouseholdSavings["mara"] = HouseholdAccount{}
	if err := w.Burgle("mara"); err != nil {
		t.Fatal(err)
	}
	cue := burglaryCue(t, w)
	if !cue.Success || cue.Taken != 0 || cue.HealthLost != 0 || cue.Fatal {
		t.Fatalf("empty search incorrectly presented: %+v", cue)
	}
}

func TestFailedBurglaryPresentationRecordsWitnessedResidentAndActualInjury(t *testing.T) {
	w := burglaryWorld()
	n := w.NPC("mara")
	n.Location = n.Home
	for seed := uint32(1); ; seed++ {
		trial := w.Clone()
		trial.RNG = seed
		if trial.Random() >= w.burglaryChance(n) {
			w.RNG = seed
			break
		}
	}
	before := w.Player.Health
	if err := w.Burgle("mara"); err != nil {
		t.Fatal(err)
	}
	cue := burglaryCue(t, w)
	if cue.Success || cue.Taken != 0 || !cue.ResidentPresent || !cue.Identified || cue.HealthLost != before-w.Player.Health || cue.HealthLost == 0 || cue.Fatal {
		t.Fatalf("wrong failed encounter: %+v", cue)
	}
	if w.HouseholdSavings["mara"].Cash != 180 {
		t.Fatal("failed attempt consumed household cash")
	}
	n.Location = "bar"
	if !burglaryCue(t, w).ResidentPresent {
		t.Fatal("recorded presence follows current movements")
	}
}

func TestFatalBurglaryCueDoesNotClaimAnEscape(t *testing.T) {
	w := burglaryWorld()
	n := w.NPC("mara")
	n.Location = n.Home
	w.Player.Health = 25
	found := false
	for seed := uint32(1); seed < 10000; seed++ {
		trial := w.Clone()
		trial.RNG = seed
		if err := trial.Burgle("mara"); err != nil {
			t.Fatal(err)
		}
		cue := burglaryCue(t, trial)
		if !trial.Player.Alive {
			if !cue.Fatal || cue.Success || cue.Taken != 0 || cue.HealthLost != 25 {
				t.Fatalf("wrong fatal outcome: %+v", cue)
			}
			found = true
			break
		}
	}
	if !found {
		t.Fatal("fatal fixture could not be constructed")
	}
}
