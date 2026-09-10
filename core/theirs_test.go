package core

import (
	"strings"
	"testing"
)

// "We could probably also extend to be able to provide cars, and armor the cars
// for people in our family to keep them more protected from attacks."
//
// A car of their own is not a number on somebody else's sheet. It is what gets
// them off the street when a job goes wrong, which is the one place in this
// game where the difference between sending somebody and going yourself is
// actually decided.

func patron(t *testing.T) *World {
	t.Helper()
	w := New(83)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 40000, 100
	// Somebody of yours, with a name and a place in the city.
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || IsOfficial(n.ID) || n.Rank >= RankLeader {
			continue
		}
		n.Faction, n.Rank, n.Car, n.Trust = w.PlayerOrganizationID(), RankSoldier, 0, 60
		w.Player.Crew = []Crew{{n.ID, n.Name, 70}}
		n.Location = w.theForecourt()
		w.Player.Location = n.Location
		if w.Player.Location == "" {
			t.Skip("this city has no forecourt")
		}
		return w
	}
	t.Fatal("nobody to put on")
	return nil
}

func TestYouCanPutOneOfYoursInACar(t *testing.T) {
	w := patron(t)
	who := w.Player.Crew[0]
	a := actionByID(w.Actions(w.Player.Location), "car:"+who.ID)
	if a == nil {
		t.Fatalf("the forecourt will not sell a car for %s", who.Name)
	}
	if a.Disabled {
		t.Fatalf("refused: %s", a.Reason)
	}
	cash := w.Player.Cash
	if err := w.BuyCarFor(who.ID); err != nil {
		t.Fatal(err)
	}
	n := w.NPC(who.ID)
	if n.Car == 0 {
		t.Fatal("they are still walking")
	}
	if w.Player.Cash >= cash {
		t.Fatal("the car was free")
	}
	// And not twice.
	if w.TheirCarReadiness(who.ID) == "" {
		t.Fatal("bought a second car for somebody who has one")
	}
}

func TestYouCanPlateTheirCarAtAGarage(t *testing.T) {
	w := patron(t)
	who := w.Player.Crew[0]
	if err := w.BuyCarFor(who.ID); err != nil {
		t.Fatal(err)
	}
	garage := w.theGarage()
	if garage == "" {
		t.Skip("this city has no garage")
	}
	w.Player.Location = garage
	w.NPC(who.ID).Location = garage
	a := actionByID(w.Actions(garage), "plate:"+who.ID)
	if a == nil {
		t.Fatalf("the garage will not plate %s's car", who.Name)
	}
	before := w.TheirPlating(who.ID)
	if err := w.PlateTheirCar(who.ID); err != nil {
		t.Fatal(err)
	}
	if w.TheirPlating(who.ID) != before+1 {
		t.Fatalf("plating went from %d to %d", before, w.TheirPlating(who.ID))
	}
	// Somebody on foot has nothing to plate.
	w.NPC(who.ID).Car = 0
	if reason := w.TheirPlateReadiness(who.ID); !strings.Contains(strings.ToLower(reason), "car") {
		t.Fatalf("plate was offered for a man with no car: %q", reason)
	}
}

func TestACarOfTheirOwnIsWhatGetsThemOut(t *testing.T) {
	const runs = 500
	away := map[string]int{}
	for _, kit := range []string{"walking", "driving", "plated"} {
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(seed * 2654435761)
			w.Event, w.District = nil, 9
			w.Player.Cash, w.Player.Health = 40000, 100
			var who *NPC
			for i := range w.NPCs {
				n := &w.NPCs[i]
				if n.Dead || IsOfficial(n.ID) || n.Rank >= RankLeader {
					continue
				}
				n.Faction, n.Rank, n.Trust = w.PlayerOrganizationID(), RankSoldier, 60
				w.Player.Crew = []Crew{{n.ID, n.Name, 70}}
				who = n
				break
			}
			if who == nil {
				t.Skip("nobody to send")
			}
			switch kit {
			case "walking":
				who.Car, who.Plate = 0, 0
			case "driving":
				who.Car, who.Plate = 1, 0
			case "plated":
				who.Car, who.Plate = 1, PlateStages
			}
			mark := ""
			for i := range w.NPCs {
				n := &w.NPCs[i]
				if !n.Dead && n.Faction != "" && n.Faction != w.PlayerOrganizationID() {
					mark = n.ID
					break
				}
			}
			if mark == "" {
				t.Skip("nobody to go after")
			}
			w.itWentWrong(w.NPC(mark), Hand{Crew: true, Name: who.Name}, "the street", nil)
			if n := w.NPC(who.ID); n != nil && !n.Dead && n.Held <= w.Minute {
				away[kit]++
			}
		}
	}
	t.Logf("of %d jobs that went wrong: walking got away %d, driving %d, plated %d",
		runs, away["walking"], away["driving"], away["plated"])
	if away["driving"] <= away["walking"] {
		t.Fatalf("a car got nobody out: %d against %d", away["driving"], away["walking"])
	}
	if away["plated"] <= away["driving"] {
		t.Fatalf("plate on their car was worth nothing: %d against %d", away["plated"], away["driving"])
	}
}
