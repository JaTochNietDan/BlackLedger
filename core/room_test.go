package core

import (
	"strings"
	"testing"
)

// Half of what a player does anywhere is done to somebody standing there. These
// guard the two halves of making that legible: the core says who each action is
// about, and it says who is in the room to be acted on.

func TestWorkAimedAtSomebodyKnowsWhoItIsAimedAt(t *testing.T) {
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Contacts = 90000, 200, 4
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 70}}
	w.ensureOfficials()
	w.OrganizationDay()
	for _, n := range w.Civilians() {
		if IsOfficial(n.ID) || w.isRoleHolder(n) || len(w.OwnPeople()) >= 2 {
			continue
		}
		w.Player.Location = n.Location
		w.SignOn(n.ID)
	}

	// Anything whose id carries somebody's name, and the handful that put the
	// name in the label instead, must say who it is about.
	personal := map[string]bool{"mug": true, "recruit": true, "contact": true,
		"delegate": true, "crew_bonus": true}
	checked, missing := 0, []string{}
	for _, l := range Locations {
		w.Player.Location = l.ID
		here := map[string]bool{}
		for _, person := range w.PeopleHere(l.ID) {
			here[person.ID] = true
		}
		for _, a := range w.Actions(l.ID) {
			aimed := personal[a.ID]
			if at := indexByte(a.ID, ':'); at >= 0 && w.NPC(a.ID[at+1:]) != nil {
				aimed = true
			}
			if !aimed {
				if a.Subject != "" && !here[a.Subject] {
					missing = append(missing, a.ID+" is about "+a.Subject+", who is not at "+l.ID)
				}
				continue
			}
			checked++
			if a.Subject == "" {
				missing = append(missing, a.ID+" at "+l.ID+" names nobody")
				continue
			}
			if w.NPC(a.Subject) == nil {
				missing = append(missing, a.ID+" is about "+a.Subject+", who does not exist")
			}
		}
	}
	if checked < 8 {
		t.Fatalf("only %d actions aimed at a person were offered anywhere; this is not measuring what it claims to", checked)
	}
	if len(missing) > 0 {
		t.Fatalf("%v", missing)
	}
	t.Logf("%d actions aimed at a person, every one of them naming who", checked)
}

func TestTheRoomKnowsWhoIsStandingInIt(t *testing.T) {
	w := proprietor(t)
	w.Player.Cash, w.Player.Respect = 40000, 200
	w.Player.Location = "bar"
	people := w.PeopleHere("bar")
	if len(people) == 0 {
		t.Fatal("nobody is at Saint Agnes, which is where this game begins")
	}
	for _, p := range people {
		if p.Name == "" || p.Standing == "" {
			t.Fatalf("%s is described as %q", p.ID, p.Standing)
		}
		if n := w.NPC(p.ID); n == nil || n.Location != "bar" {
			t.Fatalf("%s is listed at the bar and is not there", p.ID)
		}
	}
	// A stranger's character is not on display; somebody you have dealt with is.
	var stranger *Presence
	for i := range people {
		if !people[i].Known {
			stranger = &people[i]
		}
	}
	if stranger != nil && stranger.Temperament != "" {
		t.Fatalf("%s is a stranger and the player can read their character", stranger.Name)
	}

	// Somebody who owes money says so, and says whether it is late.
	own(w, "laundry")
	w.District = 2
	for _, n := range w.Civilians() {
		w.Player.Location = n.Location
		if w.LendReadiness(n.ID) != "" {
			continue
		}
		if err := w.Lend(n.ID); err != nil {
			t.Fatal(err)
		}
		for _, p := range w.PeopleHere(n.Location) {
			if p.ID != n.ID {
				continue
			}
			if p.Owes == 0 {
				t.Fatal("a man who owes money is shown owing nothing")
			}
			if p.Overdue {
				t.Fatal("a debt made this morning is already late")
			}
			return
		}
		t.Fatal("the man who took the money is not in the room he took it in")
	}
	t.Skip("nobody in this city would borrow")
}

func TestYourOwnPeopleAreListedFirst(t *testing.T) {
	w, member := testator(t)
	// Put everybody in one room so the ordering is the only thing being tested.
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead {
			n.Location = "laundry"
		}
	}
	w.Player.Location = "laundry"
	people := w.PeopleHere("laundry")
	if len(people) < 3 {
		t.Fatalf("only %d people in the room", len(people))
	}
	if people[0].ID != member.ID || !people[0].Yours {
		t.Fatalf("the first person listed is %s (%s), not the man who works for you", people[0].Name, people[0].Standing)
	}
	// And a stranger is never listed above somebody the player knows.
	seenUnknown := false
	for _, p := range people {
		if !p.Known {
			seenUnknown = true
		} else if seenUnknown {
			t.Fatalf("%s, who you know, is listed below a stranger", p.Name)
		}
	}
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// The ultimate shape of this game is a city the player can watch, and a city
// you can watch is one where every person on the screen is visibly occupied
// with something true.

func TestEverybodyInThisCityIsDoingSomething(t *testing.T) {
	w, member := testator(t)
	w.District = 2
	w.Player.Cash = 40000
	w.ensureOfficials()
	blank, checked := []string{}, 0
	for _, l := range Locations {
		for _, p := range w.PeopleHere(l.ID) {
			checked++
			if p.Doing == "" {
				blank = append(blank, p.Name+" at "+l.ID)
			}
		}
	}
	if checked < 20 {
		t.Fatalf("only %d people were placed anywhere; this is not measuring what it claims to", checked)
	}
	if len(blank) > 0 {
		t.Fatalf("people standing about doing nothing the city can name: %v", blank)
	}

	// And it is not the same sentence for everybody, or it says nothing.
	said := map[string]bool{}
	for _, l := range Locations {
		for _, p := range w.PeopleHere(l.ID) {
			said[p.Doing] = true
		}
	}
	if len(said) < 4 {
		t.Fatalf("%d people are doing %d distinguishable things", checked, len(said))
	}

	// What somebody is doing follows what is true of them, not their name.
	w.Player.Location = "laundry"
	if err := w.Post("laundry"); err != nil {
		t.Fatal(err)
	}
	for _, p := range w.PeopleHere("laundry") {
		if p.ID == member.ID && !strings.Contains(p.Doing, "door") {
			t.Fatalf("a man put on the door is %q", p.Doing)
		}
	}
	member.Held = w.Minute + 3*1440
	for _, l := range Locations {
		for _, p := range w.PeopleHere(l.ID) {
			if p.ID == member.ID && !strings.Contains(p.Doing, "Ward Street") {
				t.Fatalf("a man the police are holding is %q", p.Doing)
			}
		}
	}
	t.Logf("%d people placed across the city, %d distinguishable things being done", checked, len(said))
}
