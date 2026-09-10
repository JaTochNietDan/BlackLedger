package core

import "testing"

// The people behind your counters see the street all day. Somebody standing
// across the road from your laundry for three afternoons is the sort of thing
// they would mention, and until this the only way to learn that a family had
// commissioned an attack on a business of yours was to go and investigate it.
//
// It also gives the trust of somebody you employ a job. A man who thinks well of
// you tells you what he saw. A man who does not keeps his head down.

func counterWorld(t *testing.T, trust int) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	prop := w.Properties["laundry"]
	if len(prop.Hands) == 0 {
		t.Fatal("nobody works at the laundry")
	}
	for _, who := range prop.Hands {
		w.NPC(who).Trust = trust
	}
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	w.Plots = append(w.Plots, Plot{ID: ID(), Kind: "sabotage", Life: w.Life,
		Due: w.Minute + 4320, Actor: rival, Target: "laundry", Strength: 35})
	return w, "laundry"
}

func TestSomebodyBehindYourCounterNoticesWhatIsComing(t *testing.T) {
	told, quiet := 0, 0
	for seed := uint32(1); seed <= 300; seed++ {
		w, _ := counterWorld(t, 80)
		w.WorldRNG = seed * 2654435761
		w.WordFromTheCounter()
		if w.Plots[len(w.Plots)-1].Known {
			told++
		}
		q, _ := counterWorld(t, 0)
		q.WorldRNG = seed * 2654435761
		q.WordFromTheCounter()
		if !q.Plots[len(q.Plots)-1].Known {
			quiet++
		}
	}
	t.Logf("of 300 days: people who think well of you mentioned it %d times, people who do not kept quiet %d times", told, quiet)
	if told == 0 {
		t.Fatal("nobody who works for you ever noticed a family casing the place they stand in all day")
	}
	if told == 300 {
		t.Fatal("every single day somebody spotted it, which is not noticing, it is being told")
	}
	if quiet <= 300-told {
		t.Fatalf("somebody who thinks nothing of you is as useful as somebody who does: %d against %d", 300-quiet, told)
	}
}

// And they only see their own street. A plot against a business across the city
// is not something the man at your laundry counter knows about.
func TestTheyOnlySeeTheirOwnStreet(t *testing.T) {
	w, _ := counterWorld(t, 100)
	w.Plots[len(w.Plots)-1].Target = "butcher"
	for i := 0; i < 60; i++ {
		w.WordFromTheCounter()
	}
	if w.Plots[len(w.Plots)-1].Known {
		t.Fatal("the laundry counter reported an attack being planned on a butcher's shop across town")
	}
}

// A room with nobody in it tells you nothing, which is what being short-handed
// costs beyond the takings.
func TestAnEmptyCounterSeesNothing(t *testing.T) {
	w, id := counterWorld(t, 100)
	prop := w.Properties[id]
	prop.Hands, prop.Staff = nil, 0
	for i := 0; i < 60; i++ {
		w.WordFromTheCounter()
	}
	if w.Plots[len(w.Plots)-1].Known {
		t.Fatal("a laundry with nobody working in it noticed somebody casing it")
	}
}

// Knowing has to be worth something or the word from the counter is flavour. A
// business that is expecting it takes less: the shutters come down, the stock
// goes out the back, and whoever comes finds a room that is ready for them.
func TestAWarnedBusinessTakesLess(t *testing.T) {
	hit := func(warned bool) int {
		w, id := counterWorld(t, 80)
		p := w.Plots[len(w.Plots)-1]
		p.Known = warned
		w.Properties[id].Condition = 100
		w.Player.Crew = nil
		w.ResolveSabotage(p)
		return 100 - w.Properties[id].Condition
	}
	cold, ready := hit(false), hit(true)
	t.Logf("a sabotage against a laundry: %d condition off it cold, %d off it when the counter had said something", cold, ready)
	if cold == 0 {
		t.Fatal("sabotage does nothing at all, so this proves nothing")
	}
	if ready >= cold {
		t.Fatalf("a business that saw them coming took the same beating: %d against %d", ready, cold)
	}
}

// And an empty counter cannot be ready for anything, however much warning there
// was.
func TestAnEmptyCounterCannotBeReady(t *testing.T) {
	w, id := counterWorld(t, 80)
	p := w.Plots[len(w.Plots)-1]
	p.Known = true
	w.Properties[id].Condition, w.Player.Crew = 100, nil
	w.Properties[id].Hands, w.Properties[id].Staff = nil, 0
	w.ResolveSabotage(p)
	empty := 100 - w.Properties[id].Condition

	w2, id2 := counterWorld(t, 80)
	p2 := w2.Plots[len(w2.Plots)-1]
	p2.Known = true
	w2.Properties[id2].Condition, w2.Player.Crew = 100, nil
	w2.ResolveSabotage(p2)
	staffed := 100 - w2.Properties[id2].Condition
	if empty <= staffed {
		t.Fatalf("a laundry with nobody in it defended itself as well as one with three: %d against %d", empty, staffed)
	}
}
