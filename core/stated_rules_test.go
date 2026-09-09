package core

import (
	"strings"
	"testing"
)

// The guide prints six rules as prose. They are the game explaining itself to
// the player, and every one is a claim about the code. The page beside them
// already turned out to be wrong about itself twice, so these are worth asking
// rather than assuming.

// "The clock only moves when you commit to something. Reading, inspecting and
// choosing cost nothing."
//
// Read as a property: an action that declares no minutes must not advance the
// clock, and one that declares minutes must advance it by what it said.
func TestAnActionCostsTheTimeItSaysAndNoOther(t *testing.T) {
	offered, _ := sweepCities(t)
	if len(offered) < 50 {
		t.Fatalf("the sweep only produced %d kinds of action", len(offered))
	}
	checked, free := 0, 0
	for _, seed := range []uint32{31, 88} {
		w := New(seed)
		w.District = 2
		w.Player.Cash, w.Player.Respect, w.Player.Health = 20000, 200, 100
		for _, id := range []string{"laundry", "garage", "casino"} {
			w.Properties[id].Owner = "player:1"
		}
		w.Incorporate()
		w.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
		for i := range Locations {
			id := Locations[i].ID
			w.Player.Location = id
			for _, a := range w.Actions(id) {
				if a.Disabled || a.Minutes != 0 {
					continue
				}
				free++
				before := w.Minute
				next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: a.ID, Target: id})
				if err != nil {
					continue
				}
				checked++
				spent := next.Minute - before
				// Work whose effect runs the clock itself declares that time as
				// Away, because Minutes is what the command layer spends and
				// spending it twice would send the player away for a fortnight.
				// Either way the button has to say how long it takes.
				if a.Away > 0 {
					if spent < a.Away {
						t.Errorf("%s at %s says %d minutes away and moved the clock %d",
							a.ID, id, a.Away, spent)
					}
					continue
				}
				if spent != 0 {
					t.Errorf("%s at %s declares no time at all and moved the clock %d minutes",
						a.ID, id, spent)
				}
				w = next
				w.Player.Cash, w.Player.Health = 20000, 100
			}
		}
	}
	if checked == 0 {
		t.Fatalf("no free action was ever taken (%d offered), so this proved nothing", free)
	}
	t.Logf("took %d actions that declare no time", checked)
}

// "Attention is public and so are the thresholds. Past N the police come to the
// door; past M they take the premises." The numbers in that sentence have to be
// the numbers the game uses.
func TestTheThresholdsInTheRulesAreTheOnesTheGameUses(t *testing.T) {
	rules := GuideRules()
	stated := ""
	for _, r := range rules {
		if strings.Contains(r, "Attention is public") {
			stated = r
		}
	}
	if stated == "" {
		t.Fatal("the rule about attention is no longer printed")
	}
	if !strings.Contains(stated, itoa(RaidThreshold)) || !strings.Contains(stated, itoa(ForfeitThreshold)) {
		t.Fatalf("the rule does not quote the constants: %q", stated)
	}
	if RaidThreshold >= ForfeitThreshold {
		t.Fatalf("the rule says the police come first and take later, and the numbers say %d then %d",
			RaidThreshold, ForfeitThreshold)
	}
}

// "People are not furniture ... anybody out on the street is not at the address
// they left and cannot be dealt with until they arrive — including your own,
// when you send them somewhere."
func TestNobodyOnTheStreetCanBeDealtWith(t *testing.T) {
	w := New(29)
	w.District = 2
	w.Player.Cash, w.Player.Respect = 8000, 120
	w.Player.Location = "bar"
	// Enough standing that lending is offered at all, and somebody worth
	// lending to.
	w.Player.Respect = LendStanding + 40
	// Somebody the game already offers work about, rather than a person I
	// appended: the fixer and the driver are refused earlier and correctly for
	// reasons that have nothing to do with the street, and lending is only
	// offered to the first few people in a room, so an extra name at the end of
	// the list never gets an offer to test.
	walking := (*NPC)(nil)
	for _, a := range w.Actions("bar") {
		if strings.HasPrefix(a.ID, "lend:") && !a.Disabled {
			walking = w.NPC(a.Subject)
			break
		}
	}
	if walking == nil {
		t.Fatal("nobody in the bar can be lent to, so this proved nothing")
	}
	// Put them on the street between two addresses.
	walking.Heading = "market"
	walking.Arrives = w.Minute + 40
	if !w.Travelling(walking) {
		t.Fatal("they are not out walking")
	}
	if reason := w.OutOfReach(walking.ID); reason == "" {
		t.Fatal("somebody out on the street can be dealt with")
	}
	// The rule is kept in two places, and this asserts the whole of it rather
	// than the half I first assumed. Somebody out walking is not in the room,
	// so most work about them is never offered at all; anything that still is
	// must be refused, and refused because they are on the street.
	for i := range Locations {
		id := Locations[i].ID
		w.Player.Location = id
		for _, a := range w.Actions(id) {
			if a.Subject != walking.ID || a.Disabled {
				continue
			}
			t.Errorf("%q at %s is live for somebody out on the street", a.Label, id)
		}
	}
	// And when they arrive, the same work becomes possible again — otherwise
	// this would pass for a person nobody can ever deal with.
	walking.Heading, walking.Arrives = "", 0
	w.Player.Location = "bar"
	live := 0
	for _, a := range w.Actions("bar") {
		if a.Subject == walking.ID && !a.Disabled {
			live++
		}
	}
	if live == 0 {
		t.Fatal("nothing is possible with them even standing still, so this proved nothing")
	}
	t.Logf("%d actions become live once they stop walking", live)
}
