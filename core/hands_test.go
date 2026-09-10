package core

import "testing"

// "We want to really have a lot of characters living in this city."
//
// A business's staff was a number. You hired a pair of hands, the wage bill went
// up, and nobody in Bellwether had a job: the person behind the counter of a
// place the player owned did not exist, could not be talked to, killed, robbed,
// poached or arrested, and the city's own people had no reason to be anywhere in
// the daytime except the one the routine invented for them.

func staffed(t *testing.T) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.Player.Location = "laundry"
	return w, "laundry"
}

func TestHiringTakesOnSomebodyWhoLivesHere(t *testing.T) {
	t.Parallel()
	w, id := staffed(t)
	// A business the player has just taken over already has people in it, and
	// they are people: the day settles that before anybody is hired.
	w.EmptyChairs()
	prop0 := w.Properties[id]
	if len(prop0.Hands) != prop0.Staff {
		t.Fatalf("a laundry with %d positions filled has %d names behind the counter", prop0.Staff, len(prop0.Hands))
	}
	if err := w.LayOff(id); err != nil {
		t.Fatal(err)
	}
	before := w.Properties[id].Staff
	if err := w.Hire(id); err != nil {
		t.Fatalf("hiring was refused: %v", err)
	}
	prop := w.Properties[id]
	if prop.Staff != before+1 {
		t.Fatalf("a position was filled and the count says %d of %d", prop.Staff, before)
	}
	if len(prop.Hands) != prop.Staff {
		t.Fatalf("%d positions are filled and %d of them have a name", prop.Staff, len(prop.Hands))
	}
	n := w.NPC(prop.Hands[len(prop.Hands)-1])
	if n == nil {
		t.Fatalf("the person hired is nobody in this city: %q", prop.Hands[0])
	}
	if n.Dead {
		t.Fatalf("%s was hired and is dead", n.Name)
	}
	if w.EmployerOf(n.ID) != id {
		t.Fatalf("%s works at %q", n.Name, w.EmployerOf(n.ID))
	}
	// And they are behind that counter in the daytime, which is what having a
	// job means to everything else in this city.
	if n.Post != id {
		t.Fatalf("%s works at the laundry and spends the day at %q", n.Name, n.Post)
	}
}

func TestNobodyHoldsTwoJobs(t *testing.T) {
	t.Parallel()
	w, id := staffed(t)
	own(w, "butcher")
	w.EmptyChairs()
	for i := 0; i < 3; i++ {
		_ = w.LayOff(id)
		_ = w.LayOff("butcher")
		_ = w.Hire(id)
		_ = w.Hire("butcher")
	}
	seen := map[string]string{}
	for _, place := range []string{id, "butcher"} {
		for _, who := range w.Properties[place].Hands {
			if was, twice := seen[who]; twice {
				t.Fatalf("%s works at %s and at %s", w.NPC(who).Name, was, place)
			}
			seen[who] = place
		}
	}
	if len(seen) == 0 {
		t.Fatal("nobody was hired anywhere")
	}
}

func TestLettingSomebodyGoLetsThatPersonGo(t *testing.T) {
	t.Parallel()
	w, id := staffed(t)
	w.EmptyChairs()
	who := w.Properties[id].Hands[len(w.Properties[id].Hands)-1]
	if err := w.LayOff(id); err != nil {
		t.Fatalf("letting somebody go was refused: %v", err)
	}
	if len(w.Properties[id].Hands) != w.Properties[id].Staff {
		t.Fatalf("%d positions and %d names", w.Properties[id].Staff, len(w.Properties[id].Hands))
	}
	if w.EmployerOf(who) != "" {
		t.Fatalf("%s was let go and still works at %s", w.NPC(who).Name, w.EmployerOf(who))
	}
}

// The point of the staff being people: a position held by somebody who is dead
// is not a position that is filled.
func TestAPositionHeldByADeadPersonIsAnEmptyPosition(t *testing.T) {
	t.Parallel()
	w, id := staffed(t)
	w.EmptyChairs()
	prop := w.Properties[id]
	filled := prop.Staff
	n := w.NPC(prop.Hands[0])
	n.Dead = true
	w.BusinessDay()
	if prop.Staff != filled-1 {
		t.Fatalf("%s is dead and the laundry still has %d of its positions filled", n.Name, prop.Staff)
	}
	for _, who := range prop.Hands {
		if who == n.ID {
			t.Fatalf("%s is dead and still on the books", n.Name)
		}
	}
}

// Having a job has to mean something to the rest of the city, or it is a list
// of names in a save file. The people the player employs are behind that
// counter in the daytime, where they can be found, talked to and taken.
func TestThePeopleYouEmployAreWhereTheyWorkInTheDaytime(t *testing.T) {
	t.Parallel()
	w, id := staffed(t)
	w.EmptyChairs()
	hands := append([]string{}, w.Properties[id].Hands...)
	if len(hands) == 0 {
		t.Fatal("nobody works at the laundry")
	}
	seen := map[string]bool{}
	for step := 0; step < 24*7; step++ {
		w.Event = nil
		w.Advance(60)
		w.Event = nil
		if Evening(w.Minute) {
			continue
		}
		for _, who := range hands {
			if n := w.NPC(who); n != nil && n.Location == id {
				seen[who] = true
			}
		}
	}
	if len(seen) < len(hands) {
		t.Fatalf("%d of the %d people on the laundry's books were never once behind its counter in a week of daytimes", len(hands)-len(seen), len(hands))
	}
}

