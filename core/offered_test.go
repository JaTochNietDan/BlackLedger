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
			c := Command{Kind: a.ID, Target: a.Target, Revision: here.Revision, RequestID: "offered-check"}
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
