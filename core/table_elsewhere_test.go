package core

import "testing"

// Reading the casino and the club side by side: both refused to deal, both for
// the same reason — "There is a hand on the table already" — and only one of
// them offered anything to do about it. The hand was at the club, so the club
// showed "Take another card" and "Stand on 12" beside the refusal. The casino
// showed the refusal alone. A player standing at the casino was told a hand was
// open, given no way to finish it, and not told where it was, while the world
// had been storing the room on the hand the whole time.

func TestARefusedTableSaysWhereTheHandIs(t *testing.T) {
	w := New(5)
	w.Player.Cash = 400
	w.District = 2 // the casino is not in the first district
	w.Player.Location = "club"
	if _, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "play", Target: "club", Amount: 50}); err != nil {
		t.Fatalf("could not sit down at the club: %v", err)
	}
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "play", Target: "club", Amount: 50})
	if err == nil {
		w = next
	}
	if w.Hand == nil || w.Hand.Done {
		t.Fatal("no hand is in play")
	}
	club, _ := PlaceByID(w.Hand.Place)

	// Where the hand is, the player can finish it, and the refusal can assume
	// they can see how.
	here := map[string]Action{}
	w.Player.Location = w.Hand.Place
	for _, a := range w.Actions(w.Hand.Place) {
		here[a.ID] = a
	}
	if _, ok := here["hit"]; !ok {
		t.Fatal("the room holding the hand does not offer to play it")
	}

	// Anywhere else, the refusal is the only thing the player gets, so it has
	// to carry the room.
	elsewhere := ""
	for i := range Locations {
		if HasTables(Locations[i].ID) && Locations[i].ID != w.Hand.Place {
			elsewhere = Locations[i].ID
		}
	}
	if elsewhere == "" {
		t.Fatal("the city has only one place with tables")
	}
	// Only the room the player is standing in carries its full list, which is
	// exactly the situation: they have walked to the other house to play.
	w.Player.Location = elsewhere
	for _, a := range w.Actions(elsewhere) {
		if a.ID != "play" {
			continue
		}
		if !a.Disabled {
			t.Fatal("a second hand could be dealt elsewhere")
		}
		if !contains(a.Reason, club.Name) {
			t.Fatalf("the refusal does not say where the hand is: %q", a.Reason)
		}
		return
	}
	t.Fatalf("%s does not offer the small tables at all", elsewhere)
}

// And where the hand is, the refusal does not send the player to the room they
// are already standing in.
func TestTheRefusalWhereTheHandIsDoesNotSendYouAnywhere(t *testing.T) {
	w := New(5)
	w.Player.Cash = 400
	w.District = 2
	w.Player.Location = "club"
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "play", Target: "club", Amount: 50})
	if err != nil {
		t.Fatalf("could not sit down: %v", err)
	}
	w = next
	if w.Hand == nil || w.Hand.Done {
		t.Skip("the hand settled immediately")
	}
	w.Player.Location = w.Hand.Place
	here, _ := PlaceByID(w.Hand.Place)
	for _, a := range w.Actions(w.Hand.Place) {
		if a.ID != "play" || !a.Disabled {
			continue
		}
		if contains(a.Reason, here.Name) {
			t.Fatalf("the refusal sends the player to the room they are in: %q", a.Reason)
		}
		return
	}
}
