package core

import (
	"strings"
	"testing"
)

func TestTravelArrivalAtDecisionBoundaryIsNotLost(t *testing.T) {
	w := pressureWorld()
	w.District = 1
	w.Player.Location = "bar"
	duration := TravelMinutes("bar", "apartment")
	w.NextPressure = w.Minute + duration
	next, err := Execute(w, Command{Kind: "travel", Target: "apartment", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Event == nil || next.Event.Kind != "business_pressure" || next.Player.Location != "apartment" || next.Minute != w.Minute+duration {
		t.Fatal("complete journey lost when decision arrived")
	}
	choice(t, &next, "pay")
	if next.Player.Location != "apartment" {
		t.Fatal("resolving demand moved character back")
	}
}
func TestPartialTravelExplainsWherePlayerRemains(t *testing.T) {
	w := pressureWorld()
	w.District = 1
	w.Player.Location = "bar"
	w.NextPressure = w.Minute + 10
	next, err := Execute(w, Command{Kind: "travel", Target: "apartment", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Location != "bar" || next.Minute != w.Minute+10 || next.Event == nil {
		t.Fatal("partial travel completed incorrectly")
	}
	r := next.History[len(next.History)-1]
	if r.Title != "Journey interrupted" || !strings.Contains(r.Text, "10 minutes") || !strings.Contains(r.Text, "Saint Agnes") {
		t.Fatal("missing partial travel explanation")
	}
}

func TestPreparedEncounterDeliveredOnArrival(t *testing.T) {
	for _, dueAtArrival := range []bool{false, true} {
		w := New(27)
		duration := TravelMinutes("room", "bar")
		ready := w.Minute
		if dueAtArrival {
			ready += duration
		}
		e, err := w.ValidateProposal(Proposal{Title: "A quiet delivery", Body: "Mara has a sealed message to deliver.", Speaker: "mara", Operation: "courier", Outcome: "The message was delivered."})
		if err != nil {
			t.Fatal(err)
		}
		w.Offers = []Offer{{Ready: ready, Event: e}}
		next, err := Execute(w, Command{Kind: "travel", Target: "bar", Revision: w.Revision})
		if err != nil {
			t.Fatal(err)
		}
		if next.Player.Location != "bar" || next.Minute != w.Minute+duration || next.Event == nil || next.Event.ID != e.ID || len(next.Offers) != 0 {
			t.Fatal("arrival failed to deliver due encounter without extra time")
		}
		if w.Event != nil || len(w.Offers) != 1 {
			t.Fatal("command mutated original state")
		}
		choice(t, &next, "decline")
		if next.Event != nil || len(next.Offers) != 0 {
			t.Fatal("delivered encounter repeated")
		}
	}
}

func TestArrivalEncounterWaitsBehindUrgentIncident(t *testing.T) {
	w := pressureWorld()
	w.District = 1
	w.Player.Location = "bar"
	duration := TravelMinutes("bar", "apartment")
	w.NextPressure = w.Minute + duration
	e, err := w.ValidateProposal(Proposal{Title: "A quiet delivery", Body: "Mara has a sealed message to deliver.", Speaker: "mara", Operation: "courier", Outcome: "The message was delivered."})
	if err != nil {
		t.Fatal(err)
	}
	w.Offers = []Offer{{Ready: w.Minute, Event: e}}
	act(t, &w, "travel", "apartment")
	if w.Event == nil || w.Event.Kind != "business_pressure" || len(w.Offers) != 1 || w.Offers[0].Event.ID != e.ID {
		t.Fatal("routine offer displaced urgent incident or was consumed")
	}
}
