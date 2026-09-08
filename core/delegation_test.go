package core

import "testing"

func withCrew(t *testing.T, loyalty int) *World {
	t.Helper()
	w := New(59)
	w.MigrateLivingWorld()
	w.Player.Location, w.Player.Health, w.Player.Respect, w.Player.Cash = "club", 100, 40, 3000
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: loyalty}}
	return w
}

func TestYouCannotSendSomebodyYouDoNotHave(t *testing.T) {
	w := withCrew(t, 80)
	w.Player.Crew = nil
	if _, ok := w.CrewHands(); ok {
		t.Fatal("there was somebody to send")
	}
	if w.DelegateReadiness() == "" {
		t.Fatal("sent nobody")
	}
	if err := w.RobBy("club", Hand{Crew: true, Name: "Nobody"}); err == nil {
		t.Fatal("robbed a casino with nobody")
	}
}

func TestNobodyGoesOutOnAJobForSomebodyTheyDoNotTrust(t *testing.T) {
	w := withCrew(t, HandLoyalty-1)
	if w.DelegateReadiness() == "" {
		t.Fatal("a man on 39 loyalty went out on an errand like this")
	}
	w.Player.Crew[0].Loyalty = HandLoyalty
	if w.DelegateReadiness() != "" {
		t.Fatal("refused at exactly the stated loyalty:", w.DelegateReadiness())
	}
	// Nor while he is already out.
	w.Tasks = append(w.Tasks, Task{ID(), "Leo · collections", w.Minute + 120})
	if w.DelegateReadiness() == "" {
		t.Fatal("sent a man who was already on an assignment")
	}
}

func TestYourOwnHandsAreBetterAtItThanAnybodyElses(t *testing.T) {
	w := withCrew(t, 100)
	w.Player.Respect, w.Player.Weapon = 80, 2
	hand, _ := w.CrewHands()
	if w.HandEdge(w.OwnHands()) <= w.HandEdge(hand) {
		t.Fatalf("your own hands were worth %.3f against %.3f for sending somebody", w.HandEdge(w.OwnHands()), w.HandEdge(hand))
	}
	if w.robberyOdds("club", w.OwnHands()) <= w.robberyOdds("club", hand) {
		t.Fatal("sending somebody was better odds than going")
	}
}

func TestSendingSomebodyBuysLessStandingAndLessAttention(t *testing.T) {
	w := withCrew(t, 100)
	hand, _ := w.CrewHands()
	if w.HandRespectFor(w.OwnHands(), 10) <= w.HandRespectFor(hand, 10) {
		t.Fatal("a man who sends people was respected as much as one who goes")
	}
	if w.HandHeat(w.OwnHands(), 15) <= w.HandHeat(hand, 15) {
		t.Fatal("sending somebody drew the same attention as going")
	}
	if w.HandRespectFor(hand, 2) < 1 {
		t.Fatal("sending somebody was worth no standing at all")
	}
}

func TestTheRiskLandsOnWhoeverWasStandingThere(t *testing.T) {
	// The player takes it in health.
	mine := withCrew(t, 100)
	mine.HandHurt(mine.OwnHands(), 20, "something")
	if mine.Player.Health >= 100 {
		t.Fatal("the player walked away from their own beating")
	}
	if mine.Player.Crew[0].Loyalty != 100 {
		t.Fatal("somebody else paid for it")
	}

	// Somebody sent takes it in loyalty, and the player takes nothing.
	theirs := withCrew(t, 100)
	hand, _ := theirs.CrewHands()
	theirs.HandHurt(hand, 20, "something")
	if theirs.Player.Health != 100 {
		t.Fatal("the player was hurt by something they were not at")
	}
	if len(theirs.Player.Crew) == 0 || theirs.Player.Crew[0].Loyalty != 100-HandLoyaltyCost {
		t.Fatalf("loyalty went to %v", theirs.Player.Crew)
	}
}

func TestSometimesTheManYouSentDoesNotComeBack(t *testing.T) {
	const runs = 400
	lost := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := withCrew(t, 100)
		w.RNG = seed * 2654435761
		hand, _ := w.CrewHands()
		w.HandHurt(hand, 30, "take a till")
		if len(w.Player.Crew) == 0 {
			lost++
		}
	}
	if lost == 0 {
		t.Fatal("nobody sent on a job that went badly was ever killed")
	}
	if lost == runs {
		t.Fatal("everybody sent on a job that went badly was killed")
	}
	t.Logf("across %d jobs that went wrong badly enough, the man who was sent did not come back %d times", runs, lost)
}

func TestASmallSetbackNeverKillsTheManYouSent(t *testing.T) {
	for seed := uint32(1); seed <= 400; seed++ {
		w := withCrew(t, 100)
		w.RNG = seed * 2654435761
		hand, _ := w.CrewHands()
		w.HandHurt(hand, 24, "carry a message")
		if len(w.Player.Crew) == 0 {
			t.Fatal("a scuffle killed somebody")
		}
	}
}

func TestBothHandsAreOfferedAndBothAreRefusedForTheSameReasons(t *testing.T) {
	w := withCrew(t, 100)
	actions := map[string]Action{}
	for _, a := range w.Actions("club") {
		actions[a.ID] = a
	}
	if _, ok := actions["rob"]; !ok {
		t.Fatal("the player could not rob a casino themselves")
	}
	if _, ok := actions["rob:crew"]; !ok {
		t.Fatal("the player was not offered somebody to send")
	}
	// Both refuse for the reason that applies to both.
	w.Player.Health = 10
	for _, a := range w.Actions("club") {
		if (a.ID == "rob" || a.ID == "rob:crew") && !a.Disabled {
			t.Fatalf("%s was offered to somebody who could barely stand", a.ID)
		}
	}
	// And only the delegated one refuses for the reason that applies to it.
	w.Player.Health = 100
	w.Player.Crew[0].Loyalty = 10
	for _, a := range w.Actions("club") {
		if a.ID == "rob" && a.Disabled {
			t.Fatal("a disloyal crewman stopped the player doing it themselves")
		}
		if a.ID == "rob:crew" && !a.Disabled {
			t.Fatal("a disloyal crewman went anyway")
		}
	}
}
