package core

import "testing"

// walksOut runs the clock on to the moment a person actually leaves. Deciding
// to go and going are no longer the same minute: a shift change is spread over
// the hours a shift change takes, so a test that wants somebody on the street
// has to wait for them to pick up their coat.
func walksOut(w *World, n *NPC) {
	if n.Sets > w.Minute {
		w.Minute = n.Sets
	}
	w.Arrivals()
}

// The city's people stood at fixed addresses for their whole lives. Somebody
// who ran a business was permanently inside it; somebody who hated a man across
// town never went to find him. The only thing that ever changed an address was
// a takeover, and it changed instantly — a man was in one building and then, in
// the same minute, in another.
//
// Nobody in this city walks anywhere, which is the largest thing still missing
// from a city that is supposed to be lived in.

func TestSomebodyWhoIsNotWhereTheyWorkGoesToWork(t *testing.T) {
	w := New(4)
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	w.SetOut()
	if n.Heading != "laundry" {
		t.Fatalf("the woman who runs the laundry is at the bar and heading for %q", n.Heading)
	}
	if n.Arrives <= w.Minute {
		t.Fatalf("she arrives at minute %d and it is minute %d", n.Arrives, w.Minute)
	}
	walksOut(w, n)
	if !w.Travelling(n) {
		t.Fatal("she has set off and the city does not think she is going anywhere")
	}
	// While she is walking she is in neither building. A room shows who is in
	// it, and she is not in it.
	for _, p := range w.PeopleHere("bar") {
		if p.ID == n.ID {
			t.Fatal("she is on the street and still standing in the bar")
		}
	}
	for _, p := range w.PeopleHere("laundry") {
		if p.ID == n.ID {
			t.Fatal("she has arrived somewhere she has not reached yet")
		}
	}
}

func TestAJourneyTakesTheTimeAJourneyTakes(t *testing.T) {
	w := New(4)
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	w.SetOut()
	walksOut(w, n)
	want := TravelMinutes("bar", "laundry")
	if n.Arrives-w.Minute != want {
		t.Fatalf("the walk from the bar to the laundry takes %d minutes and she has allowed %d", want, n.Arrives-w.Minute)
	}
	// Half way there she is still on the street.
	w.Minute += want / 2
	w.Arrivals()
	if !w.Travelling(n) || n.Location != "bar" {
		t.Fatalf("half way through the walk she is at %q, travelling %v", n.Location, w.Travelling(n))
	}
	w.Minute += want
	w.Arrivals()
	if w.Travelling(n) || n.Location != "laundry" {
		t.Fatalf("after the walk she is at %q, travelling %v", n.Location, w.Travelling(n))
	}
	if len(w.PeopleHere("laundry")) == 0 {
		t.Fatal("she arrived and the laundry is empty")
	}
}

func TestNobodyWalksOutOfACell(t *testing.T) {
	w := New(4)
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	n.Held = w.Minute + 2880
	w.SetOut()
	if n.Heading != "" {
		t.Fatalf("a woman in a cell set off for %q", n.Heading)
	}
	dead := w.NPC("leo")
	dead.Role, dead.Location, dead.Dead = "Runs Bluebird Laundry", "bar", true
	w.SetOut()
	if dead.Heading != "" {
		t.Fatal("the dead have somewhere to be")
	}
}

// What the city view needs: who is out there, where they came from, where they
// are going and why, in the words somebody watching the street would use.
func TestTheStreetSaysWhoIsOnItAndWhy(t *testing.T) {
	w := New(4)
	n := w.NPC("mara")
	n.Role, n.Location = "Runs Bluebird Laundry", "bar"
	w.SetOut()
	walksOut(w, n)
	street := w.OnTheStreet()
	if len(street) == 0 {
		t.Fatal("somebody set off and the street is empty")
	}
	var her Journeying
	for _, j := range street {
		if j.ID == n.ID {
			her = j
		}
	}
	if her.Name != n.Name || her.From != "Saint Agnes" || her.To != "Bluebird Laundry" {
		t.Fatalf("the street says %q is going from %q to %q", her.Name, her.From, her.To)
	}
	if her.Because == "" || her.Minutes <= 0 {
		t.Fatalf("she is going somewhere for %q, %d minutes out", her.Because, her.Minutes)
	}
	if her.FromID != "bar" || her.ToID != "laundry" {
		t.Fatal("the street cannot say which addresses the journey is between")
	}
}

// A grudge that has never made anybody do anything is a number in a save file.
func TestSomebodyCarryingSomethingGoesToFindTheManTheyBlame(t *testing.T) {
	w := New(4)
	holder, target := w.NPC("mara"), w.NPC("leo")
	holder.Location, target.Location = "bar", "club"
	// Somebody standing still. You go to where a man is; if he is himself out
	// walking there is nothing to go to, which is a rule of its own below.
	target.Role = ""
	w.Resent(holder.ID, target.ID, GrudgeActs+10, "what happened at the docks")
	w.SetOut()
	if holder.Heading != "club" {
		t.Fatalf("she is carrying enough to act on and is heading for %q", holder.Heading)
	}
	walksOut(w, holder)
	street := w.OnTheStreet()
	if len(street) == 0 || street[0].Because == "" {
		t.Fatal("the city cannot say why she is walking across town")
	}
}

