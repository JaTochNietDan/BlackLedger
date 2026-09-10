package main

import (
	"blackledger/core"
	"reflect"
	"testing"
)

func TestDirectorContactsGrowWithOrganization(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	w := core.New(27)
	w.Life = 2
	w.Player.Crew = []core.Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 65}}
	w.Arrangements = []core.ArrangementMemory{{Life: 1, Speaker: "mara"}, {Life: 1, Speaker: "leo"}}
	if got := directorSpeakers(w, nil); !reflect.DeepEqual(got, []string{"mara", "leo"}) {
		t.Fatal(got)
	}
}

func TestTheDirectorDoesNotSpeakThroughTheDead(t *testing.T) {
	t.Parallel()
	w := core.New(21)
	w.Factions[0].Goodwill = 40 // the Bellandi leader would otherwise be eligible
	before := eligibleDirectorSpeakers(w)
	if !before["vittorio"] {
		t.Fatal("a living family leader with standing was not eligible")
	}
	w.Kill("vittorio", "Shot leaving the club.")
	after := eligibleDirectorSpeakers(w)
	if after["vittorio"] {
		t.Fatal("a dead leader is still offered as a speaker")
	}
	for _, id := range directorSpeakers(w, nil) {
		if npc := w.NPC(id); npc == nil || npc.Dead {
			t.Fatal("the speaker list contains someone who is dead:", id)
		}
	}
	// Their successor can speak in their place.
	w.Factions[0].Goodwill = 40
	successor := w.Factions[0].Leader
	found := false
	for id := range eligibleDirectorSpeakers(w) {
		if npc := w.NPC(id); npc != nil && npc.Name == successor {
			found = true
		}
	}
	if !found {
		t.Fatal("the new head of the family cannot speak for it")
	}
}

func TestMaraCanDieAndStopsBeingAvailable(t *testing.T) {
	t.Parallel()
	w := core.New(22)
	if !eligibleDirectorSpeakers(w)["mara"] {
		t.Fatal("the starting contact was not eligible")
	}
	w.Kill("mara", "Found behind Saint Agnes.")
	if eligibleDirectorSpeakers(w)["mara"] {
		t.Fatal("a dead fixer is still taking work")
	}
}
