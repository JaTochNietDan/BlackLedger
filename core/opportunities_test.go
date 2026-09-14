package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestOpeningGuidanceFindsPaidWorkWhenEnvelopesRunOut(t *testing.T) {
	w := New(81)
	w.Event = nil
	w.Player.JobCount = 1
	w.Player.Respect = 2
	w.Player.CarriedDay = w.Minute / 1440
	w.Player.Carried = CourierADay
	next := w.NextOpportunity()
	if next == nil || next.Target != "laundry" || !strings.Contains(next.Detail, "$80") {
		t.Fatalf("no limited daytime alternative: %+v", next)
	}
	w.Properties["laundry"].RushOrderDay = w.Minute/1440 + 1
	next = w.NextOpportunity()
	if next == nil || next.Target != "docks" || !strings.Contains(next.Detail, "$75") || !strings.Contains(next.Detail, "90 minutes") {
		t.Fatalf("no actionable cargo fallback: %+v", next)
	}
	// The next morning's real envelope availability changes the recommendation.
	w.Minute = (w.Minute/1440+1)*1440 + 480
	next = w.NextOpportunity()
	if next == nil || next.Target != "bar" || !strings.Contains(next.Detail, "45 minutes") {
		t.Fatalf("new day's jobs not available: %+v", next)
	}
}

func TestOpeningGuidanceSavesBeforeOfferingUnaffordableProperty(t *testing.T) {
	w := New(81)
	w.Event = nil
	w.Player.JobCount = 3
	w.Player.Respect = PremisesRespect
	w.Player.Cash = 20
	next := w.NextOpportunity()
	if next == nil || next.Title != "Save for your first business" || !strings.Contains(next.Detail, "you have $20") {
		t.Fatalf("sent broke player to a purchase: %+v", next)
	}
	w.Player.Cash = AcquisitionCost(w, "laundry")
	next = w.NextOpportunity()
	a, ok := w.opportunityAction("laundry", "acquire")
	if !ok || a.Disabled || next == nil || next.Target != "laundry" || !strings.Contains(next.Detail, a.Detail) {
		t.Fatalf("funded purchase not grounded in actual terms: %+v / %+v", next, a)
	}
}

func TestOpeningGuidanceDoesNotMoveThePlayerOrOfferWorkInCustody(t *testing.T) {
	w := New(81)
	w.Event = nil
	w.Player.JobCount = 3
	w.Player.Respect = 6
	w.Player.Cash = 20
	before, _ := json.Marshal(w)
	for i := 0; i < 3; i++ {
		_ = w.NextOpportunity()
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("reading next move changed world")
	}
	w.Confine(3, "test detention")
	if w.NextOpportunity() != nil {
		t.Fatal("detained player offered street work")
	}
}

func TestOpeningGuidanceQuotesActualRecruitmentAndSkipsUnavailableDriver(t *testing.T) {
	w := New(81)
	w.Event = nil
	w.Player.JobCount = 3
	w.Player.Respect = 6
	w.Player.Cash = 89
	w.Properties["laundry"].Owner = "player:" + itoa(w.Life)
	w.Properties["laundry"].Condition = 100
	next := w.NextOpportunity()
	if next == nil || next.Title != "Set aside the hiring money" {
		t.Fatalf("unfunded recruitment: %+v", next)
	}
	w.Player.Cash = 90
	next = w.NextOpportunity()
	a, ok := w.opportunityAction("bar", "recruit")
	if !ok || a.Disabled || next == nil || !strings.Contains(next.Detail, a.Label) || !strings.Contains(next.Detail, a.Detail) {
		t.Fatalf("recruitment drifted from button: %+v / %+v", next, a)
	}
	if driver := w.Holder("driver"); driver != nil {
		driver.Dead = true
	}
	next = w.NextOpportunity()
	if next != nil && next.Title == "Bring someone into the fold" {
		t.Fatalf("offered unavailable driver: %+v", next)
	}
}