// You cannot go and find a man who is himself out on the street. Waiting until
// he is somewhere is what stops two people walking past each other forever.
func TestYouDoNotSetOutAfterSomebodyWhoIsAlreadyWalking(t *testing.T) {
	w := New(4)
	holder, target := w.NPC("mara"), w.NPC("leo")
	holder.Location, target.Location = "bar", "club"
	target.Heading, target.Arrives, target.Errand = "docks", w.Minute+40, "somewhere else"
	w.Resent(holder.ID, target.ID, GrudgeActs+10, "what happened at the docks")
	w.SetOut()
	if holder.Heading == "club" {
		t.Fatal("she set off for a building the man had already left")
	}
}

// At scale: a city where everybody is always walking is as wrong as one where
// nobody ever does. Twenty seeds, three weeks each, counting how often the
// street has somebody on it and how many journeys finish.
func TestTheCityWalksWithoutBeingAlwaysInMotion(t *testing.T) {
	journeys, sampled, occupied, mostAtOnce, crowd := 0, 0, 0, 0, 0
	for seed := 0; seed < 20; seed++ {
		w := New(uint32(seed) * 2654435761)
		crowd = max(crowd, len(w.People()))
		heading := map[string]string{}
		for day := 0; day < 21; day++ {
			for step := 0; step < 8; step++ {
				// Step through the half-day an hour at a time. People no
				// longer all leave on the boundary — a shift change is spread
				// over the hours a shift change takes — so sampling only at
				// the boundary would find the street empty and prove nothing.
				boundary := (w.Minute/720 + 1) * 720
				for at := w.Minute + 60; at <= boundary; at += 60 {
					w.Minute = at
					w.Arrivals()
					if w.Minute%720 == 0 {
						w.SetOut()
					}
					street := w.OnTheStreet()
					sampled++
					if len(street) > 0 {
						occupied++
					}
					if len(street) > mostAtOnce {
						mostAtOnce = len(street)
					}
					for i := range w.NPCs {
						n := &w.NPCs[i]
						if was := heading[n.ID]; was != "" && n.Heading == "" && n.Location == was {
							journeys++
						}
						heading[n.ID] = n.Heading
					}
				}
			}
			w.PeopleDay()
			w.FillRoles()
		}
	}
	share := float64(occupied) / float64(sampled)
	t.Logf("journeys completed %d · street occupied on %.0f%% of %d samples · most at once %d",
		journeys, share*100, sampled, mostAtOnce)
	if journeys == 0 {
		t.Fatal("three weeks in twenty cities and nobody walked anywhere")
	}
	// Staggering departures once made this test sample only moments when
	// nobody had set off yet, so it saw an empty street and passed on nothing.
	// A city with journeys in it has somebody walking some of the time.
	if share < .05 {
		t.Fatalf("the street was occupied on %.0f%% of samples: this test is not looking at anything", share*100)
	}
	if share > .9 {
		t.Fatalf("somebody is on the street %.0f%% of the time: the city is permanently in motion", share*100)
	}
	// A share of the city, not a count. This was "more than twelve", written
	// when the city held about fifty people, and it meant "a quarter of them".
	// The city holds nearly ninety now and twelve stopped meaning that: the
	// figure was a proportion all along, and the message beside it said so.
	if crowd == 0 {
		t.Fatal("nobody lives here")
	}
	if share := float64(mostAtOnce) / float64(crowd); share > .25 {
		t.Fatalf("%d of %d people were out on the street at once, which is %.0f%% of the city",
			mostAtOnce, crowd, share*100)
	}
}

// The point of minding ground: when a place changes hands, the city's people
// move because of it. This is the war becoming visible on the street rather
// than only in the ledger.
func TestWhenGroundChangesHandsSomebodyWalks(t *testing.T) {
	w := New(4)
	w.SetOut() // let the city settle into who minds what
	for i := 0; i < 6; i++ {
		w.Minute += 720
		w.Arrivals()
		w.SetOut()
	}
	if len(w.OnTheStreet()) != 0 {
		t.Fatal("the city never settles: somebody is always walking")
	}
	// A holding of the Bellandi family passes to the Russo Outfit.
	taken := ""
	for _, id := range w.FamilyHoldings("bellandi") {
		if w.Properties[id] != nil {
			taken = id
			break
		}
	}
	if taken == "" {
		t.Fatal("the family holds nothing to lose")
	}
	w.Properties[taken].Owner = "russo"
	w.Minute += 720
	w.SetOut()
	// Everybody who decided to go has a departure time of their own; run on to
	// the last of them so the whole shift change is on the street at once.
	latest := w.Minute
	for i := range w.NPCs {
		if s := w.NPCs[i].Sets; s > latest {
			latest = s
		}
	}
	w.Minute = latest
	w.Arrivals()
	street := w.OnTheStreet()
	if len(street) == 0 {
		t.Fatalf("%s changed hands and nobody in the city moved", taken)
	}
	going := false
	for _, j := range street {
		if j.ToID == taken {
			going = true
		}
	}
	if !going {
		t.Fatalf("ground changed hands and nobody went to it: %+v", street)
	}
}
