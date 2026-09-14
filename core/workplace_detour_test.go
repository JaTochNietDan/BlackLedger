package core

import "testing"

func TestTemporaryErrandsDoNotReplaceTheDayJob(t *testing.T) {
	for _, reason := range []string{"wanting petrol at Kessler's", "taking the car in to Russo", "buying another car at Ferris", "looking for somebody over a debt", "on collections at Kessler's", "coming back from the round"} {
		w := New(404)
		w.Minute = 540
		n := w.NPC("street-47")
		n.Faction = ""
		n.Role = "Net mender"
		n.Rank = RankAssociate
		n.Post = "docks"
		n.Home = "room"
		n.Location = "docks"
		n.Heading = "filling"
		n.Arrives = w.Minute
		n.Sets = 0
		n.Errand = reason
		w.Arrivals()
		if n.Post != "docks" {
			t.Fatalf("%s replaced work with %s", reason, n.Post)
		}
		if n.Location != "filling" {
			t.Fatal("arrival teleported back to work")
		}
	}
}
func TestCompletedFuelErrandStartsARealReturnJourney(t *testing.T) {
	w := New(404)
	w.Minute = 540
	n := w.NPC("street-47")
	n.Faction = ""
	n.Role = "Net mender"
	n.Rank = RankAssociate
	n.Post = "docks"
	n.Home = "room"
	n.Location = "filling"
	n.Heading = "filling"
	n.Arrives = w.Minute
	n.Sets = 0
	n.Errand = "wanting petrol at Kessler's"
	n.Car = 1
	n.Dry = true
	n.Purse = 500
	w.Arrivals()
	if n.Dry || n.Purse != 500-FuelPrice {
		t.Fatal("fuel was not purchased before returning")
	}
	if n.Post != "docks" || n.Heading != "docks" || n.Arrives <= w.Minute || n.Location != "filling" {
		t.Fatal("temporary visit did not lead back to the actual job")
	}
}
