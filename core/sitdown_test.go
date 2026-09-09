package core

import "testing"

// mediator builds a city with a quarrel bad enough to need somebody, and a
// player the two sides would come for.
func mediator(t *testing.T, hostility int, hotLeader bool) *World {
	t.Helper()
	w := New(97)
	w.MigrateLivingWorld()
	w.Player.Location = SitdownGround
	w.Player.Cash, w.Player.Respect, w.Player.Health = 5000, 60, 100
	c := &w.Conflicts[0]
	c.Hostility = hostility
	c.State = classify(c)
	// Make the leader of one side the kind of person who does or does not wait.
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Name != w.faction(c.A).Leader {
			continue
		}
		if hotLeader {
			n.Name = nameWith(t, "hot")
		} else {
			n.Name = nameWith(t, "careful")
		}
		w.faction(c.A).Leader = n.Name
	}
	// And the other side too, so the fixture is unambiguous.
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Name != w.faction(c.B).Leader {
			continue
		}
		n.Name = nameWith(t, "careful")
		w.faction(c.B).Leader = n.Name
	}
	return w
}

func TestNobodyComesToARoomForANobody(t *testing.T) {
	w := mediator(t, 70, false)
	w.Player.Respect = 0
	if w.SitdownReadiness() == "" {
		t.Fatal("two families crossed the city for somebody nobody had heard of")
	}
	w.Player.Respect = 60
	if w.SitdownReadiness() != "" {
		t.Fatal("they would not come for somebody with standing:", w.SitdownReadiness())
	}
	// Nor into a room arranged by somebody one of them wants dead.
	w.Factions[0].Goodwill = -60
	if w.SitdownReadiness() == "" {
		t.Fatal("a family that hates you sat in a room you arranged")
	}
	w.Factions[0].Goodwill = 0
	// Nor anywhere but neutral ground.
	w.Player.Location = "club"
	if w.SitdownReadiness() == "" {
		t.Fatal("a meeting was held on one side's own floor")
	}
}

func TestThereIsNothingToMediateInAQuietCity(t *testing.T) {
	w := mediator(t, 10, false)
	for i := range w.Conflicts {
		w.Conflicts[i].Hostility, w.Conflicts[i].State = 5, "cold"
	}
	if _, ok := w.OpenQuarrel(); ok {
		t.Fatal("a city at peace produced a quarrel to settle")
	}
	if w.SitdownReadiness() == "" {
		t.Fatal("a meeting was called about nothing")
	}
}

func TestPressingASettleableQuarrelEndsIt(t *testing.T) {
	w := mediator(t, 55, false)
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	if w.Event == nil || w.Event.Kind != "sitdown" || len(w.Event.Choices) != 4 {
		t.Fatalf("the room was %v", w.Event)
	}
	c := w.Conflict(w.Event.Actor, w.Event.Target)
	before, respect := c.Hostility, w.Player.Respect
	if err := w.ResolveSitdown(w.Event, "press"); err != nil {
		t.Fatal(err)
	}
	if c.Hostility >= before {
		t.Fatalf("hostility went %d to %d", before, c.Hostility)
	}
	if w.Player.Respect <= respect || w.faction(w.Event.Actor).Goodwill <= 0 {
		t.Fatal("settling a war bought the player nothing")
	}
	if !w.hasRecord("It is settled") {
		t.Fatal("nothing was recorded")
	}
	if w.Player.Health != 100 {
		t.Fatal("a settled meeting hurt somebody")
	}
}

func TestPressingAQuarrelSomebodyCameToFinishIsTheWorstRoomInTheCity(t *testing.T) {
	w := mediator(t, 80, true)
	q, ok := w.OpenQuarrel()
	if !ok || !q.Trap {
		t.Fatalf("a hostility-80 quarrel led by a hot-headed man was not a trap: %+v", q)
	}
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	before := 0
	for i := range w.NPCs {
		if w.NPCs[i].Dead {
			before++
		}
	}
	if err := w.ResolveSitdown(w.Event, "press"); err != nil {
		t.Fatal(err)
	}
	after := 0
	for i := range w.NPCs {
		if w.NPCs[i].Dead {
			after++
		}
	}
	if after <= before {
		t.Fatal("a bloodbath killed nobody")
	}
	if w.Player.Health >= 100 {
		t.Fatalf("the player walked out of it untouched, at %d health", w.Player.Health)
	}
	if !w.hasRecord("It was never a meeting") || !w.hasNewsKind("killing") {
		t.Fatal("the city never heard about it")
	}
}

