package core

import "testing"

// "Why can't I attempt to take people out? ... either send a family member
// after someone to kill them or to attempt to kill them myself, where doing it
// myself comes with much greater risk of my own injury or death based on my
// skills and equipment."
//
// Only the paid contract existed. You could put a price on a name and wait, and
// that was the whole of violence against a person in this city.

func striker(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := New(37)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 3000, 60, 100
	w.Player.Location = "bar"
	var mark *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Faction == w.PlayerOrganizationID() {
			continue
		}
		n.Location = "bar"
		w.MeetPerson(n.ID)
		mark = w.NPC(n.ID)
		break
	}
	if mark == nil {
		t.Fatal("nobody to go after")
	}
	return w, mark
}

func TestYouCanGoAfterSomebodyYourself(t *testing.T) {
	w, mark := striker(t)
	if reason := w.StrikeReadiness(mark.ID); reason != "" {
		t.Fatalf("standing in a room with them: %q", reason)
	}
	dead, hurt, survived, killedYou := 0, 0, 0, 0
	for seed := uint32(1); seed <= 60; seed++ {
		run, target := striker(t)
		run.RNG, run.WorldRNG = seed*2654435761, seed*2654435761
		run.Player.Weapon, run.Player.Armour = 1, 1
		health := run.Player.Health
		if err := run.Strike(target.ID, run.OwnHands()); err != nil {
			t.Fatal(err)
		}
		if n := run.NPC(target.ID); n == nil || n.Dead {
			dead++
		} else {
			survived++
			if n.Sore == 0 {
				t.Error("somebody survived being shot at and holds nothing against you")
			}
		}
		if run.Player.Health < health || !run.Player.Alive {
			hurt++
		}
		if !run.Player.Alive {
			killedYou++
		}
	}
	t.Logf("sixty attempts, armed and wearing something: %d finished, %d survived, and you were hurt or killed in %d",
		dead, survived, hurt)
	if dead == 0 {
		t.Error("sixty attempts and nobody ever died")
	}
	if survived == 0 {
		t.Error("going after somebody always works, which is not a risk")
	}
	if hurt == 0 {
		t.Error("going yourself never cost you anything")
	}
	// "much greater risk of my own injury or death" — the death half of that is
	// the whole difference between going and sending, so it is asked for on its
	// own rather than folded into "you were hurt".
	if killedYou == 0 {
		t.Error("sixty attempts in person and it was never the end of you")
	}
	t.Logf("and it killed you %d times", killedYou)
}

// Being armed and wearing something has to matter, or "based on my skills and
// equipment" is a sentence about nothing.
func TestGoingInArmedIsBetterThanGoingInEmptyHanded(t *testing.T) {
	armed, bare := 0, 0
	for seed := uint32(1); seed <= 120; seed++ {
		run, target := striker(t)
		run.RNG, run.WorldRNG = seed*40503, seed*40503
		run.Player.Weapon, run.Player.Armour = 2, 2
		_ = run.Strike(target.ID, run.OwnHands())
		if n := run.NPC(target.ID); n == nil || n.Dead {
			armed++
		}
		plain, mark := striker(t)
		plain.RNG, plain.WorldRNG = seed*40503, seed*40503
		_ = plain.Strike(mark.ID, plain.OwnHands())
		if n := plain.NPC(mark.ID); n == nil || n.Dead {
			bare++
		}
	}
	t.Logf("armed and armoured finished %d of 120; empty-handed finished %d", armed, bare)
	if armed <= bare {
		t.Errorf("a gun made no difference: %d against %d", armed, bare)
	}
}

// "If you send someone of your own then there's a chance they are killed or
// captured and then they could be interrogated and give you up."
func TestSendingYourOwnRisksThemRatherThanYou(t *testing.T) {
	w, mark := striker(t)
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
	if n := w.NPC("leo"); n != nil {
		n.Location = w.Player.Location
	}
	if reason := w.SendReadiness(mark.ID); reason != "" {
		t.Fatalf("with a loyal crew in the room: %q", reason)
	}
	dead, lost, gaveYouUp, hurtYou := 0, 0, 0, 0
	for seed := uint32(1); seed <= 60; seed++ {
		run, target := striker(t)
		run.RNG, run.WorldRNG = seed*2246822519, seed*2246822519
		run.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
		if n := run.NPC("leo"); n != nil {
			n.Location = run.Player.Location
		}
		hand, ok := run.CrewHands()
		if !ok {
			t.Fatal("nobody to send")
		}
		heat, health := run.Player.Heat, run.Player.Health
		if err := run.Strike(target.ID, hand); err != nil {
			t.Fatal(err)
		}
		if n := run.NPC(target.ID); n == nil || n.Dead {
			dead++
		}
		if n := run.NPC("leo"); n == nil || n.Dead || n.Held > run.Minute {
			lost++
		}
		if run.Player.Heat > heat+10 {
			gaveYouUp++
		}
		if run.Player.Health < health {
			hurtYou++
		}
	}
	t.Logf("sixty sent: %d finished, your own was killed or taken in %d, it came back to you in %d",
		dead, lost, gaveYouUp)
	if dead == 0 {
		t.Error("sixty attempts by somebody else and nobody ever died")
	}
	if lost == 0 {
		t.Error("sending somebody never costs them anything")
	}
	if gaveYouUp == 0 {
		t.Error("your own person is taken and it never comes back to you")
	}
	if hurtYou > 0 {
		t.Errorf("you were hurt %d times by a job you sent somebody else on", hurtYou)
	}
}

// The rules that stop it being a way to tidy up your own house.
func TestThereAreNamesYouCannotGoAfterThisWay(t *testing.T) {
	w, mark := striker(t)
	if reason := w.StrikeReadiness("nobody-at-all"); reason == "" {
		t.Error("you can go after somebody who does not exist")
	}
	mark.Location = "docks"
	if reason := w.StrikeReadiness(mark.ID); reason == "" {
		t.Error("you can go after somebody on the other side of the city")
	}
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
	if n := w.NPC("leo"); n != nil {
		n.Location = w.Player.Location
		if reason := w.StrikeReadiness("leo"); reason == "" {
			t.Error("your own crew is a target for you")
		}
	}
}

// The wiring: both buttons are offered where somebody is standing, and pressing
// them plays out there and then.
func TestBothWaysOfGoingAfterSomebodyAreOfferedInTheRoom(t *testing.T) {
	w, mark := striker(t)
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
	if n := w.NPC("leo"); n != nil {
		n.Location = w.Player.Location
	}
	found := map[string]bool{}
	for _, a := range w.Actions(w.Player.Location) {
		if a.ID == "strike:"+mark.ID || a.ID == "send:"+mark.ID {
			found[a.ID] = true
			if a.Subject != mark.ID {
				t.Errorf("%s is not filed under %s", a.ID, mark.Name)
			}
			if a.Disabled {
				t.Errorf("%s is refused: %s", a.ID, a.Reason)
			}
		}
	}
	if !found["strike:"+mark.ID] || !found["send:"+mark.ID] {
		t.Fatalf("the room offers %v", found)
	}
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "strike:" + mark.ID})
	if err != nil {
		t.Fatal(err)
	}
	if next.Minute != w.Minute+StrikeMinutes {
		t.Errorf("an attempt took %d minutes", next.Minute-w.Minute)
	}
	if next.LastResult == nil || len(next.LastResult.Records) == 0 {
		t.Error("an attempt was made and the city has no account of it")
	}
}
