package core

import "testing"

// From the inbox: people steal car parts and sell them to garages, and a garage
// does better when there is more of it about. Both halves need a car that
// belongs to somebody, which the city has now.
//
// A stripped car is four links in one act: somebody loses what they drive, the
// parts have a buyer, the buyer's trade picks up, and the forecourt eventually
// sells that person another one.

func stripper(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 2000, 60, 100
	w.Player.Location = "bar"
	var mark *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; w.WouldDrive(n) && !n.Dead && !IsOfficial(n.ID) {
			n.Location, n.Car = "bar", 1
			w.MeetPerson(n.ID)
			mark = w.NPC(n.ID)
			break
		}
	}
	if mark == nil {
		t.Fatal("nobody in this city drives")
	}
	return w, mark
}

func TestStrippingACarTakesItOffWhoeverOwnedIt(t *testing.T) {
	w, mark := stripper(t)
	cash, sore := w.Player.Cash, mark.Sore
	if reason := w.StripReadiness("bar"); reason != "" {
		t.Fatalf("a car standing outside cannot be touched: %q", reason)
	}
	if err := w.StripCar("bar"); err != nil {
		t.Fatal(err)
	}
	if mark.Car != 0 {
		t.Errorf("%s was stripped and still drives", mark.Name)
	}
	if w.Player.Cash <= cash {
		t.Errorf("the parts sold for %d", w.Player.Cash-cash)
	}
	if mark.Sore <= sore {
		t.Errorf("somebody lost their car and holds %d against you, up from %d", mark.Sore, sore)
	}
	if w.Player.Heat == 0 {
		t.Error("stripping a car in the street drew no attention at all")
	}
	t.Logf("the parts made $%d; they hold %d against you and attention is %d",
		w.Player.Cash-cash, mark.Sore, w.Player.Heat)
}

// Nothing to strip where nobody drives.
func TestThereIsNothingToStripWhereNobodyDrives(t *testing.T) {
	w, mark := stripper(t)
	mark.Car = 0
	if reason := w.StripReadiness("bar"); reason == "" {
		t.Error("a room with no car in it offered one to strip")
	}
	if err := w.StripCar("bar"); err == nil {
		t.Error("a car nobody owns was stripped anyway")
	}
}

// The buyer's half: parts go to a garage, and a garage with more of them coming
// through does better trade. This is the link the inbox asked for.
func TestGaragesDoBetterWhenThereAreMorePartsAbout(t *testing.T) {
	w, _ := stripper(t)
	before := map[string]int{}
	for _, l := range Locations {
		if l.Kind == "garage" {
			before[l.ID] = w.Custom(l.ID)
		}
	}
	if len(before) < 2 {
		t.Fatal("the city has fewer than two garages")
	}
	if err := w.StripCar("bar"); err != nil {
		t.Fatal(err)
	}
	lifted := 0
	for id, was := range before {
		if w.Custom(id) > was {
			lifted++
		}
	}
	if lifted == 0 {
		t.Error("a car was stripped and no garage in the city saw any more work")
	}
	t.Logf("%d of %d garages picked up trade", lifted, len(before))
}

// And the city says it. A number nobody sees is not a link.
func TestTheCitySaysACarWasStripped(t *testing.T) {
	w, mark := stripper(t)
	if err := w.StripCar("bar"); err != nil {
		t.Fatal(err)
	}
	said := false
	for _, text := range writings(w) {
		if contains(text, mark.Name) && (contains(text, "parts") || contains(text, "stripped")) {
			said = true
		}
	}
	if !said {
		t.Errorf("a car was stripped off %s and the city wrote nothing about it", mark.Name)
	}
}

// The wiring. A function that works and is never reached is the seam that has
// caught me more than any other, so the button and the command are checked as
// the interface actually uses them.
func TestTakingACarApartIsSomethingThePlayerCanActuallyDo(t *testing.T) {
	w, mark := stripper(t)
	var button *Action
	for i, a := range w.Actions("bar") {
		if a.ID == "strip" {
			button = &w.Actions("bar")[i]
		}
	}
	if button == nil {
		t.Fatal("a car standing outside is offered no button at all")
	}
	if button.Disabled {
		t.Fatalf("the button is refused: %s", button.Reason)
	}
	if button.Cost != 0 {
		t.Errorf("the button declares a cost of %d for something that pays out", button.Cost)
	}

	before := w.Player.Cash
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "strip", Target: "bar"})
	if err != nil {
		t.Fatalf("the command was refused: %v", err)
	}
	if got := next.NPC(mark.ID); got == nil || got.Car != 0 {
		t.Error("the command went through and the car is still theirs")
	}
	if next.Player.Cash <= before {
		t.Errorf("the command paid %d", next.Player.Cash-before)
	}
}

// And your own bench is worth having: parts are worth more when you hold a
// garage to take them to.
func TestPartsAreWorthMoreWithAGarageOfYourOwn(t *testing.T) {
	w, mark := stripper(t)
	loose := w.PartsWorth(mark)
	for _, l := range Locations {
		if l.Kind == "garage" {
			w.Properties[l.ID].Owner = "player:1"
			break
		}
	}
	if own := w.PartsWorth(mark); own <= loose {
		t.Errorf("parts are worth %d with a garage of your own and %d without", own, loose)
	}
}
