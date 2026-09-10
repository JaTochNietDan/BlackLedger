package core

import (
	"strings"
	"testing"
)

// "I go to The Monarch and request a sit down with the controlling family and it
// sits me down with the family lead but the lead is not present at this
// location. That lead is at Kesseler Filling Station..."
//
// And the other half of the same message: "Vitor Bellendi is always at the
// kessler filling station for some reason. That seems odd?"

func caller(t *testing.T) *World {
	t.Helper()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 3000, 100
	return w
}

func TestAFamilyLeadKeepsTheirOwnSeat(t *testing.T) {
	t.Parallel()
	w := caller(t)
	lead := w.Leader("bellandi")
	if lead == nil {
		t.Fatal("no Bellandi lead")
	}
	// Put him somewhere he has no business being and let the city run.
	lead.Location, lead.Post, lead.Heading, lead.Arrives = "filling", "filling", "", 0
	seat := w.homeOf("bellandi")
	if seat == "filling" {
		t.Skip("the filling station is this family's seat in this city")
	}
	for day := 0; day < 4 && w.Leader("bellandi").Location != seat; day++ {
		w.Advance(1440)
	}
	if got := w.Leader("bellandi").Location; got != seat {
		t.Fatalf("four days later the lead is still at %s rather than %s", got, seat)
	}
}

func TestYouAreReceivedByWhoeverIsActuallyThere(t *testing.T) {
	t.Parallel()
	w := caller(t)
	lead := w.Leader("bellandi")
	if lead == nil {
		t.Fatal("no Bellandi lead")
	}
	// Nobody of theirs at the club at all.
	for i := range w.NPCs {
		if w.NPCs[i].Faction == "bellandi" && w.NPCs[i].Location == "club" {
			w.NPCs[i].Location = "docks"
		}
	}
	lead.Location, lead.Heading = "filling", ""
	w.Player.Location = "club"
	reason := w.AudienceReadiness("club")
	if reason == "" {
		t.Fatal("granted a sit-down in a room with nobody from the family in it")
	}
	place, _ := PlaceByID("filling")
	if !strings.Contains(reason, place.Name) {
		t.Fatalf("the refusal does not say where they are: %q", reason)
	}
	// With the lead in the room it is granted, and the chair holds him.
	lead.Location = "club"
	if reason := w.AudienceReadiness("club"); reason != "" {
		t.Fatalf("refused with the lead standing there: %s", reason)
	}
	w.OpenAudience("club")
	if w.Event == nil || w.Event.Speaker != lead.ID {
		t.Fatalf("the chair does not hold the lead: %+v", w.Event)
	}
}

func TestSomebodyOfTheirsCanSpeakForThem(t *testing.T) {
	t.Parallel()
	w := caller(t)
	lead := w.Leader("bellandi")
	lead.Location = "filling"
	deputy := ""
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Faction == "bellandi" && !n.Dead && n.ID != lead.ID && n.Rank >= RankLieutenant {
			n.Location, n.Heading = "club", ""
			deputy = n.ID
			break
		}
	}
	if deputy == "" {
		t.Skip("this family has nobody below the lead who could speak for it")
	}
	w.Player.Location = "club"
	if reason := w.AudienceReadiness("club"); reason != "" {
		t.Fatalf("nobody would sit down with a lieutenant of theirs: %s", reason)
	}
	w.OpenAudience("club")
	if w.Event == nil {
		t.Fatal("no sit-down happened")
	}
	if w.Event.Speaker != deputy {
		t.Fatalf("somebody who is not in the room took the chair: %q", w.Event.Speaker)
	}
}
