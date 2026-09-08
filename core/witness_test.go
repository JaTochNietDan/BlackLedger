package core

import "testing"

// The four loudest things that happen in this game produced nothing to look at.
// These hold the fix: the moments that matter say what they looked like and
// where, and nothing claims to have happened somewhere that does not exist.

func TestTheLoudestMomentsAreWorthLookingAt(t *testing.T) {
	for _, kind := range []string{"killing", "explosion", "gunfight", "raid", "seizure", "arrest", "attack"} {
		if Gravity(kind) == 0 {
			t.Errorf("%q is worth nothing to stop for", kind)
		}
	}
	if Gravity("killing") <= Gravity("attack") || Gravity("explosion") <= Gravity("seizure") {
		t.Fatal("the ordering does not put the loudest thing first")
	}
}

func TestAKillingIsSomethingYouSee(t *testing.T) {
	w := proprietor(t)
	victim := w.NPC("mara")
	if victim == nil {
		t.Fatal("nobody to lose")
	}
	victim.Location = "bar"
	w.VisualCues = nil
	if !w.Kill(victim.ID, "Shot twice at the counter.") {
		t.Fatal("nobody died")
	}
	cue, ok := w.Worth()
	if !ok {
		t.Fatal("a killing produced nothing to look at")
	}
	if cue.Kind != "killing" || cue.Target != "bar" {
		t.Fatalf("the city says it was a %q at %q", cue.Kind, cue.Target)
	}
	if cue.Headline == "" {
		t.Fatal("the scene has no headline to follow it")
	}
	// The scene and the paper must not be able to disagree.
	found := false
	for _, s := range w.News {
		if s.Headline == cue.Headline {
			found = true
		}
	}
	if !found {
		t.Fatalf("the scene promises %q and the paper never carried it", cue.Headline)
	}
	if len(cue.Actors) == 0 || cue.Actors[0] != victim.Name {
		t.Fatalf("the scene is about %v", cue.Actors)
	}
}

func TestARaidAndASeizureAreSomethingYouSee(t *testing.T) {
	w := chargeable(t)
	w.Player.Retainers = nil
	w.VisualCues = nil
	w.Player.Heat = 60
	w.search()
	cue, ok := w.Worth()
	if !ok {
		t.Fatal("the police came to the door and there was nothing to look at")
	}
	if cue.Kind != "raid" && cue.Kind != "seizure" {
		t.Fatalf("the visit reads as a %q", cue.Kind)
	}
	if _, exists := PlaceByID(cue.Target); !exists {
		t.Fatalf("it happened at %q, which is not a place", cue.Target)
	}
	if cue.Caption == "" {
		t.Fatal("nobody can say what it looked like")
	}
}

func TestGroundChangingHandsIsSomethingYouSee(t *testing.T) {
	// Over enough raids one of them takes the ground, and when it does the
	// player is taken there rather than told about it in a list.
	seen := 0
	for seed := uint32(1); seed <= 300 && seen == 0; seed++ {
		w, _ := testator(t)
		w.WorldRNG, w.RNG = seed*2654435761, seed*2654435761
		attacker := &w.Factions[0]
		attacker.Power = 100
		w.Properties["laundry"].Condition = 12
		w.VisualCues = nil
		w.contestAt(attacker, w.PlayerOrganization(), "laundry")
		for _, cue := range w.VisualCues {
			if cue.Kind == "seizure" && cue.Target == "laundry" {
				seen++
				if cue.Headline == "" {
					t.Fatal("ground changed hands with no headline behind it")
				}
			}
		}
	}
	if seen == 0 {
		t.Fatal("300 raids on a laundry at 12% condition and none of them took it")
	}
}

func TestNothingHappensSomewhereThatDoesNotExist(t *testing.T) {
	w := proprietor(t)
	w.VisualCues = nil
	w.Witness("killing", "a-place-with-no-address", "Something happened.", "SOMETHING HAPPENED")
	if len(w.VisualCues) != 0 {
		t.Fatal("the interface was sent somewhere that is not on the map")
	}
	// And the worst of several is the one worth showing.
	w.Witness("robbery", "bar", "A till was emptied.", "")
	w.Witness("killing", "club", "Somebody was shot.", "SOMEBODY WAS SHOT")
	w.Witness("attack", "laundry", "A window went in.", "")
	cue, ok := w.Worth()
	if !ok || cue.Kind != "killing" {
		t.Fatalf("out of three things the city picked %q", cue.Kind)
	}
}
