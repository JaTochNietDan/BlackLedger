package core

import "testing"

// "Leo heads out" was the log line, and Leo did not head anywhere. He stood
// exactly where he had been standing, the money appeared two hours later, and
// the city said he had returned from somewhere he never went. Everybody else
// in this city has to walk; the man the player pays was the last one who did
// not have to.

func crewman(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := New(4)
	w.Player.Location = "bar"
	w.Player.Respect = 20
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "recruit", Target: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	w = next
	if len(w.Player.Crew) == 0 {
		t.Fatal("nobody was recruited")
	}
	n := w.NPC(w.Player.Crew[0].ID)
	if n == nil {
		t.Fatal("the crew member is not somebody in the city")
	}
	return w, n
}

func TestSendingSomebodyOnCollectionsSendsThemSomewhere(t *testing.T) {
	w, leo := crewman(t)
	w.Properties["laundry"].Owner = "player:1"
	w.Player.Cash = 2000
	leo.Location = "bar"
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "delegate", Target: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	him := next.NPC(leo.ID)
	if !next.Travelling(him) {
		t.Fatalf("he was sent on collections and is standing at %q", him.Location)
	}
	if him.Heading != "laundry" {
		t.Fatalf("collections are from the laundry and he is heading for %q", him.Heading)
	}
	if him.Errand == "" {
		t.Fatal("the city cannot say why he is walking")
	}
	// And he is on the street like anybody else.
	found := false
	for _, j := range next.OnTheStreet() {
		if j.ID == him.ID {
			found = true
		}
	}
	if !found {
		t.Fatal("the man the player pays is crossing the city invisibly")
	}
}

// With nothing of the player's to collect from, the work still has an address.
func TestCollectionsHaveAnAddressEvenWithNoPremises(t *testing.T) {
	w, leo := crewman(t)
	leo.Location = "club"
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "delegate", Target: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	him := next.NPC(leo.ID)
	if him.Heading == "" {
		t.Fatal("he was sent on collections with nowhere to go")
	}
	if _, ok := PlaceByID(him.Heading); !ok {
		t.Fatalf("he was sent to %q, which is not an address", him.Heading)
	}
}

// The round still pays what it paid, and he comes back from where he went.
func TestHeComesBackFromWhereHeWent(t *testing.T) {
	w, leo := crewman(t)
	w.Properties["laundry"].Owner = "player:1"
	from := "bar"
	leo.Location = from
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "delegate", Target: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	w = next
	cash := w.Player.Cash
	// Run the clock past the round.
	for i := 0; i < 12 && len(w.Tasks) > 0; i++ {
		w.Minute += 60
		w.Arrivals()
		w.settleTasks()
	}
	if len(w.Tasks) != 0 {
		t.Fatal("the round never finished")
	}
	if w.Player.Cash <= cash {
		t.Fatalf("collections paid %d", w.Player.Cash-cash)
	}
	him := w.NPC(leo.ID)
	if !w.Travelling(him) && him.Location != from {
		t.Fatalf("he finished collections at %q and is not coming back to %q", him.Location, from)
	}
	if w.Travelling(him) && him.Heading != from {
		t.Fatalf("he is walking to %q rather than back to %q", him.Heading, from)
	}
}

// While he is doing the round, the city says so rather than reciting his job.
func TestAManOnARoundIsDoingTheRoundNotHisJob(t *testing.T) {
	w, leo := crewman(t)
	w.Properties["laundry"].Owner = "player:1"
	leo.Location = "bar"
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "delegate", Target: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	him := next.NPC(leo.ID)
	next.Minute = him.Arrives
	next.Arrivals()
	doing := next.doingNow(him)
	if !contains(doing, "Collecting") {
		t.Fatalf("he is standing in the laundry collecting and the city says %q", doing)
	}
	// And when the round is over he goes back to being what he is.
	next.Minute += CollectionMinutes
	next.settleTasks()
	next.Arrivals()
	if contains(next.doingNow(him), "Collecting") {
		t.Fatalf("the round is over and he is still %q", next.doingNow(him))
	}
}
