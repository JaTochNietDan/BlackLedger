package core

import "testing"

func distiller(t *testing.T) *World {
	t.Helper()
	w := New(801)
	w.Properties["laundry"].Owner = "player:1"
	w.Player.Location = "laundry"
	w.Player.Cash = 5000
	return w
}

func TestAStillNeedsSomewhereToHideItAndSomebodyToWorkIt(t *testing.T) {
	t.Parallel()
	w := distiller(t)
	if reason := w.StillReadiness("laundry"); reason != "" {
		t.Fatal("a staffed laundry could not hide one:", reason)
	}
	// A casino floor is not a place to hide a still, nor is somebody else's.
	w.Properties["casino"].Owner = "player:1"
	if w.StillReadiness("casino") == "" {
		t.Fatal("a still was set up on a casino floor")
	}
	if w.StillReadiness("club") == "" {
		t.Fatal("a still was set up in a rival's premises")
	}
	// Nobody working means nobody to run it.
	w.Properties["laundry"].Staff = 0
	if w.StillReadiness("laundry") == "" {
		t.Fatal("a still was set up with nobody to work it")
	}
	if err := w.BuildStill("laundry"); err == nil {
		t.Fatal("the command ignored a requirement the action reports")
	}
}

func TestAStillMakesStockAndTheStockIsTheRisk(t *testing.T) {
	t.Parallel()
	w := distiller(t)
	if err := w.BuildStill("laundry"); err != nil {
		t.Fatal(err)
	}
	heat, supply := w.Player.Heat, w.Properties["laundry"].Supply
	w.StillDay()
	if w.Holding("moonshine") == 0 {
		t.Fatal("a working still produced nothing")
	}
	if w.Player.Heat <= heat {
		t.Fatal("a still drew no attention of its own")
	}
	if w.Properties["laundry"].Supply >= supply {
		t.Fatal("a still consumed nothing to run")
	}
	// The stock it makes draws attention every day it is held, through the
	// same rule that governs anything else being carried.
	held := w.Player.Heat
	w.ContrabandDay()
	if w.Player.Heat <= held {
		t.Fatal("stock from a still draws no attention")
	}
}

func TestAStillStopsWhenThereIsNowhereToPutTheOutput(t *testing.T) {
	t.Parallel()
	w := distiller(t)
	if err := w.BuildStill("laundry"); err != nil {
		t.Fatal(err)
	}
	w.Player.Stock = map[string]int{"moonshine": StillHoard}
	w.Properties["laundry"].Supply = 40
	w.StillDay()
	if w.Holding("moonshine") != StillHoard {
		t.Fatal("a still kept producing with nowhere to put it")
	}
	// And it makes nothing without supplies or people.
	w.Player.Stock["moonshine"] = 0
	w.Properties["laundry"].Supply = 0
	w.StillDay()
	if w.Holding("moonshine") != 0 {
		t.Fatal("a still ran with nothing to run on")
	}
}

func TestASearchThatFindsAStillCostsFarMore(t *testing.T) {
	t.Parallel()
	plain, distilling := distiller(t), distiller(t)
	for _, w := range []*World{plain, distilling} {
		w.Player.Heat = 60
		w.Player.Cash = 20000
	}
	if err := distilling.BuildStill("laundry"); err != nil {
		t.Fatal(err)
	}
	distilling.Player.Cash = 20000
	plainCash, stillCash := plain.Player.Cash, distilling.Player.Cash
	plain.Raid()
	distilling.Raid()
	if plainCash-plain.Player.Cash >= stillCash-distilling.Player.Cash {
		t.Fatal("finding a still cost no more than a search that found nothing")
	}
	if distilling.Properties["laundry"].Still {
		t.Fatal("the still survived the raid")
	}
	found := false
	for _, s := range distilling.Edition() {
		if s["kind"] == "police" && s["headline"] != "" {
			found = true
		}
	}
	if !found {
		t.Fatal("a seized still was never reported")
	}
}

func TestTakingItOutIsAWayToStopBeingWorthWatching(t *testing.T) {
	t.Parallel()
	w := distiller(t)
	if err := w.BuildStill("laundry"); err != nil {
		t.Fatal(err)
	}
	w.Player.Heat = 30
	if err := w.Dismantle("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.Properties["laundry"].Still {
		t.Fatal("the still is still there")
	}
	if w.Player.Heat >= 30 {
		t.Fatal("taking it out did not reduce what there is to find")
	}
	if w.DismantleReadiness("laundry") == "" {
		t.Fatal("a second dismantling was offered")
	}
}

func TestThievesPreferCratesToATill(t *testing.T) {
	t.Parallel()
	stolen := 0
	for i := uint32(1); i <= 300; i++ {
		w := New(i * 2654435761)
		w.Properties["laundry"].Owner = "player:1"
		w.Properties["laundry"].Still = true
		w.Player.Stock = map[string]int{"moonshine": 40}
		w.Player.Cash = 5000
		before := w.Holding("moonshine")
		for day := 0; day < 30; day++ {
			w.Minute += 1440
			w.PeopleDay()
		}
		if w.Holding("moonshine") < before {
			stolen++
		}
	}
	t.Logf("of 300 cities over 30 days, crates were taken from a still in %d", stolen)
	if stolen == 0 {
		t.Fatal("nobody ever steals from a still")
	}
	if stolen == 300 {
		t.Fatal("running a still means being robbed in every campaign")
	}
}
