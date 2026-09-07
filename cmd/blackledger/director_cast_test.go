package main

import (
	"blackledger/core"
	"reflect"
	"testing"
)

func TestDirectorContactsGrowWithOrganization(t *testing.T) {
	w := core.New(27)
	check := func(want ...string) {
		t.Helper()
		if got := directorSpeakers(w, nil); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
	check("mara")
	w.Player.Crew = []core.Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 65}}
	w.Arrangements = []core.ArrangementMemory{{Life: w.Life, Speaker: "mara", Status: "declined"}}
	check("leo")
	w.Offers = []core.Offer{{Event: &core.Scene{Speaker: "leo"}}}
	check("mara")
	w.Factions[1].Goodwill = 6
	check("elena")
	w.Factions[0].Goodwill = 8
	check("vittorio", "elena")
	w.Factions[0].Goodwill = -40
	w.Factions[1].Goodwill = -40
	w.Player.Crew[0].Loyalty = 29
	check("mara")
	connection := &core.ArrangementMemory{Speaker: "elena"}
	if got := directorSpeakers(w, connection); !reflect.DeepEqual(got, []string{"elena"}) {
		t.Fatal("established follow-up lost its speaker")
	}
}

func TestDirectorContactRecencyIgnoresPreviousLives(t *testing.T) {
	w := core.New(27)
	w.Life = 2
	w.Player.Crew = []core.Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 65}}
	w.Arrangements = []core.ArrangementMemory{{Life: 1, Speaker: "mara"}, {Life: 1, Speaker: "leo"}}
	if got := directorSpeakers(w, nil); !reflect.DeepEqual(got, []string{"mara", "leo"}) {
		t.Fatal(got)
	}
}
