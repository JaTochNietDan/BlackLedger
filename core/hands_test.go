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
