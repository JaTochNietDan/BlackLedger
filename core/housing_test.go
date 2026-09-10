package core

import "testing"

func TestOwnedResidenceSurvivesMovingButNotDeath(t *testing.T) {
	t.Parallel()
	w := New(27)
	w.District = 2
	w.Player.Cash = 8000
	w.Player.Location = "estate"
	act(t, &w, "move_home", "estate")
	if !w.Own("estate") || w.Player.Home != "estate" || w.Player.Respect != 3 {
		t.Fatal("purchase did not grant ownership and housing progress")
	}
	act(t, &w, "travel", "room")
	act(t, &w, "move_home", "room")
	act(t, &w, "travel", "estate")
	cash := w.Player.Cash
	act(t, &w, "move_home", "estate")
	if w.Player.Cash != cash || w.Player.Respect != 3 {
		t.Fatal("return charged another purchase or farmed respect")
	}
	w.Die("Housing persistence test")
	if w.Properties["estate"].Owner != "former:Alex Varga" || w.Properties["estate"].Income != 0 {
		t.Fatal("estate failed to pass to former organization without inventing income")
	}
	act(t, &w, "new_life", "")
	w.Player.Cash = 8000
	w.Player.Location = "estate"
	if w.Own("estate") || w.Player.BestHome != 0 {
		t.Fatal("new stranger inherited the estate or housing rank")
	}
	if _, err := Execute(w, Command{Kind: "move_home", Target: "estate", Revision: w.Revision}); err == nil {
		t.Fatal("new stranger bought an occupied estate")
	}
}

func TestRentingBackAndForthCannotFarmRespect(t *testing.T) {
	t.Parallel()
	w := New(27)
	w.District = 1
	w.Player.Cash = 1000
	w.Player.Location = "apartment"
	act(t, &w, "move_home", "apartment")
	for i := 0; i < 2; i++ {
		act(t, &w, "travel", "room")
		act(t, &w, "move_home", "room")
		act(t, &w, "travel", "apartment")
		act(t, &w, "move_home", "apartment")
	}
	if w.Player.Respect != 3 {
		t.Fatalf("repeated moves generated respect: %d", w.Player.Respect)
	}
}