// And what a player who employs people can ask about them.
func TestYouCanAskWhoWorksForYou(t *testing.T) {
	t.Parallel()
	w, id := staffed(t)
	own(w, "butcher")
	w.EmptyChairs()
	people := w.AtWork()
	if len(people) == 0 {
		t.Fatal("two businesses and nobody works for you")
	}
	places := map[string]bool{}
	for _, n := range people {
		places[w.EmployerOf(n.ID)] = true
		if n.Dead {
			t.Fatalf("%s is dead and on your books", n.Name)
		}
	}
	if !places[id] || !places["butcher"] {
		t.Fatalf("you hold two businesses and your people are at %v", places)
	}
}

// Now that the people behind a counter are people, they can be taken. Somebody
// good is somebody a rival is paying, and money is the whole of the argument.
func poachable(t *testing.T) (*World, *NPC, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	// A rival holds the butcher and has people in it.
	w.Properties["butcher"].Owner = w.Factions[0].ID
	if w.Properties["butcher"].Owner == w.PlayerOrganizationID() {
		w.Properties["butcher"].Owner = w.Factions[1].ID
	}
	w.EmptyChairs()
	if len(w.Properties["butcher"].Hands) == 0 {
		t.Fatal("a rival's butcher has nobody in it")
	}
	if err := w.LayOff("laundry"); err != nil {
		t.Fatal(err)
	}
	n := w.NPC(w.Properties["butcher"].Hands[0])
	w.Player.Location = n.Location
	return w, n, "laundry"
}

func TestYouCanTakeSomebodyOffARivalsCounter(t *testing.T) {
	t.Parallel()
	w, n, mine := poachable(t)
	was := w.EmployerOf(n.ID)
	cash := w.Player.Cash
	if reason := w.PoachReadiness(n.ID, mine); reason != "" {
		t.Fatalf("taking %s on was refused: %s", n.Name, reason)
	}
	if err := w.Poach(n.ID, mine); err != nil {
		t.Fatalf("taking %s on failed: %v", n.Name, err)
	}
	if w.EmployerOf(n.ID) != mine {
		t.Fatalf("%s was taken on and works at %q", n.Name, w.EmployerOf(n.ID))
	}
	if w.Player.Cash >= cash {
		t.Fatal("taking somebody off a rival's books cost nothing")
	}
	// The rival is short-handed for it.
	if prop := w.Properties[was]; prop.Staff != len(prop.Hands) {
		t.Fatalf("the butcher says %d positions filled and has %d people", prop.Staff, len(prop.Hands))
	}
	for _, who := range w.Properties[was].Hands {
		if who == n.ID {
			t.Fatalf("%s works for you and is still on the butcher's books", n.Name)
		}
	}
}

func TestARivalMindsHavingSomebodyTakenOffThem(t *testing.T) {
	t.Parallel()
	w, n, mine := poachable(t)
	rival := w.faction(w.Properties[w.EmployerOf(n.ID)].Owner)
	goodwill := rival.Goodwill
	if err := w.Poach(n.ID, mine); err != nil {
		t.Fatal(err)
	}
	if rival.Goodwill >= goodwill {
		t.Fatalf("%s lost a pair of hands to you and thinks exactly as much of you as before", rival.Name)
	}
}

func TestYouCannotTakeOnSomebodyWithNowhereToPutThem(t *testing.T) {
	t.Parallel()
	w, n, _ := poachable(t)
	// Fill the laundry back up, so there is no position going.
	for w.HireReadiness("laundry") == "" {
		if err := w.Hire("laundry"); err != nil {
			break
		}
	}
	if reason := w.PoachReadiness(n.ID, "laundry"); reason == "" {
		t.Fatal("somebody was taken on into a business with every position filled")
	}
}

// An action nobody can press is a feature in a file. This walks the same path
// the interface walks.
func TestTakingSomebodyOnIsOfferedWhereTheyAreStanding(t *testing.T) {
	t.Parallel()
	w, n, _ := poachable(t)
	var offered *Action
	for _, a := range w.Actions(w.Player.Location) {
		if a.ID == "poach:"+n.ID {
			offered = &a
		}
	}
	if offered == nil {
		t.Fatalf("%s works for a rival and is standing in front of you, and there is no way to offer them anything", n.Name)
	}
	if offered.Disabled {
		t.Fatalf("taking %s on is refused where it should work: %s", n.Name, offered.Reason)
	}
	if offered.Subject != n.ID {
		t.Fatalf("an offer aimed at %s is filed under %q", n.Name, offered.Subject)
	}
	if err := w.apply(Command{Kind: "poach:" + n.ID, Target: w.Player.Location, RequestID: "poachthemover1"}); err != nil {
		t.Fatalf("taking them on through the same path as everything else failed: %v", err)
	}
	if w.EmployerOf(n.ID) == "" || !w.Own(w.EmployerOf(n.ID)) {
		t.Fatalf("%s was taken on and works at %q", n.Name, w.EmployerOf(n.ID))
	}
	// And is no longer offered, because they already work for you.
	for _, a := range w.Actions(w.Player.Location) {
		if a.ID == "poach:"+n.ID && !a.Disabled {
			t.Fatalf("%s already works for you and is still being offered a place", n.Name)
		}
	}
}
