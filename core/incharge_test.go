package core

import "testing"

// Somebody comes to run a business in this city already: one of the player's
// own people, when they have had enough of being badly paid, walks off with a
// holding and the paper prints it. The player cannot do it on purpose. Every
// business they hold is a place they have to walk to and restock by hand, which
// is the work a manager exists to take off them.

func toRun(t *testing.T) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	w.Player.Location = "laundry"
	return w, "laundry"
}

func TestYouCanPutSomebodyInChargeOfABusiness(t *testing.T) {
	t.Parallel()
	w, id := toRun(t)
	who := w.Properties[id].Hands[0]
	place, _ := PlaceByID(id)
	if reason := w.InChargeReadiness(id, who); reason != "" {
		t.Fatalf("putting one of your own in charge was refused: %s", reason)
	}
	if err := w.PutInCharge(id, who); err != nil {
		t.Fatalf("it failed: %v", err)
	}
	n := w.NPC(who)
	if n.Role != "Runs "+place.Name {
		t.Fatalf("%s was put in charge and their role reads %q", n.Name, n.Role)
	}
	if w.RunsIt(id) == nil || w.RunsIt(id).ID != who {
		t.Fatalf("the laundry is run by %v", w.RunsIt(id))
	}
	// They are still one of the hands: a manager stands behind the counter too.
	if w.EmployerOf(who) != id {
		t.Fatalf("%s runs the place and is not on its books", n.Name)
	}
	// And there is only one of them.
	other := w.Properties[id].Hands[1]
	if w.InChargeReadiness(id, other) == "" {
		t.Fatal("a business can have two people running it")
	}
}

// What a manager is for: the place keeps itself stocked instead of the player
// walking over to do it.
func TestAManagerKeepsThePlaceStocked(t *testing.T) {
	t.Parallel()
	managed, id := toRun(t)
	alone, _ := toRun(t)
	who := managed.Properties[id].Hands[0]
	if err := managed.PutInCharge(id, who); err != nil {
		t.Fatal(err)
	}
	for _, w := range []*World{managed, alone} {
		w.Properties[id].Supply = 0
		w.Player.Location = "bar"
	}
	for day := 0; day < 10; day++ {
		for _, w := range []*World{managed, alone} {
			w.Event = nil
			w.Advance(1440)
			w.Event = nil
		}
	}
	t.Logf("after ten days away: run by somebody it holds %d supplies, run by nobody %d",
		managed.Properties[id].Supply, alone.Properties[id].Supply)
	if managed.Properties[id].Supply <= alone.Properties[id].Supply {
		t.Fatal("somebody running the place kept it no better stocked than nobody did")
	}
	// And it is the player's money that bought them.
	if managed.Player.Cash >= alone.Player.Cash {
		t.Fatal("the supplies were bought with nobody's money")
	}
}

// It is offered where the business is, and only to somebody who works there.
func TestPuttingSomebodyInChargeIsOfferedAtTheBusiness(t *testing.T) {
	t.Parallel()
	w, id := toRun(t)
	who := w.Properties[id].Hands[0]
	var put *Action
	for _, a := range w.Actions(id) {
		if a.ID == "incharge:"+who {
			put = &a
		}
	}
	if put == nil {
		t.Fatalf("%s works here and cannot be put in charge of it", w.NPC(who).Name)
	}
	if put.Disabled {
		t.Fatalf("it is refused where it should work: %s", put.Reason)
	}
	if err := w.apply(Command{Kind: "incharge:" + who, Target: id, RequestID: "putthemincharge1"}); err != nil {
		t.Fatalf("through the same path as everything else it failed: %v", err)
	}
	if w.RunsIt(id) == nil {
		t.Fatal("nobody runs the laundry")
	}
}
