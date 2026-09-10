package core

import "testing"

func TestNewLifeCanEarnAndBuyFormerBusinessWithoutInheritance(t *testing.T) {
	t.Parallel()
	w := New(27)
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["laundry"].Condition = 60
	w.Die("Previous life ended")
	act(t, &w, "new_life", "")
	if w.Own("laundry") || w.Player.Cash != 90 {
		t.Fatal("new life inherited assets")
	}
	w.Player.Location = "laundry"
	w.Player.Respect = 6
	// A laundry off a former organization is twice its listed price, and the
	// listed price is four times what it was: premises are meant to be built up
	// to. The figures are read off the game rather than typed, so this test
	// stays about paying the full price rather than about what the price is.
	full := AcquisitionCost(w, "laundry")
	w.Player.Cash = full - 1
	if _, err := Execute(w, Command{Kind: "acquire", Target: "laundry", Revision: w.Revision}); err == nil {
		t.Fatal("buyout accepted without full payment")
	}
	w.Player.Cash = full
	act(t, &w, "acquire", "laundry")
	if !w.Own("laundry") || w.Player.Cash != 0 || w.Properties["laundry"].Condition != 60 || w.Properties["laundry"].Income != 12 {
		t.Fatal("buyout failed to preserve the changed city or charge its price")
	}
}
func TestBuyoutDoesNotSeizeFactionProperty(t *testing.T) {
	t.Parallel()
	w := New(27)
	w.Player.Location = "laundry"
	w.Player.Respect = 100
	w.Player.Cash = 10000
	w.Properties["laundry"].Owner = "bellandi"
	if _, err := Execute(w, Command{Kind: "acquire", Target: "laundry"}); err == nil {
		t.Fatal("ordinary buyout seized faction property")
	}
}
