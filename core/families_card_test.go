package core

import "testing"

// A second reading of the Families screen. Each card lays out four figures —
// strength, money, people, ground — and two of the five cards had only three.
// The People cell simply was not there.
//
// Strength and money degrade when the player has nobody inside: they read "not
// much" or "nobody will say" instead of a number, which is the right shape.
// People did not degrade, it disappeared, so the same ignorance was expressed
// two different ways in one card and the reader saw a table whose columns
// change from row to row.
//
// Worse, it cannot tell two things apart. A family the player knows well that
// has nobody left reads exactly like a family they know nothing about.

func TestEveryFamilyCardSaysHowManyPeopleItHas(t *testing.T) {
	t.Parallel()
	w := New(43)
	// Somebody inside one family and nobody inside the other. Knowledge is
	// earned rather than stored: a contact network hears the ordinary things,
	// somebody inside who trusts you is worth more, and asking around directly
	// is worth a third.
	w.Player.Contacts = 5
	for _, n := range w.Members("bellandi") {
		n.Trust = 40
	}
	w.Player.Enquiries = map[string]int{"bellandi": w.Minute + 1440}
	seen := map[int]string{}
	for _, f := range w.PublicFactions() {
		if f.Hands == "" {
			t.Errorf("%s has no people figure at all: knowledge %d", f.Name, f.Knowledge)
		}
		seen[f.Knowledge] = f.Hands
	}

	// And with nobody to ask at all, it says so rather than going quiet — the
	// same way strength and money do.
	cold := New(43)
	for _, f := range cold.PublicFactions() {
		if f.Knowledge != 0 {
			t.Fatalf("%s is known without any contacts: %d", f.Name, f.Knowledge)
		}
		if f.Hands != "nobody will say" {
			t.Errorf("with nobody to ask, %s reports %q", f.Name, f.Hands)
		}
		if f.Strength != "nobody will say" {
			t.Errorf("strength has stopped behaving this way: %q", f.Strength)
		}
	}

	// A crowd is described in its own words. "as much as anybody" belongs to
	// money and strength, and reads strangely about eight people.
	for _, c := range []struct {
		n    int
		want string
	}{{0, "nobody worth naming"}, {1, "a handful"}, {3, "a handful"},
		{4, "a fair few"}, {7, "a fair few"}, {8, "a great many"}} {
		if got := handful(c.n); got != c.want {
			t.Errorf("handful(%d) = %q, want %q", c.n, got, c.want)
		}
		if got := handful(c.n); got == roughly(c.n*10) {
			t.Errorf("handful(%d) borrows the money scale: %q", c.n, got)
		}
	}
	t.Logf("by what the player knows: %v", seen)
}

func TestAFamilyWithNobodyLeftSaysSoRatherThanGoingBlank(t *testing.T) {
	t.Parallel()
	w := New(43)
	// Known inside out, and then emptied. Bellandi stays a stranger.
	w.Player.Contacts = 5
	for _, n := range w.Members("russo") {
		n.Trust = 40
	}
	w.Player.Enquiries = map[string]int{"russo": w.Minute + 1440}
	// Everybody in that family is gone.
	for i := range w.NPCs {
		if w.NPCs[i].Faction == "russo" {
			w.Kill(w.NPCs[i].ID, "A bad week.")
		}
	}
	var russo, bellandi PublicFaction
	for _, f := range w.PublicFactions() {
		switch f.ID {
		case "russo":
			russo = f
		case "bellandi":
			bellandi = f
		}
	}
	if len(w.Members("russo")) != 0 {
		t.Fatalf("somebody is still in the family: %d", len(w.Members("russo")))
	}
	if russo.Hands == "" {
		t.Fatal("a family with nobody left says nothing about its people")
	}
	if russo.Hands == bellandi.Hands {
		t.Fatalf("a family with nobody left reads the same as one nobody will talk about: both %q", russo.Hands)
	}
	t.Logf("known and empty: %q; unknown: %q", russo.Hands, bellandi.Hands)
}
