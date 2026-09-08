package core

import "testing"

func cityWithJobs(t *testing.T) *World {
	t.Helper()
	w := New(179)
	w.MigrateLivingWorld()
	return w
}

func TestSomebodyIsDoingEveryJobOnTheFirstMorning(t *testing.T) {
	w := cityWithJobs(t)
	for _, r := range roles {
		holder := w.Holder(r.ID)
		if holder == nil {
			t.Fatalf("nobody was doing the %s's job", r.Title)
		}
		if w.HolderID(r.ID) != holder.ID {
			t.Fatal("the id and the person disagreed")
		}
	}
	if w.Holder("fixer").ID != "mara" {
		t.Fatalf("the campaign did not start with Mara as the fixer: %s", w.Holder("fixer").Name)
	}
}

func TestNobodyIsFixed(t *testing.T) {
	w := cityWithJobs(t)
	for _, r := range roles {
		before := w.Holder(r.ID)
		if !w.Kill(before.ID, "Shot.") {
			t.Fatalf("%s could not be killed", before.Name)
		}
		if w.Holder(r.ID) != nil {
			t.Fatalf("%s was still doing the job after dying", before.Name)
		}
		w.FillRoles()
		after := w.Holder(r.ID)
		if after == nil {
			t.Fatalf("nobody replaced the %s", r.Title)
		}
		if after.ID == before.ID {
			t.Fatal("the dead took their own job back")
		}
		if after.Location != r.Where || after.Role != r.Title {
			t.Fatalf("the new %s was %s at %s", r.Title, after.Role, after.Location)
		}
		if after.Trust != 0 {
			t.Fatalf("a stranger arrived already trusting the player: %d", after.Trust)
		}
		if !w.hasRecord("Somebody else is doing that job now") {
			t.Fatal("nobody was told")
		}
	}
}

func TestNobodyDoesTwoJobsAtOnce(t *testing.T) {
	w := cityWithJobs(t)
	// Empty every job at once and see who fills them.
	for _, r := range roles {
		w.Kill(w.Holder(r.ID).ID, "Shot.")
	}
	w.FillRoles()
	seen := map[string]bool{}
	for _, r := range roles {
		holder := w.Holder(r.ID)
		if holder == nil {
			t.Fatalf("the %s's job went unfilled", r.Title)
		}
		if seen[holder.ID] {
			t.Fatalf("%s was doing two jobs", holder.Name)
		}
		seen[holder.ID] = true
	}
}

func TestFillingIsIdempotent(t *testing.T) {
	w := cityWithJobs(t)
	before := len(w.NPCs)
	holders := map[string]string{}
	for _, r := range roles {
		holders[r.ID] = w.Holder(r.ID).ID
	}
	for i := 0; i < 20; i++ {
		w.FillRoles()
	}
	if len(w.NPCs) != before {
		t.Fatalf("the city grew from %d to %d by filling jobs that were filled", before, len(w.NPCs))
	}
	for _, r := range roles {
		if w.Holder(r.ID).ID != holders[r.ID] {
			t.Fatalf("the %s was replaced while still alive", r.Title)
		}
	}
}

func TestScenesSpeakThroughWhoeverHasTheJob(t *testing.T) {
	w := cityWithJobs(t)
	w.Kill("mara", "Shot.")
	w.FillRoles()
	fixer := w.Holder("fixer")

	// A warning at home is the fixer's line.
	w.Player.Location, w.Player.Contacts = w.Player.Home, 3
	w.RetaliationFrom("bellandi")
	for i := 0; i < 40 && w.Event == nil; i++ {
		w.Advance(120)
	}
	if w.Event == nil {
		t.Skip("no scene was raised in this city")
	}
	if w.Event.Speaker == "mara" {
		t.Fatal("a dead fixer was still speaking")
	}
	if w.Event.Speaker != "" && w.Event.Speaker != fixer.ID && w.NPC(w.Event.Speaker) == nil {
		t.Fatalf("the scene spoke through %q, who does not exist", w.Event.Speaker)
	}
}

func TestNobodyIsDoingAJobNobodyCanDo(t *testing.T) {
	w := cityWithJobs(t)
	// Kill everybody who is not in an organization, so the only candidates are
	// people the city would have to invent.
	for _, n := range w.Civilians() {
		n.Dead = true
	}
	for _, r := range roles {
		if h := w.Holder(r.ID); h != nil {
			w.Kill(h.ID, "Shot.")
		}
	}
	w.FillRoles()
	for _, r := range roles {
		if w.Holder(r.ID) == nil {
			t.Fatalf("the city could not find anybody to be the %s", r.Title)
		}
	}
}
