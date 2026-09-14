package core

import (
	"fmt"
	"strings"
	"testing"
)

func withCrew(t *testing.T, loyalty int) *World {
	t.Helper()
	w := New(59)
	w.MigrateLivingWorld()
	w.Player.Location, w.Player.Health, w.Player.Respect, w.Player.Cash = "club", 100, 40, 3000
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: loyalty}}
	return w
}

func TestYouCannotSendSomebodyYouDoNotHave(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

// Sending one of your own was strictly worse than going yourself and there was
// nothing to do about it: four attempts in sixty against eleven, and the only
// thing that moved the number was how they felt about you. Loyalty is earned
// slowly and cannot be bought at a counter, so anybody who preferred not to be
// shot at had no way to make the safe option any good.
func attempts(t *testing.T, weapon int, own bool) int {
	t.Helper()
	done := 0
	for seed := uint32(1); seed <= 200; seed++ {
		w := New(seed * 2654435761)
		w.Event, w.District = nil, 9
		w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 40, 20000
		hand := w.Holder("driver")
		if hand == nil {
			t.Fatal("this city has nobody who could be sent")
		}
		w.Player.Crew = append(w.Player.Crew, Crew{ID: hand.ID, Name: hand.Name, Loyalty: 65})
		var mark *NPC
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if !n.Dead && n.Faction != "" && n.Rank < RankLeader {
				mark = n
				break
			}
		}
		if mark == nil {
			t.Fatal("this city has nobody in a family to go after")
		}
		w.Player.Location = mark.Location
		h := w.OwnHands()
		if own {
			w.Player.Weapon = weapon
		} else {
			hand.Weapon = weapon
			sent, ok := w.CrewHands()
			if !ok {
				t.Fatal("nobody to send")
			}
			h = sent
		}
		if err := w.Strike(mark.ID, h); err != nil {
			t.Fatalf("going after somebody standing in front of you was refused: %v", err)
		}
		if mark.Dead {
			done++
		}
	}
	return done
}

func TestArmingTheOneYouSendIsWorthSomething(t *testing.T) {
	heavy(t)
	bare := attempts(t, 0, false)
	revolver := attempts(t, 1, false)
	shotgun := attempts(t, 2, false)
	thompson := attempts(t, 3, false)
	t.Logf("200 sent each: empty handed %d, a revolver %d, a shotgun %d, a Thompson %d",
		bare, revolver, shotgun, thompson)
	if revolver <= bare || shotgun <= revolver || thompson <= shotgun {
		t.Fatalf("putting something in their hand is worth nothing: %d, %d, %d, %d",
			bare, revolver, shotgun, thompson)
	}

	// And going yourself is still better, with the same gun, or there would be
	// no reason ever to take the risk.
	mine := attempts(t, 3, true)
	t.Logf("with the same Thompson: %d sent against %d gone yourself", thompson, mine)
	if mine <= thompson {
		t.Fatalf("going yourself buys nothing over sending an armed man: %d against %d", mine, thompson)
	}
	// Which is the trade: a man you armed does it nearly as well and you are
	// not the one who was there.
	if thompson*2 < mine {
		t.Fatalf("an armed man you sent is worth %d against %d, which is not a choice", thompson, mine)
	}
}

// The counter sells it to them, and only to your own.
func TestTheCounterWillArmYourOwnPeople(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 40, 20000
	hand := w.Holder("driver")
	if hand == nil {
		t.Fatal("nobody to sign on")
	}
	hand.Faction, hand.Location = w.PlayerOrganizationID(), "docks"
	w.Player.Location, w.Event = "docks", nil

	offered := map[string]*Action{}
	for i, a := range w.Actions("docks") {
		if strings.HasPrefix(a.ID, "give:") {
			offered[a.ID] = &w.Actions("docks")[i]
		}
	}
	top := Armaments("weapon")[len(Armaments("weapon"))-1]
	id := fmt.Sprintf("give:%s:%d", hand.ID, top.Tier)
	if _, ok := offered[id]; !ok {
		t.Fatalf("%s is standing on the dock and there is no way to put anything in their hand", hand.Name)
	}
	if err := w.BuyArmsFor(hand.ID, top.Tier); err != nil {
		t.Fatal(err)
	}
	if hand.Weapon != top.Tier {
		t.Fatalf("bought them %s and they are carrying tier %d", top.Label, hand.Weapon)
	}
	// Not for somebody who is not yours.
	stranger := ""
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Faction != w.PlayerOrganizationID() {
			stranger = n.ID
			break
		}
	}
	if err := w.BuyArmsFor(stranger, 1); err == nil {
		t.Fatal("bought a gun for somebody who does not answer to the player")
	}
	// And nothing worse than what they carry.
	if err := w.BuyArmsFor(hand.ID, 1); err == nil {
		t.Fatal("bought them a revolver while they are carrying a Thompson")
	}
}

func TestFailedDelegatedStrikeNamesTheTargetAndPlace(t *testing.T) {
	seen := 0
	for seed := uint32(1); seed <= 64; seed++ {
		w := withCrew(t, 80)
		w.RNG = seed * 2654435761
		mark := &NPC{ID: "target", Name: "Morgan Dale", Location: "bar"}
		w.NPCs = append(w.NPCs, *mark)
		hand, _ := w.CrewHands()
		w.itWentWrong(w.NPC(mark.ID), hand, "Saint Agnes", nil)
		for _, record := range w.History {
			if !strings.HasSuffix(record.Title, " took it instead of you") {
				continue
			}
			seen++
			if !strings.Contains(record.Text, "go after Morgan Dale at Saint Agnes") {
				t.Fatalf("lost the delegated task: %s", record.Text)
			}
		}
	}
	if seen == 0 {
		t.Fatal("did not exercise an injured returning crew member")
	}
}

func TestFailedStrikeDoesNotReportADeadCrewMemberEscaping(t *testing.T) {
	injured, dead := 0, 0
	for seed := uint32(1); seed <= 128; seed++ {
		w := withCrew(t, 80)
		w.RNG = seed * 2654435761
		w.NPCs = append(w.NPCs, NPC{ID: "strike-target", Name: "Morgan Dale", Location: "bar"})
		hand, _ := w.CrewHands()
		w.itWentWrong(w.NPC("strike-target"), hand, "Saint Agnes", nil)
		escaped := false
		for _, record := range w.History {
			if record.Title == "They got away with nothing" {
				escaped = true
			}
		}
		crew := w.NPC("leo")
		if crew != nil && crew.Dead {
			dead++
			if escaped {
				t.Fatalf("seed %d: dead crew member reported escaping", seed)
			}
		}
		if len(w.Player.Crew) > 0 {
			injured++
			if !escaped {
				t.Fatal("surviving crew lost its escape report")
			}
		}
	}
	if dead == 0 || injured == 0 {
		t.Fatalf("missing outcomes: dead=%d injured=%d", dead, injured)
	}
}
