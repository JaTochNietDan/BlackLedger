package core

import "testing"

// "It's weird that you need to go to where Russo is to commission a hit on her.
// It feels like a lot of the game is based around where you are for ALL actions
// which is not the right way."
//
// Putting a price on a name is done through a contact, and the list of names it
// can reach was always the whole city rather than whoever happened to be in the
// room. The only thing tying it to the exchange was the panel it was printed in.

func standing(t *testing.T, at string) *World {
	t.Helper()
	w := New(19)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 9000, 80, 100
	w.Player.Contacts = 2
	w.Player.Location = at
	return w
}

func offered(w *World, kind string) (Action, bool) {
	for _, a := range w.Actions(w.Player.Location) {
		if a.ID == kind {
			return a, true
		}
	}
	return Action{}, false
}

func TestTheWorkThatIsNotAboutTheRoomIsOfferedWhereverYouStand(t *testing.T) {
	rooms := []string{"bar", "laundry", "docks", "room", "garage", "casino", "market"}
	for _, kind := range []string{"contract", "investigate", "enquire:bellandi", "pact:russo"} {
		for _, room := range rooms {
			w := standing(t, room)
			a, ok := offered(w, kind)
			if !ok {
				t.Errorf("%s cannot be reached from %s", kind, room)
				continue
			}
			if !a.Anywhere {
				t.Errorf("%s at %s is still filed as work belonging to the room", kind, room)
			}
		}
	}
}

// And it is not only offered but pressable: the command layer looks the action
// up at the place the player is standing, so an action that is not there is an
// action that cannot be pressed however it is drawn.
func TestAPriceCanBePutOnANameFromAnywhere(t *testing.T) {
	w := standing(t, "docks")
	if len(w.ContractTargets()) == 0 {
		t.Fatal("nobody in this city can be named")
	}
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "contract"})
	if err != nil {
		t.Fatalf("standing at the docks, a name cannot be raised: %v", err)
	}
	if next.Event == nil {
		t.Error("nothing came of it")
	}
}

// The room's own panel must not carry them. A fishmonger offering an
// understanding with the Russo family is the fault this is fixing.
func TestTheRoomsOwnPanelLeavesThatWorkOut(t *testing.T) {
	w := standing(t, "market")
	for _, a := range w.Actions("market") {
		if a.Anywhere && a.Group == GroupOf("repair") {
			t.Errorf("%s is filed as premises work", a.ID)
		}
	}
	if _, ok := offered(w, "contract"); !ok {
		t.Error("the exchange lost the button it used to be the only home of")
	}
}

// Sending your own crew against a family's premises, and pointing one family at
// another, were offered only while standing in the building it was aimed at.
// Neither rule ever asked where the player was: the readiness reads respect, a
// crew, their loyalty and what you know, and nothing else. Walking across the
// city to give an order you are not going to carry out yourself is not a
// decision, it is a journey.
func TestOrdersAimedAtAPlaceCanBeGivenFromAnywhere(t *testing.T) {
	w := standing(t, "bar")
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
	w.Player.Contacts = 3
	target := ""
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && l.ID != w.Player.Location && l.District <= w.District {
			if f := w.faction(prop.Owner); f != nil && f.ID != w.PlayerOrganizationID() {
				target = l.ID
				break
			}
		}
	}
	if target == "" {
		t.Fatal("no rival holds anything reachable")
	}
	want := map[string]bool{"sabotage:crew": false, "incite": false}
	for _, a := range w.Actions(target) {
		if _, ours := want[a.ID]; !ours {
			continue
		}
		want[a.ID] = true
		if a.Disabled {
			t.Errorf("%s at %s is refused from across the city: %s", a.ID, target, a.Reason)
		}
		if a.Target != target {
			t.Errorf("%s is aimed at %q rather than at %s", a.ID, a.Target, target)
		}
	}
	for id, found := range want {
		if !found {
			t.Errorf("%s cannot be given about %s without walking there", id, target)
		}
	}
	// And the order actually goes through from where the player is standing.
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "sabotage:crew", Target: target})
	if err != nil {
		t.Fatalf("the order could not be given from the bar: %v", err)
	}
	if next.Player.Location != "bar" {
		t.Errorf("giving an order moved the player to %s", next.Player.Location)
	}
	// It resolves there and then rather than leaving a task behind, the same
	// way going in yourself does, so what proves it happened is the time it
	// took and the account of it.
	if next.Minute != w.Minute+90 {
		t.Errorf("the order took %d minutes", next.Minute-w.Minute)
	}
	if next.LastResult == nil || len(next.LastResult.Records) == 0 {
		t.Error("the order was given and the city has no account of it")
	}
}

// Going in yourself is the one that needs you there. It is the difference
// between the two buttons, and it must not quietly become a third way to
// teleport.
func TestGoingInYourselfStillMeansGoingThere(t *testing.T) {
	w := standing(t, "bar")
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && l.ID != w.Player.Location {
			if f := w.faction(prop.Owner); f != nil && f.ID != w.PlayerOrganizationID() {
				for _, a := range w.Actions(l.ID) {
					if a.ID == "sabotage" {
						t.Fatalf("%s offers going in yourself from across the city", l.ID)
					}
				}
			}
		}
	}
}
