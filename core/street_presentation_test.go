package core

import (
	"encoding/json"
	"math"
	"testing"
)

func TestTravelRecordsDepartureAndArrivalWithinOneCommand(t *testing.T) {
	w := New(2)
	w.District = 2
	w.Player.Location = "bar"
	w.Player.Cash = 10000
	w.Player.Heat = 0
	w.Event = nil
	w.Plots = nil
	w.Tasks = nil
	w.Contracts = nil
	w.NextPressure = 0
	w.NPCs = w.NPCs[:1]
	n := &w.NPCs[0]
	n.Dead = false
	n.Location = "bar"
	n.Heading = "market"
	n.Sets = w.Minute + 2
	n.Arrives = n.Sets + TravelMinutes("bar", "market")
	n.Car = 1
	n.Dry = false
	n.Hurt = false
	start, arrival, id := n.Sets, n.Arrives, n.ID
	next, err := Execute(w, Command{Revision: w.Revision, RequestID: ID(), Kind: "travel", Target: "dealer"})
	if err != nil {
		t.Fatal(err)
	}
	segments := next.LastResult.StreetTravel
	if len(segments) != 1 {
		t.Fatalf("wanted one observed leg, got %+v", segments)
	}
	leg := segments[0]
	if leg.ID != id || leg.FromMinute != start || leg.ToMinute != arrival || leg.Progress != 0 || leg.EndProgress != 1 || leg.Vehicle == "" {
		t.Fatalf("incorrect observed journey: %+v", leg)
	}
	if leg.ToMinute > next.Minute {
		t.Fatal("future travel disclosed")
	}
	data, err := json.Marshal(next)
	if err != nil {
		t.Fatal(err)
	}
	var restored World
	if err = json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if len(restored.LastResult.StreetTravel) != 1 || restored.recordStreet || restored.streetTravel != nil {
		t.Fatal("receipt must persist, transient recorder must not")
	}
}

func TestStreetSegmentsCoalesceAndStopAtCommandBoundary(t *testing.T) {
	w := New(2)
	j := Journeying{ID: "test", FromID: "bar", ToID: "market", Minutes: 20, Progress: .2}
	w.recordStreetSegment([]Journeying{j}, 100, 105)
	j.Minutes = 15
	j.Progress = .4
	w.recordStreetSegment([]Journeying{j}, 105, 110)
	if len(w.streetTravel) != 1 || w.streetTravel[0].FromMinute != 100 || w.streetTravel[0].ToMinute != 110 || math.Abs(w.streetTravel[0].EndProgress-.6) > 1e-9 {
		t.Fatalf("wrong continuous leg: %+v", w.streetTravel)
	}
}

func TestEmptyRecordedStreetIsDistinctFromLegacyMissingTrace(t *testing.T) {
	w := New(2)
	w.Player.Location = "bar"
	w.NPCs = nil
	w.Event = nil
	w.Plots = nil
	w.Tasks = nil
	w.NextPressure = 0
	next, err := Execute(w, Command{Revision: w.Revision, RequestID: ID(), Kind: "travel", Target: "market"})
	if err != nil {
		t.Fatal(err)
	}
	if next.LastResult.StreetTravel == nil || len(next.LastResult.StreetTravel) != 0 {
		t.Fatal("new travel must publish an explicit empty street trace")
	}
}
