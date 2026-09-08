package core

import "testing"

func lieutenant(t *testing.T) (*World, *Faction, *NPC) {
	t.Helper()
	w, f := recruit(t)
	if err := w.Serve(f.ID); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < PromotionWork*2; i++ {
		w.ServeWork(f.ID)
	}
	if w.ServiceRank() != RankLieutenant {
		t.Fatal("the fixture did not come up")
	}
	leader := w.Leader(f.ID)
	if leader == nil {
		t.Fatal("nobody was in the chair")
	}
	w.Player.Location, w.Player.Respect, w.Player.Health = leader.Location, TakeoverStanding, 100
	f.Peak, f.Power = 100, 100*TakeoverWeak/5
	return w, f, leader
}

func TestNobodyBelowALieutenantGetsCloseEnough(t *testing.T) {
	w, f := recruit(t)
	w.Serve(f.ID)
	leader := w.Leader(f.ID)
	w.Player.Location, w.Player.Respect, w.Player.Health = leader.Location, TakeoverStanding, 100
	f.Peak, f.Power = 100, 40
	if w.TakeoverReadiness() == "" {
		t.Fatal("an associate moved on the man in the chair")
	}
	for i := 0; i < PromotionWork*2; i++ {
		w.ServeWork(f.ID)
	}
	if w.TakeoverReadiness() != "" {
		t.Fatal("a lieutenant could not:", w.TakeoverReadiness())
	}
}

func TestNobodyMovesOnAManWhoseOrganizationIsWinning(t *testing.T) {
	w, f, _ := lieutenant(t)
	f.Power = 100
	if w.TakeoverReadiness() == "" {
		t.Fatal("somebody moved on the head of an organization at full strength")
	}
}

func TestYouHaveToBeInTheRoomAndWorthFollowing(t *testing.T) {
	w, _, leader := lieutenant(t)
	w.Player.Location = "room"
	if w.TakeoverReadiness() == "" {
		t.Fatal("moved on a man across the city")
	}
	w.Player.Location = leader.Location
	w.Player.Respect = 0
	if w.TakeoverReadiness() == "" {
		t.Fatal("nobody would have followed them and they went anyway")
	}
	w.Player.Respect = TakeoverStanding
	w.Player.Health = 30
	if w.TakeoverReadiness() == "" {
		t.Fatal("moved on him while barely standing")
	}
}

func TestTakingItMeansTakingAllOfIt(t *testing.T) {
	found := false
	for seed := uint32(1); seed <= 400 && !found; seed++ {
		w, f, leader := lieutenant(t)
		w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
		name, id := f.Name, f.ID
		held := len(w.FamilyHoldings(id))
		people := len(w.Members(id))
		rival := w.Factions[1].ID
		c := w.Conflict(id, rival)
		c.Hostility, c.State = 70, "war"
		if err := w.TakeOver(); err != nil {
			t.Fatal(err)
		}
		if !leader.Dead {
			continue
		}
		found = true
		me := w.PlayerOrganizationID()
		if len(w.FamilyHoldings(me)) != held {
			t.Fatalf("took %d of the %d premises", len(w.FamilyHoldings(me)), held)
		}
		if !w.Incorporated() {
			t.Fatal("taking an organization did not make them one")
		}
		if w.Serving() != "" {
			t.Fatal("they still answer to somebody")
		}
		if w.faction(id) != nil {
			t.Fatalf("%s still exists", name)
		}
		kept := len(w.Members(me))
		if kept == 0 || kept == people {
			t.Fatalf("%d of %d people stayed", kept, people)
		}
		after := w.Conflict(rival, me)
		if after == nil || after.State != "war" {
			t.Fatal("the quarrels did not come with it")
		}
		if !w.hasRecord("It is yours") || !w.hasNewsKind("politics") {
			t.Fatal("nobody heard about it")
		}
	}
	if !found {
		t.Fatal("no attempt in 400 ever succeeded")
	}
}

func TestBeingExpectedIsTheEndOfIt(t *testing.T) {
	found := false
	for seed := uint32(1); seed <= 400 && !found; seed++ {
		w, f, leader := lieutenant(t)
		w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
		goodwill := f.Goodwill
		if err := w.TakeOver(); err != nil {
			t.Fatal(err)
		}
		if leader.Dead {
			continue
		}
		found = true
		if w.Serving() != "" {
			t.Fatal("they still answer to the man they moved on")
		}
		if w.Player.Health >= 100 {
			t.Fatal("being expected cost nothing")
		}
		if f.Goodwill >= goodwill {
			t.Fatal("the organization did not mind")
		}
		if !w.hasRecord("He was expecting it") {
			t.Fatal("nothing was recorded")
		}
	}
	if !found {
		t.Fatal("no attempt in 400 ever failed")
	}
}

func TestTheOddsAreWhatYouBringAgainstWhatHeHas(t *testing.T) {
	w, f, leader := lieutenant(t)
	base := w.takeoverOdds(leader, f)
	w.Player.Respect = 100
	if w.takeoverOdds(leader, f) <= base {
		t.Fatal("a bigger name was worth nothing in the room")
	}
	w.Player.Respect = TakeoverStanding
	w.Player.Weapon = 3
	if w.takeoverOdds(leader, f) <= base {
		t.Fatal("what you are carrying was worth nothing")
	}
	w.Player.Weapon = 0
	f.Power = 100
	if w.takeoverOdds(leader, f) >= base {
		t.Fatal("what the organization can put behind him was worth nothing")
	}
}
