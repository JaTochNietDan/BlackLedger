package core

import "testing"

// The city's people walk to the ground their family holds. The player's people
// are members of an organization like any other, so the same rule should reach
// them — and a man told to stand on a door across the city should have to get
// there, rather than being in two places in the same minute.

func ownOrganization(t *testing.T, seed uint32) *World {
	t.Helper()
	w := New(seed)
	w.Player.Cash, w.Player.Respect, w.Player.Contacts = 9000, OrganizationStanding, 3
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["garage"].Owner = "player:1"
	w.OrganizationDay()
	if !w.Incorporated() {
		t.Fatal("the player has no organization to have people in")
	}
	for _, n := range w.Civilians() {
		if IsOfficial(n.ID) {
			continue
		}
		w.Player.Location = n.Location
		if w.SignOn(n.ID) == nil {
			break
		}
	}
	w.Player.Location = "laundry"
	return w
}

func TestYourOwnPeopleMindYourGroundToo(t *testing.T) {
	w := ownOrganization(t, 4)
	mine := []*NPC{}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Faction == w.PlayerOrganizationID() {
			mine = append(mine, n)
		}
	}
	if len(mine) == 0 {
		t.Fatal("nobody signed on")
	}
	// Put them all somewhere that is not one of the player's premises.
	for _, n := range mine {
		n.Location = "bar"
	}
	w.SetOut()
	for _, n := range mine {
		if n.Heading == "laundry" || n.Heading == "garage" {
			return
		}
	}
	t.Fatalf("the player holds two addresses and none of their own people went to either: %v",
		func() []string {
			out := []string{}
			for _, n := range mine {
				out = append(out, n.Name+"→"+n.Heading)
			}
			return out
		}())
}

// Posting a man on a door put him there in the same minute the order was
// given, however far away he was standing. Everybody else in this city has to
// walk; the player's own people were the only ones who could be in two places
// at once, and the door was defended from the moment of the decision rather
// than from the moment somebody was standing in it.
func TestSomebodySentToADoorHasToGetThere(t *testing.T) {
	w := ownOrganization(t, 4)
	free := w.Unposted()
	if len(free) == 0 {
		t.Fatal("nobody to post")
	}
	for _, n := range free {
		n.Location = "bar" // across the city from the laundry
	}
	if err := w.Post("laundry"); err != nil {
		t.Fatal(err)
	}
	man := w.NPC(w.Properties["laundry"].Posted)
	if man == nil {
		t.Fatal("nobody was given the door")
	}
	if man.Location == "laundry" {
		t.Fatal("he was standing at the laundry in the same minute he was told to go there")
	}
	if !w.Travelling(man) || man.Heading != "laundry" {
		t.Fatalf("he was told to go to the laundry and is %q heading %q", man.Location, man.Heading)
	}
	// And the door is not defended by somebody who has not arrived.
	if w.PostingDefenceAt("laundry") != 0 {
		t.Fatalf("a door with nobody standing in it is worth %d", w.PostingDefenceAt("laundry"))
	}
	// He is on the street like anybody else, and the city can say why.
	found := false
	for _, j := range w.OnTheStreet() {
		if j.ID == man.ID {
			found = true
			if j.ToID != "laundry" || j.Because == "" {
				t.Fatalf("the street says %+v", j)
			}
		}
	}
	if !found {
		t.Fatal("he is walking across the city and is not on the street")
	}
	// Once he arrives the door is worth something and stays his.
	w.Minute = man.Arrives
	w.Arrivals()
	if w.PostedAt("laundry") == nil || w.PostingDefenceAt("laundry") <= 0 {
		t.Fatal("he arrived and the door is still worth nothing")
	}
}

// A man already standing in the place is on the door at once: there is nowhere
// for him to walk.
func TestSomebodyAlreadyThereIsOnTheDoorAtOnce(t *testing.T) {
	w := ownOrganization(t, 4)
	free := w.Unposted()
	if len(free) == 0 {
		t.Fatal("nobody to post")
	}
	for _, n := range free {
		n.Location = "laundry"
	}
	if err := w.Post("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.PostingDefenceAt("laundry") <= 0 {
		t.Fatal("somebody standing in the room was sent on a journey to it")
	}
}

// And the address does not claim to be held by somebody who is not there yet.
func TestADoorDoesNotClaimAManWhoIsStillWalking(t *testing.T) {
	w := ownOrganization(t, 4)
	for _, n := range w.Unposted() {
		n.Location = "bar"
	}
	if err := w.Post("laundry"); err != nil {
		t.Fatal(err)
	}
	note := w.PlaceNote("laundry")
	if note == "" || !contains(note, "on the way") {
		t.Fatalf("the laundry says %q while the man is still crossing the city", note)
	}
	if contains(note, "is on the door") {
		t.Fatalf("the laundry claims a door it does not have: %q", note)
	}
	man := w.NPC(w.Properties["laundry"].Posted)
	w.Minute = man.Arrives
	w.Arrivals()
	if !contains(w.PlaceNote("laundry"), "is on the door") {
		t.Fatalf("he arrived and the laundry says %q", w.PlaceNote("laundry"))
	}
}
