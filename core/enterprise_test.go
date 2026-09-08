package core

import "testing"

func ownedBusiness(t *testing.T) *World {
	t.Helper()
	w := New(3)
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["laundry"].Condition = 100
	return w
}

func TestHowABusinessIsRunChangesWhatItEarns(t *testing.T) {
	earned := map[string]int{}
	for _, mode := range []string{"clean", "standard", "hard"} {
		w := ownedBusiness(t)
		if mode != "standard" {
			if err := w.SetMode("laundry", mode); err != nil {
				t.Fatal(err)
			}
		}
		before := w.Player.Cash
		w.Advance(1440)
		earned[mode] = w.Player.Cash - before
	}
	if !(earned["hard"] > earned["standard"] && earned["standard"] > earned["clean"]) {
		t.Fatalf("running a business harder did not pay more: %v", earned)
	}
}

func TestSkimmingCostsAttentionAndCondition(t *testing.T) {
	w := ownedBusiness(t)
	if err := w.SetMode("laundry", "hard"); err != nil {
		t.Fatal(err)
	}
	heat, condition := w.Player.Heat, w.Properties["laundry"].Condition
	w.BusinessDay()
	if w.Player.Heat <= heat {
		t.Fatal("skimming drew no police attention")
	}
	if w.Properties["laundry"].Condition >= condition {
		t.Fatal("skimming did not wear the premises")
	}

	clean := ownedBusiness(t)
	if err := clean.SetMode("laundry", "clean"); err != nil {
		t.Fatal(err)
	}
	cleanHeat, cleanCondition := clean.Player.Heat, clean.Properties["laundry"].Condition
	clean.BusinessDay()
	if clean.Player.Heat != cleanHeat {
		t.Fatal("running clean still drew attention")
	}
	if clean.Properties["laundry"].Condition != cleanCondition {
		t.Fatal("running clean still wore the premises")
	}
}

func TestSkimmingBringsAFamilyDemandSooner(t *testing.T) {
	hard, quiet := ownedBusiness(t), ownedBusiness(t)
	if err := hard.SetMode("laundry", "hard"); err != nil {
		t.Fatal(err)
	}
	if err := quiet.SetMode("laundry", "clean"); err != nil {
		t.Fatal(err)
	}
	hard.BusinessPressure()
	quiet.BusinessPressure()
	if hard.NextPressure >= quiet.NextPressure {
		t.Fatalf("skimming did not shorten the wait for a demand: %d vs %d", hard.NextPressure, quiet.NextPressure)
	}
}

func TestABusinessIsOnlyRunByWhoeverOwnsIt(t *testing.T) {
	w := New(4)
	if err := w.SetMode("club", "hard"); err == nil {
		t.Fatal("the player set the operating mode of a family's premises")
	}
	if err := w.SetMode("laundry", "hard"); err == nil {
		t.Fatal("the player ran a business they do not own")
	}
	owned := ownedBusiness(t)
	if err := owned.SetMode("laundry", "nonsense"); err == nil {
		t.Fatal("an unknown way of running a business was accepted")
	}
	if err := owned.SetMode("laundry", "standard"); err == nil {
		t.Fatal("setting the mode it already has was accepted")
	}
	// Premises that do not trade have nothing to decide about.
	owned.Properties["estate"].Owner = "player:1"
	if err := owned.SetMode("estate", "hard"); err == nil {
		t.Fatal("a residence was given an operating mode")
	}
}

func TestTheChoiceIsOfferedWhereItApplies(t *testing.T) {
	w := ownedBusiness(t)
	w.Player.Location = "laundry"
	offered := map[string]Action{}
	for _, a := range w.Actions("laundry") {
		offered[a.ID] = a
	}
	for _, m := range operatingModes {
		a, ok := offered["operate:"+m.ID]
		if !ok {
			t.Fatalf("%s was not offered at a business the player runs", m.ID)
		}
		if m.ID == "standard" && !a.Disabled {
			t.Fatal("the mode already in use was offered as a change")
		}
		if a.Minutes != 0 {
			t.Fatal("deciding how to run a business should not consume time")
		}
	}
	// Not offered on premises the player does not own.
	w.Player.Location = "club"
	for _, a := range w.Actions("club") {
		if len(a.ID) > 8 && a.ID[:8] == "operate:" {
			t.Fatal("an operating decision was offered for a family's premises")
		}
	}
}

func TestOlderSavesKeepEarningWhatTheyEarned(t *testing.T) {
	// A property with no recorded mode runs the ordinary way.
	w := ownedBusiness(t)
	w.Properties["laundry"].Mode = ""
	if w.Mode("laundry").ID != "standard" {
		t.Fatal("a property with no recorded mode did not default to the ordinary way")
	}
	if operatingMode("").Take != 1 {
		t.Fatal("the default way of running a business changed what it earns")
	}
}
