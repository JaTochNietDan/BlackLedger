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