func TestAWarmQuarrelIsNotATrapHoweverBadTheLeader(t *testing.T) {
	w := mediator(t, TrapHostility-5, true)
	q, ok := w.OpenQuarrel()
	if !ok || q.Trap {
		t.Fatal("a hot-headed man in a cooling quarrel was treated as an ambush")
	}
}

func TestYouOnlyKnowItIsATrapIfSomebodyToldYou(t *testing.T) {
	w := mediator(t, 80, true)
	w.Player.Contacts = 0
	q, _ := w.OpenQuarrel()
	if !q.Trap || q.Suspected {
		t.Fatal("a stranger with no contacts walked in knowing")
	}
	w.Player.Contacts = 3
	q, _ = w.OpenQuarrel()
	if !q.Suspected {
		t.Fatal("a network of contacts warned nobody")
	}
	// And the warning reaches the room itself.
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	if !containsName(w.Event.Body, "more people than the room needs") {
		t.Fatalf("the scene said: %s", w.Event.Body)
	}
}

func TestListeningIsSafeAndSettlesLittle(t *testing.T) {
	w := mediator(t, 80, true) // even in the worst room
	w.CallSitdown()
	c := w.Conflict(w.Event.Actor, w.Event.Target)
	before := c.Hostility
	if err := w.ResolveSitdown(w.Event, "listen"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Health != 100 {
		t.Fatal("saying nothing got the player shot")
	}
	if c.Hostility >= before || before-c.Hostility > 10 {
		t.Fatalf("hostility went %d to %d, which is either nothing or too much for saying nothing", before, c.Hostility)
	}
}

func TestTakingASideMakesOneEnemyAndOneFriend(t *testing.T) {
	w := mediator(t, 55, false)
	w.CallSitdown()
	a, b := w.faction(w.Event.Actor), w.faction(w.Event.Target)
	beforeA, beforeB := a.Goodwill, b.Goodwill
	c := w.Conflict(a.ID, b.ID)
	hostility := c.Hostility
	if err := w.ResolveSitdown(w.Event, "side"); err != nil {
		t.Fatal(err)
	}
	if a.Goodwill <= beforeA || b.Goodwill >= beforeB {
		t.Fatalf("standing went %d→%d and %d→%d", beforeA, a.Goodwill, beforeB, b.Goodwill)
	}
	if c.Hostility <= hostility {
		t.Fatal("taking a side calmed the quarrel down")
	}
	if len(w.Plots) == 0 {
		t.Fatal("the side you crossed did nothing about it")
	}
}

func TestLeavingCostsStandingWithBoth(t *testing.T) {
	w := mediator(t, 55, false)
	w.CallSitdown()
	a, b := w.faction(w.Event.Actor), w.faction(w.Event.Target)
	beforeA, beforeB := a.Goodwill, b.Goodwill
	if err := w.ResolveSitdown(w.Event, "leave"); err != nil {
		t.Fatal(err)
	}
	if a.Goodwill >= beforeA || b.Goodwill >= beforeB {
		t.Fatal("walking out cost nothing")
	}
	if w.Player.Health != 100 {
		t.Fatal("walking out got the player hurt")
	}
}

func TestTheFeeIsPaidWhateverHappens(t *testing.T) {
	w := mediator(t, 55, false)
	cash := w.Player.Cash
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-SitdownFee {
		t.Fatalf("the room cost $%d", cash-w.Player.Cash)
	}
	w.ResolveSitdown(w.Event, "leave")
	if w.Player.Cash != cash-SitdownFee {
		t.Fatal("leaving refunded the room")
	}
}
