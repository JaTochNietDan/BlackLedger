package core

import "testing"

func doorman(t *testing.T) (*World, *NPC) {
	t.Helper()
	w, member := testator(t)
	return w, member
}

// postAndArrive puts somebody on a door and waits for them to walk there.
// Posting stopped being instant when the city's people started having to walk
// places: the door is worth nothing until somebody is standing in it, which is
// the point of these tests rather than a detail of them.
func postAndArrive(t *testing.T, w *World, id string) {
	t.Helper()
	if err := w.Post(id); err != nil {
		t.Fatal(err)
	}
	if n := w.NPC(w.Properties[id].Posted); n != nil && w.Travelling(n) {
		w.Minute = n.Arrives
		w.Arrivals()
	}
}

func TestNobodyStandsOnADoorThatIsNotYours(t *testing.T) {
	w, _ := doorman(t)
	for _, id := range []string{"club", "market", "room"} {
		if w.PostReadiness(id) == "" {
			t.Fatalf("somebody was put on the door at %s", id)
		}
	}
	if w.PostReadiness("laundry") != "" {
		t.Fatal("nobody could be put on your own laundry:", w.PostReadiness("laundry"))
	}
}

func TestAManCanOnlyBeInOnePlace(t *testing.T) {
	w, member := doorman(t)
	postAndArrive(t, w, "laundry")
	if posted := w.PostedAt("laundry"); posted == nil || posted.ID != member.ID {
		t.Fatal("nobody was on the door")
	}
	if member.Location != "laundry" {
		t.Fatalf("he is standing at %s", member.Location)
	}
	if len(w.Unposted()) != 0 {
		t.Fatal("he was available for somewhere else as well")
	}
	if w.PostReadiness("garage") == "" {
		t.Fatal("the same man was put on two doors")
	}
	if w.PostReadiness("laundry") == "" {
		t.Fatal("a second man was put on the same door")
	}
	// Taking him off frees him.
	if err := w.Unpost("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.PostedAt("laundry") != nil || len(w.Unposted()) != 1 {
		t.Fatal("he was still on it")
	}
	if err := w.Unpost("laundry"); err == nil {
		t.Fatal("took nobody off a door twice")
	}
}

func TestADeadManComesOffTheDoorByHimself(t *testing.T) {
	w, member := doorman(t)
	w.Post("laundry")
	member.Dead = true
	if w.PostedAt("laundry") != nil {
		t.Fatal("a dead man was still on the door")
	}
	if w.Properties["laundry"].Posted != "" {
		t.Fatal("the door still names him")
	}
	// And somebody who stops answering to you comes off it too.
	live, other := doorman(t)
	live.Post("laundry")
	other.Faction = ""
	if live.PostedAt("laundry") != nil {
		t.Fatal("somebody who left was still standing there")
	}
}

func TestAManOnTheDoorIsWorthSomethingAgainstARaid(t *testing.T) {
	w, member := doorman(t)
	if w.PostingDefenceAt("laundry") != 0 {
		t.Fatal("an empty door was worth something")
	}
	postAndArrive(t, w, "laundry")
	worth := w.PostingDefenceAt("laundry")
	if worth <= PostingDefence {
		t.Fatalf("a man on the door was worth %d", worth)
	}
	member.Trust = 100
	if w.PostingDefenceAt("laundry") <= worth {
		t.Fatal("somebody who means it was worth no more than somebody who does not")
	}

	// And it shows up where it matters: the same raid, with and without.
	const runs = 400
	held := func(post bool) int {
		kept := 0
		for seed := uint32(1); seed <= runs; seed++ {
			probe, _ := doorman(t)
			probe.WorldRNG = seed * 2654435761
			if post {
				postAndArrive(t, probe, "laundry")
			}
			attacker := &probe.Factions[0]
			attacker.Power = 80
			condition := probe.Properties["laundry"].Condition
			probe.contestAt(attacker, probe.PlayerOrganization(), "laundry")
			if probe.Properties["laundry"].Condition == condition {
				kept++
			}
		}
		return kept
	}
	bare, manned := held(false), held(true)
	if manned <= bare {
		t.Fatalf("a manned door turned away %d raids in %d against %d for an empty one", manned, runs, bare)
	}
	t.Logf("%d raids on the same laundry: an empty door turns away %d, a manned one %d", runs, bare, manned)
}

func TestAManOnTheDoorTurnsAwayTheStreet(t *testing.T) {
	const runs = 400
	robbed := func(post bool) int {
		hit := 0
		for seed := uint32(1); seed <= runs; seed++ {
			w, _ := doorman(t)
			w.WorldRNG = seed * 2654435761
			w.Properties["laundry"].Income = 30
			if post {
				postAndArrive(t, w, "laundry")
			}
			cash := w.Player.Cash
			var thief *NPC
			for _, n := range w.Civilians() {
				if !IsOfficial(n.ID) && !w.isRoleHolder(n) {
					thief = n
					break
				}
			}
			thief.Skill = 90
			for i := 0; i < 6 && w.Player.Cash == cash; i++ {
				w.takeFromSomebody(thief)
			}
			if w.Player.Cash < cash {
				hit++
			}
		}
		return hit
	}
	bare, manned := robbed(false), robbed(true)
	if manned >= bare {
		t.Fatalf("a manned laundry was robbed %d times in %d against %d for an empty one", manned, runs, bare)
	}
	t.Logf("%d attempts on the same laundry: an empty one loses its takings %d times, a manned one %d", runs, bare, manned)
}

func TestTheManOnTheDoorIsTheOneWhoPaysForIt(t *testing.T) {
	// A raid that gets through reaches him before it reaches anybody else. Over
	// enough raids on a place that falls, he is the one it costs.
	const runs = 300
	lost := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, member := doorman(t)
		w.WorldRNG = seed * 2654435761
		postAndArrive(t, w, "laundry")
		attacker := &w.Factions[0]
		attacker.Power = 100
		w.contestAt(attacker, w.PlayerOrganization(), "laundry")
		if member.Dead {
			if w.Properties["laundry"].Posted != "" {
				t.Fatal("a dead man is still named on the door")
			}
			lost++
		}
	}
	if lost == 0 {
		t.Fatal("300 raids on a manned door and nobody standing in it was ever hurt")
	}
	if lost > runs/2 {
		t.Fatalf("a man on the door died in %d of %d raids, which is not a job anybody takes", lost, runs)
	}
	t.Logf("%d raids on a manned laundry cost the man on the door his life %d times", runs, lost)
}
