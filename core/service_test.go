package core

import "testing"

func recruit(t *testing.T) (*World, *Faction) {
	t.Helper()
	w := New(281)
	w.MigrateLivingWorld()
	w.Player.Respect, w.Player.Cash, w.Player.Health = ServiceRespect, 2000, 100
	f := &w.Factions[0]
	f.Goodwill, f.Cash = ServiceGoodwill, 20000
	w.Player.Location = w.homeOf(f.ID)
	return w, f
}

func TestNobodyTakesOnSomebodyTheyDoNotWant(t *testing.T) {
	t.Parallel()
	w, f := recruit(t)
	if w.ServeReadiness(f.ID) != "" {
		t.Fatal("could not go to work for them:", w.ServeReadiness(f.ID))
	}
	f.Goodwill = ServiceGoodwill - 1
	if w.ServeReadiness(f.ID) == "" {
		t.Fatal("somebody they barely tolerate signed on")
	}
	f.Goodwill = ServiceGoodwill
	w.Player.Respect = ServiceRespect - 1
	if w.ServeReadiness(f.ID) == "" {
		t.Fatal("a nobody signed on")
	}
	w.Player.Respect = ServiceRespect
	w.Player.Location = "room"
	if w.ServeReadiness(f.ID) == "" {
		t.Fatal("that conversation happened across the city")
	}
}

func TestAManWithHisOwnThingDoesNotComeUpThroughSomebodyElses(t *testing.T) {
	t.Parallel()
	w, f := recruit(t)
	own(w, "laundry", "garage")
	w.Player.Respect = OrganizationStanding
	w.OrganizationDay()
	if !w.Incorporated() {
		t.Fatal("the fixture did not become an organization")
	}
	if w.ServeReadiness(f.ID) == "" {
		t.Fatal("an organization went to work for a family")
	}
	// And the other way: somebody in service does not become one.
	quiet := proprietor(t)
	own(quiet, "laundry", "garage")
	quiet.Player.Respect = OrganizationStanding
	quiet.Player.Serves = quiet.Factions[0].ID
	quiet.OrganizationDay()
	if quiet.Incorporated() {
		t.Fatal("somebody's soldier started their own thing")
	}
}

func TestAnsweringToSomebodyPaysAndMakesTheirEnemiesYours(t *testing.T) {
	t.Parallel()
	w, f := recruit(t)
	other := &w.Factions[1]
	c := w.Conflict(f.ID, other.ID)
	c.Hostility, c.State = 70, "war"
	goodwill := other.Goodwill
	if err := w.Serve(f.ID); err != nil {
		t.Fatal(err)
	}
	if w.Serving() != f.ID || w.ServiceTitle() != "Associate" {
		t.Fatalf("they are a %q of %q", w.ServiceTitle(), w.Serving())
	}
	if other.Goodwill >= goodwill {
		t.Fatal("going to work for somebody was not noticed by the people they fight")
	}
	cash := w.Player.Cash
	w.ServiceDay()
	if w.Player.Cash != cash+SoldierWage {
		t.Fatalf("a day's work paid $%d", w.Player.Cash-cash)
	}
	if f.Cash >= 20000 {
		t.Fatal("the wage came out of nowhere")
	}
	// Nobody stays where the money stops.
	f.Cash = 0
	cash = w.Player.Cash
	w.ServiceDay()
	if w.Player.Cash != cash {
		t.Fatal("they were paid out of an empty till")
	}
	if !w.hasRecord("Nothing came this week") {
		t.Fatal("nobody mentioned it")
	}
}

func TestComingUpTakesWorkAndPaysBetter(t *testing.T) {
	t.Parallel()
	w, f := recruit(t)
	w.Serve(f.ID)
	base := w.ServicePay()
	for i := 0; i < PromotionWork; i++ {
		w.ServeWork(f.ID)
	}
	if w.ServiceTitle() != "Soldier" {
		t.Fatalf("three pieces of work made them a %s", w.ServiceTitle())
	}
	for i := 0; i < PromotionWork; i++ {
		w.ServeWork(f.ID)
	}
	if w.ServiceTitle() != "Lieutenant" {
		t.Fatalf("six pieces of work made them a %s", w.ServiceTitle())
	}
	if w.ServicePay() <= base {
		t.Fatalf("a lieutenant is paid %d against %d at the bottom", w.ServicePay(), base)
	}
	if !w.hasRecord("They have moved you up") {
		t.Fatal("nobody was told")
	}
	// Work for somebody else is not work for them.
	before := w.Player.Service
	w.ServeWork(w.Factions[1].ID)
	if w.Player.Service != before {
		t.Fatal("work for a rival counted toward coming up")
	}
}

func TestWalkingOutIsNotForgiven(t *testing.T) {
	t.Parallel()
	w, f := recruit(t)
	w.Serve(f.ID)
	goodwill := f.Goodwill
	if err := w.LeaveService(); err != nil {
		t.Fatal(err)
	}
	if w.Serving() != "" || w.ServiceRank() != 0 {
		t.Fatal("they still answer to somebody")
	}
	if f.Goodwill != goodwill-LeavingCost {
		t.Fatalf("walking out cost %d standing", goodwill-f.Goodwill)
	}
	if len(w.Plots) == 0 {
		t.Fatal("nobody answered it")
	}
	if err := w.LeaveService(); err == nil {
		t.Fatal("left an organization twice")
	}
}

func TestWhoeverYouAnsweredToCanStopExisting(t *testing.T) {
	t.Parallel()
	w, f := recruit(t)
	w.Serve(f.ID)
	w.Dissolve(f.ID)
	w.ServiceDay()
	if w.Serving() != "" {
		t.Fatal("they still answer to an organization that does not exist")
	}
}

func TestNobodyInheritsAJob(t *testing.T) {
	t.Parallel()
	w, f := recruit(t)
	w.Serve(f.ID)
	w.Player.Alive = false
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	if next.Serving() != "" || next.ServicePay() != 0 {
		t.Fatal("the next arrival woke up working for somebody")
	}
}

// Most people are never offered an arrangement with a beneficiary. Harm done to
// somebody's enemies has to count, or coming up is a path only the lucky find.
func TestHarmToTheirEnemiesIsWorkForThem(t *testing.T) {
	t.Parallel()
	w, f := recruit(t)
	w.Serve(f.ID)
	other := w.Factions[1].ID
	before := w.Player.Service

	// Somebody they are on civil terms with is not an enemy.
	if c := w.Conflict(f.ID, other); c != nil {
		c.Hostility, c.State = 5, "cold"
	}
	w.ServeAgainst(other)
	if w.Player.Service != before {
		t.Fatal("harm to somebody they get on with counted as work")
	}

	if c := w.Conflict(f.ID, other); c != nil {
		c.Hostility, c.State = 70, "war"
	}
	w.ServeAgainst(other)
	if w.Player.Service != before+1 {
		t.Fatal("harm to somebody they are at war with counted as nothing")
	}
	// And harm to themselves is not work for them.
	w.ServeAgainst(f.ID)
	if w.Player.Service != before+1 {
		t.Fatal("attacking your own organization counted as work for it")
	}
}
