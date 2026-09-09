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
