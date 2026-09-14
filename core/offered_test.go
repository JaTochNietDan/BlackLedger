package core

import (
	"strings"
	"testing"
)

// An action the interface offers as available must actually work.
//
// This is the fault shape this project keeps hitting and had no guard for. The
// car button lived in the garage's case while its rule asked for a forecourt:
// offered where it was always refused. The holder's button to set the house
// limit was offered with no field on it, so it sent an amount of nothing and
// was refused every single time. Both were found by a person pressing them.
//
// A card that says you can do something and a rule that then refuses costs the
// player the trip, and it is checkable: take every action the city offers as
// available, press it on a copy, and require it to succeed or to fail for a
// reason the card had already given.
func TestEveryActionOfferedCanActuallyBeTaken(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry", "garage", "casino")
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Contacts = 90000, 200, 4
	w.Player.Health = 100
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 70}}
	w.Player.Car, w.Player.CarWear = 2, 100
	w.Player.Fuel, w.Player.Fuelled = FuelFull, max(1, w.Minute)
	w.Player.Dress, w.Player.DressWear = 2, 100
	w.Player.Stock = map[string]int{"moonshine": 4}
	w.ensureOfficials()
	w.OrganizationDay()

	tried, refused := 0, []string{}
	for _, l := range Locations {
		if l.District > w.District {
			continue
		}
		here := w.Clone()
		here.Player.Location = l.ID
		here.Event = nil
		for _, a := range here.Actions(l.ID) {
			if a.Disabled {
				continue
			}
			// A figure the card asks for is a figure the player would type, so
			// press it with what the field starts on.
			// The second name, where an action has one. Asking somebody where a
			// third person is needs both, and replaying it without the second
			// asked about nobody — which the guard reported as the action being
			// offered and then refused, when what was missing was the command.
			c := Command{Kind: a.ID, Target: a.Target, Choice: a.Choice,
				Revision: here.Revision, RequestID: "offered-check"}
			if a.Sum != nil {
				c.Amount = a.Sum.Preset
			}
			tried++
			if _, err := Execute(here, c); err != nil {
				refused = append(refused, l.ID+"/"+a.ID+": "+err.Error())
			}
		}
	}
	if tried < 60 {
		t.Fatalf("only %d actions were offered anywhere; this is not measuring what it claims to", tried)
	}
	if len(refused) > 0 {
		t.Fatalf("%d of %d actions were offered as available and then refused:\n  %s",
			len(refused), tried, strings.Join(refused, "\n  "))
	}
	t.Logf("pressed %d actions the city offered as available; every one of them worked", tried)
}

// And the other half: a card that is offered everywhere it appears and refused
// everywhere it appears is a card nobody can ever press.
//
// This is the first fault shape in the loop's own list. The car button lived in
// the garage's case while its rule asked for a forecourt, so it was drawn in the
// one room where it could never work and absent from the room where it would
// have. Guard A above cannot see that: the card is honestly greyed out with a
// reason.
//
// The trick is what "can never" means. The first version of this asked one rich
// campaign and reported fourteen faults, every one of them false — the tank was
// full, the premises fully staffed, nobody was looking at the player. Those are
// refusals that clear the moment the state changes, which is a button doing its
// job. So it asks established campaigns that have or need everything, plus a new
// arrival with no premises. An action must be usable in at least one.
func TestNoActionIsOfferedOnlyWhereItIsRefused(t *testing.T) {
	t.Parallel()
	seen, usable := map[string]string{}, map[string]bool{}
	for _, shape := range []string{"comfortable", "needy", "starting"} {
		w := proprietor(t)
		own(w, "laundry", "garage", "casino")
		w.District = 2
		w.Player.Cash, w.Player.Respect, w.Player.Contacts = 90000, 200, 4
		w.Player.Health = 100
		w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 70}}
		w.Player.Car, w.Player.CarWear = 1, 100
		w.Player.Fuel, w.Player.Fuelled = FuelFull, max(1, w.Minute)
		w.Player.Dress, w.Player.DressWear = 1, 100
		w.ensureOfficials()
		w.OrganizationDay()
		if shape == "comfortable" {
			// Everything in hand: a crew, premises trading well, an account
			// abroad that already answers to you.
			w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 70}}
			w.Player.Offshore, w.Offshore = true, 2000
			// Somebody who answers to you and is not standing anywhere yet, so
			// there is a person to put on a door.
			for _, n := range w.Civilians() {
				if IsOfficial(n.ID) || w.isRoleHolder(n) || len(w.OwnPeople()) >= 2 {
					continue
				}
				w.Player.Location = n.Location
				w.SignOn(n.ID)
			}
			for _, l := range Locations {
				if prop := w.Properties[l.ID]; prop != nil && w.Own(l.ID) {
					prop.Custom = 95
				}
			}
		}
		if shape == "needy" {
			// A city where everything wants doing: the premises run down and
			// short-handed, the tank dry, the suit ruined, the police
			// interested, and money out there to be brought home.
			for _, l := range Locations {
				if prop := w.Properties[l.ID]; prop != nil && w.Own(l.ID) {
					prop.Condition, prop.Staff, prop.Supply, prop.Trouble = 40, 0, 0, true
					prop.Bankroll, prop.Custom = 4000, 90
				}
			}
			// And nothing in hand: nobody working for you, a car that wants
			// looking at, money abroad you have not yet proved is yours.
			w.Player.Fuel, w.Player.DressWear, w.Player.Heat = 0, 10, 60
			w.Player.CarWear = 45
			w.Player.Crew = nil
			for _, l := range Locations {
				if prop := w.Properties[l.ID]; prop != nil {
					prop.Posted = ""
				}
			}
			w.Player.Offshore, w.Offshore = false, 2000
			w.Player.Stock = map[string]int{"moonshine": 4}
		}
		// A newcomer can take helper work that an owner cannot. Include a
		// real starting state instead of exempting those actions from coverage.
		if shape == "starting" {
			w = New(27)
			w.Event = nil
		}
		for _, l := range Locations {
			if l.District > w.District {
				continue
			}
			here := w.Clone()
			here.Player.Location = l.ID
			here.Event = nil
			for _, a := range here.Actions(l.ID) {
				if _, had := seen[a.ID]; !had {
					seen[a.ID] = l.ID + ": " + a.Reason
				}
				if !a.Disabled {
					usable[a.ID] = true
				}
			}
		}
	}
	stuck := []string{}
	for id, where := range seen {
		if usable[id] || situational[id] || strings.Contains(id, ":") {
			continue
		}
		stuck = append(stuck, id+" ("+where+")")
	}
	if len(stuck) > 0 {
		t.Fatalf("%d actions are drawn in every room they appear in and refused in every one of them: %v",
			len(stuck), stuck)
	}
	t.Logf("%d distinct actions across three campaigns, every one of them usable somewhere", len(seen))
}

// situational names the work that is deliberately drawn before it can be done,
// because the room is where the player learns it exists. A cell offers the three
// ways out of it whether or not anybody is in one.
var situational = map[string]bool{
	"sit_out": true, "lawyer": true, "talk": true, // only from inside a cell
	"rise": true, "hit": true, "stand": true, "roll": true, // only mid-game at a table
	"stock_arms": true, "dismantle": true, "unpost": true, // only once the thing exists
	// The paper takes your calls once you hold it or pay somebody there, and
	// neither campaign here does. Offered at the Herald so that walking in is
	// how a player learns the two things a newspaper can be made to do.
	"spike": true, "puff": true,
}
